# Anki Flashcard Builder
 
> **Readme history**
> - [# Installation](#installation)
> - [# Usage](#usage)

> **See more topics:**
> - [ARCHITECTURE.md](./docs/ARCHITECTURE.md) — Features, Architecture, API, Error Handling
> - [DEVELOPMENT.md](./docs/DEVELOPMENT.md) — Development, Testing, Code Style
> - [RELEASE.md](./docs/RELEASE.md) — Release & Distribution
> - [CONTRIBUTING.md](./docs/CONTRIBUTING.md) — Contributing 
> - [TROUBLESHOOTING.md](./docs/TROUBLESHOOTING.md) — Troubleshooting & Debugging

A CLI tool in Go that reads Russian-English word pairs from Excel, enriches them with dictionary data and images, and generates Anki packages (.apkg) for language learning.


# Acknowledgments

- [Free Dictionary API](https://dictionaryapi.dev/) for word definitions and audio
- [Unsplash](https://unsplash.com/) for high-quality images
- [genanki](https://github.com/kerrickstaley/genanki) for Anki package generation
- [Zap](https://github.com/uber-go/zap) for structured logging


# Installation

## Prerequisites

### Go Requirements
- Go 1.22.x or later
- All Go dependencies will be downloaded automatically

### Python Requirements
```bash
pip install genanki
```

### API Keys
- **Unsplash API Key**: Get a free API key from [Unsplash Developers](https://unsplash.com/developers)
- **Free Dictionary API**: No API key required (free service)

## Option 1: Download Pre-built Binary (Recommended)

### Quick Install Script
```bash
curl -fsSL https://raw.githubusercontent.com/commedesvlados/anki-flashcards-autogen-cli/main/scripts/install.sh | bash
```

### Manual Installation

1. **Download the appropriate binary** for your platform from the [latest release](https://github.com/commedesvlados/anki-flashcards-autogen-cli/releases/latest):

   **Linux (x86_64):**
   ```bash
   wget https://github.com/commedesvlados/anki-flashcards-autogen-cli/releases/latest/download/anki-builder_linux_amd64.tar.gz
   tar -xzf anki-builder_linux_amd64.tar.gz
   sudo mv anki-builder /usr/local/bin/
   ```

   **Linux (ARM64):**
   ```bash
   wget https://github.com/commedesvlados/anki-flashcards-autogen-cli/releases/latest/download/anki-builder_linux_arm64.tar.gz
   tar -xzf anki-builder_linux_arm64.tar.gz
   sudo mv anki-builder /usr/local/bin/
   ```

   **macOS (Intel):**
   ```bash
   curl -L https://github.com/commedesvlados/anki-flashcards-autogen-cli/releases/latest/download/anki-builder_darwin_amd64.tar.gz | tar xz
   sudo mv anki-builder /usr/local/bin/
   ```

   **macOS (Apple Silicon):**
   ```bash
   curl -L https://github.com/commedesvlados/anki-flashcards-autogen-cli/releases/latest/download/anki-builder_darwin_arm64.tar.gz | tar xz
   sudo mv anki-builder /usr/local/bin/
   ```

   **Windows:**
   - Download `anki-builder_windows_amd64.zip`
   - Extract and add the directory to your PATH

2. **Install Python dependencies**:
   ```bash
   pip install genanki
   ```

3. **Verify installation**:
   ```bash
   anki-builder --version
   ```

## Option 2: Build from Source

1. **Clone the repository**:
   ```bash
   git clone https://github.com/commedesvlados/anki-flashcards-autogen-cli.git
   cd anki-flashcards-autogen-cli
   ```

2. **Install Go dependencies**:
   ```bash
   go mod tidy
   ```

3. **Install Python dependencies**:
   ```bash
   pip install genanki
   ```

4. **Build the application**:
   ```bash
   # Build for current platform
   make build
   
   # Or build manually
   go build -o build/anki-builder cmd/cli/main.go
   ```

5. **Install locally**:
   ```bash
   # System-wide installation (requires sudo)
   make install
   
   # User installation (no sudo required)
   make install-user
   ```


# Usage

## Basic Usage

```bash
anki-builder make-apkg --input your_sheet.xlsx --output completed_flashcards.apkg --unsplash YOUR_UNSPLASH_API_KEY
```

#### Optional: Clean up media files after build
To save disk space (especially for large decks), use the `--no-media-cache` flag to automatically delete all files in the media directory after the .apkg is built:

```bash
anki-builder make-apkg --input your_sheet.xlsx --output completed_flashcards.apkg --unsplash YOUR_UNSPLASH_API_KEY --no-media-cache
```

**Note:** Only `--input` (`-i`), `--output` (`-o`), and `--verbose` (`-v`) have short flag forms. All other flags must use the double-dash long form (e.g. `--deck`, `--unsplash`).

### Command Line Options

| Flag | Short | Description | Default | Required |
|------|-------|-------------|---------|----------|
| `--input` | `-i` | Path to Excel file with word pairs | `data/words.xlsx` | No |
| `--output` | `-o` | Output Anki package file | `output/vocab.apkg` | No |
| `--media` |  | Directory for downloaded media files | `media` | No |
| `--enriched` |  | Directory for enriched JSON data | `enriched` | No |
| `--unsplash` |  | Unsplash API access key | - | **Yes** |
| `--progress` |  | Show progress bar during enrichment | `true` | No |
| `--verbose` | `-v` | Enable verbose logging | `false` | No |
| `--deck` |  | Name of the Anki deck | `Designed Autogenerated RU-EN Vocabulary` | No |
| `--no-media-cache` |  | Delete all files in media/ after .apkg is built | `false` | No |
| `--help` | `-h` | Show help message | - | No |

### Excel File Format 

The Excel file should have the following structure:

| Column A (Russian) | Column B (English) | Column С (PartOfSpeech) |
|--------------------|--------------------|-------------------------|
|   яблоко           | apple              | noun                    |
|   книга            | book               | noun                    |
|   они              | they               | pronoun                 |
|   пожалуйста       | please             | adverb                  |
|   ехать            | drive              | verb                    |
|   красивый         | pretty             | adjective               |
|   с                | with               | preposition             |
|   в то время как   | while              | conjunction             |
|   вау              | wow                | interjection            |
|   иметь ввиду      | keep in mind       | phrase                  |

**Requirements**:
- First row should be headers (e.g., "Russian", "English", "PartOfSpeech")
- At least 3 columns: Russian words in column A, English words in column B, PartOfSpeech type in column C
- No empty rows between data

### Example Workflow

1. **Prepare Excel file** with Russian-English word pairs
2. **Get Unsplash API key** from [Unsplash Developers](https://unsplash.com/developers)
3. **Run the application**:
```bash
anki-builder make-apkg \
  --input data/vocabulary.xlsx \
  --output my_deck.apkg \
  --unsplash YOUR_API_KEY \
  --verbose \
  --no-media-cache   # Optional: clean up media files after build
```

4. **Import the .apkg file** into Anki

## Output Structure

The application creates the following directory structure:

```
project/
├── data/
│   └── words.xlsx          # Input Excel file
├── media/
│   ├── 1.jpg              # Downloaded images (deleted if --no-media-cache is used)
│   ├── 1_uk.mp3           # Downloaded audio files (deleted if --no-media-cache is used)
│   └── ...
├── enriched/
│   └── enriched.json      # Enriched flashcards in JSON format
└── output/
    └── vocab.apkg         # Generated Anki package
```

**Note:** For large decks (e.g., 5,000+ words), media files can consume significant disk space (hundreds of MBs to several GBs). Use `--no-media-cache` to automatically clean up after building.

## Generated Flashcard Structure

Each flashcard includes:

- **Russian word** (front of card)
- **English word** (back of card)
- **Part of speech** (noun, verb, etc.)
- **Definition** (English definition)
- **Example sentence** (if available)
- **IPA pronunciation** (UK and US)
- **Audio pronunciation** (UK and US)
- **Relevant image**

## PDF Word Extraction (extract-pdf)

Extract all unique words that are highlighted or underlined in a PDF file and write them to an Excel (.xlsx) file. Useful for quickly building vocabulary lists from annotated e-books or study materials.

### Usage

```bash
anki-builder extract-pdf \
  --uni-api-key YOUR_UNIPDF_API_KEY \
  --input-pdf-book-path /path/to/book.pdf \
  --output-excel-path /path/to/output.xlsx
```

### Command Line Options

| Flag | Description | Required |
|------|-------------|----------|
| `--uni-api-key` | UniPDF API key | **Yes** |
| `--input-pdf-book-path` | Path to input PDF file | **Yes** |
| `--output-excel-path` | Path to output Excel file | **Yes** |
| `--verbose` | Enable verbose logging | No |

### Excel Output Format

The generated Excel file will have:

| Column A (Russian) | Column B (English) | Column C (PartOfSpeech) |
|--------------------|--------------------|------------------------|
| (empty)            | extracted_word     | (empty)                |

- Only unique words are written (deduplicated)
- All words are lowercased
- Only column B is filled; columns A and C are left empty for later enrichment

### Example

Suppose you highlight or underline words in a PDF e-book. Run:

```bash
anki-builder extract-pdf --uni-api-key=YOUR_KEY --input-pdf-book-path=book.pdf --output-excel-path=words.xlsx
```

This will create an Excel file with all unique highlighted/underlined words in column B.

**Note:** Requires a [UniPDF API key](https://unidoc.io/pricing/). Free tier available for limited use.

> **Developer Note:**
> The `extract-pdf` command is orchestrated via `app.NewPDFExtractor` and `PDFExtractorConfig`, mirroring the architecture of `make-apkg` (which uses `NewApkgMaker`). This ensures modularity and testability.
