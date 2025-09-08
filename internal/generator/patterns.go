package generator

import (
    "strconv"
    "strings"

    "github.com/omarelshopky/craftlist/internal/config"
)

// PatternGenerator handles different password pattern generation
type PatternGenerator struct {
    config config.GeneratorConfig
}

// NewPatternGenerator creates a new PatternGenerator
func NewPatternGenerator(cfg config.GeneratorConfig) *PatternGenerator {
    return &PatternGenerator{config: cfg}
}

// GenerateVariations creates all variations for a base word
func (pg *PatternGenerator) GenerateVariations(baseWord string) []string {
    var passwords []string

    // Clean the base word (remove spaces, special chars)
    cleanWords := pg.cleanWord(baseWord)

    for _, word := range cleanWords {
        caseVariations := pg.generateCaseVariations(word)
        
        for _, caseWord := range caseVariations {
            substitutionVariations := pg.applySubstitutions(caseWord)
            
            for _, subWord := range substitutionVariations {
                // Base word alone
                passwords = append(passwords, subWord)
                
                // With years
                passwords = append(passwords, pg.generateYearCombinations(subWord)...)
                
                // With common words
                passwords = append(passwords, pg.generateCommonWordCombinations(subWord)...)
                
                // With numbers
                passwords = append(passwords, pg.generateNumberCombinations(subWord)...)
            }
        }
    }

    return passwords
}

func (pg *PatternGenerator) cleanWord(word string) []string {
    var cleaned []string
    
    // Original
    cleaned = append(cleaned, word)
    
    // Without spaces
    noSpaces := strings.ReplaceAll(word, " ", "")
    if noSpaces != word {
        cleaned = append(cleaned, noSpaces)
    }
    
    // Individual words
    parts := strings.Fields(word)
    cleaned = append(cleaned, parts...)
    
    return cleaned
}

func (pg *PatternGenerator) generateCaseVariations(word string) []string {
    return []string{
        word,
        strings.ToLower(word),
        strings.ToUpper(word),
        strings.Title(word),
    }
}

func (pg *PatternGenerator) applySubstitutions(word string) []string {
    variations := []string{word}
    
    // Apply each substitution
    for original, replacement := range pg.config.Substitutions {
        if strings.Contains(word, original) {
            substituted := strings.ReplaceAll(word, original, replacement)
            if substituted != word {
                variations = append(variations, substituted)
            }
        }
    }
    
    return variations
}

func (pg *PatternGenerator) generateYearCombinations(word string) []string {
    var combinations []string
    
    for year := pg.config.MinYear; year <= pg.config.MaxYear; year++ {
        yearStr := strconv.Itoa(year)
        shortYear := yearStr[2:]
        
        for _, sep := range pg.config.Separators {
            combinations = append(combinations,
                word+sep+yearStr,
                word+sep+shortYear,
                yearStr+sep+word,
                shortYear+sep+word,
            )
        }
    }
    
    return combinations
}

func (pg *PatternGenerator) generateCommonWordCombinations(word string) []string {
    var combinations []string
    
    for _, commonWord := range pg.config.CommonWords {
        commonVariations := pg.generateCaseVariations(commonWord)
        for _, common := range commonVariations {
            for _, sep := range pg.config.Separators {
                combinations = append(combinations,
                    word+sep+common,
                    common+sep+word,
                )
            }
        }
    }
    
    return combinations
}

func (pg *PatternGenerator) generateNumberCombinations(word string) []string {
    var combinations []string
    
    for _, number := range pg.config.NumberPatterns {
        for _, sep := range pg.config.Separators {
            combinations = append(combinations,
                word+sep+number,
                number+sep+word,
            )
        }
    }
    
    return combinations
}