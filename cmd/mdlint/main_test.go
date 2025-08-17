// Copyright 2025 Sam Caldwell
//
// Tests for the mdlint command-line interface.
package main_test

import (
	"bytes"
	"os/exec"
	"testing"

	"github.com/sam-caldwell/mdlint/internal/version"
)

// TestVersionFlag verifies that the --version flag prints the semantic version.
func TestVersionFlag(t *testing.T) {
	cmd := exec.Command("go", "run", ".", "--version")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	if err := cmd.Run(); err != nil {
		t.Fatalf("command failed: %v output: %s", err, out.String())
	}
	expected := version.Version + "\n"
	if out.String() != expected {
		t.Fatalf("expected %q got %q", expected, out.String())
	}
}
