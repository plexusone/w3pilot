package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/spf13/cobra"

	w3pilot "github.com/plexusone/w3pilot"
)

var (
	mapTimeout       time.Duration
	mapIncludeHidden bool
	mapMaxElements   int
	mapScope         string
)

var mapCmd = &cobra.Command{
	Use:   "map",
	Short: "Map interactive elements to refs",
	Long: `Map interactive elements on the page to human-friendly references (@e1, @e2, etc.).

These refs can be used in place of CSS selectors in subsequent commands.
This is especially useful for AI agents to identify clickable elements.

Examples:
  # Map all interactive elements
  w3pilot map

  # Map with scope (only elements within #main-content)
  w3pilot map --scope "#main-content"

  # Include hidden elements
  w3pilot map --include-hidden

  # Use refs in commands
  w3pilot click @e1
  w3pilot fill @e3 "john@example.com"`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), mapTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		opts := &w3pilot.MapOptions{
			IncludeHidden: mapIncludeHidden,
			MaxElements:   mapMaxElements,
			Scope:         mapScope,
		}

		refs, err := pilot.MapElements(ctx, opts)
		if err != nil {
			return fmt.Errorf("mapping failed: %w", err)
		}

		// Save refs to disk
		if err := saveRefs(refs); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not save refs: %v\n", err)
		}

		Output(refs, func(data interface{}) string {
			r := data.([]w3pilot.ElementRef)
			return formatRefs(r)
		})
		return nil
	},
}

func formatRefs(refs []w3pilot.ElementRef) string {
	if len(refs) == 0 {
		return "No interactive elements found."
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Found %d interactive elements:\n\n", len(refs)))

	for _, ref := range refs {
		sb.WriteString(ref.FormatRef())
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("\nUse refs with commands: w3pilot click @e1, w3pilot fill @e2 \"value\"\n"))

	return sb.String()
}

// getRefsPath returns the path to the refs file
func getRefsPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".w3pilot-refs.json"
	}
	return filepath.Join(home, ".w3pilot", "refs.json")
}

// saveRefs saves refs to disk
func saveRefs(refs []w3pilot.ElementRef) error {
	path := getRefsPath()

	// Ensure directory exists
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create refs directory: %w", err)
	}

	data, err := json.MarshalIndent(refs, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal refs: %w", err)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write refs file: %w", err)
	}

	return nil
}

// loadRefs loads refs from disk
func loadRefs() ([]w3pilot.ElementRef, error) {
	path := getRefsPath()

	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("no element refs found (run 'w3pilot map' first)")
		}
		return nil, fmt.Errorf("failed to read refs file: %w", err)
	}

	var refs []w3pilot.ElementRef
	if err := json.Unmarshal(data, &refs); err != nil {
		return nil, fmt.Errorf("failed to parse refs file: %w", err)
	}

	return refs, nil
}

// resolveRef resolves a @ref to a CSS selector, or returns the input if not a ref
func resolveRef(selectorOrRef string) (string, error) {
	selectorOrRef = strings.TrimSpace(selectorOrRef)

	if !w3pilot.IsRef(selectorOrRef) {
		// Not a ref, return as-is
		return selectorOrRef, nil
	}

	refs, err := loadRefs()
	if err != nil {
		return "", err
	}

	for _, ref := range refs {
		if ref.Ref == selectorOrRef {
			return ref.Selector, nil
		}
	}

	return "", fmt.Errorf("unknown element reference: %s (run 'w3pilot map' to refresh)", selectorOrRef)
}

// clearRefs removes the refs file
func clearRefs() error {
	path := getRefsPath()
	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove refs file: %w", err)
	}
	return nil
}

