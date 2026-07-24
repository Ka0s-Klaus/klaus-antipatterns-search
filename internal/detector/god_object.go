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

// GodObject flags any Go type whose method count exceeds cfg.Thresholds.GodObject.Methods.
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
