package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// Ollama API types (copied from main.go)
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

// Book Structure Types
type Section struct {
	Title           string
	WordCount       int
	Events          []string
	QuoteIntegration string
}

type Chapter struct {
	Title   string
	Sections []Section
}

// Book Context Types
type Character struct {
	Name string
	Description string
}

type Motivation struct {
	Entity string
	Motivation string
}

type StyleAnalysis struct {
	WritingStyle string
	WordChoice string
	PlotlineChoices string
	CharacterDevelopmentStyle string
	DetailLevel string
	EmotionalRollercoaster string
}

type BookContext struct {
	BookTitle string
	MainPlot string
	Characters []Character
	Motivations []Motivation
	StyleAnalysis StyleAnalysis
	Chapters []Chapter
}

func main() {
	fmt.Println("Starting 100Monkeys.go - AI-Driven Chapter Generation")

	bookName := "The_Quiet_Bloom_of_Winter" // Hardcoding for the first book
	bookPath := filepath.Join("bullet_books", "100monkeys", bookName)
	outDir := filepath.Join(bookPath, "generated_chapters_ai") // Corrected variable name

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outDir, 0755); err != nil { // Corrected variable name
		fmt.Printf("Error creating output directory: %v\n", err)
		return
	}

	// --- Load all book context ---
	var bookCtx BookContext
	bookCtx.BookTitle = strings.ReplaceAll(bookName, "_", " ")

	// Load Main Plot
	mainPlotContent, err := ioutil.ReadFile(filepath.Join(bookPath, "plot_lines", "main_plot.md"))
	if err != nil {
		fmt.Printf("Error reading main_plot.md: %v\n", err)
		return
	}
	bookCtx.MainPlot = string(mainPlotContent)

	// Load Characters (simplified parsing for now)
	charContent, err := ioutil.ReadFile(filepath.Join(bookPath, "character_development", "main_characters.md"))
	if err != nil {
		fmt.Printf("Error reading main_characters.md: %v\n", err)
		return
	}
	bookCtx.Characters = parseCharacters(string(charContent))

	// Load Motivations (simplified parsing for now)
	motivContent, err := ioutil.ReadFile(filepath.Join(bookPath, "motiv", "core_motivations.md"))
	if err != nil {
		fmt.Printf("Error reading core_motivations.md: %v\n", err)
		return
	}
	bookCtx.Motivations = parseMotivations(string(motivContent))

	// Load Style Analysis (simplified parsing for now)
	styleContent, err := ioutil.ReadFile(filepath.Join(bookPath, "analysis", "style_analysis.md"))
	if err != nil {
		fmt.Printf("Error reading style_analysis.md: %v\n", err)
		return 
	}
	bookCtx.StyleAnalysis = parseStyleAnalysis(string(styleContent))

	// Load Chapter Outline
	outlineContent, err := ioutil.ReadFile(filepath.Join(bookPath, "chapters_parts", "chapter_outline.md"))
	if err != nil {
		fmt.Printf("Error reading chapter_outline.md: %v\n", err)
		return
	}
	bookCtx.Chapters, err = parseDetailedChapterOutline(string(outlineContent))
	if err != nil {
		fmt.Printf("Error parsing detailed chapter outline: %v\n", err)
		return
	}

	fmt.Printf("Successfully loaded context for \"%s\". Found %d chapters.\n", bookCtx.BookTitle, len(bookCtx.Chapters))

	// DEBUG: Print parsed chapters and sections
	fmt.Fprintf(os.Stderr, "\n--- PARSED CHAPTERS AND SECTIONS ---\n")
	for i, ch := range bookCtx.Chapters {
		fmt.Fprintf(os.Stderr, "  Chapter %d: %s (Sections: %d)\n", i+1, ch.Title, len(ch.Sections))
		for j, sec := range ch.Sections {
			fmt.Fprintf(os.Stderr, "    Section %d: %s (Words: %d, Events: %v, Quote: %s)\n", j+1, sec.Title, sec.WordCount, sec.Events, sec.QuoteIntegration)
		}
	}
	fmt.Fprintf(os.Stderr, "--- END PARSED CHAPTERS AND SECTIONS ---\n\n")

	// --- Content Generation Loop ---
	selectedModel := "llama2" // Changed default model to llama2

	for chapterIndex, chapter := range bookCtx.Chapters {
		fmt.Printf("\nGenerating Chapter %d: %s\n", chapterIndex+1, chapter.Title)
		var chapterContentBuilder strings.Builder
		chapterContentBuilder.WriteString(fmt.Sprintf("# %s\n\n", chapter.Title))

		prevSectionContent := ""

		for sectionIndex, section := range chapter.Sections {
			fmt.Printf("  Generating Section %d: %s (Target: %d words)\n", sectionIndex+1, section.Title, section.WordCount)

			nextSectionPreview := ""
			if sectionIndex+1 < len(chapter.Sections) {
				nextSectionPreview = chapter.Sections[sectionIndex+1].Title
			}

			prompt := buildSectionPrompt(
				bookCtx,
				chapter,
				section,
				prevSectionContent,
				nextSectionPreview,
			)
			
			// Call Ollama
			generatedSection, err := callOllama(selectedModel, prompt)
			if err != nil {
				fmt.Printf("    Error calling Ollama for section %s: %v\n", section.Title, err)
				generatedSection = fmt.Sprintf("[[ERROR: Could not generate content for %s - %v]]\n", section.Title, err)
			} else if generatedSection == "" {
				fmt.Printf("    WARNING: Ollama returned empty content for section %s. Prompt was:\n%s\n", section.Title, prompt)
				generatedSection = fmt.Sprintf("[[WARNING: Empty content generated for %s]]\n", section.Title)
			}

			chapterContentBuilder.WriteString(fmt.Sprintf("\n## %s\n\n", section.Title))
			chapterContentBuilder.WriteString(generatedSection)
			chapterContentBuilder.WriteString("\n")

			// Update prevSectionContent for the next iteration
			prevSectionContent = generatedSection
		}

		// Save the completed chapter
		chapterFileName := fmt.Sprintf("Chapter_%02d_%s.md", chapterIndex+1, strings.ReplaceAll(chapter.Title, " ", "_"))
		chapterFilePath := filepath.Join(outDir, chapterFileName) // Corrected variable name
		if err := ioutil.WriteFile(chapterFilePath, []byte(chapterContentBuilder.String()), 0644); err != nil {
			fmt.Printf("Error writing chapter %s: %v\n", chapterFileName, err)
			continue
		}
		fmt.Printf("  ✅ Chapter %s saved to %s\n", chapter.Title, chapterFilePath)
	}

	fmt.Println("AI-driven chapter generation complete.")
}

