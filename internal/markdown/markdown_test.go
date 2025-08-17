// Copyright (c) 2024 Asymmetric Effort

package markdown

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yuin/goldmark/text"
)

// helper to load testdata
func readFixture(t *testing.T, name string) []byte {
	t.Helper()
	path := filepath.Join("testdata", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("failed to read fixture %s: %v", name, err)
	}
	return data
}

func TestFrontMatterRange(t *testing.T) {
	src := readFixture(t, "basic.md")
	rng, ok := FrontMatterRange(src)
	if !ok {
		t.Fatalf("expected front-matter range")
	}
	if rng.StartLine != 1 || rng.EndLine != 3 {
		t.Fatalf("unexpected range: %#v", rng)
	}
}

func TestCodeBlockRanges(t *testing.T) {
	src := readFixture(t, "basic.md")
	ranges := CodeBlockRanges(src)
	if len(ranges) != 1 {
		t.Fatalf("expected one code block, got %d", len(ranges))
	}
	if ranges[0].StartLine != 9 || ranges[0].EndLine != 11 {
		t.Fatalf("unexpected range: %#v", ranges[0])
	}
}

func TestParserSupportsExtensions(t *testing.T) {
	src := readFixture(t, "basic.md")
	md := Parser()
	if md == nil {
		t.Fatalf("Parser returned nil")
	}
	// parsing should succeed without panic
	reader := text.NewReader(src)
	if md.Parser().Parse(reader) == nil {
		t.Fatalf("expected parsed document")
	}
}
