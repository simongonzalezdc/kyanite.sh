package ai

import (
	"context"
	"fmt"
)

// RapidBrainstormAgent handles AI-powered rapid brainstorming
type RapidBrainstormAgent struct {
	// In a full implementation, this would include:
	// - Ollama client
	// - Prompt templates
	// - Configuration for temperature, max tokens, etc.
}

// NewRapidBrainstormAgent creates a new rapid brainstorm agent
func NewRapidBrainstormAgent() *RapidBrainstormAgent {
	return &RapidBrainstormAgent{}
}

// RapidBrainstormRequest represents a request for rapid brainstorming
type RapidBrainstormRequest struct {
	Theme     string `json:"theme"`
	MaxAngles int    `json:"max_angles"`
}

// RapidBrainstormResponse represents the response from the rapid brainstorm agent
type RapidBrainstormResponse struct {
	Angles []string `json:"angles"`
}

// GenerateBrainstorm generates brainstorming angles for a theme
func (a *RapidBrainstormAgent) GenerateBrainstorm(ctx context.Context, req *RapidBrainstormRequest) (*RapidBrainstormResponse, error) {
	// For now, we'll use placeholder responses
	// In a full implementation, this would:
	// 1. Build a prompt from the request
	// 2. Call Ollama with the prompt
	// 3. Parse the JSON response
	// 4. Return the structured response

	// Placeholder implementation - generate 3 angles max
	maxAngles := req.MaxAngles
	if maxAngles <= 0 || maxAngles > 3 {
		maxAngles = 3
	}

	angles := []string{
		fmt.Sprintf("Explore %s through personal memories", req.Theme),
		fmt.Sprintf("Use nature imagery to symbolize %s", req.Theme),
		fmt.Sprintf("Focus on sensory details related to %s", req.Theme),
	}

	// Trim to max angles
	if len(angles) > maxAngles {
		angles = angles[:maxAngles]
	}

	return &RapidBrainstormResponse{
		Angles: angles,
	}, nil
}

// GenerateOpeningLine generates an opening line based on a selected angle
func (a *RapidBrainstormAgent) GenerateOpeningLine(ctx context.Context, theme, angle string) (string, error) {
	// For now, we'll use a placeholder response
	// In a full implementation, this would:
	// 1. Build a prompt from the theme and angle
	// 2. Call Ollama with the prompt
	// 3. Return the generated opening line

	return fmt.Sprintf("Opening line for %s: %s", theme, angle), nil
}

// RapidBrainstormPrompt represents the prompt template for rapid brainstorming
const RapidBrainstormPrompt = `You are helping a songwriter start writing FAST.

Theme: {{.Theme}}

Generate EXACTLY 3 specific angles. Each must:
- Be one sentence
- Include a concrete sensory detail
- Avoid clichés
- Suggest emotional direction

JSON format only:
{
  "angles": ["angle1", "angle2", "angle3"]
}`

// OpeningLinePrompt represents the prompt template for generating opening lines
const OpeningLinePrompt = `Write an opening line for a song.

Theme: {{.Theme}}
Angle: {{.Angle}}

Requirements:
- One line only
- Concrete imagery
- No clichés
- 8-12 syllables

Return only the line.`
