# CDP Features (Chrome DevTools Protocol)

W3Pilot provides direct access to Chrome DevTools Protocol (CDP) for advanced profiling and emulation features not available through WebDriver BiDi.

## Dual-Protocol Architecture

W3Pilot connects to a single Chrome browser using **both** protocols:

```
┌──────────────────────────────────────────────────────────┐
│                     w3pilot                              │
│  ┌─────────────────┐    ┌─────────────────┐              │
│  │   BiDi Client   │    │   CDP Client    │              │
│  │  (automation)   │    │  (profiling)    │              │
│  └────────┬────────┘    └────────┬────────┘              │
│           │                      │                        │
│           ▼                      ▼                        │
│    VibiumDev Clicker      Chrome DevTools                 │
│    (WebDriver BiDi)       (CDP WebSocket)                 │
│           │                      │                        │
│           └──────────┬───────────┘                        │
│                      ▼                                    │
│                   Chrome                                  │
│           (single browser instance)                       │
└──────────────────────────────────────────────────────────┘
```

| Protocol | Use Case |
|----------|----------|
| **BiDi** | Page automation, element interactions, screenshots, tracing |
| **CDP** | Heap profiling, network response bodies, CPU/network emulation |

## Checking CDP Availability

CDP is automatically connected when launching a browser:

```go
pilot, err := w3pilot.Launch(ctx)
if err != nil {
    log.Fatal(err)
}

if pilot.HasCDP() {
    fmt.Printf("CDP connected on port %d\n", pilot.CDPPort())
} else {
    fmt.Println("CDP not available")
}
```

## Heap Snapshots

Capture V8 heap snapshots for memory profiling. These can be loaded in Chrome DevTools Memory tab.

```go
// Take heap snapshot (auto-generates path if empty)
snapshot, err := pilot.TakeHeapSnapshot(ctx, "")
if err != nil {
    log.Fatal(err)
}
fmt.Printf("Snapshot saved: %s (%d bytes)\n", snapshot.Path, snapshot.Size)

// Specify a path
snapshot, err := pilot.TakeHeapSnapshot(ctx, "/tmp/memory-leak.heapsnapshot")
```

### Analyzing Heap Snapshots

1. Open Chrome DevTools
2. Go to **Memory** tab
3. Click **Load** and select the `.heapsnapshot` file
4. Analyze memory usage, retained size, and object allocations

## Network Emulation

Simulate various network conditions for testing performance under degraded networks.

### Using Presets

```go
import "github.com/plexusone/w3pilot/cdp"

// Slow 3G: 400ms latency, 400 Kbps download, 400 Kbps upload
err := pilot.EmulateNetwork(ctx, cdp.NetworkSlow3G)

// Fast 3G: 150ms latency, 1.5 Mbps download, 750 Kbps upload
err := pilot.EmulateNetwork(ctx, cdp.NetworkFast3G)

// 4G: 50ms latency, 10 Mbps download, 5 Mbps upload
err := pilot.EmulateNetwork(ctx, cdp.Network4G)
```

### Custom Conditions

```go
err := pilot.EmulateNetwork(ctx, cdp.NetworkConditions{
    Offline:            false,
    Latency:            200,          // 200ms
    DownloadThroughput: 1024 * 1024,  // 1 MB/s
    UploadThroughput:   512 * 1024,   // 512 KB/s
})
```

### Offline Mode

```go
err := pilot.EmulateNetwork(ctx, cdp.NetworkConditions{
    Offline: true,
})
```

### Clearing Network Emulation

```go
err := pilot.ClearNetworkEmulation(ctx)
```

## CPU Emulation

Throttle CPU to simulate lower-powered devices for performance testing.

### Using Presets

```go
import "github.com/plexusone/w3pilot/cdp"

// 2x slowdown
err := pilot.EmulateCPU(ctx, cdp.CPU2xSlowdown)

// 4x slowdown (mid-tier mobile)
err := pilot.EmulateCPU(ctx, cdp.CPU4xSlowdown)

// 6x slowdown (low-end mobile)
err := pilot.EmulateCPU(ctx, cdp.CPU6xSlowdown)
```

### Custom Rate

```go
// 10x slowdown
err := pilot.EmulateCPU(ctx, 10)
```

### Clearing CPU Emulation

```go
err := pilot.ClearCPUEmulation(ctx)
```

## Direct CDP Access

For advanced use cases, access the CDP client directly to send any CDP command.

```go
if pilot.HasCDP() {
    cdpClient := pilot.CDP()

    // Get performance metrics
    result, err := cdpClient.Send(ctx, "Performance.getMetrics", nil)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("Metrics: %s\n", result)

    // Enable a domain
    _, err = cdpClient.Send(ctx, "DOM.enable", nil)

    // Send command with parameters
    result, err = cdpClient.Send(ctx, "DOM.getDocument", map[string]interface{}{
        "depth": 1,
    })
}
```

### CDP Event Handling

```go
cdpClient := pilot.CDP()

// Subscribe to events
cdpClient.OnEvent("Network.requestWillBeSent", func(params json.RawMessage) {
    fmt.Printf("Request: %s\n", params)
})

// Enable the domain to start receiving events
cdpClient.Send(ctx, "Network.enable", nil)
```

## Performance Testing Workflow

A typical performance testing workflow combining CDP features:

```go
func runPerformanceTest(ctx context.Context, pilot *w3pilot.Pilot, url string) error {
    // 1. Set network conditions
    if err := pilot.EmulateNetwork(ctx, cdp.NetworkSlow3G); err != nil {
        return err
    }

    // 2. Set CPU throttling
    if err := pilot.EmulateCPU(ctx, cdp.CPU4xSlowdown); err != nil {
        return err
    }

    // 3. Navigate and measure
    start := time.Now()
    if err := pilot.Go(ctx, url); err != nil {
        return err
    }
    loadTime := time.Since(start)
    fmt.Printf("Load time: %v\n", loadTime)

    // 4. Take heap snapshot
    snapshot, err := pilot.TakeHeapSnapshot(ctx, "")
    if err != nil {
        return err
    }
    fmt.Printf("Memory snapshot: %s (%d bytes)\n", snapshot.Path, snapshot.Size)

    // 5. Clear emulation
    pilot.ClearNetworkEmulation(ctx)
    pilot.ClearCPUEmulation(ctx)

    return nil
}
```

## Troubleshooting

### CDP Not Available

If `HasCDP()` returns false:

1. Ensure Chrome was launched with remote debugging enabled
2. Check that the `DevToolsActivePort` file exists in Chrome's user data directory
3. Verify no other process is using the CDP port

### Connection Refused

If CDP connection fails:

```go
pilot, err := w3pilot.Launch(ctx)
if err != nil {
    // Check error for CDP-specific issues
    if strings.Contains(err.Error(), "cdp") {
        log.Printf("CDP connection issue: %v", err)
        // Continue with BiDi-only mode if acceptable
    }
}
```

### Heap Snapshot Timeout

Large heap snapshots can take time. Increase the context timeout:

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
defer cancel()

snapshot, err := pilot.TakeHeapSnapshot(ctx, "")
```
