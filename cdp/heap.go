package cdp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// HeapSnapshot represents a captured heap snapshot.
type HeapSnapshot struct {
	Path string // File path where snapshot is saved
	Size int64  // Size in bytes
}

// TakeHeapSnapshot captures a V8 heap snapshot and saves it to a file.
// If path is empty, a temporary file is created.
func (c *Client) TakeHeapSnapshot(ctx context.Context, path string) (*HeapSnapshot, error) {
	// Enable HeapProfiler
	if _, err := c.Send(ctx, HeapProfilerEnable, nil); err != nil {
		return nil, fmt.Errorf("cdp: failed to enable HeapProfiler: %w", err)
	}

	// Collect chunks
	var chunks []string
	var chunksMu sync.Mutex
	done := make(chan struct{})

	c.OnEvent(HeapProfilerAddHeapSnapshotChunk, func(params json.RawMessage) {
		var chunk HeapSnapshotChunk
		if err := json.Unmarshal(params, &chunk); err == nil {
			chunksMu.Lock()
			chunks = append(chunks, chunk.Chunk)
			chunksMu.Unlock()
		}
	})

	defer c.RemoveEventHandlers(HeapProfilerAddHeapSnapshotChunk)

	// Take snapshot
	snapshotDone := make(chan error, 1)
	go func() {
		_, err := c.Send(ctx, HeapProfilerTakeHeapSnapshot, map[string]interface{}{
			"reportProgress":            false,
			"treatGlobalObjectsAsRoots": true,
			"captureNumericValue":       true,
		})
		snapshotDone <- err
		close(done)
	}()

	// Wait for snapshot to complete
	select {
	case err := <-snapshotDone:
		if err != nil {
			return nil, fmt.Errorf("cdp: failed to take heap snapshot: %w", err)
		}
	case <-ctx.Done():
		return nil, ctx.Err()
	}

	// Small delay to ensure all chunks are received
	time.Sleep(100 * time.Millisecond)

	// Disable HeapProfiler
	_, _ = c.Send(ctx, HeapProfilerDisable, nil)

	// Combine chunks
	chunksMu.Lock()
	data := strings.Join(chunks, "")
	chunksMu.Unlock()

	if len(data) == 0 {
		return nil, fmt.Errorf("cdp: no heap snapshot data received")
	}

	// Determine output path
	if path == "" {
		tmpDir := os.TempDir()
		path = filepath.Join(tmpDir, fmt.Sprintf("heap-%d.heapsnapshot", time.Now().UnixNano()))
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, fmt.Errorf("cdp: failed to create directory: %w", err)
	}

	// Write to file
	if err := os.WriteFile(path, []byte(data), 0644); err != nil {
		return nil, fmt.Errorf("cdp: failed to write heap snapshot: %w", err)
	}

	info, _ := os.Stat(path)
	size := int64(0)
	if info != nil {
		size = info.Size()
	}

	return &HeapSnapshot{
		Path: path,
		Size: size,
	}, nil
}
