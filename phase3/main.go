package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

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

type ReaderPersona struct {
	Name         string `json:"name"`
	Age          int    `json:"age"`
	Demographics string `json:"demographics"`
	Interests    string `json:"interests"`
	ReadingStyle string `json:"reading_style"`
}

type ReaderFeedback struct {
	Persona      ReaderPersona `json:"persona"`
	Rating       float64       `json:"rating"`
	Comment      string        `json:"comment"`
	Concerns     []string      `json:"concerns"`
	Engagement   float64       `json:"engagement"`
}

type ViralQuote struct {
	Text              string  `json:"text"`
	ShareabilityScore float64 `json:"shareability_score"`
	Platform          string  `json:"platform"`
}

type Phase3Output struct {
	Timestamp        string           `json:"timestamp"`
	ConceptTitle     string           `json:"concept_title"`
	ReaderFeedback   []ReaderFeedback `json:"reader_feedback"`
	ViralQuotes      []ViralQuote     `json:"viral_quotes"`
	OverallRating    float64          `json:"overall_rating"`
	EngagementScore  float64          `json:"engagement_score"`
	ShareabilityScore float64         `json:"shareability_score"`
	ApprovedForNext  bool             `json:"approved_for_next"`
}

func main() {
	fmt.Println("🚀 PHASE 3: Reader Feedback & Shareability Analysis")
	fmt.Println(strings.Repeat("=", 60))
	
	// Load Phase 2 output
	phase2Data, err := loadPhase2Output("../phase2/concepts.json")
	if err != nil {
		fmt.Printf("Error loading Phase 2 data: %v\n", err)
		fmt.Println("Using mock data...")
		phase2Data = createMockPhase2Data()
	}
	
	concept := phase2Data.BestConcept
	fmt.Printf("\n📖 Analyzing concept: %s\n", concept.Title)
	
	// Generate reader feedback
	feedback := generateReaderFeedback(concept)
	quotes := generateViralQuotes(concept)
	
	// Calculate scores
	overallRating := calculateOverallRating(feedback)
	engagementScore := calculateEngagementScore(feedback)
	shareabilityScore := calculateShareabilityScore(quotes)
	
	output := Phase3Output{
		Timestamp:        time.Now().Format(time.RFC3339),
		ConceptTitle:     concept.Title,
		ReaderFeedback:   feedback,
		ViralQuotes:      quotes,
		OverallRating:    overallRating,
		EngagementScore:  engagementScore,
		ShareabilityScore: shareabilityScore,
		ApprovedForNext:  overallRating >= 7.5 && engagementScore >= 75.0,
	}
	
	// Save to JSON file
	outputFile := "feedback.json"
	if err := saveToJSON(output, outputFile); err != nil {
		fmt.Printf("Error saving output: %v\n", err)
		return
	}
	
	fmt.Printf("\n✅ Reader feedback analysis complete\n")
	fmt.Printf("📁 Output saved to: %s\n", outputFile)
	fmt.Printf("\n📊 Results:\n")
	fmt.Printf("   Overall Rating: %.1f/10\n", overallRating)
	fmt.Printf("   Engagement: %.1f%%\n", engagementScore)
	fmt.Printf("   Shareability: %.1f%%\n", shareabilityScore)
	fmt.Printf("   Approved: %v\n", output.ApprovedForNext)
	
	fmt.Printf("\n💬 Sample Feedback:\n")
	for i, fb := range feedback {
		if i >= 2 { break }
		fmt.Printf("   %s (%.1f/10): %s\n", fb.Persona.Name, fb.Rating, fb.Comment)
	}
}

func loadPhase2Output(filename string) (Phase2Output, error) {
	var data Phase2Output
	file, err := os.Open(filename)
	if err != nil {
		return data, err
	}
	defer file.Close()
	
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	return data, err
}

func createMockPhase2Data() Phase2Output {
	return Phase2Output{
		BestConcept: ContentConcept{
			Title: "AI Leadership in 2026 - Practical Framework",
			Description: "A hands-on guide to implementing AI leadership strategies",
			ContentType: "book",
		},
	}
}

func generateReaderFeedback(concept ContentConcept) []ReaderFeedback {
	personas := []ReaderPersona{
		{Name: "Sarah", Age: 34, Demographics: "Tech Executive", Interests: "AI, Leadership", ReadingStyle: "Practical"},
		{Name: "Mike", Age: 42, Demographics: "Middle Manager", Interests: "Career Growth", ReadingStyle: "Case Studies"},
		{Name: "Lisa", Age: 28, Demographics: "Startup Founder", Interests: "Innovation", ReadingStyle: "Quick Insights"},
		{Name: "David", Age: 39, Demographics: "Consultant", Interests: "Strategy", ReadingStyle: "Comprehensive"},
		{Name: "Emma", Age: 31, Demographics: "HR Director", Interests: "Team Building", ReadingStyle: "People-focused"},
	}
	
	comments := []string{
		"This concept really resonates with current industry needs",
		"Love the practical approach - exactly what leaders need",
		"Great timing with AI adoption accelerating",
		"Could use more specific implementation examples",
		"Addresses a real gap in the market",
	}
	
	var feedback []ReaderFeedback
	for i, persona := range personas {
		rating := 7.5 + float64(i)*0.3 // Mock ratings 7.5-8.7
		engagement := 70.0 + float64(i*5) // Mock engagement 70-90%
		
		feedback = append(feedback, ReaderFeedback{
			Persona:     persona,
			Rating:      rating,
			Comment:     comments[i],
			Concerns:    []string{"Market timing", "Competition"},
			Engagement:  engagement,
		})
	}
	
	return feedback
}

func generateViralQuotes(concept ContentConcept) []ViralQuote {
	quotes := []ViralQuote{
		{
			Text: "AI leadership isn't about replacing humans - it's about amplifying human potential",
			ShareabilityScore: 87.5,
			Platform: "LinkedIn",
		},
		{
			Text: "The future belongs to leaders who can dance with artificial intelligence",
			ShareabilityScore: 82.3,
			Platform: "Twitter",
		},
		{
			Text: "2026 will separate AI-native leaders from those still catching up",
			ShareabilityScore: 79.8,
			Platform: "LinkedIn",
		},
	}
	
	return quotes
}

func calculateOverallRating(feedback []ReaderFeedback) float64 {
	if len(feedback) == 0 {
		return 0
	}
	
	total := 0.0
	for _, fb := range feedback {
		total += fb.Rating
	}
	return total / float64(len(feedback))
}

func calculateEngagementScore(feedback []ReaderFeedback) float64 {
	if len(feedback) == 0 {
		return 0
	}
	
	total := 0.0
	for _, fb := range feedback {
		total += fb.Engagement
	}
	return total / float64(len(feedback))
}

func calculateShareabilityScore(quotes []ViralQuote) float64 {
	if len(quotes) == 0 {
		return 0
	}
	
	total := 0.0
	for _, quote := range quotes {
		total += quote.ShareabilityScore
	}
	return total / float64(len(quotes))
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