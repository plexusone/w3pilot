# Page Capture

The `pagecapture` subpackage captures normalized, reusable evidence about a web
page — rendered HTML, a full-page screenshot, visible text, and structural
signals (title, language, links, headings). It is domain-neutral: accessibility
evaluation, functional/journey testing, i18n checks, and security review all
consume the same `PageEvidence` rather than each driving the browser themselves.

## SPA-aware settling

Client-rendered single-page apps fire the `load` event before their content
exists, and their network may never go idle. `pagecapture` therefore defaults to
a **DOM-stability** settle: it polls a fingerprint of the DOM (element count +
visible text length) until it stops changing across several samples *and* holds
real content, gated behind a **minimum-settle floor** so an early render lull
(before async data arrives) doesn't cause a premature capture.

`SettleLoad` (wait for the `load` event) and `SettleNone` (capture immediately)
are available for static pages.

## API

```go
import "github.com/plexusone/w3pilot/pagecapture"

// Convenience: launch a headless browser, capture one page, close it.
ev, err := pagecapture.CaptureURL(ctx, "https://example.com", pagecapture.DefaultOptions())

// Reuse an existing Pilot.
ev, err := pagecapture.Capture(ctx, pilot, url, pagecapture.DefaultOptions())

// Capture the page the Pilot is already on, without navigating or settling —
// e.g. from an engine that just ran against it (avoids a redundant navigation).
ev := pagecapture.CaptureCurrent(ctx, pilot, url, true /* screenshot */)

// Breadth-first, same-host crawl.
site, err := pagecapture.CaptureSite(ctx, pilot, startURL, pagecapture.DefaultCrawlOptions())
```

## PageEvidence

| Field | Description |
|-------|-------------|
| `URL`, `Title`, `Language` | Page identity and `<html lang>` |
| `HTML` | Rendered outer HTML, bounded to `MaxHTMLBytes` |
| `Text` | Visible body text, bounded to `MaxTextBytes` |
| `ScreenshotPNG` | Full-page screenshot (base64 when serialized) |
| `Links`, `Headings` | Structural signals |
| `Truncated`, `LoadWarning` | Whether content was clipped; whether settling was confirmed |

`SiteEvidence` wraps a slice of `PageEvidence` for multi-page crawls.

> **Note on SPA navigation:** anchor-based crawling only follows `<a href>`
> links. Apps that navigate via client-side routing (buttons, `onClick`) expose
> no anchors to follow — capturing across their states requires interaction
> (e.g. a journey runner driving the app), not a link crawl.
