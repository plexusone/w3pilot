package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var elementAttrTimeout time.Duration

// ElementAttrResult represents the result of getting element attribute.
type ElementAttrResult struct {
	Selector  string `json:"selector"`
	Attribute string `json:"attribute"`
	Value     string `json:"value"`
}

var elementAttrCmd = &cobra.Command{
	Use:   "attr <selector> <attribute>",
	Short: "Get element attribute value",
	Long: `Get the value of an element's attribute.

Examples:
  w3pilot element attr "#link" href
  w3pilot element attr "img" src
  w3pilot element attr "#btn" data-id --format json`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		selector := args[0]
		attr := args[1]

		ctx, cancel := context.WithTimeout(context.Background(), elementAttrTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		el, err := pilot.Find(ctx, selector, nil)
		if err != nil {
			return fmt.Errorf("element not found: %w", err)
		}

		value, err := el.GetAttribute(ctx, attr)
		if err != nil {
			return fmt.Errorf("failed to get attribute: %w", err)
		}

		Output(ElementAttrResult{Selector: selector, Attribute: attr, Value: value}, func(data interface{}) string {
			return data.(ElementAttrResult).Value
		})
		return nil
	},
}

func init() {
	elementCmd.AddCommand(elementAttrCmd)
	elementAttrCmd.Flags().DurationVar(&elementAttrTimeout, "timeout", 10*time.Second, "Timeout")
}
