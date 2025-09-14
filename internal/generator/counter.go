package generator

import (
	"strings"
	"sort"

	"github.com/omarelshopky/craftlist/internal/config"
)

type Counter struct {
	config       config.GeneratorConfig
	placeholders config.PlaceholdersConfig
}

func NewCounter(cfg config.GeneratorConfig, placeholders config.PlaceholdersConfig) *Counter {
	return &Counter{
		config:       cfg,
		placeholders: placeholders,
	}
}

type DistributionInfo struct {
	MinLength     int
	MaxLength     int
	LengthDist map[int]int // length -> count
	TotalCount int
}

type PatternComponent struct {
	Type   ComponentType // Type of component
	Length int           // For base components, this is the fixed length
}

type ComponentType string

const (
	ComponentBase      ComponentType = "base"
	ComponentCustom    ComponentType = "custom"
	ComponentCommon    ComponentType = "common"
	ComponentSSID      ComponentType = "ssid"
	ComponentNumber    ComponentType = "number"
	ComponentYear      ComponentType = "year"
	ComponentShortYear ComponentType = "shortyear"
	ComponentSeparator ComponentType = "separator"
)

const (
	defaultMaxLength = 1000000 // Default max length when no upper limit is specified
	lengthBuffer     = 50      // Buffer for optimization in combination counting
	yearLength       = 4       // Standard year length
	shortYearLength  = 2       // Short year length
)

// CountPasswords calculates the total number of possible passwords and per-pattern statistics
func (c *Counter) CountPasswords(customWords, commonWords, ssids, numbers []string) (int, map[string]int) {
	stats := make(map[string]int)
	total := 0

	// Pre-compute statistics for all word lists
	wordListStats := c.buildWordListStats(customWords, commonWords, ssids, numbers)
	
	// Calculate statistics for each pattern
	for _, pattern := range c.config.Patterns {
		count := c.calculatePatternCount(pattern, wordListStats)
		total += count
		stats[pattern] = count
	}

	return total, stats
}

type wordListStats struct {
	custom    *DistributionInfo
	common    *DistributionInfo
	ssid      *DistributionInfo
	number    *DistributionInfo
	year      *DistributionInfo
	shortYear *DistributionInfo
	separator *DistributionInfo
}

func (c *Counter) buildWordListStats(customWords, commonWords, ssids, numbers []string) *wordListStats {
	yearCount := c.config.MaxYear - c.config.MinYear + 1
	
	return &wordListStats{
		custom:    c.buildDistributionInfo(customWords),
		common:    c.buildDistributionInfo(commonWords),
		ssid:      c.buildDistributionInfo(ssids),
		number:    c.buildDistributionInfo(numbers),
		separator: c.buildDistributionInfo(c.config.Separators),
		year:      c.createFixedDistributionInfo(yearLength, yearCount),
		shortYear: c.createFixedDistributionInfo(shortYearLength, yearCount),
	}
}

func (c *Counter) buildDistributionInfo(words []string) *DistributionInfo {
	if len(words) == 0 {
		return &DistributionInfo{LengthDist: make(map[int]int)}
	}
	
	stats := &DistributionInfo{
		MinLength:    len(words[0]),
		MaxLength:    len(words[0]),
		LengthDist: make(map[int]int),
		TotalCount:   len(words),
	}
	
	for _, word := range words {
		length := len(word)
		stats.LengthDist[length]++
		
		if length < stats.MinLength {
			stats.MinLength = length
		}
		if length > stats.MaxLength {
			stats.MaxLength = length
		}
	}
	
	return stats
}

func (c *Counter) createFixedDistributionInfo(length, count int) *DistributionInfo {
	return &DistributionInfo{
		MinLength:    length,
		MaxLength:    length,
		LengthDist: map[int]int{length: count},
		TotalCount:   count,
	}
}

// calculatePatternCount calculates the number of valid passwords for a given pattern
func (c *Counter) calculatePatternCount(pattern string, wordStats *wordListStats) int {
	components := c.parsePattern(pattern)
	componentInfos := c.buildComponentInfos(components, wordStats)
	
	return c.countValidCombinations(componentInfos, c.config.MinPasswordLen, c.config.MaxPasswordLen)
}

// buildComponentInfos converts pattern components to DistributionInfo structs
func (c *Counter) buildComponentInfos(components []PatternComponent, wordStats *wordListStats) []*DistributionInfo {
	var componentInfos []*DistributionInfo
	
	for _, comp := range components {
		var compInfo *DistributionInfo
		
		switch comp.Type {
		case ComponentBase:
			compInfo = c.createFixedDistributionInfo(comp.Length, 1)
		case ComponentCustom:
			compInfo = wordStats.custom
		case ComponentCommon:
			compInfo = wordStats.common
		case ComponentSSID:
			compInfo = wordStats.ssid
		case ComponentNumber:
			compInfo = wordStats.number
		case ComponentYear:
			compInfo = wordStats.year
		case ComponentShortYear:
			compInfo = wordStats.shortYear
		case ComponentSeparator:
			compInfo = wordStats.separator
		default:
			continue
		}
		
		if compInfo != nil {
			componentInfos = append(componentInfos, compInfo)
		}
	}
	
	return componentInfos
}

