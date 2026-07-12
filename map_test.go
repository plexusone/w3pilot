package w3pilot

import (
	"testing"
)

func TestDiffRefs_NoChanges(t *testing.T) {
	refs := []ElementRef{
		{Ref: "@e1", Tag: "button", Text: "Submit", Selector: "#submit"},
		{Ref: "@e2", Tag: "input", Placeholder: "Email", Selector: "#email"},
	}

	diff := DiffRefs(refs, refs)

	if diff.HasChanges() {
		t.Error("expected no changes")
	}
	if diff.Summary.Unchanged != 2 {
		t.Errorf("expected 2 unchanged, got %d", diff.Summary.Unchanged)
	}
	if diff.Summary.Added != 0 {
		t.Errorf("expected 0 added, got %d", diff.Summary.Added)
	}
	if diff.Summary.Removed != 0 {
		t.Errorf("expected 0 removed, got %d", diff.Summary.Removed)
	}
}

func TestDiffRefs_Added(t *testing.T) {
	before := []ElementRef{
		{Ref: "@e1", Tag: "button", Text: "Submit", Selector: "#submit"},
	}
	after := []ElementRef{
		{Ref: "@e1", Tag: "button", Text: "Submit", Selector: "#submit"},
		{Ref: "@e2", Tag: "button", Text: "Cancel", Selector: "#cancel"},
	}

	diff := DiffRefs(before, after)

	if !diff.HasChanges() {
		t.Error("expected changes")
	}
	if diff.Summary.Added != 1 {
		t.Errorf("expected 1 added, got %d", diff.Summary.Added)
	}
	if len(diff.Added) != 1 {
		t.Fatalf("expected 1 added element, got %d", len(diff.Added))
	}
	if diff.Added[0].Text != "Cancel" {
		t.Errorf("expected added element to have text 'Cancel', got %q", diff.Added[0].Text)
	}
}

func TestDiffRefs_Removed(t *testing.T) {
	before := []ElementRef{
		{Ref: "@e1", Tag: "button", Text: "Submit", Selector: "#submit"},
		{Ref: "@e2", Tag: "button", Text: "Cancel", Selector: "#cancel"},
	}
	after := []ElementRef{
		{Ref: "@e1", Tag: "button", Text: "Submit", Selector: "#submit"},
	}

	diff := DiffRefs(before, after)

	if !diff.HasChanges() {
		t.Error("expected changes")
	}
	if diff.Summary.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", diff.Summary.Removed)
	}
	if len(diff.Removed) != 1 {
		t.Fatalf("expected 1 removed element, got %d", len(diff.Removed))
	}
	if diff.Removed[0].Text != "Cancel" {
		t.Errorf("expected removed element to have text 'Cancel', got %q", diff.Removed[0].Text)
	}
}

func TestDiffRefs_Changed(t *testing.T) {
	before := []ElementRef{
		{Ref: "@e1", Tag: "button", Text: "Submit", Selector: "#submit"},
	}
	after := []ElementRef{
		{Ref: "@e1", Tag: "button", Text: "Submit", Selector: "#new-submit"}, // Same element, different selector
	}

	diff := DiffRefs(before, after)

	if !diff.HasChanges() {
		t.Error("expected changes")
	}
	if diff.Summary.Changed != 1 {
		t.Errorf("expected 1 changed, got %d", diff.Summary.Changed)
	}
	if len(diff.Changed) != 1 {
		t.Fatalf("expected 1 changed element, got %d", len(diff.Changed))
	}
	if diff.Changed[0].Before.Selector != "#submit" {
		t.Errorf("expected before selector '#submit', got %q", diff.Changed[0].Before.Selector)
	}
	if diff.Changed[0].After.Selector != "#new-submit" {
		t.Errorf("expected after selector '#new-submit', got %q", diff.Changed[0].After.Selector)
	}
}

