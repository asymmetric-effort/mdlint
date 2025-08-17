package md1000_test

import (
	"testing"

	"github.com/asymmetric-effort/mdlint/internal/engine"
	"github.com/asymmetric-effort/mdlint/internal/rules/md1000"
)

func TestMD1000Rule_Basic(t *testing.T) {
	rule, ok := engine.GetRule("MD1000")
	if !ok {
		t.Fatalf("MD1000 rule not registered")
	}
	cfg := md1000.Config{LineLength: 10}
	content := "short\nthis line is way too long\n"
	findings := rule.Apply(content, cfg)
	if len(findings) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(findings))
	}
	if findings[0].Line != 2 {
		t.Fatalf("expected finding on line 2, got %d", findings[0].Line)
	}
}

func TestMD1000Rule_CodeBlockOption(t *testing.T) {
	rule, _ := engine.GetRule("MD1000")
	content := "```\nlong line inside code block that should not trigger\n```\n"
	cfg := md1000.Config{LineLength: 10}
	if f := rule.Apply(content, cfg); len(f) != 0 {
		t.Fatalf("expected no findings when code blocks ignored, got %d", len(f))
	}
	cfg.CodeBlocks = true
	if f := rule.Apply(content, cfg); len(f) != 1 {
		t.Fatalf("expected finding when code blocks checked, got %d", len(f))
	}
}

func TestMD1000Rule_TablesOption(t *testing.T) {
	rule, _ := engine.GetRule("MD1000")
	content := "|h1|h2|\n|-|-|\n| longlongline |ok|\n"
	cfg := md1000.Config{LineLength: 10}
	if f := rule.Apply(content, cfg); len(f) != 0 {
		t.Fatalf("expected no findings when tables ignored, got %d", len(f))
	}
	cfg.Tables = true
	if f := rule.Apply(content, cfg); len(f) != 1 {
		t.Fatalf("expected finding when tables checked, got %d", len(f))
	}
}

func TestMD1000Rule_Boundary(t *testing.T) {
	rule, _ := engine.GetRule("MD1000")
	cfg := md1000.Config{LineLength: 10}
	content := "0123456789\n01234567890\n"
	f := rule.Apply(content, cfg)
	if len(f) != 1 {
		t.Fatalf("expected 1 finding, got %d", len(f))
	}
	if f[0].Line != 2 {
		t.Fatalf("expected finding on line 2, got %d", f[0].Line)
	}
}
