package renderer

import (
	"encoding/json"
	"io"
	"path/filepath"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/model"
)

const (
	sarifSchema      = "https://raw.githubusercontent.com/oasis-tcs/sarif-spec/master/Schemata/sarif-schema-2.1.0.json"
	sarifSpecVersion = "2.1.0"
	toolName         = "Klaus-antipatterns-search"
	toolInfoURI      = "https://github.com/Ka0s-Klaus/Klaus-antipatterns-search"
)

// sarifReport is the top-level SARIF 2.1.0 object.
type sarifReport struct {
	Schema  string     `json:"$schema"`
	Version string     `json:"version"`
	Runs    []sarifRun `json:"runs"`
}

type sarifRun struct {
	Tool    sarifTool     `json:"tool"`
	Results []sarifResult `json:"results"`
}

type sarifTool struct {
	Driver sarifDriver `json:"driver"`
}

type sarifDriver struct {
	Name           string                      `json:"name"`
	Version        string                      `json:"version"`
	InformationURI string                      `json:"informationUri"`
	Rules          []sarifReportingDescriptor  `json:"rules"`
}

type sarifReportingDescriptor struct {
	ID               string           `json:"id"`
	Name             string           `json:"name"`
	ShortDescription sarifMessage     `json:"shortDescription"`
	Properties       sarifProperties  `json:"properties,omitempty"`
}

type sarifMessage struct {
	Text string `json:"text"`
}

type sarifProperties struct {
	Tags []string `json:"tags,omitempty"`
}

type sarifResult struct {
	RuleID    string          `json:"ruleId"`
	Level     string          `json:"level"`
	Message   sarifMessage    `json:"message"`
	Locations []sarifLocation `json:"locations"`
}

type sarifLocation struct {
	PhysicalLocation sarifPhysicalLocation `json:"physicalLocation"`
}

type sarifPhysicalLocation struct {
	ArtifactLocation sarifArtifactLocation `json:"artifactLocation"`
	Region           *sarifRegion          `json:"region,omitempty"`
}

type sarifArtifactLocation struct {
	URI       string `json:"uri"`
	URIBaseID string `json:"uriBaseId,omitempty"`
}

type sarifRegion struct {
	StartLine   int `json:"startLine,omitempty"`
	StartColumn int `json:"startColumn,omitempty"`
}

// knownRules is the static catalog of rules emitted by Klaus-antipatterns-search.
// GitHub Code Scanning uses this to render rule metadata in the Security tab.
var knownRules = []sarifReportingDescriptor{
	{
		ID: "large_function", Name: "LargeFunction",
		ShortDescription: sarifMessage{Text: "Function or method exceeds the configured LOC threshold"},
		Properties:       sarifProperties{Tags: []string{"maintainability"}},
	},
	{
		ID: "god_object", Name: "GodObject",
		ShortDescription: sarifMessage{Text: "Class or file has too many methods or excessive lines of code"},
		Properties:       sarifProperties{Tags: []string{"maintainability", "solid"}},
	},
	{
		ID: "magic_number", Name: "MagicNumber",
		ShortDescription: sarifMessage{Text: "Numeric literal used directly instead of a named constant"},
		Properties:       sarifProperties{Tags: []string{"maintainability", "readability"}},
	},
	{
		ID: "duplication", Name: "Duplication",
		ShortDescription: sarifMessage{Text: "Code clone detected (type-1 or type-2 duplication)"},
		Properties:       sarifProperties{Tags: []string{"maintainability", "duplication"}},
	},
	{
		ID: "circular_dependency", Name: "CircularDependency",
		ShortDescription: sarifMessage{Text: "Circular import or module dependency detected"},
		Properties:       sarifProperties{Tags: []string{"architecture", "coupling"}},
	},
	{
		ID: "cyclomatic_complexity", Name: "CyclomaticComplexity",
		ShortDescription: sarifMessage{Text: "Function cyclomatic complexity exceeds the configured threshold"},
		Properties:       sarifProperties{Tags: []string{"maintainability", "complexity"}},
	},
}

// SARIFRenderer writes findings as a SARIF 2.1.0 JSON document to w.
// root is used to compute paths relative to the repository root (%SRCROOT%).
type SARIFRenderer struct {
	w       io.Writer
	version string
	root    string
}

func NewSARIF(w io.Writer, version, root string) *SARIFRenderer {
	return &SARIFRenderer{w: w, version: version, root: root}
}

func (r *SARIFRenderer) Render(findings []model.Finding) error {
	if findings == nil {
		findings = []model.Finding{}
	}

	results := make([]sarifResult, 0, len(findings))
	for _, f := range findings {
		results = append(results, r.toSARIFResult(f))
	}

	report := sarifReport{
		Schema:  sarifSchema,
		Version: sarifSpecVersion,
		Runs: []sarifRun{
			{
				Tool: sarifTool{
					Driver: sarifDriver{
						Name:           toolName,
						Version:        r.version,
						InformationURI: toolInfoURI,
						Rules:          knownRules,
					},
				},
				Results: results,
			},
		},
	}

	enc := json.NewEncoder(r.w)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}

func (r *SARIFRenderer) toSARIFResult(f model.Finding) sarifResult {
	uri := toRelativeURI(f.Location.File, r.root)

	loc := sarifLocation{
		PhysicalLocation: sarifPhysicalLocation{
			ArtifactLocation: sarifArtifactLocation{
				URI:       uri,
				URIBaseID: "%SRCROOT%",
			},
		},
	}

	if f.Location.Line > 0 {
		region := &sarifRegion{StartLine: f.Location.Line}
		if f.Location.Column > 0 {
			region.StartColumn = f.Location.Column
		}
		loc.PhysicalLocation.Region = region
	}

	return sarifResult{
		RuleID:    f.Rule,
		Level:     severityToSARIFLevel(f.Severity),
		Message:   sarifMessage{Text: f.Message},
		Locations: []sarifLocation{loc},
	}
}

// severityToSARIFLevel maps internal severity to SARIF result level.
// SARIF levels: error, warning, note, none.
func severityToSARIFLevel(s model.Severity) string {
	switch s {
	case model.SeverityCritical, model.SeverityHigh:
		return "error"
	case model.SeverityMedium:
		return "warning"
	default:
		return "note"
	}
}

// toRelativeURI makes path relative to root and converts to forward-slash URI form.
func toRelativeURI(path, root string) string {
	rel, err := filepath.Rel(root, path)
	if err != nil {
		return filepath.ToSlash(path)
	}
	return filepath.ToSlash(rel)
}
