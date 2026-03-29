package mcp

import (
	"context"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	w3pilot "github.com/plexusone/w3pilot"
)

// WorkflowLogin tool - automated login workflow

type WorkflowLoginInput struct {
	UsernameSelector string `json:"username_selector" jsonschema:"CSS selector for username/email field,required"`
	PasswordSelector string `json:"password_selector" jsonschema:"CSS selector for password field,required"`
	SubmitSelector   string `json:"submit_selector" jsonschema:"CSS selector for submit button,required"`
	Username         string `json:"username" jsonschema:"Username or email to enter,required"`
	Password         string `json:"password" jsonschema:"Password to enter,required"`
	SuccessIndicator string `json:"success_indicator,omitempty" jsonschema:"CSS selector or URL pattern indicating successful login"`
	TimeoutMS        int    `json:"timeout_ms,omitempty" jsonschema:"Timeout in milliseconds (default 30000)"`
}

type WorkflowLoginOutput struct {
	Success     bool   `json:"success"`
	URL         string `json:"url"`
	Title       string `json:"title"`
	Message     string `json:"message"`
	ErrorReason string `json:"error_reason,omitempty"`
}

func (s *Server) handleWorkflowLogin(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input WorkflowLoginInput,
) (*mcp.CallToolResult, WorkflowLoginOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, WorkflowLoginOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	timeout := 30 * time.Second
	if input.TimeoutMS > 0 {
		timeout = time.Duration(input.TimeoutMS) * time.Millisecond
	}

	opts := &w3pilot.LoginOptions{
		UsernameSelector: input.UsernameSelector,
		PasswordSelector: input.PasswordSelector,
		SubmitSelector:   input.SubmitSelector,
		Username:         input.Username,
		Password:         input.Password,
		SuccessIndicator: input.SuccessIndicator,
		Timeout:          timeout,
	}

	result, err := pilot.Login(ctx, opts)
	if err != nil {
		return nil, WorkflowLoginOutput{}, fmt.Errorf("login failed: %w", err)
	}

	return nil, WorkflowLoginOutput{
		Success:     result.Success,
		URL:         result.URL,
		Title:       result.Title,
		Message:     result.Message,
		ErrorReason: result.ErrorReason,
	}, nil
}

// WorkflowExtractTable tool - extract table data to JSON

type WorkflowExtractTableInput struct {
	Selector       string `json:"selector" jsonschema:"CSS selector for the table element,required"`
	IncludeHeaders *bool  `json:"include_headers,omitempty" jsonschema:"Treat first row as headers (default true)"`
	MaxRows        int    `json:"max_rows,omitempty" jsonschema:"Maximum number of rows to extract (default 1000)"`
	HeaderSelector string `json:"header_selector,omitempty" jsonschema:"Custom selector for header cells (default: th)"`
	RowSelector    string `json:"row_selector,omitempty" jsonschema:"Custom selector for data rows (default: tbody tr)"`
	CellSelector   string `json:"cell_selector,omitempty" jsonschema:"Custom selector for cells (default: td)"`
}

type WorkflowExtractTableOutput struct {
	Headers  []string            `json:"headers,omitempty"`
	Rows     [][]string          `json:"rows"`
	RowsJSON []map[string]string `json:"rows_json,omitempty"`
	RowCount int                 `json:"row_count"`
}

func (s *Server) handleWorkflowExtractTable(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input WorkflowExtractTableInput,
) (*mcp.CallToolResult, WorkflowExtractTableOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, WorkflowExtractTableOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	opts := &w3pilot.ExtractTableOptions{
		IncludeHeaders: true,
		MaxRows:        input.MaxRows,
		HeaderSelector: input.HeaderSelector,
		RowSelector:    input.RowSelector,
		CellSelector:   input.CellSelector,
	}

	if input.IncludeHeaders != nil {
		opts.IncludeHeaders = *input.IncludeHeaders
	}

	result, err := pilot.ExtractTable(ctx, input.Selector, opts)
	if err != nil {
		return nil, WorkflowExtractTableOutput{}, fmt.Errorf("table extraction failed: %w", err)
	}

	return nil, WorkflowExtractTableOutput{
		Headers:  result.Headers,
		Rows:     result.Rows,
		RowsJSON: result.RowsJSON,
		RowCount: result.RowCount,
	}, nil
}
