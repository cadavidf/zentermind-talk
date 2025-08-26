package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

// Book theme structure
type BookTheme struct {
	ID             int    `json:"id"`
	Title          string `json:"title"`
	MemorablePhrase string `json:"memorable_phrase"`
	Category       string `json:"category"`
	Description    string `json:"description"`
}

// Generation job structure
type GenerationJob struct {
	ID          int       `json:"id"`
	Theme       BookTheme `json:"theme"`
	Status      string    `json:"status"`
	StartTime   time.Time `json:"start_time"`
	EndTime     time.Time `json:"end_time"`
	Duration    string    `json:"duration"`
	OutputFiles []string  `json:"output_files"`
	Error       string    `json:"error,omitempty"`
}

// Batch generation config
type BatchConfig struct {
	TotalBooks       int    `json:"total_books"`
	MaxConcurrent    int    `json:"max_concurrent"`
	OutputDirectory  string `json:"output_directory"`
	LogFile          string `json:"log_file"`
	ThemesFile       string `json:"themes_file"`
	UsePhase1Beta    bool   `json:"use_phase1_beta"`
	ProgressTracking bool   `json:"progress_tracking"`
}

// Generation statistics
type GenerationStats struct {
	TotalBooks     int           `json:"total_books"`
	Successful     int           `json:"successful"`
	Failed         int           `json:"failed"`
	StartTime      time.Time     `json:"start_time"`
	EndTime        time.Time     `json:"end_time"`
	TotalDuration  string        `json:"total_duration"`
	AverageDuration string       `json:"average_duration"`
	Jobs           []GenerationJob `json:"jobs"`
}

var (
	mu    sync.Mutex
	stats GenerationStats
)

// ASCII Art Banner
func showBanner() {
	fmt.Println(`
██████╗  ██████╗  ██████╗ ██╗  ██╗     ██████╗ ███████╗███╗   ██╗███████╗██████╗  █████╗ ████████╗██╗ ██████╗ ███╗   ██╗
██╔══██╗██╔═══██╗██╔═══██╗██║ ██╔╝    ██╔════╝ ██╔════╝████╗  ██║██╔════╝██╔══██╗██╔══██╗╚══██╔══╝██║██╔═══██╗████╗  ██║
██████╔╝██║   ██║██║   ██║█████╔╝     ██║  ███╗█████╗  ██╔██╗ ██║█████╗  ██████╔╝███████║   ██║   ██║██║   ██║██╔██╗ ██║
██╔══██╗██║   ██║██║   ██║██╔═██╗     ██║   ██║██╔══╝  ██║╚██╗██║██╔══╝  ██╔══██╗██╔══██║   ██║   ██║██║   ██║██║╚██╗██║
██████╔╝╚██████╔╝╚██████╔╝██║  ██╗    ╚██████╔╝███████╗██║ ╚████║███████╗██║  ██║██║  ██║   ██║   ██║╚██████╔╝██║ ╚████║
╚═════╝  ╚═════╝  ╚═════╝ ╚═╝  ╚═╝     ╚═════╝ ╚══════╝╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝   ╚═╝ ╚═════╝ ╚═╝  ╚═══╝

                             █████╗  ██████╗ ███████╗███╗   ██╗████████╗
                            ██╔══██╗██╔════╝ ██╔════╝████╗  ██║╚══██╔══╝
                            ███████║██║  ███╗█████╗  ██╔██╗ ██║   ██║   
                            ██╔══██║██║   ██║██╔══╝  ██║╚██╗██║   ██║   
                            ██║  ██║╚██████╔╝███████╗██║ ╚████║   ██║   
                            ╚═╝  ╚═╝ ╚═════╝ ╚══════╝╚═╝  ╚═══╝   ╚═╝   

            🤖 AUTOMATED BOOK GENERATION AGENT - BATCH PROCESSING 🤖
                               Co-authored by animality.ai
`)
}

