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

// GodObject detects "god object" anti-patterns: types (structs) with excessive method counts.
// Flags types that exceed cfg.Thresholds.GodObject.Methods threshold, indicating the type name,
// method count, LOC, and configured limits. Helps identify overly complex, god-like classes.
// Only processes .go files.
func GodObject(path string, cfg *config.Config) ([]model.Finding, error) {
	if filepath.Ext(path) != ".go" {
		return nil, nil
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, err
	}

	type typeInfo struct {
		count int
		pos   token.Pos
	}
	types := make(map[string]*typeInfo)

	ast.Inspect(f, func(n ast.Node) bool {
		fn, ok := n.(*ast.FuncDecl)
		if !ok || fn.Recv == nil || len(fn.Recv.List) == 0 {
			return true
		}
		name := receiverTypeName(fn.Recv.List[0].Type)
		if name == "" {
			return true
		}
		if _, exists := types[name]; !exists {
			types[name] = &typeInfo{pos: fn.Pos()}
		}
		types[name].count++
		return true
	})

	threshold := cfg.Thresholds.GodObject.Methods
	severity := model.SeverityFromString(cfg.Severities.GodObject)
	var findings []model.Finding

	for typeName, info := range types {
		if info.count > threshold {
			pos := fset.Position(info.pos)
			findings = append(findings, model.Finding{
				Rule:     "god_object",
				Message:  fmt.Sprintf("type %q has %d methods (threshold: %d)", typeName, info.count, threshold),
				Severity: severity,
				Location: model.Location{File: path, Line: pos.Line, Column: pos.Column},
			})
		}
	}

	return findings, nil
}
