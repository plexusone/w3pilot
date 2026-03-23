package mcp

import (
	"context"
	"fmt"

	vibium "github.com/plexusone/webpilot"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Console Messages Tools

// GetConsoleMessagesInput for retrieving console messages.
type GetConsoleMessagesInput struct {
	Level string `json:"level,omitempty" jsonschema:"Filter by message level (log/info/warn/error/debug). Empty for all levels.,enum=log,enum=info,enum=warn,enum=error,enum=debug"`
	Clear bool   `json:"clear,omitempty" jsonschema:"Clear messages after retrieving them"`
}

// GetConsoleMessagesOutput contains console messages.
type GetConsoleMessagesOutput struct {
	Messages []ConsoleMessageInfo `json:"messages"`
	Count    int                  `json:"count"`
}

// ConsoleMessageInfo represents a console message.
type ConsoleMessageInfo struct {
	Type string   `json:"type"`
	Text string   `json:"text"`
	Args []string `json:"args,omitempty"`
	URL  string   `json:"url,omitempty"`
	Line int      `json:"line,omitempty"`
}

func (s *Server) handleGetConsoleMessages(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetConsoleMessagesInput,
) (*mcp.CallToolResult, GetConsoleMessagesOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, GetConsoleMessagesOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	messages, err := pilot.ConsoleMessages(ctx, input.Level)
	if err != nil {
		return nil, GetConsoleMessagesOutput{}, fmt.Errorf("failed to get console messages: %w", err)
	}

	// Convert to output format
	msgInfos := make([]ConsoleMessageInfo, len(messages))
	for i, msg := range messages {
		msgInfos[i] = ConsoleMessageInfo{
			Type: msg.Type,
			Text: msg.Text,
			Args: msg.Args,
			URL:  msg.URL,
			Line: msg.Line,
		}
	}

	// Clear messages if requested
	if input.Clear {
		_ = pilot.ClearConsoleMessages(ctx)
	}

	return nil, GetConsoleMessagesOutput{
		Messages: msgInfos,
		Count:    len(msgInfos),
	}, nil
}

// ClearConsoleMessagesInput for clearing console messages.
type ClearConsoleMessagesInput struct{}

// ClearConsoleMessagesOutput confirms the clear operation.
type ClearConsoleMessagesOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleClearConsoleMessages(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ClearConsoleMessagesInput,
) (*mcp.CallToolResult, ClearConsoleMessagesOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, ClearConsoleMessagesOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	err = pilot.ClearConsoleMessages(ctx)
	if err != nil {
		return nil, ClearConsoleMessagesOutput{}, fmt.Errorf("failed to clear console messages: %w", err)
	}

	return nil, ClearConsoleMessagesOutput{Message: "Console messages cleared"}, nil
}

// Network Requests Tools

// GetNetworkRequestsInput for retrieving network requests.
type GetNetworkRequestsInput struct {
	URLPattern   string `json:"url_pattern,omitempty" jsonschema:"Filter by URL pattern (glob or regex)"`
	Method       string `json:"method,omitempty" jsonschema:"Filter by HTTP method (GET/POST/PUT/DELETE/etc)"`
	ResourceType string `json:"resource_type,omitempty" jsonschema:"Filter by resource type (document/script/xhr/fetch/stylesheet/image/font/other)"`
	Clear        bool   `json:"clear,omitempty" jsonschema:"Clear requests after retrieving them"`
}

// GetNetworkRequestsOutput contains network requests.
type GetNetworkRequestsOutput struct {
	Requests []NetworkRequestInfo `json:"requests"`
	Count    int                  `json:"count"`
}

// NetworkRequestInfo represents a network request.
type NetworkRequestInfo struct {
	URL          string `json:"url"`
	Method       string `json:"method"`
	ResourceType string `json:"resource_type"`
	Status       int    `json:"status,omitempty"`
	StatusText   string `json:"status_text,omitempty"`
	ResponseSize int64  `json:"response_size,omitempty"`
}

func (s *Server) handleGetNetworkRequests(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetNetworkRequestsInput,
) (*mcp.CallToolResult, GetNetworkRequestsOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, GetNetworkRequestsOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	var opts *vibium.NetworkRequestsOptions
	if input.URLPattern != "" || input.Method != "" || input.ResourceType != "" {
		opts = &vibium.NetworkRequestsOptions{
			URLPattern:   input.URLPattern,
			Method:       input.Method,
			ResourceType: input.ResourceType,
		}
	}

	requests, err := pilot.NetworkRequests(ctx, opts)
	if err != nil {
		return nil, GetNetworkRequestsOutput{}, fmt.Errorf("failed to get network requests: %w", err)
	}

	// Convert to output format
	reqInfos := make([]NetworkRequestInfo, len(requests))
	for i, r := range requests {
		reqInfos[i] = NetworkRequestInfo{
			URL:          r.URL,
			Method:       r.Method,
			ResourceType: r.ResourceType,
			Status:       r.Status,
			StatusText:   r.StatusText,
			ResponseSize: r.ResponseSize,
		}
	}

	// Clear requests if requested
	if input.Clear {
		_ = pilot.ClearNetworkRequests(ctx)
	}

	return nil, GetNetworkRequestsOutput{
		Requests: reqInfos,
		Count:    len(reqInfos),
	}, nil
}

// ClearNetworkRequestsInput for clearing network requests.
type ClearNetworkRequestsInput struct{}

// ClearNetworkRequestsOutput confirms the clear operation.
type ClearNetworkRequestsOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleClearNetworkRequests(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ClearNetworkRequestsInput,
) (*mcp.CallToolResult, ClearNetworkRequestsOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, ClearNetworkRequestsOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	err = pilot.ClearNetworkRequests(ctx)
	if err != nil {
		return nil, ClearNetworkRequestsOutput{}, fmt.Errorf("failed to clear network requests: %w", err)
	}

	return nil, ClearNetworkRequestsOutput{Message: "Network requests cleared"}, nil
}
