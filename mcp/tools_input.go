package mcp

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	vibium "github.com/plexusone/w3pilot"
)

// KeyboardPress tool

type KeyboardPressInput struct {
	Key string `json:"key" jsonschema:"Key to press (e.g. Enter Tab ArrowDown),required"`
}

type KeyboardPressOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleKeyboardPress(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input KeyboardPressInput,
) (*mcp.CallToolResult, KeyboardPressOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, KeyboardPressOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	keyboard, err := pilot.Keyboard(ctx)
	if err != nil {
		return nil, KeyboardPressOutput{}, fmt.Errorf("keyboard not available: %w", err)
	}

	err = keyboard.Press(ctx, input.Key)
	if err != nil {
		return nil, KeyboardPressOutput{}, fmt.Errorf("keyboard press failed: %w", err)
	}

	return nil, KeyboardPressOutput{Message: fmt.Sprintf("Pressed key: %s", input.Key)}, nil
}

// KeyboardDown tool

type KeyboardDownInput struct {
	Key string `json:"key" jsonschema:"Key to hold down,required"`
}

type KeyboardDownOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleKeyboardDown(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input KeyboardDownInput,
) (*mcp.CallToolResult, KeyboardDownOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, KeyboardDownOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	keyboard, err := pilot.Keyboard(ctx)
	if err != nil {
		return nil, KeyboardDownOutput{}, fmt.Errorf("keyboard not available: %w", err)
	}

	err = keyboard.Down(ctx, input.Key)
	if err != nil {
		return nil, KeyboardDownOutput{}, fmt.Errorf("keyboard down failed: %w", err)
	}

	return nil, KeyboardDownOutput{Message: fmt.Sprintf("Holding key: %s", input.Key)}, nil
}

// KeyboardUp tool

type KeyboardUpInput struct {
	Key string `json:"key" jsonschema:"Key to release,required"`
}

type KeyboardUpOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleKeyboardUp(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input KeyboardUpInput,
) (*mcp.CallToolResult, KeyboardUpOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, KeyboardUpOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	keyboard, err := pilot.Keyboard(ctx)
	if err != nil {
		return nil, KeyboardUpOutput{}, fmt.Errorf("keyboard not available: %w", err)
	}

	err = keyboard.Up(ctx, input.Key)
	if err != nil {
		return nil, KeyboardUpOutput{}, fmt.Errorf("keyboard up failed: %w", err)
	}

	return nil, KeyboardUpOutput{Message: fmt.Sprintf("Released key: %s", input.Key)}, nil
}

// KeyboardType tool

type KeyboardTypeInput struct {
	Text string `json:"text" jsonschema:"Text to type,required"`
}

type KeyboardTypeOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleKeyboardType(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input KeyboardTypeInput,
) (*mcp.CallToolResult, KeyboardTypeOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, KeyboardTypeOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	keyboard, err := pilot.Keyboard(ctx)
	if err != nil {
		return nil, KeyboardTypeOutput{}, fmt.Errorf("keyboard not available: %w", err)
	}

	err = keyboard.Type(ctx, input.Text)
	if err != nil {
		return nil, KeyboardTypeOutput{}, fmt.Errorf("keyboard type failed: %w", err)
	}

	return nil, KeyboardTypeOutput{Message: fmt.Sprintf("Typed: %s", truncateString(input.Text, 50))}, nil
}

// MouseClick tool

type MouseClickInput struct {
	X          float64 `json:"x" jsonschema:"X coordinate,required"`
	Y          float64 `json:"y" jsonschema:"Y coordinate,required"`
	Button     string  `json:"button" jsonschema:"Mouse button: left right middle"`
	ClickCount int     `json:"click_count" jsonschema:"Number of clicks (default: 1)"`
}

type MouseClickOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleMouseClick(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input MouseClickInput,
) (*mcp.CallToolResult, MouseClickOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, MouseClickOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	mouse, err := pilot.Mouse(ctx)
	if err != nil {
		return nil, MouseClickOutput{}, fmt.Errorf("mouse not available: %w", err)
	}

	opts := &vibium.ClickOptions{}
	if input.Button != "" {
		opts.Button = vibium.MouseButton(input.Button)
	}
	if input.ClickCount > 0 {
		opts.ClickCount = input.ClickCount
	}

	err = mouse.Click(ctx, input.X, input.Y, opts)
	if err != nil {
		return nil, MouseClickOutput{}, fmt.Errorf("mouse click failed: %w", err)
	}

	return nil, MouseClickOutput{Message: fmt.Sprintf("Clicked at (%f, %f)", input.X, input.Y)}, nil
}

// MouseMove tool

type MouseMoveInput struct {
	X float64 `json:"x" jsonschema:"X coordinate,required"`
	Y float64 `json:"y" jsonschema:"Y coordinate,required"`
}

type MouseMoveOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleMouseMove(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input MouseMoveInput,
) (*mcp.CallToolResult, MouseMoveOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, MouseMoveOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	mouse, err := pilot.Mouse(ctx)
	if err != nil {
		return nil, MouseMoveOutput{}, fmt.Errorf("mouse not available: %w", err)
	}

	err = mouse.Move(ctx, input.X, input.Y)
	if err != nil {
		return nil, MouseMoveOutput{}, fmt.Errorf("mouse move failed: %w", err)
	}

	return nil, MouseMoveOutput{Message: fmt.Sprintf("Moved mouse to (%f, %f)", input.X, input.Y)}, nil
}

// MouseDown tool

type MouseDownInput struct {
	Button string `json:"button" jsonschema:"Mouse button: left right middle"`
}

type MouseDownOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleMouseDown(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input MouseDownInput,
) (*mcp.CallToolResult, MouseDownOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, MouseDownOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	mouse, err := pilot.Mouse(ctx)
	if err != nil {
		return nil, MouseDownOutput{}, fmt.Errorf("mouse not available: %w", err)
	}

	button := vibium.MouseButton(input.Button)
	err = mouse.Down(ctx, button)
	if err != nil {
		return nil, MouseDownOutput{}, fmt.Errorf("mouse down failed: %w", err)
	}

	return nil, MouseDownOutput{Message: "Mouse button pressed"}, nil
}

// MouseUp tool

type MouseUpInput struct {
	Button string `json:"button" jsonschema:"Mouse button: left right middle"`
}

type MouseUpOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleMouseUp(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input MouseUpInput,
) (*mcp.CallToolResult, MouseUpOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, MouseUpOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	mouse, err := pilot.Mouse(ctx)
	if err != nil {
		return nil, MouseUpOutput{}, fmt.Errorf("mouse not available: %w", err)
	}

	button := vibium.MouseButton(input.Button)
	err = mouse.Up(ctx, button)
	if err != nil {
		return nil, MouseUpOutput{}, fmt.Errorf("mouse up failed: %w", err)
	}

	return nil, MouseUpOutput{Message: "Mouse button released"}, nil
}

// MouseWheel tool

type MouseWheelInput struct {
	DeltaX float64 `json:"delta_x" jsonschema:"Horizontal scroll amount"`
	DeltaY float64 `json:"delta_y" jsonschema:"Vertical scroll amount"`
}

type MouseWheelOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleMouseWheel(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input MouseWheelInput,
) (*mcp.CallToolResult, MouseWheelOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, MouseWheelOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	mouse, err := pilot.Mouse(ctx)
	if err != nil {
		return nil, MouseWheelOutput{}, fmt.Errorf("mouse not available: %w", err)
	}

	err = mouse.Wheel(ctx, input.DeltaX, input.DeltaY)
	if err != nil {
		return nil, MouseWheelOutput{}, fmt.Errorf("mouse wheel failed: %w", err)
	}

	return nil, MouseWheelOutput{Message: fmt.Sprintf("Scrolled (%f, %f)", input.DeltaX, input.DeltaY)}, nil
}

// TouchTap tool

type TouchTapInput struct {
	X float64 `json:"x" jsonschema:"X coordinate,required"`
	Y float64 `json:"y" jsonschema:"Y coordinate,required"`
}

type TouchTapOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleTouchTap(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input TouchTapInput,
) (*mcp.CallToolResult, TouchTapOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, TouchTapOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	touch, err := pilot.Touch(ctx)
	if err != nil {
		return nil, TouchTapOutput{}, fmt.Errorf("touch not available: %w", err)
	}

	err = touch.Tap(ctx, input.X, input.Y)
	if err != nil {
		return nil, TouchTapOutput{}, fmt.Errorf("touch tap failed: %w", err)
	}

	return nil, TouchTapOutput{Message: fmt.Sprintf("Tapped at (%f, %f)", input.X, input.Y)}, nil
}

// TouchSwipe tool

type TouchSwipeInput struct {
	StartX float64 `json:"start_x" jsonschema:"Starting X coordinate,required"`
	StartY float64 `json:"start_y" jsonschema:"Starting Y coordinate,required"`
	EndX   float64 `json:"end_x" jsonschema:"Ending X coordinate,required"`
	EndY   float64 `json:"end_y" jsonschema:"Ending Y coordinate,required"`
}

type TouchSwipeOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleTouchSwipe(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input TouchSwipeInput,
) (*mcp.CallToolResult, TouchSwipeOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, TouchSwipeOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	touch, err := pilot.Touch(ctx)
	if err != nil {
		return nil, TouchSwipeOutput{}, fmt.Errorf("touch not available: %w", err)
	}

	err = touch.Swipe(ctx, input.StartX, input.StartY, input.EndX, input.EndY)
	if err != nil {
		return nil, TouchSwipeOutput{}, fmt.Errorf("touch swipe failed: %w", err)
	}

	return nil, TouchSwipeOutput{Message: fmt.Sprintf("Swiped from (%f, %f) to (%f, %f)", input.StartX, input.StartY, input.EndX, input.EndY)}, nil
}

// MouseDrag tool

type MouseDragInput struct {
	StartX float64 `json:"start_x" jsonschema:"Starting X coordinate,required"`
	StartY float64 `json:"start_y" jsonschema:"Starting Y coordinate,required"`
	EndX   float64 `json:"end_x" jsonschema:"Ending X coordinate,required"`
	EndY   float64 `json:"end_y" jsonschema:"Ending Y coordinate,required"`
	Steps  int     `json:"steps,omitempty" jsonschema:"Number of intermediate steps (default: 10)"`
}

type MouseDragOutput struct {
	Message string `json:"message"`
}

func (s *Server) handleMouseDrag(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input MouseDragInput,
) (*mcp.CallToolResult, MouseDragOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, MouseDragOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	mouse, err := pilot.Mouse(ctx)
	if err != nil {
		return nil, MouseDragOutput{}, fmt.Errorf("mouse not available: %w", err)
	}

	steps := input.Steps
	if steps == 0 {
		steps = 10
	}

	// Move to start position
	err = mouse.Move(ctx, input.StartX, input.StartY)
	if err != nil {
		return nil, MouseDragOutput{}, fmt.Errorf("mouse move failed: %w", err)
	}

	// Press mouse button
	err = mouse.Down(ctx, "left")
	if err != nil {
		return nil, MouseDragOutput{}, fmt.Errorf("mouse down failed: %w", err)
	}

	// Move to end position in steps
	deltaX := (input.EndX - input.StartX) / float64(steps)
	deltaY := (input.EndY - input.StartY) / float64(steps)

	for i := 1; i <= steps; i++ {
		x := input.StartX + deltaX*float64(i)
		y := input.StartY + deltaY*float64(i)
		err = mouse.Move(ctx, x, y)
		if err != nil {
			// Release button on error
			_ = mouse.Up(ctx, "left")
			return nil, MouseDragOutput{}, fmt.Errorf("mouse move failed: %w", err)
		}
	}

	// Release mouse button
	err = mouse.Up(ctx, "left")
	if err != nil {
		return nil, MouseDragOutput{}, fmt.Errorf("mouse up failed: %w", err)
	}

	return nil, MouseDragOutput{
		Message: fmt.Sprintf("Dragged from (%.0f, %.0f) to (%.0f, %.0f)", input.StartX, input.StartY, input.EndX, input.EndY),
	}, nil
}
