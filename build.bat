@echo off

REM Build script for NEON CLI on Windows

echo Installing dependencies...
go mod tidy

echo Building NEON CLI application...
go build -o neon.exe ./cmd/neon

echo Build complete! Run neon.exe to start using your AI-powered task manager.