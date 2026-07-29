package scanner

import (
	"errors"
	"fmt"
	"io"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/adapter"
	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/config"
	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/detector"
	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/model"
)

// detectorFunc is the signature for native detectors: one invocation per file.
// path is the absolute file path; returns findings or an error if the file cannot be analyzed.
type detectorFunc func(path string, cfg *config.Config) ([]model.Finding, error)

// dirAdapterFunc is the signature for OSS adapters: one invocation per directory scan.
// root is the scan root; returns findings or an error if the tool fails or is not installed.
type dirAdapterFunc func(root string, cfg *config.Config) ([]model.Finding, error)

// namedDetector pairs a human-readable name with a detectorFunc for progress reporting.
type namedDetector struct {
	name string
	fn   detectorFunc
}

// namedAdapter pairs a human-readable name with a dirAdapterFunc for progress reporting.
type namedAdapter struct {
	name string
	fn   dirAdapterFunc
}

// Scanner orchestrates directory walking and anti-pattern detection.
// It runs native detectors on each file and OSS adapters on the entire directory.
type Scanner struct {
	cfg         *config.Config
	detectors   []namedDetector
	dirAdapters []namedAdapter
	log         io.Writer // io.Discard = silent; os.Stderr when --verbose
}

// New creates a Scanner with default native detectors and OSS adapters, logging disabled.
// Use NewVerbose to enable progress output.
func New(cfg *config.Config) *Scanner {
	return &Scanner{
		cfg: cfg,
		detectors: []namedDetector{
			{name: "large_function", fn: detector.LargeFunction},
			{name: "god_object", fn: detector.GodObject},
			{name: "magic_numbers", fn: detector.MagicNumbers},
		},
		dirAdapters: []namedAdapter{
			{name: "jscpd", fn: adapter.Jscpd},
			{name: "madge", fn: adapter.Madge},
			{name: "radon", fn: adapter.Radon},
			{name: "gocyclo", fn: adapter.Gocyclo},
		},
		log: io.Discard,
	}
}

// NewVerbose creates a Scanner that writes progress reports and detector names to w.
// Typically w is os.Stderr; use io.Discard to suppress output. Progress is printed
// to show which native detectors and OSS adapters ran and their match counts.
func NewVerbose(cfg *config.Config, w io.Writer) *Scanner {
	s := New(cfg)
	s.log = w
	return s
}

func (s *Scanner) logf(format string, args ...any) {
	fmt.Fprintf(s.log, format+"\n", args...)
}

// Run walks the directory tree rooted at root and executes all native detectors on each .go file,
// then executes OSS adapters once on the root directory. Returns aggregated findings from all detectors.
// Non-fatal errors (e.g., unreadable files, missing OSS tools) are silently skipped with ErrToolNotFound.
func (s *Scanner) Run(root string) ([]model.Finding, error) {
	var findings []model.Finding

	// Count Go files once for the verbose summary line
	goFiles := countGoFiles(root, s.cfg)
	detectorNames := make([]string, len(s.detectors))
	for i, d := range s.detectors {
		detectorNames[i] = d.name
	}
	s.logf("[native] %s → %d .go file(s)", strings.Join(detectorNames, ", "), goFiles)

	// File-level native detectors
	nativeFindings := 0
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
		for _, det := range s.detectors {
			ff, err := det.fn(path, s.cfg)
			if err != nil {
				// Non-fatal: one bad file must not abort the scan.
				return nil
			}
			nativeFindings += len(ff)
			findings = append(findings, ff...)
		}
		return nil
	})
	if err != nil {
		return findings, err
	}
	s.logf("[native] done → %d finding(s)", nativeFindings)

	// Dir-level OSS adapters (run once over root, skip gracefully if tool absent)
	for _, adp := range s.dirAdapters {
		s.logf("[oss]    %s: running...", adp.name)
		ff, err := adp.fn(root, s.cfg)
		if err != nil {
			if errors.Is(err, adapter.ErrToolNotFound) {
				s.logf("[skip]   %s: not found in PATH", adp.name)
			} else {
				s.logf("[oss]    %s: error — %v", adp.name, err)
			}
			continue
		}
		s.logf("[oss]    %s: %d finding(s)", adp.name, len(ff))
		findings = append(findings, ff...)
	}

	return findings, nil
}

// countGoFiles returns the number of .go files under root (respecting excluded dirs).
func countGoFiles(root string, cfg *config.Config) int {
	count := 0
	_ = filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if d.IsDir() {
			if cfg.IsExcludedDir(d.Name()) {
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(path, ".go") {
			count++
		}
		return nil
	})
	return count
}
