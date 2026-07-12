package w3pilot

import (
	"fmt"
	"regexp"
	"strings"
	"sync"
)

// RefStore stores element references from page mapping.
// It is thread-safe for concurrent access.
type RefStore struct {
	mu   sync.RWMutex
	refs map[string]ElementRef // @e1 -> ElementRef
}

// NewRefStore creates a new ref store.
func NewRefStore() *RefStore {
	return &RefStore{
		refs: make(map[string]ElementRef),
	}
}

// Store stores a list of element refs, replacing any existing refs.
func (s *RefStore) Store(refs []ElementRef) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Clear existing refs
	s.refs = make(map[string]ElementRef)

	// Store new refs
	for _, ref := range refs {
		s.refs[ref.Ref] = ref
	}
}

// Get retrieves an element ref by its reference string (e.g., "@e1").
func (s *RefStore) Get(ref string) (ElementRef, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	r, ok := s.refs[ref]
	return r, ok
}

// GetSelector retrieves the CSS selector for a ref.
func (s *RefStore) GetSelector(ref string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if r, ok := s.refs[ref]; ok {
		return r.Selector, true
	}
	return "", false
}

// Count returns the number of stored refs.
func (s *RefStore) Count() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.refs)
}

// All returns all stored refs.
func (s *RefStore) All() []ElementRef {
	s.mu.RLock()
	defer s.mu.RUnlock()

	result := make([]ElementRef, 0, len(s.refs))
	for _, ref := range s.refs {
		result = append(result, ref)
	}
	return result
}

// Clear removes all stored refs.
func (s *RefStore) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.refs = make(map[string]ElementRef)
}

// refPattern matches @e1, @e2, etc.
var refPattern = regexp.MustCompile(`^@e\d+$`)

// IsRef checks if a string is a valid element reference.
func IsRef(s string) bool {
	return refPattern.MatchString(s)
}

// ResolveRef resolves an element reference to its CSS selector.
// If the input is not a ref (doesn't match @e\d+), it is returned as-is.
// This allows seamless use of both refs and regular selectors.
func (s *RefStore) ResolveRef(selectorOrRef string) (string, error) {
	selectorOrRef = strings.TrimSpace(selectorOrRef)

	if !IsRef(selectorOrRef) {
		// Not a ref, return as-is (it's a regular selector)
		return selectorOrRef, nil
	}

	selector, ok := s.GetSelector(selectorOrRef)
	if !ok {
		return "", fmt.Errorf("unknown element reference: %s (run 'map' to refresh element references)", selectorOrRef)
	}

	return selector, nil
}

// MustResolveRef resolves a ref, panicking on error.
// Use only when you're certain the ref exists.
func (s *RefStore) MustResolveRef(selectorOrRef string) string {
	selector, err := s.ResolveRef(selectorOrRef)
	if err != nil {
		panic(err)
	}
	return selector
}
