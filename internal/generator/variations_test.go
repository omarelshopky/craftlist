package generator

import (
	"reflect"
	"sort"
	"testing"

	"github.com/omarelshopky/craftlist/internal/config"
)

func TestGenerateWordVariations(t *testing.T) {
	vg := NewVariationGenerator(config.GeneratorConfig{})

	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single word",
			input:    "evil",
			expected: []string{"evil"},
		},
		{
			name:     "word with spaces",
			input:    "evil corp",
			expected: []string{"evil corp", "evilcorp", "evil_corp", "evil-corp", "evil", "corp"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []string{""},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := vg.GenerateWordVariations(tt.input)
			sort.Strings(got)
			sort.Strings(tt.expected)

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestGenerateCaseVariations(t *testing.T) {
	vg := NewVariationGenerator(config.GeneratorConfig{})

	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single letter",
			input:    "a",
			expected: []string{"a", "A"},
		},
		{
			name:     "two letters",
			input:    "ab",
			expected: []string{"ab", "aB", "Ab", "AB"},
		},
		{
			name:     "three letters",
			input:    "evi",
			expected: []string{"evi", "evI", "eVi", "Evi", "EVi", "eVI", "EvI", "EVI"},
		},
		{
			name:     "empty string",
			input:    "",
			expected: []string{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := vg.GenerateCaseVariations(tt.input)
			sort.Strings(got)
			sort.Strings(tt.expected)

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestApplyAllSubstitutions(t *testing.T) {
	cfg := config.GeneratorConfig{
		Substitutions: map[string][]string{
			"a": {"4", "@"},
			"e": {"3"},
		},
	}
	vg := NewVariationGenerator(cfg)

	tests := []struct {
		name  string
		input string
		expected []string
	}{
		{
			name:  "word with substitutions",
			input: "ae",
			expected: []string{"ae", "4e", "43", "@e", "@3", "a3"},
		},
		{
			name:  "word without substitutions",
			input: "xyz",
			expected: []string{"xyz"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := vg.ApplyAllSubstitutions(tt.input)
			sort.Strings(got)
			sort.Strings(tt.expected)

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}

func TestConvertSetToSlice(t *testing.T) {
	vg := NewVariationGenerator(config.GeneratorConfig{})

	set := map[string]struct{}{
		"a": {},
		"b": {},
	}
	got := vg.ConvertSetToSlice(set)

	if len(got) != 2 {
		t.Errorf("expected 2 elements, got %d", len(got))
	}
}

func TestDeduplicate(t *testing.T) {
	vg := NewVariationGenerator(config.GeneratorConfig{})

	input := []string{"a", "b", "a"}
	got := vg.deduplicate(input)

	if len(got) != 2 {
		t.Errorf("expected 2 unique elements, got %d", len(got))
	}
}
