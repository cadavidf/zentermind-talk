# Phase 4: Media Coverage Prediction

## Overview
Analyzes potential media coverage and PR opportunities for the approved concept.

## Input
- `../phase3/feedback.json` (from Phase 3)

## Output
- `media.json` - Media analysis, coverage estimates, and PR assessment

## Usage
```bash
cd phase4
go run main.go
```

## Process
1. Load concept feedback from Phase 3
2. Analyze media outlets and their audience fit
3. Calculate coverage likelihood based on concept scores
4. Estimate potential reach and timeframes
5. Generate overall media score and PR worthiness

## Output Format
```json
{
  "timestamp": "2024-07-09T14:30:45Z",
  "concept_title": "AI Leadership in 2026 - Practical Framework",
  "media_analysis": [
    {
      "outlet": {
        "name": "Harvard Business Review",
        "type": "Business Magazine",
        "audience": "Executives",
        "reach": 800000
      },
      "likelihood": 78.5,
      "timeframe": "2-4 months",
      "coverage_type": "Feature Article",
      "estimated_reach": 628000
    }
  ],
  "total_estimated_reach": 2450000,
  "media_score": 72.3,
  "pr_worthiness": "Good - Moderate media interest",
  "approved_for_next": true
}
```

## Media Outlets Analyzed
- TechCrunch (Tech Blog)
- Harvard Business Review (Business Magazine)
- Wired (Tech Magazine)
- Forbes (Business Magazine)
- MIT Technology Review (Research)
- Fast Company (Innovation)

## Scoring Criteria
- Engagement Score → Likelihood boost
- Shareability Score → Media appeal
- Overall Rating → Quality indicator
- Audience fit → Outlet-specific adjustment

## Thresholds
- Media Score: ≥70/100 for next phase approval

## Next Phase
Pass `media.json` to Phase 5 for title optimization.