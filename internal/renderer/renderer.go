package renderer

import (
	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/model"
)

// Renderer defines the interface for output renderers
type Renderer interface {
	Render(result *model.ScanResult) error
}
