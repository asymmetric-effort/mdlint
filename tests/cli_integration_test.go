// Copyright 2024 MdLint Authors

package tests

import (
	"bytes"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// run executes the CLI with given arguments.
func run(args ...string) (string, int, error) {
	cmdArgs := append([]string{"run", "../cmd/mdlint"}, args...)
	cmd := exec.Command("go", cmdArgs...)
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out
	err := cmd.Run()
	exit := 0
	if err != nil {
		if ee, ok := err.(*exec.ExitError); ok {
			exit = ee.ExitCode()
		} else {
			return out.String(), -1, err
		}
	}
	return out.String(), exit, nil
}

// TestCLI_NoFindings ensures exit code 0 when no issues.
func TestCLI_NoFindings(t *testing.T) {
	out, code, err := run("--format", "json", filepath.Join("..", "testdata", "good.md"))
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if code != 0 {
		t.Fatalf("expected exit 0 got %d output %s", code, out)
	}
	if out != "" {
		t.Fatalf("expected no output got %s", out)
	}
}

// TestCLI_Findings ensures exit code 1 when findings exist.
func TestCLI_Findings(t *testing.T) {
	out, code, err := run("--format", "json", filepath.Join("..", "testdata", "bad.md"))
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if code != 1 {
		t.Fatalf("expected exit 1 got %d output %s", code, out)
	}
}

// TestCLI_ListRules ensures rules are listed.
func TestCLI_ListRules(t *testing.T) {
	out, code, err := run("--list-rules")
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if code != 0 || !strings.Contains(out, "MD9000") {
		t.Fatalf("unexpected: code %d output %s", code, out)
	}
}

// TestCLI_Version ensures version flag prints version.
func TestCLI_Version(t *testing.T) {
	out, code, err := run("--version")
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if code != 0 || !strings.Contains(out, "mdlint") {
		t.Fatalf("unexpected: code %d output %s", code, out)
	}
}

// TestCLI_ConfigFormat ensures config file controls formatting.
func TestCLI_ConfigFormat(t *testing.T) {
	out, code, err := run("--config", filepath.Join("..", "testdata", "config.yaml"), filepath.Join("..", "testdata", "bad.md"))
	if err != nil {
		t.Fatalf("run: %v", err)
	}
	if code != 1 {
		t.Fatalf("expected exit 1 got %d output %s", code, out)
	}
	if !strings.Contains(out, "MD9000") || strings.HasPrefix(strings.TrimSpace(out), "[") {
		t.Fatalf("expected text findings got %s", out)
	}
}
