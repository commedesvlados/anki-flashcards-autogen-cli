# Anki Flashcard Builder

A CLI tool in Go that reads Russian-English word pairs from Excel, enriches them with dictionary data and images, and generates Anki packages (.apkg) for language learning.

## Features

- 📊 **Excel Integration**: Reads Russian-English word pairs from Excel files
- 📚 **Dictionary Enrichment**: Uses Free Dictionary API to get definitions, pronunciations, and audio
- 🖼️ **Image Enrichment**: Uses Unsplash API to get relevant images for each word
- 🎵 **Audio Download**: Downloads pronunciation audio files
- 📦 **Anki Export**: Generates .apkg files compatible with Anki
- 📈 **Progress Tracking**: Shows progress bars during enrichment
- 📝 **Structured Logging**: Comprehensive logging with Zap
- 🔄 **Retry Logic**: Robust error handling with exponential backoff
- 🧪 **Testable**: Clean architecture with unit test support

## Architecture

```
anki-flashcards/
├── cmd/                   # CLI entrypoint
│   └── app/
│       └── main.go
├── internal/
│   ├── core/              # Domain logic (pure, testable)
│   │   ├── model.go       # Flashcard structs
│   │   └── enrich.go      # Enrichment orchestrator
│   ├── excel/             # Excel file reader
│   │   └── reader.go
│   ├── api/               # API clients
│   │   ├── dictionary.go  # Free Dictionary API
│   │   └── image.go       # Unsplash API
│   ├── media/             # Media downloaders
│   │   └── downloader.go
│   ├── storage/           # JSON export
│   │   └── json_export.go
│   ├── util/              # Utilities
│   │   └── retry.go
│   └── app/               # Application orchestrator
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

## Prerequisites

### Go Requirements
- Go 1.23.8 or later
- All Go dependencies will be downloaded automatically

### Python Requirements
```bash
pip install genanki
```

### API Keys
- **Unsplash API Key**: Get a free API key from [Unsplash Developers](https://unsplash.com/developers)
- **Free Dictionary API**: No API key required (free service)

## Installation

### Option 1: Download Pre-built Binary (Recommended)

#### Quick Install Script
```bash
curl -fsSL https://raw.githubusercontent.com/commedesvlados/anki-flashcards-autogen-cli/main/scripts/install.sh | bash
```

#### Manual Installation

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

### Option 2: Build from Source

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
   go build -o build/anki-builder cmd/app/main.go
   ```

5. **Install locally**:
   ```bash
   # System-wide installation (requires sudo)
   make install
   
   # User installation (no sudo required)
   make install-user
   ```

## Usage

### Basic Usage

```bash
anki-builder --excel your_sheet.xlsx --output completed_flashcards.apkg --unsplash YOUR_UNSPLASH_API_KEY
```

#### Optional: Clean up media files after build
To save disk space (especially for large decks), use the `--no-media-cache` flag to automatically delete all files in the media directory after the .apkg is built:

```bash
anki-builder --excel your_sheet.xlsx --output completed_flashcards.apkg --unsplash YOUR_UNSPLASH_API_KEY --no-media-cache
```

**Note:** Only `--excel` (`-e`), `--output` (`-o`), and `--verbose` (`-v`) have short flag forms. All other flags must use the double-dash long form (e.g. `--deck`, `--unsplash`).

### Command Line Options

| Flag | Short | Description | Default | Required |
|------|-------|-------------|---------|----------|
| `--excel` | `-e` | Path to Excel file with word pairs | `data/words.xlsx` | No |
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
anki-builder \
  --excel data/vocabulary.xlsx \
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

## API Integration

### Free Dictionary API
- **URL**: https://api.dictionaryapi.dev/
- **Rate Limit**: No limits (free service)
- **Data**: Definitions, pronunciations, audio URLs, examples

### Unsplash API
- **URL**: https://api.unsplash.com/
- **Rate Limit**: 50 requests per hour (free tier)
- **Data**: High-quality images for each English word

## Error Handling

The application includes robust error handling:

- **Retry Logic**: Exponential backoff for API calls
- **Graceful Degradation**: Continues processing even if some enrichments fail
- **Comprehensive Logging**: Detailed logs for debugging
- **Context Cancellation**: Proper timeout handling

## Development

### Running Tests
```bash
go test ./...
```

### Building for Different Platforms

#### Using Makefile (Recommended)
```bash
# Build for current platform
make build

# Build for all platforms (Linux, macOS, Windows)
make release

# Clean build artifacts
make clean
```

#### Manual Cross-Compilation
```bash
# Linux (x86_64)
GOOS=linux GOARCH=amd64 go build -o build/anki-builder_linux_amd64 cmd/app/main.go

# Linux (ARM64)
GOOS=linux GOARCH=arm64 go build -o build/anki-builder_linux_arm64 cmd/app/main.go

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o build/anki-builder_darwin_amd64 cmd/app/main.go

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o build/anki-builder_darwin_arm64 cmd/app/main.go

# Windows (x86_64)
GOOS=windows GOARCH=amd64 go build -o build/anki-builder_windows_amd64.exe cmd/app/main.go
```

