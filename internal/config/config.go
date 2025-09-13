package config

import (
	"encoding/json"
	"fmt"
	"os"
	"time"
	"regexp"
	"strings"
	"reflect"
)

var Colors = struct {
    Reset  string
    Red    string  
    Yellow string
    Green  string
    Cyan   string
    Bold   string
}{
    Reset:  "\033[0m",
    Red:    "\033[31m", 
    Yellow: "\033[33m",
    Green:  "\033[32m",
    Cyan:   "\033[36m",
    Bold:   "\033[1m",
}

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

type PlaceholdersConfig struct {
	CustomWord Placeholder `mapstructure:"custom_word" json:"custom_word"`
	CommonWord Placeholder `mapstructure:"common_word" json:"common_word"`
	SSID       Placeholder `mapstructure:"ssid" json:"ssid"`
	Separator  Placeholder `mapstructure:"separator" json:"separator"`
	Year       Placeholder `mapstructure:"year" json:"year"`
	ShortYear  Placeholder `mapstructure:"short_year" json:"short_year"`
	Number     Placeholder `mapstructure:"number" json:"number"`
}

type OutputConfig struct {
	Filename string `mapstructure:"filename" json:"filename"`
}

type Placeholder struct {
	Format      string `mapstructure:"format" json:"format"`
	Description string `mapstructure:"description" json:"description"`
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
		Generator:    newDefaultGeneratorConfig(),
		Placeholders: newDefaultPlaceholdersConfig(),
		Output:       newDefaultOutputConfig(),
	}

	if jsonConfigPath != "" {
		if err := cfg.loadFromJSON(jsonConfigPath); err != nil {
			return nil, fmt.Errorf("failed to load JSON config: %w", err)
		}
	}

	return cfg, nil
}

func newDefaultGeneratorConfig() GeneratorConfig {
	return GeneratorConfig{
		MinYear:        1990,
		MaxYear:        time.Now().Year(),
		MinPasswordLen: 8,
		MaxPasswordLen: 64,
		CommonWords:    getDefaultCommonWords(),
		Separators:     getDefaultSeparators(),
		Substitutions:  getDefaultSubstitutions(),
		NumberPatterns: getDefaultNumberPatterns(),
		Patterns:       getDefaultPatterns(),
	}
}

func newDefaultPlaceholdersConfig() PlaceholdersConfig {
	return PlaceholdersConfig{
		CustomWord: Placeholder{
			Format:      "<CUSTOM>",
			Description: "Inserts custom word variations from the file specified with the --words flag",
		},
		CommonWord: Placeholder{
			Format:      "<COMMON>",
			Description: "Inserts common word variations based on the list defined in your config file",
		},
		SSID: Placeholder{
			Format:      "<SSID>",
			Description: "Inserts SSID variations from the file specified with the --ssids flag",
		},
		Separator: Placeholder{
			Format:      "<SEP>",
			Description: "Inserts separators based on the list defined in your config file",
		},
		Year: Placeholder{
			Format:      "<YEAR>",
			Description: "Inserts full year based on the range defined in flags or config file (e.g., 2025)",
		},
		ShortYear: Placeholder{
			Format:      "<SHORTYEAR>",
			Description: "Inserts two-digit year based on the range defined in flags or config file (e.g., 25)",
		},
		Number: Placeholder{
			Format:      "<NUM>",
			Description: "Inserts numbers based on the list defined in your config file",
		},
	}
}

func newDefaultOutputConfig() OutputConfig {
	return OutputConfig{
		Filename: "passwords.txt",
	}
}

func getDefaultCommonWords() []string {
	return []string{
		"password", "admin", "guest", "wifi", "wireless", "IT", "tech",
		"pass", "login", "user", "root", "default", "access",
		"network", "internet", "secure", "temp", "test",
	}
}

func getDefaultSeparators() []string {
	return []string{"", "@", "_", "-", ".", "#", "!", "*", "+", "~", "%", "&", "^"}
}

func getDefaultSubstitutions() map[string][]string {
	return map[string][]string{
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
	}
}

func getDefaultNumberPatterns() []string {
	return []string{
		"d", "dd", "ddd", "dddd", "12345", "123456",
	}
}

