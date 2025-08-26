#!/bin/bash

# BULLET BOOKS - Sequential Multi-Phase Book Generation Pipeline
# Complete book generation using outline-first approach with all 7 phases

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
███████╗███████╗ ██████╗ ██╗   ██╗███████╗███╗   ██╗████████╗██╗ █████╗ ██╗     
██╔════╝██╔════╝██╔═══██╗██║   ██║██╔════╝████╗  ██║╚══██╔══╝██║██╔══██╗██║     
███████╗█████╗  ██║   ██║██║   ██║█████╗  ██╔██╗ ██║   ██║   ██║███████║██║     
╚════██║██╔══╝  ██║▄▄ ██║██║   ██║██╔══╝  ██║╚██╗██║   ██║   ██║██╔══██║██║     
███████║███████╗╚██████╔╝╚██████╔╝███████╗██║ ╚████║   ██║   ██║██║  ██║███████╗
╚══════╝╚══════╝ ╚══▀▀═╝  ╚═════╝ ╚══════╝╚═╝  ╚═══╝   ╚═╝   ╚═╝╚═╝  ╚═╝╚══════╝

██████╗  ██████╗  ██████╗ ██╗  ██╗     ██████╗ ███████╗███╗   ██╗███████╗██████╗  █████╗ ████████╗ ██████╗ ██████╗ 
██╔══██╗██╔═══██╗██╔═══██╗██║ ██╔╝    ██╔════╝ ██╔════╝████╗  ██║██╔════╝██╔══██╗██╔══██╗╚══██╔══╝██╔═══██╗██╔══██╗
██████╔╝██║   ██║██║   ██║█████╔╝     ██║  ███╗█████╗  ██╔██╗ ██║█████╗  ██████╔╝███████║   ██║   ██║   ██║██████╔╝
██╔══██╗██║   ██║██║   ██║██╔═██╗     ██║   ██║██╔══╝  ██║╚██╗██║██╔══╝  ██╔══██╗██╔══██║   ██║   ██║   ██║██╔══██╗
██████╔╝╚██████╔╝╚██████╔╝██║  ██╗    ╚██████╔╝███████╗██║ ╚████║███████╗██║  ██║██║  ██║   ██║   ╚██████╔╝██║  ██║
╚═════╝  ╚═════╝  ╚═════╝ ╚═╝  ╚═╝     ╚═════╝ ╚══════╝╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝    ╚═════╝ ╚═╝  ╚═╝

        🔄 SEQUENTIAL MULTI-PHASE BOOK GENERATION PIPELINE 🔄
                     Co-authored by animality.ai
EOF
echo -e "${NC}"

# Start time
START_TIME=$(date +%s)
print_status "Starting sequential multi-phase book generation..."

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
if [ ! -f "sequential_book_generator.go" ]; then
    print_error "Please run this script from the Dev37 directory"
    exit 1
fi

# Check if themes file exists
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
mkdir -p output/books
mkdir -p output/progress_tracking

print_success "Output directories created"

# Build dependencies
print_status "Building Go modules..."

# Build sequential book generator (includes all dependencies)
if ! go build -o sequential_book_generator sequential_book_generator.go outline_generator.go phase_coordinator.go progress_tracker.go; then
    print_error "Failed to build sequential book generator"
    exit 1
fi

print_success "Build successful"

# Set number of books to generate
NUM_BOOKS=3
print_status "Generating $NUM_BOOKS books..."

# Run the sequential book generator
print_status "Starting sequential book generation pipeline..."
echo ""
echo "════════════════════════════════════════════════════════════════════════"
print_info "About to generate $NUM_BOOKS books through complete 7-phase pipeline"
print_info "Each book will:"
print_info "  1. Generate 1000-word outline"
print_info "  2. Process through Phase 1-7 with progressive improvements"
print_info "  3. Save progress after each phase"
print_info "  4. Create final book in multiple formats (.epub, .txt, .md, .json)"
print_info "This process may take 45-90 minutes depending on your system"
print_info "Progress will be displayed in real-time"
echo "════════════════════════════════════════════════════════════════════════"
echo ""

# Set environment variables
export OLLAMA_MODEL="llama3.2"
export AUTOMATED_MODE="true"
export TOTAL_BOOKS="$NUM_BOOKS"

# Run the sequential generator
if ./sequential_book_generator; then
    # Calculate execution time
    END_TIME=$(date +%s)
    DURATION=$((END_TIME - START_TIME))
    HOURS=$((DURATION / 3600))
    MINUTES=$(((DURATION % 3600) / 60))
    SECONDS=$((DURATION % 60))
    
    echo ""
    print_success "🎉 Sequential book generation completed successfully!"
    print_info "Total execution time: ${HOURS}h ${MINUTES}m ${SECONDS}s"
    
    # Check generated files
    TOTAL_BOOKS_CREATED=$(find output/books -name "book_summary.json" | wc -l)
    TOTAL_FILES=$(find output/books -type f | wc -l)
    
    print_info "📊 Generation summary:"
    print_info "   📚 Books completed: $TOTAL_BOOKS_CREATED"
    print_info "   📁 Total files created: $TOTAL_FILES"
    
    # Display file breakdown by book
    if [ $TOTAL_BOOKS_CREATED -gt 0 ]; then
        print_info "📖 Books generated:"
        for book_dir in output/books/book_*; do
            if [ -d "$book_dir" ] && [ -f "$book_dir/book_summary.json" ]; then
                book_name=$(basename "$book_dir")
                book_title=$(jq -r '.title // "Unknown Title"' "$book_dir/book_summary.json" 2>/dev/null || echo "Unknown Title")
                file_count=$(find "$book_dir" -type f | wc -l)
                print_info "   $book_name: $book_title ($file_count files)"
            fi
        done
    fi
    
    print_info "📁 Output structure:"
    print_info "   - output/books/book_XXX/ (individual book directories)"
    print_info "   - Each book contains: outline, phase results, final book files"
    print_info "   - Formats: .epub, .txt, .md, .json per book"
    
    # Show sample book structure
    if [ -d "output/books/book_001" ]; then
        print_info ""
        print_info "📂 Sample book structure (book_001):"
        ls -la output/books/book_001/ | head -10 | while read line; do
            print_info "     $line"
        done
    fi
    
else
    print_error "Sequential book generation failed. Check the output above for error details."
    exit 1
fi

echo ""
print_success "🎊 Sequential multi-phase book generation complete! 🎊"
print_info "Each book was processed through:"
print_info "  ✅ 1000-word outline generation"
print_info "  ✅ Phase 1: Market Intelligence & USP Optimization"
print_info "  ✅ Phase 2: Concept Generation & Validation"
print_info "  ✅ Phase 3: Reader Feedback & Shareability"
print_info "  ✅ Phase 4: Media Coverage & PR Analysis"
print_info "  ✅ Phase 5: Title Optimization & A/B Testing"
print_info "  ✅ Phase 6: Complete Content Generation"
print_info "  ✅ Phase 7: Marketing Assets & Campaign"
print_info ""
print_info "🎉 Books are ready for publication with complete marketing packages!"
print_info "Co-authored by animality.ai"

exit 0
