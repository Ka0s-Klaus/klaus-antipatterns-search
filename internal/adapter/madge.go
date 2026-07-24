package adapter

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/config"
	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/model"
)

// Madge runs madge --circular on root and returns circular-dependency findings.
// Returns nil silently if madge is not installed or fails.
func Madge(root string, cfg *config.Config) ([]model.Finding, error) {
	madgePath, err := exec.LookPath("madge")
	if err != nil {
		return nil, nil
	}
	data, err := runMadge(madgePath, root)
	if err != nil {
		return nil, nil
	}
	return parseMadgeOutput(data, root, cfg)
}

func runMadge(madgePath, root string) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var out bytes.Buffer
	cmd := exec.CommandContext(ctx, madgePath, "--circular", root)
	cmd.Stdout = &out
	// madge exits non-zero when circular deps are found — ignore exit code
	_ = cmd.Run()
	return out.Bytes(), nil
}

func parseMadgeOutput(data []byte, root string, cfg *config.Config) ([]model.Finding, error) {
	sev := model.SeverityFromString(cfg.Severities.CircularDeps)
	var findings []model.Finding

	lines := strings.Split(string(data), "\n")
	cycleNum := 0
	for _, line := range lines {
		line = strings.TrimSpace(line)
		// madge output lines for cycles: "1) a -> b -> a"
		if len(line) == 0 || line[0] < '1' || line[0] > '9' {
			continue
		}
		// strip leading "N) "
		idx := strings.Index(line, ") ")
		if idx < 0 {
			continue
		}
		cycle := strings.TrimSpace(line[idx+2:])
		if cycle == "" {
			continue
		}
		cycleNum++
		findings = append(findings, model.Finding{
			Rule:     "circular_dependency",
			Severity: sev,
			Message:  fmt.Sprintf("circular dependency: %s", cycle),
			Location: model.Location{File: root},
		})
	}
	return findings, nil
}
