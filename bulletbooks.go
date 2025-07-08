package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type OllamaResponse struct {
	Response string `json:"response"`
	Done     bool   `json:"done"`
}

type ModelInfo struct {
	Name     string `json:"name"`
	Size     int64  `json:"size"`
	Modified string `json:"modified"`
}

type ModelsResponse struct {
	Models []ModelInfo `json:"models"`
}

type BookRecommendation struct {
	Title           string
	Confidence      float64
	Rank            int
	MarketGapScore  float64
	ImprovedVersion string
}

func getAvailableModels() ([]ModelInfo, error) {
	resp, err := http.Get("http://localhost:11434/api/tags")
	if err != nil {
		return nil, fmt.Errorf("error fetching models: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	var modelsResp ModelsResponse
	if err := json.Unmarshal(body, &modelsResp); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	sort.Slice(modelsResp.Models, func(i, j int) bool {
		return modelsResp.Models[i].Size > modelsResp.Models[j].Size
	})

	return modelsResp.Models, nil
}

func selectModel() (string, error) {
	models, err := getAvailableModels()
	if err != nil {
		return "", err
	}

	if len(models) == 0 {
		return "", fmt.Errorf("no models available")
	}

	fmt.Println("Available Ollama models:")
	for i, model := range models {
		sizeGB := float64(model.Size) / (1024 * 1024 * 1024)
		fmt.Printf("%d. %s (%.1f GB)\n", i+1, model.Name, sizeGB)
	}
	fmt.Printf("\nPress Enter to use the largest model (%s) or select a number: ", models[0].Name)

	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", fmt.Errorf("error reading input: %v", err)
	}

	input = strings.TrimSpace(input)
	if input == "" {
		return models[0].Name, nil
	}

	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > len(models) {
		return "", fmt.Errorf("invalid selection")
	}

	return models[choice-1].Name, nil
}

func parseBookRecommendations(response string) []BookRecommendation {
	var books []BookRecommendation
	
	lines := strings.Split(response, "\n")
	re := regexp.MustCompile(`^\d+\.\s*(.+?)\s*-\s*(\d+(?:\.\d+)?)\s*%`)
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		matches := re.FindStringSubmatch(line)
		if len(matches) >= 3 {
			title := strings.TrimSpace(matches[1])
			confidenceStr := matches[2]
			
			confidence, err := strconv.ParseFloat(confidenceStr, 64)
			if err != nil {
				continue
			}
			
			books = append(books, BookRecommendation{
				Title:      title,
				Confidence: confidence,
			})
		}
	}
	
	return books
}

func rankBooksByConfidence(books []BookRecommendation) []BookRecommendation {
	sort.Slice(books, func(i, j int) bool {
		return books[i].Confidence > books[j].Confidence
	})
	
	for i := range books {
		books[i].Rank = i + 1
	}
	
	return books
}

func getUserConfidenceCutoff() (float64, error) {
	fmt.Print("\nEnter confidence percentage cutoff (default 85%): ")
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return 0, fmt.Errorf("error reading input: %v", err)
	}
	
	input = strings.TrimSpace(input)
	if input == "" {
		return 85.0, nil
	}
	
	cutoff, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid percentage: %v", err)
	}
	
	if cutoff < 0 || cutoff > 100 {
		return 0, fmt.Errorf("percentage must be between 0 and 100")
	}
	
	return cutoff, nil
}

func filterBooksByConfidence(books []BookRecommendation, cutoff float64) ([]BookRecommendation, []BookRecommendation) {
	var passing, failing []BookRecommendation
	
	for _, book := range books {
		if book.Confidence >= cutoff {
			passing = append(passing, book)
		} else {
			failing = append(failing, book)
		}
	}
	
	return passing, failing
}

