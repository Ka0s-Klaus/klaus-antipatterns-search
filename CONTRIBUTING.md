# 🤝 Contributing to Klaus antipatterns-search

We welcome contributions! This guide shows how to add new detectors, adapters, tests, and fixes.

---

## 📋 Workflow

All work follows the **Issue → Branch → PR** pattern. Never commit directly to `main`.

1. **Open an Issue** with full context (what, why, expected behavior)
2. **Create a branch** from `main`: `git checkout -b GH-N-short-description`
3. **Make your changes** (code, tests, docs)
4. **Open a PR** against `main` and link the issue
5. **Code review** and CI checks pass
6. **Merge** to `main`

---

## 🔍 Adding a Native Detector

### 1. Create the detector file

Add a new file in `internal/detector/` (e.g., `internal/detector/unused_vars.go`):

```go
package detector

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/config"
	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/model"
)

// UnusedVars detects local variables that are assigned but never read.
// Only processes .go files; works via go/ast inspection.
func UnusedVars(path string, cfg *config.Config) ([]model.Finding, error) {
	if filepath.Ext(path) != ".go" {
		return nil, nil
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, err
	}

	var findings []model.Finding

	// Example: inspect AST for unused variables
	ast.Inspect(f, func(n ast.Node) bool {
		// TODO: Implement your logic here
		return true
	})

	return findings, nil
}
```

**Key patterns:**
- Use `go/ast` (stdlib) — no external AST libraries needed
- Return `nil, nil` for non-.go files (not an error)
- Return `[]model.Finding{}` for zero findings (not nil)
- Severity from config: `model.SeverityFromString(cfg.Severities.YourRuleHere)`
- Location: file path, line number, and optional column

### 2. Register in Scanner

Edit `internal/scanner/scanner.go` and add your detector to the slice:

```go
detectors: []namedDetector{
	{name: "large_function", fn: detector.LargeFunction},
	{name: "god_object", fn: detector.GodObject},
	{name: "magic_numbers", fn: detector.MagicNumbers},
	{name: "unused_vars", fn: detector.UnusedVars},  // NEW
}
```

### 3. Add tests

Create `internal/detector/unused_vars_test.go`:

```go
package detector

import (
	"path/filepath"
	"testing"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/config"
)

func TestUnusedVars(t *testing.T) {
	cfg := config.Default()
	
	// Use testdata fixtures
	findings, err := UnusedVars("testdata/unused_vars_example.go", cfg)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Assert findings count, rule, message, etc.
	if len(findings) == 0 {
		t.Error("expected findings, got none")
	}
	if findings[0].Rule != "unused_vars" {
		t.Errorf("rule: want unused_vars, got %s", findings[0].Rule)
	}
}

func TestUnusedVarsNonGoFile(t *testing.T) {
	cfg := config.Default()
	findings, err := UnusedVars("example.txt", cfg)
	if err != nil || len(findings) != 0 {
		t.Error("should skip non-.go files")
	}
}
```

Create test fixtures in `internal/detector/testdata/`:
```bash
mkdir -p internal/detector/testdata
```

Then add `internal/detector/testdata/unused_vars_example.go` with Go code that has unused variables.

### 4. Add config (if needed)

If your detector needs thresholds, add them to `internal/config/config.go`:

```go
type ThresholdsConfig struct {
	// ... existing fields ...
	UnusedVars int `yaml:"unused_vars"`  // NEW
}
```

And defaults in `Default()`:
```go
Thresholds: ThresholdsConfig{
	// ... existing ...
	UnusedVars: 5,  // Flag if > 5 unused vars
}
```

And add severity in `SeveritiesConfig` and defaults.

### 5. Document the detector

