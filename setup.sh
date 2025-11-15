#!/bin/bash

# Blitz Quick Setup Script
# This script helps you get started with Blitz quickly

set -e

echo "=================================="
echo "  ⚡ Blitz Quick Setup ⚡"
echo "=================================="
echo ""

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed!"
    echo "Please install Go 1.21 or higher from https://golang.org/dl/"
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
echo "✅ Go version: $GO_VERSION"

# Check if git is installed
if ! command -v git &> /dev/null; then
    echo "❌ Git is not installed!"
    echo "Please install Git from https://git-scm.com/"
    exit 1
fi

echo "✅ Git is installed"
echo ""

# Download dependencies
echo "📦 Downloading dependencies..."
go mod download
go mod verify
echo "✅ Dependencies downloaded"
echo ""

# Build the project
echo "🔨 Building Blitz..."
go build -ldflags="-s -w" -o blitz .
echo "✅ Build complete!"
echo ""

# Make executable (Unix-like systems)
if [[ "$OSTYPE" != "msys" && "$OSTYPE" != "win32" ]]; then
    chmod +x blitz
fi

# Check if wordlists exist
if [ ! -f "usernames.txt" ] || [ ! -f "passwords.txt" ]; then
    echo "⚠️  Warning: Default wordlists not found"
    echo "Make sure usernames.txt and passwords.txt are in the current directory"
else
    echo "✅ Wordlists found"
fi

echo ""
echo "=================================="
echo "  ⚡ Setup Complete! ⚡"
echo "=================================="
echo ""
echo "You can now run Blitz:"
echo "  ./blitz -url http://example.com/login"
echo ""
echo "For help:"
echo "  ./blitz --help"
echo ""
echo "For examples:"
echo "  cat EXAMPLES.md"
echo ""
echo "⚠️  REMEMBER: Only test systems you have permission to test!"
echo ""
