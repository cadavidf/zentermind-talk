package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Phase6Output struct {
	Timestamp        string   `json:"timestamp"`
	ConceptTitle     string   `json:"concept_title"`
	ContentFiles     []string `json:"content_files"`
	WordCount        int      `json:"word_count"`
	QuotableCount    int      `json:"quotable_count"`
	GenerationComplete bool   `json:"generation_complete"`
}

type MarketingAsset struct {
	Type        string `json:"type"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	Platform    string `json:"platform"`
	Filename    string `json:"filename"`
	WordCount   int    `json:"word_count"`
	TargetAudience string `json:"target_audience"`
}

type SocialMediaPost struct {
	Platform    string `json:"platform"`
	Content     string `json:"content"`
	Hashtags    []string `json:"hashtags"`
	CallToAction string `json:"call_to_action"`
	CharCount   int    `json:"char_count"`
}

type Phase7Output struct {
	Timestamp         string            `json:"timestamp"`
	BookTitle         string            `json:"book_title"`
	MarketingAssets   []MarketingAsset  `json:"marketing_assets"`
	SocialMediaPosts  []SocialMediaPost `json:"social_media_posts"`
	CreatedFiles      []string          `json:"created_files"`
	AssetCount        int               `json:"asset_count"`
	TotalWordCount    int               `json:"total_word_count"`
	CampaignReady     bool              `json:"campaign_ready"`
}

func main() {
	fmt.Println("🚀 PHASE 7: Marketing Assets Generation")
	fmt.Println(strings.Repeat("=", 50))
	
	// Try to load Phase 6 output or use Phase 5 for title
	bookTitle, err := getBookTitle()
	if err != nil {
		fmt.Printf("Error getting book title: %v\n", err)
		bookTitle = "AI Leadership in 2026: What Every Leader Needs to Know"
		fmt.Printf("Using default title: %s\n", bookTitle)
	}
	
	fmt.Printf("\n📖 Generating marketing assets for: %s\n", bookTitle)
	
	// Create output directory
	outputDir := "marketing_assets"
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		fmt.Printf("Error creating output directory: %v\n", err)
		return
	}
	
	// Generate marketing assets
	assets := generateMarketingAssets(bookTitle)
	socialPosts := generateSocialMediaPosts(bookTitle)
	
	// Save assets to files
	createdFiles := saveMarketingAssets(assets, outputDir)
	totalWordCount := calculateTotalWordCount(assets)
	
	output := Phase7Output{
		Timestamp:         time.Now().Format(time.RFC3339),
		BookTitle:         bookTitle,
		MarketingAssets:   assets,
		SocialMediaPosts:  socialPosts,
		CreatedFiles:      createdFiles,
		AssetCount:        len(assets),
		TotalWordCount:    totalWordCount,
		CampaignReady:     len(assets) >= 5,
	}
	
	// Save to JSON file
	outputFile := "marketing.json"
	if err := saveToJSON(output, outputFile); err != nil {
		fmt.Printf("Error saving output: %v\n", err)
		return
	}
	
	fmt.Printf("\n✅ Marketing assets generation complete\n")
	fmt.Printf("📁 Output saved to: %s\n", outputFile)
	fmt.Printf("📂 Assets directory: %s\n", outputDir)
	fmt.Printf("\n📊 Results:\n")
	fmt.Printf("   Assets Created: %d\n", len(assets))
	fmt.Printf("   Social Posts: %d\n", len(socialPosts))
	fmt.Printf("   Total Words: %d\n", totalWordCount)
	fmt.Printf("   Campaign Ready: %v\n", output.CampaignReady)
	
	fmt.Printf("\n📄 Generated Assets:\n")
	for i, asset := range assets {
		if i >= 5 { break }
		fmt.Printf("   %s: %s (%d words)\n", asset.Type, asset.Title, asset.WordCount)
	}
	
	fmt.Printf("\n💬 Sample Social Posts:\n")
	for i, post := range socialPosts {
		if i >= 2 { break }
		fmt.Printf("   %s: %.50s... (%d chars)\n", post.Platform, post.Content, post.CharCount)
	}
}

func getBookTitle() (string, error) {
	// Try Phase 6 first
	if phase6Data, err := loadPhase6Output("../phase6/content.json"); err == nil {
		return phase6Data.ConceptTitle, nil
	}
	
	// Try Phase 5
	if phase5Data, err := loadPhase5Output("../phase5/titles.json"); err == nil {
		return phase5Data.OptimizedTitle, nil
	}
	
	return "", fmt.Errorf("no phase data found")
}

func loadPhase6Output(filename string) (Phase6Output, error) {
	var data Phase6Output
	file, err := os.Open(filename)
	if err != nil {
		return data, err
	}
	defer file.Close()
	
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	return data, err
}

func loadPhase5Output(filename string) (struct {
	OptimizedTitle string `json:"optimized_title"`
}, error) {
	var data struct {
		OptimizedTitle string `json:"optimized_title"`
	}
	file, err := os.Open(filename)
	if err != nil {
		return data, err
	}
	defer file.Close()
	
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&data)
	return data, err
}

func generateMarketingAssets(bookTitle string) []MarketingAsset {
	assets := []MarketingAsset{
		{
			Type:     "Press Release",
			Title:    fmt.Sprintf("New Book '%s' Launches to Transform Business Leadership", bookTitle),
			Content:  generatePressRelease(bookTitle),
			Platform: "Media Outlets",
			Filename: "press_release.txt",
			TargetAudience: "Media & Journalists",
		},
		{
			Type:     "Email Campaign",
			Title:    "Exclusive Preview: Revolutionary Leadership Guide",
			Content:  generateEmailCampaign(bookTitle),
			Platform: "Email Marketing",
			Filename: "email_campaign.html",
			TargetAudience: "Business Executives",
		},
		{
			Type:     "Blog Post",
			Title:    "5 Key Insights from " + bookTitle,
			Content:  generateBlogPost(bookTitle),
			Platform: "Company Blog",
			Filename: "blog_post.md",
			TargetAudience: "Industry Professionals",
		},
		{
			Type:     "Product Description",
			Title:    "Amazon Product Description",
			Content:  generateProductDescription(bookTitle),
			Platform: "E-commerce",
			Filename: "product_description.txt",
			TargetAudience: "Book Buyers",
		},
		{
			Type:     "Landing Page",
			Title:    "Book Launch Landing Page",
			Content:  generateLandingPage(bookTitle),
			Platform: "Website",
			Filename: "landing_page.html",
			TargetAudience: "Potential Readers",
		},
		{
			Type:     "LinkedIn Article",
			Title:    "Why " + bookTitle + " Matters for Today's Leaders",
			Content:  generateLinkedInArticle(bookTitle),
			Platform: "LinkedIn",
			Filename: "linkedin_article.md",
			TargetAudience: "Professional Network",
		},
	}
	
	// Calculate word counts
	for i := range assets {
		assets[i].WordCount = len([]rune(assets[i].Content)) / 5 // Rough word count
	}
	
	return assets
}

func generateSocialMediaPosts(bookTitle string) []SocialMediaPost {
	posts := []SocialMediaPost{
		{
			Platform:     "Twitter",
			Content:      fmt.Sprintf("🚀 Excited to share insights from '%s' - the future of leadership is here! What's your biggest AI leadership challenge?", bookTitle),
			Hashtags:     []string{"#AILeadership", "#Leadership2026", "#FutureOfWork", "#Innovation"},
			CallToAction: "Pre-order now",
			CharCount:    0,
		},
		{
			Platform:     "LinkedIn",
			Content:      fmt.Sprintf("After months of research, I'm thrilled to introduce '%s'. This book addresses the critical gap in preparing leaders for an AI-driven future. What resonates most with your leadership journey?", bookTitle),
			Hashtags:     []string{"#Leadership", "#ArtificialIntelligence", "#BusinessStrategy"},
			CallToAction: "Download free chapter",
			CharCount:    0,
		},
		{
			Platform:     "Facebook",
			Content:      fmt.Sprintf("📚 New book alert! '%s' is now available. Perfect for executives, managers, and anyone leading teams in our rapidly changing world. Early reviews are incredible!", bookTitle),
			Hashtags:     []string{"#NewBook", "#Leadership", "#ProfessionalDevelopment"},
			CallToAction: "Get your copy",
			CharCount:    0,
		},
		{
			Platform:     "Instagram",
			Content:      fmt.Sprintf("✨ The future of leadership starts today. '%s' - your guide to navigating the AI revolution. Swipe for key insights! 📖", bookTitle),
			Hashtags:     []string{"#BookLaunch", "#LeadershipBooks", "#AI", "#Success", "#Inspiration"},
			CallToAction: "Link in bio",
			CharCount:    0,
		},
	}
	
	// Calculate character counts
	for i := range posts {
		posts[i].CharCount = len(posts[i].Content)
	}
	
	return posts
}

func generatePressRelease(title string) string {
	return fmt.Sprintf(`FOR IMMEDIATE RELEASE

Revolutionary Leadership Guide '%s' Launches to Address Critical Skills Gap in AI Era

New book provides practical framework for executives navigating artificial intelligence transformation

[City, Date] - Business leaders worldwide are facing unprecedented challenges as artificial intelligence reshapes entire industries. A new book, '%s', offers the first comprehensive guide for executives seeking to lead effectively in this rapidly evolving landscape.

The book addresses the critical gap between traditional leadership approaches and the skills needed to guide organizations through AI transformation. Based on extensive research and real-world case studies, it provides actionable insights for leaders at all levels.

"The future belongs to leaders who can effectively integrate AI capabilities while maintaining human-centered values," says the author. "This book provides the roadmap for that integration."

Key features include:
- Practical frameworks for AI implementation
- Case studies from successful AI transformations  
- Tools for building AI-ready teams
- Strategies for maintaining competitive advantage

The book is available in multiple formats and has already received endorsements from leading business executives and AI researchers.

For more information, visit [website] or contact [contact information].

###`, title, title)
}

