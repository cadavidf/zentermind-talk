package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

// Book outline structure
type BookOutline struct {
	Title            string              `json:"title"`
	Subtitle         string              `json:"subtitle"`
	MemorablePhrase  string              `json:"memorable_phrase"`
	Category         string              `json:"category"`
	TargetAudience   string              `json:"target_audience"`
	CoreThesis       string              `json:"core_thesis"`
	Chapters         []ChapterOutline    `json:"chapters"`
	KeyConcepts      []string            `json:"key_concepts"`
	SupportingPoints []string            `json:"supporting_points"`
	CallToAction     string              `json:"call_to_action"`
	WordCount        int                 `json:"word_count"`
	Generated        string              `json:"generated"`
	ThemeSource      BookTheme           `json:"theme_source"`
}

type ChapterOutline struct {
	Number       int      `json:"number"`
	Title        string   `json:"title"`
	Purpose      string   `json:"purpose"`
	KeyPoints    []string `json:"key_points"`
	Examples     []string `json:"examples"`
	WordTarget   int      `json:"word_target"`
	Takeaways    []string `json:"takeaways"`
}

type BookTheme struct {
	ID             int    `json:"id"`
	Title          string `json:"title"`
	MemorablePhrase string `json:"memorable_phrase"`
	Category       string `json:"category"`
	Description    string `json:"description"`
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

// Generate comprehensive 1000-word book outline
func GenerateBookOutline(theme BookTheme, model string) (*BookOutline, error) {
	fmt.Printf("📝 Generating 1000-word outline for: %s\n", theme.Title)
	fmt.Printf("💡 Memorable phrase: %s\n", theme.MemorablePhrase)
	fmt.Printf("🏷️  Category: %s\n", theme.Category)

	// Generate the core book structure
	coreStructure, err := generateCoreStructure(theme, model)
	if err != nil {
		return nil, fmt.Errorf("error generating core structure: %v", err)
	}

	// Generate detailed chapter outlines
	chapters, err := generateChapterOutlines(theme, coreStructure, model)
	if err != nil {
		return nil, fmt.Errorf("error generating chapter outlines: %v", err)
	}

	// Generate supporting elements
	concepts, points, err := generateSupportingElements(theme, coreStructure, model)
	if err != nil {
		return nil, fmt.Errorf("error generating supporting elements: %v", err)
	}

	// Generate call to action
	callToAction, err := generateCallToAction(theme, coreStructure, model)
	if err != nil {
		return nil, fmt.Errorf("error generating call to action: %v", err)
	}

	// Calculate total word count (should be approximately 1000 words)
	totalWords := calculateOutlineWordCount(coreStructure, chapters, concepts, points, callToAction)

	outline := &BookOutline{
		Title:            theme.Title,
		Subtitle:         extractSubtitle(coreStructure),
		MemorablePhrase:  theme.MemorablePhrase,
		Category:         theme.Category,
		TargetAudience:   extractTargetAudience(coreStructure),
		CoreThesis:       extractCoreThesis(coreStructure),
		Chapters:         chapters,
		KeyConcepts:      concepts,
		SupportingPoints: points,
		CallToAction:     callToAction,
		WordCount:        totalWords,
		Generated:        time.Now().Format(time.RFC3339),
		ThemeSource:      theme,
	}

	fmt.Printf("✅ Generated outline with %d words and %d chapters\n", totalWords, len(chapters))
	return outline, nil
}

// Generate core book structure and thesis
func generateCoreStructure(theme BookTheme, model string) (string, error) {
	prompt := fmt.Sprintf(`Create a book structure for "%s" (memorable phrase: "%s").

Category: %s
Description: %s

Generate in 150 words:
1. Subtitle
2. Target audience 
3. Core thesis
4. Value proposition

Keep it concise and focused.`, 
		theme.Title, theme.MemorablePhrase, theme.Category, theme.Description)

	response, err := callOllama(model, prompt)
	if err != nil {
		return "", err
	}

	return response, nil
}

// Generate detailed chapter outlines
func generateChapterOutlines(theme BookTheme, coreStructure string, model string) ([]ChapterOutline, error) {
	fmt.Printf("📚 Generating chapter structure...\n")
	
	prompt := fmt.Sprintf(`Create 6 chapters for "%s" (phrase: "%s").

Format:
Chapter 1: [Title]
Purpose: [Brief purpose]
Key Points: [2-3 points]

Chapter 2: [Title]
Purpose: [Brief purpose]  
Key Points: [2-3 points]

Continue for all 6 chapters. Keep each chapter description under 50 words.`, 
		theme.Title, theme.MemorablePhrase)

	response, err := callOllama(model, prompt)
	if err != nil {
		return nil, err
	}

	// Parse the response into structured chapter outlines
	chapters := parseChapterOutlines(response)
	
	// Ensure we have 6 chapters and set word targets
	if len(chapters) < 6 {
		// Generate missing chapters
		for i := len(chapters); i < 6; i++ {
			chapters = append(chapters, ChapterOutline{
				Number:     i + 1,
				Title:      fmt.Sprintf("Chapter %d: Advanced Applications", i+1),
				Purpose:    "Provide additional insights and applications",
				KeyPoints:  []string{"Key point 1", "Key point 2", "Key point 3"},
				Examples:   []string{"Example 1", "Example 2"},
				WordTarget: 1800,
				Takeaways:  []string{"Takeaway 1", "Takeaway 2"},
			})
		}
	}

	// Set word targets (aiming for ~1800 words per chapter)
	for i := range chapters {
		chapters[i].WordTarget = 1800
	}

	fmt.Printf("📖 Generated %d chapters\n", len(chapters))
	return chapters, nil
}

// Generate supporting concepts and points
func generateSupportingElements(theme BookTheme, coreStructure string, model string) ([]string, []string, error) {
	fmt.Printf("🔍 Generating supporting elements...\n")
	
	prompt := fmt.Sprintf(`For "%s":

KEY CONCEPTS (5 items):
- [Concept 1]
- [Concept 2]
- [Concept 3]
- [Concept 4]
- [Concept 5]

SUPPORTING POINTS (5 items):
- [Point 1]
- [Point 2]
- [Point 3]
- [Point 4]
- [Point 5]

Keep each item under 15 words.`, theme.Title)

	response, err := callOllama(model, prompt)
	if err != nil {
		return nil, nil, err
	}

	concepts, points := parseSupportingElements(response)
	
	fmt.Printf("💡 Generated %d key concepts and %d supporting points\n", len(concepts), len(points))
	return concepts, points, nil
}

// Generate call to action
func generateCallToAction(theme BookTheme, coreStructure string, model string) (string, error) {
	prompt := fmt.Sprintf(`Write a 50-word call to action for "%s" that motivates readers to take action. Include the phrase "%s".`, theme.Title, theme.MemorablePhrase)

	response, err := callOllama(model, prompt)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(response), nil
}

// Helper functions for parsing AI responses
func parseChapterOutlines(response string) []ChapterOutline {
	var chapters []ChapterOutline
	lines := strings.Split(response, "\n")
	
	var currentChapter *ChapterOutline
	currentSection := ""
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Detect chapter headers
		if strings.Contains(strings.ToLower(line), "chapter") && (strings.Contains(line, "1") || strings.Contains(line, "2") || strings.Contains(line, "3") || strings.Contains(line, "4") || strings.Contains(line, "5") || strings.Contains(line, "6")) {
			if currentChapter != nil {
				chapters = append(chapters, *currentChapter)
			}
			
			chapterNum := len(chapters) + 1
			currentChapter = &ChapterOutline{
				Number:     chapterNum,
				Title:      line,
				KeyPoints:  []string{},
				Examples:   []string{},
				Takeaways:  []string{},
				WordTarget: 1800,
			}
			continue
		}
		
		if currentChapter == nil {
			continue
		}
		
		// Detect sections
		if strings.Contains(strings.ToLower(line), "purpose") {
			currentSection = "purpose"
			continue
		} else if strings.Contains(strings.ToLower(line), "key point") || strings.Contains(strings.ToLower(line), "points") {
			currentSection = "keypoints"
			continue
		} else if strings.Contains(strings.ToLower(line), "example") {
			currentSection = "examples"
			continue
		} else if strings.Contains(strings.ToLower(line), "takeaway") {
			currentSection = "takeaways"
			continue
		}
		
		// Add content to appropriate section
		if strings.HasPrefix(line, "-") || strings.HasPrefix(line, "•") {
			content := strings.TrimLeft(line, "- •")
			content = strings.TrimSpace(content)
			
			switch currentSection {
			case "keypoints":
				currentChapter.KeyPoints = append(currentChapter.KeyPoints, content)
			case "examples":
				currentChapter.Examples = append(currentChapter.Examples, content)
			case "takeaways":
				currentChapter.Takeaways = append(currentChapter.Takeaways, content)
			}
		} else if currentSection == "purpose" && len(line) > 10 {
			currentChapter.Purpose = line
		}
	}
	
	// Add the last chapter
	if currentChapter != nil {
		chapters = append(chapters, *currentChapter)
	}
	
	return chapters
}

