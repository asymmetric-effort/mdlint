// Copyright 2024

package engine

// Rule defines a linting rule that can be applied to a node within a file.
// ID must return a unique identifier for the rule. Apply evaluates the rule
// against the provided node and returns any findings.
type Rule interface {
	// ID returns the unique identifier of the rule.
	ID() string
	// Apply evaluates the rule against the given node and context. It returns
	// zero or more findings describing rule violations.
	Apply(node any, ctx *Context) []Finding
}

// Context carries information about the file being processed.
type Context struct {
	// FilePath is the absolute path to the file currently being linted.
	FilePath string
}

// Finding describes a rule violation discovered during linting.
type Finding struct {
	// RuleID is the identifier of the rule that produced the finding.
	RuleID string
	// Location describes where the finding occurred. Typically this is a file
	// path optionally followed by a line number.
	Location string
	// Message explains the nature of the finding.
	Message string
	// Severity indicates the importance of the finding.
	Severity string
}
