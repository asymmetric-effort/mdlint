// Copyright (c) 2024 MdLint contributors.
// SPDX-License-Identifier: MIT

// Package md1400 provides a rule ensuring fenced code blocks use recognized languages.
package md1400

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// Config configures the MD1400 rule.
type Config struct {
	// Allowed lists the permitted language identifiers. If empty, any language
	// recognized by Chroma is allowed.
	Allowed []string
}

// Finding reports a code fence language issue.
type Finding struct {
	// Line is the line number where the fence starts.
	Line int
	// Message describes the violation.
	Message string
}

// CheckCodeBlockLanguages verifies that fenced code blocks specify allowed and
// recognized language identifiers.
func CheckCodeBlockLanguages(src []byte, cfg Config) []Finding {
	md := goldmark.New()
	root := md.Parser().Parse(text.NewReader(src))

	allowed := map[string]struct{}{}
	for _, lang := range cfg.Allowed {
		allowed[strings.ToLower(lang)] = struct{}{}
	}

	var findings []Finding
	ast.Walk(root, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		block, ok := n.(*ast.FencedCodeBlock)
		if !ok || !entering {
			return ast.WalkContinue, nil
		}
		lang := strings.ToLower(string(block.Language(src)))
		seg := block.Lines().At(0)
		line := bytes.Count(src[:seg.Start], []byte("\n")) + 1

		if lang == "" {
			findings = append(findings, Finding{Line: line, Message: "code fence is missing a language identifier"})
			return ast.WalkContinue, nil
		}
		if lexers.Get(lang) == nil {
			findings = append(findings, Finding{Line: line, Message: fmt.Sprintf("unknown language %q", lang)})
			return ast.WalkContinue, nil
		}
		if len(allowed) > 0 {
			if _, ok := allowed[lang]; !ok {
				findings = append(findings, Finding{Line: line, Message: fmt.Sprintf("language %q not allowed", lang)})
			}
		}
		return ast.WalkContinue, nil
	})

	return findings
}
