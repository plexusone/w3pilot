package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	w3pilot "github.com/plexusone/w3pilot"
	"github.com/spf13/cobra"
)

var (
	pagePDFTimeout    time.Duration
	pagePDFLandscape  bool
	pagePDFBackground bool
	pagePDFFormat     string
)

var pagePDFCmd = &cobra.Command{
	Use:   "pdf <filename>",
	Short: "Generate PDF of the page",
	Long: `Generate a PDF of the current page.

Examples:
  w3pilot page pdf page.pdf
  w3pilot page pdf page.pdf --landscape
  w3pilot page pdf page.pdf --format A4 --background`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		filename := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), pagePDFTimeout)
		defer cancel()

		pilot := mustGetVibe(ctx)

		opts := &w3pilot.PDFOptions{
			Landscape:       pagePDFLandscape,
			PrintBackground: pagePDFBackground,
		}
		if pagePDFFormat != "" {
			opts.Format = pagePDFFormat
		}

		data, err := pilot.PDF(ctx, opts)
		if err != nil {
			return fmt.Errorf("PDF generation failed: %w", err)
		}

		if err := os.WriteFile(filename, data, 0600); err != nil {
			return fmt.Errorf("failed to write file: %w", err)
		}

		fmt.Printf("PDF saved: %s (%d bytes)\n", filename, len(data))
		return nil
	},
}

func init() {
	pageCmd.AddCommand(pagePDFCmd)
	pagePDFCmd.Flags().DurationVar(&pagePDFTimeout, "timeout", 60*time.Second, "Timeout")
	pagePDFCmd.Flags().BoolVar(&pagePDFLandscape, "landscape", false, "Landscape orientation")
	pagePDFCmd.Flags().BoolVar(&pagePDFBackground, "background", false, "Print background graphics")
	pagePDFCmd.Flags().StringVar(&pagePDFFormat, "format", "", "Paper format (A4, Letter, etc.)")
}
