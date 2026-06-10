package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client handles communication with Ollama API
type Client struct {
	config     Config
	httpClient *http.Client
}

// NewClient creates a new AI client
func NewClient(config Config) *Client {
	return &Client{
		config: config,
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
	}
}

// OllamaRequest represents a request to Ollama API
type OllamaRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
	Options map[string]interface{} `json:"options,omitempty"`
}

// OllamaResponse represents a response from Ollama API
type OllamaResponse struct {
	Model     string `json:"model"`
	CreatedAt string `json:"created_at"`
	Response  string `json:"response"`
	Done      bool   `json:"done"`
}

// IsEnabled returns whether the AI assistant is enabled
func (c *Client) IsEnabled() bool {
	return c.config.Enabled
}

// GetSuggestion gets an AI suggestion based on the context
func (c *Client) GetSuggestion(ctx context.Context, suggestionType SuggestionType, content string, context string) (*Suggestion, error) {
	if !c.config.Enabled {
		return nil, fmt.Errorf("AI assistant is disabled")
	}

	prompt := c.buildPrompt(suggestionType, content, context)

	response, err := c.generate(ctx, prompt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate suggestion: %w", err)
	}

	return &Suggestion{
		Type:      suggestionType,
		Content:   response,
		Timestamp: time.Now(),
	}, nil
}

// buildPrompt constructs the prompt based on suggestion type
func (c *Client) buildPrompt(suggestionType SuggestionType, content string, context string) string {
	var systemPrompt string

	switch suggestionType {
	case SuggestionContinue:
		systemPrompt = "You are a creative writing assistant. Continue the following story naturally and engagingly. Match the writing style and tone. Write 2-3 paragraphs."
	case SuggestionImprove:
		systemPrompt = "You are an editing assistant. Improve the following text while maintaining the author's voice. Make it more vivid and engaging."
	case SuggestionDialogue:
		systemPrompt = "You are a dialogue specialist. Suggest realistic dialogue that fits the scene and characters. Make it natural and character-driven."
	case SuggestionDescription:
		systemPrompt = "You are a descriptive writing expert. Enhance this scene with vivid sensory details. Show, don't tell."
	case SuggestionCharacter:
		systemPrompt = "You are a character development expert. Suggest character behaviors, motivations, or reactions that would make this scene more compelling."
	default:
		systemPrompt = "You are a creative writing assistant. Help improve this text."
	}

	prompt := fmt.Sprintf("%s\n\n", systemPrompt)

	if context != "" {
		prompt += fmt.Sprintf("Context:\n%s\n\n", context)
	}

	prompt += fmt.Sprintf("Text:\n%s\n\nSuggestion:", content)

	return prompt
}

// generate sends a request to Ollama and returns the response
func (c *Client) generate(ctx context.Context, prompt string) (string, error) {
	reqBody := OllamaRequest{
		Model:  c.config.Model,
		Prompt: prompt,
		Stream: false,
		Options: map[string]interface{}{
			"temperature": c.config.Temperature,
			"num_predict": c.config.MaxTokens,
		},
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("failed to marshal request: %w", err)
	}

	url := fmt.Sprintf("%s/api/generate", c.config.APIEndpoint)
	req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to read response: %w", err)
	}

	var ollamaResp OllamaResponse
	if err := json.Unmarshal(body, &ollamaResp); err != nil {
		return "", fmt.Errorf("failed to parse response: %w", err)
	}

	return ollamaResp.Response, nil
}

// CheckConnection verifies that Ollama is running and accessible
func (c *Client) CheckConnection(ctx context.Context) error {
	url := fmt.Sprintf("%s/api/tags", c.config.APIEndpoint)
	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("failed to connect to Ollama: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("Ollama API returned status %d", resp.StatusCode)
	}

	return nil
}
