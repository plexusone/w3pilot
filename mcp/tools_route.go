package mcp

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	vibium "github.com/plexusone/webpilot"
)

// Route tool - register a mock response for a URL pattern

type RouteInput struct {
	Pattern     string            `json:"pattern" jsonschema:"URL pattern to match (glob or regex e.g. **/api/* or /api/.*),required"`
	Status      int               `json:"status" jsonschema:"HTTP status code (default: 200)"`
	Body        string            `json:"body" jsonschema:"Response body content"`
	ContentType string            `json:"content_type" jsonschema:"Content-Type header (default: application/json)"`
	Headers     map[string]string `json:"headers" jsonschema:"Additional response headers"`
}

type RouteOutput struct {
	Message string `json:"message"`
	Pattern string `json:"pattern"`
}

func (s *Server) handleRoute(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input RouteInput,
) (*mcp.CallToolResult, RouteOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, RouteOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	opts := vibium.MockRouteOptions{
		Status:      input.Status,
		Body:        input.Body,
		ContentType: input.ContentType,
		Headers:     input.Headers,
	}

	err = pilot.MockRoute(ctx, input.Pattern, opts)
	if err != nil {
		return nil, RouteOutput{}, fmt.Errorf("route failed: %w", err)
	}

	status := input.Status
	if status == 0 {
		status = 200
	}

	return nil, RouteOutput{
		Message: fmt.Sprintf("Route registered for %s (status: %d)", input.Pattern, status),
		Pattern: input.Pattern,
	}, nil
}

// RouteList tool - list all active routes

type RouteListInput struct{}

type RouteListOutput struct {
	Routes []RouteInfoOutput `json:"routes"`
	Count  int               `json:"count"`
}

type RouteInfoOutput struct {
	Pattern     string `json:"pattern"`
	Status      int    `json:"status,omitempty"`
	ContentType string `json:"content_type,omitempty"`
}

func (s *Server) handleRouteList(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input RouteListInput,
) (*mcp.CallToolResult, RouteListOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, RouteListOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	routes, err := pilot.ListRoutes(ctx)
	if err != nil {
		return nil, RouteListOutput{}, fmt.Errorf("list routes failed: %w", err)
	}

	output := make([]RouteInfoOutput, len(routes))
	for i, r := range routes {
		output[i] = RouteInfoOutput{
			Pattern:     r.Pattern,
			Status:      r.Status,
			ContentType: r.ContentType,
		}
	}

	return nil, RouteListOutput{
		Routes: output,
		Count:  len(routes),
	}, nil
}

// Unroute tool - remove a route handler

type UnrouteInput struct {
	Pattern string `json:"pattern" jsonschema:"URL pattern to unregister,required"`
}

type UnrouteOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleUnroute(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input UnrouteInput,
) (*mcp.CallToolResult, UnrouteOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, UnrouteOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	err = pilot.Unroute(ctx, input.Pattern)
	if err != nil {
		return nil, UnrouteOutput{}, fmt.Errorf("unroute failed: %w", err)
	}

	return nil, UnrouteOutput{
		Message: fmt.Sprintf("Route removed for %s", input.Pattern),
	}, nil
}

// NetworkStateSet tool - set offline mode

type NetworkStateSetInput struct {
	Offline bool `json:"offline" jsonschema:"Set to true to enable offline mode,required"`
}

type NetworkStateSetOutput struct {
	Message string `json:"message"`
	Offline bool   `json:"offline"`
}

func (s *Server) handleNetworkStateSet(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input NetworkStateSetInput,
) (*mcp.CallToolResult, NetworkStateSetOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, NetworkStateSetOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	err = pilot.SetOffline(ctx, input.Offline)
	if err != nil {
		return nil, NetworkStateSetOutput{}, fmt.Errorf("set network state failed: %w", err)
	}

	msg := "Network online"
	if input.Offline {
		msg = "Network offline"
	}

	return nil, NetworkStateSetOutput{
		Message: msg,
		Offline: input.Offline,
	}, nil
}
