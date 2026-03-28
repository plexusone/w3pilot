package main

import (
	"os"

	"github.com/plexusone/w3pilot/cmd/w3pilot-rpa/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
