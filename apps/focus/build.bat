@echo off

REM Build script for focus.sh CLI on Windows

echo Installing dependencies...
go mod tidy

echo Building focus.sh CLI application...
go build -o focus.exe ./cmd/focus

echo Build complete! Run focus.exe to start using your AI-powered task manager.