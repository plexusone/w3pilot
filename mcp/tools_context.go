package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	vibium "github.com/plexusone/vibium-go"
)

// NewPage tool

type NewPageInput struct{}

type NewPageOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleNewPage(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input NewPageInput,
) (*mcp.CallToolResult, NewPageOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, NewPageOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	_, err = vibe.NewPage(ctx)
	if err != nil {
		return nil, NewPageOutput{}, fmt.Errorf("new page failed: %w", err)
	}

	return nil, NewPageOutput{Message: "New page created"}, nil
}

// GetPages tool

type GetPagesInput struct{}

type GetPagesOutput struct {
	Count int `json:"count"`
}

func (s *Server) handleGetPages(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetPagesInput,
) (*mcp.CallToolResult, GetPagesOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, GetPagesOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	pages, err := vibe.Pages(ctx)
	if err != nil {
		return nil, GetPagesOutput{}, fmt.Errorf("get pages failed: %w", err)
	}

	return nil, GetPagesOutput{Count: len(pages)}, nil
}

// GetCookies tool

type GetCookiesInput struct {
	URLs []string `json:"urls" jsonschema:"description=URLs to get cookies for (optional)"`
}

type GetCookiesOutput struct {
	Cookies []CookieOutput `json:"cookies"`
}

type CookieOutput struct {
	Name     string  `json:"name"`
	Value    string  `json:"value"`
	Domain   string  `json:"domain"`
	Path     string  `json:"path"`
	Expires  float64 `json:"expires"`
	HTTPOnly bool    `json:"httpOnly"`
	Secure   bool    `json:"secure"`
	SameSite string  `json:"sameSite"`
}

func (s *Server) handleGetCookies(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetCookiesInput,
) (*mcp.CallToolResult, GetCookiesOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, GetCookiesOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	browserCtx, err := vibe.NewContext(ctx)
	if err != nil {
		return nil, GetCookiesOutput{}, fmt.Errorf("context not available: %w", err)
	}

	cookies, err := browserCtx.Cookies(ctx, input.URLs...)
	if err != nil {
		return nil, GetCookiesOutput{}, fmt.Errorf("get cookies failed: %w", err)
	}

	output := make([]CookieOutput, len(cookies))
	for i, c := range cookies {
		output[i] = CookieOutput{
			Name:     c.Name,
			Value:    c.Value,
			Domain:   c.Domain,
			Path:     c.Path,
			Expires:  c.Expires,
			HTTPOnly: c.HTTPOnly,
			Secure:   c.Secure,
			SameSite: c.SameSite,
		}
	}

	return nil, GetCookiesOutput{Cookies: output}, nil
}

// SetCookies tool

type SetCookiesInput struct {
	Cookies []SetCookieInput `json:"cookies" jsonschema:"description=Cookies to set,required"`
}

type SetCookieInput struct {
	Name     string  `json:"name"`
	Value    string  `json:"value"`
	URL      string  `json:"url,omitempty"`
	Domain   string  `json:"domain,omitempty"`
	Path     string  `json:"path,omitempty"`
	Expires  float64 `json:"expires,omitempty"`
	HTTPOnly bool    `json:"httpOnly,omitempty"`
	Secure   bool    `json:"secure,omitempty"`
	SameSite string  `json:"sameSite,omitempty"`
}

type SetCookiesOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleSetCookies(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SetCookiesInput,
) (*mcp.CallToolResult, SetCookiesOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, SetCookiesOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	browserCtx, err := vibe.NewContext(ctx)
	if err != nil {
		return nil, SetCookiesOutput{}, fmt.Errorf("context not available: %w", err)
	}

	cookies := make([]vibium.SetCookieParam, len(input.Cookies))
	for i, c := range input.Cookies {
		cookies[i] = vibium.SetCookieParam{
			Name:     c.Name,
			Value:    c.Value,
			URL:      c.URL,
			Domain:   c.Domain,
			Path:     c.Path,
			Expires:  c.Expires,
			HTTPOnly: c.HTTPOnly,
			Secure:   c.Secure,
			SameSite: c.SameSite,
		}
	}

	err = browserCtx.SetCookies(ctx, cookies)
	if err != nil {
		return nil, SetCookiesOutput{}, fmt.Errorf("set cookies failed: %w", err)
	}

	return nil, SetCookiesOutput{Message: fmt.Sprintf("Set %d cookies", len(input.Cookies))}, nil
}

