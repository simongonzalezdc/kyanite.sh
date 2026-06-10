// Package app provides application services for the noise.sh application.
// It includes services for AI assistance, auto-saving, theory analysis, and editor functionality.
package app

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/kyanite/noise/internal/app/ai"
	"github.com/kyanite/noise/internal/config"
	"github.com/kyanite/noise/internal/constants"
	"github.com/kyanite/noise/internal/domain"
	"github.com/kyanite/noise/internal/infra/brain"
	"github.com/kyanite/noise/internal/infra/glm"
)

// AIService handles AI-powered assistance
type AIService struct {
	config       *config.Config
	quickAgent   *ai.QuickIdeaAgent
	brainClient  *brain.Client
	glmClient    *glm.Client
}

// NewAIService creates a new AI service
func NewAIService(cfg *config.Config) *AIService {
	s := &AIService{
		config: cfg,
	}

	// Initialize infrastructure clients
	s.brainClient = brain.NewClient()
	glmTimeout := 60 * time.Second
	s.glmClient = glm.NewClient(cfg.GLM.APIKey, glmTimeout)

	// Configure QuickIdeaAgent based on provider
	var client ai.QuickLLMClient
	var model string

	switch cfg.AI.Provider {
	case "glm":
		client = s.glmClient
		model = cfg.GLM.Model
	case "hybrid":
		// QuickIdeaAgent uses Brain for fast feedback
		client = s.brainClient
		model = cfg.AI.Model
	default: // "ollama" — now routes through Brain
		client = s.brainClient
		model = cfg.AI.Model
	}

	s.quickAgent = ai.NewQuickIdeaAgent().
		WithClient(client, cfg.AI.Timeout).
		WithModel(model)

	return s
}


// BrainClient returns the underlying brain client for session/memory operations.
func (s *AIService) BrainClient() *brain.Client {
	return s.brainClient
}

// GetQuickAgent returns the underlying QuickIdeaAgent
func (s *AIService) GetQuickAgent() *ai.QuickIdeaAgent {
	return s.quickAgent
}

// Brainstorm generates creative angles for a song
func (s *AIService) Brainstorm(theme string) ([]string, error) {
	// If provider is GLM or Hybrid, use GLM for high-quality structural brainstorm
	if s.config.AI.Provider == "glm" || s.config.AI.Provider == "hybrid" {
		ctx, cancel := context.WithTimeout(context.Background(), constants.BrainstormTimeout)
		defer cancel()

		prompt := fmt.Sprintf("Generate 5 unique and evocative songwriting angles for the theme: '%s'. Return as a plain list of bullet points.", theme)
		resp, err := s.glmClient.Generate(ctx, s.config.GLM.Model, prompt, nil)
		if err == nil {
			lines := strings.Split(resp, "\n")
			var angles []string
			for _, line := range lines {
				line = strings.TrimSpace(line)
				if line != "" {
					angles = append(angles, strings.TrimPrefix(line, "- "))
				}
			}
			if len(angles) > 0 {
				return angles, nil
			}
		}
	}

	// Fallback to placeholder/local logic if GLM fails or is disabled
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
func (s *AIService) GenerateSection(sectionType domain.SectionType, contentContext string) (string, error) {
	// In a full implementation, this would use Ollama or GLM with specific prompts
	// For now, using placeholders as we focus on the Rapid Prototyping (QuickIdeaAgent)

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
	refined := "Refined: " + lyrics + " (based on: " + feedback + ")"
	return refined, nil
}

// AnalyzeQuality analyzes the quality of lyrics
func (s *AIService) AnalyzeQuality(song *domain.Song) (*domain.QualityScore, error) {
	return song.CalculateQualityScore(), nil
}

// Chat handles conversational AI interaction
// Implementation moved to ai_chat.go for better organization
// This method now serves as a clean entry point

// IsAvailable checks if AI services are available
func (s *AIService) IsAvailable() bool {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	return s.brainClient.IsAvailable(ctx)
}

// GetModelStatus returns the status of AI models
func (s *AIService) GetModelStatus() map[string]interface{} {
	return s.brainClient.GetModelStatus()
}

// Rapid prototyping methods (Delegated to QuickIdeaAgent)

func (s *AIService) GenerateContinuations(ctx context.Context, previousLines []string, sectionType string) ([]string, error) {
	content := strings.Join(previousLines, "\n")
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

func (s *AIService) GenerateVariations(ctx context.Context, line, section, constraint string) ([]string, error) {
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

func (s *AIService) GenerateRapidBrainstorm(ctx context.Context, theme string, maxAngles int) ([]string, error) {
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
	suggestions := resp.Suggestions
	if len(suggestions) > maxAngles && maxAngles > 0 {
		return suggestions[:maxAngles], nil
	}
	return suggestions, nil
}

func (s *AIService) GenerateOpeningLine(ctx context.Context, theme, angle string) (string, error) {
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

func (s *AIService) CheckQualityRedFlags(ctx context.Context, content string) (*ai.QuickResponse, error) {
	return s.quickAgent.Generate(ctx, ai.QuickRequest{
		Mode:    ai.QuickIdeaModeCheck,
		Context: content,
		Options: map[string]string{
			"mode": "sketch",
		},
	})
}

func (s *AIService) CheckQualityByMode(ctx context.Context, content, mode string) (*ai.QuickResponse, error) {
	return s.quickAgent.Generate(ctx, ai.QuickRequest{
		Mode:    ai.QuickIdeaModeCheck,
		Context: content,
		Options: map[string]string{
			"mode": mode,
		},
	})
}

// GenerateHarmonySuggestions generates chord progressions based on a mood
func (s *AIService) GenerateHarmonySuggestions(ctx context.Context, mood string) ([]string, error) {
	resp, err := s.quickAgent.Generate(ctx, ai.QuickRequest{
		Mode:    ai.QuickIdeaModeHarmony,
		Context: mood,
		Options: map[string]string{
			"theme": mood,
		},
	})
	if err != nil {
		return nil, err
	}
	return resp.Suggestions, nil
}

// FindRhymes finds words that rhyme with the given word
func (s *AIService) FindRhymes(word string) []string {
	// Use the basic rhyme finder from the ai package
	return ai.FindBasicRhymes(word)
}
