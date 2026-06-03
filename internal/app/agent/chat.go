package agent

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/Kyanite/noise/internal/app"
)

// ChatHandler manages conversations with the user
type ChatHandler struct {
	memory    *MemoryManager
	aiService *app.AIService
	tools     *ToolRegistry
	config    *AgentConfig

	// Conversation state
	isActive  bool
	lastQuery string

	mutex sync.RWMutex
}

// NewChatHandler creates a new chat handler
func NewChatHandler(memory *MemoryManager, aiService *app.AIService, config *AgentConfig) *ChatHandler {
	if config == nil {
		config = DefaultAgentConfig()
	}

	tools := NewToolRegistry()
	registerDefaultTools(tools, memory, aiService)

	return &ChatHandler{
		memory:    memory,
		aiService: aiService,
		tools:     tools,
		config:    config,
		isActive:  false,
	}
}

// Chat sends a message to the AI and returns a response
func (c *ChatHandler) Chat(ctx context.Context, userMessage string) (string, error) {
	c.mutex.Lock()
	c.isActive = true
	c.lastQuery = userMessage
	c.mutex.Unlock()

	defer func() {
		c.mutex.Lock()
		c.isActive = false
		c.mutex.Unlock()
	}()

	// Record user message
	if c.memory != nil {
		_ = c.memory.RecordChatMessage(ChatMessage{
			Role:    "user",
			Content: userMessage,
			Context: c.getCurrentContext(),
		})
	}

	// Build context from memory
	contextStr := c.buildContext()

	// Check if this is a tool request
	toolResponse, toolUsed := c.checkToolRequest(userMessage)
	if toolUsed {
		// Record assistant response with tool call
		if c.memory != nil {
			_ = c.memory.RecordChatMessage(ChatMessage{
				Role:    "assistant",
				Content: toolResponse,
				ToolCalls: []ToolCall{{
					Name:   extractToolName(userMessage),
					Result: toolResponse,
				}},
			})
		}
		return toolResponse, nil
	}

	// Generate AI response
	response, err := c.generateResponse(ctx, userMessage, contextStr)
	if err != nil {
		return "", err
	}

	// Record assistant response
	if c.memory != nil {
		_ = c.memory.RecordChatMessage(ChatMessage{
			Role:    "assistant",
			Content: response,
		})
	}

	return response, nil
}

// getCurrentContext returns the current context as a map
func (c *ChatHandler) getCurrentContext() map[string]string {
	ctx := make(map[string]string)

	if c.memory != nil {
		wm := c.memory.GetWorkingMemory()
		if wm.CurrentSong != nil {
			ctx["song_title"] = wm.CurrentSong.Metadata.Title
			ctx["song_id"] = fmt.Sprintf("%d", wm.CurrentSong.ID)
		}
		ctx["section"] = wm.CurrentSection
		ctx["progress_state"] = wm.ProgressState.String()
		ctx["words_written"] = fmt.Sprintf("%d", wm.WordsWritten)
	}

	return ctx
}

// buildContext builds a context string from memory
func (c *ChatHandler) buildContext() string {
	var parts []string

	if c.memory != nil {
		wm := c.memory.GetWorkingMemory()

		// Add current song info
		if wm.CurrentSong != nil {
			parts = append(parts, fmt.Sprintf("Currently working on: %s", wm.CurrentSong.Metadata.Title))
			if wm.CurrentSection != "" {
				parts = append(parts, fmt.Sprintf("Current section: %s", wm.CurrentSection))
			}
		}

		// Add progress info
		parts = append(parts, fmt.Sprintf("Progress state: %s", wm.ProgressState.String()))
		parts = append(parts, fmt.Sprintf("Words written this session: %d", wm.WordsWritten))

		// Add recent chat history
		history, err := c.memory.GetChatHistory(c.config.ContextWindowSize)
		if err == nil && len(history) > 0 {
			parts = append(parts, "\nRecent conversation:")
			for _, msg := range history {
				if len(msg.Content) > 200 {
					msg.Content = msg.Content[:200] + "..."
				}
				parts = append(parts, fmt.Sprintf("%s: %s", msg.Role, msg.Content))
			}
		}
	}

	return strings.Join(parts, "\n")
}

// checkToolRequest checks if the message is a tool request and executes it
func (c *ChatHandler) checkToolRequest(message string) (string, bool) {
	message = strings.ToLower(message)

	// Check for tool triggers
	if strings.Contains(message, "rhyme") || strings.Contains(message, "rhymes with") {
		// Extract the word to rhyme
		word := extractWord(message, "rhyme")
		if word != "" {
			result, err := c.tools.Execute("rhyme_finder", map[string]string{"word": word})
			if err == nil {
				return result, true
			}
		}
	}

	if strings.Contains(message, "analyze") || strings.Contains(message, "analysis") {
		result, err := c.tools.Execute("lyrics_analyzer", nil)
		if err == nil {
			return result, true
		}
	}

	if strings.Contains(message, "search") || strings.Contains(message, "find") {
		query := extractQuery(message)
		if query != "" {
			result, err := c.tools.Execute("search_songs", map[string]string{"query": query})
			if err == nil {
				return result, true
			}
		}
	}

	if strings.Contains(message, "history") || strings.Contains(message, "versions") {
		result, err := c.tools.Execute("version_history", nil)
		if err == nil {
			return result, true
		}
	}

	return "", false
}

