package detector

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/config"
)

func cfgGod(threshold int) *config.Config {
	c := config.Default()
	c.Thresholds.GodObject.Methods = threshold
	return c
}

func godSrc(typeName string, methodCount int) string {
	var sb strings.Builder
	sb.WriteString("package p\ntype ")
	sb.WriteString(typeName)
	sb.WriteString(" struct{}\n")
	for i := 0; i < methodCount; i++ {
		fmt.Fprintf(&sb, "func (t *%s) M%d() {}\n", typeName, i)
	}
	return sb.String()
}

func TestGodObjectSkipsNonGo(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "f.py")
	if err := os.WriteFile(path, []byte("class Big: pass"), 0600); err != nil {
		t.Fatal(err)
	}
	findings, err := GodObject(path, cfgGod(20))
	if err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings for non-Go file, got %v %v", findings, err)
	}
}

func TestGodObjectBelowThreshold(t *testing.T) {
	path := writeTempGo(t, godSrc("Small", 5))
	findings, err := GodObject(path, cfgGod(20))
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings, got %d", len(findings))
	}
}

func TestGodObjectAboveThreshold(t *testing.T) {
	path := writeTempGo(t, godSrc("Giant", 25))
	findings, err := GodObject(path, cfgGod(20))
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Rule != "god_object" {
		t.Errorf("unexpected rule %q", findings[0].Rule)
	}
	if !strings.Contains(findings[0].Message, "Giant") {
		t.Errorf("expected message to mention Giant, got %q", findings[0].Message)
	}
}

func TestGodObjectPointerReceiver(t *testing.T) {
	// All methods use pointer receiver — should still count correctly.
	path := writeTempGo(t, godSrc("Ptr", 22))
	findings, err := GodObject(path, cfgGod(20))
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
}

func TestGodObjectMultipleTypes(t *testing.T) {
	// Build a single source with Big (25 methods) and Small (3 methods).
	var sb strings.Builder
	sb.WriteString("package p\n")
	sb.WriteString("type Big struct{}\n")
	for i := 0; i < 25; i++ {
		fmt.Fprintf(&sb, "func (b *Big) M%d() {}\n", i)
	}
	sb.WriteString("type Small struct{}\n")
	for i := 0; i < 3; i++ {
		fmt.Fprintf(&sb, "func (s *Small) S%d() {}\n", i)
	}
	path := writeTempGo(t, sb.String())
	findings, err := GodObject(path, cfgGod(20))
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding (only Big), got %d", len(findings))
	}
	if !strings.Contains(findings[0].Message, "Big") {
		t.Errorf("expected finding for Big, got %q", findings[0].Message)
	}
}
