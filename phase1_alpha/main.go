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
)

// Enhanced USP Analysis Data Structures
type USPVariant struct {
	USP                    string  `json:"usp"`
	NoveltyScore          float64 `json:"novelty_score"`
	ReaderAppealScore     float64 `json:"reader_appeal_score"`
	DifferentiationScore  float64 `json:"differentiation_score"`
	OverallScore          float64 `json:"overall_score"`
	Iteration             int     `json:"iteration"`
	PassesThresholds      bool    `json:"passes_thresholds"`
	AnalysisDetails       string  `json:"analysis_details"`
}

type CompetitiveAnalysis struct {
	CompetitorTitle       string   `json:"competitor_title"`
	CompetitorUSP         string   `json:"competitor_usp"`
	DifferentiatingKeywords []string `json:"differentiating_keywords"`
	SimilarityScore       float64  `json:"similarity_score"`
}

type MarketIntelligence struct {
	TrendingTopics        []string `json:"trending_topics"`
	MarketGaps            []string `json:"market_gaps"`
	TargetAudience        string   `json:"target_audience"`
	MarketSaturationLevel float64  `json:"market_saturation_level"`
	PotentialReach        int      `json:"potential_reach"`
}

type Phase1Output struct {
	Timestamp             string                `json:"timestamp"`
	OriginalConcept       string                `json:"original_concept"`
	FinalOptimizedUSP     string                `json:"final_optimized_usp"`
	SelectionJustification string              `json:"selection_justification"`
	USPVariants           []USPVariant          `json:"usp_variants"`
	CompetitiveAnalysis   []CompetitiveAnalysis `json:"competitive_analysis"`
	MarketIntelligence    MarketIntelligence    `json:"market_intelligence"`
	CriteriaResults       CriteriaResults       `json:"criteria_results"`
	SelectionRubric       CriteriaRubric        `json:"selection_rubric"`
	OptimalAuthor         AuthorPersona         `json:"optimal_author"`
	TotalIterations       int                   `json:"total_iterations"`
	OptimizationSuccess   bool                  `json:"optimization_success"`
}

type CriteriaResults struct {
	NoveltyThreshold        float64 `json:"novelty_threshold"`
	ReaderAppealThreshold   float64 `json:"reader_appeal_threshold"`
	DifferentiationThreshold float64 `json:"differentiation_threshold"`
	NoveltyMet              bool    `json:"novelty_met"`
	ReaderAppealMet         bool    `json:"reader_appeal_met"`
	DifferentiationMet      bool    `json:"differentiation_met"`
	AllCriteriaMet          bool    `json:"all_criteria_met"`
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

type ConceptSuggestion struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Timestamp   string `json:"timestamp"`
}

type ConceptLog struct {
	PreviousConcepts []ConceptSuggestion `json:"previous_concepts"`
	LastGenerated    string              `json:"last_generated"`
}

type SelectionCriteria struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Weight      float64 `json:"weight"`
	Threshold   float64 `json:"threshold"`
	Category    string  `json:"category"`
}

type CriteriaRubric struct {
	Timestamp     string              `json:"timestamp"`
	TotalCriteria int                 `json:"total_criteria"`
	Criteria      []SelectionCriteria `json:"criteria"`
	Purpose       string              `json:"purpose"`
	TargetMarket  string              `json:"target_market"`
}

type AuthorPersona struct {
	Name          string   `json:"name"`
	Background    string   `json:"background"`
	Credentials   []string `json:"credentials"`
	MarketAppeal  float64  `json:"market_appeal"`
	GenreMatch    float64  `json:"genre_match"`
	Trustability  float64  `json:"trustability"`
	Biography     string   `json:"biography"`
}

type ConceptScore struct {
	CriteriaName string  `json:"criteria_name"`
	Score        float64 `json:"score"`
	Weight       float64 `json:"weight"`
	WeightedScore float64 `json:"weighted_score"`
	Threshold    float64 `json:"threshold"`
	Passes       bool    `json:"passes"`
}

type ConceptEvaluation struct {
	Concept       ConceptSuggestion `json:"concept"`
	Scores        []ConceptScore    `json:"scores"`
	TotalWeighted float64           `json:"total_weighted"`
	Rank          int               `json:"rank"`
	PassingCount  int               `json:"passing_count"`
	Analysis      string            `json:"analysis"`
}

func main() {
	showBanner()
	
	// Get model selection from user
	selectedModel, err := getModelInput()
	if err != nil {
		fmt.Printf("Error getting model: %v\n", err)
		return
	}
	
	// Get original concept from user
	originalConcept, selectionRubric, authorPersona, err := getOriginalConcept(selectedModel)
	if err != nil {
		fmt.Printf("Error getting original concept: %v\n", err)
		return
	}
	
	// Create output directories
	createOutputDirectories()
	
	// Start USP optimization process
	result, err := processUSPOptimization(selectedModel, originalConcept, selectionRubric, authorPersona)
	if err != nil {
		fmt.Printf("Error in USP optimization: %v\n", err)
		return
	}
	
	// Save results
	outputFile := "usp_optimization.json"
	if err := saveToJSON(result, outputFile); err != nil {
		fmt.Printf("Error saving output: %v\n", err)
		return
	}
	
	// Display results
	displayResults(result, outputFile)
	
	fmt.Println("\n🎉 Phase 1 Alpha USP Optimization completed successfully! 🎉")
}

func showBanner() {
	fmt.Println(`
██████╗ ██╗   ██╗██╗     ██╗     ███████╗████████╗    ██████╗  ██████╗  ██████╗ ██╗  ██╗███████╗
██╔══██╗██║   ██║██║     ██║     ██╔════╝╚══██╔══╝    ██╔══██╗██╔═══██╗██╔═══██╗██║ ██╔╝██╔════╝
██████╔╝██║   ██║██║     ██║     █████╗     ██║       ██████╔╝██║   ██║██║   ██║█████╔╝ ███████╗
██╔══██╗██║   ██║██║     ██║     ██╔══╝     ██║       ██╔══██╗██║   ██║██║   ██║██╔═██╗ ╚════██║
██████╔╝╚██████╔╝███████╗███████╗███████╗   ██║       ██████╔╝╚██████╔╝╚██████╔╝██║  ██╗███████║
╚═════╝  ╚═════╝ ╚══════╝╚══════╝╚══════╝   ╚═╝       ╚═════╝  ╚═════╝  ╚═════╝ ╚═╝  ╚═╝╚══════╝

            🎯 PHASE 1 ALPHA: Market Intelligence & USP Optimization 🎯
`)
}

func getModelInput() (string, error) {
	fmt.Println("=== MODEL SELECTION ===")
	
	// Check Ollama connectivity first
	fmt.Print("🔍 Checking Ollama connectivity... ")
	if !checkOllamaConnection() {
		fmt.Println("❌ FAILED")
		fmt.Println("⚠️  Ollama is not running or not accessible on localhost:11434")
		fmt.Println("   Please start Ollama and try again.")
		return "", fmt.Errorf("ollama not accessible")
	}
	fmt.Println("✅ Connected")
	
	// Get available models
	models := getAvailableModels()
	
	fmt.Println("\n📋 Available Models:")
	for i, model := range models {
		fmt.Printf("%d. %s\n", i+1, model)
	}
	
	fmt.Print("\nSelect model (1-", len(models), ") or press Enter for default: ")
	
	reader := bufio.NewReader(os.Stdin)
	choice, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	
	choice = strings.TrimSpace(choice)
	if choice == "" {
		return models[0], nil // Default to first model
	}
	
	var choiceNum int
	if _, err := fmt.Sscanf(choice, "%d", &choiceNum); err == nil && choiceNum >= 1 && choiceNum <= len(models) {
		selectedModel := models[choiceNum-1]
		fmt.Printf("✅ Selected: %s\n", selectedModel)
		return selectedModel, nil
	}
	
	fmt.Printf("Invalid choice, using default: %s\n", models[0])
	return models[0], nil
}

func getOriginalConcept(selectedModel string) (string, CriteriaRubric, AuthorPersona, error) {
	fmt.Println("\n=== BOOK DEVELOPMENT STRATEGY ===")
	fmt.Println("Choose your approach:")
	fmt.Println("1. Generate concepts first, then market intelligence-driven criteria (RECOMMENDED)")
	fmt.Println("2. Generate AI-powered concept suggestions directly")
	fmt.Println("3. Enter your own concept")
	fmt.Print("Choice (1, 2, or 3): ")
	
	reader := bufio.NewReader(os.Stdin)
	choice, err := reader.ReadString('\n')
	if err != nil {
		return "", CriteriaRubric{}, AuthorPersona{}, err
	}
	
	choice = strings.TrimSpace(choice)
	
	var rubric CriteriaRubric
	var selectedConcept string
	var authorPersona AuthorPersona
	
	if choice == "1" || choice == "" {
		// NEW ENHANCED FLOW: Concepts → Market Intelligence → Criteria → Author
		fmt.Printf("\n🎯 Step 1: Generating initial concept suggestions...\n")
		conceptSuggestions, err := generateMultipleConceptSuggestions(selectedModel)
		if err != nil {
			return "", CriteriaRubric{}, AuthorPersona{}, fmt.Errorf("error generating concepts: %v", err)
		}
		
		fmt.Printf("\n📊 Step 2: Gathering market intelligence for generated concepts...\n")
		fmt.Println("Analyzing market data to create optimal selection criteria...")
		marketData, err := gatherMarketIntelligenceForConcepts(selectedModel, conceptSuggestions)
		if err != nil {
			return "", CriteriaRubric{}, AuthorPersona{}, fmt.Errorf("error gathering market intelligence: %v", err)
		}
		
		fmt.Printf("\n🎯 Step 3: Creating market intelligence-driven selection criteria...\n")
		rubric, err = generateMarketIntelligenceDrivenCriteria(selectedModel, marketData)
		if err != nil {
			return "", CriteriaRubric{}, AuthorPersona{}, fmt.Errorf("error generating criteria: %v", err)
		}
		
		fmt.Printf("\n📋 Step 4: Final concept selection based on criteria...\n")
		selectedConcept, err = selectOptimalConcept(conceptSuggestions, rubric, selectedModel)
		if err != nil {
			return "", rubric, AuthorPersona{}, fmt.Errorf("error selecting concept: %v", err)
		}
		
		fmt.Printf("\n👤 Step 5: Generating optimal author persona...\n")
		authorPersona, err = generateAuthorPersona(selectedModel, selectedConcept, rubric)
		if err != nil {
			return selectedConcept, rubric, AuthorPersona{}, fmt.Errorf("error generating author: %v", err)
		}
		
	} else if choice == "2" {
		// Generate concepts directly (legacy mode)
		selectedConcept, err = generateConceptSuggestions(selectedModel)
		if err != nil {
			return "", CriteriaRubric{}, AuthorPersona{}, err
		}
		
		// Generate basic rubric and author after concept selection
		rubric = generateBasicRubric()
		authorPersona, _ = generateAuthorPersona(selectedModel, selectedConcept, rubric)
		
	} else {
		// Manual input
		fmt.Println("\n=== MANUAL CONCEPT INPUT ===")
		fmt.Print("Enter your concept: ")
		concept, err := reader.ReadString('\n')
		if err != nil {
			return "", CriteriaRubric{}, AuthorPersona{}, err
		}
		
		selectedConcept = strings.TrimSpace(concept)
		if selectedConcept == "" {
			selectedConcept = "The Empathy Economy" // Default fallback
		}
		
		// Generate basic rubric and author for manual concepts
		rubric = generateBasicRubric()
		authorPersona, _ = generateAuthorPersona(selectedModel, selectedConcept, rubric)
	}
	
	return selectedConcept, rubric, authorPersona, nil
}

