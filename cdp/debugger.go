package cdp

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

// Console debugger event names.
const (
	RuntimeConsoleAPICalled = "Runtime.consoleAPICalled"
	RuntimeExceptionThrown  = "Runtime.exceptionThrown"
	LogEntryAdded           = "Log.entryAdded"
)

// ConsoleMessageType represents the type of console message.
type ConsoleMessageType string

const (
	ConsoleLog     ConsoleMessageType = "log"
	ConsoleDebug   ConsoleMessageType = "debug"
	ConsoleInfo    ConsoleMessageType = "info"
	ConsoleError   ConsoleMessageType = "error"
	ConsoleWarning ConsoleMessageType = "warning"
	ConsoleDir     ConsoleMessageType = "dir"
	ConsoleDirXML  ConsoleMessageType = "dirxml"
	ConsoleTable   ConsoleMessageType = "table"
	ConsoleTrace   ConsoleMessageType = "trace"
	ConsoleClear   ConsoleMessageType = "clear"
	ConsoleAssert  ConsoleMessageType = "assert"
)

// CallFrame represents a single frame in a call stack.
type CallFrame struct {
	FunctionName string `json:"functionName"`
	ScriptID     string `json:"scriptId"`
	URL          string `json:"url"`
	LineNumber   int    `json:"lineNumber"`
	ColumnNumber int    `json:"columnNumber"`
}

// StackTrace represents a JavaScript stack trace.
type StackTrace struct {
	Description string      `json:"description,omitempty"`
	CallFrames  []CallFrame `json:"callFrames"`
	Parent      *StackTrace `json:"parent,omitempty"`
}

// RemoteObject represents a JavaScript object.
type RemoteObject struct {
	Type        string `json:"type"`
	Subtype     string `json:"subtype,omitempty"`
	ClassName   string `json:"className,omitempty"`
	Value       any    `json:"value,omitempty"`
	Description string `json:"description,omitempty"`
	ObjectID    string `json:"objectId,omitempty"`
}

// ConsoleEntry represents a console API call with full stack trace.
type ConsoleEntry struct {
	Type       ConsoleMessageType `json:"type"`
	Args       []RemoteObject     `json:"args"`
	StackTrace *StackTrace        `json:"stackTrace,omitempty"`
	Timestamp  float64            `json:"timestamp"`
	// Formatted message extracted from args
	Text string `json:"text"`
}

// ExceptionDetails contains information about an exception.
type ExceptionDetails struct {
	ExceptionID        int           `json:"exceptionId"`
	Text               string        `json:"text"`
	LineNumber         int           `json:"lineNumber"`
	ColumnNumber       int           `json:"columnNumber"`
	ScriptID           string        `json:"scriptId,omitempty"`
	URL                string        `json:"url,omitempty"`
	StackTrace         *StackTrace   `json:"stackTrace,omitempty"`
	Exception          *RemoteObject `json:"exception,omitempty"`
	ExecutionContextID int           `json:"executionContextId,omitempty"`
}

// LogEntry represents a browser log entry (deprecations, interventions, etc).
type LogEntry struct {
	Source     string      `json:"source"` // network, violation, intervention, etc.
	Level      string      `json:"level"`  // verbose, info, warning, error
	Text       string      `json:"text"`
	Timestamp  float64     `json:"timestamp"`
	URL        string      `json:"url,omitempty"`
	LineNumber int         `json:"lineNumber,omitempty"`
	StackTrace *StackTrace `json:"stackTrace,omitempty"`
}

// ConsoleDebugger captures console messages with stack traces.
type ConsoleDebugger struct {
	client   *Client
	mu       sync.RWMutex
	enabled  bool
	entries  []ConsoleEntry
	errors   []ExceptionDetails
	logs     []LogEntry
	handlers struct {
		console   func(*ConsoleEntry)
		exception func(*ExceptionDetails)
		log       func(*LogEntry)
	}
}

// NewConsoleDebugger creates a new console debugger.
func NewConsoleDebugger(client *Client) *ConsoleDebugger {
	return &ConsoleDebugger{
		client:  client,
		entries: make([]ConsoleEntry, 0),
		errors:  make([]ExceptionDetails, 0),
		logs:    make([]LogEntry, 0),
	}
}

