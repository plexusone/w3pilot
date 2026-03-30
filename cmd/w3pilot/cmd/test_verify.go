package cmd

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/spf13/cobra"

	w3pilot "github.com/plexusone/w3pilot"
)

var (
	verifyTimeout time.Duration
	verifyExact   bool
)

// VerifyResult represents the result of a verification.
type VerifyResult struct {
	Passed  bool   `json:"passed"`
	Message string `json:"message"`
	Actual  string `json:"actual,omitempty"`
}

// elementVerifyFunc is a function that verifies an element state.
type elementVerifyFunc func(ctx context.Context, elem *w3pilot.Element) error

// runElementVerify is a helper that runs element state verification commands.
// It reduces duplication across verify-visible, verify-enabled, verify-checked, etc.
func runElementVerify(
	selector string,
	verifyFn elementVerifyFunc,
	successMsg string,
	notFoundPassesAs string, // if non-empty, element not found is treated as pass with this message
) error {
	ctx, cancel := context.WithTimeout(context.Background(), verifyTimeout)
	defer cancel()

	pilot := mustGetVibe(ctx)

	elem, err := pilot.Find(ctx, selector, &w3pilot.FindOptions{Timeout: verifyTimeout})
	if err != nil {
		if notFoundPassesAs != "" {
			result := VerifyResult{
				Passed:  true,
				Message: fmt.Sprintf(notFoundPassesAs, selector),
			}
			outputVerifyResult(result)
			return nil
		}
		result := VerifyResult{
			Passed:  false,
			Message: fmt.Sprintf("Element not found: %s", selector),
		}
		outputVerifyResult(result)
		return fmt.Errorf("verification failed")
	}

	verifyErr := verifyFn(ctx, elem)

	result := VerifyResult{Passed: verifyErr == nil}
	if verifyErr != nil {
		var vErr *w3pilot.VerificationError
		if errors.As(verifyErr, &vErr) {
			result.Message = vErr.Message
		} else {
			result.Message = verifyErr.Error()
		}
	} else {
		result.Message = fmt.Sprintf(successMsg, selector)
	}

	outputVerifyResult(result)
	if !result.Passed {
		return fmt.Errorf("verification failed")
	}
	return nil
}

// verifyValueCmd verifies an input element's value
var verifyValueCmd = &cobra.Command{
	Use:   "verify-value <selector> <expected>",
	Short: "Verify input value matches",
	Long: `Verify that an input element's value matches the expected value.

Examples:
  w3pilot test verify-value "#email" "user@example.com"
  w3pilot test verify-value "input[name=username]" "john"`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		selector := args[0]
		expected := args[1]

		ctx, cancel := context.WithTimeout(context.Background(), verifyTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		elem, err := pilot.Find(ctx, selector, &w3pilot.FindOptions{Timeout: verifyTimeout})
		if err != nil {
			result := VerifyResult{
				Passed:  false,
				Message: fmt.Sprintf("Element not found: %s", selector),
			}
			outputVerifyResult(result)
			return fmt.Errorf("verification failed")
		}

		actual, _ := elem.Value(ctx)
		verifyErr := elem.VerifyValue(ctx, expected)

		result := VerifyResult{Passed: verifyErr == nil, Actual: actual}
		if verifyErr != nil {
			var vErr *w3pilot.VerificationError
			if errors.As(verifyErr, &vErr) {
				result.Message = vErr.Message
			} else {
				result.Message = verifyErr.Error()
			}
		} else {
			result.Message = fmt.Sprintf("Value matches: %q", actual)
		}

		outputVerifyResult(result)
		if !result.Passed {
			return fmt.Errorf("verification failed")
		}
		return nil
	},
}

// verifyTextCmd verifies an element's text content
var verifyTextCmd = &cobra.Command{
	Use:   "verify-text <selector> <expected>",
	Short: "Verify element text matches",
	Long: `Verify that an element's text content matches the expected value.

By default, checks if text contains the expected value.
Use --exact for exact match.

Examples:
  w3pilot test verify-text "#heading" "Welcome"
  w3pilot test verify-text ".message" "Success" --exact`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		selector := args[0]
		expected := args[1]

		ctx, cancel := context.WithTimeout(context.Background(), verifyTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		elem, err := pilot.Find(ctx, selector, &w3pilot.FindOptions{Timeout: verifyTimeout})
		if err != nil {
			result := VerifyResult{
				Passed:  false,
				Message: fmt.Sprintf("Element not found: %s", selector),
			}
			outputVerifyResult(result)
			return fmt.Errorf("verification failed")
		}

		actual, _ := elem.Text(ctx)
		verifyErr := elem.VerifyText(ctx, expected, &w3pilot.VerifyTextOptions{Exact: verifyExact})

		result := VerifyResult{Passed: verifyErr == nil, Actual: actual}
		if verifyErr != nil {
			var vErr *w3pilot.VerificationError
			if errors.As(verifyErr, &vErr) {
				result.Message = vErr.Message
			} else {
				result.Message = verifyErr.Error()
			}
		} else {
			result.Message = fmt.Sprintf("Text matches: %q", actual)
		}

		outputVerifyResult(result)
		if !result.Passed {
			return fmt.Errorf("verification failed")
		}
		return nil
	},
}

