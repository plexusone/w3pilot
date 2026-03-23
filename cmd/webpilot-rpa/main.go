package main

import (
	"os"

	"github.com/plexusone/webpilot/cmd/webpilot-rpa/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
