package renderer

import (
	"encoding/json"
	"io"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/model"
)

// OrgJSONRenderer serializes an OrgReport as indented JSON.
type OrgJSONRenderer struct {
	w io.Writer
}

// NewOrgJSON returns an OrgJSONRenderer writing to w.
func NewOrgJSON(w io.Writer) *OrgJSONRenderer {
	return &OrgJSONRenderer{w: w}
}

// Render encodes the org report as indented JSON.
func (r *OrgJSONRenderer) Render(report *model.OrgReport) error {
	enc := json.NewEncoder(r.w)
	enc.SetIndent("", "  ")
	return enc.Encode(report)
}
