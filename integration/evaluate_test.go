//go:build integration

package integration

import (
	"testing"
)

// TestEvaluateAsyncIIFE tests that async IIFEs properly return their resolved value.
// This was a bug where async IIFEs returned null because the script wrapping
// used block syntax (losing the return value) when semicolons were present.
func TestEvaluateAsyncIIFE(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	// Navigate to a simple page
	bt.go_(`data:text/html,<!DOCTYPE html><html><body>Test</body></html>`)

	// Test async IIFE with await - this was the reported bug
	result, err := bt.pilot.Evaluate(bt.ctx, `(async () => {
		const data = await Promise.resolve({status: 200, message: "ok"});
		return data;
	})()`)
	if err != nil {
		t.Fatalf("Failed to evaluate async IIFE: %v", err)
	}

	// Result should be the resolved object, not null
	if result == nil {
		t.Fatal("Async IIFE returned nil, expected resolved value")
	}

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected map result, got %T: %v", result, result)
	}

	if status, ok := resultMap["status"].(float64); !ok || status != 200 {
		t.Errorf("Expected status 200, got %v", resultMap["status"])
	}
	if message, ok := resultMap["message"].(string); !ok || message != "ok" {
		t.Errorf("Expected message 'ok', got %v", resultMap["message"])
	}
}

// TestEvaluateSyncIIFE tests that sync IIFEs also work correctly.
func TestEvaluateSyncIIFE(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(`data:text/html,<!DOCTYPE html><html><body>Test</body></html>`)

	// Sync IIFE with semicolons inside
	result, err := bt.pilot.Evaluate(bt.ctx, `(function() {
		const x = 1;
		const y = 2;
		return x + y;
	})()`)
	if err != nil {
		t.Fatalf("Failed to evaluate sync IIFE: %v", err)
	}

	if result == nil {
		t.Fatal("Sync IIFE returned nil")
	}

	// JSON numbers are float64
	if val, ok := result.(float64); !ok || val != 3 {
		t.Errorf("Expected 3, got %v (%T)", result, result)
	}
}

// TestEvaluateArrowIIFE tests arrow function IIFEs.
func TestEvaluateArrowIIFE(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(`data:text/html,<!DOCTYPE html><html><body>Test</body></html>`)

	// Arrow IIFE with semicolons
	result, err := bt.pilot.Evaluate(bt.ctx, `(() => {
		const items = [1, 2, 3];
		return items.reduce((a, b) => a + b, 0);
	})()`)
	if err != nil {
		t.Fatalf("Failed to evaluate arrow IIFE: %v", err)
	}

	if result == nil {
		t.Fatal("Arrow IIFE returned nil")
	}

	if val, ok := result.(float64); !ok || val != 6 {
		t.Errorf("Expected 6, got %v (%T)", result, result)
	}
}

// TestEvaluateSimpleExpression verifies simple expressions still work.
func TestEvaluateSimpleExpression(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(`data:text/html,<!DOCTYPE html><html><body>Test</body></html>`)

	// Simple expression without semicolons
	result, err := bt.pilot.Evaluate(bt.ctx, `1 + 2 + 3`)
	if err != nil {
		t.Fatalf("Failed to evaluate simple expression: %v", err)
	}

	if val, ok := result.(float64); !ok || val != 6 {
		t.Errorf("Expected 6, got %v", result)
	}
}

// TestEvaluateWithReturn verifies explicit return statements work.
func TestEvaluateWithReturn(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(`data:text/html,<!DOCTYPE html><html><body>Test</body></html>`)

	// Script with explicit return
	result, err := bt.pilot.Evaluate(bt.ctx, `return document.body.textContent`)
	if err != nil {
		t.Fatalf("Failed to evaluate with return: %v", err)
	}

	if val, ok := result.(string); !ok || val != "Test" {
		t.Errorf("Expected 'Test', got %v", result)
	}
}

// TestEvaluatePromiseThen verifies .then() chains work (the workaround).
func TestEvaluatePromiseThen(t *testing.T) {
	bt := newBrowserTest(t)
	defer bt.cleanup()

	bt.go_(`data:text/html,<!DOCTYPE html><html><body>Test</body></html>`)

	// Promise with .then() - this was the workaround
	result, err := bt.pilot.Evaluate(bt.ctx, `Promise.resolve(42).then(x => x * 2)`)
	if err != nil {
		t.Fatalf("Failed to evaluate promise.then: %v", err)
	}

	if val, ok := result.(float64); !ok || val != 84 {
		t.Errorf("Expected 84, got %v", result)
	}
}
