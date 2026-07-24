package detector

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/config"
)

func writeTempGo(t *testing.T, src string) string {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "f.go")
	if err := os.WriteFile(path, []byte(src), 0600); err != nil {
		t.Fatal(err)
	}
	return path
}

func cfg80() *config.Config {
	c := config.Default()
	c.Thresholds.FunctionLOC = 80
	return c
}

func TestLargeFunctionSkipsNonGo(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "f.py")
	if err := os.WriteFile(path, []byte("def foo(): pass"), 0600); err != nil {
		t.Fatal(err)
	}
	findings, err := LargeFunction(path, cfg80())
	if err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings for non-Go file, got %v %v", findings, err)
	}
}

func TestLargeFunctionBelowThreshold(t *testing.T) {
	// A function with 10 lines — well under threshold of 80.
	src := "package p\nfunc Small() {\n" + strings.Repeat("\t_ = 0\n", 10) + "}\n"
	path := writeTempGo(t, src)
	findings, err := LargeFunction(path, cfg80())
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings, got %d", len(findings))
	}
}

func TestLargeFunctionAboveThreshold(t *testing.T) {
	// A function with 85 lines — exceeds threshold of 80.
	src := "package p\nfunc Huge() {\n" + strings.Repeat("\t_ = 0\n", 85) + "}\n"
	path := writeTempGo(t, src)
	findings, err := LargeFunction(path, cfg80())
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Rule != "large_function" {
		t.Errorf("unexpected rule %q", findings[0].Rule)
	}
}

func TestLargeFunctionMethod(t *testing.T) {
	// A method on T with 85 lines — finding message must include "T.BigMethod".
	src := "package p\ntype T struct{}\nfunc (t *T) BigMethod() {\n" +
		strings.Repeat("\t_ = 0\n", 85) + "}\n"
	path := writeTempGo(t, src)
	findings, err := LargeFunction(path, cfg80())
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if !strings.Contains(findings[0].Message, "T.BigMethod") {
		t.Errorf("expected message to contain T.BigMethod, got %q", findings[0].Message)
	}
}

func TestLargeFunctionExternalDecl(t *testing.T) {
	// External (CGO/linkname) function with no body — must not panic or find anything.
	src := "package p\nfunc External() // no body\n"
	path := writeTempGo(t, src)
	findings, err := LargeFunction(path, cfg80())
	if err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings for extern decl, got %v %v", findings, err)
	}
}
