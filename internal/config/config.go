package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

// ValidationError is returned when config validation fails.
type ValidationError struct {
	Field   string
	Message string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("invalid config field %q: %s", e.Field, e.Message)
}

const filename = ".antipatterns.yml"

type GodObjectConfig struct {
	Methods int `yaml:"methods"`
	LOC     int `yaml:"loc"`
}

type ThresholdsConfig struct {
	GodObject      GodObjectConfig `yaml:"god_object"`
	FunctionLOC    int             `yaml:"function_loc"`
	FunctionParams int             `yaml:"function_params"`
	Cyclomatic     int             `yaml:"cyclomatic"`
	Cognitive      int             `yaml:"cognitive"`
	DuplicationPct float64         `yaml:"duplication_pct"`
	MagicMinCount  int             `yaml:"magic_min_count"`
}

type SeveritiesConfig struct {
	GodObject    string `yaml:"god_object"`
	Cyclomatic   string `yaml:"cyclomatic"`
	Duplication  string `yaml:"duplication"`
	CircularDeps string `yaml:"circular_deps"`
	MagicNumbers string `yaml:"magic_numbers"`
	LargeFunction string `yaml:"large_function"`
	DeadCode     string `yaml:"dead_code"`
}

type OutputConfig struct {
	Formats []string `yaml:"formats"`
	Dir     string   `yaml:"dir"`
}

type OrgConfig struct {
	TokenEnv        string   `yaml:"token_env"`
	Output          string   `yaml:"output"`
	Publish         bool     `yaml:"publish"`
	ExcludeRepos    []string `yaml:"exclude_repos"`
	IncludeForks    bool     `yaml:"include_forks"`
	IncludeArchived bool     `yaml:"include_archived"`
}

type ExcludeConfig struct {
	Dirs  []string `yaml:"dirs"`
	Files []string `yaml:"files"`
}

type MagicNumbersConfig struct {
	Enabled bool     `yaml:"enabled"`
	Exclude []string `yaml:"exclude"`
}

type Config struct {
	Thresholds   ThresholdsConfig    `yaml:"thresholds"`
	MagicNumbers MagicNumbersConfig  `yaml:"magic_numbers"`
	Severities   SeveritiesConfig    `yaml:"severities"`
	Output       OutputConfig        `yaml:"output"`
	Orgs         map[string]OrgConfig `yaml:"orgs"`
	Exclude      ExcludeConfig       `yaml:"exclude"`
}

// Default returns a Config with the recommended out-of-the-box thresholds.
func Default() *Config {
	return &Config{
		MagicNumbers: MagicNumbersConfig{Enabled: true},
		Thresholds: ThresholdsConfig{
			GodObject:      GodObjectConfig{Methods: 20, LOC: 400},
			FunctionLOC:    80,
			FunctionParams: 7,
			Cyclomatic:     15,
			Cognitive:      20,
			DuplicationPct: 5,
			MagicMinCount:  3,
		},
		Severities: SeveritiesConfig{
			GodObject:     "high",
			Cyclomatic:    "medium",
			Duplication:   "medium",
			CircularDeps:  "high",
			MagicNumbers:  "low",
			LargeFunction: "low",
			DeadCode:      "low",
		},
		Output: OutputConfig{
			Formats: []string{"console"},
			Dir:     "reports/",
		},
		Exclude: ExcludeConfig{
			Dirs:  []string{"vendor/", "node_modules/", ".git/", "testdata/"},
			Files: []string{"**/*_test.go", "**/*.pb.go"},
		},
	}
}

// Load reads .antipatterns.yml from root, merging with defaults for missing fields.
func Load(root string) (*Config, error) {
	path := filepath.Join(root, filename)
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return Default(), nil
		}
		return nil, err
	}

	cfg := Default()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	// Validate config after loading to catch invalid user-provided values.
	if err := cfg.Validate(); err != nil {
		return nil, err
	}
	return cfg, nil
}

// Validate checks that the configuration has valid thresholds, severities, and ranges.
func (c *Config) Validate() error {
	// Validate thresholds are non-negative.
	if c.Thresholds.GodObject.Methods < 0 {
		return &ValidationError{Field: "thresholds.god_object.methods", Message: "must be non-negative"}
	}
	if c.Thresholds.GodObject.LOC < 0 {
		return &ValidationError{Field: "thresholds.god_object.loc", Message: "must be non-negative"}
	}
	if c.Thresholds.FunctionLOC < 0 {
		return &ValidationError{Field: "thresholds.function_loc", Message: "must be non-negative"}
	}
	if c.Thresholds.FunctionParams < 0 {
		return &ValidationError{Field: "thresholds.function_params", Message: "must be non-negative"}
	}
	if c.Thresholds.Cyclomatic < 0 {
		return &ValidationError{Field: "thresholds.cyclomatic", Message: "must be non-negative"}
	}
	if c.Thresholds.Cognitive < 0 {
		return &ValidationError{Field: "thresholds.cognitive", Message: "must be non-negative"}
	}
	if c.Thresholds.DuplicationPct < 0 || c.Thresholds.DuplicationPct > 100 {
		return &ValidationError{Field: "thresholds.duplication_pct", Message: "must be between 0 and 100"}
	}
	if c.Thresholds.MagicMinCount < 0 {
		return &ValidationError{Field: "thresholds.magic_min_count", Message: "must be non-negative"}
	}

	// Validate severities are valid.
	validSeverities := map[string]bool{"info": true, "low": true, "medium": true, "high": true, "critical": true}
	if c.Severities.GodObject != "" && !validSeverities[c.Severities.GodObject] {
		return &ValidationError{Field: "severities.god_object", Message: "invalid severity: " + c.Severities.GodObject}
	}
	if c.Severities.Cyclomatic != "" && !validSeverities[c.Severities.Cyclomatic] {
		return &ValidationError{Field: "severities.cyclomatic", Message: "invalid severity: " + c.Severities.Cyclomatic}
	}
	if c.Severities.Duplication != "" && !validSeverities[c.Severities.Duplication] {
		return &ValidationError{Field: "severities.duplication", Message: "invalid severity: " + c.Severities.Duplication}
	}
	if c.Severities.CircularDeps != "" && !validSeverities[c.Severities.CircularDeps] {
		return &ValidationError{Field: "severities.circular_deps", Message: "invalid severity: " + c.Severities.CircularDeps}
	}
	if c.Severities.MagicNumbers != "" && !validSeverities[c.Severities.MagicNumbers] {
		return &ValidationError{Field: "severities.magic_numbers", Message: "invalid severity: " + c.Severities.MagicNumbers}
	}
	if c.Severities.LargeFunction != "" && !validSeverities[c.Severities.LargeFunction] {
		return &ValidationError{Field: "severities.large_function", Message: "invalid severity: " + c.Severities.LargeFunction}
	}
	if c.Severities.DeadCode != "" && !validSeverities[c.Severities.DeadCode] {
		return &ValidationError{Field: "severities.dead_code", Message: "invalid severity: " + c.Severities.DeadCode}
	}

	return nil
}

// IsExcludedDir reports whether a directory name matches the exclusion list.
func (c *Config) IsExcludedDir(name string) bool {
	for _, d := range c.Exclude.Dirs {
		if d == name || d == name+"/" {
			return true
		}
	}
	return false
}
