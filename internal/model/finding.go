package model

import "strings"

// Severity ranks the impact of a detected anti-pattern.
type Severity string

const (
	SeverityInfo     Severity = "info"
	SeverityLow      Severity = "low"
	SeverityMedium   Severity = "medium"
	SeverityHigh     Severity = "high"
	SeverityCritical Severity = "critical"
)

// Location pinpoints where in the source the anti-pattern was detected.
type Location struct {
	File   string `json:"file"`
	Line   int    `json:"line"`
	Column int    `json:"column,omitempty"`
}

// Finding is the canonical output unit — one per detected anti-pattern instance.
type Finding struct {
	Rule     string   `json:"rule"`
	Message  string   `json:"message"`
	Severity Severity `json:"severity"`
	Location Location `json:"location"`
}

// SeverityFromString parses a severity string (case-insensitive), defaulting to SeverityMedium.
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
