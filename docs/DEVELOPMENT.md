# Development

## Running Tests
```bash
go test ./...
```

## Building for Different Platforms

### Using Makefile (Recommended)
```bash
# Build for current platform
make build

# Build for all platforms (Linux, macOS, Windows)
make release

# Clean build artifacts
make clean
```

### Manual Cross-Compilation
```bash
# Linux (x86_64)
GOOS=linux GOARCH=amd64 go build -o build/anki-builder_linux_amd64 cmd/cli/main.go

# Linux (ARM64)
GOOS=linux GOARCH=arm64 go build -o build/anki-builder_linux_arm64 cmd/cli/main.go

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o build/anki-builder_darwin_amd64 cmd/cli/main.go

# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o build/anki-builder_darwin_arm64 cmd/cli/main.go

# Windows (x86_64)
GOOS=windows GOARCH=amd64 go build -o build/anki-builder_windows_amd64.exe cmd/cli/main.go
```

## Development Commands

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

## Code Style
- Follows Go conventions and idioms
- Use meaningful variable and function names
- Add comments for complex logic
- Keep functions focused and single-purpose
- Use proper error handling with wrapped errors
- Prefer composition over inheritance
- Use `go fmt` and `gofmt -s` for formatting
- Use Go naming conventions (camelCase for variables, PascalCase for exported)
- Use `context.Context` for cancellation and timeouts
- Implement proper logging with structured logging (zap)
- Use interfaces for testability
- Handle errors explicitly, don't ignore them 

## PDF Extraction Command (extract-pdf)

- Code: `cmd/cli/extract_pdf.go`
- Implements the `extract-pdf` CLI command for extracting highlighted/underlined words from PDFs to Excel.
- Orchestrated via `app.NewPDFExtractor` and `PDFExtractorConfig` (mirrors `make-apkg`/`NewApkgMaker` pattern).
- Uses UniPDF API (requires API key).
- To add new flags, edit the `extractPdfOptions` struct and `NewExtractPdfCmd()`.
- To test, run:
  ```bash
  go run cmd/cli/main.go extract-pdf --uni-api-key=... --input-pdf-book-path=... --output-excel-path=...
  ``` 