func main() {
	showBanner()
	
	// Default configuration
	config := BatchConfig{
		TotalBooks:       10, // Generate 10 books as requested
		MaxConcurrent:    1, // Sequential processing for now
		OutputDirectory:  "output/generated_books",
		LogFile:          "output/batch_generation.log",
		ThemesFile:       "future_book_themes.md",
		UsePhase1Beta:    false, // Direct generation to avoid complexity
		ProgressTracking: true,
	}
	
	fmt.Printf("🎯 Starting automated book generation...\n")
	fmt.Printf("📊 Configuration:\n")
	fmt.Printf("   📚 Total books: %d\n", config.TotalBooks)
	fmt.Printf("   🔄 Max concurrent: %d\n", config.MaxConcurrent)
	fmt.Printf("   📁 Output directory: %s\n", config.OutputDirectory)
	fmt.Printf("   📝 Themes file: %s\n", config.ThemesFile)
	
	// Initialize statistics
	stats = GenerationStats{
		TotalBooks: config.TotalBooks,
		StartTime:  time.Now(),
		Jobs:       make([]GenerationJob, 0, config.TotalBooks),
	}
	
	// Load book themes
	themes, err := loadBookThemes(config.ThemesFile)
	if err != nil {
		fmt.Printf("❌ Error loading themes: %v\n", err)
		return
	}
	
	fmt.Printf("✅ Loaded %d book themes\n", len(themes))
	
	// Select diverse themes for generation
	selectedThemes := selectDiverseThemes(themes, config.TotalBooks)
	
	fmt.Printf("\n🎲 Selected themes for generation:\n")
	for i, theme := range selectedThemes {
		fmt.Printf("   %d. %s (%s)\n", i+1, theme.Title, theme.Category)
	}
	
	// Create output directories
	if err := os.MkdirAll(config.OutputDirectory, 0755); err != nil {
		fmt.Printf("❌ Error creating output directory: %v\n", err)
		return
	}
	
	// Start generation process
	fmt.Printf("\n🚀 Starting book generation process...\n")
	fmt.Printf("════════════════════════════════════════════════════════════════════\n")
	
	// Process books sequentially for now
	for i, theme := range selectedThemes {
		job := GenerationJob{
			ID:        i + 1,
			Theme:     theme,
			Status:    "starting",
			StartTime: time.Now(),
		}
		
		fmt.Printf("\n📖 BOOK %d/%d: %s\n", i+1, config.TotalBooks, theme.Title)
		fmt.Printf("📝 Phrase: \"%s\"\n", theme.MemorablePhrase)
		fmt.Printf("🏷️  Category: %s\n", theme.Category)
		fmt.Printf("⏰ Started: %s\n", job.StartTime.Format("15:04:05"))
		
		// Generate the book
		success, outputFiles, errMsg := generateBook(theme, i+1)
		
		// Update job status
		job.EndTime = time.Now()
		job.Duration = job.EndTime.Sub(job.StartTime).String()
		
		if success {
			job.Status = "completed"
			job.OutputFiles = outputFiles
			stats.Successful++
			
			fmt.Printf("✅ COMPLETED in %s\n", job.Duration)
			fmt.Printf("📁 Files generated:\n")
			for _, file := range outputFiles {
				fmt.Printf("   - %s\n", filepath.Base(file))
			}
		} else {
			job.Status = "failed"
			job.Error = errMsg
			stats.Failed++
			
			fmt.Printf("❌ FAILED after %s\n", job.Duration)
			fmt.Printf("💥 Error: %s\n", errMsg)
		}
		
		// Add job to stats
		mu.Lock()
		stats.Jobs = append(stats.Jobs, job)
		mu.Unlock()
		
		// Progress update
		progress := float64(i+1) / float64(config.TotalBooks) * 100
		fmt.Printf("📈 Progress: %.1f%% (%d/%d books)\n", progress, i+1, config.TotalBooks)
		
		// Brief pause between books
		if i < len(selectedThemes)-1 {
			fmt.Printf("⏸️  Pausing for 2 seconds...\n")
			time.Sleep(2 * time.Second)
		}
	}
	
	// Finalize statistics
	stats.EndTime = time.Now()
	stats.TotalDuration = stats.EndTime.Sub(stats.StartTime).String()
	
	if stats.Successful > 0 {
		avgDuration := stats.EndTime.Sub(stats.StartTime) / time.Duration(stats.Successful)
		stats.AverageDuration = avgDuration.String()
	}
	
	// Save final report
	if err := saveFinalReport(config, stats); err != nil {
		fmt.Printf("⚠️  Warning: Could not save final report: %v\n", err)
	}
	
	// Display final summary
	displayFinalSummary(stats)
}

