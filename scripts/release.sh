#!/bin/bash
# Script to create release builds for Windows and Linux
# Creates archives with executables and plugins

set -e

EXE_NAME="likhis"
RELEASE_DIR="release"
VERSION="$1"

if [ -z "$VERSION" ]; then
    echo "Usage: ./release.sh [version]"
    echo "Example: ./release.sh v1.0.0"
    exit 1
fi

echo "========================================"
echo "Creating Release: $VERSION"
echo "========================================"
echo ""

# Create release directory
if [ -d "$RELEASE_DIR" ]; then
    echo "Cleaning existing release directory..."
    rm -rf "$RELEASE_DIR"
fi
mkdir -p "$RELEASE_DIR"
echo ""

# Build Windows amd64
echo "[1/4] Building Windows amd64..."
WINDOWS_DIR="$RELEASE_DIR/likhis-windows-amd64"
mkdir -p "$WINDOWS_DIR"
GOOS=windows GOARCH=amd64 go build -o "$WINDOWS_DIR/$EXE_NAME.exe" main.go
if [ $? -ne 0 ]; then
    echo "Error: Windows build failed!"
    exit 1
fi
echo "  ✓ Windows executable built"
echo ""

# Build Linux amd64
echo "[2/4] Building Linux amd64..."
LINUX_DIR="$RELEASE_DIR/likhis-linux-amd64"
mkdir -p "$LINUX_DIR"
GOOS=linux GOARCH=amd64 go build -o "$LINUX_DIR/$EXE_NAME" main.go
if [ $? -ne 0 ]; then
    echo "Error: Linux build failed!"
    exit 1
fi
echo "  ✓ Linux executable built"
echo ""

# Copy plugins to Windows build
echo "[3/4] Copying plugins..."
cp -r plugins "$WINDOWS_DIR/"
if [ ! -d "$WINDOWS_DIR/plugins" ]; then
    echo "Error: Failed to copy plugins to Windows build!"
    exit 1
fi
echo "  ✓ Plugins copied to Windows build"

# Copy plugins to Linux build
cp -r plugins "$LINUX_DIR/"
if [ ! -d "$LINUX_DIR/plugins" ]; then
    echo "Error: Failed to copy plugins to Linux build!"
    exit 1
fi
echo "  ✓ Plugins copied to Linux build"
echo ""

# Create archives
echo "[4/4] Creating archives..."

# Create Windows ZIP
echo "  Creating Windows ZIP archive..."
cd "$RELEASE_DIR"
zip -r "likhis-windows-amd64-$VERSION.zip" "likhis-windows-amd64/" > /dev/null
cd ..
echo "  ✓ Windows ZIP created: likhis-windows-amd64-$VERSION.zip"

# Create Linux tar.gz
echo "  Creating Linux tar.gz archive..."
cd "$RELEASE_DIR"
tar -czf "likhis-linux-amd64-$VERSION.tar.gz" "likhis-linux-amd64/"
cd ..
echo "  ✓ Linux tar.gz created: likhis-linux-amd64-$VERSION.tar.gz"
echo ""

# Cleanup build directories (optional - comment out if you want to keep them)
echo "Cleaning up build directories..."
rm -rf "$WINDOWS_DIR"
rm -rf "$LINUX_DIR"
echo "  ✓ Build directories cleaned"
echo ""

echo "========================================"
echo "Release $VERSION created successfully!"
echo "========================================"
echo ""
echo "Release files:"
echo "  - $RELEASE_DIR/likhis-windows-amd64-$VERSION.zip"
echo "  - $RELEASE_DIR/likhis-linux-amd64-$VERSION.tar.gz"
echo ""

