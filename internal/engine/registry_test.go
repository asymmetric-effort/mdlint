package engine

import (
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// testRule is a simple rule used for testing registration and execution.
type testRule struct{}

func (testRule) ID() string { return "test-rule" }

func (testRule) Apply(node any, ctx *Context) []Finding {
	return []Finding{{
		RuleID:   "test-rule",
		Location: ctx.FilePath,
		Message:  "ok",
		Severity: "info",
	}}
}

func init() { Register(testRule{}) }

// TestRules ensures that rules are registered via init and retrieved in a
// deterministic order.
func TestRules(t *testing.T) {
	rules := Rules()
	if len(rules) != 1 || rules[0].ID() != "test-rule" {
		t.Fatalf("unexpected rules: %#v", rules)
	}
}

// TestRun verifies that Run respects include/exclude patterns and returns
// deterministic results.
func TestRun(t *testing.T) {
	dir := t.TempDir()
	mustWrite(t, filepath.Join(dir, "a.md"), "alpha")
	mustWrite(t, filepath.Join(dir, "b.md"), "beta")
	mustWrite(t, filepath.Join(dir, "c.txt"), "gamma")

	cfg := Config{Include: []string{"*.md"}, Exclude: []string{"b.md"}, Workers: 2}

	got1, err := Run(context.Background(), dir, cfg)
	if err != nil {
		t.Fatalf("run 1: %v", err)
	}
	got2, err := Run(context.Background(), dir, cfg)
	if err != nil {
		t.Fatalf("run 2: %v", err)
	}
	if !reflect.DeepEqual(got1, got2) {
		t.Fatalf("non-deterministic results: %v vs %v", got1, got2)
	}
	if len(got1) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(got1))
	}
	if got1[0].Location != filepath.Join(dir, "a.md") {
		t.Fatalf("unexpected finding location: %s", got1[0].Location)
	}
}

func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write %s: %v", path, err)
	}
}
