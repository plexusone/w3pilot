package w3pilot

import (
	"context"
	"fmt"
	"strings"
)

// ElementRef represents a mapped element with a human-friendly reference.
type ElementRef struct {
	// Ref is the human-friendly reference (e.g., "@e1", "@e2")
	Ref string `json:"ref"`

	// Tag is the HTML tag name (e.g., "button", "input", "a")
	Tag string `json:"tag"`

	// Role is the ARIA role (e.g., "button", "textbox", "link")
	Role string `json:"role,omitempty"`

	// Text is the visible text content
	Text string `json:"text,omitempty"`

	// Label is the accessible label (aria-label or associated label)
	Label string `json:"label,omitempty"`

	// Placeholder is the placeholder text for inputs
	Placeholder string `json:"placeholder,omitempty"`

	// Type is the input type for input elements
	Type string `json:"type,omitempty"`

	// Selector is the CSS selector that can locate this element
	Selector string `json:"selector"`

	// Visible indicates if the element is currently visible
	Visible bool `json:"visible"`

	// Enabled indicates if the element is enabled (not disabled)
	Enabled bool `json:"enabled"`
}

// MapOptions configures element mapping behavior.
type MapOptions struct {
	// IncludeHidden includes hidden elements in the mapping
	IncludeHidden bool

	// MaxElements limits the number of elements to map (0 = no limit)
	MaxElements int

	// Scope limits mapping to descendants of this selector
	Scope string
}

