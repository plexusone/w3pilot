package cdp

import (
	"sync"
	"testing"
)

func TestNewConsoleDebugger(t *testing.T) {
	c := NewClient()
	d := NewConsoleDebugger(c)

	if d == nil {
		t.Fatal("NewConsoleDebugger() returned nil")
	}
	if d.client != c {
		t.Error("NewConsoleDebugger().client should be the provided client")
	}
	if d.enabled {
		t.Error("NewConsoleDebugger() should not be enabled")
	}
	if len(d.entries) != 0 {
		t.Error("NewConsoleDebugger().entries should be empty")
	}
	if len(d.errors) != 0 {
		t.Error("NewConsoleDebugger().errors should be empty")
	}
	if len(d.logs) != 0 {
		t.Error("NewConsoleDebugger().logs should be empty")
	}
}

func TestConsoleDebuggerIsEnabled(t *testing.T) {
	c := NewClient()
	d := NewConsoleDebugger(c)

	if d.IsEnabled() {
		t.Error("IsEnabled() should return false initially")
	}

	d.mu.Lock()
	d.enabled = true
	d.mu.Unlock()

	if !d.IsEnabled() {
		t.Error("IsEnabled() should return true when enabled")
	}
}

func TestConsoleDebuggerEntries(t *testing.T) {
	c := NewClient()
	d := NewConsoleDebugger(c)

	// Add some entries
	d.mu.Lock()
	d.entries = []ConsoleEntry{
		{Type: ConsoleLog, Text: "log1"},
		{Type: ConsoleError, Text: "error1"},
	}
	d.mu.Unlock()

	entries := d.Entries()
	if len(entries) != 2 {
		t.Errorf("Entries() returned %d entries, want 2", len(entries))
	}

	// Verify it's a copy
	entries[0].Text = "modified"
	origEntries := d.Entries()
	if origEntries[0].Text == "modified" {
		t.Error("Entries() should return a copy")
	}
}

func TestConsoleDebuggerErrors(t *testing.T) {
	c := NewClient()
	d := NewConsoleDebugger(c)

	// Add some errors
	d.mu.Lock()
	d.errors = []ExceptionDetails{
		{ExceptionID: 1, Text: "error1"},
		{ExceptionID: 2, Text: "error2"},
	}
	d.mu.Unlock()

	errors := d.Errors()
	if len(errors) != 2 {
		t.Errorf("Errors() returned %d errors, want 2", len(errors))
	}

	// Verify it's a copy
	errors[0].Text = "modified"
	origErrors := d.Errors()
	if origErrors[0].Text == "modified" {
		t.Error("Errors() should return a copy")
	}
}

func TestConsoleDebuggerLogs(t *testing.T) {
	c := NewClient()
	d := NewConsoleDebugger(c)

	// Add some logs
	d.mu.Lock()
	d.logs = []LogEntry{
		{Level: "warning", Text: "deprecation warning"},
		{Level: "error", Text: "intervention"},
	}
	d.mu.Unlock()

	logs := d.Logs()
	if len(logs) != 2 {
		t.Errorf("Logs() returned %d logs, want 2", len(logs))
	}

	// Verify it's a copy
	logs[0].Text = "modified"
	origLogs := d.Logs()
	if origLogs[0].Text == "modified" {
		t.Error("Logs() should return a copy")
	}
}

func TestConsoleDebuggerClear(t *testing.T) {
	c := NewClient()
	d := NewConsoleDebugger(c)

	// Add entries, errors, and logs
	d.mu.Lock()
	d.entries = []ConsoleEntry{{Text: "entry"}}
	d.errors = []ExceptionDetails{{Text: "error"}}
	d.logs = []LogEntry{{Text: "log"}}
	d.mu.Unlock()

	d.Clear()

	if len(d.Entries()) != 0 {
		t.Error("Clear() should empty entries")
	}
	if len(d.Errors()) != 0 {
		t.Error("Clear() should empty errors")
	}
	if len(d.Logs()) != 0 {
		t.Error("Clear() should empty logs")
	}
}

func TestConsoleDebuggerOnConsole(t *testing.T) {
	c := NewClient()
	d := NewConsoleDebugger(c)

	var received *ConsoleEntry
	d.OnConsole(func(entry *ConsoleEntry) {
		received = entry
	})

	d.mu.RLock()
	handler := d.handlers.console
	d.mu.RUnlock()

	if handler == nil {
		t.Error("OnConsole() should set the handler")
	}

	// Call the handler
	testEntry := &ConsoleEntry{Text: "test"}
	handler(testEntry)

	if received == nil || received.Text != "test" {
		t.Error("handler should receive the entry")
	}
}

func TestConsoleDebuggerOnException(t *testing.T) {
	c := NewClient()
	d := NewConsoleDebugger(c)

	var received *ExceptionDetails
	d.OnException(func(details *ExceptionDetails) {
		received = details
	})

	d.mu.RLock()
	handler := d.handlers.exception
	d.mu.RUnlock()

	if handler == nil {
		t.Error("OnException() should set the handler")
	}

	// Call the handler
	testDetails := &ExceptionDetails{Text: "test exception"}
	handler(testDetails)

	if received == nil || received.Text != "test exception" {
		t.Error("handler should receive the exception details")
	}
}