func createOutputDirectories() {
	directories := []string{
		"../output/phase1_results",
		"../output/market_intelligence",
		"../output/competitive_analysis",
		"../output/concept_logs",
	}
	
	for _, dir := range directories {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("Warning: Could not create directory %s: %v\n", dir, err)
		}
	}
}

// Generate AI-powered concept suggestions with proper LLM integration
func generateConceptSuggestions(selectedModel string) (string, error) {
	fmt.Printf("\n🤖 Generating AI-powered concept suggestions with %s...\n", selectedModel)
	
	// Load previous concepts to avoid repetition
	conceptLog := loadConceptLog()
	
	// Create exclusion list from previous concepts
	exclusionList := ""
	if len(conceptLog.PreviousConcepts) > 0 {
		exclusionList = "\n\nIMPORTANT: EXCLUDE these previously generated concepts to ensure uniqueness:\n"
		for _, prev := range conceptLog.PreviousConcepts {
			exclusionList += fmt.Sprintf("- %s\n", prev.Title)
		}
		exclusionList += "\nGenerate COMPLETELY DIFFERENT concepts than those listed above.\n"
	}
	
	// Enhanced prompt for better LLM response
	currentYear := time.Now().Year()
	prompt := fmt.Sprintf(`You are an innovative business content strategist. Generate exactly 5 unique, cutting-edge book/content concept titles for %d-%d.

FOCUS AREAS:
- Emerging technologies (AI, quantum computing, biotech, web3)
- Business transformation and leadership evolution
- Future-of-work and organizational design
- Sustainability, climate tech, and circular economy
- Psychology, neuroscience, and human performance
- Cross-industry innovation and convergence

REQUIREMENTS:
- Each concept must be highly novel and distinctive
- Target: business professionals, executives, and thought leaders
- Make titles compelling, marketable, and memorable
- Include trending keywords and emerging themes
- Avoid generic or overused business concepts

%s

RESPONSE FORMAT (follow exactly):
CONCEPT_1: [Unique Title] | [Compelling 10-15 word description] | [Category]
CONCEPT_2: [Unique Title] | [Compelling 10-15 word description] | [Category]
CONCEPT_3: [Unique Title] | [Compelling 10-15 word description] | [Category]
CONCEPT_4: [Unique Title] | [Compelling 10-15 word description] | [Category]
CONCEPT_5: [Unique Title] | [Compelling 10-15 word description] | [Category]

Generate fresh, innovative concepts now:`, currentYear, currentYear+1, exclusionList)
	
	response, err := callOllama(selectedModel, prompt)
	if err != nil {
		fmt.Printf("⚠️  LLM generation failed: %v\n", err)
		fmt.Printf("🔄 Using fallback suggestions...\n")
		suggestions := getFallbackSuggestions()
		return displayAndSelectConcept(suggestions, conceptLog)
	}
	
	// Parse LLM response
	suggestions := parseConceptSuggestions(response)
	if len(suggestions) == 0 {
		fmt.Printf("⚠️  Could not parse LLM response, using fallback suggestions...\n")
		suggestions = getFallbackSuggestions()
	} else {
		fmt.Printf("✅ Successfully generated %d unique concepts\n", len(suggestions))
	}
	
	return displayAndSelectConcept(suggestions, conceptLog)
}

// Display concepts and handle user selection
func displayAndSelectConcept(suggestions []ConceptSuggestion, conceptLog ConceptLog) (string, error) {
	fmt.Printf("\n🎯 Generated Concept Suggestions:\n")
	fmt.Printf("════════════════════════════════════════════\n")
	
	for i, suggestion := range suggestions {
		fmt.Printf("%d. %s\n", i+1, suggestion.Title)
		fmt.Printf("   📝 %s\n", suggestion.Description)
		fmt.Printf("   🏷️  Category: %s\n\n", suggestion.Category)
	}
	
	fmt.Printf("Select concept (1-%d) or press Enter for concept 1: ", len(suggestions))
	reader := bufio.NewReader(os.Stdin)
	choice, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	
	choice = strings.TrimSpace(choice)
	selectedIndex := 0 // Default to first suggestion
	
	if choice != "" {
		var choiceNum int
		if _, err := fmt.Sscanf(choice, "%d", &choiceNum); err == nil && choiceNum >= 1 && choiceNum <= len(suggestions) {
			selectedIndex = choiceNum - 1
		}
	}
	
	selectedConcept := suggestions[selectedIndex]
	fmt.Printf("✅ Selected: %s\n", selectedConcept.Title)
	
	// Log the selected concept
	conceptLog.PreviousConcepts = append(conceptLog.PreviousConcepts, selectedConcept)
	conceptLog.LastGenerated = time.Now().Format(time.RFC3339)
	saveConceptLog(conceptLog)
	
	return selectedConcept.Title, nil
}

// Parse concept suggestions from AI response
func parseConceptSuggestions(response string) []ConceptSuggestion {
	var suggestions []ConceptSuggestion
	lines := strings.Split(response, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "CONCEPT_") {
			parts := strings.Split(line, "|")
			if len(parts) >= 3 {
				titlePart := strings.TrimSpace(parts[0])
				title := strings.TrimSpace(strings.Split(titlePart, ":")[1])
				description := strings.TrimSpace(parts[1])
				category := strings.TrimSpace(parts[2])
				
				suggestion := ConceptSuggestion{
					Title:       title,
					Description: description,
					Category:    category,
					Timestamp:   time.Now().Format(time.RFC3339),
				}
				suggestions = append(suggestions, suggestion)
			}
		}
	}
	
	return suggestions
}

// Load concept log from file
func loadConceptLog() ConceptLog {
	logFile := "../output/concept_logs/concept_history.json"
	
	data, err := os.ReadFile(logFile)
	if err != nil {
		// Return empty log if file doesn't exist
		return ConceptLog{
			PreviousConcepts: []ConceptSuggestion{},
			LastGenerated:    "",
		}
	}
	
	var log ConceptLog
	if err := json.Unmarshal(data, &log); err != nil {
		// Return empty log if parsing fails
		return ConceptLog{
			PreviousConcepts: []ConceptSuggestion{},
			LastGenerated:    "",
		}
	}
	
	return log
}

// Save concept log to file
func saveConceptLog(log ConceptLog) {
	logFile := "../output/concept_logs/concept_history.json"
	
	// Keep only last 50 concepts to prevent infinite growth
	if len(log.PreviousConcepts) > 50 {
		log.PreviousConcepts = log.PreviousConcepts[len(log.PreviousConcepts)-50:]
	}
	
	data, err := json.MarshalIndent(log, "", "  ")
	if err != nil {
		fmt.Printf("Warning: Could not marshal concept log: %v\n", err)
		return
	}
	
	if err := os.WriteFile(logFile, data, 0644); err != nil {
		fmt.Printf("Warning: Could not save concept log: %v\n", err)
	}
}

