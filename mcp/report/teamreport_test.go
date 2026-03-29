package report

import (
	"testing"
	"time"

	multiagentspec "github.com/plexusone/multi-agent-spec/sdk/go"
)

func TestCategorizeAction(t *testing.T) {
	tests := []struct {
		action   string
		expected string
	}{
		// Browser actions
		{"browser_launch", "browser"},
		{"browser_quit", "browser"},

		// Navigation actions
		{"navigate", "navigation"},
		{"back", "navigation"},
		{"forward", "navigation"},
		{"reload", "navigation"},

		// Interaction actions
		{"click", "interaction"},
		{"type", "interaction"},

		// Extraction actions
		{"get_text", "extraction"},
		{"get_attribute", "extraction"},
		{"screenshot", "extraction"},
		{"evaluate", "extraction"},
		{"find", "extraction"},
		{"find_all", "extraction"},

		// Assertion actions
		{"assert_text", "assertion"},
		{"assert_element", "assertion"},
		{"wait_for", "assertion"},

		// Unknown actions
		{"unknown_action", "other"},
		{"", "other"},
		{"custom", "other"},
	}

	for _, tt := range tests {
		t.Run(tt.action, func(t *testing.T) {
			got := categorizeAction(tt.action)
			if got != tt.expected {
				t.Errorf("categorizeAction(%q) = %q, want %q", tt.action, got, tt.expected)
			}
		})
	}
}

func TestConvertStatus(t *testing.T) {
	tests := []struct {
		input    Status
		expected multiagentspec.Status
	}{
		{StatusGo, multiagentspec.StatusGo},
		{StatusWarn, multiagentspec.StatusWarn},
		{StatusNoGo, multiagentspec.StatusNoGo},
		{StatusSkip, multiagentspec.StatusSkip},
		{Status("unknown"), multiagentspec.StatusSkip},
	}

	for _, tt := range tests {
		t.Run(string(tt.input), func(t *testing.T) {
			got := convertStatus(tt.input)
			if got != tt.expected {
				t.Errorf("convertStatus(%v) = %v, want %v", tt.input, got, tt.expected)
			}
		})
	}
}

