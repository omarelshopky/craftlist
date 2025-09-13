package generator

import (
	"testing"

	"github.com/omarelshopky/craftlist/internal/config"
)

func TestEstimatePasswordCount(t *testing.T) {
	placeholders := config.PlaceholdersConfig{
		CustomWord:  config.Placeholder{Format: "<CUSTOM>"},
		CommonWord:  config.Placeholder{Format: "<COMMON>"},
		SSID:        config.Placeholder{Format: "<SSID>"},
		Number:      config.Placeholder{Format: "<NUM>"},
		Year:        config.Placeholder{Format: "<YEAR>"},
		ShortYear:   config.Placeholder{Format: "<SHORTYEAR>"},
		Separator:   config.Placeholder{Format: "<SEP>"},
	}

	cfg := config.GeneratorConfig{
		MinYear:     2020,
		MaxYear:     2021,
		Separators:  []string{"-", "_"},
		Patterns:    []string{"<CUSTOM><SEP><YEAR>", "<COMMON><SEP><NUM>", "<SSID>", "<CUSTOM><COMMON><NUM>"},
	}

	tests := []struct {
		name              string
		customWordsCount  int
		commonWordsCount  int
		ssidsCount        int
		numbersCount      int
		expected          int
	}{
		{
			name:             "basic counts",
			customWordsCount: 2,
			commonWordsCount: 2,
			ssidsCount:       2,
			numbersCount:     2,
			expected:         (2 * 2 * 2) + (2 * 2 * 2) + (2) + (2 * 2 * 2),
		},
		{
			name:             "no custom words",
			customWordsCount: 0,
			commonWordsCount: 2,
			ssidsCount:       1,
			numbersCount:     2,
			expected:         (1 * 2 * 2) + (2 * 2 * 2) + (1) + (1 * 2 * 2),
		},
		{
			name:             "no patterns",
			customWordsCount: 2,
			commonWordsCount: 2,
			ssidsCount:       2,
			numbersCount:     2,
			expected:         0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clone config to reset patterns for each test
			testCfg := cfg
			if tt.name == "no patterns" {
				testCfg.Patterns = []string{}
			}

			c := NewCounter(testCfg, placeholders)
			result := c.EstimatePasswordCount(tt.customWordsCount, tt.commonWordsCount, tt.ssidsCount, tt.numbersCount)
			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}
