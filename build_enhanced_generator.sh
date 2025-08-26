#!/bin/bash

# Enhanced Sequential Book Generator Build Script
echo "🔨 Building Enhanced Sequential Book Generator..."

# Set executable name
BINARY_NAME="enhanced_sequential_generator"

# Build the enhanced generator
echo "📦 Compiling Go binary..."
go build -o "$BINARY_NAME" enhanced_sequential_generator.go progress_tracker.go outline_generator.go

# Check if build succeeded
if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo "📋 Binary created: $BINARY_NAME"
    echo ""
    echo "🚀 Usage:"
    echo "   ./$BINARY_NAME                    # Start fresh generation"
    echo "   kill -USR1 \$(pgrep $BINARY_NAME) # Pause generation"
    echo "   kill -USR2 \$(pgrep $BINARY_NAME) # Resume generation"
    echo "   kill -INT \$(pgrep $BINARY_NAME)  # Graceful stop"
    echo ""
    echo "💡 Interactive commands while running:"
    echo "   pause   - Pause generation"
    echo "   resume  - Resume generation"
    echo "   stop    - Graceful stop"
    echo "   status  - Show current status"
    echo "   save    - Manual state save"
    echo "   help    - Show available commands"
    echo ""
    
    # Make binary executable
    chmod +x "$BINARY_NAME"
    
else
    echo "❌ Build failed!"
    exit 1
fi