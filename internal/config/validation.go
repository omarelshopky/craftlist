package config

import (
	"fmt"
	"reflect"
	"regexp"
	"strings"

	"github.com/omarelshopky/craftlist/internal/ui"
)

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

	colors := ui.DefaultColors
	knownPlaceholders := c.getKnownPlaceholders()
	var validationErrors []string
	var hasErrors bool

	for idx, pattern := range c.Generator.Patterns {
		if unknownPlaceholders := c.findUnknownPlaceholders(pattern, knownPlaceholders); len(unknownPlaceholders) > 0 {
			hasErrors = true
			errorMsg := fmt.Sprintf("Pattern %d: %s contains unknown placeholders: %s%s%s",
				idx+1,
				c.highlightPattern(pattern, unknownPlaceholders, colors),
				colors.Red,
				strings.Join(unknownPlaceholders, ", "),
				colors.Reset)
			validationErrors = append(validationErrors, errorMsg)
		}
	}

	if hasErrors {
		fmt.Printf("%sPattern validation failed:%s\n", colors.Red, colors.Reset)
		for _, err := range validationErrors {
			fmt.Printf("  %s\n", err)
		}

		return fmt.Errorf("%d pattern(s) contain unknown placeholders", len(validationErrors))
	}

	fmt.Printf("%sAll patterns validated successfully%s\n\n", colors.Green, colors.Reset)

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

func (c *Config) highlightPattern(pattern string, unknownPlaceholders []string, colors ui.Colors) string {
	highlighted := pattern

	for _, unknown := range unknownPlaceholders {
		// Escape special regex characters in the placeholder
		escapedUnknown := regexp.QuoteMeta(unknown)
		re := regexp.MustCompile(escapedUnknown)
		highlighted = re.ReplaceAllString(highlighted, colors.Red+unknown+colors.Reset)
	}

	return highlighted
}