Add an entry to `docs/fase-1-detectores-nativos.md` (or create one if it doesn't exist):

```markdown
## UnusedVars

**What?** Detects local variables that are assigned but never read.

**Config:**
```yaml
thresholds:
  unused_vars: 5
severities:
  unused_vars: low
```

**Why?** Dead code wastes maintenance burden and confuses readers.
```

---

## 🔌 Adding an OSS Adapter

### 1. Create the adapter file

Add a new file in `internal/adapter/` (e.g., `internal/adapter/eslint.go`):

```go
package adapter

import (
	"context"
	"encoding/json"
	"os/exec"
	"time"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/config"
	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/model"
)

// Eslint adapts eslint (JavaScript linter) for style/lint findings.
// Returns ErrToolNotFound if eslint is not installed.
func Eslint(root string, cfg *config.Config) ([]model.Finding, error) {
	eslintPath, err := exec.LookPath("eslint")
	if err != nil {
		return nil, ErrToolNotFound
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, eslintPath, "--format", "json", root)
	output, err := cmd.Output()
	if err != nil {
		// eslint returns non-zero exit if findings exist; parse anyway
		if exitErr, ok := err.(*exec.ExitError); !ok {
			return nil, nil // Truly failed (not found, permission denied, etc)
		} else if exitErr.ExitCode() != 1 {
			return nil, nil // Unexpected exit code
		}
	}

	// Parse output
	var findings []model.Finding
	// TODO: implement parsing logic

	return findings, nil
}
```

**Key patterns:**
- Always check `exec.LookPath` first
- Return `ErrToolNotFound` if tool not found (enables graceful degradation)
- Use context with timeout to prevent hangs
- Return `nil, nil` on fatal errors (non-tool-not-found)
- Parse tool output and normalize to `model.Finding`

### 2. Register in Scanner

Edit `internal/scanner/scanner.go`:

```go
dirAdapters: []namedAdapter{
	{name: "jscpd", fn: adapter.Jscpd},
	{name: "madge", fn: adapter.Madge},
	{name: "radon", fn: adapter.Radon},
	{name: "gocyclo", fn: adapter.Gocyclo},
	{name: "eslint", fn: adapter.Eslint},  // NEW
}
```

### 3. Add tests

Create `internal/adapter/eslint_test.go`:

```go
package adapter

import (
	"testing"

	"github.com/Ka0s-Klaus/Klaus-antipatterns-search/internal/config"
)

func TestEslintNotInstalled(t *testing.T) {
	cfg := config.Default()
	findings, err := Eslint("/tmp", cfg)
	
	// If eslint truly not installed, expect ErrToolNotFound
	if err == ErrToolNotFound {
		t.Skip("eslint not installed (expected)")
	}
	
	// If installed, should return either findings or nil on error
	if err != nil && err != ErrToolNotFound {
		t.Errorf("unexpected error: %v", err)
	}
}
```

### 4. Document the adapter

Add to `docs/fase-2-adaptadores-oss.md`:

```markdown
## Eslint

**What?** JavaScript linter findings (style, best practices, errors).

**Install:**
```bash
npm install -g eslint
```

**Why?** Detects style issues and bugs in JavaScript/TypeScript codebases.

**Graceful degradation:** If eslint is not installed, scanning continues without error.
```

---

## 📝 Adding a Test

### For detectors/adapters:

1. Create testdata fixtures in `internal/{detector,adapter}/testdata/`
2. Write tests in `internal/{detector,adapter}/*_test.go`
3. Use `config.Default()` for standard thresholds
4. Test edge cases: empty files, missing tools, malformed input

### Run tests locally:

```bash
go test ./internal/detector -v
go test ./internal/adapter -v
go test ./... -race       # All tests with race detector
```

---

## 🐛 Reporting Bugs

When opening an issue for a bug:

1. **Title:** One-line description (e.g., "large_function reports wrong line numbers")
2. **Context:** What were you doing? What did you expect?
3. **Reproduction:** Minimal example (code snippet, repo URL, or testdata file)
4. **Output:** Show actual vs. expected findings
5. **Environment:** Go version, OS, tool versions (for adapters)

Example:
```
Title: jscpd adapter fails with permission denied

## Context
Running `antipatterns scan ./` on a repo with restricted vendor/ files.

## Repro
```bash
antipatterns scan ./example-repo --verbose
```

## Expected
Gracefully skip vendor/ or report jscpd skipped.

## Actual
Error: "jscpd: permission denied vendor/..."
```

---

## 📦 PR Checklist

Before opening a PR:

- [ ] Branch named `GH-N-short-description`
- [ ] Issue linked (in PR description)
- [ ] Tests added and pass locally (`go test ./...`)
- [ ] Code follows existing patterns (scanning, error handling, logging)
- [ ] Godocs added for public types/functions
- [ ] No `log.Fatal()` or `panic()` (use errors)
- [ ] Testdata committed (fixtures in `testdata/` dirs)
- [ ] Docs updated (fase-X.md or new doc)

---

## 🎯 Architecture Overview

### Plugin Pattern

Detectors and adapters are registered as slices of named functions in `Scanner`:

```go
detectors []namedDetector       // File-level native detectors
dirAdapters []namedAdapter      // Directory-level OSS adapters
```

Adding a new detector/adapter is just adding a function matching the signature:
- `detectorFunc(path, cfg) (findings, error)`
- `dirAdapterFunc(root, cfg) (findings, error)`

### Error Handling

- **Fatal errors** (file unreadable, parse error): return error, skip file/dir
- **Tool not found** (OSS adapter missing): return `ErrToolNotFound` (graceful degradation)
- **Tool fails** (non-zero exit, timeout): return `nil, nil` (log warning in verbose, continue)

### Config & Thresholds

All thresholds live in `internal/config/config.go`:
```go
type ThresholdsConfig struct {
	FunctionLOC int  // Flag functions > N lines
	Cyclomatic  int  // Flag functions > N complexity
	// Add yours here
}
```

Defaults in `Default()`, can be overridden by `.antipatterns.yml`.

---

## 🔗 Links

- **Issues:** https://github.com/Ka0s-Klaus/Klaus-antipatterns-search/issues
- **Project:** https://github.com/orgs/Ka0s-Klaus/projects/21
- **Roadmap:** [README.md](README.md#-roadmap)
- **Docs:** `docs/fase-*.md`

---

## Questions?

Open an issue or ask in a PR comment. We'll help!

Happy coding! 🚀
