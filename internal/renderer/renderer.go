package renderer

import "github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/model"

// Renderer writes a slice of findings to an output destination.
type Renderer interface {
	Render(findings []model.Finding) error
}
