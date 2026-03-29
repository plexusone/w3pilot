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
	jsAddScriptTimeout time.Duration
	jsAddScriptFile    bool
)

var jsAddScriptCmd = &cobra.Command{
	Use:   "add-script <source>",
	Short: "Add script to page",
	Long: `Add a script that will be evaluated in the page context.

Examples:
  w3pilot js add-script "console.log('Hello')"
  w3pilot js add-script ./script.js --file`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		source := args[0]

		// If --file flag is set, read from file
		if jsAddScriptFile {
			content, err := os.ReadFile(source)
			if err != nil {
				return fmt.Errorf("failed to read script file: %w", err)
			}
			source = string(content)
		}

		ctx, cancel := context.WithTimeout(context.Background(), jsAddScriptTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		if err := pilot.AddScript(ctx, source); err != nil {
			return fmt.Errorf("add script failed: %w", err)
		}

		fmt.Println("Script added")
		return nil
	},
}

func init() {
	jsCmd.AddCommand(jsAddScriptCmd)
	jsAddScriptCmd.Flags().DurationVar(&jsAddScriptTimeout, "timeout", 10*time.Second, "Timeout")
	jsAddScriptCmd.Flags().BoolVar(&jsAddScriptFile, "file", false, "Read script from file")
}
