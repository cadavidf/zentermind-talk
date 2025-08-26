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

func main() {
	fmt.Println("🚀 PHASE 1: Content Recommendations Generator")
	fmt.Println(strings.Repeat("=", 50))
	
	contentType := getContentType()
	recommendations := generateMockRecommendations(contentType)
	
	output := Phase1Output{
		Timestamp:       time.Now().Format(time.RFC3339),
		ContentType:     contentType,
		Recommendations: recommendations,
		TotalGenerated:  len(recommendations),
		PassingThreshold: 5,
	}
	
	// Save to JSON file
	outputFile := "recommendations.json"
	if err := saveToJSON(output, outputFile); err != nil {
		fmt.Printf("Error saving output: %v\n", err)
		return
	}
	
	fmt.Printf("\n✅ Generated %d %s recommendations\n", len(recommendations), contentType)
	fmt.Printf("📁 Output saved to: %s\n", outputFile)
	fmt.Println("\n📋 Top 3 Recommendations:")
	for i := 0; i < 3 && i < len(recommendations); i++ {
		rec := recommendations[i]
		fmt.Printf("  %d. %s (%.1f%% confidence, %.1f/10 market gap)\n", 
			rec.Rank, rec.Title, rec.Confidence, rec.MarketGapScore)
	}
}

func getContentType() string {
	// Simple mock - in real implementation would get from user input
	return "book"
}

func generateMockRecommendations(contentType string) []ContentRecommendation {
	mockTitles := map[string][]string{
		"book": {
			"AI Leadership in 2026: The Future Executive's Guide",
			"Climate Solutions That Actually Work",
			"The Remote Work Revolution: Building Better Teams",
			"Cryptocurrency Beyond Bitcoin: Next Generation Finance",
			"Mental Health in the Digital Age",
			"Space Economy: The New Business Frontier",
			"Sustainable Cities: Urban Planning for Tomorrow",
			"The Psychology of Virtual Reality",
		},
		"podcast": {
			"Future Tech Talks: AI and Beyond",
			"Climate Action Heroes",
			"Remote Work Mastery",
			"Crypto Conversations",
			"Digital Wellness Journey",
			"Space Business Weekly",
			"Smart City Solutions",
			"VR Psychology Insights",
		},
		"meditation": {
			"AI Anxiety Relief: Meditation for the Digital Age",
			"Climate Grief and Hope: Healing Our Planet Wounds",
			"Remote Work Balance: Mindful Productivity",
			"Financial Stress Freedom",
			"Digital Detox Meditations",
			"Cosmic Perspective: Space-Inspired Mindfulness",
			"Urban Mindfulness: City Living Peace",
			"Virtual Reality Calm",
		},
	}
	
	titles := mockTitles[contentType]
	if titles == nil {
		titles = mockTitles["book"] // Default
	}
	
	var recommendations []ContentRecommendation
	for i, title := range titles {
		recommendations = append(recommendations, ContentRecommendation{
			Title:           title,
			Confidence:      75.0 + float64(i*3), // Mock confidence scores
			Rank:            i + 1,
			MarketGapScore:  6.5 + float64(i)*0.3, // Mock market gap scores
			ContentType:     contentType,
		})
	}
	
	return recommendations
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
func processUSPOptimization(selectedModel, originalConcept string) (Phase1Output, error) {
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
	prompt := fmt.Sprintf(`Evaluate the USP: "%s" against the following criteria:

	CRITERIA 1 - Novelty Score (0.0-1.0): How conceptually novel is this USP compared to existing concepts?
	CRITERIA 2 - Reader Appeal Score (0.0-1.0): How likely is this USP to resonate with readers and drive engagement?
	CRITERIA 3 - Differentiation Score (0.0-1.0): How well does this USP differentiate from these competitors:
	%s

	Provide detailed analysis and assign scores:
	NOVELTY_SCORE: [0.0-1.0]
	READER_APPEAL_SCORE: [0.0-1.0]
	DIFFERENTIATION_SCORE: [0.0-1.0]
	ANALYSIS: [detailed explanation of scores]`, usp.USP, formatCompetitorsForPrompt(competitors))

	response, err := callOllama(selectedModel, prompt)
	if err != nil {
		return usp, err
	}

	// Parse scores from response
	lines := strings.Split(response, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if strings.HasPrefix(line, "NOVELTY_SCORE:") {
			scoreStr := strings.TrimPrefix(line, "NOVELTY_SCORE:")
			fmt.Sscanf(strings.TrimSpace(scoreStr), "%f", &usp.NoveltyScore)
		} else if strings.HasPrefix(line, "READER_APPEAL_SCORE:") {
			scoreStr := strings.TrimPrefix(line, "READER_APPEAL_SCORE:")
			fmt.Sscanf(strings.TrimSpace(scoreStr), "%f", &usp.ReaderAppealScore)
		} else if strings.HasPrefix(line, "DIFFERENTIATION_SCORE:") {
			scoreStr := strings.TrimPrefix(line, "DIFFERENTIATION_SCORE:")
			fmt.Sscanf(strings.TrimSpace(scoreStr), "%f", &usp.DifferentiationScore)
		} else if strings.HasPrefix(line, "ANALYSIS:") {
			usp.AnalysisDetails = strings.TrimPrefix(line, "ANALYSIS:")
		}
	}
	
	// Calculate overall score and check thresholds
	usp.OverallScore = (usp.NoveltyScore + usp.ReaderAppealScore + usp.DifferentiationScore) / 3
	usp.PassesThresholds = usp.NoveltyScore >= 0.60 && usp.ReaderAppealScore >= 0.75 && usp.DifferentiationScore >= 0.70
	
	fmt.Printf("   📊 Scores - Novelty: %.3f, Appeal: %.3f, Differentiation: %.3f\n", 
		usp.NoveltyScore, usp.ReaderAppealScore, usp.DifferentiationScore)
	
	return usp, nil
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