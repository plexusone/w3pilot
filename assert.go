package w3pilot

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// AssertOptions configures assertion behavior.
type AssertOptions struct {
	// Timeout for finding elements or waiting for conditions.
	Timeout time.Duration

	// Selector to scope the search (for AssertText).
	Selector string
}

// AssertionError is returned when an assertion fails.
type AssertionError struct {
	Type     string `json:"type"`
	Message  string `json:"message"`
	Expected string `json:"expected,omitempty"`
	Actual   string `json:"actual,omitempty"`
	Selector string `json:"selector,omitempty"`
}

func (e *AssertionError) Error() string {
	return e.Message
}

// AssertText asserts that the specified text exists on the page.
func (p *Pilot) AssertText(ctx context.Context, text string, opts *AssertOptions) error {
	if opts == nil {
		opts = &AssertOptions{}
	}
	if opts.Timeout == 0 {
		opts.Timeout = 5 * time.Second
	}

	var script string
	if opts.Selector != "" {
		script = fmt.Sprintf(`
			(function() {
				const el = document.querySelector(%q);
				return el && el.textContent.includes(%q);
			})()
		`, opts.Selector, text)
	} else {
		script = fmt.Sprintf(`document.body.textContent.includes(%q)`, text)
	}

	result, err := p.Evaluate(ctx, script)
	if err != nil {
		return fmt.Errorf("assertion evaluation failed: %w", err)
	}

	found, ok := result.(bool)
	if !ok || !found {
		return &AssertionError{
			Type:     "AssertTextFailed",
			Message:  fmt.Sprintf("text %q not found on page", text),
			Expected: text,
			Selector: opts.Selector,
		}
	}

	return nil
}

// AssertElement asserts that an element matching the selector exists.
func (p *Pilot) AssertElement(ctx context.Context, selector string, opts *AssertOptions) error {
	if opts == nil {
		opts = &AssertOptions{}
	}
	if opts.Timeout == 0 {
		opts.Timeout = 5 * time.Second
	}

	_, err := p.Find(ctx, selector, &FindOptions{Timeout: opts.Timeout})
	if err != nil {
		return &AssertionError{
			Type:     "AssertElementFailed",
			Message:  fmt.Sprintf("element not found: %s", selector),
			Selector: selector,
		}
	}

	return nil
}

// AssertURL asserts that the current URL matches the expected pattern.
// The pattern can be an exact URL, a glob pattern (with *), or a regex (wrapped in /).
func (p *Pilot) AssertURL(ctx context.Context, pattern string, _ *AssertOptions) error {
	currentURL, err := p.URL(ctx)
	if err != nil {
		return fmt.Errorf("failed to get current URL: %w", err)
	}

	matched := matchURLPattern(currentURL, pattern)
	if !matched {
		return &AssertionError{
			Type:     "AssertURLFailed",
			Message:  fmt.Sprintf("URL does not match: expected %q, got %q", pattern, currentURL),
			Expected: pattern,
			Actual:   currentURL,
		}
	}

	return nil
}

// matchURLPattern checks if the URL matches the pattern.
// Supports exact match, glob patterns (*), and regex (wrapped in /).
func matchURLPattern(url, pattern string) bool {
	// Check for regex pattern (wrapped in /)
	if strings.HasPrefix(pattern, "/") && strings.HasSuffix(pattern, "/") && len(pattern) > 2 {
		regexPattern := pattern[1 : len(pattern)-1]
		re, err := regexp.Compile(regexPattern)
		if err != nil {
			return false
		}
		return re.MatchString(url)
	}

	// Check for glob pattern (contains *)
	if strings.Contains(pattern, "*") {
		// Convert glob to regex
		regexPattern := "^" + regexp.QuoteMeta(pattern) + "$"
		regexPattern = strings.ReplaceAll(regexPattern, `\*\*`, ".*")
		regexPattern = strings.ReplaceAll(regexPattern, `\*`, "[^/]*")
		re, err := regexp.Compile(regexPattern)
		if err != nil {
			return false
		}
		return re.MatchString(url)
	}

	// Exact match
	return url == pattern
}

// LocatorInfo contains information about a generated locator.
type LocatorInfo struct {
	// Locator is the generated locator string.
	Locator string `json:"locator"`

	// Strategy is the locator strategy used (css, xpath, testid, role, text).
	Strategy string `json:"strategy"`

	// Metadata contains additional information about the element.
	Metadata map[string]string `json:"metadata,omitempty"`
}

// GenerateLocatorOptions configures locator generation.
type GenerateLocatorOptions struct {
	// Strategy specifies the locator strategy: css, xpath, testid, role, text.
	// Default is "css".
	Strategy string

	// Timeout for finding the element.
	Timeout time.Duration
}

