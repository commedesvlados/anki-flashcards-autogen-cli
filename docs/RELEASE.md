# Release & Distribution

## Creating a Release

### Automated Release (Recommended)

1. **Create and push a git tag**:
   ```bash
   git tag v1.0.0
   git push origin v1.0.0
   ```

2. **GitHub Actions will automatically**:
   - Build binaries for all platforms
   - Create a GitHub release
   - Upload all artifacts

### Manual Release

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

## Supported Platforms

| Platform | Architecture | Binary Name |
|----------|--------------|-------------|
| Linux | x86_64 | `anki-builder_linux_amd64` |
| Linux | ARM64 | `anki-builder_linux_arm64` |
| macOS | Intel | `anki-builder_darwin_amd64` |
| macOS | Apple Silicon | `anki-builder_darwin_arm64` |
| Windows | x86_64 | `anki-builder_windows_amd64.exe` |

## Release Process

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

## Package Managers

#### Manual Installation Script
```bash
# Download and run installer
curl -fsSL https://raw.githubusercontent.com/commedesvlados/anki-flashcards-autogen-cli/main/scripts/install.sh | bash
```

## CI/CD Pipeline

The project uses GitHub Actions for automated builds:

- **Trigger**: Push to tags starting with `v*`
- **Builds**: All supported platforms
- **Artifacts**: Binary files, checksums, release notes
- **Deployment**: Automatic GitHub release creation 