package main

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

// Book generation state
type BookGenerationState struct {
	BookNumber      int                    `json:"book_number"`
	Theme           BookTheme              `json:"theme"`
	Outline         *BookOutline           `json:"outline"`
	CurrentPhase    int                    `json:"current_phase"`
	PhaseResults    map[string]interface{} `json:"phase_results"`
	StartTime       time.Time              `json:"start_time"`
	LastUpdateTime  time.Time              `json:"last_update_time"`
	Status          string                 `json:"status"`
	OutputDirectory string                 `json:"output_directory"`
}

// Phase configuration for the pipeline
type PhaseInfo struct {
	Number      int    `json:"number"`
	Name        string `json:"name"`
	Directory   string `json:"directory"`
	Description string `json:"description"`
	InputFile   string `json:"input_file"`
	OutputFile  string `json:"output_file"`
}

// Sequential book generator
type SequentialBookGenerator struct {
	Model        string      `json:"model"`
	BookQueue    []BookTheme `json:"book_queue"`
	CurrentBook  int         `json:"current_book"`
	TotalBooks   int         `json:"total_books"`
	OutputDir    string      `json:"output_dir"`
	Phases       []PhaseInfo `json:"phases"`
}

// ASCII Art Banner
func showSequentialBanner() {
	fmt.Println(`
███████╗███████╗ ██████╗ ██╗   ██╗███████╗███╗   ██╗████████╗██╗ █████╗ ██╗     
██╔════╝██╔════╝██╔═══██╗██║   ██║██╔════╝████╗  ██║╚══██╔══╝██║██╔══██╗██║     
███████╗█████╗  ██║   ██║██║   ██║█████╗  ██╔██╗ ██║   ██║   ██║███████║██║     
╚════██║██╔══╝  ██║▄▄ ██║██║   ██║██╔══╝  ██║╚██╗██║   ██║   ██║██╔══██║██║     
███████║███████╗╚██████╔╝╚██████╔╝███████╗██║ ╚████║   ██║   ██║██║  ██║███████╗
╚══════╝╚══════╝ ╚══▀▀═╝  ╚═════╝ ╚══════╝╚═╝  ╚═══╝   ╚═╝   ╚═╝╚═╝  ╚═╝╚══════╝

██████╗  ██████╗  ██████╗ ██╗  ██╗     ██████╗ ███████╗███╗   ██╗███████╗██████╗  █████╗ ████████╗ ██████╗ ██████╗ 
██╔══██╗██╔═══██╗██╔═══██╗██║ ██╔╝    ██╔════╝ ██╔════╝████╗  ██║██╔════╝██╔══██╗██╔══██╗╚══██╔══╝██╔═══██╗██╔══██╗
██████╔╝██║   ██║██║   ██║█████╔╝     ██║  ███╗█████╗  ██╔██╗ ██║█████╗  ██████╔╝███████║   ██║   ██║   ██║██████╔╝
██╔══██╗██║   ██║██║   ██║██╔═██╗     ██║   ██║██╔══╝  ██║╚██╗██║██╔══╝  ██╔══██╗██╔══██║   ██║   ██║   ██║██╔══██╗
██████╔╝╚██████╔╝╚██████╔╝██║  ██╗    ╚██████╔╝███████╗██║ ╚████║███████╗██║  ██║██║  ██║   ██║   ╚██████╔╝██║  ██║
╚═════╝  ╚═════╝  ╚═════╝ ╚═╝  ╚═╝     ╚═════╝ ╚══════╝╚═╝  ╚═══╝╚══════╝╚═╝  ╚═╝╚═╝  ╚═╝   ╚═╝    ╚═════╝ ╚═╝  ╚═╝

          🔄 SEQUENTIAL MULTI-PHASE BOOK GENERATION SYSTEM 🔄
                        Co-authored by animality.ai
`)
}

