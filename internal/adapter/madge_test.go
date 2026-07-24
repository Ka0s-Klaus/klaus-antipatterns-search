package adapter

import (
	"testing"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/config"
)

func TestParseMadgeOutputNoCycles(t *testing.T) {
	data := []byte("✓ No circular dependency found!\n")
	findings, err := parseMadgeOutput(data, ".", config.Default())
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings, got %d", len(findings))
	}
}

func TestParseMadgeOutputOneCycle(t *testing.T) {
	data := []byte("✖ Found 1 circular dependency!\n\n1) a -> b -> a\n")
	findings, err := parseMadgeOutput(data, ".", config.Default())
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Rule != "circular_dependency" {
		t.Errorf("unexpected rule %q", findings[0].Rule)
	}
	expected := "circular dependency: a -> b -> a"
	if findings[0].Message != expected {
		t.Errorf("expected message %q, got %q", expected, findings[0].Message)
	}
}

func TestParseMadgeOutputMultipleCycles(t *testing.T) {
	data := []byte("✖ Found 2 circular dependencies!\n\n1) a -> b -> a\n2) c -> d -> e -> c\n")
	findings, err := parseMadgeOutput(data, ".", config.Default())
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(findings))
	}
}

func TestMadgeSkipsWhenNotInstalled(t *testing.T) {
	findings, err := Madge(t.TempDir(), config.Default())
	if err != nil {
		t.Fatalf("Madge must not return error when tool is missing, got: %v", err)
	}
	_ = findings
}
