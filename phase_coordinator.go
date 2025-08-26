package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Phase coordinator manages data flow and transitions between phases
type PhaseCoordinator struct {
	BookOutline     *BookOutline           `json:"book_outline"`
	PhaseData       map[string]interface{} `json:"phase_data"`
	BookDirectory   string                 `json:"book_directory"`
	CurrentPhase    int                    `json:"current_phase"`
	PhaseHistory    []PhaseExecution       `json:"phase_history"`
}

// Phase execution record
type PhaseExecution struct {
	PhaseNumber   int                    `json:"phase_number"`
	PhaseName     string                 `json:"phase_name"`
	StartTime     time.Time              `json:"start_time"`
	EndTime       time.Time              `json:"end_time"`
	Duration      string                 `json:"duration"`
	InputData     map[string]interface{} `json:"input_data"`
	OutputData    map[string]interface{} `json:"output_data"`
	Status        string                 `json:"status"`
	OutputFiles   []string               `json:"output_files"`
	Improvements  []string               `json:"improvements"`
}

// Enhanced data structures for phase transitions
type USPOptimization struct {
	OptimizedUSP      string   `json:"optimized_usp"`
	MarketGapScore    float64  `json:"market_gap_score"`
	CompetitorAnalysis []string `json:"competitor_analysis"`
	MarketInsights    []string `json:"market_insights"`
	RecommendedFocus  []string `json:"recommended_focus"`
}

type ConceptValidation struct {
	ValidatedConcepts  []ValidatedConcept `json:"validated_concepts"`
	BestConcept        ValidatedConcept   `json:"best_concept"`
	ViabilityScore     float64           `json:"viability_score"`
	ReaderAppeal       float64           `json:"reader_appeal"`
	MarketPotential    float64           `json:"market_potential"`
}

type ValidatedConcept struct {
	Title           string  `json:"title"`
	Description     string  `json:"description"`
	OverallScore    float64 `json:"overall_score"`
	Strengths       []string `json:"strengths"`
	Improvements    []string `json:"improvements"`
}

type ReaderFeedback struct {
	ReaderPersonas    []ReaderPersona `json:"reader_personas"`
	EngagementScore   float64        `json:"engagement_score"`
	ShareabilityScore float64        `json:"shareability_score"`
	ViralQuotes       []string       `json:"viral_quotes"`
	FeedbackSummary   []string       `json:"feedback_summary"`
}

type ReaderPersona struct {
	Name           string   `json:"name"`
	Demographics   string   `json:"demographics"`
	Feedback       []string `json:"feedback"`
	Rating         float64  `json:"rating"`
	Recommendations []string `json:"recommendations"`
}

type MediaAnalysis struct {
	PRWorthiness      float64      `json:"pr_worthiness"`
	MediaOutlets      []MediaOutlet `json:"media_outlets"`
	CoverageEstimate  string       `json:"coverage_estimate"`
	PRRecommendations []string     `json:"pr_recommendations"`
	MediaAngles       []string     `json:"media_angles"`
}

type MediaOutlet struct {
	Name           string  `json:"name"`
	Type           string  `json:"type"`
	Audience       string  `json:"audience"`
	CoverageChance float64 `json:"coverage_chance"`
	Angle          string  `json:"angle"`
}

type TitleOptimization struct {
	OptimizedTitle    string         `json:"optimized_title"`
	TitleVariations   []TitleVariant `json:"title_variations"`
	ABTestResults     []ABTest       `json:"ab_test_results"`
	SEOScore          float64        `json:"seo_score"`
	MarketabilityScore float64       `json:"marketability_score"`
}

type TitleVariant struct {
	Title           string  `json:"title"`
	Score           float64 `json:"score"`
	Strengths       []string `json:"strengths"`
	TestResults     string  `json:"test_results"`
}

type ABTest struct {
	TitleA          string  `json:"title_a"`
	TitleB          string  `json:"title_b"`
	Winner          string  `json:"winner"`
	ConfidenceLevel float64 `json:"confidence_level"`
	Metrics         map[string]float64 `json:"metrics"`
}

// Initialize phase coordinator
func NewPhaseCoordinator(outline *BookOutline, bookDir string) *PhaseCoordinator {
	return &PhaseCoordinator{
		BookOutline:   outline,
		PhaseData:     make(map[string]interface{}),
		BookDirectory: bookDir,
		CurrentPhase:  0,
		PhaseHistory:  []PhaseExecution{},
	}
}

