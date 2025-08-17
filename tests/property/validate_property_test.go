// Copyright 2024
package mdlint_test

import (
	"strings"
	"testing"
	"testing/quick"

	"mdlint"
)

// TestValidateMarkdownProperties uses property-based testing to ensure
// ValidateMarkdown matches expectations for empty and non-empty strings.
func TestValidateMarkdownProperties(t *testing.T) {
	t.Parallel()

	property := func(s string) bool {
		err := mdlint.ValidateMarkdown(s)
		if strings.TrimSpace(s) == "" {
			return err != nil
		}
		return err == nil
	}

	if err := quick.Check(property, nil); err != nil {
		t.Fatalf("property test failed: %v", err)
	}
}