// Enhanced fallback suggestions that avoid previous concepts
func getFallbackSuggestions() []ConceptSuggestion {
	// Load previous concepts to avoid duplication
	conceptLog := loadConceptLog()
	previousTitles := make(map[string]bool)
	for _, prev := range conceptLog.PreviousConcepts {
		previousTitles[prev.Title] = true
	}
	
	// Full pool of fallback suggestions
	allSuggestions := []ConceptSuggestion{
		{
			Title:       "The Quantum Mindset",
			Description: "Thinking beyond traditional business paradigms for breakthrough results",
			Category:    "Innovation & Strategy",
			Timestamp:   time.Now().Format(time.RFC3339),
		},
		{
			Title:       "The Convergence Economy",
			Description: "How industry boundaries are dissolving to create new market opportunities",
			Category:    "Business Transformation",
			Timestamp:   time.Now().Format(time.RFC3339),
		},
		{
			Title:       "Neurometic Leadership",
			Description: "Using brain science to optimize decision-making and team performance",
			Category:    "Psychology & Leadership",
			Timestamp:   time.Now().Format(time.RFC3339),
		},
		{
			Title:       "The Symbiotic Enterprise",
			Description: "Building business ecosystems that thrive through mutual benefit",
			Category:    "Future of Business",
			Timestamp:   time.Now().Format(time.RFC3339),
		},
		{
			Title:       "The Antifragile Organization",
			Description: "How companies can gain strength from volatility and uncertainty",
			Category:    "Organizational Resilience",
			Timestamp:   time.Now().Format(time.RFC3339),
		},
		{
			Title:       "The Invisible Revolution",
			Description: "Hidden technologies reshaping every aspect of business and society",
			Category:    "Technology & Innovation",
			Timestamp:   time.Now().Format(time.RFC3339),
		},
		{
			Title:       "The Compassionate Capitalist",
			Description: "Profit with purpose in the age of stakeholder expectations",
			Category:    "Sustainable Business",
			Timestamp:   time.Now().Format(time.RFC3339),
		},
		{
			Title:       "The Network Effect Playbook",
			Description: "Building platforms that create exponential value through connections",
			Category:    "Digital Strategy",
			Timestamp:   time.Now().Format(time.RFC3339),
		},
	}
	
	// Filter out previously used suggestions
	var availableSuggestions []ConceptSuggestion
	for _, suggestion := range allSuggestions {
		if !previousTitles[suggestion.Title] {
			availableSuggestions = append(availableSuggestions, suggestion)
		}
	}
	
	// If we need more suggestions, add some time-stamped variants
	if len(availableSuggestions) < 5 {
		currentYear := time.Now().Year()
		additionalSuggestions := []ConceptSuggestion{
			{
				Title:       fmt.Sprintf("The %d Transformation", currentYear+1),
				Description: "Navigating the next wave of technological and social change",
				Category:    "Future Trends",
				Timestamp:   time.Now().Format(time.RFC3339),
			},
			{
				Title:       "The Acceleration Paradox",
				Description: "Why faster isn't always better in business and life",
				Category:    "Business Philosophy",
				Timestamp:   time.Now().Format(time.RFC3339),
			},
		}
		
		for _, suggestion := range additionalSuggestions {
			if !previousTitles[suggestion.Title] {
				availableSuggestions = append(availableSuggestions, suggestion)
			}
		}
	}
	
	// Return first 5 available suggestions
	if len(availableSuggestions) >= 5 {
		return availableSuggestions[:5]
	}
	
	return availableSuggestions
}

// Check Ollama connectivity
func checkOllamaConnection() bool {
	resp, err := http.Get("http://localhost:11434/api/version")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

// Get available models from Ollama
func getAvailableModels() []string {
	resp, err := http.Get("http://localhost:11434/api/tags")
	if err != nil {
		// Return default models if API call fails
		return []string{"llama3.1", "llama3.1:70b", "mistral", "codellama", "qwen2.5"}
	}
	defer resp.Body.Close()
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return []string{"llama3.1", "llama3.1:70b", "mistral", "codellama", "qwen2.5"}
	}
	
	var result struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}
	
	if err := json.Unmarshal(body, &result); err != nil {
		return []string{"llama3.1", "llama3.1:70b", "mistral", "codellama", "qwen2.5"}
	}
	
	var models []string
	for _, model := range result.Models {
		models = append(models, model.Name)
	}
	
	if len(models) == 0 {
		return []string{"llama3.1", "llama3.1:70b", "mistral", "codellama", "qwen2.5"}
	}
	
	return models
}

// Call Ollama API with improved error handling
func callOllama(model, prompt string) (string, error) {
	fmt.Printf("   🤖 Calling %s...\n", model)
	
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
	
	if resp.StatusCode != 200 {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("ollama API error (status %d): %s", resp.StatusCode, string(body))
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("error reading response: %v", err)
	}
	
	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", fmt.Errorf("error unmarshaling response: %v", err)
	}
	
	if ollamaResp.Response == "" {
		return "", fmt.Errorf("empty response from Ollama")
	}
	
	fmt.Printf("   ✅ Response received (%d chars)\n", len(ollamaResp.Response))
	return ollamaResp.Response, nil
}

func displayResults(result Phase1Output, outputFile string) {
	fmt.Printf("\n🎯 USP OPTIMIZATION RESULTS\n")
	fmt.Printf("══════════════════════════════════════════════════════════════════════\n")
	fmt.Printf("📊 Original Concept: %s\n", result.OriginalConcept)
	fmt.Printf("🚀 Final Optimized USP: %s\n", result.FinalOptimizedUSP)
	fmt.Printf("📈 Total Iterations: %d\n", result.TotalIterations)
	fmt.Printf("✅ Optimization Success: %v\n", result.OptimizationSuccess)
	
	fmt.Printf("\n📋 CRITERIA RESULTS:\n")
	fmt.Printf("   Novelty Score: %.3f (threshold: %.2f) %s\n", 
		result.CriteriaResults.NoveltyThreshold, 0.60, 
		getPassFailIcon(result.CriteriaResults.NoveltyMet))
	fmt.Printf("   Reader Appeal: %.3f (threshold: %.2f) %s\n", 
		result.CriteriaResults.ReaderAppealThreshold, 0.75, 
		getPassFailIcon(result.CriteriaResults.ReaderAppealMet))
	fmt.Printf("   Differentiation: %.3f (threshold: %.2f) %s\n", 
		result.CriteriaResults.DifferentiationThreshold, 0.70, 
		getPassFailIcon(result.CriteriaResults.DifferentiationMet))
	
	fmt.Printf("\n🎲 USP VARIANTS EXPLORED:\n")
	for i, variant := range result.USPVariants {
		if i >= 3 { // Show top 3
			break
		}
		fmt.Printf("   %d. %s (Overall: %.3f)\n", i+1, variant.USP, variant.OverallScore)
	}
	
	fmt.Printf("\n💡 SELECTION JUSTIFICATION:\n")
	fmt.Printf("   %s\n", result.SelectionJustification)
	
	outputPath, _ := filepath.Abs(outputFile)
	fmt.Printf("\n💾 Full results saved to:\n")
	fmt.Printf("   📁 %s\n", outputPath)
	fmt.Printf("   🔗 file://%s\n", outputPath)
}

func getPassFailIcon(passed bool) string {
	if passed {
		return "✅"
	}
	return "❌"
}

func saveToJSON(data interface{}, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()
	
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(data)
}

// Main USP optimization process
func processUSPOptimization(selectedModel, originalConcept string, selectionRubric CriteriaRubric, authorPersona AuthorPersona) (Phase1Output, error) {
	fmt.Printf("\n🎯 Starting USP Optimization for \"%s\"\n", originalConcept)
	fmt.Println(strings.Repeat("=", 70))
	
	startTime := time.Now()
	
	// Step 1: Market Intelligence Gathering
	fmt.Printf("\n📊 Step 1: Gathering market intelligence...\n")
	marketIntel, err := gatherMarketIntelligence(selectedModel, originalConcept)
	if err != nil {
		return Phase1Output{}, fmt.Errorf("error gathering market intelligence: %v", err)
	}
	
	// Step 2: Competitive Analysis
	fmt.Printf("\n🔍 Step 2: Analyzing competitive landscape...\n")
	competitiveAnalysis, err := performCompetitiveAnalysis(selectedModel, originalConcept)
	if err != nil {
		return Phase1Output{}, fmt.Errorf("error performing competitive analysis: %v", err)
	}
	
	// Step 3: Initial USP Evaluation
	fmt.Printf("\n🎯 Step 3: Evaluating initial USP...\n")
	initialUSP := USPVariant{
		USP:       originalConcept,
		Iteration: 0,
	}
	
	initialUSP, err = evaluateUSP(selectedModel, initialUSP, competitiveAnalysis)
	if err != nil {
		return Phase1Output{}, fmt.Errorf("error evaluating initial USP: %v", err)
	}
	
	uspVariants := []USPVariant{initialUSP}
	
	// Step 4: Recursive Optimization Loop
	fmt.Printf("\n🔄 Step 4: Recursive USP optimization...\n")
	finalUSP, allVariants, err := optimizeUSPRecursively(selectedModel, originalConcept, competitiveAnalysis, uspVariants)
	if err != nil {
		return Phase1Output{}, fmt.Errorf("error in recursive optimization: %v", err)
	}
	
	// Step 5: Final Decision and Justification
	fmt.Printf("\n🎯 Step 5: Final USP selection and justification...\n")
	justification, err := generateSelectionJustification(selectedModel, finalUSP, allVariants)
	if err != nil {
		return Phase1Output{}, fmt.Errorf("error generating justification: %v", err)
	}
	
	// Compile results
	result := Phase1Output{
		Timestamp:              time.Now().Format(time.RFC3339),
		OriginalConcept:        originalConcept,
		FinalOptimizedUSP:      finalUSP.USP,
		SelectionJustification: justification,
		USPVariants:            allVariants,
		CompetitiveAnalysis:    competitiveAnalysis,
		MarketIntelligence:     marketIntel,
		CriteriaResults: CriteriaResults{
			NoveltyThreshold:         0.60,
			ReaderAppealThreshold:    0.75,
			DifferentiationThreshold: 0.70,
			NoveltyMet:              finalUSP.NoveltyScore >= 0.60,
			ReaderAppealMet:         finalUSP.ReaderAppealScore >= 0.75,
			DifferentiationMet:      finalUSP.DifferentiationScore >= 0.70,
			AllCriteriaMet:          finalUSP.PassesThresholds,
		},
		SelectionRubric:     selectionRubric,
		OptimalAuthor:       authorPersona,
		TotalIterations:     len(allVariants),
		OptimizationSuccess: finalUSP.PassesThresholds,
	}
	
	fmt.Printf("\n⏱️ Total optimization time: %v\n", time.Since(startTime))
	
	return result, nil
}

