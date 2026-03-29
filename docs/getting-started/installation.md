# Installation

## Go Client SDK

Install the Go module:

```bash
go get github.com/grokify/w3pilot
```

## CLI Tool

Build and install the CLI:

```bash
go install github.com/grokify/w3pilot/cmd/vibium@latest
```

Or build from source:

```bash
git clone https://github.com/grokify/w3pilot
cd w3pilot
go build -o vibium ./cmd/vibium
```

## Prerequisites

### W3Pilot Clicker Binary

The Go client requires the W3Pilot clicker binary. Install via npm:

```bash
npm install -g vibium
```

Or download from [W3Pilot releases](https://github.com/W3PilotDev/vibium/releases).

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `W3PILOT_CLICKER_PATH` | Path to clicker binary | Auto-detected |
| `W3PILOT_DEBUG` | Enable debug logging | `false` |
| `W3PILOT_HEADLESS` | Run headless by default | `false` |

## Verify Installation

```bash
# Check CLI
vibium --help

# Check clicker
w3pilot mcp --help
```
