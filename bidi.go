package webpilot

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
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

// BiDiClient manages WebSocket communication with Chrome's BiDi server.
type BiDiClient struct {
	conn      *websocket.Conn
	url       string
	nextID    atomic.Int64
	pending   map[int64]chan *BiDiResponse
	pendingMu sync.Mutex
	handlers  map[string][]EventHandler // Event method -> handlers
	handlerMu sync.RWMutex
	closed    atomic.Bool
	closeCh   chan struct{}
}

// NewBiDiClient creates a new BiDi client.
func NewBiDiClient() *BiDiClient {
	return &BiDiClient{
		pending:  make(map[int64]chan *BiDiResponse),
		handlers: make(map[string][]EventHandler),
		closeCh:  make(chan struct{}),
	}
}

// OnEvent registers a handler for events matching the given method pattern.
// The method can be an exact match (e.g., "log.entryAdded") or a prefix
// (e.g., "log." to match all log events).
func (c *BiDiClient) OnEvent(method string, handler EventHandler) {
	c.handlerMu.Lock()
	defer c.handlerMu.Unlock()
	c.handlers[method] = append(c.handlers[method], handler)
}

// RemoveEventHandlers removes all handlers for the given method.
func (c *BiDiClient) RemoveEventHandlers(method string) {
	c.handlerMu.Lock()
	defer c.handlerMu.Unlock()
	delete(c.handlers, method)
}

// dispatchEvent routes an event to all matching handlers.
func (c *BiDiClient) dispatchEvent(event *BiDiEvent) {
	c.handlerMu.RLock()
	defer c.handlerMu.RUnlock()

	// Try exact match first
	if handlers, ok := c.handlers[event.Method]; ok {
		for _, h := range handlers {
			go h(event)
		}
	}

	// Also check for wildcard/prefix handlers (e.g., "*" for all events)
	if handlers, ok := c.handlers["*"]; ok {
		for _, h := range handlers {
			go h(event)
		}
	}
}

// Connect establishes a WebSocket connection to Chrome's BiDi server.
func (c *BiDiClient) Connect(ctx context.Context, url string) error {
	dialer := websocket.Dialer{}
	conn, _, err := dialer.DialContext(ctx, url, nil)
	if err != nil {
		return &ConnectionError{URL: url, Cause: err}
	}

	c.conn = conn
	c.url = url

	// Start message receiver
	go c.receiveLoop()

	return nil
}

// Close closes the WebSocket connection.
func (c *BiDiClient) Close() error {
	if c.closed.Swap(true) {
		return nil // Already closed
	}

	close(c.closeCh)

	// Reject all pending requests
	c.pendingMu.Lock()
	for _, ch := range c.pending {
		close(ch)
	}
	c.pending = make(map[int64]chan *BiDiResponse)
	c.pendingMu.Unlock()

	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// Send sends a command and waits for the response.
func (c *BiDiClient) Send(ctx context.Context, method string, params interface{}) (json.RawMessage, error) {
	if c.closed.Load() {
		return nil, ErrConnectionClosed
	}

	id := c.nextID.Add(1)
	cmd := BiDiCommand{
		ID:     id,
		Method: method,
		Params: params,
	}

	// Create response channel
	respCh := make(chan *BiDiResponse, 1)
	c.pendingMu.Lock()
	c.pending[id] = respCh
	c.pendingMu.Unlock()

	// Clean up on exit
	defer func() {
		c.pendingMu.Lock()
		delete(c.pending, id)
		c.pendingMu.Unlock()
	}()

	// Send command
	if err := c.conn.WriteJSON(cmd); err != nil {
		return nil, fmt.Errorf("failed to send command: %w", err)
	}

	// Wait for response
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-c.closeCh:
		return nil, ErrConnectionClosed
	case resp, ok := <-respCh:
		if !ok {
			return nil, ErrConnectionClosed
		}
		if resp.Type == "error" || resp.Error != "" {
			return nil, &BiDiError{
				ErrorType: resp.Error,
				Message:   resp.Message,
			}
		}
		return resp.Result, nil
	}
}

func (c *BiDiClient) receiveLoop() {
	for {
		select {
		case <-c.closeCh:
			return
		default:
		}

		var resp BiDiResponse
		if err := c.conn.ReadJSON(&resp); err != nil {
			if c.closed.Load() {
				return
			}
			// Connection error - close everything
			_ = c.Close()
			return
		}

		// Check if this is an event (type="event" or has method but no ID)
		if resp.Type == "event" || (resp.Method != "" && resp.ID == 0) {
			event := &BiDiEvent{
				Method: resp.Method,
				Params: resp.Params,
			}
			c.dispatchEvent(event)
			continue
		}

		// Route response to waiting request
		c.pendingMu.Lock()
		ch, ok := c.pending[resp.ID]
		c.pendingMu.Unlock()

		if ok {
			select {
			case ch <- &resp:
			default:
			}
		}
	}
}
