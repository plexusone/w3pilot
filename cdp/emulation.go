package cdp

import (
	"context"
	"fmt"
)

// CPU throttling rate presets.
const (
	CPUNoThrottle   = 1 // No throttling
	CPU2xSlowdown   = 2 // 2x slowdown
	CPU4xSlowdown   = 4 // 4x slowdown (mid-tier mobile)
	CPU6xSlowdown   = 6 // 6x slowdown (low-end mobile)
)

// SetCPUThrottlingRate sets CPU throttling.
// rate=1 means no throttling, rate=4 means 4x slowdown.
func (c *Client) SetCPUThrottlingRate(ctx context.Context, rate int) error {
	if rate < 1 {
		rate = 1
	}

	_, err := c.Send(ctx, EmulationSetCPUThrottlingRate, map[string]interface{}{
		"rate": rate,
	})
	if err != nil {
		return fmt.Errorf("cdp: failed to set CPU throttling: %w", err)
	}
	return nil
}

// ClearCPUThrottling clears CPU throttling (returns to normal).
func (c *Client) ClearCPUThrottling(ctx context.Context) error {
	return c.SetCPUThrottlingRate(ctx, CPUNoThrottle)
}
