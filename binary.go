package w3pilot

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// ClickerBinaryEnvVar is the environment variable for specifying the clicker binary path.
const ClickerBinaryEnvVar = "CLICKER_BIN_PATH"

// FindClickerBinary locates the clicker binary.
//
// Search order:
//  1. CLICKER_BIN_PATH environment variable
//  2. PATH lookup for "clicker"
//  3. Go bin directory
//  4. Common installation paths
func FindClickerBinary() (string, error) {
	// 1. Check environment variable
	if envPath := os.Getenv(ClickerBinaryEnvVar); envPath != "" {
		cleanPath := filepath.Clean(envPath)
		if _, err := os.Stat(cleanPath); err == nil {
			return cleanPath, nil
		}
		return "", fmt.Errorf("%s is set to %q but file does not exist", ClickerBinaryEnvVar, envPath)
	}

	// 2. Check PATH
	if path, err := exec.LookPath("clicker"); err == nil {
		return path, nil
	}

	// 3. Check Go bin directory
	home, _ := os.UserHomeDir()
	if home != "" {
		goBin := filepath.Join(home, "go", "bin", "clicker")
		if runtime.GOOS == "windows" {
			goBin = filepath.Join(home, "go", "bin", "clicker.exe")
		}
		if _, err := os.Stat(goBin); err == nil {
			return goBin, nil
		}
	}

	// 4. Check common installation paths
	commonPaths := getCommonClickerPaths()
	for _, path := range commonPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	return "", fmt.Errorf("clicker binary not found; install from: https://github.com/anthropics/vibium, or set %s", ClickerBinaryEnvVar)
}

// getCommonClickerPaths returns common installation paths for the clicker binary.
func getCommonClickerPaths() []string {
	var paths []string

	home, _ := os.UserHomeDir()

	switch runtime.GOOS {
	case "darwin", "linux":
		paths = []string{
			"/usr/local/bin/clicker",
			"/usr/bin/clicker",
		}
		if home != "" {
			paths = append(paths,
				filepath.Join(home, ".local", "bin", "clicker"),
				filepath.Join(home, "bin", "clicker"),
			)
		}
	case "windows":
		paths = []string{
			"C:\\Program Files\\vibium\\clicker.exe",
			"C:\\Program Files (x86)\\vibium\\clicker.exe",
		}
		if home != "" {
			paths = append(paths,
				filepath.Join(home, "AppData", "Local", "Programs", "vibium", "clicker.exe"),
			)
		}
	}

	return paths
}

// ClickerVersion returns the version of the clicker binary.
func ClickerVersion(clickerPath string) (string, error) {
	cmd := exec.Command(clickerPath, "--version")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("failed to get clicker version: %w", err)
	}
	return string(output), nil
}
