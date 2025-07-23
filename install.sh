#!/bin/bash

set -e

# btui installer for Linux
# This script builds and installs the btui Bluetooth TUI application

INSTALL_DIR="/usr/local/bin"
BINARY_NAME="btui"
BUILD_DIR="."

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Print colored output
print_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

print_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Check if running as root for system-wide install
check_permissions() {
    if [[ $EUID -ne 0 ]]; then
        print_warn "Not running as root. Will attempt to install to $INSTALL_DIR using sudo."
        print_warn "You may be prompted for your password."
    fi
}

# Check if Go is installed
check_go() {
    if ! command -v go &> /dev/null; then
        print_error "Go is not installed or not in PATH"
        print_error "Please install Go 1.24.1 or later from https://golang.org/dl/"
        exit 1
    fi
    
    GO_VERSION=$(go version | grep -oE 'go[0-9]+\.[0-9]+' | cut -c3-)
    REQUIRED_VERSION="1.24"
    
    if ! awk -v ver="$GO_VERSION" -v req="$REQUIRED_VERSION" 'BEGIN{exit(ver<req)}'; then
        print_error "Go version $GO_VERSION found, but Go $REQUIRED_VERSION or later is required"
        exit 1
    fi
    
    print_info "Go version $GO_VERSION detected"
}

# Check if bluetoothctl is available
check_bluetooth() {
    if ! command -v bluetoothctl &> /dev/null; then
        print_warn "bluetoothctl not found. btui requires bluetoothctl to function."
        print_warn "Install it with: sudo apt install bluez-utils (Ubuntu/Debian) or sudo dnf install bluez (Fedora)"
    else
        print_info "bluetoothctl found"
    fi
}

# Remove existing installation
remove_existing() {
    if [[ -f "$INSTALL_DIR/$BINARY_NAME" ]]; then
        print_info "Removing existing btui installation..."
        if [[ $EUID -eq 0 ]]; then
            rm -f "$INSTALL_DIR/$BINARY_NAME"
        else
            sudo rm -f "$INSTALL_DIR/$BINARY_NAME"
        fi
        print_info "Existing installation removed"
    fi
}

# Build the application
build_app() {
    print_info "Building btui..."
    
    # Clean any existing binary in build directory
    if [[ -f "$BINARY_NAME" ]]; then
        rm "$BINARY_NAME"
    fi
    
    # Ensure dependencies are up to date
    go mod tidy
    
    # Build the binary
    go build -o "$BINARY_NAME" .
    
    if [[ ! -f "$BINARY_NAME" ]]; then
        print_error "Build failed - binary not created"
        exit 1
    fi
    
    print_info "Build successful"
}

# Install the binary
install_binary() {
    print_info "Installing btui to $INSTALL_DIR..."
    
    # Create install directory if it doesn't exist
    if [[ ! -d "$INSTALL_DIR" ]]; then
        if [[ $EUID -eq 0 ]]; then
            mkdir -p "$INSTALL_DIR"
        else
            sudo mkdir -p "$INSTALL_DIR"
        fi
    fi
    
    # Copy binary to install directory
    if [[ $EUID -eq 0 ]]; then
        cp "$BINARY_NAME" "$INSTALL_DIR/"
        chmod +x "$INSTALL_DIR/$BINARY_NAME"
    else
        sudo cp "$BINARY_NAME" "$INSTALL_DIR/"
        sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
    fi
    
    # Clean up build artifact
    rm "$BINARY_NAME"
    
    print_info "Installation complete"
}

# Verify installation
verify_install() {
    if command -v btui &> /dev/null; then
        VERSION_OUTPUT=$(btui --version 2>/dev/null || btui --help | head -n1 || echo "btui installed")
        print_info "btui installed successfully and is in PATH"
        print_info "$VERSION_OUTPUT"
        print_info "Run 'btui --help' to get started"
    else
        print_warn "btui installed but may not be in PATH"
        print_warn "You may need to add $INSTALL_DIR to your PATH or run $INSTALL_DIR/btui directly"
    fi
}

# Main installation process
main() {
    print_info "Starting btui installation..."
    
    check_permissions
    check_go
    check_bluetooth
    remove_existing
    build_app
    install_binary
    verify_install
    
    print_info "Installation completed successfully!"
}

# Run main function
main "$@"