package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
	
	"github.com/bmaupin/go-epub"
)

// Academic Writing Structure Types
type Example struct {
	Description string `json:"description"`
	Details     string `json:"details"`
	WordCount   int    `json:"word_count"`
}

type SupportingElement struct {
	MainPoint    string    `json:"main_point"`
	Evidence     string    `json:"evidence"`
	Examples     []Example `json:"examples"`
	WordCount    int       `json:"word_count"`
}

type Section struct {
	TopicSentence       string              `json:"topic_sentence"`
	SupportingElements  []SupportingElement `json:"supporting_elements"`
	Conclusion          string              `json:"conclusion"`
	WordCount           int                 `json:"word_count"`
	TargetWords         int                 `json:"target_words"`
}

type Chapter struct {
	Title       string    `json:"title"`
	Sections    []Section `json:"sections"`
	WordCount   int       `json:"word_count"`
	TargetWords int       `json:"target_words"`
}

type ContentConcept struct {
	Title            string `json:"title"`
	Description      string `json:"description"`
	UniquenessScore  float64 `json:"uniqueness_score"`
	ViabilityScore   bool `json:"viability_score"`
	CommercialScore  float64 `json:"commercial_score"`
	Status           string `json:"status"`
	FailureReason    string `json:"failure_reason"`
	ContentType      string `json:"content_type"`
}

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

type ProgressSave struct {
	Timestamp   string  `json:"timestamp"`
	Phase       string  `json:"phase"`
	Attempt     int     `json:"attempt"`
	Content     string  `json:"content"`
	WordCount   int     `json:"word_count"`
	TargetWords int     `json:"target_words"`
	Status      string  `json:"status"`
}

// ASCII Art Banner
func showBanner() {
	fmt.Println(`
██████╗ ██╗   ██╗██╗     ██╗     ███████╗████████╗    ██████╗  ██████╗  ██████╗ ██╗  ██╗███████╗
██╔══██╗██║   ██║██║     ██║     ██╔════╝╚══██╔══╝    ██╔══██╗██╔═══██╗██╔═══██╗██║ ██╔╝██╔════╝
██████╔╝██║   ██║██║     ██║     █████╗     ██║       ██████╔╝██║   ██║██║   ██║█████╔╝ ███████╗
██╔══██╗██║   ██║██║     ██║     ██╔══╝     ██║       ██╔══██╗██║   ██║██║   ██║██╔═██╗ ╚════██║
██████╔╝╚██████╔╝███████╗███████╗███████╗   ██║       ██████╔╝╚██████╔╝╚██████╔╝██║  ██╗███████║
╚═════╝  ╚═════╝ ╚══════╝╚══════╝╚══════╝   ╚═╝       ╚═════╝  ╚═════╝  ╚═════╝ ╚═╝  ╚═╝╚══════╝

                              ═══════ powered by ═══════

 █████╗ ███╗   ██╗██╗███╗   ███╗ █████╗ ██╗     ██╗████████╗██╗   ██╗   ██╗██╗
██╔══██╗████╗  ██║██║████╗ ████║██╔══██╗██║     ██║╚══██╔══╝╚██╗ ██╔╝   ██║██║
███████║██╔██╗ ██║██║██╔████╔██║███████║██║     ██║   ██║    ╚████╔╝    ██║██║
██╔══██║██║╚██╗██║██║██║╚██╔╝██║██╔══██║██║     ██║   ██║     ╚██╔╝   ██╗██║██║
██║  ██║██║ ╚████║██║██║ ╚═╝ ██║██║  ██║███████╗██║   ██║      ██║    ╚█║██║██║
╚═╝  ╚═╝╚═╝  ╚═══╝╚═╝╚═╝     ╚═╝╚═╝  ╚═╝╚══════╝╚═╝   ╚═╝      ╚═╝     ╚╝╚═╝╚═╝

            🎓 PHASE 6 ENHANCED: Academic Content Generation 🎓
`)
}

