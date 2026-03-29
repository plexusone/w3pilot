package mcp

import (
	"encoding/json"
	"testing"

	"github.com/plexusone/w3pilot/script"
)

func TestNewRecorder(t *testing.T) {
	r := NewRecorder()

	if r == nil {
		t.Fatal("NewRecorder() returned nil")
	}
	if r.IsRecording() {
		t.Error("NewRecorder() should not be recording by default")
	}
	if len(r.Steps()) != 0 {
		t.Error("NewRecorder() should have empty steps")
	}
}

func TestRecorderStartStop(t *testing.T) {
	r := NewRecorder()

	// Start recording
	r.Start(RecorderMetadata{Name: "Test"})
	if !r.IsRecording() {
		t.Error("Start() should set recording to true")
	}

	// Add a step while recording
	r.AddStep(script.Step{Action: script.ActionClick})
	if len(r.Steps()) != 1 {
		t.Error("Step should be added while recording")
	}

	// Start again should clear steps
	r.Start(RecorderMetadata{Name: "Test 2"})
	if len(r.Steps()) != 0 {
		t.Error("Start() should clear previous steps")
	}

	// Stop recording
	r.Stop()
	if r.IsRecording() {
		t.Error("Stop() should set recording to false")
	}
}

func TestRecorderAddStep(t *testing.T) {
	r := NewRecorder()

	// Add step when not recording - should not add
	r.AddStep(script.Step{Action: script.ActionClick})
	if len(r.Steps()) != 0 {
		t.Error("AddStep() should not add when not recording")
	}

	// Start recording and add step
	r.Start(RecorderMetadata{})
	r.AddStep(script.Step{Action: script.ActionClick})
	r.AddStep(script.Step{Action: script.ActionType})
	if len(r.Steps()) != 2 {
		t.Errorf("Steps() should have 2 steps, got %d", len(r.Steps()))
	}

	// Stop and try to add - should not add
	r.Stop()
	r.AddStep(script.Step{Action: script.ActionNavigate})
	if len(r.Steps()) != 2 {
		t.Error("AddStep() should not add when stopped")
	}
}

func TestRecorderSteps(t *testing.T) {
	r := NewRecorder()
	r.Start(RecorderMetadata{})
	r.AddStep(script.Step{Action: script.ActionClick})

	// Get steps and verify it's a copy
	steps := r.Steps()
	if len(steps) != 1 {
		t.Fatalf("Steps() should return 1 step, got %d", len(steps))
	}

	// Modify returned slice
	steps[0].Action = script.ActionType

	// Original should be unchanged
	original := r.Steps()
	if original[0].Action != script.ActionClick {
		t.Error("Steps() should return a copy, not the original slice")
	}
}

func TestRecorderClear(t *testing.T) {
	r := NewRecorder()
	r.Start(RecorderMetadata{})
	r.AddStep(script.Step{Action: script.ActionClick})
	r.AddStep(script.Step{Action: script.ActionType})

	if len(r.Steps()) != 2 {
		t.Fatalf("Should have 2 steps before clear")
	}

	r.Clear()
	if len(r.Steps()) != 0 {
		t.Error("Clear() should remove all steps")
	}

	// Should still be recording
	if !r.IsRecording() {
		t.Error("Clear() should not affect recording state")
	}
}

func TestRecorderStepCount(t *testing.T) {
	r := NewRecorder()
	r.Start(RecorderMetadata{})

	if r.StepCount() != 0 {
		t.Error("StepCount() should be 0 initially")
	}

	r.AddStep(script.Step{Action: script.ActionClick})
	r.AddStep(script.Step{Action: script.ActionType})

	if r.StepCount() != 2 {
		t.Errorf("StepCount() = %d, want 2", r.StepCount())
	}
}

