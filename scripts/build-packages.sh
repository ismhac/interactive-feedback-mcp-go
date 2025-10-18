#!/bin/bash

# Build packages script for Interactive Feedback MCP Go
# Creates .zip and .tar.gz packages for all platforms

set -e

PROJECT_NAME="interactive-feedback-mcp"
VERSION="1.0.0"
BUILD_DIR="build"
PACKAGES_DIR="packages"

echo "🚀 Building packages for Interactive Feedback MCP Go v${VERSION}"

# Clean and create directories
rm -rf ${BUILD_DIR} ${PACKAGES_DIR}
mkdir -p ${BUILD_DIR} ${PACKAGES_DIR}

# Build for different platforms
echo "📦 Building binaries for all platforms..."

# Linux AMD64
echo "Building for Linux AMD64..."
GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ${BUILD_DIR}/interactive-feedback-mcp-linux-amd64 cmd/mcp-server-single/main.go

# Windows AMD64
echo "Building for Windows AMD64..."
GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ${BUILD_DIR}/interactive-feedback-mcp-windows-amd64.exe cmd/mcp-server-single/main.go

# macOS AMD64
echo "Building for macOS AMD64..."
GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o ${BUILD_DIR}/interactive-feedback-mcp-darwin-amd64 cmd/mcp-server-single/main.go

# macOS ARM64 (Apple Silicon)
echo "Building for macOS ARM64..."
GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o ${BUILD_DIR}/interactive-feedback-mcp-darwin-arm64 cmd/mcp-server-single/main.go

echo "📦 Creating packages..."

# Create Linux package
echo "Creating Linux package..."
mkdir -p ${PACKAGES_DIR}/interactive-feedback-mcp-linux-amd64
cp ${BUILD_DIR}/interactive-feedback-mcp-linux-amd64 ${PACKAGES_DIR}/interactive-feedback-mcp-linux-amd64/
cp desktop_gui_single.py ${PACKAGES_DIR}/interactive-feedback-mcp-linux-amd64/
cp scripts/install-linux.sh ${PACKAGES_DIR}/interactive-feedback-mcp-linux-amd64/install.sh
cp scripts/README-package.txt ${PACKAGES_DIR}/interactive-feedback-mcp-linux-amd64/README.txt
chmod +x ${PACKAGES_DIR}/interactive-feedback-mcp-linux-amd64/interactive-feedback-mcp-linux-amd64
chmod +x ${PACKAGES_DIR}/interactive-feedback-mcp-linux-amd64/install.sh
cd ${PACKAGES_DIR}
tar -czf interactive-feedback-mcp-linux-amd64.tar.gz interactive-feedback-mcp-linux-amd64/
cd ..

# Create Windows package
echo "Creating Windows package..."
mkdir -p ${PACKAGES_DIR}/interactive-feedback-mcp-windows-amd64
cp ${BUILD_DIR}/interactive-feedback-mcp-windows-amd64.exe ${PACKAGES_DIR}/interactive-feedback-mcp-windows-amd64/interactive-feedback-mcp.exe
cp desktop_gui_single.py ${PACKAGES_DIR}/interactive-feedback-mcp-windows-amd64/
cp scripts/install-windows.bat ${PACKAGES_DIR}/interactive-feedback-mcp-windows-amd64/install.bat
cp scripts/README-package.txt ${PACKAGES_DIR}/interactive-feedback-mcp-windows-amd64/README.txt
cd ${PACKAGES_DIR}
zip -r interactive-feedback-mcp-windows-amd64.zip interactive-feedback-mcp-windows-amd64/
cd ..

# Create macOS AMD64 package
echo "Creating macOS AMD64 package..."
mkdir -p ${PACKAGES_DIR}/interactive-feedback-mcp-darwin-amd64
cp ${BUILD_DIR}/interactive-feedback-mcp-darwin-amd64 ${PACKAGES_DIR}/interactive-feedback-mcp-darwin-amd64/interactive-feedback-mcp
cp desktop_gui_single.py ${PACKAGES_DIR}/interactive-feedback-mcp-darwin-amd64/
cp scripts/install-macos.sh ${PACKAGES_DIR}/interactive-feedback-mcp-darwin-amd64/install.sh
cp scripts/README-package.txt ${PACKAGES_DIR}/interactive-feedback-mcp-darwin-amd64/README.txt
chmod +x ${PACKAGES_DIR}/interactive-feedback-mcp-darwin-amd64/interactive-feedback-mcp
chmod +x ${PACKAGES_DIR}/interactive-feedback-mcp-darwin-amd64/install.sh
cd ${PACKAGES_DIR}
tar -czf interactive-feedback-mcp-darwin-amd64.tar.gz interactive-feedback-mcp-darwin-amd64/
cd ..

# Create macOS ARM64 package
echo "Creating macOS ARM64 package..."
mkdir -p ${PACKAGES_DIR}/interactive-feedback-mcp-darwin-arm64
cp ${BUILD_DIR}/interactive-feedback-mcp-darwin-arm64 ${PACKAGES_DIR}/interactive-feedback-mcp-darwin-arm64/interactive-feedback-mcp
cp desktop_gui_single.py ${PACKAGES_DIR}/interactive-feedback-mcp-darwin-arm64/
cp scripts/install-macos.sh ${PACKAGES_DIR}/interactive-feedback-mcp-darwin-arm64/install.sh
cp scripts/README-package.txt ${PACKAGES_DIR}/interactive-feedback-mcp-darwin-arm64/README.txt
chmod +x ${PACKAGES_DIR}/interactive-feedback-mcp-darwin-arm64/interactive-feedback-mcp
chmod +x ${PACKAGES_DIR}/interactive-feedback-mcp-darwin-arm64/install.sh
cd ${PACKAGES_DIR}
tar -czf interactive-feedback-mcp-darwin-arm64.tar.gz interactive-feedback-mcp-darwin-arm64/
cd ..

echo "✅ Package creation completed successfully!"
echo "📁 Packages are in the ${PACKAGES_DIR}/ directory"
echo ""
echo "📋 Available packages:"
ls -la ${PACKAGES_DIR}/*.tar.gz ${PACKAGES_DIR}/*.zip 2>/dev/null || true
echo ""
echo "🚀 Ready for GitHub Releases!"
