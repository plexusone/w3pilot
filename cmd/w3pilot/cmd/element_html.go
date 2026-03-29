package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var (
	elementHTMLTimeout time.Duration
	elementHTMLOuter   bool
)

var elementHTMLCmd = &cobra.Command{
	Use:   "html <selector>",
	Short: "Get element HTML",
	Long: `Get the HTML of an element.

By default returns innerHTML. Use --outer for outerHTML.

Examples:
  w3pilot element html "#content"
  w3pilot element html "#container" --outer`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		selector := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), elementHTMLTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		el, err := pilot.Find(ctx, selector, nil)
		if err != nil {
			return fmt.Errorf("element not found: %w", err)
		}

		var html string
		if elementHTMLOuter {
			html, err = el.HTML(ctx)
		} else {
			html, err = el.InnerHTML(ctx)
		}
		if err != nil {
			return fmt.Errorf("failed to get HTML: %w", err)
		}

		fmt.Println(html)
		return nil
	},
}

func init() {
	elementCmd.AddCommand(elementHTMLCmd)
	elementHTMLCmd.Flags().DurationVar(&elementHTMLTimeout, "timeout", 10*time.Second, "Timeout")
	elementHTMLCmd.Flags().BoolVar(&elementHTMLOuter, "outer", false, "Get outerHTML instead of innerHTML")
}
