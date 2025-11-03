package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/goccy/go-yaml"
)

type Source struct {
	Name string `yaml:"name"`
	URL  string `yaml:"url"`
}

const DefaultSourceName = "working_dir"

type Strategy string

const (
	StrategyFlatten  Strategy = "flatten"
	StrategyPreserve Strategy = "preserve"
	StrategyConcat   Strategy = "concat"
)

type Include struct {
	Source string
	Files  []string
}

type Target struct {
	Name          string    `yaml:"name"`
	Output        string    `yaml:"output"`
	Strategy      Strategy  `yaml:"strategy,omitempty"`
	Include       []string  `yaml:"include"`
	IncludeParsed []Include `yaml:"-"`
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

	if err := cfg.normalize(); err != nil {
		return nil, err
	}

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
		if source.Name == DefaultSourceName {
			hasWorkingDir = true
			break
		}
	}

	if !hasWorkingDir {
		c.Sources = append([]Source{{Name: DefaultSourceName, URL: wd}}, c.Sources...)
	}

	return nil
}

func (c *Config) normalize() error {
	if err := c.parseIncludes(); err != nil {
		return err
	}
	c.setDefaultSourceForIncludes()
	c.setDefaultStrategy()
	return nil
}

func (c *Config) parseIncludes() error {
	for i := range c.Targets {
		var includes []Include
		for _, includeStr := range c.Targets[i].Include {
			include, err := ParseInclude(includeStr)
			if err != nil {
				return fmt.Errorf("failed to parse include in target '%s': %w", c.Targets[i].Name, err)
			}
			includes = append(includes, include)
		}
		c.Targets[i].IncludeParsed = includes
	}
	return nil
}

func (c *Config) setDefaultSourceForIncludes() {
	for i := range c.Targets {
		for j := range c.Targets[i].IncludeParsed {
			include := &c.Targets[i].IncludeParsed[j]
			if include.Source == "" {
				include.Source = DefaultSourceName
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

		for _, include := range target.IncludeParsed {
			if !sourceKeys[include.Source] {
				return fmt.Errorf("target '%s' references unknown source: %s", target.Name, include.Source)
			}
		}
	}

	return nil
}

func ParseInclude(includeStr string) (Include, error) {
	// if includeStr starts with @, it's structure is "@source/path1,path2,..."
	if len(includeStr) > 0 && includeStr[0] == '@' {
		parts := strings.SplitN(includeStr[1:], "/", 2)
		if len(parts) != 2 {
			return Include{}, fmt.Errorf("invalid include format: %s", includeStr)
		}

		source := parts[0]
		filePaths := strings.Split(parts[1], ",")
		for i := range filePaths {
			filePaths[i] = strings.TrimSpace(filePaths[i])
		}

		return Include{
			Source: source,
			Files:  filePaths,
		}, nil
	} else {
		// otherwise, it's just a path in the working_dir source
		filePaths := strings.Split(includeStr, ",")
		for i := range filePaths {
			filePaths[i] = strings.TrimSpace(filePaths[i])
		}

		return Include{
			Source: DefaultSourceName,
			Files:  filePaths,
		}, nil
	}
}
