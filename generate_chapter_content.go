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
)

type BookOutlineData struct {
	Title           string                 `json:"title"`
	Subtitle        string                 `json:"subtitle"`
	MemorablePhrase string                 `json:"memorable_phrase"`
	Category        string                 `json:"category"`
	TargetAudience  string                 `json:"target_audience"`
	CoreThesis      string                 `json:"core_thesis"`
	Chapters        []ChapterOutlineData   `json:"chapters"`
	KeyConcepts     []string               `json:"key_concepts"`
	SupportingPoints []string              `json:"supporting_points"`
	CallToAction    string                 `json:"call_to_action"`
	WordCount       int                    `json:"word_count"`
	Generated       string                 `json:"generated"`
}

type ChapterOutlineData struct {
	Number     int      `json:"number"`
	Title      string   `json:"title"`
	Purpose    string   `json:"purpose"`
	KeyPoints  []string `json:"key_points"`
	Examples   []string `json:"examples"`
	WordTarget int      `json:"word_target"`
	Takeaways  []string `json:"takeaways"`
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

func main() {
	fmt.Println("🎓 CHAPTER CONTENT GENERATOR 🎓")
	fmt.Println("Generating individual chapter content for 'Your AI, Your Mirror'")
	
	// Load book outline
	outline, err := loadBookOutline()
	if err != nil {
		fmt.Printf("❌ Error loading outline: %v\n", err)
		return
	}
	
	fmt.Printf("✅ Loaded outline: %s (%d chapters)\n", outline.Title, len(outline.Chapters))
	
	// Generate content for each chapter
	for _, chapter := range outline.Chapters {
		fmt.Printf("\n📝 Generating Chapter %d: %s\n", chapter.Number, strings.TrimPrefix(chapter.Title, "**"))
		
		content, err := generateChapterContent(chapter, outline, "llama3.2")
		if err != nil {
			fmt.Printf("❌ Error generating chapter %d: %v\n", chapter.Number, err)
			continue
		}
		
		// Save chapter content to its folder
		err = saveChapterToFolder(chapter, content)
		if err != nil {
			fmt.Printf("❌ Error saving chapter %d: %v\n", chapter.Number, err)
			continue
		}
		
		wordCount := len(strings.Fields(content))
		fmt.Printf("✅ Chapter %d completed: %d words\n", chapter.Number, wordCount)
	}
	
	fmt.Printf("\n🎉 All chapters generated successfully!\n")
}

func loadBookOutline() (*BookOutlineData, error) {
	outlineFile := "output/books/book_001/outline_1000_words.json"
	
	data, err := os.ReadFile(outlineFile)
	if err != nil {
		return nil, fmt.Errorf("could not read outline file: %v", err)
	}
	
	var outline BookOutlineData
	if err := json.Unmarshal(data, &outline); err != nil {
		return nil, fmt.Errorf("could not parse outline JSON: %v", err)
	}
	
	return &outline, nil
}

func generateChapterContent(chapter ChapterOutlineData, outline *BookOutlineData, model string) (string, error) {
	// Clean up chapter title
	cleanTitle := strings.TrimPrefix(chapter.Title, "**")
	cleanTitle = strings.TrimSuffix(cleanTitle, "**")
	
	prompt := fmt.Sprintf(`Write %s of the book "%s" with the memorable phrase "%s".

CHAPTER DETAILS:
Title: %s
Key Points to Cover:
%s

BOOK CONTEXT:
- Category: %s
- Core Message: %s
- Target Audience: Personal development enthusiasts, AI users, and self-reflection practitioners

CHAPTER REQUIREMENTS:
1. Write 1800-2000 words
2. Include an engaging introduction that hooks the reader
3. Thoroughly explore each key point with practical examples
4. Include real-world scenarios and actionable insights
5. Connect concepts back to the memorable phrase: "%s"
6. End with clear takeaways and reflection questions
7. Use an accessible, engaging writing style
8. Include subheadings to organize the content
9. Provide practical exercises or self-reflection prompts

Write comprehensive, valuable content that helps readers understand how AI truly reflects their own patterns, biases, and potential for growth. Focus on practical application and personal transformation.`,
		cleanTitle,
		outline.Title,
		outline.MemorablePhrase,
		cleanTitle,
		"- " + strings.Join(chapter.KeyPoints, "\n- "),
		outline.Category,
		outline.CallToAction,
		outline.MemorablePhrase)
	
	return callOllama(model, prompt)
}

func saveChapterToFolder(chapter ChapterOutlineData, content string) error {
	// Create folder name
	cleanTitle := strings.TrimPrefix(chapter.Title, "**")
	cleanTitle = strings.TrimSuffix(cleanTitle, "**")
	cleanTitle = strings.ReplaceAll(cleanTitle, "Chapter ", "")
	cleanTitle = strings.ReplaceAll(cleanTitle, ": ", "_")
	cleanTitle = strings.ReplaceAll(cleanTitle, " ", "_")
	cleanTitle = strings.ReplaceAll(cleanTitle, "'", "")
	
	folderName := fmt.Sprintf("Chapter_%02d_%s", chapter.Number, cleanTitle)
	folderPath := filepath.Join("output/books/book_001", folderName)
	
	// Create the folder if it doesn't exist
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return err
	}
	
	// Save content as markdown
	mdFile := filepath.Join(folderPath, fmt.Sprintf("chapter_%02d_content.md", chapter.Number))
	if err := os.WriteFile(mdFile, []byte(content), 0644); err != nil {
		return err
	}
	
	// Save content as text
	txtFile := filepath.Join(folderPath, fmt.Sprintf("chapter_%02d_content.txt", chapter.Number))
	if err := os.WriteFile(txtFile, []byte(content), 0644); err != nil {
		return err
	}
	
	// Save chapter metadata
	metadata := map[string]interface{}{
		"chapter_number": chapter.Number,
		"title":         strings.TrimPrefix(strings.TrimSuffix(chapter.Title, "**"), "**"),
		"key_points":    chapter.KeyPoints,
		"word_target":   chapter.WordTarget,
		"word_count":    len(strings.Fields(content)),
		"generated_at":  "2025-07-12T16:15:00-05:00",
	}
	
	metadataJson, _ := json.MarshalIndent(metadata, "", "  ")
	metadataFile := filepath.Join(folderPath, fmt.Sprintf("chapter_%02d_metadata.json", chapter.Number))
	if err := os.WriteFile(metadataFile, metadataJson, 0644); err != nil {
		return err
	}
	
	return nil
}

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