func validateMarketGap(selectedModel, title string) (float64, error) {
	prompt := fmt.Sprintf(`Analyze the market gap opportunity for this book title: "%s"

Rate the market gap opportunity on a scale of 1-10 based on:
- Market saturation in this topic
- Current demand trends
- Competition level
- Unique positioning potential
- Target audience size

Respond with only a number between 1-10 (can include decimals like 7.5).`, title)

	reqBody := OllamaRequest{
		Model:  selectedModel,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return 0, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return 0, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, fmt.Errorf("error reading response: %v", err)
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return 0, fmt.Errorf("error unmarshaling response: %v", err)
	}

	scoreStr := strings.TrimSpace(ollamaResp.Response)
	re := regexp.MustCompile(`(\d+(?:\.\d+)?)`)
	matches := re.FindStringSubmatch(scoreStr)
	
	if len(matches) == 0 {
		return 0, fmt.Errorf("no valid score found in response")
	}
	
	score, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing score: %v", err)
	}
	
	if score < 1 || score > 10 {
		return 0, fmt.Errorf("score out of range: %f", score)
	}
	
	return score, nil
}

func improveBookTitle(selectedModel, originalTitle string, attempt int) (string, float64, error) {
	if attempt > 3 {
		return originalTitle, 0, fmt.Errorf("max attempts reached")
	}
	
	prompt := fmt.Sprintf(`Improve this book title to make it more trendy and marketable for 2026: "%s"

Create a more compelling, trend-focused title that:
- Incorporates current/future trends (AI, sustainability, digital transformation, etc.)
- Has stronger market appeal
- Sounds more contemporary and engaging
- Maintains the core concept but makes it more attractive

Respond with only the improved title and a specific confidence percentage (like 87.3%%) in this format:
[Improved Title] - [Confidence%%]`, originalTitle)

	reqBody := OllamaRequest{
		Model:  selectedModel,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", 0, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", 0, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", 0, fmt.Errorf("error reading response: %v", err)
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", 0, fmt.Errorf("error unmarshaling response: %v", err)
	}

	re := regexp.MustCompile(`(.+?)\s*-\s*(\d+(?:\.\d+)?)\s*%`)
	matches := re.FindStringSubmatch(strings.TrimSpace(ollamaResp.Response))
	
	if len(matches) < 3 {
		return originalTitle, 0, fmt.Errorf("invalid response format")
	}
	
	improvedTitle := strings.TrimSpace(matches[1])
	confidenceStr := matches[2]
	
	confidence, err := strconv.ParseFloat(confidenceStr, 64)
	if err != nil {
		return originalTitle, 0, fmt.Errorf("error parsing confidence: %v", err)
	}
	
	return improvedTitle, confidence, nil
}

func recursivelyImproveBook(selectedModel string, book BookRecommendation, confidenceCutoff float64, attempt int) BookRecommendation {
	fmt.Printf("Improving: %s (Confidence: %.1f%%, Attempt: %d)\n", book.Title, book.Confidence, attempt)
	
	improvedTitle, newConfidence, err := improveBookTitle(selectedModel, book.Title, attempt)
	if err != nil {
		fmt.Printf("Error improving title: %v\n", err)
		return book
	}
	
	marketGapScore, err := validateMarketGap(selectedModel, improvedTitle)
	if err != nil {
		fmt.Printf("Error validating market gap: %v\n", err)
		marketGapScore = 0
	}
	
	improvedBook := BookRecommendation{
		Title:           improvedTitle,
		Confidence:      newConfidence,
		Rank:            book.Rank,
		MarketGapScore:  marketGapScore,
		ImprovedVersion: improvedTitle,
	}
	
	if newConfidence >= confidenceCutoff && marketGapScore > 7.0 {
		fmt.Printf("✓ Success: %s (Confidence: %.1f%%, Market Gap: %.1f/10)\n", improvedTitle, newConfidence, marketGapScore)
		return improvedBook
	}
	
	if attempt < 3 {
		return recursivelyImproveBook(selectedModel, improvedBook, confidenceCutoff, attempt+1)
	}
	
	fmt.Printf("✗ Failed to meet threshold after 3 attempts\n")
	return improvedBook
}

