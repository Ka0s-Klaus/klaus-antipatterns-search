package renderer

import (
	"encoding/json"
	"io"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/model"
)

// JSONRenderer writes findings as a JSON array to w.
type JSONRenderer struct {
	w io.Writer
}

func NewJSON(w io.Writer) *JSONRenderer {
	return &JSONRenderer{w: w}
}

func (r *JSONRenderer) Render(findings []model.Finding) error {
	if findings == nil {
		findings = []model.Finding{}
	}
	enc := json.NewEncoder(r.w)
	enc.SetIndent("", "  ")
	return enc.Encode(findings)
}
