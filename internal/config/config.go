// Copyright 2024 The mdlint Authors
// SPDX-License-Identifier: MIT

// Package config provides loading and validation for mdlint configuration.
package config

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"gopkg.in/yaml.v3"
)

// Severity represents a rule severity level.
type Severity string

// Config holds the top-level configuration for mdlint.
type Config struct {
	Version          int                   `yaml:"version"`
	Ignored          []string              `yaml:"ignored"`
	Severity         map[string]Severity   `yaml:"severity"`
	Paths            map[string]PathConfig `yaml:"paths"`
	Spell            SpellConfig           `yaml:"spell"`
	Heading          HeadingConfig         `yaml:"heading"`
	Output           OutputConfig          `yaml:"output"`
	FailureThreshold Severity              `yaml:"failure_threshold"`
}

// PathConfig defines per-path overrides.
type PathConfig struct {
	Ignored  []string            `yaml:"ignored"`
	Severity map[string]Severity `yaml:"severity"`
}

// SpellConfig defines options for the spelling rule MD1000.
type SpellConfig struct {
	Lang        string   `yaml:"lang"`
	AddWords    []string `yaml:"add_words"`
	RejectWords []string `yaml:"reject_words"`
	Filters     []string `yaml:"filters"`
}

// HeadingConfig defines options for heading style checks.
type HeadingConfig struct {
	Style      string `yaml:"style"`
	AllowMixed *bool  `yaml:"allow_mixed"`
}

// OutputConfig defines formatting options for findings output.
type OutputConfig struct {
	Format string `yaml:"format"`
	Color  string `yaml:"color"`
}

// DefaultConfig returns configuration with built-in defaults.
func DefaultConfig() Config {
	allowMixed := false
	return Config{
		Version:          1,
		Output:           OutputConfig{Format: "json", Color: "auto"},
		Heading:          HeadingConfig{AllowMixed: &allowMixed},
		FailureThreshold: "warning",
	}
}

// Load resolves configuration from user, project and CLI sources in precedence order.
// CLI overrides are provided via the cli parameter; projectDir determines where the
// project configuration file is looked up.
func Load(cli Config, projectDir string) (Config, error) {
	cfg := DefaultConfig()

	if userCfg, err := readConfigFile(userConfigPath()); err == nil {
		merge(&cfg, userCfg)
	} else if !errors.Is(err, os.ErrNotExist) {
		return Config{}, err
	}

	if projectDir != "" {
		projPath := filepath.Join(projectDir, ".mdlintrc.yaml")
		if projCfg, err := readConfigFile(projPath); err == nil {
			merge(&cfg, projCfg)
		} else if !errors.Is(err, os.ErrNotExist) {
			return Config{}, err
		}
	}

	merge(&cfg, cli)

	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}

	return cfg, nil
}

func readConfigFile(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, err
	}
	cfg, err := parseYAML(data)
	if err != nil {
		return Config{}, err
	}
	if err := cfg.Validate(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func parseYAML(data []byte) (Config, error) {
	var cfg Config
	dec := yaml.NewDecoder(bytes.NewReader(data))
	dec.KnownFields(true)
	if err := dec.Decode(&cfg); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

// Validate performs custom schema checks beyond YAML decoding.
func (c Config) Validate() error {
	if c.Version != 1 {
		return fmt.Errorf("unsupported version %d", c.Version)
	}

	validSev := map[Severity]bool{"suggestion": true, "warning": true, "error": true}
	checkSev := func(sev Severity) error {
		if sev == "" {
			return nil
		}
		if !validSev[sev] {
			return fmt.Errorf("invalid severity %q", sev)
		}
		return nil
	}

	for _, sev := range c.Severity {
		if err := checkSev(sev); err != nil {
			return err
		}
	}
	for _, pc := range c.Paths {
		for _, sev := range pc.Severity {
			if err := checkSev(sev); err != nil {
				return err
			}
		}
	}

	if c.Heading.Style != "" {
		switch c.Heading.Style {
		case "atx", "setext", "consistent":
		default:
			return fmt.Errorf("invalid heading style %q", c.Heading.Style)
		}
	}

	switch c.Output.Format {
	case "", "json", "text":
	default:
		return fmt.Errorf("invalid output format %q", c.Output.Format)
	}
	switch c.Output.Color {
	case "", "auto", "always", "never":
	default:
		return fmt.Errorf("invalid output color %q", c.Output.Color)
	}

	if err := checkSev(c.FailureThreshold); err != nil {
		return fmt.Errorf("invalid failure threshold: %w", err)
	}
	return nil
}

func merge(dst *Config, src Config) {
	if src.Version != 0 {
		dst.Version = src.Version
	}
	if len(src.Ignored) > 0 {
		dst.Ignored = append(dst.Ignored, src.Ignored...)
	}
	if src.Severity != nil {
		if dst.Severity == nil {
			dst.Severity = make(map[string]Severity)
		}
		for k, v := range src.Severity {
			dst.Severity[k] = v
		}
	}
	if src.Paths != nil {
		if dst.Paths == nil {
			dst.Paths = make(map[string]PathConfig)
		}
		for p, pc := range src.Paths {
			existing := dst.Paths[p]
			if len(pc.Ignored) > 0 {
				existing.Ignored = append(existing.Ignored, pc.Ignored...)
			}
			if pc.Severity != nil {
				if existing.Severity == nil {
					existing.Severity = make(map[string]Severity)
				}
				for rk, rv := range pc.Severity {
					existing.Severity[rk] = rv
				}
			}
			dst.Paths[p] = existing
		}
	}
	if src.Spell.Lang != "" {
		dst.Spell.Lang = src.Spell.Lang
	}
	if len(src.Spell.AddWords) > 0 {
		dst.Spell.AddWords = append(dst.Spell.AddWords, src.Spell.AddWords...)
	}
	if len(src.Spell.RejectWords) > 0 {
		dst.Spell.RejectWords = append(dst.Spell.RejectWords, src.Spell.RejectWords...)
	}
	if len(src.Spell.Filters) > 0 {
		dst.Spell.Filters = append(dst.Spell.Filters, src.Spell.Filters...)
	}
	if src.Heading.Style != "" {
		dst.Heading.Style = src.Heading.Style
	}
	if src.Heading.AllowMixed != nil {
		dst.Heading.AllowMixed = src.Heading.AllowMixed
	}
	if src.Output.Format != "" {
		dst.Output.Format = src.Output.Format
	}
	if src.Output.Color != "" {
		dst.Output.Color = src.Output.Color
	}
	if src.FailureThreshold != "" {
		dst.FailureThreshold = src.FailureThreshold
	}
}

func userConfigPath() string {
	if xdg := os.Getenv("XDG_CONFIG_HOME"); xdg != "" {
		return filepath.Join(xdg, "mdlint", "config.yaml")
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return "config.yaml"
	}
	return filepath.Join(home, ".config", "mdlint", "config.yaml")
}

