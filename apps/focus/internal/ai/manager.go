package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	ai "github.com/kyanite/ai"
	"github.com/kyanite/cache"
	"github.com/kyanite/config"
)

// Manager handles AI interactions via the shared pkg/ai Brain.
type Manager struct {
	brain          *ai.Brain
	promptBuilder  *PromptBuilder
	validator      *TaskValidator
	cache          *cache.LRU
}

// ParsedTask represents the structured output from the LLM.
type ParsedTask struct {
	Description string    `json:"description"`
	Deadline    time.Time `json:"deadline,omitempty"`
	Priority    string    `json:"priority"`
	Categories  []string  `json:"categories,omitempty"`
}

// New creates a new AI manager backed by the shared pkg/ai Brain.
func New() *Manager {
	root, _ := config.Load()
	cfg := ai.ConfigFromRoot(root, "focus")
	brain, err := ai.New(cfg)
	if err != nil {
		fmt.Fprintf(os.Stderr, "warning: brain init failed (AI features offline): %v\n", err)
	}

	home, _ := os.UserHomeDir()
	cachePath := filepath.Join(home, ".focus", "ai_cache.json")

	return &Manager{
		brain:         brain,
		promptBuilder: NewPromptBuilder(),
		validator:     NewTaskValidator(),
		cache:         cache.NewLRU(500, 24*time.Hour, cachePath),
	}
}

// NewWithBrain creates a Manager backed by an explicit Brain. Used by
// tests with pkg/testutil.MockBrain so the AI path is exercised
// without hitting a real Ollama. The cache is memory-only (no file
// path) so tests don't pollute a shared on-disk cache. (T7-05)
func NewWithBrain(brain *ai.Brain) *Manager {
	return &Manager{
		brain:         brain,
		promptBuilder: NewPromptBuilder(),
		validator:     NewTaskValidator(),
		cache:         cache.NewLRU(500, 24*time.Hour, ""),
	}
}

// Close releases Brain resources and flushes the cache.
func (m *Manager) Close() {
	m.cache.Close()
	if m.brain != nil {
		m.brain.Close()
	}
}

// ParseTask converts natural language to structured task using the LLM.
func (m *Manager) ParseTask(ctx context.Context, input string) (*ParsedTask, error) {
	cacheKey := cache.GenerateKey("parse", input)
	if cached, found := m.cache.Get(cacheKey); found {
		if task, ok := cached.(*ParsedTask); ok {
			return task, nil
		}
	}

	if m.brain != nil && m.brain.IsLLMAvailable(ctx) {
		prompt := m.promptBuilder.BuildParsePrompt(input)
		resp, err := m.brain.Generate(ctx, prompt, ai.WithJSONMode())
		if err == nil {
		cleaned := extractJSONFromResponse(resp)
			// Two-phase unmarshal: LLM may return "deadline":""
			// which fails time.Time parsing. Strip it first.
			var raw map[string]json.RawMessage
			if jsonErr := json.Unmarshal([]byte(cleaned), &raw); jsonErr == nil {
				if dl, ok := raw["deadline"]; ok && string(dl) == `""` {
					delete(raw, "deadline")
				}
				fixed, _ := json.Marshal(raw)
				var task ParsedTask
				if jsonErr = json.Unmarshal(fixed, &task); jsonErr == nil {
					if validated, ok := m.validateResponse(&task); ok {
						m.cache.Set(cacheKey, validated)
						return validated, nil
					}
				}
			}
		}
	}

	basicResult := m.basicParse(input)
	m.cache.Set(cacheKey, basicResult)
	return basicResult, nil
}

// SuggestTasks generates contextual task suggestions using the LLM.
func (m *Manager) SuggestTasks(ctx context.Context, existingTasks []string) ([]string, error) {
	cacheKey := cache.GenerateKey("suggest", strings.Join(existingTasks, "|"))
	if cached, found := m.cache.Get(cacheKey); found {
		if suggestions, ok := cached.([]string); ok {
			return suggestions, nil
		}
	}

	if m.brain != nil && m.brain.IsLLMAvailable(ctx) {
		prompt := m.promptBuilder.BuildSuggestPrompt(existingTasks)
		resp, err := m.brain.Generate(ctx, prompt)
		if err == nil {
			suggestions := parseSuggestionList(resp)
			if len(suggestions) > 0 {
				m.cache.Set(cacheKey, suggestions)
				return suggestions, nil
			}
		}
	}

	fallback := []string{
		"Review today's highest-priority task",
		"Break one blocked item into a smaller next step",
	}
	m.cache.Set(cacheKey, fallback)
	return fallback, nil
}

// SummarizeTasks generates a summary of tasks using the LLM.
func (m *Manager) SummarizeTasks(ctx context.Context, tasks []string) (string, error) {
	cacheKey := cache.GenerateKey("summary", strings.Join(tasks, "|"))
	if cached, found := m.cache.Get(cacheKey); found {
		if summary, ok := cached.(string); ok {
			return summary, nil
		}
	}

	if m.brain != nil && m.brain.IsLLMAvailable(ctx) {
		prompt := m.promptBuilder.BuildSummaryPrompt(tasks)
		resp, err := m.brain.Generate(ctx, prompt)
		if err == nil {
			summary := strings.TrimSpace(resp)
			if summary != "" {
				m.cache.Set(cacheKey, summary)
				return summary, nil
			}
		}
	}

	summary := m.basicSummary(tasks)
	m.cache.Set(cacheKey, summary)
	return summary, nil
}

