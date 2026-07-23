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
