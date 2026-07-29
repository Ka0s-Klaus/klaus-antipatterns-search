package model

import (
	"testing"
)

func TestNewFinding(t *testing.T) {
	f := NewFinding("test.go", "god_object", "class too large", SeverityHigh)
	if f.Path != "test.go" {
		t.Errorf("expected path 'test.go', got %s", f.Path)
	}
	if f.Rule != "god_object" {
		t.Errorf("expected rule 'god_object', got %s", f.Rule)
	}
	if f.Severity != SeverityHigh {
		t.Errorf("expected severity high, got %s", f.Severity)
	}
}

func TestCalculateSummary(t *testing.T) {
	findings := []Finding{
		{Rule: "god_object", Severity: SeverityCritical},
		{Rule: "god_object", Severity: SeverityHigh},
		{Rule: "magic_numbers", Severity: SeverityLow},
		{Rule: "magic_numbers", Severity: SeverityLow},
	}

	summary := CalculateSummary(findings)
	if summary.Total != 4 {
		t.Errorf("expected total 4, got %d", summary.Total)
	}
	if summary.Critical != 1 {
		t.Errorf("expected 1 critical, got %d", summary.Critical)
	}
	if summary.High != 1 {
		t.Errorf("expected 1 high, got %d", summary.High)
	}
	if summary.Low != 2 {
		t.Errorf("expected 2 low, got %d", summary.Low)
	}
	if summary.ByRule["god_object"] != 2 {
		t.Errorf("expected 2 god_object findings, got %d", summary.ByRule["god_object"])
	}
	if summary.ByRule["magic_numbers"] != 2 {
		t.Errorf("expected 2 magic_numbers findings, got %d", summary.ByRule["magic_numbers"])
	}
}

func TestSeverityColor(t *testing.T) {
	tests := []struct {
		sev    Severity
		hasANSI bool
	}{
		{SeverityCritical, true},
		{SeverityHigh, true},
		{SeverityMedium, true},
		{SeverityLow, true},
		{SeverityInfo, true},
	}

	for _, tt := range tests {
		color := tt.sev.Color()
		if tt.hasANSI && color == "" {
			t.Errorf("expected ANSI code for %s, got empty", tt.sev)
		}
	}
}
