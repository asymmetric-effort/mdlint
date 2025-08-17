// Tests for MD1100 sequential heading rule.
//
// Copyright (c) 2025 Sam Caldwell
package md1100

import "testing"

func TestCheckSequentialHeadingsValid(t *testing.T) {
	src := []byte("# H1\n## H2\n### H3\n")
	if f := CheckSequentialHeadings(src, Config{}); len(f) != 0 {
		t.Fatalf("expected no findings, got %d", len(f))
	}
}

func TestCheckSequentialHeadingsInvalid(t *testing.T) {
	src := []byte("# H1\n### H3\n")
	if f := CheckSequentialHeadings(src, Config{}); len(f) == 0 {
		t.Fatalf("expected findings for non-sequential headings")
	}
}

func TestCheckSequentialHeadingsExclude(t *testing.T) {
	src := []byte("# Intro\n### Jump\n")
	if f := CheckSequentialHeadings(src, Config{Exclude: []string{"Intro"}}); len(f) != 0 {
		t.Fatalf("expected exclusion to skip findings, got %d", len(f))
	}
}
