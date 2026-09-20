package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/spf13/cobra"
)

var (
	fetchHTML     bool
	fetchSelector string
	fetchOutput   string
	fetchWait     string
	fetchHeadless bool
	fetchTimeout  time.Duration
	fetchDelay    time.Duration
)

// FetchResult represents the result of a one-shot page fetch.
type FetchResult struct {
	URL     string `json:"url"`
	Title   string `json:"title,omitempty"`
	Content string `json:"content"`
}

var fetchCmd = &cobra.Command{
	Use:   "fetch <url>",
	Short: "Fetch a page's rendered text or HTML in one shot",
	Long: `Fetch launches a headless browser, navigates to a URL, waits for it to
load, and prints the page's rendered text (or HTML) — then quits. It is
self-contained: unlike the 'page' subcommands, it does not need a browser
started by 'browser launch', because each CLI invocation is its own process.

Think of it as "curl for JavaScript-rendered pages": the output reflects the
DOM after scripts run, not the raw server HTML.

Examples:
  w3pilot fetch https://example.com                       # visible text to stdout
  w3pilot fetch https://example.com --html                # full rendered HTML
  w3pilot fetch https://example.com -s "main" -O out.txt  # text of <main> to a file
  w3pilot fetch https://spa.example.com --wait networkidle # wait for network to settle
  w3pilot fetch https://example.com --format json         # {url,title,content}`,
	Args: cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		url := args[0]

		ctx, cancel := context.WithTimeout(context.Background(), fetchTimeout)
		defer cancel()

		pilot, err := launchBrowser(ctx, fetchHeadless)
		if err != nil {
			return err
		}
		defer func() {
			_ = pilot.Quit(context.Background())
			_ = clearSession()
		}()

		if err := pilot.Go(ctx, url); err != nil {
			return fmt.Errorf("navigation failed: %w", err)
		}

		// Go already waits for the default load event; wait again only when the
		// caller asked for a stricter state (e.g. networkidle) than the default.
		if fetchWait != "" && fetchWait != "load" {
			if err := pilot.WaitForLoad(ctx, fetchWait, fetchTimeout); err != nil {
				return fmt.Errorf("wait for %q failed: %w", fetchWait, err)
			}
		}

		// An optional settle delay lets client-rendered pages hydrate and lets
		// interstitials (e.g. a bot-challenge redirect) resolve before we read.
		if fetchDelay > 0 {
			select {
			case <-time.After(fetchDelay):
			case <-ctx.Done():
				return fmt.Errorf("delay interrupted: %w", ctx.Err())
			}
		}

		content, err := extractContent(ctx, pilot)
		if err != nil {
			return err
		}

		if fetchOutput != "" {
			if err := os.WriteFile(fetchOutput, []byte(content), 0600); err != nil {
				return fmt.Errorf("failed to write %s: %w", fetchOutput, err)
			}
			if verbose {
				fmt.Fprintf(os.Stderr, "Wrote %d bytes to %s\n", len(content), fetchOutput)
			}
		}

		title, _ := pilot.Title(ctx)
		result := FetchResult{URL: url, Title: title, Content: content}

		if fetchOutput != "" && GetOutputFormat() != FormatJSON {
			// Content already went to the file; keep stdout quiet in text mode.
			return nil
		}

		Output(result, func(data interface{}) string {
			return data.(FetchResult).Content
		})
		return nil
	},
}

// extractContent returns the requested slice of the loaded page: visible text
// or HTML, for the whole document or a single selector.
func extractContent(ctx context.Context, pilot interface {
	Content(context.Context) (string, error)
	Evaluate(context.Context, string) (interface{}, error)
}) (string, error) {
	// Whole-page HTML has a dedicated Pilot method that returns the serialized DOM.
	if fetchHTML && fetchSelector == "" {
		return pilot.Content(ctx)
	}

	var script string
	switch {
	case fetchHTML && fetchSelector != "":
		script = fmt.Sprintf(
			`(() => { const el = document.querySelector(%s); return el ? el.outerHTML : ""; })()`,
			jsString(fetchSelector))
	case !fetchHTML && fetchSelector != "":
		script = fmt.Sprintf(
			`(() => { const el = document.querySelector(%s); return el ? el.innerText : ""; })()`,
			jsString(fetchSelector))
	default: // text, whole page
		script = `(document.body && document.body.innerText) || ""`
	}

	result, err := pilot.Evaluate(ctx, script)
	if err != nil {
		return "", fmt.Errorf("content extraction failed: %w", err)
	}
	if result == nil {
		return "", nil
	}
	if s, ok := result.(string); ok {
		return s, nil
	}
	return fmt.Sprintf("%v", result), nil
}

// jsString renders s as a safe JavaScript string literal (JSON strings are a
// valid subset of JS string literals), so selectors can't break out of the
// evaluated snippet.
func jsString(s string) string {
	b, err := json.Marshal(s)
	if err != nil {
		// json.Marshal of a string cannot fail; fall back defensively.
		return `""`
	}
	return string(b)
}

func init() {
	rootCmd.AddCommand(fetchCmd)
	fetchCmd.Flags().BoolVar(&fetchHTML, "html", false, "Output rendered HTML instead of visible text")
	fetchCmd.Flags().StringVarP(&fetchSelector, "selector", "s", "", "Restrict extraction to a CSS selector")
	fetchCmd.Flags().StringVarP(&fetchOutput, "output", "O", "", "Write result to a file instead of stdout")
	fetchCmd.Flags().StringVarP(&fetchWait, "wait", "w", "load", "Load state to wait for: load, domcontentloaded, networkidle")
	fetchCmd.Flags().BoolVar(&fetchHeadless, "headless", true, "Run the browser headless")
	fetchCmd.Flags().DurationVar(&fetchTimeout, "timeout", 45*time.Second, "Total fetch timeout")
	fetchCmd.Flags().DurationVar(&fetchDelay, "delay", 0, "Settle delay after load before extracting (e.g. 5s), for SPAs/interstitials")
}
