package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// Progress tracker manages state, resumption, and detailed logging
type ProgressTracker struct {
	SessionID        string                 `json:"session_id"`
	StartTime        time.Time              `json:"start_time"`
	LastUpdateTime   time.Time              `json:"last_update_time"`
	TotalBooks       int                    `json:"total_books"`
	CompletedBooks   int                    `json:"completed_books"`
	CurrentBook      *BookProgress          `json:"current_book"`
	BookHistory      []BookProgress         `json:"book_history"`
	SystemMetrics    SystemMetrics          `json:"system_metrics"`
	Configuration    TrackerConfiguration   `json:"configuration"`
	StateFile        string                 `json:"state_file"`
}

// Individual book progress
type BookProgress struct {
	BookNumber       int                     `json:"book_number"`
	BookTitle        string                  `json:"book_title"`
	Theme            BookTheme               `json:"theme"`
	Status           string                  `json:"status"` // starting, outline_generated, phase_processing, completed, failed
	StartTime        time.Time               `json:"start_time"`
	EndTime          time.Time               `json:"end_time"`
	TotalDuration    string                  `json:"total_duration"`
	CurrentPhase     int                     `json:"current_phase"`
	PhasesCompleted  []PhaseProgress         `json:"phases_completed"`
	OutlineGenerated bool                    `json:"outline_generated"`
	OutlineWordCount int                     `json:"outline_word_count"`
	OutputDirectory  string                  `json:"output_directory"`
	GeneratedFiles   []string                `json:"generated_files"`
	ErrorMessages    []string                `json:"error_messages"`
	QualityMetrics   BookQualityMetrics      `json:"quality_metrics"`
}

// Phase-level progress tracking
type PhaseProgress struct {
	PhaseNumber     int                    `json:"phase_number"`
	PhaseName       string                 `json:"phase_name"`
	Status          string                 `json:"status"` // pending, running, completed, failed
	StartTime       time.Time              `json:"start_time"`
	EndTime         time.Time              `json:"end_time"`
	Duration        string                 `json:"duration"`
	InputFiles      []string               `json:"input_files"`
	OutputFiles     []string               `json:"output_files"`
	Improvements    []string               `json:"improvements"`
	QualityScores   map[string]float64     `json:"quality_scores"`
	ResourceUsage   PhaseResourceUsage     `json:"resource_usage"`
	AttemptCount    int                    `json:"attempt_count"`
	LastError       string                 `json:"last_error,omitempty"`
}

// Resource usage tracking per phase
type PhaseResourceUsage struct {
	CPUTime       float64 `json:"cpu_time_seconds"`
	MemoryPeak    int64   `json:"memory_peak_mb"`
	DiskIORead    int64   `json:"disk_io_read_mb"`
	DiskIOWrite   int64   `json:"disk_io_write_mb"`
	NetworkCalls  int     `json:"network_calls"`
	LLMTokens     int     `json:"llm_tokens_used"`
}

// Book quality metrics
type BookQualityMetrics struct {
	OutlineQuality    float64 `json:"outline_quality"`
	ContentQuality    float64 `json:"content_quality"`
	MarketViability   float64 `json:"market_viability"`
	ReaderEngagement  float64 `json:"reader_engagement"`
	SEOScore          float64 `json:"seo_score"`
	OverallScore      float64 `json:"overall_score"`
	WordCount         int     `json:"word_count"`
	ChapterCount      int     `json:"chapter_count"`
}

// System-wide metrics
type SystemMetrics struct {
	TotalExecutionTime  string  `json:"total_execution_time"`
	AverageBookTime     string  `json:"average_book_time"`
	SuccessRate         float64 `json:"success_rate"`
	TotalOutputFiles    int     `json:"total_output_files"`
	TotalDiskUsage      int64   `json:"total_disk_usage_mb"`
	PeakMemoryUsage     int64   `json:"peak_memory_usage_mb"`
	TotalLLMTokens      int     `json:"total_llm_tokens"`
	SystemLoad          float64 `json:"system_load"`
}

