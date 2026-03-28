package mcp

import (
	"context"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	w3pilot "github.com/plexusone/w3pilot"
	"github.com/plexusone/w3pilot/cdp"
	"github.com/plexusone/w3pilot/mcp/report"
)

// GetPerformanceMetrics tool - get Core Web Vitals and navigation timing

type GetPerformanceMetricsInput struct{}

type GetPerformanceMetricsOutput struct {
	// Core Web Vitals
	LCP  float64 `json:"lcp,omitempty"`  // Largest Contentful Paint (ms)
	CLS  float64 `json:"cls,omitempty"`  // Cumulative Layout Shift
	INP  float64 `json:"inp,omitempty"`  // Interaction to Next Paint (ms)
	FID  float64 `json:"fid,omitempty"`  // First Input Delay (ms)
	FCP  float64 `json:"fcp,omitempty"`  // First Contentful Paint (ms)
	TTFB float64 `json:"ttfb,omitempty"` // Time to First Byte (ms)

	// Navigation Timing
	DOMContentLoaded float64 `json:"domContentLoaded,omitempty"` // DOMContentLoaded event (ms)
	Load             float64 `json:"load,omitempty"`             // Load event (ms)
	DOMInteractive   float64 `json:"domInteractive,omitempty"`   // DOM interactive (ms)

	// Resource Timing
	ResourceCount int `json:"resourceCount,omitempty"` // Number of resources loaded
}

func (s *Server) handleGetPerformanceMetrics(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetPerformanceMetricsInput,
) (*mcp.CallToolResult, GetPerformanceMetricsOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, GetPerformanceMetricsOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	start := time.Now()
	metrics, err := pilot.GetPerformanceMetrics(ctx)
	duration := time.Since(start)

	result := report.StepResult{
		ID:         s.session.NextStepID("get_performance_metrics"),
		Action:     "get_performance_metrics",
		DurationMS: duration.Milliseconds(),
	}

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityMedium
		result.Error = &report.StepError{
			Type:    "PerformanceMetricsError",
			Message: err.Error(),
		}
		s.session.RecordStep(result)
		return nil, GetPerformanceMetricsOutput{}, fmt.Errorf("failed to get performance metrics: %w", err)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	result.Result = map[string]any{
		"lcp":              metrics.LCP,
		"cls":              metrics.CLS,
		"fcp":              metrics.FCP,
		"ttfb":             metrics.TTFB,
		"domContentLoaded": metrics.DOMContentLoaded,
		"load":             metrics.Load,
	}
	s.session.RecordStep(result)

	return nil, GetPerformanceMetricsOutput{
		LCP:              metrics.LCP,
		CLS:              metrics.CLS,
		INP:              metrics.INP,
		FID:              metrics.FID,
		FCP:              metrics.FCP,
		TTFB:             metrics.TTFB,
		DOMContentLoaded: metrics.DOMContentLoaded,
		Load:             metrics.Load,
		DOMInteractive:   metrics.DOMInteractive,
		ResourceCount:    metrics.ResourceCount,
	}, nil
}

// GetMemoryStats tool - get JavaScript heap memory information

type GetMemoryStatsInput struct{}

type GetMemoryStatsOutput struct {
	UsedJSHeapSize  int64  `json:"usedJSHeapSize"`  // Used JS heap size in bytes
	TotalJSHeapSize int64  `json:"totalJSHeapSize"` // Total JS heap size in bytes
	JSHeapSizeLimit int64  `json:"jsHeapSizeLimit"` // JS heap size limit in bytes
	UsedMB          string `json:"usedMB"`          // Human-readable used size
	TotalMB         string `json:"totalMB"`         // Human-readable total size
	LimitMB         string `json:"limitMB"`         // Human-readable limit
	UsagePercent    string `json:"usagePercent"`    // Usage percentage
}

func (s *Server) handleGetMemoryStats(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetMemoryStatsInput,
) (*mcp.CallToolResult, GetMemoryStatsOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, GetMemoryStatsOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	start := time.Now()
	stats, err := pilot.GetMemoryStats(ctx)
	duration := time.Since(start)

	result := report.StepResult{
		ID:         s.session.NextStepID("get_memory_stats"),
		Action:     "get_memory_stats",
		DurationMS: duration.Milliseconds(),
	}

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityMedium
		result.Error = &report.StepError{
			Type:    "MemoryStatsError",
			Message: err.Error(),
		}
		s.session.RecordStep(result)
		return nil, GetMemoryStatsOutput{}, fmt.Errorf("failed to get memory stats: %w", err)
	}

	usagePercent := float64(0)
	if stats.JSHeapSizeLimit > 0 {
		usagePercent = float64(stats.UsedJSHeapSize) / float64(stats.JSHeapSizeLimit) * 100
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	result.Result = map[string]any{
		"usedJSHeapSize":  stats.UsedJSHeapSize,
		"totalJSHeapSize": stats.TotalJSHeapSize,
		"jsHeapSizeLimit": stats.JSHeapSizeLimit,
		"usagePercent":    usagePercent,
	}
	s.session.RecordStep(result)

	return nil, GetMemoryStatsOutput{
		UsedJSHeapSize:  stats.UsedJSHeapSize,
		TotalJSHeapSize: stats.TotalJSHeapSize,
		JSHeapSizeLimit: stats.JSHeapSizeLimit,
		UsedMB:          fmt.Sprintf("%.2f MB", float64(stats.UsedJSHeapSize)/(1024*1024)),
		TotalMB:         fmt.Sprintf("%.2f MB", float64(stats.TotalJSHeapSize)/(1024*1024)),
		LimitMB:         fmt.Sprintf("%.2f MB", float64(stats.JSHeapSizeLimit)/(1024*1024)),
		UsagePercent:    fmt.Sprintf("%.1f%%", usagePercent),
	}, nil
}