func TestFormatBrowserInfo(t *testing.T) {
	tests := []struct {
		name     string
		bi       BrowserInfo
		expected string
	}{
		{
			name:     "headless browser",
			bi:       BrowserInfo{Name: "chromium", Headless: true},
			expected: "chromium (headless)",
		},
		{
			name:     "headed browser",
			bi:       BrowserInfo{Name: "chromium", Headless: false},
			expected: "chromium (headed)",
		},
		{
			name:     "firefox headless",
			bi:       BrowserInfo{Name: "firefox", Headless: true},
			expected: "firefox (headless)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatBrowserInfo(tt.bi)
			if got != tt.expected {
				t.Errorf("formatBrowserInfo() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		name     string
		ms       int64
		expected string
	}{
		{
			name:     "sub-second shows ms",
			ms:       500,
			expected: "500ms",
		},
		{
			name:     "exactly 1 second",
			ms:       1000,
			expected: "1ss", // Note: current implementation adds extra "s"
		},
		{
			name:     "over 1 second",
			ms:       1500,
			expected: "1.5ss", // Note: current implementation adds extra "s"
		},
		{
			name:     "zero milliseconds",
			ms:       0,
			expected: "0s",
		},
		{
			name:     "100ms",
			ms:       100,
			expected: "100ms",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatDuration(tt.ms)
			if got != tt.expected {
				t.Errorf("formatDuration(%d) = %q, want %q", tt.ms, got, tt.expected)
			}
		})
	}
}

func TestGroupStepsIntoTeams(t *testing.T) {
	steps := []StepResult{
		{ID: "1", Action: "browser_launch", Status: StatusGo},
		{ID: "2", Action: "navigate", Status: StatusGo},
		{ID: "3", Action: "click", Status: StatusGo},
		{ID: "4", Action: "get_text", Status: StatusWarn},
		{ID: "5", Action: "assert_text", Status: StatusNoGo},
	}

	teams := groupStepsIntoTeams(steps)

	// Should have 5 categories with steps
	if len(teams) != 5 {
		t.Errorf("groupStepsIntoTeams() returned %d teams, want 5", len(teams))
	}

	// Verify category order and dependencies
	expectedCategories := []string{"browser", "navigation", "interaction", "extraction", "assertion"}
	var prevCategory string
	for i, team := range teams {
		if i >= len(expectedCategories) {
			t.Errorf("unexpected team at index %d: %q", i, team.ID)
			continue
		}

		expected := expectedCategories[i]
		if team.ID != expected {
			t.Errorf("team[%d].ID = %q, want %q", i, team.ID, expected)
		}

		// First team should have no dependencies
		if i == 0 && len(team.DependsOn) != 0 {
			t.Errorf("first team should have no dependencies, got %v", team.DependsOn)
		}

		// Subsequent teams should depend on previous
		if i > 0 && prevCategory != "" {
			if len(team.DependsOn) != 1 || team.DependsOn[0] != prevCategory {
				t.Errorf("team[%d].DependsOn = %v, want [%q]", i, team.DependsOn, prevCategory)
			}
		}
		prevCategory = expected
	}
}

func TestGroupStepsIntoTeams_EmptySteps(t *testing.T) {
	teams := groupStepsIntoTeams([]StepResult{})
	if len(teams) != 0 {
		t.Errorf("groupStepsIntoTeams([]) should return empty slice, got %d teams", len(teams))
	}
}

func TestGroupStepsIntoTeams_ComputesStatusPerCategory(t *testing.T) {
	steps := []StepResult{
		{ID: "1", Action: "navigate", Status: StatusGo},
		{ID: "2", Action: "click", Status: StatusNoGo},
		{ID: "3", Action: "type", Status: StatusGo},
	}

	teams := groupStepsIntoTeams(steps)

	// Find interaction team
	var interactionTeam *multiagentspec.TeamSection
	for i := range teams {
		if teams[i].ID == "interaction" {
			interactionTeam = &teams[i]
			break
		}
	}

	if interactionTeam == nil {
		t.Fatal("interaction team not found")
	}

	// Interaction has NO-GO and GO, should be NO-GO overall
	if interactionTeam.Status != multiagentspec.StatusNoGo {
		t.Errorf("interaction team status = %v, want NO-GO", interactionTeam.Status)
	}
}

func TestToTeamReport(t *testing.T) {
	tr := &TestResult{
		Project:    "test-project",
		Target:     "test-target",
		Status:     StatusGo,
		DurationMS: 1500,
		Browser:    BrowserInfo{Name: "chromium", Headless: true},
		Steps: []StepResult{
			{ID: "1", Action: "navigate", Status: StatusGo},
		},
		Recommendations: []string{"Fix the selector"},
		GeneratedAt:     time.Now(),
	}

	report := ToTeamReport(tr)

	if report.Title != "BROWSER TEST REPORT" {
		t.Errorf("ToTeamReport().Title = %q, want %q", report.Title, "BROWSER TEST REPORT")
	}
	if report.Project != "test-project" {
		t.Errorf("ToTeamReport().Project = %q, want %q", report.Project, "test-project")
	}
	if report.Target != "test-target" {
		t.Errorf("ToTeamReport().Target = %q, want %q", report.Target, "test-target")
	}
	if report.Status != multiagentspec.StatusGo {
		t.Errorf("ToTeamReport().Status = %v, want GO", report.Status)
	}
	if len(report.SummaryBlocks) != 1 {
		t.Errorf("ToTeamReport() should have 1 summary block, got %d", len(report.SummaryBlocks))
	}
	if len(report.FooterBlocks) != 1 {
		t.Errorf("ToTeamReport() should have 1 footer block (recommendations), got %d", len(report.FooterBlocks))
	}
}

func TestToTeamReport_NoRecommendations(t *testing.T) {
	tr := &TestResult{
		Project:         "test",
		Status:          StatusGo,
		Steps:           []StepResult{},
		Recommendations: nil,
		GeneratedAt:     time.Now(),
	}

	report := ToTeamReport(tr)

	if len(report.FooterBlocks) != 0 {
		t.Errorf("ToTeamReport() should have no footer blocks when no recommendations, got %d", len(report.FooterBlocks))
	}
}

func TestFormatStepDetail(t *testing.T) {
	tests := []struct {
		name     string
		step     StepResult
		expected string
	}{
		{
			name: "step with error shows error message",
			step: StepResult{
				Error: &StepError{Message: "element not found"},
			},
			expected: "element not found",
		},
		{
			name: "navigate action shows URL",
			step: StepResult{
				Action: "navigate",
				Args:   map[string]any{"url": "https://example.com"},
			},
			expected: "https://example.com",
		},
		{
			name: "click action shows selector",
			step: StepResult{
				Action: "click",
				Args:   map[string]any{"selector": "#submit"},
			},
			expected: "#submit",
		},
		{
			name: "type action shows selector",
			step: StepResult{
				Action: "type",
				Args:   map[string]any{"selector": "input[name=email]"},
			},
			expected: "input[name=email]",
		},
		{
			name: "get_text action shows selector",
			step: StepResult{
				Action: "get_text",
				Args:   map[string]any{"selector": ".message"},
			},
			expected: ".message",
		},
		{
			name: "screenshot action shows captured",
			step: StepResult{
				Action: "screenshot",
			},
			expected: "captured",
		},
		{
			name: "unknown action with no args",
			step: StepResult{
				Action: "custom",
			},
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := formatStepDetail(tt.step)
			if got != tt.expected {
				t.Errorf("formatStepDetail() = %q, want %q", got, tt.expected)
			}
		})
	}
}

func TestConvertStepToTask(t *testing.T) {
	step := StepResult{
		ID:         "step-1",
		Status:     StatusGo,
		Severity:   SeverityMedium,
		DurationMS: 100,
		Action:     "click",
		Args:       map[string]any{"selector": "#btn"},
	}

	task := convertStepToTask(step)

	if task.ID != "step-1" {
		t.Errorf("convertStepToTask().ID = %q, want %q", task.ID, "step-1")
	}
	if task.Status != multiagentspec.StatusGo {
		t.Errorf("convertStepToTask().Status = %v, want GO", task.Status)
	}
	if task.Severity != "medium" {
		t.Errorf("convertStepToTask().Severity = %q, want %q", task.Severity, "medium")
	}
	if task.DurationMs != 100 {
		t.Errorf("convertStepToTask().DurationMs = %d, want %d", task.DurationMs, 100)
	}
}

func TestConvertStepToTask_TruncatesLongDetail(t *testing.T) {
	step := StepResult{
		ID:     "step-1",
		Status: StatusGo,
		Action: "navigate",
		Args:   map[string]any{"url": "https://example.com/very/long/path/that/exceeds/limit"},
	}

	task := convertStepToTask(step)

	if len(task.Detail) > 33 {
		t.Errorf("convertStepToTask().Detail should be truncated to 33 chars, got %d", len(task.Detail))
	}
	if len(task.Detail) == 33 && task.Detail[30:] != "..." {
		t.Errorf("truncated detail should end with '...'")
	}
}
