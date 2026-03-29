//nolint:dupl // grouped command intentionally mirrors flat command for backward compatibility
package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/spf13/cobra"
)

var (
	browserLaunchHeadless    bool
	browserLaunchInitScripts []string
)

var browserLaunchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Launch a browser instance",
	Long: `Launch a new browser instance for automation.

The browser will stay open until you run 'w3pilot browser quit' or press Ctrl+C.

Examples:
  w3pilot browser launch              # Launch visible browser
  w3pilot browser launch --headless   # Launch headless browser
  w3pilot browser launch --init-script ./setup.js`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Handle interrupt
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

		vibe, err := launchBrowser(ctx, browserLaunchHeadless)
		if err != nil {
			return err
		}

		// Load and apply init scripts
		for _, scriptPath := range browserLaunchInitScripts {
			content, err := os.ReadFile(scriptPath)
			if err != nil {
				return fmt.Errorf("failed to read init script %s: %w", scriptPath, err)
			}
			if err := vibe.AddInitScript(ctx, string(content)); err != nil {
				return fmt.Errorf("failed to add init script %s: %w", scriptPath, err)
			}
			if verbose {
				fmt.Printf("Loaded init script: %s\n", scriptPath)
			}
		}

		mode := "visible"
		if browserLaunchHeadless {
			mode = "headless"
		}
		fmt.Printf("Browser launched (%s mode)\n", mode)
		if len(browserLaunchInitScripts) > 0 {
			fmt.Printf("Loaded %d init script(s)\n", len(browserLaunchInitScripts))
		}
		fmt.Println("Press Ctrl+C to quit or use 'w3pilot browser quit'")

		// Wait for interrupt
		<-sigCh
		fmt.Println("\nShutting down...")

		if err := vibe.Quit(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
		}
		if err := clearSession(); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: %v\n", err)
		}

		return nil
	},
}

func init() {
	browserCmd.AddCommand(browserLaunchCmd)
	browserLaunchCmd.Flags().BoolVar(&browserLaunchHeadless, "headless", false, "Run browser in headless mode")
	browserLaunchCmd.Flags().StringArrayVar(&browserLaunchInitScripts, "init-script", nil, "JavaScript file to inject before page scripts (can be repeated)")
}
