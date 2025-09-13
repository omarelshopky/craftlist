package config

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestLoad(t *testing.T) {
	t.Run("valid JSON file", func(t *testing.T) {
		tmpFile := filepath.Join(t.TempDir(), "config.json")
		jsonData := `{
			"common_words": ["admin", "test"],
			"separators": ["-", "_"],
			"number_patterns": ["123", "456"],
			"substitutions": {"a":["@"],"e":["3"]},
			"patterns": ["<CUSTOM><SEP><NUM>"]
		}`
		if err := os.WriteFile(tmpFile, []byte(jsonData), 0644); err != nil {
			t.Fatalf("Failed to create temp JSON file: %v", err)
		}

		cfg, err := Load(tmpFile)
		if err != nil {
			t.Fatalf("Load() returned error: %v", err)
		}

		expected := []string{"admin", "test"}
		if !reflect.DeepEqual(cfg.Generator.CommonWords, expected) {
			t.Errorf("expected CommonWords=%v, got %v", expected, cfg.Generator.CommonWords)
		}

		expected = []string{"-", "_"}
		if !reflect.DeepEqual(cfg.Generator.Separators, []string{"-", "_"}) {
			t.Errorf("expected Separators=%v, got %v", expected, cfg.Generator.Separators)
		}

		expected = []string{"123", "456"}
		if !reflect.DeepEqual(cfg.Generator.NumberPatterns, expected) {
			t.Errorf("expected NumberPatterns=%v, got %v", expected, cfg.Generator.NumberPatterns)
		}

		if len(cfg.Generator.Substitutions["a"]) == 0 || cfg.Generator.Substitutions["a"][0] != "@" {
			t.Errorf("expected substitutions not applied correctly")
		}

		expected = []string{"<CUSTOM><SEP><NUM>"}
		if !reflect.DeepEqual(cfg.Generator.Patterns, expected) {
			t.Errorf("expected Patterns=%v, got %v", expected, cfg.Generator.Patterns)
		}
	})

	t.Run("non existent JSON file", func(t *testing.T) {
		_, err := Load("non_existent.json")
		if err == nil {
			t.Error("expected error for non-existent file, got nil")
		}
	})

	t.Run("invalid JSON file", func(t *testing.T) {
		tmpFile := filepath.Join(t.TempDir(), "invalid.json")
		if err := os.WriteFile(tmpFile, []byte(`{invalid json`), 0644); err != nil {
			t.Fatalf("Failed to write invalid JSON file: %v", err)
		}

		_, err := Load(tmpFile)
		if err == nil {
			t.Error("expected error for invalid JSON, got nil")
		}
	})
}