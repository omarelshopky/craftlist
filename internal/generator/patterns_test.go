package generator

import (
	"reflect"
	"sort"
	"testing"

	"github.com/omarelshopky/craftlist/internal/config"
)

func TestProcessPattern(t *testing.T) {
	pp := NewPatternProcessor(config.GeneratorConfig{})

	tests := []struct {
		name     string
		job      PasswordJob
		expected string
	}{
		{
			name: "all placeholders replaced",
			job: PasswordJob{
				Pattern:    "<CUSTOM><SEP><COMMON><SEP><SSID><SEP><YEAR><SEP><SHORTYEAR><SEP><NUM>",
				CustomWord: "custom",
				CommonWord: "common",
				SSID:       "ssid",
				Year:       2025,
				Number:     "123",
				Separators: []string{"-", "_", "+", "!", "#"},
			},
			expected: "custom-common_ssid+2025!25#123",
		},
		{
			name: "no year and separators",
			job: PasswordJob{
				Pattern:    "<CUSTOM><COMMON><SSID><NUM>",
				CustomWord: "c",
				CommonWord: "m",
				SSID:       "s",
				Number:     "1",
			},
			expected: "cms1",
		},
		{
			name: "more <SEP> than separators provided",
			job: PasswordJob{
				Pattern:    "<CUSTOM><SEP><COMMON><SEP><SSID><SEP>",
				CustomWord: "A",
				CommonWord: "B",
				SSID:       "C",
				Separators: []string{"-"},
			},
			expected: "A-B<SEP>C<SEP>",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pp.ProcessPattern(tt.job)
			if got != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, got)
			}
		})
	}
}

func TestGenerateAllNumberPatterns(t *testing.T) {
	cfg := config.GeneratorConfig{
		NumberPatterns: []string{"12", "dd"},
	}
	pp := NewPatternProcessor(cfg)

	got := pp.GenerateAllNumberPatterns()

	// We expect "12" + all two-digit combinations (100 combinations)
	if len(got) != 1+100 {
		t.Errorf("expected 101 patterns, got %d", len(got))
	}
}

func TestExpandDigitPattern(t *testing.T) {
	pp := NewPatternProcessor(config.GeneratorConfig{})

	tests := []struct {
		name     string
		pattern  string
		expected int
	}{
		{"no digits", "1", 1},
		{"one digit", "d", 10},
		{"two digits", "dd", 100},
		{"two digits with only one expandable", "5d", 10},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := pp.expandDigitPattern(tt.pattern)
			if len(got) != tt.expected {
				t.Errorf("expected %d results, got %d", tt.expected, len(got))
			}
		})
	}
}

func TestGenerateDigitCombinations(t *testing.T) {
	pp := NewPatternProcessor(config.GeneratorConfig{})
	results := []string{}
	pp.generateDigitCombinations("d", "", 0, &results)

	sort.Strings(results)
	expected := []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9"}

	if !reflect.DeepEqual(results, expected) {
		t.Errorf("expected %v, got %v", expected, results)
	}
}