func main() {
	showSequentialBanner()
	
	// Initialize the generator
	generator := &SequentialBookGenerator{
		Model:      "llama3.2",
		OutputDir:  "output/books",
		TotalBooks: 10,
	}
	
	// Setup phases configuration
	generator.setupPhases()
	
	// Load book themes
	themes, err := loadBookThemes()
	if err != nil {
		fmt.Printf("❌ Error loading themes: %v\n", err)
		return
	}
	
	// Select themes for book generation
	generator.BookQueue = selectDiverseThemes(themes, generator.TotalBooks)
	
	fmt.Printf("🎯 Sequential Book Generation Configuration:\n")
	fmt.Printf("   📚 Total books to generate: %d\n", generator.TotalBooks)
	fmt.Printf("   🤖 AI Model: %s\n", generator.Model)
	fmt.Printf("   📁 Output directory: %s\n", generator.OutputDir)
	fmt.Printf("   🔄 Phases per book: %d\n", len(generator.Phases))
	
	fmt.Printf("\n🎲 Selected book themes:\n")
	for i, theme := range generator.BookQueue {
		fmt.Printf("   %d. %s (%s)\n", i+1, theme.Title, theme.Category)
		fmt.Printf("      \"%s\"\n", theme.MemorablePhrase)
	}
	
	// Start sequential generation
	fmt.Printf("\n🚀 Starting sequential book generation...\n")
	fmt.Printf("════════════════════════════════════════════════════════════════════\n")
	
	// Process books one at a time
	for i, theme := range generator.BookQueue {
		generator.CurrentBook = i + 1
		
		fmt.Printf("\n📖 STARTING BOOK %d/%d: %s\n", generator.CurrentBook, generator.TotalBooks, theme.Title)
		fmt.Printf("═══════════════════════════════════════════════════════════════════════════════\n")
		
		success, outputDir := generator.processBook(theme)
		
		if success {
			fmt.Printf("✅ BOOK %d COMPLETED SUCCESSFULLY!\n", generator.CurrentBook)
			fmt.Printf("📁 Output saved to: %s\n", outputDir)
		} else {
			fmt.Printf("❌ BOOK %d FAILED\n", generator.CurrentBook)
		}
		
		fmt.Printf("═══════════════════════════════════════════════════════════════════════════════\n")
		
		// Brief pause between books
		if i < len(generator.BookQueue)-1 {
			fmt.Printf("⏸️  Pausing 3 seconds before next book...\n")
			time.Sleep(3 * time.Second)
		}
	}
	
	fmt.Printf("\n🎊 SEQUENTIAL BOOK GENERATION COMPLETE! 🎊\n")
	fmt.Printf("📚 Generated %d books through complete 7-phase pipeline\n", generator.TotalBooks)
	fmt.Printf("📁 All books available in: %s\n", generator.OutputDir)
}

// Setup phase configuration
func (g *SequentialBookGenerator) setupPhases() {
	g.Phases = []PhaseInfo{
		{1, "Phase 1 Beta", "phase1_beta", "Market Intelligence & USP Optimization", "", "usp_optimization.json"},
		{2, "Phase 2", "phase2", "Concept Generation & Validation", "usp_optimization.json", "concepts.json"},
		{3, "Phase 3", "phase3", "Reader Feedback & Shareability", "concepts.json", "feedback.json"},
		{4, "Phase 4", "phase4", "Media Coverage & PR Analysis", "feedback.json", "media.json"},
		{5, "Phase 5", "phase5", "Title Optimization & A/B Testing", "media.json", "titles.json"},
		{6, "Phase 6 Beta", "phase6_beta", "Complete Content Generation", "titles.json", "content.json"},
		{7, "Phase 7", "phase7", "Marketing Assets & Campaign", "content.json", "marketing.json"},
	}
}