// Tracker configuration
type TrackerConfiguration struct {
	AutoSaveInterval    time.Duration `json:"auto_save_interval"`
	DetailedLogging     bool          `json:"detailed_logging"`
	ResourceMonitoring  bool          `json:"resource_monitoring"`
	BackupStateFiles    bool          `json:"backup_state_files"`
	MaxRetryAttempts    int           `json:"max_retry_attempts"`
	OutputDirectory     string        `json:"output_directory"`
}

// Initialize progress tracker
func NewProgressTracker(totalBooks int, outputDir string) *ProgressTracker {
	sessionID := fmt.Sprintf("session_%d", time.Now().Unix())
	
	tracker := &ProgressTracker{
		SessionID:      sessionID,
		StartTime:      time.Now(),
		LastUpdateTime: time.Now(),
		TotalBooks:     totalBooks,
		CompletedBooks: 0,
		BookHistory:    []BookProgress{},
		Configuration: TrackerConfiguration{
			AutoSaveInterval:   30 * time.Second,
			DetailedLogging:    true,
			ResourceMonitoring: true,
			BackupStateFiles:   true,
			MaxRetryAttempts:   3,
			OutputDirectory:    outputDir,
		},
		StateFile: filepath.Join(outputDir, "progress_tracker_state.json"),
	}
	
	// Ensure output directory exists
	os.MkdirAll(outputDir, 0755)
	
	return tracker
}

// Start tracking a new book
func (pt *ProgressTracker) StartBook(bookNumber int, theme BookTheme) {
	bookProgress := BookProgress{
		BookNumber:       bookNumber,
		BookTitle:        theme.Title,
		Theme:            theme,
		Status:           "starting",
		StartTime:        time.Now(),
		CurrentPhase:     0,
		PhasesCompleted:  []PhaseProgress{},
		OutlineGenerated: false,
		OutputDirectory:  filepath.Join(pt.Configuration.OutputDirectory, fmt.Sprintf("book_%03d", bookNumber)),
		GeneratedFiles:   []string{},
		ErrorMessages:    []string{},
		QualityMetrics:   BookQualityMetrics{},
	}
	
	pt.CurrentBook = &bookProgress
	pt.LastUpdateTime = time.Now()
	
	// Create book output directory
	os.MkdirAll(bookProgress.OutputDirectory, 0755)
	
	// Save initial state
	pt.saveState()
	
	fmt.Printf("📊 Started tracking book %d: %s\n", bookNumber, theme.Title)
}

// Update outline generation progress
func (pt *ProgressTracker) UpdateOutlineGenerated(wordCount int, qualityScore float64) {
	if pt.CurrentBook == nil {
		return
	}
	
	pt.CurrentBook.OutlineGenerated = true
	pt.CurrentBook.OutlineWordCount = wordCount
	pt.CurrentBook.QualityMetrics.OutlineQuality = qualityScore
	pt.CurrentBook.Status = "outline_generated"
	pt.LastUpdateTime = time.Now()
	
	// Add outline file to generated files
	outlineFile := filepath.Join(pt.CurrentBook.OutputDirectory, "outline_1000_words.md")
	pt.CurrentBook.GeneratedFiles = append(pt.CurrentBook.GeneratedFiles, outlineFile)
	
	pt.saveState()
	
	fmt.Printf("📝 Outline generated: %d words (quality: %.1f/10)\n", wordCount, qualityScore)
}

// Start tracking a phase
func (pt *ProgressTracker) StartPhase(phaseNumber int, phaseName string) {
	if pt.CurrentBook == nil {
		return
	}
	
	phaseProgress := PhaseProgress{
		PhaseNumber:   phaseNumber,
		PhaseName:     phaseName,
		Status:        "running",
		StartTime:     time.Now(),
		InputFiles:    []string{},
		OutputFiles:   []string{},
		Improvements:  []string{},
		QualityScores: make(map[string]float64),
		ResourceUsage: PhaseResourceUsage{},
		AttemptCount:  1,
	}
	
	pt.CurrentBook.CurrentPhase = phaseNumber
	pt.CurrentBook.Status = "phase_processing"
	pt.LastUpdateTime = time.Now()
	
	// Add to current book's phases
	pt.CurrentBook.PhasesCompleted = append(pt.CurrentBook.PhasesCompleted, phaseProgress)
	
	pt.saveState()
	
	fmt.Printf("⚡ Started Phase %d: %s\n", phaseNumber, phaseName)
}

