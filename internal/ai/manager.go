package ai

import (
	"bufio"
	"bytes"
	"context"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

// Manager handles AI interactions with local and remote models
type Manager struct {
	ollamaURL       string
	openRouterKey   string
	model           string
	availableModels []string

	// Shared HTTP client for connection pooling
	httpClient *http.Client

	// Cache for AI responses with LRU eviction
	cache         map[string]cacheEntry
	cacheHits     map[string]time.Time // Track last access time for LRU
	cacheMaxSize  int
	cacheMutex    sync.RWMutex
	cachePath     string
	cacheDirty    bool // Track if cache needs saving
	lastCacheSave time.Time

	// Model management
	modelAvailable map[string]bool
	modelMutex     sync.RWMutex
}

// cacheEntry represents a cached AI response
type cacheEntry struct {
	Response  any       `json:"response"`
	Timestamp time.Time `json:"timestamp"`
	ExpiresAt time.Time `json:"expires_at"`
	AccessCnt int       `json:"access_count"` // Track access count for LRU
}

// ParsedTask represents the structured output from AI
type ParsedTask struct {
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline,omitempty"`
	Priority    string    `json:"priority"`
	Categories  []string  `json:"categories,omitempty"`
}

// New creates a new AI manager
func New() *Manager {
	manager := &Manager{
		ollamaURL:     "http://localhost:11434/api/generate",
		openRouterKey: os.Getenv("OPENROUTER_API_KEY"),
		model:         "qwen2.5:1.5b", // Use Qwen 2.5 1.5B for everything
		availableModels: []string{
			"qwen2.5:1.5b", // Primary - fast and efficient
		},
		// Create shared HTTP client with connection pooling
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		},
		cache:          make(map[string]cacheEntry),
		cacheHits:      make(map[string]time.Time),
		cacheMaxSize:   500, // Limit cache to 500 entries
		modelAvailable: make(map[string]bool),
	}

	// Set cache path
	home, err := os.UserHomeDir()
	if err != nil {
		manager.cachePath = "./ai_cache.json"
	} else {
		manager.cachePath = filepath.Join(home, ".focus", "ai_cache.json")
	}

	// Load cache if it exists
	manager.loadCache()

	// Start background cache saver (saves every 5 minutes)
	go manager.periodicCacheSaver()

	return manager
}

// ParseTask converts natural language to structured task using AI
func (m *Manager) ParseTask(ctx context.Context, input string) (*ParsedTask, error) {
	// Check cache first
	cacheKey := m.generateCacheKey("parse", input)
	if cached, found := m.getFromCache(cacheKey); found {
		// Type assertion with safety check to prevent runtime panic
		if task, ok := cached.(*ParsedTask); ok {
			return task, nil
		}
		// Cache corruption detected - remove invalid entry and continue processing
		// (removeFromCache method will be added later with logging framework)
		// For now, just continue with normal processing
	}

	// Use only qwen2.5:1.5b model
	m.model = "qwen2.5:1.5b"
	if result, err := m.parseWithOllama(ctx, input); err == nil {
		if validated, ok := m.validateResponse(result); ok {
			// Cache the result
			m.saveToCache(cacheKey, validated)
			return validated, nil
		}
	}

	// Fallback to OpenRouter with user approval
	if m.openRouterKey != "" {
		if m.requestRemoteApproval("task parsing") {
			if result, err := m.parseWithOpenRouter(ctx, input); err == nil {
				if validated, ok := m.validateResponse(result); ok {
					// Cache the result
					m.saveToCache(cacheKey, validated)
					return validated, nil
				}
			}
		}
	}

	// If both fail, return basic parsing
	basicResult := m.basicParse(input)
	m.saveToCache(cacheKey, basicResult)
	return basicResult, nil
}