// Prepare input data for a phase based on previous results and outline
func (pc *PhaseCoordinator) PreparePhaseInput(phaseNumber int, phaseName string) map[string]interface{} {
	inputData := make(map[string]interface{})
	
	// Always include the book outline
	inputData["book_outline"] = pc.BookOutline
	inputData["book_title"] = pc.BookOutline.Title
	inputData["memorable_phrase"] = pc.BookOutline.MemorablePhrase
	inputData["category"] = pc.BookOutline.Category
	inputData["target_audience"] = pc.BookOutline.TargetAudience
	inputData["core_thesis"] = pc.BookOutline.CoreThesis
	
	// Add phase-specific input preparation
	switch phaseNumber {
	case 1: // Phase 1 Beta: Market Intelligence & USP Optimization
		inputData["content_type"] = "book"
		inputData["initial_usp"] = pc.BookOutline.MemorablePhrase
		inputData["key_concepts"] = pc.BookOutline.KeyConcepts
		inputData["supporting_points"] = pc.BookOutline.SupportingPoints
		
	case 2: // Phase 2: Concept Generation & Validation
		if uspData, exists := pc.PhaseData["phase1"]; exists {
			inputData["usp_optimization"] = uspData
		}
		inputData["chapter_outlines"] = pc.BookOutline.Chapters
		
	case 3: // Phase 3: Reader Feedback & Shareability
		if conceptData, exists := pc.PhaseData["phase2"]; exists {
			inputData["validated_concepts"] = conceptData
		}
		inputData["viral_potential_phrases"] = []string{pc.BookOutline.MemorablePhrase}
		
	case 4: // Phase 4: Media Coverage & PR Analysis
		if feedbackData, exists := pc.PhaseData["phase3"]; exists {
			inputData["reader_feedback"] = feedbackData
		}
		inputData["book_category"] = pc.BookOutline.Category
		
	case 5: // Phase 5: Title Optimization & A/B Testing
		if mediaData, exists := pc.PhaseData["phase4"]; exists {
			inputData["media_analysis"] = mediaData
		}
		inputData["current_title"] = pc.BookOutline.Title
		inputData["subtitle_options"] = []string{pc.BookOutline.Subtitle}
		
	case 6: // Phase 6 Enhanced: Complete Content Generation
		// Compile all previous insights
		inputData["usp_insights"] = pc.PhaseData["phase1"]
		inputData["concept_validation"] = pc.PhaseData["phase2"]
		inputData["reader_insights"] = pc.PhaseData["phase3"]
		inputData["media_insights"] = pc.PhaseData["phase4"]
		inputData["title_optimization"] = pc.PhaseData["phase5"]
		inputData["content_generation_mode"] = "outline_driven"
		
	case 7: // Phase 7: Marketing Assets & Campaign
		if contentData, exists := pc.PhaseData["phase6"]; exists {
			inputData["book_content"] = contentData
		}
		// Include all previous phase data for comprehensive marketing
		inputData["all_phase_insights"] = pc.PhaseData
	}
	
	return inputData
}

// Process phase output and extract improvements for next phases
func (pc *PhaseCoordinator) ProcessPhaseOutput(phaseNumber int, phaseName string, outputFile string) error {
	execution := PhaseExecution{
		PhaseNumber: phaseNumber,
		PhaseName:   phaseName,
		StartTime:   time.Now(),
		Status:      "processing",
	}
	
	// Read and parse phase output
	outputData, err := pc.readPhaseOutput(outputFile)
	if err != nil {
		execution.Status = "failed"
		execution.EndTime = time.Now()
		execution.Duration = execution.EndTime.Sub(execution.StartTime).String()
		pc.PhaseHistory = append(pc.PhaseHistory, execution)
		return fmt.Errorf("error reading phase output: %v", err)
	}
	
	// Extract phase-specific insights and improvements
	improvements := pc.extractPhaseImprovements(phaseNumber, outputData)
	
	// Store processed data
	pc.PhaseData[fmt.Sprintf("phase%d", phaseNumber)] = outputData
	pc.CurrentPhase = phaseNumber
	
	// Update execution record
	execution.EndTime = time.Now()
	execution.Duration = execution.EndTime.Sub(execution.StartTime).String()
	execution.OutputData = outputData
	execution.Status = "completed"
	execution.Improvements = improvements
	execution.OutputFiles = []string{outputFile}
	
	pc.PhaseHistory = append(pc.PhaseHistory, execution)
	
	// Save coordinator state
	pc.saveCoordinatorState()
	
	fmt.Printf("🔄 Phase %d output processed: %d improvements identified\n", phaseNumber, len(improvements))
	for i, improvement := range improvements {
		if i < 3 { // Show first 3 improvements
			fmt.Printf("   💡 %s\n", improvement)
		}
	}
	if len(improvements) > 3 {
		fmt.Printf("   ... and %d more\n", len(improvements)-3)
	}
	
	return nil
}

