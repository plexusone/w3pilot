# MCP Tools Reference

Complete reference for all 80+ MCP tools.

## Browser Management

### browser_launch

Launch a browser instance.

**Input:**

| Field | Type | Description |
|-------|------|-------------|
| `headless` | boolean | Run without GUI |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `message` | string | Status message |

### browser_quit

Close the browser.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `message` | string | Status message |

## Navigation

### navigate

Navigate to a URL.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `url` | string | ✅ | Target URL |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `url` | string | Current URL |
| `title` | string | Page title |

### back

Navigate back in history.

### forward

Navigate forward in history.

### reload

Reload the current page.

### scroll

Scroll the page or a specific element.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `direction` | string | ✅ | Scroll direction: up, down, left, right |
| `amount` | integer | | Pixels to scroll (0 for full page) |
| `selector` | string | | CSS selector to scroll within |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `message` | string | Status message |

## Element Interactions

### click

Click an element.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |
| `timeout_ms` | integer | | Timeout (default: 5000) |

### dblclick

Double-click an element.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |
| `timeout_ms` | integer | | Timeout |

### type

Type text into an element (appends).

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |
| `text` | string | ✅ | Text to type |
| `timeout_ms` | integer | | Timeout |

### fill

Fill an input (replaces content).

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |
| `value` | string | ✅ | Value to fill |
| `timeout_ms` | integer | | Timeout |

### fill_form

Fill multiple form fields at once.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `fields` | array | ✅ | Array of {selector, value} objects |
| `timeout_ms` | integer | | Timeout per field |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `message` | string | Status message |
| `filled` | integer | Number of fields filled |
| `errors` | array | Any errors encountered |

### clear

Clear an input element.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |

### press

Press a key on an element.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |
| `key` | string | ✅ | Key (e.g., "Enter") |

## Form Controls

### check

Check a checkbox.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |

### uncheck

Uncheck a checkbox.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |

### select_option

Select dropdown option(s).

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |
| `values` | array | | Option values |
| `labels` | array | | Option labels |
| `indexes` | array | | Option indexes |

### set_files

Set files on a file input.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |
| `files` | array | ✅ | File paths |

## Element State

### get_text

Get element text content.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `text` | string | Text content |

### get_value

Get input element value.

### get_inner_html

Get element innerHTML.

### get_outer_html

Get element outerHTML (including the element itself).

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `html` | string | Element's outer HTML |

### get_inner_text

Get element innerText.

### get_attribute

Get element attribute.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |
| `name` | string | ✅ | Attribute name |

### get_bounding_box

Get element bounding box.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `x` | number | X position |
| `y` | number | Y position |
| `width` | number | Width |
| `height` | number | Height |

### is_visible

Check if element is visible.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `visible` | boolean | Visibility state |

### is_hidden

Check if element is hidden.

### is_enabled

Check if element is enabled.

### is_checked

Check if checkbox/radio is checked.

### is_editable

Check if element is editable.

### get_role

Get ARIA role.

### get_label

Get accessible label.

## Page State

### get_title

Get page title.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `title` | string | Page title |

### get_url

Get current URL.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `url` | string | Current URL |

### get_content

Get page HTML content.

### set_content

Set page HTML content.

### get_viewport

Get viewport dimensions.

### set_viewport

Set viewport dimensions.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `width` | integer | ✅ | Width |
| `height` | integer | ✅ | Height |

## Screenshots & PDF

### screenshot

Capture page screenshot.

**Input:**

| Field | Type | Description |
|-------|------|-------------|
| `format` | string | "base64" or "file" |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `data` | string | Base64 image data |

### element_screenshot

Capture element screenshot.

### pdf

Generate PDF.

## JavaScript

### evaluate

Execute JavaScript.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `script` | string | ✅ | JavaScript code |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `result` | any | Evaluation result |

### element_eval

Evaluate JavaScript with element.

### add_script

Inject JavaScript.

### add_style

Inject CSS.

## Waiting

### wait_until

Wait for element state.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |
| `state` | string | ✅ | visible/hidden/attached/detached |

### wait_for_url

Wait for URL pattern.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `pattern` | string | ✅ | URL pattern |

