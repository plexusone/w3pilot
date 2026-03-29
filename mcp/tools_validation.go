package mcp

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	w3pilot "github.com/plexusone/w3pilot"
)

// ValidateSelectors tool - validate multiple selectors before use

type ValidateSelectorsInput struct {
	Selectors []string `json:"selectors" jsonschema:"List of CSS selectors to validate,required"`
}

type ValidateSelectorsOutput struct {
	Results []w3pilot.SelectorValidation `json:"results"`
	Summary ValidationSummary            `json:"summary"`
}

type ValidationSummary struct {
	Total   int `json:"total"`
	Found   int `json:"found"`
	Missing int `json:"missing"`
	Visible int `json:"visible"`
}

func (s *Server) handleValidateSelectors(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ValidateSelectorsInput,
) (*mcp.CallToolResult, ValidateSelectorsOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, ValidateSelectorsOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	results, err := pilot.ValidateSelectors(ctx, input.Selectors)
	if err != nil {
		return nil, ValidateSelectorsOutput{}, fmt.Errorf("validation failed: %w", err)
	}

	// Build summary
	summary := ValidationSummary{Total: len(results)}
	for _, r := range results {
		if r.Found {
			summary.Found++
			if r.Visible {
				summary.Visible++
			}
		} else {
			summary.Missing++
		}
	}

	return nil, ValidateSelectorsOutput{
		Results: results,
		Summary: summary,
	}, nil
}
