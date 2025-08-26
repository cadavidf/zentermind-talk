# Phase 2: Concept Generation & Validation

## Overview
Takes the best recommendation from Phase 1 and generates multiple detailed concepts, validates them, and selects the best one.

## Input
- `../phase1/recommendations.json` (from Phase 1)

## Output
- `concepts.json` - Generated and validated concepts with best selection

## Usage
```bash
cd phase2
go run main.go
```

## Process
1. Load Phase 1 recommendations
2. Select highest-ranked recommendation
3. Generate multiple concept variations
4. Validate concepts (viability, uniqueness, commercial scores)
5. Select best concept using scoring algorithm

## Output Format
```json
{
  "timestamp": "2024-07-09T14:30:45Z",
  "selected_title": "AI Leadership in 2026",
  "generated_concepts": [...],
  "valid_concepts": [...],
  "best_concept": {
    "title": "AI Leadership in 2026 - Practical Framework",
    "description": "A hands-on guide...",
    "uniqueness_score": 8.2,
    "viability_score": true,
    "commercial_score": 7.8,
    "status": "approved"
  },
  "total_processed": 3
}
```

## Next Phase
Pass `concepts.json` to Phase 3 for reader feedback simulation.