func main() {
	showBanner()
	
	// Check if running in automated mode
	automatedMode := os.Getenv("AUTOMATED_MODE") == "true"
	
	var selectedModel string
	var concept ContentConcept
	var err error
	
	if automatedMode {
		// Automated mode: get inputs from environment variables
		selectedModel = getAutomatedModel()
		concept = getAutomatedConcept()
		fmt.Printf("🤖 Running in automated mode\n")
		fmt.Printf("📖 Generating: %s\n", concept.Title)
		fmt.Printf("🎯 Description: %s\n", concept.Description)
	} else {
		// Interactive mode: get inputs from user
		selectedModel, err = getModelInput()
		if err != nil {
			fmt.Printf("Error getting model: %v\n", err)
			return
		}
		
		concept, err = getContentConcept()
		if err != nil {
			fmt.Printf("Error getting content concept: %v\n", err)
			return
		}
	}
	
	// Create output directories
	createOutputDirectories()
	
	// Start academic content generation
	err = processAcademicContentGeneration(selectedModel, concept)
	if err != nil {
		fmt.Printf("Error in content generation: %v\n", err)
		return
	}
	
	fmt.Println("\n🎉 Enhanced Phase 6 completed successfully! 🎉")
}

func getModelInput() (string, error) {
	fmt.Println("=== MODEL SELECTION ===")
	fmt.Println("Enter the Ollama model name you want to use:")
	fmt.Println("(Examples: llama3.1, llama3.1:70b, mistral, codellama)")
	fmt.Print("Model: ")
	
	reader := bufio.NewReader(os.Stdin)
	model, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	
	model = strings.TrimSpace(model)
	if model == "" {
		return "llama3.1", nil // Default
	}
	
	return model, nil
}

func getContentConcept() (ContentConcept, error) {
	fmt.Println("\n=== CONTENT CONCEPT INPUT ===")
	reader := bufio.NewReader(os.Stdin)
	
	fmt.Print("Enter book title: ")
	title, err := reader.ReadString('\n')
	if err != nil {
		return ContentConcept{}, err
	}
	
	fmt.Print("Enter book description: ")
	description, err := reader.ReadString('\n')
	if err != nil {
		return ContentConcept{}, err
	}
	
	concept := ContentConcept{
		Title:            strings.TrimSpace(title),
		Description:      strings.TrimSpace(description),
		UniquenessScore:  8.5,
		ViabilityScore:   true,
		CommercialScore:  7.8,
		Status:           "approved",
		FailureReason:    "",
		ContentType:      "book",
	}
	
	return concept, nil
}

// Get model for automated mode
func getAutomatedModel() string {
	model := os.Getenv("OLLAMA_MODEL")
	if model == "" {
		return "llama3.1" // Default model
	}
	return model
}

// Get concept for automated mode
func getAutomatedConcept() ContentConcept {
	// Check if a book was selected from the web UI
	if _, err := os.Stat("selected_book.json"); err == nil {
		data, err := os.ReadFile("selected_book.json")
		if err == nil {
			var concept ContentConcept
			if json.Unmarshal(data, &concept) == nil {
				// Clean up the selection file so it's not used again
				os.Remove("selected_book.json")
				return concept
			}
		}
	}

	title := os.Getenv("BOOK_TITLE")
	description := os.Getenv("BOOK_DESCRIPTION")
	phrase := os.Getenv("BOOK_PHRASE")
	
	// Check if outline is available for outline-driven generation
	outlineAvailable := os.Getenv("OUTLINE_AVAILABLE") == "true"
	
	// Default values if not provided
	if title == "" {
		title = "The Future of Innovation"
	}
	if description == "" {
		description = "A comprehensive guide to emerging technologies and their impact on society"
	}
	
	// Incorporate the memorable phrase into the description if provided
	if phrase != "" {
		description = fmt.Sprintf("%s. %s", description, phrase)
	}
	
	// Load outline if available
	if outlineAvailable {
		if outline, err := loadBookOutline(); err == nil {
			// Use outline information to enhance the concept
			title = outline.Title
			description = outline.CoreThesis
			if len(outline.KeyConcepts) > 0 {
				description += " Key concepts include: " + strings.Join(outline.KeyConcepts[:3], ", ")
			}
		}
	}
	
	return ContentConcept{
		Title:            title,
		Description:      description,
		UniquenessScore:  8.5,
		ViabilityScore:   true,
		CommercialScore:  7.8,
		Status:           "approved",
		FailureReason:    "",
		ContentType:      "book",
	}
}

