package w3pilot

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

// LoginOptions configures the login workflow.
type LoginOptions struct {
	// UsernameSelector is the CSS selector for the username/email field
	UsernameSelector string
	// PasswordSelector is the CSS selector for the password field
	PasswordSelector string
	// SubmitSelector is the CSS selector for the submit button
	SubmitSelector string
	// Username is the username/email to fill
	Username string
	// Password is the password to fill
	Password string
	// SuccessIndicator is a CSS selector or URL pattern that indicates successful login
	SuccessIndicator string
	// Timeout for the entire login process
	Timeout time.Duration
}

// LoginResult contains the result of a login attempt.
type LoginResult struct {
	Success     bool   `json:"success"`
	URL         string `json:"url"`
	Title       string `json:"title"`
	Message     string `json:"message"`
	ErrorReason string `json:"error_reason,omitempty"`
}

// Login performs an automated login workflow.
// It fills the username and password fields, submits the form, and waits for success.
func (p *Pilot) Login(ctx context.Context, opts *LoginOptions) (*LoginResult, error) {
	if p.closed {
		return nil, ErrConnectionClosed
	}

	if opts == nil {
		return nil, fmt.Errorf("login options required")
	}

	timeout := opts.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	result := &LoginResult{}

	// Find and fill username
	usernameEl, err := p.Find(ctx, opts.UsernameSelector, nil)
	if err != nil {
		result.ErrorReason = fmt.Sprintf("username field not found: %s", opts.UsernameSelector)
		return result, nil
	}

	if err := usernameEl.Fill(ctx, opts.Username, nil); err != nil {
		result.ErrorReason = fmt.Sprintf("failed to fill username: %v", err)
		return result, nil
	}

	// Find and fill password
	passwordEl, err := p.Find(ctx, opts.PasswordSelector, nil)
	if err != nil {
		result.ErrorReason = fmt.Sprintf("password field not found: %s", opts.PasswordSelector)
		return result, nil
	}

	if err := passwordEl.Fill(ctx, opts.Password, nil); err != nil {
		result.ErrorReason = fmt.Sprintf("failed to fill password: %v", err)
		return result, nil
	}

	// Find and click submit
	submitEl, err := p.Find(ctx, opts.SubmitSelector, nil)
	if err != nil {
		result.ErrorReason = fmt.Sprintf("submit button not found: %s", opts.SubmitSelector)
		return result, nil
	}

	if err := submitEl.Click(ctx, nil); err != nil {
		result.ErrorReason = fmt.Sprintf("failed to click submit: %v", err)
		return result, nil
	}

	// Wait for navigation and success indicator
	if opts.SuccessIndicator != "" {
		// Check if it's a URL pattern (contains / or **)
		if len(opts.SuccessIndicator) > 0 && (opts.SuccessIndicator[0] == '/' ||
			opts.SuccessIndicator[0] == '*' ||
			opts.SuccessIndicator[0] == 'h') {
			// URL pattern
			if err := p.WaitForURL(ctx, opts.SuccessIndicator, timeout); err != nil {
				result.ErrorReason = fmt.Sprintf("success URL not reached: %v", err)
				return result, nil
			}
		} else {
			// CSS selector
			_, err := p.Find(ctx, opts.SuccessIndicator, &FindOptions{Timeout: timeout})
			if err != nil {
				result.ErrorReason = fmt.Sprintf("success indicator not found: %s", opts.SuccessIndicator)
				return result, nil
			}
		}
	} else {
		// Default: wait for navigation to complete
		if err := p.WaitForNavigation(ctx, timeout); err != nil {
			result.ErrorReason = "navigation did not complete"
			return result, nil
		}
	}

	// Get final page state
	result.Success = true
	result.URL, _ = p.URL(ctx)
	result.Title, _ = p.Title(ctx)
	result.Message = "Login successful"

	return result, nil
}

// ExtractTableOptions configures the table extraction.
type ExtractTableOptions struct {
	// IncludeHeaders determines if the first row is treated as headers
	IncludeHeaders bool
	// MaxRows limits the number of rows to extract (0 = no limit)
	MaxRows int
	// HeaderSelector is a custom selector for header cells (default: th, thead td)
	HeaderSelector string
	// RowSelector is a custom selector for data rows (default: tbody tr, tr)
	RowSelector string
	// CellSelector is a custom selector for cells (default: td)
	CellSelector string
}