// SuggestTasks generates contextual task suggestions
func (m *Manager) SuggestTasks(ctx context.Context, existingTasks []string) ([]string, error) {
	// Create cache key from tasks
	taskStr := strings.Join(existingTasks, "|")
	cacheKey := m.generateCacheKey("suggest", taskStr)

	// Check cache first
	if cached, found := m.getFromCache(cacheKey); found {
		if suggestions, ok := cached.([]string); ok {
			return suggestions, nil
		}
	}

	// Use only qwen2.5:1.5b model
	m.model = "qwen2.5:1.5b"
	if result, err := m.suggestWithOllama(ctx, existingTasks); err == nil {
		// Cache the result
		m.saveToCache(cacheKey, result)
		return result, nil
	}

	// Fallback to OpenRouter with user approval
	if m.openRouterKey != "" {
		if m.requestRemoteApproval("task suggestions") {
			if result, err := m.suggestWithOpenRouter(ctx, existingTasks); err == nil {
				// Cache the result
				m.saveToCache(cacheKey, result)
				return result, nil
			}
		}
	}

	// If both fail, return empty suggestions
	return []string{}, nil
}

// SummarizeTasks generates a summary of tasks using AI
func (m *Manager) SummarizeTasks(ctx context.Context, tasks []string) (string, error) {
	// Create cache key from tasks
	taskStr := strings.Join(tasks, "|")
	cacheKey := m.generateCacheKey("summary", taskStr)

	// Check cache first
	if cached, found := m.getFromCache(cacheKey); found {
		if summary, ok := cached.(string); ok {
			return summary, nil
		}
	}

	// Use only qwen2.5:1.5b model
	m.model = "qwen2.5:1.5b"
	if result, err := m.summarizeWithOllama(ctx, tasks); err == nil {
		// Cache the result
		m.saveToCache(cacheKey, result)
		return result, nil
	}

	// Fallback to OpenRouter with user approval
	if m.openRouterKey != "" {
		if m.requestRemoteApproval("task summary") {
			if result, err := m.summarizeWithOpenRouter(ctx, tasks); err == nil {
				// Cache the result
				m.saveToCache(cacheKey, result)
				return result, nil
			}
		}
	}

	// If both fail, return basic summary
	summary := m.basicSummary(tasks)
	m.saveToCache(cacheKey, summary)
	return summary, nil
}

// ChatAssistant provides help and answers questions about tasks and app usage
func (m *Manager) ChatAssistant(ctx context.Context, question string, tasks []string) (string, error) {
	// Create cache key from question and tasks
	taskStr := strings.Join(tasks, "|")
	cacheKey := m.generateCacheKey("chat", question+taskStr)

	// Check cache first
	if cached, found := m.getFromCache(cacheKey); found {
		if response, ok := cached.(string); ok {
			return response, nil
		}
	}

	// Use only qwen2.5:1.5b model
	m.model = "qwen2.5:1.5b"
	if result, err := m.chatWithOllama(ctx, question, tasks); err == nil {
		// Cache the result
		m.saveToCache(cacheKey, result)
		return result, nil
	}

	// Fallback to OpenRouter with user approval
	if m.openRouterKey != "" {
		if m.requestRemoteApproval("chat assistant") {
			if result, err := m.chatWithOpenRouter(ctx, question, tasks); err == nil {
				// Cache the result
				m.saveToCache(cacheKey, result)
				return result, nil
			}
		}
	}

	// If both fail, return basic response with more context
	if strings.Contains(strings.ToLower(question), "help") || strings.Contains(strings.ToLower(question), "hi") {
		return "Hello! I'm your focus.sh AI assistant. I can help you with:\n• Task management and organization\n• Productivity tips\n• App usage guidance\n• Smart suggestions\n\nTry asking me about your tasks or how to use specific features!", nil
	}

	if strings.Contains(strings.ToLower(question), "task") {
		if len(tasks) > 0 {
			return fmt.Sprintf("I see you have %d tasks. I can help you organize, prioritize, or complete them. Try using 'focus inspire' for suggestions or 'focus list' to see all your tasks!", len(tasks)), nil
		} else {
			return "You don't have any tasks yet! Start by adding a mission with 'focus add \"your task here\"' and I can help you manage them.", nil
		}
	}

	return fmt.Sprintf("I'm having trouble connecting to AI services right now (Ollama not available). However, I can tell you that you have %d tasks. Try 'focus --help' to see available commands!", len(tasks)), nil
}