// generateResponse generates an AI response
func (c *ChatHandler) generateResponse(ctx context.Context, userMessage, contextStr string) (string, error) {
	// Mark ctx and contextStr as used for future AI integration
	_ = ctx
	_ = contextStr

	if c.aiService == nil {
		return c.generateLocalResponse(userMessage)
	}

	// TODO: Use AI service to generate response when external AI is available
	// For now, return a helpful local response
	return c.generateLocalResponse(userMessage)
}

// generateLocalResponse generates a response without external AI
func (c *ChatHandler) generateLocalResponse(message string) (string, error) {
	message = strings.ToLower(message)

	// Pattern matching for common questions
	if strings.Contains(message, "help") || strings.Contains(message, "what can you do") {
		return `I'm Muse, your AI songwriting companion! I can help you with:

- **Rhymes**: Ask "what rhymes with love?" 
- **Analysis**: Say "analyze my lyrics" to get feedback
- **Structure**: I can suggest song structures and sections
- **Ideas**: Tell me a theme and I'll help brainstorm
- **History**: Say "show history" to see your versions

What would you like help with?`, nil
	}

	if strings.Contains(message, "stuck") || strings.Contains(message, "block") {
		return `Writer's block is totally normal! Here are some ideas:

1. **Change perspective**: Write from a different character's view
2. **Sense details**: What do you see, hear, smell in the scene?
3. **Freewrite**: Just write anything for 5 minutes, no editing
4. **Skip ahead**: Write a different section and come back
5. **Take a break**: Sometimes stepping away helps

Which approach sounds interesting to try?`, nil
	}

	if strings.Contains(message, "verse") || strings.Contains(message, "chorus") || strings.Contains(message, "bridge") {
		return `Here are some tips for that section:

**Verse**: Set up the story, use specific details
**Chorus**: The emotional core, keep it memorable and singable  
**Bridge**: Shift perspective or add new insight

What's the theme or emotion you're exploring?`, nil
	}

	if strings.Contains(message, "thank") {
		return "You're welcome! Keep writing, you're doing great. Let me know if you need anything else!", nil
	}

	// Default response
	return "I'm here to help with your songwriting! You can ask me about rhymes, song structure, or just tell me what you're working on.", nil
}

// GetHistory returns the chat history
func (c *ChatHandler) GetHistory(limit int) ([]ChatMessage, error) {
	if c.memory == nil {
		return nil, fmt.Errorf("memory not available")
	}
	return c.memory.GetChatHistory(limit)
}

// ClearHistory clears the chat history
func (c *ChatHandler) ClearHistory() error {
	if c.memory == nil {
		return fmt.Errorf("memory not available")
	}
	return c.memory.ClearMemory(false, false, true)
}

// IsActive returns whether a chat is currently in progress
func (c *ChatHandler) IsActive() bool {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.isActive
}

// GetTools returns the tool registry
func (c *ChatHandler) GetTools() *ToolRegistry {
	return c.tools
}

// Helper functions

func extractToolName(message string) string {
	message = strings.ToLower(message)
	if strings.Contains(message, "rhyme") {
		return "rhyme_finder"
	}
	if strings.Contains(message, "analyze") {
		return "lyrics_analyzer"
	}
	if strings.Contains(message, "search") {
		return "search_songs"
	}
	if strings.Contains(message, "history") {
		return "version_history"
	}
	return "unknown"
}

func extractWord(message, trigger string) string {
	message = strings.ToLower(message)

	// Pattern: "rhymes with X" or "rhyme for X"
	patterns := []string{
		trigger + "s with ",
		trigger + " for ",
		trigger + " ",
	}

	for _, pattern := range patterns {
		if idx := strings.Index(message, pattern); idx != -1 {
			rest := strings.TrimSpace(message[idx+len(pattern):])
			// Get first word
			if words := strings.Fields(rest); len(words) > 0 {
				return strings.Trim(words[0], "?!.,")
			}
		}
	}

	return ""
}

func extractQuery(message string) string {
	message = strings.ToLower(message)

	patterns := []string{
		"search for ",
		"find ",
		"search ",
	}

	for _, pattern := range patterns {
		if idx := strings.Index(message, pattern); idx != -1 {
			rest := strings.TrimSpace(message[idx+len(pattern):])
			return strings.Trim(rest, "?!.,")
		}
	}

	return ""
}