// Market Intelligence Gathering
func gatherMarketIntelligence(selectedModel, concept string) (MarketIntelligence, error) {
	prompt := fmt.Sprintf(`Analyze the market intelligence for the concept: "%s"

	As a market research expert, provide analysis on:
	1. Current trending topics in this domain
	2. Identified market gaps and opportunities
	3. Target audience demographics and psychographics
	4. Market saturation level (0.0-1.0)
	5. Potential market reach estimation

	Format your response as:
	TRENDING_TOPICS: [list 3-5 trending topics]
	MARKET_GAPS: [list 3-4 market gaps]
	TARGET_AUDIENCE: [detailed description]
	SATURATION_LEVEL: [0.0-1.0 with explanation]
	POTENTIAL_REACH: [estimated number with rationale]`, concept)

	response, err := callOllama(selectedModel, prompt)
	if err != nil {
		return MarketIntelligence{}, err
	}

	// Parse response into structured data
	intel := MarketIntelligence{}
	lines := strings.Split(response, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "TRENDING_TOPICS:") {
			topicsStr := strings.TrimPrefix(line, "TRENDING_TOPICS:")
			intel.TrendingTopics = parseListFromString(topicsStr)
		} else if strings.HasPrefix(line, "MARKET_GAPS:") {
			gapsStr := strings.TrimPrefix(line, "MARKET_GAPS:")
			intel.MarketGaps = parseListFromString(gapsStr)
		} else if strings.HasPrefix(line, "TARGET_AUDIENCE:") {
			intel.TargetAudience = strings.TrimPrefix(line, "TARGET_AUDIENCE:")
		} else if strings.HasPrefix(line, "SATURATION_LEVEL:") {
			saturationStr := strings.TrimPrefix(line, "SATURATION_LEVEL:")
			// Extract numeric value (simplified parsing)
			if strings.Contains(saturationStr, "0.") {
				fmt.Sscanf(saturationStr, "%f", &intel.MarketSaturationLevel)
			} else {
				intel.MarketSaturationLevel = 0.5 // Default moderate saturation
			}
		} else if strings.HasPrefix(line, "POTENTIAL_REACH:") {
			reachStr := strings.TrimPrefix(line, "POTENTIAL_REACH:")
			// Extract numeric value (simplified parsing)
			fmt.Sscanf(reachStr, "%d", &intel.PotentialReach)
		}
	}
	
	// Set defaults if parsing failed
	if len(intel.TrendingTopics) == 0 {
		intel.TrendingTopics = []string{"Digital transformation", "Sustainability", "Remote work"}
	}
	if len(intel.MarketGaps) == 0 {
		intel.MarketGaps = []string{"Practical implementation", "Measurable results", "Scalable solutions"}
	}
	if intel.TargetAudience == "" {
		intel.TargetAudience = "Business leaders and professionals seeking actionable insights"
	}
	if intel.PotentialReach == 0 {
		intel.PotentialReach = 50000
	}
	
	return intel, nil
}

// Competitive Analysis
func performCompetitiveAnalysis(selectedModel, concept string) ([]CompetitiveAnalysis, error) {
	prompt := fmt.Sprintf(`Analyze the competitive landscape for the concept: "%s"

	As a competitive intelligence expert, identify the top 5 bestselling books in similar categories and analyze:
	1. Competitor title
	2. Their unique selling proposition
	3. Key differentiating keywords/phrases
	4. Similarity score to our concept (0.0-1.0)

	Format your response as:
	COMPETITOR_1: [Title] | [USP] | [Keywords] | [Similarity Score]
	COMPETITOR_2: [Title] | [USP] | [Keywords] | [Similarity Score]
	COMPETITOR_3: [Title] | [USP] | [Keywords] | [Similarity Score]
	COMPETITOR_4: [Title] | [USP] | [Keywords] | [Similarity Score]
	COMPETITOR_5: [Title] | [USP] | [Keywords] | [Similarity Score]`, concept)

	response, err := callOllama(selectedModel, prompt)
	if err != nil {
		return nil, err
	}

	var competitors []CompetitiveAnalysis
	lines := strings.Split(response, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "COMPETITOR_") {
			parts := strings.Split(line, "|")
			if len(parts) >= 4 {
				titlePart := strings.TrimSpace(parts[0])
				title := strings.TrimSpace(strings.Split(titlePart, ":")[1])
				usp := strings.TrimSpace(parts[1])
				keywords := parseListFromString(parts[2])
				
				var similarity float64
				fmt.Sscanf(strings.TrimSpace(parts[3]), "%f", &similarity)
				
				competitor := CompetitiveAnalysis{
					CompetitorTitle:         title,
					CompetitorUSP:           usp,
					DifferentiatingKeywords: keywords,
					SimilarityScore:         similarity,
				}
				competitors = append(competitors, competitor)
			}
		}
	}
	
	// Add default competitors if parsing failed
	if len(competitors) == 0 {
		competitors = []CompetitiveAnalysis{
			{
				CompetitorTitle:         "Emotional Intelligence 2.0",
				CompetitorUSP:           "Practical strategies for developing emotional intelligence",
				DifferentiatingKeywords: []string{"emotional intelligence", "practical strategies", "leadership"},
				SimilarityScore:         0.6,
			},
			{
				CompetitorTitle:         "The Culture Code",
				CompetitorUSP:           "Secrets of highly successful groups",
				DifferentiatingKeywords: []string{"culture", "teamwork", "success"},
				SimilarityScore:         0.5,
			},
			{
				CompetitorTitle:         "Good to Great",
				CompetitorUSP:           "Why some companies make the leap... and others don't",
				DifferentiatingKeywords: []string{"business transformation", "leadership", "performance"},
				SimilarityScore:         0.4,
			},
		}
	}
	
	return competitors, nil
}

// USP Evaluation with Scoring
func evaluateUSP(selectedModel string, usp USPVariant, competitors []CompetitiveAnalysis) (USPVariant, error) {
	prompt := fmt.Sprintf(`Rate this book concept: "%s"

Score from 0.0 to 1.0 for each category:

Novelty: How unique and innovative is this concept?
Reader Appeal: How engaging will this be for business readers?
Differentiation: How different is this from existing business books?

Respond EXACTLY in this format:
Novelty: 0.75
Reader Appeal: 0.80
Differentiation: 0.70
Analysis: Brief explanation of the scores`, usp.USP)

	response, err := callOllama(selectedModel, prompt)
	if err != nil {
		return usp, err
	}

	// Debug output to see what LLM is returning
	fmt.Printf("   🔍 DEBUG - LLM Response:\n%s\n", response)
	
	// Parse scores from response with robust parsing
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		lineLower := strings.ToLower(line)
		
		// Look for novelty score
		if strings.Contains(lineLower, "novelty") && strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				scoreStr := strings.TrimSpace(parts[1])
				extractedScore := extractNumericScore(scoreStr)
				fmt.Printf("   🔍 Novelty parsing: '%s' -> '%s'\n", scoreStr, extractedScore)
				if score := parseFloat(extractedScore); score > 0 {
					usp.NoveltyScore = score
					fmt.Printf("   ✅ Novelty score parsed: %.3f\n", score)
				}
			}
		}
		
		// Look for reader appeal score
		if (strings.Contains(lineLower, "reader appeal") || strings.Contains(lineLower, "appeal")) && strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				scoreStr := strings.TrimSpace(parts[1])
				extractedScore := extractNumericScore(scoreStr)
				fmt.Printf("   🔍 Appeal parsing: '%s' -> '%s'\n", scoreStr, extractedScore)
				if score := parseFloat(extractedScore); score > 0 {
					usp.ReaderAppealScore = score
					fmt.Printf("   ✅ Appeal score parsed: %.3f\n", score)
				}
			}
		}
		
		// Look for differentiation score
		if strings.Contains(lineLower, "differentiation") && strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				scoreStr := strings.TrimSpace(parts[1])
				extractedScore := extractNumericScore(scoreStr)
				fmt.Printf("   🔍 Differentiation parsing: '%s' -> '%s'\n", scoreStr, extractedScore)
				if score := parseFloat(extractedScore); score > 0 {
					usp.DifferentiationScore = score
					fmt.Printf("   ✅ Differentiation score parsed: %.3f\n", score)
				}
			}
		}
		
		// Look for analysis
		if strings.Contains(lineLower, "analysis") && strings.Contains(line, ":") {
			parts := strings.Split(line, ":")
			if len(parts) >= 2 {
				usp.AnalysisDetails = strings.TrimSpace(strings.Join(parts[1:], ":"))
			}
		}
	}
	
	// If parsing failed, report error but don't use fake numbers
	if usp.NoveltyScore == 0 && usp.ReaderAppealScore == 0 && usp.DifferentiationScore == 0 {
		fmt.Printf("   ❌ ERROR: Could not parse any scores from LLM response\n")
		fmt.Printf("   🔍 Raw response length: %d characters\n", len(response))
		return usp, fmt.Errorf("failed to parse scores from LLM response")
	}
	
	// Calculate overall score and check thresholds
	usp.OverallScore = (usp.NoveltyScore + usp.ReaderAppealScore + usp.DifferentiationScore) / 3
	usp.PassesThresholds = usp.NoveltyScore >= 0.60 && usp.ReaderAppealScore >= 0.75 && usp.DifferentiationScore >= 0.70
	
	fmt.Printf("   📊 Scores - Novelty: %.3f, Appeal: %.3f, Differentiation: %.3f\n", 
		usp.NoveltyScore, usp.ReaderAppealScore, usp.DifferentiationScore)
	
	return usp, nil
}

