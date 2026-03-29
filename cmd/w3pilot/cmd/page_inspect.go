package cmd

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/spf13/cobra"

	w3pilot "github.com/plexusone/w3pilot"
)

var (
	pageInspectTimeout         time.Duration
	pageInspectIncludeButtons  bool
	pageInspectIncludeLinks    bool
	pageInspectIncludeInputs   bool
	pageInspectIncludeSelects  bool
	pageInspectIncludeHeadings bool
	pageInspectIncludeImages   bool
	pageInspectMaxItems        int
)

var pageInspectCmd = &cobra.Command{
	Use:   "inspect",
	Short: "Inspect page elements",
	Long: `Inspect the current page to discover interactive elements.

This command helps AI agents understand the page structure by listing:
  - Buttons (button, input[type=submit], [role=button])
  - Links (a[href])
  - Inputs (input, textarea)
  - Selects (select)
  - Headings (h1-h6)
  - Images with alt text

Examples:
  w3pilot page inspect
  w3pilot page inspect --format json
  w3pilot page inspect --no-links --no-images
  w3pilot page inspect --max-items 100`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), pageInspectTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		opts := &w3pilot.InspectOptions{
			IncludeButtons:  pageInspectIncludeButtons,
			IncludeLinks:    pageInspectIncludeLinks,
			IncludeInputs:   pageInspectIncludeInputs,
			IncludeSelects:  pageInspectIncludeSelects,
			IncludeHeadings: pageInspectIncludeHeadings,
			IncludeImages:   pageInspectIncludeImages,
			MaxItems:        pageInspectMaxItems,
		}

		result, err := pilot.Inspect(ctx, opts)
		if err != nil {
			return fmt.Errorf("inspection failed: %w", err)
		}

		Output(result, func(data interface{}) string {
			r := data.(*w3pilot.InspectResult)
			return formatInspectResult(r)
		})
		return nil
	},
}

