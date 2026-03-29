//nolint:dupl // grouped command intentionally mirrors flat command for backward compatibility
package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var jsEvalTimeout time.Duration

// JSEvalResult represents the result of JavaScript evaluation.
type JSEvalResult struct {
	Result interface{} `json:"result"`
}

var jsEvalCmd = &cobra.Command{
	Use:   "eval <javascript>",
	Short: "Execute JavaScript",
	Long: `Execute JavaScript on the page and print the result.

Examples:
  w3pilot js eval "document.title"
  w3pilot js eval "document.querySelectorAll('a').length"
  w3pilot js eval "window.location.href"
  w3pilot js eval "document.title" --format json`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		script := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), jsEvalTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		result, err := pilot.Evaluate(ctx, script)
		if err != nil {
			return fmt.Errorf("eval failed: %w", err)
		}

		Output(JSEvalResult{Result: result}, func(data interface{}) string {
			r := data.(JSEvalResult).Result
			// Pretty print result for text mode
			if r == nil {
				return "undefined"
			}
			if s, ok := r.(string); ok {
				return s
			}
			jsonBytes, err := json.MarshalIndent(r, "", "  ")
			if err != nil {
				return fmt.Sprintf("%v", r)
			}
			return string(jsonBytes)
		})
		return nil
	},
}

func init() {
	jsCmd.AddCommand(jsEvalCmd)
	jsEvalCmd.Flags().DurationVar(&jsEvalTimeout, "timeout", 30*time.Second, "Evaluation timeout")
}
