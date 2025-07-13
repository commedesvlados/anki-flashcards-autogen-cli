#!/bin/bash

# Release script for anki-builder
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if we're in a git repository
if ! git rev-parse --git-dir > /dev/null 2>&1; then
    echo -e "${RED}Error: Not in a git repository${NC}"
    exit 1
fi

# Check if we have uncommitted changes
if ! git diff-index --quiet HEAD --; then
    echo -e "${YELLOW}Warning: You have uncommitted changes${NC}"
    read -p "Continue anyway? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

# Get current version from git tags
CURRENT_VERSION=$(git describe --tags --abbrev=0 2>/dev/null || echo "v0.1.0")
echo -e "${GREEN}Current version: ${CURRENT_VERSION}${NC}"

# Prompt for new version
read -p "Enter new version (e.g., v1.0.0): " NEW_VERSION

if [[ -z "$NEW_VERSION" ]]; then
    echo -e "${RED}Error: Version cannot be empty${NC}"
    exit 1
fi

# Validate version format
if [[ ! $NEW_VERSION =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo -e "${RED}Error: Version must be in format vX.Y.Z${NC}"
    exit 1
fi

echo -e "${GREEN}Building release for version: ${NEW_VERSION}${NC}"

# Build for all platforms
echo -e "${YELLOW}Building binaries for all platforms...${NC}"
make release

# Create release directory
RELEASE_DIR="release_${NEW_VERSION}"
mkdir -p "$RELEASE_DIR"

# Copy binaries to release directory
echo -e "${YELLOW}Preparing release files...${NC}"
cp build/anki-builder_* "$RELEASE_DIR/"

# Create checksums
echo -e "${YELLOW}Creating checksums...${NC}"
cd "$RELEASE_DIR"
for file in anki-builder_*; do
    if [[ "$OSTYPE" == "darwin"* ]]; then
        shasum -a 256 "$file" > "${file}.sha256"
    else
        sha256sum "$file" > "${file}.sha256"
    fi
done
cd ..

# Create release notes
cat > "$RELEASE_DIR/RELEASE_NOTES.md" << EOF
# Anki Builder ${NEW_VERSION}

## Installation

### Linux
\`\`\`bash
# Download the appropriate binary for your architecture
# For x86_64:
wget https://github.com/yourusername/anki-flashcards-autogen-cli/releases/download/${NEW_VERSION}/anki-builder_linux_amd64
chmod +x anki-builder_linux_amd64
sudo mv anki-builder_linux_amd64 /usr/local/bin/anki-builder

# For ARM64:
wget https://github.com/yourusername/anki-flashcards-autogen-cli/releases/download/${NEW_VERSION}/anki-builder_linux_arm64
chmod +x anki-builder_linux_arm64
sudo mv anki-builder_linux_arm64 /usr/local/bin/anki-builder
\`\`\`

### macOS
\`\`\`bash
# For Intel Macs:
wget https://github.com/yourusername/anki-flashcards-autogen-cli/releases/download/${NEW_VERSION}/anki-builder_darwin_amd64
chmod +x anki-builder_darwin_amd64
sudo mv anki-builder_darwin_amd64 /usr/local/bin/anki-builder

# For Apple Silicon:
wget https://github.com/yourusername/anki-flashcards-autogen-cli/releases/download/${NEW_VERSION}/anki-builder_darwin_arm64
chmod +x anki-builder_darwin_arm64
sudo mv anki-builder_darwin_arm64 /usr/local/bin/anki-builder
\`\`\`

### Windows
Download the appropriate .exe file and run it from command prompt or PowerShell.

## Usage
\`\`\`bash
anki-builder --help
\`\`\`

## Requirements
- Python 3 with genanki library
- Unsplash API key
- Excel file with Russian-English word pairs

## Changes in this release
- [Add your changes here]

## Checksums
EOF

# Add checksums to release notes
for file in "$RELEASE_DIR"/anki-builder_*.sha256; do
    echo "  - $(basename "$file" .sha256): $(cat "$file" | cut -d' ' -f1)" >> "$RELEASE_DIR/RELEASE_NOTES.md"
done

echo -e "${GREEN}Release prepared in: ${RELEASE_DIR}${NC}"
echo -e "${YELLOW}Next steps:${NC}"
echo "1. Review the release files in $RELEASE_DIR"
echo "2. Create a git tag: git tag $NEW_VERSION"
echo "3. Push the tag: git push origin $NEW_VERSION"
echo "4. Create a GitHub release with the files from $RELEASE_DIR"
echo "5. Update the RELEASE_NOTES.md with actual changes"

# Ask if user wants to create git tag
read -p "Create git tag $NEW_VERSION now? (y/N): " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    git tag "$NEW_VERSION"
    echo -e "${GREEN}Tag created: $NEW_VERSION${NC}"
    echo -e "${YELLOW}Push with: git push origin $NEW_VERSION${NC}"
fi 