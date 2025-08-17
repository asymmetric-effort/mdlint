// Copyright (c) 2024

package format

import (
	"bytes"
	"fmt"

	"github.com/asymmetric-effort/mdlint/internal/findings"
)

// Text outputs findings in human-readable text.
type Text struct {
	Threshold findings.Severity
}

// NewText creates a Text formatter.
func NewText(threshold findings.Severity) *Text { return &Text{Threshold: threshold} }

// Format implements Formatter.
func (t *Text) Format(fs []findings.Finding) ([]byte, error) {
	// Filter by threshold.
	filtered := make([]findings.Finding, 0, len(fs))
	for _, f := range fs {
		if f.Severity.AtLeast(t.Threshold) {
			filtered = append(filtered, f)
		}
	}
	sortFindings(filtered)
	var buf bytes.Buffer
	for i, f := range filtered {
		if i > 0 {
			buf.WriteByte('\n')
		}
		fmt.Fprintf(&buf, "%s:%d:%d %s[%s] %s", f.File, f.Line, f.Column, f.Rule, f.Severity.String(), f.Message)
	}
	if buf.Len() > 0 {
		buf.WriteByte('\n')
	}
	return buf.Bytes(), nil
}
