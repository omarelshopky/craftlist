package generator

import (
	"strings"

	"github.com/omarelshopky/craftlist/internal/config"
)

type VariationGenerator struct {
	config config.GeneratorConfig
}

func NewVariationGenerator(cfg config.GeneratorConfig) *VariationGenerator {
	return &VariationGenerator{config: cfg}
}

func (vg *VariationGenerator) GenerateWordVariations(baseWord string) []string {
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

	return vg.deduplicate(variations)
}

func (vg *VariationGenerator) GenerateCaseVariations(word string) []string {
	if len(word) == 0 {
		return []string{}
	}

	variations := []string{""}

	for _, ch := range word {
		var newVariations []string
		lower := strings.ToLower(string(ch))
		upper := strings.ToUpper(string(ch))

		for _, v := range variations {
			newVariations = append(newVariations, v+lower)
			newVariations = append(newVariations, v+upper)
		}

		variations = newVariations
	}

	return vg.deduplicate(variations)
}

func (vg *VariationGenerator) ApplyAllSubstitutions(word string) []string {
	variations := make(map[string]struct{})
	variations[word] = struct{}{} // Original word

	// Get all possible substitution combinations
	vg.generateSubstitutionCombinations(word, "", 0, variations)

	return vg.ConvertSetToSlice(variations)
}

func (vg *VariationGenerator) generateSubstitutionCombinations(original, current string, index int, variations map[string]struct{}) {
	if index == len(original) {
		variations[current] = struct{}{}
		return
	}

	char := string(original[index])

	// Option 1: Keep original character
	vg.generateSubstitutionCombinations(original, current+char, index+1, variations)

	// Option 2: Apply substitutions if available
	if substitutes, exists := vg.config.Substitutions[char]; exists {
		for _, substitute := range substitutes {
			vg.generateSubstitutionCombinations(original, current+substitute, index+1, variations)
		}
	}
}

func (vg *VariationGenerator) ConvertSetToSlice(set map[string]struct{}) []string {
	result := make([]string, 0, len(set))
	for value := range set {
		result = append(result, value)
	}

	return result
}

func (vg *VariationGenerator) deduplicate(slice []string) []string {
	set := make(map[string]struct{})
	for _, value := range slice {
		set[value] = struct{}{}
	}
	return vg.ConvertSetToSlice(set)
}