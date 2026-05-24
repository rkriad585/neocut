#!/usr/bin/env bash
set -euo pipefail

PROJECT_NAME="neocut"
GITHUB_USER="rkriad585"
CONFIG_DIR="${HOME}/.config/neostore/${PROJECT_NAME}"
BIN_DIR="${CONFIG_DIR}/bin"
BINARY_PATH="${BIN_DIR}/${PROJECT_NAME}"

# ── Helpers ─────────────────────────────────────────────────
log()   { printf "  %s %s\n" "·" "$1"; }
ok()    { printf "  \xE2\x9C\x93 %s\n" "$1"; }
fail()  { printf "  \xE2\x9C\x97 %s\n" "$1"; exit 1; }
skip()  { printf "  \xE2\x80\xA3 %s\n" "$1"; }

cleanup() {
    [ -n "${TMP_FILE:-}" ] && [ -f "$TMP_FILE" ] && rm -f "$TMP_FILE"
    [ -n "${TMP_DIR:-}" ] && [ -d "$TMP_DIR" ] && rm -rf "$TMP_DIR"
}
trap cleanup EXIT

add_to_path() {
    local dir="$1"
    local shell_profile=""
    local profile_updated=false

    # Detect shell profile
    case "${SHELL:-}" in
        */zsh) shell_profile="${ZDOTDIR:-$HOME}/.zshrc" ;;
        */bash)
            if [[ "$(uname -s)" == "Darwin" ]]; then
                shell_profile="$HOME/.zshrc"
            else
                shell_profile="$HOME/.bashrc"
            fi
            ;;
    esac

    # Already in PATH for this session?
    case ":$PATH:" in
        *":$dir:"*) return 0 ;;
    esac

    export PATH="$dir:$PATH"

    if [ -n "$shell_profile" ] && [ -w "$shell_profile" ]; then
        if ! grep -qF "$dir" "$shell_profile" 2>/dev/null; then
            printf '\nexport PATH="%s:$PATH"\n' "$dir" >> "$shell_profile"
            profile_updated=true
        fi
    fi

    if $profile_updated; then
        ok "Added to PATH in $shell_profile"
        log "Restart your terminal or run: source $shell_profile"
    else
        ok "Added to PATH for this session"
        log "Add the following to your shell profile:"
        log "  export PATH=\"$dir:\$PATH\""
    fi
}

remove_from_path() {
    local dir="$1"
    local removed=false

    # Remove from current session
    local new_path=""
    local IFS=:
    for p in $PATH; do
        [ "$p" != "$dir" ] && new_path="${new_path:+$new_path:}$p"
    done
    PATH="$new_path"

    # Remove from shell profiles
    for profile in "$HOME/.bashrc" "$HOME/.zshrc" "$HOME/.bash_profile" "$HOME/.profile"; do
        if [ -f "$profile" ]; then
            if grep -qF "$dir" "$profile" 2>/dev/null; then
                grep -vF "$dir" "$profile" > "${profile}.tmp" && mv "${profile}.tmp" "$profile"
                ok "Removed PATH entry from $profile"
                removed=true
            fi
        fi
    done

    $removed && return 0 || return 1
}

# ── Self-uninstall ──────────────────────────────────────────
if [ "${1:-}" = "--selfuninstall" ]; then
    echo ""
    echo "╭──────────────── ${PROJECT_NAME} uninstall ───────────────╮"
    echo "│"
    
    local_removed=false
    if [ -f "$BINARY_PATH" ]; then
        rm -f "$BINARY_PATH"
        ok "${PROJECT_NAME} binary removed"
        local_removed=true
    else
        skip "No binary found at $BINARY_PATH"
    fi

    if [ -d "$BIN_DIR" ]; then
        if [ -z "$(ls -A "$BIN_DIR" 2>/dev/null)" ]; then
            rmdir "$BIN_DIR" 2>/dev/null || true
            ok "Bin directory removed"
        fi
    fi

    if remove_from_path "$BIN_DIR"; then
        :  # already logged
    else
        skip "No PATH entry found"
    fi

    echo "│"
    if $local_removed; then
        echo "╰──────────────── ${PROJECT_NAME} uninstalled ───────────────╯"
    else
        echo "╰────────── ${PROJECT_NAME} not found — nothing to do ───────╯"
    fi
    echo ""
    exit 0