// --- Parsing Functions (Simplified for initial implementation) ---

func parseCharacters(content string) []Character {
	var chars []Character
	re := regexp.MustCompile(`\*\s*\*\*(.+?):\*\*\s*(.+)`)
	matches := re.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) == 3 {
			chars = append(chars, Character{Name: strings.TrimSpace(match[1]), Description: strings.TrimSpace(match[2])})
		}
	}
	return chars
}

func parseMotivations(content string) []Motivation {
	var motivs []Motivation
	re := regexp.MustCompile(`\*\s*\*\*(.+?):\*\*\s*(.+)`) // Corrected regex for entity and motivation
	matches := re.FindAllStringSubmatch(content, -1)
	for _, match := range matches {
		if len(match) == 3 {
			motivs = append(motivs, Motivation{Entity: strings.TrimSpace(match[1]), Motivation: strings.TrimSpace(match[2])})
		}
	}
	return motivs
}

func parseStyleAnalysis(content string) StyleAnalysis {
	var sa StyleAnalysis
	// This parsing is very basic and assumes specific markdown formatting.
	// A more robust solution would use a markdown parser library.

	lines := strings.Split(content, "\n")
	for _, line := range lines {
		if strings.HasPrefix(line, "**Author's Way of Writing/Style:**") {
			sa.WritingStyle = strings.TrimSpace(strings.TrimPrefix(line, "**Author's Way of Writing/Style:**"))
		} else if strings.HasPrefix(line, "**Word Choice:**") {
			sa.WordChoice = strings.TrimSpace(strings.TrimPrefix(line, "**Word Choice:**"))
		} else if strings.HasPrefix(line, "**Plotline Choices:**") {
			sa.PlotlineChoices = strings.TrimSpace(strings.TrimPrefix(line, "**Plotline Choices:**"))
		} else if strings.HasPrefix(line, "**Character Development Style:**") {
			sa.CharacterDevelopmentStyle = strings.TrimSpace(strings.TrimPrefix(line, "**Character Development Style:**"))
		} else if strings.HasPrefix(line, "**Detail Level on Each Section:**") {
			sa.DetailLevel = strings.TrimSpace(strings.TrimPrefix(line, "**Detail Level on Each Section:**"))
		} else if strings.HasPrefix(line, "**Emotional Rollercoaster:**") {
			sa.EmotionalRollercoaster = strings.TrimSpace(strings.TrimPrefix(line, "**Emotional Rollercoaster:**"))
		}
	}
	return sa
}

