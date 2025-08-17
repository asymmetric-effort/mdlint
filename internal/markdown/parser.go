// Copyright (c) 2024 Asymmetric Effort

package markdown

import (
	"sync"

	"github.com/yuin/goldmark"
	meta "github.com/yuin/goldmark-meta"
	"github.com/yuin/goldmark/extension"
)

// parserOnce ensures singleton initialization of the Markdown parser.
var parserOnce sync.Once

// parserInstance holds the configured goldmark Markdown parser.
var parserInstance goldmark.Markdown

// Parser returns a goldmark Markdown parser configured with GFM, footnote,
// table, and front-matter extensions. The parser is lazily initialized and
// safe for concurrent use.
func Parser() goldmark.Markdown {
	parserOnce.Do(func() {
		parserInstance = goldmark.New(
			goldmark.WithExtensions(
				extension.GFM,
				extension.Footnote,
				extension.Table,
				meta.Meta,
			),
		)
	})
	return parserInstance
}
