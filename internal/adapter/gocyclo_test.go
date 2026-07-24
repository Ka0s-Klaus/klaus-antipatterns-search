package adapter

import (
	"strings"
	"testing"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/config"
)

func TestParseGocycloOutputEmpty(t *testing.T) {
	findings, err := parseGocycloOutput([]byte(""), config.Default())
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings for empty output, got %d", len(findings))
	}
}

func TestParseGocycloOutputSingleLine(t *testing.T) {
	// format: "<complexity> <package> <func> <file>:<line>:<col>"
	data := []byte("16 main bigFunc main.go:10:1\n")
	cfg := config.Default()
	findings, err := parseGocycloOutput(data, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Rule != "cyclomatic_complexity" {
		t.Errorf("unexpected rule %q", findings[0].Rule)
	}
	if !strings.Contains(findings[0].Message, "bigFunc") {
		t.Errorf("expected function name in message, got %q", findings[0].Message)
	}
	if findings[0].Location.File != "main.go" {
		t.Errorf("expected file main.go, got %q", findings[0].Location.File)
	}
	if findings[0].Location.Line != 10 {
		t.Errorf("expected line 10, got %d", findings[0].Location.Line)
	}
}

func TestParseGocycloOutputMultipleLines(t *testing.T) {
	data := []byte("20 pkg alpha pkg/a.go:5:1\n18 pkg beta pkg/b.go:30:1\n")
	findings, err := parseGocycloOutput(data, config.Default())
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(findings))
	}
}

func TestParseGocycloOutputMalformedLine(t *testing.T) {
	// Malformed lines must be silently skipped, not crash.
	data := []byte("not-a-number pkg func file.go:1:1\nthis is garbage\n")
	findings, err := parseGocycloOutput(data, config.Default())
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings from malformed input, got %d", len(findings))
	}
}

func TestSplitLocation(t *testing.T) {
	file, line := splitLocation("internal/pkg/file.go:42:1")
	if file != "internal/pkg/file.go" {
		t.Errorf("expected file, got %q", file)
	}
	if line != 42 {
		t.Errorf("expected line 42, got %d", line)
	}
}

func TestGocycloSkipsWhenNotInstalled(t *testing.T) {
	findings, err := Gocyclo(t.TempDir(), config.Default())
	if findings != nil {
		t.Fatalf("Gocyclo must return nil findings when tool is missing")
	}
	if err != ErrToolNotFound {
		t.Fatalf("Gocyclo must return ErrToolNotFound when tool is missing, got: %v", err)
	}
}
