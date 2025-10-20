#!/bin/bash

# Build script for AI Todo application

echo "Installing dependencies..."
go mod tidy

echo "Building application..."
go build -o todo.exe cmd/todo/main.go

echo "Build complete! Run ./todo.exe to start using your AI-powered todo app."