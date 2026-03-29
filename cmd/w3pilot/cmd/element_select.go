package cmd

import (
	"context"
	"fmt"
	"time"

	w3pilot "github.com/plexusone/w3pilot"
	"github.com/spf13/cobra"
)

var elementSelectTimeout time.Duration

var elementSelectCmd = &cobra.Command{
	Use:   "select <selector> <value>",
	Short: "Select an option from a dropdown",
	Long: `Select an option from a select element by value.

Examples:
  w3pilot element select "#country" "US"
  w3pilot element select "select[name='size']" "large"`,
	Args: cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		selector := args[0]
		value := args[1]

		ctx, cancel := context.WithTimeout(context.Background(), elementSelectTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		el, err := pilot.Find(ctx, selector, nil)
		if err != nil {
			return fmt.Errorf("element not found: %w", err)
		}

		selectValues := w3pilot.SelectOptionValues{Values: []string{value}}
		if err := el.SelectOption(ctx, selectValues, nil); err != nil {
			return fmt.Errorf("select failed: %w", err)
		}

		fmt.Printf("Selected %s in: %s\n", value, selector)
		return nil
	},
}

func init() {
	elementCmd.AddCommand(elementSelectCmd)
	elementSelectCmd.Flags().DurationVar(&elementSelectTimeout, "timeout", 10*time.Second, "Timeout")
}