func TestConsoleDebuggerOnLog(t *testing.T) {
	c := NewClient()
	d := NewConsoleDebugger(c)

	var received *LogEntry
	d.OnLog(func(entry *LogEntry) {
		received = entry
	})

	d.mu.RLock()
	handler := d.handlers.log
	d.mu.RUnlock()

	if handler == nil {
		t.Error("OnLog() should set the handler")
	}

	// Call the handler
	testEntry := &LogEntry{Text: "test log"}
	handler(testEntry)

	if received == nil || received.Text != "test log" {
		t.Error("handler should receive the log entry")
	}
}

func TestFormatConsoleArgs(t *testing.T) {
	tests := []struct {
		name     string
		args     []RemoteObject
		expected string
	}{
		{
			name:     "empty args",
			args:     []RemoteObject{},
			expected: "",
		},
		{
			name: "single string",
			args: []RemoteObject{
				{Type: "string", Value: "hello"},
			},
			expected: "hello",
		},
		{
			name: "multiple strings",
			args: []RemoteObject{
				{Type: "string", Value: "hello"},
				{Type: "string", Value: "world"},
			},
			expected: "hello world",
		},
		{
			name: "number",
			args: []RemoteObject{
				{Type: "number", Value: 42.0},
			},
			expected: "42",
		},
		{
			name: "boolean",
			args: []RemoteObject{
				{Type: "boolean", Value: true},
			},
			expected: "true",
		},
		{
			name: "undefined",
			args: []RemoteObject{
				{Type: "undefined"},
			},
			expected: "undefined",
		},
		{
			name: "null object",
			args: []RemoteObject{
				{Type: "object", Subtype: "null"},
			},
			expected: "null",
		},
		{
			name: "object with description",
			args: []RemoteObject{
				{Type: "object", Description: "Array(3)"},
			},
			expected: "Array(3)",
		},
		{
			name: "object with className",
			args: []RemoteObject{
				{Type: "object", ClassName: "HTMLDivElement"},
			},
			expected: "HTMLDivElement",
		},
		{
			name: "object without description or className",
			args: []RemoteObject{
				{Type: "object"},
			},
			expected: "[object]",
		},
		{
			name: "function with description",
			args: []RemoteObject{
				{Type: "function", Description: "function foo() {}"},
			},
			expected: "function foo() {}",
		},
		{
			name: "function without description",
			args: []RemoteObject{
				{Type: "function"},
			},
			expected: "[function]",
		},
		{
			name: "unknown type with description",
			args: []RemoteObject{
				{Type: "symbol", Description: "Symbol(foo)"},
			},
			expected: "Symbol(foo)",
		},
		{
			name: "unknown type without description",
			args: []RemoteObject{
				{Type: "bigint"},
			},
			expected: "[bigint]",
		},
		{
			name: "mixed types",
			args: []RemoteObject{
				{Type: "string", Value: "Count:"},
				{Type: "number", Value: 42.0},
				{Type: "boolean", Value: true},
			},
			expected: "Count: 42 true",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatConsoleArgs(tt.args)
			if got != tt.expected {
				t.Errorf("formatConsoleArgs() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestConsoleDebuggerConcurrency(t *testing.T) {
	c := NewClient()
	d := NewConsoleDebugger(c)

	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 100; j++ {
				_ = d.IsEnabled()
				_ = d.Entries()
				_ = d.Errors()
				_ = d.Logs()
				d.OnConsole(func(e *ConsoleEntry) {})
				d.OnException(func(e *ExceptionDetails) {})
				d.OnLog(func(e *LogEntry) {})
			}
		}()
	}
	wg.Wait()
}

func TestConsoleMessageTypes(t *testing.T) {
	// Verify all console message types have correct values
	types := map[ConsoleMessageType]string{
		ConsoleLog:     "log",
		ConsoleDebug:   "debug",
		ConsoleInfo:    "info",
		ConsoleError:   "error",
		ConsoleWarning: "warning",
		ConsoleDir:     "dir",
		ConsoleDirXML:  "dirxml",
		ConsoleTable:   "table",
		ConsoleTrace:   "trace",
		ConsoleClear:   "clear",
		ConsoleAssert:  "assert",
	}

	for typ, expected := range types {
		if string(typ) != expected {
			t.Errorf("ConsoleMessageType %q should be %q", typ, expected)
		}
	}
}

func TestCallFrame(t *testing.T) {
	frame := CallFrame{
		FunctionName: "testFunc",
		ScriptID:     "123",
		URL:          "http://example.com/script.js",
		LineNumber:   10,
		ColumnNumber: 5,
	}

	if frame.FunctionName != "testFunc" {
		t.Errorf("FunctionName = %q, want %q", frame.FunctionName, "testFunc")
	}
	if frame.LineNumber != 10 {
		t.Errorf("LineNumber = %d, want %d", frame.LineNumber, 10)
	}
}

func TestStackTrace(t *testing.T) {
	parent := &StackTrace{
		Description: "parent stack",
		CallFrames: []CallFrame{
			{FunctionName: "parentFunc"},
		},
	}

	stack := StackTrace{
		Description: "child stack",
		CallFrames: []CallFrame{
			{FunctionName: "childFunc"},
		},
		Parent: parent,
	}

	if stack.Parent == nil {
		t.Error("Parent should not be nil")
	}
	if stack.Parent.Description != "parent stack" {
		t.Errorf("Parent.Description = %q, want %q", stack.Parent.Description, "parent stack")
	}
}