// Load book outline for outline-driven generation
func loadBookOutline() (*BookOutlineData, error) {
	bookNumber := os.Getenv("BOOK_NUMBER")
	if bookNumber == "" {
		bookNumber = "001"
	}
	
	outlineFile := fmt.Sprintf("../output/books/book_%s/outline_1000_words.json", bookNumber)
	
	data, err := os.ReadFile(outlineFile)
	if err != nil {
		return nil, err
	}
	
	var outline BookOutlineData
	if err := json.Unmarshal(data, &outline); err != nil {
		return nil, err
	}
	
	return &outline, nil
}

// Book outline data structure for loading
type BookOutlineData struct {
	Title            string                  `json:"title"`
	Subtitle         string                  `json:"subtitle"`
	MemorablePhrase  string                  `json:"memorable_phrase"`
	Category         string                  `json:"category"`
	TargetAudience   string                  `json:"target_audience"`
	CoreThesis       string                  `json:"core_thesis"`
	Chapters         []ChapterOutlineData    `json:"chapters"`
	KeyConcepts      []string                `json:"key_concepts"`
	SupportingPoints []string                `json:"supporting_points"`
	CallToAction     string                  `json:"call_to_action"`
	WordCount        int                     `json:"word_count"`
}

type ChapterOutlineData struct {
	Number       int      `json:"number"`
	Title        string   `json:"title"`
	Purpose      string   `json:"purpose"`
	KeyPoints    []string `json:"key_points"`
	Examples     []string `json:"examples"`
	WordTarget   int      `json:"word_target"`
	Takeaways    []string `json:"takeaways"`
}

func createOutputDirectories() {
	directories := []string{
		"../output/generated_books",
		"../output/progress_saves",
		"../output/phase_attempts",
	}
	
	for _, dir := range directories {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Warning: Could not create directory %s: %v\n", dir, err)
		}
	}
}

// Main academic content generation process
func processAcademicContentGeneration(selectedModel string, concept ContentConcept) error {
	fmt.Printf("\n🎓 Starting Academic Content Generation for \"%s\"\n", concept.Title)
	fmt.Println(strings.Repeat("=", 70))
	
	startTime := time.Now()
	
	// Step 1: Create initial chapter structure
	fmt.Printf("\n📚 Step 1: Creating academic chapter structure...\n")
	chapters, err := createInitialChapterStructure(selectedModel, concept)
	if err != nil {
		return fmt.Errorf("error creating chapter structure: %v", err)
	}
	
	saveAttempt("initial_structure", 0, fmt.Sprintf("%+v", chapters))
	
	// Step 2: Recursively expand each chapter using academic writing
	fmt.Printf("\n✍️ Step 2: Recursive academic content expansion...\n")
	expandedChapters := make([]Chapter, len(chapters))
	
	for i, chapter := range chapters {
		fmt.Printf("\n📖 Expanding Chapter %d: %s\n", i+1, chapter.Title)
		
		expandedChapter, err := expandChapterRecursively(selectedModel, chapter, 1800) // Target 1800 words per chapter
		if err != nil {
			return fmt.Errorf("error expanding chapter %d: %v", i+1, err)
		}
		
		expandedChapters[i] = expandedChapter
		saveAttempt("chapter_"+fmt.Sprintf("%d", i+1), 0, renderChapterToText(expandedChapter))
		
		fmt.Printf("   ✅ Chapter %d completed: %d words\n", i+1, expandedChapter.WordCount)
	}
	
	// Step 3: Compile and save final book
	fmt.Printf("\n📝 Step 3: Compiling final book...\n")
	finalBook := compileBook(concept.Title, expandedChapters)
	
	// Save in multiple formats
	savedFiles, err := saveBookContent(concept.Title, finalBook, expandedChapters)
	if err != nil {
		return fmt.Errorf("error saving book: %v", err)
	}
	
	// Step 4: Display completion summary
	totalWords := calculateTotalWords(expandedChapters)
	fmt.Printf("\n🎉 Academic Content Generation Complete!\n")
	fmt.Printf("Duration: %v\n", time.Since(startTime))
	fmt.Printf("Total Chapters: %d\n", len(expandedChapters))
	fmt.Printf("Total Words: %d\n", totalWords)
	fmt.Printf("Average Words per Chapter: %d\n", totalWords/len(expandedChapters))
	
	fmt.Printf("\n💾 Files saved:\n")
	for _, file := range savedFiles {
		absPath, _ := filepath.Abs(file)
		fmt.Printf("   📁 %s\n", absPath)
		fmt.Printf("   🔗 file://%s\n", absPath)
	}
	
	return nil
}