// Recursive USP Optimization Loop
func optimizeUSPRecursively(selectedModel, originalConcept string, competitors []CompetitiveAnalysis, variants []USPVariant) (USPVariant, []USPVariant, error) {
	maxIterations := 3
	currentIteration := 1
	
	// Check if initial USP already passes thresholds
	if len(variants) > 0 && variants[0].PassesThresholds {
		fmt.Printf("   ✅ Initial USP already meets all thresholds!\n")
		return variants[0], variants, nil
	}
	
	for currentIteration <= maxIterations {
		fmt.Printf("\n   🔄 Iteration %d/%d: Generating USP variants...\n", currentIteration, maxIterations)
		
		// Generate 5 new USP variants
		newVariants, err := generateUSPVariants(selectedModel, originalConcept, competitors, currentIteration)
		if err != nil {
			return USPVariant{}, variants, err
		}
		
		// Evaluate each variant
		for i, variant := range newVariants {
			fmt.Printf("      📝 Evaluating variant %d...\n", i+1)
			evaluatedVariant, err := evaluateUSP(selectedModel, variant, competitors)
			if err != nil {
				continue // Skip failed evaluation
			}
			
			variants = append(variants, evaluatedVariant)
			
			// Check if this variant passes all thresholds
			if evaluatedVariant.PassesThresholds {
				fmt.Printf("   ✅ Found optimal USP in iteration %d!\n", currentIteration)
				return evaluatedVariant, variants, nil
			}
		}
		
		currentIteration++
	}
	
	// If no variant passes all thresholds, return the best one
	bestVariant := findBestUSPVariant(variants)
	fmt.Printf("   ⚠️ No variant met all thresholds. Using best variant (Overall: %.3f)\n", bestVariant.OverallScore)
	
	return bestVariant, variants, nil
}

// Generate 5 USP variants based on current analysis
func generateUSPVariants(selectedModel, originalConcept string, competitors []CompetitiveAnalysis, iteration int) ([]USPVariant, error) {
	prompt := fmt.Sprintf(`Generate 5 distinct USP variations for the concept: "%s" (Iteration %d)

	Based on competitive analysis of:
	%s

	Create variations that prioritize:
	1. Higher conceptual novelty (>0.60)
	2. Stronger reader appeal (>0.75)
	3. Better differentiation from competitors (>0.70)

	Generate variations that explore different angles:
	- Emotional benefits
	- Practical applications
	- Unique methodologies
	- Target audience focus
	- Problem-solution framing

	Format your response as:
	VARIANT_1: [USP text]
	VARIANT_2: [USP text]
	VARIANT_3: [USP text]
	VARIANT_4: [USP text]
	VARIANT_5: [USP text]`, originalConcept, iteration, formatCompetitorsForPrompt(competitors))

	response, err := callOllama(selectedModel, prompt)
	if err != nil {
		return nil, err
	}

	var variants []USPVariant
	lines := strings.Split(response, "\n")
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "VARIANT_") {
			uspText := strings.TrimSpace(strings.Split(line, ":")[1])
			if uspText != "" {
				variant := USPVariant{
					USP:       uspText,
					Iteration: iteration,
				}
				variants = append(variants, variant)
			}
		}
	}
	
	// Generate default variants if parsing failed
	if len(variants) == 0 {
		variants = []USPVariant{
			{USP: originalConcept + ": A Revolutionary Approach", Iteration: iteration},
			{USP: originalConcept + ": Transforming Business Through Human Connection", Iteration: iteration},
			{USP: originalConcept + ": The Future of Sustainable Business", Iteration: iteration},
			{USP: originalConcept + ": Evidence-Based Strategies for Modern Leaders", Iteration: iteration},
			{USP: originalConcept + ": Building Competitive Advantage Through Empathy", Iteration: iteration},
		}
	}
	
	return variants, nil
}

// Find the best USP variant based on overall score
func findBestUSPVariant(variants []USPVariant) USPVariant {
	if len(variants) == 0 {
		return USPVariant{}
	}
	
	best := variants[0]
	for _, variant := range variants {
		if variant.OverallScore > best.OverallScore {
			best = variant
		}
	}
	
	return best
}

// Generate final selection justification
func generateSelectionJustification(selectedModel string, finalUSP USPVariant, allVariants []USPVariant) (string, error) {
	prompt := fmt.Sprintf(`Generate a comprehensive justification for selecting the final USP: "%s"

	Compared to %d other variants considered, explain why this USP is the optimal choice:
	
	Scores:
	- Novelty Score: %.3f (threshold: 0.60)
	- Reader Appeal Score: %.3f (threshold: 0.75)
	- Differentiation Score: %.3f (threshold: 0.70)
	- Overall Score: %.3f

	Provide a detailed justification covering:
	1. How it meets/exceeds the scoring thresholds
	2. Its competitive advantages
	3. Market positioning benefits
	4. Target audience appeal
	5. Implementation feasibility

	Format as a professional business justification.`, 
		finalUSP.USP, len(allVariants), 
		finalUSP.NoveltyScore, finalUSP.ReaderAppealScore, finalUSP.DifferentiationScore, finalUSP.OverallScore)

	response, err := callOllama(selectedModel, prompt)
	if err != nil {
		return fmt.Sprintf("Selected USP '%s' with overall score %.3f as the best variant from %d options evaluated.", 
			finalUSP.USP, finalUSP.OverallScore, len(allVariants)), nil
	}

	return response, nil
}

// Helper functions
func parseListFromString(str string) []string {
	str = strings.TrimSpace(str)
	if str == "" {
		return []string{}
	}
	
	// Remove brackets and split by comma
	str = strings.Trim(str, "[]")
	items := strings.Split(str, ",")
	
	var result []string
	for _, item := range items {
		item = strings.TrimSpace(item)
		if item != "" {
			result = append(result, item)
		}
	}
	
	return result
}

func formatCompetitorsForPrompt(competitors []CompetitiveAnalysis) string {
	var result strings.Builder
	for i, comp := range competitors {
		result.WriteString(fmt.Sprintf("\n%d. %s: %s", i+1, comp.CompetitorTitle, comp.CompetitorUSP))
	}
	return result.String()
}

// Extract numeric score from text (e.g., "0.75" from "0.75 (high novelty)")
func extractNumericScore(text string) string {
	// Remove common prefixes and suffixes
	text = strings.TrimSpace(text)
	text = strings.ReplaceAll(text, "[", "")
	text = strings.ReplaceAll(text, "]", "")
	text = strings.ReplaceAll(text, "(", "")
	text = strings.ReplaceAll(text, ")", "")
	
	// Look for decimal number pattern (0.xx)
	start := -1
	end := -1
	dotFound := false
	
	for i, char := range text {
		if char >= '0' && char <= '9' {
			if start == -1 {
				start = i
			}
			end = i + 1
		} else if char == '.' && start != -1 && !dotFound {
			dotFound = true
			end = i + 1
		} else if start != -1 {
			// Found end of number
			break
		}
	}
	
	if start != -1 && end > start {
		return text[start:end]
	}
	
	return ""
}

// Parse float with error handling
func parseFloat(text string) float64 {
	if text == "" {
		return 0
	}
	
	var result float64
	n, err := fmt.Sscanf(text, "%f", &result)
	if err != nil || n == 0 {
		// Try parsing with strconv as backup
		if val, err2 := fmt.Sscanf(strings.Fields(text)[0], "%f", &result); err2 == nil && val > 0 {
			// Successfully parsed
		} else {
			return 0
		}
	}
	
	// Ensure value is in valid range
	if result < 0 {
		return 0
	}
	if result > 1 {
		return 1
	}
	
	return result
}

// Generate quantitative selection criteria rubric
func generateSelectionCriteria(selectedModel string) (CriteriaRubric, error) {
	fmt.Printf("🎯 Generating quantitative selection criteria for optimal book selection...\n")
	
	prompt := `As a market research and publishing expert, create a comprehensive quantitative rubric for selecting the most profitable and successful book concept to write.

TASK: Generate 8-12 specific, measurable criteria that can be used to evaluate book concepts objectively.

Each criterion should have:
- A clear name
- Detailed description of what it measures
- Weight (importance factor 0.1-1.0, total should sum to ~1.0)
- Minimum threshold score (0.0-1.0)
- Category (Market, Content, Author, Competition)

Focus on criteria that predict:
- Commercial success and profitability
- Market demand and reader appeal
- Competitive advantage and differentiation
- Author credibility and platform strength
- Content quality and uniqueness

RESPONSE FORMAT:
CRITERIA_1: [Name] | [Description] | [Weight] | [Threshold] | [Category]
CRITERIA_2: [Name] | [Description] | [Weight] | [Threshold] | [Category]
...continue for 8-12 criteria...

PURPOSE: Book concept optimization for maximum market success
TARGET: Business professionals and thought leaders seeking practical insights`

	response, err := callOllama(selectedModel, prompt)
	if err != nil {
		fmt.Printf("⚠️  Error generating criteria, using default rubric: %v\n", err)
		return generateBasicRubric(), nil
	}

	// Parse response into criteria rubric
	rubric := CriteriaRubric{
		Timestamp:     time.Now().Format(time.RFC3339),
		Purpose:       "Book concept optimization for maximum market success",
		TargetMarket:  "Business professionals and thought leaders",
		Criteria:      []SelectionCriteria{},
	}

	lines := strings.Split(response, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "CRITERIA_") {
			parts := strings.Split(line, "|")
			if len(parts) >= 5 {
				namePart := strings.TrimSpace(parts[0])
				name := strings.TrimSpace(strings.Split(namePart, ":")[1])
				description := strings.TrimSpace(parts[1])
				
				var weight, threshold float64
				fmt.Sscanf(strings.TrimSpace(parts[2]), "%f", &weight)
				fmt.Sscanf(strings.TrimSpace(parts[3]), "%f", &threshold)
				category := strings.TrimSpace(parts[4])

				criteria := SelectionCriteria{
					Name:        name,
					Description: description,
					Weight:      weight,
					Threshold:   threshold,
					Category:    category,
				}
				rubric.Criteria = append(rubric.Criteria, criteria)
			}
		}
	}

	// Fallback to basic rubric if parsing failed
	if len(rubric.Criteria) == 0 {
		fmt.Printf("⚠️  Could not parse criteria response, using default rubric\n")
		return generateBasicRubric(), nil
	}

	rubric.TotalCriteria = len(rubric.Criteria)
	
	fmt.Printf("✅ Generated %d quantitative selection criteria\n", len(rubric.Criteria))
	fmt.Printf("📋 Criteria categories: Market, Content, Author, Competition\n")
	
	return rubric, nil
}

