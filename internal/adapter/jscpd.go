package adapter

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/config"
	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/model"
)

type jscpdReport struct {
	Statistics struct {
		Percentage float64 `json:"percentage"`
	} `json:"statistics"`
	Duplicates []jscpdDuplicate `json:"duplicates"`
}

type jscpdDuplicate struct {
	Lines      int          `json:"lines"`
	FirstFile  jscpdFile   `json:"firstFile"`
	SecondFile jscpdFile   `json:"secondFile"`
}

type jscpdFile struct {
	Name  string `json:"name"`
	Start int    `json:"start"`
}

// Jscpd runs jscpd on root and returns duplication findings.
// Returns nil silently if jscpd is not installed or fails.
func Jscpd(root string, cfg *config.Config) ([]model.Finding, error) {
	jscpdPath, err := exec.LookPath("jscpd")
	if err != nil {
		return nil, nil
	}
	data, err := runJscpd(jscpdPath, root)
	if err != nil {
		return nil, nil
	}
	return parseJscpdOutput(data, cfg)
}

func runJscpd(jscpdPath, root string) ([]byte, error) {
	tmpDir, err := os.MkdirTemp("", "jscpd-*")
	if err != nil {
		return nil, err
	}
	defer os.RemoveAll(tmpDir)

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, jscpdPath, root, "--reporters", "json", "--output", tmpDir, "--silent")
	if err := cmd.Run(); err != nil {
		return nil, err
	}
	return os.ReadFile(filepath.Join(tmpDir, "jscpd-report.json"))
}

func parseJscpdOutput(data []byte, cfg *config.Config) ([]model.Finding, error) {
	var report jscpdReport
	if err := json.Unmarshal(data, &report); err != nil {
		return nil, fmt.Errorf("jscpd: parse error: %w", err)
	}

	sev := model.SeverityFromString(cfg.Severities.Duplication)
	var findings []model.Finding

	if report.Statistics.Percentage >= cfg.Thresholds.DuplicationPct && cfg.Thresholds.DuplicationPct > 0 {
		findings = append(findings, model.Finding{
			Rule:     "duplication",
			Severity: sev,
			Message:  fmt.Sprintf("%.1f%% code duplication detected (threshold %.1f%%)", report.Statistics.Percentage, cfg.Thresholds.DuplicationPct),
			Location: model.Location{File: "."},
		})
	}

	for _, dup := range report.Duplicates {
		findings = append(findings, model.Finding{
			Rule:     "duplication",
			Severity: sev,
			Message:  fmt.Sprintf("clone: %d lines duplicated in %s", dup.Lines, dup.SecondFile.Name),
			Location: model.Location{File: dup.FirstFile.Name, Line: dup.FirstFile.Start},
		})
	}
	return findings, nil
}
