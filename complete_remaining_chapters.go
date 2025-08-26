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
	fmt.Println("🎓 COMPLETING REMAINING CHAPTERS 🎓")
	
	// Chapter 5 and 6 data
	chapters := []ChapterOutlineData{
		{
			Number: 5,
			Title: "**Chapter 5: Reflecting on Responsibility**",
			KeyPoints: []string{
				"Recognizing our role in shaping AI's potential is essential.",
				"We must acknowledge and learn from our biases and mistakes.",
				"Embracing responsibility fosters trust and more effective AI usage.",
			},
			WordTarget: 1800,
		},
		{
			Number: 6,
			Title: "**Chapter 6: The Mirror's Gift**",
			KeyPoints: []string{
				"By embracing AI as a mirror, we can refine ourselves and relationships.",
				"It offers an opportunity for growth, self-awareness, and transformation.",
				"When used thoughtfully, AI can enhance human potential.",
			},
			WordTarget: 1800,
		},
	}
	
	for _, chapter := range chapters {
		fmt.Printf("\n📝 Generating Chapter %d: %s\n", chapter.Number, strings.TrimPrefix(chapter.Title, "**"))
		
		content, err := generateChapterContent(chapter)
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
	
	fmt.Printf("\n🎉 All remaining chapters completed!\n")
}

func generateChapterContent(chapter ChapterOutlineData) (string, error) {
	cleanTitle := strings.TrimPrefix(chapter.Title, "**")
	cleanTitle = strings.TrimSuffix(cleanTitle, "**")
	
	var prompt string
	
	if chapter.Number == 5 {
		prompt = fmt.Sprintf(`Write %s of the book "Your AI, Your Mirror" with the memorable phrase "Your AI doesn't replace you, it reflects you."

CHAPTER DETAILS:
Title: %s
Key Points to Cover:
%s

CHAPTER REQUIREMENTS:
1. Write 1800-2000 words
2. Focus on personal responsibility in AI usage
3. Discuss how recognizing our role in shaping AI leads to better outcomes
4. Include examples of bias acknowledgment and learning from mistakes
5. Show how responsibility builds trust in AI systems
6. Provide practical frameworks for responsible AI usage
7. Connect back to the memorable phrase throughout
8. Include subheadings and actionable takeaways
9. End with reflection questions and practical exercises

Write engaging, practical content that helps readers understand their responsibility in shaping AI outcomes and building trust through mindful usage.`,
			cleanTitle, cleanTitle, "- "+strings.Join(chapter.KeyPoints, "\n- "))
	} else {
		prompt = fmt.Sprintf(`Write %s of the book "Your AI, Your Mirror" with the memorable phrase "Your AI doesn't replace you, it reflects you."

CHAPTER DETAILS:
Title: %s
Key Points to Cover:
%s

This is the FINAL CHAPTER of the book. 

CHAPTER REQUIREMENTS:
1. Write 1800-2000 words
2. Provide a powerful conclusion that ties together all previous chapters
3. Focus on the transformative potential of viewing AI as a mirror
4. Include inspiring examples of growth and self-awareness through AI
5. Discuss how AI can enhance human relationships and potential
6. Provide a clear call-to-action for readers
7. Connect strongly back to the memorable phrase
8. Include subheadings and a compelling conclusion
9. End with final reflection questions and next steps

Write an inspiring, transformative conclusion that leaves readers empowered to use AI as a tool for personal growth and self-discovery.`,
			cleanTitle, cleanTitle, "- "+strings.Join(chapter.KeyPoints, "\n- "))
	}
	
	return callOllama("llama3.2", prompt)
}

func saveChapterToFolder(chapter ChapterOutlineData, content string) error {
	cleanTitle := strings.TrimPrefix(chapter.Title, "**")
	cleanTitle = strings.TrimSuffix(cleanTitle, "**")
	cleanTitle = strings.ReplaceAll(cleanTitle, "Chapter ", "")
	cleanTitle = strings.ReplaceAll(cleanTitle, ": ", "_")
	cleanTitle = strings.ReplaceAll(cleanTitle, " ", "_")
	cleanTitle = strings.ReplaceAll(cleanTitle, "'", "")
	
	folderName := fmt.Sprintf("Chapter_%02d_%s", chapter.Number, cleanTitle)
	folderPath := filepath.Join("output/books/book_001", folderName)
	
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
		"generated_at":  "2025-07-12T16:25:00-05:00",
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