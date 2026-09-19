package pagecapture

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/plexusone/w3pilot"
)

// Bounds on captured text to keep evidence bundles a sane size.
const (
	MaxHTMLBytes = 200_000
	MaxTextBytes = 40_000
)

// SettleMode selects how Capture waits for a page to be ready.
type SettleMode int

const (
	// SettleDOMStable (default) polls the DOM until it stops changing and has
	// real content. This is the robust choice for client-rendered SPAs, whose
	// "load" event fires before content exists and whose network may never idle.
	SettleDOMStable SettleMode = iota
	// SettleLoad waits only for the "load" event. Fine for static pages.
	SettleLoad
	// SettleNone captures immediately after navigation.
	SettleNone
)

// Options configure a capture.
type Options struct {
	// LoadWait bounds how long to wait for the page to settle. The wait is
	// best-effort: capture proceeds against the rendered DOM regardless.
	LoadWait time.Duration
	// MinSettle is the minimum time to wait before accepting a stable DOM.
	// SPAs render in phases (shell, then async data) with quiet gaps between;
	// this floor prevents settling in an early gap before data arrives.
	// Defaults to 2.5s when zero.
	MinSettle time.Duration
	// Settle selects the readiness strategy (default SettleDOMStable).
	Settle SettleMode
	// Screenshot toggles screenshot capture (the most expensive part).
	Screenshot bool
}

// DefaultOptions returns sensible defaults: DOM-stability settle, screenshot on.
func DefaultOptions() Options {
	return Options{LoadWait: 12 * time.Second, MinSettle: 2500 * time.Millisecond, Settle: SettleDOMStable, Screenshot: true}
}

// settlePollInterval is how often DOM stability is sampled.
const settlePollInterval = 300 * time.Millisecond

// domStableChecks is the number of consecutive unchanged samples (after the
// MinSettle floor) that count as "settled".
const domStableChecks = 3

// settleSignalScript returns a compact fingerprint of the DOM's size and text,
// so equality across samples means the page has stopped changing.
const settleSignalScript = `(() => {
  const els = document.querySelectorAll('*').length;
  const len = document.body ? document.body.innerText.length : 0;
  return els + ':' + len;
})()`

// settle waits for the page to be ready per opts.Settle. It is best-effort and
// returns a non-empty warning if readiness could not be confirmed in time.
func settle(ctx context.Context, pilot *w3pilot.Pilot, opts Options) string {
	if opts.LoadWait <= 0 {
		opts.LoadWait = 8 * time.Second
	}
	switch opts.Settle {
	case SettleNone:
		return ""
	case SettleLoad:
		if err := pilot.WaitForLoad(ctx, "load", opts.LoadWait); err != nil {
			return "page-load wait did not settle; captured rendered DOM anyway"
		}
		return ""
	default: // SettleDOMStable
		minSettle := opts.MinSettle
		if minSettle <= 0 {
			minSettle = 2500 * time.Millisecond
		}
		start := time.Now()
		deadline := start.Add(opts.LoadWait)
		last := ""
		stable := 0
		for time.Now().Before(deadline) {
			raw, err := pilot.Evaluate(ctx, settleSignalScript)
			if err == nil {
				if sig, ok := raw.(string); ok {
					// Require the DOM to hold real content (not just <html><body>).
					hasContent := sig != "" && !strings.HasSuffix(sig, ":0")
					if sig == last && hasContent {
						stable++
						// Only accept stability after the MinSettle floor, so an
						// early render lull before async data doesn't fool us.
						if stable >= domStableChecks && time.Since(start) >= minSettle {
							return ""
						}
					} else {
						stable = 0
						last = sig
					}
				}
			}
			select {
			case <-ctx.Done():
				return "settle cancelled"
			case <-time.After(settlePollInterval):
			}
		}
		return "DOM did not stabilize within the settle window; captured current DOM"
	}
}

// domData is the structural data extracted from the live DOM in one JS round-trip.
type domData struct {
	Lang     string    `json:"lang"`
	Text     string    `json:"text"`
	Links    []Link    `json:"links"`
	Headings []Heading `json:"headings"`
}

const extractScript = `(() => {
  const links = [...document.querySelectorAll('a[href]')].slice(0, 300).map(a => ({
    text: (a.innerText || a.getAttribute('aria-label') || '').trim(),
    href: a.href,
  }));
  const headings = [...document.querySelectorAll('h1,h2,h3,h4,h5,h6')].slice(0, 150).map(h => ({
    level: parseInt(h.tagName.substring(1), 10),
    text: (h.innerText || '').trim(),
  }));
  return JSON.stringify({
    lang: document.documentElement.getAttribute('lang') || '',
    text: document.body ? document.body.innerText.slice(0, 40000) : '',
    links: links,
    headings: headings,
  });
})()`

