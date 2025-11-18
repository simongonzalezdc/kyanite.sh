// Package app provides application services for the noise.sh application.
// It includes services for AI assistance, auto-saving, theory analysis, and editor functionality.
package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/Kyanite/noise/internal/app/ai"
	"github.com/Kyanite/noise/internal/domain"
)

// AIService handles AI-powered assistance
type AIService struct {
	quickAgent *ai.QuickIdeaAgent
}

// NewAIService creates a new AI service
func NewAIService() *AIService {
	return &AIService{
		quickAgent: ai.NewQuickIdeaAgent(),
	}
}

// Brainstorm generates creative angles for a song
func (s *AIService) Brainstorm(theme string) ([]string, error) {
	// Placeholder implementation
	// In a full implementation, this would use Ollama and the RAG system

	angles := []string{
		"Explore the theme through personal memories and nostalgia",
		"Use nature imagery to symbolize the emotional journey",
		"Focus on sensory details to make the theme more vivid",
		"Contrast the theme with unexpected elements for surprise",
		"Build the theme through a narrative story arc",
	}

	return angles, nil
}

// GenerateSection generates lyrics for a specific section
func (s *AIService) GenerateSection(sectionType domain.SectionType, context string) (string, error) {
	// Placeholder implementation
	// In a full implementation, this would use Ollama with specific prompts

	sections := map[domain.SectionType]string{
		domain.SectionVerse:     "Sample verse lyrics would go here...",
		domain.SectionChorus:    "Sample chorus lyrics would go here...",
		domain.SectionBridge:    "Sample bridge lyrics would go here...",
		domain.SectionPreChorus: "Sample pre-chorus lyrics would go here...",
		domain.SectionIntro:     "Sample intro lyrics would go here...",
		domain.SectionOutro:     "Sample outro lyrics would go here...",
	}

	if lyrics, exists := sections[sectionType]; exists {
		return lyrics, nil
	}

	return "Sample lyrics would go here...", nil
}

// RefineLyrics refines existing lyrics based on feedback
func (s *AIService) RefineLyrics(lyrics, feedback string) (string, error) {
	// Placeholder implementation
	// In a full implementation, this would use Ollama to refine the lyrics

	refined := "Refined: " + lyrics + " (based on: " + feedback + ")"
	return refined, nil
}

// AnalyzeQuality analyzes the quality of lyrics
func (s *AIService) AnalyzeQuality(song *domain.Song) (*domain.QualityScore, error) {
	// Use the domain's built-in quality calculation
	return song.CalculateQualityScore(), nil
}

// Chat handles conversational AI interaction
// Chat is now implemented in ai_chat.go
// Kept here for backward compatibility, delegates to new implementation
func (s *AIService) Chat(ctx context.Context, message string) (<-chan string, error) {
	// Implementation moved to ai_chat.go for better organization
	// This method now checks availability and streams from Ollama
	return s.Chat(ctx, message)
}

