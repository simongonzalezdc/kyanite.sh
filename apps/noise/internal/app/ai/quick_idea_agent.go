package ai

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/kyanite/noise/internal/app/knowledge"
)

// QuickIdeaMode represents the supported AI assistance modes.
type QuickIdeaMode string

const (
	// QuickIdeaModeUnstick generates next-line suggestions based on recent context.
	QuickIdeaModeUnstick QuickIdeaMode = "unstick"
	// QuickIdeaModeSpark generates new opening ideas based on a theme keyword.
	QuickIdeaModeSpark QuickIdeaMode = "spark"
	// QuickIdeaModeTweak rewrites a single line in multiple ways.
	QuickIdeaModeTweak QuickIdeaMode = "tweak"
	// QuickIdeaModeCheck performs a lightweight quality check and returns a rating.
	QuickIdeaModeCheck QuickIdeaMode = "check"
	// QuickIdeaModeHarmony generates chord progressions based on a mood keyword.
	QuickIdeaModeHarmony QuickIdeaMode = "harmony"

	// defaultQuickIdeaTimeout configures the maximum time an AI request may take.
	defaultQuickIdeaTimeout = 1700 * time.Millisecond
	// defaultQuickIdeaModel documents the preferred lightweight model for this agent.
	defaultQuickIdeaModel = "qwen2.5:3b"

	// Retry configuration
	defaultMaxRetries     = 3
	defaultBaseRetryDelay = 100 * time.Millisecond
	defaultMaxRetryDelay  = 1000 * time.Millisecond
)

var supportedQuickIdeaModes = map[QuickIdeaMode]struct{}{
	QuickIdeaModeUnstick: {},
	QuickIdeaModeSpark:   {},
	QuickIdeaModeTweak:   {},
	QuickIdeaModeCheck:   {},
	QuickIdeaModeHarmony: {},
}

// QuickRequest encapsulates the data required to run a QuickIdeaAgent interaction.
type QuickRequest struct {
	Mode    QuickIdeaMode
	Context string
	Options map[string]string
}

// QuickResponse captures the structured result of an AI interaction.
type QuickResponse struct {
	Suggestions  []string
	Rating       string
	Tip          string
	ResponseTime time.Duration
}

// RetryConfig configures the retry behavior for AI requests.
type RetryConfig struct {
	MaxRetries     int
	BaseDelay      time.Duration
	MaxDelay       time.Duration
	RetryableCheck func(error) bool
}

// DefaultRetryConfig returns the default retry configuration.
func DefaultRetryConfig() RetryConfig {
	return RetryConfig{
		MaxRetries:     defaultMaxRetries,
		BaseDelay:      defaultBaseRetryDelay,
		MaxDelay:       defaultMaxRetryDelay,
		RetryableCheck: isRetryableError,
	}
}

// QuickIdeaAgent consolidates previously distinct AI helpers into a single entry point.
// By default it uses a stubbed client so the application can run without Ollama,
// while providing hooks to integrate the real qwen2.5:3b model when available.
type QuickIdeaAgent struct {
	client          QuickLLMClient
	model           string
	timeout         time.Duration
	prompts         quickIdeaPrompts
	contextDetector *ContextDetector
	contextPrompts  *ContextAwarePrompts
	knowledgeBase   knowledge.EnhancementProvider
	retryConfig     RetryConfig
}

// QuickLLMClient is a minimal interface that can be satisfied by the Ollama or GLM client.
type QuickLLMClient interface {
	Generate(ctx context.Context, model, prompt string, options map[string]any) (string, error)
}

