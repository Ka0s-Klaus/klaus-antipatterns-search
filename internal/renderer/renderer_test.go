package renderer

import (
	"bytes"
	"testing"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/model"
)

func TestConsoleRenderer(t *testing.T) {
	result := &model.ScanResult{
		Findings: []model.Finding{
			{Path: "test.go", Rule: "god_object", Message: "test", Severity: model.SeverityHigh, Line: 1},
		},
		Summary: model.Summary{Total: 1, High: 1, ByRule: map[string]int{"god_object": 1}},
	}

	buf := &bytes.Buffer{}
	renderer := NewConsoleRenderer(buf)
	if err := renderer.Render(result); err != nil {
		t.Fatalf("render failed: %v", err)
	}

	if !bytes.Contains(buf.Bytes(), []byte("god_object")) {
		t.Error("expected 'god_object' in output")
	}
	if !bytes.Contains(buf.Bytes(), []byte("test.go")) {
		t.Error("expected 'test.go' in output")
	}
}

func TestJSONRenderer(t *testing.T) {
	result := &model.ScanResult{
		Findings: []model.Finding{
			{Path: "test.go", Rule: "god_object", Message: "test", Severity: model.SeverityHigh, Line: 1},
		},
		Summary: model.Summary{Total: 1, High: 1, ByRule: map[string]int{"god_object": 1}},
	}

	buf := &bytes.Buffer{}
	renderer := NewJSONRenderer(buf)
	if err := renderer.Render(result); err != nil {
		t.Fatalf("render failed: %v", err)
	}

	if !bytes.Contains(buf.Bytes(), []byte("god_object")) {
		t.Error("expected 'god_object' in JSON output")
	}
	if !bytes.Contains(buf.Bytes(), []byte("test.go")) {
		t.Error("expected 'test.go' in JSON output")
	}
}

func TestMarkdownRenderer(t *testing.T) {
	result := &model.ScanResult{
		Findings: []model.Finding{
			{Path: "test.go", Rule: "god_object", Message: "test", Severity: model.SeverityHigh, Line: 1},
		},
		Summary: model.Summary{Total: 1, High: 1, ByRule: map[string]int{"god_object": 1}},
	}

	buf := &bytes.Buffer{}
	renderer := NewMarkdownRenderer(buf)
	if err := renderer.Render(result); err != nil {
		t.Fatalf("render failed: %v", err)
	}

	if !bytes.Contains(buf.Bytes(), []byte("# 🔍")) {
		t.Error("expected markdown header in output")
	}
}

func TestEmptyResult(t *testing.T) {
	result := &model.ScanResult{
		Findings: []model.Finding{},
		Summary:  model.Summary{Total: 0},
	}

	buf := &bytes.Buffer{}
	renderer := NewConsoleRenderer(buf)
	if err := renderer.Render(result); err != nil {
		t.Fatalf("render failed: %v", err)
	}

	output := buf.String()
	if !bytes.Contains(buf.Bytes(), []byte("No antipatterns")) {
		t.Errorf("expected 'No antipatterns' message, got: %s", output)
	}
}
