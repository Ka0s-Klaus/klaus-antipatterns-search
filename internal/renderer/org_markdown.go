package renderer

import (
	"fmt"
	"io"
	"time"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/model"
)

// OrgMarkdownRenderer writes an OrgReport as an actionable Markdown document.
// Each repo gets its own section with per-rule findings as checkboxes.
type OrgMarkdownRenderer struct {
	w       io.Writer
	version string
}

func NewOrgMarkdown(w io.Writer, version string) *OrgMarkdownRenderer {
	return &OrgMarkdownRenderer{w: w, version: version}
}

func (r *OrgMarkdownRenderer) Render(report *model.OrgReport) error {
	date := time.Now().Format("2006-01-02")
	total := report.TotalFindings()

	fmt.Fprintf(r.w, "# 🌐 Org Anti-pattern Report\n\n")
	fmt.Fprintf(r.w, "> **Org:** `%s`  \n", report.Org)
	fmt.Fprintf(r.w, "> **Date:** %s  \n", date)
	fmt.Fprintf(r.w, "> **Tool:** antipatterns %s  \n", r.version)
	fmt.Fprintf(r.w, "> **Repos:** %d  |  **Total findings:** %d\n\n", len(report.Repos), total)

	fmt.Fprintf(r.w, "## 📊 Resumen por repositorio\n\n")
	fmt.Fprintf(r.w, "| Repo | Estado | Findings | Top rule |\n")
	fmt.Fprintf(r.w, "|---|---|---|---|\n")
	for _, rr := range report.Repos {
		if rr.Err != "" {
			fmt.Fprintf(r.w, "| `%s` | ❌ error | — | %s |\n", rr.Repo, rr.Err)
			continue
		}
		top := model.TopRuleFor(rr)
		if top == "" {
			top = "—"
		}
		status := "✅"
		if len(rr.Findings) > 0 {
			status = "⚠️"
		}
		fmt.Fprintf(r.w, "| `%s` | %s | %d | `%s` |\n", rr.Repo, status, len(rr.Findings), top)
	}
	fmt.Fprintf(r.w, "\n---\n\n")

	for _, rr := range report.Repos {
		fmt.Fprintf(r.w, "## 📦 `%s`\n\n", rr.Repo)

		if rr.Err != "" {
			fmt.Fprintf(r.w, "> ❌ **Error:** %s\n\n", rr.Err)
			continue
		}

		if len(rr.Findings) == 0 {
			fmt.Fprintf(r.w, "> ✅ No anti-patterns found.\n\n")
			continue
		}

		ruleOrder, ruleFindings := groupByRule(rr.Findings)

		for _, rule := range ruleOrder {
			ff := ruleFindings[rule]
			sev := ff[0].Severity
			fmt.Fprintf(r.w, "### %s `%s` — %d finding(s)\n\n", severityEmoji[sev], rule, len(ff))
			for _, f := range ff {
				if f.Location.Line > 0 {
					fmt.Fprintf(r.w, "- [ ] **`%s:%d`** — %s\n", f.Location.File, f.Location.Line, f.Message)
				} else {
					fmt.Fprintf(r.w, "- [ ] **`%s`** — %s\n", f.Location.File, f.Message)
				}
			}
			fmt.Fprintf(r.w, "\n")
		}
	}

	return nil
}

