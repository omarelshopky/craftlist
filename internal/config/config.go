package config

import (
    "fmt"
    "time"
)

// Config holds the application configuration
type Config struct {
    Generator GeneratorConfig `mapstructure:"generator"`
    Output    OutputConfig    `mapstructure:"output"`
}

// GeneratorConfig holds password generation settings
type GeneratorConfig struct {
    MinYear         int               `mapstructure:"min_year"`
    MaxYear         int               `mapstructure:"max_year"`
    CommonWords     []string          `mapstructure:"common_words"`
    Separators      []string          `mapstructure:"separators"`
    Substitutions   map[string]string `mapstructure:"substitutions"`
    NumberPatterns  []string          `mapstructure:"number_patterns"`
    MaxPasswordLen  int               `mapstructure:"max_password_length"`
}

// OutputConfig holds output settings
type OutputConfig struct {
    Filename     string `mapstructure:"filename"`
    ShowExamples bool   `mapstructure:"show_examples"`
    ExampleCount int    `mapstructure:"example_count"`
    Format       string `mapstructure:"format"`
}

// Load returns the default configuration
func Load() (*Config, error) {
    currentYear := time.Now().Year()
    
    return &Config{
        Generator: GeneratorConfig{
            MinYear:        2000,
            MaxYear:        currentYear + 2,
            MaxPasswordLen: 64,
            CommonWords: []string{
                "password", "admin", "guest", "wifi", "wireless", 
                "IT", "tech", "pass", "login", "user", "123", 
                "root", "default", "access", "network", "internet",
            },
            Separators: []string{"", "@", "_", "-", ".", "#", "!", "*", "+"},
            Substitutions: map[string]string{
                "a": "4", "A": "4",
                "e": "3", "E": "3",
                "i": "1", "I": "1",
                "o": "0", "O": "0",
                "s": "5", "S": "5",
                "t": "7", "T": "7",
                "l": "1", "L": "1",
                "g": "9", "G": "9",
            },
            NumberPatterns: []string{
                "1", "12", "123", "1234", "12345", "123456",
                "01", "001", "100", "1000", "2024", "2023", "2022",
            },
        },
        Output: OutputConfig{
            Filename:     "passwords.ls",
            ShowExamples: true,
            ExampleCount: 15,
            Format:       "text",
        },
    }, nil
}