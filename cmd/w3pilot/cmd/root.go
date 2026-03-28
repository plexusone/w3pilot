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
	Long: `WebPilot is a browser automation tool that provides:

  - MCP server for AI-assisted browser automation
  - CLI commands for scripted browser control
  - YAML script runner for batch operations

Examples:
  # Start MCP server
  w3pilot mcp --headless

  # Launch browser and run commands
  w3pilot launch --headless
  w3pilot go https://example.com
  w3pilot click "#submit"
  w3pilot screenshot output.png
  w3pilot quit

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
