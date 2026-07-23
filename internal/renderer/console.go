package renderer

import (
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/model"
)

const (
	ansiReset    = "\033[0m"
	ansiCyan     = "\033[36m"
	ansiGreen    = "\033[32m"
	ansiYellow   = "\033[33m"
	ansiRed      = "\033[31m"
	ansiMagenta  = "\033[35m"
	ansiBold     = "\033[1m"
)

var severityColor = map[model.Severity]string{
	model.SeverityInfo:     ansiCyan,
	model.SeverityLow:      ansiGreen,
	model.SeverityMedium:   ansiYellow,
	model.SeverityHigh:     ansiRed,
	model.SeverityCritical: ansiMagenta,
}

// ConsoleRenderer writes a human-readable, color-coded table to w.
type ConsoleRenderer struct {
	w io.Writer
}

func NewConsole(w io.Writer) *ConsoleRenderer {
	return &ConsoleRenderer{w: w}
}

func (r *ConsoleRenderer) Render(findings []model.Finding) error {
	if len(findings) == 0 {
		fmt.Fprintf(r.w, "%s✅ No anti-patterns found.%s\n", ansiGreen, ansiReset)
		return nil
	}

	counts := map[model.Severity]int{}
	for _, f := range findings {
		counts[f.Severity]++
	}

	fmt.Fprintf(r.w, "\n%s🔍 %d finding(s) detected%s — critical:%d high:%d medium:%d low:%d info:%d\n\n",
		ansiBold, len(findings), ansiReset,
		counts[model.SeverityCritical],
		counts[model.SeverityHigh],
		counts[model.SeverityMedium],
		counts[model.SeverityLow],
		counts[model.SeverityInfo],
	)

	tw := tabwriter.NewWriter(r.w, 0, 0, 2, ' ', 0)
	fmt.Fprintf(tw, "%sSEVERITY\tRULE\tLOCATION\tMESSAGE%s\n", ansiBold, ansiReset)
	fmt.Fprintf(tw, "────────\t────\t────────\t───────\n")

	for _, f := range findings {
		color := severityColor[f.Severity]
		fmt.Fprintf(tw, "%s%s%s\t%s\t%s:%d\t%s\n",
			color, f.Severity, ansiReset,
			f.Rule,
			f.Location.File, f.Location.Line,
			f.Message,
		)
	}
	return tw.Flush()
}
