// Package md1100 provides a rule checking sequential heading levels in Markdown.
//
// Copyright (c) 2025 Sam Caldwell
package md1100

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// Config configures the MD1100 rule.
type Config struct {
	// Exclude lists section headings to ignore during validation.
	Exclude []string
}

// Finding represents a heading level violation.
type Finding struct {
	// Line is the line number of the offending heading.
	Line int
	// Message describes the violation.
	Message string
}

// CheckSequentialHeadings scans the Markdown document and reports headings that
// increase by more than one level at a time. Sections whose heading text matches
// cfg.Exclude are skipped entirely.
func CheckSequentialHeadings(src []byte, cfg Config) []Finding {
	md := goldmark.New()
	root := md.Parser().Parse(text.NewReader(src))

	var findings []Finding
	var prevLevel int
	var skipLevel int
	var skipping bool

	ast.Walk(root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		heading, ok := n.(*ast.Heading)
		if !ok || !entering {
			return ast.WalkContinue, nil
		}

		title := string(heading.Text(src))
		level := heading.Level

		if skipping {
			if level <= skipLevel {
				skipping = false
			} else {
				return ast.WalkContinue, nil
			}
		}

		if contains(cfg.Exclude, title) {
			skipping = true
			skipLevel = level
			prevLevel = level
			return ast.WalkContinue, nil
		}

		if prevLevel != 0 && level > prevLevel+1 {
			seg := heading.Lines().At(0)
			line := bytes.Count(src[:seg.Start], []byte("\n")) + 1
			findings = append(findings, Finding{
				Line:    line,
				Message: "heading level should only increment by one level at a time",
			})
		}

		prevLevel = level
		return ast.WalkContinue, nil
	})

	return findings
}

func contains(list []string, s string) bool {
	for _, v := range list {
		if v == s {
			return true
		}
	}
	return false
}
