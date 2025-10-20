#!/bin/bash

# Setup script for AI Todo application

echo "🔧 Setting up AI Todo Assistant..."

# Check if Go is installed
if ! command -v go &> /dev/null
then
    echo "❌ Go is not installed. Please install Go 1.21+ from https://golang.org/dl/"
    exit 1
fi

echo "✅ Go is installed"

# Check if Ollama is installed
if ! command -v ollama &> /dev/null
then
    echo "❌ Ollama is not installed. Please install Ollama from https://ollama.ai/"
    exit 1
fi

echo "✅ Ollama is installed"

# Pull required model
echo "📥 Pulling llama3 model..."
ollama pull llama3

# Install Go dependencies
echo "📦 Installing Go dependencies..."
go mod tidy

echo "✅ Setup complete!"
echo ""
echo "To build the application, run:"
echo "  go build -o todo cmd/todo/main.go"
echo ""
echo "To run directly, use:"
echo "  go run cmd/todo/main.go"