// GenerateLocator generates a robust locator for the element at the given selector.
func (p *Pilot) GenerateLocator(ctx context.Context, selector string, opts *GenerateLocatorOptions) (*LocatorInfo, error) {
	if opts == nil {
		opts = &GenerateLocatorOptions{}
	}
	if opts.Strategy == "" {
		opts.Strategy = "css"
	}
	if opts.Timeout == 0 {
		opts.Timeout = 5 * time.Second
	}

	elem, err := p.Find(ctx, selector, &FindOptions{Timeout: opts.Timeout})
	if err != nil {
		return nil, fmt.Errorf("element not found: %s", selector)
	}

	metadata := make(map[string]string)
	var locator string

	switch opts.Strategy {
	case "css":
		locator, err = p.generateCSSLocator(ctx, selector)
		if err != nil {
			return nil, err
		}

	case "xpath":
		locator, err = p.generateXPathLocator(ctx, selector)
		if err != nil {
			return nil, err
		}

	case "testid":
		testID, err := elem.GetAttribute(ctx, "data-testid")
		if err != nil || testID == "" {
			return nil, fmt.Errorf("element has no data-testid attribute")
		}
		locator = fmt.Sprintf("[data-testid=\"%s\"]", testID)
		metadata["testid"] = testID

	case "role":
		role, err := elem.Role(ctx)
		if err != nil || role == "" {
			return nil, fmt.Errorf("element has no ARIA role")
		}
		label, _ := elem.Label(ctx)
		if label != "" {
			locator = fmt.Sprintf("role=%s[name=%q]", role, label)
			metadata["label"] = label
		} else {
			locator = fmt.Sprintf("role=%s", role)
		}
		metadata["role"] = role

	case "text":
		text, err := elem.Text(ctx)
		if err != nil || text == "" {
			return nil, fmt.Errorf("element has no text content")
		}
		// Truncate long text
		if len(text) > 50 {
			text = text[:50]
		}
		locator = fmt.Sprintf("text=%q", text)
		metadata["text"] = text

	default:
		return nil, fmt.Errorf("unknown strategy: %s", opts.Strategy)
	}

	return &LocatorInfo{
		Locator:  locator,
		Strategy: opts.Strategy,
		Metadata: metadata,
	}, nil
}

// generateCSSLocator generates a unique CSS selector for the element.
func (p *Pilot) generateCSSLocator(ctx context.Context, selector string) (string, error) {
	script := `
		(function(selector) {
			const el = document.querySelector(selector);
			if (!el) return null;

			// Try to generate a unique selector
			// Priority: id > data-testid > class combination > tag with index

			if (el.id) {
				return '#' + CSS.escape(el.id);
			}

			if (el.dataset.testid) {
				return '[data-testid="' + el.dataset.testid + '"]';
			}

			// Generate path from element
			let path = [];
			let current = el;
			while (current && current.nodeType === Node.ELEMENT_NODE) {
				let selector = current.tagName.toLowerCase();
				if (current.id) {
					selector = '#' + CSS.escape(current.id);
					path.unshift(selector);
					break;
				}

				let sibling = current;
				let nth = 1;
				while (sibling = sibling.previousElementSibling) {
					if (sibling.tagName === current.tagName) nth++;
				}

				if (nth > 1 || current.nextElementSibling?.tagName === current.tagName) {
					selector += ':nth-of-type(' + nth + ')';
				}

				path.unshift(selector);
				current = current.parentElement;
			}

			return path.join(' > ');
		})(%q)
	`
	result, err := p.Evaluate(ctx, fmt.Sprintf(script, selector))
	if err != nil {
		return "", fmt.Errorf("generate CSS locator failed: %w", err)
	}
	if result == nil {
		return selector, nil
	}
	return fmt.Sprintf("%v", result), nil
}

// generateXPathLocator generates an XPath selector for the element.
func (p *Pilot) generateXPathLocator(ctx context.Context, selector string) (string, error) {
	script := `
		(function(selector) {
			const el = document.querySelector(selector);
			if (!el) return null;

			if (el.id) {
				return '//*[@id="' + el.id + '"]';
			}

			let path = [];
			let current = el;
			while (current && current.nodeType === Node.ELEMENT_NODE) {
				let tag = current.tagName.toLowerCase();
				let sibling = current;
				let index = 1;
				while (sibling = sibling.previousElementSibling) {
					if (sibling.tagName.toLowerCase() === tag) index++;
				}
				path.unshift(tag + '[' + index + ']');
				current = current.parentElement;
			}

			return '/' + path.join('/');
		})(%q)
	`
	result, err := p.Evaluate(ctx, fmt.Sprintf(script, selector))
	if err != nil {
		return "", fmt.Errorf("generate XPath locator failed: %w", err)
	}
	if result == nil {
		return "", fmt.Errorf("element not found")
	}
	return fmt.Sprintf("%v", result), nil
}
