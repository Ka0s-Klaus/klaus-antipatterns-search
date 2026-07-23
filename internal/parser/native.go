//go:build native

package parser

// NativeParser uses smacker/go-tree-sitter (CGO) for in-process parsing.
// Full CGO implementation is delivered in Fase 1.
// In Fase 0 this stub delegates to the CLI parser so the binary compiles
// and the architecture is validated end-to-end before adding the CGO dep.
type nativeParser struct {
	delegate *cliDelegate
}

func newImpl() Parser {
	return &nativeParser{delegate: &cliDelegate{}}
}

func (p *nativeParser) ParseFile(path, language string) (*Tree, error) {
	return p.delegate.ParseFile(path, language)
}

func (p *nativeParser) Close() {}

// cliDelegate is a copy of the CLI subprocess logic used as the Fase 0 fallback
// inside the native build. Avoids importing the !native file across build tags.
type cliDelegate struct{}

func (d *cliDelegate) ParseFile(path, language string) (*Tree, error) {
	// Fase 0 stub: CGO binding added in Fase 1.
	// Returns ErrParserUnavailable so callers can degrade gracefully.
	return nil, ErrParserUnavailable
}