func formatInspectResult(r *w3pilot.InspectResult) string {
	var sb strings.Builder

	sb.WriteString(fmt.Sprintf("Page: %s\n", r.Title))
	sb.WriteString(fmt.Sprintf("URL: %s\n\n", r.URL))

	// Headings
	if len(r.Headings) > 0 {
		sb.WriteString("HEADINGS:\n")
		for _, h := range r.Headings {
			visibility := ""
			if !h.Visible {
				visibility = " (hidden)"
			}
			sb.WriteString(fmt.Sprintf("  H%d: %s%s\n", h.Level, truncate(h.Text, 60), visibility))
		}
		sb.WriteString("\n")
	}

	// Buttons
	if len(r.Buttons) > 0 {
		sb.WriteString("BUTTONS:\n")
		for _, btn := range r.Buttons {
			status := ""
			if btn.Disabled {
				status = " (disabled)"
			} else if !btn.Visible {
				status = " (hidden)"
			}
			sb.WriteString(fmt.Sprintf("  [%s] %s%s\n", btn.Selector, truncate(btn.Text, 40), status))
		}
		sb.WriteString("\n")
	}

	// Links
	if len(r.Links) > 0 {
		sb.WriteString("LINKS:\n")
		for _, link := range r.Links {
			visibility := ""
			if !link.Visible {
				visibility = " (hidden)"
			}
			sb.WriteString(fmt.Sprintf("  [%s] %s -> %s%s\n",
				link.Selector, truncate(link.Text, 30), truncate(link.Href, 40), visibility))
		}
		sb.WriteString("\n")
	}

	// Inputs
	if len(r.Inputs) > 0 {
		sb.WriteString("INPUTS:\n")
		for _, input := range r.Inputs {
			label := input.Label
			if label == "" {
				label = input.Placeholder
			}
			if label == "" {
				label = input.Name
			}
			status := ""
			if input.Disabled {
				status = " (disabled)"
			} else if input.ReadOnly {
				status = " (readonly)"
			} else if !input.Visible {
				status = " (hidden)"
			}
			required := ""
			if input.Required {
				required = "*"
			}
			sb.WriteString(fmt.Sprintf("  [%s] %s (%s)%s%s\n",
				input.Selector, label, input.Type, required, status))
		}
		sb.WriteString("\n")
	}

	// Selects
	if len(r.Selects) > 0 {
		sb.WriteString("SELECTS:\n")
		for _, sel := range r.Selects {
			label := sel.Label
			if label == "" {
				label = sel.Name
			}
			status := ""
			if sel.Disabled {
				status = " (disabled)"
			} else if !sel.Visible {
				status = " (hidden)"
			}
			optionsStr := ""
			if len(sel.Options) > 0 {
				optionsStr = fmt.Sprintf(" [%s]", strings.Join(sel.Options[:min(3, len(sel.Options))], ", "))
				if len(sel.Options) > 3 {
					optionsStr += "..."
				}
			}
			sb.WriteString(fmt.Sprintf("  [%s] %s%s%s\n", sel.Selector, label, optionsStr, status))
		}
		sb.WriteString("\n")
	}

	// Images
	if len(r.Images) > 0 {
		sb.WriteString("IMAGES:\n")
		for _, img := range r.Images {
			visibility := ""
			if !img.Visible {
				visibility = " (hidden)"
			}
			sb.WriteString(fmt.Sprintf("  [%s] %s%s\n", img.Selector, truncate(img.Alt, 50), visibility))
		}
		sb.WriteString("\n")
	}

	// Summary
	sb.WriteString(fmt.Sprintf("SUMMARY: %d buttons, %d links, %d inputs, %d selects, %d headings, %d images\n",
		r.Summary.TotalButtons, r.Summary.TotalLinks, r.Summary.TotalInputs,
		r.Summary.TotalSelects, r.Summary.TotalHeadings, r.Summary.TotalImages))

	return sb.String()
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func init() {
	pageCmd.AddCommand(pageInspectCmd)
	pageInspectCmd.Flags().DurationVar(&pageInspectTimeout, "timeout", 30*time.Second, "Timeout")
	pageInspectCmd.Flags().BoolVar(&pageInspectIncludeButtons, "buttons", true, "Include buttons")
	pageInspectCmd.Flags().BoolVar(&pageInspectIncludeLinks, "links", true, "Include links")
	pageInspectCmd.Flags().BoolVar(&pageInspectIncludeInputs, "inputs", true, "Include inputs")
	pageInspectCmd.Flags().BoolVar(&pageInspectIncludeSelects, "selects", true, "Include selects")
	pageInspectCmd.Flags().BoolVar(&pageInspectIncludeHeadings, "headings", true, "Include headings")
	pageInspectCmd.Flags().BoolVar(&pageInspectIncludeImages, "images", true, "Include images")
	pageInspectCmd.Flags().IntVar(&pageInspectMaxItems, "max-items", 50, "Max items per category")

	// Convenience flags for exclusion
	pageInspectCmd.Flags().Bool("no-buttons", false, "Exclude buttons")
	pageInspectCmd.Flags().Bool("no-links", false, "Exclude links")
	pageInspectCmd.Flags().Bool("no-inputs", false, "Exclude inputs")
	pageInspectCmd.Flags().Bool("no-selects", false, "Exclude selects")
	pageInspectCmd.Flags().Bool("no-headings", false, "Exclude headings")
	pageInspectCmd.Flags().Bool("no-images", false, "Exclude images")

	// Pre-run to handle --no-* flags
	pageInspectCmd.PreRun = func(cmd *cobra.Command, args []string) {
		if v, _ := cmd.Flags().GetBool("no-buttons"); v {
			pageInspectIncludeButtons = false
		}
		if v, _ := cmd.Flags().GetBool("no-links"); v {
			pageInspectIncludeLinks = false
		}
		if v, _ := cmd.Flags().GetBool("no-inputs"); v {
			pageInspectIncludeInputs = false
		}
		if v, _ := cmd.Flags().GetBool("no-selects"); v {
			pageInspectIncludeSelects = false
		}
		if v, _ := cmd.Flags().GetBool("no-headings"); v {
			pageInspectIncludeHeadings = false
		}
		if v, _ := cmd.Flags().GetBool("no-images"); v {
			pageInspectIncludeImages = false
		}
	}
}
