#!/bin/bash
set -e

# Vibe Skills Installer
# Usage: curl -sSL https://raw.githubusercontent.com/cuongtl1992/vibe-skills/main/scripts/install.sh | bash

REPO="cuongtl/vibe-skills"
BINARY_NAME="vibe-skills"
INSTALL_DIR="/usr/local/bin"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

# Detect OS and architecture
detect_platform() {
    OS=$(uname -s | tr '[:upper:]' '[:lower:]')
    ARCH=$(uname -m)

    case "$ARCH" in
        x86_64|amd64)
            ARCH="x86_64"
            ;;
        arm64|aarch64)
            ARCH="arm64"
            ;;
        *)
            error "Unsupported architecture: $ARCH"
            ;;
    esac

    case "$OS" in
        linux)
            OS="linux"
            ;;
        darwin)
            OS="darwin"
            ;;
        mingw*|msys*|cygwin*)
            OS="windows"
            ;;
        *)
            error "Unsupported OS: $OS"
            ;;
    esac

    echo "${OS}_${ARCH}"
}

# Get latest release version
get_latest_version() {
    curl -s "https://api.github.com/repos/${REPO}/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/'
}

# Download and install
install() {
    PLATFORM=$(detect_platform)
    VERSION=$(get_latest_version)

    if [ -z "$VERSION" ]; then
        error "Failed to get latest version"
    fi

    info "Installing vibe-skills ${VERSION} for ${PLATFORM}..."

    # Determine file extension
    EXT="tar.gz"
    if [[ "$PLATFORM" == windows* ]]; then
        EXT="zip"
    fi

    DOWNLOAD_URL="https://github.com/${REPO}/releases/download/${VERSION}/${BINARY_NAME}_${PLATFORM}.${EXT}"

    info "Downloading from ${DOWNLOAD_URL}..."

    # Create temp directory
    TMP_DIR=$(mktemp -d)
    trap "rm -rf $TMP_DIR" EXIT

    # Download
    if command -v curl &> /dev/null; then
        curl -sL "$DOWNLOAD_URL" -o "$TMP_DIR/archive.${EXT}"
    elif command -v wget &> /dev/null; then
        wget -q "$DOWNLOAD_URL" -O "$TMP_DIR/archive.${EXT}"
    else
        error "Neither curl nor wget found. Please install one of them."
    fi

    # Extract
    cd "$TMP_DIR"
    if [ "$EXT" = "tar.gz" ]; then
        tar -xzf "archive.${EXT}"
    else
        unzip -q "archive.${EXT}"
    fi

    # Install
    if [ -w "$INSTALL_DIR" ]; then
        mv "$BINARY_NAME" "$INSTALL_DIR/"
    else
        info "Need sudo to install to $INSTALL_DIR"
        sudo mv "$BINARY_NAME" "$INSTALL_DIR/"
    fi

    chmod +x "$INSTALL_DIR/$BINARY_NAME"

    info "Successfully installed vibe-skills to $INSTALL_DIR/$BINARY_NAME"
    info "Run 'vibe-skills --help' to get started"
}

# Check if already installed
check_existing() {
    if command -v vibe-skills &> /dev/null; then
        CURRENT_VERSION=$(vibe-skills version 2>/dev/null | head -1 || echo "unknown")
        warn "vibe-skills is already installed: $CURRENT_VERSION"
        read -p "Do you want to update? [y/N] " -n 1 -r
        echo
        if [[ ! $REPLY =~ ^[Yy]$ ]]; then
            info "Installation cancelled"
            exit 0
        fi
    fi
}

# Main
main() {
    info "Vibe Skills Installer"
    check_existing
    install
}

main
