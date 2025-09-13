package generator

import (
	"reflect"
	"testing"
	"sort"

	"github.com/omarelshopky/craftlist/internal/config"
)

func TestGetVariations(t *testing.T) {
	tests := []struct {
		name         string
		words        []string
		substitutions map[string][]string
		expected []string
	}{
		{
			name:    "empty input returns empty slice",
			words:   []string{},
		},
		{
			name:         "single word no substitutions",
			words:        []string{"evi"},
			substitutions: map[string][]string{},
			expected: []string{"evi", "evI", "eVi", "Evi", "EVi", "eVI", "EvI", "EVI"},
		},
		{
			name:  "word with spaces triggers word variations",
			words: []string{"e c"},
			substitutions: map[string][]string{},
			expected: []string{
				"e", "c", "e-c", "e_c", "e c", "E", "C", "E-c", "e-C", "E-C", "e_C", "E_c", "E_C", "e C", "E c", "E C", "ec", "eC", "Ec", "EC",
			},
		},
		{
			name:  "substitutions applied",
			words: []string{"evi"},
			substitutions: map[string][]string{
				"e": {"3"},
				"t": {"7"},
			},
			expected: []string{"evi", "evI", "eVi", "Evi", "EVi", "eVI", "EvI", "EVI", "3vi", "3vI", "3Vi", "3VI",},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &Generator{
				variations: NewVariationGenerator(config.GeneratorConfig{
					Substitutions: tt.substitutions,
				}),
			}

			got, _ := g.getVariations(tt.words)

			// If no expected values, skip validation
			if len(tt.expected) == 0 {
				if len(got) != 0 {
					t.Errorf("expected empty result, got %v", got)
				}
				return
			}

			sort.Strings(got)
			sort.Strings(tt.expected)

			if !reflect.DeepEqual(got, tt.expected) {
				t.Errorf("expected %v, got %v", tt.expected, got)
			}
		})
	}
}
