package w3pilot

import (
	"testing"
)

func TestIsRef(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"@e1", true},
		{"@e2", true},
		{"@e10", true},
		{"@e100", true},
		{"@e0", true},
		{"#submit", false},
		{".button", false},
		{"button", false},
		{"@button", false},
		{"@e", false},
		{"e1", false},
		{"@E1", false}, // Case sensitive
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := IsRef(tt.input)
			if result != tt.expected {
				t.Errorf("IsRef(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestRefStore_StoreAndGet(t *testing.T) {
	store := NewRefStore()

	refs := []ElementRef{
		{Ref: "@e1", Tag: "button", Text: "Submit", Selector: "#submit"},
		{Ref: "@e2", Tag: "input", Placeholder: "Email", Selector: "#email"},
		{Ref: "@e3", Tag: "a", Text: "Learn more", Selector: "a.learn-more"},
	}

	store.Store(refs)

	// Test Get
	ref, ok := store.Get("@e1")
	if !ok {
		t.Fatal("expected to find @e1")
	}
	if ref.Tag != "button" {
		t.Errorf("ref.Tag = %q, want %q", ref.Tag, "button")
	}
	if ref.Text != "Submit" {
		t.Errorf("ref.Text = %q, want %q", ref.Text, "Submit")
	}

	// Test GetSelector
	selector, ok := store.GetSelector("@e2")
	if !ok {
		t.Fatal("expected to find @e2")
	}
	if selector != "#email" {
		t.Errorf("selector = %q, want %q", selector, "#email")
	}

	// Test not found
	_, ok = store.Get("@e99")
	if ok {
		t.Error("expected @e99 not to be found")
	}

	// Test Count
	if store.Count() != 3 {
		t.Errorf("Count() = %d, want 3", store.Count())
	}
}

func TestRefStore_Clear(t *testing.T) {
	store := NewRefStore()

	refs := []ElementRef{
		{Ref: "@e1", Tag: "button", Selector: "#submit"},
	}
	store.Store(refs)

	if store.Count() != 1 {
		t.Errorf("Count() = %d, want 1", store.Count())
	}

	store.Clear()

	if store.Count() != 0 {
		t.Errorf("Count() = %d, want 0 after Clear()", store.Count())
	}
}

func TestRefStore_ResolveRef(t *testing.T) {
	store := NewRefStore()

	refs := []ElementRef{
		{Ref: "@e1", Tag: "button", Selector: "#submit"},
		{Ref: "@e2", Tag: "input", Selector: "#email"},
	}
	store.Store(refs)

	tests := []struct {
		input       string
		expected    string
		expectError bool
	}{
		{"@e1", "#submit", false},
		{"@e2", "#email", false},
		{"#submit", "#submit", false}, // Regular selector passed through
		{".button", ".button", false}, // Regular selector passed through
		{"button", "button", false},   // Regular selector passed through
		{"@e99", "", true},            // Unknown ref
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result, err := store.ResolveRef(tt.input)
			if tt.expectError {
				if err == nil {
					t.Errorf("expected error for %q", tt.input)
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error for %q: %v", tt.input, err)
				}
				if result != tt.expected {
					t.Errorf("ResolveRef(%q) = %q, want %q", tt.input, result, tt.expected)
				}
			}
		})
	}
}

func TestRefStore_All(t *testing.T) {
	store := NewRefStore()

	refs := []ElementRef{
		{Ref: "@e1", Tag: "button", Selector: "#submit"},
		{Ref: "@e2", Tag: "input", Selector: "#email"},
	}
	store.Store(refs)

	all := store.All()
	if len(all) != 2 {
		t.Errorf("All() returned %d refs, want 2", len(all))
	}

	// Verify refs are present (order not guaranteed)
	found := make(map[string]bool)
	for _, ref := range all {
		found[ref.Ref] = true
	}
	if !found["@e1"] {
		t.Error("@e1 not found in All()")
	}
	if !found["@e2"] {
		t.Error("@e2 not found in All()")
	}
}

func TestRefStore_StoreReplaces(t *testing.T) {
	store := NewRefStore()

	// Store first set
	refs1 := []ElementRef{
		{Ref: "@e1", Tag: "button", Selector: "#old"},
	}
	store.Store(refs1)

	// Store second set (should replace)
	refs2 := []ElementRef{
		{Ref: "@e1", Tag: "input", Selector: "#new"},
		{Ref: "@e2", Tag: "a", Selector: "#link"},
	}
	store.Store(refs2)

	if store.Count() != 2 {
		t.Errorf("Count() = %d, want 2", store.Count())
	}

	ref, _ := store.Get("@e1")
	if ref.Tag != "input" {
		t.Errorf("ref.Tag = %q, want %q", ref.Tag, "input")
	}
	if ref.Selector != "#new" {
		t.Errorf("ref.Selector = %q, want %q", ref.Selector, "#new")
	}
}

func TestElementRef_FormatRef(t *testing.T) {
	tests := []struct {
		name     string
		ref      ElementRef
		contains []string
	}{
		{
			name: "button with text",
			ref: ElementRef{
				Ref:     "@e1",
				Tag:     "button",
				Text:    "Submit",
				Visible: true,
				Enabled: true,
			},
			contains: []string{"@e1", "[button]", `"Submit"`},
		},
		{
			name: "input with placeholder",
			ref: ElementRef{
				Ref:         "@e2",
				Tag:         "input",
				Type:        "email",
				Placeholder: "Enter email",
				Visible:     true,
				Enabled:     true,
			},
			contains: []string{"@e2", "[input]", "Enter email", "type=email"},
		},
		{
			name: "disabled element",
			ref: ElementRef{
				Ref:     "@e3",
				Tag:     "button",
				Text:    "Disabled",
				Visible: true,
				Enabled: false,
			},
			contains: []string{"@e3", "[button]", "[disabled]"},
		},
		{
			name: "hidden element",
			ref: ElementRef{
				Ref:     "@e4",
				Tag:     "div",
				Visible: false,
				Enabled: true,
			},
			contains: []string{"@e4", "[div]", "[hidden]"},
		},
		{
			name: "element with role",
			ref: ElementRef{
				Ref:     "@e5",
				Tag:     "div",
				Role:    "button",
				Text:    "Click me",
				Visible: true,
				Enabled: true,
			},
			contains: []string{"@e5", "[div]", "role=button", "Click me"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.ref.FormatRef()
			for _, s := range tt.contains {
				if !containsStr(result, s) {
					t.Errorf("FormatRef() = %q, should contain %q", result, s)
				}
			}
		})
	}
}

func containsStr(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
