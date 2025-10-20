package utils

import (
	"time"
	"net/http"
	
	"github.com/kyanite/focus/pkg/styles"
	"github.com/charmbracelet/lipgloss"
)

// AIStatus represents the status of AI services
type AIStatus struct {
	OllamaAvailable bool
	OpenRouterKey  bool
	LastCheck      time.Time
}

var openRouterKey = "" // This would be os.Getenv("OPENROUTER_API_KEY") in real use

// CheckAIStatus checks the availability of AI services
func CheckAIStatus() AIStatus {
	status := AIStatus{
		LastCheck: time.Now(),
		OpenRouterKey: len(openRouterKey) > 0,
	}
	
	// Check Ollama availability
	client := &http.Client{Timeout: 3 * time.Second}
	resp, err := client.Get("http://localhost:11434/api/tags")
	if err == nil && resp.StatusCode == 200 {
		status.OllamaAvailable = true
		resp.Body.Close()
	}
	
	return status
}

// GetStatusIndicator returns a styled status indicator
func GetStatusIndicator(status AIStatus) string {
	if status.OllamaAvailable {
		return lipgloss.NewStyle().
			Foreground(styles.GetSuccess()).
			Bold(true).
			Render("🤖 AI: Online")
	} else if status.OpenRouterKey {
		return lipgloss.NewStyle().
			Foreground(styles.GetWarning()).
			Bold(true).
			Render("🤖 AI: Remote")
	} else {
		return lipgloss.NewStyle().
			Foreground(styles.GetError()).
			Bold(true).
			Render("🤖 AI: Offline")
	}
}

// StreamStatusWithIndicator shows status with streaming effect
func StreamStatusWithIndicator(status AIStatus) {
	statusText := ""
	var color lipgloss.Color
	
	if status.OllamaAvailable {
		statusText = "🤖 AI Systems Online - Local LLM Connected"
		color = styles.GetSuccess()
	} else if status.OpenRouterKey {
		statusText = "⚠️ AI Systems Limited - Remote API Only"
		color = styles.GetWarning()
	} else {
		statusText = "❌ AI Systems Offline - Basic Responses Only"
		color = styles.GetError()
	}
	
	// Stream the status with appropriate color
	StreamWithTypingEffect(statusText, color)
}

// Check OpenRouter API key from environment
func init() {
	// In a real implementation, this would check the actual environment
	// For now, we'll assume no key is set
	openRouterKey = "" // Would be os.Getenv("OPENROUTER_API_KEY")
}
