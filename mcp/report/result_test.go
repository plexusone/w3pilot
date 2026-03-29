package report

import (
	"testing"
)

func TestComputeOverallStatus(t *testing.T) {
	tests := []struct {
		name     string
		steps    []StepResult
		expected Status
	}{
		{
			name:     "empty steps returns GO",
			steps:    []StepResult{},
			expected: StatusGo,
		},
		{
			name: "all GO returns GO",
			steps: []StepResult{
				{Status: StatusGo},
				{Status: StatusGo},
				{Status: StatusGo},
			},
			expected: StatusGo,
		},
		{
			name: "has WARN no NO-GO returns WARN",
			steps: []StepResult{
				{Status: StatusGo},
				{Status: StatusWarn},
				{Status: StatusGo},
			},
			expected: StatusWarn,
		},
		{
			name: "has NO-GO returns NO-GO",
			steps: []StepResult{
				{Status: StatusGo},
				{Status: StatusNoGo},
				{Status: StatusGo},
			},
			expected: StatusNoGo,
		},
		{
			name: "NO-GO takes priority over WARN",
			steps: []StepResult{
				{Status: StatusWarn},
				{Status: StatusNoGo},
				{Status: StatusWarn},
			},
			expected: StatusNoGo,
		},
		{
			name: "all SKIP returns SKIP",
			steps: []StepResult{
				{Status: StatusSkip},
				{Status: StatusSkip},
			},
			expected: StatusSkip,
		},
		{
			name: "mix with SKIP and GO returns GO",
			steps: []StepResult{
				{Status: StatusSkip},
				{Status: StatusGo},
			},
			expected: StatusGo,
		},
		{
			name: "single GO",
			steps: []StepResult{
				{Status: StatusGo},
			},
			expected: StatusGo,
		},
		{
			name: "single NO-GO",
			steps: []StepResult{
				{Status: StatusNoGo},
			},
			expected: StatusNoGo,
		},
		{
			name: "single WARN",
			steps: []StepResult{
				{Status: StatusWarn},
			},
			expected: StatusWarn,
		},
		{
			name: "single SKIP",
			steps: []StepResult{
				{Status: StatusSkip},
			},
			expected: StatusSkip,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ComputeOverallStatus(tt.steps)
			if got != tt.expected {
				t.Errorf("ComputeOverallStatus() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestComputeTotalDuration(t *testing.T) {
	tests := []struct {
		name     string
		steps    []StepResult
		expected int64
	}{
		{
			name:     "empty steps returns 0",
			steps:    []StepResult{},
			expected: 0,
		},
		{
			name: "single step returns that duration",
			steps: []StepResult{
				{DurationMS: 100},
			},
			expected: 100,
		},
		{
			name: "multiple steps returns sum",
			steps: []StepResult{
				{DurationMS: 100},
				{DurationMS: 200},
				{DurationMS: 300},
			},
			expected: 600,
		},
		{
			name: "includes zero durations",
			steps: []StepResult{
				{DurationMS: 100},
				{DurationMS: 0},
				{DurationMS: 50},
			},
			expected: 150,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ComputeTotalDuration(tt.steps)
			if got != tt.expected {
				t.Errorf("ComputeTotalDuration() = %v, want %v", got, tt.expected)
			}
		})
	}
}
