# Phase 6: Content Generation

## Overview
Phase 6 is the final and most comprehensive phase of the BULLET BOOKS system, responsible for generating complete book content from approved concepts.

## Features
- **3-Phase Content Generation**:
  - Phase 1: Detailed chapter outlines (1,100+ words each)
  - Phase 2: Full chapter generation from outlines  
  - Phase 3: Chapter polishing (2,000+ words each)
- **Quality Assurance**: Recursive improvement until word count targets are met
- **Engagement Optimization**: Reader retention analysis and optimization
- **Quotable Moments**: Extraction and integration of shareable quotes
- **Multi-Format Output**: TXT, Markdown, and HTML files
- **Progress Tracking**: Real-time progress updates and timestamps

## Requirements
- Go 1.21 or higher
- Ollama running on localhost:11434
- Available LLM model (e.g., llama3.1, mistral, codellama)

## Installation
```bash
cd phase6
go build -o phase6 phase6.go
```

## Usage
```bash
./phase6
```

### Input Required
1. **Model Selection**: Choose your Ollama model
2. **Book Title**: Enter the title of your book
3. **Book Description**: Provide a detailed description of the book concept
4. **Outline Approval**: Review and approve the generated chapter outline

### Output
The system generates three file formats in the `generated_books/` directory:
- `.txt` - Plain text format
- `.md` - Markdown format with proper formatting
- `.html` - Styled HTML format for web viewing

### File Naming
Files are saved with the pattern: `[Book_Title]_[YYYY-MM-DD_HH-MM-SS].[ext]`

## Process Flow
1. **Chapter Outline Creation**: Generate comprehensive chapter structure
2. **Engagement Analysis**: Analyze and optimize chapter engagement scores
3. **Content Generation**: 
   - Generate detailed outlines (1,100+ words)
   - Transform outlines into full chapters
   - Polish chapters to 2,000+ words
4. **Retention Optimization**: Optimize content for reader retention
5. **Quote Integration**: Extract and integrate quotable moments
6. **File Generation**: Save in multiple formats with clickable links

## Target Metrics
- **Total Words**: 11,100+ words minimum
- **Chapter Outlines**: 1,100+ words each
- **Final Chapters**: 2,000+ words each
- **Engagement**: 65%+ threshold
- **Retention**: 80%+ threshold
- **Quote Shareability**: 75%+ threshold

## Example Output Structure
```
generated_books/
├── My_Book_Title_2024-07-09_14-30-45.txt
├── My_Book_Title_2024-07-09_14-30-45.md
└── My_Book_Title_2024-07-09_14-30-45.html
```

## Error Handling
- Automatic retry logic for content generation
- Graceful degradation if targets aren't met after 3 attempts
- Progress saving for recovery from interruptions

## Dependencies
- Standard Go libraries only
- Ollama API for LLM interactions
- File system for output generation

## Notes
- Requires active Ollama service
- Processing time depends on model size and complexity
- Generated content is optimized for publication readiness
- All files include clickable links for easy access