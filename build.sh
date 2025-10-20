#!/bin/bash

# Build script for AI Focus application

echo "Installing dependencies..."
go mod tidy

echo "Building application..."
go build -o focus cmd/focus/main.go

echo "Build complete! Run ./focus to start using your AI-powered focus app."