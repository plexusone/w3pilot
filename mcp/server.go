package mcp

import (
	"context"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Server is the WebPilot MCP server.
type Server struct {
	session   *Session
	mcpServer *mcp.Server
	config    Config
}

// NewServer creates a new MCP server.
func NewServer(config Config) *Server {
	s := &Server{
		config: config,
		session: NewSession(SessionConfig{
			Headless:       config.Headless,
			DefaultTimeout: config.DefaultTimeout,
			Project:        config.Project,
			InitScripts:    config.InitScripts,
		}),
	}

	s.mcpServer = mcp.NewServer(
		&mcp.Implementation{
			Name:    "w3pilot-mcp",
			Version: "0.2.0",
		},
		nil,
	)

	s.registerTools()
	return s
}

// registerTools registers all MCP tools.
func (s *Server) registerTools() {
	// === Browser Management ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "browser_launch",
		Description: "Launch a browser instance. Call this before any other browser operations.",
	}, s.handleBrowserLaunch)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "browser_quit",
		Description: "Close the browser and cleanup resources.",
	}, s.handleBrowserQuit)

	// === Navigation ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "page_navigate",
		Description: "Navigate to a URL.",
	}, s.handleNavigate)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "page_go_back",
		Description: "Navigate back in browser history.",
	}, s.handleBack)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "page_go_forward",
		Description: "Navigate forward in browser history.",
	}, s.handleForward)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "page_reload",
		Description: "Reload the current page.",
	}, s.handleReload)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "page_scroll",
		Description: "Scroll the page or a specific element in a direction (up, down, left, right).",
	}, s.handleScroll)

	// === Basic Interactions ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "element_click",
		Description: "Click an element by CSS selector.",
	}, s.handleClick)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "element_double_click",
		Description: "Double-click an element by CSS selector.",
	}, s.handleDblClick)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "element_type",
		Description: "Type text into an input element (appends to existing content).",
	}, s.handleType)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "element_fill",
		Description: "Clear an input and fill it with text (replaces existing content).",
	}, s.handleFill)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "element_fill_form",
		Description: "Fill multiple form fields at once. Provide an array of {selector, value} pairs.",
	}, s.handleFillForm)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "element_clear",
		Description: "Clear the content of an input element.",
	}, s.handleClear)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "element_press",
		Description: "Press a key on an element (e.g., Enter, Tab, ArrowDown).",
	}, s.handlePress)

	// === Form Controls ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "element_check",
		Description: "Check a checkbox element.",
	}, s.handleCheck)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "element_uncheck",
		Description: "Uncheck a checkbox element.",
	}, s.handleUncheck)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "element_select",
		Description: "Select option(s) in a <select> element by value, label, or index.",
	}, s.handleSelectOption)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "element_set_files",
		Description: "Set files on a file input element.",
	}, s.handleSetFiles)

	// === Element Interaction ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "element_hover",
		Description: "Hover over an element.",
	}, s.handleHover)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "element_focus",
		Description: "Focus an element.",
	}, s.handleFocus)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "element_scroll_into_view",
		Description: "Scroll an element into view.",
	}, s.handleScrollIntoView)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "element_drag_to",
		Description: "Drag an element to another element.",
	}, s.handleDragTo)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "element_tap",
		Description: "Tap an element (touch gesture).",
	}, s.handleTap)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "element_dispatch_event",
		Description: "Dispatch a DOM event on an element.",
	}, s.handleDispatchEvent)

	// === Element State ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "element_get_text",
		Description: "Get the text content of an element.",
	}, s.handleGetText)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "element_get_value",
		Description: "Get the value of an input element.",
	}, s.handleGetValue)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "element_get_inner_html",
		Description: "Get the innerHTML of an element.",
	}, s.handleGetInnerHTML)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "element_get_outer_html",
		Description: "Get the outerHTML of an element (including the element itself).",
	}, s.handleGetOuterHTML)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "element_get_inner_text",
		Description: "Get the innerText of an element.",
	}, s.handleGetInnerText)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "element_get_attribute",
		Description: "Get an attribute value of an element.",
	}, s.handleGetAttribute)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "element_get_bounding_box",
		Description: "Get the bounding box of an element.",
	}, s.handleGetBoundingBox)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "element_is_visible",
		Description: "Check if an element is visible.",
	}, s.handleIsVisible)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "element_is_hidden",
		Description: "Check if an element is hidden.",
	}, s.handleIsHidden)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "element_is_enabled",
		Description: "Check if an element is enabled.",
	}, s.handleIsEnabled)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "element_is_checked",
		Description: "Check if a checkbox/radio is checked.",
	}, s.handleIsChecked)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "element_is_editable",
		Description: "Check if an element is editable.",
	}, s.handleIsEditable)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "element_get_role",
		Description: "Get the ARIA role of an element.",
	}, s.handleGetRole)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "element_get_label",
		Description: "Get the accessible label of an element.",
	}, s.handleGetLabel)

	// === Page State ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "page_get_title",
		Description: "Get the current page title.",
	}, s.handleGetTitle)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "page_get_url",
		Description: "Get the current page URL.",
	}, s.handleGetURL)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "page_get_content",
		Description: "Get the full HTML content of the page.",
	}, s.handleGetContent)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "page_set_content",
		Description: "Set the HTML content of the page.",
	}, s.handleSetContent)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "page_get_viewport",
		Description: "Get the viewport dimensions.",
	}, s.handleGetViewport)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "page_set_viewport",
		Description: "Set the viewport dimensions.",
	}, s.handleSetViewport)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "get_frames",
		Description: "Get all frames on the page.",
	}, s.handleGetFrames)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "frame_select",
		Description: "Switch to a frame by name or URL pattern. Subsequent commands will target this frame.",
	}, s.handleSelectFrame)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "frame_select_main",
		Description: "Switch back to the main frame (top-level page).",
	}, s.handleSelectMainFrame)

	// === Screenshots & PDF ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "page_screenshot",
		Description: "Capture a screenshot of the current page.",
	}, s.handleScreenshot)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "element_screenshot",
		Description: "Capture a screenshot of a specific element.",
	}, s.handleElementScreenshot)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "page_pdf",
		Description: "Generate a PDF of the page.",
	}, s.handlePDF)

	// === JavaScript ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "js_evaluate",
		Description: "Execute JavaScript on the page and return the result.",
	}, s.handleEvaluate)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "element_evaluate",
		Description: "Evaluate JavaScript with an element as the first argument.",
	}, s.handleElementEval)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "js_add_script",
		Description: "Inject JavaScript into the page.",
	}, s.handleAddScript)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "js_add_style",
		Description: "Inject CSS into the page.",
	}, s.handleAddStyle)

	// === Waiting ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "wait_for_state",
		Description: "Wait for an element to reach a state (attached, detached, visible, hidden).",
	}, s.handleWaitUntil)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "wait_for_url",
		Description: "Wait for the URL to match a pattern.",
	}, s.handleWaitForURL)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "wait_for_load",
		Description: "Wait for page load state (load, domcontentloaded, networkidle).",
	}, s.handleWaitForLoad)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "wait_for_function",
		Description: "Wait for a JavaScript function to return truthy.",
	}, s.handleWaitForFunction)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "wait_for_text",
		Description: "Wait for text to appear on the page. Optionally scope to a specific element.",
	}, s.handleWaitForText)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "accessibility_snapshot",
		Description: "Get an accessibility tree snapshot of the page. Useful for understanding page structure and testing accessibility.",
	}, s.handleAccessibilitySnapshot)

	// === Input Controllers ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "input_keyboard_press",
		Description: "Press a key on the keyboard.",
	}, s.handleKeyboardPress)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "input_keyboard_down",
		Description: "Hold down a key.",
	}, s.handleKeyboardDown)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "input_keyboard_up",
		Description: "Release a held key.",
	}, s.handleKeyboardUp)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "input_keyboard_type",
		Description: "Type text using the keyboard.",
	}, s.handleKeyboardType)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "input_mouse_click",
		Description: "Click at coordinates.",
	}, s.handleMouseClick)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "input_mouse_move",
		Description: "Move the mouse to coordinates.",
	}, s.handleMouseMove)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "input_mouse_down",
		Description: "Press the mouse button.",
	}, s.handleMouseDown)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "input_mouse_up",
		Description: "Release the mouse button.",
	}, s.handleMouseUp)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "input_mouse_wheel",
		Description: "Scroll the mouse wheel.",
	}, s.handleMouseWheel)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "input_touch_tap",
		Description: "Tap at coordinates (touch).",
	}, s.handleTouchTap)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "input_touch_swipe",
		Description: "Swipe from one point to another (touch).",
	}, s.handleTouchSwipe)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "input_mouse_drag",
		Description: "Drag from one point to another using the mouse.",
	}, s.handleMouseDrag)

	// === Page Management ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "page_new",
		Description: "Create a new page/tab.",
	}, s.handleNewPage)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "page_get_count",
		Description: "Get the number of open pages.",
	}, s.handleGetPages)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "page_close",
		Description: "Close the current page.",
	}, s.handleClosePage)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "page_bring_to_front",
		Description: "Bring the page to the front.",
	}, s.handleBringToFront)

	// === Tab Management ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "tab_list",
		Description: "List all open browser tabs with their index, ID, URL, and title.",
	}, s.handleListTabs)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "tab_select",
		Description: "Switch to a specific tab by index (0-based) or tab ID.",
	}, s.handleSelectTab)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "tab_close",
		Description: "Close a specific tab by index or ID. Defaults to current tab if not specified.",
	}, s.handleCloseTab)

	// === Emulation ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "page_emulate_media",
		Description: "Emulate CSS media features for accessibility testing: color scheme (dark/light mode), reduced motion (disable animations), forced colors (high contrast mode), and contrast preferences.",
	}, s.handleEmulateMedia)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "page_set_geolocation",
		Description: "Set the browser's geolocation.",
	}, s.handleSetGeolocation)

	// === Cookies & Storage ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "storage_get_cookies",
		Description: "Get browser cookies.",
	}, s.handleGetCookies)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "storage_set_cookies",
		Description: "Set browser cookies.",
	}, s.handleSetCookies)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "storage_clear_cookies",
		Description: "Clear all cookies.",
	}, s.handleClearCookies)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "storage_delete_cookie",
		Description: "Delete a specific cookie by name. Optionally filter by domain and path.",
	}, s.handleDeleteCookie)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "storage_get_state",
		Description: "Get complete browser storage state (cookies, localStorage, and sessionStorage) as JSON.",
	}, s.handleGetStorageState)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "storage_set_state",
		Description: "Restore browser storage from JSON (output of get_storage_state). Restores cookies, localStorage, and sessionStorage.",
	}, s.handleSetStorageState)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "storage_clear_all",
		Description: "Clear all browser storage (cookies, localStorage, and sessionStorage).",
	}, s.handleClearStorage)

	// === LocalStorage ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "storage_local_get",
		Description: "Get a value from localStorage by key.",
	}, s.handleLocalStorageGet)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "storage_local_set",
		Description: "Set a value in localStorage.",
	}, s.handleLocalStorageSet)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "storage_local_delete",
		Description: "Delete a key from localStorage.",
	}, s.handleLocalStorageDelete)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "storage_local_clear",
		Description: "Clear all localStorage data for the current origin.",
	}, s.handleLocalStorageClear)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "storage_local_list",
		Description: "List all keys and values in localStorage.",
	}, s.handleLocalStorageList)

	// === SessionStorage ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "storage_session_get",
		Description: "Get a value from sessionStorage by key.",
	}, s.handleSessionStorageGet)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "storage_session_set",
		Description: "Set a value in sessionStorage.",
	}, s.handleSessionStorageSet)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "storage_session_delete",
		Description: "Delete a key from sessionStorage.",
	}, s.handleSessionStorageDelete)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "storage_session_clear",
		Description: "Clear all sessionStorage data for the current origin.",
	}, s.handleSessionStorageClear)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "storage_session_list",
		Description: "List all keys and values in sessionStorage.",
	}, s.handleSessionStorageList)

	// === Dialog Handling ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "dialog_handle",
		Description: "Handle a browser dialog (alert, confirm, prompt, beforeunload) by accepting or dismissing it.",
	}, s.handleHandleDialog)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "dialog_get",
		Description: "Get information about the current dialog, if any is open.",
	}, s.handleGetDialog)

	// === Console Messages ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "console_get_messages",
		Description: "Get console messages from the page. Optionally filter by level (log, info, warn, error, debug).",
	}, s.handleGetConsoleMessages)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "console_clear",
		Description: "Clear the buffered console messages.",
	}, s.handleClearConsoleMessages)

	// === Network Requests ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "network_get_requests",
		Description: "Get captured network requests. Optionally filter by URL pattern, HTTP method, or resource type.",
	}, s.handleGetNetworkRequests)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "network_clear",
		Description: "Clear the buffered network requests.",
	}, s.handleClearNetworkRequests)

	// === Network Mocking ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "network_route",
		Description: "Register a mock response for requests matching a URL pattern. Use glob patterns (e.g., **/api/*) or regex (e.g., /api/.*).",
	}, s.handleRoute)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "network_list_routes",
		Description: "List all active route handlers.",
	}, s.handleRouteList)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "network_unroute",
		Description: "Remove a previously registered route handler.",
	}, s.handleUnroute)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "network_set_offline",
		Description: "Set the browser's network state. Use offline=true to simulate offline mode for testing.",
	}, s.handleNetworkStateSet)

	// === Human-in-the-Loop ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "human_pause",
		Description: "Pause automation and wait for human to complete an action (e.g., SSO login, CAPTCHA). Shows a visual overlay that the human dismisses when done.",
	}, s.handlePauseForHuman)

	// === Assertions ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "test_assert_text",
		Description: "Assert that text exists on the page.",
	}, s.handleAssertText)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "test_assert_element",
		Description: "Assert that an element exists on the page.",
	}, s.handleAssertElement)

	// === Testing Tools ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "test_verify_value",
		Description: "Verify that an input element has the expected value.",
	}, s.handleVerifyValue)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "test_verify_list",
		Description: "Verify that a list of text items are all visible on the page.",
	}, s.handleVerifyListVisible)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "test_generate_locator",
		Description: "Generate a locator string for a given element using a specific strategy (css, xpath, testid, role, text).",
	}, s.handleGenerateLocator)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "wait_for_selector",
		Description: "Wait for an element to reach a specific state (attached, detached, visible, hidden).",
	}, s.handleWaitForSelector)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "test_verify_text",
		Description: "Verify that an element's text content matches the expected value.",
	}, s.handleVerifyText)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "test_verify_visible",
		Description: "Verify that an element is visible on the page.",
	}, s.handleVerifyVisible)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "test_verify_enabled",
		Description: "Verify that an element is enabled (not disabled).",
	}, s.handleVerifyEnabled)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "test_verify_checked",
		Description: "Verify that a checkbox or radio button is checked or unchecked.",
	}, s.handleVerifyChecked)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "test_verify_hidden",
		Description: "Verify that an element is hidden (not visible) on the page.",
	}, s.handleVerifyHidden)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "test_verify_disabled",
		Description: "Verify that an element is disabled.",
	}, s.handleVerifyDisabled)

	// === Test Reporting ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "test_get_report",
		Description: "Get the test execution report in the specified format (box, diagnostic, or json).",
	}, s.handleGetTestReport)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "test_reset",
		Description: "Clear test results and start a new test session.",
	}, s.handleResetSession)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "test_set_target",
		Description: "Set the test target description for reports.",
	}, s.handleSetTarget)

	// === Script Recording ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "record_start",
		Description: "Start recording browser actions to create a replayable test script.",
	}, s.handleStartRecording)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "record_stop",
		Description: "Stop recording browser actions.",
	}, s.handleStopRecording)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "record_export",
		Description: "Export recorded actions as a JSON test script that can be run with 'w3pilot run'.",
	}, s.handleExportScript)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "record_get_status",
		Description: "Check if recording is active and how many steps have been recorded.",
	}, s.handleRecordingStatus)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "record_clear",
		Description: "Clear all recorded steps without stopping recording.",
	}, s.handleClearRecording)

	// === Tracing ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "trace_start",
		Description: "Start trace recording with screenshots and DOM snapshots for debugging. The trace can be viewed with 'npx playwright show-trace <trace.zip>'.",
	}, s.handleStartTrace)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "trace_stop",
		Description: "Stop trace recording and save or return the trace data as a ZIP file.",
	}, s.handleStopTrace)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "trace_chunk_start",
		Description: "Start a new trace chunk within an active trace. Useful for segmenting traces into logical sections.",
	}, s.handleStartTraceChunk)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "trace_chunk_stop",
		Description: "Stop the current trace chunk and optionally save it to a file.",
	}, s.handleStopTraceChunk)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "trace_group_start",
		Description: "Start a trace group for logical grouping of actions in the trace viewer.",
	}, s.handleStartTraceGroup)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "trace_group_stop",
		Description: "Stop the current trace group.",
	}, s.handleStopTraceGroup)

	// === Video Recording ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "video_start",
		Description: "Start recording video of the browser page. Video is saved when stop_video is called.",
	}, s.handleStartVideo)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "video_stop",
		Description: "Stop video recording and return the path to the video file.",
	}, s.handleStopVideo)

	// === Init Scripts ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "js_init_script",
		Description: "Add JavaScript that runs before page scripts on every navigation. Useful for mocking APIs, injecting test helpers, or setting up authentication.",
	}, s.handleAddInitScript)

	// === Configuration ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "config_get",
		Description: "Get the resolved MCP server configuration including headless mode, project name, and timeouts.",
	}, s.handleGetConfig)

	// === Performance & Profiling (CDP) ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "get_performance_metrics",
		Description: "Get Core Web Vitals and navigation timing metrics (LCP, CLS, FCP, TTFB, etc.).",
	}, s.handleGetPerformanceMetrics)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "get_memory_stats",
		Description: "Get JavaScript heap memory statistics (used, total, limit).",
	}, s.handleGetMemoryStats)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "take_heap_snapshot",
		Description: "Capture a V8 heap snapshot for memory profiling. Requires CDP connection. Output can be loaded in Chrome DevTools Memory tab.",
	}, s.handleTakeHeapSnapshot)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "emulate_network",
		Description: "Simulate network conditions (slow3g, fast3g, 4g, wifi, offline) for performance testing. Requires CDP connection.",
	}, s.handleEmulateNetwork)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "clear_network_emulation",
		Description: "Remove network throttling and return to normal network conditions.",
	}, s.handleClearNetworkEmulation)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "emulate_cpu",
		Description: "Simulate slower CPU (rate: 1=none, 2=2x slower, 4=4x slower/mid-tier mobile, 6=6x slower/low-end mobile). Requires CDP connection.",
	}, s.handleEmulateCPU)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "clear_cpu_emulation",
		Description: "Remove CPU throttling and return to normal CPU speed.",
	}, s.handleClearCPUEmulation)

	// === Quality Auditing ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "lighthouse_audit",
		Description: "Run a Lighthouse quality audit on the current page. Returns scores for accessibility, SEO, best-practices, and optionally performance. Requires lighthouse CLI (npm install -g lighthouse).",
	}, s.handleLighthouseAudit)

	// === Network Request Bodies (CDP) ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "get_network_request_body",
		Description: "Get the response body for a network request by ID. Use get_network_requests to find request IDs. Requires CDP connection. Can save binary content to file.",
	}, s.handleGetNetworkRequestBody)

	// === Screencast (CDP) ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "start_screencast",
		Description: "Start capturing screen frames. Frames are captured as base64-encoded images. Requires CDP connection.",
	}, s.handleStartScreencast)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "stop_screencast",
		Description: "Stop capturing screen frames. Requires CDP connection.",
	}, s.handleStopScreencast)

	// === Extensions Management (CDP) ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "install_extension",
		Description: "Load an unpacked extension from a directory. Returns the extension ID. Requires CDP connection.",
	}, s.handleInstallExtension)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "uninstall_extension",
		Description: "Remove an extension by ID. Requires CDP connection.",
	}, s.handleUninstallExtension)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "list_extensions",
		Description: "Get all installed browser extensions with their IDs, names, and status. Requires CDP connection.",
	}, s.handleListExtensions)

	// === Code Coverage (CDP) ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "start_coverage",
		Description: "Start collecting JavaScript and CSS code coverage. Navigate to pages after starting to capture coverage. Requires CDP connection.",
	}, s.handleStartCoverage)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "stop_coverage",
		Description: "Stop collecting code coverage and return results with summary statistics. Requires CDP connection.",
	}, s.handleStopCoverage)

	// === Enhanced Console Debugging (CDP) ===

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "enable_console_debugger",
		Description: "Start capturing console messages with full stack traces. Uses CDP for enhanced debugging beyond standard BiDi console events. Requires CDP connection.",
	}, s.handleEnableConsoleDebugger)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "get_console_entries_with_stack",
		Description: "Get console messages with full stack traces. Call enable_console_debugger first. Optionally filter by type (log, error, warning, etc.).",
	}, s.handleGetConsoleEntriesWithStack)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "get_browser_logs",
		Description: "Get browser log entries including deprecations, interventions, and violations. Call enable_console_debugger first.",
	}, s.handleGetBrowserLogs)

	mcp.AddTool(s.mcpServer, &mcp.Tool{
		Name:        "disable_console_debugger",
		Description: "Stop capturing console messages with stack traces.",
	}, s.handleDisableConsoleDebugger)
}

// Run starts the MCP server.
func (s *Server) Run(ctx context.Context) error {
	return s.mcpServer.Run(ctx, &mcp.StdioTransport{})
}

// Close closes the server and browser session.
func (s *Server) Close(ctx context.Context) error {
	return s.session.Close(ctx)
}
