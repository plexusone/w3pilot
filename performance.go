package w3pilot

import (
	"context"
	"encoding/json"
	"fmt"
)

// PerformanceMetrics contains Core Web Vitals and navigation timing.
type PerformanceMetrics struct {
	// Core Web Vitals
	LCP  float64 `json:"lcp,omitempty"`  // Largest Contentful Paint (ms)
	CLS  float64 `json:"cls,omitempty"`  // Cumulative Layout Shift
	INP  float64 `json:"inp,omitempty"`  // Interaction to Next Paint (ms)
	FID  float64 `json:"fid,omitempty"`  // First Input Delay (ms)
	FCP  float64 `json:"fcp,omitempty"`  // First Contentful Paint (ms)
	TTFB float64 `json:"ttfb,omitempty"` // Time to First Byte (ms)

	// Navigation Timing
	DOMContentLoaded float64 `json:"domContentLoaded,omitempty"` // DOMContentLoaded event (ms)
	Load             float64 `json:"load,omitempty"`             // Load event (ms)
	DOMInteractive   float64 `json:"domInteractive,omitempty"`   // DOM interactive (ms)

	// Resource Timing
	ResourceCount int `json:"resourceCount,omitempty"` // Number of resources loaded
}

// MemoryStats contains JavaScript heap memory information.
type MemoryStats struct {
	UsedJSHeapSize  int64 `json:"usedJSHeapSize"`  // Used JS heap size in bytes
	TotalJSHeapSize int64 `json:"totalJSHeapSize"` // Total JS heap size in bytes
	JSHeapSizeLimit int64 `json:"jsHeapSizeLimit"` // JS heap size limit in bytes
}

// GetPerformanceMetrics retrieves Core Web Vitals and navigation timing.
// Note: Some metrics (LCP, CLS, INP) require user interaction or time to measure.
func (p *Pilot) GetPerformanceMetrics(ctx context.Context) (*PerformanceMetrics, error) {
	script := `
(function() {
	const metrics = {};

	// Navigation Timing
	const nav = performance.getEntriesByType('navigation')[0];
	if (nav) {
		metrics.ttfb = nav.responseStart - nav.requestStart;
		metrics.domContentLoaded = nav.domContentLoadedEventEnd - nav.startTime;
		metrics.load = nav.loadEventEnd - nav.startTime;
		metrics.domInteractive = nav.domInteractive - nav.startTime;
	}

	// First Contentful Paint
	const fcpEntry = performance.getEntriesByName('first-contentful-paint')[0];
	if (fcpEntry) {
		metrics.fcp = fcpEntry.startTime;
	}

	// Largest Contentful Paint (if available)
	const lcpEntries = performance.getEntriesByType('largest-contentful-paint');
	if (lcpEntries.length > 0) {
		metrics.lcp = lcpEntries[lcpEntries.length - 1].startTime;
	}

	// Layout Shift (CLS)
	let cls = 0;
	const layoutShiftEntries = performance.getEntriesByType('layout-shift');
	for (const entry of layoutShiftEntries) {
		if (!entry.hadRecentInput) {
			cls += entry.value;
		}
	}
	if (layoutShiftEntries.length > 0) {
		metrics.cls = cls;
	}

	// Resource count
	metrics.resourceCount = performance.getEntriesByType('resource').length;

	return metrics;
})()
`

	result, err := p.Evaluate(ctx, script)
	if err != nil {
		return nil, fmt.Errorf("w3pilot: failed to get performance metrics: %w", err)
	}

	// Convert result to JSON and unmarshal
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("w3pilot: failed to marshal performance metrics: %w", err)
	}

	var metrics PerformanceMetrics
	if err := json.Unmarshal(jsonBytes, &metrics); err != nil {
		return nil, fmt.Errorf("w3pilot: failed to parse performance metrics: %w", err)
	}

	return &metrics, nil
}

