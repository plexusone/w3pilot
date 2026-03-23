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
	launchHeadless    bool
	launchInitScripts []string
)

var launchCmd = &cobra.Command{
	Use:   "launch",
	Short: "Launch a browser instance",
	Long: `Launch a new browser instance for automation.

The browser will stay open until you run 'webpilot quit' or press Ctrl+C.

Examples:
  webpilot launch              # Launch visible browser
  webpilot launch --headless   # Launch headless browser
  webpilot launch --init-script ./setup.js --init-script ./mock-api.js`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Handle interrupt
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

		vibe, err := launchBrowser(ctx, launchHeadless)
		if err != nil {
			return err
		}

		// Load and apply init scripts
		for _, scriptPath := range launchInitScripts {
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
		if launchHeadless {
			mode = "headless"
		}
		fmt.Printf("Browser launched (%s mode)\n", mode)
		if len(launchInitScripts) > 0 {
			fmt.Printf("Loaded %d init script(s)\n", len(launchInitScripts))
		}
		fmt.Println("Press Ctrl+C to quit or use 'webpilot quit'")

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
	rootCmd.AddCommand(launchCmd)
	launchCmd.Flags().BoolVar(&launchHeadless, "headless", false, "Run browser in headless mode")
	launchCmd.Flags().StringArrayVar(&launchInitScripts, "init-script", nil, "JavaScript file to inject before page scripts (can be repeated)")
}
