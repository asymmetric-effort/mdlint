// Copyright 2024 the mdlint authors.

package rules

import "testing"

func TestCheckTrailingWhitespace(t *testing.T) {
	tests := []struct {
		name    string
		content string
		cfg     TrailingWhitespaceConfig
		want    []int
	}{
		{
			name:    "mixed whitespace with code block ignored",
			content: "no trailing\nline with space \n```go\ncode line with space  \n```\nline with tab\t\n",
			cfg:     TrailingWhitespaceConfig{IgnoreCodeBlocks: true},
			want:    []int{2, 6},
		},
		{
			name:    "check code blocks when not ignored",
			content: "no trailing\nline with space \n```go\ncode line with space  \n```\nline with tab\t\n",
			cfg:     TrailingWhitespaceConfig{IgnoreCodeBlocks: false},
			want:    []int{2, 4, 6},
		},
		{
			name:    "only code block when ignored",
			content: "```go\ncode  \n```\n",
			cfg:     TrailingWhitespaceConfig{IgnoreCodeBlocks: true},
			want:    nil,
		},
		{
			name:    "blank line with spaces",
			content: "line\n   \nnext\n",
			cfg:     TrailingWhitespaceConfig{IgnoreCodeBlocks: true},
			want:    []int{2},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckTrailingWhitespace([]byte(tt.content), tt.cfg)
			if len(got) != len(tt.want) {
				t.Fatalf("got %d findings, want %d", len(got), len(tt.want))
			}
			for i, f := range got {
				if f.Line != tt.want[i] {
					t.Errorf("finding %d line = %d, want %d", i, f.Line, tt.want[i])
				}
			}
		})
	}
}