func (c *Counter) parsePattern(pattern string) []PatternComponent {
	var components []PatternComponent
	remaining := pattern
	
	placeholderMap := c.buildPlaceholderMap()
	
	for len(remaining) > 0 {
		if comp, consumed := c.tryMatchPlaceholder(remaining, placeholderMap); consumed > 0 {
			components = append(components, comp)
			remaining = remaining[consumed:]
		} else if baseLength := c.countBaseCharacters(remaining, placeholderMap); baseLength > 0 {
			components = append(components, PatternComponent{
				Type:   ComponentBase,
				Length: baseLength,
			})
			remaining = remaining[baseLength:]
		} else {
			// Skip unrecognized character
			remaining = remaining[1:]
		}
	}
	
	return components
}

func (c *Counter) buildPlaceholderMap() map[ComponentType]string {
	return map[ComponentType]string{
		ComponentCustom:    c.placeholders.CustomWord.Format,
		ComponentCommon:    c.placeholders.CommonWord.Format,
		ComponentSSID:      c.placeholders.SSID.Format,
		ComponentNumber:    c.placeholders.Number.Format,
		ComponentYear:      c.placeholders.Year.Format,
		ComponentShortYear: c.placeholders.ShortYear.Format,
		ComponentSeparator: c.placeholders.Separator.Format,
	}
}

func (c *Counter) tryMatchPlaceholder(text string, placeholderMap map[ComponentType]string) (PatternComponent, int) {
	// Create deterministic order by sorting component types
	var types []ComponentType
	for compType := range placeholderMap {
		types = append(types, compType)
	}
	sort.Slice(types, func(i, j int) bool {
		return string(types[i]) < string(types[j])
	})
	
	// Try matching in deterministic order
	for _, compType := range types {
		placeholder := placeholderMap[compType]
		if strings.HasPrefix(text, placeholder) {
			return PatternComponent{Type: compType, Length: 0}, len(placeholder)
		}
	}

	return PatternComponent{}, 0
}

func (c *Counter) countBaseCharacters(text string, placeholderMap map[ComponentType]string) int {
	// Create sorted list of placeholders for deterministic behavior
	var placeholders []string
	for _, placeholder := range placeholderMap {
		placeholders = append(placeholders, placeholder)
	}
	sort.Strings(placeholders)
	
	for i := range text {
		for _, placeholder := range placeholders {
			if strings.HasPrefix(text[i:], placeholder) {
				return i
			}
		}
	}

	return len(text)
}

func (c *Counter) combineDistributions(dist1 map[int]int, dist2 map[int]int, maxLen int) map[int]int {
	nextDist := make(map[int]int)
	
	// Get sorted keys for deterministic iteration
	var keys1, keys2 []int
	for len1 := range dist1 {
		keys1 = append(keys1, len1)
	}
	for len2 := range dist2 {
		keys2 = append(keys2, len2)
	}
	sort.Ints(keys1)
	sort.Ints(keys2)
	
	// Iterate in sorted order
	for _, len1 := range keys1 {
		count1 := dist1[len1]
		for _, len2 := range keys2 {
			count2 := dist2[len2]
			totalLen := len1 + len2
			if totalLen <= maxLen+lengthBuffer {
				nextDist[totalLen] += count1 * count2
			}
		}
	}
	
	return nextDist
}

// countValidCombinations uses dynamic programming to count combinations
// that result in lengths within the specified range
func (c *Counter) countValidCombinations(components []*DistributionInfo, minLen, maxLen int) int {
	if len(components) == 0 {
		return 0
	}
	
	if maxLen == 0 {
		maxLen = defaultMaxLength
	}
	
	// Initialize with first component's distribution
	currentDist := make(map[int]int)
	for length, count := range components[0].LengthDist {
		currentDist[length] = count
	}
	
	// Combine with remaining components using convolution
	for i := 1; i < len(components); i++ {
		currentDist = c.combineDistributions(currentDist, components[i].LengthDist, maxLen)
	}
	
	return c.countValidLengths(currentDist, minLen, maxLen)
}

// countValidLengths counts combinations with lengths in the valid range
func (c *Counter) countValidLengths(distribution map[int]int, minLen, maxLen int) int {
	validCount := 0
	
	for length, count := range distribution {
		if length >= minLen && length <= maxLen {
			validCount += count
		}
	}

	return validCount
}