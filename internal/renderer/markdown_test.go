package renderer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/model"
)

func renderMarkdown(t *testing.T, findings []model.Finding, root string) string {
	t.Helper()
	var buf bytes.Buffer
	r := NewMarkdown(&buf, "test", root)
	if err := r.Render(findings); err != nil {
		t.Fatalf("Render error: %v", err)
	}
	return buf.String()
}

func TestMarkdownRendererEmpty(t *testing.T) {
	out := renderMarkdown(t, []model.Finding{}, ".")
	if !strings.Contains(out, "# 🔍 Anti-pattern Report") {
		t.Error("expected report header")
	}
	if !strings.Contains(out, "No anti-patterns found") {
		t.Error("expected no-findings message")
	}
	if strings.Contains(out, "- [ ]") {
		t.Error("expected no checkboxes for empty findings")
	}
}

func TestMarkdownRendererSingleFinding(t *testing.T) {
	findings := []model.Finding{
		newFinding("large_function", model.SeverityMedium, "/repo/pkg/foo.go", 42, "function too large"),
	}
	out := renderMarkdown(t, findings, "/repo")

	if !strings.Contains(out, "- [ ]") {
		t.Error("expected checkbox for finding")
	}
	if !strings.Contains(out, "pkg/foo.go:42") {
		t.Error("expected relative path with line number")
	}
	if !strings.Contains(out, "function too large") {
		t.Error("expected finding message")
	}
	if !strings.Contains(out, "`large_function`") {
		t.Error("expected rule name in section header")
	}
}

func TestMarkdownRendererRulesSorted(t *testing.T) {
	findings := []model.Finding{
		newFinding("magic_numbers", model.SeverityLow, "b.go", 5, "magic number"),
		newFinding("large_function", model.SeverityMedium, "a.go", 10, "too large"),
	}
	out := renderMarkdown(t, findings, ".")

	idxLarge := strings.Index(out, "`large_function`")
	idxMagic := strings.Index(out, "`magic_numbers`")
	if idxLarge == -1 || idxMagic == -1 {
		t.Fatal("expected both rule sections in output")
	}
	if idxLarge > idxMagic {
		t.Error("expected large_function before magic_numbers (alphabetical order)")
	}
}

func TestMarkdownRendererSummaryTable(t *testing.T) {
	findings := []model.Finding{
		newFinding("large_function", model.SeverityMedium, "a.go", 1, "too large"),
		newFinding("large_function", model.SeverityMedium, "b.go", 2, "too large"),
		newFinding("god_object", model.SeverityHigh, "c.go", 3, "god object"),
	}
	out := renderMarkdown(t, findings, ".")

	if !strings.Contains(out, "## 📊 Resumen") {
		t.Error("expected summary section")
	}
	if !strings.Contains(out, "Total: 3 finding(s)") {
		t.Error("expected total count in summary")
	}
}

func TestMarkdownRendererSeverityEmoji(t *testing.T) {
	cases := []struct {
		sev   model.Severity
		emoji string
	}{
		{model.SeverityInfo, "🔵"},
		{model.SeverityLow, "🟢"},
		{model.SeverityMedium, "🟡"},
		{model.SeverityHigh, "🟠"},
		{model.SeverityCritical, "🔴"},
	}
	for _, tc := range cases {
		t.Run(string(tc.sev), func(t *testing.T) {
			findings := []model.Finding{newFinding("rule", tc.sev, "f.go", 1, "msg")}
			out := renderMarkdown(t, findings, ".")
			if !strings.Contains(out, tc.emoji) {
				t.Errorf("severity %s: expected emoji %s in output", tc.sev, tc.emoji)
			}
		})
	}
}

func TestMarkdownRendererZeroLineFinding(t *testing.T) {
	findings := []model.Finding{
		{Rule: "duplication", Severity: model.SeverityMedium, Message: "10 duplicated lines", Location: model.Location{File: "src/a.js"}},
	}
	out := renderMarkdown(t, findings, ".")
	if !strings.Contains(out, "- [ ]") {
		t.Error("expected checkbox for finding without line number")
	}
	if strings.Contains(out, ":0") {
		t.Error("must not print ':0' for findings without line number")
	}
}

func TestMarkdownRendererVersion(t *testing.T) {
	var buf bytes.Buffer
	r := NewMarkdown(&buf, "v1.2.0", ".")
	_ = r.Render([]model.Finding{})
	if !strings.Contains(buf.String(), "antipatterns v1.2.0") {
		t.Error("expected version in report header")
	}
}