// GetMemoryStats retrieves JavaScript heap memory information.
// Note: This requires Chrome with --enable-precise-memory-info flag.
func (p *Pilot) GetMemoryStats(ctx context.Context) (*MemoryStats, error) {
	script := `
(function() {
	if (!performance.memory) {
		return null;
	}
	return {
		usedJSHeapSize: performance.memory.usedJSHeapSize,
		totalJSHeapSize: performance.memory.totalJSHeapSize,
		jsHeapSizeLimit: performance.memory.jsHeapSizeLimit
	};
})()
`

	result, err := p.Evaluate(ctx, script)
	if err != nil {
		return nil, fmt.Errorf("w3pilot: failed to get memory stats: %w", err)
	}

	if result == nil {
		return nil, fmt.Errorf("w3pilot: performance.memory not available (requires --enable-precise-memory-info)")
	}

	// Convert result to JSON and unmarshal
	jsonBytes, err := json.Marshal(result)
	if err != nil {
		return nil, fmt.Errorf("w3pilot: failed to marshal memory stats: %w", err)
	}

	var stats MemoryStats
	if err := json.Unmarshal(jsonBytes, &stats); err != nil {
		return nil, fmt.Errorf("w3pilot: failed to parse memory stats: %w", err)
	}

	return &stats, nil
}

// ObserveWebVitals starts observing Core Web Vitals in real-time.
// Returns a channel that receives metrics as they are measured.
// Call the returned cancel function to stop observing.
func (p *Pilot) ObserveWebVitals(ctx context.Context) (<-chan *PerformanceMetrics, func(), error) {
	// Inject the observer script
	observerScript := `
window.__w3pilotMetrics = window.__w3pilotMetrics || { lcp: 0, cls: 0, inp: 0, fid: 0 };

// LCP Observer
new PerformanceObserver((list) => {
	const entries = list.getEntries();
	const lastEntry = entries[entries.length - 1];
	window.__w3pilotMetrics.lcp = lastEntry.startTime;
}).observe({ type: 'largest-contentful-paint', buffered: true });

// CLS Observer
new PerformanceObserver((list) => {
	for (const entry of list.getEntries()) {
		if (!entry.hadRecentInput) {
			window.__w3pilotMetrics.cls += entry.value;
		}
	}
}).observe({ type: 'layout-shift', buffered: true });

// FID Observer
new PerformanceObserver((list) => {
	const firstEntry = list.getEntries()[0];
	if (firstEntry) {
		window.__w3pilotMetrics.fid = firstEntry.processingStart - firstEntry.startTime;
	}
}).observe({ type: 'first-input', buffered: true });

// INP Observer (via event timing)
let maxINP = 0;
new PerformanceObserver((list) => {
	for (const entry of list.getEntries()) {
		if (entry.interactionId) {
			const duration = entry.duration;
			if (duration > maxINP) {
				maxINP = duration;
				window.__w3pilotMetrics.inp = duration;
			}
		}
	}
}).observe({ type: 'event', buffered: true, durationThreshold: 16 });

true
`

	if _, err := p.Evaluate(ctx, observerScript); err != nil {
		return nil, nil, fmt.Errorf("w3pilot: failed to start web vitals observer: %w", err)
	}

	// Create channel and start polling
	ch := make(chan *PerformanceMetrics, 10)
	done := make(chan struct{})

	go func() {
		defer close(ch)
		for {
			select {
			case <-done:
				return
			case <-ctx.Done():
				return
			default:
				// Poll metrics
				result, err := p.Evaluate(ctx, "window.__w3pilotMetrics")
				if err == nil && result != nil {
					jsonBytes, err := json.Marshal(result)
					if err == nil {
						var metrics PerformanceMetrics
						if json.Unmarshal(jsonBytes, &metrics) == nil {
							select {
							case ch <- &metrics:
							default: // Don't block if channel is full
							}
						}
					}
				}
				// Wait before next poll
				select {
				case <-done:
					return
				case <-ctx.Done():
					return
				}
			}
		}
	}()

	cancel := func() {
		close(done)
	}

	return ch, cancel, nil
}