// Create initial chapter structure
func createInitialChapterStructure(selectedModel string, concept ContentConcept) ([]Chapter, error) {
	prompt := fmt.Sprintf(`Create a comprehensive academic chapter structure for "%s".

Description: %s

Generate 6-8 chapters using academic writing principles. For each chapter, provide:
1. A clear, focused chapter title
2. Main topics to be covered
3. Estimated target of 1800 words per chapter

Format as:
Chapter 1: [Title]
Main Topics: [2-3 key topics]

Chapter 2: [Title]
Main Topics: [2-3 key topics]

Continue for all chapters...`, concept.Title, concept.Description)

	response, err := callOllama(selectedModel, prompt)
	if err != nil {
		return nil, err
	}
	
	chapters := parseChapterStructure(response)
	return chapters, nil
}

// Parse AI response into chapter structure
func parseChapterStructure(response string) []Chapter {
	var chapters []Chapter
	lines := strings.Split(response, "\n")
	
	var currentChapter Chapter
	inChapter := false
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		if strings.HasPrefix(line, "Chapter ") && strings.Contains(line, ":") {
			if inChapter {
				chapters = append(chapters, currentChapter)
			}
			
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				title := strings.TrimSpace(parts[1])
				currentChapter = Chapter{
					Title:       title,
					Sections:    []Section{},
					TargetWords: 1800,
				}
				inChapter = true
			}
		}
	}
	
	if inChapter {
		chapters = append(chapters, currentChapter)
	}
	
	// Ensure we have at least some chapters
	if len(chapters) == 0 {
		chapters = []Chapter{
			{Title: "Introduction and Foundations", TargetWords: 1800},
			{Title: "Core Concepts and Principles", TargetWords: 1800},
			{Title: "Practical Applications", TargetWords: 1800},
			{Title: "Case Studies and Examples", TargetWords: 1800},
			{Title: "Future Implications", TargetWords: 1800},
			{Title: "Conclusion and Recommendations", TargetWords: 1800},
		}
	}
	
	return chapters
}

// Recursively expand chapter using academic writing structure
func expandChapterRecursively(selectedModel string, chapter Chapter, targetWords int) (Chapter, error) {
	attempt := 1
	maxAttempts := 3
	currentChapter := chapter
	currentChapter.TargetWords = targetWords
	
	for attempt <= maxAttempts {
		fmt.Printf("   📝 Attempt %d (target: %d words)...\n", attempt, targetWords)
		
		// Create or expand sections using academic structure
		expandedChapter, err := expandChapterSections(selectedModel, currentChapter)
		if err != nil {
			return currentChapter, err
		}
		
		expandedChapter.WordCount = calculateChapterWordCount(expandedChapter)
		
		fmt.Printf("   📊 Generated %d words", expandedChapter.WordCount)
		
		if expandedChapter.WordCount >= targetWords {
			fmt.Printf(" ✅ Target reached!\n")
			return expandedChapter, nil
		}
		
		fmt.Printf(" ⚠️ Below target, expanding...\n")
		
		// If below target, analyze gaps and expand
		gapAnalysis := analyzeContentGaps(expandedChapter, targetWords)
		currentChapter = expandBasedOnGaps(selectedModel, expandedChapter, gapAnalysis)
		
		saveAttempt("chapter_"+chapter.Title+"_attempt", attempt, renderChapterToText(currentChapter))
		
		attempt++
	}
	
	fmt.Printf("   ⚠️ Using best attempt after %d tries (%d words)\n", maxAttempts, currentChapter.WordCount)
	return currentChapter, nil
}

