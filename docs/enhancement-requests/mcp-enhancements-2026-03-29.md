# MCP Server Enhancement Requests

**Date**: 2026-03-29
**Source**: Real-world usage feedback
**Status**: Accepted for implementation

---

## Summary

| # | Enhancement | Priority | Type | Status |
|---|-------------|----------|------|--------|
| 1 | Await async IIFEs in `js_evaluate` | High | Bug fix | Planned |
| 2 | Fix `state_save`/`state_load` compatibility | Medium | Bug fix | Planned |
| 3 | Add `http_request` tool | High | New feature | Planned |
| 4 | Response truncation for `js_evaluate` | Medium | Enhancement | Planned |
| 5 | Explicit `js_evaluate_async` tool | Low | New feature | Planned (if #1 difficult) |
| 6 | Batch tool execution | Low | Enhancement | Planned |

---

## 1. `js_evaluate` should await async IIFEs

**Priority**: High
**Category**: Bug / behavior gap

`js_evaluate` returns `null` when the script is an async IIFE:

```javascript
// Returns null
(async () => {
  const resp = await fetch('/api/test', {credentials: 'include'});
  return {status: resp.status};
})()
```

The workaround is `.then()` chaining, which does resolve:

```javascript
// Works
fetch('/api/test', {credentials: 'include'}).then(r => ({status: r.status}))
```

**Ask**: If the evaluated expression returns a `Promise`, `await` it before returning the result to the caller. This is how Playwright's `page.evaluate()` works natively — the MCP layer appears to be dropping the promise.

---

## 2. `state_save` / `state_load` fails on some browser contexts

**Priority**: Medium
**Category**: Bug / compatibility

`state_save` fails with:

```
Unknown command 'vibium:context.storageState'
```

This prevents saving and restoring authenticated sessions across browser restarts, which is important for long-running sessions where the browser may crash or need to be relaunched.

**Ask**: Either fall back to a manual cookie/localStorage extraction when the native storage state command isn't available, or document which browser configurations support this feature.

---

## 3. Add a `fetch` / `http_request` tool for authenticated requests

**Priority**: High
**Category**: New feature

A common pattern is: send an HTTP request from the authenticated browser context and inspect the response. Currently this requires wrapping everything in `js_evaluate` + `fetch()` + `.then()`:

```javascript
fetch('/ECM/some/endpoint', {
  method: 'POST',
  credentials: 'include',
  headers: {'Content-Type': 'application/x-www-form-urlencoded'},
  body: 'param=value'
}).then(async r => ({status: r.status, body: (await r.text()).substring(0, 4000)}))
```

A dedicated tool would be cleaner and less error-prone:

```
http_request
  url="/ECM/some/endpoint"
  method="POST"
  content_type="application/x-www-form-urlencoded"
  body="param=value"
  max_body_length=4000
```

**Benefits**:

- No need to manually handle promise resolution quirks
- Automatic credential inclusion from browser context
- Built-in response truncation (large responses can blow up `js_evaluate` results)
- Structured response object (status, headers, body) without manual parsing

---

## 4. Add response body truncation / size control to `js_evaluate`

**Priority**: Medium
**Category**: Enhancement

When fetching endpoint responses via `js_evaluate`, we manually truncate with `.substring(0, N)` to avoid oversized results. If a script accidentally returns a large DOM or response body, it can overwhelm the MCP response channel.

**Ask**: Add an optional `max_result_size` parameter to `js_evaluate` that truncates the serialized result and appends a `[truncated]` indicator.

---

## 5. Add a `js_evaluate_async` tool (or explicit async mode)

**Priority**: Low (if #1 is fixed)
**Category**: New feature

If fixing the async IIFE behavior in `js_evaluate` is complex, an alternative is a separate `js_evaluate_async` tool that explicitly awaits the returned promise with a configurable timeout:

```
js_evaluate_async
  script="(async () => { ... })()"
  timeout_ms=15000
```

This would also be useful for scripts that need to wait for DOM mutations, network responses, or timers.

---

## 6. Batch tool execution

**Priority**: Low
**Category**: Enhancement

Several workflows follow the pattern: navigate → screenshot → evaluate → screenshot. Each step is a separate MCP round-trip. A batch/sequence tool would reduce latency:

```
batch_execute
  steps=[
    {tool: "page_navigate", args: {url: "..."}},
    {tool: "page_screenshot", args: {format: "file", path: "..."}},
    {tool: "js_evaluate", args: {script: "..."}},
    {tool: "page_screenshot", args: {format: "file", path: "..."}}
  ]
```

This is a nice-to-have — the current per-call latency (1-3s) is acceptable.
