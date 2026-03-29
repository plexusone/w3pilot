# Quick Start

## Go Client SDK

```go
package main

import (
    "context"
    "fmt"
    "log"

github.com/plexusone/w3pilot
)

func main() {
    ctx := context.Background()

    // Launch browser
    pilot, err := w3pilot.Launch(ctx)
    if err != nil {
        log.Fatal(err)
    }
    defer pilot.Quit(ctx)

    // Navigate
    if err := pilot.Go(ctx, "https://example.com"); err != nil {
        log.Fatal(err)
    }

    // Get page title
    title, _ := pilot.Title(ctx)
    fmt.Println("Title:", title)

    // Find and click a link
    link, err := pilot.Find(ctx, "a", nil)
    if err != nil {
        log.Fatal(err)
    }

    if err := link.Click(ctx, nil); err != nil {
        log.Fatal(err)
    }

    // Take screenshot
    data, _ := pilot.Screenshot(ctx)
    os.WriteFile("screenshot.png", data, 0644)
}
```

## MCP Server

Start the server:

```bash
w3pilot mcp --headless
```

Configure in Claude Desktop (`~/Library/Application Support/Claude/claude_desktop_config.json`):

```json
{
  "mcpServers": {
    "w3pilot": {
      "command": "w3pilot",
      "args": ["mcp", "--headless"]
    }
  }
}
```

Then ask Claude: "Navigate to example.com and take a screenshot"

## CLI

Interactive browser control:

```bash
# Launch browser
w3pilot launch

# Navigate
w3pilot go https://example.com

# Interact
w3pilot fill "#search" "hello world"
w3pilot click "#submit"

# Capture
w3pilot screenshot result.png

# Cleanup
w3pilot quit
```

## Script Runner

Create `test.json`:

```json
{
  "name": "Example Test",
  "steps": [
    {"action": "navigate", "url": "https://example.com"},
    {"action": "assertTitle", "expected": "Example Domain"},
    {"action": "screenshot", "file": "result.png"}
  ]
}
```

Run:

```bash
w3pilot run test.json
```
