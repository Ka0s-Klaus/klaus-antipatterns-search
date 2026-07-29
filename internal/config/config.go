package config

import (
	"fmt"
	"os"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/model"
)

// Config holds the complete antipatterns configuration
type Config struct {
	Thresholds Thresholds         `yaml:"thresholds"`
	Severities map[string]string  `yaml:"severities"`
	Output     OutputConfig       `yaml:"output"`
	Exclude    ExcludeConfig      `yaml:"exclude"`
	Orgs       map[string]OrgConfig `yaml:"orgs,omitempty"`
}

// GodObjectThreshold defines thresholds for God Object detection
type GodObjectThreshold struct {
	Methods int
	LOC     int
}

// Thresholds defines detection thresholds for various rules
type Thresholds struct {
	GodObject      GodObjectThreshold
	FunctionLOC    int
	FunctionParams int
	Cyclomatic     int
	Cognitive      int
	DuplicationPct float64
	MagicMinCount  int
}

// OutputConfig defines output format and directory
type OutputConfig struct {
	Formats []string `yaml:"formats"`
	Dir     string   `yaml:"dir"`
}

// ExcludeConfig defines exclusion patterns
type ExcludeConfig struct {
	Dirs  []string `yaml:"dirs"`
	Files []string `yaml:"files"`
}

// OrgConfig holds per-organization settings for multi-org scans
type OrgConfig struct {
	TokenEnv string `yaml:"token_env"`
	Output   string `yaml:"output"`
	Publish  bool   `yaml:"publish"`
}

// Default returns a config with sensible defaults
func Default() *Config {
	return &Config{
		Thresholds: Thresholds{
			GodObject: GodObjectThreshold{Methods: 20, LOC: 400},
			FunctionLOC:    80,
			FunctionParams: 7,
			Cyclomatic:     15,
			Cognitive:      20,
			DuplicationPct: 5,
			MagicMinCount:  3,
		},
		Severities: map[string]string{
			"god_object":    "high",
			"cyclomatic":    "medium",
			"duplication":   "medium",
			"circular_deps": "high",
			"magic_numbers": "low",
			"large_function": "low",
			"dead_code":     "low",
		},
		Output: OutputConfig{
			Formats: []string{"console", "json"},
			Dir:     "reports/",
		},
		Exclude: ExcludeConfig{
			Dirs: []string{"vendor/", "node_modules/", ".git/", "testdata/"},
			Files: []string{"**/*_test.go", "**/*.pb.go"},
		},
	}
}

// Load reads a config file and merges with defaults
// For Phase 0, this is a placeholder that returns defaults
// Phase 1 will add YAML parsing
func Load(path string) (*Config, error) {
	cfg := Default()

	if path == "" {
		return cfg, nil
	}

	// For now, just verify the file exists if specified
	if _, err := os.Stat(path); err != nil {
		if !os.IsNotExist(err) {
			return nil, fmt.Errorf("failed to access config: %w", err)
		}
		// File doesn't exist, return defaults
	}

	// TODO: Phase 1 — Add YAML parsing with gopkg.in/yaml.v3

	return cfg, nil
}

// GetSeverity returns the severity level for a rule
func (c *Config) GetSeverity(rule string) model.Severity {
	if sevStr, ok := c.Severities[rule]; ok {
		switch sevStr {
		case "critical":
			return model.SeverityCritical
		case "high":
			return model.SeverityHigh
		case "medium":
			return model.SeverityMedium
		case "low":
			return model.SeverityLow
		case "info":
			return model.SeverityInfo
		}
	}
	return model.SeverityMedium // Default fallback
}