// NewQuickIdeaAgent constructs an agent with a fallback stub client.
// Use WithClient to inject a real Ollama client at runtime.
func NewQuickIdeaAgent() *QuickIdeaAgent {
	contextDetector := NewContextDetector()
	contextPrompts := NewContextAwarePrompts()
	contextPrompts.Initialize()
	kbProvider := knowledge.NewStubEnhancementProvider()

	return &QuickIdeaAgent{
		client:          &stubQuickClient{},
		model:           defaultQuickIdeaModel,
		timeout:         defaultQuickIdeaTimeout,
		prompts:         defaultQuickIdeaPrompts(),
		contextDetector: contextDetector,
		contextPrompts:  contextPrompts,
		knowledgeBase:   kbProvider,
		retryConfig:     DefaultRetryConfig(),
	}
}

// WithClient returns a copy of the agent configured to use the provided LLM client and timeout.
func (a *QuickIdeaAgent) WithClient(client QuickLLMClient, timeout time.Duration) *QuickIdeaAgent {
	if client == nil {
		return a
	}

	if timeout <= 0 {
		timeout = defaultQuickIdeaTimeout
	}

	return &QuickIdeaAgent{
		client:          client,
		model:           a.model,
		timeout:         timeout,
		prompts:         a.prompts,
		contextDetector: a.contextDetector,
		contextPrompts:  a.contextPrompts,
		knowledgeBase:   a.knowledgeBase,
		retryConfig:     a.retryConfig,
	}
}

// WithModel returns a copy of the agent configured to use the provided model string.
func (a *QuickIdeaAgent) WithModel(model string) *QuickIdeaAgent {
	if model == "" {
		return a
	}

	return &QuickIdeaAgent{
		client:          a.client,
		model:           model,
		timeout:         a.timeout,
		prompts:         a.prompts,
		contextDetector: a.contextDetector,
		contextPrompts:  a.contextPrompts,
		knowledgeBase:   a.knowledgeBase,
		retryConfig:     a.retryConfig,
	}
}

// WithKnowledgeBase returns a copy of the agent configured to use the provided knowledge base.
func (a *QuickIdeaAgent) WithKnowledgeBase(kb knowledge.EnhancementProvider) *QuickIdeaAgent {
	if kb == nil {
		return a
	}

	return &QuickIdeaAgent{
		client:          a.client,
		model:           a.model,
		timeout:         a.timeout,
		prompts:         a.prompts,
		contextDetector: a.contextDetector,
		contextPrompts:  a.contextPrompts,
		knowledgeBase:   kb,
		retryConfig:     a.retryConfig,
	}
}

// WithRetryConfig returns a copy of the agent configured with the provided retry settings.
func (a *QuickIdeaAgent) WithRetryConfig(cfg RetryConfig) *QuickIdeaAgent {
	return &QuickIdeaAgent{
		client:          a.client,
		model:           a.model,
		timeout:         a.timeout,
		prompts:         a.prompts,
		contextDetector: a.contextDetector,
		contextPrompts:  a.contextPrompts,
		knowledgeBase:   a.knowledgeBase,
		retryConfig:     cfg,
	}
}

// Generate runs the selected quick idea mode and returns structured suggestions.
// A strict timeout is enforced; on timeout or client failure the agent falls back to deterministic suggestions.
func (a *QuickIdeaAgent) Generate(ctx context.Context, req QuickRequest) (*QuickResponse, error) {
	if req.Mode == "" {
		return nil, errors.New("quick idea request missing mode")
	}

	if !isSupportedQuickIdeaMode(req.Mode) {
		return nil, fmt.Errorf("quick idea mode %q is not supported", req.Mode)
	}

	ctx, cancel := context.WithTimeout(ctx, a.timeout)
	defer cancel()

	// Detect content type for context-aware prompts
	contentType := a.contextDetector.AnalyzeContent(req.Context)

	// Enhance with knowledge base if available
	kbEnhanced := a.enhanceWithKnowledgeBase(ctx, req, contentType)

	// Use context-aware prompts if available
	var prompt string
	if a.contextPrompts != nil {
		prompt = a.contextPrompts.RenderPrompt(contentType, req.Mode, kbEnhanced.context, req.Options)
	} else {
		// Fallback to original prompts
		prompt = a.prompts.render(req.Mode, kbEnhanced.context, req.Options)
	}

	start := time.Now()
	output, err := a.invoke(ctx, prompt)
	elapsed := ensurePositiveDuration(time.Since(start))

	if err != nil || len(strings.TrimSpace(output)) == 0 {
		fallback := a.generateContextAwareFallback(req, contentType)
		// Apply knowledge base enhancement to fallback
		a.applyKnowledgeBaseToFallback(fallback, kbEnhanced, req.Mode)
		fallback.ResponseTime = elapsed
		return fallback, nil
	}

	resp := a.prompts.parse(req.Mode, output)
	if resp == nil {
		resp = a.generateContextAwareFallback(req, contentType)
		// Apply knowledge base enhancement to fallback
		a.applyKnowledgeBaseToFallback(resp, kbEnhanced, req.Mode)
	}
	resp.ResponseTime = elapsed
	return resp, nil
}

