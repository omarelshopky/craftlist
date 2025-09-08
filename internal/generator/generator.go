package generator

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"runtime"
	"strings"
	"sync"

	"github.com/omarelshopky/craftlist/internal/config"
	"github.com/omarelshopky/craftlist/pkg/wordlist"
)

type Generator struct {
	config   config.GeneratorConfig
	wordlist *wordlist.Wordlist
	patterns *PatternGenerator
}

func New(cfg config.GeneratorConfig) *Generator {
	return &Generator{
		config:   cfg,
		wordlist: wordlist.New(),
		patterns: NewPatternGenerator(cfg),
	}
}

func (g *Generator) LoadWordsFromFile(filePath, category string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file %s: %w", filePath, err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	count := 0

	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if word == "" {
			continue
		}

		switch category {
		case "words":
			g.wordlist.AddWord(word)
		case "ssids":
			g.wordlist.AddSSID(word)
		default:
			return fmt.Errorf("unknown category: %s", category)
		}
		count++
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	fmt.Printf("Loaded %d words for %s\n", count, category)

	return nil
}

type PasswordJob struct {
	Pattern    string
	CustomWord string
	CommonWord string
	SSID       string
	Year       int
	Number     string
	Separators  []string
}

func (g *Generator) Generate(ctx context.Context, outputFile string) error {
	customWords, err := g.getVariations(g.wordlist.GetWords())
	if err != nil {
		return fmt.Errorf("failed to get custom word variations: %w", err)
	}

	commonWords, err := g.getVariations(g.config.CommonWords)
	if err != nil {
		return fmt.Errorf("failed to get common word variations: %w", err)
	}

	ssids, err := g.getVariations(g.wordlist.GetSSIDs())

	// Generate all number patterns including digit replacements
	allNumbers := g.patterns.GenerateAllNumberPatterns()

	// Create output file
	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer file.Close()

	writer := bufio.NewWriter(file)
	defer writer.Flush()

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
			fmt.Fprintln(writer, password)
			passwordCount++

			if passwordCount%10000 == 0 {
				fmt.Printf("\rGenerated %d unique passwords...", passwordCount)
				writer.Flush() // Flush periodically
			}
		}
	}()

	// Generate jobs based on patterns
	go func() {
		defer close(jobChan)
		g.generateJobs(ctx, jobChan, customWords, commonWords, ssids, allNumbers)
	}()

	// Wait for all workers to finish
	wg.Wait()
	close(resultChan)

	// Wait for writer to finish
	writerWg.Wait()

	fmt.Printf("\n\nGenerated %d total unique passwords\n", passwordCount)

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

func (g *Generator) generateJobs(ctx context.Context, jobChan chan<- PasswordJob, customWords, commonWords, ssids, allNumbers []string) {
	for _, pattern := range g.config.Patterns {
		select {
		case <-ctx.Done():
			return
		default:
		}

		// Ignore patterns with SSID if not entered
		if len(ssids) == 0 && strings.Contains(pattern, "<SSID>") {
			continue
		}

		g.generateJobsForPattern(ctx, jobChan, pattern, customWords, commonWords, ssids, allNumbers)
	}
}

func (g *Generator) generateJobsForPattern(ctx context.Context, jobChan chan<- PasswordJob, pattern string, customWords, commonWords, ssids, allNumbers []string) {
	// Determine what placeholders exist in the pattern
	hasCustom := strings.Contains(pattern, "<CUSTOM>")
	hasCommon := strings.Contains(pattern, "<COMMON>")
	hasSSID := strings.Contains(pattern, "<SSID>")
	hasYear := strings.Contains(pattern, "<YEAR>") || strings.Contains(pattern, "<SHORTYEAR>")
	hasNum := strings.Contains(pattern, "<NUM>")

	separatorCount  := strings.Count(pattern, "<SEP>")

	// Generate combinations based on pattern requirements
	customRange := []string{""}
	if hasCustom {
		customRange = customWords
	}

	commonRange := []string{""}
	if hasCommon {
		commonRange = commonWords
	}

	ssidRange := []string{""}
	if hasSSID {
		ssidRange = ssids
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
		numberRange = allNumbers
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
	wordVariations, err := g.generateVariations(words, g.patterns.GenerateWordVariations)
	if err != nil {
		return nil, err
	}

	caseVariations, err := g.generateVariations(wordVariations, g.patterns.GenerateCaseVariations)
	if err != nil {
		return nil, err
	}

	subVariations, err := g.generateVariations(caseVariations, g.patterns.ApplyAllSubstitutions)
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

	return g.patterns.ConvertSetToSlice(variations), nil
}