// TakeHeapSnapshot tool - capture V8 heap snapshot for memory profiling

type TakeHeapSnapshotInput struct {
	Path string `json:"path" jsonschema:"File path to save snapshot (optional, auto-generated if empty)"`
}

type TakeHeapSnapshotOutput struct {
	Path    string `json:"path"`
	Size    int64  `json:"size"`
	SizeMB  string `json:"sizeMB"`
	Message string `json:"message"`
}

func (s *Server) handleTakeHeapSnapshot(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input TakeHeapSnapshotInput,
) (*mcp.CallToolResult, TakeHeapSnapshotOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, TakeHeapSnapshotOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if !pilot.HasCDP() {
		return nil, TakeHeapSnapshotOutput{}, fmt.Errorf("CDP not available - heap snapshots require CDP connection")
	}

	start := time.Now()
	snapshot, err := pilot.TakeHeapSnapshot(ctx, input.Path)
	duration := time.Since(start)

	result := report.StepResult{
		ID:         s.session.NextStepID("take_heap_snapshot"),
		Action:     "take_heap_snapshot",
		Args:       map[string]any{"path": input.Path},
		DurationMS: duration.Milliseconds(),
	}

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityMedium
		result.Error = &report.StepError{
			Type:    "HeapSnapshotError",
			Message: err.Error(),
		}
		s.session.RecordStep(result)
		return nil, TakeHeapSnapshotOutput{}, fmt.Errorf("failed to take heap snapshot: %w", err)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	result.Result = map[string]any{
		"path": snapshot.Path,
		"size": snapshot.Size,
	}
	s.session.RecordStep(result)

	return nil, TakeHeapSnapshotOutput{
		Path:    snapshot.Path,
		Size:    snapshot.Size,
		SizeMB:  fmt.Sprintf("%.2f MB", float64(snapshot.Size)/(1024*1024)),
		Message: fmt.Sprintf("Heap snapshot saved to %s", snapshot.Path),
	}, nil
}

// EmulateNetwork tool - simulate network conditions

type EmulateNetworkInput struct {
	Preset             string  `json:"preset" jsonschema:"Network preset: slow3g, fast3g, 4g, wifi, offline,enum=slow3g,enum=fast3g,enum=4g,enum=wifi,enum=offline"`
	Latency            float64 `json:"latency" jsonschema:"Custom latency in ms (overrides preset)"`
	DownloadThroughput float64 `json:"download_throughput" jsonschema:"Custom download throughput in bytes/s (overrides preset)"`
	UploadThroughput   float64 `json:"upload_throughput" jsonschema:"Custom upload throughput in bytes/s (overrides preset)"`
}

type EmulateNetworkOutput struct {
	Message string `json:"message"`
	Preset  string `json:"preset,omitempty"`
}

func (s *Server) handleEmulateNetwork(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input EmulateNetworkInput,
) (*mcp.CallToolResult, EmulateNetworkOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, EmulateNetworkOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if !pilot.HasCDP() {
		return nil, EmulateNetworkOutput{}, fmt.Errorf("CDP not available - network emulation requires CDP connection")
	}

	var conditions cdp.NetworkConditions

	// Use preset or custom values
	switch input.Preset {
	case "slow3g":
		conditions = cdp.NetworkSlow3G
	case "fast3g":
		conditions = cdp.NetworkFast3G
	case "4g":
		conditions = cdp.Network4G
	case "wifi":
		conditions = cdp.NetworkWifi
	case "offline":
		conditions = cdp.NetworkOffline
	default:
		// Use custom values if provided, otherwise default to no throttling
		conditions = cdp.NetworkConditions{
			Latency:            input.Latency,
			DownloadThroughput: input.DownloadThroughput,
			UploadThroughput:   input.UploadThroughput,
		}
		if input.Latency == 0 && input.DownloadThroughput == 0 && input.UploadThroughput == 0 {
			return nil, EmulateNetworkOutput{}, fmt.Errorf("specify a preset (slow3g, fast3g, 4g, wifi, offline) or custom values")
		}
	}

	start := time.Now()
	err = pilot.EmulateNetwork(ctx, conditions)
	duration := time.Since(start)

	result := report.StepResult{
		ID:         s.session.NextStepID("emulate_network"),
		Action:     "emulate_network",
		Args:       map[string]any{"preset": input.Preset},
		DurationMS: duration.Milliseconds(),
	}

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityMedium
		result.Error = &report.StepError{
			Type:    "EmulateNetworkError",
			Message: err.Error(),
		}
		s.session.RecordStep(result)
		return nil, EmulateNetworkOutput{}, fmt.Errorf("failed to emulate network: %w", err)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	s.session.RecordStep(result)

	msg := "Network emulation enabled"
	if input.Preset != "" {
		msg = fmt.Sprintf("Network emulation set to %s", input.Preset)
	}

	return nil, EmulateNetworkOutput{
		Message: msg,
		Preset:  input.Preset,
	}, nil
}

