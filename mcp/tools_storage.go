package mcp

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// LocalStorage tools

// LocalStorageGetInput for getting a localStorage item.
type LocalStorageGetInput struct {
	Key string `json:"key" jsonschema:"The key to get from localStorage,required"`
}

// LocalStorageGetOutput contains the retrieved value.
type LocalStorageGetOutput struct {
	Key   string  `json:"key"`
	Value *string `json:"value"` // null if key doesn't exist
}

func (s *Server) handleLocalStorageGet(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input LocalStorageGetInput,
) (*mcp.CallToolResult, LocalStorageGetOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, LocalStorageGetOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	script := fmt.Sprintf(`localStorage.getItem(%q)`, input.Key)
	result, err := pilot.Evaluate(ctx, "return "+script)
	if err != nil {
		return nil, LocalStorageGetOutput{}, fmt.Errorf("localStorage.getItem failed: %w", err)
	}

	output := LocalStorageGetOutput{Key: input.Key}
	if result != nil {
		if str, ok := result.(string); ok {
			output.Value = &str
		}
	}

	return nil, output, nil
}

// LocalStorageSetInput for setting a localStorage item.
type LocalStorageSetInput struct {
	Key   string `json:"key" jsonschema:"The key to set in localStorage,required"`
	Value string `json:"value" jsonschema:"The value to store,required"`
}

// LocalStorageSetOutput confirms the operation.
type LocalStorageSetOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleLocalStorageSet(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input LocalStorageSetInput,
) (*mcp.CallToolResult, LocalStorageSetOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, LocalStorageSetOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	script := fmt.Sprintf(`localStorage.setItem(%q, %q)`, input.Key, input.Value)
	_, err = pilot.Evaluate(ctx, script)
	if err != nil {
		return nil, LocalStorageSetOutput{}, fmt.Errorf("localStorage.setItem failed: %w", err)
	}

	return nil, LocalStorageSetOutput{Message: fmt.Sprintf("Set localStorage[%q]", input.Key)}, nil
}

// LocalStorageDeleteInput for removing a localStorage item.
type LocalStorageDeleteInput struct {
	Key string `json:"key" jsonschema:"The key to remove from localStorage,required"`
}

// LocalStorageDeleteOutput confirms the operation.
type LocalStorageDeleteOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleLocalStorageDelete(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input LocalStorageDeleteInput,
) (*mcp.CallToolResult, LocalStorageDeleteOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, LocalStorageDeleteOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	script := fmt.Sprintf(`localStorage.removeItem(%q)`, input.Key)
	_, err = pilot.Evaluate(ctx, script)
	if err != nil {
		return nil, LocalStorageDeleteOutput{}, fmt.Errorf("localStorage.removeItem failed: %w", err)
	}

	return nil, LocalStorageDeleteOutput{Message: fmt.Sprintf("Deleted localStorage[%q]", input.Key)}, nil
}

// LocalStorageClearInput for clearing all localStorage.
type LocalStorageClearInput struct{}

// LocalStorageClearOutput confirms the operation.
type LocalStorageClearOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleLocalStorageClear(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input LocalStorageClearInput,
) (*mcp.CallToolResult, LocalStorageClearOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, LocalStorageClearOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	_, err = pilot.Evaluate(ctx, "localStorage.clear()")
	if err != nil {
		return nil, LocalStorageClearOutput{}, fmt.Errorf("localStorage.clear failed: %w", err)
	}

	return nil, LocalStorageClearOutput{Message: "localStorage cleared"}, nil
}

// LocalStorageListInput for listing all localStorage items.
type LocalStorageListInput struct{}

// LocalStorageListOutput contains all localStorage items.
type LocalStorageListOutput struct {
	Items map[string]string `json:"items"`
	Count int               `json:"count"`
}

func (s *Server) handleLocalStorageList(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input LocalStorageListInput,
) (*mcp.CallToolResult, LocalStorageListOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, LocalStorageListOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	script := `
		return (function() {
			const items = {};
			for (let i = 0; i < localStorage.length; i++) {
				const key = localStorage.key(i);
				items[key] = localStorage.getItem(key);
			}
			return JSON.stringify(items);
		})()
	`
	result, err := pilot.Evaluate(ctx, script)
	if err != nil {
		return nil, LocalStorageListOutput{}, fmt.Errorf("list localStorage failed: %w", err)
	}

	items := make(map[string]string)
	if result != nil {
		if str, ok := result.(string); ok {
			if err := json.Unmarshal([]byte(str), &items); err != nil {
				return nil, LocalStorageListOutput{}, fmt.Errorf("parse localStorage failed: %w", err)
			}
		}
	}

	return nil, LocalStorageListOutput{Items: items, Count: len(items)}, nil
}

// SessionStorage tools

