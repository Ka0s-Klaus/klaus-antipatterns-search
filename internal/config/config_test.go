package config

import (
	"os"
	"testing"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/model"
)

func TestDefault(t *testing.T) {
	cfg := Default()
	if cfg == nil {
		t.Fatal("expected default config, got nil")
	}
	if cfg.Thresholds.GodObject.Methods != 20 {
		t.Errorf("expected god_object methods threshold 20, got %d", cfg.Thresholds.GodObject.Methods)
	}
	if cfg.Thresholds.FunctionLOC != 80 {
		t.Errorf("expected function_loc threshold 80, got %d", cfg.Thresholds.FunctionLOC)
	}
}

func TestGetSeverity(t *testing.T) {
	cfg := Default()
	tests := []struct {
		rule     string
		expected model.Severity
	}{
		{"god_object", model.SeverityHigh},
		{"magic_numbers", model.SeverityLow},
		{"unknown_rule", model.SeverityMedium}, // Should return default
	}

	for _, tt := range tests {
		sev := cfg.GetSeverity(tt.rule)
		if sev != tt.expected {
			t.Errorf("rule %s: expected %s, got %s", tt.rule, tt.expected, sev)
		}
	}
}

func TestLoadNonExistent(t *testing.T) {
	cfg, err := Load("/nonexistent/path/.antipatterns.yml")
	if err != nil {
		t.Fatalf("expected no error for non-existent file, got %v", err)
	}
	if cfg == nil {
		t.Fatal("expected default config returned")
	}
}

func TestLoadValidFile(t *testing.T) {
	// Create a temp file
	tmpfile, err := os.CreateTemp("", "antipatterns-*.yml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	defer os.Remove(tmpfile.Name())
	tmpfile.Close()

	cfg, err := Load(tmpfile.Name())
	if err != nil {
		t.Fatalf("failed to load config: %v", err)
	}
	// For now, returns defaults even with file present
	// YAML parsing will be added in Phase 1
	if cfg.Thresholds.GodObject.Methods != 20 {
		t.Errorf("expected god_object methods 20 (default), got %d", cfg.Thresholds.GodObject.Methods)
	}
}
