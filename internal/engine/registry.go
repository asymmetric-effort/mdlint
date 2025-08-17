// Copyright (c) 2024 The MdLint Authors.
//
// Package engine provides a simple rule registry and finding types.
package engine

// Rule defines the interface for linter rules.
type Rule interface {
	// ID returns the rule identifier, e.g. "MD1000".
	ID() string
	// Apply evaluates the rule against the provided Markdown content using
	// the supplied configuration. The configuration type is rule specific.
	Apply(content string, cfg any) []Finding
}

// Finding represents a single rule violation.
type Finding struct {
	Rule    string // Rule identifier.
	Line    int    // 1-indexed line number.
	Column  int    // 1-indexed column number.
	Message string // Human-readable description of the issue.
}

var registry = map[string]Rule{}

// RegisterRule adds a rule implementation to the registry.
func RegisterRule(r Rule) {
	registry[r.ID()] = r
}

// GetRule returns a registered rule by ID.
func GetRule(id string) (Rule, bool) {
	r, ok := registry[id]
	return r, ok
}

// Rules returns all registered rules.
func Rules() []Rule {
	rs := make([]Rule, 0, len(registry))
	for _, r := range registry {
		rs = append(rs, r)
	}
	return rs
}
