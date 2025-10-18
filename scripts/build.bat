@echo off
setlocal

set VERSION=%1
if "%VERSION%"=="" set VERSION=1.0.0

set BUILD_TIME=%date% %time%
set LDFLAGS=-s -w -X main.version=%VERSION% -X main.buildTime=%BUILD_TIME%

echo Building Interactive Feedback MCP v%VERSION%

if not exist build mkdir build

echo Building for Windows (amd64)...
set GOOS=windows
set GOARCH=amd64
go build -ldflags="%LDFLAGS%" -o build/mcp-server-single-windows-amd64.exe cmd/mcp-server-single/main.go

echo Building for Linux (amd64)...
set GOOS=linux
set GOARCH=amd64
go build -ldflags="%LDFLAGS%" -o build/mcp-server-single-linux-amd64 cmd/mcp-server-single/main.go

echo Building for macOS (amd64)...
set GOOS=darwin
set GOARCH=amd64
go build -ldflags="%LDFLAGS%" -o build/mcp-server-single-macos-amd64 cmd/mcp-server-single/main.go

echo Build completed successfully!
echo Binaries are available in the build/ directory