// requestRemoteApproval asks user for permission before using remote API
func (m *Manager) requestRemoteApproval(feature string) bool {
	fmt.Printf("\n⚠️  Remote AI required for %s\n", feature)
	fmt.Printf("   This will use OpenRouter API (costs may apply)\n")
	fmt.Print("   Continue? (y/N): ")

	reader := bufio.NewReader(os.Stdin)
	response, err := reader.ReadString('\n')
	if err != nil {
		return false
	}

	response = strings.TrimSpace(strings.ToLower(response))
	return response == "y" || response == "yes"
}

// makeOllamaRequest is a helper function to make requests to Ollama
// Returns the raw response string from the model
func (m *Manager) makeOllamaRequest(ctx context.Context, prompt string) (string, error) {
	// Check if model is available
	if !m.isModelAvailable(m.model) {
		// Try to pull model automatically
		if err := m.pullModel(m.model); err != nil {
			return "", fmt.Errorf("model %s not available and could not be pulled: %w", m.model, err)
		}
	}

	payload := map[string]any{
		"model":  m.model,
		"prompt": prompt,
		"stream": false,
		"format": "json",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", m.ollamaURL, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama request failed with status %d", resp.StatusCode)
	}

	var result struct {
		Response string `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Response, nil
}

// makeOpenRouterRequest is a helper function to make requests to OpenRouter
// Returns the raw response string from the model
func (m *Manager) makeOpenRouterRequest(ctx context.Context, prompt string) (string, error) {
	payload := map[string]any{
		"model": "meta-llama/llama-3.2-3b-instruct:free",
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+m.openRouterKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", "https://github.com/kyanite/focus")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("openrouter request failed with status %d", resp.StatusCode)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from openrouter")
	}

	return result.Choices[0].Message.Content, nil
}

// parseWithOllama uses local Ollama for task parsing
func (m *Manager) parseWithOllama(ctx context.Context, input string) (*ParsedTask, error) {
	prompt := m.buildParsePrompt(input)

	response, err := m.makeOllamaRequest(ctx, prompt)
	if err != nil {
		return nil, err
	}

	var parsedTask ParsedTask
	if err := json.Unmarshal([]byte(response), &parsedTask); err != nil {
		// Try to extract JSON from the response if it's wrapped
		content := m.extractJSON(response)
		if err := json.Unmarshal([]byte(content), &parsedTask); err != nil {
			return nil, err
		}
	}

	return &parsedTask, nil
}

// IsOllamaAvailable checks if Ollama is running and accessible
func (m *Manager) IsOllamaAvailable() bool {
	// Check if Ollama API is responding
	url := "http://localhost:11434/api/tags"
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer func() { _ = resp.Body.Close() }()
	return resp.StatusCode == http.StatusOK
}

// LaunchOllama attempts to start Ollama if it's not running
func (m *Manager) LaunchOllama() error {
	if m.IsOllamaAvailable() {
		return nil // Already running
	}

	// Use async launch with proper startup detection
	return m.launchOllamaAsync()
}

// launchOllamaAsync starts Ollama and waits for it to be available
func (m *Manager) launchOllamaAsync() error {
	var cmd *exec.Cmd

	// Choose command based on OS
	if runtime.GOOS == "windows" {
		cmd = exec.Command("powershell", "-Command", "Start-Process", "ollama")
	} else {
		cmd = exec.Command("ollama", "serve")
	}

	// Start process without context timeout since we want it to persist
	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start ollama: %w", err)
	}

	// Wait for Ollama to be available with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	return m.waitForOllamaAvailable(ctx)
}

// waitForOllamaAvailable polls until Ollama is available or context times out
func (m *Manager) waitForOllamaAvailable(ctx context.Context) error {
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return fmt.Errorf("timeout waiting for ollama to start")
		case <-ticker.C:
			if m.IsOllamaAvailable() {
				return nil
			}
		}
	}
}

// runCommand executes a command in a separate process
func (m *Manager) runCommand(ctx context.Context, cmdArgs ...string) error {
	if len(cmdArgs) == 0 {
		return fmt.Errorf("no command provided")
	}

	cmd := exec.CommandContext(ctx, cmdArgs[0], cmdArgs[1:]...)
	return cmd.Start()
}

// isModelAvailable checks if a model is available in Ollama
func (m *Manager) isModelAvailable(modelName string) bool {
	// Fast path: check cache
	m.modelMutex.RLock()
	if available, exists := m.modelAvailable[modelName]; exists {
		m.modelMutex.RUnlock()
		return available
	}
	m.modelMutex.RUnlock()

	// Slow path: fetch from Ollama without holding lock
	available := m.checkModelWithOllama(modelName)

	// Update cache
	m.modelMutex.Lock()
	m.modelAvailable[modelName] = available
	m.modelMutex.Unlock()

	return available
}

// checkModelWithOllama queries Ollama API to check if a model is available
func (m *Manager) checkModelWithOllama(modelName string) bool {
	url := "http://localhost:11434/api/tags"
	resp, err := http.Get(url)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return false
	}

	var tags struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return false
	}

	for _, model := range tags.Models {
		if strings.Contains(model.Name, modelName) {
			return true
		}
	}

	return false
}

// pullModel automatically pulls a model through Ollama
func (m *Manager) pullModel(modelName string) error {
	fmt.Printf("🤖 Pulling model %s automatically...\n", modelName)

	url := "http://localhost:11434/api/pull"
	payload := map[string]any{
		"name": modelName,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	// Note: Using shared client; for longer operations consider using context with timeout
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("ollama pull failed with status %d", resp.StatusCode)
	}

	// Mark model as available
	m.modelMutex.Lock()
	m.modelAvailable[modelName] = true
	m.modelMutex.Unlock()

	return nil
}

// suggestWithOllama uses local Ollama for task suggestions
func (m *Manager) suggestWithOllama(ctx context.Context, existingTasks []string) ([]string, error) {
	prompt := m.buildSuggestPrompt(existingTasks)

	payload := map[string]any{
		"model":  m.model,
		"prompt": prompt,
		"stream": false,
		"format": "json",
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", m.ollamaURL, bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama request failed with status %d", resp.StatusCode)
	}

	var result struct {
		Response string `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	var suggestions struct {
		Tasks []string `json:"tasks"`
	}

	if err := json.Unmarshal([]byte(result.Response), &suggestions); err != nil {
		// Try to extract JSON from the response if it's wrapped
		content := m.extractJSON(result.Response)
		if err := json.Unmarshal([]byte(content), &suggestions); err != nil {
			return nil, err
		}
	}

	// Filter out low quality suggestions
	filtered := m.filterLowQualitySuggestions(suggestions.Tasks)

	return filtered, nil
}

