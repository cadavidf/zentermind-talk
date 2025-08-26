# Phase 1: Content Recommendations

## Overview
Generates trending content recommendations based on market analysis and content type selection.

## Input
- Content type selection (book, podcast, meditation)

## Output
- `recommendations.json` - List of recommended content with confidence scores

## Usage
```bash
cd phase1
go run main.go
```

## Output Format
```json
{
  "timestamp": "2024-07-09T14:30:45Z",
  "content_type": "book",
  "recommendations": [
    {
      "title": "AI Leadership in 2026",
      "confidence": 85.5,
      "rank": 1,
      "market_gap_score": 7.2,
      "content_type": "book"
    }
  ],
  "total_generated": 8,
  "passing_threshold": 5
}
```

## Next Phase
Pass `recommendations.json` to Phase 2 for concept generation.