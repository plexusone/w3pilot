package cdp

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
)

// Screencast domain methods.
const (
	PageStartScreencast     = "Page.startScreencast"
	PageStopScreencast      = "Page.stopScreencast"
	PageScreencastFrameAck  = "Page.screencastFrameAck"
	PageScreencastFrame     = "Page.screencastFrame"
)

// ScreencastFormat represents the image format for screencast.
type ScreencastFormat string

const (
	ScreencastFormatJPEG ScreencastFormat = "jpeg"
	ScreencastFormatPNG  ScreencastFormat = "png"
)

// ScreencastOptions configures screencast capture.
type ScreencastOptions struct {
	Format      ScreencastFormat `json:"format,omitempty"`      // Image format (jpeg or png)
	Quality     int              `json:"quality,omitempty"`     // Image quality (0-100, jpeg only)
	MaxWidth    int              `json:"maxWidth,omitempty"`    // Maximum width
	MaxHeight   int              `json:"maxHeight,omitempty"`   // Maximum height
	EveryNthFrame int            `json:"everyNthFrame,omitempty"` // Skip frames (1 = every frame)
}

// ScreencastFrame represents a captured frame.
type ScreencastFrame struct {
	Data       string            `json:"data"`       // Base64-encoded image data
	Metadata   ScreencastMetadata `json:"metadata"`
	SessionID  int               `json:"sessionId"`
}

// ScreencastMetadata contains frame metadata.
type ScreencastMetadata struct {
	OffsetTop       float64 `json:"offsetTop"`
	PageScaleFactor float64 `json:"pageScaleFactor"`
	DeviceWidth     int     `json:"deviceWidth"`
	DeviceHeight    int     `json:"deviceHeight"`
	ScrollOffsetX   float64 `json:"scrollOffsetX"`
	ScrollOffsetY   float64 `json:"scrollOffsetY"`
	Timestamp       float64 `json:"timestamp,omitempty"`
}

// ScreencastFrameHandler is called for each captured frame.
type ScreencastFrameHandler func(frame *ScreencastFrame)

// Screencast manages screencast capture.
type Screencast struct {
	client    *Client
	mu        sync.RWMutex
	running   bool
	handler   ScreencastFrameHandler
	sessionID int
}

// NewScreencast creates a new screencast manager.
func NewScreencast(client *Client) *Screencast {
	return &Screencast{
		client: client,
	}
}

// Start begins screencast capture.
func (s *Screencast) Start(ctx context.Context, opts *ScreencastOptions, handler ScreencastFrameHandler) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.running {
		return fmt.Errorf("cdp: screencast already running")
	}

	if opts == nil {
		opts = &ScreencastOptions{
			Format:  ScreencastFormatJPEG,
			Quality: 80,
		}
	}

	s.handler = handler

	// Register event handler for frames
	s.client.OnEvent(PageScreencastFrame, func(params json.RawMessage) {
		var frame ScreencastFrame
		if err := json.Unmarshal(params, &frame); err != nil {
			return
		}

		s.mu.RLock()
		h := s.handler
		s.mu.RUnlock()

		if h != nil {
			h(&frame)
		}

		// Acknowledge frame received
		_, _ = s.client.Send(ctx, PageScreencastFrameAck, map[string]interface{}{
			"sessionId": frame.SessionID,
		})
	})

	// Start screencast
	_, err := s.client.Send(ctx, PageStartScreencast, opts)
	if err != nil {
		return fmt.Errorf("cdp: failed to start screencast: %w", err)
	}

	s.running = true
	return nil
}

// Stop ends screencast capture.
func (s *Screencast) Stop(ctx context.Context) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.running {
		return nil
	}

	_, err := s.client.Send(ctx, PageStopScreencast, nil)
	if err != nil {
		return fmt.Errorf("cdp: failed to stop screencast: %w", err)
	}

	s.running = false
	s.handler = nil
	return nil
}

// IsRunning returns whether screencast is active.
func (s *Screencast) IsRunning() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.running
}
