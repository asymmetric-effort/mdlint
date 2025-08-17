// Copyright 2024 MdLint Authors

package findings

// Finding represents a lint finding.
type Finding struct {
	Rule    string `json:"rule"`
	Message string `json:"message"`
	File    string `json:"file"`
	Line    int    `json:"line"`
	Column  int    `json:"column"`
}