// ClearNetworkEmulation tool - remove network throttling

type ClearNetworkEmulationInput struct{}

type ClearNetworkEmulationOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleClearNetworkEmulation(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ClearNetworkEmulationInput,
) (*mcp.CallToolResult, ClearNetworkEmulationOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, ClearNetworkEmulationOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if !pilot.HasCDP() {
		return nil, ClearNetworkEmulationOutput{}, fmt.Errorf("CDP not available")
	}

	if err := pilot.ClearNetworkEmulation(ctx); err != nil {
		return nil, ClearNetworkEmulationOutput{}, fmt.Errorf("failed to clear network emulation: %w", err)
	}

	return nil, ClearNetworkEmulationOutput{
		Message: "Network emulation cleared",
	}, nil
}

// EmulateCPU tool - simulate slower CPU

type EmulateCPUInput struct {
	Rate int `json:"rate" jsonschema:"CPU throttling rate: 1 (none), 2 (2x slower), 4 (4x slower, mid-tier mobile), 6 (6x slower, low-end mobile),required"`
}

type EmulateCPUOutput struct {
	Message string `json:"message"`
	Rate    int    `json:"rate"`
}

func (s *Server) handleEmulateCPU(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input EmulateCPUInput,
) (*mcp.CallToolResult, EmulateCPUOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, EmulateCPUOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if !pilot.HasCDP() {
		return nil, EmulateCPUOutput{}, fmt.Errorf("CDP not available - CPU emulation requires CDP connection")
	}

	if input.Rate < 1 {
		input.Rate = 1
	}

	start := time.Now()
	err = pilot.EmulateCPU(ctx, input.Rate)
	duration := time.Since(start)

	result := report.StepResult{
		ID:         s.session.NextStepID("emulate_cpu"),
		Action:     "emulate_cpu",
		Args:       map[string]any{"rate": input.Rate},
		DurationMS: duration.Milliseconds(),
	}

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityMedium
		result.Error = &report.StepError{
			Type:    "EmulateCPUError",
			Message: err.Error(),
		}
		s.session.RecordStep(result)
		return nil, EmulateCPUOutput{}, fmt.Errorf("failed to emulate CPU: %w", err)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	s.session.RecordStep(result)

	msg := "CPU throttling disabled"
	if input.Rate > 1 {
		msg = fmt.Sprintf("CPU throttled to %dx slowdown", input.Rate)
	}

	return nil, EmulateCPUOutput{
		Message: msg,
		Rate:    input.Rate,
	}, nil
}

// ClearCPUEmulation tool - remove CPU throttling

type ClearCPUEmulationInput struct{}

type ClearCPUEmulationOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleClearCPUEmulation(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ClearCPUEmulationInput,
) (*mcp.CallToolResult, ClearCPUEmulationOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, ClearCPUEmulationOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if !pilot.HasCDP() {
		return nil, ClearCPUEmulationOutput{}, fmt.Errorf("CDP not available")
	}

	if err := pilot.ClearCPUEmulation(ctx); err != nil {
		return nil, ClearCPUEmulationOutput{}, fmt.Errorf("failed to clear CPU emulation: %w", err)
	}

	return nil, ClearCPUEmulationOutput{
		Message: "CPU emulation cleared",
	}, nil
}

// LighthouseAudit tool - run Lighthouse quality audit

type LighthouseAuditInput struct {
	Categories []string `json:"categories" jsonschema:"Audit categories to run: accessibility, seo, best-practices, performance (default: all except performance)"`
	Device     string   `json:"device" jsonschema:"Device to emulate: desktop or mobile (default: desktop),enum=desktop,enum=mobile"`
	OutputDir  string   `json:"output_dir" jsonschema:"Directory to save reports (optional, uses temp dir if empty)"`
}

type LighthouseAuditOutput struct {
	URL             string                        `json:"url"`
	Device          string                        `json:"device"`
	Scores          map[string]LighthouseScoreOut `json:"scores"`
	PassedAudits    int                           `json:"passedAudits"`
	FailedAudits    int                           `json:"failedAudits"`
	TotalDurationMS float64                       `json:"totalDurationMs"`
	JSONReportPath  string                        `json:"jsonReportPath,omitempty"`
	HTMLReportPath  string                        `json:"htmlReportPath,omitempty"`
	Message         string                        `json:"message"`
}

type LighthouseScoreOut struct {
	Title    string  `json:"title"`
	Score    int     `json:"score"`    // 0-100
	RawScore float64 `json:"rawScore"` // 0-1
}