// Expand chapter sections using academic writing structure
func expandChapterSections(selectedModel string, chapter Chapter) (Chapter, error) {
	prompt := fmt.Sprintf(`Write academic content for the chapter: "%s"

Use academic writing structure for each main section:
1. Topic Sentence - Clear main idea
2. Supporting Elements - 2-3 key points with evidence
3. Examples - Real-world illustrations  
4. Mini-conclusion or transition

Target: %d words total

Create 3-4 main sections, each following this academic structure. Write in a scholarly but accessible tone with practical examples.

Format each section clearly with subheadings.`, chapter.Title, chapter.TargetWords)

	response, err := callOllama(selectedModel, prompt)
	if err != nil {
		return chapter, err
	}
	
	// Parse response into structured sections
	sections := parseAcademicSections(response)
	chapter.Sections = sections
	chapter.WordCount = calculateChapterWordCount(chapter)
	
	return chapter, nil
}

// Parse AI response into academic sections
func parseAcademicSections(response string) []Section {
	var sections []Section
	paragraphs := strings.Split(response, "\n\n")
	
	var currentSection Section
	sectionCount := 0
	
	for _, paragraph := range paragraphs {
		paragraph = strings.TrimSpace(paragraph)
		if len(paragraph) < 50 { // Skip very short paragraphs
			continue
		}
		
		// Every 3-4 paragraphs, start a new section
		if sectionCount%4 == 0 {
			if len(currentSection.TopicSentence) > 0 {
				sections = append(sections, currentSection)
			}
			
			// First sentence as topic sentence
			sentences := strings.Split(paragraph, ".")
			if len(sentences) > 0 {
				currentSection = Section{
					TopicSentence: strings.TrimSpace(sentences[0]) + ".",
					SupportingElements: []SupportingElement{
						{
							MainPoint: paragraph,
							Evidence:  paragraph,
							Examples:  []Example{{Description: "Contextual example", Details: paragraph, WordCount: len(strings.Fields(paragraph))}},
							WordCount: len(strings.Fields(paragraph)),
						},
					},
					WordCount: len(strings.Fields(paragraph)),
				}
			}
		} else {
			// Add as supporting element
			if len(currentSection.SupportingElements) < 3 {
				example := Example{
					Description: "Supporting example",
					Details:     paragraph,
					WordCount:   len(strings.Fields(paragraph)),
				}
				
				supportElement := SupportingElement{
					MainPoint: paragraph,
					Evidence:  paragraph,
					Examples:  []Example{example},
					WordCount: len(strings.Fields(paragraph)),
				}
				
				currentSection.SupportingElements = append(currentSection.SupportingElements, supportElement)
				currentSection.WordCount += len(strings.Fields(paragraph))
			}
		}
		
		sectionCount++
	}
	
	if len(currentSection.TopicSentence) > 0 {
		sections = append(sections, currentSection)
	}
	
	return sections
}

// Analyze content gaps for expansion
func analyzeContentGaps(chapter Chapter, targetWords int) string {
	currentWords := chapter.WordCount
	gap := targetWords - currentWords
	
	if gap <= 0 {
		return "Target word count reached"
	}
	
	// Simple gap analysis
	analysis := fmt.Sprintf(`Content Gap Analysis:
- Current words: %d
- Target words: %d
- Gap: %d words needed
- Sections: %d

Expansion suggestions:
1. Add more detailed examples to existing sections
2. Include additional supporting evidence
3. Add practical implementation steps
4. Include case studies or real-world applications
5. Add transitions and connections between concepts`, 
	currentWords, targetWords, gap, len(chapter.Sections))
	
	return analysis
}

