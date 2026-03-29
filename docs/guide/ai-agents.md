# AI Agent Ergonomics

W3Pilot includes features designed specifically for AI agents to interact with web pages more effectively.

## Overview

AI agents face unique challenges when automating browsers:

- **Discovery**: Finding interactive elements without prior knowledge of the page
- **Validation**: Verifying selectors before attempting interactions
- **Recovery**: Understanding what went wrong when actions fail
- **State Management**: Saving and restoring browser state across sessions

W3Pilot addresses these with purpose-built tools.

## Page Inspection

The `Inspect` method scans the page and returns a structured inventory of interactive elements:

```go
result, err := pilot.Inspect(ctx, nil)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Page: %s\n", result.Title)
fmt.Printf("URL: %s\n", result.URL)

// Discover buttons
for _, btn := range result.Buttons {
    fmt.Printf("Button: %s - %s\n", btn.Selector, btn.Text)
}

// Discover inputs
for _, input := range result.Inputs {
    fmt.Printf("Input: %s (%s) - %s\n",
        input.Selector, input.Type, input.Label)
}
```

### Inspection Options

Control what elements to include:

```go
opts := &w3pilot.InspectOptions{
    IncludeButtons:  true,
    IncludeLinks:    true,
    IncludeInputs:   true,
    IncludeSelects:  true,
    IncludeHeadings: true,
    IncludeImages:   true,
    MaxItems:        50,  // Per category
}

result, err := pilot.Inspect(ctx, opts)
```

### Inspect Result Structure

```go
type InspectResult struct {
    URL      string
    Title    string
    Buttons  []InspectButton   // button, input[type=submit], [role=button]
    Links    []InspectLink     // a[href]
    Inputs   []InspectInput    // input, textarea
    Selects  []InspectSelect   // select
    Headings []InspectHeading  // h1-h6
    Images   []InspectImage    // img[alt]
    Summary  InspectSummary
}

type InspectButton struct {
    Selector string
    Text     string
    Type     string  // button, submit, reset
    Disabled bool
    Visible  bool
}

type InspectInput struct {
    Selector    string
    Type        string  // text, password, email, etc.
    Name        string
    Placeholder string
    Value       string  // Masked for password fields
    Label       string  // Associated label text
    Required    bool
    Disabled    bool
    ReadOnly    bool
    Visible     bool
}
```

### MCP Tool

```json
{
    "name": "page_inspect",
    "arguments": {
        "include_buttons": true,
        "include_links": false,
        "max_items": 25
    }
}
```

### CLI Command

```bash
# Full inspection
w3pilot page inspect

# JSON output for parsing
w3pilot page inspect --format json

# Exclude categories
w3pilot page inspect --no-links --no-images

# Limit items
w3pilot page inspect --max-items 100
```

## Selector Validation

Before interacting with elements, validate selectors to check existence and state:

```go
results, err := pilot.ValidateSelectors(ctx, []string{
    "#login-button",
    "#username",
    "#password",
    "#nonexistent",
})

for _, v := range results {
    if v.Found {
        fmt.Printf("FOUND: %s (tag=%s, visible=%t, enabled=%t)\n",
            v.Selector, v.TagName, v.Visible, v.Enabled)
    } else {
        fmt.Printf("NOT FOUND: %s\n", v.Selector)
        if len(v.Suggestions) > 0 {
            fmt.Printf("  Suggestions: %v\n", v.Suggestions)
        }
    }
}
```

### Validation Result

```go
type SelectorValidation struct {
    Selector    string
    Found       bool
    Count       int      // Number of matching elements
    Visible     bool     // First match is visible
    Enabled     bool     // First match is enabled
    TagName     string   // Tag name of first match
    Suggestions []string // Alternative selectors if not found
}
```

### MCP Tool

```json
{
    "name": "test_validate_selectors",
    "arguments": {
        "selectors": ["#login", ".submit-btn", "button[type=submit]"]
    }
}
```

### CLI Command

```bash
w3pilot test validate-selectors "#login" ".password" "button[type=submit]"
w3pilot test validate-selectors --format json "#nonexistent"
```

## Workflow Recipes

High-level workflows combine multiple steps for common tasks.

### Login Workflow

Automated login with credential filling and success verification:

```go
result, err := pilot.Login(ctx, &w3pilot.LoginOptions{
    UsernameSelector: "#email",
    PasswordSelector: "#password",
    SubmitSelector:   "button[type=submit]",
    Username:         "user@example.com",
    Password:         "secret123",
    SuccessIndicator: "/dashboard",  // URL pattern
    Timeout:          30 * time.Second,
})

if result.Success {
    fmt.Printf("Logged in! Now at: %s\n", result.URL)
} else {
    fmt.Printf("Login failed: %s\n", result.ErrorReason)
}
```

