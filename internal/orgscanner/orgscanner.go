package orgscanner

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/config"
	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/model"
	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/scanner"
)

// GhRepo represents a GitHub repository with the minimum fields required for scanning.
type GhRepo struct {
	// Name is the repository name (e.g., "klaus-antipatterns-search").
	Name string `json:"name"`
	// CloneURL is the HTTPS clone URL for git clone operations.
	CloneURL string `json:"clone_url"`
	// HTMLURL is the web browsable URL of the repository.
	HTMLURL string `json:"html_url"`
	// Archived indicates whether the repository is archived.
	Archived bool `json:"archived"`
	// Fork indicates whether this is a fork.
	Fork bool `json:"fork"`
}

// FetcherFn retrieves all repositories in an organization, applying profile-based filters
// (exclude_repos, include_forks, include_archived). Implementations typically call GitHub API.
type FetcherFn func(orgName string, profile config.OrgConfig) ([]GhRepo, error)

// ScannerFn scans a single cloned repository and returns its scan result.
// Implementations typically invoke the Scanner on the repo root and clone it first.
type ScannerFn func(repo GhRepo, cfg *config.Config) model.RepoResult

// OrgScanner orchestrates parallel scanning of all repositories in a GitHub organization.
// Uses a worker-pool semaphore to limit concurrency and prevent resource exhaustion.
type OrgScanner struct {
	workers int
	fetchFn FetcherFn
	scanFn  ScannerFn
}

// New creates an OrgScanner with production-ready GitHub API fetcher and scanner.
// workers limits the number of concurrent repo scans.
func New(workers int) *OrgScanner {
	if workers < 1 {
		workers = 4
	}
	return &OrgScanner{
		workers: workers,
		fetchFn: fetchRepos,
		scanFn:  scanRepo,
	}
}

// WithFetcher replaces the repo-enumeration function. Used in tests to inject mock GitHub API responses.
// Returns the OrgScanner for method chaining.
func (s *OrgScanner) WithFetcher(fn FetcherFn) *OrgScanner {
	s.fetchFn = fn
	return s
}

// WithScanner replaces the per-repo scanning function. Used in tests to inject mock scanner behavior.
// Returns the OrgScanner for method chaining.
func (s *OrgScanner) WithScanner(fn ScannerFn) *OrgScanner {
	s.scanFn = fn
	return s
}

// Run enumerates all repositories in the GitHub organization matching the org profile's filters,
// scans each in parallel using the configured worker pool, and returns an aggregated report.
// Returns an error only if fetching the repo list fails; scan errors are captured in individual RepoResult.Err.
func (s *OrgScanner) Run(orgName string, profile config.OrgConfig, cfg *config.Config) (*model.OrgReport, error) {
	repos, err := s.fetchFn(orgName, profile)
	if err != nil {
		return nil, fmt.Errorf("listing repos for %s: %w", orgName, err)
	}

	sem := make(chan struct{}, s.workers)
	var mu sync.Mutex
	var wg sync.WaitGroup
	results := make([]model.RepoResult, 0, len(repos))

	for _, repo := range repos {
		repo := repo
		wg.Add(1)
		sem <- struct{}{}
		go func() {
			defer wg.Done()
			defer func() { <-sem }()

			result := s.scanFn(repo, cfg)

			mu.Lock()
			results = append(results, result)
			mu.Unlock()
		}()
	}
	wg.Wait()

	return &model.OrgReport{Org: orgName, Repos: results}, nil
}

// fetchRepos calls gh CLI to list all repos in an org, applying profile filters.
func fetchRepos(orgName string, profile config.OrgConfig) ([]GhRepo, error) {
	args := []string{
		"api",
		fmt.Sprintf("/orgs/%s/repos", orgName),
		"--paginate",
		"--jq", ".[]",
	}
	cmd := exec.Command("gh", args...)
	if token := os.Getenv(profile.TokenEnv); profile.TokenEnv != "" && token != "" {
		cmd.Env = append(os.Environ(), "GH_TOKEN="+token)
	}
	out, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("gh api: %w", err)
	}

	var repos []GhRepo
	for _, line := range strings.Split(strings.TrimSpace(string(out)), "\n") {
		if line == "" {
			continue
		}
		var r GhRepo
		if err := json.Unmarshal([]byte(line), &r); err != nil {
			continue
		}
		if r.Archived && !profile.IncludeArchived {
			continue
		}
		if r.Fork && !profile.IncludeForks {
			continue
		}
		if IsExcluded(r.Name, profile.ExcludeRepos) {
			continue
		}
		repos = append(repos, r)
	}
	return repos, nil
}

// IsExcluded reports whether name matches any of the glob patterns.
func IsExcluded(name string, patterns []string) bool {
	for _, p := range patterns {
		if matched, _ := filepath.Match(p, name); matched {
			return true
		}
	}
	return false
}

// scanRepo clones the repo to a temp dir, runs the scanner, then cleans up.
func scanRepo(repo GhRepo, cfg *config.Config) model.RepoResult {
	result := model.RepoResult{Repo: repo.Name, URL: repo.HTMLURL}

	tmpDir, err := os.MkdirTemp("", "antipatterns-*")
	if err != nil {
		result.Err = err.Error()
		return result
	}
	defer os.RemoveAll(tmpDir)

	cloneDir := filepath.Join(tmpDir, repo.Name)
	cloneCmd := exec.Command("git", "clone", "--depth=1", "--quiet", repo.CloneURL, cloneDir)
	if out, err := cloneCmd.CombinedOutput(); err != nil {
		result.Err = fmt.Sprintf("clone failed: %s", strings.TrimSpace(string(out)))
		return result
	}

	s := scanner.New(cfg)
	findings, err := s.Run(cloneDir)
	if err != nil {
		result.Err = err.Error()
	}
	result.Findings = findings
	return result
}