// Complete a phase with results
func (pt *ProgressTracker) CompletePhase(phaseNumber int, outputFiles []string, improvements []string, qualityScores map[string]float64) {
	if pt.CurrentBook == nil {
		return
	}
	
	// Find and update the phase
	for i := range pt.CurrentBook.PhasesCompleted {
		phase := &pt.CurrentBook.PhasesCompleted[i]
		if phase.PhaseNumber == phaseNumber {
			phase.Status = "completed"
			phase.EndTime = time.Now()
			phase.Duration = phase.EndTime.Sub(phase.StartTime).String()
			phase.OutputFiles = outputFiles
			phase.Improvements = improvements
			phase.QualityScores = qualityScores
			
			// Add files to book's generated files
			pt.CurrentBook.GeneratedFiles = append(pt.CurrentBook.GeneratedFiles, outputFiles...)
			
			// Update book quality metrics
			pt.updateBookQualityMetrics(qualityScores)
			
			break
		}
	}
	
	pt.LastUpdateTime = time.Now()
	pt.saveState()
	
	fmt.Printf("✅ Phase %d completed: %d files, %d improvements\n", phaseNumber, len(outputFiles), len(improvements))
}

// Fail a phase with error information
func (pt *ProgressTracker) FailPhase(phaseNumber int, errorMsg string) {
	if pt.CurrentBook == nil {
		return
	}
	
	// Find and update the phase
	for i := range pt.CurrentBook.PhasesCompleted {
		phase := &pt.CurrentBook.PhasesCompleted[i]
		if phase.PhaseNumber == phaseNumber {
			phase.Status = "failed"
			phase.EndTime = time.Now()
			phase.Duration = phase.EndTime.Sub(phase.StartTime).String()
			phase.LastError = errorMsg
			phase.AttemptCount++
			
			break
		}
	}
	
	// Add error to book errors
	pt.CurrentBook.ErrorMessages = append(pt.CurrentBook.ErrorMessages, fmt.Sprintf("Phase %d: %s", phaseNumber, errorMsg))
	pt.LastUpdateTime = time.Now()
	
	pt.saveState()
	
	fmt.Printf("❌ Phase %d failed: %s\n", phaseNumber, errorMsg)
}

// Complete current book
func (pt *ProgressTracker) CompleteBook(finalFiles []string) {
	if pt.CurrentBook == nil {
		return
	}
	
	pt.CurrentBook.Status = "completed"
	pt.CurrentBook.EndTime = time.Now()
	pt.CurrentBook.TotalDuration = pt.CurrentBook.EndTime.Sub(pt.CurrentBook.StartTime).String()
	pt.CurrentBook.GeneratedFiles = append(pt.CurrentBook.GeneratedFiles, finalFiles...)
	
	// Calculate final quality score
	pt.calculateFinalQualityScore()
	
	// Move to book history
	pt.BookHistory = append(pt.BookHistory, *pt.CurrentBook)
	pt.CompletedBooks++
	pt.CurrentBook = nil
	
	pt.updateSystemMetrics()
	pt.LastUpdateTime = time.Now()
	pt.saveState()
	
	fmt.Printf("🎉 Book completed! Total files: %d\n", len(finalFiles))
}

// Fail current book
func (pt *ProgressTracker) FailBook(errorMsg string) {
	if pt.CurrentBook == nil {
		return
	}
	
	pt.CurrentBook.Status = "failed"
	pt.CurrentBook.EndTime = time.Now()
	pt.CurrentBook.TotalDuration = pt.CurrentBook.EndTime.Sub(pt.CurrentBook.StartTime).String()
	pt.CurrentBook.ErrorMessages = append(pt.CurrentBook.ErrorMessages, errorMsg)
	
	// Move to book history
	pt.BookHistory = append(pt.BookHistory, *pt.CurrentBook)
	pt.CurrentBook = nil
	
	pt.updateSystemMetrics()
	pt.LastUpdateTime = time.Now()
	pt.saveState()
	
	fmt.Printf("💥 Book failed: %s\n", errorMsg)
}