// ClearCookies tool

type ClearCookiesInput struct{}

type ClearCookiesOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleClearCookies(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ClearCookiesInput,
) (*mcp.CallToolResult, ClearCookiesOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, ClearCookiesOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	browserCtx, err := vibe.NewContext(ctx)
	if err != nil {
		return nil, ClearCookiesOutput{}, fmt.Errorf("context not available: %w", err)
	}

	err = browserCtx.ClearCookies(ctx)
	if err != nil {
		return nil, ClearCookiesOutput{}, fmt.Errorf("clear cookies failed: %w", err)
	}

	return nil, ClearCookiesOutput{Message: "Cookies cleared"}, nil
}

// GetStorageState tool

type GetStorageStateInput struct{}

type GetStorageStateOutput struct {
	State string `json:"state"`
}

func (s *Server) handleGetStorageState(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetStorageStateInput,
) (*mcp.CallToolResult, GetStorageStateOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, GetStorageStateOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	browserCtx, err := vibe.NewContext(ctx)
	if err != nil {
		return nil, GetStorageStateOutput{}, fmt.Errorf("context not available: %w", err)
	}

	state, err := browserCtx.StorageState(ctx)
	if err != nil {
		return nil, GetStorageStateOutput{}, fmt.Errorf("get storage state failed: %w", err)
	}

	jsonBytes, err := json.Marshal(state)
	if err != nil {
		return nil, GetStorageStateOutput{}, fmt.Errorf("json marshal failed: %w", err)
	}

	return nil, GetStorageStateOutput{State: string(jsonBytes)}, nil
}

// SetStorageState tool

type SetStorageStateInput struct {
	State string `json:"state" jsonschema:"description=JSON from get_storage_state containing cookies and localStorage,required"`
}

type SetStorageStateOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleSetStorageState(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input SetStorageStateInput,
) (*mcp.CallToolResult, SetStorageStateOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, SetStorageStateOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	browserCtx, err := vibe.NewContext(ctx)
	if err != nil {
		return nil, SetStorageStateOutput{}, fmt.Errorf("context not available: %w", err)
	}

	// Parse the storage state JSON
	var state vibium.StorageState
	if err := json.Unmarshal([]byte(input.State), &state); err != nil {
		return nil, SetStorageStateOutput{}, fmt.Errorf("invalid storage state JSON: %w", err)
	}

	// Set cookies
	if len(state.Cookies) > 0 {
		cookies := make([]vibium.SetCookieParam, len(state.Cookies))
		for i, c := range state.Cookies {
			cookies[i] = vibium.SetCookieParam{
				Name:     c.Name,
				Value:    c.Value,
				Domain:   c.Domain,
				Path:     c.Path,
				Expires:  c.Expires,
				HTTPOnly: c.HTTPOnly,
				Secure:   c.Secure,
				SameSite: c.SameSite,
			}
		}
		if err := browserCtx.SetCookies(ctx, cookies); err != nil {
			return nil, SetStorageStateOutput{}, fmt.Errorf("set cookies failed: %w", err)
		}
	}

	// Set localStorage for each origin
	for _, origin := range state.Origins {
		if len(origin.LocalStorage) == 0 {
			continue
		}

		// Build JavaScript to set localStorage items
		localStorageJSON, err := json.Marshal(origin.LocalStorage)
		if err != nil {
			return nil, SetStorageStateOutput{}, fmt.Errorf("marshal localStorage failed: %w", err)
		}

		// Navigate to origin first (required to set localStorage for that origin)
		// Then set localStorage items
		script := fmt.Sprintf(`
			(function() {
				const items = %s;
				for (const [key, value] of Object.entries(items)) {
					localStorage.setItem(key, value);
				}
				return Object.keys(items).length;
			})()
		`, string(localStorageJSON))

		// We need to be on the correct origin to set localStorage
		// First check current URL
		currentURL, _ := vibe.URL(ctx)
		if currentURL == "" || currentURL == "about:blank" {
			// Navigate to the origin to set localStorage
			if err := vibe.Go(ctx, origin.Origin); err != nil {
				return nil, SetStorageStateOutput{}, fmt.Errorf("navigate to origin %s failed: %w", origin.Origin, err)
			}
		}

		if _, err := vibe.Evaluate(ctx, script); err != nil {
			return nil, SetStorageStateOutput{}, fmt.Errorf("set localStorage for %s failed: %w", origin.Origin, err)
		}
	}

	cookieCount := len(state.Cookies)
	originCount := len(state.Origins)
	return nil, SetStorageStateOutput{
		Message: fmt.Sprintf("Restored %d cookies and localStorage for %d origins", cookieCount, originCount),
	}, nil
}

