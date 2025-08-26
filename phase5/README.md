# Phase 5: Title Optimization Loop

## Overview
Optimizes the content title through A/B testing and multiple scoring criteria.

## Input
- `../phase4/media.json` (from Phase 4)

## Output
- `titles.json` - Title variations, A/B test results, and optimized title

## Usage
```bash
cd phase5
go run main.go
```

## Process
1. Load concept title from Phase 4
2. Generate multiple title variations
3. Score each variation on multiple criteria
4. Run simulated A/B tests
5. Select optimized title based on performance

## Scoring Criteria
- **Clickability**: How likely to generate clicks
- **SEO Score**: Search engine optimization potential
- **Memorability**: How easy to remember
- **Clarity**: How clear the value proposition is

## Output Format
```json
{
  "timestamp": "2024-07-09T14:30:45Z",
  "original_title": "AI Leadership in 2026 - Practical Framework",
  "title_variations": [
    {
      "title": "The Executive's Guide to AI Leadership in 2026",
      "clickability": 8.1,
      "seo_score": 7.5,
      "memorability": 7.8,
      "clarity": 8.4,
      "overall_score": 8.0,
      "test_audience": "Executives"
    }
  ],
  "ab_test_results": [
    {
      "title_a": "Original Title",
      "title_b": "Optimized Title",
      "click_rate_a": 8.6,
      "click_rate_b": 12.2,
      "winner": "Optimized Title",
      "confidence_level": 92.3
    }
  ],
  "optimized_title": "AI Leadership in 2026: What Every Leader Needs to Know",
  "improvement_score": 15.2,
  "final_confidence": 89.1,
  "approved_for_next": true
}
```

## A/B Testing
- Simulates real-world click rates
- Tests multiple variations against each other
- Calculates statistical confidence levels
- Identifies winning titles with high confidence

## Thresholds
- Final Confidence: ≥85% for next phase approval
- Minimum improvement for optimization recommendation

## Next Phase
Pass `titles.json` to Phase 6 for content generation with optimized title.