func (s *Server) handleLighthouseAudit(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input LighthouseAuditInput,
) (*mcp.CallToolResult, LighthouseAuditOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, LighthouseAuditOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	// Import w3pilot package types
	opts := &w3pilot.LighthouseOptions{
		OutputDir: input.OutputDir,
	}

	// Parse categories
	for _, cat := range input.Categories {
		switch cat {
		case "performance":
			opts.Categories = append(opts.Categories, w3pilot.LighthousePerformance)
		case "accessibility":
			opts.Categories = append(opts.Categories, w3pilot.LighthouseAccessibility)
		case "best-practices":
			opts.Categories = append(opts.Categories, w3pilot.LighthouseBestPractices)
		case "seo":
			opts.Categories = append(opts.Categories, w3pilot.LighthouseSEO)
		}
	}

	// Parse device
	switch input.Device {
	case "mobile":
		opts.Device = w3pilot.LighthouseMobile
	case "desktop", "":
		opts.Device = w3pilot.LighthouseDesktop
	}

	start := time.Now()
	result, err := pilot.LighthouseAudit(ctx, opts)
	duration := time.Since(start)

	stepResult := report.StepResult{
		ID:         s.session.NextStepID("lighthouse_audit"),
		Action:     "lighthouse_audit",
		Args:       map[string]any{"categories": input.Categories, "device": input.Device},
		DurationMS: duration.Milliseconds(),
	}

	if err != nil {
		stepResult.Status = report.StatusNoGo
		stepResult.Severity = report.SeverityMedium
		stepResult.Error = &report.StepError{
			Type:    "LighthouseError",
			Message: err.Error(),
		}
		s.session.RecordStep(stepResult)
		return nil, LighthouseAuditOutput{}, fmt.Errorf("lighthouse audit failed: %w", err)
	}

	// Convert scores
	scores := make(map[string]LighthouseScoreOut)
	for id, score := range result.Scores {
		scores[id] = LighthouseScoreOut{
			Title:    score.Title,
			Score:    int(score.Score * 100),
			RawScore: score.Score,
		}
	}

	stepResult.Status = report.StatusGo
	stepResult.Severity = report.SeverityInfo
	stepResult.Result = map[string]any{
		"url":          result.URL,
		"scores":       scores,
		"passedAudits": result.PassedAudits,
		"failedAudits": result.FailedAudits,
	}
	s.session.RecordStep(stepResult)

	// Build summary message
	var scoreStrs []string
	for id, score := range scores {
		scoreStrs = append(scoreStrs, fmt.Sprintf("%s: %d", id, score.Score))
	}

	return nil, LighthouseAuditOutput{
		URL:             result.URL,
		Device:          result.Device,
		Scores:          scores,
		PassedAudits:    result.PassedAudits,
		FailedAudits:    result.FailedAudits,
		TotalDurationMS: result.TotalDurationMS,
		JSONReportPath:  result.JSONReportPath,
		HTMLReportPath:  result.HTMLReportPath,
		Message:         fmt.Sprintf("Lighthouse audit complete. Scores: %v", scoreStrs),
	}, nil
}

// GetNetworkRequestBody tool - get response body for a network request

type GetNetworkRequestBodyInput struct {
	RequestID  string `json:"request_id" jsonschema:"The request ID from get_network_requests output,required"`
	SaveToFile string `json:"save_to_file" jsonschema:"File path to save binary content (optional, returns base64 if not set)"`
}

type GetNetworkRequestBodyOutput struct {
	Body          string `json:"body,omitempty"`        // Text body (if not binary or saved to file)
	Base64Encoded bool   `json:"base64Encoded"`         // Whether body is base64 encoded
	SavedToFile   string `json:"savedToFile,omitempty"` // File path if saved
	Size          int    `json:"size"`                  // Body size in bytes
	Message       string `json:"message"`
}

func (s *Server) handleGetNetworkRequestBody(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetNetworkRequestBodyInput,
) (*mcp.CallToolResult, GetNetworkRequestBodyOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, GetNetworkRequestBodyOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if !pilot.HasCDP() {
		return nil, GetNetworkRequestBodyOutput{}, fmt.Errorf("CDP not available - network request bodies require CDP connection")
	}

	if input.RequestID == "" {
		return nil, GetNetworkRequestBodyOutput{}, fmt.Errorf("request_id is required")
	}

	start := time.Now()
	body, err := pilot.GetNetworkResponseBody(ctx, input.RequestID, input.SaveToFile)
	duration := time.Since(start)

	result := report.StepResult{
		ID:         s.session.NextStepID("get_network_request_body"),
		Action:     "get_network_request_body",
		Args:       map[string]any{"request_id": input.RequestID},
		DurationMS: duration.Milliseconds(),
	}

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityMedium
		result.Error = &report.StepError{
			Type:    "NetworkRequestBodyError",
			Message: err.Error(),
		}
		s.session.RecordStep(result)
		return nil, GetNetworkRequestBodyOutput{}, fmt.Errorf("failed to get network request body: %w", err)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	result.Result = map[string]any{
		"size":          body.Size,
		"base64Encoded": body.Base64Encoded,
		"savedToFile":   body.Path,
	}
	s.session.RecordStep(result)

	msg := fmt.Sprintf("Retrieved response body (%d bytes)", body.Size)
	if body.Path != "" {
		msg = fmt.Sprintf("Saved response body to %s (%d bytes)", body.Path, body.Size)
	}

	return nil, GetNetworkRequestBodyOutput{
		Body:          body.Body,
		Base64Encoded: body.Base64Encoded,
		SavedToFile:   body.Path,
		Size:          body.Size,
		Message:       msg,
	}, nil
}

// StartScreencast tool - begin capturing screen frames

