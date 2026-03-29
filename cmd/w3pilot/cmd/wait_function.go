package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var waitFunctionTimeout time.Duration

var waitFunctionCmd = &cobra.Command{
	Use:   "function <javascript>",
	Short: "Wait for JavaScript condition",
	Long: `Wait for a JavaScript expression to return truthy.

Examples:
  w3pilot wait function "window.ready === true"
  w3pilot wait function "document.querySelectorAll('.item').length > 5"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		fn := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), waitFunctionTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		if err := pilot.WaitForFunction(ctx, fn, waitFunctionTimeout); err != nil {
			return fmt.Errorf("wait for function failed: %w", err)
		}

		fmt.Println("Condition satisfied")
		return nil
	},
}

func init() {
	waitCmd.AddCommand(waitFunctionCmd)
	waitFunctionCmd.Flags().DurationVar(&waitFunctionTimeout, "timeout", 30*time.Second, "Timeout")
}
