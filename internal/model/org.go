package model

// RepoResult holds the scan outcome for a single repository.
type RepoResult struct {
	Repo     string    `json:"repo"`
	URL      string    `json:"url,omitempty"`
	Findings []Finding `json:"findings"`
	Err      string    `json:"error,omitempty"`
}

// OrgReport aggregates scan results across all repos of an organization.
type OrgReport struct {
	Org   string       `json:"org"`
	Repos []RepoResult `json:"repos"`
}

// TotalFindings sums findings across all repos (errors are counted as 0 findings).
func (r *OrgReport) TotalFindings() int {
	total := 0
	for _, rr := range r.Repos {
		total += len(rr.Findings)
	}
	return total
}

// TopRule returns the most frequent rule ID across all repos, or "" if no findings.
func (r *OrgReport) TopRule() string {
	counts := make(map[string]int)
	for _, rr := range r.Repos {
		for _, f := range rr.Findings {
			counts[f.Rule]++
		}
	}
	top, maxC := "", 0
	for rule, count := range counts {
		if count > maxC {
			top, maxC = rule, count
		}
	}
	return top
}

// topRuleFor returns the most frequent rule in a single RepoResult.
func topRuleFor(rr RepoResult) string {
	counts := make(map[string]int)
	for _, f := range rr.Findings {
		counts[f.Rule]++
	}
	top, maxC := "", 0
	for rule, count := range counts {
		if count > maxC {
			top, maxC = rule, count
		}
	}
	return top
}

// TopRuleFor is the exported version of topRuleFor for use by renderers.
func TopRuleFor(rr RepoResult) string { return topRuleFor(rr) }
