package app

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"
)

// ChatMessage represents a single message in a chat conversation
type ChatMessage struct {
	Role      string    `json:"role"`      // "user" or "assistant"
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

// ChatSession manages a conversation with the AI
type ChatSession struct {
	messages []ChatMessage
	model    string
	ollamaURL string
}

// NewChatSession creates a new AI chat session
func NewChatSession(model, ollamaURL string) *ChatSession {
	return &ChatSession{
		messages:  make([]ChatMessage, 0),
		model:     model,
		ollamaURL: ollamaURL,
	}
}

// Chat sends a message and streams the response
func (s *AIService) Chat(ctx context.Context, message string) (<-chan string, error) {
	// Create response channel
	responseChan := make(chan string, 10)

	// Check if Ollama is available
	if !s.IsAvailable() {
		go func() {
			defer close(responseChan)
			responseChan <- "AI assistant is not available. Please ensure Ollama is running."
		}()
		return responseChan, nil
	}

	// Start streaming response in goroutine
	go func() {
		defer close(responseChan)

		// Prepare request
		reqBody := map[string]interface{}{
			"model":  s.model,
			"prompt": message,
			"stream": true,
			"options": map[string]interface{}{
				"temperature": 0.7,
				"top_p":       0.9,
			},
		}

		reqJSON, err := json.Marshal(reqBody)
		if err != nil {
			responseChan <- fmt.Sprintf("Error: %v", err)
			return
		}

		// Create HTTP request
		req, err := http.NewRequestWithContext(ctx, "POST", s.ollamaURL+"/api/generate", strings.NewReader(string(reqJSON)))
		if err != nil {
			responseChan <- fmt.Sprintf("Error: %v", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")

		// Send request
		client := &http.Client{Timeout: 60 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			responseChan <- fmt.Sprintf("Error connecting to AI: %v", err)
			return
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			responseChan <- fmt.Sprintf("AI returned error: %s", resp.Status)
			return
		}

		// Stream response
		decoder := json.NewDecoder(resp.Body)
		for {
			var streamResp struct {
				Response string `json:"response"`
				Done     bool   `json:"done"`
			}

			if err := decoder.Decode(&streamResp); err != nil {
				if err.Error() != "EOF" {
					responseChan <- fmt.Sprintf("Error reading response: %v", err)
				}
				return
			}

			if streamResp.Response != "" {
				responseChan <- streamResp.Response
			}

			if streamResp.Done {
				return
			}

			// Check context cancellation
			select {
			case <-ctx.Done():
				return
			default:
			}
		}
	}()

	return responseChan, nil
}

// ChatWithHistory sends a message with conversation history
func (session *ChatSession) Send(ctx context.Context, message string, ollamaURL string) (<-chan string, error) {
	// Add user message to history
	session.messages = append(session.messages, ChatMessage{
		Role:      "user",
		Content:   message,
		Timestamp: time.Now(),
	})

	// Build context from history
	contextPrompt := session.buildContextPrompt()

	// Create response channel
	responseChan := make(chan string, 10)

	// Start streaming response
	go func() {
		defer close(responseChan)

		reqBody := map[string]interface{}{
			"model":  session.model,
			"prompt": contextPrompt,
			"stream": true,
			"options": map[string]interface{}{
				"temperature": 0.7,
				"top_p":       0.9,
			},
		}

		reqJSON, _ := json.Marshal(reqBody)
		req, err := http.NewRequestWithContext(ctx, "POST", ollamaURL+"/api/generate", strings.NewReader(string(reqJSON)))
		if err != nil {
			responseChan <- fmt.Sprintf("Error: %v", err)
			return
		}

		req.Header.Set("Content-Type", "application/json")

		client := &http.Client{Timeout: 60 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			responseChan <- fmt.Sprintf("Error: %v", err)
			return
		}
		defer resp.Body.Close()

		var fullResponse strings.Builder
		decoder := json.NewDecoder(resp.Body)

		for {
			var streamResp struct {
				Response string `json:"response"`
				Done     bool   `json:"done"`
			}

			if err := decoder.Decode(&streamResp); err != nil {
				break
			}

			if streamResp.Response != "" {
				fullResponse.WriteString(streamResp.Response)
				responseChan <- streamResp.Response
			}

			if streamResp.Done {
				// Add assistant response to history
				session.messages = append(session.messages, ChatMessage{
					Role:      "assistant",
					Content:   fullResponse.String(),
					Timestamp: time.Now(),
				})
				break
			}
		}
	}()

	return responseChan, nil
}

// buildContextPrompt builds a prompt with conversation history
func (session *ChatSession) buildContextPrompt() string {
	var prompt strings.Builder

	prompt.WriteString("You are a helpful AI assistant for songwriting. Here is our conversation:\n\n")

	// Include last 5 messages for context
	start := 0
	if len(session.messages) > 10 {
		start = len(session.messages) - 10
	}

	for i := start; i < len(session.messages); i++ {
		msg := session.messages[i]
		if msg.Role == "user" {
			prompt.WriteString("User: ")
		} else {
			prompt.WriteString("Assistant: ")
		}
		prompt.WriteString(msg.Content)
		prompt.WriteString("\n\n")
	}

	prompt.WriteString("Assistant: ")
	return prompt.String()
}

// ClearHistory clears the conversation history
func (session *ChatSession) ClearHistory() {
	session.messages = make([]ChatMessage, 0)
}

// GetHistory returns the conversation history
func (session *ChatSession) GetHistory() []ChatMessage {
	return session.messages
}