func parseDetailedChapterOutline(outline string) ([]Chapter, error) {
	var chapters []Chapter
	chapterRegex := regexp.MustCompile(`(?m)^\*\s*Chapter\s+\d+:\s*(.+?)\s*\(Approx\.\s*(\d+)\s*words\)`) // Matches chapter lines
	sectionTitleRegex := regexp.MustCompile(`(?m)^\s*\*\s*Event:\s*(.+?)$`) // Corrected regex
	quoteRegex := regexp.MustCompile(`(?m)^\s*\*\s*Quote Integration:\s*(.+?)$`) // Corrected regex

	lines := strings.Split(outline, "\n")
	var currentChapter *Chapter
	var currentSection *Section

	for i, line := range lines {
		fmt.Fprintf(os.Stderr, "DEBUG PARSING: Processing line %d: %s\n", i, line) // Uncommented for debugging
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		// Match Chapter
		if chapterMatch := chapterRegex.FindStringSubmatch(line); len(chapterMatch) > 0 {
			fmt.Fprintf(os.Stderr, "DEBUG PARSING: Matched Chapter: %s\n", chapterMatch[1]) // Uncommented for debugging
			if currentChapter != nil && currentSection != nil {
				currentChapter.Sections = append(currentChapter.Sections, *currentSection)
				currentSection = nil // Reset section after appending
			}
			if currentChapter != nil {
				chapters = append(chapters, *currentChapter)
			}
			currentChapter = &Chapter{Title: strings.TrimSpace(chapterMatch[1])}
			continue
		}

		// Match Section Title (Event line)
		if sectionTitleMatch := sectionTitleRegex.FindStringSubmatch(line); len(sectionTitleMatch) > 0 {
			fmt.Fprintf(os.Stderr, "DEBUG PARSING: Matched Section Title: %s\n", sectionTitleMatch[1]) // Uncommented for debugging
			if currentChapter == nil {
				return nil, fmt.Errorf("section found before any chapter: %s", line)
			}
			// If there was a previous section, finalize and append it to the current chapter
			if currentSection != nil {
				currentChapter.Sections = append(currentChapter.Sections, *currentSection)
			}

			// Initialize a new section
			currentSection = &Section{
				Title: strings.TrimSpace(sectionTitleMatch[1]),
				WordCount: 500, // Default word count for sections
				Events: []string{strings.TrimSpace(sectionTitleMatch[1])},
			} 
			continue
		}

		// Match Quote Integration
		if quoteMatch := quoteRegex.FindStringSubmatch(line); len(quoteMatch) > 0 {
			fmt.Fprintf(os.Stderr, "DEBUG PARSING: Matched Quote: %s\n", quoteMatch[1]) // Uncommented for debugging
			if currentSection != nil {
				currentSection.QuoteIntegration = strings.TrimSpace(quoteMatch[1])
			}
			continue
		}

		// If it's not a chapter, section, or quote, assume it's part of the current section's events
		// This part needs to be more robust for multi-line events/descriptions
		// For now, we'll just append the line if it's not empty and we have a current section
		if currentSection != nil && line != "" {
			fmt.Fprintf(os.Stderr, "DEBUG PARSING: Appending Event: %s\n", line) // Uncommented for debugging
			currentSection.Events = append(currentSection.Events, line)
		}
	}

	// Append the last chapter and section if they exist
	// This is crucial for the very last section of the very last chapter
	if currentChapter != nil && currentSection != nil {
		currentChapter.Sections = append(currentChapter.Sections, *currentSection)
	}
	// Append the last chapter if it exists
	if currentChapter != nil {
		chapters = append(chapters, *currentChapter)
	}

	return chapters, nil
}

