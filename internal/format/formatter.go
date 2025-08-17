// Copyright (c) 2024

package format

import (
	"sort"

	"github.com/asymmetric-effort/mdlint/internal/findings"
)

// Formatter outputs findings in a specific representation.
type Formatter interface {
	Format([]findings.Finding) ([]byte, error)
}

// sortFindings sorts the findings deterministically.
func sortFindings(fs []findings.Finding) {
	sort.Slice(fs, func(i, j int) bool {
		a, b := fs[i], fs[j]
		if a.File != b.File {
			return a.File < b.File
		}
		if a.Line != b.Line {
			return a.Line < b.Line
		}
		if a.Column != b.Column {
			return a.Column < b.Column
		}
		if a.Rule != b.Rule {
			return a.Rule < b.Rule
		}
		return a.Message < b.Message
	})
}
