#!/bin/sh
set -e

REPO="XenophonLXH/xodo"
BINARY="xodo"
INSTALL_DIR="/usr/local/bin"

# Detect OS and arch
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)
case "$ARCH" in
    x86_64) ARCH="amd64" ;;
    arm64|aarch64) ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH"; exit 1 ;;
esac

# Get latest release tag from GitHub
VERSION=$(curl -sSf "https://api.github.com/repos/$REPO/releases/latest" \
    | grep '"tag_name"' | sed -E 's/.*"([^"]+)".*/\1/')

echo "Installing $BINARY $VERSION..."

ARCHIVE="${BINARY}_${OS}_${ARCH}.tar.gz"
URL="https://github.com/$REPO/releases/download/$VERSION/$ARCHIVE"
echo "URL: $URL"
curl -sSfL "$URL" -o "/tmp/$ARCHIVE"
tar -xzf "/tmp/$ARCHIVE" -C "/tmp" "$BINARY"
chmod +x "/tmp/$BINARY"
sudo mv "/tmp/$BINARY" "$INSTALL_DIR/$BINARY"
rm "/tmp/$ARCHIVE"

echo "Done! Run: $BINARY"
