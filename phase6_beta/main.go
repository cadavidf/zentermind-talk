package main

import (
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

// Book outline data structure from outline generator
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
	Generated        string                  `json:"generated"`
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

// Generated book structure
type GeneratedBook struct {
	Title       string             `json:"title"`
	Subtitle    string             `json:"subtitle"`
	Author      string             `json:"author"`
	Chapters    []GeneratedChapter `json:"chapters"`
	WordCount   int                `json:"word_count"`
	Generated   string             `json:"generated"`
	Metadata    BookMetadata       `json:"metadata"`
}

type GeneratedChapter struct {
	Number    int    `json:"number"`
	Title     string `json:"title"`
	Content   string `json:"content"`
	WordCount int    `json:"word_count"`
}

type BookMetadata struct {
	MemorablePhrase  string   `json:"memorable_phrase"`
	Category         string   `json:"category"`
	TargetAudience   string   `json:"target_audience"`
	CoreThesis       string   `json:"core_thesis"`
	KeyConcepts      []string `json:"key_concepts"`
	SupportingPoints []string `json:"supporting_points"`
	CallToAction     string   `json:"call_to_action"`
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

// ASCII Art Banner
func showBanner() {
	fmt.Println(`
██████╗ ██╗   ██╗██╗     ██╗     ███████╗████████╗    ██████╗  ██████╗  ██████╗ ██╗  ██╗███████╗
██╔══██╗██║   ██║██║     ██║     ██╔════╝╚══██╔══╝    ██╔══██╗██╔═══██╗██╔═══██╗██║ ██╔╝██╔════╝
██████╔╝██║   ██║██║     ██║     █████╗     ██║       ██████╔╝██║   ██║██║   ██║█████╔╝ ███████╗
██╔══██╗██║   ██║██║     ██║     ██╔══╝     ██║       ██╔══██╗██║   ██║██║   ██║██╔═██╗ ╚════██║
██████╔╝╚██████╔╝███████╗███████╗███████╗   ██║       ██████╔╝╚██████╔╝╚██████╔╝██║  ██╗███████║
╚═════╝  ╚═════╝ ╚══════╝╚══════╝╚══════╝   ╚═╝       ╚═════╝  ╚═════╝  ╚═════╝ ╚═╝  ╚═╝╚══════╝

            🎓 PHASE 6 BETA: Outline-Driven Complete Book Generation 🎓
`)
}

func main() {
	showBanner()
	
	// Check if running in automated mode
	automatedMode := os.Getenv("AUTOMATED_MODE") == "true"
	
	if !automatedMode {
		fmt.Println("❌ Phase 6 Beta requires AUTOMATED_MODE=true")
		fmt.Println("This phase is designed for sequential book generation")
		return
	}
	
	// Get environment variables
	bookNumber := os.Getenv("BOOK_NUMBER")
	if bookNumber == "" {
		bookNumber = "001"
	}
	
	model := os.Getenv("OLLAMA_MODEL")
	if model == "" {
		model = "llama3.2"
	}
	
	fmt.Printf("🤖 Running in automated mode\n")
	fmt.Printf("📖 Processing book number: %s\n", bookNumber)
	fmt.Printf("🧠 Using AI model: %s\n", model)
	
	// Load book outline
	fmt.Printf("📋 Loading book outline...\n")
	outline, err := loadBookOutline(bookNumber)
	if err != nil {
		fmt.Printf("❌ Error loading outline: %v\n", err)
		return
	}
	
	fmt.Printf("✅ Loaded outline: %s (%d words)\n", outline.Title, outline.WordCount)
	
	// Generate complete book content
	fmt.Printf("\n📚 Generating complete book content...\n")
	book, err := generateCompleteBook(outline, model)
	if err != nil {
		fmt.Printf("❌ Error generating book: %v\n", err)
		return
	}
	
	// Save book in multiple formats
	fmt.Printf("\n💾 Saving book in multiple formats...\n")
	savedFiles, err := saveBookInAllFormats(book, bookNumber)
	if err != nil {
		fmt.Printf("❌ Error saving book: %v\n", err)
		return
	}
	
	// Save content.json for phase pipeline
	contentFile := "content.json"
	contentData := map[string]interface{}{
		"book":        book,
		"generated":   time.Now().Format(time.RFC3339),
		"phase":       "6_beta",
		"status":      "completed",
		"word_count":  book.WordCount,
		"files_saved": savedFiles,
	}
	
	contentJson, _ := json.MarshalIndent(contentData, "", "  ")
	if err := os.WriteFile(contentFile, contentJson, 0644); err != nil {
		fmt.Printf("⚠️  Warning: Could not save content.json: %v\n", err)
	} else {
		fmt.Printf("✅ Saved content.json for pipeline\n")
	}
	
	// Display completion summary
	fmt.Printf("\n🎉 Phase 6 Beta completed successfully! 🎉\n")
	fmt.Printf("📖 Book: %s\n", book.Title)
	fmt.Printf("📊 Total words: %d\n", book.WordCount)
	fmt.Printf("📚 Chapters: %d\n", len(book.Chapters))
	fmt.Printf("💾 Files saved: %d\n", len(savedFiles))
	
	for _, file := range savedFiles {
		fmt.Printf("   📁 %s\n", file)
	}
}

// Load book outline from JSON file
func loadBookOutline(bookNumber string) (*BookOutlineData, error) {
	outlineFile := fmt.Sprintf("../output/books/book_%s/outline_1000_words.json", bookNumber)
	
	data, err := os.ReadFile(outlineFile)
	if err != nil {
		return nil, fmt.Errorf("could not read outline file %s: %v", outlineFile, err)
	}
	
	var outline BookOutlineData
	if err := json.Unmarshal(data, &outline); err != nil {
		return nil, fmt.Errorf("could not parse outline JSON: %v", err)
	}
	
	return &outline, nil
}

// Generate complete book content using outline
func generateCompleteBook(outline *BookOutlineData, model string) (*GeneratedBook, error) {
	book := &GeneratedBook{
		Title:    outline.Title,
		Subtitle: outline.Subtitle,
		Author:   "BULLET BOOKS AI",
		Chapters: make([]GeneratedChapter, 0, len(outline.Chapters)),
		Generated: time.Now().Format(time.RFC3339),
		Metadata: BookMetadata{
			MemorablePhrase:  outline.MemorablePhrase,
			Category:         outline.Category,
			TargetAudience:   outline.TargetAudience,
			CoreThesis:       outline.CoreThesis,
			KeyConcepts:      outline.KeyConcepts,
			SupportingPoints: outline.SupportingPoints,
			CallToAction:     outline.CallToAction,
		},
	}
	
	totalWords := 0
	
	// Generate content for each chapter
	for i, chapterOutline := range outline.Chapters {
		fmt.Printf("📝 Generating Chapter %d: %s...\n", i+1, chapterOutline.Title)
		
		content, err := generateChapterContent(chapterOutline, outline, model)
		if err != nil {
			return nil, fmt.Errorf("error generating chapter %d: %v", i+1, err)
		}
		
		wordCount := len(strings.Fields(content))
		totalWords += wordCount
		
		chapter := GeneratedChapter{
			Number:    i + 1,
			Title:     chapterOutline.Title,
			Content:   content,
			WordCount: wordCount,
		}
		
		book.Chapters = append(book.Chapters, chapter)
		
		fmt.Printf("✅ Chapter %d completed: %d words\n", i+1, wordCount)
	}
	
	book.WordCount = totalWords
	
	fmt.Printf("📚 Book generation completed: %d total words\n", totalWords)
	return book, nil
}

// Generate content for a single chapter
func generateChapterContent(chapterOutline ChapterOutlineData, bookOutline *BookOutlineData, model string) (string, error) {
	// Build detailed prompt using chapter outline
	prompt := fmt.Sprintf(`Write Chapter %d of the book "%s" with the memorable phrase "%s".

CHAPTER DETAILS:
Title: %s
Purpose: %s

BOOK CONTEXT:
- Core Thesis: %s
- Target Audience: %s
- Category: %s

CHAPTER REQUIREMENTS:
1. Write 1800-2000 words
2. Include engaging introduction
3. Cover these key points: %s
4. Include practical examples and actionable insights
5. Connect back to the memorable phrase: "%s"
6. End with clear takeaways

Write in an engaging, accessible style that delivers real value to readers. Focus on practical application and real-world relevance.`,
		chapterOutline.Number,
		bookOutline.Title,
		bookOutline.MemorablePhrase,
		chapterOutline.Title,
		chapterOutline.Purpose,
		bookOutline.CoreThesis,
		bookOutline.TargetAudience,
		bookOutline.Category,
		strings.Join(chapterOutline.KeyPoints, ", "),
		bookOutline.MemorablePhrase)
	
	content, err := callOllama(model, prompt)
	if err != nil {
		return "", err
	}
	
	return content, nil
}

// Save book in all formats
func saveBookInAllFormats(book *GeneratedBook, bookNumber string) ([]string, error) {
	var savedFiles []string
	
	// Create output directories
	generatedBooksDir := "../output/generated_books"
	bookDir := fmt.Sprintf("../output/books/book_%s", bookNumber)
	
	if err := os.MkdirAll(generatedBooksDir, 0755); err != nil {
		return nil, err
	}
	
	// Clean title for filenames
	cleanTitle := cleanFilename(book.Title)
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	
	// 1. Save as JSON
	jsonFile := filepath.Join(generatedBooksDir, fmt.Sprintf("%s_%s.json", cleanTitle, timestamp))
	jsonData, _ := json.MarshalIndent(book, "", "  ")
	if err := os.WriteFile(jsonFile, jsonData, 0644); err != nil {
		return nil, err
	}
	savedFiles = append(savedFiles, jsonFile)
	
	// Also save to book directory
	bookJsonFile := filepath.Join(bookDir, "complete_book.json")
	os.WriteFile(bookJsonFile, jsonData, 0644)
	
	// 2. Save as Markdown
	mdContent := formatBookAsMarkdown(book)
	mdFile := filepath.Join(generatedBooksDir, fmt.Sprintf("%s_%s.md", cleanTitle, timestamp))
	if err := os.WriteFile(mdFile, []byte(mdContent), 0644); err != nil {
		return nil, err
	}
	savedFiles = append(savedFiles, mdFile)
	
	// Also save to book directory
	bookMdFile := filepath.Join(bookDir, "complete_book.md")
	os.WriteFile(bookMdFile, []byte(mdContent), 0644)
	
	// 3. Save as plain text
	txtContent := formatBookAsText(book)
	txtFile := filepath.Join(generatedBooksDir, fmt.Sprintf("%s_%s.txt", cleanTitle, timestamp))
	if err := os.WriteFile(txtFile, []byte(txtContent), 0644); err != nil {
		return nil, err
	}
	savedFiles = append(savedFiles, txtFile)
	
	// 4. Save as EPUB
	epubFile := filepath.Join(generatedBooksDir, fmt.Sprintf("%s_%s.epub", cleanTitle, timestamp))
	if err := generateEPUB(book, epubFile); err != nil {
		fmt.Printf("⚠️  Warning: Could not generate EPUB: %v\n", err)
	} else {
		savedFiles = append(savedFiles, epubFile)
	}
	
	return savedFiles, nil
}

// Format book as Markdown
func formatBookAsMarkdown(book *GeneratedBook) string {
	var md strings.Builder
	
	md.WriteString(fmt.Sprintf("# %s\n\n", book.Title))
	if book.Subtitle != "" {
		md.WriteString(fmt.Sprintf("## %s\n\n", book.Subtitle))
	}
	
	md.WriteString(fmt.Sprintf("**Author**: %s\n\n", book.Author))
	md.WriteString(fmt.Sprintf("**Generated**: %s\n\n", book.Generated))
	md.WriteString(fmt.Sprintf("**Total Words**: %d\n\n", book.WordCount))
	
	if book.Metadata.MemorablePhrase != "" {
		md.WriteString(fmt.Sprintf("**Memorable Phrase**: \"%s\"\n\n", book.Metadata.MemorablePhrase))
	}
	
	md.WriteString("---\n\n")
	
	// Table of Contents
	md.WriteString("## Table of Contents\n\n")
	for _, chapter := range book.Chapters {
		md.WriteString(fmt.Sprintf("%d. [%s](#chapter-%d)\n", chapter.Number, chapter.Title, chapter.Number))
	}
	md.WriteString("\n---\n\n")
	
	// Chapters
	for _, chapter := range book.Chapters {
		md.WriteString(fmt.Sprintf("## Chapter %d: %s {#chapter-%d}\n\n", chapter.Number, chapter.Title, chapter.Number))
		md.WriteString(chapter.Content)
		md.WriteString("\n\n---\n\n")
	}
	
	// Call to Action
	if book.Metadata.CallToAction != "" {
		md.WriteString("## Call to Action\n\n")
		md.WriteString(book.Metadata.CallToAction)
		md.WriteString("\n\n")
	}
	
	md.WriteString("---\n\n")
	md.WriteString("*Co-authored by animality.ai*\n")
	
	return md.String()
}

// Format book as plain text
func formatBookAsText(book *GeneratedBook) string {
	var txt strings.Builder
	
	txt.WriteString(fmt.Sprintf("%s\n", strings.ToUpper(book.Title)))
	txt.WriteString(strings.Repeat("=", len(book.Title)) + "\n\n")
	
	if book.Subtitle != "" {
		txt.WriteString(fmt.Sprintf("%s\n\n", book.Subtitle))
	}
	
	txt.WriteString(fmt.Sprintf("Author: %s\n", book.Author))
	txt.WriteString(fmt.Sprintf("Generated: %s\n", book.Generated))
	txt.WriteString(fmt.Sprintf("Total Words: %d\n\n", book.WordCount))
	
	if book.Metadata.MemorablePhrase != "" {
		txt.WriteString(fmt.Sprintf("Memorable Phrase: \"%s\"\n\n", book.Metadata.MemorablePhrase))
	}
	
	txt.WriteString(strings.Repeat("-", 60) + "\n\n")
	
	// Chapters
	for _, chapter := range book.Chapters {
		txt.WriteString(fmt.Sprintf("CHAPTER %d: %s\n", chapter.Number, strings.ToUpper(chapter.Title)))
		txt.WriteString(strings.Repeat("-", 40) + "\n\n")
		txt.WriteString(chapter.Content)
		txt.WriteString("\n\n" + strings.Repeat("-", 60) + "\n\n")
	}
	
	// Call to Action
	if book.Metadata.CallToAction != "" {
		txt.WriteString("CALL TO ACTION\n")
		txt.WriteString(strings.Repeat("-", 15) + "\n\n")
		txt.WriteString(book.Metadata.CallToAction)
		txt.WriteString("\n\n")
	}
	
	txt.WriteString("Co-authored by animality.ai\n")
	
	return txt.String()
}

// Generate EPUB format
func generateEPUB(book *GeneratedBook, outputFile string) error {
	// Create a new EPUB
	e := epub.NewEpub(book.Title)
	
	// Set metadata
	e.SetAuthor(book.Author)
	if book.Subtitle != "" {
		e.SetDescription(book.Subtitle)
	} else {
		e.SetDescription(fmt.Sprintf("Generated by BULLET BOOKS AI - %s", book.Title))
	}
	e.SetLang("en")
	
	// Add CSS
	css := `
		body { font-family: Georgia, serif; line-height: 1.6; margin: 40px; }
		h1 { color: #333; border-bottom: 2px solid #333; padding-bottom: 10px; page-break-before: always; }
		h2 { color: #666; margin-top: 30px; }
		p { margin-bottom: 15px; text-align: justify; }
		.metadata { font-style: italic; color: #666; }
	`
	_, err := e.AddCSS(css, "styles.css")
	if err != nil {
		return err
	}
	
	// Add title page
	titlePageHTML := fmt.Sprintf(`
		<html>
		<head><title>%s</title><link rel="stylesheet" type="text/css" href="styles.css"/></head>
		<body>
			<div style="text-align: center; margin-top: 100px;">
				<h1>%s</h1>
				%s
				<p class="metadata">by %s</p>
				<p class="metadata">Co-authored by animality.ai</p>
				%s
			</div>
		</body>
		</html>
	`, book.Title, book.Title, 
		func() string { if book.Subtitle != "" { return fmt.Sprintf("<h2>%s</h2>", book.Subtitle) }; return "" }(),
		book.Author,
		func() string { if book.Metadata.MemorablePhrase != "" { return fmt.Sprintf("<p class=\"metadata\">\"%s\"</p>", book.Metadata.MemorablePhrase) }; return "" }())
	
	_, err = e.AddSection(titlePageHTML, "Title Page", "", "")
	if err != nil {
		return err
	}
	
	// Add each chapter
	for _, chapter := range book.Chapters {
		chapterHTML := fmt.Sprintf(`
			<html>
			<head><title>%s</title><link rel="stylesheet" type="text/css" href="styles.css"/></head>
			<body>
				<h1>Chapter %d: %s</h1>
				%s
			</body>
			</html>
		`, chapter.Title, chapter.Number, chapter.Title, convertToHTML(chapter.Content))
		
		_, err = e.AddSection(chapterHTML, fmt.Sprintf("Chapter %d", chapter.Number), "", "")
		if err != nil {
			return err
		}
	}
	
	// Write EPUB file
	return e.Write(outputFile)
}

// Convert text to basic HTML
func convertToHTML(text string) string {
	html := text
	
	// Convert line breaks to paragraphs
	paragraphs := strings.Split(html, "\n\n")
	var htmlParagraphs []string
	
	for _, para := range paragraphs {
		para = strings.TrimSpace(para)
		if para != "" {
			htmlParagraphs = append(htmlParagraphs, "<p>"+para+"</p>")
		}
	}
	
	return strings.Join(htmlParagraphs, "\n")
}

// Clean filename for safe file system usage
func cleanFilename(title string) string {
	clean := strings.ReplaceAll(title, " ", "_")
	clean = strings.ReplaceAll(clean, ":", "")
	clean = strings.ReplaceAll(clean, "?", "")
	clean = strings.ReplaceAll(clean, "!", "")
	clean = strings.ReplaceAll(clean, "\"", "")
	clean = strings.ReplaceAll(clean, "'", "")
	clean = strings.ReplaceAll(clean, "/", "_")
	clean = strings.ReplaceAll(clean, "\\", "_")
	clean = strings.ReplaceAll(clean, "<", "")
	clean = strings.ReplaceAll(clean, ">", "")
	clean = strings.ReplaceAll(clean, "|", "")
	clean = strings.ReplaceAll(clean, "*", "")
	return clean
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