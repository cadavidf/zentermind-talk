#!/bin/bash

# BULLET BOOKS - Automated Book Generation Script
# Generates 10 books using the book generation agent

set -e  # Exit on any error

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
CYAN='\033[0;36m'
NC='\033[0m' # No Color

# Function to print colored output
print_status() {
    echo -e "${BLUE}[$(date +'%H:%M:%S')]${NC} $1"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_info() {
    echo -e "${CYAN}ℹ️  $1${NC}"
}

# Banner
echo -e "${PURPLE}"
cat << "EOF"
██████╗  ██████╗  ██████╗ ██╗  ██╗     ██████╗ ███████╗███╗   ██╗███████╗██████╗  █████╗ ████████╗██╗ ██████╗ ███╗   ██╗
██╔══██╗██╔═══██╗██╔═══██╗██║ ██╔╝    ██╔════╝ ██╔════╝████╗  ██║██╔════╝██╔══██╗██╔══██╗╚══██╔══╝██║██╔═══██╗████╗  ██║
██████╔╝██║   ██║██║   ██║█████╔╝     ██║  ███╗█████╗  ██╔██╗ ██║█████╗  ██████╔╝███████║   ██║   ██║██║   ██║██╔██╗ ██║
██╔══██╗██║   ██║██║   ██║██╔═██╗     ██║   ██║██╔══╝  ██║╚██╗██║██╔══╝  ██╔══██╗██╔══██║   ██║   ██║██║   ██║██║╚██╗██║
██████╔╝╚██████╔╝╚██████╔╝██║  ██╗    ╚██████╔╝███████╗██║ ╚████║███████╗██║  ██║██║  ██║   ██║   ██║╚██████╔╝██║ ╚████║
╚═════╝  ╚═════╝  ╚═════╝ ╚═╝  ╚═╝     ╚═════╝ ╚══════╝╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝   ╚═╝ ╚═════╝ ╚═╝  ╚═══╝

            🤖 AUTOMATED BOOK GENERATION - 10 BOOKS IN MULTIPLE FORMATS 🤖
                                    Co-authored by animality.ai
EOF
echo -e "${NC}"

# Start time
START_TIME=$(date +%s)
print_status "Starting automated book generation process..."

# Check prerequisites
print_status "Checking prerequisites..."

# Check Go installation
if ! command -v go &> /dev/null; then
    print_error "Go is not installed. Please install Go 1.21+ first."
    exit 1
fi

# Check Go version
GO_VERSION=$(go version | grep -o 'go[0-9]\+\.[0-9]\+' | sed 's/go//')
MAJOR_VERSION=$(echo $GO_VERSION | cut -d. -f1)
MINOR_VERSION=$(echo $GO_VERSION | cut -d. -f2)

if [ "$MAJOR_VERSION" -eq 1 ] && [ "$MINOR_VERSION" -lt 21 ]; then
    print_error "Go version $GO_VERSION is too old. Please upgrade to Go 1.21+"
    exit 1
fi

print_success "Go version $GO_VERSION is compatible"

# Check if we're in the right directory
if [ ! -f "book_generation_agent.go" ]; then
    print_error "Please run this script from the Dev37 directory"
    exit 1
fi

# Check if future_book_themes.md exists
if [ ! -f "future_book_themes.md" ]; then
    print_error "future_book_themes.md not found. Please ensure the themes file exists."
    exit 1
fi

print_success "Prerequisites check passed"

# Check Ollama service
print_status "Checking Ollama service..."
if ! curl -s http://localhost:11434/api/tags >/dev/null 2>&1; then
    print_warning "Ollama service is not accessible at localhost:11434"
    print_info "Please start Ollama service before running book generation"
    print_info "Run: ollama serve"
    
    read -p "Do you want to continue anyway? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        print_info "Exiting. Start Ollama service and try again."
        exit 1
    fi
else
    print_success "Ollama service is running"
fi

# Create output directories
print_status "Creating output directories..."
mkdir -p output/generated_books
mkdir -p output/progress_saves
mkdir -p output/batch_reports

print_success "Output directories created"

# Build Phase 6 Enhanced if needed
print_status "Preparing Phase 6 Enhanced..."
cd phase6_enhanced

# Install dependencies if needed
if [ ! -f "go.sum" ]; then
    print_status "Installing Go dependencies for EPUB generation..."
    go mod tidy
fi

# Build binary if it doesn't exist
if [ ! -f "phase6_enhanced" ]; then
    print_status "Building Phase 6 Enhanced..."
    go build -o phase6_enhanced main.go
    if [ $? -eq 0 ]; then
        print_success "Phase 6 Enhanced built successfully"
    else
        print_error "Failed to build Phase 6 Enhanced"
        exit 1
    fi
else
    print_success "Phase 6 Enhanced binary exists"
fi

cd ..

# Run the book generation agent
print_status "Starting book generation agent..."
echo ""
echo "════════════════════════════════════════════════════════════════════════"
print_info "About to generate 10 books in 4 formats each (.epub, .txt, .md, .json)"
print_info "This process may take 30-60 minutes depending on your system"
print_info "Progress will be displayed in real-time"
echo "════════════════════════════════════════════════════════════════════════"
echo ""

# Set environment variables for automated mode
export OLLAMA_MODEL="llama3.2"
export AUTOMATED_MODE="true"

# Run the book generation agent
if go run book_generation_agent.go; then
    # Calculate execution time
    END_TIME=$(date +%s)
    DURATION=$((END_TIME - START_TIME))
    MINUTES=$((DURATION / 60))
    SECONDS=$((DURATION % 60))
    
    echo ""
    print_success "🎉 Automated book generation completed successfully!"
    print_info "Total execution time: ${MINUTES}m ${SECONDS}s"
    
    # Check generated files
    BOOK_COUNT=$(find output/generated_books -name "*.epub" | wc -l)
    TOTAL_FILES=$(find output/generated_books -type f | wc -l)
    
    print_info "📊 Generated files summary:"
    print_info "   📚 Books: $BOOK_COUNT"
    print_info "   📁 Total files: $TOTAL_FILES"
    
    # Display file breakdown
    if [ $BOOK_COUNT -gt 0 ]; then
        EPUB_COUNT=$(find output/generated_books -name "*.epub" | wc -l)
        TXT_COUNT=$(find output/generated_books -name "*.txt" | wc -l)
        MD_COUNT=$(find output/generated_books -name "*.md" | wc -l)
        JSON_COUNT=$(find output/generated_books -name "*.json" | wc -l)
        
        print_info "   📖 EPUB files: $EPUB_COUNT"
        print_info "   📄 TXT files: $TXT_COUNT"
        print_info "   📝 MD files: $MD_COUNT"
        print_info "   📊 JSON files: $JSON_COUNT"
    fi
    
    print_info "📁 Files are available in:"
    print_info "   - output/generated_books/ (all book files)"
    print_info "   - output/batch_generation_report.json (detailed report)"
    
else
    print_error "Book generation failed. Check the output above for error details."
    exit 1
fi

echo ""
print_success "🎊 All done! Enjoy your 10 generated books! 🎊"
print_info "Co-authored by animality.ai"