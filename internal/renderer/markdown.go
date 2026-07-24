package renderer

import (
	"fmt"
	"io"
	"sort"
	"time"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/model"
)

var severityEmoji = map[model.Severity]string{
	model.SeverityInfo:     "🔵",
	model.SeverityLow:      "🟢",
	model.SeverityMedium:   "🟡",
	model.SeverityHigh:     "🟠",
	model.SeverityCritical: "🔴",
}

// MarkdownRenderer writes findings as an actionable Markdown report.
// Each finding becomes a checkbox "- [ ]" ready to use as a refactoring backlog.
type MarkdownRenderer struct {
	w       io.Writer
	version string
	root    string
}

func NewMarkdown(w io.Writer, version, root string) *MarkdownRenderer {
	return &MarkdownRenderer{w: w, version: version, root: root}
}

func (r *MarkdownRenderer) Render(findings []model.Finding) error {
	date := time.Now().Format("2006-01-02")

	fmt.Fprintf(r.w, "# 🔍 Anti-pattern Report\n\n")
	fmt.Fprintf(r.w, "> **Path:** `%s`  \n", r.root)
	fmt.Fprintf(r.w, "> **Date:** %s  \n", date)
	fmt.Fprintf(r.w, "> **Tool:** antipatterns %s\n\n", r.version)

	if len(findings) == 0 {
		fmt.Fprintf(r.w, "✅ **No anti-patterns found.**\n")
		return nil
	}

	ruleOrder, ruleFindings := groupByRule(findings)

	fmt.Fprintf(r.w, "## 📊 Resumen\n\n")
	fmt.Fprintf(r.w, "**Total: %d finding(s)**\n\n", len(findings))
	fmt.Fprintf(r.w, "| Regla | Severidad | Findings |\n")
	fmt.Fprintf(r.w, "|---|---|---|\n")
	for _, rule := range ruleOrder {
		ff := ruleFindings[rule]
		sev := ff[0].Severity
		fmt.Fprintf(r.w, "| `%s` | %s %s | %d |\n", rule, severityEmoji[sev], sev, len(ff))
	}
	fmt.Fprintf(r.w, "\n---\n\n")

	for _, rule := range ruleOrder {
		ff := ruleFindings[rule]
		sev := ff[0].Severity
		fmt.Fprintf(r.w, "## %s `%s` — %d finding(s)\n\n", severityEmoji[sev], rule, len(ff))
		for _, f := range ff {
			loc := toRelativeURI(f.Location.File, r.root)
			if f.Location.Line > 0 {
				fmt.Fprintf(r.w, "- [ ] **`%s:%d`** — %s\n", loc, f.Location.Line, f.Message)
			} else {
				fmt.Fprintf(r.w, "- [ ] **`%s`** — %s\n", loc, f.Message)
			}
		}
		fmt.Fprintf(r.w, "\n")
	}

	return nil
}

// groupByRule groups findings by rule, returning a sorted rule order and the map.
func groupByRule(findings []model.Finding) ([]string, map[string][]model.Finding) {
	order := []string{}
	byRule := map[string][]model.Finding{}
	for _, f := range findings {
		if _, exists := byRule[f.Rule]; !exists {
			order = append(order, f.Rule)
		}
		byRule[f.Rule] = append(byRule[f.Rule], f)
	}
	sort.Strings(order)
	return order, byRule
}
