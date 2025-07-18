package app

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/commedesvlados/anki-flashcards-autogen-cli/internal/core"
	"github.com/commedesvlados/anki-flashcards-autogen-cli/internal/downloader"
	"github.com/commedesvlados/anki-flashcards-autogen-cli/internal/excel"
	"github.com/commedesvlados/anki-flashcards-autogen-cli/internal/storage"
	free_dictionary "github.com/commedesvlados/anki-flashcards-autogen-cli/pkg/clients/free-dictionary"
	"github.com/commedesvlados/anki-flashcards-autogen-cli/pkg/clients/unsplash"

	"github.com/schollz/progressbar/v3"
	"go.uber.org/zap"
)

// ApkgMakerConfig holds application configuration
type ApkgMakerConfig struct {
	ProgressBar  bool
	ExcelFile    string
	OutputFile   string
	DeckName     string
	UnsplashKey  string
	MediaDir     string
	EnrichedDir  string
	NoMediaCache bool
}

// ApkgMaker is the main application orchestrator
type ApkgMaker struct {
	config            *ApkgMakerConfig
	logger            *zap.Logger
	excelReader       *excel.Reader
	enrichmentService *core.EnrichmentService
	jsonExporter      *storage.JSONExporter
}

// NewApkgMaker creates a new application instance
func NewApkgMaker(config *ApkgMakerConfig, logger *zap.Logger) *ApkgMaker {
	// Initialize components
	excelReader := excel.NewReader(logger)
	dictionaryAPI := free_dictionary.NewAPI(logger)
	imageAPI := unsplash.NewAPI(config.UnsplashKey, logger)
	downloader := downloader.NewDownloader(config.MediaDir, logger)
	enrichmentService := core.NewEnrichmentService(dictionaryAPI, imageAPI, downloader, logger)
	jsonExporter := storage.NewJSONExporter(logger)

	return &ApkgMaker{
		config:            config,
		logger:            logger,
		excelReader:       excelReader,
		enrichmentService: enrichmentService,
		jsonExporter:      jsonExporter,
	}
}

// Run executes the complete flashcard generation process
func (a *ApkgMaker) Run(ctx context.Context) error {
	a.logger.Info("Starting Anki Flashcard Builder",
		zap.String("excel_file", a.config.ExcelFile),
		zap.String("output_file", a.config.OutputFile))

	// Step 1: Validate Excel file
	a.logger.Info("Step 1: Validating Excel file")
	if err := a.excelReader.ValidateExcelFile(a.config.ExcelFile); err != nil {
		return fmt.Errorf("Excel validation failed: %w", err) //nolint:stylecheck
	}

	// Step 2: Read word pairs from Excel
	a.logger.Info("Step 2: Reading word pairs from Excel")
	rawFlashcards, err := a.excelReader.ReadWordPairs(a.config.ExcelFile)
	if err != nil {
		return fmt.Errorf("failed to read Excel file: %w", err)
	}

	if len(rawFlashcards) == 0 {
		return fmt.Errorf("no word pairs found in Excel file")
	}

	// Step 3: Enrich flashcards with progress bar
	a.logger.Info("Step 3: Enriching flashcards")
	var enrichedFlashcards []*core.Flashcard

	if a.config.ProgressBar {
		enrichedFlashcards, err = a.enrichWithProgress(ctx, rawFlashcards)
	} else {
		enrichedFlashcards, err = a.enrichmentService.EnrichFlashcards(ctx, rawFlashcards)
	}

	if err != nil {
		return fmt.Errorf("failed to enrich flashcards: %w", err)
	}

	// Step 4: Export to JSON
	a.logger.Info("Step 4: Exporting to JSON")
	jsonPath := filepath.Join(a.config.EnrichedDir, "enriched.json")
	if err := a.jsonExporter.ExportFlashcards(enrichedFlashcards, jsonPath); err != nil {
		return fmt.Errorf("failed to export JSON: %w", err)
	}

	// Step 5: Generate Anki package
	a.logger.Info("Step 5: Generating Anki package")
	if err := a.generateAnkiPackage(jsonPath, a.config.OutputFile); err != nil {
		return fmt.Errorf("failed to generate Anki package: %w", err)
	}

	// Step 6: Optionally clean up media directory
	if a.config.NoMediaCache {
		a.logger.Info("Cleaning up media directory after .apkg build", zap.String("media_dir", a.config.MediaDir))
		err := cleanupMediaDir(a.config.MediaDir, a.logger)
		if err != nil {
			a.logger.Warn("Failed to clean up media directory", zap.Error(err))
		}
	}

	a.logger.Info("Successfully completed flashcard generation",
		zap.String("output_file", a.config.OutputFile),
		zap.Int("flashcards", len(enrichedFlashcards)))

	return nil
}