type StartScreencastInput struct {
	Format        string `json:"format" jsonschema:"Image format: jpeg or png (default: jpeg),enum=jpeg,enum=png"`
	Quality       int    `json:"quality" jsonschema:"Image quality 0-100 for jpeg (default: 80)"`
	MaxWidth      int    `json:"max_width" jsonschema:"Maximum frame width in pixels (optional)"`
	MaxHeight     int    `json:"max_height" jsonschema:"Maximum frame height in pixels (optional)"`
	EveryNthFrame int    `json:"every_nth_frame" jsonschema:"Capture every Nth frame, 1=every frame (default: 1)"`
}

type StartScreencastOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleStartScreencast(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input StartScreencastInput,
) (*mcp.CallToolResult, StartScreencastOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, StartScreencastOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if !pilot.HasCDP() {
		return nil, StartScreencastOutput{}, fmt.Errorf("CDP not available - screencast requires CDP connection")
	}

	// Build options
	opts := &cdp.ScreencastOptions{
		Quality: 80, // Default quality
	}

	if input.Format == "png" {
		opts.Format = cdp.ScreencastFormatPNG
	} else {
		opts.Format = cdp.ScreencastFormatJPEG
	}

	if input.Quality > 0 && input.Quality <= 100 {
		opts.Quality = input.Quality
	}
	if input.MaxWidth > 0 {
		opts.MaxWidth = input.MaxWidth
	}
	if input.MaxHeight > 0 {
		opts.MaxHeight = input.MaxHeight
	}
	if input.EveryNthFrame > 0 {
		opts.EveryNthFrame = input.EveryNthFrame
	}

	start := time.Now()

	// Store frames in session for retrieval (MCP tools can't use callbacks directly)
	// The frames are available via CDP events but we just start the capture
	err = pilot.StartScreencast(ctx, opts, nil)
	duration := time.Since(start)

	result := report.StepResult{
		ID:         s.session.NextStepID("start_screencast"),
		Action:     "start_screencast",
		Args:       map[string]any{"format": string(opts.Format), "quality": opts.Quality},
		DurationMS: duration.Milliseconds(),
	}

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityMedium
		result.Error = &report.StepError{
			Type:    "ScreencastError",
			Message: err.Error(),
		}
		s.session.RecordStep(result)
		return nil, StartScreencastOutput{}, fmt.Errorf("failed to start screencast: %w", err)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	s.session.RecordStep(result)

	return nil, StartScreencastOutput{
		Message: fmt.Sprintf("Screencast started (format: %s, quality: %d)", opts.Format, opts.Quality),
	}, nil
}

// StopScreencast tool - stop capturing screen frames

type StopScreencastInput struct{}

type StopScreencastOutput struct {
	Message string `json:"message"`
}

//nolint:dupl // Handler pattern intentionally similar to other CDP tool handlers
func (s *Server) handleStopScreencast(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input StopScreencastInput,
) (*mcp.CallToolResult, StopScreencastOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, StopScreencastOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if !pilot.HasCDP() {
		return nil, StopScreencastOutput{}, fmt.Errorf("CDP not available")
	}

	start := time.Now()
	err = pilot.StopScreencast(ctx)
	duration := time.Since(start)

	result := report.StepResult{
		ID:         s.session.NextStepID("stop_screencast"),
		Action:     "stop_screencast",
		DurationMS: duration.Milliseconds(),
	}

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityMedium
		result.Error = &report.StepError{
			Type:    "ScreencastError",
			Message: err.Error(),
		}
		s.session.RecordStep(result)
		return nil, StopScreencastOutput{}, fmt.Errorf("failed to stop screencast: %w", err)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	s.session.RecordStep(result)

	return nil, StopScreencastOutput{
		Message: "Screencast stopped",
	}, nil
}

// InstallExtension tool - load an unpacked extension

type InstallExtensionInput struct {
	Path string `json:"path" jsonschema:"Path to unpacked extension directory,required"`
}

type InstallExtensionOutput struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

func (s *Server) handleInstallExtension(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input InstallExtensionInput,
) (*mcp.CallToolResult, InstallExtensionOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, InstallExtensionOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if !pilot.HasCDP() {
		return nil, InstallExtensionOutput{}, fmt.Errorf("CDP not available - extension management requires CDP connection")
	}

	if input.Path == "" {
		return nil, InstallExtensionOutput{}, fmt.Errorf("path is required")
	}

	start := time.Now()
	id, err := pilot.InstallExtension(ctx, input.Path)
	duration := time.Since(start)

	result := report.StepResult{
		ID:         s.session.NextStepID("install_extension"),
		Action:     "install_extension",
		Args:       map[string]any{"path": input.Path},
		DurationMS: duration.Milliseconds(),
	}

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityMedium
		result.Error = &report.StepError{
			Type:    "ExtensionError",
			Message: err.Error(),
		}
		s.session.RecordStep(result)
		return nil, InstallExtensionOutput{}, fmt.Errorf("failed to install extension: %w", err)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	result.Result = map[string]any{"id": id}
	s.session.RecordStep(result)

	return nil, InstallExtensionOutput{
		ID:      id,
		Message: fmt.Sprintf("Extension installed with ID: %s", id),
	}, nil
}

// UninstallExtension tool - remove an extension by ID

type UninstallExtensionInput struct {
	ID string `json:"id" jsonschema:"Extension ID to uninstall,required"`
}

type UninstallExtensionOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleUninstallExtension(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input UninstallExtensionInput,
) (*mcp.CallToolResult, UninstallExtensionOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, UninstallExtensionOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if !pilot.HasCDP() {
		return nil, UninstallExtensionOutput{}, fmt.Errorf("CDP not available - extension management requires CDP connection")
	}

	if input.ID == "" {
		return nil, UninstallExtensionOutput{}, fmt.Errorf("id is required")
	}

	start := time.Now()
	err = pilot.UninstallExtension(ctx, input.ID)
	duration := time.Since(start)

	result := report.StepResult{
		ID:         s.session.NextStepID("uninstall_extension"),
		Action:     "uninstall_extension",
		Args:       map[string]any{"id": input.ID},
		DurationMS: duration.Milliseconds(),
	}

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityMedium
		result.Error = &report.StepError{
			Type:    "ExtensionError",
			Message: err.Error(),
		}
		s.session.RecordStep(result)
		return nil, UninstallExtensionOutput{}, fmt.Errorf("failed to uninstall extension: %w", err)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	s.session.RecordStep(result)

	return nil, UninstallExtensionOutput{
		Message: fmt.Sprintf("Extension %s uninstalled", input.ID),
	}, nil
}

// ListExtensions tool - get all installed extensions

type ListExtensionsInput struct{}

type ListExtensionsOutput struct {
	Extensions []ExtensionInfoOutput `json:"extensions"`
	Count      int                   `json:"count"`
}

type ExtensionInfoOutput struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Version     string `json:"version,omitempty"`
	Description string `json:"description,omitempty"`
	Enabled     bool   `json:"enabled"`
}

func (s *Server) handleListExtensions(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ListExtensionsInput,
) (*mcp.CallToolResult, ListExtensionsOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, ListExtensionsOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if !pilot.HasCDP() {
		return nil, ListExtensionsOutput{}, fmt.Errorf("CDP not available - extension management requires CDP connection")
	}

	start := time.Now()
	extensions, err := pilot.ListExtensions(ctx)
	duration := time.Since(start)

	result := report.StepResult{
		ID:         s.session.NextStepID("list_extensions"),
		Action:     "list_extensions",
		DurationMS: duration.Milliseconds(),
	}

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityMedium
		result.Error = &report.StepError{
			Type:    "ExtensionError",
			Message: err.Error(),
		}
		s.session.RecordStep(result)
		return nil, ListExtensionsOutput{}, fmt.Errorf("failed to list extensions: %w", err)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	result.Result = map[string]any{"count": len(extensions)}
	s.session.RecordStep(result)

	output := ListExtensionsOutput{
		Extensions: make([]ExtensionInfoOutput, len(extensions)),
		Count:      len(extensions),
	}

	for i, ext := range extensions {
		output.Extensions[i] = ExtensionInfoOutput{
			ID:          ext.ID,
			Name:        ext.Name,
			Version:     ext.Version,
			Description: ext.Description,
			Enabled:     ext.Enabled,
		}
	}

	return nil, output, nil
}

// StartCoverage tool - begin collecting JS and CSS coverage

type StartCoverageInput struct {
	JS        bool `json:"js" jsonschema:"Enable JavaScript coverage (default: true)"`
	CSS       bool `json:"css" jsonschema:"Enable CSS coverage (default: true)"`
	CallCount bool `json:"call_count" jsonschema:"Collect execution counts per block (default: true)"`
	Detailed  bool `json:"detailed" jsonschema:"Collect block-level coverage vs function-level (default: true)"`
}

type StartCoverageOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleStartCoverage(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input StartCoverageInput,
) (*mcp.CallToolResult, StartCoverageOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, StartCoverageOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if !pilot.HasCDP() {
		return nil, StartCoverageOutput{}, fmt.Errorf("CDP not available - coverage requires CDP connection")
	}

	// Default to both JS and CSS if neither specified
	enableJS := input.JS
	enableCSS := input.CSS
	if !enableJS && !enableCSS {
		enableJS = true
		enableCSS = true
	}

	// Default to detailed coverage
	callCount := input.CallCount
	detailed := input.Detailed
	if !callCount && !detailed {
		callCount = true
		detailed = true
	}

	start := time.Now()
	var coverageErr error

	if enableJS {
		coverageErr = pilot.StartJSCoverage(ctx, callCount, detailed)
	}
	if coverageErr == nil && enableCSS {
		coverageErr = pilot.StartCSSCoverage(ctx)
	}

	duration := time.Since(start)

	result := report.StepResult{
		ID:         s.session.NextStepID("start_coverage"),
		Action:     "start_coverage",
		Args:       map[string]any{"js": enableJS, "css": enableCSS},
		DurationMS: duration.Milliseconds(),
	}

	if coverageErr != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityMedium
		result.Error = &report.StepError{
			Type:    "CoverageError",
			Message: coverageErr.Error(),
		}
		s.session.RecordStep(result)
		return nil, StartCoverageOutput{}, fmt.Errorf("failed to start coverage: %w", coverageErr)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	s.session.RecordStep(result)

	types := []string{}
	if enableJS {
		types = append(types, "JS")
	}
	if enableCSS {
		types = append(types, "CSS")
	}

	return nil, StartCoverageOutput{
		Message: fmt.Sprintf("Coverage started for: %v", types),
	}, nil
}

