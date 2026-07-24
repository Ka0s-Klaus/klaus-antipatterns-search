package config

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

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
	return cfg, nil
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