func displayRankedBooks(books []BookRecommendation) {
	fmt.Println("\n=== VERIFIED & RANKED BULLETBOOKS: Top Trendy 2026 Books ===")
	fmt.Println()
	
	if len(books) == 0 {
		fmt.Println("No valid book recommendations found.")
		return
	}
	
	fmt.Printf("%-4s %-50s %-12s %-10s\n", "Rank", "Title", "Confidence", "Market Gap")
	fmt.Println(strings.Repeat("-", 85))
	
	for _, book := range books {
		gapStr := ""
		if book.MarketGapScore > 0 {
			gapStr = fmt.Sprintf("%.1f/10", book.MarketGapScore)
		}
		fmt.Printf("%-4d %-50s %-12s %-10s\n", book.Rank, book.Title, fmt.Sprintf("%.1f%%", book.Confidence), gapStr)
	}
	
	fmt.Printf("\nTotal recommendations verified: %d\n", len(books))
	if len(books) > 0 {
		fmt.Printf("Highest confidence: %.1f%% - %s\n", books[0].Confidence, books[0].Title)
		fmt.Printf("Lowest confidence: %.1f%% - %s\n", books[len(books)-1].Confidence, books[len(books)-1].Title)
	}
}

func bulletbooks() error {
	selectedModel, err := selectModel()
	if err != nil {
		return fmt.Errorf("error selecting model: %v", err)
	}

	fmt.Printf("\nUsing model: %s\n", selectedModel)
	fmt.Println("Generating trendy book recommendations...")

	prompt := `Generate exactly 10 trendy book title recommendations that would be popular in 2026. For each title, provide a very specific confidence percentage (like 73.2%, 89.7%, etc.) for how trendy it will be in 2026. Format as:

1. [Title] - [Confidence%]
2. [Title] - [Confidence%]
...and so on.

Focus on current trends like AI, climate change, space exploration, virtual reality, and social issues that will be relevant in 2026.`

	reqBody := OllamaRequest{
		Model:  selectedModel,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("error making request to Ollama: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response: %v", err)
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return fmt.Errorf("error unmarshaling response: %v", err)
	}

	fmt.Println("=== BULLETBOOKS: Raw Ollama Response ===")
	fmt.Println()
	fmt.Println(ollamaResp.Response)

	books := parseBookRecommendations(ollamaResp.Response)
	rankedBooks := rankBooksByConfidence(books)
	displayRankedBooks(rankedBooks)

	confidenceCutoff, err := getUserConfidenceCutoff()
	if err != nil {
		return fmt.Errorf("error getting confidence cutoff: %v", err)
	}

	passingBooks, failingBooks := filterBooksByConfidence(rankedBooks, confidenceCutoff)
	
	fmt.Printf("\n=== CONFIDENCE FILTERING (%.1f%% cutoff) ===\n", confidenceCutoff)
	fmt.Printf("Books passing confidence test: %d\n", len(passingBooks))
	fmt.Printf("Books failing confidence test: %d\n", len(failingBooks))

	var finalBooks []BookRecommendation

	for _, book := range passingBooks {
		fmt.Printf("\nValidating market gap for: %s\n", book.Title)
		marketGapScore, err := validateMarketGap(selectedModel, book.Title)
		if err != nil {
			fmt.Printf("Error validating market gap: %v\n", err)
			continue
		}
		
		book.MarketGapScore = marketGapScore
		fmt.Printf("Market gap score: %.1f/10\n", marketGapScore)
		
		if marketGapScore > 7.0 {
			fmt.Printf("✓ Approved: %s (Market Gap: %.1f/10)\n", book.Title, marketGapScore)
			finalBooks = append(finalBooks, book)
		} else {
			fmt.Printf("✗ Rejected: Market gap score %.1f/10 too low\n", marketGapScore)
		}
	}

	fmt.Printf("\n=== PROCESSING FAILED BOOKS ===\n")
	for _, book := range failingBooks {
		fmt.Printf("\nProcessing failed book: %s (%.1f%%)\n", book.Title, book.Confidence)
		improvedBook := recursivelyImproveBook(selectedModel, book, confidenceCutoff, 1)
		
		if improvedBook.Confidence >= confidenceCutoff && improvedBook.MarketGapScore > 7.0 {
			finalBooks = append(finalBooks, improvedBook)
		}
	}

	fmt.Printf("\n=== FINAL RESULTS ===\n")
	sort.Slice(finalBooks, func(i, j int) bool {
		return finalBooks[i].Confidence > finalBooks[j].Confidence
	})
	
	for i := range finalBooks {
		finalBooks[i].Rank = i + 1
	}
	
	displayRankedBooks(finalBooks)

	return nil
}

func main() {
	if err := bulletbooks(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}