// Update book quality metrics from phase scores
func (pt *ProgressTracker) updateBookQualityMetrics(phaseScores map[string]float64) {
	if pt.CurrentBook == nil {
		return
	}
	
	metrics := &pt.CurrentBook.QualityMetrics
	
	// Update specific metrics based on phase scores
	for key, value := range phaseScores {
		switch key {
		case "market_gap_score", "market_viability":
			metrics.MarketViability = value
		case "engagement_score", "reader_appeal":
			metrics.ReaderEngagement = value
		case "seo_score":
			metrics.SEOScore = value
		case "content_quality":
			metrics.ContentQuality = value
		case "word_count":
			metrics.WordCount = int(value)
		case "chapter_count":
			metrics.ChapterCount = int(value)
		}
	}
}

// Calculate final quality score
func (pt *ProgressTracker) calculateFinalQualityScore() {
	if pt.CurrentBook == nil {
		return
	}
	
	metrics := &pt.CurrentBook.QualityMetrics
	
	// Calculate weighted overall score
	scores := []float64{
		metrics.OutlineQuality * 0.15,   // 15%
		metrics.ContentQuality * 0.25,   // 25%
		metrics.MarketViability * 0.20,  // 20%
		metrics.ReaderEngagement * 0.20, // 20%
		metrics.SEOScore * 0.20,         // 20%
	}
	
	total := 0.0
	count := 0
	for _, score := range scores {
		if score > 0 {
			total += score
			count++
		}
	}
	
	if count > 0 {
		metrics.OverallScore = total / float64(count) * 10.0 // Scale to 10
	}
}

// Update system-wide metrics
func (pt *ProgressTracker) updateSystemMetrics() {
	totalBooks := len(pt.BookHistory)
	if totalBooks == 0 {
		return
	}
	
	// Calculate success rate
	successCount := 0
	totalFiles := 0
	totalDuration := time.Duration(0)
	
	for _, book := range pt.BookHistory {
		if book.Status == "completed" {
			successCount++
		}
		totalFiles += len(book.GeneratedFiles)
		
		if duration, err := time.ParseDuration(book.TotalDuration); err == nil {
			totalDuration += duration
		}
	}
	
	pt.SystemMetrics.SuccessRate = float64(successCount) / float64(totalBooks) * 100.0
	pt.SystemMetrics.TotalOutputFiles = totalFiles
	pt.SystemMetrics.TotalExecutionTime = time.Since(pt.StartTime).String()
	
	if successCount > 0 {
		avgDuration := totalDuration / time.Duration(successCount)
		pt.SystemMetrics.AverageBookTime = avgDuration.String()
	}
}

// Generate comprehensive progress report
func (pt *ProgressTracker) GenerateProgressReport() map[string]interface{} {
	report := map[string]interface{}{
		"session_info": map[string]interface{}{
			"session_id":        pt.SessionID,
			"start_time":        pt.StartTime.Format(time.RFC3339),
			"last_update":       pt.LastUpdateTime.Format(time.RFC3339),
			"total_duration":    time.Since(pt.StartTime).String(),
		},
		"overall_progress": map[string]interface{}{
			"total_books":       pt.TotalBooks,
			"completed_books":   pt.CompletedBooks,
			"books_remaining":   pt.TotalBooks - pt.CompletedBooks,
			"completion_rate":   float64(pt.CompletedBooks) / float64(pt.TotalBooks) * 100.0,
		},
		"current_book": nil,
		"book_history": pt.BookHistory,
		"system_metrics": pt.SystemMetrics,
		"quality_summary": pt.generateQualitySummary(),
		"performance_summary": pt.generatePerformanceSummary(),
	}
	
	if pt.CurrentBook != nil {
		report["current_book"] = pt.CurrentBook
	}
	
	return report
}

