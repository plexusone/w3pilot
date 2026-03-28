package cdp

import (
	"context"
	"encoding/json"
	"fmt"
)

// CSS coverage domain methods.
const (
	CSSStartRuleUsageTracking = "CSS.startRuleUsageTracking"
	CSSStopRuleUsageTracking  = "CSS.stopRuleUsageTracking"
	CSSTakeCoverageDelta      = "CSS.takeCoverageDelta"
)

// CoverageRange represents a range of code that was covered.
type CoverageRange struct {
	StartOffset int `json:"startOffset"`
	EndOffset   int `json:"endOffset"`
	Count       int `json:"count"`
}

// FunctionCoverage represents coverage for a function.
type FunctionCoverage struct {
	FunctionName    string          `json:"functionName"`
	Ranges          []CoverageRange `json:"ranges"`
	IsBlockCoverage bool            `json:"isBlockCoverage"`
}

// ScriptCoverage represents coverage for a script.
type ScriptCoverage struct {
	ScriptID  string             `json:"scriptId"`
	URL       string             `json:"url"`
	Functions []FunctionCoverage `json:"functions"`
}

// CSSRuleUsage represents usage info for a CSS rule.
type CSSRuleUsage struct {
	StyleSheetID string  `json:"styleSheetId"`
	StartOffset  float64 `json:"startOffset"`
	EndOffset    float64 `json:"endOffset"`
	Used         bool    `json:"used"`
}

// CoverageReport contains JS and CSS coverage data.
type CoverageReport struct {
	JS  []ScriptCoverage `json:"js"`
	CSS []CSSRuleUsage   `json:"css"`
}

// CoverageSummary provides a high-level summary of coverage.
type CoverageSummary struct {
	JSScripts       int     `json:"jsScripts"`
	JSFunctions     int     `json:"jsFunctions"`
	JSCoveredRanges int     `json:"jsCoveredRanges"`
	CSSRules        int     `json:"cssRules"`
	CSSUsedRules    int     `json:"cssUsedRules"`
	CSSUnusedRules  int     `json:"cssUnusedRules"`
	CSSUsagePercent float64 `json:"cssUsagePercent"`
}

// Coverage manages code coverage collection.
type Coverage struct {
	client     *Client
	jsEnabled  bool
	cssEnabled bool
}

// NewCoverage creates a new coverage manager.
func NewCoverage(client *Client) *Coverage {
	return &Coverage{
		client: client,
	}
}

// StartJS enables JavaScript coverage collection.
// callCount: collect execution counts per block
// detailed: collect block-level coverage (vs function-level)
func (c *Coverage) StartJS(ctx context.Context, callCount, detailed bool) error {
	// Enable profiler domain
	if _, err := c.client.Send(ctx, ProfilerEnable, nil); err != nil {
		return fmt.Errorf("cdp: failed to enable Profiler: %w", err)
	}

	// Start precise coverage
	params := map[string]interface{}{
		"callCount": callCount,
		"detailed":  detailed,
	}
	if _, err := c.client.Send(ctx, ProfilerStartPreciseCoverage, params); err != nil {
		return fmt.Errorf("cdp: failed to start JS coverage: %w", err)
	}

	c.jsEnabled = true
	return nil
}

// StartCSS enables CSS coverage collection.
func (c *Coverage) StartCSS(ctx context.Context) error {
	if _, err := c.client.Send(ctx, CSSStartRuleUsageTracking, nil); err != nil {
		return fmt.Errorf("cdp: failed to start CSS coverage: %w", err)
	}

	c.cssEnabled = true
	return nil
}

// Start enables both JS and CSS coverage collection.
func (c *Coverage) Start(ctx context.Context) error {
	if err := c.StartJS(ctx, true, true); err != nil {
		return err
	}
	return c.StartCSS(ctx)
}

// StopJS stops JavaScript coverage and returns the results.
func (c *Coverage) StopJS(ctx context.Context) ([]ScriptCoverage, error) {
	if !c.jsEnabled {
		return nil, nil
	}

	// Take coverage snapshot
	result, err := c.client.Send(ctx, ProfilerTakePreciseCoverage, nil)
	if err != nil {
		return nil, fmt.Errorf("cdp: failed to take JS coverage: %w", err)
	}

	var resp struct {
		Result []ScriptCoverage `json:"result"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("cdp: failed to parse JS coverage: %w", err)
	}

	// Stop coverage
	if _, err := c.client.Send(ctx, ProfilerStopPreciseCoverage, nil); err != nil {
		// Log but don't fail
	}

	// Disable profiler
	if _, err := c.client.Send(ctx, ProfilerDisable, nil); err != nil {
		// Log but don't fail
	}

	c.jsEnabled = false
	return resp.Result, nil
}

// StopCSS stops CSS coverage and returns the results.
func (c *Coverage) StopCSS(ctx context.Context) ([]CSSRuleUsage, error) {
	if !c.cssEnabled {
		return nil, nil
	}

	// Take coverage delta
	result, err := c.client.Send(ctx, CSSTakeCoverageDelta, nil)
	if err != nil {
		return nil, fmt.Errorf("cdp: failed to take CSS coverage: %w", err)
	}

	var resp struct {
		Coverage []CSSRuleUsage `json:"coverage"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("cdp: failed to parse CSS coverage: %w", err)
	}

	// Stop tracking
	if _, err := c.client.Send(ctx, CSSStopRuleUsageTracking, nil); err != nil {
		// Log but don't fail
	}

	c.cssEnabled = false
	return resp.Coverage, nil
}

// Stop stops all coverage collection and returns the results.
func (c *Coverage) Stop(ctx context.Context) (*CoverageReport, error) {
	js, jsErr := c.StopJS(ctx)
	css, cssErr := c.StopCSS(ctx)

	if jsErr != nil {
		return nil, jsErr
	}
	if cssErr != nil {
		return nil, cssErr
	}

	return &CoverageReport{
		JS:  js,
		CSS: css,
	}, nil
}

// Summary returns a high-level summary of coverage data.
func (r *CoverageReport) Summary() CoverageSummary {
	summary := CoverageSummary{
		JSScripts: len(r.JS),
	}

	for _, script := range r.JS {
		summary.JSFunctions += len(script.Functions)
		for _, fn := range script.Functions {
			for _, rng := range fn.Ranges {
				if rng.Count > 0 {
					summary.JSCoveredRanges++
				}
			}
		}
	}

	summary.CSSRules = len(r.CSS)
	for _, rule := range r.CSS {
		if rule.Used {
			summary.CSSUsedRules++
		} else {
			summary.CSSUnusedRules++
		}
	}

	if summary.CSSRules > 0 {
		summary.CSSUsagePercent = float64(summary.CSSUsedRules) / float64(summary.CSSRules) * 100
	}

	return summary
}

// IsRunning returns whether coverage collection is active.
func (c *Coverage) IsRunning() bool {
	return c.jsEnabled || c.cssEnabled
}
