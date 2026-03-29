# Architecture Overview

## System Architecture

W3Pilot uses a **dual-protocol architecture** connecting to a single Chrome browser via both WebDriver BiDi and Chrome DevTools Protocol (CDP):

```
┌─────────────────────────────────────────────────────────────────────────┐
│                              User Layer                                  │
├─────────────────┬─────────────────┬─────────────────┬──────────────────┤
│    Go Client    │   MCP Server    │      CLI        │  Script Runner   │
│      SDK        │   (75+ tools)   │    (vibium)     │  (w3pilot run)  │
├─────────────────┴─────────────────┴─────────────────┴──────────────────┤
│                           w3pilot Core                                │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐     │
│  │  Pilot   │ │ Element  │ │ Keyboard │ │  Mouse   │ │  Touch   │     │
│  │ (page)   │ │ (DOM)    │ │ (input)  │ │ (input)  │ │ (input)  │     │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘     │
│  ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐ ┌──────────┐     │
│  │ Context  │ │  Clock   │ │ Tracing  │ │  Route   │ │   CDP    │     │
│  │(session) │ │ (time)   │ │(capture) │ │(network) │ │(profiling)│    │
│  └──────────┘ └──────────┘ └──────────┘ └──────────┘ └──────────┘     │
├─────────────────────────────────────────────────────────────────────────┤
│                      Dual Protocol Layer                                │
│  ┌─────────────────────────────────┐  ┌────────────────────────────┐   │
│  │        BiDi Client              │  │       CDP Client           │   │
│  │  (page automation, DOM, events) │  │ (profiling, emulation)     │   │
│  └─────────────────┬───────────────┘  └─────────────┬──────────────┘   │
│                    │                                │                   │
├────────────────────┼────────────────────────────────┼───────────────────┤
│                    ▼                                ▼                   │
│         VibiumDev Clicker                   Chrome DevTools             │
│     (vibium:* + WebDriver BiDi)            (CDP WebSocket)              │
├─────────────────────────────────────────────────────────────────────────┤
│                    Chrome / Chromium                                    │
│             (single browser instance)                                   │
└─────────────────────────────────────────────────────────────────────────┘
```

### Protocol Responsibilities

| Protocol | Layer | Features |
|----------|-------|----------|
| **WebDriver BiDi** | Via clicker | Page automation, element interactions, screenshots, tracing, events |
| **Chrome DevTools Protocol** | Direct | Heap profiling, network response bodies, CPU/network emulation |

## Component Descriptions

### Go Client SDK

The core programmatic API for browser automation:

- **Pilot**: Page-level operations (navigation, screenshots, JS evaluation)
- **Element**: DOM element interactions (click, type, fill, state queries)
- **Input Controllers**: Low-level keyboard, mouse, touch control
- **Context**: Isolated browser sessions with separate cookies/storage
- **Network**: Request interception and modification
- **CDP**: Direct Chrome DevTools Protocol access for profiling and emulation

### CDP Client

Direct Chrome DevTools Protocol access for advanced features:

- **Heap Profiler**: Capture V8 heap snapshots for memory analysis
- **Network Emulation**: Simulate Slow 3G, Fast 3G, 4G, or custom conditions
- **CPU Emulation**: Throttle CPU for performance testing (2x, 4x, 6x slowdown)
- **Direct Commands**: Send any CDP command for advanced use cases

### MCP Server

Model Context Protocol server for AI assistant integration:

- 75+ browser automation tools
- Session management with test reporting
- Script recording capability
- Structured error messages with suggestions

### CLI

Command-line interface for scripted automation:

- Subcommand structure (`w3pilot launch`, `w3pilot click`, etc.)
- Session persistence between commands
- YAML/JSON script execution

### Script Runner

Deterministic test execution:

- JSON/YAML script format with JSON Schema
- Variable interpolation
- Assertions and data extraction
- Error handling with `continueOnError`

## Data Flow

### MCP Tool Call Flow

