// Package main provides the extract-pdf command for the CLI.
package main

import (
	"fmt"
	"os"

	"github.com/commedesvlados/anki-flashcards-autogen-cli/internal/app"
	"github.com/commedesvlados/anki-flashcards-autogen-cli/pkg/logger"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

type extractPdfOptions struct {
	uniPDFAPIKey     string
	inputPDFBookPath string
	outputExcelPath  string
}

// NewExtractPdfCmd returns the extract-pdf cobra command.
func NewExtractPdfCmd() *cobra.Command {
	opts := &extractPdfOptions{}
	cmd := &cobra.Command{
		Use:   "extract-pdf",
		Short: "Extract highlighted/underlined words from PDF to Excel",
		Run: func(cmd *cobra.Command, _ []string) {
			// Get global flags
			progressBar, _ = cmd.Root().PersistentFlags().GetBool("progress")
			verbose, _ = cmd.Root().PersistentFlags().GetBool("verbose")
			runExtractPdf(cmd, opts, progressBar, verbose)
		},
	}
	cmd.Flags().StringVar(&opts.uniPDFAPIKey, "uni-api-key", "", "UniPDF API key (required)")
	cmd.Flags().StringVar(&opts.inputPDFBookPath, "input-pdf-book-path", "", "Input PDF file path (required)")
	cmd.Flags().StringVar(&opts.outputExcelPath, "output-excel-path", "", "Output Excel file path (required)")
	cmd.MarkFlagRequired("uni-api-key")         //nolint:errcheck
	cmd.MarkFlagRequired("input-pdf-book-path") //nolint:errcheck
	cmd.MarkFlagRequired("output-excel-path")   //nolint:errcheck
	return cmd
}

func runExtractPdf(_ *cobra.Command, opts *extractPdfOptions, progressBar, verbose bool) {
	log := logger.SetupLogger(verbose)
	defer func() {
		if err := log.Sync(); err != nil {
			fmt.Fprintf(os.Stderr, "logger.Sync error: %v\n", err)
		}
	}()

	config := &app.PDFExtractorConfig{
		ProgressBar:  progressBar,
		UniPDFAPIKey: opts.uniPDFAPIKey,
		PDFPath:      opts.inputPDFBookPath,
		SheetPath:    opts.outputExcelPath,
	}
	extractor := app.NewPDFExtractor(config, log)
	if err := extractor.Run(); err != nil {
		log.Fatal("PDF extraction failed", zap.Error(err)) //nolint:gocritic
	}
	log.Info("Wrote words to Excel", zap.String("output", opts.outputExcelPath))
}
