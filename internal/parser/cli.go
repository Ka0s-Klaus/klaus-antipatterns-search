//go:build !native

package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"os/exec"
)

type cliParser struct{}

func newImpl() Parser { return &cliParser{} }

func (p *cliParser) ParseFile(path, language string) (*Tree, error) {
	cmd := exec.Command("tree-sitter", "parse", "--json", path)
	out, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return nil, fmt.Errorf("%w: tree-sitter CLI not installed or parse failed", ErrParserUnavailable)
		}
		return nil, fmt.Errorf("tree-sitter parse: %w", err)
	}
	return parseCLIOutput(out, language, path)
}

func (p *cliParser) Close() {}

// parseCLIOutput converts the JSON output of `tree-sitter parse --json` into a *Tree.
func parseCLIOutput(data []byte, language, path string) (*Tree, error) {
	var raw []json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("parse tree-sitter json: %w", err)
	}
	if len(raw) == 0 {
		return &Tree{Language: language, FilePath: path}, nil
	}
	root, err := parseNode(raw[0])
	if err != nil {
		return nil, err
	}
	return &Tree{Root: root, Language: language, FilePath: path}, nil
}

type rawNode struct {
	Type     string       `json:"type"`
	StartPosition [2]int  `json:"startPosition"`
	EndPosition   [2]int  `json:"endPosition"`
	Children []json.RawMessage `json:"children"`
	Text     string       `json:"text"`
}

func parseNode(data json.RawMessage) (*Node, error) {
	var rn rawNode
	if err := json.Unmarshal(data, &rn); err != nil {
		return nil, err
	}
	n := &Node{
		Type:     rn.Type,
		Text:     rn.Text,
		StartRow: rn.StartPosition[0],
		StartCol: rn.StartPosition[1],
		EndRow:   rn.EndPosition[0],
		EndCol:   rn.EndPosition[1],
	}
	for _, child := range rn.Children {
		cn, err := parseNode(child)
		if err != nil {
			return nil, err
		}
		n.Children = append(n.Children, cn)
	}
	return n, nil
}