// chatWithOllama uses local Ollama for chat assistance
func (m *Manager) chatWithOllama(ctx context.Context, question string, tasks []string) (string, error) {
	prompt := m.buildChatPrompt(question, tasks)

	payload := map[string]any{
		"model":  m.model,
		"prompt": prompt,
		"stream": false,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", m.ollamaURL, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama request failed with status %d", resp.StatusCode)
	}

	var result struct {
		Response string `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return strings.TrimSpace(result.Response), nil
}

// summarizeWithOllama uses local Ollama for task summary
func (m *Manager) summarizeWithOllama(ctx context.Context, tasks []string) (string, error) {
	prompt := m.buildSummaryPrompt(tasks)

	payload := map[string]any{
		"model":  m.model,
		"prompt": prompt,
		"stream": false,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", m.ollamaURL, bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")

	// Use shared HTTP client with connection pooling
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ollama request failed with status %d", resp.StatusCode)
	}

	var result struct {
		Response string `json:"response"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Response, nil
}

// parseWithOpenRouter uses OpenRouter as fallback
func (m *Manager) parseWithOpenRouter(ctx context.Context, input string) (*ParsedTask, error) {
	prompt := m.buildParsePrompt(input)

	messages := []map[string]string{
		{
			"role":    "user",
			"content": prompt,
		},
	}

	payload := map[string]any{
		"model":    "openai/gpt-3.5-turbo",
		"messages": messages,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+m.openRouterKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", "https://github.com/kyanite/focus")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openrouter request failed with status %d", resp.StatusCode)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("no response from openrouter")
	}

	// Extract JSON from the response
	content := result.Choices[0].Message.Content

	// Handle markdown code blocks if present
	if len(content) > 6 && content[:3] == "```" {
		lines := strings.Split(content, "\n")
		if len(lines) > 2 {
			content = strings.Join(lines[1:len(lines)-1], "\n")
		}
	}

	var parsedTask ParsedTask
	if err := json.Unmarshal([]byte(content), &parsedTask); err != nil {
		return nil, err
	}

	return &parsedTask, nil
}