// IsAvailable checks if AI services (Ollama) are available
func (s *AIService) IsAvailable() bool {
	// Check if Ollama is reachable
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(s.ollamaURL + "/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

// GetModelStatus returns the status of AI models
func (s *AIService) GetModelStatus() map[string]interface{} {
	// Placeholder implementation
	return map[string]interface{}{
		"ollama_running": false,
		"models_loaded":  []string{},
		"last_check":     time.Now(),
	}
}

// Rapid prototyping methods

// GenerateContinuations generates line continuations for rapid prototyping
func (s *AIService) GenerateContinuations(ctx context.Context, previousLines []string, sectionType string) ([]string, error) {
	content := strings.Join(previousLines, "\n")

	// Detect content type for context-aware suggestions
	if s.quickAgent != nil {
		resp, err := s.quickAgent.Generate(ctx, ai.QuickRequest{
			Mode:    ai.QuickIdeaModeUnstick,
			Context: content,
			Options: map[string]string{
				"section": sectionType,
			},
		})
		if err != nil {
			return nil, err
		}
		return resp.Suggestions, nil
	}

	// Fallback to basic suggestions
	return []string{
		content + " and the story continues",
		content + " while the music plays on",
		content + " as the rhythm takes hold",
	}, nil
}

// GenerateVariations generates line variations for rapid prototyping
func (s *AIService) GenerateVariations(ctx context.Context, line, section, constraint string) ([]string, error) {
	// Detect content type for context-aware suggestions
	if s.quickAgent != nil {
		resp, err := s.quickAgent.Generate(ctx, ai.QuickRequest{
			Mode:    ai.QuickIdeaModeTweak,
			Context: line,
			Options: map[string]string{
				"section":    section,
				"constraint": constraint,
			},
		})
		if err != nil {
			return nil, err
		}
		return resp.Suggestions, nil
	}

	// Fallback to basic variations
	return []string{
		line,
		"Rewrite with stronger imagery",
		"Add more emotional depth",
	}, nil
}

// GenerateRapidBrainstorm generates brainstorming angles for rapid prototyping
func (s *AIService) GenerateRapidBrainstorm(ctx context.Context, theme string, maxAngles int) ([]string, error) {
	// Detect content type for context-aware suggestions
	if s.quickAgent != nil {
		resp, err := s.quickAgent.Generate(ctx, ai.QuickRequest{
			Mode:    ai.QuickIdeaModeSpark,
			Context: theme,
			Options: map[string]string{
				"theme": theme,
				"limit": fmt.Sprintf("%d", maxAngles),
			},
		})
		if err != nil {
			return nil, err
		}
		if len(resp.Suggestions) > maxAngles && maxAngles > 0 {
			return resp.Suggestions[:maxAngles], nil
		}
		return resp.Suggestions, nil
	}

	// Fallback to basic brainstorming
	suggestions := []string{
		"Explore " + theme + " through personal memories",
		"Use nature imagery to symbolize " + theme,
		"Focus on sensory details related to " + theme,
	}
	if len(suggestions) > maxAngles && maxAngles > 0 {
		return suggestions[:maxAngles], nil
	}
	return suggestions, nil
}

// GenerateOpeningLine generates an opening line based on a selected angle
func (s *AIService) GenerateOpeningLine(ctx context.Context, theme, angle string) (string, error) {
	// Detect content type for context-aware suggestions
	if s.quickAgent != nil {
		resp, err := s.quickAgent.Generate(ctx, ai.QuickRequest{
			Mode:    ai.QuickIdeaModeSpark,
			Context: theme,
			Options: map[string]string{
				"theme": theme,
				"angle": angle,
			},
		})
		if err != nil {
			return "", err
		}
		if len(resp.Suggestions) == 0 {
			return "", nil
		}
		return resp.Suggestions[0], nil
	}

	// Fallback to basic opening line
	return "In the heart of " + theme + ", a story begins", nil
}

// Tiered quality checking methods

// CheckQualityRedFlags performs red flag checking for sketch mode
func (s *AIService) CheckQualityRedFlags(ctx context.Context, content string) (*ai.QuickResponse, error) {
	if s.quickAgent != nil {
		return s.quickAgent.Generate(ctx, ai.QuickRequest{
			Mode:    ai.QuickIdeaModeCheck,
			Context: content,
			Options: map[string]string{
				"mode": "sketch",
			},
		})
	}

	// Fallback to basic quality check
	return &ai.QuickResponse{
		Rating: "OKAY",
		Tip:    "Add more specific details",
	}, nil
}

// CheckQualityBasic performs basic quality checking for draft mode
func (s *AIService) CheckQualityBasic(ctx context.Context, content string) (*ai.QuickResponse, error) {
	if s.quickAgent != nil {
		return s.quickAgent.Generate(ctx, ai.QuickRequest{
			Mode:    ai.QuickIdeaModeCheck,
			Context: content,
			Options: map[string]string{
				"mode": "draft",
			},
		})
	}

	// Fallback to basic quality check
	return &ai.QuickResponse{
		Rating: "OKAY",
		Tip:    "Strengthen the main idea",
	}, nil
}

// CheckQualityFull performs full quality checking for polish mode
func (s *AIService) CheckQualityFull(ctx context.Context, content string) (*ai.QuickResponse, error) {
	if s.quickAgent != nil {
		return s.quickAgent.Generate(ctx, ai.QuickRequest{
			Mode:    ai.QuickIdeaModeCheck,
			Context: content,
			Options: map[string]string{
				"mode": "polish",
			},
		})
	}

	// Fallback to basic quality check
	return &ai.QuickResponse{
		Rating: "OKAY",
		Tip:    "Refine the emotional impact",
	}, nil
}

// CheckQualityByMode performs quality checking based on the editor mode
func (s *AIService) CheckQualityByMode(ctx context.Context, content, mode string) (*ai.QuickResponse, error) {
	if s.quickAgent != nil {
		return s.quickAgent.Generate(ctx, ai.QuickRequest{
			Mode:    ai.QuickIdeaModeCheck,
			Context: content,
			Options: map[string]string{
				"mode": mode,
			},
		})
	}

	// Fallback to basic quality check
	return &ai.QuickResponse{
		Rating: "OKAY",
		Tip:    "Continue developing this idea",
	}, nil
}
