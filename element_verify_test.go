package w3pilot

import (
	"testing"
)

func TestVerificationError(t *testing.T) {
	err := &VerificationError{
		Type:     "VerifyValueFailed",
		Message:  "value mismatch: expected \"foo\", got \"bar\"",
		Selector: "#input",
		Expected: "foo",
		Actual:   "bar",
	}

	if err.Error() != "value mismatch: expected \"foo\", got \"bar\"" {
		t.Errorf("Error() = %q, want %q", err.Error(), "value mismatch: expected \"foo\", got \"bar\"")
	}
}

func TestVerifyTextOptions(t *testing.T) {
	// Test that VerifyTextOptions can be constructed
	opts := &VerifyTextOptions{
		Exact: true,
	}

	if !opts.Exact {
		t.Error("Exact should be true")
	}
}
