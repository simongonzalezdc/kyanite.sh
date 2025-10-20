#!/bin/bash
# NEON AI Todo CLI - Deployment Test

echo "🧪 Starting NEON deployment validation..."

# Build the application
echo "🔨 Building NEON CLI..."
go build -o neon ./cmd/todo

if [ $? -ne 0 ]; then
    echo "❌ Build failed"
    exit 1
fi

echo "✅ Build successful"

# Basic functionality test
echo "🎯 Testing basic functionality..."

# Add a test task
echo "Adding test mission..."
./neon add "Test deployment functionality" > /dev/null 2>&1

if [ $? -ne 0 ]; then
    echo "❌ Add command failed"
    exit 1
fi
echo "✅ Add command works"

# List tasks
echo "Listing missions..."
./neon list > /dev/null 2>&1

if [ $? -ne 0 ]; then
    echo "❌ List command failed"
    exit 1
fi
echo "✅ List command works"

# Run unit tests
echo "🧪 Running unit tests..."
go test ./internal/ai ./internal/engine ./internal/store ./pkg/models ./pkg/utils -v > test_output.txt 2>&1

if [ $? -ne 0 ]; then
    echo "❌ Unit tests failed"
    echo "Test output:"
    cat test_output.txt
    exit 1
fi
echo "✅ All unit tests pass"

# Clean up
rm -f neon test_output.txt

echo ""
echo "🎉 NEON CLI deployment validation complete!"
echo "   All systems operational - ready for production! 🚀"