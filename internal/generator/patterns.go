package generator

import (
	"strconv"
	"strings"

	"github.com/omarelshopky/craftlist/internal/config"
)

type PatternProcessor struct {
	config config.GeneratorConfig
}

type PasswordJob struct {
	Pattern    string
	CustomWord string
	CommonWord string
	SSID       string
	Year       int
	Number     string
	Separators []string
}

func NewPatternProcessor(cfg config.GeneratorConfig) *PatternProcessor {
	return &PatternProcessor{config: cfg}
}

func (pp *PatternProcessor) ProcessPattern(job PasswordJob) string {
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

func (pp *PatternProcessor) GenerateAllNumberPatterns() []string {
	var allNumbers []string

	for _, pattern := range pp.config.NumberPatterns {
		if strings.Contains(pattern, "d") {
			// Handle digit replacement
			expanded := pp.expandDigitPattern(pattern)
			allNumbers = append(allNumbers, expanded...)
		} else {
			allNumbers = append(allNumbers, pattern)
		}
	}

	return allNumbers
}

func (pp *PatternProcessor) expandDigitPattern(pattern string) []string {
	if !strings.Contains(pattern, "d") {
		return []string{pattern}
	}

	var results []string
	pp.generateDigitCombinations(pattern, "", 0, &results)
	
	return results
}

func (pp *PatternProcessor) generateDigitCombinations(pattern, current string, index int, results *[]string) {
	if index == len(pattern) {
		*results = append(*results, current)
		return
	}

	if pattern[index] == 'd' {
		// Replace with digits 0-9
		for digit := 0; digit <= 9; digit++ {
			pp.generateDigitCombinations(pattern, current+strconv.Itoa(digit), index+1, results)
		}
	} else {
		// Keep the original character
		pp.generateDigitCombinations(pattern, current+string(pattern[index]), index+1, results)
	}
}