func getDefaultPatterns() []string {
	return []string{
		"<CUSTOM>",
		"<COMMON>",
		"<SSID>",
		"<CUSTOM><SEP><YEAR>",
		"<CUSTOM><SEP><SHORTYEAR>",
		"<CUSTOM><SEP><NUM>",
		"<CUSTOM><SEP><COMMON>",
		"<COMMON><SEP><YEAR>",
		"<COMMON><SEP><SHORTYEAR>",
		"<COMMON><SEP><NUM>",
		"<COMMON><SEP><CUSTOM>",
		"<SSID><SEP><YEAR>",
		"<SSID><SEP><SHORTYEAR>",
		"<SSID><SEP><NUM>",
		"<SSID><SEP><CUSTOM>",
		"<YEAR><SEP><CUSTOM>",
		"<YEAR><SEP><COMMON>",
		"<YEAR><SEP><SSID>",
		"<SHORTYEAR><SEP><CUSTOM>",
		"<SHORTYEAR><SEP><COMMON>",
		"<SHORTYEAR><SEP><SSID>",
		"<NUM><SEP><CUSTOM>",
		"<NUM><SEP><COMMON>",
		"<NUM><SEP><SSID>",
		"<SEP><CUSTOM><SEP><SSID><SEP><YEAR>",
		"<COMMON><SEP><CUSTOM><YEAR>",
	}
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

func (c *Config) Validate() error {
	if c.Generator.MinYear > c.Generator.MaxYear {
		return fmt.Errorf("min year (%d) cannot be greater than max year (%d)", 
			c.Generator.MinYear, c.Generator.MaxYear)
	}
	
	if c.Generator.MinPasswordLen > c.Generator.MaxPasswordLen {
		return fmt.Errorf("min password length (%d) cannot be greater than max password length (%d)", 
			c.Generator.MinPasswordLen, c.Generator.MaxPasswordLen)
	}
	
	if c.Generator.MinPasswordLen < 1 {
		return fmt.Errorf("min password length must be at least 1")
	}
	
	if c.Output.Filename == "" {
		return fmt.Errorf("output filename cannot be empty")
	}

	if err := c.validatePatterns(); err != nil {
		return err
	}
	
	return nil
}

func (c *Config) validatePatterns() error {
	if len(c.Generator.Patterns) == 0 {
		return fmt.Errorf("no patterns defined in configuration")
	}

	knownPlaceholders := c.getKnownPlaceholders()
	var validationErrors []string
	var hasErrors bool

	for idx, pattern := range c.Generator.Patterns {
		if unknownPlaceholders := c.findUnknownPlaceholders(pattern, knownPlaceholders); len(unknownPlaceholders) > 0 {
			hasErrors = true
			errorMsg := fmt.Sprintf("Pattern %d: %s contains unknown placeholders: %s%s%s",
				idx+1,
				c.highlightPattern(pattern, unknownPlaceholders),
				Colors.Red,
				strings.Join(unknownPlaceholders, ", "),
				Colors.Reset)
			validationErrors = append(validationErrors, errorMsg)
		}
	}

	if hasErrors {
		fmt.Printf("%sPattern validation failed:%s\n\n", Colors.Red, Colors.Reset)
		for _, err := range validationErrors {
			fmt.Printf("  %s\n", err)
		}
		
		return fmt.Errorf("%d pattern(s) contain unknown placeholders", len(validationErrors))
	}

	fmt.Printf("%sAll patterns validated successfully%s\n\n", Colors.Green, Colors.Reset)
	
	return nil
}

func (c *Config) getKnownPlaceholders() map[string]bool {
	known := make(map[string]bool)
	
	values := reflect.ValueOf(c.Placeholders)
	for idx := 0; idx < values.NumField(); idx++ {
		field := values.Field(idx)
		if field.CanInterface() {
			if placeholder, ok := field.Interface().(Placeholder); ok {
				known[placeholder.Format] = true
			}
		}
	}
	
	return known
}

func (c *Config) findUnknownPlaceholders(pattern string, knownPlaceholders map[string]bool) []string {
	placeholderRegex := regexp.MustCompile(`<[^>]+>`)
	matches := placeholderRegex.FindAllString(pattern, -1)
	
	var unknown []string
	seen := make(map[string]bool)
	
	for _, match := range matches {
		if !knownPlaceholders[match] && !seen[match] {
			unknown = append(unknown, match)
			seen[match] = true
		}
	}
	
	return unknown
}

func (c *Config) highlightPattern(pattern string, unknownPlaceholders []string) string {
	highlighted := pattern
	
	for _, unknown := range unknownPlaceholders {
		// Escape special regex characters in the placeholder
		escapedUnknown := regexp.QuoteMeta(unknown)
		re := regexp.MustCompile(escapedUnknown)
		highlighted = re.ReplaceAllString(highlighted, Colors.Red + unknown + Colors.Reset)
	}
	
	return highlighted
}