// Load book themes from markdown file
func loadBookThemes(filename string) ([]BookTheme, error) {
	// For now, we'll use predefined themes instead of parsing the markdown
	// This could be enhanced to actually parse the markdown file
	themes := []BookTheme{
		{1, "Your AI, Your Mirror", "Your AI doesn't replace you, it reflects you.", "AI & Personal Development", "How personal AI assistants reveal our patterns, biases, and potential for growth"},
		{2, "The Quiet Machine", "Tech that whispers instead of shouts changes everything.", "Technology & Human Experience", "The future of calm technology and mindful digital design"},
		{3, "Conscious Code", "We programmed machines, then learned to program values.", "AI Ethics & Philosophy", "Building ethical frameworks into artificial intelligence systems"},
		{4, "Live Long, Live Well", "Lifespan is nothing without healthspan.", "Longevity & Health", "Extending healthy years, not just total years of life"},
		{5, "Preventive Everything", "The future of health is prevention, not correction.", "Preventive Medicine & AI", "Using AI and data to prevent disease before it occurs"},
		{6, "The Regenerative Revolution", "Regeneration became the cost of entry.", "Sustainable Business & Environment", "How regenerative practices become essential for business survival"},
		{7, "Waste to Wealth", "Nothing is wasted in a living system.", "Circular Economy & Innovation", "Transforming waste streams into profitable resources"},
		{8, "Own Your Data, Own Your Life", "Your data is your dignity.", "Data Privacy & Digital Rights", "Taking control of personal data as a fundamental right"},
		{9, "The Infinite Canvas", "Reality became our new interface.", "Augmented Reality & User Experience", "How spatial computing transforms how we interact with information"},
		{10, "Focus as Superpower", "Attention is your scarcest resource.", "Attention Economy & Mental Training", "Developing focus as the key skill for the digital age"},
		{11, "Electrify Everything", "Clean electrons powered the future.", "Clean Energy & Electrification", "The complete electrification of transportation, heating, and industry"},
		{12, "Asteroid Harvesters", "We mined space to heal Earth.", "Space Mining & Resource Extraction", "Using space resources to reduce environmental impact on Earth"},
		{13, "Code Meets Biology", "We programmed life like software.", "Synthetic Biology & Programming", "Applying software development principles to biological systems"},
		{14, "Beyond Scarcity", "Abundance became the default.", "Post-Scarcity Economics & Society", "Moving beyond scarcity-based thinking to abundance-based systems"},
		{15, "Stewardship Economy", "Caring became the core of value creation.", "Care Economy & Values-Based Business", "Economic systems based on stewardship and care rather than extraction"},
	}
	
	return themes, nil
}

// Select diverse themes across categories
func selectDiverseThemes(themes []BookTheme, count int) []BookTheme {
	if len(themes) <= count {
		return themes
	}
	
	// Shuffle themes for random selection
	rand.Seed(time.Now().UnixNano())
	shuffled := make([]BookTheme, len(themes))
	copy(shuffled, themes)
	
	for i := range shuffled {
		j := rand.Intn(i + 1)
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	}
	
	return shuffled[:count]
}

// Generate a single book
func generateBook(theme BookTheme, bookNumber int) (bool, []string, string) {
	// Create input for Phase 6 Enhanced
	inputData := map[string]interface{}{
		"title":       theme.Title,
		"description": theme.Description,
		"phrase":      theme.MemorablePhrase,
		"category":    theme.Category,
		"automated":   true,
		"book_number": bookNumber,
	}
	
	// Save input to temporary file
	inputFile := fmt.Sprintf("temp_input_%d.json", bookNumber)
	inputData_json, _ := json.MarshalIndent(inputData, "", "  ")
	if err := os.WriteFile(inputFile, inputData_json, 0644); err != nil {
		return false, nil, fmt.Sprintf("Could not create input file: %v", err)
	}
	defer os.Remove(inputFile) // Clean up
	
	// Execute Phase 6 Enhanced
	cmd := exec.Command("go", "run", "main.go")
	cmd.Dir = "phase6_enhanced"
	
	// Set up environment to use our theme data
	cmd.Env = append(os.Environ(), 
		fmt.Sprintf("BOOK_TITLE=%s", theme.Title),
		fmt.Sprintf("BOOK_DESCRIPTION=%s", theme.Description),
		fmt.Sprintf("BOOK_PHRASE=%s", theme.MemorablePhrase),
		"AUTOMATED_MODE=true",
		"OLLAMA_MODEL=llama3.2",
	)
	
	// Capture output
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		return false, nil, fmt.Sprintf("Phase 6 execution failed: %v\nOutput: %s", err, string(output))
	}
	
	// Find generated files
	outputFiles, err := findGeneratedFiles(theme.Title)
	if err != nil {
		return false, nil, fmt.Sprintf("Could not find generated files: %v", err)
	}
	
	if len(outputFiles) == 0 {
		return false, nil, "No output files were generated"
	}
	
	return true, outputFiles, ""
}

