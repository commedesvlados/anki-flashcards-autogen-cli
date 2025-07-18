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
./build/anki-builder --verbose --unsplash YOUR_API_KEY
``` 