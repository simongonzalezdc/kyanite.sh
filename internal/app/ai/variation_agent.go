package ai

import (
	"context"
	"fmt"
)

// VariationAgent handles AI-powered line variation generation
type VariationAgent struct {
	// In a full implementation, this would include:
	// - Ollama client
	// - Prompt templates
	// - Configuration for temperature, max tokens, etc.
}

// NewVariationAgent creates a new variation agent
func NewVariationAgent() *VariationAgent {
	return &VariationAgent{}
}

// VariationRequest represents a request for line variations
type VariationRequest struct {
	Line       string `json:"line"`
	Section    string `json:"section"`
	Constraint string `json:"constraint"`
	Syllables  int    `json:"syllables"`
}

// VariationResponse represents the response from the variation agent
type VariationResponse struct {
	Variations []string `json:"variations"`
	Reasoning  string   `json:"reasoning"`
}

// GenerateVariations generates line variations based on the original line
func (a *VariationAgent) GenerateVariations(ctx context.Context, req *VariationRequest) (*VariationResponse, error) {
	// For now, we'll use placeholder responses
	// In a full implementation, this would:
	// 1. Build a prompt from the request
	// 2. Call Ollama with the prompt
	// 3. Parse the JSON response
	// 4. Return the structured response
	
	// Placeholder implementation
	variations := []string{
		fmt.Sprintf("Variation 1 of: %s", req.Line),
		fmt.Sprintf("Variation 2 of: %s", req.Line),
		fmt.Sprintf("Variation 3 of: %s", req.Line),
	}
	
	// Adjust variations based on constraint if provided
	if req.Constraint != "" {
		switch req.Constraint {
		case "more concrete":
			variations[0] = fmt.Sprintf("More concrete version of: %s", req.Line)
		case "less abstract":
			variations[1] = fmt.Sprintf("Less abstract version of: %s", req.Line)
		case "different POV":
			variations[2] = fmt.Sprintf("Different POV version of: %s", req.Line)
		}
	}
	
	return &VariationResponse{
		Variations: variations,
		Reasoning:  "Generated 3 variations with different approaches",
	}, nil
}

// VariationPrompt represents the prompt template for variation
const VariationPrompt = `Rewrite this line 3 ways.

Original: {{.Line}}
Section: {{.Section}}
Constraint: {{.Constraint}}

Each variation must:
- Convey the same core emotion
- Use different words/structure
- Maintain meter ({{.Syllables}} syllables ±1)
- Add concrete details

If constraint is "more concrete": Add specific sensory details
If constraint is "less abstract": Replace abstractions with images
If constraint is "different POV": Change perspective (1st/2nd/3rd person)

Return JSON only:
{
  "variations": [
    "variation1",
    "variation2",
    "variation3"
  ],
  "reasoning": "brief explanation of approach"
}`