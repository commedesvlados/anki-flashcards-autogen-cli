#!/bin/bash

# Installation script for anki-builder
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
BINARY_NAME="anki-builder"
REPO_URL="https://github.com/yourusername/anki-flashcards-autogen-cli"
LATEST_RELEASE_URL="https://api.github.com/repos/yourusername/anki-flashcards-autogen-cli/releases/latest"

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64|arm64)
        ARCH="arm64"
        ;;
    *)
        echo -e "${RED}Unsupported architecture: $ARCH${NC}"
        exit 1
        ;;
esac

# Set binary suffix
if [[ "$OS" == "windows"* ]]; then
    BINARY_SUFFIX="${OS}_${ARCH}.exe"
    DOWNLOAD_NAME="${BINARY_NAME}_${OS}_${ARCH}.exe"
else
    BINARY_SUFFIX="${OS}_${ARCH}"
    DOWNLOAD_NAME="${BINARY_NAME}_${OS}_${ARCH}"
fi

echo -e "${BLUE}Anki Builder CLI Installer${NC}"
echo -e "${BLUE}========================${NC}"
echo -e "OS: $OS"
echo -e "Architecture: $ARCH"
echo -e "Binary: $DOWNLOAD_NAME"
echo

# Function to get latest version
get_latest_version() {
    if command -v curl >/dev/null 2>&1; then
        curl -s "$LATEST_RELEASE_URL" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
    elif command -v wget >/dev/null 2>&1; then
        wget -qO- "$LATEST_RELEASE_URL" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
    else
        echo "v1.0.0"  # fallback
    fi
}

# Function to download binary
download_binary() {
    local version=$1
    local download_url="https://github.com/yourusername/anki-flashcards-autogen-cli/releases/download/${version}/${DOWNLOAD_NAME}"
    
    echo -e "${YELLOW}Downloading ${BINARY_NAME} ${version}...${NC}"
    
    if command -v curl >/dev/null 2>&1; then
        curl -L -o "$BINARY_NAME" "$download_url"
    elif command -v wget >/dev/null 2>&1; then
        wget -O "$BINARY_NAME" "$download_url"
    else
        echo -e "${RED}Error: Neither curl nor wget is installed${NC}"
        exit 1
    fi
    
    if [[ ! -f "$BINARY_NAME" ]]; then
        echo -e "${RED}Error: Failed to download binary${NC}"
        exit 1
    fi
}

# Function to install binary
install_binary() {
    local install_dir=$1
    
    echo -e "${YELLOW}Installing to $install_dir...${NC}"
    
    # Create directory if it doesn't exist
    mkdir -p "$install_dir"
    
    # Move binary
    mv "$BINARY_NAME" "$install_dir/"
    chmod +x "$install_dir/$BINARY_NAME"
    
    echo -e "${GREEN}Installation complete!${NC}"
    echo -e "${GREEN}Binary installed at: $install_dir/$BINARY_NAME${NC}"
    
    # Check if directory is in PATH
    if [[ ":$PATH:" != *":$install_dir:"* ]]; then
        echo -e "${YELLOW}Note: $install_dir is not in your PATH${NC}"
        echo -e "${YELLOW}Add this line to your shell profile (.bashrc, .zshrc, etc.):${NC}"
        echo -e "${BLUE}export PATH=\"\$PATH:$install_dir\"${NC}"
    else
        echo -e "${GREEN}You can now run: $BINARY_NAME --help${NC}"
    fi
}

# Main installation logic
main() {
    # Get latest version
    VERSION=$(get_latest_version)
    echo -e "${BLUE}Latest version: $VERSION${NC}"
    
    # Ask for installation directory
    echo -e "${YELLOW}Choose installation method:${NC}"
    echo "1) System-wide installation (requires sudo) - /usr/local/bin"
    echo "2) User installation - ~/.local/bin"
    echo "3) Custom directory"
    read -p "Enter choice (1-3): " choice
    
    case $choice in
        1)
            if [[ "$EUID" -eq 0 ]]; then
                install_dir="/usr/local/bin"
            else
                echo -e "${YELLOW}System installation requires sudo privileges${NC}"
                read -p "Continue with sudo? (y/N): " -n 1 -r
                echo
                if [[ ! $REPLY =~ ^[Yy]$ ]]; then
                    exit 1
                fi
                install_dir="/usr/local/bin"
            fi
            ;;
        2)
            install_dir="$HOME/.local/bin"
            ;;
        3)
            read -p "Enter custom installation directory: " install_dir
            if [[ -z "$install_dir" ]]; then
                echo -e "${RED}Error: Installation directory cannot be empty${NC}"
                exit 1
            fi
            ;;
        *)
            echo -e "${RED}Invalid choice${NC}"
            exit 1
            ;;
    esac
    
    # Download and install
    download_binary "$VERSION"
    
    if [[ "$choice" == "1" && "$EUID" -ne 0 ]]; then
        sudo install_binary "$install_dir"
    else
        install_binary "$install_dir"
    fi
    
    echo
    echo -e "${GREEN}Installation completed successfully!${NC}"
    echo -e "${BLUE}Run '$BINARY_NAME --help' to get started${NC}"
}

# Run main function
main "$@" 