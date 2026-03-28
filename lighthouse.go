package webpilot

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// LighthouseCategory represents a Lighthouse audit category.
type LighthouseCategory string

const (
	LighthousePerformance    LighthouseCategory = "performance"
	LighthouseAccessibility  LighthouseCategory = "accessibility"
	LighthouseBestPractices  LighthouseCategory = "best-practices"
	LighthouseSEO            LighthouseCategory = "seo"
)

// LighthouseDevice represents the device to emulate.
type LighthouseDevice string

const (
	LighthouseDesktop LighthouseDevice = "desktop"
	LighthouseMobile  LighthouseDevice = "mobile"
)

// LighthouseOptions configures a Lighthouse audit.
type LighthouseOptions struct {
	// Categories to audit (default: all)
	Categories []LighthouseCategory
	// Device to emulate (default: desktop)
	Device LighthouseDevice
	// OutputDir for reports (default: temp directory)
	OutputDir string
	// Port for Chrome debugging (default: use Pilot's CDP port)
	Port int
}

// LighthouseResult contains the results of a Lighthouse audit.
type LighthouseResult struct {
	// URL that was audited
	URL string `json:"url"`
	// Device used for emulation
	Device string `json:"device"`
	// Scores by category (0-100)
	Scores map[string]LighthouseScore `json:"scores"`
	// Timing information
	TotalDurationMS float64 `json:"totalDurationMs"`
	// Report paths
	JSONReportPath string `json:"jsonReportPath,omitempty"`
	HTMLReportPath string `json:"htmlReportPath,omitempty"`
	// Summary of audits
	PassedAudits int `json:"passedAudits"`
	FailedAudits int `json:"failedAudits"`
}

// LighthouseScore represents a category score.
type LighthouseScore struct {
	Title string  `json:"title"`
	Score float64 `json:"score"` // 0-1
}

// LighthouseAudit runs a Lighthouse audit on the current page.
// Requires lighthouse CLI: npm install -g lighthouse
func (p *Pilot) LighthouseAudit(ctx context.Context, opts *LighthouseOptions) (*LighthouseResult, error) {
	if opts == nil {
		opts = &LighthouseOptions{}
	}

	// Get current URL
	currentURL, err := p.URL(ctx)
	if err != nil {
		return nil, fmt.Errorf("webpilot: failed to get current URL: %w", err)
	}

	// Set defaults
	if len(opts.Categories) == 0 {
		opts.Categories = []LighthouseCategory{
			LighthouseAccessibility,
			LighthouseBestPractices,
			LighthouseSEO,
		}
	}
	if opts.Device == "" {
		opts.Device = LighthouseDesktop
	}
	if opts.Port == 0 {
		opts.Port = p.CDPPort()
	}
	if opts.OutputDir == "" {
		opts.OutputDir = os.TempDir()
	}

	// Find lighthouse binary
	lighthousePath, err := findLighthouseBinary()
	if err != nil {
		return nil, err
	}

	// Build category list
	var categories []string
	for _, cat := range opts.Categories {
		categories = append(categories, string(cat))
	}

	// Create output paths
	baseName := fmt.Sprintf("lighthouse-%d", os.Getpid())
	jsonPath := filepath.Join(opts.OutputDir, baseName+".json")
	htmlPath := filepath.Join(opts.OutputDir, baseName+".html")

	// Build command
	args := []string{
		currentURL,
		"--port=" + fmt.Sprintf("%d", opts.Port),
		"--output=json,html",
		"--output-path=" + filepath.Join(opts.OutputDir, baseName),
		"--only-categories=" + strings.Join(categories, ","),
		"--chrome-flags=--headless",
		"--no-enable-error-reporting",
		"--quiet",
	}

	// Add device-specific flags
	if opts.Device == LighthouseDesktop {
		args = append(args, "--preset=desktop")
	} else {
		args = append(args, "--preset=mobile", "--emulated-form-factor=mobile")
	}

	// Execute lighthouse
	cmd := exec.CommandContext(ctx, lighthousePath, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return nil, fmt.Errorf("webpilot: lighthouse failed: %w\nstderr: %s", err, stderr.String())
	}

	// Parse JSON report
	jsonData, err := os.ReadFile(jsonPath)
	if err != nil {
		return nil, fmt.Errorf("webpilot: failed to read lighthouse report: %w", err)
	}

	var lhr lighthouseReport
	if err := json.Unmarshal(jsonData, &lhr); err != nil {
		return nil, fmt.Errorf("webpilot: failed to parse lighthouse report: %w", err)
	}

	// Build result
	result := &LighthouseResult{
		URL:             lhr.FinalURL,
		Device:          string(opts.Device),
		Scores:          make(map[string]LighthouseScore),
		TotalDurationMS: lhr.Timing.Total,
		JSONReportPath:  jsonPath,
		HTMLReportPath:  htmlPath,
	}

	// Extract category scores
	for id, cat := range lhr.Categories {
		score := float64(0)
		if cat.Score != nil {
			score = *cat.Score
		}
		result.Scores[id] = LighthouseScore{
			Title: cat.Title,
			Score: score,
		}
	}

	// Count audits
	for _, audit := range lhr.Audits {
		if audit.Score != nil {
			if *audit.Score >= 0.9 {
				result.PassedAudits++
			} else {
				result.FailedAudits++
			}
		}
	}

	return result, nil
}

// lighthouseReport is the internal structure for parsing Lighthouse JSON output.
type lighthouseReport struct {
	FinalURL   string `json:"finalUrl"`
	Categories map[string]struct {
		Title string   `json:"title"`
		Score *float64 `json:"score"`
	} `json:"categories"`
	Audits map[string]struct {
		Score *float64 `json:"score"`
	} `json:"audits"`
	Timing struct {
		Total float64 `json:"total"`
	} `json:"timing"`
}

// findLighthouseBinary finds the lighthouse CLI binary.
func findLighthouseBinary() (string, error) {
	// Try direct lighthouse command first
	if path, err := exec.LookPath("lighthouse"); err == nil {
		return path, nil
	}

	// Try npx
	if path, err := exec.LookPath("npx"); err == nil {
		// Check if lighthouse is available via npx
		cmd := exec.Command(path, "lighthouse", "--version")
		if err := cmd.Run(); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("webpilot: lighthouse CLI not found. Install with: npm install -g lighthouse")
}

// runLighthouseWithNpx runs lighthouse via npx.
func runLighthouseWithNpx(ctx context.Context, args []string) ([]byte, error) {
	npxPath, err := exec.LookPath("npx")
	if err != nil {
		return nil, fmt.Errorf("npx not found: %w", err)
	}

	fullArgs := append([]string{"lighthouse"}, args...)
	cmd := exec.CommandContext(ctx, npxPath, fullArgs...)
	return cmd.Output()
}