### wait_for_load

Wait for load state.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `state` | string | ✅ | load/domcontentloaded/networkidle |

### wait_for_function

Wait for JavaScript function.

## Input Controllers

### keyboard_press

Press a key.

### keyboard_down

Hold a key.

### keyboard_up

Release a key.

### keyboard_type

Type text via keyboard.

### mouse_click

Click at coordinates.

### mouse_move

Move mouse.

### mouse_down

Press mouse button.

### mouse_up

Release mouse button.

### mouse_wheel

Scroll mouse wheel.

### touch_tap

Tap at coordinates.

### touch_swipe

Swipe gesture.

### mouse_drag

Drag from one point to another using the mouse.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `start_x` | number | ✅ | Starting X coordinate |
| `start_y` | number | ✅ | Starting Y coordinate |
| `end_x` | number | ✅ | Ending X coordinate |
| `end_y` | number | ✅ | Ending Y coordinate |
| `steps` | integer | | Number of intermediate steps (default: 10) |

## Page Management

### new_page

Create new page/tab.

### get_pages

Get page count.

### close_page

Close current page.

### bring_to_front

Activate page.

## Emulation

### emulate_media

Emulate CSS media features for accessibility testing.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `media` | string | | Media type: "screen" or "print" |
| `color_scheme` | string | | Color scheme: "light", "dark", "no-preference" |
| `reduced_motion` | string | | Reduced motion: "reduce", "no-preference" |
| `forced_colors` | string | | Forced colors: "active", "none" |
| `contrast` | string | | Contrast: "more", "less", "no-preference" |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `message` | string | Status message |
| `settings` | array | Applied settings |

**Accessibility Testing Use Cases:**

- **color_scheme**: Test dark/light mode support for users with light sensitivity
- **reduced_motion**: Test that animations are disabled for users with vestibular disorders
- **forced_colors**: Test Windows High Contrast Mode compatibility
- **contrast**: Test increased/decreased contrast for users with low vision

### set_geolocation

Set the browser's geolocation.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `latitude` | number | Yes | Latitude coordinate |
| `longitude` | number | Yes | Longitude coordinate |
| `accuracy` | number | | Accuracy in meters |

## Tab Management

### list_tabs

List all open browser tabs.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `tabs` | array | Tab information (index, id, url, title) |
| `count` | integer | Number of open tabs |
| `current_tab` | integer | Index of the current tab |

### select_tab

Switch to a specific tab.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `index` | integer | | Tab index (0-based) |
| `id` | string | | Tab ID from list_tabs |

### close_tab

Close a specific tab.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `index` | integer | | Tab index to close |
| `id` | string | | Tab ID to close |

## Dialog Handling

### handle_dialog

Handle a browser dialog (alert, confirm, prompt, beforeunload).

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `action` | string | Yes | Action: "accept" or "dismiss" |
| `prompt_text` | string | | Text for prompt dialogs |

### get_dialog

Get information about the current dialog.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `has_dialog` | boolean | Whether a dialog is open |
| `dialog_type` | string | alert, confirm, prompt, beforeunload |
| `message` | string | Dialog message |
| `default_value` | string | Default value for prompt dialogs |

## Console Messages

### get_console_messages

Get console messages from the page.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `level` | string | | Filter by level: log, info, warn, error, debug |
| `clear` | boolean | | Clear messages after retrieving |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `messages` | array | Console message objects |
| `count` | integer | Number of messages |

Each message contains:

| Field | Type | Description |
|-------|------|-------------|
| `type` | string | Message type (log, info, warn, error, debug) |
| `text` | string | Message text |
| `args` | array | Additional arguments |
| `url` | string | Source URL |
| `line` | integer | Source line number |

### clear_console_messages

Clear all buffered console messages.

## Network Requests

### get_network_requests

Get captured network requests.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `url_pattern` | string | | Filter by URL pattern (glob or regex) |
| `method` | string | | Filter by HTTP method |
| `resource_type` | string | | Filter by resource type |
| `clear` | boolean | | Clear requests after retrieving |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `requests` | array | Network request objects |
| `count` | integer | Number of requests |

Each request contains:

