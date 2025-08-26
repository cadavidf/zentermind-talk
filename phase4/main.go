package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

type Phase3Output struct {
	Timestamp        string           `json:"timestamp"`
	ConceptTitle     string           `json:"concept_title"`
	OverallRating    float64          `json:"overall_rating"`
	EngagementScore  float64          `json:"engagement_score"`
	ShareabilityScore float64         `json:"shareability_score"`
	ApprovedForNext  bool             `json:"approved_for_next"`
}

type MediaOutlet struct {
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Audience    string  `json:"audience"`
	Reach       int     `json:"reach"`
	Probability float64 `json:"probability"`
}

type CoverageEstimate struct {
	Outlet         MediaOutlet `json:"outlet"`
	Likelihood     float64     `json:"likelihood"`
	TimeFrame      string      `json:"timeframe"`
	CoverageType   string      `json:"coverage_type"`
	EstimatedReach int         `json:"estimated_reach"`
}

type Phase4Output struct {
	Timestamp         string              `json:"timestamp"`
	ConceptTitle      string              `json:"concept_title"`
	MediaAnalysis     []CoverageEstimate  `json:"media_analysis"`
	TotalEstimatedReach int               `json:"total_estimated_reach"`
	MediaScore        float64             `json:"media_score"`
	PrWorthiness      string              `json:"pr_worthiness"`
	ApprovedForNext   bool                `json:"approved_for_next"`
}

func main() {
	fmt.Println("🚀 PHASE 4: Media Coverage Prediction")
	fmt.Println(strings.Repeat("=", 50))
	
	// Load Phase 3 output
	phase3Data, err := loadPhase3Output("../phase3/feedback.json")
	if err != nil {
		fmt.Printf("Error loading Phase 3 data: %v\n", err)
		fmt.Println("Using mock data...")
		phase3Data = createMockPhase3Data()
	}
	
	fmt.Printf("\n📖 Analyzing media potential for: %s\n", phase3Data.ConceptTitle)
	
	// Generate media analysis
	mediaAnalysis := generateMediaAnalysis(phase3Data)
	totalReach := calculateTotalReach(mediaAnalysis)
	mediaScore := calculateMediaScore(mediaAnalysis)
	prWorthiness := assessPrWorthiness(mediaScore)
	
	output := Phase4Output{
		Timestamp:         time.Now().Format(time.RFC3339),
		ConceptTitle:      phase3Data.ConceptTitle,
		MediaAnalysis:     mediaAnalysis,
		TotalEstimatedReach: totalReach,
		MediaScore:        mediaScore,
		PrWorthiness:      prWorthiness,
		ApprovedForNext:   mediaScore >= 70.0,
	}
	
	// Save to JSON file
	outputFile := "media.json"
	if err := saveToJSON(output, outputFile); err != nil {
		fmt.Printf("Error saving output: %v\n", err)
		return
	}
	
	fmt.Printf("\n✅ Media coverage analysis complete\n")
	fmt.Printf("📁 Output saved to: %s\n", outputFile)
	fmt.Printf("\n📊 Results:\n")
	fmt.Printf("   Media Score: %.1f/100\n", mediaScore)
	fmt.Printf("   PR Worthiness: %s\n", prWorthiness)
	fmt.Printf("   Total Reach: %d people\n", totalReach)
	fmt.Printf("   Approved: %v\n", output.ApprovedForNext)
	
	fmt.Printf("\n📰 Top Coverage Opportunities:\n")
	for i, analysis := range mediaAnalysis {
		if i >= 3 { break }
		fmt.Printf("   %s: %.1f%% chance (%s)\n", 
			analysis.Outlet.Name, analysis.Likelihood, analysis.TimeFrame)
	}
}

func loadPhase3Output(filename string) (Phase3Output, error) {
	var data Phase3Output
	file, err := os.Open(filename)
	if err != nil {
		return data, err
	}
	defer file.Close()
	
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	return data, err
}

func createMockPhase3Data() Phase3Output {
	return Phase3Output{
		ConceptTitle:     "AI Leadership in 2026 - Practical Framework",
		OverallRating:    8.1,
		EngagementScore:  78.5,
		ShareabilityScore: 83.2,
		ApprovedForNext:  true,
	}
}

func generateMediaAnalysis(phase3Data Phase3Output) []CoverageEstimate {
	outlets := []MediaOutlet{
		{Name: "TechCrunch", Type: "Tech Blog", Audience: "Startups/Tech", Reach: 500000},
		{Name: "Harvard Business Review", Type: "Business Magazine", Audience: "Executives", Reach: 800000},
		{Name: "Wired", Type: "Tech Magazine", Audience: "Tech Enthusiasts", Reach: 300000},
		{Name: "Forbes", Type: "Business Magazine", Audience: "Business Leaders", Reach: 1200000},
		{Name: "MIT Technology Review", Type: "Research Publication", Audience: "Academics/Tech", Reach: 200000},
		{Name: "Fast Company", Type: "Business Magazine", Audience: "Innovation", Reach: 600000},
	}
	
	var estimates []CoverageEstimate
	
	for _, outlet := range outlets {
		// Calculate likelihood based on concept scores
		baseLikelihood := 20.0 // Base 20% chance
		
		// Boost based on engagement and shareability
		engagementBoost := phase3Data.EngagementScore * 0.3
		shareabilityBoost := phase3Data.ShareabilityScore * 0.2
		ratingBoost := phase3Data.OverallRating * 5.0
		
		likelihood := baseLikelihood + engagementBoost + shareabilityBoost + ratingBoost
		if likelihood > 95.0 {
			likelihood = 95.0
		}
		
		// Determine coverage type and timeframe
		coverageType := "Article"
		timeFrame := "3-6 months"
		if likelihood > 80 {
			coverageType = "Feature Article"
			timeFrame = "1-3 months"
		} else if likelihood > 60 {
			timeFrame = "2-4 months"
		}
		
		estimates = append(estimates, CoverageEstimate{
			Outlet:         outlet,
			Likelihood:     likelihood,
			TimeFrame:      timeFrame,
			CoverageType:   coverageType,
			EstimatedReach: int(float64(outlet.Reach) * likelihood / 100.0),
		})
	}
	
	return estimates
}

func calculateTotalReach(estimates []CoverageEstimate) int {
	total := 0
	for _, estimate := range estimates {
		total += estimate.EstimatedReach
	}
	return total
}

func calculateMediaScore(estimates []CoverageEstimate) float64 {
	if len(estimates) == 0 {
		return 0
	}
	
	totalScore := 0.0
	for _, estimate := range estimates {
		// Weight by outlet reach and likelihood
		score := estimate.Likelihood * (float64(estimate.Outlet.Reach) / 10000.0)
		totalScore += score
	}
	
	// Normalize to 0-100 scale
	normalized := totalScore / float64(len(estimates)) / 100.0
	if normalized > 100.0 {
		normalized = 100.0
	}
	
	return normalized
}

func assessPrWorthiness(mediaScore float64) string {
	if mediaScore >= 80 {
		return "Excellent - High media appeal"
	} else if mediaScore >= 70 {
		return "Good - Moderate media interest"
	} else if mediaScore >= 50 {
		return "Fair - Limited media potential"
	} else {
		return "Poor - Low media appeal"
	}
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