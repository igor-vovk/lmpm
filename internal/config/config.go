package config

import (
	"fmt"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
)

type Source struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

type Include struct {
	Source string   `yaml:"source"`
	Files  []string `yaml:"files"`
}

type Strategy string

const (
	StrategyFlatten  Strategy = "flatten"
	StrategyPreserve Strategy = "preserve"
	StrategyConcat   Strategy = "concat"
)

type Target struct {
	Name     string    `yaml:"name"`
	Output   string    `yaml:"output"`
	Strategy Strategy  `yaml:"strategy,omitempty"`
	Include  []Include `yaml:"include"`
}

type Config struct {
	Version int      `yaml:"version"`
	Sources []Source `yaml:"sources"`
	Targets []Target `yaml:"targets"`
}

func NewConfig() *Config {
	return &Config{
		Version: 1,
		Sources: []Source{},
		Targets: []Target{},
	}
}

func LoadConfig(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	cfg := NewConfig()
	if err := yaml.Unmarshal(data, cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if err := cfg.addWorkingDirSource(); err != nil {
		return nil, err
	}

	cfg.normalize()

	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (c *Config) addWorkingDirSource() error {
	wd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("failed to get working directory: %w", err)
	}

	hasWorkingDir := false
	for _, source := range c.Sources {
		if source.Name == "working_dir" {
			hasWorkingDir = true
			break
		}
	}

	if !hasWorkingDir {
		c.Sources = append([]Source{{Name: "working_dir", URL: wd}}, c.Sources...)
	}

	return nil
}

func (c *Config) normalize() {
	c.setDefaultSourceForIncludes()
	c.setDefaultStrategy()
}

func (c *Config) setDefaultSourceForIncludes() {
	for i := range c.Targets {
		for j := range c.Targets[i].Include {
			include := &c.Targets[i].Include[j]
			if include.Source == "" {
				include.Source = "working_dir"
			}
		}
	}
}

func (c *Config) setDefaultStrategy() {
	for i := range c.Targets {
		target := &c.Targets[i]
		if target.Strategy == "" {
			// If output has .md or .txt extension, default to concat
			if hasTextExtension(target.Output) {
				target.Strategy = StrategyConcat
			} else {
				target.Strategy = StrategyFlatten
			}
		}
	}
}

func hasTextExtension(path string) bool {
	ext := filepath.Ext(path)
	return ext == ".md" || ext == ".txt"
}

func (c *Config) Validate() error {
	sourceKeys := make(map[string]bool)
	for _, source := range c.Sources {
		if source.Name == "" {
			return fmt.Errorf("source key cannot be empty")
		}
		if sourceKeys[source.Name] {
			return fmt.Errorf("duplicate source key: %s", source.Name)
		}
		sourceKeys[source.Name] = true
	}

	for _, target := range c.Targets {
		if target.Strategy != StrategyFlatten && target.Strategy != StrategyPreserve && target.Strategy != StrategyConcat {
			return fmt.Errorf("target '%s' has invalid strategy: %s (must be 'flatten', 'preserve', or 'concat')", target.Name, target.Strategy)
		}

		for _, include := range target.Include {
			if !sourceKeys[include.Source] {
				return fmt.Errorf("target '%s' references unknown source: %s", target.Name, include.Source)
			}
		}
	}

	return nil
}
