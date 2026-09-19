// Package pagecapture captures normalized, reusable evidence about web pages —
// rendered HTML, a screenshot, visible text, and structural signals (title,
// language, links, headings) — using a w3pilot browser session. It is
// domain-neutral: accessibility evaluation, functional/journey testing, i18n
// checks, and security review all consume the same PageEvidence rather than
// each driving the browser themselves.
package pagecapture

import "time"

// Heading is a heading element captured from the page.
type Heading struct {
	Level int    `json:"level"`
	Text  string `json:"text"`
}

// Link is an anchor captured from the page.
type Link struct {
	Text string `json:"text"`
	Href string `json:"href"`
}

// PageEvidence is normalized evidence for a single page.
type PageEvidence struct {
	URL      string `json:"url"`
	Title    string `json:"title"`
	Language string `json:"language,omitempty"` // <html lang>

	// HTML is the rendered outer HTML, bounded to MaxHTMLBytes.
	HTML string `json:"html,omitempty"`
	// Text is the visible body text, bounded to MaxTextBytes.
	Text string `json:"text,omitempty"`
	// ScreenshotPNG is a full-page screenshot; base64-encoded when serialized.
	ScreenshotPNG []byte `json:"screenshotPng,omitempty"`

	Links    []Link    `json:"links,omitempty"`
	Headings []Heading `json:"headings,omitempty"`

	// Truncated reports whether HTML/Text were clipped to their bounds.
	Truncated bool `json:"truncated,omitempty"`
	// LoadWarning is set when the page-load wait did not fully settle.
	LoadWarning string    `json:"loadWarning,omitempty"`
	CapturedAt  time.Time `json:"capturedAt"`
}

// SiteEvidence is evidence captured across multiple pages of one site.
type SiteEvidence struct {
	StartURL string         `json:"startUrl"`
	Pages    []PageEvidence `json:"pages"`
}

// HasScreenshot reports whether a screenshot was captured (used by consumers to
// decide whether visual criteria can be judged).
func (p PageEvidence) HasScreenshot() bool { return len(p.ScreenshotPNG) > 0 }