// StopCoverage tool - stop collecting coverage and return results

type StopCoverageInput struct{}

type StopCoverageOutput struct {
	Summary   CoverageSummaryOutput `json:"summary"`
	JSScripts int                   `json:"jsScripts"`
	CSSRules  int                   `json:"cssRules"`
	Message   string                `json:"message"`
}

type CoverageSummaryOutput struct {
	JSScripts       int     `json:"jsScripts"`
	JSFunctions     int     `json:"jsFunctions"`
	JSCoveredRanges int     `json:"jsCoveredRanges"`
	CSSRules        int     `json:"cssRules"`
	CSSUsedRules    int     `json:"cssUsedRules"`
	CSSUnusedRules  int     `json:"cssUnusedRules"`
	CSSUsagePercent float64 `json:"cssUsagePercent"`
}

func (s *Server) handleStopCoverage(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input StopCoverageInput,
) (*mcp.CallToolResult, StopCoverageOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, StopCoverageOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if !pilot.HasCDP() {
		return nil, StopCoverageOutput{}, fmt.Errorf("CDP not available")
	}

	start := time.Now()
	coverageReport, err := pilot.StopCoverage(ctx)
	duration := time.Since(start)

	stepResult := report.StepResult{
		ID:         s.session.NextStepID("stop_coverage"),
		Action:     "stop_coverage",
		DurationMS: duration.Milliseconds(),
	}

	if err != nil {
		stepResult.Status = report.StatusNoGo
		stepResult.Severity = report.SeverityMedium
		stepResult.Error = &report.StepError{
			Type:    "CoverageError",
			Message: err.Error(),
		}
		s.session.RecordStep(stepResult)
		return nil, StopCoverageOutput{}, fmt.Errorf("failed to stop coverage: %w", err)
	}

	summary := coverageReport.Summary()

	stepResult.Status = report.StatusGo
	stepResult.Severity = report.SeverityInfo
	stepResult.Result = map[string]any{
		"jsScripts":       summary.JSScripts,
		"jsFunctions":     summary.JSFunctions,
		"cssRules":        summary.CSSRules,
		"cssUsagePercent": summary.CSSUsagePercent,
	}
	s.session.RecordStep(stepResult)

	return nil, StopCoverageOutput{
		Summary: CoverageSummaryOutput{
			JSScripts:       summary.JSScripts,
			JSFunctions:     summary.JSFunctions,
			JSCoveredRanges: summary.JSCoveredRanges,
			CSSRules:        summary.CSSRules,
			CSSUsedRules:    summary.CSSUsedRules,
			CSSUnusedRules:  summary.CSSUnusedRules,
			CSSUsagePercent: summary.CSSUsagePercent,
		},
		JSScripts: summary.JSScripts,
		CSSRules:  summary.CSSRules,
		Message:   fmt.Sprintf("Coverage stopped. JS: %d scripts, CSS: %.1f%% used (%d/%d rules)", summary.JSScripts, summary.CSSUsagePercent, summary.CSSUsedRules, summary.CSSRules),
	}, nil
}

// EnableConsoleDebugger tool - start capturing console messages with stack traces

type EnableConsoleDebuggerInput struct{}

type EnableConsoleDebuggerOutput struct {
	Message string `json:"message"`
}

//nolint:dupl // Handler pattern intentionally similar to other CDP tool handlers
func (s *Server) handleEnableConsoleDebugger(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input EnableConsoleDebuggerInput,
) (*mcp.CallToolResult, EnableConsoleDebuggerOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, EnableConsoleDebuggerOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if !pilot.HasCDP() {
		return nil, EnableConsoleDebuggerOutput{}, fmt.Errorf("CDP not available - console debugger requires CDP connection")
	}

	start := time.Now()
	err = pilot.EnableConsoleDebugger(ctx)
	duration := time.Since(start)

	result := report.StepResult{
		ID:         s.session.NextStepID("enable_console_debugger"),
		Action:     "enable_console_debugger",
		DurationMS: duration.Milliseconds(),
	}

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityMedium
		result.Error = &report.StepError{
			Type:    "ConsoleDebuggerError",
			Message: err.Error(),
		}
		s.session.RecordStep(result)
		return nil, EnableConsoleDebuggerOutput{}, fmt.Errorf("failed to enable console debugger: %w", err)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	s.session.RecordStep(result)

	return nil, EnableConsoleDebuggerOutput{
		Message: "Console debugger enabled - capturing messages with stack traces",
	}, nil
}

// GetConsoleEntriesWithStack tool - get console messages with full stack traces

type GetConsoleEntriesWithStackInput struct {
	Type  string `json:"type" jsonschema:"Filter by message type: log, debug, info, error, warning, trace (optional)"`
	Limit int    `json:"limit" jsonschema:"Maximum number of entries to return (default: all)"`
}

type GetConsoleEntriesWithStackOutput struct {
	Entries []ConsoleEntryOutput `json:"entries"`
	Count   int                  `json:"count"`
}

type ConsoleEntryOutput struct {
	Type       string            `json:"type"`
	Text       string            `json:"text"`
	Timestamp  float64           `json:"timestamp"`
	StackTrace []CallFrameOutput `json:"stackTrace,omitempty"`
}

