package main

import (
    "context"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/omarelshopky/craftlist/internal/config"
    "github.com/omarelshopky/craftlist/internal/generator"
    "github.com/omarelshopky/craftlist/internal/output"
)

const (
    appName    = "craftlist"
    appVersion = "1.0.0"
)

func main() {
    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer cancel()

    if err := run(ctx); err != nil {
        log.Fatal(err)
    }
}

func run(ctx context.Context) error {
    cfg, err := config.Load()
    if err != nil {
        return fmt.Errorf("failed to load config: %w", err)
    }

    gen := generator.New(cfg.Generator)
    
    fmt.Printf("%s v%s - Professional Password List Generator\n", appName, appVersion)
    fmt.Println()

    if err := promptForInput(gen); err != nil {
        return fmt.Errorf("failed to get input: %w", err)
    }

    fmt.Println("\nGenerating password combinations...")
    
    passwords, err := gen.Generate(ctx)
    if err != nil {
        return fmt.Errorf("failed to generate passwords: %w", err)
    }

    writer := output.NewWriter(cfg.Output)
    if err := writer.Write(passwords); err != nil {
        return fmt.Errorf("failed to write passwords: %w", err)
    }

    fmt.Printf("Generated %d unique passwords\n", len(passwords))
    fmt.Printf("Saved to: %s\n", cfg.Output.Filename)
    
    if cfg.Output.ShowExamples {
        showExamples(passwords, cfg.Output.ExampleCount)
    }

    return nil
}

func promptForInput(gen *generator.Generator) error {
    return gen.PromptForTargetInfo()
}

func showExamples(passwords []string, count int) {
    if len(passwords) == 0 {
        return
    }

    if count > len(passwords) {
        count = len(passwords)
    }

    fmt.Println("\n📝 Example passwords:")
    for i := 0; i < count; i++ {
        fmt.Printf("   %s\n", passwords[i])
    }

    if len(passwords) > count {
        fmt.Printf("   ... and %d more\n", len(passwords)-count)
    }
}