// ChatAssistant provides help and answers questions about tasks and app usage.
func (m *Manager) ChatAssistant(ctx context.Context, question string, tasks []string) (string, error) {
	cacheKey := cache.GenerateKey("chat", question+strings.Join(tasks, "|"))
	if cached, found := m.cache.Get(cacheKey); found {
		if response, ok := cached.(string); ok {
			return response, nil
		}
	}

	if m.brain != nil && m.brain.IsLLMAvailable(ctx) {
		prompt := m.promptBuilder.BuildChatPrompt(question, tasks)
		result, err := m.brain.Generate(ctx, prompt)
		if err == nil {
			result = strings.TrimSpace(result)
			if result != "" {
				m.cache.Set(cacheKey, result)
				return result, nil
			}
		}
	}

	response := m.fallbackChatResponse(question, tasks)
	m.cache.Set(cacheKey, response)
	return response, nil
}

// fallbackChatResponse returns a deterministic answer when the LLM is unavailable.
func (m *Manager) fallbackChatResponse(question string, tasks []string) string {
	lower := strings.ToLower(question)
	if strings.Contains(lower, "help") || strings.Contains(lower, "hi") {
		return "Hello! I'm your focus.sh AI assistant. I can help you with:\n• Task management and organization\n• Productivity tips\n• App usage guidance\n• Smart suggestions\n\nTry asking me about your tasks or how to use specific features!"
	}

	if strings.Contains(lower, "task") {
		if len(tasks) > 0 {
			return fmt.Sprintf("I see you have %d tasks. I can help you organize, prioritize, or complete them. Try using 'focus inspire' for suggestions or 'focus list' to see all your tasks!", len(tasks))
		}
		return "You don't have any tasks yet! Start by adding a mission with 'focus add \"your task here\"' and I can help you manage them."
	}

	return fmt.Sprintf("AI is currently unavailable. You have %d tasks. Try 'focus --help' to see available commands!", len(tasks))
}

// GetCrossAppContext retrieves recent context from other kyanite apps.
func (m *Manager) GetCrossAppContext(ctx context.Context, limit int) []ai.CrossAppContext {
	if m.brain == nil {
		return nil
	}
	contexts, err := m.brain.GetCrossAppContext(ctx, limit)
	if err != nil {
		return nil
	}
	return contexts
}

// IsOllamaAvailable reports whether the AI backend is reachable.
func (m *Manager) IsOllamaAvailable() bool {
	return m.brain != nil && m.brain.IsLLMAvailable(context.Background())
}

// validateResponse delegates to TaskValidator
func (m *Manager) validateResponse(task *ParsedTask) (*ParsedTask, bool) {
	return m.validator.Validate(task)
}

// basicParse provides fallback parsing when AI fails
func (m *Manager) basicParse(input string) *ParsedTask {
	task := &ParsedTask{
		Description: input,
		Priority:    "medium",
	}

	lowerInput := strings.ToLower(input)
	if strings.Contains(lowerInput, "urgent") || strings.Contains(lowerInput, "asap") ||
		strings.Contains(lowerInput, "critical") || strings.Contains(lowerInput, "emergency") {
		task.Priority = "high"
	} else if strings.Contains(lowerInput, "low priority") || strings.Contains(lowerInput, "when possible") ||
		strings.Contains(lowerInput, "sometime") {
		task.Priority = "low"
	}

	var categories []string
	if strings.Contains(lowerInput, "work") || strings.Contains(lowerInput, "job") {
		categories = append(categories, "work")
	}
	if strings.Contains(lowerInput, "personal") || strings.Contains(lowerInput, "home") {
		categories = append(categories, "personal")
	}
	if strings.Contains(lowerInput, "meeting") {
		categories = append(categories, "meetings")
	}
	task.Categories = categories

	if strings.Contains(lowerInput, "today") {
		task.Deadline = time.Now()
	} else if strings.Contains(lowerInput, "tomorrow") {
		task.Deadline = time.Now().AddDate(0, 0, 1)
	} else if strings.Contains(lowerInput, "next week") {
		task.Deadline = time.Now().AddDate(0, 0, 7)
	}

	words := strings.Fields(task.Description)
	if len(words) > 10 {
		task.Description = strings.Join(words[:10], " ") + "..."
	}

	return task
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

// extractJSONFromResponse extracts the first JSON object from a response.
func extractJSONFromResponse(response string) string {
	start := strings.Index(response, "{")
	end := strings.LastIndex(response, "}")

	if start != -1 && end != -1 && end > start {
		return response[start : end+1]
	}

	return response
}

// parseSuggestionList splits a multi-line LLM response into a clean slice.
func parseSuggestionList(response string) []string {
	var out []string
	for _, line := range strings.Split(response, "\n") {
		cleaned := strings.TrimSpace(line)
		cleaned = strings.TrimLeft(cleaned, "-*•\t ")
		if idx := strings.IndexAny(cleaned, ".)"); idx == len(cleaned)-1 {
			cleaned = ""
		}
		cleaned = strings.TrimSpace(cleaned)
		if cleaned != "" {
			out = append(out, cleaned)
		}
	}
	return out
}

// Brain returns the underlying Brain instance for direct access (e.g., AI panel streaming).
func (m *Manager) Brain() *ai.Brain {
	if m == nil {
		return nil
	}
	return m.brain
}
