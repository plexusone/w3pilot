package w3pilot

import (
	"context"
	"encoding/json"
	"fmt"
)

// SelectorValidation represents the validation result for a single selector.
type SelectorValidation struct {
	Selector    string   `json:"selector"`
	Found       bool     `json:"found"`
	Count       int      `json:"count"`
	Visible     bool     `json:"visible"`
	Enabled     bool     `json:"enabled,omitempty"`
	TagName     string   `json:"tag_name,omitempty"`
	Suggestions []string `json:"suggestions,omitempty"`
}

// ValidateSelectors checks multiple selectors and returns validation results.
// This helps AI agents verify selectors before attempting interactions.
func (p *Pilot) ValidateSelectors(ctx context.Context, selectors []string) ([]SelectorValidation, error) {
	if p.closed {
		return nil, ErrConnectionClosed
	}

	if len(selectors) == 0 {
		return []SelectorValidation{}, nil
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return nil, err
	}

	// Build selector list as JSON for the script
	selectorsJSON, err := json.Marshal(selectors)
	if err != nil {
		return nil, fmt.Errorf("failed to encode selectors: %w", err)
	}

	script := fmt.Sprintf(`
		(function() {
			const selectors = %s;
			const results = [];

			function isVisible(el) {
				if (!el) return false;
				const rect = el.getBoundingClientRect();
				const style = window.getComputedStyle(el);
				return rect.width > 0 && rect.height > 0 &&
					   style.visibility !== 'hidden' &&
					   style.display !== 'none' &&
					   style.opacity !== '0';
			}

			function isEnabled(el) {
				if (!el) return false;
				return !el.disabled && !el.getAttribute('aria-disabled');
			}

			function findSuggestions(selector) {
				const suggestions = [];

				// Extract the base name from the selector
				let baseName = selector;
				if (baseName.startsWith('#') || baseName.startsWith('.')) {
					baseName = baseName.substring(1);
				}
				// Remove attribute selectors
				baseName = baseName.replace(/\[.*\]/g, '');
				// Remove pseudo-selectors
				baseName = baseName.replace(/:.*$/g, '');

				if (!baseName) return suggestions;

				// Try ID variations
				['#' + baseName, '#' + baseName + '-btn', '#' + baseName + 'Btn', '#' + baseName + '-button'].forEach(sel => {
					try { if (document.querySelector(sel)) suggestions.push(sel); } catch {}
				});

				// Try class variations
				['.' + baseName, '.' + baseName + '-btn', '.' + baseName + '-button', '.' + baseName.toLowerCase()].forEach(sel => {
					try { if (document.querySelector(sel)) suggestions.push(sel); } catch {}
				});

				// Try data-testid
				try {
					const testId = document.querySelector('[data-testid="' + baseName + '"]');
					if (testId) suggestions.push('[data-testid="' + baseName + '"]');
				} catch {}

				// Find elements with similar text content
				const lowerBase = baseName.toLowerCase();
				document.querySelectorAll('button, a, input[type="submit"], [role="button"]').forEach(el => {
					const text = (el.textContent || el.value || '').toLowerCase();
					if (text.includes(lowerBase) && suggestions.length < 5) {
						if (el.id) {
							const idSel = '#' + el.id;
							if (!suggestions.includes(idSel)) suggestions.push(idSel);
						} else if (el.className && typeof el.className === 'string') {
							const classSel = '.' + el.className.split(' ')[0];
							if (!suggestions.includes(classSel)) suggestions.push(classSel);
						}
					}
				});

				// Remove the original selector from suggestions
				return [...new Set(suggestions)].filter(s => s !== selector).slice(0, 5);
			}

			for (const selector of selectors) {
				let found = false;
				let count = 0;
				let visible = false;
				let enabled = false;
				let tagName = '';
				let suggestions = [];

				try {
					const elements = document.querySelectorAll(selector);
					count = elements.length;
					found = count > 0;

					if (found && elements[0]) {
						const el = elements[0];
						visible = isVisible(el);
						enabled = isEnabled(el);
						tagName = el.tagName.toLowerCase();
					} else {
						suggestions = findSuggestions(selector);
					}
				} catch (e) {
					// Invalid selector syntax
					found = false;
					suggestions = ['Invalid selector syntax: ' + e.message];
				}

				results.push({
					selector: selector,
					found: found,
					count: count,
					visible: visible,
					enabled: enabled,
					tag_name: tagName,
					suggestions: suggestions
				});
			}

			return JSON.stringify(results);
		})()
	`, string(selectorsJSON))

	// Execute via Evaluate
	rawResult, err := p.client.Send(ctx, "script.callFunction", map[string]interface{}{
		"functionDeclaration": "() => { " + script + " }",
		"target":              map[string]interface{}{"context": browsingCtx},
		"arguments":           []interface{}{},
		"awaitPromise":        true,
		"resultOwnership":     "root",
	})
	if err != nil {
		return nil, fmt.Errorf("validation script failed: %w", err)
	}

	// Parse the BiDi response
	var resp struct {
		Result struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"result"`
	}
	if err := json.Unmarshal(rawResult, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse validation response: %w", err)
	}

	// Parse the validation results
	var results []SelectorValidation
	if err := json.Unmarshal([]byte(resp.Result.Value), &results); err != nil {
		return nil, fmt.Errorf("failed to parse validation data: %w", err)
	}

	return results, nil
}

// ValidateSelector validates a single selector.
func (p *Pilot) ValidateSelector(ctx context.Context, selector string) (*SelectorValidation, error) {
	results, err := p.ValidateSelectors(ctx, []string{selector})
	if err != nil {
		return nil, err
	}
	if len(results) == 0 {
		return nil, fmt.Errorf("no validation result")
	}
	return &results[0], nil
}
