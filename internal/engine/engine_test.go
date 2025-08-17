// Copyright 2025 The mdlint Authors
// SPDX-License-Identifier: MIT

package engine

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/asymmetric-effort/mdlint/internal/config"
)

func writeFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write file: %v", err)
	}
	return path
}

// TestRun_MaxFiles verifies that the engine enforces the maximum file count.
func TestRun_MaxFiles(t *testing.T) {
	dir := t.TempDir()
	p1 := writeFile(t, dir, "a.md", "ok")
	p2 := writeFile(t, dir, "b.md", "ok")
	eng := Engine{Limits: config.LimitsConfig{MaxFiles: 1, MaxFileSize: 1 << 25, Concurrency: 1}}
	if _, err := eng.Run([]string{p1, p2}); err == nil {
		t.Fatalf("expected error for too many files")
	}
}

// TestRun_MaxFileSize verifies that oversized files are rejected.
func TestRun_MaxFileSize(t *testing.T) {
	dir := t.TempDir()
	big := make([]byte, 1024)
	p := writeFile(t, dir, "big.md", string(big))
	eng := Engine{Limits: config.LimitsConfig{MaxFiles: 10, MaxFileSize: 10, Concurrency: 1}}
	if _, err := eng.Run([]string{p}); err == nil {
		t.Fatalf("expected error for large file")
	}
}

// TestRun_DisablesNetwork ensures network access is blocked during execution.
func TestRun_DisablesNetwork(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	}))
	defer server.Close()

	dir := t.TempDir()
	big := bytes.Repeat([]byte("line\n"), 1<<20)
	p := writeFile(t, dir, "a.md", string(big))
	eng := Engine{Limits: config.LimitsConfig{MaxFiles: 1, MaxFileSize: 1 << 25, Concurrency: 1}}

	done := make(chan struct{})
	go func() {
		if _, err := eng.Run([]string{p}); err != nil {
			t.Errorf("run: %v", err)
		}
		close(done)
	}()
	var blocked bool
	for {
		select {
		case <-done:
			if !blocked {
				t.Fatalf("expected network error during run")
			}
			if _, err := http.Get(server.URL); err != nil {
				t.Fatalf("expected network restored: %v", err)
			}
			return
		default:
			if _, err := http.Get(server.URL); err != nil {
				blocked = true
			}
		}
	}
}
