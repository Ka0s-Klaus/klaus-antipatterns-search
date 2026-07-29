package model

import (
	"fmt"
	"time"
)

// Severity level for findings
type Severity string

const (
	SeverityInfo     Severity = "info"
	SeverityLow      Severity = "low"
	SeverityMedium   Severity = "medium"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"
)

// Finding represents a detected antipattern
type Finding struct {
	Path      string    `json:"path"`
	Line      int       `json:"line"`
	Column    int       `json:"column"`
	Rule      string    `json:"rule"`
	Message   string    `json:"message"`
	Severity  Severity  `json:"severity"`
	Snippet   string    `json:"snippet,omitempty"`
	Language  string    `json:"language,omitempty"`
	StartLine int       `json:"start_line,omitempty"`
	EndLine   int       `json:"end_line,omitempty"`
	Context   string    `json:"context,omitempty"`
	Detector  string    `json:"detector"`
	Timestamp time.Time `json:"timestamp"`
}

// ScanResult aggregates all findings from a scan
type ScanResult struct {
	Findings       []Finding              `json:"findings"`
	Summary        Summary                `json:"summary"`
	Duration       time.Duration          `json:"duration"`
	RulesSeverity  map[string]Severity    `json:"rules_severity,omitempty"`
	FilesScanned   int                    `json:"files_scanned"`
	Errors         []string               `json:"errors,omitempty"`
}

// Summary provides aggregated metrics
type Summary struct {
	Total      int `json:"total"`
	Critical   int `json:"critical"`
	High       int `json:"high"`
	Medium     int `json:"medium"`
	Low        int `json:"low"`
	Info       int `json:"info"`
	ByRule     map[string]int `json:"by_rule"`
	BySeverity map[Severity]int `json:"by_severity"`
}

// String returns a string representation of a severity
func (s Severity) String() string {
	return string(s)
}

// Color returns ANSI color code for the severity (for terminal output)
func (s Severity) Color() string {
	switch s {
	case SeverityCritical:
		return "\033[91m" // Bright red
	case SeverityHigh:
		return "\033[31m" // Red
	case SeverityMedium:
		return "\033[93m" // Bright yellow
	case SeverityLow:
		return "\033[94m" // Bright blue
	default:
		return "\033[0m" // Reset
	}
}

// ColorReset returns the ANSI reset code
func ColorReset() string {
	return "\033[0m"
}

// NewFinding creates a new Finding with a timestamp
func NewFinding(path, rule, message string, severity Severity) *Finding {
	return &Finding{
		Path:      path,
		Rule:      rule,
		Message:   message,
		Severity:  severity,
		Timestamp: time.Now(),
		Line:      1,
		Column:    0,
	}
}

// CalculateSummary generates a summary from findings
func CalculateSummary(findings []Finding) Summary {
	summary := Summary{
		Total:      len(findings),
		ByRule:     make(map[string]int),
		BySeverity: make(map[Severity]int),
	}

	for _, f := range findings {
		summary.ByRule[f.Rule]++
		summary.BySeverity[f.Severity]++

		switch f.Severity {
		case SeverityCritical:
			summary.Critical++
		case SeverityHigh:
			summary.High++
		case SeverityMedium:
			summary.Medium++
		case SeverityLow:
			summary.Low++
		case SeverityInfo:
			summary.Info++
		}
	}

	return summary
}

// String returns a concise string representation of a Finding
func (f *Finding) String() string {
	return fmt.Sprintf("%s:%d:%d — %s: %s", f.Path, f.Line, f.Column, f.Rule, f.Message)
}
