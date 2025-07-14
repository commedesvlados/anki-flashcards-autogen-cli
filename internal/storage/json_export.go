package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/commedesvlados/anki-flashcards-autogen-cli/internal/core"

	"go.uber.org/zap"
)

// JSONExporter handles exporting flashcards to JSON format
type JSONExporter struct {
	logger *zap.Logger
}

// NewJSONExporter creates a new JSON exporter
func NewJSONExporter(logger *zap.Logger) *JSONExporter {
	return &JSONExporter{
		logger: logger,
	}
}

// ExportFlashcards exports enriched flashcards to JSON file
func (e *JSONExporter) ExportFlashcards(flashcards []*core.Flashcard, outputPath string) error {
	e.logger.Info("Exporting flashcards to JSON", zap.String("path", outputPath), zap.Int("count", len(flashcards)))

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0700); err != nil { //nolint:mnd
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Convert to export format
	exportFlashcards := make([]*core.ExportFlash, len(flashcards))
	for i, flashcard := range flashcards {
		exportFlashcards[i] = flashcard.ToExportFlash()
	}

	// Marshal to JSON with pretty formatting
	jsonData, err := json.MarshalIndent(exportFlashcards, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal flashcards to JSON: %w", err)
	}

	// Write to file
	if err := os.WriteFile(outputPath, jsonData, 0600); err != nil { //nolint:mnd
		return fmt.Errorf("failed to write JSON file: %w", err)
	}

	e.logger.Info("Successfully exported flashcards to JSON", zap.String("path", outputPath))
	return nil
}

// LoadFlashcards loads flashcards from JSON file
func (e *JSONExporter) LoadFlashcards(inputPath string) ([]*core.ExportFlash, error) {
	e.logger.Info("Loading flashcards from JSON", zap.String("path", inputPath))

	// Read file
	jsonData, err := os.ReadFile(inputPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read JSON file: %w", err)
	}

	// Unmarshal JSON
	var flashcards []*core.ExportFlash
	if err := json.Unmarshal(jsonData, &flashcards); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	e.logger.Info("Successfully loaded flashcards from JSON", zap.String("path", inputPath), zap.Int("count", len(flashcards)))
	return flashcards, nil
}