var mapClearCmd = &cobra.Command{
	Use:   "clear",
	Short: "Clear stored element refs",
	Long:  `Clear all stored element references. Run 'w3pilot map' to refresh.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := clearRefs(); err != nil {
			return err
		}
		fmt.Println("Element references cleared.")
		return nil
	},
}

var mapGetCmd = &cobra.Command{
	Use:   "get <ref>",
	Short: "Get details about a ref",
	Long:  `Get details about a specific element reference (e.g., @e1).`,
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		refStr := args[0]

		if !w3pilot.IsRef(refStr) {
			return fmt.Errorf("invalid ref format: %s (expected @e1, @e2, etc.)", refStr)
		}

		refs, err := loadRefs()
		if err != nil {
			return err
		}

		for _, ref := range refs {
			if ref.Ref == refStr {
				Output(ref, func(data interface{}) string {
					r := data.(w3pilot.ElementRef)
					var sb strings.Builder
					sb.WriteString(fmt.Sprintf("Ref:         %s\n", r.Ref))
					sb.WriteString(fmt.Sprintf("Tag:         %s\n", r.Tag))
					if r.Role != "" {
						sb.WriteString(fmt.Sprintf("Role:        %s\n", r.Role))
					}
					if r.Text != "" {
						sb.WriteString(fmt.Sprintf("Text:        %s\n", r.Text))
					}
					if r.Label != "" {
						sb.WriteString(fmt.Sprintf("Label:       %s\n", r.Label))
					}
					if r.Placeholder != "" {
						sb.WriteString(fmt.Sprintf("Placeholder: %s\n", r.Placeholder))
					}
					if r.Type != "" {
						sb.WriteString(fmt.Sprintf("Type:        %s\n", r.Type))
					}
					sb.WriteString(fmt.Sprintf("Selector:    %s\n", r.Selector))
					sb.WriteString(fmt.Sprintf("Visible:     %t\n", r.Visible))
					sb.WriteString(fmt.Sprintf("Enabled:     %t\n", r.Enabled))
					return sb.String()
				})
				return nil
			}
		}

		return fmt.Errorf("element reference not found: %s", refStr)
	},
}

var mapListCmd = &cobra.Command{
	Use:   "list",
	Short: "List stored element refs",
	Long:  `List all stored element references without re-scanning the page.`,
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		refs, err := loadRefs()
		if err != nil {
			return err
		}

		Output(refs, func(data interface{}) string {
			r := data.([]w3pilot.ElementRef)
			return formatRefs(r)
		})
		return nil
	},
}

var mapDiffCmd = &cobra.Command{
	Use:   "diff",
	Short: "Compare current page against stored refs",
	Long: `Compare the current page's interactive elements against the previously stored refs.

This is useful for detecting what changed after an action:
  w3pilot map                  # Store initial refs
  w3pilot click @e1            # Perform action
  w3pilot map diff             # See what changed

The diff shows:
  - Added: New elements that appeared
  - Removed: Elements that disappeared
  - Changed: Elements that moved (selector changed)
  - Unchanged: Elements that stayed the same`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), mapTimeout)
		defer cancel()

		// Load previous refs
		before, err := loadRefs()
		if err != nil {
			return fmt.Errorf("no previous mapping found: %w", err)
		}

		pilot := mustGetVibe(ctx)

		opts := &w3pilot.MapOptions{
			IncludeHidden: mapIncludeHidden,
			MaxElements:   mapMaxElements,
			Scope:         mapScope,
		}

		after, err := pilot.MapElements(ctx, opts)
		if err != nil {
			return fmt.Errorf("mapping failed: %w", err)
		}

		// Calculate diff
		diff := w3pilot.DiffRefs(before, after)

		// Save new refs
		if err := saveRefs(after); err != nil {
			fmt.Fprintf(os.Stderr, "Warning: could not save refs: %v\n", err)
		}

		Output(diff, func(data interface{}) string {
			d := data.(*w3pilot.RefDiff)
			return formatDiff(d)
		})
		return nil
	},
}

func formatDiff(diff *w3pilot.RefDiff) string {
	var sb strings.Builder

	if !diff.HasChanges() {
		sb.WriteString(fmt.Sprintf("No changes detected (%d elements unchanged)\n", diff.Summary.Unchanged))
		return sb.String()
	}

	sb.WriteString("Page Changes:\n\n")

	if len(diff.Added) > 0 {
		sb.WriteString(fmt.Sprintf("ADDED (%d):\n", len(diff.Added)))
		for _, ref := range diff.Added {
			sb.WriteString(fmt.Sprintf("  + %s\n", ref.FormatRef()))
		}
		sb.WriteString("\n")
	}

	if len(diff.Removed) > 0 {
		sb.WriteString(fmt.Sprintf("REMOVED (%d):\n", len(diff.Removed)))
		for _, ref := range diff.Removed {
			sb.WriteString(fmt.Sprintf("  - %s\n", ref.FormatRef()))
		}
		sb.WriteString("\n")
	}

	if len(diff.Changed) > 0 {
		sb.WriteString(fmt.Sprintf("MOVED (%d):\n", len(diff.Changed)))
		for _, change := range diff.Changed {
			sb.WriteString(fmt.Sprintf("  ~ %s\n", change.After.FormatRef()))
			sb.WriteString(fmt.Sprintf("    was: %s\n", change.Before.Selector))
			sb.WriteString(fmt.Sprintf("    now: %s\n", change.After.Selector))
		}
		sb.WriteString("\n")
	}

	sb.WriteString(fmt.Sprintf("Summary: %d added, %d removed, %d moved, %d unchanged\n",
		diff.Summary.Added, diff.Summary.Removed, diff.Summary.Changed, diff.Summary.Unchanged))

	return sb.String()
}

func init() {
	rootCmd.AddCommand(mapCmd)
	mapCmd.AddCommand(mapClearCmd)
	mapCmd.AddCommand(mapGetCmd)
	mapCmd.AddCommand(mapListCmd)
	mapCmd.AddCommand(mapDiffCmd)

	mapCmd.Flags().DurationVar(&mapTimeout, "timeout", 30*time.Second, "Timeout")
	mapCmd.Flags().BoolVar(&mapIncludeHidden, "include-hidden", false, "Include hidden elements")
	mapCmd.Flags().IntVar(&mapMaxElements, "max", 100, "Maximum elements to map (0 = no limit)")
	mapCmd.Flags().StringVar(&mapScope, "scope", "", "CSS selector to limit mapping scope")

	// Diff command inherits the same flags
	mapDiffCmd.Flags().DurationVar(&mapTimeout, "timeout", 30*time.Second, "Timeout")
	mapDiffCmd.Flags().BoolVar(&mapIncludeHidden, "include-hidden", false, "Include hidden elements")
	mapDiffCmd.Flags().IntVar(&mapMaxElements, "max", 100, "Maximum elements to map (0 = no limit)")
	mapDiffCmd.Flags().StringVar(&mapScope, "scope", "", "CSS selector to limit mapping scope")
}
