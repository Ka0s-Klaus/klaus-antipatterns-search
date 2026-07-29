package renderer

import (
	"encoding/json"
	"io"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/model"
)

// JSONRenderer outputs the scan result as JSON
type JSONRenderer struct {
	w io.Writer
}

// NewJSONRenderer creates a new JSON renderer
func NewJSONRenderer(w io.Writer) *JSONRenderer {
	return &JSONRenderer{w: w}
}

// Render outputs the scan result as formatted JSON
func (r *JSONRenderer) Render(result *model.ScanResult) error {
	encoder := json.NewEncoder(r.w)
	encoder.SetIndent("", "  ")
	return encoder.Encode(result)
}
