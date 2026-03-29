package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	w3pilot "github.com/plexusone/w3pilot"
)

var testValidateTimeout time.Duration

// TestValidateResult represents the validation results for output.
type TestValidateResult struct {
	Results []w3pilot.SelectorValidation `json:"results"`
	Summary struct {
		Total   int `json:"total"`
		Found   int `json:"found"`
		Missing int `json:"missing"`
		Visible int `json:"visible"`
	} `json:"summary"`
}

var testValidateCmd = &cobra.Command{
	Use:   "validate-selectors <selector>...",
	Short: "Validate CSS selectors",
	Long: `Validate one or more CSS selectors to check if they exist and are usable.

This command helps verify selectors before using them in automation scripts.
It checks each selector and reports:
  - Whether the element exists
  - How many elements match
  - Whether the first match is visible
  - Whether the first match is enabled
  - Suggestions for similar selectors if not found

Examples:
  w3pilot test validate-selectors "#submit"
  w3pilot test validate-selectors "#login" ".password" "button[type=submit]"
  w3pilot test validate-selectors --format json "#nonexistent"`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), testValidateTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		results, err := pilot.ValidateSelectors(ctx, args)
		if err != nil {
			return fmt.Errorf("validation failed: %w", err)
		}

		// Build output
		output := TestValidateResult{
			Results: results,
		}
		output.Summary.Total = len(results)
		for _, r := range results {
			if r.Found {
				output.Summary.Found++
				if r.Visible {
					output.Summary.Visible++
				}
			} else {
				output.Summary.Missing++
			}
		}

		Output(output, func(data interface{}) string {
			r := data.(TestValidateResult)
			return formatValidateResult(r)
		})
		return nil
	},
}

func formatValidateResult(r TestValidateResult) string {
	var sb strings.Builder

	for _, v := range r.Results {
		if v.Found {
			status := "FOUND"
			if !v.Visible {
				status = "FOUND (hidden)"
			}
			sb.WriteString(fmt.Sprintf("%s %s\n", status, v.Selector))
			sb.WriteString(fmt.Sprintf("  tag: %s, count: %d, visible: %t, enabled: %t\n",
				v.TagName, v.Count, v.Visible, v.Enabled))
		} else {
			sb.WriteString(fmt.Sprintf("NOT FOUND %s\n", v.Selector))
			if len(v.Suggestions) > 0 {
				sb.WriteString(fmt.Sprintf("  suggestions: %s\n", strings.Join(v.Suggestions, ", ")))
			}
		}
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("Summary: %d/%d found, %d visible, %d missing\n",
		r.Summary.Found, r.Summary.Total, r.Summary.Visible, r.Summary.Missing))

	return sb.String()
}

func init() {
	testCmd.AddCommand(testValidateCmd)
	testValidateCmd.Flags().DurationVar(&testValidateTimeout, "timeout", 10*time.Second, "Timeout")
}