func TestRecorderExport(t *testing.T) {
	r := NewRecorder()
	r.Start(RecorderMetadata{
		Name:        "My Test",
		Description: "Test description",
		BaseURL:     "https://example.com",
	})
	r.AddStep(script.Step{Action: script.ActionNavigate, URL: "/path"})
	r.AddStep(script.Step{Action: script.ActionClick, Selector: "#btn"})

	s := r.Export()

	if s.Name != "My Test" {
		t.Errorf("Export().Name = %q, want %q", s.Name, "My Test")
	}
	if s.Description != "Test description" {
		t.Errorf("Export().Description = %q, want %q", s.Description, "Test description")
	}
	if s.Version != 1 {
		t.Errorf("Export().Version = %d, want 1", s.Version)
	}
	if s.BaseURL != "https://example.com" {
		t.Errorf("Export().BaseURL = %q, want %q", s.BaseURL, "https://example.com")
	}
	if len(s.Steps) != 2 {
		t.Errorf("Export().Steps length = %d, want 2", len(s.Steps))
	}
}

func TestRecorderExport_DefaultName(t *testing.T) {
	r := NewRecorder()
	r.Start(RecorderMetadata{})

	s := r.Export()

	if s.Name != "Recorded Test" {
		t.Errorf("Export().Name = %q, want %q", s.Name, "Recorded Test")
	}
}

func TestRecorderExportJSON(t *testing.T) {
	r := NewRecorder()
	r.Start(RecorderMetadata{Name: "JSON Test"})
	r.AddStep(script.Step{Action: script.ActionClick, Selector: "#btn"})

	data, err := r.ExportJSON()
	if err != nil {
		t.Fatalf("ExportJSON() error = %v", err)
	}

	var s script.Script
	if err := json.Unmarshal(data, &s); err != nil {
		t.Fatalf("ExportJSON() produced invalid JSON: %v", err)
	}

	if s.Name != "JSON Test" {
		t.Errorf("Unmarshaled Name = %q, want %q", s.Name, "JSON Test")
	}
	if len(s.Steps) != 1 {
		t.Errorf("Unmarshaled Steps length = %d, want 1", len(s.Steps))
	}
}

