package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var jsInitScriptTimeout time.Duration

var jsInitScriptCmd = &cobra.Command{
	Use:   "init-script <file>",
	Short: "Add init script (runs before page scripts)",
	Long: `Add a script that will be evaluated in every page before any page scripts.

This is useful for mocking APIs, injecting test helpers, or setting up authentication.

Examples:
  w3pilot js init-script ./setup.js
  w3pilot js init-script ./mock-api.js`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		scriptPath := args[0]

		content, err := os.ReadFile(scriptPath)
		if err != nil {
			return fmt.Errorf("failed to read script file: %w", err)
		}

		ctx, cancel := context.WithTimeout(context.Background(), jsInitScriptTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		if err := pilot.AddInitScript(ctx, string(content)); err != nil {
			return fmt.Errorf("add init script failed: %w", err)
		}

		fmt.Printf("Init script added: %s\n", scriptPath)
		return nil
	},
}

func init() {
	jsCmd.AddCommand(jsInitScriptCmd)
	jsInitScriptCmd.Flags().DurationVar(&jsInitScriptTimeout, "timeout", 10*time.Second, "Timeout")
}
