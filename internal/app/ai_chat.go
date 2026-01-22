package app

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// ChatMessage represents a single message in a chat conversation
type ChatMessage struct {
	Role      string    `json:"role"` // "user" or "assistant"
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// ChatSession manages a conversation with the AI
type ChatSession struct {
	messages []ChatMessage
	provider string
	client   any // ollama.Client or glm.Client
	model    string
}

// NewChatSession creates a new AI chat session
func (s *AIService) NewChatSession() *ChatSession {
	var client any
	var model string

	switch s.config.AI.Provider {
	case "glm":
		client = s.glmClient
		model = s.config.GLM.Model
	default:
		client = s.ollamaClient
		model = s.config.AI.Model
	}

	return &ChatSession{
		messages: make([]ChatMessage, 0),
		provider: s.config.AI.Provider,
		client:   client,
		model:    model,
	}
}

// Chat sends a message and streams the response
func (s *AIService) StreamChat(ctx context.Context, message string) (<-chan string, error) {
	// Create response channel
	responseChan := make(chan string, 10)

	// Determine provider
	provider := s.config.AI.Provider
	if provider == "hybrid" {
		provider = "ollama" // Default hybrid chat to local for speed/privacy
	}

	// Start streaming response in goroutine
	go func() {
		defer close(responseChan)

		if provider == "glm" {
			s.streamGLM(ctx, message, responseChan)
		} else {
			s.streamOllama(ctx, message, responseChan)
		}
	}()

	return responseChan, nil
}

func (s *AIService) streamOllama(ctx context.Context, message string, ch chan<- string) {
	reqBody := map[string]interface{}{
		"model":  s.config.AI.Model,
		"prompt": message,
		"stream": true,
	}

	reqJSON, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, "POST", s.config.AI.BaseURL+"/api/generate", bytes.NewReader(reqJSON))
	if err != nil {
		ch <- fmt.Sprintf("Error creating request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		ch <- fmt.Sprintf("Error connecting to Ollama: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ch <- fmt.Sprintf("Ollama returned error: %s", resp.Status)
		return
	}

	decoder := json.NewDecoder(resp.Body)
	for {
		var streamResp struct {
			Response string `json:"response"`
			Done     bool   `json:"done"`
		}

		if err := decoder.Decode(&streamResp); err != nil {
			return
		}

		if streamResp.Response != "" {
			ch <- streamResp.Response
		}

		if streamResp.Done {
			return
		}

		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

func (s *AIService) streamGLM(ctx context.Context, message string, ch chan<- string) {
	reqBody := map[string]interface{}{
		"model": s.config.GLM.Model,
		"messages": []map[string]string{
			{"role": "user", "content": message},
		},
		"stream": true,
	}

	reqJSON, _ := json.Marshal(reqBody)
	req, err := http.NewRequestWithContext(ctx, "POST", "https://open.bigmodel.cn/api/paas/v4/chat/completions", bytes.NewReader(reqJSON))
	if err != nil {
		ch <- fmt.Sprintf("Error creating request: %v", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+s.config.GLM.APIKey)

	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		ch <- fmt.Sprintf("Error connecting to GLM: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		ch <- fmt.Sprintf("GLM returned error: %s", resp.Status)
		return
	}

	reader := bufio.NewReader(resp.Body)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return
		}

		line = strings.TrimSpace(line)
		if line == "" || !strings.HasPrefix(line, "data: ") {
			continue
		}

		data := strings.TrimPrefix(line, "data: ")
		if data == "[DONE]" {
			return
		}

		var streamResp struct {
			Choices []struct {
				Delta struct {
					Content string `json:"content"`
				} `json:"delta"`
			} `json:"choices"`
		}

		if err := json.Unmarshal([]byte(data), &streamResp); err != nil {
			continue
		}

		if len(streamResp.Choices) > 0 {
			content := streamResp.Choices[0].Delta.Content
			if content != "" {
				ch <- content
			}
		}

		select {
		case <-ctx.Done():
			return
		default:
		}
	}
}

// buildContextPrompt builds a prompt with conversation history (for non-streaming or simple use)
func (session *ChatSession) buildContextPrompt() string {
	var prompt strings.Builder
	prompt.WriteString("You are a helpful AI assistant for songwriting. History:\n\n")

	start := 0
	if len(session.messages) > 10 {
		start = len(session.messages) - 10
	}

	for i := start; i < len(session.messages); i++ {
		msg := session.messages[i]
		prompt.WriteString(fmt.Sprintf("%s: %s\n\n", strings.Title(msg.Role), msg.Content))
	}

	prompt.WriteString("Assistant: ")
	return prompt.String()
}

// Send (Draft) - Simplified for now as we focus on streaming
func (session *ChatSession) Send(ctx context.Context, message string) (<-chan string, error) {
	// This would mirror StreamChat but maintain internal state
	// Placeholder for now
	ch := make(chan string, 1)
	ch <- "Chat history support is coming soon in the unified AI update."
	close(ch)
	return ch, nil
}
