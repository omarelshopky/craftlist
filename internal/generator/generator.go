package generator

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"sync"

	"github.com/omarelshopky/craftlist/internal/config"
	"github.com/omarelshopky/craftlist/internal/interfaces"
)

type Generator struct {
	config      	config.GeneratorConfig
	placeholders 	config.PlaceholdersConfig
	customWords 	[]string
	commonWords 	[]string
	ssids       	[]string
	numbers 		[]string
	patterns    	*PatternProcessor
	variations  	*VariationGenerator
	output      	*OutputManager
}

func New(cfg config.GeneratorConfig, placeholders config.PlaceholdersConfig) *Generator {
	return &Generator{
		config:     	cfg,
		placeholders: 	placeholders,
		patterns:   	NewPatternProcessor(cfg, placeholders),
		variations: 	NewVariationGenerator(cfg),
		output:     	NewOutputManager(),
	}
}

func (g *Generator) SetCustomWords(words []string) {
	g.customWords = words
}

func (g *Generator) SetSSIDs(ssids []string) {
	g.ssids = ssids
}

func (g *Generator) GetCustomWords() []string {
	return g.customWords
}

func (g *Generator) GetCommonWords() []string {
	return g.commonWords
}

func (g *Generator) GetSSIDs() []string {
	return g.ssids
}

func (g *Generator) GetNumbers() []string {
	return g.numbers
}

func (g *Generator) PrepareVariations() error {
	var err error

	g.customWords, err = g.getVariations(g.customWords)
	if err != nil {
		return fmt.Errorf("failed to get custom word variations: %w", err)
	}

	g.commonWords, err = g.getVariations(g.config.CommonWords)
	if err != nil {
		return fmt.Errorf("failed to get common word variations: %w", err)
	}

	g.ssids, err = g.getVariations(g.ssids)
	if err != nil {
		return fmt.Errorf("failed to get SSID variations: %w", err)
	}

	g.numbers = g.patterns.GenerateAllNumberPatterns()

	return nil
}

func (g *Generator) Generate(ctx context.Context, outputFile string, printer interfaces.Printer) error {
	// Create output file
	writer, err := g.output.CreateWriter(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer writer.Close()

	// Setup concurrent processing
	numWorkers := runtime.NumCPU()
	jobChan := make(chan PasswordJob, 1000)
	resultChan := make(chan string, 1000)

	var wg sync.WaitGroup
	var writerWg sync.WaitGroup

	// Start workers
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			g.worker(ctx, jobChan, resultChan)
		}()
	}

	// Start writer goroutine
	writerWg.Add(1)
	passwordCount := 0

	go func() {
		defer writerWg.Done()
		for password := range resultChan {
			if err := writer.WritePassword(password); err != nil {
				return
			}
			passwordCount++

			if passwordCount%10000 == 0 {
				printer.PrintProgress(passwordCount)
				writer.Flush()
			}
		}
	}()

	// Generate jobs based on patterns
	go func() {
		defer close(jobChan)
		g.generateJobs(ctx, jobChan)
	}()

	// Wait for all workers to finish
	wg.Wait()
	close(resultChan)

	// Wait for writer to finish
	writerWg.Wait()

	printer.PrintFinalCount(passwordCount)

	return nil
}

func (g *Generator) worker(ctx context.Context, jobs <-chan PasswordJob, results chan<- string) {
	for {
		select {
		case <-ctx.Done():
			return
		case job, ok := <-jobs:
			if !ok {
				return
			}

			password := g.patterns.ProcessPattern(job)
			if password != "" && len(password) >= g.config.MinPasswordLen && len(password) <= g.config.MaxPasswordLen {
				select {
				case results <- password:
				case <-ctx.Done():
					return
				}
			}
		}
	}
}

func (g *Generator) generateJobs(ctx context.Context, jobChan chan<- PasswordJob) {
	for _, pattern := range g.config.Patterns {
		select {
		case <-ctx.Done():
			return
		default:
		}

		// Ignore patterns with SSID if not entered
		if len(g.ssids) == 0 && strings.Contains(pattern, g.placeholders.SSID.Format) {
			continue
		}

		g.generateJobsForPattern(ctx, jobChan, pattern)
	}
}

