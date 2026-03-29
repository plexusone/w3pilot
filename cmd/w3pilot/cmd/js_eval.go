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

var jsEvalCmd = &cobra.Command{
	Use:   "eval <javascript>",
	Short: "Execute JavaScript",
	Long: `Execute JavaScript on the page and print the result.

Examples:
  w3pilot js eval "document.title"
  w3pilot js eval "document.querySelectorAll('a').length"
  w3pilot js eval "window.location.href"`,
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

		// Pretty print result
		if result == nil {
			fmt.Println("undefined")
		} else if s, ok := result.(string); ok {
			fmt.Println(s)
		} else {
			jsonBytes, err := json.MarshalIndent(result, "", "  ")
			if err != nil {
				fmt.Printf("%v\n", result)
			} else {
				fmt.Println(string(jsonBytes))
			}
		}

		return nil
	},
}

func init() {
	jsCmd.AddCommand(jsEvalCmd)
	jsEvalCmd.Flags().DurationVar(&jsEvalTimeout, "timeout", 30*time.Second, "Evaluation timeout")
}
