#!/bin/bash
#
# Knit Installer - Universal AppImage installer for Linux
# https://github.com/Cheviiot/Knit
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/Cheviiot/Knit/main/install.sh | bash
#   curl -fsSL https://raw.githubusercontent.com/Cheviiot/Knit/main/install.sh | bash -s -- --icons
#   wget -qO- https://raw.githubusercontent.com/Cheviiot/Knit/main/install.sh | bash
#

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Configuration
REPO="Cheviiot/Knit"
APP_NAME="knit"
INSTALL_DIR="${HOME}/.local/bin"
ICON_DIR="${HOME}/.local/share/icons/hicolor"
DESKTOP_DIR="${HOME}/.local/share/applications"

# Flags
INSTALL_ICONS=false

print_banner() {
    echo -e "${BLUE}"
    echo "  _  __      _ _   "
    echo " | |/ /_ __ (_) |_ "
    echo " | ' /| '_ \| | __|"
    echo " | . \| | | | | |_ "
    echo " |_|\_\_| |_|_|\__|"
    echo -e "${NC}"
    echo "Movie Torrent Search App"
    echo ""
}

info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

success() {
    echo -e "${GREEN}[OK]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1"
    exit 1
}

# Check dependencies
check_deps() {
    local missing=""
    
    if ! command -v curl &> /dev/null && ! command -v wget &> /dev/null; then
        missing="curl or wget"
    fi
    
    if [ -n "$missing" ]; then
        error "Missing required tools: $missing"
    fi
}

# Detect architecture
detect_arch() {
    local arch=$(uname -m)
    case "$arch" in
        x86_64|amd64)
            echo "x86_64"
            ;;
        aarch64|arm64)
            echo "aarch64"
            ;;
        *)
            error "Unsupported architecture: $arch"
            ;;
    esac
}

# Download file
download() {
    local url="$1"
    local output="$2"
    
    if command -v curl &> /dev/null; then
        curl -fsSL -g "$url" -o "$output"
    elif command -v wget &> /dev/null; then
        wget -qO "$output" "$url"
    else
        error "No download tool available (curl or wget)"
    fi
}

# Get latest release URL
get_latest_release() {
    local arch="$1"
    local api_url="https://api.github.com/repos/${REPO}/releases/latest"
    local release_info
    
    echo -e "${BLUE}[INFO]${NC} Fetching latest release info..." >&2
    
    if command -v curl &> /dev/null; then
        release_info=$(curl -fsSL -g "$api_url")
    else
        release_info=$(wget -qO- "$api_url")
    fi
    
    # Try to find AppImage for architecture
    local download_url=$(echo "$release_info" | grep -oE "https://[^\"]+${APP_NAME}[^\"]*${arch}[^\"]*\.AppImage" | head -1)
    
    # If not found, try case-insensitive
    if [ -z "$download_url" ]; then
        download_url=$(echo "$release_info" | grep -oiE "https://[^\"]+\.AppImage" | grep -i "${arch}\|amd64\|x86_64" | head -1)
    fi
    
    # If still not found, get any AppImage
    if [ -z "$download_url" ]; then
        download_url=$(echo "$release_info" | grep -oE "https://[^\"]+\.AppImage" | head -1)
    fi
    
    if [ -z "$download_url" ]; then
        error "Could not find AppImage in latest release. Please download manually from https://github.com/${REPO}/releases"
    fi
    
    echo "$download_url"
}

# Install AppImage
install_appimage() {
    local arch=$(detect_arch)
    local download_url=$(get_latest_release "$arch")
    local appimage_path="${INSTALL_DIR}/${APP_NAME}.AppImage"
    
    echo -e "${BLUE}[INFO]${NC} Architecture: $arch" >&2
    echo -e "${BLUE}[INFO]${NC} Download URL: $download_url" >&2
    
    # Create install directory
    mkdir -p "$INSTALL_DIR"
    
    # Download AppImage
    echo -e "${BLUE}[INFO]${NC} Downloading Knit..." >&2
    download "$download_url" "$appimage_path"
    
    # Make executable
    chmod +x "$appimage_path"
    echo -e "${GREEN}[OK]${NC} AppImage downloaded to $appimage_path" >&2
    
    echo "$appimage_path"
}

# Extract and install icons
install_icons() {
    local appimage_path="$1"
    local tmpdir=$(mktemp -d)
    
    info "Installing icons..."
    
    # Extract AppImage
    cd "$tmpdir"
    "$appimage_path" --appimage-extract > /dev/null 2>&1 || true
    
    # Find icon
    local icon_src=""
    if [ -f "$tmpdir/squashfs-root/.DirIcon" ]; then
        icon_src="$tmpdir/squashfs-root/.DirIcon"
    elif [ -f "$tmpdir/squashfs-root/${APP_NAME}.png" ]; then
        icon_src="$tmpdir/squashfs-root/${APP_NAME}.png"
    fi
    
    if [ -n "$icon_src" ] && [ -f "$icon_src" ]; then
        # Create icon directories
        for size in 256 128 64 48 32 16; do
            mkdir -p "${ICON_DIR}/${size}x${size}/apps"
        done
        
        # Resize icons using ImageMagick if available
        if command -v magick &> /dev/null; then
            for size in 256 128 64 48 32 16; do
                magick "$icon_src" -resize ${size}x${size} "${ICON_DIR}/${size}x${size}/apps/${APP_NAME}.png" 2>/dev/null || true
            done
            success "Icons installed (resized)"
        elif command -v convert &> /dev/null; then
            for size in 256 128 64 48 32 16; do
                convert "$icon_src" -resize ${size}x${size} "${ICON_DIR}/${size}x${size}/apps/${APP_NAME}.png" 2>/dev/null || true
            done
            success "Icons installed (resized)"
        else
            # Just copy original to all sizes
            for size in 256 128 64 48 32 16; do
                cp "$icon_src" "${ICON_DIR}/${size}x${size}/apps/${APP_NAME}.png" 2>/dev/null || true
            done
            success "Icons installed (original size)"
        fi
        
        # Update icon cache
        if command -v gtk-update-icon-cache &> /dev/null; then
            gtk-update-icon-cache -f -t "$ICON_DIR" 2>/dev/null || true
        fi
    else
        warn "Could not extract icons from AppImage"
    fi
    
    # Cleanup
    rm -rf "$tmpdir"
}

