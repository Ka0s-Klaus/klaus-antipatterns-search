package renderer

import (
	"fmt"
	"io"
	"sort"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/model"
)

// MarkdownRenderer outputs the scan result as a Markdown report
type MarkdownRenderer struct {
	w io.Writer
}

// NewMarkdownRenderer creates a new Markdown renderer
func NewMarkdownRenderer(w io.Writer) *MarkdownRenderer {
	return &MarkdownRenderer{w: w}
}

// Render outputs the scan result as a Markdown report
func (r *MarkdownRenderer) Render(result *model.ScanResult) error {
	fmt.Fprintf(r.w, "# 🔍 Antipatterns Report\n\n")

	// Summary section
	fmt.Fprintf(r.w, "## Summary\n\n")
	fmt.Fprintf(r.w, "- **Total Findings:** %d\n", result.Summary.Total)
	fmt.Fprintf(r.w, "- **Files Scanned:** %d\n", result.FilesScanned)
	fmt.Fprintf(r.w, "- **Duration:** %v\n", result.Duration)

	if result.Summary.Total > 0 {
		fmt.Fprintf(r.w, "- **Critical:** %d\n", result.Summary.Critical)
		fmt.Fprintf(r.w, "- **High:** %d\n", result.Summary.High)
		fmt.Fprintf(r.w, "- **Medium:** %d\n", result.Summary.Medium)
		fmt.Fprintf(r.w, "- **Low:** %d\n", result.Summary.Low)
	}

	fmt.Fprintf(r.w, "\n")

	// By rule breakdown
	if len(result.Summary.ByRule) > 0 {
		fmt.Fprintf(r.w, "## Breakdown by Rule\n\n")
		rules := make([]string, 0, len(result.Summary.ByRule))
		for rule := range result.Summary.ByRule {
			rules = append(rules, rule)
		}
		sort.Strings(rules)

		for _, rule := range rules {
			count := result.Summary.ByRule[rule]
			fmt.Fprintf(r.w, "- **%s:** %d\n", rule, count)
		}
		fmt.Fprintf(r.w, "\n")
	}

	// Findings table
	if result.Summary.Total > 0 {
		fmt.Fprintf(r.w, "## Findings\n\n")
		fmt.Fprintf(r.w, "| Severity | Rule | Path | Line:Col | Message |\n")
		fmt.Fprintf(r.w, "|----------|------|------|----------|----------|\n")

		sorted := make([]model.Finding, len(result.Findings))
		copy(sorted, result.Findings)
		sort.Slice(sorted, func(i, j int) bool {
			if sorted[i].Severity != sorted[j].Severity {
				return severityOrder(sorted[i].Severity) < severityOrder(sorted[j].Severity)
			}
			return sorted[i].Path < sorted[j].Path
		})

		for _, f := range sorted {
			fmt.Fprintf(r.w, "| %s | `%s` | `%s` | %d:%d | %s |\n",
				f.Severity, f.Rule, f.Path, f.Line, f.Column, f.Message)
		}
		fmt.Fprintf(r.w, "\n")
	}

	return nil
}