// invoke delegates prompt generation to the configured client with retry logic.
func (a *QuickIdeaAgent) invoke(ctx context.Context, prompt string) (string, error) {
	if a.client == nil {
		return "", errors.New("no llm client configured")
	}

	return a.invokeWithRetry(ctx, prompt)
}

// invokeWithRetry implements exponential backoff retry logic for AI requests.
// This follows 2026 best practices for resilient LLM client interactions.
func (a *QuickIdeaAgent) invokeWithRetry(ctx context.Context, prompt string) (string, error) {
	opts := map[string]any{
		"temperature": 0.7,
		"top_p":       0.9,
		"num_predict": 120,
	}

	var lastErr error
	maxRetries := a.retryConfig.MaxRetries
	if maxRetries <= 0 {
		maxRetries = 1 // At least one attempt
	}

	for attempt := 0; attempt < maxRetries; attempt++ {
		// Check context before attempt
		if ctx.Err() != nil {
			return "", ctx.Err()
		}

		result, err := a.client.Generate(ctx, a.model, prompt, opts)
		if err == nil {
			return result, nil
		}

		lastErr = err

		// Check if error is retryable
		retryable := true
		if a.retryConfig.RetryableCheck != nil {
			retryable = a.retryConfig.RetryableCheck(err)
		}

		if !retryable {
			return "", err
		}

		// Don't sleep after the last attempt
		if attempt < maxRetries-1 {
			delay := a.calculateBackoff(attempt)
			select {
			case <-ctx.Done():
				return "", ctx.Err()
			case <-time.After(delay):
				// Continue to next attempt
			}
		}
	}

	return "", fmt.Errorf("max retries (%d) exceeded: %w", maxRetries, lastErr)
}

// calculateBackoff computes the exponential backoff delay for a given attempt.
func (a *QuickIdeaAgent) calculateBackoff(attempt int) time.Duration {
	// Exponential backoff: baseDelay * 2^attempt
	delay := a.retryConfig.BaseDelay * time.Duration(1<<uint(attempt))

	// Cap at maxDelay
	if delay > a.retryConfig.MaxDelay {
		delay = a.retryConfig.MaxDelay
	}

	return delay
}

// isRetryableError determines if an error should trigger a retry.
// Returns true for transient errors like timeouts and connection issues.
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	errMsg := strings.ToLower(err.Error())

	// Non-retryable: stub client, model not found, invalid request
	if strings.Contains(errMsg, "stub") ||
		strings.Contains(errMsg, "not found") ||
		strings.Contains(errMsg, "invalid") ||
		strings.Contains(errMsg, "unsupported") {
		return false
	}

	// Retryable: timeouts, connection errors, server errors
	retryablePatterns := []string{
		"timeout",
		"deadline exceeded",
		"connection refused",
		"connection reset",
		"temporary failure",
		"service unavailable",
		"503",
		"502",
		"500",
		"rate limit",
		"too many requests",
		"429",
	}

	for _, pattern := range retryablePatterns {
		if strings.Contains(errMsg, pattern) {
			return true
		}
	}

	// Default: retry on unknown errors (conservative approach)
	return false
}

