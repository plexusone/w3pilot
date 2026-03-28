package cdp

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"runtime"
	"strings"
)

// DiscoverFromRunningChrome finds the CDP port and page endpoint from a running Chrome process.
// This works by examining process command lines for --user-data-dir.
func DiscoverFromRunningChrome() (port int, wsEndpoint string, err error) {
	userDataDir, err := findChromeUserDataDir()
	if err != nil {
		return 0, "", err
	}

	port, _, err = DiscoverFromUserDataDir(userDataDir)
	if err != nil {
		return 0, "", err
	}

	// Find a page target to connect to (browser endpoint doesn't support most CDP domains)
	targets, err := ListTargets(nil, port)
	if err != nil {
		return 0, "", fmt.Errorf("cdp: failed to list targets: %w", err)
	}

	// Find the first page target
	for _, t := range targets {
		if t.Type == "page" && t.WebSocketDebuggerURL != "" {
			return port, t.WebSocketDebuggerURL, nil
		}
	}

	// Fallback to browser endpoint if no page found
	_, browserEndpoint, _ := DiscoverFromUserDataDir(userDataDir)
	return port, browserEndpoint, nil
}

// findChromeUserDataDir finds the user data directory from running Chrome processes.
func findChromeUserDataDir() (string, error) {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin", "linux":
		// Use ps to find Chrome processes with --user-data-dir
		cmd = exec.Command("sh", "-c", `ps aux | grep -E "[c]hrome.*--user-data-dir" | head -1`)
	case "windows":
		// Use wmic on Windows
		cmd = exec.Command("wmic", "process", "where", "name like '%chrome%'", "get", "commandline")
	default:
		return "", fmt.Errorf("cdp: unsupported platform: %s", runtime.GOOS)
	}

	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("cdp: failed to find Chrome process: %w", err)
	}

	output := out.String()
	if output == "" {
		return "", fmt.Errorf("cdp: no Chrome process found")
	}

	// Extract --user-data-dir value
	re := regexp.MustCompile(`--user-data-dir=([^\s]+)`)
	matches := re.FindStringSubmatch(output)
	if len(matches) < 2 {
		return "", fmt.Errorf("cdp: no user-data-dir found in Chrome process")
	}

	userDataDir := strings.TrimSpace(matches[1])
	return userDataDir, nil
}

// DiscoverAny attempts to discover CDP using multiple methods.
// It tries running Chrome processes first, then falls back to common locations.
func DiscoverAny() (port int, wsEndpoint string, err error) {
	// Try running Chrome first
	port, wsEndpoint, err = DiscoverFromRunningChrome()
	if err == nil {
		return port, wsEndpoint, nil
	}

	// Could add more discovery methods here (e.g., common ports)
	return 0, "", fmt.Errorf("cdp: could not discover Chrome DevTools: %w", err)
}
