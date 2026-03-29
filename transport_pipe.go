package w3pilot

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sync"
	"sync/atomic"
	"time"
)

// pipeTransport implements BiDiTransport using stdin/stdout pipes to clicker.
type pipeTransport struct {
	cmd             *exec.Cmd
	stdin           io.WriteCloser
	stdout          *bufio.Reader
	stderr          io.ReadCloser
	nextID          atomic.Int64
	pending         map[int64]chan *BiDiResponse
	pendingMu       sync.RWMutex
	handlers        map[string][]EventHandler
	handlerMu       sync.RWMutex
	closed          bool
	closedMu        sync.RWMutex
	closeCh         chan struct{}
	writeMu         sync.Mutex // Serialize writes to stdin
	browsingContext string     // Captured from contextCreated event
}

// PipeOptions configures the pipe transport.
type PipeOptions struct {
	// Headless runs the browser in headless mode.
	Headless bool

	// ExecutablePath is the path to the clicker binary.
	// If empty, it will be discovered automatically.
	ExecutablePath string

	// StartupTimeout is the maximum time to wait for clicker to be ready.
	// Default: 30 seconds.
	StartupTimeout time.Duration
}

// newPipeTransport creates a new pipe transport.
func newPipeTransport() *pipeTransport {
	return &pipeTransport{
		pending:  make(map[int64]chan *BiDiResponse),
		handlers: make(map[string][]EventHandler),
		closeCh:  make(chan struct{}),
	}
}

// Start spawns the clicker process and establishes pipe communication.
func (t *pipeTransport) Start(ctx context.Context, opts *PipeOptions) error {
	if opts == nil {
		opts = &PipeOptions{}
	}

	// Find clicker binary
	clickerPath := opts.ExecutablePath
	if clickerPath == "" {
		var err error
		clickerPath, err = FindClickerBinary()
		if err != nil {
			return fmt.Errorf("clicker binary not found: %w", err)
		}
	}

	// Build command arguments
	args := []string{"pipe"}
	if opts.Headless {
		args = append(args, "--headless")
	}

	// Create command WITHOUT CommandContext to prevent process termination
	// when the request context is cancelled. The clicker process should live
	// beyond the initial launch request. We manage its lifecycle via Close().
	t.cmd = exec.Command(clickerPath, args...)

	// Set up pipes
	var err error
	t.stdin, err = t.cmd.StdinPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdin pipe: %w", err)
	}

	stdout, err := t.cmd.StdoutPipe()
	if err != nil {
		return fmt.Errorf("failed to create stdout pipe: %w", err)
	}
	t.stdout = bufio.NewReader(stdout)

	t.stderr, err = t.cmd.StderrPipe()
	if err != nil {
		return fmt.Errorf("failed to create stderr pipe: %w", err)
	}

	// Start the process
	if err := t.cmd.Start(); err != nil {
		return fmt.Errorf("failed to start clicker: %w", err)
	}

	// Start reading stderr (for debugging)
	go t.readStderr()

	// Start reading responses
	go t.readLoop()

	// Wait for ready signal
	timeout := opts.StartupTimeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}
	if err := t.waitForReady(ctx, timeout); err != nil {
		_ = t.Close()
		return err
	}

	return nil
}

// waitForReady waits for the vibium:lifecycle.ready event.
// It also captures the browsingContext from the contextCreated event.
func (t *pipeTransport) waitForReady(ctx context.Context, timeout time.Duration) error {
	readyCh := make(chan struct{}, 1)

	// Capture browsingContext from contextCreated event (sent before lifecycle.ready)
	t.OnEvent("browsingContext.contextCreated", func(event *BiDiEvent) {
		var params struct {
			Context string `json:"context"`
		}
		if err := json.Unmarshal(event.Params, &params); err == nil && params.Context != "" {
			t.browsingContext = params.Context
		}
	})

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
		return fmt.Errorf("timeout waiting for clicker ready")
	case <-ctx.Done():
		return ctx.Err()
	}
}

// BrowsingContext returns the captured browsing context ID.
func (t *pipeTransport) BrowsingContext() string {
	return t.browsingContext
}

// readStderr reads stderr and logs it when W3PILOT_DEBUG is enabled.
// Previously this was discarded, making clicker errors invisible.
func (t *pipeTransport) readStderr() {
	debug := Debug()
	scanner := bufio.NewScanner(t.stderr)
	for scanner.Scan() {
		if debug {
			// Log clicker stderr to our stderr for debugging
			fmt.Fprintf(os.Stderr, "[clicker] %s\n", scanner.Text())
		}
	}
}

// readLoop continuously reads messages from stdout.
func (t *pipeTransport) readLoop() {
	for {
		select {
		case <-t.closeCh:
			return
		default:
		}

		line, err := t.stdout.ReadBytes('\n')
		if err != nil {
			t.closedMu.Lock()
			closed := t.closed
			t.closedMu.Unlock()
			if closed {
				return
			}
			_ = t.Close()
			return
		}

		// Parse the message
		var resp BiDiResponse
		if err := json.Unmarshal(line, &resp); err != nil {
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
func (t *pipeTransport) dispatchEvent(event *BiDiEvent) {
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
func (t *pipeTransport) Send(ctx context.Context, method string, params interface{}) (json.RawMessage, error) {
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

	// Write with newline delimiter
	t.writeMu.Lock()
	_, err = t.stdin.Write(append(data, '\n'))
	t.writeMu.Unlock()
	if err != nil {
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
func (t *pipeTransport) OnEvent(method string, handler EventHandler) {
	t.handlerMu.Lock()
	t.handlers[method] = append(t.handlers[method], handler)
	t.handlerMu.Unlock()
}

// RemoveEventHandlers removes all handlers for the given method.
func (t *pipeTransport) RemoveEventHandlers(method string) {
	t.handlerMu.Lock()
	delete(t.handlers, method)
	t.handlerMu.Unlock()
}

// Close closes the pipe transport and terminates the clicker process.
func (t *pipeTransport) Close() error {
	t.closedMu.Lock()
	if t.closed {
		t.closedMu.Unlock()
		return nil
	}
	t.closed = true
	t.closedMu.Unlock()

	close(t.closeCh)

	if t.stdin != nil {
		_ = t.stdin.Close()
	}

	if t.cmd != nil && t.cmd.Process != nil {
		_ = t.cmd.Process.Kill()
		_ = t.cmd.Wait()
	}

	return nil
}
