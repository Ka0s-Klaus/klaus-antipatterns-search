package detector

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/config"
	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/model"
)

// LargeFunction detects Go functions and methods whose bodies exceed the configured threshold.
// It flags functions with more lines than cfg.Thresholds.FunctionLOC, reporting the function
// name, actual line count, and the configured limit. Only processes .go files.
func LargeFunction(path string, cfg *config.Config) ([]model.Finding, error) {
	if filepath.Ext(path) != ".go" {
		return nil, nil
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, err
	}

	threshold := cfg.Thresholds.FunctionLOC
	severity := model.SeverityFromString(cfg.Severities.LargeFunction)
	var findings []model.Finding

	ast.Inspect(f, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Body == nil {
			return true
		}

		start := fset.Position(fn.Body.Lbrace)
		end := fset.Position(fn.Body.Rbrace)
		loc := end.Line - start.Line + 1

		if loc > threshold {
			name := fn.Name.Name
			if fn.Recv != nil && len(fn.Recv.List) > 0 {
				name = receiverTypeName(fn.Recv.List[0].Type) + "." + name
			}
			findings = append(findings, model.Finding{
				Rule:     "large_function",
				Message:  fmt.Sprintf("%q has %d lines (threshold: %d)", name, loc, threshold),
				Severity: severity,
				Location: model.Location{File: path, Line: start.Line, Column: start.Column},
			})
		}
		return true
	})

	return findings, nil
}