func (a *QuickIdeaAgent) generateFallback(req QuickRequest) *QuickResponse {
	// Detect content type for better fallback suggestions
	contentType := a.contextDetector.AnalyzeContent(req.Context)
	return a.generateContextAwareFallback(req, contentType)
}

func (a *QuickIdeaAgent) generateContextAwareFallback(req QuickRequest, contentType ContentType) *QuickResponse {
	switch req.Mode {
	case QuickIdeaModeUnstick:
		lines := nonEmptyTail(strings.Split(req.Context, "\n"), 2)
		prefix := ""
		if len(lines) > 0 {
			prefix = lines[len(lines)-1]
		}

		// Generate context-aware suggestions
		switch contentType {
		case ContentTypeLyrics:
			return &QuickResponse{
				Suggestions: []string{
					fmt.Sprintf("%s and the stars begin to fall", strings.TrimSpace(prefix)),
					fmt.Sprintf("%s while the city sleeps below", strings.TrimSpace(prefix)),
					fmt.Sprintf("%s as the morning light breaks through", strings.TrimSpace(prefix)),
				},
			}
		case ContentTypePatterns:
			return &QuickResponse{
				Suggestions: []string{
					fmt.Sprintf("%s | Am - G - C - F", strings.TrimSpace(prefix)),
					fmt.Sprintf("%s | I - V - vi - IV", strings.TrimSpace(prefix)),
					fmt.Sprintf("%s | Dm - G - C - Am", strings.TrimSpace(prefix)),
				},
			}
		case ContentTypeMixed:
			return &QuickResponse{
				Suggestions: []string{
					fmt.Sprintf("%s with a gentle C major progression", strings.TrimSpace(prefix)),
					fmt.Sprintf("%s building to an emotional chorus", strings.TrimSpace(prefix)),
					fmt.Sprintf("%s with a driving rhythm section", strings.TrimSpace(prefix)),
				},
			}
		default:
			return &QuickResponse{
				Suggestions: []string{
					fmt.Sprintf("%s and I keep moving forward", strings.TrimSpace(prefix)),
					fmt.Sprintf("%s while the night hums softly", strings.TrimSpace(prefix)),
					fmt.Sprintf("%s as the skyline flickers", strings.TrimSpace(prefix)),
				},
			}
		}

	case QuickIdeaModeHarmony:
		theme := strings.TrimSpace(req.Options["theme"])
		if theme == "" {
			theme = "all"
		}
		return &QuickResponse{
			Suggestions: []string{
				"C - F - G - C",
				"Am - F - C - G",
				"Dm - G7 - Cmaj7",
			},
			Tip: fmt.Sprintf("Fallback harmony for %s mood", theme),
		}

	case QuickIdeaModeSpark:
		theme := strings.TrimSpace(req.Options["theme"])
		if theme == "" {
			theme = "creativity"
		}

		// Generate context-aware starting ideas
		switch contentType {
		case ContentTypeLyrics:
			return &QuickResponse{
				Suggestions: []string{
					fmt.Sprintf("In the heart of %s, I found my way", theme),
					fmt.Sprintf("%s whispers through the window pane", theme),
					fmt.Sprintf("Chasing %s through the pouring rain", theme),
				},
			}
		case ContentTypePatterns:
			return &QuickResponse{
				Suggestions: []string{
					fmt.Sprintf("%s theme: C - G - Am - F progression", theme),
					fmt.Sprintf("%s rhythm: driving 4/4 with syncopation", theme),
					fmt.Sprintf("%s mood: minor key with descending bassline", theme),
				},
			}
		case ContentTypeMixed:
			return &QuickResponse{
				Suggestions: []string{
					fmt.Sprintf("%s: acoustic guitar with intimate vocals", theme),
					fmt.Sprintf("%s: electronic beat with reflective lyrics", theme),
					fmt.Sprintf("%s: piano ballad with swelling strings", theme),
				},
			}
		default:
			return &QuickResponse{
				Suggestions: []string{
					fmt.Sprintf("I woke up chasing %s shadows", theme),
					fmt.Sprintf("Streetlights paint %s horizons", theme),
					fmt.Sprintf("You left fingerprints on %s dawn", theme),
				},
			}
		}

	case QuickIdeaModeTweak:
		base := strings.TrimSpace(req.Context)
		if base == "" {
			base = "Hold on to the quiet light"
		}

		// Generate context-aware variations
		switch contentType {
		case ContentTypeLyrics:
			return &QuickResponse{
				Suggestions: []string{
					base,
					"Rewrite with stronger imagery and emotion",
					"Replace clichÃ©s with fresh, specific details",
				},
			}
		case ContentTypePatterns:
			return &QuickResponse{
				Suggestions: []string{
					base,
					"Add sophisticated voice leading",
					"Incorporate rhythmic variation",
				},
			}
		case ContentTypeMixed:
			return &QuickResponse{
				Suggestions: []string{
					base,
					"Enhance both lyrical and musical flow",
					"Strengthen the connection between words and music",
				},
			}
		default:
			return &QuickResponse{
				Suggestions: []string{
					base,
					"Rewrite the lines your heartbeat drew",
					"Trade every clichÃ© for neon sparks",
				},
			}
		}

	case QuickIdeaModeCheck:
		// Generate context-aware quality feedback
		switch contentType {
		case ContentTypeLyrics:
			return &QuickResponse{
				Rating: "OKAY",
				Tip:    "Add vivid sensory details",
			}
		case ContentTypePatterns:
			return &QuickResponse{
				Rating: "OKAY",
				Tip:    "Strengthen harmonic resolution",
			}
		case ContentTypeMixed:
			return &QuickResponse{
				Rating: "OKAY",
				Tip:    "Balance lyrics and music",
			}
		default:
			return &QuickResponse{
				Rating: "OKAY",
				Tip:    "Add vivid sensory image",
			}
		}

	default:
		return &QuickResponse{
			Suggestions: []string{"No suggestion available."},
		}
	}
}