// Find generated files for a book
func findGeneratedFiles(title string) ([]string, error) {
	outputDir := "output/generated_books"
	
	// Clean title for pattern matching
	cleanTitle := strings.ReplaceAll(title, " ", "_")
	cleanTitle = strings.ReplaceAll(cleanTitle, ":", "")
	cleanTitle = strings.ReplaceAll(cleanTitle, "?", "")
	cleanTitle = strings.ReplaceAll(cleanTitle, "!", "")
	cleanTitle = strings.ReplaceAll(cleanTitle, "\"", "")
	cleanTitle = strings.ReplaceAll(cleanTitle, "'", "")
	
	files, err := filepath.Glob(filepath.Join(outputDir, cleanTitle+"*"))
	if err != nil {
		return nil, err
	}
	
	return files, nil
}

// Save final generation report
func saveFinalReport(config BatchConfig, stats GenerationStats) error {
	reportFile := "output/batch_generation_report.json"
	
	report := map[string]interface{}{
		"timestamp":    time.Now().Format(time.RFC3339),
		"config":       config,
		"statistics":   stats,
		"success_rate": float64(stats.Successful) / float64(stats.TotalBooks) * 100,
	}
	
	reportData, err := json.MarshalIndent(report, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(reportFile, reportData, 0644)
}

// Display final summary
func displayFinalSummary(stats GenerationStats) {
	fmt.Printf("\n════════════════════════════════════════════════════════════════════\n")
	fmt.Printf("🎊 BATCH GENERATION COMPLETE! 🎊\n")
	fmt.Printf("════════════════════════════════════════════════════════════════════\n")
	
	fmt.Printf("📊 FINAL STATISTICS:\n")
	fmt.Printf("   📚 Total books attempted: %d\n", stats.TotalBooks)
	fmt.Printf("   ✅ Successfully generated: %d\n", stats.Successful)
	fmt.Printf("   ❌ Failed: %d\n", stats.Failed)
	
	successRate := float64(stats.Successful) / float64(stats.TotalBooks) * 100
	fmt.Printf("   📈 Success rate: %.1f%%\n", successRate)
	
	fmt.Printf("   ⏱️  Total time: %s\n", stats.TotalDuration)
	if stats.AverageDuration != "" {
		fmt.Printf("   ⏱️  Average per book: %s\n", stats.AverageDuration)
	}
	
	fmt.Printf("\n📁 Generated files can be found in:\n")
	fmt.Printf("   - output/generated_books/ (book files)\n")
	fmt.Printf("   - output/batch_generation_report.json (detailed report)\n")
	
	if stats.Successful > 0 {
		fmt.Printf("\n🎉 %d books successfully generated in multiple formats:\n", stats.Successful)
		fmt.Printf("   📖 .epub (e-reader compatible)\n")
		fmt.Printf("   📄 .txt (plain text)\n") 
		fmt.Printf("   📝 .md (markdown)\n")
		fmt.Printf("   📊 .json (structured data)\n")
	}
	
	if stats.Failed > 0 {
		fmt.Printf("\n⚠️  %d books failed to generate. Check the detailed report for error information.\n", stats.Failed)
	}
	
	fmt.Printf("\n🤖 Thank you for using BULLET BOOKS Generation Agent!\n")
	fmt.Printf("Co-authored by animality.ai\n")
	fmt.Printf("════════════════════════════════════════════════════════════════════\n")
}