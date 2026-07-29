package config

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDefault(t *testing.T) {
	cfg := Default()

	if cfg.Thresholds.GodObject.Methods != 20 {
		t.Errorf("god_object.methods: want 20, got %d", cfg.Thresholds.GodObject.Methods)
	}
	if cfg.Thresholds.FunctionLOC != 80 {
		t.Errorf("function_loc: want 80, got %d", cfg.Thresholds.FunctionLOC)
	}
	if len(cfg.Exclude.Dirs) == 0 {
		t.Error("exclude.dirs should not be empty in defaults")
	}
}

func TestLoadMissingFile(t *testing.T) {
	cfg, err := Load(t.TempDir())
	if err != nil {
		t.Fatalf("missing config file should return defaults, got error: %v", err)
	}
	if cfg.Thresholds.Cyclomatic != 15 {
		t.Errorf("cyclomatic default: want 15, got %d", cfg.Thresholds.Cyclomatic)
	}
}

func TestLoadFromFile(t *testing.T) {
	dir := t.TempDir()
	content := `
thresholds:
  function_loc: 50
  cyclomatic: 10
`
	if err := os.WriteFile(filepath.Join(dir, filename), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	cfg, err := Load(dir)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.Thresholds.FunctionLOC != 50 {
		t.Errorf("function_loc: want 50, got %d", cfg.Thresholds.FunctionLOC)
	}
	if cfg.Thresholds.Cyclomatic != 10 {
		t.Errorf("cyclomatic: want 10, got %d", cfg.Thresholds.Cyclomatic)
	}
	// Fields not in file should keep defaults
	if cfg.Thresholds.GodObject.Methods != 20 {
		t.Errorf("god_object.methods should keep default 20, got %d", cfg.Thresholds.GodObject.Methods)
	}
}

func TestIsExcludedDir(t *testing.T) {
	cfg := Default()
	cases := []struct {
		name string
		want bool
	}{
		{"vendor", true},
		{"node_modules", true},
		{".git", true},
		{"src", false},
		{"internal", false},
	}
	for _, tc := range cases {
		if got := cfg.IsExcludedDir(tc.name); got != tc.want {
			t.Errorf("IsExcludedDir(%q): want %v, got %v", tc.name, tc.want, got)
		}
	}
}

func TestValidateNegativeThresholds(t *testing.T) {
	cases := []struct {
		name  string
		setup func(*Config)
	}{
		{
			"negative god_object.methods",
			func(c *Config) { c.Thresholds.GodObject.Methods = -1 },
		},
		{
			"negative god_object.loc",
			func(c *Config) { c.Thresholds.GodObject.LOC = -1 },
		},
		{
			"negative function_loc",
			func(c *Config) { c.Thresholds.FunctionLOC = -1 },
		},
		{
			"negative function_params",
			func(c *Config) { c.Thresholds.FunctionParams = -1 },
		},
		{
			"negative cyclomatic",
			func(c *Config) { c.Thresholds.Cyclomatic = -1 },
		},
		{
			"negative cognitive",
			func(c *Config) { c.Thresholds.Cognitive = -1 },
		},
		{
			"negative magic_min_count",
			func(c *Config) { c.Thresholds.MagicMinCount = -1 },
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := Default()
			tc.setup(cfg)
			if err := cfg.Validate(); err == nil {
				t.Error("expected validation error, got nil")
			}
		})
	}
}

func TestValidateDuplicationPct(t *testing.T) {
	cases := []struct {
		value float64
		valid bool
	}{
		{-1, false},
		{0, true},
		{5, true},
		{100, true},
		{101, false},
	}

	for _, tc := range cases {
		cfg := Default()
		cfg.Thresholds.DuplicationPct = tc.value
		err := cfg.Validate()
		if (err == nil) != tc.valid {
			t.Errorf("duplication_pct=%v: want valid=%v, got error=%v", tc.value, tc.valid, err)
		}
	}
}

func TestValidateSeverities(t *testing.T) {
	cases := []struct {
		name     string
		severity string
		valid    bool
	}{
		{"info", "info", true},
		{"low", "low", true},
		{"medium", "medium", true},
		{"high", "high", true},
		{"critical", "critical", true},
		{"invalid", "invalid_severity", false},
		{"empty", "", true}, // empty is valid (uses default)
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg := Default()
			cfg.Severities.GodObject = tc.severity
			err := cfg.Validate()
			if (err == nil) != tc.valid {
				t.Errorf("severity=%q: want valid=%v, got error=%v", tc.severity, tc.valid, err)
			}
		})
	}
}

func TestLoadValidatesConfig(t *testing.T) {
	dir := t.TempDir()
	content := `
thresholds:
  cyclomatic: -1
`
	if err := os.WriteFile(filepath.Join(dir, filename), []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}

	_, err := Load(dir)
	if err == nil {
		t.Error("expected validation error when loading config with negative threshold, got nil")
	}
}
