# CLI Reference

The `vibium` CLI provides command-line browser automation.

## Installation

```bash
go install github.com/grokify/w3pilot/cmd/vibium@latest
```

## Global Flags

| Flag | Description |
|------|-------------|
| `--session` | Session file path (default: `~/.vibium/session.json`) |
| `-v, --verbose` | Verbose output |

## Commands

### launch

Launch a browser instance.

```bash
w3pilot launch [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--headless` | Run in headless mode |

**Example:**

```bash
w3pilot launch --headless
```

### go

Navigate to a URL.

```bash
w3pilot go <url> [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--timeout` | Navigation timeout (default: 30s) |

**Example:**

```bash
w3pilot go https://example.com
```

### click

Click an element.

```bash
w3pilot click <selector> [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--timeout` | Timeout (default: 10s) |

**Example:**

```bash
w3pilot click "#submit"
w3pilot click "button.login"
```

### type

Type text into an element (appends).

```bash
vibium type <selector> <text> [flags]
```

**Example:**

```bash
vibium type "#search" "hello world"
```

### fill

Fill an input (replaces existing content).

```bash
w3pilot fill <selector> <text> [flags]
```

**Example:**

```bash
w3pilot fill "#email" "user@example.com"
```

### screenshot

Capture a screenshot.

```bash
w3pilot screenshot <filename> [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--selector` | Capture specific element |
| `--timeout` | Timeout (default: 30s) |

**Example:**

```bash
w3pilot screenshot page.png
w3pilot screenshot button.png --selector "#submit"
```

### eval

Execute JavaScript.

```bash
vibium eval <javascript> [flags]
```

**Example:**

```bash
vibium eval "document.title"
vibium eval "document.querySelectorAll('a').length"
```

### quit

Close the browser.

```bash
w3pilot quit
```

### mcp

Start MCP server.

```bash
w3pilot mcp [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--headless` | Run headless |
| `--timeout` | Default timeout |
| `--project` | Project name for reports |

### run

Run a YAML/JSON script.

```bash
w3pilot run <script> [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--headless` | Run headless |
| `--timeout` | Total script timeout |

**Example:**

```bash
w3pilot run test.yaml
w3pilot run login.json --headless
```

## Session Management

The CLI maintains session state in `~/.vibium/session.json`. This allows running commands across multiple invocations:

```bash
w3pilot launch
w3pilot go https://example.com
# ... later ...
w3pilot screenshot result.png
w3pilot quit
```

## Examples

### Login Flow

```bash
w3pilot launch --headless
w3pilot go https://example.com/login
w3pilot fill "#email" "user@example.com"
w3pilot fill "#password" "secret123"
w3pilot click "#submit"
w3pilot screenshot dashboard.png
w3pilot quit
```

### Form Automation

```bash
w3pilot launch
w3pilot go https://example.com/form
w3pilot fill "#name" "John Doe"
w3pilot fill "#email" "john@example.com"
w3pilot click "input[type='checkbox']"
w3pilot click "#submit"
w3pilot quit
```
