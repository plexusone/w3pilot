//nolint:dupl // MCP tool handlers intentionally follow a consistent pattern
package mcp

import (
	"context"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	vibium "github.com/plexusone/w3pilot"
	"github.com/plexusone/w3pilot/mcp/report"
)

// ClockInstallInput represents input for the clock_install tool.
type ClockInstallInput struct {
	Time string `json:"time" jsonschema:"Optional ISO 8601 timestamp to set as initial time (e.g. '2024-12-25T00:00:00Z')"`
}

// ClockInstallOutput represents output from the clock_install tool.
type ClockInstallOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleClockInstall(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ClockInstallInput,
) (*mcp.CallToolResult, ClockInstallOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, ClockInstallOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	clock, err := pilot.Clock(ctx)
	if err != nil {
		return nil, ClockInstallOutput{}, fmt.Errorf("failed to get clock: %w", err)
	}

	start := time.Now()

	var opts *vibium.ClockInstallOptions
	if input.Time != "" {
		t, err := time.Parse(time.RFC3339, input.Time)
		if err != nil {
			return nil, ClockInstallOutput{}, fmt.Errorf("invalid time format (use ISO 8601): %w", err)
		}
		opts = &vibium.ClockInstallOptions{Time: t}
	}

	err = clock.Install(ctx, opts)
	duration := time.Since(start)

	result := report.StepResult{
		ID:         s.session.NextStepID("clock_install"),
		Action:     "clock_install",
		Args:       map[string]any{"time": input.Time},
		DurationMS: duration.Milliseconds(),
	}

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityMedium
		result.Error = &report.StepError{
			Type:    "ClockError",
			Message: err.Error(),
		}
		s.session.RecordStep(result)
		return nil, ClockInstallOutput{}, fmt.Errorf("failed to install clock: %w", err)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	s.session.RecordStep(result)

	msg := "Fake clock installed"
	if input.Time != "" {
		msg = fmt.Sprintf("Fake clock installed at %s", input.Time)
	}

	return nil, ClockInstallOutput{Message: msg}, nil
}

// ClockSetTimeInput represents input for the clock_set_time tool.
type ClockSetTimeInput struct {
	Time string `json:"time" jsonschema:"ISO 8601 timestamp to set (e.g. '2024-12-25T00:00:00Z'),required"`
}

// ClockSetTimeOutput represents output from the clock_set_time tool.
type ClockSetTimeOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleClockSetTime(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ClockSetTimeInput,
) (*mcp.CallToolResult, ClockSetTimeOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, ClockSetTimeOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	clock, err := pilot.Clock(ctx)
	if err != nil {
		return nil, ClockSetTimeOutput{}, fmt.Errorf("failed to get clock: %w", err)
	}

	t, err := time.Parse(time.RFC3339, input.Time)
	if err != nil {
		return nil, ClockSetTimeOutput{}, fmt.Errorf("invalid time format (use ISO 8601): %w", err)
	}

	start := time.Now()
	err = clock.SetFixedTime(ctx, t)
	duration := time.Since(start)

	result := report.StepResult{
		ID:         s.session.NextStepID("clock_set_time"),
		Action:     "clock_set_time",
		Args:       map[string]any{"time": input.Time},
		DurationMS: duration.Milliseconds(),
	}

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityMedium
		result.Error = &report.StepError{
			Type:    "ClockError",
			Message: err.Error(),
		}
		s.session.RecordStep(result)
		return nil, ClockSetTimeOutput{}, fmt.Errorf("failed to set time: %w", err)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	s.session.RecordStep(result)

	return nil, ClockSetTimeOutput{Message: fmt.Sprintf("Clock set to %s", input.Time)}, nil
}

// ClockFastForwardInput represents input for the clock_fast_forward tool.
type ClockFastForwardInput struct {
	Milliseconds int64 `json:"milliseconds" jsonschema:"Number of milliseconds to advance,required"`
}

// ClockFastForwardOutput represents output from the clock_fast_forward tool.
type ClockFastForwardOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleClockFastForward(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ClockFastForwardInput,
) (*mcp.CallToolResult, ClockFastForwardOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, ClockFastForwardOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	clock, err := pilot.Clock(ctx)
	if err != nil {
		return nil, ClockFastForwardOutput{}, fmt.Errorf("failed to get clock: %w", err)
	}

	start := time.Now()
	err = clock.FastForward(ctx, input.Milliseconds)
	duration := time.Since(start)

	result := report.StepResult{
		ID:         s.session.NextStepID("clock_fast_forward"),
		Action:     "clock_fast_forward",
		Args:       map[string]any{"milliseconds": input.Milliseconds},
		DurationMS: duration.Milliseconds(),
	}

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityMedium
		result.Error = &report.StepError{
			Type:    "ClockError",
			Message: err.Error(),
		}
		s.session.RecordStep(result)
		return nil, ClockFastForwardOutput{}, fmt.Errorf("failed to fast forward: %w", err)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	s.session.RecordStep(result)

	return nil, ClockFastForwardOutput{
		Message: fmt.Sprintf("Advanced clock by %dms (timers not fired)", input.Milliseconds),
	}, nil
}

// ClockRunForInput represents input for the clock_run_for tool.
type ClockRunForInput struct {
	Milliseconds int64 `json:"milliseconds" jsonschema:"Number of milliseconds to advance (fires pending timers),required"`
}

