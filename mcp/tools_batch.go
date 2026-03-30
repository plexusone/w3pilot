package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	w3pilot "github.com/plexusone/w3pilot"
)

// BatchStep represents a single step in a batch execution.
type BatchStep struct {
	Tool string         `json:"tool" jsonschema:"Tool name to execute,required"`
	Args map[string]any `json:"args" jsonschema:"Arguments for the tool"`
}

// BatchExecuteInput defines the input for batch_execute tool.
type BatchExecuteInput struct {
	Steps           []BatchStep `json:"steps" jsonschema:"Array of tool steps to execute sequentially,required"`
	StopOnError     bool        `json:"stop_on_error" jsonschema:"Stop execution on first error (default: true)"`
	ContinueOnError bool        `json:"continue_on_error" jsonschema:"Continue execution even if a step fails"`
}

// BatchStepResult represents the result of a single batch step.
type BatchStepResult struct {
	Tool       string `json:"tool"`
	Success    bool   `json:"success"`
	Result     any    `json:"result,omitempty"`
	Error      string `json:"error,omitempty"`
	DurationMS int64  `json:"duration_ms"`
}

// BatchExecuteOutput defines the output for batch_execute tool.
type BatchExecuteOutput struct {
	Results       []BatchStepResult `json:"results"`
	TotalSteps    int               `json:"total_steps"`
	SuccessCount  int               `json:"success_count"`
	FailureCount  int               `json:"failure_count"`
	StoppedEarly  bool              `json:"stopped_early,omitempty"`
	TotalDuration int64             `json:"total_duration_ms"`
}

func (s *Server) handleBatchExecute(
	ctx context.Context,
	req *mcp.CallToolRequest,
	input BatchExecuteInput,
) (*mcp.CallToolResult, BatchExecuteOutput, error) {
	pilot, err := s.session.Pilot(ctx)
	if err != nil {
		return nil, BatchExecuteOutput{}, fmt.Errorf("browser not available: %w", err)
	}

	// Default to stop on error unless continue_on_error is set
	stopOnError := !input.ContinueOnError

	output := BatchExecuteOutput{
		Results:    make([]BatchStepResult, 0, len(input.Steps)),
		TotalSteps: len(input.Steps),
	}

	startTime := time.Now()

	for _, step := range input.Steps {
		stepStart := time.Now()
		result, err := s.executeStep(ctx, pilot, step)
		stepDuration := time.Since(stepStart).Milliseconds()

		stepResult := BatchStepResult{
			Tool:       step.Tool,
			DurationMS: stepDuration,
		}

		if err != nil {
			stepResult.Success = false
			stepResult.Error = err.Error()
			output.FailureCount++
		} else {
			stepResult.Success = true
			stepResult.Result = result
			output.SuccessCount++
		}

		output.Results = append(output.Results, stepResult)

		if err != nil && stopOnError {
			output.StoppedEarly = true
			break
		}
	}

	output.TotalDuration = time.Since(startTime).Milliseconds()

	return nil, output, nil
}

