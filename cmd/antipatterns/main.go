package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/config"
	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/core"
	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/renderer"
)

var (
	Version = "dev"
	GitHash = "unknown"
)

func main() {
	fs := flag.NewFlagSet("antipatterns", flag.ExitOnError)
	configPath := fs.String("config", ".antipatterns.yml", "Path to config file")
	outputFormat := fs.String("format", "console", "Output format: console, json, markdown")
	outputFile := fs.String("output", "", "Output file (default: stdout)")
	showVersion := fs.Bool("version", false, "Show version")

	fs.Usage = func() {
		fmt.Fprintf(os.Stderr, `Klaus-antipatterns-search — Multi-language antipattern detector

Usage:
  antipatterns [flags] <path>
  antipatterns [flags] scan <path>
  antipatterns version

Flags:
`)
		fs.PrintDefaults()
	}

	if len(os.Args) < 2 {
		fs.Usage()
		os.Exit(1)
	}

	// Handle version subcommand
	if os.Args[1] == "version" {
		fmt.Printf("antipatterns version %s (%s)\n", Version, GitHash)
		os.Exit(0)
	}

	// Handle scan subcommand or path directly
	args := os.Args[1:]
	if args[0] == "scan" {
		args = args[1:]
	}

	if err := fs.Parse(args); err != nil {
		log.Fatalf("parse error: %v", err)
	}

	if *showVersion {
		fmt.Printf("antipatterns version %s (%s)\n", Version, GitHash)
		os.Exit(0)
	}

	if fs.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "Error: path argument required\n")
		fs.Usage()
		os.Exit(1)
	}

	scanPath := fs.Arg(0)

	// Load configuration
	cfg, err := config.Load(*configPath)
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// Create scanner
	scanner := core.NewScanner(cfg)

	// Perform scan
	result, err := scanner.Scan(scanPath)
	if err != nil {
		log.Fatalf("scan failed: %v", err)
	}

	// Select renderer
	var r renderer.Renderer
	switch *outputFormat {
	case "console":
		r = renderer.NewConsoleRenderer(os.Stdout)
	case "json":
		r = renderer.NewJSONRenderer(os.Stdout)
	case "markdown":
		r = renderer.NewMarkdownRenderer(os.Stdout)
	default:
		log.Fatalf("unknown format: %s", *outputFormat)
	}

	// Handle output file redirection
	var out *os.File = os.Stdout
	if *outputFile != "" {
		var err error
		out, err = os.Create(*outputFile)
		if err != nil {
			log.Fatalf("failed to create output file: %v", err)
		}
		defer out.Close()

		// Update renderer with file output
		switch *outputFormat {
		case "console":
			r = renderer.NewConsoleRenderer(out)
		case "json":
			r = renderer.NewJSONRenderer(out)
		case "markdown":
			r = renderer.NewMarkdownRenderer(out)
		}
	}

	// Render output
	if err := r.Render(result); err != nil {
		log.Fatalf("render failed: %v", err)
	}
}