### Development Commands

```bash
# Run tests
make test

# Run tests with coverage
make test-coverage

# Format code
make fmt

# Lint code
make lint

# Show all available commands
make help
```

### Code Structure

The application follows clean architecture principles:

- **Domain Layer** (`internal/core/`): Pure business logic, no dependencies
- **Infrastructure Layer** (`internal/api/`, `internal/media/`): External API clients
- **Application Layer** (`internal/app/`): Use cases and orchestration
- **Interface Layer** (`cmd/app/`): CLI interface

## Troubleshooting

### Common Issues

1. **Python genanki not found**:
   ```bash
   pip install genanki
   ```

2. **Unsplash API rate limit exceeded**:
   - Wait for rate limit reset (1 hour)
   - Consider upgrading to paid plan

3. **Excel file not found**:
   - Ensure the Excel file exists and path is correct
   - Check file permissions

4. **Audio/images not downloading**:
   - Check internet connection
   - Verify API keys are valid
   - Check media directory permissions

5. **Binary not found in PATH**:
   ```bash
   # Check if binary is installed
   which anki-builder
   
   # Add to PATH if needed
   export PATH="$PATH:/usr/local/bin"
   # or
   export PATH="$PATH:$HOME/.local/bin"
   ```

6. **Permission denied on binary**:
   ```bash
   chmod +x /usr/local/bin/anki-builder
   # or
   chmod +x ~/.local/bin/anki-builder
   ```

7. **Cross-platform build issues**:
   ```bash
   # Ensure Go is properly installed
   go version
   
   # Clean and rebuild
   make clean
   make release
   ```

### Debug Mode

Enable verbose logging for detailed debugging:
```bash
./build/anki-builder --verbose --unsplash YOUR_API_KEY
```

## Distribution & Releases

### Creating a Release

#### Automated Release (Recommended)

1. **Create and push a git tag**:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. **GitHub Actions will automatically**:
   - Build binaries for all platforms
   - Create a GitHub release
   - Upload all artifacts

#### Manual Release

1. **Build all platforms**:
   ```bash
   make release
   ```

2. **Create release package**:
   ```bash
   ./scripts/release.sh
   ```

3. **Upload to GitHub Releases**:
   - Go to GitHub repository
   - Create a new release with tag
   - Upload files from `release_v1.0.0/` directory

### Supported Platforms

| Platform | Architecture | Binary Name |
|----------|--------------|-------------|
| Linux | x86_64 | `anki-builder_linux_amd64` |
| Linux | ARM64 | `anki-builder_linux_arm64` |
| macOS | Intel | `anki-builder_darwin_amd64` |
| macOS | Apple Silicon | `anki-builder_darwin_arm64` |
| Windows | x86_64 | `anki-builder_windows_amd64.exe` |

### Release Process

1. **Version Management**:
   - Version is automatically extracted from git tags
   - Build time is embedded in binary
   - Use semantic versioning (v1.0.0, v1.1.0, etc.)

2. **Build Process**:
   - Cross-compilation for all target platforms
   - Stripped binaries for smaller size
   - Checksums generated for verification

3. **Distribution**:
   - GitHub Releases for binary distribution
   - Installation script for easy setup
   - Documentation updated for each release

### Package Managers

#### Manual Installation Script
```bash
# Download and run installer
curl -fsSL https://raw.githubusercontent.com/commedesvlados/anki-flashcards-autogen-cli/main/scripts/install.sh | bash
```

### CI/CD Pipeline

The project uses GitHub Actions for automated builds:

- **Trigger**: Push to tags starting with `v*`
- **Builds**: All supported platforms
- **Artifacts**: Binary files, checksums, release notes
- **Deployment**: Automatic GitHub release creation

### Development Workflow

1. **Local Development**:
   ```bash
   make build    # Build for current platform
   make test     # Run tests
   make fmt      # Format code
   ```

2. **Pre-release Testing**:
   ```bash
   make release  # Build all platforms
   make test-coverage  # Check test coverage
   ```

3. **Release**:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

### Development Setup

```bash
# Clone repository
git clone https://github.com/commedesvlados/anki-flashcards-autogen-cli.git
cd anki-flashcards-autogen-cli

# Install dependencies
go mod tidy
pip install genanki

# Build and install locally
make build
make install-user

# Run tests
make test
```

## Acknowledgments

- [Free Dictionary API](https://dictionaryapi.dev/) for word definitions and audio
- [Unsplash](https://unsplash.com/) for high-quality images
- [genanki](https://github.com/kerrickstaley/genanki) for Anki package generation
- [Zap](https://github.com/uber-go/zap) for structured logging 