// executeStep dispatches a single batch step to the appropriate handler.
func (s *Server) executeStep(ctx context.Context, pilot *w3pilot.Pilot, step BatchStep) (any, error) {
	args := step.Args
	if args == nil {
		args = make(map[string]any)
	}

	switch step.Tool {
	// === Navigation ===
	case "page_navigate":
		url, _ := args["url"].(string)
		if url == "" {
			return nil, fmt.Errorf("url is required")
		}
		err := pilot.Go(ctx, url)
		if err != nil {
			return nil, err
		}
		return map[string]any{"url": url, "navigated": true}, nil

	case "page_go_back":
		err := pilot.Back(ctx)
		return map[string]any{"action": "back"}, err

	case "page_go_forward":
		err := pilot.Forward(ctx)
		return map[string]any{"action": "forward"}, err

	case "page_reload":
		err := pilot.Reload(ctx)
		return map[string]any{"action": "reload"}, err

	// === Page Info ===
	case "page_get_title":
		title, err := pilot.Title(ctx)
		if err != nil {
			return nil, err
		}
		return map[string]any{"title": title}, nil

	case "page_get_url":
		url, err := pilot.URL(ctx)
		if err != nil {
			return nil, err
		}
		return map[string]any{"url": url}, nil

	// === Screenshots ===
	case "page_screenshot":
		data, err := pilot.Screenshot(ctx)
		if err != nil {
			return nil, err
		}
		return map[string]any{"captured": true, "size": len(data)}, nil

	// === Element Interactions ===
	case "element_click":
		selector, _ := args["selector"].(string)
		if selector == "" {
			return nil, fmt.Errorf("selector is required")
		}
		elem, err := pilot.Find(ctx, selector, nil)
		if err != nil {
			return nil, err
		}
		err = elem.Click(ctx, nil)
		if err != nil {
			return nil, err
		}
		return map[string]any{"selector": selector, "clicked": true}, nil

	case "element_fill":
		selector, _ := args["selector"].(string)
		value, _ := args["value"].(string)
		if selector == "" {
			return nil, fmt.Errorf("selector is required")
		}
		elem, err := pilot.Find(ctx, selector, nil)
		if err != nil {
			return nil, err
		}
		err = elem.Fill(ctx, value, nil)
		if err != nil {
			return nil, err
		}
		return map[string]any{"selector": selector, "filled": true}, nil

	case "element_type":
		selector, _ := args["selector"].(string)
		text, _ := args["text"].(string)
		if selector == "" {
			return nil, fmt.Errorf("selector is required")
		}
		elem, err := pilot.Find(ctx, selector, nil)
		if err != nil {
			return nil, err
		}
		err = elem.Type(ctx, text, nil)
		if err != nil {
			return nil, err
		}
		return map[string]any{"selector": selector, "typed": true}, nil

	case "element_get_text":
		selector, _ := args["selector"].(string)
		if selector == "" {
			return nil, fmt.Errorf("selector is required")
		}
		elem, err := pilot.Find(ctx, selector, nil)
		if err != nil {
			return nil, err
		}
		text, err := elem.Text(ctx)
		if err != nil {
			return nil, err
		}
		return map[string]any{"selector": selector, "text": text}, nil

	// === JavaScript ===
	case "js_evaluate":
		script, _ := args["script"].(string)
		if script == "" {
			return nil, fmt.Errorf("script is required")
		}
		result, err := pilot.Evaluate(ctx, script)
		if err != nil {
			return nil, err
		}
		return map[string]any{"result": result}, nil

	// === Waiting ===
	case "wait_for_selector":
		selector, _ := args["selector"].(string)
		if selector == "" {
			return nil, fmt.Errorf("selector is required")
		}
		timeout := time.Duration(30) * time.Second
		if t, ok := args["timeout_ms"].(float64); ok && t > 0 {
			timeout = time.Duration(t) * time.Millisecond
		}
		// Use Find with timeout to wait for selector
		_, err := pilot.Find(ctx, selector, &w3pilot.FindOptions{Timeout: timeout})
		if err != nil {
			return nil, err
		}
		return map[string]any{"selector": selector, "found": true}, nil

	case "wait_for_load":
		timeout := time.Duration(30) * time.Second
		if t, ok := args["timeout_ms"].(float64); ok && t > 0 {
			timeout = time.Duration(t) * time.Millisecond
		}
		err := pilot.WaitForNavigation(ctx, timeout)
		if err != nil {
			return nil, err
		}
		return map[string]any{"loaded": true}, nil

	// === HTTP Request ===
	case "http_request":
		return s.executeHTTPRequest(ctx, pilot, args)

	default:
		return nil, fmt.Errorf("unsupported tool in batch: %s", step.Tool)
	}
}

// executeHTTPRequest handles the http_request tool within a batch.
func (s *Server) executeHTTPRequest(ctx context.Context, pilot *w3pilot.Pilot, args map[string]any) (any, error) {
	url, _ := args["url"].(string)
	if url == "" {
		return nil, fmt.Errorf("url is required")
	}

	method, _ := args["method"].(string)
	if method == "" {
		method = "GET"
	}

	body, _ := args["body"].(string)
	contentType, _ := args["content_type"].(string)

	maxLen := 8192
	if m, ok := args["max_body_length"].(float64); ok && m > 0 {
		maxLen = int(m)
	}

	// Build headers
	headers := make(map[string]string)
	if h, ok := args["headers"].(map[string]any); ok {
		for k, v := range h {
			if s, ok := v.(string); ok {
				headers[k] = s
			}
		}
	}
	if contentType != "" {
		headers["Content-Type"] = contentType
	}

	headersJSON, _ := json.Marshal(headers)

	script := fmt.Sprintf(`(async () => {
		const url = %q;
		const method = %q;
		const headers = %s;
		const body = %q;
		const maxLen = %d;

		const options = {
			method: method,
			credentials: 'include',
			headers: headers
		};

		if (body && (method === 'POST' || method === 'PUT' || method === 'PATCH')) {
			options.body = body;
		}

		const response = await fetch(url, options);

		const responseHeaders = {};
		response.headers.forEach((value, key) => {
			responseHeaders[key] = value;
		});

		let responseBody = await response.text();
		let truncated = false;
		if (maxLen > 0 && responseBody.length > maxLen) {
			responseBody = responseBody.substring(0, maxLen);
			truncated = true;
		}

		return {
			status: response.status,
			statusText: response.statusText,
			headers: responseHeaders,
			body: responseBody,
			truncated: truncated,
			url: response.url
		};
	})()`, url, method, string(headersJSON), body, maxLen)

	result, err := pilot.Evaluate(ctx, script)
	if err != nil {
		return nil, fmt.Errorf("fetch failed: %w", err)
	}

	return result, nil
}
