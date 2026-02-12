#!/bin/bash
#
# Knit Installer - Binary installer for Linux
# https://github.com/Cheviiot/Knit
#
# Usage:
#   curl -fsSL https://raw.githubusercontent.com/Cheviiot/Knit/main/install.sh | bash
#   wget -qO- https://raw.githubusercontent.com/Cheviiot/Knit/main/install.sh | bash
#

set -e

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Configuration
REPO="Cheviiot/Knit"
APP_NAME="knit"
INSTALL_DIR="${HOME}/.local/bin"
ICON_DIR="${HOME}/.local/share/icons/hicolor"
DESKTOP_DIR="${HOME}/.local/share/applications"
VERSION="1.2.2"

print_banner() {
    echo -e "${CYAN}"
    echo "  ╔═══════════════════════════════════════╗"
    echo "  ║                                       ║"
    echo "  ║   █▄▀ █▄░█ █ ▀█▀                      ║"
    echo "  ║   █░█ █░▀█ █ ░█░                      ║"
    echo "  ║                                       ║"
    echo "  ║   Movie Torrent Search App            ║"
    echo "  ║                                       ║"
    echo "  ╚═══════════════════════════════════════╝"
    echo -e "${NC}"
}

info() {
    echo -e "${BLUE}►${NC} $1" >&2
}

success() {
    echo -e "${GREEN}✓${NC} $1" >&2
}

warn() {
    echo -e "${YELLOW}⚠${NC} $1" >&2
}

error() {
    echo -e "${RED}✗${NC} $1" >&2
    exit 1
}

# Check dependencies
check_deps() {
    if ! command -v curl &> /dev/null && ! command -v wget &> /dev/null; then
        error "Требуется curl или wget для загрузки"
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
            error "Неподдерживаемая архитектура: $arch"
            ;;
    esac
}

# Download file
download() {
    local url="$1"
    local output="$2"
    
    if command -v curl &> /dev/null; then
        curl -fsSL -o "$output" "$url"
    elif command -v wget &> /dev/null; then
        wget -qO "$output" "$url"
    else
        error "Нет инструмента для загрузки (curl или wget)"
    fi
}

# Get latest release binary URL
get_binary_url() {
    local api_url="https://api.github.com/repos/${REPO}/releases/latest"
    local release_info
    
    info "Получение информации о релизе..."
    
    if command -v curl &> /dev/null; then
        release_info=$(curl -fsSL "$api_url" 2>/dev/null)
    else
        release_info=$(wget -qO- "$api_url" 2>/dev/null)
    fi
    
    # Find binary named "knit" (not AppImage, not deb, not rpm)
    local download_url
    download_url=$(echo "$release_info" | grep -oE '"browser_download_url":\s*"[^"]+"' | grep -oE 'https://[^"]+' | grep -E '/knit$' | head -1 | tr -d '\r\n')
    
    echo "$download_url"
}

# Get icon URL from repo
get_icon_url() {
    echo "https://raw.githubusercontent.com/${REPO}/main/build/appicon.png"
}

# Install binary
install_binary() {
    local binary_url
    binary_url=$(get_binary_url)
    binary_url=$(echo "$binary_url" | tr -d '[:space:]')
    local binary_path="${INSTALL_DIR}/${APP_NAME}-bin"
    
    # Create install directory
    mkdir -p "$INSTALL_DIR"
    
    if [ -n "$binary_url" ]; then
        info "Загрузка Knit из $binary_url..."
        download "$binary_url" "$binary_path"
    else
        warn "Бинарник не найден в релизе."
        error "Скачайте вручную: https://github.com/${REPO}/releases"
    fi
    
    # Make executable
    chmod +x "$binary_path"
    success "Бинарник установлен: $binary_path"
}

# Install icon
install_icon() {
    local icon_url=$(get_icon_url)
    local tmpfile=$(mktemp)
    
    info "Установка иконок..."
    
    # Download icon
    download "$icon_url" "$tmpfile" 2>/dev/null || {
        warn "Не удалось загрузить иконку"
        rm -f "$tmpfile"
        return 0
    }
    
    # Create icon directories and copy
    for size in 512 256 128 64 48 32 16; do
        mkdir -p "${ICON_DIR}/${size}x${size}/apps"
        
        # Resize if ImageMagick available, otherwise copy original
        if command -v magick &> /dev/null; then
            magick "$tmpfile" -resize ${size}x${size} "${ICON_DIR}/${size}x${size}/apps/${APP_NAME}.png" 2>/dev/null || \
            cp "$tmpfile" "${ICON_DIR}/${size}x${size}/apps/${APP_NAME}.png" 2>/dev/null || true
        elif command -v convert &> /dev/null; then
            convert "$tmpfile" -resize ${size}x${size} "${ICON_DIR}/${size}x${size}/apps/${APP_NAME}.png" 2>/dev/null || \
            cp "$tmpfile" "${ICON_DIR}/${size}x${size}/apps/${APP_NAME}.png" 2>/dev/null || true
        else
            cp "$tmpfile" "${ICON_DIR}/${size}x${size}/apps/${APP_NAME}.png" 2>/dev/null || true
        fi
    done
    
    # Also put in pixmaps
    mkdir -p "${HOME}/.local/share/pixmaps"
    cp "$tmpfile" "${HOME}/.local/share/pixmaps/${APP_NAME}.png" 2>/dev/null || true
    
    rm -f "$tmpfile"
    
    # Update icon cache
    if command -v gtk-update-icon-cache &> /dev/null; then
        gtk-update-icon-cache -f -t "$ICON_DIR" 2>/dev/null || true
    fi
    
    success "Иконки установлены"
}