// Read phase output file
func (pc *PhaseCoordinator) readPhaseOutput(outputFile string) (map[string]interface{}, error) {
	var data map[string]interface{}
	
	// Try to read from book directory first
	fullPath := filepath.Join(pc.BookDirectory, outputFile)
	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		// Try current directory
		fullPath = outputFile
	}
	
	fileData, err := os.ReadFile(fullPath)
	if err != nil {
		return nil, err
	}
	
	if err := json.Unmarshal(fileData, &data); err != nil {
		return nil, err
	}
	
	return data, nil
}

// Extract improvements and insights from each phase
func (pc *PhaseCoordinator) extractPhaseImprovements(phaseNumber int, data map[string]interface{}) []string {
	var improvements []string
	
	switch phaseNumber {
	case 1: // USP Optimization improvements
		improvements = append(improvements, "Market positioning refined based on competitive analysis")
		improvements = append(improvements, "USP strengthened with market gap insights")
		if marketGap, ok := data["market_gap_score"].(float64); ok && marketGap > 7.0 {
			improvements = append(improvements, "High market gap identified - significant opportunity")
		}
		
	case 2: // Concept Validation improvements
		improvements = append(improvements, "Core concepts validated with target audience")
		improvements = append(improvements, "Book structure optimized for reader engagement")
		if viability, ok := data["viability_score"].(float64); ok && viability > 8.0 {
			improvements = append(improvements, "High viability score confirms strong concept")
		}
		
	case 3: // Reader Feedback improvements
		improvements = append(improvements, "Content optimized based on reader persona feedback")
		improvements = append(improvements, "Shareability elements identified and enhanced")
		if engagement, ok := data["engagement_score"].(float64); ok && engagement > 7.5 {
			improvements = append(improvements, "High engagement potential confirmed")
		}
		
	case 4: // Media Analysis improvements
		improvements = append(improvements, "PR angles identified for media coverage")
		improvements = append(improvements, "Media outreach strategy developed")
		if prScore, ok := data["pr_worthiness"].(float64); ok && prScore > 7.0 {
			improvements = append(improvements, "Strong PR potential identified")
		}
		
	case 5: // Title Optimization improvements
		improvements = append(improvements, "Title optimized for discoverability and appeal")
		improvements = append(improvements, "A/B testing insights incorporated")
		if seoScore, ok := data["seo_score"].(float64); ok && seoScore > 8.0 {
			improvements = append(improvements, "Excellent SEO optimization achieved")
		}
		
	case 6: // Content Generation improvements
		improvements = append(improvements, "Complete book content generated with all phase insights")
		improvements = append(improvements, "Multi-format output created (EPUB, TXT, MD, JSON)")
		if wordCount, ok := data["total_words"].(float64); ok && wordCount > 10000 {
			improvements = append(improvements, fmt.Sprintf("Comprehensive content: %.0f words", wordCount))
		}
		
	case 7: // Marketing improvements
		improvements = append(improvements, "Complete marketing campaign assets created")
		improvements = append(improvements, "Multi-platform marketing strategy developed")
		improvements = append(improvements, "Launch-ready promotional materials generated")
	}
	
	return improvements
}

