package output

import (
    "bufio"
    "fmt"
    "os"

    "github.com/omarelshopky/craftlist/internal/config"
)

// Writer handles password output
type Writer struct {
    config config.OutputConfig
}

// NewWriter creates a new Writer instance
func NewWriter(cfg config.OutputConfig) *Writer {
    return &Writer{config: cfg}
}

// Write saves passwords to file
func (w *Writer) Write(passwords []string) error {
    file, err := os.Create(w.config.Filename)
    if err != nil {
        return fmt.Errorf("failed to create output file: %w", err)
    }
    defer file.Close()

    writer := bufio.NewWriter(file)
    defer writer.Flush()

    for _, password := range passwords {
        if _, err := fmt.Fprintln(writer, password); err != nil {
            return fmt.Errorf("failed to write password: %w", err)
        }
    }

    return nil
}