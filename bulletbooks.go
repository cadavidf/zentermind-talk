package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
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

type BookConcept struct {
	Title            string
	Description      string
	UniquenessScore  float64
	ViabilityScore   bool
	CommercialScore  float64
	Status           string
	FailureReason    string
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

	input, err := readInputWithTimeout(fmt.Sprintf("\nPress Enter to use the largest model (%s) or select a number: ", models[0].Name), models[0].Name, 11*time.Second)
	if err != nil {
		return "", fmt.Errorf("error reading input: %v", err)
	}

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

func readInputWithTimeout(prompt string, defaultValue string, timeout time.Duration) (string, error) {
	fmt.Printf("%s", prompt)
	
	inputChan := make(chan string, 1)
	errorChan := make(chan error, 1)
	
	go func() {
		reader := bufio.NewReader(os.Stdin)
		input, err := reader.ReadString('\n')
		if err != nil {
			errorChan <- err
			return
		}
		inputChan <- strings.TrimSpace(input)
	}()
	
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()
	
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()
	
	remaining := int(timeout.Seconds())
	
	for {
		select {
		case input := <-inputChan:
			return input, nil
		case err := <-errorChan:
			return "", err
		case <-ticker.C:
			remaining--
			if remaining > 0 {
				fmt.Printf("\rAuto-selecting '%s' in %d seconds... ", defaultValue, remaining)
			}
		case <-ctx.Done():
			fmt.Printf("\rTimeout reached. Using default: %s\n", defaultValue)
			return "", nil
		}
	}
}

func getUserConfidenceCutoff() (float64, error) {
	input, err := readInputWithTimeout("\nEnter confidence percentage cutoff (default 85%): ", "85%", 11*time.Second)
	if err != nil {
		return 0, fmt.Errorf("error reading input: %v", err)
	}
	
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
	
	prompt := fmt.Sprintf(`Improve this book title to make it more popular and marketable for 2026: "%s"

Create a more compelling, trend-focused title that:
- Incorporates current/future trends (AI, sustainability, digital transformation, etc.)
- Has stronger market appeal
- Sounds more contemporary and engaging
- Maintains the core concept but makes it more attractive

Respond ONLY with the improved title and confidence percentage in this EXACT format:
[Improved Title] - [XX.X%%]

Example: The Future of AI Leadership - 89.5%%`, originalTitle)

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

	response := strings.TrimSpace(ollamaResp.Response)
	lines := strings.Split(response, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		re := regexp.MustCompile(`^(.+?)\s*-\s*(\d+(?:\.\d+)?)\s*%\s*$`)
		matches := re.FindStringSubmatch(line)
		
		if len(matches) >= 3 {
			improvedTitle := strings.TrimSpace(matches[1])
			confidenceStr := matches[2]
			
			confidence, err := strconv.ParseFloat(confidenceStr, 64)
			if err != nil {
				continue
			}
			
			return improvedTitle, confidence, nil
		}
	}
	
	return originalTitle, 0, fmt.Errorf("invalid response format: %s", response)
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
	fmt.Println("Generating popular book recommendations...")

	prompt := `Generate exactly 10 popular book title recommendations that would be successful in 2026. For each title, provide a very specific confidence percentage (like 73.2%, 89.7%, etc.) for how popular it will be in 2026. Format as:

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

	if len(finalBooks) > 0 {
		selectedBook, err := selectBook(finalBooks)
		if err != nil {
			return fmt.Errorf("error selecting book: %v", err)
		}
		
		if selectedBook != nil {
			fmt.Printf("\n=== PHASE 2: CONCEPT REFINEMENT ===\n")
			err = processBookConcepts(selectedModel, *selectedBook)
			if err != nil {
				return fmt.Errorf("error processing concepts: %v", err)
			}
		}
	}

	return nil
}

func selectBook(books []BookRecommendation) (*BookRecommendation, error) {
	fmt.Println("\nSelect a book for concept refinement:")
	for i, book := range books {
		fmt.Printf("%d. %s (%.1f%%, Gap: %.1f/10)\n", i+1, book.Title, book.Confidence, book.MarketGapScore)
	}
	
	input, err := readInputWithTimeout("Enter number (or press Enter to skip): ", "skip", 11*time.Second)
	if err != nil {
		return nil, fmt.Errorf("error reading input: %v", err)
	}
	
	if input == "" {
		return nil, nil
	}
	
	choice, err := strconv.Atoi(input)
	if err != nil || choice < 1 || choice > len(books) {
		return nil, fmt.Errorf("invalid selection")
	}
	
	return &books[choice-1], nil
}

func generateConcepts(selectedModel, bookTitle string) ([]BookConcept, error) {
	prompt := fmt.Sprintf(`Generate exactly 5 unique book concepts for: "%s"

For each concept, provide:
1. A clear concept title
2. A brief description (1-2 sentences)

Format as:
1. [Concept Title] - [Description]
2. [Concept Title] - [Description]
...and so on`, bookTitle)

	reqBody := OllamaRequest{
		Model:  selectedModel,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	var concepts []BookConcept
	lines := strings.Split(ollamaResp.Response, "\n")
	re := regexp.MustCompile(`^\d+\.\s*(.+?)\s*-\s*(.+)$`)
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		matches := re.FindStringSubmatch(line)
		if len(matches) >= 3 {
			title := strings.TrimSpace(matches[1])
			description := strings.TrimSpace(matches[2])
			
			concepts = append(concepts, BookConcept{
				Title:       title,
				Description: description,
			})
		}
	}
	
	return concepts, nil
}

func checkUniqueness(selectedModel, conceptTitle, conceptDescription string) (float64, error) {
	prompt := fmt.Sprintf(`Rate the uniqueness of this book concept on a scale of 1-100%%:

Title: "%s"
Description: "%s"

Consider:
- How original is this concept?
- How different is it from existing books?
- Does it offer a fresh perspective?

Respond with only a number between 1-100.`, conceptTitle, conceptDescription)

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

	re := regexp.MustCompile(`(\d+(?:\.\d+)?)`)
	matches := re.FindStringSubmatch(strings.TrimSpace(ollamaResp.Response))
	
	if len(matches) == 0 {
		return 0, fmt.Errorf("no valid score found")
	}
	
	score, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing score: %v", err)
	}
	
	return score, nil
}

func checkViability(selectedModel, conceptTitle, conceptDescription string) (bool, error) {
	prompt := fmt.Sprintf(`Is this book concept viable and feasible to write and publish?

Title: "%s"
Description: "%s"

Consider:
- Is it practical to research and write?
- Are there available sources and information?
- Is it publishable in the market?

Respond with only YES or NO.`, conceptTitle, conceptDescription)

	reqBody := OllamaRequest{
		Model:  selectedModel,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return false, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return false, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return false, fmt.Errorf("error reading response: %v", err)
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return false, fmt.Errorf("error unmarshaling response: %v", err)
	}

	response := strings.ToUpper(strings.TrimSpace(ollamaResp.Response))
	return strings.Contains(response, "YES"), nil
}

func checkCommercial(selectedModel, conceptTitle, conceptDescription string) (float64, error) {
	prompt := fmt.Sprintf(`Rate the commercial potential of this book concept on a scale of 1-100%%:

Title: "%s"
Description: "%s"

Consider:
- Market demand for this topic
- Target audience size
- Revenue potential
- Competition level

Respond with only a number between 1-100.`, conceptTitle, conceptDescription)

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

	re := regexp.MustCompile(`(\d+(?:\.\d+)?)`)
	matches := re.FindStringSubmatch(strings.TrimSpace(ollamaResp.Response))
	
	if len(matches) == 0 {
		return 0, fmt.Errorf("no valid score found")
	}
	
	score, err := strconv.ParseFloat(matches[1], 64)
	if err != nil {
		return 0, fmt.Errorf("error parsing score: %v", err)
	}
	
	return score, nil
}

func generateMoreConcepts(selectedModel, bookTitle string, existingConcepts []BookConcept) ([]BookConcept, error) {
	existingTitles := make([]string, len(existingConcepts))
	for i, concept := range existingConcepts {
		existingTitles[i] = concept.Title
	}
	
	prompt := fmt.Sprintf(`Generate exactly 5 NEW and DIFFERENT book concepts for: "%s"

Avoid these existing concepts:
%s

For each NEW concept, provide:
1. A clear concept title
2. A brief description (1-2 sentences)

Format as:
1. [Concept Title] - [Description]
2. [Concept Title] - [Description]
...and so on`, bookTitle, strings.Join(existingTitles, "\n- "))

	reqBody := OllamaRequest{
		Model:  selectedModel,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	var concepts []BookConcept
	lines := strings.Split(ollamaResp.Response, "\n")
	re := regexp.MustCompile(`^\d+\.\s*(.+?)\s*-\s*(.+)$`)
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		matches := re.FindStringSubmatch(line)
		if len(matches) >= 3 {
			title := strings.TrimSpace(matches[1])
			description := strings.TrimSpace(matches[2])
			
			concepts = append(concepts, BookConcept{
				Title:       title,
				Description: description,
			})
		}
	}
	
	return concepts, nil
}

func recursivelyImproveConcepts(selectedModel, bookTitle string, validConcepts []BookConcept, attemptCount int) ([]BookConcept, error) {
	if attemptCount >= 3 {
		return validConcepts, fmt.Errorf("max attempts reached")
	}
	
	if len(validConcepts) >= 3 {
		return validConcepts, nil
	}
	
	fmt.Printf("Need more concepts. Generating 5 more (attempt %d)...\n", attemptCount+1)
	
	newConcepts, err := generateMoreConcepts(selectedModel, bookTitle, validConcepts)
	if err != nil {
		return validConcepts, fmt.Errorf("error generating more concepts: %v", err)
	}
	
	fmt.Printf("Generated %d new concepts\n", len(newConcepts))
	
	for i, concept := range newConcepts {
		fmt.Printf("\nProcessing new concept %d: %s\n", i+1, concept.Title)
		
		uniqueness, err := checkUniqueness(selectedModel, concept.Title, concept.Description)
		if err != nil {
			fmt.Printf("Error checking uniqueness: %v\n", err)
			continue
		}
		concept.UniquenessScore = uniqueness
		fmt.Printf("Uniqueness: %.1f%%\n", uniqueness)
		
		if uniqueness < 80 {
			fmt.Printf("❌ Concept rejected (uniqueness < 80%%)\n")
			continue
		}
		
		viable, err := checkViability(selectedModel, concept.Title, concept.Description)
		if err != nil {
			fmt.Printf("Error checking viability: %v\n", err)
			continue
		}
		concept.ViabilityScore = viable
		fmt.Printf("Viable: %t\n", viable)
		
		if !viable {
			fmt.Printf("❌ Concept rejected (not viable)\n")
			continue
		}
		
		commercial, err := checkCommercial(selectedModel, concept.Title, concept.Description)
		if err != nil {
			fmt.Printf("Error checking commercial potential: %v\n", err)
			continue
		}
		concept.CommercialScore = commercial
		fmt.Printf("Commercial: %.1f%%\n", commercial)
		
		if commercial < 70 {
			fmt.Printf("❌ Concept rejected (commercial < 70%%)\n")
			continue
		}
		
		fmt.Printf("✅ Concept approved!\n")
		validConcepts = append(validConcepts, concept)
		
		if len(validConcepts) >= 3 {
			break
		}
	}
	
	if len(validConcepts) < 3 {
		return recursivelyImproveConcepts(selectedModel, bookTitle, validConcepts, attemptCount+1)
	}
	
	return validConcepts, nil
}

func showWorkSummary(allConcepts []BookConcept) {
	fmt.Printf("\n=== WORK SUMMARY ===\n")
	fmt.Printf("Total concepts processed: %d\n\n", len(allConcepts))
	
	for i, concept := range allConcepts {
		fmt.Printf("%d. %s\n", i+1, concept.Title)
		fmt.Printf("   Description: %s\n", concept.Description)
		fmt.Printf("   Uniqueness: %.1f%% | Viable: %t | Commercial: %.1f%%\n", 
			concept.UniquenessScore, concept.ViabilityScore, concept.CommercialScore)
		fmt.Printf("   Status: %s", concept.Status)
		if concept.FailureReason != "" {
			fmt.Printf(" - %s", concept.FailureReason)
		}
		fmt.Printf("\n\n")
	}
}

func calculateConceptScore(concept BookConcept) float64 {
	if !concept.ViabilityScore || concept.UniquenessScore == 0 || concept.CommercialScore == 0 {
		return 0
	}
	
	return (concept.UniquenessScore * concept.CommercialScore) / (concept.UniquenessScore + concept.CommercialScore)
}

func selectBestConcept(concepts []BookConcept) *BookConcept {
	var bestConcept *BookConcept
	var bestScore float64 = 0
	
	for i, concept := range concepts {
		if concept.ViabilityScore {
			score := calculateConceptScore(concept)
			if score > bestScore {
				bestScore = score
				bestConcept = &concepts[i]
			}
		}
	}
	
	return bestConcept
}

func getUserDecision() (string, string, error) {
	fmt.Println("What would you like to do?")
	fmt.Println("1. 'phase3' - Proceed to Phase 3 with best concept (quantitative selection)")
	fmt.Println("2. 'continue' - Keep refining automatically")
	fmt.Println("3. 'steer [keyword]' - Guide refinement direction (e.g., 'steer technology')")
	fmt.Println("4. 'restart' - Start concept generation again")
	fmt.Println("5. 'exit' - Exit the program")
	
	input, err := readInputWithTimeout("Enter your choice: ", "phase3", 11*time.Second)
	if err != nil {
		return "", "", fmt.Errorf("error reading input: %v", err)
	}
	
	if input == "" {
		return "phase3", "", nil
	}
	
	parts := strings.SplitN(input, " ", 2)
	
	command := strings.ToLower(parts[0])
	keyword := ""
	if len(parts) > 1 {
		keyword = parts[1]
	}
	
	return command, keyword, nil
}

func generateDirectedConcepts(selectedModel, bookTitle, direction string) ([]BookConcept, error) {
	prompt := fmt.Sprintf(`Generate exactly 5 unique book concepts for: "%s"

Focus the concepts around this direction: %s

For each concept, provide:
1. A clear concept title
2. A brief description (1-2 sentences)

Format as:
1. [Concept Title] - [Description]
2. [Concept Title] - [Description]
...and so on`, bookTitle, direction)

	reqBody := OllamaRequest{
		Model:  selectedModel,
		Prompt: prompt,
		Stream: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("error marshaling request: %v", err)
	}

	resp, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("error making request: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response: %v", err)
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return nil, fmt.Errorf("error unmarshaling response: %v", err)
	}

	var concepts []BookConcept
	lines := strings.Split(ollamaResp.Response, "\n")
	re := regexp.MustCompile(`^\d+\.\s*(.+?)\s*-\s*(.+)$`)
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		
		matches := re.FindStringSubmatch(line)
		if len(matches) >= 3 {
			title := strings.TrimSpace(matches[1])
			description := strings.TrimSpace(matches[2])
			
			concepts = append(concepts, BookConcept{
				Title:       title,
				Description: description,
			})
		}
	}
	
	return concepts, nil
}

func processConceptWithTracking(selectedModel string, concept *BookConcept) bool {
	fmt.Printf("\nProcessing: %s\n", concept.Title)
	
	uniqueness, err := checkUniqueness(selectedModel, concept.Title, concept.Description)
	if err != nil {
		fmt.Printf("Error checking uniqueness: %v\n", err)
		concept.Status = "❌ FAILED"
		concept.FailureReason = "uniqueness check error"
		return false
	}
	concept.UniquenessScore = uniqueness
	fmt.Printf("Uniqueness: %.1f%%\n", uniqueness)
	
	if uniqueness < 80 {
		fmt.Printf("❌ Concept rejected (uniqueness < 80%%)\n")
		concept.Status = "❌ FAILED"
		concept.FailureReason = "uniqueness too low"
		return false
	}
	
	viable, err := checkViability(selectedModel, concept.Title, concept.Description)
	if err != nil {
		fmt.Printf("Error checking viability: %v\n", err)
		concept.Status = "❌ FAILED"
		concept.FailureReason = "viability check error"
		return false
	}
	concept.ViabilityScore = viable
	fmt.Printf("Viable: %t\n", viable)
	
	if !viable {
		fmt.Printf("❌ Concept rejected (not viable)\n")
		concept.Status = "❌ FAILED"
		concept.FailureReason = "not viable"
		return false
	}
	
	commercial, err := checkCommercial(selectedModel, concept.Title, concept.Description)
	if err != nil {
		fmt.Printf("Error checking commercial potential: %v\n", err)
		concept.Status = "❌ FAILED"
		concept.FailureReason = "commercial check error"
		return false
	}
	concept.CommercialScore = commercial
	fmt.Printf("Commercial: %.1f%%\n", commercial)
	
	if commercial < 70 {
		fmt.Printf("❌ Concept rejected (commercial < 70%%)\n")
		concept.Status = "❌ FAILED"
		concept.FailureReason = "commercial potential too low"
		return false
	}
	
	fmt.Printf("✅ Concept approved!\n")
	concept.Status = "✅ APPROVED"
	return true
}

func processBookConcepts(selectedModel string, book BookRecommendation) error {
	startTime := time.Now()
	timeout := 5 * time.Minute
	
	fmt.Printf("Processing concepts for: %s\n", book.Title)
	
	var allConcepts []BookConcept
	var validConcepts []BookConcept
	
	for {
		if time.Since(startTime) > timeout {
			fmt.Printf("\n⏰ Timeout reached (5 minutes). Showing partial results...\n")
			break
		}
		
		var concepts []BookConcept
		var err error
		
		if len(allConcepts) == 0 {
			concepts, err = generateConcepts(selectedModel, book.Title)
			if err != nil {
				return fmt.Errorf("error generating concepts: %v", err)
			}
			fmt.Printf("Generated %d initial concepts\n", len(concepts))
		}
		
		for i := range concepts {
			if time.Since(startTime) > timeout {
				break
			}
			
			if processConceptWithTracking(selectedModel, &concepts[i]) {
				validConcepts = append(validConcepts, concepts[i])
			}
			allConcepts = append(allConcepts, concepts[i])
		}
		
		if len(validConcepts) >= 3 {
			break
		}
		
		if len(validConcepts) == 0 && len(allConcepts) > 0 {
			showWorkSummary(allConcepts)
			
			command, keyword, err := getUserDecision()
			if err != nil {
				return fmt.Errorf("error getting user decision: %v", err)
			}
			
			switch command {
			case "phase3":
				fmt.Printf("\n=== PHASE 3: QUANTITATIVE CONCEPT SELECTION ===\n")
				bestConcept := selectBestConcept(allConcepts)
				if bestConcept == nil {
					fmt.Println("No viable concepts found for Phase 3.")
					fmt.Println("Please try 'continue' or 'steer' to generate more concepts.")
					continue
				}
				
				score := calculateConceptScore(*bestConcept)
				fmt.Printf("Selected best concept using quantitative scoring:\n")
				fmt.Printf("Title: %s\n", bestConcept.Title)
				fmt.Printf("Description: %s\n", bestConcept.Description)
				fmt.Printf("Uniqueness: %.1f%% | Commercial: %.1f%%\n", bestConcept.UniquenessScore, bestConcept.CommercialScore)
				fmt.Printf("Quantitative Score: %.2f\n", score)
				fmt.Printf("Formula: (%.1f * %.1f) / (%.1f + %.1f) = %.2f\n", 
					bestConcept.UniquenessScore, bestConcept.CommercialScore, 
					bestConcept.UniquenessScore, bestConcept.CommercialScore, score)
				
				fmt.Printf("\nProceeding to Phase 3 with selected concept...\n")
				// TODO: Implement Phase 3 functionality here
				fmt.Println("Phase 3 functionality coming soon!")
				return nil
				
			case "continue":
				fmt.Printf("\nContinuing automatic refinement...\n")
				newConcepts, err := generateMoreConcepts(selectedModel, book.Title, allConcepts)
				if err != nil {
					return fmt.Errorf("error generating more concepts: %v", err)
				}
				concepts = newConcepts
				fmt.Printf("Generated %d new concepts\n", len(concepts))
				
			case "steer":
				if keyword == "" {
					fmt.Println("Please provide a keyword for steering direction.")
					continue
				}
				fmt.Printf("\nGenerating concepts focused on: %s\n", keyword)
				newConcepts, err := generateDirectedConcepts(selectedModel, book.Title, keyword)
				if err != nil {
					return fmt.Errorf("error generating directed concepts: %v", err)
				}
				concepts = newConcepts
				fmt.Printf("Generated %d directed concepts\n", len(concepts))
				
			case "restart":
				fmt.Printf("\nRestarting concept generation...\n")
				allConcepts = []BookConcept{}
				validConcepts = []BookConcept{}
				concepts, err = generateConcepts(selectedModel, book.Title)
				if err != nil {
					return fmt.Errorf("error generating concepts: %v", err)
				}
				fmt.Printf("Generated %d fresh concepts\n", len(concepts))
				
			case "exit":
				fmt.Println("Exiting concept refinement.")
				return nil
				
			default:
				fmt.Println("Invalid command. Please try again.")
				continue
			}
		} else {
			break
		}
	}
	
	fmt.Printf("\n=== FINAL CONCEPTS ===\n")
	if len(validConcepts) == 0 {
		fmt.Println("No concepts passed all validations.")
		showWorkSummary(allConcepts)
	} else {
		for i, concept := range validConcepts {
			fmt.Printf("\n%d. %s\n", i+1, concept.Title)
			fmt.Printf("   Description: %s\n", concept.Description)
			fmt.Printf("   Uniqueness: %.1f%% | Viable: %t | Commercial: %.1f%%\n", 
				concept.UniquenessScore, concept.ViabilityScore, concept.CommercialScore)
		}
	}
	
	return nil
}

func main() {
	if err := bulletbooks(); err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}