func TestDiffRefs_MixedChanges(t *testing.T) {
	before := []ElementRef{
		{Ref: "@e1", Tag: "button", Text: "Submit", Selector: "#submit"},
		{Ref: "@e2", Tag: "button", Text: "Cancel", Selector: "#cancel"},
		{Ref: "@e3", Tag: "input", Placeholder: "Email", Selector: "#email"},
	}
	after := []ElementRef{
		{Ref: "@e1", Tag: "button", Text: "Submit", Selector: "#submit"},     // Unchanged
		{Ref: "@e2", Tag: "input", Placeholder: "Email", Selector: "#email"}, // Unchanged (different ref number)
		{Ref: "@e3", Tag: "button", Text: "Next", Selector: "#next"},         // Added
	}

	diff := DiffRefs(before, after)

	if !diff.HasChanges() {
		t.Error("expected changes")
	}
	if diff.Summary.Unchanged != 2 {
		t.Errorf("expected 2 unchanged, got %d", diff.Summary.Unchanged)
	}
	if diff.Summary.Removed != 1 {
		t.Errorf("expected 1 removed (Cancel), got %d", diff.Summary.Removed)
	}
	if diff.Summary.Added != 1 {
		t.Errorf("expected 1 added (Next), got %d", diff.Summary.Added)
	}
}

func TestDiffRefs_MatchesByFingerprint(t *testing.T) {
	// Elements should match by content, not by ref number
	before := []ElementRef{
		{Ref: "@e1", Tag: "button", Text: "Submit", Selector: "#submit"},
	}
	after := []ElementRef{
		{Ref: "@e5", Tag: "button", Text: "Submit", Selector: "#submit"}, // Same element, different ref number
	}

	diff := DiffRefs(before, after)

	if diff.HasChanges() {
		t.Error("expected no changes (same element, different ref number)")
	}
	if diff.Summary.Unchanged != 1 {
		t.Errorf("expected 1 unchanged, got %d", diff.Summary.Unchanged)
	}
}

func TestDiffRefs_EmptyBefore(t *testing.T) {
	before := []ElementRef{}
	after := []ElementRef{
		{Ref: "@e1", Tag: "button", Text: "Submit", Selector: "#submit"},
	}

	diff := DiffRefs(before, after)

	if !diff.HasChanges() {
		t.Error("expected changes")
	}
	if diff.Summary.Added != 1 {
		t.Errorf("expected 1 added, got %d", diff.Summary.Added)
	}
}

func TestDiffRefs_EmptyAfter(t *testing.T) {
	before := []ElementRef{
		{Ref: "@e1", Tag: "button", Text: "Submit", Selector: "#submit"},
	}
	after := []ElementRef{}

	diff := DiffRefs(before, after)

	if !diff.HasChanges() {
		t.Error("expected changes")
	}
	if diff.Summary.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", diff.Summary.Removed)
	}
}

func TestDiffRefs_InputTypeDistinguished(t *testing.T) {
	// Inputs with different types should be considered different elements
	before := []ElementRef{
		{Ref: "@e1", Tag: "input", Placeholder: "Value", Type: "text", Selector: "#text-input"},
	}
	after := []ElementRef{
		{Ref: "@e1", Tag: "input", Placeholder: "Value", Type: "password", Selector: "#password-input"},
	}

	diff := DiffRefs(before, after)

	if !diff.HasChanges() {
		t.Error("expected changes (different input types)")
	}
	// text input was removed, password input was added
	if diff.Summary.Removed != 1 {
		t.Errorf("expected 1 removed, got %d", diff.Summary.Removed)
	}
	if diff.Summary.Added != 1 {
		t.Errorf("expected 1 added, got %d", diff.Summary.Added)
	}
}

func TestRefFingerprint(t *testing.T) {
	tests := []struct {
		name     string
		ref      ElementRef
		contains string
	}{
		{
			name:     "button with text",
			ref:      ElementRef{Tag: "button", Text: "Submit"},
			contains: "button:Submit",
		},
		{
			name:     "input with label",
			ref:      ElementRef{Tag: "input", Label: "Email"},
			contains: "input:label=Email",
		},
		{
			name:     "input with placeholder",
			ref:      ElementRef{Tag: "input", Placeholder: "Enter email"},
			contains: "input:placeholder=Enter email",
		},
		{
			name:     "input with type",
			ref:      ElementRef{Tag: "input", Placeholder: "Password", Type: "password"},
			contains: "type=password",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fp := refFingerprint(tt.ref)
			if !containsStr(fp, tt.contains) {
				t.Errorf("fingerprint %q should contain %q", fp, tt.contains)
			}
		})
	}
}