// Generate concepts based on selection criteria
func generateConceptsWithCriteria(selectedModel string, rubric CriteriaRubric) (string, error) {
	fmt.Printf("🎯 Generating book concepts optimized for the selection criteria...\n")
	
	// Format criteria for prompt
	criteriaList := ""
	for i, criteria := range rubric.Criteria {
		criteriaList += fmt.Sprintf("%d. %s (Weight: %.2f, Min: %.2f)\n   %s\n\n", 
			i+1, criteria.Name, criteria.Weight, criteria.Threshold, criteria.Description)
	}
	
	prompt := fmt.Sprintf(`Based on the following quantitative selection criteria, generate exactly 5 book concept titles that are specifically optimized to score highly on these metrics:

SELECTION CRITERIA:
%s

TARGET MARKET: %s

REQUIREMENTS:
- Each concept must be designed to excel in the highest-weighted criteria
- Focus on concepts that can realistically achieve threshold scores
- Prioritize commercial viability and market appeal
- Consider author platform requirements and feasibility
- Ensure concepts are innovative yet market-proven

RESPONSE FORMAT:
CONCEPT_1: [Title] | [Brief description optimized for criteria] | [Primary strength category]
CONCEPT_2: [Title] | [Brief description optimized for criteria] | [Primary strength category]
CONCEPT_3: [Title] | [Brief description optimized for criteria] | [Primary strength category]
CONCEPT_4: [Title] | [Brief description optimized for criteria] | [Primary strength category]
CONCEPT_5: [Title] | [Brief description optimized for criteria] | [Primary strength category]

Generate concepts now:`, criteriaList, rubric.TargetMarket)

	response, err := callOllama(selectedModel, prompt)
	if err != nil {
		return "", fmt.Errorf("error generating concepts: %v", err)
	}

	// Parse and display concepts
	suggestions := parseConceptSuggestions(response)
	if len(suggestions) == 0 {
		return "", fmt.Errorf("could not parse concept suggestions from LLM response")
	}

	// Load concept log for tracking
	conceptLog := loadConceptLog()
	
	return displayAndSelectConcept(suggestions, conceptLog)
}

// Generate optimal author persona for the selected concept
func generateAuthorPersona(selectedModel string, concept string, rubric CriteriaRubric) (AuthorPersona, error) {
	fmt.Printf("👤 Generating optimal author persona for '%s'...\n", concept)
	
	// Extract author-related criteria
	authorCriteria := ""
	for _, criteria := range rubric.Criteria {
		if criteria.Category == "Author" {
			authorCriteria += fmt.Sprintf("- %s: %s\n", criteria.Name, criteria.Description)
		}
	}
	
	prompt := fmt.Sprintf(`Create an optimal author persona for the book concept: "%s"

Based on these author-related criteria:
%s

TARGET MARKET: %s

Generate a fictional but realistic author profile that would maximize credibility, market appeal, and sales potential for this specific book concept.

Consider:
1. Professional background and expertise relevant to the topic
2. Educational credentials that support authority
3. Industry experience and career achievements
4. Public platform and media presence potential
5. Demographic factors that appeal to target market
6. Unique value proposition as an author

RESPONSE FORMAT:
NAME: [Full professional name]
BACKGROUND: [2-3 sentence professional background]
CREDENTIALS: [List 3-5 relevant credentials/achievements]
MARKET_APPEAL: [Score 0.0-1.0 with justification]
GENRE_MATCH: [Score 0.0-1.0 with justification]
TRUSTABILITY: [Score 0.0-1.0 with justification]
BIOGRAPHY: [2-3 paragraph professional biography suitable for book jacket]

Create author persona now:`, concept, authorCriteria, rubric.TargetMarket)

	response, err := callOllama(selectedModel, prompt)
	if err != nil {
		return generateDefaultAuthor(concept), nil
	}

	// Parse response
	persona := AuthorPersona{}
	lines := strings.Split(response, "\n")
	
	var bioLines []string
	inBio := false
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		
		if strings.HasPrefix(line, "NAME:") {
			persona.Name = strings.TrimSpace(strings.TrimPrefix(line, "NAME:"))
		} else if strings.HasPrefix(line, "BACKGROUND:") {
			persona.Background = strings.TrimSpace(strings.TrimPrefix(line, "BACKGROUND:"))
		} else if strings.HasPrefix(line, "CREDENTIALS:") {
			credStr := strings.TrimSpace(strings.TrimPrefix(line, "CREDENTIALS:"))
			persona.Credentials = parseListFromString(credStr)
		} else if strings.HasPrefix(line, "MARKET_APPEAL:") {
			scoreStr := strings.TrimSpace(strings.TrimPrefix(line, "MARKET_APPEAL:"))
			fmt.Sscanf(scoreStr, "%f", &persona.MarketAppeal)
		} else if strings.HasPrefix(line, "GENRE_MATCH:") {
			scoreStr := strings.TrimSpace(strings.TrimPrefix(line, "GENRE_MATCH:"))
			fmt.Sscanf(scoreStr, "%f", &persona.GenreMatch)
		} else if strings.HasPrefix(line, "TRUSTABILITY:") {
			scoreStr := strings.TrimSpace(strings.TrimPrefix(line, "TRUSTABILITY:"))
			fmt.Sscanf(scoreStr, "%f", &persona.Trustability)
		} else if strings.HasPrefix(line, "BIOGRAPHY:") {
			inBio = true
			bioText := strings.TrimSpace(strings.TrimPrefix(line, "BIOGRAPHY:"))
			if bioText != "" {
				bioLines = append(bioLines, bioText)
			}
		} else if inBio && line != "" {
			bioLines = append(bioLines, line)
		}
	}
	
	persona.Biography = strings.Join(bioLines, " ")
	
	// Set defaults if parsing failed
	if persona.Name == "" {
		return generateDefaultAuthor(concept), nil
	}
	
	fmt.Printf("✅ Generated author persona: %s\n", persona.Name)
	fmt.Printf("📊 Scores - Market Appeal: %.2f, Genre Match: %.2f, Trust: %.2f\n", 
		persona.MarketAppeal, persona.GenreMatch, persona.Trustability)
	
	return persona, nil
}

// Generate basic fallback rubric
func generateBasicRubric() CriteriaRubric {
	return CriteriaRubric{
		Timestamp:     time.Now().Format(time.RFC3339),
		TotalCriteria: 8,
		Purpose:       "Basic book concept evaluation framework",
		TargetMarket:  "Business professionals and thought leaders",
		Criteria: []SelectionCriteria{
			{
				Name:        "Market Demand",
				Description: "Current market demand and trending interest in topic",
				Weight:      0.20,
				Threshold:   0.70,
				Category:    "Market",
			},
			{
				Name:        "Commercial Viability",
				Description: "Potential for commercial success and profitability",
				Weight:      0.18,
				Threshold:   0.65,
				Category:    "Market",
			},
			{
				Name:        "Content Uniqueness",
				Description: "Novelty and differentiation from existing content",
				Weight:      0.15,
				Threshold:   0.60,
				Category:    "Content",
			},
			{
				Name:        "Reader Appeal",
				Description: "Emotional resonance and engagement potential",
				Weight:      0.15,
				Threshold:   0.75,
				Category:    "Content",
			},
			{
				Name:        "Author Credibility",
				Description: "Author expertise and authority in subject matter",
				Weight:      0.12,
				Threshold:   0.70,
				Category:    "Author",
			},
			{
				Name:        "Platform Potential",
				Description: "Author's ability to promote and market the book",
				Weight:      0.10,
				Threshold:   0.60,
				Category:    "Author",
			},
			{
				Name:        "Competitive Advantage",
				Description: "Differentiation from competing titles",
				Weight:      0.08,
				Threshold:   0.70,
				Category:    "Competition",
			},
			{
				Name:        "Implementation Feasibility",
				Description: "Practical feasibility of writing and producing the book",
				Weight:      0.02,
				Threshold:   0.80,
				Category:    "Content",
			},
		},
	}
}

// Generate default author for fallback
func generateDefaultAuthor(concept string) AuthorPersona {
	return AuthorPersona{
		Name:          "Dr. Alexandra Mitchell",
		Background:    "Business strategist and organizational development expert with 15+ years of consulting experience.",
		Credentials:   []string{"PhD Business Strategy", "Former McKinsey Partner", "Harvard Business Review Contributor", "TEDx Speaker"},
		MarketAppeal:  0.75,
		GenreMatch:    0.80,
		Trustability:  0.85,
		Biography:     "Dr. Alexandra Mitchell is a renowned business strategist and organizational development expert who has spent over 15 years helping Fortune 500 companies navigate complex transformations. As a former partner at McKinsey & Company, she has worked with leaders across industries to drive sustainable growth and innovation. Dr. Mitchell holds a PhD in Business Strategy from Wharton and has authored numerous articles for Harvard Business Review. Her insights on leadership and organizational change have been featured in major business publications, and she is a sought-after speaker at international conferences.",
	}
}

