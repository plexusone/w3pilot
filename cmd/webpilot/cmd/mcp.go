package cmd

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/plexusone/webpilot/mcp"
	"github.com/spf13/cobra"
)

var (
	mcpHeadless       bool
	mcpDefaultTimeout time.Duration
	mcpProject        string
	mcpInitScripts    []string
)

var mcpCmd = &cobra.Command{
	Use:   "mcp",
	Short: "Start MCP server",
	Long: `Start the Vibium MCP (Model Context Protocol) server.

The MCP server provides browser automation tools for AI assistants.
It communicates via stdio using the MCP protocol.

Examples:
  webpilot mcp
  webpilot mcp --headless
  webpilot mcp --timeout 60s
  webpilot mcp --init-script ./mock-api.js`,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Handle interrupt
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

		go func() {
			<-sigCh
			fmt.Fprintln(os.Stderr, "\nShutting down MCP server...")
			cancel()
		}()

		// Load init scripts from files
		var initScripts []string
		for _, scriptPath := range mcpInitScripts {
			content, err := os.ReadFile(scriptPath)
			if err != nil {
				return fmt.Errorf("failed to read init script %s: %w", scriptPath, err)
			}
			initScripts = append(initScripts, string(content))
			if verbose {
				fmt.Fprintf(os.Stderr, "Loaded init script: %s\n", scriptPath)
			}
		}

		config := mcp.Config{
			Headless:       mcpHeadless,
			DefaultTimeout: mcpDefaultTimeout,
			Project:        mcpProject,
			InitScripts:    initScripts,
		}

		server := mcp.NewServer(config)
		defer func() {
			if err := server.Close(context.Background()); err != nil {
				fmt.Fprintf(os.Stderr, "Warning: cleanup error: %v\n", err)
			}
		}()

		if verbose {
			fmt.Fprintln(os.Stderr, "Starting Vibium MCP server...")
			if mcpHeadless {
				fmt.Fprintln(os.Stderr, "Mode: headless")
			}
			if len(initScripts) > 0 {
				fmt.Fprintf(os.Stderr, "Init scripts: %d\n", len(initScripts))
			}
		}

		return server.Run(ctx)
	},
}

func init() {
	rootCmd.AddCommand(mcpCmd)
	mcpCmd.Flags().BoolVar(&mcpHeadless, "headless", false, "Run browser in headless mode")
	mcpCmd.Flags().DurationVar(&mcpDefaultTimeout, "timeout", 30*time.Second, "Default timeout for operations")
	mcpCmd.Flags().StringVar(&mcpProject, "project", "", "Project name for test reports")
	mcpCmd.Flags().StringArrayVar(&mcpInitScripts, "init-script", nil, "JavaScript file to inject before page scripts (can be repeated)")
}
