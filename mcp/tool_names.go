package mcp

// Tool naming convention:
//
//   {namespace}_{verb}_{target}
//
// Principles:
//   1. Keep verbs explicit: el_get_text not el_text
//   2. Avoid abbreviations: wait_function not wait_fn
//   3. Use full words: human_pause not hitl_pause
//   4. Consistent verb patterns: get/set, start/stop, is/has
//
// Namespaces (20 total):
//   browser_      - Browser lifecycle
//   page_         - Page navigation, state, screenshots
//   element_      - Element interactions and state
//   input_        - Low-level keyboard/mouse/touch
//   js_           - JavaScript execution
//   wait_         - Waiting operations
//   dialog_       - Dialog handling
//   tab_          - Tab management
//   frame_        - Frame selection
//   network_      - Network requests and mocking
//   storage_      - Cookies, localStorage, sessionStorage
//   console_      - Console messages
//   trace_        - Tracing
//   record_       - Script recording
//   test_         - Assertions, verification, reporting
//   accessibility_- Accessibility
//   video_        - Video recording
//   human_        - Human-in-the-loop
//   config_       - Configuration
//   cdp_          - Chrome DevTools Protocol tools
//
// Examples:
//   browser_launch, browser_quit
//   page_navigate, page_go_back, page_screenshot
//   element_click, element_fill, element_get_text, element_is_visible
//   input_keyboard_press, input_mouse_click
//   cdp_take_heap_snapshot, cdp_run_lighthouse

