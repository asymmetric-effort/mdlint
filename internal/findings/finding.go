// Copyright (c) 2024

package findings

// Finding represents a single rule violation.
type Finding struct {
	Rule     string   `json:"rule"`
	Severity Severity `json:"severity"`
	Message  string   `json:"message"`
	File     string   `json:"file"`
	Line     int      `json:"line"`
	Column   int      `json:"column"`
}
