package session

import (
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestSaveAndLoad(t *testing.T) {
	// Create a temp directory for the test
	tmpDir := t.TempDir()
	sessionPath := filepath.Join(tmpDir, "session.json")

	// Create session info
	info := &Info{
		WebSocketURL: "ws://localhost:9222",
		ClickerPort:  9222,
		CDPPort:      9223,
		Headless:     true,
		ClickerPID:   12345,
		LaunchedAt:   time.Now().UTC().Truncate(time.Second),
		InitScripts:  []string{"console.log('init')"},
	}

	// Save session
	err := Save(sessionPath, info)
	if err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file exists
	if !Exists(sessionPath) {
		t.Fatal("Session file should exist after save")
	}

	// Load session
	loaded, err := Load(sessionPath)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}

	// Verify loaded data
	if loaded.WebSocketURL != info.WebSocketURL {
		t.Errorf("WebSocketURL mismatch: got %s, want %s", loaded.WebSocketURL, info.WebSocketURL)
	}
	if loaded.ClickerPort != info.ClickerPort {
		t.Errorf("ClickerPort mismatch: got %d, want %d", loaded.ClickerPort, info.ClickerPort)
	}
	if loaded.CDPPort != info.CDPPort {
		t.Errorf("CDPPort mismatch: got %d, want %d", loaded.CDPPort, info.CDPPort)
	}
	if loaded.Headless != info.Headless {
		t.Errorf("Headless mismatch: got %v, want %v", loaded.Headless, info.Headless)
	}
	if loaded.ClickerPID != info.ClickerPID {
		t.Errorf("ClickerPID mismatch: got %d, want %d", loaded.ClickerPID, info.ClickerPID)
	}
	if !loaded.LaunchedAt.Equal(info.LaunchedAt) {
		t.Errorf("LaunchedAt mismatch: got %v, want %v", loaded.LaunchedAt, info.LaunchedAt)
	}
	if len(loaded.InitScripts) != len(info.InitScripts) || loaded.InitScripts[0] != info.InitScripts[0] {
		t.Errorf("InitScripts mismatch: got %v, want %v", loaded.InitScripts, info.InitScripts)
	}
}

func TestLoadNonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	sessionPath := filepath.Join(tmpDir, "nonexistent.json")

	// Load should return nil, nil for non-existent file
	loaded, err := Load(sessionPath)
	if err != nil {
		t.Fatalf("Load should not return error for non-existent file: %v", err)
	}
	if loaded != nil {
		t.Error("Load should return nil for non-existent file")
	}
}

func TestClear(t *testing.T) {
	tmpDir := t.TempDir()
	sessionPath := filepath.Join(tmpDir, "session.json")

	// Create a file
	info := &Info{WebSocketURL: "ws://localhost:9222"}
	if err := Save(sessionPath, info); err != nil {
		t.Fatalf("Save failed: %v", err)
	}

	// Verify file exists
	if !Exists(sessionPath) {
		t.Fatal("Session file should exist after save")
	}

	// Clear the file
	if err := Clear(sessionPath); err != nil {
		t.Fatalf("Clear failed: %v", err)
	}

	// Verify file is gone
	if Exists(sessionPath) {
		t.Error("Session file should not exist after clear")
	}
}

func TestClearNonExistent(t *testing.T) {
	tmpDir := t.TempDir()
	sessionPath := filepath.Join(tmpDir, "nonexistent.json")

	// Clear should not return error for non-existent file
	if err := Clear(sessionPath); err != nil {
		t.Errorf("Clear should not return error for non-existent file: %v", err)
	}
}

func TestDefaultSessionPath(t *testing.T) {
	path := DefaultSessionPath()
	if path == "" {
		t.Error("DefaultSessionPath should not return empty string")
	}
	if !filepath.IsAbs(path) && !filepath.HasPrefix(path, ".") {
		t.Errorf("DefaultSessionPath should return absolute or relative path: %s", path)
	}
	if filepath.Base(path) != DefaultSessionFile {
		t.Errorf("DefaultSessionPath should use default file name: got %s, want %s", filepath.Base(path), DefaultSessionFile)
	}
}

func TestSaveCreatesDirectory(t *testing.T) {
	tmpDir := t.TempDir()
	sessionPath := filepath.Join(tmpDir, "nested", "dir", "session.json")

	info := &Info{WebSocketURL: "ws://localhost:9222"}
	if err := Save(sessionPath, info); err != nil {
		t.Fatalf("Save should create nested directories: %v", err)
	}

	// Verify file was created
	if _, err := os.Stat(sessionPath); err != nil {
		t.Errorf("Session file should exist: %v", err)
	}
}
