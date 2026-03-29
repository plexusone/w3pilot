package cdp

import (
	"encoding/json"
	"sync"
	"testing"
	"time"
)

func TestNewClient(t *testing.T) {
	c := NewClient()

	if c == nil {
		t.Fatal("NewClient() returned nil")
	}
	if c.pending == nil {
		t.Error("NewClient().pending should be initialized")
	}
	if c.handlers == nil {
		t.Error("NewClient().handlers should be initialized")
	}
	if c.closeCh == nil {
		t.Error("NewClient().closeCh should be initialized")
	}
	if c.closed {
		t.Error("NewClient() should not be closed")
	}
}

func TestClientOnEvent(t *testing.T) {
	c := NewClient()

	callCount := 0
	c.OnEvent("Test.event", func(params json.RawMessage) {
		callCount++
	})

	c.handlerMu.RLock()
	handlers := c.handlers["Test.event"]
	c.handlerMu.RUnlock()

	if len(handlers) != 1 {
		t.Errorf("Expected 1 handler, got %d", len(handlers))
	}

	// Add another handler for the same event
	c.OnEvent("Test.event", func(params json.RawMessage) {
		callCount++
	})

	c.handlerMu.RLock()
	handlers = c.handlers["Test.event"]
	c.handlerMu.RUnlock()

	if len(handlers) != 2 {
		t.Errorf("Expected 2 handlers, got %d", len(handlers))
	}
}

func TestClientRemoveEventHandlers(t *testing.T) {
	c := NewClient()

	// Add handlers for two events
	c.OnEvent("Test.event1", func(params json.RawMessage) {})
	c.OnEvent("Test.event1", func(params json.RawMessage) {})
	c.OnEvent("Test.event2", func(params json.RawMessage) {})

	// Remove handlers for event1
	c.RemoveEventHandlers("Test.event1")

	c.handlerMu.RLock()
	handlers1 := c.handlers["Test.event1"]
	handlers2 := c.handlers["Test.event2"]
	c.handlerMu.RUnlock()

	if len(handlers1) != 0 {
		t.Errorf("Event1 handlers should be removed, got %d", len(handlers1))
	}
	if len(handlers2) != 1 {
		t.Errorf("Event2 handler should remain, got %d", len(handlers2))
	}
}

func TestClientDispatchEvent(t *testing.T) {
	c := NewClient()

	var mu sync.Mutex
	callCount := 0
	receivedParams := make([]json.RawMessage, 0)

	c.OnEvent("Test.event", func(params json.RawMessage) {
		mu.Lock()
		callCount++
		receivedParams = append(receivedParams, params)
		mu.Unlock()
	})

	c.OnEvent("Test.event", func(params json.RawMessage) {
		mu.Lock()
		callCount++
		mu.Unlock()
	})

	testParams := json.RawMessage(`{"key": "value"}`)
	c.dispatchEvent("Test.event", testParams)

	// Allow time for goroutines to execute
	time.Sleep(10 * time.Millisecond)

	mu.Lock()
	if callCount != 2 {
		t.Errorf("Expected 2 calls, got %d", callCount)
	}
	mu.Unlock()
}

func TestClientDispatchEvent_NoHandlers(t *testing.T) {
	c := NewClient()

	// Should not panic when no handlers registered
	c.dispatchEvent("Nonexistent.event", json.RawMessage(`{}`))
}

func TestClientClose(t *testing.T) {
	c := NewClient()

	// First close should succeed
	err := c.Close()
	if err != nil {
		t.Errorf("Close() error = %v", err)
	}

	c.closedMu.RLock()
	if !c.closed {
		t.Error("Close() should set closed to true")
	}
	c.closedMu.RUnlock()

	// Second close should be idempotent
	err = c.Close()
	if err != nil {
		t.Errorf("Second Close() error = %v", err)
	}
}

func TestClientIsConnected(t *testing.T) {
	c := NewClient()

	// No connection, not closed
	if c.IsConnected() {
		t.Error("IsConnected() should return false when conn is nil")
	}

	// Close the client
	c.Close()
	if c.IsConnected() {
		t.Error("IsConnected() should return false when closed")
	}
}

func TestClientURL(t *testing.T) {
	c := NewClient()
	c.url = "ws://localhost:9222/devtools/page/ABC123"

	if c.URL() != "ws://localhost:9222/devtools/page/ABC123" {
		t.Errorf("URL() = %q, want %q", c.URL(), "ws://localhost:9222/devtools/page/ABC123")
	}
}

func TestMessageError(t *testing.T) {
	tests := []struct {
		name string
		err  Error
		want string
	}{
		{
			name: "message only",
			err:  Error{Code: -32600, Message: "Invalid Request"},
			want: "Invalid Request",
		},
		{
			name: "message with data",
			err:  Error{Code: -32601, Message: "Method not found", Data: "unknown method: Test.foo"},
			want: "Method not found: unknown method: Test.foo",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.err.Error()
			if got != tt.want {
				t.Errorf("Error.Error() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestClientPendingMap(t *testing.T) {
	c := NewClient()

	// Verify pending map is initialized
	if c.pending == nil {
		t.Fatal("pending map should be initialized")
	}

	// Add a pending request
	ch := make(chan *Message, 1)
	c.pendingMu.Lock()
	c.pending[1] = ch
	c.pendingMu.Unlock()

	// Verify it was added
	c.pendingMu.RLock()
	if _, ok := c.pending[1]; !ok {
		t.Error("pending request should exist")
	}
	c.pendingMu.RUnlock()

	// Remove it
	c.pendingMu.Lock()
	delete(c.pending, 1)
	c.pendingMu.Unlock()

	// Verify it was removed
	c.pendingMu.RLock()
	if _, ok := c.pending[1]; ok {
		t.Error("pending request should be removed")
	}
	c.pendingMu.RUnlock()
}

func TestClientNextID(t *testing.T) {
	c := NewClient()

	id1 := c.nextID.Add(1)
	id2 := c.nextID.Add(1)
	id3 := c.nextID.Add(1)

	if id1 != 1 {
		t.Errorf("first ID = %d, want 1", id1)
	}
	if id2 != 2 {
		t.Errorf("second ID = %d, want 2", id2)
	}
	if id3 != 3 {
		t.Errorf("third ID = %d, want 3", id3)
	}
}