// ToolNames defines the canonical tool names.
// This serves as the single source of truth for all tool names.
var ToolNames = struct {
	// Browser
	BrowserLaunch string
	BrowserQuit   string

	// Page - Navigation
	PageNavigate  string
	PageGoBack    string
	PageGoForward string
	PageReload    string
	PageScroll    string

	// Page - State
	PageGetTitle       string
	PageGetURL         string
	PageGetContent     string
	PageSetContent     string
	PageGetViewport    string
	PageSetViewport    string
	PageScreenshot     string
	PagePDF            string
	PageNew            string
	PageGetCount       string
	PageClose          string
	PageBringToFront   string
	PageEmulateMedia   string
	PageSetGeolocation string

	// Element - Interactions
	ElementClick          string
	ElementDoubleClick    string
	ElementType           string
	ElementFill           string
	ElementFillForm       string
	ElementClear          string
	ElementPress          string
	ElementCheck          string
	ElementUncheck        string
	ElementSelect         string
	ElementSetFiles       string
	ElementHover          string
	ElementFocus          string
	ElementScrollIntoView string
	ElementDragTo         string
	ElementTap            string
	ElementDispatchEvent  string
	ElementScreenshot     string
	ElementEvaluate       string

	// Element - State
	ElementGetText        string
	ElementGetValue       string
	ElementGetInnerHTML   string
	ElementGetOuterHTML   string
	ElementGetInnerText   string
	ElementGetAttribute   string
	ElementGetBoundingBox string
	ElementIsVisible      string
	ElementIsHidden       string
	ElementIsEnabled      string
	ElementIsChecked      string
	ElementIsEditable     string
	ElementGetRole        string
	ElementGetLabel       string

	// Input
	InputKeyboardPress string
	InputKeyboardDown  string
	InputKeyboardUp    string
	InputKeyboardType  string
	InputMouseClick    string
	InputMouseMove     string
	InputMouseDown     string
	InputMouseUp       string
	InputMouseWheel    string
	InputMouseDrag     string
	InputTouchTap      string
	InputTouchSwipe    string

	// JavaScript
	JSEvaluate   string
	JSAddScript  string
	JSAddStyle   string
	JSInitScript string

	// Wait
	WaitForState    string
	WaitForURL      string
	WaitForLoad     string
	WaitForFunction string
	WaitForSelector string
	WaitForText     string

	// Tab
	TabList   string
	TabSelect string
	TabClose  string

	// Frame
	FrameSelect     string
	FrameSelectMain string

	// Dialog
	DialogHandle string
	DialogGet    string

	// Console
	ConsoleGetMessages string
	ConsoleClear       string

	// Network
	NetworkGetRequests string
	NetworkClear       string
	NetworkRoute       string
	NetworkListRoutes  string
	NetworkUnroute     string
	NetworkSetOffline  string

	// Storage - Cookies
	StorageGetCookies   string
	StorageSetCookies   string
	StorageClearCookies string
	StorageDeleteCookie string

	// Storage - State
	StorageGetState string
	StorageSetState string
	StorageClearAll string

	// Storage - Local
	StorageLocalGet    string
	StorageLocalSet    string
	StorageLocalDelete string
	StorageLocalClear  string
	StorageLocalList   string

	// Storage - Session
	StorageSessionGet    string
	StorageSessionSet    string
	StorageSessionDelete string
	StorageSessionClear  string
	StorageSessionList   string

	// Trace
	TraceStart      string
	TraceStop       string
	TraceChunkStart string
	TraceChunkStop  string
	TraceGroupStart string
	TraceGroupStop  string

	// Record
	RecordStart  string
	RecordStop   string
	RecordExport string
	RecordStatus string
	RecordClear  string

	// Test - Assertions
	TestAssertText    string
	TestAssertElement string
	TestAssertURL     string

	// Test - Verification
	TestVerifyValue     string
	TestVerifyList      string
	TestVerifyText      string
	TestVerifyVisible   string
	TestVerifyEnabled   string
	TestVerifyChecked   string
	TestVerifyHidden    string
	TestVerifyDisabled  string
	TestGenerateLocator string
	TestGetReport       string
	TestReset           string
	TestSetTarget       string

	// Accessibility
	AccessibilitySnapshot string

	// Video
	VideoStart string
	VideoStop  string

	// Human-in-the-loop
	HumanPause string

	// Config
	ConfigGet string

	// CDP - Performance
	CDPGetPerformanceMetrics string
	CDPGetMemoryStats        string
	CDPTakeHeapSnapshot      string

	// CDP - Emulation
	CDPEmulateNetwork        string
	CDPClearNetworkEmulation string
	CDPEmulateCPU            string
	CDPClearCPUEmulation     string

	// CDP - Quality
	CDPRunLighthouse string

	// CDP - Coverage
	CDPStartCoverage string
	CDPStopCoverage  string

	// CDP - Console
	CDPEnableConsoleDebugger  string
	CDPGetConsoleEntries      string
	CDPGetBrowserLogs         string
	CDPDisableConsoleDebugger string

	// CDP - Screencast
	CDPStartScreencast string
	CDPStopScreencast  string

	// CDP - Extensions
	CDPInstallExtension   string
	CDPUninstallExtension string
	CDPListExtensions     string

	// CDP - Network
	CDPGetResponseBody string
}{
	// Browser
	BrowserLaunch: "browser_launch",
	BrowserQuit:   "browser_quit",

	// Page - Navigation
	PageNavigate:  "page_navigate",
	PageGoBack:    "page_go_back",
	PageGoForward: "page_go_forward",
	PageReload:    "page_reload",
	PageScroll:    "page_scroll",

	// Page - State
	PageGetTitle:       "page_get_title",
	PageGetURL:         "page_get_url",
	PageGetContent:     "page_get_content",
	PageSetContent:     "page_set_content",
	PageGetViewport:    "page_get_viewport",
	PageSetViewport:    "page_set_viewport",
	PageScreenshot:     "page_screenshot",
	PagePDF:            "page_pdf",
	PageNew:            "page_new",
	PageGetCount:       "page_get_count",
	PageClose:          "page_close",
	PageBringToFront:   "page_bring_to_front",
	PageEmulateMedia:   "page_emulate_media",
	PageSetGeolocation: "page_set_geolocation",

	// Element - Interactions
	ElementClick:          "element_click",
	ElementDoubleClick:    "element_double_click",
	ElementType:           "element_type",
	ElementFill:           "element_fill",
	ElementFillForm:       "element_fill_form",
	ElementClear:          "element_clear",
	ElementPress:          "element_press",
	ElementCheck:          "element_check",
	ElementUncheck:        "element_uncheck",
	ElementSelect:         "element_select",
	ElementSetFiles:       "element_set_files",
	ElementHover:          "element_hover",
	ElementFocus:          "element_focus",
	ElementScrollIntoView: "element_scroll_into_view",
	ElementDragTo:         "element_drag_to",
	ElementTap:            "element_tap",
	ElementDispatchEvent:  "element_dispatch_event",
	ElementScreenshot:     "element_screenshot",
	ElementEvaluate:       "element_evaluate",

	// Element - State
	ElementGetText:        "element_get_text",
	ElementGetValue:       "element_get_value",
	ElementGetInnerHTML:   "element_get_inner_html",
	ElementGetOuterHTML:   "element_get_outer_html",
	ElementGetInnerText:   "element_get_inner_text",
	ElementGetAttribute:   "element_get_attribute",
	ElementGetBoundingBox: "element_get_bounding_box",
	ElementIsVisible:      "element_is_visible",
	ElementIsHidden:       "element_is_hidden",
	ElementIsEnabled:      "element_is_enabled",
	ElementIsChecked:      "element_is_checked",
	ElementIsEditable:     "element_is_editable",
	ElementGetRole:        "element_get_role",
	ElementGetLabel:       "element_get_label",

	// Input
	InputKeyboardPress: "input_keyboard_press",
	InputKeyboardDown:  "input_keyboard_down",
	InputKeyboardUp:    "input_keyboard_up",
	InputKeyboardType:  "input_keyboard_type",
	InputMouseClick:    "input_mouse_click",
	InputMouseMove:     "input_mouse_move",
	InputMouseDown:     "input_mouse_down",
	InputMouseUp:       "input_mouse_up",
	InputMouseWheel:    "input_mouse_wheel",
	InputMouseDrag:     "input_mouse_drag",
	InputTouchTap:      "input_touch_tap",
	InputTouchSwipe:    "input_touch_swipe",

	// JavaScript
	JSEvaluate:   "js_evaluate",
	JSAddScript:  "js_add_script",
	JSAddStyle:   "js_add_style",
	JSInitScript: "js_init_script",

	// Wait
	WaitForState:    "wait_for_state",
	WaitForURL:      "wait_for_url",
	WaitForLoad:     "wait_for_load",
	WaitForFunction: "wait_for_function",
	WaitForSelector: "wait_for_selector",
	WaitForText:     "wait_for_text",

	// Tab
	TabList:   "tab_list",
	TabSelect: "tab_select",
	TabClose:  "tab_close",

	// Frame
	FrameSelect:     "frame_select",
	FrameSelectMain: "frame_select_main",

	// Dialog
	DialogHandle: "dialog_handle",
	DialogGet:    "dialog_get",

	// Console
	ConsoleGetMessages: "console_get_messages",
	ConsoleClear:       "console_clear",

	// Network
	NetworkGetRequests: "network_get_requests",
	NetworkClear:       "network_clear",
	NetworkRoute:       "network_route",
	NetworkListRoutes:  "network_list_routes",
	NetworkUnroute:     "network_unroute",
	NetworkSetOffline:  "network_set_offline",

	// Storage - Cookies
	StorageGetCookies:   "storage_get_cookies",
	StorageSetCookies:   "storage_set_cookies",
	StorageClearCookies: "storage_clear_cookies",
	StorageDeleteCookie: "storage_delete_cookie",

	// Storage - State
	StorageGetState: "storage_get_state",
	StorageSetState: "storage_set_state",
	StorageClearAll: "storage_clear_all",

	// Storage - Local
	StorageLocalGet:    "storage_local_get",
	StorageLocalSet:    "storage_local_set",
	StorageLocalDelete: "storage_local_delete",
	StorageLocalClear:  "storage_local_clear",
	StorageLocalList:   "storage_local_list",

	// Storage - Session
	StorageSessionGet:    "storage_session_get",
	StorageSessionSet:    "storage_session_set",
	StorageSessionDelete: "storage_session_delete",
	StorageSessionClear:  "storage_session_clear",
	StorageSessionList:   "storage_session_list",

	// Trace
	TraceStart:      "trace_start",
	TraceStop:       "trace_stop",
	TraceChunkStart: "trace_chunk_start",
	TraceChunkStop:  "trace_chunk_stop",
	TraceGroupStart: "trace_group_start",
	TraceGroupStop:  "trace_group_stop",

	// Record
	RecordStart:  "record_start",
	RecordStop:   "record_stop",
	RecordExport: "record_export",
	RecordStatus: "record_get_status",
	RecordClear:  "record_clear",

	// Test - Assertions
	TestAssertText:    "test_assert_text",
	TestAssertElement: "test_assert_element",
	TestAssertURL:     "test_assert_url",

	// Test - Verification
	TestVerifyValue:     "test_verify_value",
	TestVerifyList:      "test_verify_list",
	TestVerifyText:      "test_verify_text",
	TestVerifyVisible:   "test_verify_visible",
	TestVerifyEnabled:   "test_verify_enabled",
	TestVerifyChecked:   "test_verify_checked",
	TestVerifyHidden:    "test_verify_hidden",
	TestVerifyDisabled:  "test_verify_disabled",
	TestGenerateLocator: "test_generate_locator",
	TestGetReport:       "test_get_report",
	TestReset:           "test_reset",
	TestSetTarget:       "test_set_target",

	// Accessibility
	AccessibilitySnapshot: "accessibility_snapshot",

	// Video
	VideoStart: "video_start",
	VideoStop:  "video_stop",

	// Human-in-the-loop
	HumanPause: "human_pause",

	// Config
	ConfigGet: "config_get",

	// CDP - Performance
	CDPGetPerformanceMetrics: "cdp_get_performance_metrics",
	CDPGetMemoryStats:        "cdp_get_memory_stats",
	CDPTakeHeapSnapshot:      "cdp_take_heap_snapshot",

	// CDP - Emulation
	CDPEmulateNetwork:        "cdp_emulate_network",
	CDPClearNetworkEmulation: "cdp_clear_network_emulation",
	CDPEmulateCPU:            "cdp_emulate_cpu",
	CDPClearCPUEmulation:     "cdp_clear_cpu_emulation",

	// CDP - Quality
	CDPRunLighthouse: "cdp_run_lighthouse",

	// CDP - Coverage
	CDPStartCoverage: "cdp_start_coverage",
	CDPStopCoverage:  "cdp_stop_coverage",

	// CDP - Console
	CDPEnableConsoleDebugger:  "cdp_enable_console_debugger",
	CDPGetConsoleEntries:      "cdp_get_console_entries",
	CDPGetBrowserLogs:         "cdp_get_browser_logs",
	CDPDisableConsoleDebugger: "cdp_disable_console_debugger",

	// CDP - Screencast
	CDPStartScreencast: "cdp_start_screencast",
	CDPStopScreencast:  "cdp_stop_screencast",

	// CDP - Extensions
	CDPInstallExtension:   "cdp_install_extension",
	CDPUninstallExtension: "cdp_uninstall_extension",
	CDPListExtensions:     "cdp_list_extensions",

	// CDP - Network
	CDPGetResponseBody: "cdp_get_response_body",
}
