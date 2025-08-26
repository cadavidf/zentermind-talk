#!/bin/bash

# Book Selection Interface
# Allows user to select which book to expand with detailed chapters

echo "📚 BOOK EXPANSION SELECTOR"
echo "========================="
echo ""
echo "Available books for chapter expansion:"
echo ""

# Display completed books (Books 1-2 are already expanded)
echo "✅ COMPLETED (Already Expanded):"
echo "   1. Your AI, Your Mirror - Personal Transformation & Technology"
echo "   2. The Quiet Machine - Calm Technology & Mindful Design"
echo ""

echo "📖 AVAILABLE FOR EXPANSION:"
echo "   3. Conscious Code - AI Ethics & Philosophy"
echo "   4. Live Long, Live Well - Longevity & Health"
echo "   5. Preventive Everything - Preventive Medicine & AI"
echo "   6. The Regenerative Revolution - Sustainable Business & Environment"
echo "   7. Waste to Wealth - Circular Economy & Innovation"
echo "   8. Own Your Data, Own Your Life - Data Privacy & Digital Rights"
echo "   9. The Infinite Canvas - Augmented Reality & User Experience"
echo "   10. Focus as Superpower - Attention Economy & Mental Training"
echo ""
echo "   11. Custom Title (Enter your own book title)"
echo ""

read -p "Select book number (3-11): " choice

case $choice in
    3)
        BOOK_NUM=3
        BOOK_TITLE="Conscious Code"
        BOOK_CATEGORY="AI Ethics & Philosophy"
        ;;
    4)
        BOOK_NUM=4
        BOOK_TITLE="Live Long, Live Well"
        BOOK_CATEGORY="Longevity & Health"
        ;;
    5)
        BOOK_NUM=5
        BOOK_TITLE="Preventive Everything"
        BOOK_CATEGORY="Preventive Medicine & AI"
        ;;
    6)
        BOOK_NUM=6
        BOOK_TITLE="The Regenerative Revolution"
        BOOK_CATEGORY="Sustainable Business & Environment"
        ;;
    7)
        BOOK_NUM=7
        BOOK_TITLE="Waste to Wealth"
        BOOK_CATEGORY="Circular Economy & Innovation"
        ;;
    8)
        BOOK_NUM=8
        BOOK_TITLE="Own Your Data, Own Your Life"
        BOOK_CATEGORY="Data Privacy & Digital Rights"
        ;;
    9)
        BOOK_NUM=9
        BOOK_TITLE="The Infinite Canvas"
        BOOK_CATEGORY="Augmented Reality & User Experience"
        ;;
    10)
        BOOK_NUM=10
        BOOK_TITLE="Focus as Superpower"
        BOOK_CATEGORY="Attention Economy & Mental Training"
        ;;
    11)
        echo ""
        read -p "Enter custom book title: " BOOK_TITLE
        read -p "Enter book category: " BOOK_CATEGORY
        BOOK_NUM="custom"
        ;;
    *)
        echo "❌ Invalid selection. Please choose 3-11."
        exit 1
        ;;
esac

echo ""
echo "🎯 SELECTED BOOK:"
echo "   Title: $BOOK_TITLE"
echo "   Category: $BOOK_CATEGORY"
if [ "$BOOK_NUM" != "custom" ]; then
    echo "   Book Number: $BOOK_NUM"
fi
echo ""

read -p "Proceed with expanding this book? (y/n): " confirm

if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
    echo ""
    echo "🚀 Starting book expansion for: $BOOK_TITLE"
    echo ""
    
    if [ "$BOOK_NUM" != "custom" ]; then
        BOOK_DIR="output/books/book_$(printf "%03d" $BOOK_NUM)"
        
        # Check if book already has expanded chapters
        if [ -f "$BOOK_DIR/EXPANDED_BOOK_SUMMARY.md" ]; then
            echo "✅ Book already expanded! Skipping."
            echo "📄 Expansion summary: $BOOK_DIR/EXPANDED_BOOK_SUMMARY.md"
            exit 0
        fi
        
        # Check if book directory exists
        if [ -d "$BOOK_DIR" ]; then
            echo "📁 Found book directory: $BOOK_DIR"
            
            if [ -f "$BOOK_DIR/outline_1000_words.md" ]; then
                echo "✅ Outline found. Starting automatic expansion..."
                echo ""
                
                # Create single book config for enhanced generator
                cat > single_book_config.json << EOF
{
    "target_book": $BOOK_NUM,
    "book_title": "$BOOK_TITLE",
    "book_category": "$BOOK_CATEGORY",
    "mode": "expansion_only",
    "start_immediately": true
}
EOF
                
                echo "🔧 Created expansion config for book $BOOK_NUM"
                echo "🚀 Launching enhanced generator..."
                
                # Launch the enhanced generator with the specific book
                ./enhanced_sequential_generator --book=$BOOK_NUM --mode=expand
                
            else
                echo "❌ Outline not found. Book needs completion first."
                exit 1
            fi
        else
            echo "❌ Book directory not found. Book needs generation first."
            exit 1
        fi
    else
        echo "❌ Custom books not supported for automatic expansion."
        exit 1
    fi
else
    echo "❌ Book expansion cancelled."
fi