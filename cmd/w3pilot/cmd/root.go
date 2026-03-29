package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	// Global flags
	sessionFile string
	verbose     bool
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "w3pilot",
	Short: "Browser automation CLI",
	Long: `W3Pilot is a browser automation tool that provides:

  - MCP server for AI-assisted browser automation
  - CLI commands for scripted browser control
  - YAML script runner for batch operations

Examples:
  # Start MCP server
  w3pilot mcp --headless

  # Browser lifecycle (grouped commands)
  w3pilot browser launch --headless
  w3pilot browser quit

  # Page navigation and management
  w3pilot page navigate https://example.com
  w3pilot page back
  w3pilot page screenshot output.png
  w3pilot page title

  # Element interactions
  w3pilot element click "#submit"
  w3pilot element fill "#email" "test@example.com"
  w3pilot element text "#result"

  # Wait conditions
  w3pilot wait selector "#modal"
  w3pilot wait url "**/success"

  # JavaScript execution
  w3pilot js eval "document.title"

  # Run a script file
  w3pilot run test.yaml`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	rootCmd.PersistentFlags().StringVar(&sessionFile, "session", "", "Session file path (default: ~/.w3pilot/session.json)")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Verbose output")
}

// getSessionPath returns the session file path
func getSessionPath() string {
	if sessionFile != "" {
		return sessionFile
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ".w3pilot-session.json"
	}
	return fmt.Sprintf("%s/.w3pilot/session.json", home)
}