// TableResult contains the extracted table data.
type TableResult struct {
	Headers  []string            `json:"headers,omitempty"`
	Rows     [][]string          `json:"rows"`
	RowsJSON []map[string]string `json:"rows_json,omitempty"` // Rows as objects with header keys
	RowCount int                 `json:"row_count"`
}

// ExtractTable extracts data from an HTML table into structured JSON.
func (p *Pilot) ExtractTable(ctx context.Context, selector string, opts *ExtractTableOptions) (*TableResult, error) {
	if p.closed {
		return nil, ErrConnectionClosed
	}

	if opts == nil {
		opts = &ExtractTableOptions{
			IncludeHeaders: true,
		}
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return nil, err
	}

	maxRows := opts.MaxRows
	if maxRows <= 0 {
		maxRows = 1000 // Reasonable default limit
	}

	headerSelector := opts.HeaderSelector
	if headerSelector == "" {
		headerSelector = "th"
	}

	rowSelector := opts.RowSelector
	if rowSelector == "" {
		rowSelector = "tbody tr, tr"
	}

	cellSelector := opts.CellSelector
	if cellSelector == "" {
		cellSelector = "td"
	}

	script := fmt.Sprintf(`
		(function() {
			const table = document.querySelector(%q);
			if (!table) {
				return JSON.stringify({error: 'Table not found'});
			}

			const includeHeaders = %t;
			const maxRows = %d;
			const headerSelector = %q;
			const rowSelector = %q;
			const cellSelector = %q;

			const result = {
				headers: [],
				rows: [],
				rows_json: []
			};

			// Extract headers
			const headerCells = table.querySelectorAll(headerSelector);
			if (headerCells.length > 0) {
				headerCells.forEach(cell => {
					result.headers.push(cell.textContent.trim());
				});
			}

			// If no th elements, try the first row
			if (result.headers.length === 0 && includeHeaders) {
				const firstRow = table.querySelector('tr');
				if (firstRow) {
					const cells = firstRow.querySelectorAll('td, th');
					cells.forEach(cell => {
						result.headers.push(cell.textContent.trim());
					});
				}
			}

			// Extract data rows
			const rows = table.querySelectorAll(rowSelector);
			let rowCount = 0;

			rows.forEach((row, index) => {
				if (rowCount >= maxRows) return;

				const cells = row.querySelectorAll(cellSelector);
				if (cells.length === 0) return;

				// Skip header row if it's the same as first data row
				if (index === 0 && result.headers.length > 0) {
					const firstCellText = cells[0]?.textContent.trim();
					if (firstCellText === result.headers[0]) {
						return;
					}
				}

				const rowData = [];
				const rowObj = {};

				cells.forEach((cell, cellIndex) => {
					const text = cell.textContent.trim();
					rowData.push(text);
					if (result.headers[cellIndex]) {
						rowObj[result.headers[cellIndex]] = text;
					}
				});

				result.rows.push(rowData);
				if (Object.keys(rowObj).length > 0) {
					result.rows_json.push(rowObj);
				}
				rowCount++;
			});

			result.row_count = result.rows.length;
			return JSON.stringify(result);
		})()
	`, selector, opts.IncludeHeaders, maxRows, headerSelector, rowSelector, cellSelector)

	// Execute via Evaluate
	rawResult, err := p.client.Send(ctx, "script.callFunction", map[string]interface{}{
		"functionDeclaration": "() => { " + script + " }",
		"target":              map[string]interface{}{"context": browsingCtx},
		"arguments":           []interface{}{},
		"awaitPromise":        true,
		"resultOwnership":     "root",
	})
	if err != nil {
		return nil, fmt.Errorf("table extraction script failed: %w", err)
	}

	// Parse the BiDi response
	var resp struct {
		Result struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"result"`
	}
	if err := json.Unmarshal(rawResult, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse table extraction response: %w", err)
	}

	// Parse the table result
	var tableData struct {
		Error    string              `json:"error,omitempty"`
		Headers  []string            `json:"headers"`
		Rows     [][]string          `json:"rows"`
		RowsJSON []map[string]string `json:"rows_json"`
		RowCount int                 `json:"row_count"`
	}
	if err := json.Unmarshal([]byte(resp.Result.Value), &tableData); err != nil {
		return nil, fmt.Errorf("failed to parse table data: %w", err)
	}

	if tableData.Error != "" {
		return nil, fmt.Errorf("table extraction: %s", tableData.Error)
	}

	return &TableResult{
		Headers:  tableData.Headers,
		Rows:     tableData.Rows,
		RowsJSON: tableData.RowsJSON,
		RowCount: tableData.RowCount,
	}, nil
}
