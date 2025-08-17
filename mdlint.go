// Copyright 2024
// Package mdlint provides utilities for validating Markdown content.
package mdlint

import (
	"errors"
	"strings"
)

// ValidateMarkdown checks that the provided Markdown content is non-empty.
// It returns an error when the content is blank and nil otherwise.
func ValidateMarkdown(md string) error {
	if strings.TrimSpace(md) == "" {
		return errors.New("markdown content is empty")
	}
	return nil
}
