package ai

import (
	"context"
	"fmt"
	"strings"
	"time"
)

// FallbackProvider implements a simple rule-based provider when AI is unavailable
type FallbackProvider struct{}

// NewFallbackProvider creates a new fallback provider
func NewFallbackProvider() *FallbackProvider {
	return &FallbackProvider{}
}

// ParseTask provides basic rule-based task parsing
func (p *FallbackProvider) ParseTask(ctx context.Context, input string) (*ParsedTask, error) {
	task := &ParsedTask{
		Description: input,
		Priority:    "medium", // Default priority
	}

	// Simple keyword-based priority detection
	lowerInput := strings.ToLower(input)
	if strings.Contains(lowerInput, "urgent") || strings.Contains(lowerInput, "asap") ||
		strings.Contains(lowerInput, "critical") || strings.Contains(lowerInput, "emergency") {
		task.Priority = "high"
	} else if strings.Contains(lowerInput, "low priority") || strings.Contains(lowerInput, "when possible") ||
		strings.Contains(lowerInput, "sometime") {
		task.Priority = "low"
	}

	// Simple category detection
	categories := []string{}
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

	// Simple deadline detection (very basic)
	if strings.Contains(lowerInput, "today") {
		task.Deadline = time.Now()
	} else if strings.Contains(lowerInput, "tomorrow") {
		task.Deadline = time.Now().AddDate(0, 0, 1)
	} else if strings.Contains(lowerInput, "next week") {
		task.Deadline = time.Now().AddDate(0, 0, 7)
	}

	return task, nil
}

// ChatAssistant provides simple canned responses
func (p *FallbackProvider) ChatAssistant(ctx context.Context, message string, taskContext []string) (string, error) {
	lowerMsg := strings.ToLower(message)

	// Simple pattern matching for common questions
	if strings.Contains(lowerMsg, "hello") || strings.Contains(lowerMsg, "hi") {
		return "Hello! I'm operating in fallback mode. AI features are limited. Try starting Ollama with: ollama serve", nil
	}

	if strings.Contains(lowerMsg, "help") {
		return fmt.Sprintf("You have %d tasks. Use 'focus list' to see them all, or 'focus add' to create new tasks.", len(taskContext)), nil
	}

	if strings.Contains(lowerMsg, "how many") || strings.Contains(lowerMsg, "count") {
		return fmt.Sprintf("You currently have %d tasks.", len(taskContext)), nil
	}

	return fmt.Sprintf("AI is currently unavailable. You have %d tasks. For full AI features, start Ollama with: ollama serve", len(taskContext)), nil
}

// IsAvailable always returns true for fallback
func (p *FallbackProvider) IsAvailable() bool {
	return true
}

// GetName returns the provider name
func (p *FallbackProvider) GetName() string {
	return "Fallback (Rule-based)"
}
