#!/bin/bash

# Install script for Interactive Feedback MCP Go - macOS
# This script installs the MCP server to /usr/local/bin/

set -e

echo "🚀 Installing Interactive Feedback MCP Go for macOS..."

# Check if running as root
if [ "$EUID" -ne 0 ]; then
    echo "❌ This script needs to be run as root (use sudo)"
    echo "Usage: sudo ./install.sh"
    exit 1
fi

# Check if binary exists
if [ ! -f "./interactive-feedback-mcp" ]; then
    echo "❌ Binary not found! Make sure you're in the extracted package directory."
    exit 1
fi

# Install binary
echo "📦 Installing binary to /usr/local/bin/..."
cp interactive-feedback-mcp /usr/local/bin/interactive-feedback-mcp
chmod +x /usr/local/bin/interactive-feedback-mcp

# Install Python GUI script
echo "📦 Installing Python GUI script..."
cp desktop_gui_single.py /usr/local/bin/interactive-feedback-gui.py
chmod +x /usr/local/bin/interactive-feedback-gui.py

# Verify installation
if command -v interactive-feedback-mcp >/dev/null 2>&1; then
    echo "✅ Installation successful!"
    echo ""
    echo "📋 Usage:"
    echo "  interactive-feedback-mcp"
    echo ""
    echo "📝 Add to your Cursor MCP configuration:"
    echo '  {'
    echo '    "mcpServers": {'
    echo '      "interactive-feedback": {'
    echo '        "command": "interactive-feedback-mcp",'
    echo '        "args": []'
    echo '      }'
    echo '    }'
    echo '  }'
    echo ""
    echo "🎉 Ready to use!"
else
    echo "❌ Installation failed!"
    exit 1
fi