// Generate multiple concept suggestions for market analysis
func generateMultipleConceptSuggestions(selectedModel string) ([]ConceptSuggestion, error) {
	fmt.Printf("🤖 Generating multiple concept options for market analysis...\n")
	
	// Load previous concepts to avoid repetition
	conceptLog := loadConceptLog()
	
	// Create exclusion list from previous concepts
	exclusionList := ""
	if len(conceptLog.PreviousConcepts) > 0 {
		exclusionList = "\n\nIMPORTANT: EXCLUDE these previously generated concepts to ensure uniqueness:\n"
		for _, prev := range conceptLog.PreviousConcepts {
			exclusionList += fmt.Sprintf("- %s\n", prev.Title)
		}
		exclusionList += "\nGenerate COMPLETELY DIFFERENT concepts than those listed above.\n"
	}
	
	// Enhanced prompt for better LLM response
	currentYear := time.Now().Year()
	prompt := fmt.Sprintf(`You are an innovative business content strategist. Generate exactly 5 unique, cutting-edge book/content concept titles for %d-%d that will be analyzed for market intelligence.

FOCUS AREAS:
- Emerging technologies (AI, quantum computing, biotech, web3)
- Business transformation and leadership evolution
- Future-of-work and organizational design
- Sustainability, climate tech, and circular economy
- Psychology, neuroscience, and human performance
- Cross-industry innovation and convergence

REQUIREMENTS:
- Each concept must be highly novel and distinctive
- Target: business professionals, executives, and thought leaders
- Make titles compelling, marketable, and memorable
- Include trending keywords and emerging themes
- Avoid generic or overused business concepts
- Focus on concepts that can be analyzed for market demand

%s

RESPONSE FORMAT (follow exactly):
CONCEPT_1: [Unique Title] | [Compelling 10-15 word description] | [Category]
CONCEPT_2: [Unique Title] | [Compelling 10-15 word description] | [Category]
CONCEPT_3: [Unique Title] | [Compelling 10-15 word description] | [Category]
CONCEPT_4: [Unique Title] | [Compelling 10-15 word description] | [Category]
CONCEPT_5: [Unique Title] | [Compelling 10-15 word description] | [Category]

Generate fresh, innovative concepts now:`, currentYear, currentYear+1, exclusionList)
	
	response, err := callOllama(selectedModel, prompt)
	if err != nil {
		fmt.Printf("⚠️  LLM generation failed: %v\n", err)
		fmt.Printf("🔄 Using fallback suggestions...\n")
		return getFallbackSuggestions(), nil
	}
	
	// Parse LLM response
	suggestions := parseConceptSuggestions(response)
	if len(suggestions) == 0 {
		fmt.Printf("⚠️  Could not parse LLM response, using fallback suggestions...\n")
		suggestions = getFallbackSuggestions()
	} else {
		fmt.Printf("✅ Successfully generated %d concepts for market analysis\n", len(suggestions))
	}
	
	// Display the generated concepts
	fmt.Printf("\n📋 Generated Concepts for Market Analysis:\n")
	fmt.Printf("════════════════════════════════════════════\n")
	for i, suggestion := range suggestions {
		fmt.Printf("%d. %s\n", i+1, suggestion.Title)
		fmt.Printf("   📝 %s\n", suggestion.Description)
		fmt.Printf("   🏷️  Category: %s\n\n", suggestion.Category)
	}
	
	return suggestions, nil
}

// Gather market intelligence for multiple concepts
func gatherMarketIntelligenceForConcepts(selectedModel string, concepts []ConceptSuggestion) ([]MarketIntelligence, error) {
	var marketData []MarketIntelligence
	
	fmt.Printf("📊 Analyzing market intelligence for %d concepts...\n", len(concepts))
	
	for i, concept := range concepts {
		fmt.Printf("   🔍 Analyzing concept %d: %s\n", i+1, concept.Title)
		
		intel, err := gatherMarketIntelligence(selectedModel, concept.Title)
		if err != nil {
			fmt.Printf("   ⚠️  Warning: Could not analyze concept %d, using defaults\n", i+1)
			// Use basic default intelligence if analysis fails
			intel = MarketIntelligence{
				TrendingTopics:        []string{"Digital transformation", "Business innovation"},
				MarketGaps:            []string{"Practical implementation", "Measurable results"},
				TargetAudience:        "Business professionals and executives",
				MarketSaturationLevel: 0.5,
				PotentialReach:        50000,
			}
		}
		
		marketData = append(marketData, intel)
		
		// Small delay to be respectful to the API
		if i < len(concepts)-1 {
			fmt.Printf("   ⏳ Processing next concept...\n")
		}
	}
	
	fmt.Printf("✅ Market intelligence gathered for all %d concepts\n", len(concepts))
	return marketData, nil
}

// Generate market intelligence-driven criteria
func generateMarketIntelligenceDrivenCriteria(selectedModel string, marketData []MarketIntelligence) (CriteriaRubric, error) {
	fmt.Printf("🎯 Creating selection criteria based on market intelligence...\n")
	
	// Aggregate market intelligence insights
	allTopics := []string{}
	allGaps := []string{}
	avgSaturation := 0.0
	totalReach := 0
	
	for _, data := range marketData {
		allTopics = append(allTopics, data.TrendingTopics...)
		allGaps = append(allGaps, data.MarketGaps...)
		avgSaturation += data.MarketSaturationLevel
		totalReach += data.PotentialReach
	}
	
	if len(marketData) > 0 {
		avgSaturation /= float64(len(marketData))
		totalReach /= len(marketData)
	}
	
	// Create comprehensive prompt with market insights
	prompt := fmt.Sprintf(`Based on market intelligence analysis of multiple book concepts, create a comprehensive quantitative rubric for selecting the optimal book concept.

MARKET INTELLIGENCE INSIGHTS:
- Trending Topics Identified: %v
- Market Gaps Discovered: %v
- Average Market Saturation: %.2f
- Average Potential Reach: %d

TASK: Generate 10-12 specific, measurable criteria that leverage these market insights to predict commercial success.

Each criterion should be based on the market intelligence data and include:
- A clear name reflecting market realities
- Detailed description incorporating market insights
- Weight (importance factor 0.05-0.20, total should sum to ~1.0)
- Minimum threshold score (0.0-1.0) based on market conditions
- Category (Market, Content, Author, Competition)

Prioritize criteria that address:
- Market gaps and opportunities identified
- Trending topics alignment
- Competitive positioning in current market saturation
- Audience reach optimization
- Commercial viability in current market conditions
- Time-to-market considerations for rapid launch capability
- Strategic fit with brand strengths and content pipeline

RESPONSE FORMAT:
CRITERIA_1: [Name] | [Description with market insights] | [Weight] | [Threshold] | [Category]
CRITERIA_2: [Name] | [Description with market insights] | [Weight] | [Threshold] | [Category]
...continue for 10-12 criteria...

PURPOSE: Market intelligence-driven book concept optimization
TARGET: Business professionals and thought leaders seeking market-validated insights`, 
		allTopics, allGaps, avgSaturation, totalReach)

	response, err := callOllama(selectedModel, prompt)
	if err != nil {
		fmt.Printf("⚠️  Error generating market-driven criteria, using enhanced default rubric: %v\n", err)
		return generateEnhancedRubricFromMarketData(marketData), nil
	}

	// Parse response into criteria rubric
	rubric := CriteriaRubric{
		Timestamp:     time.Now().Format(time.RFC3339),
		Purpose:       "Market intelligence-driven book concept optimization",
		TargetMarket:  "Business professionals and thought leaders",
		Criteria:      []SelectionCriteria{},
	}

	lines := strings.Split(response, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "CRITERIA_") {
			parts := strings.Split(line, "|")
			if len(parts) >= 5 {
				namePart := strings.TrimSpace(parts[0])
				name := strings.TrimSpace(strings.Split(namePart, ":")[1])
				description := strings.TrimSpace(parts[1])
				
				var weight, threshold float64
				fmt.Sscanf(strings.TrimSpace(parts[2]), "%f", &weight)
				fmt.Sscanf(strings.TrimSpace(parts[3]), "%f", &threshold)
				category := strings.TrimSpace(parts[4])

				criteria := SelectionCriteria{
					Name:        name,
					Description: description,
					Weight:      weight,
					Threshold:   threshold,
					Category:    category,
				}
				rubric.Criteria = append(rubric.Criteria, criteria)
			}
		}
	}

	// Fallback to enhanced rubric if parsing failed
	if len(rubric.Criteria) == 0 {
		fmt.Printf("⚠️  Could not parse criteria response, using enhanced market-based rubric\n")
		return generateEnhancedRubricFromMarketData(marketData), nil
	}

	rubric.TotalCriteria = len(rubric.Criteria)
	
	fmt.Printf("✅ Generated %d market intelligence-driven criteria\n", len(rubric.Criteria))
	fmt.Printf("📊 Criteria based on market saturation: %.2f, potential reach: %d\n", avgSaturation, totalReach)
	
	return rubric, nil
}

// Select optimal concept based on comprehensive criteria scoring
func selectOptimalConcept(concepts []ConceptSuggestion, rubric CriteriaRubric, selectedModel string) (string, error) {
	fmt.Printf("📋 Comprehensive evaluation of %d concepts against %d criteria...\n", len(concepts), len(rubric.Criteria))
	
	// Score all concepts against all criteria
	evaluations, err := scoreAllConcepts(concepts, rubric, selectedModel)
	if err != nil {
		return "", fmt.Errorf("error scoring concepts: %v", err)
	}
	
	// Display comprehensive scoring results
	displayConceptScores(evaluations, rubric)
	
	// Select concept (auto-select highest scoring or allow user override)
	selectedConcept := selectFromScoredConcepts(evaluations)
	
	// Log the selected concept
	conceptLog := loadConceptLog()
	conceptLog.PreviousConcepts = append(conceptLog.PreviousConcepts, selectedConcept.Concept)
	conceptLog.LastGenerated = time.Now().Format(time.RFC3339)
	saveConceptLog(conceptLog)
	
	return selectedConcept.Concept.Title, nil
}

// Score all concepts against all criteria
func scoreAllConcepts(concepts []ConceptSuggestion, rubric CriteriaRubric, selectedModel string) ([]ConceptEvaluation, error) {
	var evaluations []ConceptEvaluation
	
	fmt.Printf("\n🎯 Scoring each concept against all %d criteria...\n", len(rubric.Criteria))
	
	for i, concept := range concepts {
		fmt.Printf("\n📊 Evaluating Concept %d: %s\n", i+1, concept.Title)
		
		evaluation := ConceptEvaluation{
			Concept: concept,
			Scores:  []ConceptScore{},
		}
		
		// Score against each criterion
		for _, criteria := range rubric.Criteria {
			score, err := scoreConcept(concept, criteria, selectedModel)
			if err != nil {
				fmt.Printf("   ⚠️  Warning: Could not score '%s', using 0.5 default\n", criteria.Name)
				score = 0.5 // Use neutral score if evaluation fails
			}
			
			weightedScore := score * criteria.Weight
			passes := score >= criteria.Threshold
			
			conceptScore := ConceptScore{
				CriteriaName:  criteria.Name,
				Score:         score,
				Weight:        criteria.Weight,
				WeightedScore: weightedScore,
				Threshold:     criteria.Threshold,
				Passes:        passes,
			}
			
			evaluation.Scores = append(evaluation.Scores, conceptScore)
			evaluation.TotalWeighted += weightedScore
			
			if passes {
				evaluation.PassingCount++
			}
			
			fmt.Printf("   %s %s: %.3f (weighted: %.3f)\n", 
				getPassFailIcon(passes), criteria.Name, score, weightedScore)
		}
		
		evaluations = append(evaluations, evaluation)
	}
	
	// Rank concepts by total weighted score
	for i := 0; i < len(evaluations); i++ {
		rank := 1
		for j := 0; j < len(evaluations); j++ {
			if evaluations[j].TotalWeighted > evaluations[i].TotalWeighted {
				rank++
			}
		}
		evaluations[i].Rank = rank
	}
	
	return evaluations, nil
}

