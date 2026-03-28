package cdp

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"sync/atomic"

	"github.com/gorilla/websocket"
)

// Client is a Chrome DevTools Protocol client.
type Client struct {
	conn      *websocket.Conn
	url       string
	nextID    atomic.Int64
	pending   map[int64]chan *Message
	pendingMu sync.RWMutex
	handlers  map[string][]EventHandler
	handlerMu sync.RWMutex
	closed    bool
	closedMu  sync.RWMutex
	closeCh   chan struct{}
}

// NewClient creates a new CDP client.
func NewClient() *Client {
	return &Client{
		pending:  make(map[int64]chan *Message),
		handlers: make(map[string][]EventHandler),
		closeCh:  make(chan struct{}),
	}
}

// Connect establishes a WebSocket connection to the CDP endpoint.
func (c *Client) Connect(ctx context.Context, url string) error {
	dialer := websocket.DefaultDialer

	conn, _, err := dialer.DialContext(ctx, url, nil)
	if err != nil {
		return fmt.Errorf("cdp: failed to connect: %w", err)
	}

	c.conn = conn
	c.url = url

	// Start reading messages
	go c.readLoop()

	return nil
}

// URL returns the WebSocket URL this client is connected to.
func (c *Client) URL() string {
	return c.url
}

// readLoop continuously reads messages from the WebSocket.
func (c *Client) readLoop() {
	for {
		select {
		case <-c.closeCh:
			return
		default:
		}

		_, data, err := c.conn.ReadMessage()
		if err != nil {
			c.closedMu.Lock()
			closed := c.closed
			c.closedMu.Unlock()
			if closed {
				return
			}
			// Connection error - close and exit
			_ = c.Close()
			return
		}

		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			continue
		}

		// Check if it's an event (has method, no id)
		if msg.Method != "" && msg.ID == 0 {
			c.dispatchEvent(msg.Method, msg.Params)
			continue
		}

		// It's a response to a command
		c.pendingMu.RLock()
		ch, ok := c.pending[msg.ID]
		c.pendingMu.RUnlock()

		if ok {
			ch <- &msg
		}
	}
}

// dispatchEvent sends an event to all registered handlers.
func (c *Client) dispatchEvent(method string, params json.RawMessage) {
	c.handlerMu.RLock()
	defer c.handlerMu.RUnlock()

	if handlers, ok := c.handlers[method]; ok {
		for _, h := range handlers {
			go h(params)
		}
	}
}

// Send sends a CDP command and waits for the response.
func (c *Client) Send(ctx context.Context, method string, params interface{}) (json.RawMessage, error) {
	c.closedMu.RLock()
	if c.closed {
		c.closedMu.RUnlock()
		return nil, fmt.Errorf("cdp: connection closed")
	}
	c.closedMu.RUnlock()

	id := c.nextID.Add(1)

	// Build message
	msg := Message{
		ID:     id,
		Method: method,
	}

	if params != nil {
		data, err := json.Marshal(params)
		if err != nil {
			return nil, fmt.Errorf("cdp: failed to marshal params: %w", err)
		}
		msg.Params = data
	}

	// Create response channel
	respCh := make(chan *Message, 1)
	c.pendingMu.Lock()
	c.pending[id] = respCh
	c.pendingMu.Unlock()

	defer func() {
		c.pendingMu.Lock()
		delete(c.pending, id)
		c.pendingMu.Unlock()
	}()

	// Send the message
	data, err := json.Marshal(msg)
	if err != nil {
		return nil, fmt.Errorf("cdp: failed to marshal message: %w", err)
	}

	if err := c.conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return nil, fmt.Errorf("cdp: failed to send message: %w", err)
	}

	// Wait for response
	select {
	case resp := <-respCh:
		if resp.Error != nil {
			return nil, resp.Error
		}
		return resp.Result, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-c.closeCh:
		return nil, fmt.Errorf("cdp: connection closed")
	}
}

// OnEvent registers a handler for CDP events.
func (c *Client) OnEvent(method string, handler EventHandler) {
	c.handlerMu.Lock()
	c.handlers[method] = append(c.handlers[method], handler)
	c.handlerMu.Unlock()
}

// RemoveEventHandlers removes all handlers for the given method.
func (c *Client) RemoveEventHandlers(method string) {
	c.handlerMu.Lock()
	delete(c.handlers, method)
	c.handlerMu.Unlock()
}

// Close closes the CDP connection.
func (c *Client) Close() error {
	c.closedMu.Lock()
	if c.closed {
		c.closedMu.Unlock()
		return nil
	}
	c.closed = true
	c.closedMu.Unlock()

	close(c.closeCh)

	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// IsConnected returns true if the client is connected.
func (c *Client) IsConnected() bool {
	c.closedMu.RLock()
	defer c.closedMu.RUnlock()
	return !c.closed && c.conn != nil
}
