package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/omarelshopky/craftlist/internal/config"
	"github.com/omarelshopky/craftlist/internal/generator"
	"github.com/spf13/cobra"
)

const (
	appName    = "craftlist"
	appVersion = "0.1.1"
)

var (
	cfgFile    string
	wordsFile  string
	ssidsFile  string
	outputFile string
	minLength  int
	maxLength  int
	minYear    int
	maxYear    int
)

var rootCmd = &cobra.Command{
	Use:     fmt.Sprintf("%s -w words.txt", appName),
	Long:    "A tool for generating customized wordlists tailored to a company's specific details.",
	Version: appVersion,
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx := cmd.Context()

		printIntro()

		return run(ctx)
	},
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	rootCmd.SetContext(ctx)
	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "", "config file")

	rootCmd.Flags().StringVarP(&wordsFile, "words", "w", "", "path to company names, and abbreviations file (one per line)")
	rootCmd.Flags().StringVarP(&ssidsFile, "ssids", "s", "", "path to SSIDs file (one per line)")

	rootCmd.Flags().StringVarP(&outputFile, "output", "o", "passwords.ls", "output file path")

	rootCmd.Flags().IntVar(&minLength, "min-length", 8, "minimum password length")
	rootCmd.Flags().IntVar(&maxLength, "max-length", 64, "maximum password length")

	rootCmd.Flags().IntVar(&minYear, "min-year", 1990, "minimum year for combinations")
	rootCmd.Flags().IntVar(&maxYear, "max-year", time.Now().Year(), "maximum year for combinations")

	rootCmd.MarkFlagRequired("words")
}

func run(ctx context.Context) error {
	cfg, err := getConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Validate

	gen := generator.New(cfg.Generator)

	if err := loadWordlistsFromFiles(gen); err != nil {
		return fmt.Errorf("failed to load wordlists: %w", err)
	}

	fmt.Println("\nGenerating password combinations...")

	if err := gen.Generate(ctx, cfg.Output.Filename); err != nil {
		return fmt.Errorf("failed to generate passwords: %w", err)
	}

	fmt.Printf("Output saved to: %s\n", cfg.Output.Filename)

	return nil
}

func getConfig() (*config.Config, error) {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return nil, err
	}

	cfg.Output.Filename = outputFile
	cfg.Generator.MinPasswordLen = minLength
	cfg.Generator.MaxPasswordLen = maxLength
	cfg.Generator.MinYear = minYear
	cfg.Generator.MaxYear = maxYear

	return cfg, nil
}

func loadWordlistsFromFiles(gen *generator.Generator) error {
	if wordsFile != "" {
		if err := gen.LoadWordsFromFile(wordsFile, "words"); err != nil {
			return fmt.Errorf("failed to load words file: %w", err)
		}
	}

	if ssidsFile != "" {
		if err := gen.LoadWordsFromFile(ssidsFile, "ssids"); err != nil {
			return fmt.Errorf("failed to load SSIDs file: %w", err)
		}
	}

	return nil
}

func printIntro() {
	fmt.Printf(
`                 __ _   _ _     _   
                / _| | | (_)   | |  
  ___ _ __ __ _| |_| |_| |_ ___| |_ 
 / __| '__/ _' |  _| __| | / __| __|
| (__| | | (_| | | | |_| | \__ \ |_ 
 \___|_|  \__,_|_|  \__|_|_|___/\__|
                                      
v%s             By Omar Elshopky`, appVersion)
	fmt.Println()
	fmt.Println()
}
