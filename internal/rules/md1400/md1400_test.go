// Copyright (c) 2024 MdLint contributors.

package md1400

import "testing"

func TestCheckCodeBlockLanguages(t *testing.T) {
	tests := []struct {
		name string
		src  string
		cfg  Config
		want int
	}{
		{
			name: "unknown language",
			src:  "```foobar\ncode\n```\n",
			cfg:  Config{},
			want: 1,
		},
		{
			name: "missing language",
			src:  "```\ncode\n```\n",
			cfg:  Config{},
			want: 1,
		},
		{
			name: "allowed language",
			src:  "```go\ncode\n```\n",
			cfg:  Config{Allowed: []string{"go"}},
			want: 0,
		},
		{
			name: "disallowed language",
			src:  "```python\ncode\n```\n",
			cfg:  Config{Allowed: []string{"go"}},
			want: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckCodeBlockLanguages([]byte(tt.src), tt.cfg)
			if len(got) != tt.want {
				t.Fatalf("got %d findings, want %d", len(got), tt.want)
			}
		})
	}
}
