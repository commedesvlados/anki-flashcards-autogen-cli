package excel

import (
	"fmt"

	"github.com/commedesvlados/anki-flashcards-autogen-cli/internal/core"

	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

// Reader handles reading Excel files with word pairs
type Reader struct {
	logger *zap.Logger
}

// NewReader creates a new Excel reader
func NewReader(logger *zap.Logger) *Reader {
	return &Reader{
		logger: logger,
	}
}

// ReadWordPairs reads Russian-English word pairs from an Excel file
func (r *Reader) ReadWordPairs(filePath string) ([]*core.RawFlashcard, error) {
	r.logger.Info("Reading Excel file", zap.String("path", filePath))

	// Open Excel file
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer f.Close()

	// Get the first sheet
	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, fmt.Errorf("no sheets found in Excel file")
	}

	sheetName := sheets[0]
	r.logger.Debug("Reading sheet", zap.String("sheet", sheetName))

	// Get all rows from the sheet
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return nil, fmt.Errorf("failed to get rows: %w", err)
	}

	if len(rows) < 2 {
		return nil, fmt.Errorf("Excel file must have at least 2 rows (header + data)") //nolint:stylecheck
	}

	var wordPairs []*core.RawFlashcard

	// Skip header row, process data rows
	for i, row := range rows[1:] {
		if len(row) < 2 {
			r.logger.Warn("Skipping row with insufficient columns", zap.Int("row", i+2), zap.Int("columns", len(row)))
			continue
		}

		// Read columns: Russian, English, PartOfSpeech (optional)
		russian := row[0]
		english := row[1]
		partOfSpeech := ""
		if len(row) > 2 {
			partOfSpeech = row[2]
		}

		// Skip empty rows
		if russian == "" || english == "" {
			r.logger.Debug("Skipping empty row", zap.Int("row", i+2))
			continue
		}

		wordPair := &core.RawFlashcard{
			Russian:      russian,
			English:      english,
			PartOfSpeech: partOfSpeech,
		}

		wordPairs = append(wordPairs, wordPair)
		r.logger.Debug("Read word pair", zap.String("russian", russian), zap.String("english", english))
	}

	r.logger.Info("Successfully read word pairs", zap.Int("count", len(wordPairs)))
	return wordPairs, nil
}

// ValidateExcelFile validates that an Excel file has the correct format
func (r *Reader) ValidateExcelFile(filePath string) error {
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to open Excel file: %w", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return fmt.Errorf("no sheets found in Excel file")
	}

	sheetName := sheets[0]
	rows, err := f.GetRows(sheetName)
	if err != nil {
		return fmt.Errorf("failed to get rows: %w", err)
	}

	if len(rows) < 2 {
		return fmt.Errorf("Excel file must have at least 2 rows (header + data)") //nolint:stylecheck
	}

	// Check header row
	header := rows[0]
	if len(header) < 2 {
		return fmt.Errorf("header row must have at least 2 columns")
	}

	// Validate that we have some data rows
	dataRows := 0
	for _, row := range rows[1:] {
		if len(row) >= 2 && row[0] != "" && row[1] != "" {
			dataRows++
		}
	}

	if dataRows == 0 {
		return fmt.Errorf("no valid data rows found in Excel file")
	}

	r.logger.Info("Excel file validation passed",
		zap.String("file", filePath),
		zap.String("sheet", sheetName),
		zap.Int("total_rows", len(rows)),
		zap.Int("data_rows", dataRows))

	return nil
}
