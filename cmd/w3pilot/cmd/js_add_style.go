//nolint:dupl // js commands share similar structure
package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	jsAddStyleTimeout time.Duration
	jsAddStyleFile    bool
)

var jsAddStyleCmd = &cobra.Command{
	Use:   "add-style <css>",
	Short: "Add CSS styles to page",
	Long: `Add a stylesheet to the page.

Examples:
  w3pilot js add-style "body { background: red }"
  w3pilot js add-style ./styles.css --file`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		source := args[0]

		// If --file flag is set, read from file
		if jsAddStyleFile {
			content, err := os.ReadFile(source)
			if err != nil {
				return fmt.Errorf("failed to read style file: %w", err)
			}
			source = string(content)
		}

		ctx, cancel := context.WithTimeout(context.Background(), jsAddStyleTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		if err := pilot.AddStyle(ctx, source); err != nil {
			return fmt.Errorf("add style failed: %w", err)
		}

		fmt.Println("Style added")
		return nil
	},
}

func init() {
	jsCmd.AddCommand(jsAddStyleCmd)
	jsAddStyleCmd.Flags().DurationVar(&jsAddStyleTimeout, "timeout", 10*time.Second, "Timeout")
	jsAddStyleCmd.Flags().BoolVar(&jsAddStyleFile, "file", false, "Read CSS from file")
}
