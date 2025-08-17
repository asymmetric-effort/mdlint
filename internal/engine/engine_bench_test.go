// Copyright 2025 The mdlint Authors
// SPDX-License-Identifier: MIT

package engine

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/asymmetric-effort/mdlint/internal/config"
)

// BenchmarkRun measures performance of the engine across multiple files.
func BenchmarkRun(b *testing.B) {
	dir := b.TempDir()
	for i := 0; i < 10; i++ {
		name := filepath.Join(dir, fmt.Sprintf("f%d.md", i))
		if err := os.WriteFile(name, []byte("line\n"), 0o644); err != nil {
			b.Fatal(err)
		}
	}
	paths, err := filepath.Glob(filepath.Join(dir, "*.md"))
	if err != nil {
		b.Fatal(err)
	}
	eng := Engine{Limits: config.LimitsConfig{MaxFiles: 100, MaxFileSize: 1 << 20, Concurrency: 4}}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if _, err := eng.Run(paths); err != nil {
			b.Fatal(err)
		}
	}
}
