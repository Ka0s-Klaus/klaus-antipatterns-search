package scanner

import (
	"io/fs"
	"path/filepath"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/config"
	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/model"
)

// detectorFunc is the signature all Fase-1+ detectors will implement.
type detectorFunc func(path string, cfg *config.Config) ([]model.Finding, error)

// Scanner orchestrates directory walking and anti-pattern detection.
type Scanner struct {
	cfg       *config.Config
	detectors []detectorFunc
}

func New(cfg *config.Config) *Scanner {
	return &Scanner{cfg: cfg}
}

// Run walks root recursively and runs all registered detectors on each file.
// Returns the aggregated slice of findings.
func (s *Scanner) Run(root string) ([]model.Finding, error) {
	var findings []model.Finding

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			if s.cfg.IsExcludedDir(d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		for _, detect := range s.detectors {
			ff, err := detect(path, s.cfg)
			if err != nil {
				// Non-fatal: log and continue so one bad file doesn't abort the scan.
				return nil
			}
			findings = append(findings, ff...)
		}
		return nil
	})

	return findings, err
}