// --- Ollama API Call ---
func callOllama(model, prompt string) (string, error) {
	// Log the request body for debugging
	fmt.Fprintf(os.Stderr, "\n--- OLLAMA REQUEST BODY ---\nModel: %s\nPrompt Length: %d\n---\n", model, len(prompt))
	os.Stderr.Sync()

	reqBody := OllamaRequest{
		Model:  model,
		Prompt: prompt,
		Stream: false,
	}
	
	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshaling request: %v\n", err)
		os.Stderr.Sync()
		return "", fmt.Errorf("error marshaling request: %v", err)
	}
	
	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Fprintf(os.Stderr, "HTTP POST Error: %v\n", err) // Log the actual HTTP error
		os.Stderr.Sync()
		return "", fmt.Errorf("error calling Ollama API: %v", err)
	}
	defer resp.Body.Close()
	
	// Log the HTTP status code
	fmt.Fprintf(os.Stderr, "OLLAMA RESPONSE STATUS: %s\n", resp.Status)
	os.Stderr.Sync()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Read Body Error: %v\n", err) // Log the actual read error
		os.Stderr.Sync()
		return "", fmt.Errorf("error reading response: %v", err)
	}
	
	// Log the raw response body
	fmt.Fprintf(os.Stderr, "RAW OLLAMA RESPONSE BODY: %s\n---\n", string(body))
	os.Stderr.Sync()

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		fmt.Fprintf(os.Stderr, "JSON Unmarshal Error: %v\n", err) // Log the actual unmarshal error
		os.Stderr.Sync()
		return "", fmt.Errorf("error unmarshaling response: %v", err)
	}
	
	return ollamaResp.Response, nil
}

// --- Prompt Construction ---
func buildSectionPrompt(
	bookCtx BookContext,
	chapter Chapter,
	currentSection Section,
	prevSectionContent string,
	nextSectionPreview string,
) string {
	var prompt strings.Builder

	prompt.WriteString(fmt.Sprintf("You are an expert author writing a book titled \"%s\".\n", bookCtx.BookTitle))
	prompt.WriteString(fmt.Sprintf("The book's main plot is: %s\n", bookCtx.MainPlot))

	prompt.WriteString("\n--- Author's Style Guide ---\n")
	prompt.WriteString(fmt.Sprintf("Writing Style: %s\n", bookCtx.StyleAnalysis.WritingStyle))
	prompt.WriteString(fmt.Sprintf("Word Choice: %s\n", bookCtx.StyleAnalysis.WordChoice))
	prompt.WriteString(fmt.Sprintf("Plotline Choices: %s\n", bookCtx.StyleAnalysis.PlotlineChoices))
	prompt.WriteString(fmt.Sprintf("Character Development Style: %s\n", bookCtx.StyleAnalysis.CharacterDevelopmentStyle))
	prompt.WriteString(fmt.Sprintf("Detail Level: %s\n", bookCtx.StyleAnalysis.DetailLevel))
	prompt.WriteString(fmt.Sprintf("Emotional Rollercoaster: %s\n", bookCtx.StyleAnalysis.EmotionalRollercoaster))
	prompt.WriteString("---------------------------\n\n")

	prompt.WriteString(fmt.Sprintf("You are currently writing Chapter: \"%s\".\n", chapter.Title))
	prompt.WriteString(fmt.Sprintf("The current section you need to write is: \"%s\".\n", currentSection.Title))
	prompt.WriteString(fmt.Sprintf("This section's key events are: %s\n", strings.Join(currentSection.Events, "; ")))
	
	if currentSection.QuoteIntegration != "" {
		prompt.WriteString(fmt.Sprintf("Integrate the following quote naturally into this section: \"%s\"\n", currentSection.QuoteIntegration))
	}

	// Character and Motivation Context
	if len(bookCtx.Characters) > 0 {
		prompt.WriteString("\n--- Character Context ---\n")
		for _, char := range bookCtx.Characters {
			prompt.WriteString(fmt.Sprintf("%s: %s\n", char.Name, char.Description))
		}
		prompt.WriteString("---------------------------\n")
	}

	if len(bookCtx.Motivations) > 0 {
		prompt.WriteString("\n--- Motivation Context ---\n")
		for _, motiv := range bookCtx.Motivations {
			prompt.WriteString(fmt.Sprintf("%s: %s\n", motiv.Entity, motiv.Motivation))
		}
		prompt.WriteString("---------------------------\n")
	}

	// Continuity with previous section
	if prevSectionContent != "" {
		prompt.WriteString("\n--- Previous Section Content (for continuity) ---\n")
		prompt.WriteString(prevSectionContent)
		prompt.WriteString("--------------------------------------------------\n")
	}

	// Preview of next section
	if nextSectionPreview != "" {
		prompt.WriteString(fmt.Sprintf("\n--- Next Section Preview (for smooth transition) ---\n")) 
		prompt.WriteString(fmt.Sprintf("The next section will cover: %s\n", nextSectionPreview))
		prompt.WriteString("--------------------------------------------------\n")
	}

	prompt.WriteString(fmt.Sprintf("\nWrite approximately %d words for this section. Focus on the events and integrate the quote if provided. Maintain the overall book's style and continuity.\n", currentSection.WordCount))
	prompt.WriteString("\nBegin writing the section now:\n")

	return prompt.String()
}