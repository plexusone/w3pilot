# MCP Server

The MCP (Model Context Protocol) server provides **169 browser automation tools across 24 namespaces** for AI assistants like Claude.

## Installation

The MCP server can be run two ways:

1. **Standalone binary** (recommended for MCP clients):

   ```bash
   go install github.com/plexusone/w3pilot/cmd/w3pilot-mcp@latest
   ```

2. **Via the w3pilot CLI**:

   ```bash
   go install github.com/plexusone/w3pilot/cmd/w3pilot@latest
   ```

## Starting the Server

### Standalone Binary

```bash
# Default (headless browser)
w3pilot-mcp

# Visible browser (for debugging)
w3pilot-mcp -headless=false

# Custom timeout
w3pilot-mcp -timeout=60s
```

### Via CLI

```bash
# Default (visible browser)
w3pilot mcp

# Headless mode
w3pilot mcp --headless

# Custom timeout
w3pilot mcp --timeout 60s
```

## Client Configuration

### Claude Desktop

Edit the config file:

- **macOS**: `~/Library/Application Support/Claude/claude_desktop_config.json`
- **Windows**: `%APPDATA%\Claude\claude_desktop_config.json`

```json
{
  "mcpServers": {
    "w3pilot": {
      "command": "w3pilot-mcp",
      "args": ["-headless=false"]
    }
  }
}
```

### Claude Code (CLI)

Add to your Claude Code MCP settings:

```json
{
  "mcpServers": {
    "w3pilot": {
      "command": "w3pilot-mcp",
      "args": ["-headless=false"]
    }
  }
}
```

Or use the CLI command:

```bash
claude mcp add w3pilot w3pilot-mcp -- -headless=false
```

### Kiro CLI

```bash
kiro-cli mcp add --name w3pilot --command w3pilot-mcp --args "-headless=false"
```

### Cursor

Edit `.cursor/mcp.json` in your project or home directory:

```json
{
  "mcpServers": {
    "w3pilot": {
      "command": "w3pilot-mcp",
      "args": ["-headless=false"]
    }
  }
}
```

### Windsurf

Edit the MCP configuration in Windsurf settings:

```json
{
  "mcpServers": {
    "w3pilot": {
      "command": "w3pilot-mcp",
      "args": ["-headless=false"]
    }
  }
}
```

### Generic MCP Client

For any MCP-compatible client, use:

- **Command**: `w3pilot-mcp`
- **Args**: `["-headless=false"]` (visible browser) or `[]` (headless)

## Command-Line Options

| Option | Default | Description |
|--------|---------|-------------|
| `-headless` | `true` | Run browser without GUI |
| `-project` | `"w3pilot-tests"` | Project name for reports |
| `-timeout` | `30s` | Default timeout for operations |
| `-init-script` | | JavaScript file to inject before page scripts (repeatable) |
| `--list-tools` | | Export all tools as JSON and exit |

### Init Scripts

Inject JavaScript that runs before any page scripts on every navigation:

```bash
w3pilot-mcp -init-script=./mock-api.js -init-script=./test-helpers.js
```

Use cases:

- Mock APIs before page loads
- Disable analytics/tracking
- Inject test utilities
- Set up authentication tokens

## Environment Variables

| Variable | Description |
|----------|-------------|
| `W3PILOT_DEBUG` | Enable debug logging |
| `W3PILOT_CLICKER_PATH` | Path to clicker binary |

## Tool Categories

See [MCP Tools Reference](../reference/mcp-tools.md) for complete documentation.

### Browser Management

| Tool | Description |
|------|-------------|
| `browser_launch` | Launch browser instance |
| `browser_quit` | Close browser |

### Navigation

| Tool | Description |
|------|-------------|
| `page_navigate` | Go to URL |
| `page_go_back` | Navigate back |
| `page_go_forward` | Navigate forward |
| `page_reload` | Reload page |
| `page_scroll` | Scroll page or element |

### Element Interactions

| Tool | Description |
|------|-------------|
| `element_click` | Click element |
| `element_double_click` | Double-click element |
| `element_type` | Type text (append) |
| `element_fill` | Fill input (replace) |
| `element_clear` | Clear input |
| `element_press` | Press key on element |
| `element_hover` | Hover over element |
| `element_focus` | Focus element |

### Form Controls

| Tool | Description |
|------|-------------|
| `element_check` | Check checkbox |
| `element_uncheck` | Uncheck checkbox |
| `element_select` | Select dropdown option |
| `element_set_files` | Set file input |

### Element State

| Tool | Description |
|------|-------------|
| `element_get_text` | Get element text |
| `element_get_value` | Get input value |
| `element_get_attribute` | Get attribute |
| `element_is_visible` | Check visibility |
| `element_is_enabled` | Check enabled state |
| `element_is_checked` | Check checkbox state |

