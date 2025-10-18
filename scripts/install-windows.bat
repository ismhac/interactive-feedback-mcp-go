@echo off
REM Install script for Interactive Feedback MCP Go - Windows
REM This script installs the MCP server to Program Files

echo 🚀 Installing Interactive Feedback MCP Go for Windows...

REM Check if running as administrator
net session >nul 2>&1
if %errorLevel% neq 0 (
    echo ❌ This script needs to be run as administrator
    echo Right-click and select "Run as administrator"
    pause
    exit /b 1
)

REM Check if binary exists
if not exist "interactive-feedback-mcp.exe" (
    echo ❌ Binary not found! Make sure you're in the extracted package directory.
    pause
    exit /b 1
)

REM Create installation directory
echo 📦 Creating installation directory...
if not exist "C:\Program Files\InteractiveFeedbackMCP" mkdir "C:\Program Files\InteractiveFeedbackMCP"

REM Install binary
echo 📦 Installing binary...
copy "interactive-feedback-mcp.exe" "C:\Program Files\InteractiveFeedbackMCP\"
copy "desktop_gui_single.py" "C:\Program Files\InteractiveFeedbackMCP\"

REM Add to PATH
echo 📦 Adding to system PATH...
setx PATH "%PATH%;C:\Program Files\InteractiveFeedbackMCP" /M

REM Verify installation
echo ✅ Installation successful!
echo.
echo 📋 Usage:
echo   interactive-feedback-mcp
echo.
echo 📝 Add to your Cursor MCP configuration:
echo   {
echo     "mcpServers": {
echo       "interactive-feedback": {
echo         "command": "interactive-feedback-mcp",
echo         "args": []
echo       }
echo     }
echo   }
echo.
echo 🎉 Ready to use!
echo.
echo ⚠️  You may need to restart your terminal/IDE for PATH changes to take effect.
pause
