package generator

import (
	"math"
	"strings"

	"github.com/omarelshopky/craftlist/internal/config"
)

type Counter struct {
	config config.GeneratorConfig
	placeholders config.PlaceholdersConfig
}

func NewCounter(cfg config.GeneratorConfig, placeholders config.PlaceholdersConfig) *Counter {
	return &Counter{config: cfg, placeholders: placeholders}
}

func (c *Counter) EstimatePasswordCount(customWordsCount, commonWordsCount, ssidsCount, numbersCount int) (int, map[string]int) {
	stats := make(map[string]int)
	total := 0
	yearCount := c.config.MaxYear - c.config.MinYear + 1
	separatorsCount := len(c.config.Separators)

    for _, pattern := range c.config.Patterns {
        count := 1

        if strings.Contains(pattern, c.placeholders.CustomWord.Format) {
            count *= max(customWordsCount, 1)
        }

        if strings.Contains(pattern, c.placeholders.CommonWord.Format) {
            count *= max(commonWordsCount, 1)
        }

        if strings.Contains(pattern, c.placeholders.SSID.Format) {
            count *= max(ssidsCount, 1)
        }

        if strings.Contains(pattern, c.placeholders.Number.Format) {
            count *= max(numbersCount, 1)
        }

        if strings.Contains(pattern, c.placeholders.Year.Format) {
            count *= max(yearCount, 1)
        }

		if strings.Contains(pattern, c.placeholders.ShortYear.Format) {
            count *= max(yearCount, 1)
        }

        sepCount := strings.Count(pattern, c.placeholders.Separator.Format)
        if sepCount > 0 {
            count *= pow(separatorsCount, sepCount)
        }

        total += count
		stats[pattern] = count
    }

    return total, stats
}

func pow(base, exp int) int {
	return int(math.Pow(float64(base), float64(exp)))
}
