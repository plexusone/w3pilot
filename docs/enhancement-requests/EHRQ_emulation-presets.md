# Enhancement Request: Network and CPU Emulation Presets

**ID**: EHRQ_emulation-presets
**Status**: Proposed
**Priority**: P1
**Target**: VibiumDev/vibium (clicker)
**Date**: 2026-03-24

## Summary

Add `vibium:emulate.network` and `vibium:emulate.cpu` commands for performance testing under constrained conditions.

## Motivation

Chrome DevTools MCP provides `emulate` tool which enables:

- Testing on slow networks (3G, 4G)
- Simulating mobile CPU performance
- Performance regression detection
- Real-world condition testing

Essential for accessibility testing where users may have older devices or slow connections.

## Current State

**Available**:
```go
// Media emulation only
pilot.EmulateMedia(ctx, webpilot.EmulateMediaOptions{
    ColorScheme: "dark",
    ReducedMotion: "reduce",
})

// Offline mode
pilot.SetOffline(ctx, true)
```

**Not available**: Network throttling, CPU throttling.

## Proposed Commands

### 1. Network Emulation

#### Request (Preset)

```json
{
  "id": 1,
  "method": "vibium:emulate.network",
  "params": {
    "preset": "slow-3g"
  }
}
```

#### Request (Custom)

```json
{
  "id": 1,
  "method": "vibium:emulate.network",
  "params": {
    "offline": false,
    "latency": 150,
    "downloadThroughput": 1500000,
    "uploadThroughput": 750000
  }
}
```

#### Presets

| Preset | Latency (ms) | Download (bps) | Upload (bps) |
|--------|-------------|----------------|--------------|
| `offline` | 0 | 0 | 0 |
| `slow-3g` | 400 | 400000 | 400000 |
| `fast-3g` | 150 | 1500000 | 750000 |
| `4g` | 50 | 4000000 | 3000000 |
| `wifi` | 10 | 30000000 | 15000000 |

### 2. CPU Emulation

#### Request

```json
{
  "id": 2,
  "method": "vibium:emulate.cpu",
  "params": {
    "rate": 4
  }
}
```

#### Rate Values

| Rate | Description | Use Case |
|------|-------------|----------|
| 1 | No throttling | Baseline |
| 2 | 2x slowdown | Fast mobile |
| 4 | 4x slowdown | Mid-tier mobile |
| 6 | 6x slowdown | Low-end mobile |

## CDP Implementation

### Network Throttling

```javascript
await cdp.send('Network.emulateNetworkConditions', {
  offline: false,
  latency: 150,           // ms
  downloadThroughput: 1500000,  // bytes/sec
  uploadThroughput: 750000,     // bytes/sec
  connectionType: 'cellular3g'
});
```

### CPU Throttling

```javascript
await cdp.send('Emulation.setCPUThrottlingRate', {
  rate: 4  // 4x slowdown
});
```

## Use Cases

1. **Mobile testing**: Simulate mobile network and CPU
2. **Performance budgets**: Verify load times under constraints
3. **Accessibility**: Test for users with slow connections
4. **CI/CD**: Automated performance regression testing

## WebPilot Integration

Once available in clicker:

```go
// SDK - Network
err := pilot.EmulateNetwork(ctx, webpilot.NetworkPresetSlow3G)
// or
err := pilot.EmulateNetwork(ctx, webpilot.NetworkConditions{
    Latency: 150,
    DownloadThroughput: 1500000,
})

// SDK - CPU
err := pilot.EmulateCPU(ctx, 4) // 4x slowdown

// MCP Tools
// tool: emulate_network
// params: { preset?: string, latency?: number, download?: number, upload?: number }

// tool: emulate_cpu
// params: { rate: number }
```

## References

- [CDP Network.emulateNetworkConditions](https://chromedevtools.github.io/devtools-protocol/tot/Network/#method-emulateNetworkConditions)
- [CDP Emulation.setCPUThrottlingRate](https://chromedevtools.github.io/devtools-protocol/tot/Emulation/#method-setCPUThrottlingRate)
- [Chrome DevTools MCP - emulate](https://github.com/anthropics/chrome-devtools-mcp)
