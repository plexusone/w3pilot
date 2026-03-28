# Enhancement Request: Heap Snapshot Support

**ID**: EHRQ_heap-snapshot
**Status**: Proposed
**Priority**: P0
**Target**: VibiumDev/vibium (clicker)
**Date**: 2026-03-24

## Summary

Add `vibium:heap.takeSnapshot` command to capture V8 heap snapshots for memory profiling.

## Motivation

Chrome DevTools MCP provides `take_memory_snapshot` which enables:

- Memory leak detection
- Heap analysis in Chrome DevTools
- Memory regression testing
- Production debugging

WebPilot needs this for feature parity and to support accessibility testing tools that may have memory leaks.

## Current State

**Available via JavaScript**:
```javascript
performance.memory.usedJSHeapSize     // Basic heap size
performance.memory.totalJSHeapSize    // Total allocated
performance.memory.jsHeapSizeLimit    // Max heap size
```

**Not available**: Full `.heapsnapshot` file for DevTools analysis.

## Proposed Command

### Request

```json
{
  "id": 1,
  "method": "vibium:heap.takeSnapshot",
  "params": {
    "reportProgress": false
  }
}
```

### Response

```json
{
  "id": 1,
  "type": "success",
  "result": {
    "path": "/tmp/heap-1711234567.heapsnapshot",
    "size": 1234567
  }
}
```

## CDP Implementation

Uses `HeapProfiler` domain:

```javascript
// Enable HeapProfiler
await cdp.send('HeapProfiler.enable');

// Take snapshot
const chunks = [];
cdp.on('HeapProfiler.addHeapSnapshotChunk', (params) => {
  chunks.push(params.chunk);
});

await cdp.send('HeapProfiler.takeHeapSnapshot', {
  reportProgress: false
});

// Write to file
const snapshot = chunks.join('');
fs.writeFileSync(path, snapshot);
```

## Use Cases

1. **Memory leak detection**: Compare snapshots before/after operations
2. **Accessibility testing**: Ensure a11y tools don't cause memory issues
3. **Performance testing**: Track memory growth over time
4. **CI/CD integration**: Fail builds on memory regression

## WebPilot Integration

Once available in clicker:

```go
// SDK
snapshot, err := pilot.TakeHeapSnapshot(ctx)
fmt.Println(snapshot.Path)  // "/tmp/heap-xxx.heapsnapshot"
fmt.Println(snapshot.Size)  // 1234567

// MCP Tool
// tool: take_heap_snapshot
// returns: { path: string, size: number }
```

## References

- [CDP HeapProfiler](https://chromedevtools.github.io/devtools-protocol/tot/HeapProfiler/)
- [Chrome DevTools MCP - take_memory_snapshot](https://github.com/anthropics/chrome-devtools-mcp)