# Create desktop entry
create_desktop_entry() {
    local appimage_path="$1"
    
    info "Creating desktop entry..."
    
    mkdir -p "$DESKTOP_DIR"
    
    cat > "${DESKTOP_DIR}/${APP_NAME}.desktop" << EOF
[Desktop Entry]
Type=Application
Name=Knit
GenericName=Movie Torrent Search
Comment=Elegant desktop application for searching movie torrents
Exec=${appimage_path} %U
Icon=${APP_NAME}
Categories=AudioVideo;Video;Network;
Terminal=false
Keywords=movies;torrents;films;video;download;
Version=1.0
StartupNotify=true
StartupWMClass=knit
MimeType=x-scheme-handler/magnet;
EOF
    
    # Update desktop database
    if command -v update-desktop-database &> /dev/null; then
        update-desktop-database "$DESKTOP_DIR" 2>/dev/null || true
    fi
    
    success "Desktop entry created"
}

# Create symlink for CLI access
create_symlink() {
    local appimage_path="$1"
    local symlink_path="${INSTALL_DIR}/${APP_NAME}"
    
    # Create symlink if not exists
    if [ ! -L "$symlink_path" ] && [ ! -f "$symlink_path" ]; then
        ln -sf "$appimage_path" "$symlink_path"
        success "Symlink created: $symlink_path"
    fi
}

# Check if PATH includes install dir
check_path() {
    if [[ ":$PATH:" != *":${INSTALL_DIR}:"* ]]; then
        echo ""
        warn "Add ${INSTALL_DIR} to your PATH for CLI access:"
        echo ""
        echo "  # Add to ~/.bashrc or ~/.zshrc:"
        echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
        echo ""
    fi
}

# Uninstall function
uninstall() {
    info "Uninstalling Knit..."
    
    rm -f "${INSTALL_DIR}/${APP_NAME}.AppImage"
    rm -f "${INSTALL_DIR}/${APP_NAME}"
    rm -f "${DESKTOP_DIR}/${APP_NAME}.desktop"
    
    for size in 256 128 64 48 32 16; do
        rm -f "${ICON_DIR}/${size}x${size}/apps/${APP_NAME}.png"
    done
    
    # Update caches
    if command -v gtk-update-icon-cache &> /dev/null; then
        gtk-update-icon-cache -f -t "$ICON_DIR" 2>/dev/null || true
    fi
    if command -v update-desktop-database &> /dev/null; then
        update-desktop-database "$DESKTOP_DIR" 2>/dev/null || true
    fi
    
    success "Knit uninstalled"
}

# Install icons only (for distros where icons don't show)
install_icons_only() {
    local appimage_path="${INSTALL_DIR}/${APP_NAME}.AppImage"
    
    if [ ! -f "$appimage_path" ]; then
        error "Knit not installed. Run installer first without --icons flag."
    fi
    
    info "Installing icons for Knit..."
    install_icons "$appimage_path"
    
    echo ""
    success "Icons installed!"
    echo "  Restart your desktop session or run: gtk-update-icon-cache -f -t ~/.local/share/icons/hicolor"
    echo ""
}

# Show help
show_help() {
    echo "Knit Installer"
    echo ""
    echo "Usage:"
    echo "  curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | bash"
    echo ""
    echo "Options:"
    echo "  --icons, -i      Install icons only (for distros where icon doesn't show)"
    echo "  --uninstall, -u  Uninstall Knit"
    echo "  --help, -h       Show this help"
    echo ""
    echo "Examples:"
    echo "  # Install Knit"
    echo "  curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | bash"
    echo ""
    echo "  # Install icons (if icon doesn't appear in menu)"
    echo "  curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | bash -s -- --icons"
    echo ""
    echo "  # Uninstall"
    echo "  curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | bash -s -- --uninstall"
    echo ""
}

# Main
main() {
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --icons|-i)
                INSTALL_ICONS=true
                shift
                ;;
            --uninstall|-u)
                print_banner
                uninstall
                exit 0
                ;;
            --help|-h)
                show_help
                exit 0
                ;;
            *)
                shift
                ;;
        esac
    done
    
    print_banner
    
    # Handle icons-only install
    if [ "$INSTALL_ICONS" = true ]; then
        install_icons_only
        exit 0
    fi
    
    # Check dependencies
    check_deps
    
    # Install
    local appimage_path=$(install_appimage)
    create_desktop_entry "$appimage_path"
    create_symlink "$appimage_path"
    
    echo ""
    success "Knit installed successfully!"
    echo ""
    echo "  Run from terminal:  ${APP_NAME}"
    echo "  Or find 'Knit' in your application menu"
    echo ""
    echo -e "  ${YELLOW}If icon doesn't appear:${NC}"
    echo "  curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | bash -s -- --icons"
    echo ""
    echo "  Uninstall: curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | bash -s -- --uninstall"
    echo ""
    
    check_path
}

main "$@"
