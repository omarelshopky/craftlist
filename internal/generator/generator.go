package generator

import (
    "bufio"
    "context"
    "fmt"
    "os"
    "strconv"
    "strings"

    "github.com/omarelshopky/craftlist/internal/config"
    "github.com/omarelshopky/craftlist/pkg/wordlist"
)

// Generator handles password generation logic
type Generator struct {
    config    config.GeneratorConfig
    wordlist  *wordlist.Wordlist
    patterns  *PatternGenerator
}

// New creates a new Generator instance
func New(cfg config.GeneratorConfig) *Generator {
    return &Generator{
        config:   cfg,
        wordlist: wordlist.New(),
        patterns: NewPatternGenerator(cfg),
    }
}

// PromptForTargetInfo collects target information from user
func (g *Generator) PromptForTargetInfo() error {
    scanner := bufio.NewScanner(os.Stdin)

    // Company names
    fmt.Print("Enter company names (comma-separated): ")
    if !scanner.Scan() {
        return fmt.Errorf("failed to read company names")
    }
    if input := strings.TrimSpace(scanner.Text()); input != "" {
        companies := strings.Split(input, ",")
        for _, company := range companies {
            g.wordlist.AddCompany(strings.TrimSpace(company))
        }
    }

    // Abbreviations
    fmt.Print("Enter abbreviations (comma-separated): ")
    if !scanner.Scan() {
        return fmt.Errorf("failed to read abbreviations")
    }
    if input := strings.TrimSpace(scanner.Text()); input != "" {
        abbreviations := strings.Split(input, ",")
        for _, abbr := range abbreviations {
            g.wordlist.AddAbbreviation(strings.TrimSpace(abbr))
        }
    }

    // SSIDs
    fmt.Print("Enter SSIDs (comma-separated): ")
    if !scanner.Scan() {
        return fmt.Errorf("failed to read SSIDs")
    }
    if input := strings.TrimSpace(scanner.Text()); input != "" {
        ssids := strings.Split(input, ",")
        for _, ssid := range ssids {
            g.wordlist.AddSSID(strings.TrimSpace(ssid))
        }
    }

    return scanner.Err()
}

// Generate creates the password list
func (g *Generator) Generate(ctx context.Context) ([]string, error) {
    baseWords := g.wordlist.GetAllWords()
    if len(baseWords) == 0 {
        return nil, fmt.Errorf("no base words provided")
    }

    passwordSet := make(map[string]struct{})
    
    for _, baseWord := range baseWords {
        select {
        case <-ctx.Done():
            return nil, ctx.Err()
        default:
        }

        passwords := g.patterns.GenerateVariations(baseWord)
        for _, pwd := range passwords {
            if len(pwd) <= g.config.MaxPasswordLen && len(pwd) > 0 {
                passwordSet[pwd] = struct{}{}
            }
        }
    }

    // Convert set to slice
    result := make([]string, 0, len(passwordSet))
    for pwd := range passwordSet {
        result = append(result, pwd)
    }

    return result, nil
}