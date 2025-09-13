package wordlist

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoadFromFile(t *testing.T) {
	loader := NewLoader()

	t.Run("successfully loads words", func(t *testing.T) {
		tmpFile := filepath.Join(t.TempDir(), "words.ls")
		content := "evil corp\nECO \n EC\n"
		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}

		words, err := loader.LoadFromFile(tmpFile)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := []string{"evil corp", "ECO", "EC"}
		if !reflect.DeepEqual(words, expected) {
			t.Errorf("expected %v, got %v", expected, words)
		}
	})

	t.Run("ignores empty lines", func(t *testing.T) {
		tmpFile := filepath.Join(t.TempDir(), "words_with_empty.ls")
		content := "evil corp\n\nECO\n \nEC\n"
		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}

		words, err := loader.LoadFromFile(tmpFile)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		expected := []string{"evil corp", "ECO", "EC"}
		if !reflect.DeepEqual(words, expected) {
			t.Errorf("expected %v, got %v", expected, words)
		}
	})

	t.Run("file not found", func(t *testing.T) {
		_, err := loader.LoadFromFile("non_existent.ls")
		if err == nil {
			t.Errorf("expected error for non-existent file, got nil")
		}
	})

	t.Run("file with only empty lines returns empty slice", func(t *testing.T) {
		tmpFile := filepath.Join(t.TempDir(), "empty_lines.ls")
		content := "\n \n\n"
		if err := os.WriteFile(tmpFile, []byte(content), 0644); err != nil {
			t.Fatalf("failed to create temp file: %v", err)
		}

		words, err := loader.LoadFromFile(tmpFile)
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}

		if len(words) != 0 {
			t.Errorf("expected empty slice, got %v", words)
		}
	})
}