| Field | Type | Description |
|-------|------|-------------|
| `url` | string | Request URL |
| `method` | string | HTTP method |
| `resource_type` | string | Resource type |
| `status` | integer | Response status code |
| `status_text` | string | Response status text |
| `response_size` | integer | Response size in bytes |

### clear_network_requests

Clear all buffered network requests.

## Network Mocking

### route

Register a mock response for requests matching a URL pattern.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `pattern` | string | ✅ | URL pattern (glob or regex, e.g., `**/api/*`) |
| `status` | integer | | HTTP status code (default: 200) |
| `body` | string | | Response body content |
| `content_type` | string | | Content-Type header (default: application/json) |
| `headers` | object | | Additional response headers |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `message` | string | Status message |
| `pattern` | string | Registered pattern |

### route_list

List all active route handlers.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `routes` | array | Route info objects |
| `count` | integer | Number of active routes |

Each route contains:

| Field | Type | Description |
|-------|------|-------------|
| `pattern` | string | URL pattern |
| `status` | integer | Response status code |
| `content_type` | string | Content-Type header |

### unroute

Remove a previously registered route handler.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `pattern` | string | ✅ | URL pattern to unregister |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `message` | string | Status message |

### network_state_set

Set the browser's network state for offline mode testing.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `offline` | boolean | ✅ | Set to true for offline mode |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `message` | string | Status message |
| `offline` | boolean | Current offline state |

## Cookies & Storage

### get_cookies

Get browser cookies.

### set_cookies

Set browser cookies.

### clear_cookies

Clear all cookies.

### get_storage_state

Get complete browser storage state including cookies, localStorage, and sessionStorage as JSON. This can be saved to a file and later restored using `set_storage_state` to resume a session.

**Output:**

Returns JSON containing:

