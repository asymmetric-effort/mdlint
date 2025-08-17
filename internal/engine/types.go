// Copyright 2024

package engine

// Rule defines a linting rule applied to Markdown content.
// ID must return a unique identifier for the rule. Apply evaluates the rule
// against the provided content using optional configuration and returns any findings.
type Rule interface {
	// ID returns the unique identifier of the rule.
	ID() string
	// Apply evaluates the rule against the given content and optional configuration.
	// It returns zero or more findings describing rule violations.
	Apply(content string, cfg any) []Finding
}

// Finding describes a rule violation discovered during linting.
type Finding struct {
	// Rule is the identifier of the rule that produced the finding.
	Rule string
	// Line is the 1-based line number of the violation.
	Line int
	// Column is the 1-based column number of the violation.
	Column int
	// Message explains the nature of the finding.
	Message string
}