```
Claude                    MCP Server              Vibe              Browser
  │                           │                    │                   │
  │──── navigate ────────────▶│                    │                   │
  │                           │──── Go(url) ──────▶│                   │
  │                           │                    │── BiDi request ──▶│
  │                           │                    │◀── BiDi event ────│
  │                           │◀─── url, title ────│                   │
  │◀─── NavigateOutput ───────│                    │                   │
```

### Session Recording Flow

```
Claude                    MCP Server              Recorder
  │                           │                      │
  │── start_recording ───────▶│                      │
  │                           │──── Start() ────────▶│
  │                           │                      │
  │──── navigate ────────────▶│                      │
  │                           │── RecordNavigate() ─▶│
  │                           │                      │
  │──── click ───────────────▶│                      │
  │                           │── RecordClick() ────▶│
  │                           │                      │
  │──── export_script ───────▶│                      │
  │                           │◀── ExportJSON() ────│
  │◀─── JSON script ──────────│                      │
```

## Feature Origin

| Component | Origin | Notes |
|-----------|--------|-------|
| BiDi client | Upstream | WebDriver BiDi protocol |
| Vibe API | Upstream | Parity with JS/Python |
| Element API | Upstream | Parity with JS/Python |
| Input controllers | Upstream | Parity with JS/Python |
| MCP server | Go-specific | AI assistant integration |
| CLI | Go-specific | Command-line automation |
| Script runner | Go-specific | Deterministic replay |
| Session recording | Go-specific | LLM action capture |
| JSON Schema | Go-specific | Script validation |
| Test reporting | Go-specific | Structured diagnostics |

## Key Design Decisions

### Dual Protocol Architecture

W3Pilot uses **both** WebDriver BiDi and Chrome DevTools Protocol (CDP):

**WebDriver BiDi (via VibiumDev clicker)** for:

- Standardization across browsers
- Bidirectional events (no polling)
- Future-proof design
- Page automation (navigation, DOM, interactions)

**Chrome DevTools Protocol (CDP)** for:

- Heap profiling (not available in BiDi)
- Network response bodies (not exposed in BiDi)
- CPU/network emulation presets
- Any Chrome-specific DevTools feature

Both protocols connect to the same Chrome browser instance, discovered via the `DevToolsActivePort` file in Chrome's user data directory.

### Protocol-Agnostic Methods

Some SDK methods are **protocol-agnostic** - they try BiDi first and automatically fall back to CDP when BiDi doesn't support the feature:

```
User calls SetOffline(true)
    │
    ▼
┌─────────────────────────────┐
│ Try BiDi command            │
│ vibium:network.setOffline   │
└─────────────┬───────────────┘
              │
              ▼ (fails: "unknown command")
┌─────────────────────────────┐
│ Detect unsupported command  │
│ IsUnsupportedCommand(err)   │
└─────────────┬───────────────┘
              │
              ▼
┌─────────────────────────────┐
│ Fall back to CDP            │
│ EmulateNetwork(Offline)     │
└─────────────┬───────────────┘
              │
              ▼
         Success
```

**Protocol-agnostic methods:**

| Method | BiDi Command | CDP Fallback |
|--------|--------------|--------------|
| `SetOffline()` | `vibium:network.setOffline` | `EmulateNetwork(NetworkOffline)` |
| `ConsoleMessages()` | `vibium:console.messages` | `ConsoleEntries()` |
| `ClearConsoleMessages()` | `vibium:console.clear` | `consoleDebugger.Clear()` |

This design ensures:

- **Stable API** - Users don't need to change code when BiDi support is added
- **Best available backend** - Automatically uses the most capable implementation
- **Graceful degradation** - Clear errors when neither protocol supports the feature

### Custom Commands

W3Pilot extends BiDi with `vibium:*` commands for:

- High-level actions (fill, check, selectOption)
- Actionability checks (wait for visible, enabled, stable)
- Page-level operations (screenshot, PDF, evaluate)

### MCP Architecture

The MCP server uses the Model Context Protocol for:

- Standardized tool definitions
- Structured input/output
- Easy AI assistant integration

### Session Recording

Recording captures tool calls (not raw BiDi) for:

- Human-readable scripts
- Portability (same format as CLI scripts)
- Easy editing and customization
