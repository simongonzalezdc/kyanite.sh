// Package glm provides a client for the GLM AI service.
package glm

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// Client handles communication with the Zhipu AI GLM API
type Client struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new GLM client
func NewClient(apiKey string, timeout time.Duration) *Client {
	return &Client{
		APIKey:  apiKey,
		BaseURL: "https://open.bigmodel.cn/api/paas/v4", // Standard GLM V4 endpoint
		HTTPClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// ChatRequest represents a request to the chat/completions endpoint
type ChatRequest struct {
	Model       string    `json:"model"`
	Messages    []Message `json:"messages"`
	Temperature float64   `json:"temperature,omitempty"`
	MaxTokens   int       `json:"max_tokens,omitempty"`
	Stream      bool      `json:"stream"`
}

// Message represents a message in the chat conversation
type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// ChatResponse represents a response from the chat/completions endpoint
type ChatResponse struct {
	Choices []struct {
		Message Message `json:"message"`
	} `json:"choices"`
	Usage struct {
		TotalTokens int `json:"total_tokens"`
	} `json:"usage"`
}

// Generate sends a prompt to GLM-4.7 and returns the completion
func (c *Client) Generate(ctx context.Context, model, prompt string, options map[string]any) (string, error) {
	if c.APIKey == "" {
		return "", fmt.Errorf("GLM API key is required")
	}

	// Default model for songwriting brainstorm if not specified
	if model == "" {
		model = "glm-4.7-plus"
	}

	reqBody := ChatRequest{
		Model: model,
		Messages: []Message{
			{Role: "user", Content: prompt},
		},
		Stream: false,
	}

	// Map generic options to GLM specific ones
	if t, ok := options["temperature"].(float64); ok {
		reqBody.Temperature = t
	}
	if m, ok := options["max_tokens"].(int); ok {
		reqBody.MaxTokens = m
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL+"/chat/completions", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.APIKey)

	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", fmt.Errorf("GLM error (status %d): failed to read body: %w", resp.StatusCode, err)
		}
		return "", fmt.Errorf("GLM error (status %d): %s", resp.StatusCode, string(body))
	}

	var chatResp ChatResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatResp); err != nil {
		return "", fmt.Errorf("decode response: %w", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("empty response from GLM")
	}

	return chatResp.Choices[0].Message.Content, nil
}
