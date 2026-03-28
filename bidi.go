package webpilot

import (
	"context"
	"encoding/json"
	"sync"
)

// BiDiCommand represents a WebDriver BiDi command.
type BiDiCommand struct {
	ID     int64       `json:"id"`
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

// BiDiResponse represents a WebDriver BiDi response.
type BiDiResponse struct {
	ID      int64           `json:"id"`
	Type    string          `json:"type"`
	Method  string          `json:"method,omitempty"` // For events
	Result  json.RawMessage `json:"result,omitempty"`
	Params  json.RawMessage `json:"params,omitempty"` // For events
	Error   string          `json:"error,omitempty"`
	Message string          `json:"message,omitempty"`
}

// BiDiEvent represents a WebDriver BiDi event.
type BiDiEvent struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

// EventHandler is a callback for handling BiDi events.
type EventHandler func(event *BiDiEvent)

// BiDiTransport is the interface for BiDi communication.
// wsTransport implements this interface using WebSocket.
type BiDiTransport interface {
	// Send sends a command and waits for the response.
	Send(ctx context.Context, method string, params interface{}) (json.RawMessage, error)

	// OnEvent registers a handler for events matching the given method pattern.
	OnEvent(method string, handler EventHandler)

	// RemoveEventHandlers removes all handlers for the given method.
	RemoveEventHandlers(method string)

	// Close closes the transport connection.
	Close() error
}

// BiDiClient wraps a BiDiTransport with convenience methods.
// This provides a stable interface for the rest of the codebase.
type BiDiClient struct {
	transport BiDiTransport
	handlers  map[string][]EventHandler // Event method -> handlers
	handlerMu sync.RWMutex
}

// NewBiDiClient creates a new BiDi client wrapping the given transport.
func NewBiDiClient(transport BiDiTransport) *BiDiClient {
	return &BiDiClient{
		transport: transport,
		handlers:  make(map[string][]EventHandler),
	}
}

// OnEvent registers a handler for events matching the given method pattern.
// The method can be an exact match (e.g., "log.entryAdded") or a prefix
// (e.g., "log." to match all log events).
func (c *BiDiClient) OnEvent(method string, handler EventHandler) {
	c.handlerMu.Lock()
	c.handlers[method] = append(c.handlers[method], handler)
	c.handlerMu.Unlock()

	// Also register with the underlying transport
	c.transport.OnEvent(method, handler)
}

// RemoveEventHandlers removes all handlers for the given method.
func (c *BiDiClient) RemoveEventHandlers(method string) {
	c.handlerMu.Lock()
	delete(c.handlers, method)
	c.handlerMu.Unlock()

	// Also remove from the underlying transport
	c.transport.RemoveEventHandlers(method)
}

// Close closes the connection.
func (c *BiDiClient) Close() error {
	return c.transport.Close()
}

// Send sends a command and waits for the response.
func (c *BiDiClient) Send(ctx context.Context, method string, params interface{}) (json.RawMessage, error) {
	return c.transport.Send(ctx, method, params)
}
