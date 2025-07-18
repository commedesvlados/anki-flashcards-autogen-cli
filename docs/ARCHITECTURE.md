# Features

- 📊 **Excel Integration**: Reads Russian-English word pairs from Excel files
- 📚 **Dictionary Enrichment**: Uses Free Dictionary API to get definitions, pronunciations, and audio
- 🖼️ **Image Enrichment**: Uses Unsplash API to get relevant images for each word
- 🎵 **Audio Download**: Downloads pronunciation audio files
- 📦 **Anki Export**: Generates .apkg files compatible with Anki
- 📈 **Progress Tracking**: Shows progress bars during enrichment
- 📝 **Structured Logging**: Comprehensive logging with Zap
- 🔄 **Retry Logic**: Robust error handling with exponential backoff
- 🧪 **Testable**: Clean architecture with unit test support
- 🆕 **PDF Word Extraction**: Extracts highlighted/underlined words from PDF files to Excel (see `extract-pdf` command, orchestrated via `app.NewPDFExtractor`)

---

# Architecture

```
anki-flashcards/
├── cmd/                   # CLI entrypoint
│   └── cli/
│       └── main.go
├── internal/
│   ├── core/              # Domain logic (pure, testable)
│   │   ├── model.go       # Flashcard structs
│   │   └── enrich.go      # Enrichment orchestrator
│   ├── excel/             # Excel file reader
│   │   └── reader.go
│   ├── downloader/        # Media downloaders
│   │   └── downloader.go
│   ├── storage/           # JSON export
│   │   └── json_export.go
│   ├── util/              # Utilities
│   │   └── retry.go
│   └── cli/               # Application orchestrator
│       └── app.go
├── scripts/
│   └── make_apkg.py       # Python genanki script
├── pkg/
│   └── clients/           # API response models
├── data/                  # Excel input files
├── media/                 # Downloaded audio/images
├── enriched/              # Exported enriched flashcards
└── output/                # Generated Anki packages
```

## Code Structure

The application follows clean architecture principles:

- **Domain Layer** (`internal/core/`): Pure business logic, no dependencies
- **Infrastructure Layer** (`pkg/clients/`, `internal/downloader/`): External API clients and media downloaders
- **Application Layer** (`internal/cli/`): Use cases and orchestration
- **Interface Layer** (`cmd/cli/`): CLI interface

### Folder Descriptions
- `cmd/cli/`: CLI entry point
- `cmd/cli/extract_pdf.go`: Implements the `extract-pdf` command for extracting annotated words from PDFs to Excel. Uses UniPDF API and is orchestrated via `app.NewPDFExtractor` and `PDFExtractorConfig` (mirrors `make-apkg`/`NewApkgMaker`).
- `internal/core/`: Business logic and models
- `internal/excel/`: Excel file reading
- `internal/downloader/`: Media downloaders (audio, images)
- `internal/storage/`: JSON export logic
- `internal/util/`: Utilities (e.g., retry logic)
- `internal/cli/`: Application orchestrator
- `pkg/clients/`: API response models and clients
- `scripts/`: Python scripts (e.g., genanki)
- `data/`: Input Excel files
- `media/`: Downloaded media files
- `enriched/`: Enriched JSON data
- `output/`: Generated Anki packages

---

# API Integration

## Free Dictionary API
- **URL**: https://api.dictionaryapi.dev/
- **Rate Limit**: No limits (free service)
- **Data**: Definitions, pronunciations, audio URLs, examples

## Unsplash API
- **URL**: https://api.unsplash.com/
- **Rate Limit**: 50 requests per hour (free tier)
- **Data**: High-quality images for each English word

---

# Error Handling

The application includes robust error handling:

- **Retry Logic**: Exponential backoff for API calls
- **Graceful Degradation**: Continues processing even if some enrichments fail
- **Comprehensive Logging**: Detailed logs for debugging (Zap)
- **Context Cancellation**: Proper timeout handling 