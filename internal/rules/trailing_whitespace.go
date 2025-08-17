// Copyright 2024 the mdlint authors.

package rules

import (
	"strings"
)

// Finding represents a rule violation.
type Finding struct {
	Line    int
	Column  int
	RuleID  string
	Message string
}

// TrailingWhitespaceConfig configures the trailing whitespace rule.
type TrailingWhitespaceConfig struct {
	// IgnoreCodeBlocks controls whether lines inside fenced code blocks are checked.
	// When true, lines within fenced code blocks are skipped.
	IgnoreCodeBlocks bool
}

// CheckTrailingWhitespace scans content for trailing spaces or tabs.
// It returns a list of findings for lines containing trailing whitespace.
func CheckTrailingWhitespace(content []byte, cfg TrailingWhitespaceConfig) []Finding {
	lines := strings.Split(string(content), "\n")
	var findings []Finding
	inCode := false
	var fenceChar rune
	for i, line := range lines {
		// Remove trailing carriage return for Windows line endings.
		line = strings.TrimSuffix(line, "\r")

		trimmedLeft := strings.TrimLeft(line, " \t")
		if f, ok := fenceMarker(trimmedLeft); ok {
			if inCode {
				if f == fenceChar {
					inCode = false
				}
			} else {
				inCode = true
				fenceChar = f
			}
			if cfg.IgnoreCodeBlocks {
				continue
			}
		}
		if inCode && cfg.IgnoreCodeBlocks {
			continue
		}
		if len(line) > 0 {
			last := line[len(line)-1]
			if last == ' ' || last == '\t' {
				findings = append(findings, Finding{
					Line:    i + 1,
					Column:  len(line),
					RuleID:  "MD1800",
					Message: "line has trailing whitespace",
				})
			}
		}
	}
	return findings
}

func fenceMarker(line string) (rune, bool) {
	if strings.HasPrefix(line, "```") {
		return '`', true
	}
	if strings.HasPrefix(line, "~~~") {
		return '~', true
	}
	return 0, false
}