// PauseForHuman tool

type PauseForHumanInput struct {
	Message   string `json:"message" jsonschema:"description=Message to display to the human (e.g. 'Please complete the login')"`
	TimeoutMS int    `json:"timeout_ms" jsonschema:"description=Maximum time to wait in milliseconds (default: 300000 = 5 minutes)"`
}

type PauseForHumanOutput struct {
	Message string `json:"message"`
}

func (s *Server) handlePauseForHuman(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input PauseForHumanInput,
) (*mcp.CallToolResult, PauseForHumanOutput, error) {
	vibe, err := s.session.Vibe(ctx)
	if err != nil {
		return nil, PauseForHumanOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if input.TimeoutMS == 0 {
		input.TimeoutMS = 300000 // 5 minutes default
	}

	if input.Message == "" {
		input.Message = "Please complete the required action, then click Continue."
	}

	// Inject an overlay with a button that the human clicks when done
	// The overlay is styled to be visible and unobtrusive
	overlayScript := fmt.Sprintf(`
		(function() {
			// Create overlay container
			const overlay = document.createElement('div');
			overlay.id = '__vibium_pause_overlay__';
			overlay.style.cssText = 'position:fixed;top:0;left:0;right:0;z-index:2147483647;background:linear-gradient(135deg,#667eea 0%%,#764ba2 100%%);padding:16px;display:flex;align-items:center;justify-content:center;gap:16px;font-family:-apple-system,BlinkMacSystemFont,sans-serif;box-shadow:0 4px 12px rgba(0,0,0,0.15);';

			// Create message
			const msg = document.createElement('span');
			msg.style.cssText = 'color:white;font-size:14px;font-weight:500;';
			msg.textContent = %q;

			// Create button
			const btn = document.createElement('button');
			btn.id = '__vibium_continue_btn__';
			btn.textContent = 'Continue';
			btn.style.cssText = 'background:white;color:#667eea;border:none;padding:8px 24px;border-radius:6px;font-size:14px;font-weight:600;cursor:pointer;transition:transform 0.1s;';
			btn.onmouseover = function() { this.style.transform = 'scale(1.05)'; };
			btn.onmouseout = function() { this.style.transform = 'scale(1)'; };
			btn.onclick = function() {
				overlay.remove();
				window.__vibium_human_done__ = true;
			};

			overlay.appendChild(msg);
			overlay.appendChild(btn);
			document.body.appendChild(overlay);

			window.__vibium_human_done__ = false;
			return true;
		})()
	`, input.Message)

	// Inject the overlay
	if _, err := vibe.Evaluate(ctx, overlayScript); err != nil {
		return nil, PauseForHumanOutput{}, fmt.Errorf("inject overlay failed: %w", err)
	}

	// Wait for the human to click the Continue button
	timeout := time.Duration(input.TimeoutMS) * time.Millisecond
	checkScript := `window.__vibium_human_done__ === true`

	if err := vibe.WaitForFunction(ctx, checkScript, timeout); err != nil {
		// Clean up overlay on timeout (best-effort, ignore error since we're already returning one)
		_, _ = vibe.Evaluate(ctx, `
			const overlay = document.getElementById('__vibium_pause_overlay__');
			if (overlay) overlay.remove();
		`)
		return nil, PauseForHumanOutput{}, fmt.Errorf("timeout waiting for human action: %w", err)
	}

	return nil, PauseForHumanOutput{Message: "Human action completed"}, nil
}