// Process a single book through all 7 phases
func (g *SequentialBookGenerator) processBook(theme BookTheme) (bool, string) {
	bookNumber := g.CurrentBook
	
	// Create book state
	state := &BookGenerationState{
		BookNumber:      bookNumber,
		Theme:           theme,
		CurrentPhase:    0,
		PhaseResults:    make(map[string]interface{}),
		StartTime:       time.Now(),
		LastUpdateTime:  time.Now(),
		Status:          "starting",
		OutputDirectory: fmt.Sprintf("%s/book_%03d", g.OutputDir, bookNumber),
	}
	
	// Create output directory
	if err := os.MkdirAll(state.OutputDirectory, 0755); err != nil {
		fmt.Printf("❌ Error creating output directory: %v\n", err)
		return false, ""
	}
	
	// Step 1: Generate 1000-word outline
	fmt.Printf("\n📝 STEP 1: Generating 1000-word outline...\n")
	outline, err := GenerateBookOutline(theme, g.Model)
	if err != nil {
		fmt.Printf("❌ Error generating outline: %v\n", err)
		return false, ""
	}
	
	state.Outline = outline
	
	// Save outline
	_, err = SaveOutline(outline, bookNumber)
	if err != nil {
		fmt.Printf("❌ Error saving outline: %v\n", err)
		return false, ""
	}
	
	fmt.Printf("✅ Outline generated (%d words) and saved\n", outline.WordCount)
	
	// Save initial state
	g.saveBookState(state)
	
	// Step 2: Process through all 7 phases
	fmt.Printf("\n🔄 STEP 2: Processing through 7-phase pipeline...\n")
	
	for _, phase := range g.Phases {
		state.CurrentPhase = phase.Number
		state.LastUpdateTime = time.Now()
		state.Status = fmt.Sprintf("processing_phase_%d", phase.Number)
		
		fmt.Printf("\n⚡ PHASE %d: %s\n", phase.Number, phase.Description)
		fmt.Printf("📁 Directory: %s\n", phase.Directory)
		
		success := g.executePhase(phase, state)
		
		if success {
			fmt.Printf("✅ Phase %d completed successfully\n", phase.Number)
			
			// Save progress after each phase
			g.saveBookState(state)
			g.savePhaseProgress(state, phase)
		} else {
			fmt.Printf("❌ Phase %d failed\n", phase.Number)
			state.Status = fmt.Sprintf("failed_phase_%d", phase.Number)
			g.saveBookState(state)
			return false, state.OutputDirectory
		}
		
		// Brief pause between phases
		time.Sleep(1 * time.Second)
	}
	
	// Step 3: Finalize book
	state.Status = "completed"
	state.LastUpdateTime = time.Now()
	duration := time.Since(state.StartTime)
	
	fmt.Printf("\n📚 STEP 3: Finalizing book...\n")
	fmt.Printf("⏱️  Total generation time: %v\n", duration)
	
	// Generate final book summary
	g.generateBookSummary(state, duration)
	
	// Save final state
	g.saveBookState(state)
	
	return true, state.OutputDirectory
}

