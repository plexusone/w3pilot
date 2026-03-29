// Package session provides shared browser session management for CLI and MCP.
package session

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/plexusone/w3pilot"
)

// Manager handles browser session lifecycle with optional persistence.
// It provides lazy launch, reconnection, and graceful shutdown.
type Manager struct {
	mu     sync.Mutex
	pilot  *w3pilot.Pilot
	config Config
	state  State
}

// NewManager creates a new session manager with the given config.
func NewManager(config Config) *Manager {
	if config.DefaultTimeout == 0 {
		config.DefaultTimeout = 30 * time.Second
	}
	return &Manager{
		config: config,
		state: State{
			Status: StatusDisconnected,
		},
	}
}

// Pilot returns the browser controller, launching or reconnecting as needed.
func (m *Manager) Pilot(ctx context.Context) (*w3pilot.Pilot, error) {
	if err := m.LaunchIfNeeded(ctx); err != nil {
		return nil, err
	}
	return m.pilot, nil
}

// LaunchIfNeeded launches a new browser or reconnects to an existing session.
func (m *Manager) LaunchIfNeeded(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Already connected
	if m.pilot != nil && !m.pilot.IsClosed() {
		return nil
	}

	// Try to reconnect if auto-reconnect is enabled
	if m.config.AutoReconnect {
		if err := m.reconnectLocked(ctx); err == nil {
			return nil
		}
		// Reconnection failed, proceed to launch new browser
	}

	return m.launchLocked(ctx)
}

// launchLocked launches a new browser (caller must hold mutex).
func (m *Manager) launchLocked(ctx context.Context) error {
	opts := &w3pilot.LaunchOptions{
		Headless:     m.config.Headless,
		UseWebSocket: true, // Required for reconnection support
	}

	var err error
	m.pilot, err = w3pilot.Browser.Launch(ctx, opts)
	if err != nil {
		m.state = State{
			Status: StatusDisconnected,
			Error:  err,
		}
		return fmt.Errorf("failed to launch browser: %w", err)
	}

	// Apply init scripts
	for _, script := range m.config.InitScripts {
		if err := m.pilot.AddInitScript(ctx, script); err != nil {
			return fmt.Errorf("failed to add init script: %w", err)
		}
	}

	// Build session info for persistence
	info := &Info{
		CDPPort:     m.pilot.CDPPort(),
		Headless:    m.config.Headless,
		LaunchedAt:  time.Now(),
		InitScripts: m.config.InitScripts,
	}

	// Get WebSocket URL and port from clicker if available
	if clicker := m.pilot.Clicker(); clicker != nil {
		info.WebSocketURL = clicker.WebSocketURL()
		info.ClickerPort = clicker.Port()
		if clicker.Process() != nil {
			info.ClickerPID = clicker.Process().Pid
		}
	}

	// Update state
	m.state = State{
		Status: StatusConnected,
		Pilot:  m.pilot,
		Info:   info,
	}

	// Persist session info if session file is configured
	if m.config.SessionFile != "" || m.config.AutoReconnect {
		if err := m.persistLocked(); err != nil {
			// Log but don't fail the launch
			fmt.Fprintf(os.Stderr, "warning: failed to persist session: %v\n", err)
		}
	}

	return nil
}

// reconnectLocked attempts to reconnect to an existing session (caller must hold mutex).
func (m *Manager) reconnectLocked(ctx context.Context) error {
	m.state.Status = StatusReconnecting

	// Load saved session info
	info, err := Load(m.config.SessionFile)
	if err != nil || info == nil {
		return fmt.Errorf("no saved session to reconnect")
	}

	// Check if the clicker process is still running
	if info.ClickerPID > 0 {
		process, err := os.FindProcess(info.ClickerPID)
		if err != nil {
			_ = Clear(m.config.SessionFile)
			return fmt.Errorf("clicker process not found: %w", err)
		}
		// On Unix, FindProcess always succeeds; send signal 0 to check if process exists
		if err := process.Signal(os.Signal(nil)); err != nil {
			// Process no longer running
			_ = Clear(m.config.SessionFile)
			return fmt.Errorf("clicker process no longer running")
		}
	}

	// Try to connect via WebSocket URL
	if info.WebSocketURL != "" {
		pilot, err := w3pilot.Connect(ctx, info.WebSocketURL)
		if err != nil {
			_ = Clear(m.config.SessionFile)
			return fmt.Errorf("failed to reconnect via WebSocket: %w", err)
		}
		m.pilot = pilot
		m.state = State{
			Status: StatusConnected,
			Pilot:  pilot,
			Info:   info,
		}
		return nil
	}

	// Fallback: try to construct WebSocket URL from port
	if info.ClickerPort > 0 {
		wsURL := fmt.Sprintf("ws://localhost:%d", info.ClickerPort)
		pilot, err := w3pilot.Connect(ctx, wsURL)
		if err != nil {
			_ = Clear(m.config.SessionFile)
			return fmt.Errorf("failed to reconnect via port: %w", err)
		}
		m.pilot = pilot
		m.state = State{
			Status: StatusConnected,
			Pilot:  pilot,
			Info:   info,
		}
		return nil
	}

	return fmt.Errorf("no valid reconnection endpoint in saved session")
}

// persistLocked saves session info to disk (caller must hold mutex).
func (m *Manager) persistLocked() error {
	if m.state.Info == nil {
		return nil
	}
	return Save(m.config.SessionFile, m.state.Info)
}

// State returns the current session state.
func (m *Manager) State() State {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.state
}

// IsConnected returns true if a browser is connected.
func (m *Manager) IsConnected() bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.pilot != nil && !m.pilot.IsClosed()
}

// Detach detaches from the browser without closing it.
// The session info is preserved for later reconnection.
func (m *Manager) Detach() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.pilot == nil {
		return nil
	}

	// Save session info before detaching
	if err := m.persistLocked(); err != nil {
		return fmt.Errorf("failed to persist session before detach: %w", err)
	}

	// Clear local references without closing browser
	m.pilot = nil
	m.state = State{
		Status: StatusDisconnected,
		Info:   m.state.Info, // Keep info for reconnection
	}

	return nil
}

// Close closes the browser and clears session info.
func (m *Manager) Close(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Clear persisted session
	if m.config.SessionFile != "" || m.config.AutoReconnect {
		_ = Clear(m.config.SessionFile)
	}

	if m.pilot == nil {
		return nil
	}

	err := m.pilot.Quit(ctx)
	m.pilot = nil
	m.state = State{
		Status: StatusDisconnected,
	}

	return err
}

// Refresh updates the session info (e.g., after navigation).
func (m *Manager) Refresh(ctx context.Context) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.pilot == nil || m.state.Info == nil {
		return nil
	}

	// Update CDP port if changed
	m.state.Info.CDPPort = m.pilot.CDPPort()

	return m.persistLocked()
}
