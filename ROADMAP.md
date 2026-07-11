# Roadmap

This document outlines planned features and improvements for w3pilot.

## v0.2.0 - Enhanced Element Interaction

### Element Selectors

- [ ] XPath selector support (`FindByXPath`)
- [ ] Text-based selectors (`FindByText`, `FindByRole`)
- [ ] Chained selectors (find within element)

### Wait Utilities

- [ ] `WaitForSelector` - Wait for element to appear/disappear
- [ ] `WaitForFunction` - Wait for JavaScript condition
- [ ] `WaitForNetworkIdle` - Wait for network activity to settle

### Keyboard & Mouse

- [ ] `Press(key)` - Press keyboard key (Enter, Escape, Tab, etc.)
- [ ] `KeyboardShortcut(keys)` - Press key combinations (Ctrl+A, Cmd+C)
- [ ] `Hover()` - Mouse hover over element
- [ ] `DragAndDrop(target)` - Drag element to target

## v0.3.0 - Advanced Browser Control

### Frame Support

- [ ] `SwitchToFrame(selector)` - Switch to iframe context
- [ ] `SwitchToMainFrame()` - Return to main frame
- [ ] `Frames()` - List available frames

### Multi-Tab Support

- [ ] `NewTab()` - Open new browser tab
- [ ] `Tabs()` - List open tabs
- [ ] `SwitchToTab(index)` - Switch between tabs
- [ ] `CloseTab()` - Close current tab

### Cookie Management

- [ ] `GetCookies()` - Get all cookies
- [ ] `SetCookie(cookie)` - Set a cookie
- [ ] `DeleteCookies()` - Clear cookies

## v0.4.0 - Network & Storage

### Network Interception

- [ ] `OnRequest(handler)` - Intercept outgoing requests
- [ ] `OnResponse(handler)` - Intercept incoming responses
- [ ] `Mock(url, response)` - Mock network responses
- [ ] `Block(patterns)` - Block requests matching patterns

### File Operations

- [ ] `Upload(selector, path)` - Upload file to input
- [ ] `Download(url)` - Download file
- [ ] `SetDownloadPath(path)` - Configure download directory

### Storage

- [ ] `LocalStorage()` - Access localStorage
- [ ] `SessionStorage()` - Access sessionStorage
- [ ] `ClearStorage()` - Clear browser storage

## v0.5.0 - Output & Export

### PDF Generation

- [ ] `PDF()` - Save page as PDF
- [ ] `PDFOptions` - Configure page size, margins, headers/footers

### Video Recording

- [ ] `StartRecording()` - Record browser session
- [ ] `StopRecording()` - Stop and save recording

### Tracing

- [ ] `StartTracing()` - Start performance trace
- [ ] `StopTracing()` - Stop and export trace

## v0.6.0 - Source Mapping & Framework Detection

These features enable agentic tools (like agent-a11y) to map DOM findings to source code.

### Source Map Integration

- [ ] `LoadSourceMaps(dir)` - Load source maps from directory
- [ ] `GetSourceLocation(selector)` - Map DOM element to source file:line
- [ ] `SetSourceMapURL(url)` - Load source maps from dev server
- [ ] Support for Webpack, Vite, esbuild source map formats

### Framework Detection

- [ ] `DetectFramework()` - Detect React, Vue, Svelte, Angular from page
- [ ] `GetReactComponent(selector)` - Get React component name for element
- [ ] `GetVueComponent(selector)` - Get Vue component name for element
- [ ] `GetDevToolsData(selector)` - Extract component info from browser devtools

### Element Fingerprinting

- [ ] `GetElementFingerprint(selector)` - Generate stable ID for element
- [ ] `FindByFingerprint(fingerprint)` - Locate element by fingerprint
- [ ] `CompareElements(before, after)` - Compare element state across audits

### Use Cases

| Feature | Consumer | Purpose |
|---------|----------|---------|
| Source mapping | agent-a11y | Map a11y findings to source files for coding agents |
| Framework detection | agent-a11y | Generate framework-specific fix code |
| Element fingerprinting | agent-a11y | Track findings across before/after audits |

## Infrastructure

### CI/CD

- [ ] GitHub Actions workflow for CI
- [ ] Automated integration tests in CI
- [ ] Code coverage reporting
- [ ] GoReleaser for binary releases

### Documentation

- [ ] MkDocs documentation site
- [ ] More code examples
- [ ] Tutorial: Web scraping
- [ ] Tutorial: Form automation
- [ ] Tutorial: Testing with Vibium

### Testing

- [ ] Integration tests for additional sites
- [ ] Performance benchmarks
- [ ] Fuzz testing for edge cases

## Contributing

Contributions are welcome! If you'd like to work on any of these items:

1. Open an issue to discuss the feature
2. Fork the repository
3. Create a feature branch
4. Submit a pull request

See [CONTRIBUTING.md](CONTRIBUTING.md) for details.
