// Copyright 2024
package mdlint_test

import (
	"testing"

	"mdlint"
)

// TestValidateMarkdown ensures ValidateMarkdown returns nil for valid content
// and an error for empty content.
func TestValidateMarkdown(t *testing.T) {
	t.Parallel()

	if err := mdlint.ValidateMarkdown("# Title"); err != nil {
		t.Fatalf("expected nil, got %v", err)
	}

	if err := mdlint.ValidateMarkdown("   "); err == nil {
		t.Fatalf("expected error for empty content")
	}
}
