package mcp

import (
	"context"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	vibium "github.com/plexusone/vibium-go"
	"github.com/plexusone/vibium-go/mcp/report"
)

// VerifyValue tool - verifies that an input element has the expected value

type VerifyValueInput struct {
	Selector  string `json:"selector" jsonschema:"CSS selector for the input element,required"`
	Expected  string `json:"expected" jsonschema:"Expected value to verify,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"Timeout in milliseconds (default: 5000)"`
}

type VerifyValueOutput struct {
	Passed  bool   `json:"passed"`
	Actual  string `json:"actual"`
	Message string `json:"message"`
}

func (s *Server) handleVerifyValue(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input VerifyValueInput,
) (*mcp.CallToolResult, VerifyValueOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, VerifyValueOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 5000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	start := time.Now()
	elem, err := vibe.Find(ctx, input.Selector, &vibium.FindOptions{Timeout: timeout})

	result := report.StepResult{
		ID:     s.session.NextStepID("verify_value"),
		Action: "verify_value",
		Args:   map[string]any{"selector": input.Selector, "expected": input.Expected},
	}

	if err != nil {
		result.DurationMS = time.Since(start).Milliseconds()
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityCritical
		result.Error = &report.StepError{
			Type:        "ElementNotFoundError",
			Message:     err.Error(),
			Selector:    input.Selector,
			TimeoutMS:   int64(input.TimeoutMS),
			Suggestions: s.session.FindSimilarSelectors(ctx, input.Selector),
		}
		result.Context = s.session.CaptureContext(ctx)
		result.Screenshot = s.session.CaptureScreenshot(ctx)
		s.session.RecordStep(result)
		return nil, VerifyValueOutput{}, fmt.Errorf("element not found: %s", input.Selector)
	}

	actual, err := elem.Value(ctx)
	result.DurationMS = time.Since(start).Milliseconds()

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityCritical
		result.Error = &report.StepError{
			Type:     "GetValueError",
			Message:  err.Error(),
			Selector: input.Selector,
		}
		result.Screenshot = s.session.CaptureScreenshot(ctx)
		s.session.RecordStep(result)
		return nil, VerifyValueOutput{}, fmt.Errorf("get value failed: %w", err)
	}

	passed := actual == input.Expected

	if !passed {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityCritical
		result.Error = &report.StepError{
			Type:    "VerifyValueFailed",
			Message: fmt.Sprintf("Expected %q but got %q", input.Expected, actual),
		}
		result.Context = s.session.CaptureContext(ctx)
		result.Screenshot = s.session.CaptureScreenshot(ctx)
		s.session.RecordStep(result)

		return nil, VerifyValueOutput{
			Passed:  false,
			Actual:  actual,
			Message: fmt.Sprintf("Value mismatch: expected %q but got %q", input.Expected, actual),
		}, nil
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	result.Result = map[string]any{"actual": actual, "expected": input.Expected}
	s.session.RecordStep(result)

	return nil, VerifyValueOutput{
		Passed:  true,
		Actual:  actual,
		Message: fmt.Sprintf("Value matches: %q", actual),
	}, nil
}

// VerifyListVisible tool - verifies that a list of items are visible on the page

type VerifyListVisibleInput struct {
	Items     []string `json:"items" jsonschema:"List of text items that should be visible on the page,required"`
	Selector  string   `json:"selector" jsonschema:"Optional CSS selector to scope the search"`
	TimeoutMS int      `json:"timeout_ms" jsonschema:"Timeout in milliseconds (default: 5000)"`
}

type VerifyListVisibleOutput struct {
	Passed  bool     `json:"passed"`
	Found   []string `json:"found"`
	Missing []string `json:"missing"`
	Message string   `json:"message"`
}

