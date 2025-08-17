// Copyright (c) 2024 MdLint contributors.

package md1400

import "testing"

func TestCheckCodeBlockLanguages(t *testing.T) {
	tests := []struct {
		name string
		src  string
		cfg  Config
		want []Finding
	}{
		{
			name: "unknown language",
			src:  "```foobar\ncode\n```\n",
			cfg:  Config{},
			want: []Finding{{Line: 1, Message: "unknown language \"foobar\""}},
		},
		{
			name: "missing language",
			src:  "```\ncode\n```\n",
			cfg:  Config{},
			want: []Finding{{Line: 1, Message: "code fence is missing a language identifier"}},
		},
		{
			name: "allowed language",
			src:  "```go\ncode\n```\n",
			cfg:  Config{Allowed: []string{"go"}},
			want: nil,
		},
		{
			name: "disallowed language",
			src:  "```python\ncode\n```\n",
			cfg:  Config{Allowed: []string{"go"}},
			want: []Finding{{Line: 1, Message: "language \"python\" not allowed"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := CheckCodeBlockLanguages([]byte(tt.src), tt.cfg)
			if len(got) != len(tt.want) {
				t.Fatalf("got %d findings, want %d", len(got), len(tt.want))
			}
			for i := range tt.want {
				if got[i] != tt.want[i] {
					t.Fatalf("finding %d: got %+v, want %+v", i, got[i], tt.want[i])
				}
			}
		})
	}
}
