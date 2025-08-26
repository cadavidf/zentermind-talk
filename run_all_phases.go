package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Phase configuration
type PhaseConfig struct {
	Name        string
	Directory   string
	Description string
	InputFile   string
	OutputFile  string
	Required    bool
}

// Pipeline orchestrator
func main() {
	fmt.Println(`
██████╗ ██╗   ██╗██╗     ██╗     ███████╗████████╗    ██████╗  ██████╗  ██████╗ ██╗  ██╗███████╗
██╔══██╗██║   ██║██║     ██║     ██╔════╝╚══██╔══╝    ██╔══██╗██╔═══██╗██╔═══██╗██║ ██╔╝██╔════╝
██████╔╝██║   ██║██║     ██║     █████╗     ██║       ██████╔╝██║   ██║██║   ██║█████╔╝ ███████╗
██╔══██╗██║   ██║██║     ██║     ██╔══╝     ██║       ██╔══██╗██║   ██║██║   ██║██╔═██╗ ╚════██║
██████╔╝╚██████╔╝███████╗███████╗███████╗   ██║       ██████╔╝╚██████╔╝╚██████╔╝██║  ██╗███████║
╚═════╝  ╚═════╝ ╚══════╝╚══════╝╚══════╝   ╚═╝       ╚═════╝  ╚═════╝  ╚═════╝ ╚═╝  ╚═╝╚══════╝

                    🚀 COMPLETE PIPELINE ORCHESTRATOR 🚀
`)

	// Define phase pipeline
	phases := []PhaseConfig{
		{
			Name:        "Phase 1 Beta",
			Directory:   "phase1_beta",
			Description: "Market Intelligence & USP Optimization",
			InputFile:   "",
			OutputFile:  "usp_optimization.json",
			Required:    true,
		},
		{
			Name:        "Phase 2",
			Directory:   "phase2",
			Description: "Concept Generation & Validation",
			InputFile:   "../phase1_beta/usp_optimization.json",
			OutputFile:  "concepts.json",
			Required:    true,
		},
		{
			Name:        "Phase 3",
			Directory:   "phase3",
			Description: "Reader Feedback & Shareability",
			InputFile:   "../phase2/concepts.json",
			OutputFile:  "feedback.json",
			Required:    true,
		},
		{
			Name:        "Phase 4",
			Directory:   "phase4",
			Description: "Media Coverage & PR Analysis",
			InputFile:   "../phase3/feedback.json",
			OutputFile:  "media.json",
			Required:    true,
		},
		{
			Name:        "Phase 5",
			Directory:   "phase5",
			Description: "Title Optimization & A/B Testing",
			InputFile:   "../phase4/media.json",
			OutputFile:  "titles.json",
			Required:    true,
		},
		{
			Name:        "Phase 6",
			Directory:   "phase6_enhanced",
			Description: "Complete Content Generation",
			InputFile:   "../phase5/titles.json",
			OutputFile:  "content.json",
			Required:    true,
		},
		{
			Name:        "Phase 7",
			Directory:   "phase7",
			Description: "Marketing Assets & Campaign",
			InputFile:   "../phase6_enhanced/content.json",
			OutputFile:  "marketing.json",
			Required:    true,
		},
	}

	// Get current directory
	currentDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("❌ Error getting current directory: %v\n", err)
		return
	}

	// Pipeline execution
	fmt.Printf("📍 Starting pipeline from: %s\n", currentDir)
	fmt.Printf("🎯 Total phases to execute: %d\n\n", len(phases))

	startTime := time.Now()
	var executedPhases []string
	var failedPhases []string

	for i, phase := range phases {
		fmt.Printf("═══════════════════════════════════════════════════════════════════════\n")
		fmt.Printf("🔄 EXECUTING %s (%d/%d)\n", phase.Name, i+1, len(phases))
		fmt.Printf("📝 %s\n", phase.Description)
		fmt.Printf("📁 Directory: %s\n", phase.Directory)
		
		if phase.InputFile != "" {
			fmt.Printf("📥 Input: %s\n", phase.InputFile)
		}
		fmt.Printf("📤 Output: %s\n", phase.OutputFile)
		fmt.Printf("═══════════════════════════════════════════════════════════════════════\n")

		// Check if phase directory exists
		phaseDir := filepath.Join(currentDir, phase.Directory)
		if _, err := os.Stat(phaseDir); os.IsNotExist(err) {
			fmt.Printf("⚠️  Phase directory not found: %s\n", phaseDir)
			if phase.Required {
				fmt.Printf("❌ Required phase failed - stopping pipeline\n")
				failedPhases = append(failedPhases, phase.Name)
				break
			}
			continue
		}

		// Execute phase
		success := executePhase(phase, phaseDir)
		if success {
			fmt.Printf("✅ %s completed successfully\n\n", phase.Name)
			executedPhases = append(executedPhases, phase.Name)
		} else {
			fmt.Printf("❌ %s failed\n\n", phase.Name)
			failedPhases = append(failedPhases, phase.Name)
			if phase.Required {
				fmt.Printf("🛑 Required phase failed - stopping pipeline\n")
				break
			}
		}

		// Brief pause between phases
		time.Sleep(1 * time.Second)
	}

	// Pipeline summary
	totalTime := time.Since(startTime)
	fmt.Printf("\n🏁 PIPELINE EXECUTION SUMMARY\n")
	fmt.Printf("══════════════════════════════════════════════════════════════════════\n")
	fmt.Printf("⏱️  Total execution time: %v\n", totalTime)
	fmt.Printf("✅ Successful phases: %d\n", len(executedPhases))
	fmt.Printf("❌ Failed phases: %d\n", len(failedPhases))

	if len(executedPhases) > 0 {
		fmt.Printf("\n🎉 COMPLETED PHASES:\n")
		for _, phaseName := range executedPhases {
			fmt.Printf("   ✓ %s\n", phaseName)
		}
	}

	if len(failedPhases) > 0 {
		fmt.Printf("\n⚠️  FAILED PHASES:\n")
		for _, phaseName := range failedPhases {
			fmt.Printf("   ✗ %s\n", phaseName)
		}
	}

	// Final status
	if len(failedPhases) == 0 {
		fmt.Printf("\n🎊 ALL PHASES COMPLETED SUCCESSFULLY! 🎊\n")
		fmt.Printf("📂 Check individual phase directories for outputs\n")
	} else {
		fmt.Printf("\n⚠️  Pipeline completed with %d failures\n", len(failedPhases))
		fmt.Printf("💡 Check individual phase logs for error details\n")
	}

	fmt.Printf("══════════════════════════════════════════════════════════════════════\n")
}

