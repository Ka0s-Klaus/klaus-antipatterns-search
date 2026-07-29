package model

import "strings"

// Severity ranks the impact of a detected anti-pattern.
// Values range from info (lowest impact) to critical (highest impact).
type Severity string

const (
	// SeverityInfo indicates informational findings (lowest priority).
	SeverityInfo Severity = "info"
	// SeverityLow indicates low-priority code quality issues.
	SeverityLow Severity = "low"
	// SeverityMedium indicates medium-priority technical debt or maintainability concerns.
	SeverityMedium Severity = "medium"
	// SeverityHigh indicates high-priority issues that may impact performance or reliability.
	SeverityHigh Severity = "high"
	// SeverityCritical indicates critical issues that pose security, stability, or correctness risks.
	SeverityCritical Severity = "critical"
)

// Location pinpoints where in the source the anti-pattern was detected.
type Location struct {
	// File is the relative path to the source file (e.g., "internal/scanner/scanner.go").
	File string `json:"file"`
	// Line is the 1-indexed line number where the finding occurs.
	Line int `json:"line"`
	// Column is the 0-indexed column offset; omitted if unknown (e.g., file-level findings).
	Column int `json:"column,omitempty"`
}

// Finding is the canonical output unit — one per detected anti-pattern instance.
// All detectors and adapters normalize their output to this struct.
type Finding struct {
	// Rule is the kebab-case identifier of the detected anti-pattern (e.g., "large_function", "god_object").
	Rule string `json:"rule"`
	// Message is a human-readable description of the finding (e.g., "function has 120 lines (threshold 80)").
	Message string `json:"message"`
	// Severity indicates the priority and impact of the finding.
	Severity Severity `json:"severity"`
	// Location specifies where in the source code the finding occurs.
	Location Location `json:"location"`
}

// SeverityFromString parses a severity string (case-insensitive), defaulting to SeverityMedium.
// Valid values: "info", "low", "medium", "high", "critical". Any other value returns SeverityMedium.
func SeverityFromString(s string) Severity {
	switch strings.ToLower(s) {
	case "info":
		return SeverityInfo
	case "low":
		return SeverityLow
	case "high":
		return SeverityHigh
	case "critical":
		return SeverityCritical
	default:
		return SeverityMedium
	}
}
