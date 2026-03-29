package session

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	// DefaultSessionDir is the default directory for session files.
	DefaultSessionDir = ".w3pilot"

	// DefaultSessionFile is the default session file name.
	DefaultSessionFile = "session.json"
)

// DefaultSessionPath returns the default session file path.
func DefaultSessionPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return filepath.Join(".", DefaultSessionDir, DefaultSessionFile)
	}
	return filepath.Join(home, DefaultSessionDir, DefaultSessionFile)
}

// Save persists session info to disk.
func Save(path string, info *Info) error {
	if path == "" {
		path = DefaultSessionPath()
	}

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create session directory: %w", err)
	}

	data, err := json.MarshalIndent(info, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write session file: %w", err)
	}

	return nil
}

// Load reads session info from disk.
func Load(path string) (*Info, error) {
	if path == "" {
		path = DefaultSessionPath()
	}

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil // No session file is not an error
		}
		return nil, fmt.Errorf("failed to read session file: %w", err)
	}

	var info Info
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("failed to parse session file: %w", err)
	}

	return &info, nil
}

// Clear removes the session file.
func Clear(path string) error {
	if path == "" {
		path = DefaultSessionPath()
	}

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove session file: %w", err)
	}
	return nil
}

// Exists checks if a session file exists.
func Exists(path string) bool {
	if path == "" {
		path = DefaultSessionPath()
	}

	_, err := os.Stat(path)
	return err == nil
}
