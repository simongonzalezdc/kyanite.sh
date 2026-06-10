package ai

import (
	"fmt"
	"strings"
)

// sanitizeInput prepares user input for safe inclusion in AI prompts.
// It truncates to maxRunes runes, escapes any closing delimiter sequences,
// and wraps the content in delimited sections.
func sanitizeInput(input string) string {
	const maxRunes = 2000
	const openTag = "<user_input>"
	const closeTag = "</user_input>"
	const escapedClose = "&lt;/user_input&gt;"

	runes := []rune(input)
	if len(runes) > maxRunes {
		runes = runes[:maxRunes]
	}
	sanitized := string(runes)
	sanitized = strings.ReplaceAll(sanitized, closeTag, escapedClose)
	return openTag + sanitized + closeTag
}

// PromptBuilder handles construction of AI prompts
type PromptBuilder struct{}

// NewPromptBuilder creates a new prompt builder
func NewPromptBuilder() *PromptBuilder {
	return &PromptBuilder{}
}

// BuildParsePrompt creates a prompt for parsing task descriptions
func (pb *PromptBuilder) BuildParsePrompt(input string) string {
	return fmt.Sprintf(`You are a helpful task management assistant. Parse the following task description and extract:
1. A clear, concise task description
2. Priority level (low, medium, or high)
3. Any deadline mentioned (YYYY-MM-DD format)
4. Categories/tags mentioned
5. Recurrence pattern if mentioned (daily, weekly, monthly, yearly)

Input:
%s

Respond in JSON format:
{
  "description": "clear task description",
  "priority": "medium",
  "deadline": "YYYY-MM-DD",
  "categories": ["category1", "category2"],
  "recurrence_pattern": "daily",
  "recurrence_interval": 1
}

Only include fields that are present in the input.`, sanitizeInput(input))
}

// BuildSuggestPrompt creates a prompt for suggesting new tasks
func (pb *PromptBuilder) BuildSuggestPrompt(existingTasks []string) string {
	tasksStr := "No existing tasks"
	if len(existingTasks) > 0 {
		tasksStr = strings.Join(existingTasks, "\n- ")
	}

	return fmt.Sprintf(`You are a helpful task management assistant. Based on the following existing tasks, suggest 3-5 new tasks that would be complementary or help complete the user's goals.

Existing tasks:
- %s

Suggest new tasks that:
1. Complement existing work
2. Help complete ongoing projects
3. Address gaps in the current task list
4. Are actionable and specific

Respond with a simple list of task descriptions, one per line.`, sanitizeInput(tasksStr))
}

// BuildSummaryPrompt creates a prompt for summarizing tasks
func (pb *PromptBuilder) BuildSummaryPrompt(tasks []string) string {
	tasksStr := strings.Join(tasks, "\n- ")

	return fmt.Sprintf(`You are a helpful task management assistant. Summarize the following tasks into a brief overview highlighting:
1. Main themes or categories
2. Overall progress status
3. High-priority items
4. Suggested next steps

Tasks:
- %s

Provide a concise summary (2-3 paragraphs) that gives the user a clear picture of their task landscape.`, sanitizeInput(tasksStr))
}

// BuildChatPrompt creates a prompt for chat assistance
func (pb *PromptBuilder) BuildChatPrompt(question string, tasks []string) string {
	tasksStr := "No tasks available"
	if len(tasks) > 0 {
		tasksStr = strings.Join(tasks, "\n- ")
	}

	return fmt.Sprintf(`You are a helpful task management assistant. Answer the user's question about their tasks.

User's tasks:
- %s

User's question:
%s

Provide a helpful, actionable response.`, sanitizeInput(tasksStr), sanitizeInput(question))
}
