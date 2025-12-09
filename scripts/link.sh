#!/bin/bash
# Script to create symbolic link for likhis to /usr/local/bin

EXE_NAME="likhis"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
BUILD_DIR="$SCRIPT_DIR/../build"
SOURCE_PATH="$BUILD_DIR/$EXE_NAME"
TARGET_DIR="/usr/local/bin"
TARGET_PATH="$TARGET_DIR/$EXE_NAME"

echo "Creating symbolic link for $EXE_NAME..."
echo ""

# Check if source file exists
if [ ! -f "$SOURCE_PATH" ]; then
    echo "Error: $EXE_NAME not found in $BUILD_DIR"
    echo "Please build the executable first: scripts/build.sh"
    exit 1
fi

# Create target directory if it doesn't exist
if [ ! -d "$TARGET_DIR" ]; then
    echo "Creating directory: $TARGET_DIR"
    sudo mkdir -p "$TARGET_DIR"
fi

# Remove existing link if it exists
if [ -L "$TARGET_PATH" ] || [ -f "$TARGET_PATH" ]; then
    echo "Removing existing link..."
    sudo rm -f "$TARGET_PATH"
fi

# Create symbolic link
echo "Creating symbolic link..."
echo "  Source: $SOURCE_PATH"
echo "  Target: $TARGET_PATH"
sudo ln -s "$SOURCE_PATH" "$TARGET_PATH"

if [ $? -eq 0 ]; then
    echo ""
    echo "Success! Symbolic link created."
    echo "You can now run: $TARGET_PATH"
    echo "Or simply: likhis"
else
    echo ""
    echo "Error: Failed to create symbolic link."
    echo "Make sure you have sudo privileges."
    exit 1
fi

echo ""