// ClockRunForOutput represents output from the clock_run_for tool.
type ClockRunForOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleClockRunFor(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ClockRunForInput,
) (*mcp.CallToolResult, ClockRunForOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, ClockRunForOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	clock, err := pilot.Clock(ctx)
	if err != nil {
		return nil, ClockRunForOutput{}, fmt.Errorf("failed to get clock: %w", err)
	}

	start := time.Now()
	err = clock.RunFor(ctx, input.Milliseconds)
	duration := time.Since(start)

	result := report.StepResult{
		ID:         s.session.NextStepID("clock_run_for"),
		Action:     "clock_run_for",
		Args:       map[string]any{"milliseconds": input.Milliseconds},
		DurationMS: duration.Milliseconds(),
	}

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityMedium
		result.Error = &report.StepError{
			Type:    "ClockError",
			Message: err.Error(),
		}
		s.session.RecordStep(result)
		return nil, ClockRunForOutput{}, fmt.Errorf("failed to run for: %w", err)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	s.session.RecordStep(result)

	return nil, ClockRunForOutput{
		Message: fmt.Sprintf("Advanced clock by %dms (timers fired)", input.Milliseconds),
	}, nil
}

// ClockPauseAtInput represents input for the clock_pause_at tool.
type ClockPauseAtInput struct {
	Time string `json:"time" jsonschema:"ISO 8601 timestamp to pause at (e.g. '2024-12-25T00:00:00Z'),required"`
}

// ClockPauseAtOutput represents output from the clock_pause_at tool.
type ClockPauseAtOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleClockPauseAt(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ClockPauseAtInput,
) (*mcp.CallToolResult, ClockPauseAtOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, ClockPauseAtOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	clock, err := pilot.Clock(ctx)
	if err != nil {
		return nil, ClockPauseAtOutput{}, fmt.Errorf("failed to get clock: %w", err)
	}

	t, err := time.Parse(time.RFC3339, input.Time)
	if err != nil {
		return nil, ClockPauseAtOutput{}, fmt.Errorf("invalid time format (use ISO 8601): %w", err)
	}

	start := time.Now()
	err = clock.PauseAt(ctx, t)
	duration := time.Since(start)

	result := report.StepResult{
		ID:         s.session.NextStepID("clock_pause_at"),
		Action:     "clock_pause_at",
		Args:       map[string]any{"time": input.Time},
		DurationMS: duration.Milliseconds(),
	}

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityMedium
		result.Error = &report.StepError{
			Type:    "ClockError",
			Message: err.Error(),
		}
		s.session.RecordStep(result)
		return nil, ClockPauseAtOutput{}, fmt.Errorf("failed to pause at: %w", err)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	s.session.RecordStep(result)

	return nil, ClockPauseAtOutput{Message: fmt.Sprintf("Clock paused at %s", input.Time)}, nil
}

// ClockResumeInput represents input for the clock_resume tool.
type ClockResumeInput struct{}

// ClockResumeOutput represents output from the clock_resume tool.
type ClockResumeOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleClockResume(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ClockResumeInput,
) (*mcp.CallToolResult, ClockResumeOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, ClockResumeOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	clock, err := pilot.Clock(ctx)
	if err != nil {
		return nil, ClockResumeOutput{}, fmt.Errorf("failed to get clock: %w", err)
	}

	start := time.Now()
	err = clock.Resume(ctx)
	duration := time.Since(start)

	result := report.StepResult{
		ID:         s.session.NextStepID("clock_resume"),
		Action:     "clock_resume",
		DurationMS: duration.Milliseconds(),
	}

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityMedium
		result.Error = &report.StepError{
			Type:    "ClockError",
			Message: err.Error(),
		}
		s.session.RecordStep(result)
		return nil, ClockResumeOutput{}, fmt.Errorf("failed to resume: %w", err)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	s.session.RecordStep(result)

	return nil, ClockResumeOutput{Message: "Clock resumed"}, nil
}

// ClockSetTimezoneInput represents input for the clock_set_timezone tool.
type ClockSetTimezoneInput struct {
	Timezone string `json:"timezone" jsonschema:"IANA timezone name (e.g. 'America/New_York', 'Europe/London', 'Asia/Tokyo'),required"`
}

// ClockSetTimezoneOutput represents output from the clock_set_timezone tool.
type ClockSetTimezoneOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleClockSetTimezone(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input ClockSetTimezoneInput,
) (*mcp.CallToolResult, ClockSetTimezoneOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, ClockSetTimezoneOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	clock, err := pilot.Clock(ctx)
	if err != nil {
		return nil, ClockSetTimezoneOutput{}, fmt.Errorf("failed to get clock: %w", err)
	}

	start := time.Now()
	err = clock.SetTimezone(ctx, input.Timezone)
	duration := time.Since(start)

	result := report.StepResult{
		ID:         s.session.NextStepID("clock_set_timezone"),
		Action:     "clock_set_timezone",
		Args:       map[string]any{"timezone": input.Timezone},
		DurationMS: duration.Milliseconds(),
	}

	if err != nil {
		result.Status = report.StatusNoGo
		result.Severity = report.SeverityMedium
		result.Error = &report.StepError{
			Type:    "ClockError",
			Message: err.Error(),
		}
		s.session.RecordStep(result)
		return nil, ClockSetTimezoneOutput{}, fmt.Errorf("failed to set timezone: %w", err)
	}

	result.Status = report.StatusGo
	result.Severity = report.SeverityInfo
	s.session.RecordStep(result)

	return nil, ClockSetTimezoneOutput{Message: fmt.Sprintf("Timezone set to %s", input.Timezone)}, nil
}
