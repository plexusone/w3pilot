# Feature Comparison

This document compares vibium-go against other Vibium clients (Java, JavaScript, Python) and the Playwright MCP server.

## Overview

vibium-go is a Go port of the official Vibium clients. This comparison helps track feature parity and identify gaps.

| Client | Language | Repository |
|--------|----------|------------|
| vibium-go | Go | [plexusone/vibium-go](https://github.com/plexusone/vibium-go) |
| vibium-java | Java | [VibiumDev/vibium/clients/java](https://github.com/VibiumDev/vibium) |
| vibium-js | JavaScript | [VibiumDev/vibium/clients/javascript](https://github.com/VibiumDev/vibium) |
| vibium-py | Python | [VibiumDev/vibium/clients/python](https://github.com/VibiumDev/vibium) |
| playwright-mcp | TypeScript | [microsoft/playwright-mcp](https://github.com/microsoft/playwright-mcp) |

## SDK Feature Comparison

### Core Browser Control

| Feature | Go | Java | JS | Python | Notes |
|---------|:--:|:----:|:--:|:------:|-------|
| Launch browser | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| Headless mode | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| Connect to remote | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| Multiple pages | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| Browser contexts | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| Page events (onPage, onPopup) | :x: | :white_check_mark: | :white_check_mark: | :white_check_mark: | Planned |

### Navigation

| Feature | Go | Java | JS | Python | Notes |
|---------|:--:|:----:|:--:|:------:|-------|
| Navigate (go) | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| Back/Forward | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| Reload | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| URL/Title | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| Content (HTML) | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| setContent | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| bringToFront | :x: | :white_check_mark: | :white_check_mark: | :white_check_mark: | Planned |

### Element Finding

| Feature | Go | Java | JS | Python | Notes |
|---------|:--:|:----:|:--:|:------:|-------|
| CSS selector | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| **Semantic selectors** | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| - By role | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| - By text | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| - By label | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| - By placeholder | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| - By alt text | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| - By title | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| - By testid | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| - By xpath | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| - By proximity (near) | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| FindAll | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| Scoped find (within element) | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |

### Element Interaction

| Feature | Go | Java | JS | Python | Notes |
|---------|:--:|:----:|:--:|:------:|-------|
| Click | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| Double-click | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| Fill | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| Type | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| Clear | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| Press | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| Check/Uncheck | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| Select option | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| Hover | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| Focus/Blur | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| Drag to | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| Tap (touch) | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| Scroll into view | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| Set files | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| dispatchEvent | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| highlight | :x: | :white_check_mark: | :x: | :x: | Java only |

### Element State

| Feature | Go | Java | JS | Python | Notes |
|---------|:--:|:----:|:--:|:------:|-------|
| text() | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| innerText() | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| html() (innerHTML) | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| value() | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| getAttribute() | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| boundingBox() | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| isVisible/isHidden | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| isEnabled | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| isChecked | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| isEditable | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| role() | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| label() | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| waitUntil (state) | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |

### Input Controllers

| Feature | Go | Java | JS | Python | Notes |
|---------|:--:|:----:|:--:|:------:|-------|
| **Keyboard** | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| - press/down/up | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| - type | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| - insertText | :white_check_mark: | :x: | :x: | :x: | Go only |
| **Mouse** | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| - click/move | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| - down/up | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| - wheel | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| **Touch** | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| - tap | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| - drag/longPress | :white_check_mark: | :x: | :x: | :x: | Go only |

### Clock/Time Control

| Feature | Go | Java | JS | Python | Notes |
|---------|:--:|:----:|:--:|:------:|-------|
| install | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| fastForward | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| runFor | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| pauseAt | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| resume | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| setFixedTime | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| setSystemTime | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| setTimezone | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |

### Screenshots & PDF

| Feature | Go | Java | JS | Python | Notes |
|---------|:--:|:----:|:--:|:------:|-------|
| Screenshot (viewport) | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| Screenshot (full page) | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| Screenshot (clip region) | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| Element screenshot | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| PDF | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |

### Viewport & Window

| Feature | Go | Java | JS | Python | Notes |
|---------|:--:|:----:|:--:|:------:|-------|
| getViewport | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| setViewport | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| getWindow | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| setWindow | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |

### Media & Emulation

| Feature | Go | Java | JS | Python | Notes |
|---------|:--:|:----:|:--:|:------:|-------|
| emulateMedia (type) | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| colorScheme | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| reducedMotion | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| forcedColors | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| contrast | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| setGeolocation | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |

### Accessibility

| Feature | Go | Java | JS | Python | Notes |
|---------|:--:|:----:|:--:|:------:|-------|
| a11yTree | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| a11yTree (interestingOnly) | :x: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| a11yTree (root element) | :x: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |

### Frames

| Feature | Go | Java | JS | Python | Notes |
|---------|:--:|:----:|:--:|:------:|-------|
| frames() | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| frame(name/url) | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| mainFrame() | :x: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |

### Network Interception

| Feature | Go | Java | JS | Python | Notes |
|---------|:--:|:----:|:--:|:------:|-------|
| route(pattern, handler) | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| unroute | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| route.fulfill | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| route.continue | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| route.abort | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| setHeaders | :x: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| onRequest listener | :x: | :white_check_mark: | :white_check_mark: | :white_check_mark: | Planned |
| onResponse listener | :x: | :white_check_mark: | :white_check_mark: | :white_check_mark: | Planned |

### Event Capture (JS/Python)

| Feature | Go | Java | JS | Python | Notes |
|---------|:--:|:----:|:--:|:------:|-------|
| capture.response | :x: | :x: | :white_check_mark: | :white_check_mark: | JS/Py only |
| capture.request | :x: | :x: | :white_check_mark: | :white_check_mark: | JS/Py only |
| capture.navigation | :x: | :x: | :white_check_mark: | :white_check_mark: | JS/Py only |
| capture.download | :x: | :x: | :white_check_mark: | :white_check_mark: | JS/Py only |
| capture.dialog | :x: | :x: | :white_check_mark: | :white_check_mark: | JS/Py only |

### Dialogs

| Feature | Go | Java | JS | Python | Notes |
|---------|:--:|:----:|:--:|:------:|-------|
| onDialog listener | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| dialog.accept | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| dialog.dismiss | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| dialog.message | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| dialog.type | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| dialog.defaultValue | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |

### Downloads

| Feature | Go | Java | JS | Python | Notes |
|---------|:--:|:----:|:--:|:------:|-------|
| onDownload listener | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| download.url | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| download.suggestedFilename | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| download.saveAs | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| download.path | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |

### Console & Errors

| Feature | Go | Java | JS | Python | Notes |
|---------|:--:|:----:|:--:|:------:|-------|
| onConsole listener | :x: | :white_check_mark: | :white_check_mark: | :white_check_mark: | Planned |
| collectConsole (buffered) | :x: | :white_check_mark: | :white_check_mark: | :white_check_mark: | Planned |
| consoleMessages() | :x: | :white_check_mark: | :white_check_mark: | :white_check_mark: | Planned |
| onError listener | :x: | :white_check_mark: | :white_check_mark: | :white_check_mark: | Planned |
| collectErrors (buffered) | :x: | :white_check_mark: | :white_check_mark: | :white_check_mark: | Planned |
| errors() | :x: | :white_check_mark: | :white_check_mark: | :white_check_mark: | Planned |

### WebSocket Monitoring

| Feature | Go | Java | JS | Python | Notes |
|---------|:--:|:----:|:--:|:------:|-------|
| onWebSocket listener | :x: | :white_check_mark: | :white_check_mark: | :white_check_mark: | Planned |
| ws.onMessage | :x: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| ws.onClose | :x: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| ws.isClosed | :x: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |

### Recording/Tracing

| Feature | Go | Java | JS | Python | Notes |
|---------|:--:|:----:|:--:|:------:|-------|
| recording.start | :x: | :white_check_mark: | :white_check_mark: | :white_check_mark: | Planned |
| recording.stop | :x: | :white_check_mark: | :white_check_mark: | :white_check_mark: | Planned |
| recording.startChunk | :x: | :white_check_mark: | :white_check_mark: | :white_check_mark: | Planned |
| recording.stopChunk | :x: | :white_check_mark: | :white_check_mark: | :white_check_mark: | Planned |
| recording.startGroup | :x: | :white_check_mark: | :white_check_mark: | :white_check_mark: | Planned |
| recording.stopGroup | :x: | :white_check_mark: | :white_check_mark: | :white_check_mark: | Planned |

### Storage State

| Feature | Go | Java | JS | Python | Notes |
|---------|:--:|:----:|:--:|:------:|-------|
| Cookies | | | | | |
| - cookies() | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| - setCookies() | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| - clearCookies() | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| Storage | | | | | |
| - storage() (full state) | :x: | :white_check_mark: | :white_check_mark: | :white_check_mark: | Planned |
| - setStorage() | :x: | :white_check_mark: | :white_check_mark: | :white_check_mark: | Planned |
| - clearStorage() | :x: | :white_check_mark: | :white_check_mark: | :white_check_mark: | Planned |
| Init Scripts | | | | | |
| - addInitScript() | :x: | :white_check_mark: | :white_check_mark: | :white_check_mark: | Planned |

### Scrolling

| Feature | Go | Java | JS | Python | Notes |
|---------|:--:|:----:|:--:|:------:|-------|
| scroll() | :x: | :white_check_mark: | :white_check_mark: | :white_check_mark: | Planned |
| scroll (direction) | :x: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| scroll (amount) | :x: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| scroll (selector) | :x: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |

### JavaScript Execution

| Feature | Go | Java | JS | Python | Notes |
|---------|:--:|:----:|:--:|:------:|-------|
| evaluate() | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| addScript() | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| addStyle() | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| expose() | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |

### Waiting

| Feature | Go | Java | JS | Python | Notes |
|---------|:--:|:----:|:--:|:------:|-------|
| Wait (fixed duration) | :x: | :white_check_mark: | :white_check_mark: | :white_check_mark: | Use time.Sleep |
| waitForFunction | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| waitForURL | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| waitForLoad | :white_check_mark: | :white_check_mark: | :white_check_mark: | :white_check_mark: | |
| WaitForNavigation | :white_check_mark: | :x: | :x: | :x: | Go only |

---

## MCP Server Feature Comparison

Comparison between vibium-go MCP server and [Playwright MCP](https://github.com/microsoft/playwright-mcp).

### Browser Management

| Tool | vibium-go | Playwright MCP | Notes |
|------|:---------:|:--------------:|-------|
| Launch browser | :white_check_mark: | :white_check_mark: | |
| Quit/Close browser | :white_check_mark: | :white_check_mark: | |
| Resize viewport | :white_check_mark: | :white_check_mark: | |

### Navigation

| Tool | vibium-go | Playwright MCP | Notes |
|------|:---------:|:--------------:|-------|
| Navigate to URL | :white_check_mark: | :white_check_mark: | |
| Back | :white_check_mark: | :white_check_mark: | |
| Forward | :white_check_mark: | :white_check_mark: | |
| Reload | :white_check_mark: | :white_check_mark: | |

### Element Interaction

| Tool | vibium-go | Playwright MCP | Notes |
|------|:---------:|:--------------:|-------|
| Click | :white_check_mark: | :white_check_mark: | |
| Double-click | :white_check_mark: | :white_check_mark: | |
| Type | :white_check_mark: | :white_check_mark: | |
| Fill | :white_check_mark: | :white_check_mark: | |
| Clear | :white_check_mark: | :x: | |
| Press key | :white_check_mark: | :white_check_mark: | |
| Check/Uncheck | :white_check_mark: | :x: | |
| Select option | :white_check_mark: | :white_check_mark: | |
| Set files | :white_check_mark: | :white_check_mark: | |
| Hover | :white_check_mark: | :white_check_mark: | |
| Drag | :white_check_mark: | :white_check_mark: | |
| Fill form (multiple) | :white_check_mark: | :white_check_mark: | |

### Element State

| Tool | vibium-go | Playwright MCP | Notes |
|------|:---------:|:--------------:|-------|
| Get text | :white_check_mark: | :white_check_mark: | |
| Get value | :white_check_mark: | :white_check_mark: | |
| Get innerHTML | :white_check_mark: | :white_check_mark: | |
| Get innerText | :white_check_mark: | :x: | |
| Get attribute | :white_check_mark: | :white_check_mark: | |
| Get bounding box | :white_check_mark: | :x: | |
| Is visible | :white_check_mark: | :white_check_mark: | |
| Is hidden | :white_check_mark: | :x: | |
| Is enabled | :white_check_mark: | :x: | |
| Is checked | :white_check_mark: | :x: | |
| Is editable | :white_check_mark: | :x: | |
| Get role | :white_check_mark: | :x: | |
| Get label | :white_check_mark: | :x: | |

### Page State

| Tool | vibium-go | Playwright MCP | Notes |
|------|:---------:|:--------------:|-------|
| Get title | :white_check_mark: | :white_check_mark: | |
| Get URL | :white_check_mark: | :white_check_mark: | |
| Get content | :white_check_mark: | :white_check_mark: | |
| Set content | :white_check_mark: | :x: | |
| Get viewport | :white_check_mark: | :x: | |
| Set viewport | :white_check_mark: | :white_check_mark: | |
| Accessibility snapshot | :x: | :white_check_mark: | Playwright uses a11y tree |

### Screenshots & PDF

| Tool | vibium-go | Playwright MCP | Notes |
|------|:---------:|:--------------:|-------|
| Screenshot (viewport) | :white_check_mark: | :white_check_mark: | |
| Screenshot (full page) | :white_check_mark: | :white_check_mark: | |
| Element screenshot | :white_check_mark: | :white_check_mark: | |
| PDF | :white_check_mark: | :white_check_mark: | |

### JavaScript

| Tool | vibium-go | Playwright MCP | Notes |
|------|:---------:|:--------------:|-------|
| Evaluate | :white_check_mark: | :white_check_mark: | |
| Element eval | :white_check_mark: | :x: | |
| Add script | :white_check_mark: | :x: | |
| Add style | :white_check_mark: | :x: | |
| Run Playwright code | :x: | :white_check_mark: | Playwright-specific |

### Waiting

| Tool | vibium-go | Playwright MCP | Notes |
|------|:---------:|:--------------:|-------|
| Wait for element state | :white_check_mark: | :white_check_mark: | |
| Wait for URL | :white_check_mark: | :x: | |
| Wait for load state | :white_check_mark: | :x: | |
| Wait for function | :white_check_mark: | :x: | |
| Wait for text | :x: | :white_check_mark: | Planned |

### Input Controllers

| Tool | vibium-go | Playwright MCP | Notes |
|------|:---------:|:--------------:|-------|
| Keyboard press | :white_check_mark: | :white_check_mark: | |
| Keyboard down/up | :white_check_mark: | :x: | |
| Keyboard type | :white_check_mark: | :x: | |
| Mouse click (coords) | :white_check_mark: | :white_check_mark: | |
| Mouse move | :white_check_mark: | :white_check_mark: | |
| Mouse down/up | :white_check_mark: | :white_check_mark: | |
| Mouse wheel | :white_check_mark: | :white_check_mark: | |
| Mouse drag (coords) | :white_check_mark: | :white_check_mark: | |
| Touch tap | :white_check_mark: | :x: | |
| Touch swipe | :white_check_mark: | :x: | |

### Tab/Page Management

| Tool | vibium-go | Playwright MCP | Notes |
|------|:---------:|:--------------:|-------|
| New page | :white_check_mark: | :white_check_mark: | |
| Get page count | :white_check_mark: | :white_check_mark: | |
| Close page | :white_check_mark: | :white_check_mark: | |
| Bring to front | :white_check_mark: | :x: | |
| Select tab | :white_check_mark: | :white_check_mark: | |
| List tabs | :white_check_mark: | :white_check_mark: | |
| Close tab | :white_check_mark: | :white_check_mark: | |

### Cookie Management

| Tool | vibium-go | Playwright MCP | Notes |
|------|:---------:|:--------------:|-------|
| Get cookies | :white_check_mark: | :white_check_mark: | |
| Set cookies | :white_check_mark: | :white_check_mark: | |
| Clear cookies | :white_check_mark: | :white_check_mark: | |
| Delete cookie | :x: | :white_check_mark: | Planned |

### LocalStorage

| Tool | vibium-go | Playwright MCP | Notes |
|------|:---------:|:--------------:|-------|
| Get item | :white_check_mark: | :white_check_mark: | |
| Set item | :white_check_mark: | :white_check_mark: | |
| List items | :white_check_mark: | :white_check_mark: | |
| Delete item | :white_check_mark: | :white_check_mark: | |
| Clear | :white_check_mark: | :white_check_mark: | |

### SessionStorage

| Tool | vibium-go | Playwright MCP | Notes |
|------|:---------:|:--------------:|-------|
| Get item | :white_check_mark: | :white_check_mark: | |
| Set item | :white_check_mark: | :white_check_mark: | |
| List items | :white_check_mark: | :white_check_mark: | |
| Delete item | :white_check_mark: | :white_check_mark: | |
| Clear | :white_check_mark: | :white_check_mark: | |

### Storage State

| Tool | vibium-go | Playwright MCP | Notes |
|------|:---------:|:--------------:|-------|
| Get storage state | :white_check_mark: | :white_check_mark: | |
| Set storage state | :white_check_mark: | :white_check_mark: | |

### Network

| Tool | vibium-go | Playwright MCP | Notes |
|------|:---------:|:--------------:|-------|
| Route (mock) | :x: | :white_check_mark: | Planned |
| Route list | :x: | :white_check_mark: | Planned |
| Unroute | :x: | :white_check_mark: | Planned |
| Network state (offline) | :x: | :white_check_mark: | Planned |
| List network requests | :white_check_mark: | :white_check_mark: | |
| Clear network requests | :white_check_mark: | :x: | |
| Console messages | :white_check_mark: | :white_check_mark: | |
| Clear console messages | :white_check_mark: | :x: | |

### Dialogs

| Tool | vibium-go | Playwright MCP | Notes |
|------|:---------:|:--------------:|-------|
| Handle dialog | :white_check_mark: | :white_check_mark: | |
| Get dialog info | :white_check_mark: | :x: | |

### Recording & Tracing

| Tool | vibium-go | Playwright MCP | Notes |
|------|:---------:|:--------------:|-------|
| Start recording | :white_check_mark: | :x: | Script recording |
| Stop recording | :white_check_mark: | :x: | |
| Export script | :white_check_mark: | :x: | |
| Start tracing | :x: | :white_check_mark: | Planned |
| Stop tracing | :x: | :white_check_mark: | Planned |
| Start video | :x: | :white_check_mark: | Planned |
| Stop video | :x: | :white_check_mark: | Planned |

### Testing & Assertions

| Tool | vibium-go | Playwright MCP | Notes |
|------|:---------:|:--------------:|-------|
| Assert text | :white_check_mark: | :white_check_mark: | |
| Assert element | :white_check_mark: | :white_check_mark: | |
| Assert URL | :white_check_mark: | :x: | |
| Verify value | :x: | :white_check_mark: | Planned |
| Verify list visible | :x: | :white_check_mark: | Planned |
| Generate locator | :x: | :white_check_mark: | Planned |
| Test report | :white_check_mark: | :x: | |

### Human-in-the-Loop

| Tool | vibium-go | Playwright MCP | Notes |
|------|:---------:|:--------------:|-------|
| Pause for human | :white_check_mark: | :x: | Unique to vibium-go |
| Get storage state | :white_check_mark: | :white_check_mark: | |
| Set storage state | :white_check_mark: | :white_check_mark: | |

### Configuration

| Tool | vibium-go | Playwright MCP | Notes |
|------|:---------:|:--------------:|-------|
| Get config | :x: | :white_check_mark: | Planned |

---

## Feature Parity Roadmap

### High Priority (SDK)

The following features are prioritized for SDK parity:

1. **Recording/Tracing** - Full trace recording with chunks and groups
2. ~~**Media Emulation** - colorScheme, reducedMotion, forcedColors, contrast~~ :white_check_mark: Done
3. **Console/Error Collection** - Buffered console messages and errors
4. **Request/Response Listeners** - Network observation events
5. **Init Scripts** - Per-context script injection
6. **Full Storage State** - localStorage, sessionStorage, indexedDB

### High Priority (MCP)

The following MCP tools are prioritized for Playwright MCP parity:

1. **Network Mocking** - route, route_list, unroute
2. ~~**Tab Management** - list tabs, select tab~~ :white_check_mark: Done
3. ~~**Dialog Handling** - handle_dialog tool~~ :white_check_mark: Done
4. ~~**Console Messages** - console_messages with filtering~~ :white_check_mark: Done
5. ~~**Network Requests** - network_requests listing~~ :white_check_mark: Done

### Unique vibium-go Features

Features in vibium-go not found in other clients:

| Feature | Description |
|---------|-------------|
| Script Runner | YAML/JSON deterministic test execution |
| Session Recording | Capture MCP actions as replayable scripts |
| Human-in-the-Loop | pauseForHuman for SSO, CAPTCHA, 2FA |
| Test Reports | Structured test execution reports |
| RPA Engine | Workflow automation with activities |

---

## Legend

| Symbol | Meaning |
|--------|---------|
| :white_check_mark: | Supported |
| :x: | Not supported |
| Planned | On roadmap |
