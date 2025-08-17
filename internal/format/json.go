// Copyright (c) 2024

package format

import (
	"encoding/json"

	"github.com/asymmetric-effort/mdlint/internal/findings"
)

// JSON outputs findings as JSON array.
type JSON struct{}

// NewJSON creates a JSON formatter.
func NewJSON() *JSON { return &JSON{} }

// Format implements Formatter.
func (j *JSON) Format(fs []findings.Finding) ([]byte, error) {
	dup := make([]findings.Finding, len(fs))
	copy(dup, fs)
	sortFindings(dup)
	b, err := json.MarshalIndent(dup, "", "  ")
	if err != nil {
		return nil, err
	}
	return append(b, '\n'), nil
}
