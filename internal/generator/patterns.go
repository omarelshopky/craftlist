package generator

import (
	"strconv"
	"strings"

	"github.com/omarelshopky/craftlist/internal/config"
)

type PatternGenerator struct {
	config config.GeneratorConfig
}

func NewPatternGenerator(cfg config.GeneratorConfig) *PatternGenerator {
	return &PatternGenerator{config: cfg}
}

func (pg *PatternGenerator) ProcessPattern(job PasswordJob) string {
	password := job.Pattern

	password = strings.ReplaceAll(password, "<CUSTOM>", job.CustomWord)
	password = strings.ReplaceAll(password, "<COMMON>", job.CommonWord)
	password = strings.ReplaceAll(password, "<SSID>", job.SSID)
	password = strings.ReplaceAll(password, "<NUM>", job.Number)

	if job.Year > 0 {
		yearStr := strconv.Itoa(job.Year)
		shortYear := yearStr[2:]
		password = strings.ReplaceAll(password, "<YEAR>", yearStr)
		password = strings.ReplaceAll(password, "<SHORTYEAR>", shortYear)
	}

	// Handle multiple separators by replacing each <SEP> sequentially
	if strings.Contains(password, "<SEP>") {
		separatorIndex := 0
		for strings.Contains(password, "<SEP>") && separatorIndex < len(job.Separators) {
			password = strings.Replace(password, "<SEP>", job.Separators[separatorIndex], 1)
			separatorIndex++
		}
	}

	return password
}

func (pg *PatternGenerator) GenerateAllNumberPatterns() []string {
	var allNumbers []string

	for _, pattern := range pg.config.NumberPatterns {
		if strings.Contains(pattern, "d") {
			// Handle digit replacement
			expanded := pg.expandDigitPattern(pattern)
			allNumbers = append(allNumbers, expanded...)
		} else {
			allNumbers = append(allNumbers, pattern)
		}
	}

	return allNumbers
}

func (pg *PatternGenerator) expandDigitPattern(pattern string) []string {
	if !strings.Contains(pattern, "d") {
		return []string{pattern}
	}

	var results []string
	pg.generateDigitCombinations(pattern, "", 0, &results)
	return results
}

func (pg *PatternGenerator) generateDigitCombinations(pattern, current string, index int, results *[]string) {
	if index == len(pattern) {
		*results = append(*results, current)
		return
	}

	if pattern[index] == 'd' {
		// Replace with digits 0-9
		for digit := 0; digit <= 9; digit++ {
			pg.generateDigitCombinations(pattern, current+strconv.Itoa(digit), index+1, results)
		}
	} else {
		// Keep the original character
		pg.generateDigitCombinations(pattern, current+string(pattern[index]), index+1, results)
	}
}

func (pg *PatternGenerator) GenerateWordVariations(baseWord string) []string {
	var variations []string

	// Original
	variations = append(variations, baseWord)

	// Without spaces
	noSpaces := strings.ReplaceAll(baseWord, " ", "")
	if noSpaces != baseWord && noSpaces != "" {
		variations = append(variations, noSpaces)
	}

	// With underscores instead of spaces
	withUnderscores := strings.ReplaceAll(baseWord, " ", "_")
	if withUnderscores != baseWord && withUnderscores != "" {
		variations = append(variations, withUnderscores)
	}

	// With dash instead of spaces
	withDash := strings.ReplaceAll(baseWord, " ", "-")
	if withDash != baseWord && withDash != "" {
		variations = append(variations, withDash)
	}

	// Individual words
	parts := strings.Fields(baseWord)
	for _, part := range parts {
		if len(part) > 0 {
			variations = append(variations, part)
		}
	}

	return variations
}

func (pg *PatternGenerator) GenerateCaseVariations(word string) []string {
	if len(word) == 0 {
		return []string{}
	}

	// Limit case variations for performance - use common patterns instead of all combinations
	variations := []string{
		strings.ToLower(word),
		strings.ToUpper(word),
		strings.Title(word),
	}

	// Add first letter uppercase
	if len(word) > 1 {
		firstUpper := strings.ToUpper(string(word[0])) + strings.ToLower(word[1:])
		variations = append(variations, firstUpper)
	}

	return pg.deduplicate(variations)
}

func (pg *PatternGenerator) ApplyAllSubstitutions(word string) []string {
	variations := make(map[string]struct{})
	variations[word] = struct{}{} // Original word

	// Get all possible substitution combinations
	pg.generateSubstitutionCombinations(word, "", 0, variations)

	return pg.ConvertSetToSlice(variations)
}

func (pg *PatternGenerator) generateSubstitutionCombinations(original, current string, index int, variations map[string]struct{}) {
	if index == len(original) {
		variations[current] = struct{}{}
		return
	}

	char := string(original[index])

	// Option 1: Keep original character
	pg.generateSubstitutionCombinations(original, current+char, index+1, variations)

	// Option 2: Apply substitutions if available
	if substitutes, exists := pg.config.Substitutions[char]; exists {
		for _, substitute := range substitutes {
			pg.generateSubstitutionCombinations(original, current+substitute, index+1, variations)
		}
	}
}

func (pg *PatternGenerator) ConvertSetToSlice(set map[string]struct{}) []string {
	result := make([]string, 0, len(set))
	for value := range set {
		result = append(result, value)
	}

	return result
}

func (pg *PatternGenerator) deduplicate(slice []string) []string {
	set := make(map[string]struct{})
	for _, value := range slice {
		set[value] = struct{}{}
	}
	return pg.ConvertSetToSlice(set)
}
