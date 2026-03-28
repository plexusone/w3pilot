// Package mcp provides an MCP (Model Context Protocol) server for browser automation.
package mcp

import "time"

// Config holds server configuration.
type Config struct {
	// Headless runs the browser without a GUI.
	Headless bool

	// Project is the project name for reports.
	Project string

	// DefaultTimeout is the default timeout for browser operations.
	DefaultTimeout time.Duration

	// InitScripts are JavaScript files to inject before any page scripts.
	// Each string is the content of a script (not a file path).
	InitScripts []string
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Headless:       true,
		Project:        "w3pilot-tests",
		DefaultTimeout: 30 * time.Second,
	}
}