type CallFrameOutput struct {
	FunctionName string `json:"functionName"`
	URL          string `json:"url"`
	LineNumber   int    `json:"lineNumber"`
	ColumnNumber int    `json:"columnNumber"`
}

func (s *Server) handleGetConsoleEntriesWithStack(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetConsoleEntriesWithStackInput,
) (*mcp.CallToolResult, GetConsoleEntriesWithStackOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, GetConsoleEntriesWithStackOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if !pilot.IsConsoleDebuggerEnabled() {
		return nil, GetConsoleEntriesWithStackOutput{}, fmt.Errorf("console debugger not enabled - call enable_console_debugger first")
	}

	entries := pilot.ConsoleEntries()

	// Filter by type if specified
	if input.Type != "" {
		filtered := make([]w3pilot.ConsoleEntry, 0)
		for _, e := range entries {
			if string(e.Type) == input.Type {
				filtered = append(filtered, e)
			}
		}
		entries = filtered
	}

	// Apply limit
	if input.Limit > 0 && len(entries) > input.Limit {
		entries = entries[len(entries)-input.Limit:]
	}

	output := GetConsoleEntriesWithStackOutput{
		Entries: make([]ConsoleEntryOutput, len(entries)),
		Count:   len(entries),
	}

	for i, e := range entries {
		entry := ConsoleEntryOutput{
			Type:      string(e.Type),
			Text:      e.Text,
			Timestamp: e.Timestamp,
		}

		if e.StackTrace != nil {
			for _, frame := range e.StackTrace.CallFrames {
				entry.StackTrace = append(entry.StackTrace, CallFrameOutput{
					FunctionName: frame.FunctionName,
					URL:          frame.URL,
					LineNumber:   frame.LineNumber,
					ColumnNumber: frame.ColumnNumber,
				})
			}
		}

		output.Entries[i] = entry
	}

	return nil, output, nil
}

// GetBrowserLogs tool - get browser log entries (deprecations, interventions, violations)

type GetBrowserLogsInput struct {
	Source string `json:"source" jsonschema:"Filter by source: network, violation, intervention, deprecation (optional)"`
	Level  string `json:"level" jsonschema:"Filter by level: verbose, info, warning, error (optional)"`
}

type GetBrowserLogsOutput struct {
	Logs  []BrowserLogOutput `json:"logs"`
	Count int                `json:"count"`
}

type BrowserLogOutput struct {
	Source     string            `json:"source"`
	Level      string            `json:"level"`
	Text       string            `json:"text"`
	URL        string            `json:"url,omitempty"`
	LineNumber int               `json:"lineNumber,omitempty"`
	Timestamp  float64           `json:"timestamp"`
	StackTrace []CallFrameOutput `json:"stackTrace,omitempty"`
}

func (s *Server) handleGetBrowserLogs(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input GetBrowserLogsInput,
) (*mcp.CallToolResult, GetBrowserLogsOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, GetBrowserLogsOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if !pilot.IsConsoleDebuggerEnabled() {
		return nil, GetBrowserLogsOutput{}, fmt.Errorf("console debugger not enabled - call enable_console_debugger first")
	}

	logs := pilot.BrowserLogs()

	// Filter by source
	if input.Source != "" {
		filtered := make([]w3pilot.LogEntry, 0)
		for _, l := range logs {
			if l.Source == input.Source {
				filtered = append(filtered, l)
			}
		}
		logs = filtered
	}

	// Filter by level
	if input.Level != "" {
		filtered := make([]w3pilot.LogEntry, 0)
		for _, l := range logs {
			if l.Level == input.Level {
				filtered = append(filtered, l)
			}
		}
		logs = filtered
	}

	output := GetBrowserLogsOutput{
		Logs:  make([]BrowserLogOutput, len(logs)),
		Count: len(logs),
	}

	for i, l := range logs {
		entry := BrowserLogOutput{
			Source:     l.Source,
			Level:      l.Level,
			Text:       l.Text,
			URL:        l.URL,
			LineNumber: l.LineNumber,
			Timestamp:  l.Timestamp,
		}

		if l.StackTrace != nil {
			for _, frame := range l.StackTrace.CallFrames {
				entry.StackTrace = append(entry.StackTrace, CallFrameOutput{
					FunctionName: frame.FunctionName,
					URL:          frame.URL,
					LineNumber:   frame.LineNumber,
					ColumnNumber: frame.ColumnNumber,
				})
			}
		}

		output.Logs[i] = entry
	}

	return nil, output, nil
}

// DisableConsoleDebugger tool - stop capturing console messages

type DisableConsoleDebuggerInput struct{}

type DisableConsoleDebuggerOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleDisableConsoleDebugger(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input DisableConsoleDebuggerInput,
) (*mcp.CallToolResult, DisableConsoleDebuggerOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, DisableConsoleDebuggerOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	if !pilot.HasCDP() {
		return nil, DisableConsoleDebuggerOutput{}, fmt.Errorf("CDP not available")
	}

	if err := pilot.DisableConsoleDebugger(ctx); err != nil {
		return nil, DisableConsoleDebuggerOutput{}, fmt.Errorf("failed to disable console debugger: %w", err)
	}

	return nil, DisableConsoleDebuggerOutput{
		Message: "Console debugger disabled",
	}, nil
}