// Execute a single phase
func executePhase(phase PhaseConfig, phaseDir string) bool {
	// Change to phase directory
	originalDir, _ := os.Getwd()
	defer os.Chdir(originalDir)

	if err := os.Chdir(phaseDir); err != nil {
		fmt.Printf("❌ Failed to change to directory %s: %v\n", phaseDir, err)
		return false
	}

	// Check if binary exists, if not build it
	binaryName := phase.Directory
	if _, err := os.Stat(binaryName); os.IsNotExist(err) {
		fmt.Printf("🔨 Building %s...\n", phase.Name)
		buildCmd := exec.Command("go", "build", "-o", binaryName, "main.go")
		buildOutput, err := buildCmd.CombinedOutput()
		if err != nil {
			fmt.Printf("❌ Build failed: %v\n", err)
			fmt.Printf("Build output: %s\n", string(buildOutput))
			return false
		}
		fmt.Printf("✅ Build successful\n")
	}

	// Execute the phase
	fmt.Printf("▶️  Running %s...\n", phase.Name)
	
	var cmd *exec.Cmd
	if _, err := os.Stat(binaryName); err == nil {
		// Run the binary
		cmd = exec.Command("./" + binaryName)
	} else {
		// Fall back to go run
		cmd = exec.Command("go", "run", "main.go")
	}

	// Capture output
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("❌ Execution failed: %v\n", err)
		fmt.Printf("Output:\n%s\n", string(output))
		return false
	}

	// Display condensed output (last few lines)
	outputLines := strings.Split(string(output), "\n")
	displayLines := 5
	if len(outputLines) > displayLines {
		fmt.Printf("📋 Last %d lines of output:\n", displayLines)
		for _, line := range outputLines[len(outputLines)-displayLines-1:] {
			if strings.TrimSpace(line) != "" {
				fmt.Printf("   %s\n", line)
			}
		}
	} else {
		fmt.Printf("📋 Output:\n%s\n", string(output))
	}

	// Verify output file was created
	if phase.OutputFile != "" {
		if _, err := os.Stat(phase.OutputFile); os.IsNotExist(err) {
			fmt.Printf("⚠️  Warning: Expected output file not found: %s\n", phase.OutputFile)
			// Don't fail the phase for missing output file - it might be optional
		} else {
			fmt.Printf("📄 Output file created: %s\n", phase.OutputFile)
		}
	}

	return true
}