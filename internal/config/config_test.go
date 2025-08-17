// Copyright 2025 Sam Caldwell
// SPDX-License-Identifier: MIT

package config

import (
	"os"
	"path/filepath"
	"testing"
)

// TestLoadDefault verifies that defaults are applied when no configuration files are present.
func TestLoadDefault(t *testing.T) {
	cfg, err := Load(Config{}, "")
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if cfg.Version != 1 {
		t.Fatalf("unexpected version %d", cfg.Version)
	}
	if cfg.Output.Format != "json" || cfg.Output.Color != "auto" {
		t.Fatalf("unexpected output defaults: %+v", cfg.Output)
	}
	if cfg.FailureThreshold != "warning" {
		t.Fatalf("unexpected failure threshold %q", cfg.FailureThreshold)
	}
}

// TestLoadMergePrecedence ensures CLI overrides project config which overrides user config.
func TestLoadMergePrecedence(t *testing.T) {
	tmp := t.TempDir()

	// user config
	userCfgDir := filepath.Join(tmp, "user", "mdlint")
	if err := os.MkdirAll(userCfgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(userCfgDir, "config.yaml"), []byte("version: 1\noutput:\n  color: always\nfailure_threshold: suggestion\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("XDG_CONFIG_HOME", filepath.Join(tmp, "user"))

	// project config
	projDir := filepath.Join(tmp, "proj")
	if err := os.MkdirAll(projDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(projDir, ".mdlintrc.yaml"), []byte("version: 1\noutput:\n  color: never\nfailure_threshold: warning\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	cli := Config{FailureThreshold: "error"}
	cfg, err := Load(cli, projDir)
	if err != nil {
		t.Fatalf("Load failed: %v", err)
	}
	if cfg.FailureThreshold != "error" {
		t.Fatalf("expected CLI to override failure threshold, got %q", cfg.FailureThreshold)
	}
	if cfg.Output.Color != "never" {
		t.Fatalf("expected project config to override user output color, got %q", cfg.Output.Color)
	}
}

// TestLoadValidationError ensures invalid config values are rejected.
func TestLoadValidationError(t *testing.T) {
	tmp := t.TempDir()
	userCfgDir := filepath.Join(tmp, "mdlint")
	if err := os.MkdirAll(userCfgDir, 0o755); err != nil {
		t.Fatal(err)
	}
	// invalid output format
	if err := os.WriteFile(filepath.Join(userCfgDir, "config.yaml"), []byte("version: 1\noutput:\n  format: xml\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	t.Setenv("XDG_CONFIG_HOME", tmp)

	if _, err := Load(Config{}, ""); err == nil {
		t.Fatalf("expected error for invalid output format")
	}

	// unknown field
	if err := os.WriteFile(filepath.Join(userCfgDir, "config.yaml"), []byte("version: 1\nunknown: true\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Load(Config{}, ""); err == nil {
		t.Fatalf("expected error for unknown field")
	}
}

