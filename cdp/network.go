package cdp

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// ResponseBody represents the body of a network response.
type ResponseBody struct {
	Body          string // Raw body (or base64 encoded for binary)
	Base64Encoded bool   // True if body is base64 encoded
	Size          int    // Size in bytes
	Path          string // File path if saved to disk
}

// GetResponseBody retrieves the body of a network response.
// If saveTo is provided, the body is saved to that path.
func (c *Client) GetResponseBody(ctx context.Context, requestID string, saveTo string) (*ResponseBody, error) {
	result, err := c.Send(ctx, NetworkGetResponseBody, map[string]interface{}{
		"requestId": requestID,
	})
	if err != nil {
		return nil, fmt.Errorf("cdp: failed to get response body: %w", err)
	}

	var resp GetResponseBodyResult
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("cdp: failed to parse response body: %w", err)
	}

	body := &ResponseBody{
		Body:          resp.Body,
		Base64Encoded: resp.Base64Encoded,
		Size:          len(resp.Body),
	}

	// Save to file if requested
	if saveTo != "" {
		var data []byte
		if resp.Base64Encoded {
			var err error
			data, err = base64.StdEncoding.DecodeString(resp.Body)
			if err != nil {
				return nil, fmt.Errorf("cdp: failed to decode base64 body: %w", err)
			}
		} else {
			data = []byte(resp.Body)
		}

		// Ensure directory exists
		if err := os.MkdirAll(filepath.Dir(saveTo), 0755); err != nil {
			return nil, fmt.Errorf("cdp: failed to create directory: %w", err)
		}

		if err := os.WriteFile(saveTo, data, 0644); err != nil {
			return nil, fmt.Errorf("cdp: failed to write response body: %w", err)
		}

		body.Path = saveTo
		body.Size = len(data)
	}

	return body, nil
}

// EnableNetwork enables the Network domain with response body capture.
func (c *Client) EnableNetwork(ctx context.Context) error {
	_, err := c.Send(ctx, NetworkEnable, map[string]interface{}{
		"maxResourceBufferSize": 10 * 1024 * 1024,  // 10MB per resource
		"maxTotalBufferSize":    50 * 1024 * 1024,  // 50MB total
	})
	if err != nil {
		return fmt.Errorf("cdp: failed to enable Network: %w", err)
	}
	return nil
}

// DisableNetwork disables the Network domain.
func (c *Client) DisableNetwork(ctx context.Context) error {
	_, err := c.Send(ctx, NetworkDisable, nil)
	return err
}

// SetNetworkConditions sets network emulation conditions.
func (c *Client) SetNetworkConditions(ctx context.Context, conditions NetworkConditions) error {
	_, err := c.Send(ctx, NetworkEmulateConditions, conditions)
	if err != nil {
		return fmt.Errorf("cdp: failed to set network conditions: %w", err)
	}
	return nil
}

// ClearNetworkConditions clears network emulation (returns to normal).
func (c *Client) ClearNetworkConditions(ctx context.Context) error {
	return c.SetNetworkConditions(ctx, NetworkConditions{
		Offline:            false,
		Latency:            0,
		DownloadThroughput: -1, // -1 disables throttling
		UploadThroughput:   -1,
	})
}
