#!/bin/sh
set -e

REPO="jhgundersen/imagegen"
BINARY="imagegen"
BASE="https://raw.githubusercontent.com/${REPO}/master"

OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case $ARCH in
  x86_64)        ARCH=amd64 ;;
  aarch64|arm64) ARCH=arm64 ;;
  *) echo "Unsupported architecture: $ARCH" >&2; exit 1 ;;
esac

BIN_DIR="${HOME}/.local/bin"
BASH_DIR="${HOME}/.local/share/bash-completion/completions"
ZSH_DIR="${HOME}/.local/share/zsh/site-functions"

echo "Downloading ${BINARY} for ${OS}/${ARCH}..."
mkdir -p "$BIN_DIR"
curl -fsSL "https://github.com/${REPO}/releases/latest/download/${BINARY}-${OS}-${ARCH}" \
  -o "${BIN_DIR}/${BINARY}"
chmod +x "${BIN_DIR}/${BINARY}"

echo "Installing shell completions..."
mkdir -p "$BASH_DIR" "$ZSH_DIR"
curl -fsSL "${BASE}/completions/${BINARY}.bash" -o "${BASH_DIR}/${BINARY}"
curl -fsSL "${BASE}/completions/${BINARY}.zsh"  -o "${ZSH_DIR}/_${BINARY}"

echo ""
echo "Installed to ${BIN_DIR}/${BINARY}"
echo ""
echo "To enable completions:"
echo "  bash: restart your shell (completions load automatically if ~/.local/share/bash-completion is in your setup)"
echo "  zsh:  add to ~/.zshrc if not already present:"
echo "          fpath=(~/.local/share/zsh/site-functions \$fpath)"
echo "          autoload -Uz compinit && compinit"
