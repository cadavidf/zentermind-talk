#!/bin/bash

# BULLET BOOKS - Complete Pipeline Runner
# Executes all phases in sequence with comprehensive monitoring

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
██████╗ ██╗   ██╗██╗     ██╗     ███████╗████████╗    ██████╗  ██████╗  ██████╗ ██╗  ██╗███████╗
██╔══██╗██║   ██║██║     ██║     ██╔════╝╚══██╔══╝    ██╔══██╗██╔═══██╗██╔═══██╗██║ ██╔╝██╔════╝
██████╔╝██║   ██║██║     ██║     █████╗     ██║       ██████╔╝██║   ██║██║   ██║█████╔╝ ███████╗
██╔══██╗██║   ██║██║     ██║     ██╔══╝     ██║       ██╔══██╗██║   ██║██║   ██║██╔═██╗ ╚════██║
██████╔╝╚██████╔╝███████╗███████╗███████╗   ██║       ██████╔╝╚██████╔╝╚██████╔╝██║  ██╗███████║
╚═════╝  ╚═════╝ ╚══════╝╚══════╝╚══════╝   ╚═╝       ╚═════╝  ╚═════╝  ╚═════╝ ╚═╝  ╚═╝╚══════╝

                        🚀 COMPLETE PIPELINE RUNNER 🚀
EOF
echo -e "${NC}"

# Start time
START_TIME=$(date +%s)
print_status "Starting complete pipeline execution..."

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
if [ ! -f "run_all_phases.go" ]; then
    print_error "Please run this script from the Dev37 directory"
    exit 1
fi

print_success "Prerequisites check passed"

# Option to run Go orchestrator or bash version
echo ""
echo "Choose execution method:"
echo "1) Go Orchestrator (recommended - better progress tracking)"
echo "2) Bash Sequential (simple - runs each phase directly)"
echo "3) Quick Test (Phase 1 Beta + Phase 2 only)"
echo ""
read -p "Enter your choice (1-3): " choice

case $choice in
    1)
        print_status "Running Go orchestrator..."
        if [ ! -f "run_all_phases" ]; then
            print_status "Building pipeline orchestrator..."
            go build -o run_all_phases run_all_phases.go
        fi
        ./run_all_phases
        ;;
    2)
        print_status "Running bash sequential execution..."
        
        # Define phases
        declare -a phases=("phase1_beta" "phase2" "phase3" "phase4" "phase5" "phase6_enhanced" "phase7")
        declare -a descriptions=("Market Intelligence & USP Optimization" "Concept Generation & Validation" "Reader Feedback & Shareability" "Media Coverage & PR Analysis" "Title Optimization & A/B Testing" "Complete Content Generation" "Marketing Assets & Campaign")
        
        successful_phases=0
        failed_phases=0
        
        for i in "${!phases[@]}"; do
            phase="${phases[$i]}"
            description="${descriptions[$i]}"
            phase_num=$((i + 1))
            
            echo ""
            echo "════════════════════════════════════════════════════════════════════════"
            print_status "PHASE $phase_num: $phase"
            print_info "$description"
            echo "════════════════════════════════════════════════════════════════════════"
            
            if [ -d "$phase" ]; then
                cd "$phase"
                
                # Build if binary doesn't exist
                if [ ! -f "$phase" ]; then
                    print_status "Building $phase..."
                    if go build -o "$phase" main.go; then
                        print_success "Build successful"
                    else
                        print_error "Build failed for $phase"
                        ((failed_phases++))
                        cd ..
                        continue
                    fi
                fi
                
                # Run the phase
                print_status "Executing $phase..."
                if ./"$phase"; then
                    print_success "$phase completed successfully"
                    ((successful_phases++))
                else
                    print_error "$phase execution failed"
                    ((failed_phases++))
                fi
                
                cd ..
            else
                print_warning "Directory $phase not found, skipping..."
                ((failed_phases++))
            fi
            
            sleep 1
        done
        
        # Summary
        END_TIME=$(date +%s)
        DURATION=$((END_TIME - START_TIME))
        echo ""
        echo "🏁 PIPELINE EXECUTION SUMMARY"
        echo "════════════════════════════════════════════════════════════════════════"
        print_info "Total execution time: ${DURATION}s"
        print_success "Successful phases: $successful_phases"
        if [ $failed_phases -gt 0 ]; then
            print_error "Failed phases: $failed_phases"
        else
            print_success "Failed phases: 0"
        fi
        ;;
    3)
        print_status "Running quick test (Phase 1 Beta + Phase 2)..."
        
        # Phase 1 Beta
        echo ""
        echo "════════════════════════════════════════════════════════════════════════"
        print_status "PHASE 1: Phase 1 Beta"
        print_info "Market Intelligence & USP Optimization"
        echo "════════════════════════════════════════════════════════════════════════"
        
        if [ -d "phase1_beta" ]; then
            cd "phase1_beta"
            if [ ! -f "phase1_beta" ]; then
                print_status "Building phase1_beta..."
                go build -o phase1_beta main.go
            fi
            print_status "Executing phase1_beta..."
            ./phase1_beta
            cd ..
            print_success "Phase 1 Beta completed"
        else
            print_error "phase1_beta directory not found"
            exit 1
        fi
        
        # Phase 2
        echo ""
        echo "════════════════════════════════════════════════════════════════════════"
        print_status "PHASE 2: Phase 2"
        print_info "Concept Generation & Validation"
        echo "════════════════════════════════════════════════════════════════════════"
        
        if [ -d "phase2" ]; then
            cd "phase2"
            if [ ! -f "phase2" ]; then
                print_status "Building phase2..."
                go build -o phase2 main.go
            fi
            print_status "Executing phase2..."
            ./phase2
            cd ..
            print_success "Phase 2 completed"
        else
            print_error "phase2 directory not found"
            exit 1
        fi
        
        print_success "Quick test completed successfully!"
        ;;
    *)
        print_error "Invalid choice. Please run again and select 1, 2, or 3."
        exit 1
        ;;
esac

# Final summary
END_TIME=$(date +%s)
DURATION=$((END_TIME - START_TIME))

echo ""
print_success "Pipeline execution completed in ${DURATION}s"
print_info "Check individual phase directories for generated output files"
print_info "For detailed logs, see the output above"

# Check for key output files
echo ""
print_status "Checking for key output files..."

if [ -f "phase1_beta/usp_optimization.json" ]; then
    print_success "Phase 1 Beta output found: usp_optimization.json"
fi

if [ -f "phase2/concepts.json" ]; then
    print_success "Phase 2 output found: concepts.json"
fi

if [ -f "generated_books/" ] && [ "$(ls -A generated_books/)" ]; then
    print_success "Generated books found in generated_books/"
fi

echo ""
print_info "🎉 Thank you for using BULLET BOOKS! 🎉"