func parseSupportingElements(response string) ([]string, []string) {
	var concepts, points []string
	lines := strings.Split(response, "\n")
	
	currentSection := ""
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		// Detect sections
		if strings.Contains(strings.ToUpper(line), "KEY CONCEPTS") {
			currentSection = "concepts"
			continue
		} else if strings.Contains(strings.ToUpper(line), "SUPPORTING POINTS") {
			currentSection = "points"
			continue
		}
		
		// Extract items
		if strings.HasPrefix(line, "-") || strings.HasPrefix(line, "•") {
			content := strings.TrimLeft(line, "- •")
			content = strings.TrimSpace(content)
			
			if currentSection == "concepts" {
				concepts = append(concepts, content)
			} else if currentSection == "points" {
				points = append(points, content)
			}
		}
	}
	
	return concepts, points
}

// Helper functions for extracting information
func extractSubtitle(coreStructure string) string {
	lines := strings.Split(coreStructure, "\n")
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), "subtitle") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return "A Practical Guide to Future Success"
}

func extractTargetAudience(coreStructure string) string {
	lines := strings.Split(coreStructure, "\n")
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), "target audience") || strings.Contains(strings.ToLower(line), "audience") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return "Forward-thinking professionals and leaders"
}

func extractCoreThesis(coreStructure string) string {
	lines := strings.Split(coreStructure, "\n")
	for _, line := range lines {
		if strings.Contains(strings.ToLower(line), "thesis") || strings.Contains(strings.ToLower(line), "argument") {
			parts := strings.Split(line, ":")
			if len(parts) > 1 {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return "Innovation and strategic thinking drive transformative success in the modern world"
}

func calculateOutlineWordCount(coreStructure string, chapters []ChapterOutline, concepts []string, points []string, callToAction string) int {
	wordCount := 0
	
	// Core structure words
	wordCount += len(strings.Fields(coreStructure))
	
	// Chapter outline words
	for _, chapter := range chapters {
		wordCount += len(strings.Fields(chapter.Title))
		wordCount += len(strings.Fields(chapter.Purpose))
		for _, point := range chapter.KeyPoints {
			wordCount += len(strings.Fields(point))
		}
		for _, example := range chapter.Examples {
			wordCount += len(strings.Fields(example))
		}
		for _, takeaway := range chapter.Takeaways {
			wordCount += len(strings.Fields(takeaway))
		}
	}
	
	// Supporting elements words
	for _, concept := range concepts {
		wordCount += len(strings.Fields(concept))
	}
	for _, point := range points {
		wordCount += len(strings.Fields(point))
	}
	
	// Call to action words
	wordCount += len(strings.Fields(callToAction))
	
	return wordCount
}

// Save outline to file
func SaveOutline(outline *BookOutline, bookNumber int) (string, error) {
	outputDir := fmt.Sprintf("output/books/book_%03d", bookNumber)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return "", err
	}
	
	// Save as JSON
	jsonFile := fmt.Sprintf("%s/outline_1000_words.json", outputDir)
	jsonData, err := json.MarshalIndent(outline, "", "  ")
	if err != nil {
		return "", err
	}
	
	if err := os.WriteFile(jsonFile, jsonData, 0644); err != nil {
		return "", err
	}
	
	// Save as readable markdown
	mdFile := fmt.Sprintf("%s/outline_1000_words.md", outputDir)
	mdContent := formatOutlineAsMarkdown(outline)
	if err := os.WriteFile(mdFile, []byte(mdContent), 0644); err != nil {
		return "", err
	}
	
	fmt.Printf("💾 Outline saved to %s\n", outputDir)
	return outputDir, nil
}

func formatOutlineAsMarkdown(outline *BookOutline) string {
	var md strings.Builder
	
	md.WriteString(fmt.Sprintf("# %s\n\n", outline.Title))
	if outline.Subtitle != "" {
		md.WriteString(fmt.Sprintf("## %s\n\n", outline.Subtitle))
	}
	
	md.WriteString(fmt.Sprintf("**Memorable Phrase**: \"%s\"\n\n", outline.MemorablePhrase))
	md.WriteString(fmt.Sprintf("**Category**: %s\n\n", outline.Category))
	md.WriteString(fmt.Sprintf("**Target Audience**: %s\n\n", outline.TargetAudience))
	md.WriteString(fmt.Sprintf("**Core Thesis**: %s\n\n", outline.CoreThesis))
	
	md.WriteString("## Chapter Structure\n\n")
	for _, chapter := range outline.Chapters {
		md.WriteString(fmt.Sprintf("### Chapter %d: %s\n\n", chapter.Number, chapter.Title))
		md.WriteString(fmt.Sprintf("**Purpose**: %s\n\n", chapter.Purpose))
		
		if len(chapter.KeyPoints) > 0 {
			md.WriteString("**Key Points**:\n")
			for _, point := range chapter.KeyPoints {
				md.WriteString(fmt.Sprintf("- %s\n", point))
			}
			md.WriteString("\n")
		}
		
		if len(chapter.Examples) > 0 {
			md.WriteString("**Examples**:\n")
			for _, example := range chapter.Examples {
				md.WriteString(fmt.Sprintf("- %s\n", example))
			}
			md.WriteString("\n")
		}
		
		if len(chapter.Takeaways) > 0 {
			md.WriteString("**Key Takeaways**:\n")
			for _, takeaway := range chapter.Takeaways {
				md.WriteString(fmt.Sprintf("- %s\n", takeaway))
			}
			md.WriteString("\n")
		}
	}
	
	if len(outline.KeyConcepts) > 0 {
		md.WriteString("## Key Concepts\n\n")
		for _, concept := range outline.KeyConcepts {
			md.WriteString(fmt.Sprintf("- %s\n", concept))
		}
		md.WriteString("\n")
	}
	
	if len(outline.SupportingPoints) > 0 {
		md.WriteString("## Supporting Evidence\n\n")
		for _, point := range outline.SupportingPoints {
			md.WriteString(fmt.Sprintf("- %s\n", point))
		}
		md.WriteString("\n")
	}
	
	if outline.CallToAction != "" {
		md.WriteString("## Call to Action\n\n")
		md.WriteString(outline.CallToAction)
		md.WriteString("\n\n")
	}
	
	md.WriteString(fmt.Sprintf("---\n\n"))
	md.WriteString(fmt.Sprintf("**Word Count**: %d words\n", outline.WordCount))
	md.WriteString(fmt.Sprintf("**Generated**: %s\n", outline.Generated))
	md.WriteString(fmt.Sprintf("\nCo-authored by animality.ai\n"))
	
	return md.String()
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