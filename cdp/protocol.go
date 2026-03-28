// Package cdp provides a Chrome DevTools Protocol client.
package cdp

import "encoding/json"

// Message represents a CDP message (request or response).
type Message struct {
	ID     int64           `json:"id,omitempty"`
	Method string          `json:"method,omitempty"`
	Params json.RawMessage `json:"params,omitempty"`
	Result json.RawMessage `json:"result,omitempty"`
	Error  *Error          `json:"error,omitempty"`
}

// Error represents a CDP error response.
type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    string `json:"data,omitempty"`
}

func (e *Error) Error() string {
	if e.Data != "" {
		return e.Message + ": " + e.Data
	}
	return e.Message
}

// Event represents a CDP event.
type Event struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

// EventHandler is a callback for CDP events.
type EventHandler func(params json.RawMessage)

// Common CDP domains and methods used by WebPilot.
const (
	// HeapProfiler domain
	HeapProfilerEnable             = "HeapProfiler.enable"
	HeapProfilerDisable            = "HeapProfiler.disable"
	HeapProfilerTakeHeapSnapshot   = "HeapProfiler.takeHeapSnapshot"
	HeapProfilerAddHeapSnapshotChunk = "HeapProfiler.addHeapSnapshotChunk"

	// Network domain
	NetworkEnable          = "Network.enable"
	NetworkDisable         = "Network.disable"
	NetworkGetResponseBody = "Network.getResponseBody"
	NetworkEmulateConditions = "Network.emulateNetworkConditions"

	// Emulation domain
	EmulationSetCPUThrottlingRate = "Emulation.setCPUThrottlingRate"

	// Profiler domain (for coverage)
	ProfilerEnable                = "Profiler.enable"
	ProfilerDisable               = "Profiler.disable"
	ProfilerStartPreciseCoverage  = "Profiler.startPreciseCoverage"
	ProfilerStopPreciseCoverage   = "Profiler.stopPreciseCoverage"
	ProfilerTakePreciseCoverage   = "Profiler.takePreciseCoverage"

	// Debugger domain (for source maps)
	DebuggerEnable  = "Debugger.enable"
	DebuggerDisable = "Debugger.disable"

	// Runtime domain
	RuntimeEnable  = "Runtime.enable"
	RuntimeDisable = "Runtime.disable"

	// Log domain
	LogEnable  = "Log.enable"
	LogDisable = "Log.disable"
)

// NetworkConditions represents network throttling settings.
type NetworkConditions struct {
	Offline            bool    `json:"offline"`
	Latency            float64 `json:"latency"`              // ms
	DownloadThroughput float64 `json:"downloadThroughput"`   // bytes/sec
	UploadThroughput   float64 `json:"uploadThroughput"`     // bytes/sec
	ConnectionType     string  `json:"connectionType,omitempty"`
}

// Preset network conditions.
var (
	NetworkOffline = NetworkConditions{
		Offline:            true,
		Latency:            0,
		DownloadThroughput: 0,
		UploadThroughput:   0,
	}

	NetworkSlow3G = NetworkConditions{
		Offline:            false,
		Latency:            400,
		DownloadThroughput: 400 * 1024 / 8,  // 400 Kbps
		UploadThroughput:   400 * 1024 / 8,
		ConnectionType:     "cellular3g",
	}

	NetworkFast3G = NetworkConditions{
		Offline:            false,
		Latency:            150,
		DownloadThroughput: 1500 * 1024 / 8, // 1.5 Mbps
		UploadThroughput:   750 * 1024 / 8,
		ConnectionType:     "cellular3g",
	}

	Network4G = NetworkConditions{
		Offline:            false,
		Latency:            50,
		DownloadThroughput: 4000 * 1024 / 8, // 4 Mbps
		UploadThroughput:   3000 * 1024 / 8,
		ConnectionType:     "cellular4g",
	}

	NetworkWifi = NetworkConditions{
		Offline:            false,
		Latency:            10,
		DownloadThroughput: 30000 * 1024 / 8, // 30 Mbps
		UploadThroughput:   15000 * 1024 / 8,
		ConnectionType:     "wifi",
	}
)

// GetResponseBodyResult is the result of Network.getResponseBody.
type GetResponseBodyResult struct {
	Body          string `json:"body"`
	Base64Encoded bool   `json:"base64Encoded"`
}

// HeapSnapshotChunk is sent during heap snapshot capture.
type HeapSnapshotChunk struct {
	Chunk string `json:"chunk"`
}
