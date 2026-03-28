// Command w3pilot provides a CLI for browser automation.
package main

import (
	"os"

	"github.com/plexusone/w3pilot/cmd/w3pilot/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
