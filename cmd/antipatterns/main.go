package main

import (
	"fmt"
	"io"
	"os"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/config"
	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/orgscanner"
	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/renderer"
	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/scanner"
	"github.com/spf13/cobra"
)

// version is set at build time via -ldflags.
var version = "dev"

func main() {
	if err := rootCmd().Execute(); err != nil {
		os.Exit(1)
	}
}

func rootCmd() *cobra.Command {
	root := &cobra.Command{
		Use:   "antipatterns",
		Short: "Multi-language anti-pattern detector",
		Long:  "Detects anti-patterns in source code across multiple languages and organisations.",
	}
	root.AddCommand(scanCmd(), scanOrgCmd(), versionCmd())
	return root
}

func scanCmd() *cobra.Command {
	var format string
	var output string

	cmd := &cobra.Command{
		Use:   "scan [path]",
		Short: "Scan a directory for anti-patterns",
		Args:  cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			root := "."
			if len(args) > 0 {
				root = args[0]
			}

			cfg, err := config.Load(root)
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			s := scanner.New(cfg)
			findings, err := s.Run(root)
			if err != nil {
				return fmt.Errorf("scan: %w", err)
			}

			w, closeW, err := openOutput(output)
			if err != nil {
				return err
			}
			defer closeW()

			switch format {
			case "json":
				return renderer.NewJSON(w).Render(findings)
			case "sarif":
				return renderer.NewSARIF(w, version, root).Render(findings)
			default:
				return renderer.NewConsole(w).Render(findings)
			}
		},
	}

	cmd.Flags().StringVar(&format, "format", "console", "output format: console, json, sarif")
	cmd.Flags().StringVarP(&output, "output", "o", "", "write output to file (default: stdout)")
	return cmd
}

func scanOrgCmd() *cobra.Command {
	var format string
	var output string
	var workers int

	cmd := &cobra.Command{
		Use:   "scan-org <profile>",
		Short: "Scan all repositories of a GitHub org",
		Long: `Scans all repositories of a GitHub organisation using an org profile
defined in .antipatterns.yml. Requires gh CLI authenticated with access
to the target organisation.`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			profileName := args[0]

			cfg, err := config.Load(".")
			if err != nil {
				return fmt.Errorf("loading config: %w", err)
			}

			profile, ok := cfg.Orgs[profileName]
			if !ok {
				return fmt.Errorf("org profile %q not found in .antipatterns.yml — define it under orgs:", profileName)
			}

			s := orgscanner.New(workers)
			report, err := s.Run(profileName, profile, cfg)
			if err != nil {
				return fmt.Errorf("scan-org: %w", err)
			}

			w, closeW, err := openOutput(output)
			if err != nil {
				return err
			}
			defer closeW()

			switch format {
			case "json":
				return renderer.NewOrgJSON(w).Render(report)
			default:
				return renderer.NewOrgConsole(w).Render(report)
			}
		},
	}

	cmd.Flags().StringVar(&format, "format", "console", "output format: console, json")
	cmd.Flags().StringVarP(&output, "output", "o", "", "write output to file (default: stdout)")
	cmd.Flags().IntVar(&workers, "workers", 4, "number of parallel repo scanners")
	return cmd
}

// openOutput returns a writer and a close function.
// If path is empty, returns stdout with a no-op closer.
func openOutput(path string) (io.Writer, func(), error) {
	if path == "" {
		return os.Stdout, func() {}, nil
	}
	f, err := os.Create(path)
	if err != nil {
		return nil, nil, fmt.Errorf("opening output file: %w", err)
	}
	return f, func() { f.Close() }, nil
}

func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version",
		Run: func(cmd *cobra.Command, _ []string) {
			fmt.Fprintf(cmd.OutOrStdout(), "antipatterns %s\n", version)
		},
	}
}
