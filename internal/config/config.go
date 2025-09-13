package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type Config struct {
	Generator    GeneratorConfig    `mapstructure:"generator" json:"generator"`
	Placeholders PlaceholdersConfig `mapstructure:"placeholders" json:"placeholders"`
	Output       OutputConfig       `mapstructure:"output" json:"output"`
}

type GeneratorConfig struct {
	MinYear        int                 `mapstructure:"min_year" json:"min_year"`
	MaxYear        int                 `mapstructure:"max_year" json:"max_year"`
	MinPasswordLen int                 `mapstructure:"min_password_length" json:"min_password_length"`
	MaxPasswordLen int                 `mapstructure:"max_password_length" json:"max_password_length"`
	CommonWords    []string            `mapstructure:"common_words" json:"common_words"`
	Separators     []string            `mapstructure:"separators" json:"separators"`
	Substitutions  map[string][]string `mapstructure:"substitutions" json:"substitutions"`
	NumberPatterns []string            `mapstructure:"number_patterns" json:"number_patterns"`
	Patterns       []string            `mapstructure:"patterns" json:"patterns"`
}

type OutputConfig struct {
	Filename string `mapstructure:"filename" json:"filename"`
}

type JSONConfig struct {
	CommonWords    []string            `json:"common_words,omitempty"`
	Separators     []string            `json:"separators,omitempty"`
	NumberPatterns []string            `json:"number_patterns,omitempty"`
	Substitutions  map[string][]string `json:"substitutions,omitempty"`
	Patterns       []string            `json:"patterns,omitempty"`
}

func Load(jsonConfigPath string) (*Config, error) {
	cfg := &Config{
		Generator:    NewDefaultGeneratorConfig(),
		Placeholders: NewDefaultPlaceholdersConfig(),
		Output:       NewDefaultOutputConfig(),
	}

	if jsonConfigPath != "" {
		if err := cfg.loadFromJSON(jsonConfigPath); err != nil {
			return nil, fmt.Errorf("failed to load JSON config: %w", err)
		}
	}

	return cfg, nil
}

func (c *Config) loadFromJSON(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read JSON config file: %w", err)
	}

	var jsonConfig JSONConfig
	if err := json.Unmarshal(data, &jsonConfig); err != nil {
		return fmt.Errorf("failed to parse JSON config: %w", err)
	}

	c.applyJSONConfig(&jsonConfig)
	return nil
}

func (c *Config) applyJSONConfig(jsonConfig *JSONConfig) {
	if len(jsonConfig.CommonWords) > 0 {
		c.Generator.CommonWords = jsonConfig.CommonWords
	}
	if len(jsonConfig.Separators) > 0 {
		c.Generator.Separators = jsonConfig.Separators
	}
	if len(jsonConfig.NumberPatterns) > 0 {
		c.Generator.NumberPatterns = jsonConfig.NumberPatterns
	}
	if len(jsonConfig.Substitutions) > 0 {
		c.Generator.Substitutions = jsonConfig.Substitutions
	}
	if len(jsonConfig.Patterns) > 0 {
		c.Generator.Patterns = jsonConfig.Patterns
	}
}