// suggestWithOpenRouter uses OpenRouter for task suggestions
func (m *Manager) suggestWithOpenRouter(ctx context.Context, existingTasks []string) ([]string, error) {
	prompt := m.buildSuggestPrompt(existingTasks)

	messages := []map[string]string{
		{
			"role":    "user",
			"content": prompt,
		},
	}

	payload := map[string]any{
		"model":    "openai/gpt-3.5-turbo",
		"messages": messages,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+m.openRouterKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", "https://github.com/kyanite/focus")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("openrouter request failed with status %d", resp.StatusCode)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("no response from openrouter")
	}

	// Extract JSON from the response
	content := result.Choices[0].Message.Content

	// Handle markdown code blocks if present
	if len(content) > 6 && content[:3] == "```" {
		lines := strings.Split(content, "\n")
		if len(lines) > 2 {
			content = strings.Join(lines[1:len(lines)-1], "\n")
		}
	}

	var suggestions struct {
		Tasks []string `json:"tasks"`
	}

	if err := json.Unmarshal([]byte(content), &suggestions); err != nil {
		return nil, err
	}

	// Filter out low quality suggestions
	filtered := m.filterLowQualitySuggestions(suggestions.Tasks)

	return filtered, nil
}

// chatWithOpenRouter uses OpenRouter for chat assistance
func (m *Manager) chatWithOpenRouter(ctx context.Context, question string, tasks []string) (string, error) {
	prompt := m.buildChatPrompt(question, tasks)

	messages := []map[string]string{
		{
			"role":    "user",
			"content": prompt,
		},
	}

	payload := map[string]any{
		"model":    "openai/gpt-3.5-turbo",
		"messages": messages,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+m.openRouterKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", "https://github.com/kyanite/focus")

	resp, err := m.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("openrouter request failed with status %d", resp.StatusCode)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from openrouter")
	}

	return strings.TrimSpace(result.Choices[0].Message.Content), nil
}