func generateEmailCampaign(title string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>Exclusive Preview: %s</title>
    <style>
        body { font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; }
        .header { background: #2c3e50; color: white; padding: 20px; text-align: center; }
        .content { padding: 20px; }
        .cta { background: #3498db; color: white; padding: 15px 30px; text-decoration: none; border-radius: 5px; }
    </style>
</head>
<body>
    <div class="header">
        <h1>The Future of Leadership is Here</h1>
    </div>
    <div class="content">
        <h2>Dear [First Name],</h2>
        
        <p>As a forward-thinking leader, you understand that artificial intelligence isn't just changing technology—it's revolutionizing how we lead, innovate, and compete.</p>
        
        <p>That's why I'm excited to share with you the exclusive preview of '<strong>%s</strong>'—the definitive guide for leaders navigating the AI transformation.</p>
        
        <p><strong>What you'll discover:</strong></p>
        <ul>
            <li>How to build AI-ready teams without losing human connection</li>
            <li>Practical frameworks for implementing AI in your organization</li>
            <li>Real case studies from successful AI transformations</li>
            <li>Future-proofing strategies for sustained competitive advantage</li>
        </ul>
        
        <p>The first 100 readers get access to exclusive bonus materials and a private mastermind group.</p>
        
        <center>
            <a href="#" class="cta">Get Your Copy Now</a>
        </center>
        
        <p>To your success,<br>[Your Name]</p>
    </div>
</body>
</html>`, title, title)
}

func generateBlogPost(title string) string {
	return fmt.Sprintf(`# 5 Key Insights from %s

The business world is experiencing a seismic shift. Artificial intelligence is no longer a distant future concept—it's reshaping industries, redefining roles, and revolutionizing how we think about leadership. 

After diving deep into '%s', here are the five insights that stand out as game-changers for today's leaders:

## 1. AI Leadership Requires Human-Centered Thinking

The most successful AI implementations don't replace human judgment—they amplify it. Leaders who understand this distinction create environments where technology enhances human capabilities rather than diminishing them.

## 2. Data Literacy is the New Business Literacy

Just as reading became essential for participation in society, data literacy has become crucial for effective leadership. Leaders who can interpret, question, and act on data insights will have a significant competitive advantage.

## 3. Agile Decision-Making Becomes Critical

AI systems can process information at unprecedented speeds, but they require human oversight for context and ethics. Leaders must develop the ability to make quick, informed decisions while maintaining strategic perspective.

## 4. Building AI-Ready Teams Starts with Culture

Technical implementation is only half the battle. The most successful AI transformations happen when leaders focus first on creating cultures of continuous learning, experimentation, and adaptation.

## 5. Ethics and Transparency Drive Trust

As AI becomes more prevalent, leaders who prioritize ethical implementation and transparent communication will build stronger stakeholder relationships and sustainable competitive advantages.

## The Path Forward

These insights represent just the beginning of the AI leadership journey. The leaders who invest in developing these capabilities now will be best positioned to thrive in an AI-driven future.

What resonates most with your leadership experience? Share your thoughts in the comments below.

---

*Ready to dive deeper? Get your copy of '%s' and join thousands of leaders preparing for the future.*`, title, title, title)
}

func generateProductDescription(title string) string {
	return fmt.Sprintf(`%s

★★★★★ "Essential reading for any executive serious about the future" - Business Week

THE DEFINITIVE GUIDE FOR LEADERS IN THE AI ERA

Are you prepared to lead in a world where artificial intelligence reshapes everything? 

This groundbreaking book provides the practical framework every leader needs to navigate the AI transformation successfully. Based on extensive research and real-world case studies, it delivers actionable insights that you can implement immediately.

🎯 WHAT YOU'LL LEARN:
• How to build AI-ready teams without sacrificing human connection
• Practical frameworks for ethical AI implementation
• Strategies for maintaining competitive advantage in an AI-driven market
• Real case studies from successful AI transformations
• Tools for making data-driven decisions with confidence

📈 WHO THIS IS FOR:
✓ CEOs and C-suite executives
✓ Department heads and team leaders
✓ Consultants and business advisors
✓ Anyone responsible for organizational transformation

🏆 EARLY PRAISE:
"This book bridges the gap between AI potential and practical implementation" - Tech Leadership Today

"Finally, a guide that focuses on the human side of AI transformation" - Harvard Business Review

💡 BONUS MATERIALS INCLUDED:
• Downloadable assessment tools
• Implementation checklists
• Access to exclusive online resources

📚 FORMAT OPTIONS:
• Hardcover: 312 pages
• Kindle: Available for immediate download
• Audiobook: 8 hours, narrated by the author

Order now and join thousands of leaders who are already transforming their organizations for the AI era.

⭐ 30-day money-back guarantee
🚚 Free shipping on orders over $25
📱 Instant digital delivery available`, title)
}

func generateLandingPage(title string) string {
	return fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>%s - Transform Your Leadership for the AI Era</title>
    <style>
        body { font-family: 'Helvetica', sans-serif; margin: 0; padding: 0; }
        .hero { background: linear-gradient(135deg, #667eea 0%%, #764ba2 100%%); color: white; text-align: center; padding: 100px 20px; }
        .hero h1 { font-size: 3em; margin-bottom: 20px; }
        .hero p { font-size: 1.5em; max-width: 600px; margin: 0 auto; }
        .cta-button { background: #ff6b6b; color: white; padding: 20px 40px; font-size: 1.2em; text-decoration: none; border-radius: 50px; margin: 30px; display: inline-block; }
        .section { padding: 80px 20px; max-width: 1200px; margin: 0 auto; }
        .benefits { display: grid; grid-template-columns: repeat(auto-fit, minmax(300px, 1fr)); gap: 40px; }
        .benefit { text-align: center; }
        .testimonial { background: #f8f9fa; padding: 40px; border-radius: 10px; margin: 40px 0; }
    </style>
</head>
<body>
    <div class="hero">
        <h1>%s</h1>
        <p>The definitive guide for leaders navigating the AI transformation</p>
        <a href="#order" class="cta-button">Get Your Copy Now</a>
    </div>
    
    <div class="section">
        <h2>Why This Book Matters</h2>
        <p>Artificial intelligence is reshaping business at unprecedented speed. Leaders who don't adapt will be left behind. This book provides the roadmap for successful AI transformation.</p>
        
        <div class="benefits">
            <div class="benefit">
                <h3>🎯 Practical Frameworks</h3>
                <p>Step-by-step guides for implementing AI in your organization</p>
            </div>
            <div class="benefit">
                <h3>📊 Real Case Studies</h3>
                <p>Learn from successful AI transformations across industries</p>
            </div>
            <div class="benefit">
                <h3>🚀 Future-Ready Strategies</h3>
                <p>Build sustainable competitive advantage in the AI era</p>
            </div>
        </div>
    </div>
    
    <div class="testimonial">
        <p>"This book bridges the gap between AI potential and practical implementation. Essential reading for any leader serious about the future."</p>
        <strong>- Sarah Johnson, CEO, TechCorp</strong>
    </div>
    
    <div class="section" id="order">
        <h2>Get Your Copy Today</h2>
        <p>Join thousands of leaders who are already transforming their organizations.</p>
        <a href="#" class="cta-button">Order Now - $24.99</a>
        <p><small>30-day money-back guarantee • Free shipping • Instant digital access</small></p>
    </div>
</body>
</html>`, title, title)
}

func generateLinkedInArticle(title string) string {
	return fmt.Sprintf(`# Why %s Matters for Today's Leaders

The conversation around artificial intelligence in business has shifted dramatically. We're no longer asking "if" AI will transform our industries—we're asking "how quickly" and "are we prepared?"

As someone who has spent years studying organizational transformation, I've observed a critical gap: while technology advances at breakneck speed, leadership development hasn't kept pace.

## The Leadership Gap

Traditional leadership training focuses on managing people, processes, and performance. But AI introduces new variables:

- **Decision velocity**: AI can analyze vast datasets in seconds, but still requires human judgment for context and ethics
- **Human-machine collaboration**: The most effective teams will blend human creativity with AI capabilities
- **Continuous adaptation**: AI systems learn and evolve, requiring leaders who can match that pace

## What Sets AI Leaders Apart

The leaders who will thrive in this new landscape share several characteristics:

**Data fluency**: They can interpret and act on AI-generated insights without being data scientists themselves.

**Ethical clarity**: They understand the importance of responsible AI implementation and can navigate complex ethical decisions.

**Change agility**: They view AI transformation as an ongoing journey, not a one-time project.

## The Path Forward

The question isn't whether your organization will adopt AI—it's whether you'll be ready to lead that transformation effectively.

This is why I'm excited about resources like '%s'. It provides the practical framework that leaders need to bridge the gap between AI potential and successful implementation.

The future belongs to leaders who can dance with artificial intelligence—leveraging its capabilities while maintaining the human touch that drives truly exceptional organizations.

What's your experience with AI leadership challenges? I'd love to hear your perspectives in the comments.

---

*What aspects of AI leadership are you most curious about? Let's continue this conversation.*`, title, title)
}

func saveMarketingAssets(assets []MarketingAsset, outputDir string) []string {
	var createdFiles []string
	
	for _, asset := range assets {
		fullPath := filepath.Join(outputDir, asset.Filename)
		
		if err := os.WriteFile(fullPath, []byte(asset.Content), 0644); err != nil {
			fmt.Printf("Error saving %s: %v\n", asset.Filename, err)
			continue
		}
		
		createdFiles = append(createdFiles, fullPath)
	}
	
	return createdFiles
}

func calculateTotalWordCount(assets []MarketingAsset) int {
	total := 0
	for _, asset := range assets {
		total += asset.WordCount
	}
	return total
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