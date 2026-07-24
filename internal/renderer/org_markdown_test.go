package renderer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/model"
)

func renderOrgMarkdown(t *testing.T, report *model.OrgReport) string {
	t.Helper()
	var buf bytes.Buffer
	r := NewOrgMarkdown(&buf, "test")
	if err := r.Render(report); err != nil {
		t.Fatalf("Render error: %v", err)
	}
	return buf.String()
}

func TestOrgMarkdownRendererEmpty(t *testing.T) {
	report := &model.OrgReport{Org: "my-org", Repos: []model.RepoResult{}}
	out := renderOrgMarkdown(t, report)

	if !strings.Contains(out, "# 🌐 Org Anti-pattern Report") {
		t.Error("expected org report header")
	}
	if !strings.Contains(out, "`my-org`") {
		t.Error("expected org name in header")
	}
}

func TestOrgMarkdownRendererRepoWithNoFindings(t *testing.T) {
	report := &model.OrgReport{
		Org: "my-org",
		Repos: []model.RepoResult{
			{Repo: "my-org/clean-repo", Findings: []model.Finding{}},
		},
	}
	out := renderOrgMarkdown(t, report)

	if !strings.Contains(out, "## 📦 `my-org/clean-repo`") {
		t.Error("expected repo section header")
	}
	if !strings.Contains(out, "No anti-patterns found") {
		t.Error("expected no-findings message for clean repo")
	}
	if strings.Contains(out, "- [ ]") {
		t.Error("expected no checkboxes for clean repo")
	}
}

func TestOrgMarkdownRendererRepoWithError(t *testing.T) {
	report := &model.OrgReport{
		Org: "my-org",
		Repos: []model.RepoResult{
			{Repo: "my-org/broken-repo", Err: "clone failed"},
		},
	}
	out := renderOrgMarkdown(t, report)

	if !strings.Contains(out, "❌") {
		t.Error("expected error indicator for failed repo")
	}
	if !strings.Contains(out, "clone failed") {
		t.Error("expected error message in output")
	}
}

func TestOrgMarkdownRendererRepoWithFindings(t *testing.T) {
	report := &model.OrgReport{
		Org: "my-org",
		Repos: []model.RepoResult{
			{
				Repo: "my-org/dirty-repo",
				Findings: []model.Finding{
					newFinding("large_function", model.SeverityMedium, "main.go", 42, "too large"),
					newFinding("magic_numbers", model.SeverityLow, "config.go", 10, "magic number"),
				},
			},
		},
	}
	out := renderOrgMarkdown(t, report)

	if !strings.Contains(out, "## 📦 `my-org/dirty-repo`") {
		t.Error("expected repo section")
	}
	if !strings.Contains(out, "- [ ]") {
		t.Error("expected checkboxes for findings")
	}
	if !strings.Contains(out, "main.go:42") {
		t.Error("expected finding location")
	}
}

func TestOrgMarkdownRendererSummaryTable(t *testing.T) {
	report := &model.OrgReport{
		Org: "my-org",
		Repos: []model.RepoResult{
			{Repo: "my-org/a", Findings: []model.Finding{newFinding("r", model.SeverityMedium, "f.go", 1, "msg")}},
			{Repo: "my-org/b", Findings: []model.Finding{}},
		},
	}
	out := renderOrgMarkdown(t, report)

	if !strings.Contains(out, "## 📊 Resumen por repositorio") {
		t.Error("expected summary section")
	}
	if !strings.Contains(out, "Total findings:** 1") {
		t.Error("expected total findings count")
	}
	if !strings.Contains(out, "⚠️") {
		t.Error("expected warning indicator for repo with findings")
	}
	if !strings.Contains(out, "✅") {
		t.Error("expected ok indicator for clean repo")
	}
}

func TestOrgMarkdownRendererVersion(t *testing.T) {
	var buf bytes.Buffer
	r := NewOrgMarkdown(&buf, "v1.2.0")
	_ = r.Render(&model.OrgReport{Org: "org", Repos: nil})
	if !strings.Contains(buf.String(), "antipatterns v1.2.0") {
		t.Error("expected version in org report header")
	}
}
