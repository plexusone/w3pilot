package cdp

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// BrowserInfo contains information about a Chrome browser instance.
type BrowserInfo struct {
	Port              int
	BrowserWSEndpoint string
	Browser           string
	ProtocolVersion   string
	UserAgent         string
	V8Version         string
	WebKitVersion     string
}

// DiscoverFromUserDataDir reads the DevToolsActivePort file from Chrome's user data directory.
// Returns the CDP port and browser WebSocket endpoint.
func DiscoverFromUserDataDir(userDataDir string) (port int, wsEndpoint string, err error) {
	portFile := filepath.Join(userDataDir, "DevToolsActivePort")

	data, err := os.ReadFile(portFile)
	if err != nil {
		return 0, "", fmt.Errorf("cdp: failed to read DevToolsActivePort: %w", err)
	}

	lines := strings.Split(strings.TrimSpace(string(data)), "\n")
	if len(lines) < 2 {
		return 0, "", fmt.Errorf("cdp: invalid DevToolsActivePort format")
	}

	port, err = strconv.Atoi(strings.TrimSpace(lines[0]))
	if err != nil {
		return 0, "", fmt.Errorf("cdp: invalid port in DevToolsActivePort: %w", err)
	}

	wsEndpoint = fmt.Sprintf("ws://localhost:%d%s", port, strings.TrimSpace(lines[1]))

	return port, wsEndpoint, nil
}

// DiscoverFromProcess attempts to find the DevToolsActivePort by examining
// the Chrome process command line for --user-data-dir.
func DiscoverFromProcess(pid int) (port int, wsEndpoint string, err error) {
	// Read /proc/{pid}/cmdline on Linux, or use ps on macOS
	// For now, this is a placeholder - the primary method is DiscoverFromUserDataDir
	return 0, "", fmt.Errorf("cdp: process discovery not implemented")
}

// GetBrowserInfo fetches browser information from the CDP /json/version endpoint.
func GetBrowserInfo(ctx context.Context, port int) (*BrowserInfo, error) {
	url := fmt.Sprintf("http://localhost:%d/json/version", port)

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cdp: failed to get browser info: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cdp: unexpected status %d from /json/version", resp.StatusCode)
	}

	var info struct {
		Browser              string `json:"Browser"`
		ProtocolVersion      string `json:"Protocol-Version"`
		UserAgent            string `json:"User-Agent"`
		V8Version            string `json:"V8-Version"`
		WebKitVersion        string `json:"WebKit-Version"`
		WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
	}

	if err := readJSON(resp.Body, &info); err != nil {
		return nil, fmt.Errorf("cdp: failed to parse browser info: %w", err)
	}

	return &BrowserInfo{
		Port:              port,
		BrowserWSEndpoint: info.WebSocketDebuggerURL,
		Browser:           info.Browser,
		ProtocolVersion:   info.ProtocolVersion,
		UserAgent:         info.UserAgent,
		V8Version:         info.V8Version,
		WebKitVersion:     info.WebKitVersion,
	}, nil
}

// Target represents a CDP target (page, worker, etc.).
type Target struct {
	ID                   string `json:"id"`
	Type                 string `json:"type"`
	Title                string `json:"title"`
	URL                  string `json:"url"`
	WebSocketDebuggerURL string `json:"webSocketDebuggerUrl"`
	DevtoolsFrontendURL  string `json:"devtoolsFrontendUrl"`
}

// ListTargets fetches the list of available CDP targets.
func ListTargets(ctx context.Context, port int) ([]Target, error) {
	url := fmt.Sprintf("http://localhost:%d/json/list", port)

	if ctx == nil {
		ctx = context.Background()
	}

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("cdp: failed to list targets: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("cdp: unexpected status %d from /json/list", resp.StatusCode)
	}

	var targets []Target
	if err := readJSON(resp.Body, &targets); err != nil {
		return nil, fmt.Errorf("cdp: failed to parse targets: %w", err)
	}

	return targets, nil
}

// FindPageTarget finds the first page target matching the given URL.
// If url is empty, returns the first page target.
func FindPageTarget(ctx context.Context, port int, url string) (*Target, error) {
	targets, err := ListTargets(ctx, port)
	if err != nil {
		return nil, err
	}

	for _, t := range targets {
		if t.Type == "page" {
			if url == "" || strings.Contains(t.URL, url) {
				return &t, nil
			}
		}
	}

	return nil, fmt.Errorf("cdp: no page target found")
}

// readJSON is a helper to decode JSON from an io.Reader.
func readJSON(r interface{ Read([]byte) (int, error) }, v interface{}) error {
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 1024*1024), 1024*1024)

	var data []byte
	for scanner.Scan() {
		data = append(data, scanner.Bytes()...)
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	return json.Unmarshal(data, v)
}
