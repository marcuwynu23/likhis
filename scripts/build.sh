#!/bin/bash
# Script to build likhis to build folder

EXE_NAME="likhis"
BUILD_DIR="build"
SOURCE_FILE="main.go"

echo "Building $EXE_NAME..."
echo ""

# Create build directory if it doesn't exist
if [ ! -d "$BUILD_DIR" ]; then
    echo "Creating build directory: $BUILD_DIR"
    mkdir -p "$BUILD_DIR"
fi

# Build the executable
echo "Building executable..."
go build -o "$BUILD_DIR/$EXE_NAME" "$SOURCE_FILE"

if [ $? -eq 0 ]; then
    echo ""
    echo "Success! Executable built to: $BUILD_DIR/$EXE_NAME"
    echo ""
    echo "You can now run the link script to create a symbolic link:"
    echo "  scripts/link.sh"
else
    echo ""
    echo "Error: Build failed!"
    exit 1
fi

echo ""