// Generate quality summary across all books
func (pt *ProgressTracker) generateQualitySummary() map[string]interface{} {
	if len(pt.BookHistory) == 0 {
		return map[string]interface{}{}
	}
	
	totalQuality := 0.0
	totalMarket := 0.0
	totalEngagement := 0.0
	count := 0
	
	for _, book := range pt.BookHistory {
		if book.Status == "completed" {
			totalQuality += book.QualityMetrics.OverallScore
			totalMarket += book.QualityMetrics.MarketViability
			totalEngagement += book.QualityMetrics.ReaderEngagement
			count++
		}
	}
	
	if count == 0 {
		return map[string]interface{}{}
	}
	
	return map[string]interface{}{
		"average_quality":      totalQuality / float64(count),
		"average_market_score": totalMarket / float64(count),
		"average_engagement":   totalEngagement / float64(count),
		"books_analyzed":       count,
	}
}

// Generate performance summary
func (pt *ProgressTracker) generatePerformanceSummary() map[string]interface{} {
	totalPhases := 0
	failedPhases := 0
	
	for _, book := range pt.BookHistory {
		for _, phase := range book.PhasesCompleted {
			totalPhases++
			if phase.Status == "failed" {
				failedPhases++
			}
		}
	}
	
	phaseSuccessRate := 0.0
	if totalPhases > 0 {
		phaseSuccessRate = float64(totalPhases-failedPhases) / float64(totalPhases) * 100.0
	}
	
	return map[string]interface{}{
		"total_phases_executed": totalPhases,
		"failed_phases":         failedPhases,
		"phase_success_rate":    phaseSuccessRate,
		"average_phase_time":    pt.calculateAveragePhaseTime(),
	}
}

// Calculate average phase execution time
func (pt *ProgressTracker) calculateAveragePhaseTime() string {
	totalDuration := time.Duration(0)
	count := 0
	
	for _, book := range pt.BookHistory {
		for _, phase := range book.PhasesCompleted {
			if phase.Status == "completed" {
				if duration, err := time.ParseDuration(phase.Duration); err == nil {
					totalDuration += duration
					count++
				}
			}
		}
	}
	
	if count == 0 {
		return "0s"
	}
	
	return (totalDuration / time.Duration(count)).String()
}

// Save tracker state to file
func (pt *ProgressTracker) saveState() {
	if pt.Configuration.BackupStateFiles {
		// Create backup
		backupFile := pt.StateFile + ".backup"
		if _, err := os.Stat(pt.StateFile); err == nil {
			os.Rename(pt.StateFile, backupFile)
		}
	}
	
	data, err := json.MarshalIndent(pt, "", "  ")
	if err != nil {
		fmt.Printf("⚠️  Warning: Could not marshal tracker state: %v\n", err)
		return
	}
	
	if err := os.WriteFile(pt.StateFile, data, 0644); err != nil {
		fmt.Printf("⚠️  Warning: Could not save tracker state: %v\n", err)
	}
}

// Load tracker state from file
func LoadProgressTracker(stateFile string) (*ProgressTracker, error) {
	data, err := os.ReadFile(stateFile)
	if err != nil {
		return nil, err
	}
	
	var tracker ProgressTracker
	if err := json.Unmarshal(data, &tracker); err != nil {
		return nil, err
	}
	
	tracker.StateFile = stateFile
	return &tracker, nil
}

// Generate final summary report
func (pt *ProgressTracker) GenerateFinalReport() {
	reportFile := filepath.Join(pt.Configuration.OutputDirectory, "final_generation_report.json")
	
	report := pt.GenerateProgressReport()
	report["generation_complete"] = true
	report["final_summary"] = map[string]interface{}{
		"total_execution_time": time.Since(pt.StartTime).String(),
		"books_completed":      pt.CompletedBooks,
		"success_rate":         pt.SystemMetrics.SuccessRate,
		"total_files_created":  pt.SystemMetrics.TotalOutputFiles,
		"session_id":           pt.SessionID,
	}
	
	data, _ := json.MarshalIndent(report, "", "  ")
	os.WriteFile(reportFile, data, 0644)
	
	fmt.Printf("📊 Final report saved to: %s\n", reportFile)
}