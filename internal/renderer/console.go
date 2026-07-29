package renderer

import (
	"fmt"
	"io"
	"sort"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/model"
)

// ConsoleRenderer outputs a colored, formatted table to console
type ConsoleRenderer struct {
	w io.Writer
}

// NewConsoleRenderer creates a new console renderer
func NewConsoleRenderer(w io.Writer) *ConsoleRenderer {
	return &ConsoleRenderer{w: w}
}

// Render outputs the scan result as a formatted table
func (r *ConsoleRenderer) Render(result *model.ScanResult) error {
	if result.Summary.Total == 0 {
		fmt.Fprintf(r.w, "✓ No antipatterns detected in %d files\n", result.FilesScanned)
		return nil
	}

	// Header
	fmt.Fprintf(r.w, "\n%s%s%s\n", model.ColorReset(), "ANTIPATTERNS DETECTED", model.ColorReset())
	fmt.Fprintf(r.w, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")

	// Sort findings by severity then path
	sorted := make([]model.Finding, len(result.Findings))
	copy(sorted, result.Findings)
	sort.Slice(sorted, func(i, j int) bool {
		if sorted[i].Severity != sorted[j].Severity {
			return severityOrder(sorted[i].Severity) < severityOrder(sorted[j].Severity)
		}
		return sorted[i].Path < sorted[j].Path
	})

	// Print each finding
	for _, f := range sorted {
		severity := f.Severity.Color() + f.Severity.String() + model.ColorReset()
		fmt.Fprintf(r.w, "%s %-9s %s:%d:%d  %s\n", severity, "["+f.Rule+"]", f.Path, f.Line, f.Column, f.Message)
	}

	// Summary
	fmt.Fprintf(r.w, "━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━\n")
	fmt.Fprintf(r.w, "Total: %d | ", result.Summary.Total)
	if result.Summary.Critical > 0 {
		fmt.Fprintf(r.w, "%sCritical: %d%s | ", model.SeverityCritical.Color(), result.Summary.Critical, model.ColorReset())
	}
	if result.Summary.High > 0 {
		fmt.Fprintf(r.w, "%sHigh: %d%s | ", model.SeverityHigh.Color(), result.Summary.High, model.ColorReset())
	}
	if result.Summary.Medium > 0 {
		fmt.Fprintf(r.w, "%sMedium: %d%s | ", model.SeverityMedium.Color(), result.Summary.Medium, model.ColorReset())
	}
	if result.Summary.Low > 0 {
		fmt.Fprintf(r.w, "%sLow: %d%s | ", model.SeverityLow.Color(), result.Summary.Low, model.ColorReset())
	}
	fmt.Fprintf(r.w, "Duration: %v\n\n", result.Duration)

	return nil
}

func severityOrder(s model.Severity) int {
	switch s {
	case model.SeverityCritical:
		return 0
	case model.SeverityHigh:
		return 1
	case model.SeverityMedium:
		return 2
	case model.SeverityLow:
		return 3
	case model.SeverityInfo:
		return 4
	default:
		return 5
	}
}
