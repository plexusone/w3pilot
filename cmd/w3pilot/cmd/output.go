package cmd

import (
	"encoding/json"
	"fmt"
	"os"
)

// OutputFormat specifies the output format for CLI commands.
type OutputFormat string

const (
	// FormatText is the default human-readable text format.
	FormatText OutputFormat = "text"
	// FormatJSON outputs structured JSON.
	FormatJSON OutputFormat = "json"
)

var (
	// outputFormat is the global output format flag value.
	outputFormat string
)

// GetOutputFormat returns the current output format.
func GetOutputFormat() OutputFormat {
	switch outputFormat {
	case "json":
		return FormatJSON
	default:
		return FormatText
	}
}

// Output prints data in the configured format.
// For JSON format, it marshals the data to JSON.
// For text format, it uses the textFn to generate human-readable output.
// If textFn is nil, it uses fmt.Sprintf("%v", data) for text output.
func Output(data interface{}, textFn func(interface{}) string) {
	if GetOutputFormat() == FormatJSON {
		OutputJSON(data)
		return
	}

	if textFn != nil {
		fmt.Println(textFn(data))
	} else {
		fmt.Printf("%v\n", data)
	}
}

// OutputJSON prints data as formatted JSON.
func OutputJSON(data interface{}) {
	enc := json.NewEncoder(os.Stdout)
	enc.SetIndent("", "  ")
	_ = enc.Encode(data)
}

// OutputError prints an error in the configured format.
// For JSON format, it outputs a structured error object.
// For text format, it prints the error message.
func OutputError(err error) {
	if GetOutputFormat() == FormatJSON {
		errObj := map[string]interface{}{
			"error": err.Error(),
		}

		// Include additional context for known error types
		switch e := err.(type) {
		case interface{ Unwrap() error }:
			if unwrapped := e.Unwrap(); unwrapped != nil {
				errObj["cause"] = unwrapped.Error()
			}
		}

		OutputJSON(errObj)
		return
	}

	fmt.Fprintf(os.Stderr, "Error: %v\n", err)
}

// StringResult wraps a simple string result for JSON output.
type StringResult struct {
	Value string `json:"value"`
}

// BoolResult wraps a boolean result for JSON output.
type BoolResult struct {
	Value bool `json:"value"`
}

// IntResult wraps an integer result for JSON output.
type IntResult struct {
	Value int `json:"value"`
}