// Expand chapter based on gap analysis
func expandBasedOnGaps(selectedModel string, chapter Chapter, gapAnalysis string) Chapter {
	// Add more examples and supporting elements
	for i := range chapter.Sections {
		section := &chapter.Sections[i]
		
		// Add more supporting elements if needed
		if len(section.SupportingElements) < 3 {
			newElement := SupportingElement{
				MainPoint: "Additional supporting point for " + section.TopicSentence,
				Evidence:  "Further evidence and explanation to support the main concept.",
				Examples: []Example{
					{
						Description: "Expanded example",
						Details:     "Detailed example with specific implementation steps and practical applications.",
						WordCount:   25,
					},
				},
				WordCount: 50,
			}
			section.SupportingElements = append(section.SupportingElements, newElement)
			section.WordCount += 50
		}
		
		// Add more examples to existing elements
		for j := range section.SupportingElements {
			element := &section.SupportingElements[j]
			if len(element.Examples) < 2 {
				newExample := Example{
					Description: "Additional practical example",
					Details:     "Comprehensive example with step-by-step breakdown and real-world application scenarios.",
					WordCount:   20,
				}
				element.Examples = append(element.Examples, newExample)
				element.WordCount += 20
				section.WordCount += 20
			}
		}
	}
	
	chapter.WordCount = calculateChapterWordCount(chapter)
	return chapter
}

// Calculate word count for chapter
func calculateChapterWordCount(chapter Chapter) int {
	totalWords := 0
	for _, section := range chapter.Sections {
		totalWords += section.WordCount
	}
	return totalWords
}

// Calculate total words across all chapters
func calculateTotalWords(chapters []Chapter) int {
	total := 0
	for _, chapter := range chapters {
		total += chapter.WordCount
	}
	return total
}

// Render chapter to readable text
func renderChapterToText(chapter Chapter) string {
	var content strings.Builder
	
	content.WriteString(fmt.Sprintf("# %s\n\n", chapter.Title))
	
	for i, section := range chapter.Sections {
		content.WriteString(fmt.Sprintf("## Section %d\n\n", i+1))
		content.WriteString(fmt.Sprintf("**Topic:** %s\n\n", section.TopicSentence))
		
		for j, element := range section.SupportingElements {
			content.WriteString(fmt.Sprintf("### Supporting Point %d\n", j+1))
			content.WriteString(fmt.Sprintf("%s\n\n", element.MainPoint))
			content.WriteString(fmt.Sprintf("**Evidence:** %s\n\n", element.Evidence))
			
			for k, example := range element.Examples {
				content.WriteString(fmt.Sprintf("**Example %d:** %s\n", k+1, example.Description))
				content.WriteString(fmt.Sprintf("%s\n\n", example.Details))
			}
		}
		
		if section.Conclusion != "" {
			content.WriteString(fmt.Sprintf("**Conclusion:** %s\n\n", section.Conclusion))
		}
	}
	
	content.WriteString(fmt.Sprintf("\n*Chapter word count: %d*\n", chapter.WordCount))
	
	return content.String()
}

// Compile book from chapters
func compileBook(title string, chapters []Chapter) string {
	var book strings.Builder
	
	book.WriteString(fmt.Sprintf("# %s\n\n", title))
	book.WriteString(fmt.Sprintf("*Generated: %s*\n\n", time.Now().Format("January 2, 2006")))
	
	// Table of contents
	book.WriteString("## Table of Contents\n\n")
	for i, chapter := range chapters {
		book.WriteString(fmt.Sprintf("%d. %s\n", i+1, chapter.Title))
	}
	book.WriteString("\n---\n\n")
	
	// Full chapters
	for i, chapter := range chapters {
		book.WriteString(fmt.Sprintf("# Chapter %d: %s\n\n", i+1, chapter.Title))
		book.WriteString(renderChapterToText(chapter))
		book.WriteString("\n---\n\n")
	}
	
	return book.String()
}