// Enable starts capturing console messages with stack traces.
func (d *ConsoleDebugger) Enable(ctx context.Context) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if d.enabled {
		return nil
	}

	// Enable Runtime domain
	if _, err := d.client.Send(ctx, RuntimeEnable, nil); err != nil {
		return fmt.Errorf("cdp: failed to enable Runtime: %w", err)
	}

	// Enable Log domain for deprecations/interventions
	if _, err := d.client.Send(ctx, "Log.enable", nil); err != nil {
		// Log domain may not be available, continue without it
	}

	// Register event handlers
	d.client.OnEvent(RuntimeConsoleAPICalled, func(params json.RawMessage) {
		var event struct {
			Type       string         `json:"type"`
			Args       []RemoteObject `json:"args"`
			StackTrace *StackTrace    `json:"stackTrace"`
			Timestamp  float64        `json:"timestamp"`
		}
		if err := json.Unmarshal(params, &event); err != nil {
			return
		}

		entry := ConsoleEntry{
			Type:       ConsoleMessageType(event.Type),
			Args:       event.Args,
			StackTrace: event.StackTrace,
			Timestamp:  event.Timestamp,
			Text:       formatConsoleArgs(event.Args),
		}

		d.mu.Lock()
		d.entries = append(d.entries, entry)
		handler := d.handlers.console
		d.mu.Unlock()

		if handler != nil {
			handler(&entry)
		}
	})

	d.client.OnEvent(RuntimeExceptionThrown, func(params json.RawMessage) {
		var event struct {
			Timestamp        float64          `json:"timestamp"`
			ExceptionDetails ExceptionDetails `json:"exceptionDetails"`
		}
		if err := json.Unmarshal(params, &event); err != nil {
			return
		}

		d.mu.Lock()
		d.errors = append(d.errors, event.ExceptionDetails)
		handler := d.handlers.exception
		d.mu.Unlock()

		if handler != nil {
			handler(&event.ExceptionDetails)
		}
	})

	d.client.OnEvent(LogEntryAdded, func(params json.RawMessage) {
		var event struct {
			Entry LogEntry `json:"entry"`
		}
		if err := json.Unmarshal(params, &event); err != nil {
			return
		}

		d.mu.Lock()
		d.logs = append(d.logs, event.Entry)
		handler := d.handlers.log
		d.mu.Unlock()

		if handler != nil {
			handler(&event.Entry)
		}
	})

	d.enabled = true
	return nil
}

// Disable stops capturing console messages.
func (d *ConsoleDebugger) Disable(ctx context.Context) error {
	d.mu.Lock()
	defer d.mu.Unlock()

	if !d.enabled {
		return nil
	}

	_, _ = d.client.Send(ctx, RuntimeDisable, nil)
	_, _ = d.client.Send(ctx, "Log.disable", nil)

	d.enabled = false
	return nil
}

// OnConsole registers a handler for console messages.
func (d *ConsoleDebugger) OnConsole(handler func(*ConsoleEntry)) {
	d.mu.Lock()
	d.handlers.console = handler
	d.mu.Unlock()
}

// OnException registers a handler for exceptions.
func (d *ConsoleDebugger) OnException(handler func(*ExceptionDetails)) {
	d.mu.Lock()
	d.handlers.exception = handler
	d.mu.Unlock()
}

// OnLog registers a handler for browser log entries.
func (d *ConsoleDebugger) OnLog(handler func(*LogEntry)) {
	d.mu.Lock()
	d.handlers.log = handler
	d.mu.Unlock()
}

// Entries returns all captured console entries.
func (d *ConsoleDebugger) Entries() []ConsoleEntry {
	d.mu.RLock()
	defer d.mu.RUnlock()
	result := make([]ConsoleEntry, len(d.entries))
	copy(result, d.entries)
	return result
}

// Errors returns all captured exceptions.
func (d *ConsoleDebugger) Errors() []ExceptionDetails {
	d.mu.RLock()
	defer d.mu.RUnlock()
	result := make([]ExceptionDetails, len(d.errors))
	copy(result, d.errors)
	return result
}

// Logs returns all captured browser log entries.
func (d *ConsoleDebugger) Logs() []LogEntry {
	d.mu.RLock()
	defer d.mu.RUnlock()
	result := make([]LogEntry, len(d.logs))
	copy(result, d.logs)
	return result
}

// Clear clears all captured entries.
func (d *ConsoleDebugger) Clear() {
	d.mu.Lock()
	d.entries = make([]ConsoleEntry, 0)
	d.errors = make([]ExceptionDetails, 0)
	d.logs = make([]LogEntry, 0)
	d.mu.Unlock()
}

// IsEnabled returns whether the debugger is enabled.
func (d *ConsoleDebugger) IsEnabled() bool {
	d.mu.RLock()
	defer d.mu.RUnlock()
	return d.enabled
}

// formatConsoleArgs formats console arguments into a readable string.
func formatConsoleArgs(args []RemoteObject) string {
	if len(args) == 0 {
		return ""
	}

	var parts []string
	for _, arg := range args {
		switch arg.Type {
		case "string":
			if s, ok := arg.Value.(string); ok {
				parts = append(parts, s)
			} else {
				parts = append(parts, fmt.Sprintf("%v", arg.Value))
			}
		case "number", "boolean":
			parts = append(parts, fmt.Sprintf("%v", arg.Value))
		case "undefined":
			parts = append(parts, "undefined")
		case "object":
			if arg.Subtype == "null" {
				parts = append(parts, "null")
			} else if arg.Description != "" {
				parts = append(parts, arg.Description)
			} else if arg.ClassName != "" {
				parts = append(parts, arg.ClassName)
			} else {
				parts = append(parts, "[object]")
			}
		case "function":
			if arg.Description != "" {
				parts = append(parts, arg.Description)
			} else {
				parts = append(parts, "[function]")
			}
		default:
			if arg.Description != "" {
				parts = append(parts, arg.Description)
			} else {
				parts = append(parts, fmt.Sprintf("[%s]", arg.Type))
			}
		}
	}

	result := ""
	for i, p := range parts {
		if i > 0 {
			result += " "
		}
		result += p
	}
	return result
}
