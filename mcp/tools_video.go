package mcp

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	vibium "github.com/plexusone/webpilot"
)

// StartVideo tool - start video recording

type StartVideoInput struct {
	Dir    string `json:"dir,omitempty" jsonschema:"Directory to save video to"`
	Width  int    `json:"width,omitempty" jsonschema:"Video width in pixels"`
	Height int    `json:"height,omitempty" jsonschema:"Video height in pixels"`
}

type StartVideoOutput struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

func (s *Server) handleStartVideo(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input StartVideoInput,
) (*mcp.CallToolResult, StartVideoOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, StartVideoOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	opts := &vibium.VideoOptions{}
	if input.Dir != "" {
		opts.Dir = input.Dir
	}
	if input.Width > 0 && input.Height > 0 {
		opts.Size = &vibium.VideoSize{
			Width:  input.Width,
			Height: input.Height,
		}
	}

	video, err := pilot.StartVideo(ctx, opts)
	if err != nil {
		return nil, StartVideoOutput{}, fmt.Errorf("failed to start video: %w", err)
	}

	return nil, StartVideoOutput{
		Path:    video.Path(),
		Message: fmt.Sprintf("Video recording started: %s", video.Path()),
	}, nil
}

// StopVideo tool - stop video recording and return path

type StopVideoInput struct{}

type StopVideoOutput struct {
	Path    string `json:"path"`
	Message string `json:"message"`
}

func (s *Server) handleStopVideo(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input StopVideoInput,
) (*mcp.CallToolResult, StopVideoOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, StopVideoOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	path, err := pilot.StopVideo(ctx)
	if err != nil {
		return nil, StopVideoOutput{}, fmt.Errorf("failed to stop video: %w", err)
	}

	return nil, StopVideoOutput{
		Path:    path,
		Message: fmt.Sprintf("Video recording stopped: %s", path),
	}, nil
}
