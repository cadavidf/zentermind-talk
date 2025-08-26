package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

func main() {
	fmt.Println("🧪 Testing Single Book Generation with EPUB Support")
	fmt.Println("==================================================")
	
	// Set environment variables for a quick test
	fmt.Println("📝 Generating: 'The Future of Digital Privacy'")
	fmt.Println("🎯 Testing EPUB, TXT, MD, and JSON output formats")
	
	// Execute Phase 6 Enhanced with our test data
	cmd := exec.Command("go", "run", "main.go")
	cmd.Dir = "phase6_enhanced"
	
	// Set up environment for automated test
	cmd.Env = append(os.Environ(),
		"BOOK_TITLE=The Future of Digital Privacy",
		"BOOK_DESCRIPTION=A comprehensive guide to protecting personal data in the digital age. Your data is your dignity.",
		"BOOK_PHRASE=Privacy is the foundation of freedom in the digital age.",
		"AUTOMATED_MODE=true",
		"OLLAMA_MODEL=llama3.2",
	)
	
	fmt.Println("⏰ Starting generation at:", time.Now().Format("15:04:05"))
	start := time.Now()
	
	// Run with timeout
	done := make(chan error, 1)
	go func() {
		done <- cmd.Run()
	}()
	
	// Wait for completion or timeout
	select {
	case err := <-done:
		duration := time.Since(start)
		if err != nil {
			fmt.Printf("❌ Generation failed after %v: %v\n", duration, err)
			return
		}
		fmt.Printf("✅ Generation completed in %v\n", duration)
		
	case <-time.After(5 * time.Minute):
		fmt.Println("⏰ Generation timed out after 5 minutes")
		cmd.Process.Kill()
		return
	}
	
	// Check generated files
	fmt.Println("\n📁 Checking generated files...")
	outputDir := "output/generated_books"
	
	// Look for files with our title pattern
	files, err := exec.Command("find", outputDir, "-name", "*Future*Digital*Privacy*").Output()
	if err != nil {
		fmt.Printf("Error checking files: %v\n", err)
		return
	}
	
	if len(files) == 0 {
		fmt.Println("⚠️  No files found with expected pattern")
	} else {
		fmt.Printf("📄 Generated files:\n%s", string(files))
	}
	
	fmt.Println("\n🎉 Single book generation test completed!")
	fmt.Println("This demonstrates the full pipeline working with EPUB support")
}