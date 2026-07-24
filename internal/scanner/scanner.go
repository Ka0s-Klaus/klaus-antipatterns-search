package scanner

import (
	"io/fs"
	"path/filepath"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/adapter"
	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/config"
	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/detector"
	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/model"
)

// detectorFunc operates on a single file.
type detectorFunc func(path string, cfg *config.Config) ([]model.Finding, error)

// dirAdapterFunc operates on the entire root directory (OSS tools).
type dirAdapterFunc func(root string, cfg *config.Config) ([]model.Finding, error)

// Scanner orchestrates directory walking and anti-pattern detection.
type Scanner struct {
	cfg         *config.Config
	detectors   []detectorFunc
	dirAdapters []dirAdapterFunc
}

func New(cfg *config.Config) *Scanner {
	return &Scanner{
		cfg: cfg,
		detectors: []detectorFunc{
			detector.LargeFunction,
			detector.GodObject,
			detector.MagicNumbers,
		},
		dirAdapters: []dirAdapterFunc{
			adapter.Jscpd,
			adapter.Madge,
			adapter.Radon,
			adapter.Gocyclo,
		},
	}
}

// Run walks root recursively and runs all registered detectors on each file,
// then runs OSS dir-level adapters once on the root.
// Returns the aggregated slice of findings.
func (s *Scanner) Run(root string) ([]model.Finding, error) {
	var findings []model.Finding

	// File-level native detectors
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
				// Non-fatal: one bad file must not abort the scan.
				return nil
			}
			findings = append(findings, ff...)
		}
		return nil
	})
	if err != nil {
		return findings, err
	}

	// Dir-level OSS adapters (run once over root, skip gracefully if tool absent)
	for _, adapt := range s.dirAdapters {
		ff, err := adapt(root, s.cfg)
		if err != nil {
			continue
		}
		findings = append(findings, ff...)
	}

	return findings, nil
}
