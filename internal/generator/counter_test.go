package generator

import (
	"testing"

	"github.com/omarelshopky/craftlist/internal/config"
)

func TestCountPasswords(t *testing.T) {
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
		MinYear:     2025,
		MaxYear:     2026,
		MinPasswordLen: 8,
		MaxPasswordLen: 64,
		Separators:  []string{"-", "_"},
		Patterns:    []string{"<CUSTOM><SEP><YEAR>", "<COMMON><SEP><NUM>", "<SSID>", "<CUSTOM><COMMON><NUM>"},
	}

	tests := []struct {
		name         string
		customWords  []string
		commonWords  []string
		ssids        []string
		numbers      []string
		minLength 	 int
		expected     int
	}{
		{
			name:        	"basic counts with password length limit",
			customWords: 	[]string{"evil", "corp"},
			commonWords: 	[]string{"password", "it"},
			ssids:       	[]string{"dev"},
			numbers:     	[]string{"1234", "123"},
			expected:    	(2 * 2 * 2) + (1 * 2 * 2) + (0) + (2 * 2 * 2),
		},
		{
			name:    		"basic counts without password length limit",
			customWords: 	[]string{"evil", "corp"},
			commonWords: 	[]string{"password", "it"},
			ssids:       	[]string{"dev"},
			numbers:     	[]string{"1234", "123"},
			minLength:		1,
			expected:      	(2 * 2 * 2) + (2 * 2 * 2) + (1) + (2 * 2 * 2),
		},
		{
			name:        	"no ssids",
			customWords: 	[]string{"evil", "corp"},
			commonWords: 	[]string{"password", "it"},
			ssids:       	[]string{},
			numbers:     	[]string{"1234", "123"},
			minLength:		1,
			expected:    	(2 * 2 * 2) + (2 * 2 * 2) + (0) + (2 * 2 * 2),
		},
		{
			name:        	"no patterns",
			customWords: 	[]string{"evil", "corp"},
			commonWords: 	[]string{"password", "it"},
			ssids:       	[]string{"dev"},
			numbers:     	[]string{"1234", "123"},
			expected:       0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testCfg := cfg
			if tt.name == "no patterns" {
				testCfg.Patterns = []string{}
			}

			if tt.minLength != 0 {
				testCfg.MinPasswordLen = tt.minLength
			}

			c := NewCounter(testCfg, placeholders)
			result, _ := c.CountPasswords(tt.customWords, tt.commonWords, tt.ssids, tt.numbers)

			if result != tt.expected {
				t.Errorf("expected %d, got %d", tt.expected, result)
			}
		})
	}
}