// nonEmptyTail returns up to n non-empty lines from the tail of the slice.
func nonEmptyTail(lines []string, n int) []string {
	if n <= 0 {
		return nil
	}

	result := make([]string, 0, n)
	for i := len(lines) - 1; i >= 0 && len(result) < n; i-- {
		line := strings.TrimSpace(lines[i])
		if line != "" {
			result = append(result, line)
		}
	}

	for i, j := 0, len(result)-1; i < j; i, j = i+1, j-1 {
		result[i], result[j] = result[j], result[i]
	}

	return result
}

type stubQuickClient struct{}

// enhanceWithKnowledgeBase enhances the request context with knowledge base insights
func (a *QuickIdeaAgent) enhanceWithKnowledgeBase(ctx context.Context, req QuickRequest, contentType ContentType) knowledgeEnhancedContext {
	enhanced := knowledgeEnhancedContext{
		context: req.Context,
		cards:   []knowledge.Card{},
		tip:     "",
	}

	if a.knowledgeBase == nil {
		return enhanced
	}

	// Search for relevant knowledge cards based on mode and content
	options := knowledge.SearchOptions{
		Limit:        3,
		MinRelevance: 0.6,
		UseCache:     true,
	}

	switch req.Mode {
	case QuickIdeaModeUnstick:
		// For unstick, look for continuation techniques
		if contentType == ContentTypeLyrics {
			options.Categories = []string{"lyrical-techniques", "inspiration"}
		} else if contentType == ContentTypePatterns {
			options.Categories = []string{"chord-progressions", "song-structure"}
		}

	case QuickIdeaModeSpark:
		// For spark, look for inspiration
		options.Categories = []string{"inspiration"}
		// Use theme from options if available
		if theme, ok := req.Options["theme"]; ok && theme != "" {
			enhanced.context = theme
		}

	case QuickIdeaModeTweak:
		// For tweak, look for improvement techniques
		if contentType == ContentTypeLyrics {
			options.Categories = []string{"lyrical-techniques"}
		} else if contentType == ContentTypePatterns {
			options.Categories = []string{"chord-progressions"}
		}

	case QuickIdeaModeCheck:
		// For check, look for quality guidelines
		options.Categories = []string{"lyrical-techniques", "song-structure"}
	}

	// Perform search
	result, err := a.knowledgeBase.GetInspirationCards(ctx, enhanced.context, options)
	if err == nil && len(result.Cards) > 0 {
		enhanced.cards = result.Cards

		// Create a knowledge-enhanced prompt context
		if len(result.Cards) > 0 {
			enhanced.context = fmt.Sprintf("%s\n\nKnowledge: %s", req.Context, result.Cards[0].Content)
			enhanced.tip = result.Cards[0].Title
		}
	}

	return enhanced
}