// Execute a single phase
func (g *SequentialBookGenerator) executePhase(phase PhaseInfo, state *BookGenerationState) bool {
	startTime := time.Now()
	
	// Change to phase directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)
	
	if err := os.Chdir(phase.Directory); err != nil {
		fmt.Printf("❌ Failed to change to directory %s: %v\n", phase.Directory, err)
		return false
	}
	
	// Check if binary exists, build if necessary
	binaryName := phase.Directory
	if _, err := os.Stat(binaryName); os.IsNotExist(err) {
		fmt.Printf("🔨 Building %s...\n", phase.Name)
		buildCmd := exec.Command("go", "build", "-o", binaryName, "main.go")
		if err := buildCmd.Run(); err != nil {
			fmt.Printf("❌ Build failed: %v\n", err)
			return false
		}
		fmt.Printf("✅ Build successful\n")
	}
	
	// Set up environment variables for the phase
	cmd := exec.Command("./"+binaryName)
	cmd.Env = append(os.Environ(),
		fmt.Sprintf("BOOK_TITLE=%s", state.Theme.Title),
		fmt.Sprintf("BOOK_DESCRIPTION=%s", state.Theme.Description),
		fmt.Sprintf("BOOK_PHRASE=%s", state.Theme.MemorablePhrase),
		fmt.Sprintf("BOOK_CATEGORY=%s", state.Theme.Category),
		fmt.Sprintf("BOOK_NUMBER=%03d", state.BookNumber),
		fmt.Sprintf("OUTLINE_AVAILABLE=true"),
		fmt.Sprintf("OUTPUT_DIR=%s", state.OutputDirectory),
		"AUTOMATED_MODE=true",
		fmt.Sprintf("OLLAMA_MODEL=%s", g.Model),
	)
	
	// Execute the phase
	fmt.Printf("▶️  Executing %s...\n", phase.Name)
	output, err := cmd.CombinedOutput()
	
	duration := time.Since(startTime)
	
	if err != nil {
		fmt.Printf("❌ Execution failed after %v: %v\n", duration, err)
		fmt.Printf("Output:\n%s\n", string(output))
		return false
	}
	
	fmt.Printf("⏱️  Completed in %v\n", duration)
	
	// Verify output file was created
	if phase.OutputFile != "" {
		outputFile := phase.OutputFile
		if _, err := os.Stat(outputFile); os.IsNotExist(err) {
			fmt.Printf("⚠️  Warning: Expected output file not found: %s\n", outputFile)
		} else {
			fmt.Printf("📄 Output file created: %s\n", outputFile)
			
			// Copy output file to book directory
			bookOutputFile := filepath.Join(state.OutputDirectory, fmt.Sprintf("phase%d_%s", phase.Number, outputFile))
			if err := copyFile(outputFile, bookOutputFile); err != nil {
				fmt.Printf("⚠️  Warning: Could not copy output file: %v\n", err)
			}
		}
	}
	
	return true
}

// Save book state
func (g *SequentialBookGenerator) saveBookState(state *BookGenerationState) {
	stateFile := filepath.Join(state.OutputDirectory, "book_state.json")
	stateData, _ := json.MarshalIndent(state, "", "  ")
	os.WriteFile(stateFile, stateData, 0644)
}

// Save phase progress
func (g *SequentialBookGenerator) savePhaseProgress(state *BookGenerationState, phase PhaseInfo) {
	progressFile := filepath.Join(state.OutputDirectory, fmt.Sprintf("phase%d_progress.json", phase.Number))
	
	progress := map[string]interface{}{
		"phase_number":   phase.Number,
		"phase_name":     phase.Name,
		"phase_description": phase.Description,
		"completed_at":   time.Now().Format(time.RFC3339),
		"book_number":    state.BookNumber,
		"book_title":     state.Theme.Title,
		"phase_duration": time.Since(state.LastUpdateTime).String(),
		"total_duration": time.Since(state.StartTime).String(),
	}
	
	progressData, _ := json.MarshalIndent(progress, "", "  ")
	os.WriteFile(progressFile, progressData, 0644)
}

// Generate final book summary
func (g *SequentialBookGenerator) generateBookSummary(state *BookGenerationState, duration time.Duration) {
	summaryFile := filepath.Join(state.OutputDirectory, "book_summary.json")
	
	summary := map[string]interface{}{
		"book_number":      state.BookNumber,
		"title":           state.Theme.Title,
		"memorable_phrase": state.Theme.MemorablePhrase,
		"category":        state.Theme.Category,
		"outline_words":   state.Outline.WordCount,
		"phases_completed": len(g.Phases),
		"total_duration":  duration.String(),
		"generated_at":    time.Now().Format(time.RFC3339),
		"status":          "completed",
		"output_files": g.listGeneratedFiles(state.OutputDirectory),
	}
	
	summaryData, _ := json.MarshalIndent(summary, "", "  ")
	os.WriteFile(summaryFile, summaryData, 0644)
	
	fmt.Printf("📋 Book summary saved\n")
}

// List generated files
func (g *SequentialBookGenerator) listGeneratedFiles(dir string) []string {
	var files []string
	filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err == nil && !info.IsDir() {
			relPath, _ := filepath.Rel(dir, path)
			files = append(files, relPath)
		}
		return nil
	})
	return files
}

// Helper function to copy files
func copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

// Load book themes (reuse from previous implementation)
func loadBookThemes() ([]BookTheme, error) {
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
	
	return themes[:count] // Take first 'count' themes for now
}