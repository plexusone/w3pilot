package mcp

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"github.com/plexusone/w3pilot/state"
)

// StateSave tool - save browser state to a named snapshot

type StateSaveInput struct {
	Name string `json:"name" jsonschema:"Name for the state snapshot (alphanumeric dash underscore),required"`
}

type StateSaveOutput struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func (s *Server) handleStateSave(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input StateSaveInput,
) (*mcp.CallToolResult, StateSaveOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, StateSaveOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	// Get current state
	storageState, err := pilot.StorageState(ctx)
	if err != nil {
		return nil, StateSaveOutput{}, fmt.Errorf("failed to get storage state: %w", err)
	}

	// Save to file
	mgr, err := state.NewManager("")
	if err != nil {
		return nil, StateSaveOutput{}, fmt.Errorf("failed to create state manager: %w", err)
	}

	if err := mgr.Save(input.Name, storageState); err != nil {
		return nil, StateSaveOutput{}, fmt.Errorf("failed to save state: %w", err)
	}

	return nil, StateSaveOutput{
		Name:    input.Name,
		Message: fmt.Sprintf("State saved as '%s'", input.Name),
	}, nil
}

// StateLoad tool - load browser state from a named snapshot

type StateLoadInput struct {
	Name string `json:"name" jsonschema:"Name of the state snapshot to load,required"`
}

type StateLoadOutput struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func (s *Server) handleStateLoad(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input StateLoadInput,
) (*mcp.CallToolResult, StateLoadOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, StateLoadOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	// Load from file
	mgr, err := state.NewManager("")
	if err != nil {
		return nil, StateLoadOutput{}, fmt.Errorf("failed to create state manager: %w", err)
	}

	storageState, err := mgr.Load(input.Name)
	if err != nil {
		return nil, StateLoadOutput{}, fmt.Errorf("failed to load state: %w", err)
	}

	// Apply state
	if err := pilot.SetStorageState(ctx, storageState); err != nil {
		return nil, StateLoadOutput{}, fmt.Errorf("failed to apply storage state: %w", err)
	}

	return nil, StateLoadOutput{
		Name:    input.Name,
		Message: fmt.Sprintf("State '%s' loaded", input.Name),
	}, nil
}

// StateList tool - list all saved state snapshots

type StateListInput struct{}

type StateListOutput struct {
	States []state.StateInfo `json:"states"`
	Count  int               `json:"count"`
}

func (s *Server) handleStateList(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input StateListInput,
) (*mcp.CallToolResult, StateListOutput, error) {
	mgr, err := state.NewManager("")
	if err != nil {
		return nil, StateListOutput{}, fmt.Errorf("failed to create state manager: %w", err)
	}

	states, err := mgr.List()
	if err != nil {
		return nil, StateListOutput{}, fmt.Errorf("failed to list states: %w", err)
	}

	return nil, StateListOutput{
		States: states,
		Count:  len(states),
	}, nil
}

// StateDelete tool - delete a saved state snapshot

type StateDeleteInput struct {
	Name string `json:"name" jsonschema:"Name of the state snapshot to delete,required"`
}

type StateDeleteOutput struct {
	Name    string `json:"name"`
	Message string `json:"message"`
}

func (s *Server) handleStateDelete(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input StateDeleteInput,
) (*mcp.CallToolResult, StateDeleteOutput, error) {
	mgr, err := state.NewManager("")
	if err != nil {
		return nil, StateDeleteOutput{}, fmt.Errorf("failed to create state manager: %w", err)
	}

	if err := mgr.Delete(input.Name); err != nil {
		return nil, StateDeleteOutput{}, fmt.Errorf("failed to delete state: %w", err)
	}

	return nil, StateDeleteOutput{
		Name:    input.Name,
		Message: fmt.Sprintf("State '%s' deleted", input.Name),
	}, nil
}
