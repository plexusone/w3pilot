package mcp

import (
	"encoding/json"
	"sort"
)

// ToolInfo represents an MCP tool definition for export.
type ToolInfo struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Category    string `json:"category"`
}

// ToolList represents the complete list of MCP tools.
type ToolList struct {
	Tools      []ToolInfo     `json:"tools"`
	Categories map[string]int `json:"categories"`
	Total      int            `json:"total"`
}

// toolDefinitions contains all tool definitions organized by category.
//
// Naming convention: {namespace}_{verb}_{target}
//
// Principles:
//  1. Keep verbs explicit: element_get_text not element_text
//  2. Avoid abbreviations: wait_for_function not wait_fn
//  3. Use full words: human_pause not hitl_pause
//  4. Consistent verb patterns: get/set, start/stop, is/has
var toolDefinitions = []struct {
	category string
	tools    []ToolInfo
}{
	{
		category: "browser",
		tools: []ToolInfo{
			{Name: "browser_launch", Description: "Launch a browser instance. Call this before any other browser operations."},
			{Name: "browser_quit", Description: "Close the browser and cleanup resources."},
		},
	},
	{
		category: "page",
		tools: []ToolInfo{
			{Name: "page_navigate", Description: "Navigate to a URL."},
			{Name: "page_go_back", Description: "Navigate back in browser history."},
			{Name: "page_go_forward", Description: "Navigate forward in browser history."},
			{Name: "page_reload", Description: "Reload the current page."},
			{Name: "page_scroll", Description: "Scroll the page or a specific element in a direction."},
			{Name: "page_get_title", Description: "Get the page title."},
			{Name: "page_get_url", Description: "Get the current URL."},
			{Name: "page_get_content", Description: "Get the page HTML content."},
			{Name: "page_set_content", Description: "Set the page HTML content."},
			{Name: "page_get_viewport", Description: "Get the viewport dimensions."},
			{Name: "page_set_viewport", Description: "Set the viewport dimensions."},
			{Name: "page_screenshot", Description: "Capture a page screenshot."},
			{Name: "page_pdf", Description: "Generate a PDF of the page."},
			{Name: "page_new", Description: "Create a new page/tab."},
			{Name: "page_get_count", Description: "Get the page count."},
			{Name: "page_close", Description: "Close the current page."},
			{Name: "page_bring_to_front", Description: "Activate the page."},
			{Name: "page_emulate_media", Description: "Emulate CSS media features (colorScheme, reducedMotion, forcedColors, contrast)."},
			{Name: "page_set_geolocation", Description: "Set the browser's geolocation."},
			{Name: "page_inspect", Description: "Inspect page elements to discover buttons, links, inputs, and other interactive elements. Designed for AI agents."},
		},
	},
	{
		category: "element",
		tools: []ToolInfo{
			{Name: "element_click", Description: "Click an element by CSS selector."},
			{Name: "element_double_click", Description: "Double-click an element by CSS selector."},
			{Name: "element_type", Description: "Type text into an input element (appends to existing content)."},
			{Name: "element_fill", Description: "Clear an input and fill it with text (replaces existing content)."},
			{Name: "element_fill_form", Description: "Fill multiple form fields at once."},
			{Name: "element_clear", Description: "Clear the content of an input element."},
			{Name: "element_press", Description: "Press a key on an element (e.g., Enter, Tab, ArrowDown)."},
			{Name: "element_check", Description: "Check a checkbox element."},
			{Name: "element_uncheck", Description: "Uncheck a checkbox element."},
			{Name: "element_select", Description: "Select option(s) in a <select> element."},
			{Name: "element_set_files", Description: "Set files on a file input element."},
			{Name: "element_hover", Description: "Hover over an element."},
			{Name: "element_focus", Description: "Focus an element."},
			{Name: "element_scroll_into_view", Description: "Scroll an element into view."},
			{Name: "element_drag_to", Description: "Drag an element to another element."},
			{Name: "element_tap", Description: "Tap an element (touch gesture)."},
			{Name: "element_dispatch_event", Description: "Dispatch a DOM event on an element."},
			{Name: "element_get_text", Description: "Get the text content of an element."},
			{Name: "element_get_value", Description: "Get the value of an input element."},
			{Name: "element_get_inner_html", Description: "Get the innerHTML of an element."},
			{Name: "element_get_outer_html", Description: "Get the outerHTML of an element."},
			{Name: "element_get_inner_text", Description: "Get the innerText of an element."},
			{Name: "element_get_attribute", Description: "Get an attribute value from an element."},
			{Name: "element_get_bounding_box", Description: "Get the bounding box of an element."},
			{Name: "element_is_visible", Description: "Check if an element is visible."},
			{Name: "element_is_hidden", Description: "Check if an element is hidden."},
			{Name: "element_is_enabled", Description: "Check if an element is enabled."},
			{Name: "element_is_checked", Description: "Check if a checkbox/radio is checked."},
			{Name: "element_is_editable", Description: "Check if an element is editable."},
			{Name: "element_get_role", Description: "Get the ARIA role of an element."},
			{Name: "element_get_label", Description: "Get the accessible label of an element."},
			{Name: "element_screenshot", Description: "Capture an element screenshot."},
			{Name: "element_evaluate", Description: "Evaluate JavaScript with an element context."},
		},
	},
	{
		category: "input",
		tools: []ToolInfo{
			{Name: "input_keyboard_press", Description: "Press a key."},
			{Name: "input_keyboard_down", Description: "Hold a key down."},
			{Name: "input_keyboard_up", Description: "Release a key."},
			{Name: "input_keyboard_type", Description: "Type text via keyboard."},
			{Name: "input_mouse_click", Description: "Click at coordinates."},
			{Name: "input_mouse_move", Description: "Move mouse to coordinates."},
			{Name: "input_mouse_down", Description: "Press mouse button."},
			{Name: "input_mouse_up", Description: "Release mouse button."},
			{Name: "input_mouse_wheel", Description: "Scroll mouse wheel."},
			{Name: "input_mouse_drag", Description: "Drag from one point to another."},
			{Name: "input_touch_tap", Description: "Tap at coordinates."},
			{Name: "input_touch_swipe", Description: "Swipe gesture."},
		},
	},
	{
		category: "js",
		tools: []ToolInfo{
			{Name: "js_evaluate", Description: "Execute JavaScript code."},
			{Name: "js_add_script", Description: "Inject JavaScript into the page."},
			{Name: "js_add_style", Description: "Inject CSS styles into the page."},
			{Name: "js_init_script", Description: "Add JavaScript that runs before page scripts on every navigation."},
		},
	},
	{
		category: "wait",
		tools: []ToolInfo{
			{Name: "wait_for_state", Description: "Wait for an element to reach a state (visible, hidden, attached, detached)."},
			{Name: "wait_for_url", Description: "Wait for URL to match a pattern."},
			{Name: "wait_for_load", Description: "Wait for a load state (load, domcontentloaded, networkidle)."},
			{Name: "wait_for_function", Description: "Wait for a JavaScript function to return truthy."},
			{Name: "wait_for_selector", Description: "Wait for element to appear/disappear with state option."},
			{Name: "wait_for_text", Description: "Wait for text to appear on the page."},
		},
	},
	{
		category: "tab",
		tools: []ToolInfo{
			{Name: "tab_list", Description: "List all open browser tabs."},
			{Name: "tab_select", Description: "Switch to a specific tab."},
			{Name: "tab_close", Description: "Close a specific tab."},
		},
	},
	{
		category: "frame",
		tools: []ToolInfo{
			{Name: "frame_select", Description: "Switch to a frame by name or URL pattern."},
			{Name: "frame_select_main", Description: "Switch back to the main frame."},
		},
	},
	{
		category: "dialog",
		tools: []ToolInfo{
			{Name: "dialog_handle", Description: "Handle a browser dialog (alert, confirm, prompt, beforeunload)."},
			{Name: "dialog_get", Description: "Get information about the current dialog."},
		},
	},
	{
		category: "console",
		tools: []ToolInfo{
			{Name: "console_get_messages", Description: "Get console messages from the page."},
			{Name: "console_clear", Description: "Clear all buffered console messages."},
		},
	},
	{
		category: "network",
		tools: []ToolInfo{
			{Name: "network_get_requests", Description: "Get captured network requests."},
			{Name: "network_clear", Description: "Clear all buffered network requests."},
			{Name: "network_route", Description: "Register a mock response for requests matching a URL pattern."},
			{Name: "network_list_routes", Description: "List all active route handlers."},
			{Name: "network_unroute", Description: "Remove a previously registered route handler."},
			{Name: "network_set_offline", Description: "Set the browser's network state for offline mode testing."},
		},
	},
	{
		category: "storage",
		tools: []ToolInfo{
			{Name: "storage_get_cookies", Description: "Get browser cookies."},
			{Name: "storage_set_cookies", Description: "Set browser cookies."},
			{Name: "storage_clear_cookies", Description: "Clear all cookies."},
			{Name: "storage_delete_cookie", Description: "Delete a specific cookie by name."},
			{Name: "storage_get_state", Description: "Get complete browser storage state."},
			{Name: "storage_set_state", Description: "Restore browser storage from JSON."},
			{Name: "storage_clear_all", Description: "Clear all browser storage."},
			{Name: "storage_local_get", Description: "Get a value from localStorage by key."},
			{Name: "storage_local_set", Description: "Set a value in localStorage."},
			{Name: "storage_local_delete", Description: "Delete a key from localStorage."},
			{Name: "storage_local_clear", Description: "Clear all localStorage data."},
			{Name: "storage_local_list", Description: "List all keys and values in localStorage."},
			{Name: "storage_session_get", Description: "Get a value from sessionStorage by key."},
			{Name: "storage_session_set", Description: "Set a value in sessionStorage."},
			{Name: "storage_session_delete", Description: "Delete a key from sessionStorage."},
			{Name: "storage_session_clear", Description: "Clear all sessionStorage data."},
			{Name: "storage_session_list", Description: "List all keys and values in sessionStorage."},
		},
	},
	{
		category: "trace",
		tools: []ToolInfo{
			{Name: "trace_start", Description: "Start trace recording with screenshots and DOM snapshots."},
			{Name: "trace_stop", Description: "Stop trace recording and save or return the trace data."},
			{Name: "trace_chunk_start", Description: "Start a new trace chunk within an active trace."},
			{Name: "trace_chunk_stop", Description: "Stop the current trace chunk."},
			{Name: "trace_group_start", Description: "Start a trace group for logical grouping of actions."},
			{Name: "trace_group_stop", Description: "Stop the current trace group."},
		},
	},
	{
		category: "record",
		tools: []ToolInfo{
			{Name: "record_start", Description: "Begin recording actions."},
			{Name: "record_stop", Description: "Stop recording."},
			{Name: "record_export", Description: "Export recorded script."},
			{Name: "record_get_status", Description: "Check recording state."},
			{Name: "record_clear", Description: "Clear recorded steps."},
		},
	},
	{
		category: "test",
		tools: []ToolInfo{
			{Name: "test_assert_text", Description: "Assert text exists."},
			{Name: "test_assert_element", Description: "Assert element exists."},
			{Name: "test_assert_url", Description: "Assert URL matches."},
			{Name: "test_verify_value", Description: "Verify that an input element has the expected value."},
			{Name: "test_verify_list", Description: "Verify that a list of text items are all visible."},
			{Name: "test_verify_text", Description: "Verify element text matches expected value."},
			{Name: "test_verify_visible", Description: "Verify element is visible."},
			{Name: "test_verify_enabled", Description: "Verify element is enabled."},
			{Name: "test_verify_checked", Description: "Verify checkbox/radio is checked."},
			{Name: "test_verify_hidden", Description: "Verify element is hidden."},
			{Name: "test_verify_disabled", Description: "Verify element is disabled."},
			{Name: "test_generate_locator", Description: "Generate a locator string for a given element."},
			{Name: "test_get_report", Description: "Get test execution report."},
			{Name: "test_reset", Description: "Clear test results."},
			{Name: "test_set_target", Description: "Set test target description."},
			{Name: "test_validate_selectors", Description: "Validate CSS selectors before use. Returns whether elements exist, are visible, and suggests alternatives if not found."},
		},
	},
	{
		category: "accessibility",
		tools: []ToolInfo{
			{Name: "accessibility_snapshot", Description: "Get accessibility tree snapshot."},
		},
	},
	{
		category: "video",
		tools: []ToolInfo{
			{Name: "video_start", Description: "Start video recording."},
			{Name: "video_stop", Description: "Stop video recording and save the file."},
		},
	},
	{
		category: "human",
		tools: []ToolInfo{
			{Name: "human_pause", Description: "Pause automation for human interaction (CAPTCHA, login, etc.)."},
		},
	},
	{
		category: "config",
		tools: []ToolInfo{
			{Name: "config_get", Description: "Get the resolved MCP server configuration."},
		},
	},
	{
		category: "workflow",
		tools: []ToolInfo{
			{Name: "workflow_login", Description: "Automated login workflow: fill credentials, submit form, wait for success indicator."},
			{Name: "workflow_extract_table", Description: "Extract HTML table data to structured JSON with headers and rows."},
		},
	},
	{
		category: "state",
		tools: []ToolInfo{
			{Name: "state_save", Description: "Save browser state (cookies, localStorage, sessionStorage) to a named snapshot."},
			{Name: "state_load", Description: "Load browser state from a named snapshot."},
			{Name: "state_list", Description: "List all saved state snapshots."},
			{Name: "state_delete", Description: "Delete a saved state snapshot."},
		},
	},
	{
		category: "cdp",
		tools: []ToolInfo{
			{Name: "cdp_get_performance_metrics", Description: "Get Core Web Vitals (LCP, CLS, INP) and navigation timing metrics."},
			{Name: "cdp_get_memory_stats", Description: "Get JavaScript heap memory statistics."},
			{Name: "cdp_take_heap_snapshot", Description: "Capture a V8 heap snapshot for memory profiling."},
			{Name: "cdp_emulate_network", Description: "Emulate network conditions (Slow 3G, Fast 3G, 4G, Offline, or custom)."},
			{Name: "cdp_clear_network_emulation", Description: "Clear network emulation and restore normal conditions."},
			{Name: "cdp_emulate_cpu", Description: "Emulate CPU throttling (2x, 4x, 6x slowdown)."},
			{Name: "cdp_clear_cpu_emulation", Description: "Clear CPU emulation and restore normal speed."},
			{Name: "cdp_run_lighthouse", Description: "Run Lighthouse audit for performance, accessibility, SEO, and best practices."},
			{Name: "cdp_start_coverage", Description: "Start collecting JavaScript and CSS code coverage."},
			{Name: "cdp_stop_coverage", Description: "Stop coverage collection and return the coverage report."},
			{Name: "cdp_enable_console_debugger", Description: "Enable enhanced console debugging with stack traces."},
			{Name: "cdp_get_console_entries", Description: "Get console messages with full stack traces."},
			{Name: "cdp_get_browser_logs", Description: "Get browser logs including deprecations, interventions, and violations."},
			{Name: "cdp_disable_console_debugger", Description: "Disable enhanced console debugging."},
			{Name: "cdp_start_screencast", Description: "Start streaming screen frames."},
			{Name: "cdp_stop_screencast", Description: "Stop screen streaming."},
			{Name: "cdp_install_extension", Description: "Install a Chrome extension from a local path."},
			{Name: "cdp_uninstall_extension", Description: "Uninstall a Chrome extension by ID."},
			{Name: "cdp_list_extensions", Description: "List all installed Chrome extensions."},
			{Name: "cdp_get_response_body", Description: "Get the response body for a specific network request by ID."},
		},
	},
}

// ListTools returns the complete list of MCP tools with categories.
func ListTools() *ToolList {
	var tools []ToolInfo
	categories := make(map[string]int)

	for _, cat := range toolDefinitions {
		for _, tool := range cat.tools {
			tool.Category = cat.category
			tools = append(tools, tool)
		}
		categories[cat.category] = len(cat.tools)
	}

	// Sort tools by name for consistent output
	sort.Slice(tools, func(i, j int) bool {
		return tools[i].Name < tools[j].Name
	})

	return &ToolList{
		Tools:      tools,
		Categories: categories,
		Total:      len(tools),
	}
}

// ListToolsJSON returns the tool list as formatted JSON.
func ListToolsJSON() ([]byte, error) {
	list := ListTools()
	return json.MarshalIndent(list, "", "  ")
}

// CategorySummary returns a summary of tools by category.
func CategorySummary() map[string]int {
	categories := make(map[string]int)
	for _, cat := range toolDefinitions {
		categories[cat.category] = len(cat.tools)
	}
	return categories
}
