#!/bin/bash

# Build script for Interactive Feedback MCP Go
# Cross-platform building for Windows, macOS, and Linux

set -e

PROJECT_NAME="interactive-feedback-mcp"
VERSION="1.0.0"
BUILD_DIR="build"

echo "🚀 Building Interactive Feedback MCP Go v${VERSION}"

# Create build directory
mkdir -p ${BUILD_DIR}

# Build for different platforms
echo "📦 Building for multiple platforms..."

# Linux AMD64
echo "Building for Linux AMD64..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ${BUILD_DIR}/mcp-server-single-linux-amd64 cmd/mcp-server-single/main.go

# Windows AMD64
echo "Building for Windows AMD64..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ${BUILD_DIR}/mcp-server-single-windows-amd64.exe cmd/mcp-server-single/main.go

# macOS AMD64
echo "Building for macOS AMD64..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o ${BUILD_DIR}/mcp-server-single-darwin-amd64 cmd/mcp-server-single/main.go

# macOS ARM64 (Apple Silicon)
echo "Building for macOS ARM64..."
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o ${BUILD_DIR}/mcp-server-single-darwin-arm64 cmd/mcp-server-single/main.go

echo "✅ Build completed successfully!"
echo "📁 Build artifacts are in the ${BUILD_DIR}/ directory"
echo ""
echo "📋 Available binaries:"
ls -la ${BUILD_DIR}/