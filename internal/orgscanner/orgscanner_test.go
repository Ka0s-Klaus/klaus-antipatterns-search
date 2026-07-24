package orgscanner_test

import (
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/config"
	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/model"
	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/orgscanner"
)

// TestScanOrgParallel verifies that concurrent scans do not produce data races.
// Run with: go test -race ./internal/orgscanner/...
func TestScanOrgParallel(t *testing.T) {
	const numRepos = 20

	fakeRepos := make([]orgscanner.GhRepo, numRepos)
	for i := range fakeRepos {
		fakeRepos[i] = orgscanner.GhRepo{Name: fmt.Sprintf("repo-%02d", i)}
	}

	var callCount int64
	s := orgscanner.New(4).
		WithFetcher(func(_ string, _ config.OrgConfig) ([]orgscanner.GhRepo, error) {
			return fakeRepos, nil
		}).
		WithScanner(func(repo orgscanner.GhRepo, _ *config.Config) model.RepoResult {
			atomic.AddInt64(&callCount, 1)
			return model.RepoResult{
				Repo:     repo.Name,
				Findings: []model.Finding{{Rule: "magic_number", Severity: model.SeverityLow}},
			}
		})

	report, err := s.Run("test-org", config.OrgConfig{}, config.Default())
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Repos) != numRepos {
		t.Errorf("want %d repos, got %d", numRepos, len(report.Repos))
	}
	if int(callCount) != numRepos {
		t.Errorf("want %d scan calls, got %d", numRepos, int(callCount))
	}
}

// TestScanOrgEmpty verifies empty org produces empty report without error.
func TestScanOrgEmpty(t *testing.T) {
	s := orgscanner.New(2).
		WithFetcher(func(_ string, _ config.OrgConfig) ([]orgscanner.GhRepo, error) {
			return nil, nil
		}).
		WithScanner(func(_ orgscanner.GhRepo, _ *config.Config) model.RepoResult {
			t.Fatal("scan must not be called when repo list is empty")
			return model.RepoResult{}
		})

	report, err := s.Run("empty-org", config.OrgConfig{}, config.Default())
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Repos) != 0 {
		t.Errorf("expected 0 repos, got %d", len(report.Repos))
	}
	if report.TotalFindings() != 0 {
		t.Errorf("expected 0 findings, got %d", report.TotalFindings())
	}
}

// TestRepoFilter verifies glob-based repo name exclusion.
func TestRepoFilter(t *testing.T) {
	tests := []struct {
		name     string
		repoName string
		patterns []string
		want     bool
	}{
		{"no patterns", "my-repo", nil, false},
		{"exact match", "mirror-repo", []string{"mirror-repo"}, true},
		{"glob prefix", "mirror-main", []string{"mirror-*"}, true},
		{"glob no match", "service-api", []string{"mirror-*"}, false},
		{"multi first match", "mirror-repo", []string{"mirror-*", "fork-*"}, true},
		{"multi second match", "fork-legacy", []string{"mirror-*", "fork-*"}, true},
		{"multi no match", "service-api", []string{"mirror-*", "fork-*"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := orgscanner.IsExcluded(tt.repoName, tt.patterns)
			if got != tt.want {
				t.Errorf("IsExcluded(%q, %v) = %v, want %v", tt.repoName, tt.patterns, got, tt.want)
			}
		})
	}
}

// TestAggregateReport verifies TotalFindings and TopRule across a multi-repo report.
func TestAggregateReport(t *testing.T) {
	report := &model.OrgReport{
		Org: "test-org",
		Repos: []model.RepoResult{
			{Repo: "repo-a", Findings: []model.Finding{
				{Rule: "magic_number"}, {Rule: "magic_number"}, {Rule: "large_function"},
			}},
			{Repo: "repo-b", Findings: []model.Finding{
				{Rule: "magic_number"},
			}},
			{Repo: "repo-c"},
			{Repo: "repo-err", Err: "clone failed"},
		},
	}

	if got := report.TotalFindings(); got != 4 {
		t.Errorf("TotalFindings() = %d, want 4", got)
	}
	if got := report.TopRule(); got != "magic_number" {
		t.Errorf("TopRule() = %q, want %q", got, "magic_number")
	}
}

// TestScanOrgWorkersCapped verifies that more workers than repos doesn't deadlock.
func TestScanOrgWorkersCapped(t *testing.T) {
	fakeRepos := []orgscanner.GhRepo{{Name: "only-repo"}}

	s := orgscanner.New(100). // 100 workers, 1 repo
		WithFetcher(func(_ string, _ config.OrgConfig) ([]orgscanner.GhRepo, error) {
			return fakeRepos, nil
		}).
		WithScanner(func(repo orgscanner.GhRepo, _ *config.Config) model.RepoResult {
			return model.RepoResult{Repo: repo.Name}
		})

	report, err := s.Run("test-org", config.OrgConfig{}, config.Default())
	if err != nil {
		t.Fatal(err)
	}
	if len(report.Repos) != 1 {
		t.Errorf("want 1 repo, got %d", len(report.Repos))
	}
}
