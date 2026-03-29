package w3pilot

import (
	"testing"
)

func TestMatchURLPattern(t *testing.T) {
	tests := []struct {
		name    string
		url     string
		pattern string
		want    bool
	}{
		// Exact match
		{
			name:    "exact match",
			url:     "https://example.com/page",
			pattern: "https://example.com/page",
			want:    true,
		},
		{
			name:    "exact match fail",
			url:     "https://example.com/page",
			pattern: "https://example.com/other",
			want:    false,
		},

		// Glob patterns
		{
			name:    "glob single wildcard",
			url:     "https://example.com/users/123",
			pattern: "https://example.com/users/*",
			want:    true,
		},
		{
			name:    "glob double wildcard",
			url:     "https://example.com/api/v1/users/123",
			pattern: "https://example.com/**/users/*",
			want:    true,
		},
		{
			name:    "glob wildcard in middle",
			url:     "https://example.com/users/123/profile",
			pattern: "https://example.com/users/*/profile",
			want:    true,
		},
		{
			name:    "glob no match",
			url:     "https://example.com/admin/page",
			pattern: "https://example.com/users/*",
			want:    false,
		},

		// Regex patterns
		{
			name:    "regex match",
			url:     "https://example.com/users/12345",
			pattern: "/https://example\\.com/users/\\d+/",
			want:    true,
		},
		{
			name:    "regex no match",
			url:     "https://example.com/users/abc",
			pattern: "/https://example\\.com/users/\\d+/",
			want:    false,
		},
		{
			name:    "regex partial match",
			url:     "https://example.com/dashboard",
			pattern: "/dashboard/",
			want:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := matchURLPattern(tt.url, tt.pattern)
			if got != tt.want {
				t.Errorf("matchURLPattern(%q, %q) = %v, want %v", tt.url, tt.pattern, got, tt.want)
			}
		})
	}
}

func TestAssertionError(t *testing.T) {
	err := &AssertionError{
		Type:     "TestError",
		Message:  "test message",
		Expected: "expected",
		Actual:   "actual",
	}

	if err.Error() != "test message" {
		t.Errorf("Error() = %q, want %q", err.Error(), "test message")
	}
}
