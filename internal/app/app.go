package app

import (
	"context"
	"fmt"

	"github.com/omarelshopky/craftlist/internal/config"
	"github.com/omarelshopky/craftlist/internal/generator"
	"github.com/omarelshopky/craftlist/internal/interfaces"
	"github.com/omarelshopky/craftlist/internal/ui"
	"github.com/omarelshopky/craftlist/internal/wordlist"
	"github.com/omarelshopky/craftlist/pkg/errors"
	"github.com/spf13/cobra"
)

const (
	AppName    = "craftlist"
	AppVersion = "0.2.0"
)

type App struct {
	flags   *Flags
	printer interfaces.Printer
}

func New() *App {
	return &App{
		flags:   NewFlags(),
		printer: ui.NewPrinter(),
	}
}

func Run(ctx context.Context) error {
	app := New()

	return app.Execute(ctx)
}

func (a *App) Execute(ctx context.Context) error {
	rootCmd := &cobra.Command{
		Use:           fmt.Sprintf("%s [-w words.txt|--list-placeholders]", AppName),
		Long:          "A tool for generating customized wordlists tailored to a company's specific details.",
		Version:       AppVersion,
		RunE:          a.runCommand,
		SilenceUsage:  true,
		SilenceErrors: true,
	}

	rootCmd.SetContext(ctx)
	a.setupFlags(rootCmd)
	a.setupErrorHandling(rootCmd)

	if err := rootCmd.Execute(); err != nil {
		if err != errors.SilentErr {
			a.printer.Error(fmt.Sprintf("\nError: %v", err))
		}
		return err
	}

	return nil
}

func (a *App) runCommand(cmd *cobra.Command, args []string) error {
	ctx := cmd.Context()

	a.printer.PrintIntro(AppVersion)

	if a.flags.ListPlaceholders {
		return a.handleListPlaceholders()
	}

	return a.runGeneration(ctx)
}

func (a *App) handleListPlaceholders() error {
	cfg, err := a.loadConfiguration()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	a.printer.PrintPlaceholders(cfg.Placeholders)

	return nil
}

func (a *App) runGeneration(ctx context.Context) error {
	cfg, err := a.buildConfiguration()
	if err != nil {
		return err
	}

	if err := cfg.Validate(); err != nil {
		return err
	}

	gen := generator.New(cfg.Generator, cfg.Placeholders)
	counter := generator.NewCounter(cfg.Generator, cfg.Placeholders)
	loader := wordlist.NewLoader()

	if err := a.loadWordLists(gen, loader); err != nil {
		return fmt.Errorf("failed to load word lists: %w", err)
	}

	if gen.PrepareVariations() != nil {
		return err
	}

	a.printer.PrintApproximateCount(counter.EstimatePasswordCount(gen.GetCustomWordsCount(), gen.GetCommonWordsCount(), gen.GetSSIDsCount(), gen.GetNumbersCount()))

	a.printer.Info("\nGenerating password combinations...")

	if err := gen.Generate(ctx, cfg.Output.Filename, a.printer); err != nil {
		return fmt.Errorf("password generation failed: %w", err)
	}

	a.printer.PrintOutputFile(cfg.Output.Filename)

	return nil
}

func (a *App) loadConfiguration() (*config.Config, error) {
	return config.Load(a.flags.CfgFile)
}

func (a *App) buildConfiguration() (*config.Config, error) {
	cfg, err := a.loadConfiguration()
	if err != nil {
		return nil, err
	}

	a.applyCliOverrides(cfg)

	return cfg, nil
}

func (a *App) applyCliOverrides(cfg *config.Config) {
	cfg.Output.Filename = a.flags.OutputFile
	cfg.Generator.MinPasswordLen = a.flags.MinLength
	cfg.Generator.MaxPasswordLen = a.flags.MaxLength
	cfg.Generator.MinYear = a.flags.MinYear
	cfg.Generator.MaxYear = a.flags.MaxYear
}

func (a *App) loadWordLists(gen *generator.Generator, loader *wordlist.Loader) error {
	if a.flags.WordsFile != "" {
		words, err := loader.LoadFromFile(a.flags.WordsFile)
		if err != nil {
			return fmt.Errorf("failed to load words file '%s': %w", a.flags.WordsFile, err)
		}

		gen.SetCustomWords(words)
		a.printer.PrintLoadedWords("custom words", len(words))
	}

	if a.flags.SSIDsFile != "" {
		ssids, err := loader.LoadFromFile(a.flags.SSIDsFile)
		if err != nil {
			return fmt.Errorf("failed to load SSIDs file '%s': %w", a.flags.SSIDsFile, err)
		}
		gen.SetSSIDs(ssids)
		a.printer.PrintLoadedWords("SSIDs", len(ssids))
	}

	return nil
}

func (a *App) setupFlags(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVarP(&a.flags.CfgFile, "config", "c", "", "config file path")

	cmd.Flags().StringVarP(&a.flags.WordsFile, "words", "w", "", "path to company names and abbreviations file (one per line)")
	cmd.Flags().StringVarP(&a.flags.SSIDsFile, "ssids", "s", "", "path to SSIDs file (one per line)")

	cmd.Flags().StringVarP(&a.flags.OutputFile, "output", "o", "passwords.txt", "output file path")

	cmd.Flags().IntVar(&a.flags.MinLength, "min-length", 8, "minimum password length")
	cmd.Flags().IntVar(&a.flags.MaxLength, "max-length", 64, "maximum password length")

	cmd.Flags().IntVar(&a.flags.MinYear, "min-year", 1990, "minimum year for combinations")
	cmd.Flags().IntVar(&a.flags.MaxYear, "max-year", 2025, "maximum year for combinations")

	cmd.Flags().BoolVar(&a.flags.ListPlaceholders, "list-placeholders", false, "list all available placeholders and exit")

	cmd.MarkFlagsOneRequired("words", "list-placeholders")
}

func (a *App) setupErrorHandling(cmd *cobra.Command) {
	cmd.SetFlagErrorFunc(func(cmd *cobra.Command, err error) error {
		a.printer.Error(err.Error() + "\n")
		cmd.Println(cmd.UsageString())
	
		return errors.SilentErr
	})
}