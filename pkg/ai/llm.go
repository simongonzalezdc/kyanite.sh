package ai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"sync"
	"time"
)

const (
	maxRetries       = 2
	retryBaseDelay   = 500 * time.Millisecond
	healthCacheTTL   = 10 * time.Second
)

// LLMClient handles communication with an Ollama instance.
type LLMClient struct {
	baseURL    string
	model      string
	httpClient *http.Client

	// Health check cache
	healthMu    sync.RWMutex
	healthLast  time.Time
	healthResult bool
}

// NewLLMClient creates a new LLM client for an Ollama server.
func NewLLMClient(baseURL, model string, timeout time.Duration) *LLMClient {
	return &LLMClient{
		baseURL: baseURL,
		model:   model,
		httpClient: &http.Client{
			Timeout: timeout,
		},
	}
}

// ollamaRequest is the request body for Ollama's /api/chat endpoint.
type ollamaRequest struct {
	Model    string          `json:"model"`
	Messages []ollamaMessage `json:"messages"`
	Stream   bool            `json:"stream"`
	Format   string          `json:"format,omitempty"`
	Think    bool            `json:"think"`
	Options  ollamaOptions   `json:"options,omitempty"`
}

type ollamaMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ollamaOptions struct {
	Temperature float64  `json:"temperature,omitempty"`
	NumPredict  int      `json:"num_predict,omitempty"`
	Stop        []string `json:"stop,omitempty"`
}

// ollamaResponse is the response from Ollama's /api/chat endpoint.
type ollamaResponse struct {
	Message ollamaMessage `json:"message"`
	Done    bool          `json:"done"`
	Error   string        `json:"error,omitempty"`
}

// StreamChunk represents a single chunk in a streaming response.
type StreamChunk struct {
	Content string
	Done    bool
	Error   error
}

// Generate sends a prompt to the LLM and returns the response.
// It retries transient failures (network errors, 5xx) up to maxRetries times
// with exponential backoff.
func (c *LLMClient) Generate(ctx context.Context, prompt string, opts ...GenerateOption) (string, error) {
	cfg := &GenerateOptions{
		Temperature: 0.7,
		MaxTokens:   1024,
	}
	for _, o := range opts {
		o(cfg)
	}

	// Apply per-operation timeout override
	if cfg.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cfg.Timeout)
		defer cancel()
	}

	messages := buildMessages(prompt, cfg)

	reqBody := ollamaRequest{
		Model:    c.model,
		Messages: messages,
		Stream:   false,
		Think:    false,
		Options: ollamaOptions{
			Temperature: cfg.Temperature,
			NumPredict:  cfg.MaxTokens,
			Stop:        cfg.StopWords,
		},
	}
	if cfg.JSONMode {
		reqBody.Format = "json"
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return "", fmt.Errorf("%w: %v", ErrContextCanceled, ctx.Err())
			case <-time.After(retryBaseDelay * time.Duration(1<<(attempt-1))):
			}
		}

		resp, respBody, err := c.doRequest(ctx, body)
		if err != nil {
			if isTransientError(err) {
				lastErr = err
				continue
			}
			return "", fmt.Errorf("%w: %v", ErrLLMUnavailable, err)
		}
		resp.Body.Close() // body already read by doRequest, close immediately

		if resp.StatusCode == http.StatusTooManyRequests {
			lastErr = ErrRateLimited
			continue
		}

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("%w (status %d): %s", ErrLLMRequestFailed, resp.StatusCode, string(respBody))
		}

		var ollamaResp ollamaResponse
		if err := json.Unmarshal(respBody, &ollamaResp); err != nil {
			return "", fmt.Errorf("decode response: %w", err)
		}

		if ollamaResp.Error != "" {
			return "", fmt.Errorf("%w: %s", ErrLLMRequestFailed, ollamaResp.Error)
		}

		return ollamaResp.Message.Content, nil
	}

	return "", fmt.Errorf("%w after %d retries: %v", ErrLLMRequestFailed, maxRetries, lastErr)
}

// Chat sends a conversation to the LLM and returns the assistant's response.
func (c *LLMClient) Chat(ctx context.Context, messages []Message, opts ...GenerateOption) (string, error) {
	cfg := &GenerateOptions{
		Temperature: 0.7,
		MaxTokens:   1024,
	}
	for _, o := range opts {
		o(cfg)
	}

	// Apply per-operation timeout override
	if cfg.Timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, cfg.Timeout)
		defer cancel()
	}

	ollamaMsgs := make([]ollamaMessage, 0, len(messages)+1)
	if cfg.SystemPrompt != "" {
		ollamaMsgs = append(ollamaMsgs, ollamaMessage{Role: "system", Content: cfg.SystemPrompt})
	}
	for _, m := range messages {
		ollamaMsgs = append(ollamaMsgs, ollamaMessage{Role: m.Role, Content: m.Content})
	}

	reqBody := ollamaRequest{
		Model:    c.model,
		Messages: ollamaMsgs,
		Stream:   false,
		Think:    false,
		Options: ollamaOptions{
			Temperature: cfg.Temperature,
			NumPredict:  cfg.MaxTokens,
			Stop:        cfg.StopWords,
		},
	}
	if cfg.JSONMode {
		reqBody.Format = "json"
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return "", fmt.Errorf("marshal request: %w", err)
	}

	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		if attempt > 0 {
			select {
			case <-ctx.Done():
				return "", fmt.Errorf("%w: %v", ErrContextCanceled, ctx.Err())
			case <-time.After(retryBaseDelay * time.Duration(1<<(attempt-1))):
			}
		}

		resp, respBody, err := c.doRequest(ctx, body)
		if err != nil {
			if isTransientError(err) {
				lastErr = err
				continue
			}
			return "", fmt.Errorf("%w: %v", ErrLLMUnavailable, err)
		}
		resp.Body.Close() // body already read by doRequest, close immediately

		if resp.StatusCode == http.StatusTooManyRequests {
			lastErr = ErrRateLimited
			continue
		}

		if resp.StatusCode != http.StatusOK {
			return "", fmt.Errorf("%w (status %d): %s", ErrLLMRequestFailed, resp.StatusCode, string(respBody))
		}

		var ollamaResp ollamaResponse
		if err := json.Unmarshal(respBody, &ollamaResp); err != nil {
			return "", fmt.Errorf("decode response: %w", err)
		}

		if ollamaResp.Error != "" {
			return "", fmt.Errorf("%w: %s", ErrLLMRequestFailed, ollamaResp.Error)
		}

		return ollamaResp.Message.Content, nil
	}

	return "", fmt.Errorf("%w after %d retries: %v", ErrLLMRequestFailed, maxRetries, lastErr)
}

