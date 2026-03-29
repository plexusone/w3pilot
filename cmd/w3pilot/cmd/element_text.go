//nolint:dupl // element commands share similar structure
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var elementTextTimeout time.Duration

var elementTextCmd = &cobra.Command{
	Use:   "text <selector>",
	Short: "Get element text content",
	Long: `Get the text content of an element.

Examples:
  w3pilot element text "#header"
  w3pilot element text ".message"`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		selector := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), elementTextTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		el, err := pilot.Find(ctx, selector, nil)
		if err != nil {
			return fmt.Errorf("element not found: %w", err)
		}

		text, err := el.Text(ctx)
		if err != nil {
			return fmt.Errorf("failed to get text: %w", err)
		}

		fmt.Println(text)
		return nil
	},
}

func init() {
	elementCmd.AddCommand(elementTextCmd)
	elementTextCmd.Flags().DurationVar(&elementTextTimeout, "timeout", 10*time.Second, "Timeout")
}
