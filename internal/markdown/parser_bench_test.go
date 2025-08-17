// Copyright 2025 The mdlint Authors
// SPDX-License-Identifier: MIT

package markdown

import (
	"os"
	"testing"

	"github.com/yuin/goldmark/text"
)

// BenchmarkParser measures performance of Markdown parsing.
func BenchmarkParser(b *testing.B) {
	data, err := os.ReadFile("testdata/basic.md")
	if err != nil {
		b.Fatal(err)
	}
	md := Parser()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		md.Parser().Parse(text.NewReader(data))
	}
}
