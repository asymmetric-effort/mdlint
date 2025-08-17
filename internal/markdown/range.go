// Copyright (c) 2024 Asymmetric Effort

package markdown

import (
	"bytes"

	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

// Range represents an inclusive line range within a Markdown document.
type Range struct {
	StartLine int // 1-based line number where the range starts
	EndLine   int // 1-based line number where the range ends
}

// CodeBlockRanges returns line ranges for all code blocks in the provided source.
// Ranges include fenced and indented code blocks.
func CodeBlockRanges(src []byte) []Range {
	md := Parser()
	reader := text.NewReader(src)
	doc := md.Parser().Parse(reader)

	var ranges []Range
	ast.Walk(doc, func(n ast.Node, entering bool) (ast.WalkStatus, error) {
		if !entering {
			return ast.WalkContinue, nil
		}
		switch n := n.(type) {
		case *ast.FencedCodeBlock:
			lines := n.Lines()
			if lines == nil || lines.Len() == 0 {
				return ast.WalkContinue, nil
			}
			first := lines.At(0)
			last := lines.At(lines.Len() - 1)
			start := lineNumber(src, first.Start) - 1
			if start < 1 {
				start = 1
			}
			end := lineNumber(src, last.Stop)
			ranges = append(ranges, Range{StartLine: start, EndLine: end})
			return ast.WalkSkipChildren, nil
		case *ast.CodeBlock:
			lines := n.Lines()
			if lines == nil || lines.Len() == 0 {
				return ast.WalkContinue, nil
			}
			first := lines.At(0)
			last := lines.At(lines.Len() - 1)
			ranges = append(ranges, Range{
				StartLine: lineNumber(src, first.Start),
				EndLine:   lineNumber(src, last.Stop),
			})
			return ast.WalkSkipChildren, nil
		}
		return ast.WalkContinue, nil
	})
	return ranges
}

// FrontMatterRange returns the line range for a leading YAML front-matter
// section. If no front-matter is present, ok will be false.
func FrontMatterRange(src []byte) (rng Range, ok bool) {
	lines := bytes.Split(src, []byte("\n"))
	if len(lines) == 0 {
		return Range{}, false
	}
	if !bytes.Equal(bytes.TrimSpace(lines[0]), []byte("---")) {
		return Range{}, false
	}
	for i := 1; i < len(lines); i++ {
		if bytes.Equal(bytes.TrimSpace(lines[i]), []byte("---")) {
			return Range{StartLine: 1, EndLine: i + 1}, true
		}
	}
	return Range{}, false
}

// lineNumber converts a byte offset into a 1-based line number.
func lineNumber(src []byte, pos int) int {
	return bytes.Count(src[:pos], []byte("\n")) + 1
}
