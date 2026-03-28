package w3pilot

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gorilla/websocket"
)

// wsTransport implements BiDiTransport using WebSocket.
type wsTransport struct {
	conn      *websocket.Conn
	nextID    atomic.Int64
	pending   map[int64]chan *BiDiResponse
	pendingMu sync.RWMutex
	handlers  map[string][]EventHandler
	handlerMu sync.RWMutex
	closed    bool
	closedMu  sync.RWMutex
	closeCh   chan struct{}
}

// newWSTransport creates a new WebSocket transport.
func newWSTransport() *wsTransport {
	return &wsTransport{
		pending:  make(map[int64]chan *BiDiResponse),
		handlers: make(map[string][]EventHandler),
		closeCh:  make(chan struct{}),
	}
}

// Connect establishes a WebSocket connection to the given URL.
func (t *wsTransport) Connect(ctx context.Context, url string) error {
	dialer := websocket.DefaultDialer

	conn, _, err := dialer.DialContext(ctx, url, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	t.conn = conn

	// Start reading messages
	go t.readLoop()

	return nil
}

// WaitForReady waits for the browser to be ready after connecting.
// In serve mode, clicker sends browsingContext.contextCreated when the browser is ready.
func (t *wsTransport) WaitForReady(ctx context.Context, timeout time.Duration) error {
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	readyCh := make(chan struct{}, 1)

	// In serve mode, wait for browsingContext.contextCreated event
	t.OnEvent("browsingContext.contextCreated", func(event *BiDiEvent) {
		select {
		case readyCh <- struct{}{}:
		default:
		}
	})

	// Also accept vibium:lifecycle.ready (future compatibility)
	t.OnEvent("vibium:lifecycle.ready", func(event *BiDiEvent) {
		select {
		case readyCh <- struct{}{}:
		default:
		}
	})

	select {
	case <-readyCh:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("timeout waiting for browser ready")
	case <-ctx.Done():
		return ctx.Err()
	}
}

// readLoop continuously reads messages from the WebSocket.
func (t *wsTransport) readLoop() {
	for {
		select {
		case <-t.closeCh:
			return
		default:
		}

		_, message, err := t.conn.ReadMessage()
		if err != nil {
			t.closedMu.Lock()
			closed := t.closed
			t.closedMu.Unlock()
			if closed {
				return
			}
			// Connection error - close and exit
			_ = t.Close()
			return
		}

		// Parse the message
		var resp BiDiResponse
		if err := json.Unmarshal(message, &resp); err != nil {
			continue
		}

		// Check if it's an event (no ID, has method)
		if resp.ID == 0 && resp.Method != "" {
			event := &BiDiEvent{
				Method: resp.Method,
				Params: resp.Params,
			}
			t.dispatchEvent(event)
			continue
		}

		// It's a response to a command
		t.pendingMu.RLock()
		ch, ok := t.pending[resp.ID]
		t.pendingMu.RUnlock()

		if ok {
			ch <- &resp
		}
	}
}

// dispatchEvent sends an event to all registered handlers.
func (t *wsTransport) dispatchEvent(event *BiDiEvent) {
	t.handlerMu.RLock()
	defer t.handlerMu.RUnlock()

	// Exact match handlers
	if handlers, ok := t.handlers[event.Method]; ok {
		for _, h := range handlers {
			go h(event)
		}
	}

	// Prefix match handlers (e.g., "log." matches "log.entryAdded")
	for pattern, handlers := range t.handlers {
		if len(pattern) > 0 && pattern[len(pattern)-1] == '.' {
			if len(event.Method) > len(pattern) && event.Method[:len(pattern)] == pattern {
				for _, h := range handlers {
					go h(event)
				}
			}
		}
	}
}

// Send sends a command and waits for the response.
func (t *wsTransport) Send(ctx context.Context, method string, params interface{}) (json.RawMessage, error) {
	t.closedMu.RLock()
	if t.closed {
		t.closedMu.RUnlock()
		return nil, ErrConnectionClosed
	}
	t.closedMu.RUnlock()

	id := t.nextID.Add(1)

	cmd := BiDiCommand{
		ID:     id,
		Method: method,
		Params: params,
	}

	// Create response channel
	respCh := make(chan *BiDiResponse, 1)
	t.pendingMu.Lock()
	t.pending[id] = respCh
	t.pendingMu.Unlock()

	defer func() {
		t.pendingMu.Lock()
		delete(t.pending, id)
		t.pendingMu.Unlock()
	}()

	// Send the command
	data, err := json.Marshal(cmd)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal command: %w", err)
	}

	if err := t.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return nil, fmt.Errorf("failed to send command: %w", err)
	}

	// Wait for response
	select {
	case resp := <-respCh:
		if resp.Type == "error" || resp.Error != "" {
			return nil, &BiDiError{
				ErrorType: resp.Error,
				Message:   resp.Message,
			}
		}
		return resp.Result, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-t.closeCh:
		return nil, ErrConnectionClosed
	}
}

// OnEvent registers a handler for events matching the given method pattern.
func (t *wsTransport) OnEvent(method string, handler EventHandler) {
	t.handlerMu.Lock()
	t.handlers[method] = append(t.handlers[method], handler)
	t.handlerMu.Unlock()
}

// RemoveEventHandlers removes all handlers for the given method.
func (t *wsTransport) RemoveEventHandlers(method string) {
	t.handlerMu.Lock()
	delete(t.handlers, method)
	t.handlerMu.Unlock()
}

// Close closes the WebSocket connection.
func (t *wsTransport) Close() error {
	t.closedMu.Lock()
	if t.closed {
		t.closedMu.Unlock()
		return nil
	}
	t.closed = true
	t.closedMu.Unlock()

	close(t.closeCh)

	if t.conn != nil {
		return t.conn.Close()
	}
	return nil
}