// MapElements scans the page for interactive elements and returns refs.
// Interactive elements include buttons, links, inputs, selects, textareas,
// and elements with click handlers or interactive ARIA roles.
func (p *Pilot) MapElements(ctx context.Context, opts *MapOptions) ([]ElementRef, error) {
	if opts == nil {
		opts = &MapOptions{}
	}

	script := fmt.Sprintf(`
(function() {
	const includeHidden = %t;
	const maxElements = %d;
	const scope = %q;

	// Get the scope element
	const container = scope ? document.querySelector(scope) : document;
	if (!container) return [];

	// Selectors for interactive elements
	const interactiveSelectors = [
		'button',
		'a[href]',
		'input:not([type="hidden"])',
		'select',
		'textarea',
		'[role="button"]',
		'[role="link"]',
		'[role="textbox"]',
		'[role="checkbox"]',
		'[role="radio"]',
		'[role="combobox"]',
		'[role="listbox"]',
		'[role="option"]',
		'[role="menuitem"]',
		'[role="tab"]',
		'[role="slider"]',
		'[role="spinbutton"]',
		'[role="switch"]',
		'[tabindex]:not([tabindex="-1"])',
		'[onclick]',
		'[contenteditable="true"]',
	];

	// Find all interactive elements
	const elements = [];
	for (const sel of interactiveSelectors) {
		try {
			const found = container.querySelectorAll(sel);
			for (const el of found) {
				if (!elements.includes(el)) {
					elements.push(el);
				}
			}
		} catch (e) {}
	}

	// Filter and map elements
	const refs = [];
	let refNum = 1;

	for (const el of elements) {
		// Check visibility
		const style = window.getComputedStyle(el);
		const isVisible = style.display !== 'none' &&
			style.visibility !== 'hidden' &&
			el.offsetParent !== null;

		if (!includeHidden && !isVisible) continue;

		// Check if enabled
		const isEnabled = !el.disabled && !el.hasAttribute('aria-disabled');

		// Get element info
		const tag = el.tagName.toLowerCase();
		const role = el.getAttribute('role') || getImplicitRole(el);
		const text = getElementText(el);
		const label = getElementLabel(el);
		const placeholder = el.placeholder || '';
		const inputType = el.type || '';

		// Generate a unique selector
		const selector = generateSelector(el);

		refs.push({
			ref: '@e' + refNum,
			tag: tag,
			role: role,
			text: text,
			label: label,
			placeholder: placeholder,
			type: inputType,
			selector: selector,
			visible: isVisible,
			enabled: isEnabled
		});

		refNum++;
		if (maxElements > 0 && refs.length >= maxElements) break;
	}

	return refs;

	function getImplicitRole(el) {
		const tag = el.tagName.toLowerCase();
		const type = el.type || '';

		const roleMap = {
			'button': 'button',
			'a': el.hasAttribute('href') ? 'link' : null,
			'input': getInputRole(type),
			'select': 'combobox',
			'textarea': 'textbox',
			'img': 'img',
			'nav': 'navigation',
			'header': 'banner',
			'footer': 'contentinfo',
			'main': 'main',
			'aside': 'complementary',
			'form': 'form',
			'article': 'article',
			'section': 'region'
		};

		return roleMap[tag] || null;
	}

	function getInputRole(type) {
		const inputRoles = {
			'button': 'button',
			'submit': 'button',
			'reset': 'button',
			'checkbox': 'checkbox',
			'radio': 'radio',
			'range': 'slider',
			'number': 'spinbutton',
			'search': 'searchbox',
			'email': 'textbox',
			'tel': 'textbox',
			'url': 'textbox',
			'text': 'textbox',
			'password': 'textbox'
		};
		return inputRoles[type] || 'textbox';
	}

	function getElementText(el) {
		// For inputs, get value
		if (el.tagName === 'INPUT' || el.tagName === 'TEXTAREA') {
			return el.value || '';
		}
		// For select, get selected option text
		if (el.tagName === 'SELECT') {
			return el.options[el.selectedIndex]?.text || '';
		}
		// For other elements, get trimmed text content
		return (el.textContent || '').trim().slice(0, 100);
	}

	function getElementLabel(el) {
		// Check aria-label first
		if (el.hasAttribute('aria-label')) {
			return el.getAttribute('aria-label');
		}
		// Check aria-labelledby
		const labelledBy = el.getAttribute('aria-labelledby');
		if (labelledBy) {
			const labels = labelledBy.split(' ').map(id => {
				const labelEl = document.getElementById(id);
				return labelEl ? labelEl.textContent.trim() : '';
			}).filter(Boolean);
			if (labels.length) return labels.join(' ');
		}
		// Check associated label
		if (el.id) {
			const label = document.querySelector('label[for="' + el.id + '"]');
			if (label) return label.textContent.trim();
		}
		// Check parent label
		const parentLabel = el.closest('label');
		if (parentLabel) {
			// Get label text excluding the input itself
			const clone = parentLabel.cloneNode(true);
			const inputs = clone.querySelectorAll('input, select, textarea, button');
			inputs.forEach(i => i.remove());
			return clone.textContent.trim();
		}
		// Check title attribute
		if (el.title) {
			return el.title;
		}
		return '';
	}

	function generateSelector(el) {
		// Try ID first
		if (el.id) {
			return '#' + CSS.escape(el.id);
		}
		// Try data-testid
		const testId = el.getAttribute('data-testid');
		if (testId) {
			return '[data-testid="' + testId + '"]';
		}
		// Try name attribute for form elements
		if (el.name && ['INPUT', 'SELECT', 'TEXTAREA', 'BUTTON'].includes(el.tagName)) {
			return el.tagName.toLowerCase() + '[name="' + el.name + '"]';
		}
		// Try unique class
		if (el.className) {
			const classes = el.className.split(' ').filter(c => c.trim());
			for (const cls of classes) {
				const selector = el.tagName.toLowerCase() + '.' + CSS.escape(cls);
				if (document.querySelectorAll(selector).length === 1) {
					return selector;
				}
			}
		}
		// Build a path-based selector
		const path = [];
		let current = el;
		while (current && current !== document.body) {
			let selector = current.tagName.toLowerCase();
			if (current.id) {
				selector = '#' + CSS.escape(current.id);
				path.unshift(selector);
				break;
			}
			const parent = current.parentElement;
			if (parent) {
				const siblings = Array.from(parent.children).filter(
					c => c.tagName === current.tagName
				);
				if (siblings.length > 1) {
					const index = siblings.indexOf(current) + 1;
					selector += ':nth-of-type(' + index + ')';
				}
			}
			path.unshift(selector);
			current = parent;
		}
		return path.join(' > ');
	}
})()
`, opts.IncludeHidden, opts.MaxElements, opts.Scope)

	result, err := p.Evaluate(ctx, script)
	if err != nil {
		return nil, fmt.Errorf("failed to map elements: %w", err)
	}

	// Parse result
	refs := []ElementRef{}
	if arr, ok := result.([]any); ok {
		for _, item := range arr {
			if obj, ok := item.(map[string]any); ok {
				ref := ElementRef{
					Ref:         getString(obj, "ref"),
					Tag:         getString(obj, "tag"),
					Role:        getString(obj, "role"),
					Text:        getString(obj, "text"),
					Label:       getString(obj, "label"),
					Placeholder: getString(obj, "placeholder"),
					Type:        getString(obj, "type"),
					Selector:    getString(obj, "selector"),
					Visible:     getBool(obj, "visible"),
					Enabled:     getBool(obj, "enabled"),
				}
				refs = append(refs, ref)
			}
		}
	}

	return refs, nil
}

