package adapter

import (
	"strings"
	"testing"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/config"
)

func TestParseRadonOutputEmpty(t *testing.T) {
	findings, err := parseRadonOutput([]byte(`{}`), config.Default())
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings for empty report, got %d", len(findings))
	}
}

func TestParseRadonOutputBelowThreshold(t *testing.T) {
	data := []byte(`{"path/to/file.py": [{"name":"f","type":"function","complexity":5,"lineno":1}]}`)
	cfg := config.Default() // Cyclomatic = 15
	findings, err := parseRadonOutput(data, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings below threshold, got %d", len(findings))
	}
}

func TestParseRadonOutputAboveThreshold(t *testing.T) {
	data := []byte(`{"path/to/file.py": [{"name":"heavy","type":"function","complexity":20,"lineno":10}]}`)
	cfg := config.Default() // Cyclomatic = 15
	findings, err := parseRadonOutput(data, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Rule != "cyclomatic_complexity" {
		t.Errorf("unexpected rule %q", findings[0].Rule)
	}
	if !strings.Contains(findings[0].Message, "heavy") {
		t.Errorf("expected function name in message, got %q", findings[0].Message)
	}
	if findings[0].Location.Line != 10 {
		t.Errorf("expected line 10, got %d", findings[0].Location.Line)
	}
}

func TestParseRadonOutputMultipleFiles(t *testing.T) {
	data := []byte(`{
		"a.py": [{"name":"f1","type":"function","complexity":20,"lineno":1}],
		"b.py": [{"name":"f2","type":"function","complexity":3,"lineno":5},
		         {"name":"f3","type":"method","complexity":18,"lineno":20}]
	}`)
	cfg := config.Default()
	findings, err := parseRadonOutput(data, cfg)
	if err != nil {
		t.Fatal(err)
	}
	// f1 (20>15) and f3 (18>15) should be flagged; f2 (3) should not
	if len(findings) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(findings))
	}
}

func TestParseRadonOutputInvalidJSON(t *testing.T) {
	_, err := parseRadonOutput([]byte("not json"), config.Default())
	if err == nil {
		t.Fatal("expected parse error for invalid JSON")
	}
}

func TestRadonSkipsWhenNotInstalled(t *testing.T) {
	findings, err := Radon(t.TempDir(), config.Default())
	if err != nil {
		t.Fatalf("Radon must not return error when tool is missing, got: %v", err)
	}
	_ = findings
}