func TestRecordActions(t *testing.T) {
	tests := []struct {
		name       string
		recordFunc func(r *Recorder)
		wantAction script.Action
		validate   func(t *testing.T, step script.Step)
	}{
		{
			name:       "RecordNavigate",
			recordFunc: func(r *Recorder) { r.RecordNavigate("https://example.com") },
			wantAction: script.ActionNavigate,
			validate: func(t *testing.T, step script.Step) {
				if step.URL != "https://example.com" {
					t.Errorf("URL = %q, want %q", step.URL, "https://example.com")
				}
			},
		},
		{
			name:       "RecordClick",
			recordFunc: func(r *Recorder) { r.RecordClick("#btn") },
			wantAction: script.ActionClick,
			validate: func(t *testing.T, step script.Step) {
				if step.Selector != "#btn" {
					t.Errorf("Selector = %q, want %q", step.Selector, "#btn")
				}
			},
		},
		{
			name:       "RecordDblClick",
			recordFunc: func(r *Recorder) { r.RecordDblClick("#btn") },
			wantAction: script.ActionDblClick,
			validate: func(t *testing.T, step script.Step) {
				if step.Selector != "#btn" {
					t.Errorf("Selector = %q, want %q", step.Selector, "#btn")
				}
			},
		},
		{
			name:       "RecordType",
			recordFunc: func(r *Recorder) { r.RecordType("#input", "hello") },
			wantAction: script.ActionType,
			validate: func(t *testing.T, step script.Step) {
				if step.Selector != "#input" {
					t.Errorf("Selector = %q, want %q", step.Selector, "#input")
				}
				if step.Text != "hello" {
					t.Errorf("Text = %q, want %q", step.Text, "hello")
				}
			},
		},
		{
			name:       "RecordFill",
			recordFunc: func(r *Recorder) { r.RecordFill("#input", "value") },
			wantAction: script.ActionFill,
			validate: func(t *testing.T, step script.Step) {
				if step.Selector != "#input" {
					t.Errorf("Selector = %q, want %q", step.Selector, "#input")
				}
				if step.Value != "value" {
					t.Errorf("Value = %q, want %q", step.Value, "value")
				}
			},
		},
		{
			name:       "RecordClear",
			recordFunc: func(r *Recorder) { r.RecordClear("#input") },
			wantAction: script.ActionClear,
			validate: func(t *testing.T, step script.Step) {
				if step.Selector != "#input" {
					t.Errorf("Selector = %q, want %q", step.Selector, "#input")
				}
			},
		},
		{
			name:       "RecordPress",
			recordFunc: func(r *Recorder) { r.RecordPress("#input", "Enter") },
			wantAction: script.ActionPress,
			validate: func(t *testing.T, step script.Step) {
				if step.Selector != "#input" {
					t.Errorf("Selector = %q, want %q", step.Selector, "#input")
				}
				if step.Key != "Enter" {
					t.Errorf("Key = %q, want %q", step.Key, "Enter")
				}
			},
		},
		{
			name:       "RecordCheck",
			recordFunc: func(r *Recorder) { r.RecordCheck("#checkbox") },
			wantAction: script.ActionCheck,
			validate: func(t *testing.T, step script.Step) {
				if step.Selector != "#checkbox" {
					t.Errorf("Selector = %q, want %q", step.Selector, "#checkbox")
				}
			},
		},
		{
			name:       "RecordUncheck",
			recordFunc: func(r *Recorder) { r.RecordUncheck("#checkbox") },
			wantAction: script.ActionUncheck,
			validate: func(t *testing.T, step script.Step) {
				if step.Selector != "#checkbox" {
					t.Errorf("Selector = %q, want %q", step.Selector, "#checkbox")
				}
			},
		},
		{
			name:       "RecordSelect",
			recordFunc: func(r *Recorder) { r.RecordSelect("#select", "option1") },
			wantAction: script.ActionSelect,
			validate: func(t *testing.T, step script.Step) {
				if step.Selector != "#select" {
					t.Errorf("Selector = %q, want %q", step.Selector, "#select")
				}
				if step.Value != "option1" {
					t.Errorf("Value = %q, want %q", step.Value, "option1")
				}
			},
		},
		{
			name:       "RecordHover",
			recordFunc: func(r *Recorder) { r.RecordHover("#element") },
			wantAction: script.ActionHover,
			validate: func(t *testing.T, step script.Step) {
				if step.Selector != "#element" {
					t.Errorf("Selector = %q, want %q", step.Selector, "#element")
				}
			},
		},
		{
			name:       "RecordFocus",
			recordFunc: func(r *Recorder) { r.RecordFocus("#input") },
			wantAction: script.ActionFocus,
			validate: func(t *testing.T, step script.Step) {
				if step.Selector != "#input" {
					t.Errorf("Selector = %q, want %q", step.Selector, "#input")
				}
			},
		},
		{
			name:       "RecordScrollIntoView",
			recordFunc: func(r *Recorder) { r.RecordScrollIntoView("#element") },
			wantAction: script.ActionScrollIntoView,
			validate: func(t *testing.T, step script.Step) {
				if step.Selector != "#element" {
					t.Errorf("Selector = %q, want %q", step.Selector, "#element")
				}
			},
		},
		{
			name:       "RecordScreenshot",
			recordFunc: func(r *Recorder) { r.RecordScreenshot("screenshot.png", true) },
			wantAction: script.ActionScreenshot,
			validate: func(t *testing.T, step script.Step) {
				if step.File != "screenshot.png" {
					t.Errorf("File = %q, want %q", step.File, "screenshot.png")
				}
				if !step.FullPage {
					t.Error("FullPage should be true")
				}
			},
		},
		{
			name:       "RecordEval",
			recordFunc: func(r *Recorder) { r.RecordEval("return document.title") },
			wantAction: script.ActionEval,
			validate: func(t *testing.T, step script.Step) {
				if step.Script != "return document.title" {
					t.Errorf("Script = %q, want %q", step.Script, "return document.title")
				}
			},
		},
		{
			name:       "RecordWait",
			recordFunc: func(r *Recorder) { r.RecordWait("1s") },
			wantAction: script.ActionWait,
			validate: func(t *testing.T, step script.Step) {
				if step.Duration != "1s" {
					t.Errorf("Duration = %q, want %q", step.Duration, "1s")
				}
			},
		},
		{
			name:       "RecordWaitForSelector",
			recordFunc: func(r *Recorder) { r.RecordWaitForSelector("#element", "visible") },
			wantAction: script.ActionWaitForSelector,
			validate: func(t *testing.T, step script.Step) {
				if step.Selector != "#element" {
					t.Errorf("Selector = %q, want %q", step.Selector, "#element")
				}
				if step.State != "visible" {
					t.Errorf("State = %q, want %q", step.State, "visible")
				}
			},
		},
		{
			name:       "RecordWaitForURL",
			recordFunc: func(r *Recorder) { r.RecordWaitForURL("*/success") },
			wantAction: script.ActionWaitForURL,
			validate: func(t *testing.T, step script.Step) {
				if step.Pattern != "*/success" {
					t.Errorf("Pattern = %q, want %q", step.Pattern, "*/success")
				}
			},
		},
		{
			name:       "RecordWaitForLoad",
			recordFunc: func(r *Recorder) { r.RecordWaitForLoad("networkidle") },
			wantAction: script.ActionWaitForLoad,
			validate: func(t *testing.T, step script.Step) {
				if step.LoadState != "networkidle" {
					t.Errorf("LoadState = %q, want %q", step.LoadState, "networkidle")
				}
			},
		},
		{
			name:       "RecordSetViewport",
			recordFunc: func(r *Recorder) { r.RecordSetViewport(1920, 1080) },
			wantAction: script.ActionSetViewport,
			validate: func(t *testing.T, step script.Step) {
				if step.Width != 1920 {
					t.Errorf("Width = %d, want 1920", step.Width)
				}
				if step.Height != 1080 {
					t.Errorf("Height = %d, want 1080", step.Height)
				}
			},
		},
		{
			name:       "RecordBack",
			recordFunc: func(r *Recorder) { r.RecordBack() },
			wantAction: script.ActionBack,
		},
		{
			name:       "RecordForward",
			recordFunc: func(r *Recorder) { r.RecordForward() },
			wantAction: script.ActionForward,
		},
		{
			name:       "RecordReload",
			recordFunc: func(r *Recorder) { r.RecordReload() },
			wantAction: script.ActionReload,
		},
		{
			name:       "RecordAssertText",
			recordFunc: func(r *Recorder) { r.RecordAssertText("#msg", "Hello") },
			wantAction: script.ActionAssertText,
			validate: func(t *testing.T, step script.Step) {
				if step.Selector != "#msg" {
					t.Errorf("Selector = %q, want %q", step.Selector, "#msg")
				}
				if step.Expected != "Hello" {
					t.Errorf("Expected = %q, want %q", step.Expected, "Hello")
				}
			},
		},
		{
			name:       "RecordAssertElement",
			recordFunc: func(r *Recorder) { r.RecordAssertElement("#element") },
			wantAction: script.ActionAssertElement,
			validate: func(t *testing.T, step script.Step) {
				if step.Selector != "#element" {
					t.Errorf("Selector = %q, want %q", step.Selector, "#element")
				}
			},
		},
		{
			name:       "RecordAssertVisible",
			recordFunc: func(r *Recorder) { r.RecordAssertVisible("#element") },
			wantAction: script.ActionAssertVisible,
			validate: func(t *testing.T, step script.Step) {
				if step.Selector != "#element" {
					t.Errorf("Selector = %q, want %q", step.Selector, "#element")
				}
			},
		},
		{
			name:       "RecordAssertURL",
			recordFunc: func(r *Recorder) { r.RecordAssertURL("https://example.com") },
			wantAction: script.ActionAssertURL,
			validate: func(t *testing.T, step script.Step) {
				if step.Expected != "https://example.com" {
					t.Errorf("Expected = %q, want %q", step.Expected, "https://example.com")
				}
			},
		},
		{
			name:       "RecordAssertTitle",
			recordFunc: func(r *Recorder) { r.RecordAssertTitle("Page Title") },
			wantAction: script.ActionAssertTitle,
			validate: func(t *testing.T, step script.Step) {
				if step.Expected != "Page Title" {
					t.Errorf("Expected = %q, want %q", step.Expected, "Page Title")
				}
			},
		},
		{
			name:       "RecordMouseClick",
			recordFunc: func(r *Recorder) { r.RecordMouseClick(100.5, 200.5) },
			wantAction: script.ActionMouseClick,
			validate: func(t *testing.T, step script.Step) {
				if step.X != 100.5 {
					t.Errorf("X = %v, want 100.5", step.X)
				}
				if step.Y != 200.5 {
					t.Errorf("Y = %v, want 200.5", step.Y)
				}
			},
		},
		{
			name:       "RecordMouseMove",
			recordFunc: func(r *Recorder) { r.RecordMouseMove(50, 75) },
			wantAction: script.ActionMouseMove,
			validate: func(t *testing.T, step script.Step) {
				if step.X != 50 {
					t.Errorf("X = %v, want 50", step.X)
				}
				if step.Y != 75 {
					t.Errorf("Y = %v, want 75", step.Y)
				}
			},
		},
		{
			name:       "RecordKeyboardPress",
			recordFunc: func(r *Recorder) { r.RecordKeyboardPress("Enter") },
			wantAction: script.ActionKeyboardPress,
			validate: func(t *testing.T, step script.Step) {
				if step.Key != "Enter" {
					t.Errorf("Key = %q, want %q", step.Key, "Enter")
				}
			},
		},
		{
			name:       "RecordKeyboardType",
			recordFunc: func(r *Recorder) { r.RecordKeyboardType("hello world") },
			wantAction: script.ActionKeyboardType,
			validate: func(t *testing.T, step script.Step) {
				if step.Text != "hello world" {
					t.Errorf("Text = %q, want %q", step.Text, "hello world")
				}
			},
		},
		{
			name:       "RecordDragTo",
			recordFunc: func(r *Recorder) { r.RecordDragTo("#source", "#target") },
			wantAction: script.ActionDragTo,
			validate: func(t *testing.T, step script.Step) {
				if step.Selector != "#source" {
					t.Errorf("Selector = %q, want %q", step.Selector, "#source")
				}
				if step.Target != "#target" {
					t.Errorf("Target = %q, want %q", step.Target, "#target")
				}
			},
		},
		{
			name:       "RecordTap",
			recordFunc: func(r *Recorder) { r.RecordTap("#element") },
			wantAction: script.ActionTap,
			validate: func(t *testing.T, step script.Step) {
				if step.Selector != "#element" {
					t.Errorf("Selector = %q, want %q", step.Selector, "#element")
				}
			},
		},
		{
			name:       "RecordSetFiles",
			recordFunc: func(r *Recorder) { r.RecordSetFiles("#upload", []string{"file1.txt", "file2.txt"}) },
			wantAction: script.ActionSetFiles,
			validate: func(t *testing.T, step script.Step) {
				if step.Selector != "#upload" {
					t.Errorf("Selector = %q, want %q", step.Selector, "#upload")
				}
				if len(step.Files) != 2 {
					t.Errorf("Files length = %d, want 2", len(step.Files))
				}
			},
		},
		{
			name:       "RecordAccessibilityCheck",
			recordFunc: func(r *Recorder) { r.RecordAccessibilityCheck("wcag22aa", "serious") },
			wantAction: script.ActionAssertAccessibility,
			validate: func(t *testing.T, step script.Step) {
				if step.A11y == nil {
					t.Fatal("A11y should not be nil")
				}
				if step.A11y.Standard != "wcag22aa" {
					t.Errorf("A11y.Standard = %q, want %q", step.A11y.Standard, "wcag22aa")
				}
				if step.A11y.FailOn != "serious" {
					t.Errorf("A11y.FailOn = %q, want %q", step.A11y.FailOn, "serious")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := NewRecorder()
			r.Start(RecorderMetadata{})
			tt.recordFunc(r)

			steps := r.Steps()
			if len(steps) != 1 {
				t.Fatalf("Expected 1 step, got %d", len(steps))
			}

			step := steps[0]
			if step.Action != tt.wantAction {
				t.Errorf("Action = %v, want %v", step.Action, tt.wantAction)
			}

			if tt.validate != nil {
				tt.validate(t, step)
			}
		})
	}
}

func TestRecorderConcurrency(t *testing.T) {
	r := NewRecorder()
	r.Start(RecorderMetadata{})

	done := make(chan struct{})
	for i := 0; i < 10; i++ {
		go func() {
			for j := 0; j < 100; j++ {
				r.AddStep(script.Step{Action: script.ActionClick})
				_ = r.Steps()
				_ = r.StepCount()
				_ = r.IsRecording()
			}
			done <- struct{}{}
		}()
	}

	for i := 0; i < 10; i++ {
		<-done
	}

	// Should have 1000 steps
	if r.StepCount() != 1000 {
		t.Errorf("StepCount() = %d, want 1000", r.StepCount())
	}
}
