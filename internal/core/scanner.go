package core

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/config"
	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/model"
)

// Scanner orchestrates the antipattern detection across a codebase
type Scanner struct {
	config *config.Config
}

// NewScanner creates a new scanner with the given config
func NewScanner(cfg *config.Config) *Scanner {
	return &Scanner{config: cfg}
}

// Scan performs a complete scan on the given path
func (s *Scanner) Scan(path string) (*model.ScanResult, error) {
	start := time.Now()

	// Resolve path
	absPath, err := filepath.Abs(path)
	if err != nil {
		return nil, fmt.Errorf("failed to resolve path: %w", err)
	}

	// Check path exists
	info, err := os.Stat(absPath)
	if err != nil {
		return nil, fmt.Errorf("path does not exist: %w", err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("path must be a directory, got file")
	}

	// For now, just return empty scan result
	// Phase 1 will add actual detectors
	findings := []model.Finding{}
	filesScanned := 0

	result := &model.ScanResult{
		Findings:     findings,
		Summary:      model.CalculateSummary(findings),
		Duration:     time.Since(start),
		FilesScanned: filesScanned,
	}

	return result, nil
}