// verifyVisibleCmd verifies an element is visible
var verifyVisibleCmd = &cobra.Command{
	Use:   "verify-visible <selector>",
	Short: "Verify element is visible",
	Long: `Verify that an element is visible on the page.

Examples:
  w3pilot test verify-visible "#modal"
  w3pilot test verify-visible ".success-message"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runElementVerify(args[0],
			func(ctx context.Context, elem *w3pilot.Element) error { return elem.VerifyVisible(ctx) },
			"Element is visible: %s", "")
	},
}

// verifyHiddenCmd verifies an element is hidden
var verifyHiddenCmd = &cobra.Command{
	Use:   "verify-hidden <selector>",
	Short: "Verify element is hidden",
	Long: `Verify that an element is hidden (not visible) on the page.
An element that doesn't exist is considered hidden.

Examples:
  w3pilot test verify-hidden "#loading-spinner"
  w3pilot test verify-hidden ".error-message"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runElementVerify(args[0],
			func(ctx context.Context, elem *w3pilot.Element) error { return elem.VerifyHidden(ctx) },
			"Element is hidden: %s", "Element is hidden (not found): %s")
	},
}

// verifyEnabledCmd verifies an element is enabled
var verifyEnabledCmd = &cobra.Command{
	Use:   "verify-enabled <selector>",
	Short: "Verify element is enabled",
	Long: `Verify that an element is enabled (not disabled).

Examples:
  w3pilot test verify-enabled "#submit-button"
  w3pilot test verify-enabled "button[type=submit]"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runElementVerify(args[0],
			func(ctx context.Context, elem *w3pilot.Element) error { return elem.VerifyEnabled(ctx) },
			"Element is enabled: %s", "")
	},
}

// verifyDisabledCmd verifies an element is disabled
var verifyDisabledCmd = &cobra.Command{
	Use:   "verify-disabled <selector>",
	Short: "Verify element is disabled",
	Long: `Verify that an element is disabled.

Examples:
  w3pilot test verify-disabled "#submit-button"
  w3pilot test verify-disabled "input[name=readonly-field]"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runElementVerify(args[0],
			func(ctx context.Context, elem *w3pilot.Element) error { return elem.VerifyDisabled(ctx) },
			"Element is disabled: %s", "")
	},
}

// verifyCheckedCmd verifies a checkbox/radio is checked
var verifyCheckedCmd = &cobra.Command{
	Use:   "verify-checked <selector>",
	Short: "Verify checkbox/radio is checked",
	Long: `Verify that a checkbox or radio element is checked.

Examples:
  w3pilot test verify-checked "#agree-checkbox"
  w3pilot test verify-checked "input[name=newsletter]"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runElementVerify(args[0],
			func(ctx context.Context, elem *w3pilot.Element) error { return elem.VerifyChecked(ctx) },
			"Element is checked: %s", "")
	},
}

// verifyUncheckedCmd verifies a checkbox/radio is unchecked
var verifyUncheckedCmd = &cobra.Command{
	Use:   "verify-unchecked <selector>",
	Short: "Verify checkbox/radio is unchecked",
	Long: `Verify that a checkbox or radio element is unchecked.

Examples:
  w3pilot test verify-unchecked "#agree-checkbox"
  w3pilot test verify-unchecked "input[name=newsletter]"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return runElementVerify(args[0],
			func(ctx context.Context, elem *w3pilot.Element) error { return elem.VerifyUnchecked(ctx) },
			"Element is unchecked: %s", "")
	},
}

// outputVerifyResult outputs the verification result in the configured format
func outputVerifyResult(result VerifyResult) {
	Output(result, func(data interface{}) string {
		r := data.(VerifyResult)
		if r.Passed {
			return fmt.Sprintf("PASS: %s", r.Message)
		}
		return fmt.Sprintf("FAIL: %s", r.Message)
	})
}

func init() {
	testCmd.AddCommand(verifyValueCmd)
	testCmd.AddCommand(verifyTextCmd)
	testCmd.AddCommand(verifyVisibleCmd)
	testCmd.AddCommand(verifyHiddenCmd)
	testCmd.AddCommand(verifyEnabledCmd)
	testCmd.AddCommand(verifyDisabledCmd)
	testCmd.AddCommand(verifyCheckedCmd)
	testCmd.AddCommand(verifyUncheckedCmd)

	// Common flags
	verifyValueCmd.Flags().DurationVar(&verifyTimeout, "timeout", 5*time.Second, "Timeout")

	verifyTextCmd.Flags().DurationVar(&verifyTimeout, "timeout", 5*time.Second, "Timeout")
	verifyTextCmd.Flags().BoolVar(&verifyExact, "exact", false, "Require exact match")

	verifyVisibleCmd.Flags().DurationVar(&verifyTimeout, "timeout", 5*time.Second, "Timeout")
	verifyHiddenCmd.Flags().DurationVar(&verifyTimeout, "timeout", 5*time.Second, "Timeout")
	verifyEnabledCmd.Flags().DurationVar(&verifyTimeout, "timeout", 5*time.Second, "Timeout")
	verifyDisabledCmd.Flags().DurationVar(&verifyTimeout, "timeout", 5*time.Second, "Timeout")
	verifyCheckedCmd.Flags().DurationVar(&verifyTimeout, "timeout", 5*time.Second, "Timeout")
	verifyUncheckedCmd.Flags().DurationVar(&verifyTimeout, "timeout", 5*time.Second, "Timeout")
}