// Capture navigates to pageURL, waits for it to settle, and captures evidence.
func Capture(ctx context.Context, pilot *w3pilot.Pilot, pageURL string, opts Options) (*PageEvidence, error) {
	if err := pilot.Go(ctx, pageURL); err != nil {
		return nil, fmt.Errorf("navigate %s: %w", pageURL, err)
	}
	warning := settle(ctx, pilot, opts)
	ev := CaptureCurrent(ctx, pilot, pageURL, opts.Screenshot)
	ev.LoadWarning = warning
	return ev, nil
}

// CaptureCurrent captures evidence from the page the Pilot is *already* on,
// without navigating or settling. Use it when the caller has already loaded and
// settled the page (e.g. an audit engine that just ran against it), to avoid a
// redundant navigation. pageURL labels the evidence.
func CaptureCurrent(ctx context.Context, pilot *w3pilot.Pilot, pageURL string, screenshot bool) *PageEvidence {
	ev := &PageEvidence{URL: pageURL, CapturedAt: time.Now().UTC()}

	if title, err := pilot.Title(ctx); err == nil {
		ev.Title = title
	}

	html, err := pilot.Content(ctx)
	if err == nil {
		if len(html) > MaxHTMLBytes {
			html = html[:MaxHTMLBytes]
			ev.Truncated = true
		}
		ev.HTML = html
	}

	if raw, err := pilot.Evaluate(ctx, extractScript); err == nil {
		if s, ok := raw.(string); ok {
			var d domData
			if json.Unmarshal([]byte(s), &d) == nil {
				ev.Language = d.Lang
				ev.Text = d.Text
				ev.Links = d.Links
				ev.Headings = d.Headings
			}
		}
	}

	if screenshot {
		if shot, err := pilot.Screenshot(ctx); err == nil {
			ev.ScreenshotPNG = shot
		}
	}

	return ev
}

// CaptureURL launches a headless browser, captures one page, and closes it.
// Convenience for callers that don't already hold a Pilot.
func CaptureURL(ctx context.Context, pageURL string, opts Options) (*PageEvidence, error) {
	pilot, err := w3pilot.LaunchHeadless(ctx)
	if err != nil {
		return nil, fmt.Errorf("launch browser: %w", err)
	}
	defer func() { _ = pilot.Close(ctx) }()
	return Capture(ctx, pilot, pageURL, opts)
}

// CrawlOptions configure a site capture.
type CrawlOptions struct {
	Options
	// MaxPages caps how many pages to capture (including the start page).
	MaxPages int
	// SameHostOnly restricts the crawl to the start URL's host.
	SameHostOnly bool
}

// DefaultCrawlOptions returns defaults: up to 5 same-host pages.
func DefaultCrawlOptions() CrawlOptions {
	return CrawlOptions{Options: DefaultOptions(), MaxPages: 5, SameHostOnly: true}
}

// CaptureSite crawls from startURL, breadth-first, capturing evidence for each
// page up to MaxPages. It reuses one Pilot for all pages.
func CaptureSite(ctx context.Context, pilot *w3pilot.Pilot, startURL string, opts CrawlOptions) (*SiteEvidence, error) {
	if opts.MaxPages <= 0 {
		opts.MaxPages = 1
	}
	start, err := url.Parse(startURL)
	if err != nil {
		return nil, fmt.Errorf("parse start URL: %w", err)
	}

	site := &SiteEvidence{StartURL: startURL}
	seen := map[string]bool{}
	queue := []string{startURL}

	for len(queue) > 0 && len(site.Pages) < opts.MaxPages {
		next := queue[0]
		queue = queue[1:]
		if seen[next] {
			continue
		}
		seen[next] = true

		ev, err := Capture(ctx, pilot, next, opts.Options)
		if err != nil {
			continue // skip unreachable pages; keep crawling
		}
		site.Pages = append(site.Pages, *ev)

		// Enqueue same-host, not-yet-seen links.
		for _, l := range ev.Links {
			u, err := url.Parse(l.Href)
			if err != nil {
				continue
			}
			if opts.SameHostOnly && u.Host != start.Host {
				continue
			}
			clean := u.Scheme + "://" + u.Host + u.Path
			if !seen[clean] {
				queue = append(queue, clean)
			}
		}
	}
	return site, nil
}
