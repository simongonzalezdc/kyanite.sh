// Package ai provides AI-assisted songwriting services such as rhyme finding and idea generation.
package ai

import (
	"context"
	"fmt"
)

// ContinuationAgent handles AI-powered line continuation
type ContinuationAgent struct {
	// In a full implementation, this would include:
	// - Ollama client
	// - Prompt templates
	// - Configuration for temperature, max tokens, etc.
}

// NewContinuationAgent creates a new continuation agent
func NewContinuationAgent() *ContinuationAgent {
	return &ContinuationAgent{}
}

// ContinuationRequest represents a request for line continuation
type ContinuationRequest struct {
	PreviousLines []string `json:"previous_lines"`
	SectionType   string   `json:"section_type"`
	RhymeScheme   string   `json:"rhyme_scheme"`
	Syllables     int      `json:"syllables"`
	Style         string   `json:"style"`
}

// ContinuationResponse represents the response from the continuation agent
type ContinuationResponse struct {
	Lines []string `json:"lines"`
}

// GenerateContinuations generates line continuations based on context
func (a *ContinuationAgent) GenerateContinuations(ctx context.Context, req *ContinuationRequest) (*ContinuationResponse, error) {
	// For now, we'll use placeholder responses
	// In a full implementation, this would:
	// 1. Build a prompt from the request
	// 2. Call Ollama with the prompt
	// 3. Parse the JSON response
	// 4. Return the structured response

	// Placeholder implementation
	suggestions := []string{
		"Continue with this line...",
		"Or try this alternative...",
		"Perhaps this direction...",
	}

	// Adjust suggestions based on context if available
	if len(req.PreviousLines) > 0 {
		lastLine := req.PreviousLines[len(req.PreviousLines)-1]
		suggestions[0] = fmt.Sprintf("Continuing from: %s", lastLine)
	}

	return &ContinuationResponse{
		Lines: suggestions,
	}, nil
}

// ContinuationPrompt represents the prompt template for continuation
const ContinuationPrompt = `Continue this lyric with ONE natural next line.

Context:
{{.PreviousLines}}

Current section: {{.SectionType}}
Rhyme scheme: {{.RhymeScheme}}
Style: {{.Style}}

Generate 3 variations. Each must:
- Flow naturally from the last line
- Match the meter ({{.Syllables}} syllables ±1)
- Include specific, sensory detail
- Avoid abstract language

Return only JSON:
{
  "lines": [
    "line1",
    "line2",
    "line3"
  ]
}

No explanations. Just the lines.`
