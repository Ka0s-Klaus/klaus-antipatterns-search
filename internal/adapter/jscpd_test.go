package adapter

import (
	"testing"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/config"
)

func TestParseJscpdOutputNoDuplicates(t *testing.T) {
	data := []byte(`{"statistics":{"percentage":0},"duplicates":[]}`)
	cfg := config.Default()
	findings, err := parseJscpdOutput(data, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings, got %d", len(findings))
	}
}

func TestParseJscpdOutputBelowThreshold(t *testing.T) {
	// 3% duplication, threshold 5% → no summary finding
	data := []byte(`{"statistics":{"percentage":3.0},"duplicates":[]}`)
	cfg := config.Default() // DuplicationPct = 5
	findings, err := parseJscpdOutput(data, cfg)
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings below threshold, got %d", len(findings))
	}
}

func TestParseJscpdOutputAboveThreshold(t *testing.T) {
	// 12.5% duplication + 1 clone pair
	data := []byte(`{
		"statistics": {"percentage": 12.5},
		"duplicates": [{
			"lines": 10,
			"firstFile":  {"name": "a.go", "start": 1},
			"secondFile": {"name": "b.go", "start": 5}
		}]
	}`)
	cfg := config.Default() // DuplicationPct = 5
	findings, err := parseJscpdOutput(data, cfg)
	if err != nil {
		t.Fatal(err)
	}
	// 1 summary + 1 clone finding
	if len(findings) != 2 {
		t.Fatalf("expected 2 findings, got %d", len(findings))
	}
	if findings[0].Rule != "duplication" {
		t.Errorf("unexpected rule %q", findings[0].Rule)
	}
	if findings[1].Location.File != "a.go" {
		t.Errorf("expected first file a.go, got %q", findings[1].Location.File)
	}
}

func TestParseJscpdOutputInvalidJSON(t *testing.T) {
	_, err := parseJscpdOutput([]byte("not json"), config.Default())
	if err == nil {
		t.Fatal("expected parse error for invalid JSON")
	}
}

func TestJscpdSkipsWhenNotInstalled(t *testing.T) {
	// If jscpd is absent: must return nil findings + ErrToolNotFound.
	// If jscpd is present in CI: err must be nil. Either case is valid.
	findings, err := Jscpd(t.TempDir(), config.Default())
	if err != nil && err != ErrToolNotFound {
		t.Fatalf("Jscpd must return nil or ErrToolNotFound, got: %v", err)
	}
	_ = findings
}
