package mcp

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Dialog Handling Tools

// HandleDialogInput for handling browser dialogs.
type HandleDialogInput struct {
	Action     string `json:"action" jsonschema:"Action to take: accept or dismiss,enum=accept,enum=dismiss,required"`
	PromptText string `json:"prompt_text,omitempty" jsonschema:"Text to enter for prompt dialogs (only used with accept action)"`
}

// HandleDialogOutput confirms the dialog action.
type HandleDialogOutput struct {
	Message string `json:"message"`
	Action  string `json:"action"`
}

func (s *Server) handleHandleDialog(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input HandleDialogInput,
) (*mcp.CallToolResult, HandleDialogOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, HandleDialogOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	accept := input.Action == "accept"

	err = vibe.HandleDialog(ctx, accept, input.PromptText)
	if err != nil {
		return nil, HandleDialogOutput{}, fmt.Errorf("failed to handle dialog: %w", err)
	}

	action := "accepted"
	if !accept {
		action = "dismissed"
	}

	return nil, HandleDialogOutput{
		Message: fmt.Sprintf("Dialog %s", action),
		Action:  action,
	}, nil
}

// GetDialogInput for getting the current dialog state.
type GetDialogInput struct{}

// GetDialogOutput contains dialog information.
type GetDialogOutput struct {
	HasDialog    bool   `json:"has_dialog"`
	DialogType   string `json:"dialog_type,omitempty"`
	Message      string `json:"message,omitempty"`
	DefaultValue string `json:"default_value,omitempty"`
}

func (s *Server) handleGetDialog(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetDialogInput,
) (*mcp.CallToolResult, GetDialogOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, GetDialogOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	info, err := vibe.GetDialog(ctx)
	if err != nil {
		return nil, GetDialogOutput{}, fmt.Errorf("failed to get dialog: %w", err)
	}

	return nil, GetDialogOutput{
		HasDialog:    info.HasDialog,
		DialogType:   info.Type,
		Message:      info.Message,
		DefaultValue: info.DefaultValue,
	}, nil
}
