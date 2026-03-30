# CLI Reference

The `w3pilot` CLI provides command-line browser automation.

## Installation

```bash
go install github.com/plexusone/w3pilot/cmd/w3pilot@latest
```

## Global Flags

| Flag | Description |
|------|-------------|
| `--session` | Session file path (default: `~/.w3pilot/session.json`) |
| `-o, --format` | Output format: `text` or `json` |
| `-v, --verbose` | Verbose output |

## Commands

### browser launch

Launch a browser instance.

```bash
w3pilot browser launch [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--headless` | Run in headless mode |

**Example:**

```bash
w3pilot browser launch --headless
```

### browser quit

Close the browser.

```bash
w3pilot browser quit
```

### page navigate

Navigate to a URL.

```bash
w3pilot page navigate <url> [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--timeout` | Navigation timeout (default: 30s) |

**Example:**

```bash
w3pilot page navigate https://example.com
```

### element click

Click an element.

```bash
w3pilot element click <selector> [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--timeout` | Timeout (default: 10s) |

**Example:**

```bash
w3pilot element click "#submit"
w3pilot element click "button.login"
```

### element type

Type text into an element (appends).

```bash
w3pilot element type <selector> <text> [flags]
```

**Example:**

```bash
w3pilot element type "#search" "hello world"
```

### element fill

Fill an input (replaces existing content).

```bash
w3pilot element fill <selector> <text> [flags]
```

**Example:**

```bash
w3pilot element fill "#email" "user@example.com"
```

### page screenshot

Capture a screenshot.

```bash
w3pilot page screenshot <filename> [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--selector` | Capture specific element |
| `--timeout` | Timeout (default: 30s) |

**Example:**

```bash
w3pilot page screenshot page.png
w3pilot page screenshot button.png --selector "#submit"
```

### js eval

Execute JavaScript.

```bash
w3pilot js eval <javascript> [flags]
```

**Example:**

```bash
w3pilot js eval "document.title"
w3pilot js eval "document.querySelectorAll('a').length"
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
| `--list-tools` | List all tools as JSON |

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

### test commands

Assertions and verifications for testing.

```bash
# Assertions
w3pilot test assert-text "Welcome"
w3pilot test assert-element "#login"
w3pilot test assert-url "**/dashboard"

# Verifications
w3pilot test verify-value "#email" "user@example.com"
w3pilot test verify-text ".heading" "Welcome"
w3pilot test verify-visible "#modal"
w3pilot test verify-enabled "#submit"
w3pilot test verify-checked "#agree"

# Locator generation
w3pilot test generate-locator "#submit" --strategy xpath
```

### state commands

Named state snapshots for saving/restoring browser state.

```bash
# Save current state
w3pilot state save login-session

# List saved states
w3pilot state list

# Load a state
w3pilot state load login-session

# Delete a state
w3pilot state delete login-session
```

### page inspect

Discover interactive elements on the page.

```bash
w3pilot page inspect [flags]
```

**Flags:**

| Flag | Description |
|------|-------------|
| `--type` | Element types: buttons, links, inputs, selects, headings, images |

**Example:**

```bash
w3pilot page inspect
w3pilot page inspect --type buttons
```

## Session Management

The CLI maintains session state in `~/.w3pilot/session.json`. This allows running commands across multiple invocations:

```bash
w3pilot browser launch
w3pilot page navigate https://example.com
# ... later ...
w3pilot page screenshot result.png
w3pilot browser quit
```

## Examples

### Login Flow

```bash
w3pilot browser launch --headless
w3pilot page navigate https://example.com/login
w3pilot element fill "#email" "user@example.com"
w3pilot element fill "#password" "secret123"
w3pilot element click "#submit"
w3pilot page screenshot dashboard.png
w3pilot browser quit
```

### Form Automation

```bash
w3pilot browser launch
w3pilot page navigate https://example.com/form
w3pilot element fill "#name" "John Doe"
w3pilot element fill "#email" "john@example.com"
w3pilot element click "input[type='checkbox']"
w3pilot element click "#submit"
w3pilot browser quit
```

### State Management

```bash
# First session: login and save state
w3pilot browser launch
w3pilot page navigate https://example.com/login
w3pilot element fill "#email" "user@example.com"
w3pilot element fill "#password" "secret123"
w3pilot element click "#submit"
w3pilot state save logged-in
w3pilot browser quit

# Later session: restore state
w3pilot browser launch
w3pilot state load logged-in
w3pilot page navigate https://example.com/dashboard
# Already logged in!
```
