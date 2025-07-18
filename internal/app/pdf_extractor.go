package app

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"github.com/unidoc/unipdf/v4/common/license"
	"github.com/unidoc/unipdf/v4/core"
	"github.com/unidoc/unipdf/v4/extractor"
	"github.com/unidoc/unipdf/v4/model"
	"github.com/xuri/excelize/v2"
	"go.uber.org/zap"
)

// PDFExtractorConfig holds CLI flag values for extraction
// (set by cobra command)
type PDFExtractorConfig struct {
	ProgressBar  bool
	UniPDFAPIKey string
	PDFPath      string
	SheetPath    string
}

type PDFExtractor struct {
	config *PDFExtractorConfig
	logger *zap.Logger
}

// NewPDFExtractor creates a new PDFExtractor instance
func NewPDFExtractor(config *PDFExtractorConfig, logger *zap.Logger) *PDFExtractor {
	return &PDFExtractor{
		config: config,
		logger: logger,
	}
}

// Run executes the PDF extraction and Excel writing
func (e *PDFExtractor) Run() error {
	if err := license.SetMeteredKey(e.config.UniPDFAPIKey); err != nil {
		e.logger.Fatal("License error", zap.Error(err))
	}

	reader := OpenPDF(e.config.PDFPath, e.logger)
	words := ExtractAnnotatedWords(reader, e.logger)
	WriteWordsToExcel(words, e.config.SheetPath, e.logger)
	e.logger.Info("Wrote words to Excel", zap.String("output", e.config.SheetPath))
	return nil
}

// OpenPDF opens and parses the PDF, returns a PdfReader
func OpenPDF(filename string, logger *zap.Logger) *model.PdfReader {
	f, err := os.Open(filename)
	if err != nil {
		logger.Fatal("Error opening PDF", zap.Error(err))
	}
	reader, err := model.NewPdfReader(f)
	if err != nil {
		logger.Fatal("Error creating reader", zap.Error(err))
	}
	return reader
}

// ExtractAnnotatedWords extracts all unique words from highlight/underline annotations
func ExtractAnnotatedWords(reader *model.PdfReader, logger *zap.Logger) []string { //nolint:gocyclo
	numPages, err := reader.GetNumPages()
	if err != nil {
		logger.Fatal("Error getting page count", zap.Error(err))
	}
	logger.Info("Total pages", zap.Int("pages", numPages))
	wordsSet := make(map[string]struct{})

	for i := 1; i <= numPages; i++ {
		logger.Info("Scanning page", zap.Int("page", i))
		page, err := reader.GetPage(i)
		if err != nil {
			logger.Fatal("Error loading page", zap.Int("page", i), zap.Error(err))
		}
		annots, err := page.GetAnnotations()
		if err != nil {
			logger.Info("No annotations found on page", zap.Int("page", i), zap.Error(err))
			continue
		}
		logger.Info("Found annotations", zap.Int("page", i), zap.Int("count", len(annots)))
		for idx, a := range annots {
			dictObj := a.GetContainingPdfObject()
			dict, ok := core.GetDict(dictObj)
			if !ok {
				logger.Warn("Annotation missing dictionary", zap.Int("page", i), zap.Int("annotation", idx+1))
				continue
			}
			nameObj, hasSubtype := core.GetName(dict.Get("Subtype"))
			if !hasSubtype {
				logger.Warn("Annotation missing subtype", zap.Int("page", i), zap.Int("annotation", idx+1))
				continue
			}
			annotationType := *nameObj
			content := "EMPTY"
			if a.Contents != nil {
				content = a.Contents.String()
			}
			if annotationType == "Highlight" || annotationType == "Underline" {
				if content != "EMPTY" {
					wordsSet[content] = struct{}{}
				} else {
					rectObj := dict.Get("Rect")
					if rectObj != nil {
						arr, ok := rectObj.(*core.PdfObjectArray)
						if ok && len(arr.Elements()) == 4 {
							x1, _ := core.GetNumberAsFloat(arr.Elements()[0])
							y1, _ := core.GetNumberAsFloat(arr.Elements()[1])
							x2, _ := core.GetNumberAsFloat(arr.Elements()[2])
							y2, _ := core.GetNumberAsFloat(arr.Elements()[3])
							rect := struct{ X1, Y1, X2, Y2 float64 }{x1, y1, x2, y2}
							ext, err := extractor.New(page)
							if err == nil {
								pageText, _, _, err := ext.ExtractPageText()
								if err == nil {
									var extracted string
									marks := pageText.Marks()
									for _, mark := range marks.Elements() { //nolint:gocritic
										mb := mark.BBox
										if mb.Llx < rect.X2 && mb.Urx > rect.X1 && mb.Lly < rect.Y2 && mb.Ury > rect.Y1 {
											extracted += mark.Text
										}
									}
									if len(extracted) > 0 { //nolint:gocritic
										wordsSet[extracted] = struct{}{}
									}
								}
							}
						}
					}
				}
			}
			logger.Info("Annotation", zap.Int("page", i), zap.Int("annotation", idx+1),
				zap.String("type", annotationType.String()), zap.String("word", content))
		}
	}
	words := make([]string, 0, len(wordsSet))
	for w := range wordsSet {
		words = append(words, w)
	}
	sort.Strings(words)
	return words
}

// WriteWordsToExcel writes the words to an .xlsx file in column B with headers
func WriteWordsToExcel(words []string, filename string, logger *zap.Logger) {
	fexcel := excelize.NewFile()
	sheet := "WordsSheet1"
	fexcel.SetSheetName(fexcel.GetSheetName(0), sheet) //nolint:errcheck
	fexcel.SetCellValue(sheet, "A1", "Russian")        //nolint:errcheck
	fexcel.SetCellValue(sheet, "B1", "English")        //nolint:errcheck
	fexcel.SetCellValue(sheet, "C1", "PartOfSpeech")   //nolint:errcheck
	for i, w := range words {
		cell := fmt.Sprintf("B%d", i+2)
		fexcel.SetCellValue(sheet, cell, strings.ToLower(w)) //nolint:errcheck
	}
	err := fexcel.SaveAs(filename)
	if err != nil {
		logger.Fatal("Error saving Excel file", zap.Error(err))
	}
}
