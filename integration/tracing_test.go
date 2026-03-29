//go:build integration

package integration

import (
	"testing"
)

// TODO: Tracing tests require vibium:tracing.* commands which are not implemented in clicker.
// These tests are skipped until clicker adds support for vibium:tracing.* commands.

// TestTracingBasic tests basic tracing functionality.
func TestTracingBasic(t *testing.T) {
	t.Skip("Tracing requires vibium:tracing.* commands which are not implemented in clicker")
}

// TestTracingChunks tests trace chunk functionality.
func TestTracingChunks(t *testing.T) {
	t.Skip("Tracing requires vibium:tracing.* commands which are not implemented in clicker")
}

// TestTracingGroups tests trace group functionality.
func TestTracingGroups(t *testing.T) {
	t.Skip("Tracing requires vibium:tracing.* commands which are not implemented in clicker")
}

// TestTracingFromContext tests tracing via BrowserContext.
func TestTracingFromContext(t *testing.T) {
	t.Skip("Tracing requires vibium:tracing.* commands which are not implemented in clicker")
}