// SessionStorageGetInput for getting a sessionStorage item.
type SessionStorageGetInput struct {
	Key string `json:"key" jsonschema:"The key to get from sessionStorage,required"`
}

// SessionStorageGetOutput contains the retrieved value.
type SessionStorageGetOutput struct {
	Key   string  `json:"key"`
	Value *string `json:"value"` // null if key doesn't exist
}

func (s *Server) handleSessionStorageGet(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SessionStorageGetInput,
) (*mcp.CallToolResult, SessionStorageGetOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, SessionStorageGetOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	script := fmt.Sprintf(`sessionStorage.getItem(%q)`, input.Key)
	result, err := pilot.Evaluate(ctx, "return "+script)
	if err != nil {
		return nil, SessionStorageGetOutput{}, fmt.Errorf("sessionStorage.getItem failed: %w", err)
	}

	output := SessionStorageGetOutput{Key: input.Key}
	if result != nil {
		if str, ok := result.(string); ok {
			output.Value = &str
		}
	}

	return nil, output, nil
}

// SessionStorageSetInput for setting a sessionStorage item.
type SessionStorageSetInput struct {
	Key   string `json:"key" jsonschema:"The key to set in sessionStorage,required"`
	Value string `json:"value" jsonschema:"The value to store,required"`
}

// SessionStorageSetOutput confirms the operation.
type SessionStorageSetOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleSessionStorageSet(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SessionStorageSetInput,
) (*mcp.CallToolResult, SessionStorageSetOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, SessionStorageSetOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	script := fmt.Sprintf(`sessionStorage.setItem(%q, %q)`, input.Key, input.Value)
	_, err = pilot.Evaluate(ctx, script)
	if err != nil {
		return nil, SessionStorageSetOutput{}, fmt.Errorf("sessionStorage.setItem failed: %w", err)
	}

	return nil, SessionStorageSetOutput{Message: fmt.Sprintf("Set sessionStorage[%q]", input.Key)}, nil
}

// SessionStorageDeleteInput for removing a sessionStorage item.
type SessionStorageDeleteInput struct {
	Key string `json:"key" jsonschema:"The key to remove from sessionStorage,required"`
}

// SessionStorageDeleteOutput confirms the operation.
type SessionStorageDeleteOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleSessionStorageDelete(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SessionStorageDeleteInput,
) (*mcp.CallToolResult, SessionStorageDeleteOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, SessionStorageDeleteOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	script := fmt.Sprintf(`sessionStorage.removeItem(%q)`, input.Key)
	_, err = pilot.Evaluate(ctx, script)
	if err != nil {
		return nil, SessionStorageDeleteOutput{}, fmt.Errorf("sessionStorage.removeItem failed: %w", err)
	}

	return nil, SessionStorageDeleteOutput{Message: fmt.Sprintf("Deleted sessionStorage[%q]", input.Key)}, nil
}

// SessionStorageClearInput for clearing all sessionStorage.
type SessionStorageClearInput struct{}

// SessionStorageClearOutput confirms the operation.
type SessionStorageClearOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleSessionStorageClear(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SessionStorageClearInput,
) (*mcp.CallToolResult, SessionStorageClearOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, SessionStorageClearOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	_, err = pilot.Evaluate(ctx, "sessionStorage.clear()")
	if err != nil {
		return nil, SessionStorageClearOutput{}, fmt.Errorf("sessionStorage.clear failed: %w", err)
	}

	return nil, SessionStorageClearOutput{Message: "sessionStorage cleared"}, nil
}

// SessionStorageListInput for listing all sessionStorage items.
type SessionStorageListInput struct{}

// SessionStorageListOutput contains all sessionStorage items.
type SessionStorageListOutput struct {
	Items map[string]string `json:"items"`
	Count int               `json:"count"`
}

func (s *Server) handleSessionStorageList(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SessionStorageListInput,
) (*mcp.CallToolResult, SessionStorageListOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, SessionStorageListOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	script := `
		return (function() {
			const items = {};
			for (let i = 0; i < sessionStorage.length; i++) {
				const key = sessionStorage.key(i);
				items[key] = sessionStorage.getItem(key);
			}
			return JSON.stringify(items);
		})()
	`
	result, err := pilot.Evaluate(ctx, script)
	if err != nil {
		return nil, SessionStorageListOutput{}, fmt.Errorf("list sessionStorage failed: %w", err)
	}

	items := make(map[string]string)
	if result != nil {
		if str, ok := result.(string); ok {
			if err := json.Unmarshal([]byte(str), &items); err != nil {
				return nil, SessionStorageListOutput{}, fmt.Errorf("parse sessionStorage failed: %w", err)
			}
		}
	}

	return nil, SessionStorageListOutput{Items: items, Count: len(items)}, nil
}
