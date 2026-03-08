#!/bin/bash

# say-mi installer
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/josiaranda/say-mi/main/install.sh | bash
#   curl -fsSL https://raw.githubusercontent.com/josiaranda/say-mi/main/install.sh | bash -s -- v1.0.0

set -e

REPO="josiaranda/say-mi"
BINARY_NAME="say-mi"
INSTALL_DIR="/usr/local/bin"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
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

    case "$OS" in
        darwin) OS="darwin" ;;
        linux)  OS="linux" ;;
        *)      error "Unsupported OS: $OS" ;;
    esac

    case "$ARCH" in
        x86_64|amd64)   ARCH="amd64" ;;
        arm64|aarch64)  ARCH="arm64" ;;
        *)              error "Unsupported architecture: $ARCH" ;;
    esac

    PLATFORM="${OS}-${ARCH}"
    info "Detected platform: $PLATFORM"
}

# Get latest release version
get_latest_version() {
    info "Fetching latest version..."
    VERSION=$(curl -fsSL "https://api.github.com/repos/$REPO/releases/latest" | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

    if [ -z "$VERSION" ]; then
        error "Failed to fetch latest version"
    fi

    info "Latest version: $VERSION"
}

# Get current installed version
get_installed_version() {
    if command -v say-mi &> /dev/null; then
        INSTALLED=$(say-mi --version 2>/dev/null | head -1 | awk '{print $2}')
        if [ -n "$INSTALLED" ]; then
            info "Installed version: $INSTALLED"
        fi
    fi
}

# Download binary
download_binary() {
    DOWNLOAD_URL="https://github.com/$REPO/releases/download/$VERSION/say-mi-$PLATFORM"

    if [ "$OS" = "windows" ]; then
        DOWNLOAD_URL="${DOWNLOAD_URL}.exe"
        BINARY_NAME="${BINARY_NAME}.exe"
    fi

    info "Downloading from: $DOWNLOAD_URL"

    TMP_DIR=$(mktemp -d)
    TMP_FILE="$TMP_DIR/$BINARY_NAME"

    if ! curl -fsSL --progress-bar "$DOWNLOAD_URL" -o "$TMP_FILE"; then
        error "Failed to download binary"
    fi

    chmod +x "$TMP_FILE"
    info "Downloaded to: $TMP_FILE"
}

# Install binary
install_binary() {
    if [ -w "$INSTALL_DIR" ]; then
        mv "$TMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
    else
        info "sudo required to install to $INSTALL_DIR"
        sudo mv "$TMP_FILE" "$INSTALL_DIR/$BINARY_NAME"
    fi

    info "Installed to: $INSTALL_DIR/$BINARY_NAME"
}

# Cleanup
cleanup() {
    rm -rf "$TMP_DIR"
}

# Verify installation
verify() {
    if command -v say-mi &> /dev/null; then
        NEW_VERSION=$(say-mi --version 2>/dev/null | head -1 | awk '{print $2}')
        info "✓ say-mi $NEW_VERSION installed successfully!"
    else
        warn "say-mi was installed but may not be in your PATH"
        info "Add $INSTALL_DIR to your PATH or restart your terminal"
    fi
}

# Show help
show_help() {
    echo "say-mi installer"
    echo ""
    echo "Usage:"
    echo "  curl -fsSL https://raw.githubusercontent.com/josiaranda/say-mi/main/install.sh | bash"
    echo "  curl -fsSL https://raw.githubusercontent.com/josiaranda/say-mi/main/install.sh | bash -s -- v1.0.0"
    echo ""
    echo "Options:"
    echo "  -h, --help      Show this help"
    echo "  -u, --upgrade   Upgrade to latest version (same as default)"
    echo "  v1.0.0          Install specific version"
    echo ""
    echo "Uninstall:"
    echo "  rm /usr/local/bin/say-mi"
    exit 0
}

main() {
    echo ""
    echo -e "${BLUE}       __  __       _ __  __  _ ${NC}"
    echo -e "${BLUE}  ___ / _|/ _| ___ | |\ \/ / (_)${NC}"
    echo -e "${BLUE} / __| |_| |_ / _ \| | \  /  | |${NC}"
    echo -e "${BLUE}| (__|  _|  _| (_) | | /  \  | |${NC}"
    echo -e "${BLUE} \___|_| |_|  \___/|_|/_/\_\ |_|${NC}"
    echo ""
    echo "  josiaranda/say-mi"
    echo ""

    # Parse arguments
    VERSION=""
    while [[ $# -gt 0 ]]; do
        case $1 in
            -h|--help)
                show_help
                ;;
            -u|--upgrade)
                shift
                ;;
            v*)
                VERSION=$1
                shift
                ;;
            *)
                shift
                ;;
        esac
    done

    detect_platform
    get_installed_version

    # Get version (latest or specified)
    if [ -z "$VERSION" ]; then
        get_latest_version
    else
        info "Installing version: $VERSION"
    fi

    download_binary
    install_binary
    cleanup
    verify

    echo ""
    info "Done! Run 'say-mi --help' to get started."
    info "To update, run this script again."
}

main "$@"