// getString safely extracts a string from a map.
func getString(m map[string]any, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

// getBool safely extracts a bool from a map.
func getBool(m map[string]any, key string) bool {
	if v, ok := m[key]; ok {
		if b, ok := v.(bool); ok {
			return b
		}
	}
	return false
}

// RefDiff represents the difference between two element mappings.
type RefDiff struct {
	// Added contains refs that exist in After but not in Before
	Added []ElementRef `json:"added"`

	// Removed contains refs that exist in Before but not in After
	Removed []ElementRef `json:"removed"`

	// Changed contains refs where the selector changed (element moved/recreated)
	Changed []RefChange `json:"changed"`

	// Unchanged contains refs that are identical in both mappings
	Unchanged []ElementRef `json:"unchanged"`

	// Summary provides counts for quick overview
	Summary RefDiffSummary `json:"summary"`
}

// RefChange represents a change in an element's selector.
type RefChange struct {
	Before ElementRef `json:"before"`
	After  ElementRef `json:"after"`
}

// RefDiffSummary provides counts for the diff.
type RefDiffSummary struct {
	Added     int `json:"added"`
	Removed   int `json:"removed"`
	Changed   int `json:"changed"`
	Unchanged int `json:"unchanged"`
}

// DiffRefs compares two sets of element refs and returns the differences.
// Elements are matched by their display characteristics (tag, text, label, placeholder)
// rather than by ref number, since ref numbers may change between mappings.
func DiffRefs(before, after []ElementRef) *RefDiff {
	diff := &RefDiff{
		Added:     []ElementRef{},
		Removed:   []ElementRef{},
		Changed:   []RefChange{},
		Unchanged: []ElementRef{},
	}

	// Create fingerprint maps for matching
	beforeMap := make(map[string]ElementRef)
	afterMap := make(map[string]ElementRef)

	for _, ref := range before {
		fp := refFingerprint(ref)
		beforeMap[fp] = ref
	}

	for _, ref := range after {
		fp := refFingerprint(ref)
		afterMap[fp] = ref
	}

	// Find unchanged and changed elements
	for fp, beforeRef := range beforeMap {
		if afterRef, exists := afterMap[fp]; exists {
			if beforeRef.Selector == afterRef.Selector {
				diff.Unchanged = append(diff.Unchanged, afterRef)
			} else {
				diff.Changed = append(diff.Changed, RefChange{
					Before: beforeRef,
					After:  afterRef,
				})
			}
			delete(afterMap, fp) // Mark as processed
		} else {
			diff.Removed = append(diff.Removed, beforeRef)
		}
	}

	// Remaining items in afterMap are new
	for _, ref := range afterMap {
		diff.Added = append(diff.Added, ref)
	}

	// Update summary
	diff.Summary = RefDiffSummary{
		Added:     len(diff.Added),
		Removed:   len(diff.Removed),
		Changed:   len(diff.Changed),
		Unchanged: len(diff.Unchanged),
	}

	return diff
}

// refFingerprint creates a fingerprint for matching elements across mappings.
// Uses tag + text/label/placeholder to identify the "same" element.
func refFingerprint(ref ElementRef) string {
	// Primary identifier is tag + main text content
	identifier := ref.Tag

	if ref.Text != "" {
		identifier += ":" + ref.Text
	} else if ref.Label != "" {
		identifier += ":label=" + ref.Label
	} else if ref.Placeholder != "" {
		identifier += ":placeholder=" + ref.Placeholder
	}

	// Add type for inputs to distinguish different input types
	if ref.Type != "" {
		identifier += ":type=" + ref.Type
	}

	return identifier
}

// HasChanges returns true if there are any differences.
func (d *RefDiff) HasChanges() bool {
	return d.Summary.Added > 0 || d.Summary.Removed > 0 || d.Summary.Changed > 0
}

// FormatRef formats an element ref for display.
func (r ElementRef) FormatRef() string {
	var parts []string

	// Start with ref and tag
	parts = append(parts, fmt.Sprintf("%s [%s]", r.Ref, r.Tag))

	// Add role if different from tag
	if r.Role != "" && r.Role != r.Tag {
		parts[0] += fmt.Sprintf(" role=%s", r.Role)
	}

	// Add text/label/placeholder
	if r.Text != "" {
		text := r.Text
		if len(text) > 50 {
			text = text[:50] + "..."
		}
		parts = append(parts, fmt.Sprintf(`"%s"`, text))
	} else if r.Label != "" {
		label := r.Label
		if len(label) > 50 {
			label = label[:50] + "..."
		}
		parts = append(parts, fmt.Sprintf(`label: "%s"`, label))
	} else if r.Placeholder != "" {
		parts = append(parts, fmt.Sprintf(`(placeholder: "%s")`, r.Placeholder))
	}

	// Add type for inputs
	if r.Type != "" && r.Tag == "input" {
		parts = append(parts, fmt.Sprintf("type=%s", r.Type))
	}

	// Add status indicators
	if !r.Enabled {
		parts = append(parts, "[disabled]")
	}
	if !r.Visible {
		parts = append(parts, "[hidden]")
	}

	return strings.Join(parts, " ")
}
