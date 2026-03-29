package mcp

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	w3pilot "github.com/plexusone/w3pilot"
)

// PageInspect tool - inspect page elements for AI agents

type PageInspectInput struct {
	IncludeButtons  *bool `json:"include_buttons,omitempty" jsonschema:"Include button elements (default true)"`
	IncludeLinks    *bool `json:"include_links,omitempty" jsonschema:"Include link elements (default true)"`
	IncludeInputs   *bool `json:"include_inputs,omitempty" jsonschema:"Include input elements (default true)"`
	IncludeSelects  *bool `json:"include_selects,omitempty" jsonschema:"Include select elements (default true)"`
	IncludeHeadings *bool `json:"include_headings,omitempty" jsonschema:"Include heading elements (default true)"`
	IncludeImages   *bool `json:"include_images,omitempty" jsonschema:"Include images with alt text (default true)"`
	MaxItems        int   `json:"max_items,omitempty" jsonschema:"Maximum items per category (default 50)"`
}

type PageInspectOutput struct {
	URL      string                   `json:"url"`
	Title    string                   `json:"title"`
	Buttons  []w3pilot.InspectButton  `json:"buttons,omitempty"`
	Links    []w3pilot.InspectLink    `json:"links,omitempty"`
	Inputs   []w3pilot.InspectInput   `json:"inputs,omitempty"`
	Selects  []w3pilot.InspectSelect  `json:"selects,omitempty"`
	Headings []w3pilot.InspectHeading `json:"headings,omitempty"`
	Images   []w3pilot.InspectImage   `json:"images,omitempty"`
	Summary  w3pilot.InspectSummary   `json:"summary"`
}

func (s *Server) handlePageInspect(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input PageInspectInput,
) (*mcp.CallToolResult, PageInspectOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, PageInspectOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	opts := w3pilot.DefaultInspectOptions()

	// Apply overrides
	if input.IncludeButtons != nil {
		opts.IncludeButtons = *input.IncludeButtons
	}
	if input.IncludeLinks != nil {
		opts.IncludeLinks = *input.IncludeLinks
	}
	if input.IncludeInputs != nil {
		opts.IncludeInputs = *input.IncludeInputs
	}
	if input.IncludeSelects != nil {
		opts.IncludeSelects = *input.IncludeSelects
	}
	if input.IncludeHeadings != nil {
		opts.IncludeHeadings = *input.IncludeHeadings
	}
	if input.IncludeImages != nil {
		opts.IncludeImages = *input.IncludeImages
	}
	if input.MaxItems > 0 {
		opts.MaxItems = input.MaxItems
	}

	result, err := pilot.Inspect(ctx, opts)
	if err != nil {
		return nil, PageInspectOutput{}, fmt.Errorf("inspection failed: %w", err)
	}

	return nil, PageInspectOutput{
		URL:      result.URL,
		Title:    result.Title,
		Buttons:  result.Buttons,
		Links:    result.Links,
		Inputs:   result.Inputs,
		Selects:  result.Selects,
		Headings: result.Headings,
		Images:   result.Images,
		Summary:  result.Summary,
	}, nil
}
