package renderer

import (
	"bytes"
	"encoding/json"
	"strings"
	"testing"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/model"
)

func newFinding(rule string, sev model.Severity, file string, line int, msg string) model.Finding {
	return model.Finding{
		Rule:     rule,
		Severity: sev,
		Message:  msg,
		Location: model.Location{File: file, Line: line},
	}
}

func renderSARIF(t *testing.T, findings []model.Finding, root string) map[string]interface{} {
	t.Helper()
	var buf bytes.Buffer
	r := NewSARIF(&buf, "test", root)
	if err := r.Render(findings); err != nil {
		t.Fatalf("Render error: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(buf.Bytes(), &out); err != nil {
		t.Fatalf("invalid JSON: %v — output was:\n%s", err, buf.String())
	}
	return out
}

func TestSARIFRendererEmpty(t *testing.T) {
	out := renderSARIF(t, []model.Finding{}, ".")

	if out["version"] != "2.1.0" {
		t.Errorf("expected version 2.1.0, got %v", out["version"])
	}
	runs := out["runs"].([]interface{})
	if len(runs) != 1 {
		t.Fatalf("expected 1 run, got %d", len(runs))
	}
	results := runs[0].(map[string]interface{})["results"].([]interface{})
	if len(results) != 0 {
		t.Errorf("expected 0 results for empty findings, got %d", len(results))
	}
}

func TestSARIFRendererSingleFinding(t *testing.T) {
	findings := []model.Finding{
		newFinding("large_function", model.SeverityMedium, "/repo/pkg/foo.go", 42, "function too large"),
	}
	out := renderSARIF(t, findings, "/repo")

	runs := out["runs"].([]interface{})
	results := runs[0].(map[string]interface{})["results"].([]interface{})
	if len(results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(results))
	}
	r := results[0].(map[string]interface{})
	if r["ruleId"] != "large_function" {
		t.Errorf("unexpected ruleId: %v", r["ruleId"])
	}
	if r["level"] != "warning" {
		t.Errorf("expected level warning for medium severity, got %v", r["level"])
	}
	msg := r["message"].(map[string]interface{})["text"]
	if msg != "function too large" {
		t.Errorf("unexpected message: %v", msg)
	}
	// verify relative path
	locs := r["locations"].([]interface{})
	uri := locs[0].(map[string]interface{})["physicalLocation"].(map[string]interface{})["artifactLocation"].(map[string]interface{})["uri"]
	if uri != "pkg/foo.go" {
		t.Errorf("expected relative URI pkg/foo.go, got %v", uri)
	}
}

func TestSARIFRendererSeverityMapping(t *testing.T) {
	cases := []struct {
		sev   model.Severity
		level string
	}{
		{model.SeverityCritical, "error"},
		{model.SeverityHigh, "error"},
		{model.SeverityMedium, "warning"},
		{model.SeverityLow, "note"},
		{model.SeverityInfo, "note"},
	}

	for _, tc := range cases {
		t.Run(string(tc.sev), func(t *testing.T) {
			findings := []model.Finding{newFinding("god_object", tc.sev, ".", 1, "test")}
			out := renderSARIF(t, findings, ".")
			runs := out["runs"].([]interface{})
			results := runs[0].(map[string]interface{})["results"].([]interface{})
			level := results[0].(map[string]interface{})["level"]
			if level != tc.level {
				t.Errorf("severity %s: expected level %s, got %v", tc.sev, tc.level, level)
			}
		})
	}
}

func TestSARIFRendererRelativePaths(t *testing.T) {
	findings := []model.Finding{
		newFinding("magic_number", model.SeverityLow, "/workspace/src/main.py", 10, "magic number"),
	}
	out := renderSARIF(t, findings, "/workspace")

	runs := out["runs"].([]interface{})
	results := runs[0].(map[string]interface{})["results"].([]interface{})
	locs := results[0].(map[string]interface{})["locations"].([]interface{})
	physLoc := locs[0].(map[string]interface{})["physicalLocation"].(map[string]interface{})

	uri := physLoc["artifactLocation"].(map[string]interface{})["uri"]
	if uri != "src/main.py" {
		t.Errorf("expected src/main.py, got %v", uri)
	}
	uriBase := physLoc["artifactLocation"].(map[string]interface{})["uriBaseId"]
	if uriBase != "%SRCROOT%" {
		t.Errorf("expected uriBaseId %%SRCROOT%%, got %v", uriBase)
	}
	region := physLoc["region"].(map[string]interface{})
	if region["startLine"] != float64(10) {
		t.Errorf("expected startLine 10, got %v", region["startLine"])
	}
}

func TestSARIFRendererValidSchema(t *testing.T) {
	findings := []model.Finding{
		newFinding("cyclomatic_complexity", model.SeverityHigh, "main.go", 5, "complexity 20"),
	}
	var buf bytes.Buffer
	r := NewSARIF(&buf, "1.2.3", ".")
	if err := r.Render(findings); err != nil {
		t.Fatal(err)
	}
	raw := buf.String()
	if !strings.Contains(raw, "sarif-schema-2.1.0.json") {
		t.Error("expected SARIF schema URL in output")
	}
	if !strings.Contains(raw, `"version": "2.1.0"`) {
		t.Error("expected version 2.1.0 in output")
	}
	if !strings.Contains(raw, `"name": "Klaus-antipatterns-search"`) {
		t.Error("expected tool name in output")
	}
	if !strings.Contains(raw, `"version": "1.2.3"`) {
		t.Error("expected tool version 1.2.3 in output")
	}
}

func TestSARIFRendererRulesPresent(t *testing.T) {
	out := renderSARIF(t, []model.Finding{}, ".")
	runs := out["runs"].([]interface{})
	driver := runs[0].(map[string]interface{})["tool"].(map[string]interface{})["driver"].(map[string]interface{})
	rules := driver["rules"].([]interface{})
	if len(rules) != 6 {
		t.Errorf("expected 6 rules in driver catalog, got %d", len(rules))
	}
}
