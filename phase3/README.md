# Phase 3: Reader Feedback & Shareability Analysis

## Overview
Simulates reader feedback from different personas and analyzes viral potential of the concept.

## Input
- `../phase2/concepts.json` (from Phase 2)

## Output
- `feedback.json` - Reader feedback, viral quotes, and engagement scores

## Usage
```bash
cd phase3
go run main.go
```

## Process
1. Load best concept from Phase 2
2. Generate feedback from diverse reader personas
3. Create potential viral quotes
4. Calculate engagement and shareability scores
5. Approve/reject for next phase based on thresholds

## Output Format
```json
{
  "timestamp": "2024-07-09T14:30:45Z",
  "concept_title": "AI Leadership in 2026 - Practical Framework",
  "reader_feedback": [
    {
      "persona": {
        "name": "Sarah",
        "age": 34,
        "demographics": "Tech Executive"
      },
      "rating": 8.2,
      "comment": "This concept really resonates...",
      "engagement": 85.0
    }
  ],
  "viral_quotes": [
    {
      "text": "AI leadership isn't about replacing humans...",
      "shareability_score": 87.5,
      "platform": "LinkedIn"
    }
  ],
  "overall_rating": 8.1,
  "engagement_score": 78.5,
  "shareability_score": 83.2,
  "approved_for_next": true
}
```

## Thresholds
- Overall Rating: ≥7.5/10
- Engagement Score: ≥75%
- Auto-approval if both thresholds met

## Next Phase
Pass `feedback.json` to Phase 4 for media coverage analysis.