// StreamChat sends a conversation to the LLM and streams response chunks.
// The returned channel receives StreamChunk messages until Done is true or Error is set.
func (c *LLMClient) StreamChat(ctx context.Context, messages []Message, opts ...GenerateOption) (<-chan StreamChunk, error) {
	cfg := &GenerateOptions{
		Temperature: 0.7,
		MaxTokens:   1024,
	}
	for _, o := range opts {
		o(cfg)
	}

	ollamaMsgs := make([]ollamaMessage, 0, len(messages)+1)
	if cfg.SystemPrompt != "" {
		ollamaMsgs = append(ollamaMsgs, ollamaMessage{Role: "system", Content: cfg.SystemPrompt})
	}
	for _, m := range messages {
		ollamaMsgs = append(ollamaMsgs, ollamaMessage{Role: m.Role, Content: m.Content})
	}

	reqBody := ollamaRequest{
		Model:    c.model,
		Messages: ollamaMsgs,
		Stream:   true,
		Think:    false,
		Options: ollamaOptions{
			Temperature: cfg.Temperature,
			NumPredict:  cfg.MaxTokens,
			Stop:        cfg.StopWords,
		},
	}

	body, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/chat", bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrLLMUnavailable, err)
	}

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		resp.Body.Close()
		return nil, fmt.Errorf("%w (status %d): %s", ErrLLMRequestFailed, resp.StatusCode, string(respBody))
	}

	ch := make(chan StreamChunk, 64)
	go func() {
		defer close(ch)
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Bytes()
			if len(line) == 0 {
				continue
			}

			var chunk ollamaResponse
			if err := json.Unmarshal(line, &chunk); err != nil {
				ch <- StreamChunk{Error: fmt.Errorf("decode stream chunk: %w", err)}
				return
			}

			if chunk.Error != "" {
				ch <- StreamChunk{Error: fmt.Errorf("%w: %s", ErrLLMRequestFailed, chunk.Error)}
				return
			}

			ch <- StreamChunk{
				Content: chunk.Message.Content,
				Done:    chunk.Done,
			}

			if chunk.Done {
				return
			}
		}

		if err := scanner.Err(); err != nil {
			ch <- StreamChunk{Error: fmt.Errorf("stream read: %w", err)}
		}
	}()

	return ch, nil
}

// IsAvailable checks if the Ollama server is reachable.
// Results are cached for healthCacheTTL to avoid hammering the server.
func (c *LLMClient) IsAvailable(ctx context.Context) bool {
	c.healthMu.RLock()
	if time.Since(c.healthLast) < healthCacheTTL {
		result := c.healthResult
		c.healthMu.RUnlock()
		return result
	}
	c.healthMu.RUnlock()

	ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/api/tags", nil)
	if err != nil {
		return false
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		c.setHealth(false)
		return false
	}
	defer resp.Body.Close()

	result := resp.StatusCode == http.StatusOK
	c.setHealth(result)
	return result
}

func (c *LLMClient) setHealth(result bool) {
	c.healthMu.Lock()
	c.healthResult = result
	c.healthLast = time.Now()
	c.healthMu.Unlock()
}

// doRequest executes an HTTP POST with the given body and returns the response and body.
func (c *LLMClient) doRequest(ctx context.Context, body []byte) (*http.Response, []byte, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL+"/api/chat", bytes.NewReader(body))
	if err != nil {
		return nil, nil, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, nil, err
	}

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		resp.Body.Close()
		return nil, nil, err
	}

	return resp, respBody, nil
}

// isTransientError returns true for network errors worth retrying.
func isTransientError(err error) bool {
	if errors.Is(err, context.DeadlineExceeded) {
		return false // Don't retry if caller's deadline exceeded
	}
	if errors.Is(err, context.Canceled) {
		return false
	}
	// Retry on temporary network errors (DNS, TCP timeout, connection refused)
	var netErr net.Error
	if errors.As(err, &netErr) {
		return true
	}
	return false
}

func buildMessages(prompt string, cfg *GenerateOptions) []ollamaMessage {
	msgs := make([]ollamaMessage, 0, 2)
	if cfg.SystemPrompt != "" {
		msgs = append(msgs, ollamaMessage{Role: "system", Content: cfg.SystemPrompt})
	}
	msgs = append(msgs, ollamaMessage{Role: "user", Content: prompt})
	return msgs
}
