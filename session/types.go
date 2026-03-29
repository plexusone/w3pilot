// Package session provides shared browser session management for CLI and MCP.
package session

import (
	"time"

	"github.com/plexusone/w3pilot"
)

// Info stores information about a running browser session for persistence.
type Info struct {
	// WebSocketURL is the BiDi WebSocket endpoint for reconnection.
	WebSocketURL string `json:"websocket_url,omitempty"`

	// ClickerPort is the port clicker is listening on (WebSocket mode).
	ClickerPort int `json:"clicker_port,omitempty"`

	// CDPPort is the Chrome DevTools Protocol port.
	CDPPort int `json:"cdp_port,omitempty"`

	// UserDataDir is the Chrome user data directory.
	UserDataDir string `json:"user_data_dir,omitempty"`

	// Headless indicates if the browser is running in headless mode.
	Headless bool `json:"headless"`

	// PID is the browser process ID.
	PID int `json:"pid,omitempty"`

	// ClickerPID is the clicker process ID (WebSocket mode).
	ClickerPID int `json:"clicker_pid,omitempty"`

	// LaunchedAt is when the session was started.
	LaunchedAt time.Time `json:"launched_at"`

	// InitScripts are JavaScript files injected before page scripts.
	InitScripts []string `json:"init_scripts,omitempty"`
}

// Config holds session configuration options.
type Config struct {
	// Headless runs the browser without a visible window.
	Headless bool

	// DefaultTimeout for browser operations.
	DefaultTimeout time.Duration

	// InitScripts are JavaScript files to inject before page scripts.
	InitScripts []string

	// SessionFile is the path to the session persistence file.
	// Defaults to ~/.w3pilot/session.json
	SessionFile string

	// AutoReconnect attempts to reconnect to an existing session on launch.
	AutoReconnect bool
}

// DefaultConfig returns a Config with sensible defaults.
func DefaultConfig() Config {
	return Config{
		Headless:       false,
		DefaultTimeout: 30 * time.Second,
		AutoReconnect:  true,
	}
}

// Status represents the current session state.
type Status string

const (
	// StatusDisconnected means no browser is connected.
	StatusDisconnected Status = "disconnected"

	// StatusConnected means a browser is connected and ready.
	StatusConnected Status = "connected"

	// StatusReconnecting means attempting to reconnect to an existing session.
	StatusReconnecting Status = "reconnecting"
)

// State provides current session state information.
type State struct {
	Status Status         `json:"status"`
	Info   *Info          `json:"info,omitempty"`
	Pilot  *w3pilot.Pilot `json:"-"`
	Error  error          `json:"error,omitempty"`
}
