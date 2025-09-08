package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"
)

type Config struct {
	Generator GeneratorConfig `mapstructure:"generator" json:"generator"`
	Output    OutputConfig    `mapstructure:"output" json:"output"`
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
	CommonWords    []string            `json:"common_words"`
	Separators     []string            `json:"separators"`
	NumberPatterns []string            `json:"number_patterns"`
	Substitutions  map[string][]string `json:"substitutions"`
	Patterns       []string            `json:"patterns"`
}

func Load(jsonConfigPath string) (*Config, error) {
	cfg := &Config{
		Generator: GeneratorConfig{
			MinYear:        1990,
			MaxYear:        time.Now().Year(),
			MinPasswordLen: 8,
			MaxPasswordLen: 64,
			CommonWords: []string{
				"password",
				"admin",
				"guest",
				"wifi",
				"wireless",
				"IT",
				"tech",
				"pass",
				"login",
				"user",
				"123",
				"root",
				"default",
				"access",
				"network",
				"internet",
				"secure",
				"temp",
				"test",
			},
			Separators: []string{"", "@", "_", "-", ".", "#", "!", "*", "+", "~", "%", "&", "^"},
			Substitutions: map[string][]string{
				"a": {"4", "@", "^"},
				"A": {"4", "@", "^"},
				"b": {"6"},
				"B": {"8"},
				"c": {"<", "("},
				"C": {"<", "("},
				"D": {")"},
				"e": {"3"},
				"E": {"3"},
				"g": {"9", "6", "&"},
				"G": {"9", "6", "&"},
				"h": {"#"},
				"H": {"#"},
				"i": {"1", "!", "|"},
				"I": {"1", "!", "|"},
				"l": {"1", "|", "7", "2"},
				"L": {"1", "|", "7", "2"},
				"o": {"0"},
				"O": {"0"},
				"p": {"9"},
				"P": {"9"},
				"q": {"9", "2", "&"},
				"Q": {"9", "2", "&"},
				"s": {"5", "$"},
				"S": {"5", "$"},
				"t": {"7", "+"},
				"T": {"7", "+"},
				"z": {"2"},
				"Z": {"2"},
			},
			NumberPatterns: []string{
				"ddd",
				"1",
				"12",
				"1234",
				"12345",
				"123456",
				"01",
				"1000",
			},
			Patterns: []string{
				"<CUSTOM>",
				"<COMMON>",
				"<SSID>",
				"<CUSTOM><SEP><YEAR>",
				"<CUSTOM><SEP><NUM>",
				"<CUSTOM><SEP><COMMON>",
				"<COMMON><SEP><YEAR>",
				"<COMMON><SEP><NUM>",
				"<COMMON><SEP><CUSTOM>",
				"<SSID><SEP><YEAR>",
				"<SSID><SEP><NUM>",
				"<SSID><SEP><CUSTOM>",
				"<YEAR><SEP><CUSTOM>",
				"<YEAR><SEP><COMMON>",
				"<YEAR><SEP><SSID>",
				"<NUM><SEP><CUSTOM>",
				"<NUM><SEP><COMMON>",
				"<NUM><SEP><SSID>",
				"<SEP><CUSTOM><SEP><SSID><SEP><YEAR>",
				"<COMMON><SEP><CUSTOM><YEAR>",
			},
		},
		Output: OutputConfig{
			Filename: "passwords.ls",
		},
	}

	if jsonConfigPath != "" {
		if err := cfg.loadFromJSON(jsonConfigPath); err != nil {
			return nil, fmt.Errorf("failed to load JSON config: %w", err)
		}
	}

	return cfg, nil
}

func (c *Config) loadFromJSON(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read JSON config file: %w", err)
	}

	var jsonConfig JSONConfig
	if err := json.Unmarshal(data, &jsonConfig); err != nil {
		return fmt.Errorf("failed to parse JSON config: %w", err)
	}

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

	return nil
}