- `cookies`: Array of all cookies
- `origins`: Array of origin storage, each with:
  - `origin`: The origin URL (e.g., "https://example.com")
  - `localStorage`: Key-value map of localStorage items
  - `sessionStorage`: Key-value map of sessionStorage items (for current page's origin)

### set_storage_state

Restore browser storage from JSON (output of `get_storage_state`). Restores cookies, localStorage, and sessionStorage.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `state` | string | Yes | JSON from get_storage_state |

### clear_storage

Clear all browser storage including cookies, localStorage, and sessionStorage.

## LocalStorage

### localstorage_get

Get a value from localStorage by key.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `key` | string | Yes | Key to retrieve |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `key` | string | The key |
| `value` | string | The value (null if not found) |

### localstorage_set

Set a value in localStorage.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `key` | string | Yes | Key to set |
| `value` | string | Yes | Value to store |

### localstorage_delete

Delete a key from localStorage.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `key` | string | Yes | Key to delete |

### localstorage_clear

Clear all localStorage data for the current origin.

### localstorage_list

List all keys and values in localStorage.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `items` | object | Key-value pairs |
| `count` | integer | Number of items |

## SessionStorage

### sessionstorage_get

Get a value from sessionStorage by key.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `key` | string | Yes | Key to retrieve |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `key` | string | The key |
| `value` | string | The value (null if not found) |

### sessionstorage_set

Set a value in sessionStorage.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `key` | string | Yes | Key to set |
| `value` | string | Yes | Value to store |

### sessionstorage_delete

Delete a key from sessionStorage.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `key` | string | Yes | Key to delete |

### sessionstorage_clear

Clear all sessionStorage data for the current origin.

### sessionstorage_list

List all keys and values in sessionStorage.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `items` | object | Key-value pairs |
| `count` | integer | Number of items |

## Script Recording

### start_recording

Begin recording actions.

**Input:**

| Field | Type | Description |
|-------|------|-------------|
| `name` | string | Script name |
| `description` | string | Description |
| `baseUrl` | string | Base URL |

### stop_recording

Stop recording.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `stepCount` | integer | Steps recorded |

### export_script

Export recorded script.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `script` | string | JSON script |
| `stepCount` | integer | Steps |
| `format` | string | Output format |

### recording_status

Check recording state.

### clear_recording

Clear recorded steps.

## Tracing

Tools for recording detailed execution traces for debugging and analysis. Traces include screenshots, DOM snapshots, and network activity. View traces with `npx playwright show-trace <trace.zip>`.

### start_trace

Start trace recording with screenshots and DOM snapshots.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | | Trace name (used for file naming) |
| `title` | string | | Title shown in trace viewer |
| `screenshots` | boolean | | Include screenshots (default: true) |
| `snapshots` | boolean | | Include DOM snapshots (default: true) |
| `sources` | boolean | | Include source files |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `message` | string | Status message |
| `name` | string | Trace name |

### stop_trace

Stop trace recording and save or return the trace data.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `path` | string | | File path to save trace ZIP (returns base64 if omitted) |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `message` | string | Status message |
| `path` | string | File path if saved |
| `data` | string | Base64-encoded ZIP if no path specified |
| `size_kb` | integer | Trace size in KB |
| `view_hint` | string | Command to view trace |

### start_trace_chunk

Start a new trace chunk within an active trace for segmenting recordings.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | | Chunk name |
| `title` | string | | Chunk title shown in trace viewer |

### stop_trace_chunk

Stop the current trace chunk.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `path` | string | | File path to save chunk ZIP |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `message` | string | Status message |
| `path` | string | File path if saved |
| `data` | string | Base64-encoded ZIP if no path |
| `size_kb` | integer | Chunk size in KB |

### start_trace_group

Start a trace group for logical grouping of actions in the trace viewer.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | ✅ | Group name |
| `location` | string | | Source location to associate |

### stop_trace_group

Stop the current trace group.

## Init Scripts

### add_init_script

Add JavaScript that runs before page scripts on every navigation. Useful for mocking APIs, injecting test helpers, or setting up authentication.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `script` | string | ✅ | JavaScript code to inject before page scripts |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `message` | string | Status message |

**Example Use Cases:**

```javascript
// Mock fetch API
window.fetch = async (url) => {
  if (url.includes('/api/user')) {
    return { json: () => ({ id: 1, name: 'Test User' }) };
  }
  return originalFetch(url);
};

// Disable analytics
window.gtag = () => {};
window.analytics = { track: () => {}, identify: () => {} };

// Inject test helpers
window.testHelpers = {
  fillForm: (data) => { /* ... */ },
  waitForElement: (sel) => { /* ... */ }
};
```

## Assertions

### assert_text

Assert text exists.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `text` | string | ✅ | Expected text |
| `selector` | string | | Limit to element |

### assert_element

Assert element exists.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector |

## Testing Tools

### verify_value

Verify that an input element has the expected value.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector for the input element |
| `expected` | string | ✅ | Expected value to verify |
| `timeout_ms` | integer | | Timeout in milliseconds (default: 5000) |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `passed` | boolean | Whether the verification passed |
| `actual` | string | The actual value found |
| `message` | string | Status message |

### verify_list_visible

Verify that a list of text items are all visible on the page.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `items` | array | ✅ | List of text items that should be visible |
| `selector` | string | | Optional CSS selector to scope the search |
| `timeout_ms` | integer | | Timeout in milliseconds (default: 5000) |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `passed` | boolean | Whether all items were found |
| `found` | array | Items that were found |
| `missing` | array | Items that were not found |
| `message` | string | Status message |

### generate_locator

Generate a locator string for a given element using a specific strategy.

**Input:**

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `selector` | string | ✅ | CSS selector for the element |
| `strategy` | string | | Locator strategy: css, xpath, testid, role, text (default: css) |
| `timeout_ms` | integer | | Timeout in milliseconds (default: 5000) |

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `locator` | string | Generated locator string |
| `strategy` | string | Strategy used |
| `metadata` | object | Additional metadata (role, label, text, etc.) |

## Test Reporting

### get_test_report

Get test execution report.

### reset_session

Clear test results.

### set_target

Set test target description.

## Configuration

### get_config

Get the resolved MCP server configuration.

**Output:**

| Field | Type | Description |
|-------|------|-------------|
| `headless` | boolean | Whether browser runs headless |
| `project` | string | Project name for reports |
| `default_timeout_ms` | integer | Default timeout in milliseconds |
| `browser_launched` | boolean | Whether browser has been launched |
