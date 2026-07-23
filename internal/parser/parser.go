package parser

import "errors"

// ErrParserUnavailable is returned when the underlying parser engine is not installed.
var ErrParserUnavailable = errors.New("parser unavailable")

// Node represents a single node in the syntax tree.
type Node struct {
	Type     string
	Text     string
	StartRow int
	StartCol int
	EndRow   int
	EndCol   int
	Children []*Node
}

// Tree is the parsed representation of a single source file.
type Tree struct {
	Root     *Node
	Language string
	FilePath string
}

// Parser produces a syntax Tree for a given source file.
// Implementations are selected at build time via Go build tags:
//   - default (no tags): CLIParser — subprocess, pure Go, zero CGO
//   - -tags native:      NativeParser — smacker/go-tree-sitter, CGO, no runtime deps
type Parser interface {
	// ParseFile returns the syntax tree for path in the given language.
	// Returns ErrParserUnavailable if the underlying engine is not installed.
	ParseFile(path, language string) (*Tree, error)
	Close()
}

// New returns the Parser implementation selected at compile time.
// Defined in cli.go (default) or native.go (-tags native).
func New() Parser { return newImpl() }