// summarizeWithOpenRouter uses OpenRouter for task summary
func (m *Manager) summarizeWithOpenRouter(ctx context.Context, tasks []string) (string, error) {
	prompt := m.buildSummaryPrompt(tasks)

	messages := []map[string]string{
		{
			"role":    "user",
			"content": prompt,
		},
	}

	payload := map[string]any{
		"model":    "openai/gpt-3.5-turbo",
		"messages": messages,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", "https://openrouter.ai/api/v1/chat/completions", bytes.NewBuffer(body))
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+m.openRouterKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", "https://github.com/kyanite/focus")

	// Use shared HTTP client with connection pooling
	resp, err := m.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("openrouter request failed with status %d", resp.StatusCode)
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if len(result.Choices) == 0 {
		return "", fmt.Errorf("no response from openrouter")
	}

	return result.Choices[0].Message.Content, nil
}

// buildParsePrompt creates the prompt for task parsing with enhanced deadline parsing
func (m *Manager) buildParsePrompt(input string) string {
	return fmt.Sprintf(`You are a task parsing assistant for focus.sh. Extract task information from user input.
 
Guidelines:
- Be concise and direct
- Focus on action words
- Parse natural language dates
- Set priority to medium if unclear
- Extract up to 3 meaningful categories
 
Format response as JSON:
{
  "description": "clear task description (2-200 chars)",
  "deadline": "ISO date or null",
  "priority": "low|medium|high",
  "categories": ["category1", "category2"]
}
 
User input: %s
 
JSON Response:`, input)
}

// buildSuggestPrompt creates the prompt for task suggestions with context awareness
func (m *Manager) buildSuggestPrompt(existingTasks []string) string {
	tasksStr := "[]"
	if len(existingTasks) > 0 {
		tasksStr = fmt.Sprintf("[\"%s\"]", strings.Join(existingTasks, "\", \""))
	}

	return fmt.Sprintf(`You are a task suggestion assistant for focus.sh. Generate 3-5 relevant follow-up tasks.
	
Guidelines:
- Create actionable items starting with verbs
- Make tasks specific and manageable
- Consider user's existing tasks
- Focus on one task at a time
	
Existing tasks: %s
	
Format response as JSON:
{
	 "tasks": [
	   "actionable task 1",
	   "actionable task 2",
	   "actionable task 3"
	 ]
}
	
JSON Response:`, tasksStr)
}

// buildSummaryPrompt creates the prompt for task summary with advanced insights
func (m *Manager) buildSummaryPrompt(tasks []string) string {
	tasksStr := "[]"
	if len(tasks) > 0 {
		tasksStr = fmt.Sprintf("[\"%s\"]", strings.Join(tasks, "\", \""))
	}

	return fmt.Sprintf(`You are a productivity analyst for focus.sh. Provide a concise task summary.
	
Required Analysis:
1. Total tasks and completion ratio
2. Priority distribution
3. 1-2 key insights or patterns
4. 1 specific recommendation
	
Tasks: %s
	
Provide exactly 3-4 sentences with clear insights. No filler.
	
Summary:`, tasksStr)
}

// buildChatPrompt creates the prompt for chat assistance
func (m *Manager) buildChatPrompt(question string, tasks []string) string {
	tasksStr := "[]"
	if len(tasks) > 0 {
		tasksStr = fmt.Sprintf("[\"%s\"]", strings.Join(tasks, "\", \""))
	}

	return fmt.Sprintf(`You are Focus, a helpful AI assistant for task management and productivity. You have three main roles:

1. 📋 TASK HELPER - Help with task management, planning, and productivity
2. 🎮 APP GUIDE - Help with using the focus.sh application features  
3. 🤖 GENERAL ASSISTANT - Be a helpful AI assistant for various questions

PERSONALITY:
- Helpful and professional, but friendly and approachable
- Clear and concise in your responses
- Knowledgeable about productivity, task management, and technology
- Encouraging and positive

SPEECH STYLE:
- Be direct and practical
- Use clear, simple language
- Provide actionable advice when possible
- Be encouraging without being overly casual
- Use appropriate emojis to enhance clarity

CURRENT USER TASKS (only mention if relevant to the conversation):
User tasks: %s

focus.sh APP COMMANDS:
- add "task" - Add new task
- list --filter=active/completed/all - List tasks
- done <task-id> - Complete task
- remove <task-id> - Delete task
- priority <task-id> <high/medium/low> - Set priority
- notes <task-id> - Edit task notes
- dashboard - Launch TUI dashboard
- K - Open calendar (in TUI), rad synthwave month view
- N - Edit notes (in TUI)
- / - Filter tasks (in TUI)
- C - Chat assistant (in TUI)
- help - Show all commands
	
RESPONSE GUIDELINES:
- If user asks about tasks: be helpful and professional
- If user asks about app: explain focus.sh features with synthwave vibe
- If general question: be a confident, knowledgeable AI with personality
- You can code, explain tech concepts, give advice, be creative
- Add some synthwave flavor when it fits
	
USER QUESTION: %s
	
Give a rad synthwave-cyberpunk response that's actually helpful. Be confident, engaging, and maybe a little edgy.`, tasksStr, question)
}

// validateResponse checks for hallucinations and invalid responses with enhanced detection
func (m *Manager) validateResponse(task *ParsedTask) (*ParsedTask, bool) {
	if task == nil {
		return nil, false
	}

	// Check description is not empty
	if task.Description == "" {
		return nil, false
	}

	// Check description is reasonable (not too short or too long)
	if len(task.Description) < 2 || len(task.Description) > 200 {
		return nil, false
	}

	// Enhanced description quality check
	if m.isLowQualityDescription(task.Description) {
		return nil, false
	}

	// Check priority is valid
	validPriorities := map[string]bool{
		"low":    true,
		"medium": true,
		"high":   true,
	}

	if task.Priority == "" {
		task.Priority = "medium" // Default priority
	} else if !validPriorities[task.Priority] {
		task.Priority = "medium" // Normalize invalid priority
	}

	// Basic category validation
	validCategories := []string{}
	for _, category := range task.Categories {
		if len(category) >= 2 && len(category) <= 20 {
			validCategories = append(validCategories, category)
		}
	}
	task.Categories = validCategories

	// Validate deadline if present
	if !task.Deadline.IsZero() {
		now := time.Now()
		// Check if deadline is in the past (unreasonable)
		if task.Deadline.Before(now.Add(-24 * time.Hour)) {
			task.Deadline = time.Time{} // Clear invalid deadline
		}
		// Check if deadline is too far in the future (likely hallucination)
		if task.Deadline.After(now.Add(365 * 24 * time.Hour)) {
			task.Deadline = time.Time{} // Clear unreasonable deadline
		}
	}

	return task, true
}

// isLowQualityDescription checks if a description is likely hallucinated or low quality
func (m *Manager) isLowQualityDescription(description string) bool {
	// Check for repeated characters (hallucination indicator)
	if strings.Contains(description, "aaaa") || strings.Contains(description, "bbbb") {
		return true
	}

	// Check for placeholder text
	lowerDesc := strings.ToLower(description)
	placeholderIndicators := []string{"task description", "placeholder", "example task"}
	for _, indicator := range placeholderIndicators {
		if strings.Contains(lowerDesc, indicator) {
			return true
		}
	}

	// Check for meaningless text (all same words, etc.)
	words := strings.Fields(description)
	if len(words) == 0 {
		return true
	}

	// If all words are the same (likely hallucination)
	firstWord := strings.ToLower(words[0])
	allSame := true
	for _, word := range words {
		if strings.ToLower(word) != firstWord {
			allSame = false
			break
		}
	}
	if allSame && len(words) > 3 {
		return true
	}

	return false
}

// filterLowQualitySuggestions removes low quality or hallucinated suggestions
func (m *Manager) filterLowQualitySuggestions(suggestions []string) []string {
	filtered := []string{}
	for _, suggestion := range suggestions {
		// Skip empty suggestions
		if suggestion == "" {
			continue
		}

		// Skip very short suggestions
		if len(suggestion) < 5 {
			continue
		}

		// Skip very long suggestions
		if len(suggestion) > 100 {
			continue
		}

		// Skip suggestions with placeholder text
		lowerSuggestion := strings.ToLower(suggestion)
		placeholderIndicators := []string{"example", "placeholder", "suggested task"}
		isPlaceholder := false
		for _, indicator := range placeholderIndicators {
			if strings.Contains(lowerSuggestion, indicator) {
				isPlaceholder = true
				break
			}
		}

		if !isPlaceholder {
			filtered = append(filtered, suggestion)
		}
	}

	// Return at least empty slice if all suggestions were filtered out
	if len(filtered) == 0 {
		return []string{}
	}

	return filtered
}

// basicParse provides fallback parsing when AI fails
func (m *Manager) basicParse(input string) *ParsedTask {
	// Very basic parsing - extract first few words as description
	words := strings.Fields(input)
	if len(words) > 10 {
		input = strings.Join(words[:10], " ") + "..."
	}

	return &ParsedTask{
		Description: input,
		Priority:    "medium",
	}
}

// basicSummary provides fallback summarization when AI fails
func (m *Manager) basicSummary(tasks []string) string {
	completed := 0
	pending := 0

	for _, task := range tasks {
		if strings.Contains(task, "completed") {
			completed++
		} else {
			pending++
		}
	}

	return fmt.Sprintf("You have %d tasks (%d completed, %d pending). Keep up the good work!",
		len(tasks), completed, pending)
}

// extractJSON tries to extract JSON from a response that might be wrapped in text
func (m *Manager) extractJSON(response string) string {
	// Look for JSON-like content in the response
	start := strings.Index(response, "{")
	end := strings.LastIndex(response, "}")

	if start != -1 && end != -1 && end > start {
		return response[start : end+1]
	}

	return response
}

// generateCacheKey creates a unique cache key for a request
func (m *Manager) generateCacheKey(operation, input string) string {
	key := fmt.Sprintf("%s:%s:%s", operation, m.model, input)
	hash := md5.Sum([]byte(key))
	return hex.EncodeToString(hash[:])
}

// getFromCache retrieves a response from cache if available and not expired
func (m *Manager) getFromCache(key string) (interface{}, bool) {
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()

	entry, exists := m.cache[key]
	if !exists {
		return nil, false
	}

	// Check if entry is expired
	if time.Now().After(entry.ExpiresAt) {
		return nil, false
	}

	// Update access time for LRU
	m.cacheHits[key] = time.Now()
	entry.AccessCnt++
	m.cache[key] = entry

	return entry.Response, true
}

// saveToCache stores a response in cache with expiration
func (m *Manager) saveToCache(key string, response interface{}) {
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()

	// Evict LRU entry if cache is full
	if len(m.cache) >= m.cacheMaxSize {
		m.evictLRUEntryUnsafe()
	}

	now := time.Now()
	m.cache[key] = cacheEntry{
		Response:  response,
		Timestamp: now,
		ExpiresAt: now.Add(24 * time.Hour), // Cache for 24 hours
		AccessCnt: 1,
	}
	m.cacheHits[key] = now
	m.cacheDirty = true // Mark cache as needing save
}

// evictLRUEntryUnsafe removes the least recently used cache entry
// Must be called with cacheMutex held
func (m *Manager) evictLRUEntryUnsafe() {
	if len(m.cache) == 0 {
		return
	}

	// Find LRU entry
	var oldestKey string
	var oldestTime time.Time
	first := true

	for key, accessTime := range m.cacheHits {
		if first || accessTime.Before(oldestTime) {
			oldestKey = key
			oldestTime = accessTime
			first = false
		}
	}

	// Remove oldest entry
	if oldestKey != "" {
		delete(m.cache, oldestKey)
		delete(m.cacheHits, oldestKey)
	}
}

// loadCache loads cached responses from file
func (m *Manager) loadCache() {
	m.cacheMutex.Lock()
	defer m.cacheMutex.Unlock()

	data, err := os.ReadFile(m.cachePath)
	if err != nil {
		return
	}

	var cache map[string]cacheEntry
	if err := json.Unmarshal(data, &cache); err != nil {
		return
	}

	// Filter out expired entries and initialize cacheHits
	now := time.Now()
	for key, entry := range cache {
		if now.After(entry.ExpiresAt) {
			delete(cache, key)
		} else {
			m.cacheHits[key] = entry.Timestamp
		}
	}

	m.cache = cache
	m.cacheDirty = false
}

// periodicCacheSaver saves cache to disk every 5 minutes
func (m *Manager) periodicCacheSaver() {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		m.cacheMutex.RLock()
		needsSave := m.cacheDirty
		m.cacheMutex.RUnlock()

		if needsSave {
			m.cacheMutex.Lock()
			if m.cacheDirty {
				m.saveCache()
				m.cacheDirty = false
			}
			m.cacheMutex.Unlock()
		}
	}
}

// saveCache persists cache to file
func (m *Manager) saveCache() {
	// Ensure directory exists
	dir := filepath.Dir(m.cachePath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return
	}

	data, err := json.Marshal(m.cache)
	if err != nil {
		return
	}

	os.WriteFile(m.cachePath, data, 0644)
}
