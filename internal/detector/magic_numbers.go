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

// trivialNumbers are values too ubiquitous to flag as magic.
var trivialNumbers = map[string]bool{
	"0": true, "1": true, "2": true,
}

// MagicNumbers detects numeric literals (magic numbers) that appear in non-constant contexts.
// Can be disabled via cfg.MagicNumbers.Enabled. Only flags numbers that appear at least
// cfg.Thresholds.MagicMinCount times in the file. Trivial values (0, 1, 2) are excluded.
// Only processes .go files. Helps identify repeated numeric literals that should be constants.
func MagicNumbers(path string, cfg *config.Config) ([]model.Finding, error) {
	if !cfg.MagicNumbers.Enabled {
		return nil, nil
	}
	if filepath.Ext(path) != ".go" {
		return nil, nil
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, err
	}

	// Build exclusion set: trivial numbers + user-configured extras.
	exclude := make(map[string]bool, len(trivialNumbers)+len(cfg.MagicNumbers.Exclude))
	for k := range trivialNumbers {
		exclude[k] = true
	}
	for _, v := range cfg.MagicNumbers.Exclude {
		exclude[v] = true
	}

	// Mark positions of literals inside const declarations so we skip them.
	constLitPos := make(map[token.Pos]bool)
	ast.Inspect(f, func(n ast.Node) bool {
		gd, ok := n.(*ast.GenDecl)
		if !ok || gd.Tok != token.CONST {
			return true
		}
		ast.Inspect(gd, func(inner ast.Node) bool {
			if lit, ok := inner.(*ast.BasicLit); ok {
				constLitPos[lit.Pos()] = true
			}
			return true
		})
		return false
	})

	// First pass: count occurrences of each candidate number.
	counts := make(map[string]int)
	type occurrence struct {
		value string
		pos   token.Position
	}
	var candidates []occurrence

	ast.Inspect(f, func(n ast.Node) bool {
		lit, ok := n.(*ast.BasicLit)
		if !ok || (lit.Kind != token.INT && lit.Kind != token.FLOAT) {
			return true
		}
		if constLitPos[lit.Pos()] || exclude[lit.Value] {
			return true
		}
		counts[lit.Value]++
		candidates = append(candidates, occurrence{value: lit.Value, pos: fset.Position(lit.Pos())})
		return true
	})

	minCount := cfg.Thresholds.MagicMinCount
	severity := model.SeverityFromString(cfg.Severities.MagicNumbers)
	var findings []model.Finding

	for _, c := range candidates {
		if counts[c.value] >= minCount {
			findings = append(findings, model.Finding{
				Rule:     "magic_number",
				Message:  fmt.Sprintf("magic number %s appears %d times — use a named constant", c.value, counts[c.value]),
				Severity: severity,
				Location: model.Location{File: path, Line: c.pos.Line, Column: c.pos.Column},
			})
		}
	}

	return findings, nil
}