func (s *Server) handleVerifyListVisible(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input VerifyListVisibleInput,
) (*mcp.CallToolResult, VerifyListVisibleOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, VerifyListVisibleOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 5000
	}

	start := time.Now()

	result := report.StepResult{
		ID:     s.session.NextStepID("verify_list_visible"),
		Action: "verify_list_visible",
		Args:   map[string]any{"items": input.Items, "selector": input.Selector},
	}

	// Build script to check for each item's visibility
	var found []string
	var missing []string

	for _, item := range input.Items {
		var script string
		if input.Selector != "" {
			script = fmt.Sprintf(`
				(function() {
					const el = document.querySelector(%q);
					return el && el.textContent.includes(%q);
				})()
			`, input.Selector, item)
		} else {
			script = fmt.Sprintf(`document.body.textContent.includes(%q)`, item)
		}

		evalResult, err := vibe.Evaluate(ctx, script)
		if err != nil {
			missing = append(missing, item)
			continue
		}

		if visible, ok := evalResult.(bool); ok && visible {
			found = append(found, item)
		} else {
			missing = append(missing, item)
		}
	}

	result.DurationMS = time.Since(start).Milliseconds()

	passed := len(missing) == 0

	if !passed {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityCritical
		result.Error = &report.StepError{
			Type:    "VerifyListVisibleFailed",
			Message: fmt.Sprintf("Missing items: %v", missing),
		}
		result.Context = s.session.CaptureContext(ctx)
		result.Screenshot = s.session.CaptureScreenshot(ctx)
		s.session.RecordStep(result)

		return nil, VerifyListVisibleOutput{
			Passed:  false,
			Found:   found,
			Missing: missing,
			Message: fmt.Sprintf("Found %d of %d items, missing: %v", len(found), len(input.Items), missing),
		}, nil
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	result.Result = map[string]any{"found": found}
	s.session.RecordStep(result)

	return nil, VerifyListVisibleOutput{
		Passed:  true,
		Found:   found,
		Missing: missing,
		Message: fmt.Sprintf("All %d items visible", len(input.Items)),
	}, nil
}

// GenerateLocator tool - generates a locator string for a given element

type GenerateLocatorInput struct {
	Selector  string `json:"selector" jsonschema:"CSS selector for the element,required"`
	Strategy  string `json:"strategy" jsonschema:"Locator strategy: css xpath testid role text (default: css),enum=css,enum=xpath,enum=testid,enum=role,enum=text"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"Timeout in milliseconds (default: 5000)"`
}

type GenerateLocatorOutput struct {
	Locator  string            `json:"locator"`
	Strategy string            `json:"strategy"`
	Metadata map[string]string `json:"metadata"`
}

func (s *Server) handleGenerateLocator(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GenerateLocatorInput,
) (*mcp.CallToolResult, GenerateLocatorOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, GenerateLocatorOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 5000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	if input.Strategy == "" {
		input.Strategy = "css"
	}

	elem, err := vibe.Find(ctx, input.Selector, &vibium.FindOptions{Timeout: timeout})
	if err != nil {
		return nil, GenerateLocatorOutput{}, fmt.Errorf("element not found: %s", input.Selector)
	}

	metadata := make(map[string]string)
	var locator string

	switch input.Strategy {
	case "css":
		// Generate a unique CSS selector
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
		result, err := vibe.Evaluate(ctx, fmt.Sprintf(script, input.Selector))
		if err != nil {
			return nil, GenerateLocatorOutput{}, fmt.Errorf("generate locator failed: %w", err)
		}
		if result != nil {
			locator = fmt.Sprintf("%v", result)
		} else {
			locator = input.Selector
		}

	case "xpath":
		// Generate XPath for the element
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
		result, err := vibe.Evaluate(ctx, fmt.Sprintf(script, input.Selector))
		if err != nil {
			return nil, GenerateLocatorOutput{}, fmt.Errorf("generate locator failed: %w", err)
		}
		if result != nil {
			locator = fmt.Sprintf("%v", result)
		}

	case "testid":
		testID, err := elem.GetAttribute(ctx, "data-testid")
		if err != nil {
			return nil, GenerateLocatorOutput{}, fmt.Errorf("get testid failed: %w", err)
		}
		if testID == "" {
			return nil, GenerateLocatorOutput{}, fmt.Errorf("element has no data-testid attribute")
		}
		locator = fmt.Sprintf("[data-testid=\"%s\"]", testID)
		metadata["testid"] = testID

	case "role":
		role, err := elem.Role(ctx)
		if err != nil {
			return nil, GenerateLocatorOutput{}, fmt.Errorf("get role failed: %w", err)
		}
		if role == "" {
			return nil, GenerateLocatorOutput{}, fmt.Errorf("element has no ARIA role")
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
		if err != nil {
			return nil, GenerateLocatorOutput{}, fmt.Errorf("get text failed: %w", err)
		}
		if text == "" {
			return nil, GenerateLocatorOutput{}, fmt.Errorf("element has no text content")
		}
		// Truncate long text
		if len(text) > 50 {
			text = text[:50]
		}
		locator = fmt.Sprintf("text=%q", text)
		metadata["text"] = text

	default:
		return nil, GenerateLocatorOutput{}, fmt.Errorf("unknown strategy: %s", input.Strategy)
	}

	return nil, GenerateLocatorOutput{
		Locator:  locator,
		Strategy: input.Strategy,
		Metadata: metadata,
	}, nil
}
