//nolint:dupl // page commands share similar structure
package cmd

import (
	"context"
	"fmt"
	"time"

	"github.com/spf13/cobra"
)

var pageTitleTimeout time.Duration

// TitleResult represents the result of getting the page title.
type TitleResult struct {
	Title string `json:"title"`
}

var pageTitleCmd = &cobra.Command{
	Use:   "title",
	Short: "Get the page title",
	Long: `Get the title of the current page.

Examples:
  w3pilot page title
  w3pilot page title --format json`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithTimeout(context.Background(), pageTitleTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		title, err := pilot.Title(ctx)
		if err != nil {
			return fmt.Errorf("failed to get title: %w", err)
		}

		Output(TitleResult{Title: title}, func(data interface{}) string {
			return data.(TitleResult).Title
		})
		return nil
	},
}

func init() {
	pageCmd.AddCommand(pageTitleCmd)
	pageTitleCmd.Flags().DurationVar(&pageTitleTimeout, "timeout", 10*time.Second, "Timeout")
}
