# Installation

## Go Client SDK

Install the Go module:

```bash
go get github.com/plexusone/w3pilot
```

## CLI Tool

Build and install the CLI:

```bash
go install github.com/plexusone/w3pilot/cmd/w3pilot@latest
```

Or build from source:

```bash
git clone https://github.com/plexusone/w3pilot
cd w3pilot
go build -o w3pilot ./cmd/w3pilot
```

## MCP Server

Install the standalone MCP server:

```bash
go install github.com/plexusone/w3pilot/cmd/w3pilot-mcp@latest
```

## Prerequisites

### VibiumDev Clicker Binary

W3Pilot requires the VibiumDev clicker binary for WebDriver BiDi communication. Download from [VibiumDev releases](https://github.com/VibiumDev/vibium/releases).

Alternatively, set the path manually:

```bash
export W3PILOT_CLICKER_PATH=/path/to/clicker
```

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `W3PILOT_CLICKER_PATH` | Path to clicker binary | Auto-detected |
| `W3PILOT_DEBUG` | Enable debug logging | `false` |
| `W3PILOT_HEADLESS` | Run headless by default | `false` |

## Verify Installation

```bash
# Check CLI
w3pilot --help

# Check MCP server
w3pilot-mcp --list-tools | head -5

# Launch browser
w3pilot browser launch
```