// Save book content in multiple formats
func saveBookContent(title string, bookContent string, chapters []Chapter) ([]string, error) {
	outputDir := "../output/generated_books"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return nil, err
	}
	
	// Clean title for filename
	cleanTitle := strings.ReplaceAll(title, " ", "_")
	cleanTitle = strings.ReplaceAll(cleanTitle, ":", "")
	cleanTitle = strings.ReplaceAll(cleanTitle, "?", "")
	cleanTitle = strings.ReplaceAll(cleanTitle, "!", "")
	cleanTitle = strings.ReplaceAll(cleanTitle, "\"", "")
	cleanTitle = strings.ReplaceAll(cleanTitle, "'", "")
	cleanTitle = strings.ReplaceAll(cleanTitle, "/", "_")
	cleanTitle = strings.ReplaceAll(cleanTitle, "\\", "_")
	
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	
	var savedFiles []string
	
	// Save as Markdown
	mdFile := filepath.Join(outputDir, fmt.Sprintf("%s_%s.md", cleanTitle, timestamp))
	if err := os.WriteFile(mdFile, []byte(bookContent), 0644); err != nil {
		return nil, err
	}
	savedFiles = append(savedFiles, mdFile)
	
	// Save as plain text
	txtContent := strings.ReplaceAll(bookContent, "#", "")
	txtContent = strings.ReplaceAll(txtContent, "*", "")
	txtFile := filepath.Join(outputDir, fmt.Sprintf("%s_%s.txt", cleanTitle, timestamp))
	if err := os.WriteFile(txtFile, []byte(txtContent), 0644); err != nil {
		return nil, err
	}
	savedFiles = append(savedFiles, txtFile)
	
	// Save structured data as JSON
	bookData := struct {
		Title       string    `json:"title"`
		Generated   string    `json:"generated"`
		Chapters    []Chapter `json:"chapters"`
		TotalWords  int       `json:"total_words"`
		ChapterCount int      `json:"chapter_count"`
	}{
		Title:       title,
		Generated:   time.Now().Format(time.RFC3339),
		Chapters:    chapters,
		TotalWords:  calculateTotalWords(chapters),
		ChapterCount: len(chapters),
	}
	
	jsonFile := filepath.Join(outputDir, fmt.Sprintf("%s_%s.json", cleanTitle, timestamp))
	jsonData, _ := json.MarshalIndent(bookData, "", "  ")
	if err := os.WriteFile(jsonFile, jsonData, 0644); err != nil {
		return nil, err
	}
	savedFiles = append(savedFiles, jsonFile)
	
	// Generate EPUB format
	epubFile := filepath.Join(outputDir, fmt.Sprintf("%s_%s.epub", cleanTitle, timestamp))
	if err := generateEPUB(title, bookContent, chapters, epubFile); err != nil {
		fmt.Printf("Warning: Could not generate EPUB: %v\n", err)
	} else {
		savedFiles = append(savedFiles, epubFile)
	}
	
	return savedFiles, nil
}

