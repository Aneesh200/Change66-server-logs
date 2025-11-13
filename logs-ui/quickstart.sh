#!/bin/bash

# Logs UI Quick Start Script
# This script helps you get the logs UI up and running quickly

set -e

echo "🚀 Logs UI Quick Start"
echo "====================="
echo ""

# Check if Node.js is installed
if ! command -v node &> /dev/null; then
    echo "❌ Node.js is not installed. Please install Node.js 18.x or later."
    exit 1
fi

echo "✓ Node.js $(node --version) detected"
echo ""

# Check if npm is installed
if ! command -v npm &> /dev/null; then
    echo "❌ npm is not installed. Please install npm."
    exit 1
fi

echo "✓ npm $(npm --version) detected"
echo ""

# Install dependencies if node_modules doesn't exist
if [ ! -d "node_modules" ]; then
    echo "📦 Installing dependencies..."
    npm install
    echo "✓ Dependencies installed"
    echo ""
else
    echo "✓ Dependencies already installed"
    echo ""
fi

# Check if .env.local exists
if [ ! -f ".env.local" ]; then
    echo "⚙️  Creating .env.local configuration file..."
    
    # Prompt for API URL
    read -p "Enter your API server URL (default: http://localhost:8080): " api_url
    api_url=${api_url:-http://localhost:8080}
    
    # Prompt for API key (optional)
    read -p "Enter your API key (press Enter to skip if not required): " api_key
    
    # Create .env.local
    echo "NEXT_PUBLIC_API_URL=$api_url" > .env.local
    
    if [ ! -z "$api_key" ]; then
        echo "NEXT_PUBLIC_API_KEY=$api_key" >> .env.local
    fi
    
    echo "✓ Configuration file created"
    echo ""
else
    echo "✓ Configuration file already exists"
    echo ""
fi

echo "🎉 Setup complete!"
echo ""
echo "To start the development server, run:"
echo "  npm run dev"
echo ""
echo "The UI will be available at http://localhost:3000"
echo ""
echo "📚 For more information:"
echo "  - README.md - Full documentation"
echo "  - SETUP_GUIDE.md - Detailed setup guide"
echo ""
echo "Happy logging! 🚀"


