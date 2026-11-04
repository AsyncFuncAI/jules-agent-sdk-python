#!/bin/bash
# Jules Agent VSCode Extension - Quick Setup Script

set -e

echo "🚀 Jules Agent VSCode Extension Setup"
echo "======================================"
echo ""

# Check Node.js
if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed. Please install Node.js 16+ first."
    exit 1
fi

echo "✅ Node.js version: $(node --version)"

# Check npm
if ! command -v npm &> /dev/null; then
    echo "❌ npm is not installed. Please install npm first."
    exit 1
fi

echo "✅ npm version: $(npm --version)"

# Check Python
if ! command -v python3 &> /dev/null; then
    echo "❌ Python 3 is not installed. Please install Python 3.8+ first."
    exit 1
fi

echo "✅ Python version: $(python3 --version)"

# Install npm dependencies
echo ""
echo "📦 Installing npm dependencies..."
npm install

# Install Python dependencies
echo ""
echo "🐍 Installing Python dependencies..."
cd python
pip install -r requirements.txt
cd ..

# Compile TypeScript
echo ""
echo "🔨 Compiling TypeScript..."
npm run compile

echo ""
echo "✅ Setup complete!"
echo ""
echo "Next steps:"
echo "1. Get your Jules API key from https://jules.google.com"
echo "2. Press F5 in VSCode to run the extension"
echo "3. Enter your API key when prompted"
echo "4. Run 'Jules: Start New Session' from Command Palette"
echo ""
echo "📚 See README.md for detailed usage instructions"
echo "🔧 See SETUP.md for development guide"
