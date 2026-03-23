package mcp

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	vibium "github.com/plexusone/webpilot"
)

// StartTrace tool - start trace recording

type StartTraceInput struct {
	Name        string `json:"name" jsonschema:"Trace name (used for file naming)"`
	Title       string `json:"title" jsonschema:"Title shown in trace viewer"`
	Screenshots bool   `json:"screenshots" jsonschema:"Include screenshots in trace (default: true)"`
	Snapshots   bool   `json:"snapshots" jsonschema:"Include DOM snapshots in trace (default: true)"`
	Sources     bool   `json:"sources" jsonschema:"Include source files in trace"`
}

type StartTraceOutput struct {
	Message string `json:"message"`
	Name    string `json:"name,omitempty"`
}

func (s *Server) handleStartTrace(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input StartTraceInput,
) (*mcp.CallToolResult, StartTraceOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, StartTraceOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	// Default screenshots and snapshots to true if not specified
	screenshots := input.Screenshots
	snapshots := input.Snapshots

	// If both are false and nothing was explicitly set, enable defaults
	if !screenshots && !snapshots && !input.Sources {
		screenshots = true
		snapshots = true
	}

	opts := &vibium.TracingStartOptions{
		Name:        input.Name,
		Title:       input.Title,
		Screenshots: screenshots,
		Snapshots:   snapshots,
		Sources:     input.Sources,
	}

	tracing := pilot.Tracing()
	if err := tracing.Start(ctx, opts); err != nil {
		return nil, StartTraceOutput{}, fmt.Errorf("start trace failed: %w", err)
	}

	name := input.Name
	if name == "" {
		name = "trace"
	}

	return nil, StartTraceOutput{
		Message: fmt.Sprintf("Trace recording started: %s", name),
		Name:    name,
	}, nil
}

// StopTrace tool - stop trace recording and save/return data

type StopTraceInput struct {
	Path string `json:"path" jsonschema:"File path to save the trace ZIP (optional - if not provided returns base64 data)"`
}

type StopTraceOutput struct {
	Message  string `json:"message"`
	Path     string `json:"path,omitempty"`
	Data     string `json:"data,omitempty"`
	SizeKB   int    `json:"size_kb"`
	ViewHint string `json:"view_hint,omitempty"`
}

func (s *Server) handleStopTrace(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input StopTraceInput,
) (*mcp.CallToolResult, StopTraceOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, StopTraceOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	tracing := pilot.Tracing()
	data, err := tracing.Stop(ctx, &vibium.TracingStopOptions{
		Path: input.Path,
	})
	if err != nil {
		return nil, StopTraceOutput{}, fmt.Errorf("stop trace failed: %w", err)
	}

	sizeKB := len(data) / 1024
	output := StopTraceOutput{
		SizeKB:   sizeKB,
		ViewHint: "Open trace with: npx playwright show-trace <trace.zip>",
	}

	if input.Path != "" {
		// Ensure directory exists
		dir := filepath.Dir(input.Path)
		if dir != "" && dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, StopTraceOutput{}, fmt.Errorf("failed to create directory: %w", err)
			}
		}

		// Write trace file
		if err := os.WriteFile(input.Path, data, 0600); err != nil {
			return nil, StopTraceOutput{}, fmt.Errorf("failed to write trace file: %w", err)
		}

		output.Message = fmt.Sprintf("Trace saved to %s (%d KB)", input.Path, sizeKB)
		output.Path = input.Path
	} else {
		// Return base64-encoded data
		output.Message = fmt.Sprintf("Trace recording stopped (%d KB)", sizeKB)
		output.Data = base64.StdEncoding.EncodeToString(data)
	}

	return nil, output, nil
}

// StartTraceChunk tool - start a new trace chunk

type StartTraceChunkInput struct {
	Name  string `json:"name" jsonschema:"Chunk name"`
	Title string `json:"title" jsonschema:"Chunk title shown in trace viewer"`
}

type StartTraceChunkOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleStartTraceChunk(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input StartTraceChunkInput,
) (*mcp.CallToolResult, StartTraceChunkOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, StartTraceChunkOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	tracing := pilot.Tracing()
	if err := tracing.StartChunk(ctx, &vibium.TracingChunkOptions{
		Name:  input.Name,
		Title: input.Title,
	}); err != nil {
		return nil, StartTraceChunkOutput{}, fmt.Errorf("start trace chunk failed: %w", err)
	}

	name := input.Name
	if name == "" {
		name = "chunk"
	}

	return nil, StartTraceChunkOutput{
		Message: fmt.Sprintf("Trace chunk started: %s", name),
	}, nil
}

// StopTraceChunk tool - stop the current trace chunk

type StopTraceChunkInput struct {
	Path string `json:"path" jsonschema:"File path to save the chunk ZIP (optional - if not provided returns base64 data)"`
}

type StopTraceChunkOutput struct {
	Message string `json:"message"`
	Path    string `json:"path,omitempty"`
	Data    string `json:"data,omitempty"`
	SizeKB  int    `json:"size_kb"`
}

func (s *Server) handleStopTraceChunk(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input StopTraceChunkInput,
) (*mcp.CallToolResult, StopTraceChunkOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, StopTraceChunkOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	tracing := pilot.Tracing()
	data, err := tracing.StopChunk(ctx, &vibium.TracingChunkOptions{})
	if err != nil {
		return nil, StopTraceChunkOutput{}, fmt.Errorf("stop trace chunk failed: %w", err)
	}

	sizeKB := len(data) / 1024
	output := StopTraceChunkOutput{
		SizeKB: sizeKB,
	}

	if input.Path != "" {
		// Ensure directory exists
		dir := filepath.Dir(input.Path)
		if dir != "" && dir != "." {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return nil, StopTraceChunkOutput{}, fmt.Errorf("failed to create directory: %w", err)
			}
		}

		// Write trace file
		if err := os.WriteFile(input.Path, data, 0600); err != nil {
			return nil, StopTraceChunkOutput{}, fmt.Errorf("failed to write trace file: %w", err)
		}

		output.Message = fmt.Sprintf("Trace chunk saved to %s (%d KB)", input.Path, sizeKB)
		output.Path = input.Path
	} else {
		output.Message = fmt.Sprintf("Trace chunk stopped (%d KB)", sizeKB)
		output.Data = base64.StdEncoding.EncodeToString(data)
	}

	return nil, output, nil
}

// StartTraceGroup tool - start a trace group for logical grouping

type StartTraceGroupInput struct {
	Name     string `json:"name" jsonschema:"Group name,required"`
	Location string `json:"location" jsonschema:"Source location to associate with this group"`
}

type StartTraceGroupOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleStartTraceGroup(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input StartTraceGroupInput,
) (*mcp.CallToolResult, StartTraceGroupOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, StartTraceGroupOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	tracing := pilot.Tracing()
	opts := &vibium.TracingGroupOptions{}
	if input.Location != "" {
		opts.Location = input.Location
	}

	if err := tracing.StartGroup(ctx, input.Name, opts); err != nil {
		return nil, StartTraceGroupOutput{}, fmt.Errorf("start trace group failed: %w", err)
	}

	return nil, StartTraceGroupOutput{
		Message: fmt.Sprintf("Trace group started: %s", input.Name),
	}, nil
}

// StopTraceGroup tool - stop the current trace group

type StopTraceGroupInput struct{}

type StopTraceGroupOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleStopTraceGroup(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input StopTraceGroupInput,
) (*mcp.CallToolResult, StopTraceGroupOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, StopTraceGroupOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	tracing := pilot.Tracing()
	if err := tracing.StopGroup(ctx); err != nil {
		return nil, StopTraceGroupOutput{}, fmt.Errorf("stop trace group failed: %w", err)
	}

	return nil, StopTraceGroupOutput{
		Message: "Trace group stopped",
	}, nil
}
