package vibium

import (
	"context"
	"encoding/json"
	"fmt"
)

// VideoOptions configures video recording.
type VideoOptions struct {
	// Dir is the directory to save videos to. Defaults to a temp directory.
	Dir string `json:"dir,omitempty"`
	// Size specifies the video dimensions. Defaults to viewport size.
	Size *VideoSize `json:"size,omitempty"`
}

// VideoSize specifies video dimensions.
type VideoSize struct {
	Width  int `json:"width"`
	Height int `json:"height"`
}

// Video represents an ongoing or completed video recording.
type Video struct {
	client  *BiDiClient
	context string

	// FilePath is the file path where the video will be saved.
	FilePath string `json:"path,omitempty"`
}

// StartVideo starts recording video of the page.
// The video is saved when StopVideo is called or the browser closes.
func (v *Vibe) StartVideo(ctx context.Context, opts *VideoOptions) (*Video, error) {
	if v.closed {
		return nil, ErrConnectionClosed
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return nil, err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	if opts != nil {
		if opts.Dir != "" {
			params["dir"] = opts.Dir
		}
		if opts.Size != nil {
			params["size"] = map[string]int{
				"width":  opts.Size.Width,
				"height": opts.Size.Height,
			}
		}
	}

	result, err := v.client.Send(ctx, "vibium:video.start", params)
	if err != nil {
		return nil, fmt.Errorf("failed to start video: %w", err)
	}

	var resp struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse video response: %w", err)
	}

	return &Video{
		client:   v.client,
		context:  browsingCtx,
		FilePath: resp.Path,
	}, nil
}

// StopVideo stops video recording and returns the video path.
func (v *Vibe) StopVideo(ctx context.Context) (string, error) {
	if v.closed {
		return "", ErrConnectionClosed
	}

	browsingCtx, err := v.getContext(ctx)
	if err != nil {
		return "", err
	}

	params := map[string]interface{}{
		"context": browsingCtx,
	}

	result, err := v.client.Send(ctx, "vibium:video.stop", params)
	if err != nil {
		return "", fmt.Errorf("failed to stop video: %w", err)
	}

	var resp struct {
		Path string `json:"path"`
	}
	if err := json.Unmarshal(result, &resp); err != nil {
		return "", fmt.Errorf("failed to parse video response: %w", err)
	}

	return resp.Path, nil
}

// Path returns the path where the video will be saved.
// The file may not exist until StopVideo is called.
func (vid *Video) Path() string {
	return vid.FilePath
}

// Delete deletes the video file.
func (vid *Video) Delete(ctx context.Context) error {
	params := map[string]interface{}{
		"path": vid.FilePath,
	}

	_, err := vid.client.Send(ctx, "vibium:video.delete", params)
	return err
}
