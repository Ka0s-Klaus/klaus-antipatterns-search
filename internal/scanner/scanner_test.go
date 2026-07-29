package scanner

import (
	"io"
	"testing"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/config"
	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/model"
)

// mockDetector returns a fixed finding for testing.
func mockDetector(path string, cfg *config.Config) ([]model.Finding, error) {
	return []model.Finding{
		{
			Rule:     "mock_detector",
			Message:  "test finding",
			Severity: model.SeverityMedium,
			Location: model.Location{File: path, Line: 10, Column: 0},
		},
	}, nil
}

// mockAdapter returns a fixed finding for testing.
func mockAdapter(root string, cfg *config.Config) ([]model.Finding, error) {
	return []model.Finding{
		{
			Rule:     "mock_adapter",
			Message:  "test adapter finding",
			Severity: model.SeverityLow,
			Location: model.Location{File: "test.js", Line: 5, Column: 0},
		},
	}, nil
}

func TestScannerWithMockDetector(t *testing.T) {
	cfg := config.Default()
	s := New(cfg).WithDetectors(
		namedDetector{name: "mock", fn: mockDetector},
	)

	// Since we're using a mock detector that doesn't walk files,
	// we test that the detector is registered and callable
	if len(s.detectors) != 1 || s.detectors[0].name != "mock" {
		t.Fatalf("detector not injected correctly: %v", s.detectors)
	}
}

func TestScannerWithMockAdapter(t *testing.T) {
	cfg := config.Default()
	s := New(cfg).WithAdapters(
		namedAdapter{name: "mock", fn: mockAdapter},
	)

	// Verify adapter injection
	if len(s.dirAdapters) != 1 || s.dirAdapters[0].name != "mock" {
		t.Fatalf("adapter not injected correctly: %v", s.dirAdapters)
	}
}

func TestScannerMethodChaining(t *testing.T) {
	cfg := config.Default()
	s := New(cfg).
		WithDetectors(namedDetector{name: "mock_det", fn: mockDetector}).
		WithAdapters(namedAdapter{name: "mock_adp", fn: mockAdapter})

	if len(s.detectors) != 1 || s.detectors[0].name != "mock_det" {
		t.Error("detector not injected via chaining")
	}
	if len(s.dirAdapters) != 1 || s.dirAdapters[0].name != "mock_adp" {
		t.Error("adapter not injected via chaining")
	}
}

func TestScannerVerboseLogging(t *testing.T) {
	cfg := config.Default()
	s := NewVerbose(cfg, io.Discard)

	if s.log != io.Discard {
		t.Error("verbose logging not set correctly")
	}
}

func TestScannerDefaultDetectorsAndAdapters(t *testing.T) {
	cfg := config.Default()
	s := New(cfg)

	// Default detectors
	if len(s.detectors) != 3 {
		t.Errorf("expected 3 default detectors, got %d", len(s.detectors))
	}
	wantDetectors := map[string]bool{
		"large_function": true,
		"god_object":     true,
		"magic_numbers":  true,
	}
	for _, d := range s.detectors {
		if !wantDetectors[d.name] {
			t.Errorf("unexpected default detector: %s", d.name)
		}
	}

	// Default adapters
	if len(s.dirAdapters) != 4 {
		t.Errorf("expected 4 default adapters, got %d", len(s.dirAdapters))
	}
	wantAdapters := map[string]bool{
		"jscpd":   true,
		"madge":   true,
		"radon":   true,
		"gocyclo": true,
	}
	for _, a := range s.dirAdapters {
		if !wantAdapters[a.name] {
			t.Errorf("unexpected default adapter: %s", a.name)
		}
	}
}
