package w3pilot

import (
	"context"
	"encoding/json"
	"fmt"
)

// InspectOptions configures what elements to include in the inspection.
type InspectOptions struct {
	// IncludeButtons includes button elements (button, input[type=button], input[type=submit], [role=button])
	IncludeButtons bool
	// IncludeLinks includes anchor elements (a[href])
	IncludeLinks bool
	// IncludeInputs includes input and textarea elements
	IncludeInputs bool
	// IncludeSelects includes select elements
	IncludeSelects bool
	// IncludeHeadings includes heading elements (h1-h6)
	IncludeHeadings bool
	// IncludeImages includes images with alt text
	IncludeImages bool
	// MaxItems limits the number of items per category (0 = no limit, default 50)
	MaxItems int
}

// DefaultInspectOptions returns the default inspection options.
func DefaultInspectOptions() *InspectOptions {
	return &InspectOptions{
		IncludeButtons:  true,
		IncludeLinks:    true,
		IncludeInputs:   true,
		IncludeSelects:  true,
		IncludeHeadings: true,
		IncludeImages:   true,
		MaxItems:        50,
	}
}

// InspectResult contains the results of page inspection.
type InspectResult struct {
	URL      string           `json:"url"`
	Title    string           `json:"title"`
	Buttons  []InspectButton  `json:"buttons,omitempty"`
	Links    []InspectLink    `json:"links,omitempty"`
	Inputs   []InspectInput   `json:"inputs,omitempty"`
	Selects  []InspectSelect  `json:"selects,omitempty"`
	Headings []InspectHeading `json:"headings,omitempty"`
	Images   []InspectImage   `json:"images,omitempty"`
	Summary  InspectSummary   `json:"summary"`
}

// InspectButton represents a button element.
type InspectButton struct {
	Selector string `json:"selector"`
	Text     string `json:"text"`
	Type     string `json:"type,omitempty"`     // button, submit, reset
	Disabled bool   `json:"disabled,omitempty"` // true if disabled
	Visible  bool   `json:"visible"`
}

// InspectLink represents a link element.
type InspectLink struct {
	Selector string `json:"selector"`
	Text     string `json:"text"`
	Href     string `json:"href"`
	Visible  bool   `json:"visible"`
}

// InspectInput represents an input element.
type InspectInput struct {
	Selector    string `json:"selector"`
	Type        string `json:"type"`                  // text, password, email, etc.
	Name        string `json:"name,omitempty"`        // input name attribute
	Placeholder string `json:"placeholder,omitempty"` // placeholder text
	Value       string `json:"value,omitempty"`       // current value (masked for password)
	Label       string `json:"label,omitempty"`       // associated label text
	Required    bool   `json:"required,omitempty"`    // true if required
	Disabled    bool   `json:"disabled,omitempty"`    // true if disabled
	ReadOnly    bool   `json:"readonly,omitempty"`    // true if readonly
	Visible     bool   `json:"visible"`
}

// InspectSelect represents a select element.
type InspectSelect struct {
	Selector string   `json:"selector"`
	Name     string   `json:"name,omitempty"`
	Label    string   `json:"label,omitempty"`
	Options  []string `json:"options"`            // option texts (limited)
	Selected string   `json:"selected,omitempty"` // currently selected option text
	Multiple bool     `json:"multiple,omitempty"`
	Disabled bool     `json:"disabled,omitempty"`
	Visible  bool     `json:"visible"`
}

// InspectHeading represents a heading element.
type InspectHeading struct {
	Selector string `json:"selector"`
	Level    int    `json:"level"` // 1-6
	Text     string `json:"text"`
	Visible  bool   `json:"visible"`
}

// InspectImage represents an image element.
type InspectImage struct {
	Selector string `json:"selector"`
	Alt      string `json:"alt"`
	Src      string `json:"src"`
	Visible  bool   `json:"visible"`
}

