# Phase 7: Marketing Assets Generation

## Overview
Generates comprehensive marketing materials for book launch and promotion.

## Input
- `../phase6/content.json` (from Phase 6) OR
- `../phase5/titles.json` (from Phase 5) as fallback

## Output
- `marketing.json` - Complete marketing campaign data
- `marketing_assets/` directory with all generated files

## Usage
```bash
cd phase7
go run main.go
```

## Generated Assets
1. **Press Release** - Media announcement
2. **Email Campaign** - HTML email template
3. **Blog Post** - Markdown blog content
4. **Product Description** - E-commerce copy
5. **Landing Page** - HTML landing page
6. **LinkedIn Article** - Professional content

## Social Media Posts
- Twitter/X posts with hashtags
- LinkedIn professional updates
- Facebook engagement posts  
- Instagram visual content

## Output Structure
```
phase7/
├── marketing.json          # Campaign metadata
└── marketing_assets/       # Generated files
    ├── press_release.txt
    ├── email_campaign.html
    ├── blog_post.md
    ├── product_description.txt
    ├── landing_page.html
    └── linkedin_article.md
```

## Output Format
```json
{
  "timestamp": "2024-07-09T14:30:45Z",
  "book_title": "AI Leadership in 2026: What Every Leader Needs to Know",
  "marketing_assets": [
    {
      "type": "Press Release",
      "title": "New Book Launches...",
      "content": "Full press release text...",
      "platform": "Media Outlets",
      "filename": "press_release.txt",
      "word_count": 450,
      "target_audience": "Media & Journalists"
    }
  ],
  "social_media_posts": [
    {
      "platform": "Twitter",
      "content": "🚀 Excited to share...",
      "hashtags": ["#AILeadership", "#Leadership2026"],
      "call_to_action": "Pre-order now",
      "char_count": 245
    }
  ],
  "created_files": ["marketing_assets/press_release.txt", ...],
  "asset_count": 6,
  "total_word_count": 2850,
  "campaign_ready": true
}
```

## Features
- Multiple content formats (HTML, Markdown, Plain Text)
- Platform-optimized content lengths
- SEO-friendly copy
- Call-to-action integration
- Professional copywriting style

## Target Audiences
- Media & Journalists
- Business Executives  
- Industry Professionals
- Book Buyers
- Professional Networks

## Campaign Ready Criteria
- Minimum 5 marketing assets generated
- Multiple platform coverage
- Complete social media suite

## File Usage
- Use generated HTML files for web deployment
- Copy social media posts directly to platforms
- Send press release to media contacts
- Use product description for book retailers