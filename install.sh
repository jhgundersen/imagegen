#!/bin/sh
set -e

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case $ARCH in
  x86_64)  ARCH=amd64 ;;
  aarch64|arm64) ARCH=arm64 ;;
  *) echo "Unsupported architecture: $ARCH" >&2; exit 1 ;;
esac

URL="https://github.com/jhgundersen/imagegen/releases/latest/download/imagegen-${OS}-${ARCH}"
DEST="${HOME}/.local/bin/imagegen"

mkdir -p "$(dirname "$DEST")"
echo "Downloading imagegen for ${OS}/${ARCH}..."
curl -fsSL "$URL" -o "$DEST"
chmod +x "$DEST"
echo "Installed to $DEST"
