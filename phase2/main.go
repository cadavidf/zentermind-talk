package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type ContentRecommendation struct {
	Title           string  `json:"title"`
	Confidence      float64 `json:"confidence"`
	Rank            int     `json:"rank"`
	MarketGapScore  float64 `json:"market_gap_score"`
	ContentType     string  `json:"content_type"`
}

type Phase1Output struct {
	Timestamp       string                  `json:"timestamp"`
	ContentType     string                  `json:"content_type"`
	Recommendations []ContentRecommendation `json:"recommendations"`
	TotalGenerated  int                     `json:"total_generated"`
	PassingThreshold int                    `json:"passing_threshold"`
}

type ContentConcept struct {
	Title            string  `json:"title"`
	Description      string  `json:"description"`
	UniquenessScore  float64 `json:"uniqueness_score"`
	ViabilityScore   bool    `json:"viability_score"`
	CommercialScore  float64 `json:"commercial_score"`
	Status           string  `json:"status"`
	ContentType      string  `json:"content_type"`
}

type Phase2Output struct {
	Timestamp       string            `json:"timestamp"`
	SelectedTitle   string            `json:"selected_title"`
	GeneratedConcepts []ContentConcept `json:"generated_concepts"`
	ValidConcepts   []ContentConcept  `json:"valid_concepts"`
	BestConcept     ContentConcept    `json:"best_concept"`
	TotalProcessed  int               `json:"total_processed"`
}

func main() {
	fmt.Println("🚀 PHASE 2: Concept Generation & Validation")
	fmt.Println(strings.Repeat("=", 50))
	
	// Load Phase 1 Beta output
	phase1Data, err := loadPhase1Output("../phase1_beta/usp_optimization.json")
	if err != nil {
		fmt.Printf("Error loading Phase 1 Beta data: %v\n", err)
		fmt.Println("Using mock data...")
		phase1Data = createMockPhase1Data()
	}
	
	// Select best recommendation
	selectedRec := selectBestRecommendation(phase1Data.Recommendations)
	fmt.Printf("\n📖 Selected: %s\n", selectedRec.Title)
	
	// Generate concepts
	concepts := generateMockConcepts(selectedRec)
	validConcepts := validateConcepts(concepts)
	bestConcept := selectBestConcept(validConcepts)
	
	output := Phase2Output{
		Timestamp:       time.Now().Format(time.RFC3339),
		SelectedTitle:   selectedRec.Title,
		GeneratedConcepts: concepts,
		ValidConcepts:   validConcepts,
		BestConcept:     bestConcept,
		TotalProcessed:  len(concepts),
	}
	
	// Save to JSON file
	outputFile := "concepts.json"
	if err := saveToJSON(output, outputFile); err != nil {
		fmt.Printf("Error saving output: %v\n", err)
		return
	}
	
	fmt.Printf("\n✅ Generated %d concepts, %d valid\n", len(concepts), len(validConcepts))
	fmt.Printf("📁 Output saved to: %s\n", outputFile)
	fmt.Printf("\n🎯 Best Concept: %s\n", bestConcept.Title)
	fmt.Printf("   Uniqueness: %.1f/10\n", bestConcept.UniquenessScore)
	fmt.Printf("   Commercial: %.1f/10\n", bestConcept.CommercialScore)
	fmt.Printf("   Status: %s\n", bestConcept.Status)
}

func loadPhase1Output(filename string) (Phase1Output, error) {
	var data Phase1Output
	file, err := os.Open(filename)
	if err != nil {
		return data, err
	}
	defer file.Close()
	
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	return data, err
}

func createMockPhase1Data() Phase1Output {
	return Phase1Output{
		ContentType: "book",
		Recommendations: []ContentRecommendation{
			{
				Title: "AI Leadership in 2026: The Future Executive's Guide",
				Confidence: 85.5,
				Rank: 1,
				MarketGapScore: 7.2,
				ContentType: "book",
			},
		},
	}
}

func selectBestRecommendation(recs []ContentRecommendation) ContentRecommendation {
	if len(recs) == 0 {
		return ContentRecommendation{Title: "Default Title", ContentType: "book"}
	}
	return recs[0] // Select highest ranked
}

func generateMockConcepts(rec ContentRecommendation) []ContentConcept {
	concepts := []ContentConcept{
		{
			Title: rec.Title + " - Practical Framework",
			Description: "A hands-on guide to implementing AI leadership strategies",
			UniquenessScore: 8.2,
			ViabilityScore: true,
			CommercialScore: 7.8,
			Status: "approved",
			ContentType: rec.ContentType,
		},
		{
			Title: rec.Title + " - Case Studies Edition",
			Description: "Real-world examples of AI leadership success stories",
			UniquenessScore: 7.5,
			ViabilityScore: true,
			CommercialScore: 8.1,
			Status: "approved",
			ContentType: rec.ContentType,
		},
		{
			Title: rec.Title + " - Technical Deep Dive",
			Description: "Advanced technical concepts for AI-driven organizations",
			UniquenessScore: 9.1,
			ViabilityScore: false,
			CommercialScore: 6.2,
			Status: "rejected",
			ContentType: rec.ContentType,
		},
	}
	
	return concepts
}

func validateConcepts(concepts []ContentConcept) []ContentConcept {
	var valid []ContentConcept
	for _, concept := range concepts {
		if concept.ViabilityScore && concept.UniquenessScore >= 7.0 && concept.CommercialScore >= 7.0 {
			concept.Status = "approved"
			valid = append(valid, concept)
		} else {
			concept.Status = "rejected"
		}
	}
	return valid
}

func selectBestConcept(concepts []ContentConcept) ContentConcept {
	if len(concepts) == 0 {
		return ContentConcept{Status: "none"}
	}
	
	best := concepts[0]
	bestScore := best.UniquenessScore * best.CommercialScore
	
	for _, concept := range concepts[1:] {
		score := concept.UniquenessScore * concept.CommercialScore
		if score > bestScore {
			best = concept
			bestScore = score
		}
	}
	
	return best
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