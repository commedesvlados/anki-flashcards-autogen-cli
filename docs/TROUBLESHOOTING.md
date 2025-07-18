# Troubleshooting

## Common Issues

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

## Debug Mode

Enable verbose logging for detailed debugging:
```bash
anki-builder make-apkg --input data/words.xlsx --output output.vocab.apkg --unsplash YOUR_API_KEY --verbose
``` 

## PDF Extraction (extract-pdf) Issues

- **Developer Note:** The extract-pdf command is orchestrated via `app.NewPDFExtractor` and `PDFExtractorConfig`, just like make-apkg uses `NewApkgMaker`. If you are debugging, check these orchestrators.

1. **UniPDF API key error**:
   - Ensure you have a valid UniPDF API key (https://unidoc.io/pricing/)
   - Pass it with `--uni-api-key`
2. **PDF file not found**:
   - Check the path to your PDF file
   - Ensure file permissions are correct
3. **Excel file write error**:
   - Check output path and permissions
   - Ensure the file is not open in another program
4. **No words extracted**:
   - Ensure your PDF has highlighted or underlined annotations
   - Some PDFs may use non-standard annotation formats 