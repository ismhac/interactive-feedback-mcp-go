Interactive Feedback MCP Go - Package Installation
==================================================

This package contains the Interactive Feedback MCP Go binary and installation scripts.

📦 What's included:
- interactive-feedback-mcp (binary)
- desktop_gui_single.py (Python GUI)
- install.sh / install.bat (installer script)
- README.txt (this file)

🚀 Quick Installation:

Linux/macOS:
1. Extract this package: tar -xzf interactive-feedback-mcp-*.tar.gz
2. Run installer: sudo ./install.sh
3. Done! Use 'interactive-feedback-mcp' in your MCP config

Windows:
1. Extract this package: unzip interactive-feedback-mcp-*.zip
2. Right-click install.bat → "Run as administrator"
3. Done! Use 'interactive-feedback-mcp' in your MCP config

📝 Cursor MCP Configuration:
Add this to your Cursor MCP configuration file:

{
  "mcpServers": {
    "interactive-feedback": {
      "command": "interactive-feedback-mcp",
      "args": []
    }
  }
}

🎯 Features:
- Interactive Feedback UI with conversation history
- Markdown copy support
- Cross-platform support
- High performance (12-25MB memory usage)
- Single binary, no external dependencies

📚 More Information:
- GitHub: https://github.com/your-org/interactive-feedback-mcp-go
- Documentation: See README.md in the repository
- Issues: Create an issue on GitHub

🎉 Enjoy using Interactive Feedback MCP Go!
