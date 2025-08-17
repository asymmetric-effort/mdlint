// Copyright (c) 2024 The MdLint Authors.
//
// Package md1000 implements a rule enforcing maximum line length in Markdown files.
package md1000

import (
	"fmt"
	"regexp"
	"strings"
	"unicode/utf8"

	"github.com/asymmetric-effort/mdlint/internal/engine"
)

// Config configures the MD1000 rule.
type Config struct {
	LineLength int  // Maximum allowed line length. Defaults to 80.
	CodeBlocks bool // Whether to enforce inside fenced code blocks.
	Tables     bool // Whether to enforce inside tables.
}

// Rule implements the MD1000 maximum line length rule.
type Rule struct{}

// ensure Rule satisfies engine.Rule.
var _ engine.Rule = Rule{}

const defaultLineLength = 80

// init registers the rule in the engine registry.
func init() {
	engine.RegisterRule(Rule{})
}

// ID returns the rule identifier.
func (Rule) ID() string { return "MD1000" }

// Apply checks the supplied Markdown content against the configured maximum
// line length and returns any findings.
func (Rule) Apply(content string, cfg any) []engine.Finding {
	opts := Config{LineLength: defaultLineLength}
	if c, ok := cfg.(Config); ok {
		if c.LineLength > 0 {
			opts.LineLength = c.LineLength
		}
		opts.CodeBlocks = c.CodeBlocks
		opts.Tables = c.Tables
	}

	lines := strings.Split(content, "\n")
	findings := []engine.Finding{}

	inCode := false
	inTable := false
	for i := 0; i < len(lines); i++ {
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		// Toggle fenced code blocks.
		if strings.HasPrefix(trimmed, "```") || strings.HasPrefix(trimmed, "~~~") {
			inCode = !inCode
			// Fences themselves are ignored.
			continue
		}

		// Detect start of a table.
		if !inCode && isTableHeader(i, lines) {
			inTable = true
		}
		// Detect end of a table.
		if inTable && (trimmed == "" || !strings.Contains(line, "|")) {
			inTable = false
		}

		if (!opts.CodeBlocks && inCode) || (!opts.Tables && inTable) {
			continue
		}

		if utf8.RuneCountInString(line) > opts.LineLength {
			findings = append(findings, engine.Finding{
				Rule:    "MD1000",
				Line:    i + 1,
				Column:  opts.LineLength + 1,
				Message: fmt.Sprintf("Line exceeds maximum length of %d characters", opts.LineLength),
			})
		}
	}

	return findings
}

var tableSepRE = regexp.MustCompile(`^\s*\|?\s*[:\-]+[-\s:|]*\|`)

// isTableHeader reports whether the line at index i begins a Markdown table.
func isTableHeader(i int, lines []string) bool {
	line := lines[i]
	if !strings.Contains(line, "|") {
		return false
	}
	if i+1 >= len(lines) {
		return false
	}
	next := strings.TrimSpace(lines[i+1])
	return tableSepRE.MatchString(next)
}
