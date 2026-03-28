# Enhancement Request: Network Response Body Access

**ID**: EHRQ_network-response-body
**Status**: Proposed
**Priority**: P1
**Target**: VibiumDev/vibium (clicker)
**Date**: 2026-03-24

## Summary

Add `vibium:network.getResponseBody` command to retrieve full request/response content for debugging.

## Motivation

Chrome DevTools MCP provides `get_network_request` which enables:

- API response debugging
- Test fixture generation
- Response validation
- Binary content inspection (images, files)

Current WebPilot network tools only capture metadata, not response bodies.

## Current State

**Available**:
```go
requests, _ := pilot.NetworkRequests(ctx, nil)
// Returns: URL, Method, Status, Headers, ResourceType
```

**Not available**: Response body content.

## Proposed Command

### Request

```json
{
  "id": 1,
  "method": "vibium:network.getResponseBody",
  "params": {
    "requestId": "1234.5",
    "saveTo": "/tmp/response.json"
  }
}
```

### Response

```json
{
  "id": 1,
  "type": "success",
  "result": {
    "body": "{\"data\": ...}",
    "base64Encoded": false,
    "size": 1234,
    "path": "/tmp/response.json"
  }
}
```

For binary content:
```json
{
  "id": 1,
  "type": "success",
  "result": {
    "body": "iVBORw0KGgo...",
    "base64Encoded": true,
    "size": 45678,
    "mimeType": "image/png"
  }
}
```

## CDP Implementation

Uses `Network` domain:

```javascript
// Enable Network with response body capture
await cdp.send('Network.enable', {
  maxResourceBufferSize: 10000000,
  maxTotalBufferSize: 50000000
});

// Get response body
const { body, base64Encoded } = await cdp.send('Network.getResponseBody', {
  requestId: requestId
});
```

## Prerequisites

Network interception must be enabled to capture response bodies. This may require:

1. Enabling `Network.enable` with buffer settings
2. Optionally using `Fetch.enable` for request interception

## Use Cases

1. **API debugging**: Inspect actual response payloads
2. **Test fixtures**: Save responses for mocking
3. **Validation**: Verify response content matches expected
4. **Binary inspection**: Check downloaded images/files

## WebPilot Integration

Once available in clicker:

```go
// SDK
body, err := pilot.GetNetworkResponseBody(ctx, requestId)
fmt.Println(body.Content)      // Response body
fmt.Println(body.Base64Encoded) // true for binary

// MCP Tool
// tool: get_network_response_body
// params: { request_id: string, save_to?: string }
// returns: { body: string, base64_encoded: bool, size: number }
```

## References

- [CDP Network.getResponseBody](https://chromedevtools.github.io/devtools-protocol/tot/Network/#method-getResponseBody)
- [Chrome DevTools MCP - get_network_request](https://github.com/anthropics/chrome-devtools-mcp)