// Generate summary of all phase improvements for the book
func (pc *PhaseCoordinator) GenerateBookImprovementSummary() map[string]interface{} {
	summary := map[string]interface{}{
		"book_title":        pc.BookOutline.Title,
		"total_phases":      len(pc.PhaseHistory),
		"total_improvements": 0,
		"phase_summaries":   []map[string]interface{}{},
		"key_achievements":  []string{},
		"final_metrics":     map[string]interface{}{},
	}
	
	totalImprovements := 0
	keyAchievements := []string{}
	
	for _, execution := range pc.PhaseHistory {
		phaseSummary := map[string]interface{}{
			"phase_number":   execution.PhaseNumber,
			"phase_name":     execution.PhaseName,
			"duration":       execution.Duration,
			"improvements":   len(execution.Improvements),
			"status":         execution.Status,
		}
		
		summary["phase_summaries"] = append(summary["phase_summaries"].([]map[string]interface{}), phaseSummary)
		totalImprovements += len(execution.Improvements)
		
		// Extract key achievements
		if execution.Status == "completed" {
			keyAchievements = append(keyAchievements, fmt.Sprintf("%s: %s", execution.PhaseName, "Successfully completed"))
		}
	}
	
	summary["total_improvements"] = totalImprovements
	summary["key_achievements"] = keyAchievements
	
	// Extract final metrics from last phases
	if phase5Data, exists := pc.PhaseData["phase5"]; exists {
		if data, ok := phase5Data.(map[string]interface{}); ok {
			if seoScore, ok := data["seo_score"]; ok {
				summary["final_metrics"].(map[string]interface{})["seo_score"] = seoScore
			}
		}
	}
	
	return summary
}

// Save coordinator state
func (pc *PhaseCoordinator) saveCoordinatorState() {
	stateFile := filepath.Join(pc.BookDirectory, "phase_coordinator_state.json")
	stateData, _ := json.MarshalIndent(pc, "", "  ")
	os.WriteFile(stateFile, stateData, 0644)
}

// Create enhanced input files for phases that support outline-driven generation
func (pc *PhaseCoordinator) CreateEnhancedInputFile(phaseNumber int, inputData map[string]interface{}) string {
	fileName := fmt.Sprintf("enhanced_input_phase%d.json", phaseNumber)
	filePath := filepath.Join(pc.BookDirectory, fileName)
	
	// Create comprehensive input structure
	enhancedInput := map[string]interface{}{
		"timestamp":      time.Now().Format(time.RFC3339),
		"phase_number":   phaseNumber,
		"book_outline":   pc.BookOutline,
		"input_data":     inputData,
		"previous_phases": pc.PhaseData,
		"context": map[string]interface{}{
			"generation_mode": "outline_driven",
			"target_quality":  "professional",
			"output_formats":  []string{"json", "markdown", "epub", "txt"},
		},
	}
	
	data, _ := json.MarshalIndent(enhancedInput, "", "  ")
	os.WriteFile(filePath, data, 0644)
	
	return fileName
}

// Validate phase transition readiness
func (pc *PhaseCoordinator) ValidatePhaseTransition(fromPhase, toPhase int) error {
	// Check if previous phase completed successfully
	if fromPhase > 0 {
		found := false
		for _, execution := range pc.PhaseHistory {
			if execution.PhaseNumber == fromPhase && execution.Status == "completed" {
				found = true
				break
			}
		}
		if !found {
			return fmt.Errorf("phase %d has not completed successfully", fromPhase)
		}
	}
	
	// Check if required data exists for the next phase
	requiredData := pc.getRequiredDataForPhase(toPhase)
	for _, requirement := range requiredData {
		if !pc.hasRequiredData(requirement) {
			return fmt.Errorf("required data missing for phase %d: %s", toPhase, requirement)
		}
	}
	
	return nil
}

// Get required data for each phase
func (pc *PhaseCoordinator) getRequiredDataForPhase(phaseNumber int) []string {
	switch phaseNumber {
	case 1:
		return []string{"book_outline"}
	case 2:
		return []string{"book_outline", "phase1_output"}
	case 3:
		return []string{"book_outline", "phase2_output"}
	case 4:
		return []string{"book_outline", "phase3_output"}
	case 5:
		return []string{"book_outline", "phase4_output"}
	case 6:
		return []string{"book_outline", "phase1_output", "phase2_output", "phase3_output", "phase4_output", "phase5_output"}
	case 7:
		return []string{"book_outline", "phase6_output"}
	default:
		return []string{}
	}
}

// Check if required data exists
func (pc *PhaseCoordinator) hasRequiredData(requirement string) bool {
	switch requirement {
	case "book_outline":
		return pc.BookOutline != nil
	default:
		// Check if phase output exists
		if strings.HasPrefix(requirement, "phase") && strings.HasSuffix(requirement, "_output") {
			phaseKey := strings.TrimSuffix(requirement, "_output")
			_, exists := pc.PhaseData[phaseKey]
			return exists
		}
	}
	return false
}