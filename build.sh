#!/usr/bin/env bash
set -euo pipefail

NAME="${1:-neocut}"
OUTPUT_DIR="${2:-$(dirname "$0")/bin}"

VERSION="$(cat "$(dirname "$0")/.version" | tr -d '\n\r')"
COMMIT="$(git -C "$(dirname "$0")" rev-parse --short HEAD 2>/dev/null || echo "unknown")"

PUBLISHER_NAME="rkriad585"
PUBLISHER_EMAIL="rkriad585@gmail.com"

LDFLAGS="-X 'neocut/internal/config.Commit=${COMMIT}'"
LDFLAGS="${LDFLAGS} -X 'neocut/internal/config.Version=${VERSION}'"
LDFLAGS="${LDFLAGS} -X 'neocut/internal/config.PublisherName=${PUBLISHER_NAME}'"
LDFLAGS="${LDFLAGS} -X 'neocut/internal/config.PublisherEmail=${PUBLISHER_EMAIL}'"

PLATFORMS=(
    "windows/amd64/.exe"
    "windows/arm64/.exe"
    "linux/amd64/"
    "linux/arm64/"
    "darwin/amd64/"
    "darwin/arm64/"
)

echo "  Generating embedded assets..."
go generate ./internal/config/ 2>&1

mkdir -p "$OUTPUT_DIR"

echo "╭──────────────── neocut build ────────────────╮"
printf "│  Version : %-20s│\n" "$VERSION"
printf "│  Commit  : %-20s│\n" "$COMMIT"
printf "│  Publisher: %-18s│\n" "$PUBLISHER_NAME"
printf "│  Email   : %-20s│\n" "$PUBLISHER_EMAIL"
echo "╰─────────────────────────────────────────────╯"
echo ""

count=0
total=${#PLATFORMS[@]}

for platform in "${PLATFORMS[@]}"; do
    IFS='/' read -ra parts <<< "$platform"
    os="${parts[0]}"
    arch="${parts[1]}"
    ext="${parts[2]}"

    binary="${NAME}-${os}-${arch}${ext}"
    path="${OUTPUT_DIR}/${binary}"

    printf "  [%d/%d] %s" $((count + 1)) $total "$binary"

    GOOS="$os" GOARCH="$arch" go build -ldflags "${LDFLAGS}" -o "$path" ./cmd/neocut/ 2>/tmp/neocut_build_err

    if [ $? -eq 0 ]; then
        size="$(stat -c%s "$path" 2>/dev/null || stat -f%z "$path" 2>/dev/null)"
        if [ "$size" -gt 1048576 ]; then
            size_str="$(echo "scale=1; $size / 1048576" | bc) MB"
        else
            size_str="$(echo "scale=0; $size / 1024" | bc) KB"
        fi
        echo "  ✓ $size_str"
        count=$((count + 1))
    else
        echo "  ✗ FAILED"
        cat /tmp/neocut_build_err >&2
    fi

    unset os arch ext binary path
done

echo ""
if [ "$count" -eq "$total" ]; then
    echo "  All $count binaries built successfully in $OUTPUT_DIR"
else
    echo "  $count/$total binaries built (see errors above)"
fi
echo ""
