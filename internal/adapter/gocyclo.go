package adapter

import (
	"bytes"
	"context"
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/config"
	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/model"
)

// Gocyclo runs gocyclo -over N on root and returns cyclomatic-complexity findings.
// Returns nil silently if gocyclo is not installed or fails.
func Gocyclo(root string, cfg *config.Config) ([]model.Finding, error) {
	gocycloPath, err := exec.LookPath("gocyclo")
	if err != nil {
		return nil, ErrToolNotFound
	}
	data, err := runGocyclo(gocycloPath, root, cfg.Thresholds.Cyclomatic)
	if err != nil {
		return nil, nil
	}
	return parseGocycloOutput(data, cfg)
}

func runGocyclo(gocycloPath, root string, threshold int) ([]byte, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	var out bytes.Buffer
	cmd := exec.CommandContext(ctx, gocycloPath, "-over", strconv.Itoa(threshold), root)
	cmd.Stdout = &out
	// gocyclo exits 1 when findings exist — ignore exit code
	_ = cmd.Run()
	return out.Bytes(), nil
}

func parseGocycloOutput(data []byte, cfg *config.Config) ([]model.Finding, error) {
	// gocyclo output: "<complexity> <package> <func> <file>:<line>:<col>"
	sev := model.SeverityFromString(cfg.Severities.Cyclomatic)
	var findings []model.Finding

	for _, line := range strings.Split(strings.TrimSpace(string(data)), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.Fields(line)
		if len(parts) < 4 {
			continue
		}
		complexity, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}
		fnName := parts[2]
		location := parts[3] // "file.go:10:1"

		file, lineNum := splitLocation(location)
		findings = append(findings, model.Finding{
			Rule:     "cyclomatic_complexity",
			Severity: sev,
			Message:  fmt.Sprintf("%s: cyclomatic complexity %d (threshold %d)", fnName, complexity, cfg.Thresholds.Cyclomatic),
			Location: model.Location{File: file, Line: lineNum},
		})
	}
	return findings, nil
}

// splitLocation parses "file.go:10:1" → ("file.go", 10).
func splitLocation(loc string) (string, int) {
	parts := strings.Split(loc, ":")
	if len(parts) < 2 {
		return loc, 0
	}
	file := parts[0]
	line, _ := strconv.Atoi(parts[1])
	return file, line
}
