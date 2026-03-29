package session

import (
	"path/filepath"
	"testing"
	"time"
)

func TestNewManager(t *testing.T) {
	config := Config{
		Headless: true,
	}

	manager := NewManager(config)

	if manager == nil {
		t.Fatal("NewManager should return a manager")
	}

	// Check default timeout was set
	if manager.config.DefaultTimeout != 30*time.Second {
		t.Errorf("Default timeout should be 30s, got %v", manager.config.DefaultTimeout)
	}

	// Check initial state
	state := manager.State()
	if state.Status != StatusDisconnected {
		t.Errorf("Initial status should be disconnected, got %v", state.Status)
	}

	if manager.IsConnected() {
		t.Error("Manager should not be connected initially")
	}
}

func TestManagerConfig(t *testing.T) {
	tmpDir := t.TempDir()
	sessionFile := filepath.Join(tmpDir, "session.json")

	config := Config{
		Headless:       true,
		DefaultTimeout: 60 * time.Second,
		SessionFile:    sessionFile,
		AutoReconnect:  true,
		InitScripts:    []string{"console.log('test')"},
	}

	manager := NewManager(config)

	if manager.config.Headless != config.Headless {
		t.Error("Headless config mismatch")
	}
	if manager.config.DefaultTimeout != config.DefaultTimeout {
		t.Error("DefaultTimeout config mismatch")
	}
	if manager.config.SessionFile != config.SessionFile {
		t.Error("SessionFile config mismatch")
	}
	if manager.config.AutoReconnect != config.AutoReconnect {
		t.Error("AutoReconnect config mismatch")
	}
	if len(manager.config.InitScripts) != 1 || manager.config.InitScripts[0] != "console.log('test')" {
		t.Error("InitScripts config mismatch")
	}
}

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.Headless != false {
		t.Error("Default Headless should be false")
	}
	if config.DefaultTimeout != 30*time.Second {
		t.Errorf("Default timeout should be 30s, got %v", config.DefaultTimeout)
	}
	if config.AutoReconnect != true {
		t.Error("Default AutoReconnect should be true")
	}
}
