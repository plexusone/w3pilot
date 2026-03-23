// Command webpilot provides a CLI for browser automation.
package main

import (
	"os"

	"github.com/plexusone/webpilot/cmd/webpilot/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
