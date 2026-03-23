package mcp

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Tab Management Tools

// ListTabsInput for listing all open tabs.
type ListTabsInput struct{}

// TabInfo contains information about a single tab.
type TabInfo struct {
	Index int    `json:"index"`
	ID    string `json:"id"`
	URL   string `json:"url"`
	Title string `json:"title"`
}

// ListTabsOutput contains all open tabs.
type ListTabsOutput struct {
	Tabs       []TabInfo `json:"tabs"`
	Count      int       `json:"count"`
	CurrentTab int       `json:"current_tab"`
}

func (s *Server) handleListTabs(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ListTabsInput,
) (*mcp.CallToolResult, ListTabsOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, ListTabsOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	pages, err := pilot.Pages(ctx)
	if err != nil {
		return nil, ListTabsOutput{}, fmt.Errorf("failed to get pages: %w", err)
	}

	tabs := make([]TabInfo, len(pages))
	currentTab := 0

	for i, page := range pages {
		url, _ := page.URL(ctx)
		title, _ := page.Title(ctx)

		tabs[i] = TabInfo{
			Index: i,
			ID:    page.BrowsingContext(),
			URL:   url,
			Title: title,
		}

		// Track which tab is current (matches the session's vibe)
		if page.BrowsingContext() == pilot.BrowsingContext() {
			currentTab = i
		}
	}

	return nil, ListTabsOutput{
		Tabs:       tabs,
		Count:      len(tabs),
		CurrentTab: currentTab,
	}, nil
}

// SelectTabInput for switching to a specific tab.
type SelectTabInput struct {
	Index *int   `json:"index,omitempty" jsonschema:"Tab index (0-based)"`
	ID    string `json:"id,omitempty" jsonschema:"Tab ID (from list_tabs)"`
}

// SelectTabOutput confirms the tab switch.
type SelectTabOutput struct {
	Message string `json:"message"`
	URL     string `json:"url"`
	Title   string `json:"title"`
}

func (s *Server) handleSelectTab(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SelectTabInput,
) (*mcp.CallToolResult, SelectTabOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, SelectTabOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	pages, err := pilot.Pages(ctx)
	if err != nil {
		return nil, SelectTabOutput{}, fmt.Errorf("failed to get pages: %w", err)
	}

	if len(pages) == 0 {
		return nil, SelectTabOutput{}, fmt.Errorf("no tabs available")
	}

	var targetPage *struct {
		index int
		id    string
	}

	// Find the target page by index or ID
	if input.Index != nil {
		idx := *input.Index
		if idx < 0 || idx >= len(pages) {
			return nil, SelectTabOutput{}, fmt.Errorf("tab index %d out of range (0-%d)", idx, len(pages)-1)
		}
		targetPage = &struct {
			index int
			id    string
		}{index: idx, id: pages[idx].BrowsingContext()}
	} else if input.ID != "" {
		for i, page := range pages {
			if page.BrowsingContext() == input.ID {
				targetPage = &struct {
					index int
					id    string
				}{index: i, id: input.ID}
				break
			}
		}
		if targetPage == nil {
			return nil, SelectTabOutput{}, fmt.Errorf("tab with ID %q not found", input.ID)
		}
	} else {
		return nil, SelectTabOutput{}, fmt.Errorf("either index or id must be provided")
	}

	// Switch to the target tab by updating the session's active page
	s.session.SetActiveContext(targetPage.id)

	// Get the new active page info
	newPilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, SelectTabOutput{}, fmt.Errorf("failed to get new active page: %w", err)
	}

	url, _ := newPilot.URL(ctx)
	title, _ := newPilot.Title(ctx)

	// Bring the tab to front
	_ = newPilot.BringToFront(ctx)

	return nil, SelectTabOutput{
		Message: fmt.Sprintf("Switched to tab %d", targetPage.index),
		URL:     url,
		Title:   title,
	}, nil
}

// CloseTabInput for closing a specific tab.
type CloseTabInput struct {
	Index *int   `json:"index,omitempty" jsonschema:"Tab index to close (0-based). Defaults to current tab."`
	ID    string `json:"id,omitempty" jsonschema:"Tab ID to close (from list_tabs)"`
}

// CloseTabOutput confirms the tab closure.
type CloseTabOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleCloseTab(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input CloseTabInput,
) (*mcp.CallToolResult, CloseTabOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, CloseTabOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	pages, err := pilot.Pages(ctx)
	if err != nil {
		return nil, CloseTabOutput{}, fmt.Errorf("failed to get pages: %w", err)
	}

	if len(pages) == 0 {
		return nil, CloseTabOutput{}, fmt.Errorf("no tabs available")
	}

	var targetID string
	var targetIndex int

	// Find the target page by index or ID
	if input.Index != nil {
		idx := *input.Index
		if idx < 0 || idx >= len(pages) {
			return nil, CloseTabOutput{}, fmt.Errorf("tab index %d out of range (0-%d)", idx, len(pages)-1)
		}
		targetID = pages[idx].BrowsingContext()
		targetIndex = idx
	} else if input.ID != "" {
		found := false
		for i, page := range pages {
			if page.BrowsingContext() == input.ID {
				targetID = input.ID
				targetIndex = i
				found = true
				break
			}
		}
		if !found {
			return nil, CloseTabOutput{}, fmt.Errorf("tab with ID %q not found", input.ID)
		}
	} else {
		// Close current tab
		targetID = pilot.BrowsingContext()
		for i, page := range pages {
			if page.BrowsingContext() == targetID {
				targetIndex = i
				break
			}
		}
	}

	// Find and close the target page
	for _, page := range pages {
		if page.BrowsingContext() == targetID {
			if err := page.Close(ctx); err != nil {
				return nil, CloseTabOutput{}, fmt.Errorf("failed to close tab: %w", err)
			}
			break
		}
	}

	// If we closed the current tab, switch to another
	if targetID == pilot.BrowsingContext() && len(pages) > 1 {
		// Switch to the previous tab, or the next one if this was the first
		newIndex := targetIndex - 1
		if newIndex < 0 {
			newIndex = 0
		}
		// Get fresh pages list
		newPages, err := pilot.Pages(ctx)
		if err == nil && len(newPages) > newIndex {
			s.session.SetActiveContext(newPages[newIndex].BrowsingContext())
		}
	}

	return nil, CloseTabOutput{Message: fmt.Sprintf("Closed tab %d", targetIndex)}, nil
}