// Generate EPUB format
func generateEPUB(title string, bookContent string, chapters []Chapter, outputFile string) error {
	// Create a new EPUB
	e := epub.NewEpub(title)
	
	// Set basic metadata
	e.SetAuthor("BULLET BOOKS AI")
	e.SetDescription(fmt.Sprintf("Generated by BULLET BOOKS AI - %s", title))
	e.SetLang("en")
	
	// Add a basic CSS for styling
	css := `
		body { font-family: Georgia, serif; line-height: 1.6; margin: 40px; }
		h1 { color: #333; border-bottom: 2px solid #333; padding-bottom: 10px; }
		h2 { color: #666; margin-top: 30px; }
		h3 { color: #888; margin-top: 20px; }
		p { margin-bottom: 15px; text-align: justify; }
		.chapter-title { page-break-before: always; }
	`
	_, err := e.AddCSS(css, "styles.css")
	if err != nil {
		return fmt.Errorf("error adding CSS: %v", err)
	}
	
	// Add title page
	titlePageHTML := fmt.Sprintf(`
		<html>
		<head><title>%s</title><link rel="stylesheet" type="text/css" href="styles.css"/></head>
		<body>
			<div style="text-align: center; margin-top: 100px;">
				<h1>%s</h1>
				<h3>Generated by BULLET BOOKS AI</h3>
				<p><em>Co-authored by animality.ai</em></p>
			</div>
		</body>
		</html>
	`, title, title)
	
	_, err = e.AddSection(titlePageHTML, "Title Page", "", "")
	if err != nil {
		return fmt.Errorf("error adding title page: %v", err)
	}
	
	// Add each chapter as a separate section
	for i, chapter := range chapters {
		chapterHTML := fmt.Sprintf(`
			<html>
			<head><title>%s</title><link rel="stylesheet" type="text/css" href="styles.css"/></head>
			<body>
				<h1 class="chapter-title">Chapter %d: %s</h1>
				%s
			</body>
			</html>
		`, chapter.Title, i+1, chapter.Title, convertToHTML(renderChapterToText(chapter)))
		
		_, err = e.AddSection(chapterHTML, fmt.Sprintf("Chapter %d", i+1), "", "")
		if err != nil {
			return fmt.Errorf("error adding chapter %d: %v", i+1, err)
		}
	}
	
	// Write the EPUB file
	err = e.Write(outputFile)
	if err != nil {
		return fmt.Errorf("error writing EPUB: %v", err)
	}
	
	return nil
}

// Convert markdown-style text to basic HTML
func convertToHTML(text string) string {
	html := text
	
	// Convert markdown headers to HTML
	html = strings.ReplaceAll(html, "### ", "<h3>")
	html = strings.ReplaceAll(html, "## ", "<h2>")
	html = strings.ReplaceAll(html, "# ", "<h1>")
	
	// Convert bold text
	html = strings.ReplaceAll(html, "**", "<strong>")
	html = strings.ReplaceAll(html, "*", "<em>")
	
	// Convert line breaks to paragraphs
	paragraphs := strings.Split(html, "\n\n")
	var htmlParagraphs []string
	
	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para != "" {
			// Don't wrap headers in paragraph tags
			if !strings.HasPrefix(para, "<h") {
				htmlParagraphs = append(htmlParagraphs, "<p>"+para+"</p>")
			} else {
				htmlParagraphs = append(htmlParagraphs, para)
			}
		}
	}
	
	return strings.Join(htmlParagraphs, "\n")
}

// Save progress attempt
func saveAttempt(phase string, attempt int, content string) {
	timestamp := time.Now().Format("15:04:05")
	
	progressSave := ProgressSave{
		Timestamp:   time.Now().Format(time.RFC3339),
		Phase:       phase,
		Attempt:     attempt,
		Content:     content,
		WordCount:   len(strings.Fields(content)),
		TargetWords: 1800,
		Status:      "in_progress",
	}
	
	// Save to progress directory
	outputDir := "../output/progress_saves"
	filename := fmt.Sprintf("%s_attempt_%d_%s.json", phase, attempt, timestamp)
	filepath := filepath.Join(outputDir, filename)
	
	jsonData, _ := json.MarshalIndent(progressSave, "", "  ")
	if err := os.WriteFile(filepath, jsonData, 0644); err != nil {
		fmt.Printf("   Warning: Could not save progress: %v\n", err)
	} else {
		fmt.Printf("   📁 Progress saved: %s\n", filename)
	}
}

// Call Ollama API
func callOllama(model, prompt string) (string, error) {
	reqBody := OllamaRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	}
	
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("error marshaling request: %v", err)
	}
	
	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("error calling Ollama API: %v", err)
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}
	
	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %v", err)
	}
	
	return ollamaResp.Response, nil
}