### Page State

| Tool | Description |
|------|-------------|
| `page_get_title` | Get page title |
| `page_get_url` | Get current URL |
| `page_get_content` | Get page HTML |
| `page_screenshot` | Capture screenshot |
| `page_pdf` | Generate PDF |
| `page_inspect` | Discover interactive elements |

### JavaScript

| Tool | Description |
|------|-------------|
| `js_evaluate` | Execute JavaScript (with optional `max_result_size` for truncation) |
| `js_add_script` | Inject script tag |
| `js_add_style` | Inject CSS |
| `js_init_script` | Add init script for all navigations |

### HTTP Requests

| Tool | Description |
|------|-------------|
| `http_request` | Make authenticated HTTP request in browser context |

### Batch Execution

| Tool | Description |
|------|-------------|
| `batch_execute` | Execute multiple tools in a single call |

### Waiting

| Tool | Description |
|------|-------------|
| `wait_for_selector` | Wait for element |
| `wait_for_state` | Wait for element state |
| `wait_for_url` | Wait for URL pattern |
| `wait_for_load` | Wait for load state |
| `wait_for_text` | Wait for text on page |

### Human-in-the-Loop

| Tool | Description |
|------|-------------|
| `human_pause` | Pause for human action (SSO, CAPTCHA) |

### State Management

| Tool | Description |
|------|-------------|
| `state_save` | Save browser state to named snapshot |
| `state_load` | Restore browser state from snapshot |
| `state_list` | List saved state snapshots |
| `state_delete` | Delete a state snapshot |

### Storage

| Tool | Description |
|------|-------------|
| `storage_get_state` | Export session (cookies + localStorage) |
| `storage_set_state` | Restore saved session |
| `storage_get_cookies` | Get cookies |
| `storage_set_cookies` | Set cookies |
| `storage_local_get` | Get localStorage item |
| `storage_local_set` | Set localStorage item |

### Input Controllers

| Tool | Description |
|------|-------------|
| `input_keyboard_press` | Press key |
| `input_keyboard_type` | Type text |
| `input_mouse_click` | Click at coordinates |
| `input_mouse_move` | Move mouse |

### Script Recording

| Tool | Description |
|------|-------------|
| `record_start` | Begin recording |
| `record_stop` | End recording |
| `record_export` | Export as JSON |
| `record_get_status` | Check status |

### Tracing

| Tool | Description |
|------|-------------|
| `trace_start` | Start trace with screenshots/snapshots |
| `trace_stop` | Stop and save/return trace ZIP |
| `trace_chunk_start` | Start a trace segment |
| `trace_chunk_stop` | Stop trace segment |
| `trace_group_start` | Group actions logically |
| `trace_group_stop` | End action group |

### Testing

| Tool | Description |
|------|-------------|
| `test_assert_text` | Assert text exists |
| `test_assert_element` | Assert element exists |
| `test_assert_url` | Assert URL matches |
| `test_verify_value` | Verify input value |
| `test_verify_visible` | Verify element visible |
| `test_generate_locator` | Generate robust locator |

### Workflow Automation

| Tool | Description |
|------|-------------|
| `workflow_login` | Automated login with verification |
| `workflow_extract_table` | Extract table data to JSON |

## Example Conversation

**User:** Navigate to example.com and click the "More information" link

**Claude:** I'll help you navigate and click the link.

```
[Calls browser_launch]
[Calls page_navigate with url="https://example.com"]
[Calls element_click with selector="a"]
```

Done! I've navigated to example.com and clicked the link.

## Batch Execution Example

Reduce round-trip latency by batching multiple operations:

```json
{
  "tool": "batch_execute",
  "arguments": {
    "steps": [
      {"tool": "page_navigate", "args": {"url": "https://example.com/login"}},
      {"tool": "element_fill", "args": {"selector": "#username", "value": "user"}},
      {"tool": "element_fill", "args": {"selector": "#password", "value": "pass"}},
      {"tool": "element_click", "args": {"selector": "#login"}},
      {"tool": "wait_for_selector", "args": {"selector": "#dashboard"}}
    ]
  }
}
```

## HTTP Request Example

Make authenticated requests using the browser's session:

```json
{
  "tool": "http_request",
  "arguments": {
    "url": "https://api.example.com/data",
    "method": "POST",
    "content_type": "application/json",
    "body": "{\"key\": \"value\"}",
    "max_body_length": 4096
  }
}
```

## Session Recording

Record your actions for deterministic replay:

```
[Calls record_start with name="Example Test"]
[Calls page_navigate with url="https://example.com"]
[Calls element_click with selector="a"]
[Calls record_export]
```

The exported JSON can be run with `w3pilot run`.
