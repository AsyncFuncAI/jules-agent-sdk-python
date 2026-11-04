#!/bin/bash
# Build VSIX package for distribution

set -e

echo "📦 Building Jules Agent VSCode Extension"
echo "========================================"
echo ""

# Check if vsce is installed
if ! command -v vsce &> /dev/null; then
    echo "⚠️  vsce not found. Installing..."
    npm install -g @vscode/vsce
fi

# Clean previous builds
echo "🧹 Cleaning previous builds..."
rm -rf dist/
rm -f *.vsix

# Install dependencies
echo "📦 Installing dependencies..."
npm install

# Compile
echo "🔨 Compiling TypeScript..."
npm run compile

# Run tests
echo "🧪 Running tests..."
npm test || echo "⚠️  Tests failed, continuing anyway..."

# Lint
echo "🔍 Linting code..."
npm run lint || echo "⚠️  Linting issues found, continuing anyway..."

# Package
echo "📦 Creating VSIX package..."
vsce package

echo ""
echo "✅ Build complete!"
echo ""
ls -lh *.vsix
echo ""
echo "To install:"
echo "  code --install-extension jules-agent-*.vsix"
echo ""
echo "To publish:"
echo "  vsce publish"