# Create wrapper script with env vars
create_wrapper() {
    local wrapper_path="${INSTALL_DIR}/${APP_NAME}"
    
    info "Создание launcher скрипта..."
    
    # Create wrapper script
    cat > "$wrapper_path" << 'WRAPPER_EOF'
#!/bin/bash
# Knit launcher - sets environment variables for proper cursor display
export WEBKIT_DISABLE_COMPOSITING_MODE=1
export XCURSOR_PATH="${XCURSOR_PATH:-/usr/share/icons}"
exec "$(dirname "$0")/knit-bin" "$@"
WRAPPER_EOF
    
    chmod +x "$wrapper_path"
    success "Launcher создан: $wrapper_path"
}

# Create desktop entry
create_desktop_entry() {
    local binary_path="$1"
    
    info "Создание ярлыка приложения..."
    
    mkdir -p "$DESKTOP_DIR"
    
    cat > "${DESKTOP_DIR}/${APP_NAME}.desktop" << EOF
[Desktop Entry]
Type=Application
Name=Knit
GenericName=Movie Torrent Search
Comment=Elegant desktop application for searching movie torrents
Exec=env WEBKIT_DISABLE_COMPOSITING_MODE=1 ${binary_path} %U
Icon=${APP_NAME}
Categories=AudioVideo;Video;Network;
Terminal=false
Keywords=movies;torrents;films;video;download;knit;
Version=${VERSION}
StartupNotify=true
StartupWMClass=knit
MimeType=x-scheme-handler/magnet;
EOF
    
    chmod +x "${DESKTOP_DIR}/${APP_NAME}.desktop"
    
    # Update desktop database
    if command -v update-desktop-database &> /dev/null; then
        update-desktop-database "$DESKTOP_DIR" 2>/dev/null || true
    fi
    
    success "Ярлык создан"
}

# Check if PATH includes install dir
check_path() {
    if [[ ":$PATH:" != *":${INSTALL_DIR}:"* ]]; then
        echo ""
        warn "Добавьте ${INSTALL_DIR} в PATH:"
        echo ""
        echo -e "  ${CYAN}# Добавьте в ~/.bashrc или ~/.zshrc:${NC}"
        echo "  export PATH=\"\$HOME/.local/bin:\$PATH\""
        echo ""
    fi
}

# Uninstall function
uninstall() {
    info "Удаление Knit..."
    
    rm -f "${INSTALL_DIR}/${APP_NAME}"
    rm -f "${INSTALL_DIR}/${APP_NAME}-bin"
    rm -f "${INSTALL_DIR}/${APP_NAME}.AppImage"
    rm -f "${DESKTOP_DIR}/${APP_NAME}.desktop"
    rm -f "${HOME}/.local/share/pixmaps/${APP_NAME}.png"
    
    for size in 512 256 128 64 48 32 16; do
        rm -f "${ICON_DIR}/${size}x${size}/apps/${APP_NAME}.png"
    done
    
    # Update caches
    if command -v gtk-update-icon-cache &> /dev/null; then
        gtk-update-icon-cache -f -t "$ICON_DIR" 2>/dev/null || true
    fi
    if command -v update-desktop-database &> /dev/null; then
        update-desktop-database "$DESKTOP_DIR" 2>/dev/null || true
    fi
    
    success "Knit удалён"
}

# Update function
update() {
    info "Обновление Knit..."
    
    # Remove old binary
    rm -f "${INSTALL_DIR}/${APP_NAME}-bin"
    
    # Download new
    install_binary
    create_wrapper
    
    success "Knit обновлён до последней версии"
}

# Show help
show_help() {
    echo "Knit Installer v${VERSION}"
    echo ""
    echo "Использование:"
    echo "  curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | bash"
    echo ""
    echo "Опции:"
    echo "  --update, -U     Обновить Knit"
    echo "  --uninstall, -u  Удалить Knit"
    echo "  --help, -h       Показать справку"
    echo ""
    echo "Примеры:"
    echo "  # Установка"
    echo "  curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | bash"
    echo ""
    echo "  # Обновление"  
    echo "  curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | bash -s -- --update"
    echo ""
    echo "  # Удаление"
    echo "  curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | bash -s -- --uninstall"
    echo ""
}

# Main
main() {
    # Parse arguments
    while [[ $# -gt 0 ]]; do
        case $1 in
            --update|-U)
                print_banner
                check_deps
                update
                exit 0
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
    
    # Check dependencies
    check_deps
    
    # Detect architecture
    local arch=$(detect_arch)
    info "Архитектура: $arch"
    
    # Install
    install_binary
    create_wrapper
    install_icon
    create_desktop_entry "${INSTALL_DIR}/${APP_NAME}"
    
    echo ""
    echo -e "${GREEN}╔═══════════════════════════════════════╗${NC}"
    echo -e "${GREEN}║     Knit успешно установлен!          ║${NC}"
    echo -e "${GREEN}╚═══════════════════════════════════════╝${NC}"
    echo ""
    echo -e "  Запуск из терминала:  ${CYAN}${APP_NAME}${NC}"
    echo -e "  Или найдите ${CYAN}'Knit'${NC} в меню приложений"
    echo ""
    echo -e "  ${YELLOW}Обновить:${NC} curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | bash -s -- --update"
    echo -e "  ${YELLOW}Удалить:${NC}  curl -fsSL https://raw.githubusercontent.com/${REPO}/main/install.sh | bash -s -- --uninstall"
    echo ""
    
    check_path
}

main "$@"
