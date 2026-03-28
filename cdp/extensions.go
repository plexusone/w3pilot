package cdp

import (
	"context"
	"encoding/json"
	"fmt"
)

// Extensions domain methods.
const (
	ExtensionsLoadUnpacked  = "Extensions.loadUnpacked"
	ExtensionsUninstall     = "Extensions.uninstall"
	ExtensionsGetAll        = "Extensions.getAll"
)

// ExtensionInfo contains information about a browser extension.
type ExtensionInfo struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Version     string `json:"version,omitempty"`
	Description string `json:"description,omitempty"`
	Enabled     bool   `json:"enabled"`
	Path        string `json:"path,omitempty"`
}

// LoadUnpackedExtension loads an unpacked extension from a directory.
// Returns the extension ID if successful.
func (c *Client) LoadUnpackedExtension(ctx context.Context, path string) (string, error) {
	result, err := c.Send(ctx, ExtensionsLoadUnpacked, map[string]interface{}{
		"path": path,
	})
	if err != nil {
		return "", fmt.Errorf("cdp: failed to load extension: %w", err)
	}

	var resp struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", fmt.Errorf("cdp: failed to parse extension response: %w", err)
	}

	return resp.ID, nil
}

// UninstallExtension removes an extension by ID.
func (c *Client) UninstallExtension(ctx context.Context, id string) error {
	_, err := c.Send(ctx, ExtensionsUninstall, map[string]interface{}{
		"id": id,
	})
	if err != nil {
		return fmt.Errorf("cdp: failed to uninstall extension: %w", err)
	}
	return nil
}

// GetAllExtensions returns all installed extensions.
func (c *Client) GetAllExtensions(ctx context.Context) ([]ExtensionInfo, error) {
	result, err := c.Send(ctx, ExtensionsGetAll, nil)
	if err != nil {
		return nil, fmt.Errorf("cdp: failed to get extensions: %w", err)
	}

	var resp struct {
		Extensions []ExtensionInfo `json:"extensions"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("cdp: failed to parse extensions: %w", err)
	}

	return resp.Extensions, nil
}