// Score a single concept against a single criterion
func scoreConcept(concept ConceptSuggestion, criteria SelectionCriteria, selectedModel string) (float64, error) {
	// Create a simplified prompt for scoring
	prompt := fmt.Sprintf(`Score this book concept against the specific criterion:

Concept: "%s"
Description: %s
Category: %s

Criterion: %s
Description: %s

Rate from 0.0 to 1.0 how well this concept performs on this criterion.
Consider both the concept's inherent qualities and market context.

Respond with ONLY a number between 0.0 and 1.0:`, 
		concept.Title, concept.Description, concept.Category, 
		criteria.Name, criteria.Description)

	// Use the selected model for scoring
	response, err := callOllama(selectedModel, prompt)
	if err != nil {
		return 0, err
	}
	
	// Extract numeric score from response
	scoreStr := extractNumericScore(strings.TrimSpace(response))
	score := parseFloat(scoreStr)
	
	// If no valid score found, try to parse the first line
	if score == 0 {
		lines := strings.Split(response, "\n")
		for _, line := range lines {
			line = strings.TrimSpace(line)
			if line != "" {
				scoreStr = extractNumericScore(line)
				score = parseFloat(scoreStr)
				if score > 0 {
					break
				}
			}
		}
	}
	
	return score, nil
}

// Display comprehensive concept scores
func displayConceptScores(evaluations []ConceptEvaluation, rubric CriteriaRubric) {
	fmt.Printf("\n📈 COMPREHENSIVE CONCEPT EVALUATION RESULTS\n")
	fmt.Printf("══════════════════════════════════════════════════════════════════════\n")
	
	// Sort by rank for display
	sortedEvals := make([]ConceptEvaluation, len(evaluations))
	copy(sortedEvals, evaluations)
	
	for i := 0; i < len(sortedEvals)-1; i++ {
		for j := i + 1; j < len(sortedEvals); j++ {
			if sortedEvals[i].Rank > sortedEvals[j].Rank {
				sortedEvals[i], sortedEvals[j] = sortedEvals[j], sortedEvals[i]
			}
		}
	}
	
	for _, eval := range sortedEvals {
		fmt.Printf("\n🏆 RANK %d: %s\n", eval.Rank, eval.Concept.Title)
		fmt.Printf("📝 %s\n", eval.Concept.Description)
		fmt.Printf("🏷️  Category: %s\n", eval.Concept.Category)
		fmt.Printf("📊 Total Weighted Score: %.3f\n", eval.TotalWeighted)
		fmt.Printf("✅ Criteria Passed: %d/%d\n", eval.PassingCount, len(rubric.Criteria))
		
		// Show top 5 scoring criteria
		fmt.Printf("🎯 Top Criteria Performance:\n")
		sortedScores := make([]ConceptScore, len(eval.Scores))
		copy(sortedScores, eval.Scores)
		
		// Sort by weighted score
		for i := 0; i < len(sortedScores)-1; i++ {
			for j := i + 1; j < len(sortedScores); j++ {
				if sortedScores[i].WeightedScore < sortedScores[j].WeightedScore {
					sortedScores[i], sortedScores[j] = sortedScores[j], sortedScores[i]
				}
			}
		}
		
		for i, score := range sortedScores[:5] {
			fmt.Printf("   %d. %s %s: %.3f (%.3f weighted)\n", 
				i+1, getPassFailIcon(score.Passes), score.CriteriaName, score.Score, score.WeightedScore)
		}
		fmt.Printf("   ────────────────────────────────────────────────────\n")
	}
}

// Select from scored concepts
func selectFromScoredConcepts(evaluations []ConceptEvaluation) ConceptEvaluation {
	// Find highest scoring concept
	bestEval := evaluations[0]
	for _, eval := range evaluations {
		if eval.Rank == 1 {
			bestEval = eval
			break
		}
	}
	
	fmt.Printf("\n🎯 OPTIMAL CONCEPT SELECTION\n")
	fmt.Printf("════════════════════════════════════════════\n")
	fmt.Printf("📊 Algorithm recommends: %s\n", bestEval.Concept.Title)
	fmt.Printf("⭐ Total Score: %.3f | Criteria Passed: %d/%d\n", 
		bestEval.TotalWeighted, bestEval.PassingCount, len(bestEval.Scores))
	
	fmt.Printf("\nAccept recommendation? (Y/n): ")
	reader := bufio.NewReader(os.Stdin)
	choice, err := reader.ReadString('\n')
	if err != nil || strings.ToLower(strings.TrimSpace(choice)) != "n" {
		fmt.Printf("✅ Selected optimal concept: %s\n", bestEval.Concept.Title)
		return bestEval
	}
	
	// User wants to override, show menu
	fmt.Printf("\nSelect concept manually:\n")
	for i, eval := range evaluations {
		fmt.Printf("%d. %s (Score: %.3f, Rank: %d)\n", 
			i+1, eval.Concept.Title, eval.TotalWeighted, eval.Rank)
	}
	
	fmt.Printf("Choice (1-%d): ", len(evaluations))
	choice, err = reader.ReadString('\n')
	if err != nil {
		return bestEval
	}
	
	var choiceNum int
	if _, err := fmt.Sscanf(strings.TrimSpace(choice), "%d", &choiceNum); err == nil && 
		choiceNum >= 1 && choiceNum <= len(evaluations) {
		selectedEval := evaluations[choiceNum-1]
		fmt.Printf("✅ Selected concept: %s\n", selectedEval.Concept.Title)
		return selectedEval
	}
	
	fmt.Printf("Invalid choice, using recommended concept: %s\n", bestEval.Concept.Title)
	return bestEval
}

// Generate enhanced rubric from market data when AI parsing fails
func generateEnhancedRubricFromMarketData(marketData []MarketIntelligence) CriteriaRubric {
	// Calculate market insights
	avgSaturation := 0.0
	avgReach := 0
	if len(marketData) > 0 {
		for _, data := range marketData {
			avgSaturation += data.MarketSaturationLevel
			avgReach += data.PotentialReach
		}
		avgSaturation /= float64(len(marketData))
		avgReach /= len(marketData)
	}
	
	// Adjust thresholds based on market conditions
	marketThreshold := 0.70
	if avgSaturation > 0.7 {
		marketThreshold = 0.75 // Higher bar in saturated markets
	}
	
	return CriteriaRubric{
		Timestamp:     time.Now().Format(time.RFC3339),
		TotalCriteria: 12,
		Purpose:       "Enhanced market intelligence-driven book concept evaluation",
		TargetMarket:  "Business professionals and thought leaders",
		Criteria: []SelectionCriteria{
			{
				Name:        "Market Demand Validation",
				Description: fmt.Sprintf("Validated market demand based on trending topics and audience interest (avg reach: %d)", avgReach),
				Weight:      0.16,
				Threshold:   marketThreshold,
				Category:    "Market",
			},
			{
				Name:        "Market Gap Exploitation",
				Description: "Ability to fill identified market gaps and unmet needs",
				Weight:      0.14,
				Threshold:   0.70,
				Category:    "Market",
			},
			{
				Name:        "Commercial Viability",
				Description: fmt.Sprintf("Revenue potential in market with %.2f saturation level", avgSaturation),
				Weight:      0.13,
				Threshold:   0.65,
				Category:    "Market",
			},
			{
				Name:        "Strategic Fit",
				Description: "Alignment with brand strengths, expertise, and content pipeline strategy",
				Weight:      0.12,
				Threshold:   0.70,
				Category:    "Strategic",
			},
			{
				Name:        "Content Differentiation",
				Description: "Unique positioning against existing market offerings",
				Weight:      0.11,
				Threshold:   0.65,
				Category:    "Content",
			},
			{
				Name:        "Time-to-Market",
				Description: "Speed of development and launch capability for rapid market entry",
				Weight:      0.10,
				Threshold:   0.60,
				Category:    "Strategic",
			},
			{
				Name:        "Trending Topic Alignment",
				Description: "Alignment with identified trending topics and emerging themes",
				Weight:      0.09,
				Threshold:   0.60,
				Category:    "Content",
			},
			{
				Name:        "Reader Engagement Potential",
				Description: "Emotional resonance and engagement capacity with target audience",
				Weight:      0.08,
				Threshold:   0.75,
				Category:    "Content",
			},
			{
				Name:        "Author Platform Strength",
				Description: "Author's platform and ability to reach identified target audience",
				Weight:      0.07,
				Threshold:   0.60,
				Category:    "Author",
			},
			{
				Name:        "Expertise Credibility",
				Description: "Author expertise matching content requirements and audience expectations",
				Weight:      0.06,
				Threshold:   0.70,
				Category:    "Author",
			},
			{
				Name:        "Competitive Positioning",
				Description: "Strategic positioning against competitive landscape",
				Weight:      0.03,
				Threshold:   0.65,
				Category:    "Competition",
			},
			{
				Name:        "Implementation Feasibility",
				Description: "Practical feasibility of content creation and market delivery",
				Weight:      0.01,
				Threshold:   0.80,
				Category:    "Content",
			},
		},
	}
}