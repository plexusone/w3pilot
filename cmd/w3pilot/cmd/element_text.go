//nolint:dupl // element commands share similar structure
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var elementTextTimeout time.Duration

// ElementTextResult represents the result of getting element text.
type ElementTextResult struct {
	Selector string `json:"selector"`
	Text     string `json:"text"`
}

var elementTextCmd = &cobra.Command{
	Use:   "text <selector>",
	Short: "Get element text content",
	Long: `Get the text content of an element.

Examples:
  w3pilot element text "#header"
  w3pilot element text ".message"
  w3pilot element text "#result" --format json`,
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

		Output(ElementTextResult{Selector: selector, Text: text}, func(data interface{}) string {
			return data.(ElementTextResult).Text
		})
		return nil
	},
}

func init() {
	elementCmd.AddCommand(elementTextCmd)
	elementTextCmd.Flags().DurationVar(&elementTextTimeout, "timeout", 10*time.Second, "Timeout")
}