// InspectSummary provides a summary of the inspection results.
type InspectSummary struct {
	TotalButtons  int `json:"total_buttons"`
	TotalLinks    int `json:"total_links"`
	TotalInputs   int `json:"total_inputs"`
	TotalSelects  int `json:"total_selects"`
	TotalHeadings int `json:"total_headings"`
	TotalImages   int `json:"total_images"`
}

// Inspect examines the current page and returns information about interactive elements.
// This is designed to help AI agents understand the page structure.
func (p *Pilot) Inspect(ctx context.Context, opts *InspectOptions) (*InspectResult, error) {
	if p.closed {
		return nil, ErrConnectionClosed
	}

	if opts == nil {
		opts = DefaultInspectOptions()
	}

	maxItems := opts.MaxItems
	if maxItems <= 0 {
		maxItems = 50
	}

	// Get URL and title
	url, _ := p.URL(ctx)
	title, _ := p.Title(ctx)

	result := &InspectResult{
		URL:   url,
		Title: title,
	}

	browsingCtx, err := p.getContext(ctx)
	if err != nil {
		return nil, err
	}

	// Build and execute the inspection script
	script := fmt.Sprintf(`
		(function() {
			const maxItems = %d;
			const result = {
				buttons: [],
				links: [],
				inputs: [],
				selects: [],
				headings: [],
				images: []
			};

			function isVisible(el) {
				const rect = el.getBoundingClientRect();
				const style = window.getComputedStyle(el);
				return rect.width > 0 && rect.height > 0 &&
					   style.visibility !== 'hidden' &&
					   style.display !== 'none' &&
					   style.opacity !== '0';
			}

			function getSelector(el) {
				if (el.id) return '#' + CSS.escape(el.id);
				if (el.name) return el.tagName.toLowerCase() + '[name="' + el.name + '"]';
				if (el.className && typeof el.className === 'string') {
					const classes = el.className.trim().split(/\s+/).filter(c => c).slice(0, 2);
					if (classes.length) return el.tagName.toLowerCase() + '.' + classes.join('.');
				}
				// Use nth-of-type for uniqueness
				const parent = el.parentElement;
				if (parent) {
					const siblings = Array.from(parent.children).filter(c => c.tagName === el.tagName);
					const index = siblings.indexOf(el) + 1;
					return el.tagName.toLowerCase() + ':nth-of-type(' + index + ')';
				}
				return el.tagName.toLowerCase();
			}

			function getLabel(el) {
				// Check for associated label
				if (el.id) {
					const label = document.querySelector('label[for="' + el.id + '"]');
					if (label) return label.textContent.trim();
				}
				// Check for wrapping label
				const parentLabel = el.closest('label');
				if (parentLabel) {
					return parentLabel.textContent.replace(el.value || '', '').trim();
				}
				// Check for aria-label
				if (el.getAttribute('aria-label')) {
					return el.getAttribute('aria-label');
				}
				return '';
			}

			// Buttons
			if (%t) {
				const buttons = document.querySelectorAll('button, input[type="button"], input[type="submit"], [role="button"]');
				for (const btn of buttons) {
					if (result.buttons.length >= maxItems) break;
					result.buttons.push({
						selector: getSelector(btn),
						text: (btn.textContent || btn.value || '').trim().substring(0, 100),
						type: btn.type || 'button',
						disabled: btn.disabled || false,
						visible: isVisible(btn)
					});
				}
			}

			// Links
			if (%t) {
				const links = document.querySelectorAll('a[href]');
				for (const link of links) {
					if (result.links.length >= maxItems) break;
					result.links.push({
						selector: getSelector(link),
						text: (link.textContent || '').trim().substring(0, 100),
						href: link.href,
						visible: isVisible(link)
					});
				}
			}

			// Inputs
			if (%t) {
				const inputs = document.querySelectorAll('input:not([type="button"]):not([type="submit"]):not([type="hidden"]), textarea');
				for (const input of inputs) {
					if (result.inputs.length >= maxItems) break;
					const type = input.type || 'text';
					result.inputs.push({
						selector: getSelector(input),
						type: type,
						name: input.name || '',
						placeholder: input.placeholder || '',
						value: type === 'password' ? (input.value ? '***' : '') : (input.value || '').substring(0, 50),
						label: getLabel(input),
						required: input.required || false,
						disabled: input.disabled || false,
						readonly: input.readOnly || false,
						visible: isVisible(input)
					});
				}
			}

			// Selects
			if (%t) {
				const selects = document.querySelectorAll('select');
				for (const select of selects) {
					if (result.selects.length >= maxItems) break;
					const options = Array.from(select.options).slice(0, 10).map(o => o.text.trim());
					const selectedOption = select.options[select.selectedIndex];
					result.selects.push({
						selector: getSelector(select),
						name: select.name || '',
						label: getLabel(select),
						options: options,
						selected: selectedOption ? selectedOption.text.trim() : '',
						multiple: select.multiple,
						disabled: select.disabled,
						visible: isVisible(select)
					});
				}
			}

			// Headings
			if (%t) {
				const headings = document.querySelectorAll('h1, h2, h3, h4, h5, h6');
				for (const h of headings) {
					if (result.headings.length >= maxItems) break;
					result.headings.push({
						selector: getSelector(h),
						level: parseInt(h.tagName.substring(1)),
						text: (h.textContent || '').trim().substring(0, 200),
						visible: isVisible(h)
					});
				}
			}

			// Images with alt text
			if (%t) {
				const images = document.querySelectorAll('img[alt]');
				for (const img of images) {
					if (result.images.length >= maxItems) break;
					if (!img.alt) continue;
					result.images.push({
						selector: getSelector(img),
						alt: img.alt.substring(0, 100),
						src: img.src,
						visible: isVisible(img)
					});
				}
			}

			return JSON.stringify(result);
		})()
	`, maxItems, opts.IncludeButtons, opts.IncludeLinks, opts.IncludeInputs,
		opts.IncludeSelects, opts.IncludeHeadings, opts.IncludeImages)

	// Execute via Evaluate
	rawResult, err := p.client.Send(ctx, "script.callFunction", map[string]interface{}{
		"functionDeclaration": "() => { " + script + " }",
		"target":              map[string]interface{}{"context": browsingCtx},
		"arguments":           []interface{}{},
		"awaitPromise":        true,
		"resultOwnership":     "root",
	})
	if err != nil {
		return nil, fmt.Errorf("inspection script failed: %w", err)
	}

	// Parse the BiDi response
	var resp struct {
		Result struct {
			Type  string `json:"type"`
			Value string `json:"value"`
		} `json:"result"`
	}
	if err := json.Unmarshal(rawResult, &resp); err != nil {
		return nil, fmt.Errorf("failed to parse inspection response: %w", err)
	}

	// Parse the inspection result JSON
	var inspectData struct {
		Buttons  []InspectButton  `json:"buttons"`
		Links    []InspectLink    `json:"links"`
		Inputs   []InspectInput   `json:"inputs"`
		Selects  []InspectSelect  `json:"selects"`
		Headings []InspectHeading `json:"headings"`
		Images   []InspectImage   `json:"images"`
	}
	if err := json.Unmarshal([]byte(resp.Result.Value), &inspectData); err != nil {
		return nil, fmt.Errorf("failed to parse inspection data: %w", err)
	}

	result.Buttons = inspectData.Buttons
	result.Links = inspectData.Links
	result.Inputs = inspectData.Inputs
	result.Selects = inspectData.Selects
	result.Headings = inspectData.Headings
	result.Images = inspectData.Images

	// Build summary
	result.Summary = InspectSummary{
		TotalButtons:  len(result.Buttons),
		TotalLinks:    len(result.Links),
		TotalInputs:   len(result.Inputs),
		TotalSelects:  len(result.Selects),
		TotalHeadings: len(result.Headings),
		TotalImages:   len(result.Images),
	}

	return result, nil
}