fi

# ── Install ─────────────────────────────────────────────────
echo ""
echo "╭──────────────── ${PROJECT_NAME} installer ────────────────╮"
echo "│"

# 1. Detect OS
log "Detecting OS..."
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$OS" in
    linux)  os="linux"   ;;
    darwin) os="darwin"  ;;
    *)      fail "Unsupported OS: $OS" ;;
esac
ok "OS: $os"

# 2. Detect architecture
log "Detecting architecture..."
ARCH="$(uname -m)"
case "$ARCH" in
    x86_64|amd64)  arch="amd64" ;;
    aarch64|arm64) arch="arm64" ;;
    *)             fail "Unsupported architecture: $ARCH" ;;
esac
ok "Architecture: $arch"

# 3. Fetch latest version
log "Fetching latest version..."
VERSION_URL="https://raw.githubusercontent.com/${GITHUB_USER}/${PROJECT_NAME}/main/.version"
VERSION="$(curl -fsSL "$VERSION_URL" 2>/dev/null | tr -d '\n\r')"
if [ -z "$VERSION" ]; then
    fail "Failed to fetch version from $VERSION_URL"
fi
ok "Latest version: $VERSION"

# 4. Build download URL
BINARY_NAME="${PROJECT_NAME}-${os}-${arch}"
DOWNLOAD_URL="https://github.com/${GITHUB_USER}/${PROJECT_NAME}/releases/download/${VERSION}/${BINARY_NAME}"
log "Download URL: $DOWNLOAD_URL"

# 5. Ensure bin directory
mkdir -p "$BIN_DIR"

# 6. Download binary
log "Downloading ${PROJECT_NAME} ${VERSION}..."
TMP_FILE="$(mktemp)"
if command -v curl &>/dev/null; then
    curl -fsSL -o "$TMP_FILE" "$DOWNLOAD_URL" 2>/dev/null || fail "Download failed (curl)"
elif command -v wget &>/dev/null; then
    wget -qO "$TMP_FILE" "$DOWNLOAD_URL" 2>/dev/null || fail "Download failed (wget)"
else
    fail "Neither curl nor wget found"
fi

if [ ! -s "$TMP_FILE" ]; then
    fail "Downloaded file is empty"
fi

# Check if it's a valid binary (not an HTML page)
if file "$TMP_FILE" | grep -qi "HTML\|html\|text"; then
    fail "Downloaded file is not a valid binary (got HTML — check the release URL)"
fi

SIZE="$(stat -c%s "$TMP_FILE" 2>/dev/null || stat -f%z "$TMP_FILE" 2>/dev/null || echo "0")"
if [ "$SIZE" -gt 1048576 ]; then
    SIZE_STR="$(echo "scale=1; $SIZE / 1048576" | bc) MB"
else
    SIZE_STR="$(echo "scale=0; $SIZE / 1024" | bc) KB"
fi
ok "Downloaded ($SIZE_STR)"

# 7. Install binary
mv "$TMP_FILE" "$BINARY_PATH"
chmod +x "$BINARY_PATH"
ok "Installed to $BINARY_PATH"

# 8. Test the binary
log "Verifying installation..."
"$BINARY_PATH" --version 2>/dev/null || fail "Binary verification failed"
ok "Verified: $($BINARY_PATH --version 2>/dev/null)"

# 9. Add to PATH
add_to_path "$BIN_DIR"

echo "│"
echo "╰──────────────── ${PROJECT_NAME} installed ─────────────────╯"
echo ""
echo "  Run '${PROJECT_NAME} --help' to get started."
echo ""
echo "  To uninstall later:"
echo "    curl -fsSL https://raw.githubusercontent.com/${GITHUB_USER}/${PROJECT_NAME}/main/installer.sh | sh -s -- --selfuninstall"
echo ""
