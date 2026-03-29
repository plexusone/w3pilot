# Enhancement Requests

This directory contains enhancement requests for upstream dependencies (primarily VibiumDev/vibium clicker).

> **Note (2026-03-24)**: These requests are now **optional**. W3Pilot will implement dual-protocol
> support (BiDi + CDP) connecting directly to Chrome's CDP endpoint alongside clicker's BiDi connection.
> See [TASKS.md](../../TASKS.md) for the dual-protocol architecture.

## Status Legend

- **Proposed**: Initial draft, not yet submitted
- **Deferred**: Not needed due to dual-protocol approach
- **Submitted**: PR/Issue created in upstream repo
- **Accepted**: Upstream has agreed to implement
- **Implemented**: Available in upstream release
- **Integrated**: Integrated into W3Pilot

## Current Requests

| ID | Title | Priority | Status | Target |
|----|-------|----------|--------|--------|
| [EHRQ_heap-snapshot](EHRQ_heap-snapshot.md) | Heap Snapshot Support | P0 | Deferred | clicker |
| [EHRQ_network-response-body](EHRQ_network-response-body.md) | Network Response Body | P1 | Deferred | clicker |
| [EHRQ_emulation-presets](EHRQ_emulation-presets.md) | Network/CPU Emulation | P1 | Deferred | clicker |

## Why Deferred?

**Discovery (2026-03-24)**: Chrome exposes both BiDi and CDP on the same browser instance:

```
┌─────────────────────────────────────────────────────────────┐
│                      Chrome Browser                          │
│                   (launched by clicker)                      │
├─────────────────────────────────────────────────────────────┤
│  BiDi WebSocket ◄──── clicker (existing)                    │
│  CDP WebSocket  ◄──── W3Pilot CDPClient (new)              │
└─────────────────────────────────────────────────────────────┘
```

W3Pilot can connect directly to Chrome's CDP port (discovered via `DevToolsActivePort` file)
without requiring clicker modifications.

## When to Revisit

These requests may be worth submitting if:

1. **Performance**: Direct CDP has overhead vs clicker passthrough
2. **Coordination**: Two connections cause race conditions
3. **Simplicity**: Single protocol preferred over dual

## Upstream Repository

**VibiumDev/vibium**: https://github.com/VibiumDev/vibium

Clicker source: `/clicker/`

## Original Implementation Notes

All enhancement requests require CDP (Chrome DevTools Protocol) access. Clicker currently uses WebDriver BiDi which doesn't expose these CDP domains:

- `HeapProfiler` - For heap snapshots
- `Network.getResponseBody` - For response content
- `Network.emulateNetworkConditions` - For throttling
- `Emulation.setCPUThrottlingRate` - For CPU throttling

### Original Approach (Now Optional)

1. Add CDP passthrough in clicker's BiDi proxy
2. Expose as `vibium:*` commands following existing patterns
3. W3Pilot integrates via existing transport

### Current Approach (Dual Protocol)

1. W3Pilot establishes direct CDP connection to Chrome
2. Manage two connections (BiDi via clicker + CDP direct)
3. No upstream dependency for CDP features
