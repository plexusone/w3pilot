# Roadmap

This document outlines planned features and improvements for the Vibium Go Client SDK.

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