// enrichWithProgress enriches flashcards with a progress bar
func (a *ApkgMaker) enrichWithProgress(ctx context.Context, rawFlashcards []*core.RawFlashcard) ([]*core.Flashcard, error) {
	enriched := make([]*core.Flashcard, 0, len(rawFlashcards))

	bar := progressbar.NewOptions(len(rawFlashcards),
		progressbar.OptionEnableColorCodes(true),
		progressbar.OptionShowBytes(false),
		progressbar.OptionSetWidth(15), //nolint:mnd
		progressbar.OptionSetDescription("[cyan][1/1][reset] Enriching flashcards"),
		progressbar.OptionSetTheme(progressbar.Theme{
			Saucer:        "[green]=[reset]",
			SaucerHead:    "[green]>[reset]",
			SaucerPadding: " ",
			BarStart:      "[",
			BarEnd:        "]",
		}),
	)

	for i, raw := range rawFlashcards {
		select {
		case <-ctx.Done():
			return enriched, ctx.Err()
		default:
		}

		flashcard, err := a.enrichmentService.EnrichFlashcard(ctx, raw, i+1)
		if err != nil {
			a.logger.Error("Failed to enrich flashcard", zap.String("english", raw.English), zap.Error(err))
			// Create basic flashcard without enrichment
			flashcard = &core.Flashcard{
				ID:        i + 1,
				Russian:   raw.Russian,
				English:   raw.English,
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}
		}

		enriched = append(enriched, flashcard)
		if err := bar.Add(1); err != nil {
			a.logger.Warn("Failed to update progress bar", zap.Error(err))
		}
	}

	return enriched, nil
}

// generateAnkiPackage calls the Python script to generate the Anki package
func (a *ApkgMaker) generateAnkiPackage(jsonPath, outputFile string) error {
	scriptPath := "scripts/make_apkg.py"

	// Check if Python script exists
	if _, err := os.Stat(scriptPath); os.IsNotExist(err) {
		return fmt.Errorf("Python script not found: %s", scriptPath) //nolint:stylecheck
	}

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(outputFile)
	if err := os.MkdirAll(outputDir, 0755); err != nil { //nolint:mnd
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	//nolint:gosec // Acceptable risk: controlled input for exec.Command
	cmd := exec.Command("./venv/bin/python", scriptPath, jsonPath, a.config.MediaDir, outputFile, a.config.DeckName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	a.logger.Info("Executing Python script",
		zap.String("script", scriptPath),
		zap.String("json", jsonPath),
		zap.String("media", a.config.MediaDir),
		zap.String("output", outputFile),
		zap.String("deck_name", a.config.DeckName))

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("Python script failed: %w", err) //nolint:stylecheck
	}

	return nil
}

// cleanupMediaDir deletes all files in the given directory
func cleanupMediaDir(mediaDir string, logger *zap.Logger) error {
	d, err := os.Open(mediaDir)
	if err != nil {
		return err
	}
	defer d.Close()

	files, err := d.Readdirnames(-1)
	if err != nil {
		return err
	}

	for _, name := range files {
		path := filepath.Join(mediaDir, name)
		err := os.RemoveAll(path)
		if err != nil {
			logger.Warn("Failed to delete file in media dir", zap.String("file", path), zap.Error(err))
		}
	}
	return nil
}
