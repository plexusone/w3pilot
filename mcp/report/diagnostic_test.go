package report

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewDiagnosticReport(t *testing.T) {
	tr := &TestResult{
		Project:     "test-project",
		Target:      "test-target",
		Status:      StatusGo,
		DurationMS:  1000,
		GeneratedAt: time.Now(),
		Steps: []StepResult{
			{ID: "step-1", Status: StatusGo},
		},
	}

	dr := NewDiagnosticReport(tr)

	if dr.Project != tr.Project {
		t.Errorf("NewDiagnosticReport().Project = %v, want %v", dr.Project, tr.Project)
	}
	if dr.Target != tr.Target {
		t.Errorf("NewDiagnosticReport().Target = %v, want %v", dr.Target, tr.Target)
	}
	if dr.Status != tr.Status {
		t.Errorf("NewDiagnosticReport().Status = %v, want %v", dr.Status, tr.Status)
	}
	if len(dr.Steps) != len(tr.Steps) {
		t.Errorf("NewDiagnosticReport().Steps length = %v, want %v", len(dr.Steps), len(tr.Steps))
	}
}

func TestDiagnosticReportJSON(t *testing.T) {
	tr := &TestResult{
		Project:     "test-project",
		Target:      "test-target",
		Status:      StatusGo,
		DurationMS:  500,
		GeneratedAt: time.Now(),
	}
	dr := NewDiagnosticReport(tr)

	data, err := dr.JSON()
	if err != nil {
		t.Fatalf("DiagnosticReport.JSON() error = %v", err)
	}

	var parsed map[string]any
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Errorf("DiagnosticReport.JSON() produced invalid JSON: %v", err)
	}

	if parsed["project"] != "test-project" {
		t.Errorf("JSON project = %v, want %v", parsed["project"], "test-project")
	}
}

func TestGenerateRecommendations(t *testing.T) {
	tests := []struct {
		name             string
		steps            []StepResult
		wantRecommends   int
		containsSubstr   []string
		notContainsCount int
	}{
		{
			name:           "no errors produces no recommendations",
			steps:          []StepResult{{Status: StatusGo}},
			wantRecommends: 0,
		},
		{
			name: "ElementNotFoundError with suggestions",
			steps: []StepResult{
				{
					Status: StatusNoGo,
					Error: &StepError{
						Type:        "ElementNotFoundError",
						Selector:    "#missing",
						Suggestions: []string{"#found"},
					},
				},
			},
			wantRecommends: 1,
			containsSubstr: []string{"#missing", "Try: #found"},
		},
		{
			name: "ElementNotFoundError without suggestions",
			steps: []StepResult{
				{
					Status: StatusNoGo,
					Error: &StepError{
						Type:     "ElementNotFoundError",
						Selector: "#missing",
					},
				},
			},
			wantRecommends: 1,
			containsSubstr: []string{"#missing", "Check if the element exists"},
		},
		{
			name: "TimeoutError",
			steps: []StepResult{
				{
					Status: StatusNoGo,
					Error: &StepError{
						Type: "TimeoutError",
					},
				},
			},
			wantRecommends: 1,
			containsSubstr: []string{"timed out"},
		},
		{
			name: "NavigationError",
			steps: []StepResult{
				{
					Status: StatusNoGo,
					Error: &StepError{
						Type: "NavigationError",
					},
				},
			},
			wantRecommends: 1,
			containsSubstr: []string{"Navigation failed"},
		},
		{
			name: "ClickError",
			steps: []StepResult{
				{
					Status: StatusNoGo,
					Error: &StepError{
						Type: "ClickError",
					},
				},
			},
			wantRecommends: 1,
			containsSubstr: []string{"Click failed"},
		},
		{
			name: "network error 404",
			steps: []StepResult{
				{
					Status: StatusNoGo,
					Error:  &StepError{Type: "SomeError"},
					Network: []NetworkError{
						{URL: "http://example.com/missing", StatusCode: 404},
					},
				},
			},
			wantRecommends: 1,
			containsSubstr: []string{"404", "missing resource"},
		},
		{
			name: "network error 5xx",
			steps: []StepResult{
				{
					Status: StatusNoGo,
					Error:  &StepError{Type: "SomeError"},
					Network: []NetworkError{
						{URL: "http://example.com/api", StatusCode: 500},
					},
				},
			},
			wantRecommends: 1,
			containsSubstr: []string{"Server error", "500"},
		},
		{
			name: "console errors",
			steps: []StepResult{
				{
					Status: StatusNoGo,
					Error:  &StepError{Type: "SomeError"},
					Console: []ConsoleEntry{
						{Level: "error", Message: "Uncaught TypeError: foo is undefined"},
					},
				},
			},
			wantRecommends: 1,
			containsSubstr: []string{"JavaScript error"},
		},
		{
			name: "multiple errors produce multiple recommendations",
			steps: []StepResult{
				{
					Status: StatusNoGo,
					Error:  &StepError{Type: "TimeoutError"},
				},
				{
					Status: StatusNoGo,
					Error:  &StepError{Type: "NavigationError"},
				},
			},
			wantRecommends: 2,
		},
		{
			name: "WARN status does not produce recommendations",
			steps: []StepResult{
				{
					Status: StatusWarn,
					Error:  &StepError{Type: "SomeError"},
				},
			},
			wantRecommends: 0,
		},
		{
			name: "NO-GO without error does not produce recommendations",
			steps: []StepResult{
				{
					Status: StatusNoGo,
					Error:  nil,
				},
			},
			wantRecommends: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := &TestResult{Steps: tt.steps}
			dr := NewDiagnosticReport(tr)
			dr.GenerateRecommendations()

			if len(dr.Recommendations) != tt.wantRecommends {
				t.Errorf("GenerateRecommendations() produced %d recommendations, want %d",
					len(dr.Recommendations), tt.wantRecommends)
			}

			for _, substr := range tt.containsSubstr {
				found := false
				for _, rec := range dr.Recommendations {
					if containsSubstring(rec, substr) {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("GenerateRecommendations() should contain %q", substr)
				}
			}
		})
	}
}

func TestTruncate(t *testing.T) {
	tests := []struct {
		name   string
		input  string
		maxLen int
		want   string
	}{
		{
			name:   "string shorter than max",
			input:  "hello",
			maxLen: 10,
			want:   "hello",
		},
		{
			name:   "string equal to max",
			input:  "hello",
			maxLen: 5,
			want:   "hello",
		},
		{
			name:   "string longer than max adds ellipsis",
			input:  "hello world",
			maxLen: 8,
			want:   "hello...",
		},
		{
			name:   "empty string",
			input:  "",
			maxLen: 10,
			want:   "",
		},
		{
			name:   "max length 3 edge case",
			input:  "abcdef",
			maxLen: 3,
			want:   "...",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := truncate(tt.input, tt.maxLen)
			if got != tt.want {
				t.Errorf("truncate(%q, %d) = %q, want %q", tt.input, tt.maxLen, got, tt.want)
			}
		})
	}
}

func containsSubstring(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && containsSubstringHelper(s, substr))
}

func containsSubstringHelper(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
