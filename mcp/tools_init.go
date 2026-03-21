package mcp

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// AddInitScript tool - add a script that runs before page scripts

type AddInitScriptInput struct {
	Script string `json:"script" jsonschema:"JavaScript code to inject before page scripts,required"`
}

type AddInitScriptOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleAddInitScript(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input AddInitScriptInput,
) (*mcp.CallToolResult, AddInitScriptOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, AddInitScriptOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if err := vibe.AddInitScript(ctx, input.Script); err != nil {
		return nil, AddInitScriptOutput{}, fmt.Errorf("failed to add init script: %w", err)
	}

	return nil, AddInitScriptOutput{
		Message: "Init script added successfully",
	}, nil
}
