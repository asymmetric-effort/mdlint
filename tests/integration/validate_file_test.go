// Copyright 2024
package mdlint_test

import (
	"os"
	"path/filepath"
	"testing"

	"mdlint"
)

// TestValidateMarkdownFromFile reads the project's README and validates it.
func TestValidateMarkdownFromFile(t *testing.T) {
	t.Parallel()

	path := filepath.Join("..", "..", "README.md")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read file: %v", err)
	}

	if err := mdlint.ValidateMarkdown(string(data)); err != nil {
		t.Fatalf("validate markdown: %v", err)
	}
}
