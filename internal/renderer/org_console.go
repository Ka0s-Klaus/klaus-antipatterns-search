package renderer

import (
	"fmt"
	"io"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/model"
)

const (
	orgConsoleWidth = 72
	orgColRepo      = 32
	orgColFindings  = 10
	orgColTopRule   = 22
)

// OrgConsoleRenderer renders an OrgReport as a human-readable table.
type OrgConsoleRenderer struct {
	w io.Writer
}

// NewOrgConsole returns an OrgConsoleRenderer writing to w.
func NewOrgConsole(w io.Writer) *OrgConsoleRenderer {
	return &OrgConsoleRenderer{w: w}
}

// Render writes the org report table to the underlying writer.
func (r *OrgConsoleRenderer) Render(report *model.OrgReport) error {
	sep := repeatChar('─', orgConsoleWidth)

	fmt.Fprintf(r.w, "\n🌐 Org scan: %s (%d repos)\n", report.Org, len(report.Repos))
	fmt.Fprintln(r.w, sep)
	fmt.Fprintf(r.w, " %-*s  %*s  %-*s\n",
		orgColRepo, "REPO",
		orgColFindings, "FINDINGS",
		orgColTopRule, "TOP RULE",
	)
	fmt.Fprintln(r.w, sep)

	for _, rr := range report.Repos {
		if rr.Err != "" {
			fmt.Fprintf(r.w, " %-*s  %*s  %-*s\n",
				orgColRepo, truncate(rr.Repo, orgColRepo),
				orgColFindings, "ERROR",
				orgColTopRule, truncate(rr.Err, orgColTopRule),
			)
			continue
		}
		count := len(rr.Findings)
		top := model.TopRuleFor(rr)
		if top == "" {
			top = "—"
		}
		status := "✅"
		if count > 0 {
			status = "⚠️"
		}
		fmt.Fprintf(r.w, " %-*s  %*d  %-*s  %s\n",
			orgColRepo, truncate(rr.Repo, orgColRepo),
			orgColFindings, count,
			orgColTopRule, truncate(top, orgColTopRule),
			status,
		)
	}

	fmt.Fprintln(r.w, sep)
	total := report.TotalFindings()
	topAll := report.TopRule()
	if topAll == "" {
		topAll = "—"
	}
	fmt.Fprintf(r.w, " %-*s  %*d  %-*s\n",
		orgColRepo, "TOTAL",
		orgColFindings, total,
		orgColTopRule, topAll,
	)
	fmt.Fprintln(r.w, sep)
	fmt.Fprintln(r.w)
	return nil
}

func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}
	return s[:max-1] + "…"
}

func repeatChar(ch rune, n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = ch
	}
	return string(b)
}
