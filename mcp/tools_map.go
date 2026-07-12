package mcp

import (
	"context"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	vibium "github.com/plexusone/w3pilot"
	"github.com/plexusone/w3pilot/mcp/report"
)

// PageMapInput represents input for the page_map tool.
type PageMapInput struct {
	IncludeHidden bool   `json:"include_hidden" jsonschema:"Include hidden elements in the mapping"`
	MaxElements   int    `json:"max_elements" jsonschema:"Maximum number of elements to map (0 = no limit, default: 100)"`
	Scope         string `json:"scope" jsonschema:"CSS selector to limit mapping scope (e.g. '#main-content')"`
}

// PageMapOutput represents output from the page_map tool.
type PageMapOutput struct {
	Elements []ElementRefOutput `json:"elements"`
	Count    int                `json:"count"`
	Message  string             `json:"message"`
}

// ElementRefOutput is the MCP-friendly version of ElementRef.
type ElementRefOutput struct {
	Ref         string `json:"ref"`
	Tag         string `json:"tag"`
	Role        string `json:"role,omitempty"`
	Text        string `json:"text,omitempty"`
	Label       string `json:"label,omitempty"`
	Placeholder string `json:"placeholder,omitempty"`
	Type        string `json:"type,omitempty"`
	Selector    string `json:"selector"`
	Visible     bool   `json:"visible"`
	Enabled     bool   `json:"enabled"`
	Display     string `json:"display"`
}

func (s *Server) handlePageMap(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input PageMapInput,
) (*mcp.CallToolResult, PageMapOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, PageMapOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	// Default max elements
	if input.MaxElements == 0 {
		input.MaxElements = 100
	}

	start := time.Now()
	refs, err := pilot.MapElements(ctx, &vibium.MapOptions{
		IncludeHidden: input.IncludeHidden,
		MaxElements:   input.MaxElements,
		Scope:         input.Scope,
	})
	duration := time.Since(start)

	result := report.StepResult{
		ID:     s.session.NextStepID("page_map"),
		Action: "page_map",
		Args: map[string]any{
			"include_hidden": input.IncludeHidden,
			"max_elements":   input.MaxElements,
			"scope":          input.Scope,
		},
		DurationMS: duration.Milliseconds(),
	}

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityMedium
		result.Error = &report.StepError{
			Type:    "MapError",
			Message: err.Error(),
		}
		s.session.RecordStep(result)
		return nil, PageMapOutput{}, fmt.Errorf("failed to map elements: %w", err)
	}

	// Store refs in session for later use
	s.StoreRefs(refs)

	// Convert to output format
	elements := make([]ElementRefOutput, len(refs))
	for i, ref := range refs {
		elements[i] = ElementRefOutput{
			Ref:         ref.Ref,
			Tag:         ref.Tag,
			Role:        ref.Role,
			Text:        ref.Text,
			Label:       ref.Label,
			Placeholder: ref.Placeholder,
			Type:        ref.Type,
			Selector:    ref.Selector,
			Visible:     ref.Visible,
			Enabled:     ref.Enabled,
			Display:     ref.FormatRef(),
		}
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	result.Result = map[string]any{
		"count": len(refs),
	}
	s.session.RecordStep(result)

	return nil, PageMapOutput{
		Elements: elements,
		Count:    len(refs),
		Message:  fmt.Sprintf("Mapped %d interactive elements", len(refs)),
	}, nil
}

// PageMapClearInput represents input for the page_map_clear tool.
type PageMapClearInput struct{}

// PageMapClearOutput represents output from the page_map_clear tool.
type PageMapClearOutput struct {
	Message string `json:"message"`
}

func (s *Server) handlePageMapClear(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input PageMapClearInput,
) (*mcp.CallToolResult, PageMapClearOutput, error) {
	s.ClearRefs()
	return nil, PageMapClearOutput{Message: "Element references cleared"}, nil
}

// PageMapGetInput represents input for the page_map_get tool.
type PageMapGetInput struct {
	Ref string `json:"ref" jsonschema:"Element reference (e.g. @e1, @e2),required"`
}

// PageMapGetOutput represents output from the page_map_get tool.
type PageMapGetOutput struct {
	Element *ElementRefOutput `json:"element,omitempty"`
	Found   bool              `json:"found"`
	Message string            `json:"message"`
}

func (s *Server) handlePageMapGet(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input PageMapGetInput,
) (*mcp.CallToolResult, PageMapGetOutput, error) {
	ref, ok := s.GetRef(input.Ref)
	if !ok {
		return nil, PageMapGetOutput{
			Found:   false,
			Message: fmt.Sprintf("Element reference %s not found. Run page_map to refresh.", input.Ref),
		}, nil
	}

	return nil, PageMapGetOutput{
		Element: &ElementRefOutput{
			Ref:         ref.Ref,
			Tag:         ref.Tag,
			Role:        ref.Role,
			Text:        ref.Text,
			Label:       ref.Label,
			Placeholder: ref.Placeholder,
			Type:        ref.Type,
			Selector:    ref.Selector,
			Visible:     ref.Visible,
			Enabled:     ref.Enabled,
			Display:     ref.FormatRef(),
		},
		Found:   true,
		Message: fmt.Sprintf("Found %s", input.Ref),
	}, nil
}

