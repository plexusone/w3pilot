package mcp

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	w3pilot "github.com/plexusone/w3pilot"
	"github.com/plexusone/w3pilot/mcp/report"
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
	_ *mcp.CallToolRequest,
	input VerifyValueInput,
) (*mcp.CallToolResult, VerifyValueOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, VerifyValueOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 5000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	start := time.Now()
	elem, err := pilot.Find(ctx, input.Selector, &w3pilot.FindOptions{Timeout: timeout})

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

	// Get actual value for reporting
	actual, _ := elem.Value(ctx)

	// Use SDK verification method
	verifyErr := elem.VerifyValue(ctx, input.Expected)
	result.DurationMS = time.Since(start).Milliseconds()

	if verifyErr != nil {
		var vErr *w3pilot.VerificationError
		if errors.As(verifyErr, &vErr) {
			result.Status = report.StatusNoGo
			result.Severity = report.SeverityCritical
			result.Error = &report.StepError{
				Type:    vErr.Type,
				Message: vErr.Message,
			}
			result.Context = s.session.CaptureContext(ctx)
			result.Screenshot = s.session.CaptureScreenshot(ctx)
			s.session.RecordStep(result)

			return nil, VerifyValueOutput{
				Passed:  false,
				Actual:  actual,
				Message: vErr.Message,
			}, nil
		}
		// Non-verification error (e.g., get value failed)
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityCritical
		result.Error = &report.StepError{
			Type:     "GetValueError",
			Message:  verifyErr.Error(),
			Selector: input.Selector,
		}
		result.Screenshot = s.session.CaptureScreenshot(ctx)
		s.session.RecordStep(result)
		return nil, VerifyValueOutput{}, verifyErr
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
	_ *mcp.CallToolRequest,
	input VerifyListVisibleInput,
) (*mcp.CallToolResult, VerifyListVisibleOutput, error) {
	pilot, err := s.session.Pilot(ctx)
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

	// Use SDK AssertText for each item
	var found []string
	var missing []string

	for _, item := range input.Items {
		opts := &w3pilot.AssertOptions{
			Selector: input.Selector,
		}
		err := pilot.AssertText(ctx, item, opts)
		if err != nil {
			missing = append(missing, item)
		} else {
			found = append(found, item)
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
	_ *mcp.CallToolRequest,
	input GenerateLocatorInput,
) (*mcp.CallToolResult, GenerateLocatorOutput, error) {
	pilot, err := s.session.Pilot(ctx)
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

	// Use SDK GenerateLocator method
	locatorInfo, err := pilot.GenerateLocator(ctx, input.Selector, &w3pilot.GenerateLocatorOptions{
		Strategy: input.Strategy,
		Timeout:  timeout,
	})
	if err != nil {
		return nil, GenerateLocatorOutput{}, err
	}

	return nil, GenerateLocatorOutput{
		Locator:  locatorInfo.Locator,
		Strategy: locatorInfo.Strategy,
		Metadata: locatorInfo.Metadata,
	}, nil
}

// WaitForSelector tool - waits for an element to reach a specific state

type WaitForSelectorInput struct {
	Selector  string `json:"selector" jsonschema:"CSS selector for the element,required"`
	State     string `json:"state" jsonschema:"State to wait for: attached detached visible hidden (default: visible),enum=attached,enum=detached,enum=visible,enum=hidden"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"Timeout in milliseconds (default: 30000)"`
}

type WaitForSelectorOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleWaitForSelector(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input WaitForSelectorInput,
) (*mcp.CallToolResult, WaitForSelectorOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, WaitForSelectorOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 30000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	if input.State == "" {
		input.State = "visible"
	}

	// Find the element first (for attached/visible states) or wait for condition
	switch input.State {
	case "attached", "visible":
		elem, err := pilot.Find(ctx, input.Selector, &w3pilot.FindOptions{Timeout: timeout})
		if err != nil {
			return nil, WaitForSelectorOutput{}, fmt.Errorf("wait for selector failed: %w", err)
		}
		if input.State == "visible" {
			err = elem.WaitUntil(ctx, "visible", timeout)
			if err != nil {
				return nil, WaitForSelectorOutput{}, fmt.Errorf("wait for visible failed: %w", err)
			}
		}
	case "detached", "hidden":
		err = pilot.WaitForFunction(ctx, fmt.Sprintf(`() => {
			const el = document.querySelector(%q);
			if (%q === "detached") return el === null;
			if (el === null) return true;
			const style = window.getComputedStyle(el);
			return style.display === 'none' || style.visibility === 'hidden' || el.offsetParent === null;
		}`, input.Selector, input.State), timeout)
		if err != nil {
			return nil, WaitForSelectorOutput{}, fmt.Errorf("wait for %s failed: %w", input.State, err)
		}
	default:
		return nil, WaitForSelectorOutput{}, fmt.Errorf("invalid state: %s", input.State)
	}

	return nil, WaitForSelectorOutput{
		Message: fmt.Sprintf("Element %s is %s", input.Selector, input.State),
	}, nil
}

// VerifyText tool - verifies element text matches expected value

type VerifyTextInput struct {
	Selector  string `json:"selector" jsonschema:"CSS selector for the element,required"`
	Expected  string `json:"expected" jsonschema:"Expected text content,required"`
	Exact     bool   `json:"exact" jsonschema:"Require exact match (default: false uses contains)"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"Timeout in milliseconds (default: 5000)"`
}

type VerifyTextOutput struct {
	Passed  bool   `json:"passed"`
	Actual  string `json:"actual"`
	Message string `json:"message"`
}

func (s *Server) handleVerifyText(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input VerifyTextInput,
) (*mcp.CallToolResult, VerifyTextOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, VerifyTextOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 5000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	elem, err := pilot.Find(ctx, input.Selector, &w3pilot.FindOptions{Timeout: timeout})
	if err != nil {
		return nil, VerifyTextOutput{}, fmt.Errorf("element not found: %s", input.Selector)
	}

	// Get actual text for reporting
	actual, _ := elem.Text(ctx)

	// Use SDK VerifyText method
	verifyErr := elem.VerifyText(ctx, input.Expected, &w3pilot.VerifyTextOptions{Exact: input.Exact})

	if verifyErr != nil {
		var vErr *w3pilot.VerificationError
		if errors.As(verifyErr, &vErr) {
			return nil, VerifyTextOutput{
				Passed:  false,
				Actual:  actual,
				Message: vErr.Message,
			}, nil
		}
		return nil, VerifyTextOutput{}, verifyErr
	}

	return nil, VerifyTextOutput{
		Passed:  true,
		Actual:  actual,
		Message: fmt.Sprintf("Text matches: %q", actual),
	}, nil
}

// verifyStateResult holds the result of a state verification.
type verifyStateResult struct {
	Passed  bool
	Message string
}

// elementVerifyFn is a function that verifies an element state.
type elementVerifyFn func(ctx context.Context, elem *w3pilot.Element) error

// verifyElementState is a helper that runs element state verification.
// It reduces duplication across verify-visible, verify-enabled, verify-hidden, verify-disabled handlers.
func (s *Server) verifyElementState(
	ctx context.Context,
	selector string,
	timeoutMS int,
	verifyFn elementVerifyFn,
	successMsg string,
	notFoundPassesAs string, // if non-empty, element not found is treated as pass with this message
) (verifyStateResult, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return verifyStateResult{}, fmt.Errorf("browser not available: %w", err)
	}

	if timeoutMS == 0 {
		timeoutMS = 5000
	}
	timeout := time.Duration(timeoutMS) * time.Millisecond

	elem, err := pilot.Find(ctx, selector, &w3pilot.FindOptions{Timeout: timeout})
	if err != nil {
		if notFoundPassesAs != "" {
			return verifyStateResult{
				Passed:  true,
				Message: fmt.Sprintf(notFoundPassesAs, selector),
			}, nil
		}
		return verifyStateResult{
			Passed:  false,
			Message: fmt.Sprintf("Element not found: %s", selector),
		}, nil
	}

	verifyErr := verifyFn(ctx, elem)
	if verifyErr != nil {
		var vErr *w3pilot.VerificationError
		if errors.As(verifyErr, &vErr) {
			return verifyStateResult{Passed: false, Message: vErr.Message}, nil
		}
		return verifyStateResult{}, verifyErr
	}

	return verifyStateResult{
		Passed:  true,
		Message: fmt.Sprintf(successMsg, selector),
	}, nil
}

// VerifyVisible tool - verifies element is visible

type VerifyVisibleInput struct {
	Selector  string `json:"selector" jsonschema:"CSS selector for the element,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"Timeout in milliseconds (default: 5000)"`
}

type VerifyVisibleOutput struct {
	Passed  bool   `json:"passed"`
	Message string `json:"message"`
}

func (s *Server) handleVerifyVisible(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input VerifyVisibleInput,
) (*mcp.CallToolResult, VerifyVisibleOutput, error) {
	result, err := s.verifyElementState(ctx, input.Selector, input.TimeoutMS,
		func(ctx context.Context, elem *w3pilot.Element) error { return elem.VerifyVisible(ctx) },
		"Element is visible: %s", "")
	if err != nil {
		return nil, VerifyVisibleOutput{}, err
	}
	return nil, VerifyVisibleOutput{Passed: result.Passed, Message: result.Message}, nil
}

// VerifyEnabled tool - verifies element is enabled

type VerifyEnabledInput struct {
	Selector  string `json:"selector" jsonschema:"CSS selector for the element,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"Timeout in milliseconds (default: 5000)"`
}

type VerifyEnabledOutput struct {
	Passed  bool   `json:"passed"`
	Message string `json:"message"`
}

func (s *Server) handleVerifyEnabled(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input VerifyEnabledInput,
) (*mcp.CallToolResult, VerifyEnabledOutput, error) {
	result, err := s.verifyElementState(ctx, input.Selector, input.TimeoutMS,
		func(ctx context.Context, elem *w3pilot.Element) error { return elem.VerifyEnabled(ctx) },
		"Element is enabled: %s", "")
	if err != nil {
		return nil, VerifyEnabledOutput{}, err
	}
	return nil, VerifyEnabledOutput{Passed: result.Passed, Message: result.Message}, nil
}

// VerifyHidden tool - verifies element is hidden

type VerifyHiddenInput struct {
	Selector  string `json:"selector" jsonschema:"CSS selector for the element,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"Timeout in milliseconds (default: 5000)"`
}

type VerifyHiddenOutput struct {
	Passed  bool   `json:"passed"`
	Message string `json:"message"`
}

func (s *Server) handleVerifyHidden(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input VerifyHiddenInput,
) (*mcp.CallToolResult, VerifyHiddenOutput, error) {
	result, err := s.verifyElementState(ctx, input.Selector, input.TimeoutMS,
		func(ctx context.Context, elem *w3pilot.Element) error { return elem.VerifyHidden(ctx) },
		"Element is hidden: %s", "Element is hidden (not found): %s")
	if err != nil {
		return nil, VerifyHiddenOutput{}, err
	}
	return nil, VerifyHiddenOutput{Passed: result.Passed, Message: result.Message}, nil
}

// VerifyDisabled tool - verifies element is disabled

type VerifyDisabledInput struct {
	Selector  string `json:"selector" jsonschema:"CSS selector for the element,required"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"Timeout in milliseconds (default: 5000)"`
}

type VerifyDisabledOutput struct {
	Passed  bool   `json:"passed"`
	Message string `json:"message"`
}

func (s *Server) handleVerifyDisabled(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input VerifyDisabledInput,
) (*mcp.CallToolResult, VerifyDisabledOutput, error) {
	result, err := s.verifyElementState(ctx, input.Selector, input.TimeoutMS,
		func(ctx context.Context, elem *w3pilot.Element) error { return elem.VerifyDisabled(ctx) },
		"Element is disabled: %s", "")
	if err != nil {
		return nil, VerifyDisabledOutput{}, err
	}
	return nil, VerifyDisabledOutput{Passed: result.Passed, Message: result.Message}, nil
}

// VerifyChecked tool - verifies checkbox/radio is checked

type VerifyCheckedInput struct {
	Selector  string `json:"selector" jsonschema:"CSS selector for the checkbox/radio element,required"`
	Checked   bool   `json:"checked" jsonschema:"Expected checked state (default: true)"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"Timeout in milliseconds (default: 5000)"`
}

type VerifyCheckedOutput struct {
	Passed  bool   `json:"passed"`
	Checked bool   `json:"checked"`
	Message string `json:"message"`
}

func (s *Server) handleVerifyChecked(
	ctx context.Context,
	_ *mcp.CallToolRequest,
	input VerifyCheckedInput,
) (*mcp.CallToolResult, VerifyCheckedOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, VerifyCheckedOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 5000
	}
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond

	elem, err := pilot.Find(ctx, input.Selector, &w3pilot.FindOptions{Timeout: timeout})
	if err != nil {
		return nil, VerifyCheckedOutput{
			Passed:  false,
			Message: fmt.Sprintf("Element not found: %s", input.Selector),
		}, nil
	}

	// Get actual checked state for reporting
	actualChecked, _ := elem.IsChecked(ctx)

	// Use SDK VerifyChecked or VerifyUnchecked based on expected state
	var verifyErr error
	if input.Checked {
		verifyErr = elem.VerifyChecked(ctx)
	} else {
		verifyErr = elem.VerifyUnchecked(ctx)
	}

	if verifyErr != nil {
		var vErr *w3pilot.VerificationError
		if errors.As(verifyErr, &vErr) {
			return nil, VerifyCheckedOutput{
				Passed:  false,
				Checked: actualChecked,
				Message: vErr.Message,
			}, nil
		}
		return nil, VerifyCheckedOutput{}, verifyErr
	}

	return nil, VerifyCheckedOutput{
		Passed:  true,
		Checked: actualChecked,
		Message: fmt.Sprintf("Element checked state is %v as expected", actualChecked),
	}, nil
}