func (g *Generator) generateJobsForPattern(ctx context.Context, jobChan chan<- PasswordJob, pattern string) {
	// Determine what placeholders exist in the pattern
	hasCustom := strings.Contains(pattern, g.placeholders.CustomWord.Format)
	hasCommon := strings.Contains(pattern, g.placeholders.CommonWord.Format)
	hasSSID := strings.Contains(pattern, g.placeholders.SSID.Format)
	hasYear := strings.Contains(pattern, g.placeholders.Year.Format) || strings.Contains(pattern, g.placeholders.ShortYear.Format)
	hasNum := strings.Contains(pattern, g.placeholders.Number.Format)

	separatorCount := strings.Count(pattern, g.placeholders.Separator.Format)

	// Generate combinations based on pattern requirements
	customRange := []string{""}
	if hasCustom {
		customRange = g.customWords
	}

	commonRange := []string{""}
	if hasCommon {
		commonRange = g.commonWords
	}

	ssidRange := []string{""}
	if hasSSID {
		ssidRange = g.ssids
	}

	yearRange := []int{0}
	if hasYear {
		yearRange = make([]int, 0, g.config.MaxYear-g.config.MinYear+1)
		for y := g.config.MinYear; y <= g.config.MaxYear; y++ {
			yearRange = append(yearRange, y)
		}
	}

	numberRange := []string{""}
	if hasNum {
		numberRange = g.numbers
	}

	for _, customWord := range customRange {
		for _, commonWord := range commonRange {
			for _, ssid := range ssidRange {
				for _, year := range yearRange {
					for _, number := range numberRange {
						g.generateSeparatorCombinations(ctx, jobChan, pattern, separatorCount, customWord, commonWord, ssid, year, number, []string{})
					}
				}
			}
		}
	}
}

func (g *Generator) generateSeparatorCombinations(ctx context.Context, jobChan chan<- PasswordJob, pattern string, remainingSeparators int, customWord, commonWord, ssid string, year int, number string, separators []string) {

	// Base case: no more separators to process
	if remainingSeparators == 0 {
		select {
		case <-ctx.Done():
			return
		case jobChan <- PasswordJob{
			Pattern:    pattern,
			CustomWord: customWord,
			CommonWord: commonWord,
			SSID:       ssid,
			Year:       year,
			Number:     number,
			Separators: separators,
		}:
		}
		return
	}

	// Recursive case: generate combinations for remaining separators
	for _, separator := range g.config.Separators {
		newSeparators := make([]string, len(separators)+1)
		copy(newSeparators, separators)
		newSeparators[len(separators)] = separator

		g.generateSeparatorCombinations(ctx, jobChan, pattern, remainingSeparators-1, customWord, commonWord, ssid, year, number, newSeparators)
	}
}

func (g *Generator) getVariations(words []string) ([]string, error) {
	if len(words) == 0 {
		return []string{}, nil
	}

	wordVariations, err := g.generateVariations(words, g.variations.GenerateWordVariations)
	if err != nil {
		return nil, err
	}

	caseVariations, err := g.generateVariations(wordVariations, g.variations.GenerateCaseVariations)
	if err != nil {
		return nil, err
	}

	subVariations, err := g.generateVariations(caseVariations, g.variations.ApplyAllSubstitutions)
	if err != nil {
		return nil, err
	}

	return subVariations, nil
}

func (g *Generator) generateVariations(words []string, variationFunc func(string) []string) ([]string, error) {
	if len(words) == 0 {
		return nil, fmt.Errorf("no words provided")
	}

	variations := make(map[string]struct{})
	for _, word := range words {
		for _, variation := range variationFunc(word) {
			variations[variation] = struct{}{}
		}
	}

	return g.variations.ConvertSetToSlice(variations), nil
}