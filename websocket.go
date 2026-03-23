package webpilot

import (
	"context"
	"encoding/json"
	"sync"
)

// WebSocketInfo represents a WebSocket connection made by the page.
type WebSocketInfo struct {
	client   *BiDiClient
	context  string
	socketID string

	URL      string `json:"url"`
	IsClosed bool   `json:"isClosed"`

	messageHandlers []WebSocketMessageHandler
	closeHandlers   []WebSocketCloseHandler
	handlerMu       sync.RWMutex
}

// WebSocketMessage represents a message sent or received on a WebSocket.
type WebSocketMessage struct {
	SocketID  string `json:"socketId"`
	Data      string `json:"data"`
	IsBinary  bool   `json:"isBinary"`
	Direction string `json:"direction"` // "sent" or "received"
}

// WebSocketHandler is called when a new WebSocket connection is opened.
type WebSocketHandler func(*WebSocketInfo)

// WebSocketMessageHandler is called when a WebSocket message is sent or received.
type WebSocketMessageHandler func(*WebSocketMessage)

// WebSocketCloseHandler is called when a WebSocket connection is closed.
type WebSocketCloseHandler func(code int, reason string)

// OnMessage registers a handler for WebSocket messages.
func (ws *WebSocketInfo) OnMessage(handler WebSocketMessageHandler) {
	ws.handlerMu.Lock()
	defer ws.handlerMu.Unlock()
	ws.messageHandlers = append(ws.messageHandlers, handler)
}

// OnClose registers a handler for when the WebSocket closes.
func (ws *WebSocketInfo) OnClose(handler WebSocketCloseHandler) {
	ws.handlerMu.Lock()
	defer ws.handlerMu.Unlock()
	ws.closeHandlers = append(ws.closeHandlers, handler)
}

// dispatchMessage calls all registered message handlers.
func (ws *WebSocketInfo) dispatchMessage(msg *WebSocketMessage) {
	ws.handlerMu.RLock()
	defer ws.handlerMu.RUnlock()
	for _, h := range ws.messageHandlers {
		go h(msg)
	}
}

// dispatchClose calls all registered close handlers.
func (ws *WebSocketInfo) dispatchClose(code int, reason string) {
	ws.handlerMu.RLock()
	defer ws.handlerMu.RUnlock()
	ws.IsClosed = true
	for _, h := range ws.closeHandlers {
		go h(code, reason)
	}
}

// OnWebSocket registers a handler that is called when the page opens a WebSocket connection.
func (p *Pilot) OnWebSocket(ctx context.Context, handler WebSocketHandler) error {
	if p.closed {
		return ErrConnectionClosed
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return err
	}

	// Track active WebSocket connections
	sockets := make(map[string]*WebSocketInfo)
	var socketsMu sync.RWMutex

	// Register handler for WebSocket created events
	p.client.OnEvent("network.webSocketCreated", func(event *BiDiEvent) {
		var params struct {
			SocketID string `json:"socketId"`
			URL      string `json:"url"`
			Context  string `json:"context"`
		}
		if err := json.Unmarshal(event.Params, &params); err != nil {
			debugLog(ctx, "failed to unmarshal websocket created event", "error", err)
			return
		}

		// Only handle WebSockets for our context
		if params.Context != browsingCtx {
			return
		}

		ws := &WebSocketInfo{
			client:   p.client,
			context:  browsingCtx,
			socketID: params.SocketID,
			URL:      params.URL,
			IsClosed: false,
		}

		socketsMu.Lock()
		sockets[params.SocketID] = ws
		socketsMu.Unlock()

		handler(ws)
	})

	// Register handler for WebSocket message events
	p.client.OnEvent("network.webSocketFrameSent", func(event *BiDiEvent) {
		var params struct {
			SocketID string `json:"socketId"`
			Data     string `json:"data"`
			IsBinary bool   `json:"opcode"` // opcode 2 = binary
		}
		if err := json.Unmarshal(event.Params, &params); err != nil {
			return
		}

		socketsMu.RLock()
		ws, ok := sockets[params.SocketID]
		socketsMu.RUnlock()

		if ok {
			ws.dispatchMessage(&WebSocketMessage{
				SocketID:  params.SocketID,
				Data:      params.Data,
				IsBinary:  params.IsBinary,
				Direction: "sent",
			})
		}
	})

	p.client.OnEvent("network.webSocketFrameReceived", func(event *BiDiEvent) {
		var params struct {
			SocketID string `json:"socketId"`
			Data     string `json:"data"`
			IsBinary bool   `json:"opcode"`
		}
		if err := json.Unmarshal(event.Params, &params); err != nil {
			return
		}

		socketsMu.RLock()
		ws, ok := sockets[params.SocketID]
		socketsMu.RUnlock()

		if ok {
			ws.dispatchMessage(&WebSocketMessage{
				SocketID:  params.SocketID,
				Data:      params.Data,
				IsBinary:  params.IsBinary,
				Direction: "received",
			})
		}
	})

	// Register handler for WebSocket closed events
	p.client.OnEvent("network.webSocketClosed", func(event *BiDiEvent) {
		var params struct {
			SocketID string `json:"socketId"`
			Code     int    `json:"code"`
			Reason   string `json:"reason"`
		}
		if err := json.Unmarshal(event.Params, &params); err != nil {
			return
		}

		socketsMu.Lock()
		ws, ok := sockets[params.SocketID]
		if ok {
			delete(sockets, params.SocketID)
		}
		socketsMu.Unlock()

		if ok {
			ws.dispatchClose(params.Code, params.Reason)
		}
	})

	// Subscribe to WebSocket network events
	_, err = p.client.Send(ctx, "session.subscribe", map[string]interface{}{
		"events": []string{
			"network.webSocketCreated",
			"network.webSocketFrameSent",
			"network.webSocketFrameReceived",
			"network.webSocketClosed",
		},
		"contexts": []string{browsingCtx},
	})
	return err
}
