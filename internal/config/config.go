// Copyright 2024 MdLint Authors

package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents linter configuration.
type Config struct {
	Format string `yaml:"format"`
}

// Load reads configuration from path or returns defaults when path is empty.
func Load(path string) (*Config, error) {
	cfg := &Config{Format: "json"}
	if path == "" {
		return cfg, nil
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, err
	}
	if cfg.Format == "" {
		cfg.Format = "json"
	}
	return cfg, nil
}