// PageMapDiffInput represents input for the page_map_diff tool.
type PageMapDiffInput struct {
	IncludeHidden bool   `json:"include_hidden" jsonschema:"Include hidden elements in the new mapping"`
	MaxElements   int    `json:"max_elements" jsonschema:"Maximum number of elements to map (0 = no limit, default: 100)"`
	Scope         string `json:"scope" jsonschema:"CSS selector to limit mapping scope (e.g. '#main-content')"`
}

// PageMapDiffOutput represents output from the page_map_diff tool.
type PageMapDiffOutput struct {
	Added      []ElementRefOutput `json:"added"`
	Removed    []ElementRefOutput `json:"removed"`
	Changed    []RefChangeOutput  `json:"changed"`
	Unchanged  int                `json:"unchanged"`
	HasChanges bool               `json:"has_changes"`
	Message    string             `json:"message"`
}

// RefChangeOutput represents a changed element for MCP output.
type RefChangeOutput struct {
	Before ElementRefOutput `json:"before"`
	After  ElementRefOutput `json:"after"`
}

func (s *Server) handlePageMapDiff(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input PageMapDiffInput,
) (*mcp.CallToolResult, PageMapDiffOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, PageMapDiffOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	// Get current refs (before)
	before := s.refs.All()

	// Default max elements
	if input.MaxElements == 0 {
		input.MaxElements = 100
	}

	start := time.Now()
	after, err := pilot.MapElements(ctx, &vibium.MapOptions{
		IncludeHidden: input.IncludeHidden,
		MaxElements:   input.MaxElements,
		Scope:         input.Scope,
	})
	duration := time.Since(start)

	result := report.StepResult{
		ID:     s.session.NextStepID("page_map_diff"),
		Action: "page_map_diff",
		Args: map[string]any{
			"include_hidden": input.IncludeHidden,
			"max_elements":   input.MaxElements,
			"scope":          input.Scope,
		},
		DurationMS: duration.Milliseconds(),
	}

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityMedium
		result.Error = &report.StepError{
			Type:    "MapDiffError",
			Message: err.Error(),
		}
		s.session.RecordStep(result)
		return nil, PageMapDiffOutput{}, fmt.Errorf("failed to map elements: %w", err)
	}

	// Calculate diff
	diff := vibium.DiffRefs(before, after)

	// Update stored refs with new mapping
	s.StoreRefs(after)

	// Convert to output format
	output := PageMapDiffOutput{
		Added:      make([]ElementRefOutput, len(diff.Added)),
		Removed:    make([]ElementRefOutput, len(diff.Removed)),
		Changed:    make([]RefChangeOutput, len(diff.Changed)),
		Unchanged:  diff.Summary.Unchanged,
		HasChanges: diff.HasChanges(),
	}

	for i, ref := range diff.Added {
		output.Added[i] = toElementRefOutput(ref)
	}

	for i, ref := range diff.Removed {
		output.Removed[i] = toElementRefOutput(ref)
	}

	for i, change := range diff.Changed {
		output.Changed[i] = RefChangeOutput{
			Before: toElementRefOutput(change.Before),
			After:  toElementRefOutput(change.After),
		}
	}

	// Build message
	if diff.HasChanges() {
		output.Message = fmt.Sprintf("Page changed: %d added, %d removed, %d moved, %d unchanged",
			diff.Summary.Added, diff.Summary.Removed, diff.Summary.Changed, diff.Summary.Unchanged)
	} else {
		output.Message = fmt.Sprintf("No changes detected (%d elements)", diff.Summary.Unchanged)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	result.Result = map[string]any{
		"added":     diff.Summary.Added,
		"removed":   diff.Summary.Removed,
		"changed":   diff.Summary.Changed,
		"unchanged": diff.Summary.Unchanged,
	}
	s.session.RecordStep(result)

	return nil, output, nil
}

func toElementRefOutput(ref vibium.ElementRef) ElementRefOutput {
	return ElementRefOutput{
		Ref:         ref.Ref,
		Tag:         ref.Tag,
		Role:        ref.Role,
		Text:        ref.Text,
		Label:       ref.Label,
		Placeholder: ref.Placeholder,
		Type:        ref.Type,
		Selector:    ref.Selector,
		Visible:     ref.Visible,
		Enabled:     ref.Enabled,
		Display:     ref.FormatRef(),
	}
}
