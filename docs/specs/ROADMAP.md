# w3pilot Roadmap

This document tracks planned features and enhancements for w3pilot.

## Current Status (v0.8.0)

w3pilot provides browser automation via WebDriver BiDi and Chrome DevTools Protocol (CDP).

### Completed Features

- [x] BiDi protocol support for browser automation
- [x] CDP integration for profiling and debugging
- [x] 180 MCP tools across 25 namespaces
- [x] Session management with reconnection
- [x] Network interception and mocking
- [x] Heap snapshots and memory profiling
- [x] Network emulation (3G, 4G, etc.)
- [x] CPU throttling
- [x] Code coverage analysis
- [x] Video recording
- [x] Tracing with DOM snapshots
- [x] RPA workflow engine (JSON/YAML scripts)
- [x] Storage state save/restore
- [x] Element mapping (@refs) for AI agents
- [x] Fake/controlled clock for time manipulation
- [x] Geolocation emulation
- [x] Touch emulation (tap, swipe)
- [x] Drag and drop operations

---

## v0.9.0 - Vibium Feature Parity

**Goal**: Add key features from Vibium that enhance AI agent workflows.

### Element Mapping (@refs)

Human-friendly element references for AI agents and interactive use.

**Current state**: Elements are referenced by CSS selectors or XPath.

**Proposed**: Add `map` command that labels interactive elements as @e1, @e2, etc.

```bash
# Map all interactive elements on page
w3pilot map

# Output:
# @e1 [button] "Submit Form"
# @e2 [button] "Cancel"
# @e3 [input]  "Username" (placeholder)
# @e4 [input]  "Password" (placeholder)
# @e5 [a]      "Forgot password?"

# Use refs in subsequent commands
w3pilot click @e1
w3pilot fill @e3 "john@example.com"
```

| Task | Description | Status |
|------|-------------|--------|
| Element mapper | Scan page for interactive elements | ✅ Done |
| Ref storage | Store element refs in session state | ✅ Done |
| CLI `map` command | Map elements and display refs | ✅ Done |
| CLI ref resolution | Resolve @refs in click, fill, etc. | ✅ Done |
| MCP `page_map` tool | Map elements via MCP | ✅ Done |
| Diff mapping | Compare refs between snapshots | ✅ Done |

### Fake/Controlled Clock ✅ (Already Implemented)

Manipulate page time for testing time-dependent functionality.

> **Note**: This feature is already implemented in `clock.go`. The Clock API provides Install, FastForward, RunFor, PauseAt, Resume, SetFixedTime, SetSystemTime, and SetTimezone methods.

```bash
# Install fake clock
w3pilot clock install

# Set specific time
w3pilot clock set "2024-12-25T00:00:00Z"

# Advance time
w3pilot clock tick 60000  # Advance 60 seconds

# Resume real time
w3pilot clock resume
```

| Task | Description | Status |
|------|-------------|--------|
| Clock install | Inject fake timers (Date, setTimeout, etc.) | ✅ Done |
| Clock set | Set fake time to specific value | ✅ Done |
| Clock tick | Advance fake time by milliseconds | ✅ Done |
| Clock resume | Resume real time | ✅ Done |
| MCP `clock_*` tools | Clock manipulation via MCP | ✅ Done |

### Geolocation Emulation ✅ (Already Implemented)

Set geolocation for testing location-based features.

> **Note**: This feature is already implemented via `page_set_geolocation` MCP tool.

```bash
# Set geolocation
w3pilot geolocation set --lat 37.7749 --lng -122.4194

# Set with accuracy
w3pilot geolocation set --lat 37.7749 --lng -122.4194 --accuracy 100

# Clear geolocation override
w3pilot geolocation clear
```

| Task | Description | Status |
|------|-------------|--------|
| Geolocation set | Override geolocation via CDP | ✅ Done |
| Geolocation clear | Clear geolocation override | ✅ Done |
| MCP `page_set_geolocation` | Geolocation via MCP | ✅ Done |

### Drag and Drop ✅ (Already Implemented)

Explicit drag and drop operations.

> **Note**: This feature is already implemented via `element_drag_to` MCP tool and `input_mouse_drag`.

```bash
# Drag element to target
w3pilot drag "#draggable" --to "#dropzone"

# Drag by offset
w3pilot drag "#draggable" --offset-x 100 --offset-y 50
```

| Task | Description | Status |
|------|-------------|--------|
| Drag to element | Drag source to target element | ✅ Done |
| Drag by offset | Drag by pixel offset | ✅ Done |
| MCP `element_drag_to` tool | Drag via MCP | ✅ Done |

### Touch Emulation ✅ (Already Implemented)

Mobile touch event support.

> **Note**: Touch features are already implemented via `element_tap`, `input_touch_tap`, and `input_touch_swipe` MCP tools.

```bash
# Enable touch emulation
w3pilot touch enable

# Tap element
w3pilot tap "#button"

# Swipe gesture
w3pilot swipe --start-x 100 --start-y 200 --end-x 100 --end-y 50

# Pinch zoom
w3pilot pinch --scale 2.0
```

| Task | Description | Status |
|------|-------------|--------|
| Touch enable/disable | Toggle touch emulation | Pending |
| Tap | Single touch tap | ✅ Done |
| Swipe | Swipe gesture | ✅ Done |
| Pinch | Pinch zoom gesture | Pending |
| MCP touch tools | Touch via MCP | ✅ Done |

---

## v0.10.0 - Enhanced Testing

### Assertions

Built-in assertion commands for testing workflows.

```bash
# Assert element exists
w3pilot assert exists "#submit-button"

# Assert text content
w3pilot assert text "#message" "Success!"

# Assert element visible
w3pilot assert visible "#modal"

# Assert element count
w3pilot assert count ".list-item" 5
```

### Visual Regression

Screenshot comparison for visual testing.

```bash
# Capture baseline
w3pilot visual baseline "#component" --name "header"

# Compare against baseline
w3pilot visual compare "#component" --name "header" --threshold 0.1
```

---

## v0.11.0 - Multi-Browser Support

### Browser Selection

Support for multiple browser engines.

```bash
# Launch Firefox
w3pilot launch --browser firefox

# Launch WebKit (Safari)
w3pilot launch --browser webkit

# Launch Edge
w3pilot launch --browser edge
```

### Parallel Execution

Run operations across multiple browser instances.

```bash
# Launch multiple browsers
w3pilot parallel --browsers "chrome,firefox,webkit" --script test.yaml
```

---

## Implementation Priority

| Version | Feature | Impact | Effort | Status |
|---------|---------|--------|--------|--------|
| v0.9.0 | Element Mapping (@refs) | High | Medium | ✅ Done |
| v0.9.0 | Fake Clock | Medium | Low | ✅ Done |
| v0.9.0 | Geolocation | Medium | Low | ✅ Done (page_set_geolocation) |
| v0.9.0 | Drag and Drop | Low | Low | ✅ Done (element_drag_to) |
| v0.9.0 | Touch Emulation | Medium | Medium | ✅ Done (element_tap, input_touch_*) |
| v0.10.0 | Assertions | High | Medium | Partial (test_assert_*) |
| v0.10.0 | Visual Regression | Medium | High | Pending |
| v0.11.0 | Multi-Browser | High | High | Pending |

---

## Related Documents

- [CLI Reference](../reference/cli.md)
- [MCP Server](../guide/mcp-server.md)
- [CDP Guide](../guide/cdp.md)