// applyKnowledgeBaseToFallback applies knowledge base insights to fallback responses
func (a *QuickIdeaAgent) applyKnowledgeBaseToFallback(resp *QuickResponse, enhanced knowledgeEnhancedContext, mode QuickIdeaMode) {
	if len(enhanced.cards) == 0 {
		return
	}

	// Enhance suggestions with knowledge base insights
	for i, suggestion := range resp.Suggestions {
		if i < len(enhanced.cards) {
			card := enhanced.cards[i]

			// Add knowledge-based variation to suggestion
			switch mode {
			case QuickIdeaModeUnstick:
				if card.Category == "lyrical-techniques" {
					resp.Suggestions[i] = fmt.Sprintf("%s (apply %s)", suggestion, card.Title)
				} else if card.Category == "chord-progressions" {
					if example, ok := card.Metadata["example_c"]; ok {
						resp.Suggestions[i] = fmt.Sprintf("%s | %s", suggestion, example)
					}
				}

			case QuickIdeaModeSpark:
				if card.Category == "inspiration" {
					resp.Suggestions[i] = fmt.Sprintf("%s - inspired by %s", suggestion, card.Title)
				}

			case QuickIdeaModeTweak:
				if card.Category == "lyrical-techniques" {
					resp.Suggestions[i] = fmt.Sprintf("%s (using %s technique)", suggestion, card.Title)
				}

			case QuickIdeaModeCheck:
				if resp.Tip == "" && enhanced.tip != "" {
					resp.Tip = fmt.Sprintf("KB Tip: %s", enhanced.tip)
				}
			}
		}
	}
}

// knowledgeEnhancedContext represents context enhanced with knowledge base information
type knowledgeEnhancedContext struct {
	context string
	cards   []knowledge.Card
	tip     string
}

func isSupportedQuickIdeaMode(mode QuickIdeaMode) bool {
	_, ok := supportedQuickIdeaModes[mode]
	return ok
}

func ensurePositiveDuration(d time.Duration) time.Duration {
	if d <= 0 {
		return time.Nanosecond
	}
	return d
}

// GetKnowledgeBaseStatus returns the current status of the knowledge base
func (a *QuickIdeaAgent) GetKnowledgeBaseStatus(ctx context.Context) *knowledge.KnowledgeStatus {
	if a.knowledgeBase == nil {
		return &knowledge.KnowledgeStatus{
			Available: false,
			Error:     "No knowledge base configured",
		}
	}

	kb := a.knowledgeBase.GetKnowledgeBase()
	return kb.GetStatus(ctx)
}

// IsKnowledgeBaseAvailable returns whether the knowledge base is available
func (a *QuickIdeaAgent) IsKnowledgeBaseAvailable(ctx context.Context) bool {
	if a.knowledgeBase == nil {
		return false
	}

	kb := a.knowledgeBase.GetKnowledgeBase()
	return kb.IsAvailable(ctx)
}

func (s *stubQuickClient) Generate(_ context.Context, _, _ string, _ map[string]any) (string, error) {
	return "", errors.New("stub quick idea client has no external backend")
}