The `SuccessIndicator` can be:

- A URL pattern (starts with `/`, `*`, or `http`)
- A CSS selector (element that appears after successful login)

### MCP Tool

```json
{
    "name": "workflow_login",
    "arguments": {
        "username_selector": "#email",
        "password_selector": "#password",
        "submit_selector": "button[type=submit]",
        "username": "user@example.com",
        "password": "secret123",
        "success_indicator": "#dashboard-header"
    }
}
```

### Table Extraction

Extract HTML table data to structured JSON:

```go
table, err := pilot.ExtractTable(ctx, "table.data", &w3pilot.ExtractTableOptions{
    IncludeHeaders: true,
    MaxRows:        100,
})

// Access as arrays
for _, row := range table.Rows {
    fmt.Println(row) // []string{"value1", "value2", ...}
}

// Access as objects (keyed by headers)
for _, row := range table.RowsJSON {
    fmt.Printf("Name: %s, Email: %s\n", row["Name"], row["Email"])
}
```

### MCP Tool

```json
{
    "name": "workflow_extract_table",
    "arguments": {
        "selector": "table#users",
        "include_headers": true,
        "max_rows": 50
    }
}
```

## Named State Snapshots

Save and restore browser state (cookies, localStorage, sessionStorage) with named snapshots.

### SDK Usage

```go
// Save current state
state, err := pilot.StorageState(ctx)
if err != nil {
    log.Fatal(err)
}

mgr, _ := state.NewManager("")  // Uses ~/.w3pilot/states/

// Save to named snapshot
err = mgr.Save("logged-in-session", state)

// Later, restore the state
savedState, err := mgr.Load("logged-in-session")
pilot.SetStorageState(ctx, savedState)
pilot.Reload(ctx)  // Apply the restored state
```

### MCP Tools

```json
// Save current state
{"name": "state_save", "arguments": {"name": "my-session"}}

// List saved states
{"name": "state_list"}

// Load a saved state
{"name": "state_load", "arguments": {"name": "my-session"}}

// Delete a saved state
{"name": "state_delete", "arguments": {"name": "old-session"}}
```

### CLI Commands

```bash
# Save current browser state
w3pilot state save my-session

# List all saved states
w3pilot state list

# Load a saved state
w3pilot state load my-session

# Delete a saved state
w3pilot state delete my-session
```

### State Storage Location

States are saved to `~/.w3pilot/states/{name}.json` with metadata:

```json
{
    "name": "my-session",
    "created_at": "2026-03-28T10:30:00Z",
    "state": {
        "cookies": [...],
        "origins": [
            {
                "origin": "https://example.com",
                "localStorage": {"key": "value"},
                "sessionStorage": {"key": "value"}
            }
        ]
    }
}
```

## Output Formatting

All CLI commands that return data support the `--format` flag:

```bash
# Default text output
w3pilot page title
# Example Domain

# JSON output
w3pilot page title --format json
# {"title":"Example Domain"}

# JSON for element queries
w3pilot element text "h1" --format json
# {"selector":"h1","text":"Example Domain"}
```

This enables AI agents to reliably parse command output.

## Enhanced Error Context

When operations fail, errors include page context to help diagnose issues:

```go
type ElementNotFoundError struct {
    Selector    string
    PageContext *PageContext
    Suggestions []string
}

type PageContext struct {
    URL         string
    Title       string
    VisibleText string
}
```

This helps AI agents understand:

- What URL they're on when the error occurred
- What similar selectors might work instead
- The current page state for debugging

## Best Practices for AI Agents

### 1. Inspect Before Acting

```go
// First, understand the page
result, _ := pilot.Inspect(ctx, nil)

// Then interact with discovered elements
for _, input := range result.Inputs {
    if input.Type == "email" {
        elem, _ := pilot.Find(ctx, input.Selector, nil)
        elem.Fill(ctx, "user@example.com", nil)
        break
    }
}
```

### 2. Validate Selectors

```go
// Check selectors before using them
validations, _ := pilot.ValidateSelectors(ctx, selectors)

for _, v := range validations {
    if !v.Found {
        // Use suggestions or try alternatives
        if len(v.Suggestions) > 0 {
            selectors = append(selectors, v.Suggestions[0])
        }
    }
}
```

### 3. Use State Snapshots for Repeatability

```bash
# After manual login, save the state
w3pilot state save authenticated

# In automation, load the saved state
w3pilot state load authenticated
w3pilot page navigate https://example.com/dashboard
```

### 4. Use JSON Output for Parsing

```bash
# Always use --format json when parsing output programmatically
result=$(w3pilot page inspect --format json)
buttons=$(echo "$result" | jq '.buttons')
```
