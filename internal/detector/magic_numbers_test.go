package detector

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/config"
)

func cfgMagic(enabled bool, minCount int, exclude ...string) *config.Config {
	c := config.Default()
	c.MagicNumbers.Enabled = enabled
	c.Thresholds.MagicMinCount = minCount
	c.MagicNumbers.Exclude = exclude
	return c
}

func TestMagicNumbersDisabled(t *testing.T) {
	src := "package p\nfunc F() { _ = 42 + 42 + 42 }\n"
	path := writeTempGo(t, src)
	findings, err := MagicNumbers(path, cfgMagic(false, 1))
	if err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings when disabled, got %v %v", findings, err)
	}
}

func TestMagicNumbersSkipsNonGo(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "f.js")
	if err := os.WriteFile(path, []byte("var x = 42;"), 0600); err != nil {
		t.Fatal(err)
	}
	findings, err := MagicNumbers(path, cfgMagic(true, 1))
	if err != nil || len(findings) != 0 {
		t.Fatalf("expected no findings for non-Go file, got %v %v", findings, err)
	}
}

func TestMagicNumbersTrivialSkipped(t *testing.T) {
	// 0, 1, 2 are trivial — must not be flagged regardless of count.
	src := "package p\nfunc F() { _ = 0 + 0 + 0 + 1 + 1 + 1 + 2 + 2 + 2 }\n"
	path := writeTempGo(t, src)
	findings, err := MagicNumbers(path, cfgMagic(true, 1))
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings for trivial numbers, got %d", len(findings))
	}
}

func TestMagicNumbersConstSkipped(t *testing.T) {
	// Literals inside const declarations must not be flagged.
	src := "package p\nconst X = 42\nconst Y = 42\nconst Z = 42\n"
	path := writeTempGo(t, src)
	findings, err := MagicNumbers(path, cfgMagic(true, 1))
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings for const literals, got %d", len(findings))
	}
}

func TestMagicNumbersBelowMinCount(t *testing.T) {
	// 42 appears twice, minCount = 3 → no finding.
	src := "package p\nfunc F() { _ = 42 + 42 }\n"
	path := writeTempGo(t, src)
	findings, err := MagicNumbers(path, cfgMagic(true, 3))
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings below min count, got %d", len(findings))
	}
}

func TestMagicNumbersAboveMinCount(t *testing.T) {
	// 42 appears 3 times, minCount = 3 → 3 findings (one per occurrence).
	src := "package p\nfunc F() { _ = 42 + 42 + 42 }\n"
	path := writeTempGo(t, src)
	findings, err := MagicNumbers(path, cfgMagic(true, 3))
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 3 {
		t.Fatalf("expected 3 findings, got %d", len(findings))
	}
	for _, f := range findings {
		if f.Rule != "magic_number" {
			t.Errorf("unexpected rule %q", f.Rule)
		}
	}
}

func TestMagicNumbersUserExclude(t *testing.T) {
	// 100 excluded by user config → no finding even with 3 occurrences.
	src := "package p\nfunc F() { _ = 100 + 100 + 100 }\n"
	path := writeTempGo(t, src)
	findings, err := MagicNumbers(path, cfgMagic(true, 1, "100"))
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) != 0 {
		t.Fatalf("expected 0 findings for user-excluded value, got %d", len(findings))
	}
}
