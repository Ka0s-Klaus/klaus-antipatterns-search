package adapter

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/config"
	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/model"
)

type radonEntry struct {
	Name       string `json:"name"`
	Type       string `json:"type"`
	Complexity int    `json:"complexity"`
	Lineno     int    `json:"lineno"`
}

// Radon adapts radon (Python code metrics tool) for cyclomatic complexity detection.
// Executes 'radon cc -j' on the directory tree, parses JSON output, and returns findings
// normalized to model.Finding. Returns ErrToolNotFound if radon is not installed; other errors
// are non-fatal and silently skipped to enable graceful degradation.
func Radon(root string, cfg *config.Config) ([]model.Finding, error) {
	radonPath, err := exec.LookPath("radon")
	if err != nil {
		return nil, ErrToolNotFound
	}
	data, err := runRadon(radonPath, root)
	if err != nil {
		return nil, nil
	}
	return parseRadonOutput(data, cfg)
}

func runRadon(radonPath, root string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var out bytes.Buffer
	cmd := exec.CommandContext(ctx, radonPath, "cc", "-j", root)
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func parseRadonOutput(data []byte, cfg *config.Config) ([]model.Finding, error) {
	// radon cc -j: map from filepath to []radonEntry
	var report map[string][]radonEntry
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, fmt.Errorf("radon: parse error: %w", err)
	}

	sev := model.SeverityFromString(cfg.Severities.Cyclomatic)
	var findings []model.Finding

	for file, entries := range report {
		for _, e := range entries {
			if e.Complexity > cfg.Thresholds.Cyclomatic {
				findings = append(findings, model.Finding{
					Rule:     "cyclomatic_complexity",
					Severity: sev,
					Message:  fmt.Sprintf("%s %s: cyclomatic complexity %d (threshold %d)", e.Type, e.Name, e.Complexity, cfg.Thresholds.Cyclomatic),
					Location: model.Location{File: file, Line: e.Lineno},
				})
			}
		}
	}
	return findings, nil
}
