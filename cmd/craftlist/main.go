package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"reflect"
	"strings"
	"syscall"
	"time"
	"errors"

	"github.com/omarelshopky/craftlist/internal/config"
	"github.com/omarelshopky/craftlist/internal/generator"
	"github.com/spf13/cobra"
)

const (
	appName    = "craftlist"
	appVersion = "0.2.0"
)

type flags struct {
	cfgFile          string
	wordsFile        string
	ssidsFile        string
	outputFile       string
	minLength        int
	maxLength        int
	minYear          int
	maxYear          int
	listPlaceholders bool
}

var cliFlags flags

var SilentErr = errors.New("SilentErr")

var rootCmd = &cobra.Command{
	Use:     fmt.Sprintf("%s -w words.txt", appName),
	Long:    "A tool for generating customized wordlists tailored to a company's specific details.",
	Version: appVersion,
	RunE: runCommand,
	SilenceUsage: true,
	SilenceErrors: true,
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	rootCmd.SetContext(ctx)

	if err := rootCmd.Execute(); err != nil {
		if err != SilentErr {
			fmt.Fprint(os.Stderr)
			fmt.Printf("\n%s%v%s\n", config.Colors.Red, err, config.Colors.Reset)
		}
		os.Exit(1)
	}
}

func init() {
	rootCmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		cmd.Printf("%s%v%s\n\n", config.Colors.Red, err, config.Colors.Reset)
		cmd.Println(cmd.UsageString())
	
		return SilentErr
	})

	setupFlags()
}

func setupFlags() {
	rootCmd.PersistentFlags().StringVarP(&cliFlags.cfgFile, "config", "c", "", "config file path")

	rootCmd.Flags().StringVarP(&cliFlags.wordsFile, "words", "w", "", "path to company names and abbreviations file (one per line)")
	rootCmd.Flags().StringVarP(&cliFlags.ssidsFile, "ssids", "s", "", "path to SSIDs file (one per line)")

	rootCmd.Flags().StringVarP(&cliFlags.outputFile, "output", "o", "passwords.txt", "output file path")

	rootCmd.Flags().IntVar(&cliFlags.minLength, "min-length", 8, "minimum password length")
	rootCmd.Flags().IntVar(&cliFlags.maxLength, "max-length", 64, "maximum password length")

	rootCmd.Flags().IntVar(&cliFlags.minYear, "min-year", 1990, "minimum year for combinations")
	rootCmd.Flags().IntVar(&cliFlags.maxYear, "max-year", time.Now().Year(), "maximum year for combinations")

	rootCmd.Flags().BoolVar(&cliFlags.listPlaceholders, "list-placeholders", false, "list all available placeholders and exit")

	rootCmd.MarkFlagsOneRequired("words", "list-placeholders")
}

func runCommand(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	printIntro()

	if cliFlags.listPlaceholders {
		return handleListPlaceholders()
	}

	return runGeneration(ctx)
}

func handleListPlaceholders() error {
	cfg, err := loadConfiguration()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	printPlaceholders(cfg.Placeholders)

	return nil
}

func runGeneration(ctx context.Context) error {
	cfg, err := buildConfiguration()
	if err != nil {
		return err
	}

	if err := cfg.Validate(); err != nil {
		return err
	}

	gen := generator.New(cfg.Generator)

	if err := loadWordLists(gen); err != nil {
		return fmt.Errorf("failed to load word lists: %w", err)
	}

	fmt.Printf("\n%sGenerating password combinations...%s\n", config.Colors.Cyan, config.Colors.Reset)
	if err := gen.Generate(ctx, cfg.Output.Filename); err != nil {
		return fmt.Errorf("password generation failed: %w", err)
	}

	fmt.Printf("Output saved to: %s%s%s\n", config.Colors.Bold, cfg.Output.Filename, config.Colors.Reset)

	return nil
}

func loadConfiguration() (*config.Config, error) {
	return config.Load(cliFlags.cfgFile)
}

func buildConfiguration() (*config.Config, error) {
	cfg, err := loadConfiguration()
	if err != nil {
		return nil, err
	}

	applyCliOverrides(cfg)

	return cfg, nil
}

func applyCliOverrides(cfg *config.Config) {
	cfg.Output.Filename = cliFlags.outputFile
	cfg.Generator.MinPasswordLen = cliFlags.minLength
	cfg.Generator.MaxPasswordLen = cliFlags.maxLength
	cfg.Generator.MinYear = cliFlags.minYear
	cfg.Generator.MaxYear = cliFlags.maxYear
}

func loadWordLists(gen *generator.Generator) error {
	if cliFlags.wordsFile != "" {
		if err := gen.LoadWordsFromFile(cliFlags.wordsFile, "words"); err != nil {
			return fmt.Errorf("failed to load words file '%s': %w", cliFlags.wordsFile, err)
		}
	}

	if cliFlags.ssidsFile != "" {
		if err := gen.LoadWordsFromFile(cliFlags.ssidsFile, "ssids"); err != nil {
			return fmt.Errorf("failed to load SSIDs file '%s': %w", cliFlags.ssidsFile, err)
		}
	}

	return nil
}

func printPlaceholders(placeholders config.PlaceholdersConfig) error {
	fmt.Printf("%sAvailable Placeholders:%s\n\n", config.Colors.Bold, config.Colors.Reset)
	fmt.Printf("%s%-15s %s%s\n", config.Colors.Green, "PLACEHOLDER", "DESCRIPTION", config.Colors.Reset)
	fmt.Printf("%s%-15s %s%s\n", config.Colors.Green, strings.Repeat("-", 15), strings.Repeat("-", 50), config.Colors.Reset)

	values := reflect.ValueOf(placeholders)
    
    for idx := 0; idx < values.NumField(); idx++ {        
		if placeholder, ok := values.Field(idx).Interface().(config.Placeholder); ok {
			fmt.Printf("%s%-15s %s%s\n", config.Colors.Yellow, placeholder.Format, config.Colors.Reset, placeholder.Description)
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
