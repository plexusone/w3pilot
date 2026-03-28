// Command w3pilot-mcp provides an MCP server for browser automation.
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/plexusone/w3pilot/mcp"
)

// stringSlice implements flag.Value for repeated string flags
type stringSlice []string

func (s *stringSlice) String() string {
	return fmt.Sprintf("%v", *s)
}

func (s *stringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func main() {
	headless := flag.Bool("headless", true, "Run browser in headless mode")
	project := flag.String("project", "w3pilot-tests", "Project name for reports")
	timeout := flag.Duration("timeout", 30*time.Second, "Default timeout for browser operations")
	listTools := flag.Bool("list-tools", false, "Output tool definitions as JSON and exit")

	var initScriptPaths stringSlice
	flag.Var(&initScriptPaths, "init-script", "JavaScript file to inject before page scripts (can be repeated)")

	flag.Parse()

	// If --list-tools is specified, output tools and exit
	if *listTools {
		data, err := mcp.ListToolsJSON()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(string(data))
		return
	}

	// Load init scripts from files
	var initScripts []string
	for _, path := range initScriptPaths {
		content, err := os.ReadFile(path)
		if err != nil {
			log.Fatalf("Failed to read init script %s: %v", path, err)
		}
		initScripts = append(initScripts, string(content))
	}

	config := mcp.Config{
		Headless:       *headless,
		Project:        *project,
		DefaultTimeout: *timeout,
		InitScripts:    initScripts,
	}

	server := mcp.NewServer(config)

	// Set up signal handling for graceful shutdown
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		log.Println("Shutting down...")
		cancel()
		if err := server.Close(context.Background()); err != nil {
			log.Printf("Error closing server: %v", err)
		}
	}()

	if err := server.Run(ctx); err != nil {
		log.Fatal(err)
	}
}
