package detector

import "go/ast"

// receiverTypeName extracts the base type name from a receiver expression,
// handling both value (*T → T) and pointer receivers.
func receiverTypeName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.Ident:
		return t.Name
	case *ast.StarExpr:
		return receiverTypeName(t.X)
	default:
		return ""
	}
}
