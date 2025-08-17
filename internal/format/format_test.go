// Copyright (c) 2024

package format_test

import (
	"os"
	"testing"

	"github.com/asymmetric-effort/mdlint/internal/findings"
	formatpkg "github.com/asymmetric-effort/mdlint/internal/format"
)

var testFindings = []findings.Finding{
	{File: "b.md", Line: 2, Column: 1, Rule: "MD1000", Severity: findings.Warning, Message: "warning 1"},
	{File: "a.md", Line: 1, Column: 5, Rule: "MD1000", Severity: findings.Suggestion, Message: "suggestion 1"},
	{File: "a.md", Line: 1, Column: 2, Rule: "MD1000", Severity: findings.Error, Message: "error 1"},
	{File: "a.md", Line: 1, Column: 3, Rule: "MD1001", Severity: findings.Warning, Message: "warning 2"},
	{File: "a.md", Line: 1, Column: 3, Rule: "MD1000", Severity: findings.Warning, Message: "warning 3"},
	{File: "a.md", Line: 2, Column: 1, Rule: "MD1002", Severity: findings.Error, Message: "error 2"},
}

func TestFormatters(t *testing.T) {
	t.Run("text_warning", func(t *testing.T) {
		f := formatpkg.NewText(findings.Warning)
		out, err := f.Format(testFindings)
		if err != nil {
			t.Fatalf("format: %v", err)
		}
		assertGolden(t, "testdata/text_warning.golden", out)
	})

	t.Run("text_suggestion", func(t *testing.T) {
		f := formatpkg.NewText(findings.Suggestion)
		out, err := f.Format(testFindings)
		if err != nil {
			t.Fatalf("format: %v", err)
		}
		assertGolden(t, "testdata/text_suggestion.golden", out)
	})

	t.Run("json", func(t *testing.T) {
		f := formatpkg.NewJSON()
		out, err := f.Format(testFindings)
		if err != nil {
			t.Fatalf("format: %v", err)
		}
		assertGolden(t, "testdata/findings.json", out)
	})
}

func assertGolden(t *testing.T, path string, got []byte) {
	t.Helper()
	want, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden: %v", err)
		return
	}
	if string(want) != string(got) {
		t.Fatalf("unexpected output:\nwant:\n%s\n---\ngot:\n%s", string(want), string(got))
	}
}
