package agent

import (
	"fmt"
	"sort"
	"strings"
	"sync"

	"github.com/Kyanite/noise/internal/app"
	"github.com/Kyanite/noise/internal/app/ai"
)

// Tool represents an AI-callable tool
type Tool struct {
	Name        string
	Description string
	Parameters  []ToolParameter
	Execute     func(params map[string]string) (string, error)
}

// ToolParameter describes a tool parameter
type ToolParameter struct {
	Name        string
	Description string
	Required    bool
	Default     string
}

// ToolRegistry manages available tools
type ToolRegistry struct {
	tools map[string]*Tool
	mutex sync.RWMutex
}

// NewToolRegistry creates a new tool registry
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]*Tool),
	}
}

// Register adds a tool to the registry
func (r *ToolRegistry) Register(tool *Tool) error {
	if tool == nil {
		return fmt.Errorf("cannot register nil tool")
	}
	r.mutex.Lock()
	defer r.mutex.Unlock()
	r.tools[tool.Name] = tool
	return nil
}

// Get retrieves a tool by name
func (r *ToolRegistry) Get(name string) (*Tool, bool) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	tool, ok := r.tools[name]
	return tool, ok
}

// Execute runs a tool with the given parameters
func (r *ToolRegistry) Execute(name string, params map[string]string) (string, error) {
	tool, ok := r.Get(name)
	if !ok || tool == nil {
		return "", fmt.Errorf("tool not found: %s", name)
	}
	
	// Fill in default values
	if params == nil {
		params = make(map[string]string)
	}
	for _, param := range tool.Parameters {
		if _, ok := params[param.Name]; !ok && param.Default != "" {
			params[param.Name] = param.Default
		}
	}
	
	// Check required parameters
	for _, param := range tool.Parameters {
		if param.Required {
			if _, ok := params[param.Name]; !ok {
				return "", fmt.Errorf("missing required parameter: %s", param.Name)
			}
		}
	}
	
	return tool.Execute(params)
}

// List returns all available tools
func (r *ToolRegistry) List() []*Tool {
	r.mutex.RLock()
	defer r.mutex.RUnlock()
	
	tools := make([]*Tool, 0, len(r.tools))
	for _, tool := range r.tools {
		tools = append(tools, tool)
	}
	
	// Sort by name
	sort.Slice(tools, func(i, j int) bool {
		return tools[i].Name < tools[j].Name
	})
	
	return tools
}

// GetToolDescriptions returns descriptions of all tools for AI context
func (r *ToolRegistry) GetToolDescriptions() string {
	tools := r.List()
	var sb strings.Builder
	
	sb.WriteString("Available tools:\n")
	for _, tool := range tools {
		sb.WriteString(fmt.Sprintf("- %s: %s\n", tool.Name, tool.Description))
		for _, param := range tool.Parameters {
			req := ""
			if param.Required {
				req = " (required)"
			}
			sb.WriteString(fmt.Sprintf("  - %s: %s%s\n", param.Name, param.Description, req))
		}
	}
	
	return sb.String()
}

// registerDefaultTools registers the default set of tools
func registerDefaultTools(r *ToolRegistry, memory *MemoryManager, aiService *app.AIService) {
	// Rhyme Finder Tool
	r.Register(&Tool{
		Name:        "rhyme_finder",
		Description: "Finds rhymes for a given word",
		Parameters: []ToolParameter{
			{Name: "word", Description: "The word to find rhymes for", Required: true},
			{Name: "type", Description: "Type of rhyme: perfect, near, slant", Default: "all"},
		},
		Execute: func(params map[string]string) (string, error) {
			word := params["word"]
			if word == "" {
				return "", fmt.Errorf("word parameter is required")
			}
			
			// Use AI rhyme service if available
			if aiService != nil {
				rhymes := aiService.FindRhymes(word)
				if len(rhymes) > 0 {
					return formatRhymes(word, rhymes), nil
				}
			}
			
			// Fallback to basic rhyme matching
			rhymes := ai.FindBasicRhymes(word)
			return formatRhymes(word, rhymes), nil
		},
	})
	
	// Lyrics Analyzer Tool
	r.Register(&Tool{
		Name:        "lyrics_analyzer",
		Description: "Analyzes the current lyrics for patterns, rhyme scheme, and suggestions",
		Parameters:  []ToolParameter{},
		Execute: func(params map[string]string) (string, error) {
			if memory == nil {
				return "No lyrics loaded to analyze.", nil
			}
			
			wm := memory.GetWorkingMemory()
			if wm.CurrentSong == nil || wm.CurrentSong.RawContent == "" {
				return "No lyrics loaded to analyze. Open or create a song first.", nil
			}
			
			content := wm.CurrentSong.RawContent
			return analyzeContent(content), nil
		},
	})
	
	// Search Songs Tool
	r.Register(&Tool{
		Name:        "search_songs",
		Description: "Searches through your songs",
		Parameters: []ToolParameter{
			{Name: "query", Description: "Search query", Required: true},
		},
		Execute: func(params map[string]string) (string, error) {
			query := params["query"]
			if query == "" {
				return "", fmt.Errorf("query parameter is required")
			}
			
			// This would search the database - for now return a helpful message
			return fmt.Sprintf("Searching for '%s' in your songs... (Search feature coming soon)", query), nil
		},
	})
	
	// Version History Tool
	r.Register(&Tool{
		Name:        "version_history",
		Description: "Shows the version history for the current song",
		Parameters:  []ToolParameter{},
		Execute: func(params map[string]string) (string, error) {
			if memory == nil {
				return "Version history not available.", nil
			}
			
			wm := memory.GetWorkingMemory()
			if wm.CurrentSong == nil {
				return "No song is currently open.", nil
			}
			
			// This would fetch version history from database
			return fmt.Sprintf("Version history for '%s' (coming soon)", wm.CurrentSong.Metadata.Title), nil
		},
	})
	
	// Session Stats Tool
	r.Register(&Tool{
		Name:        "session_stats",
		Description: "Shows statistics for the current writing session",
		Parameters:  []ToolParameter{},
		Execute: func(params map[string]string) (string, error) {
			if memory == nil {
				return "Session stats not available.", nil
			}
			
			stats, err := memory.GetMemoryStats()
			if err != nil {
				return fmt.Sprintf("Error getting stats: %v", err), nil
			}
			
			return formatSessionStats(stats), nil
		},
	})
	
	// Chord Suggestion Tool
	r.Register(&Tool{
		Name:        "chord_suggest",
		Description: "Suggests chord progressions based on mood or style",
		Parameters: []ToolParameter{
			{Name: "mood", Description: "The mood (happy, sad, energetic, calm)", Default: "neutral"},
			{Name: "key", Description: "The musical key (C, Am, G, etc.)", Default: "C"},
		},
		Execute: func(params map[string]string) (string, error) {
			mood := params["mood"]
			key := params["key"]
			return suggestChords(mood, key), nil
		},
	})
	
	// Structure Helper Tool
	r.Register(&Tool{
		Name:        "structure_help",
		Description: "Provides song structure suggestions and templates",
		Parameters: []ToolParameter{
			{Name: "style", Description: "Song style (pop, rock, ballad, folk)", Default: "pop"},
		},
		Execute: func(params map[string]string) (string, error) {
			style := params["style"]
			return getStructureSuggestion(style), nil
		},
	})
}

// Helper functions for tools

func formatRhymes(word string, rhymes []string) string {
	if len(rhymes) == 0 {
		return fmt.Sprintf("No rhymes found for '%s'. Try a different word or check the spelling.", word)
	}
	
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Rhymes for '%s':\n", word))
	
	// Group by first letter or show all
	if len(rhymes) <= 20 {
		for _, rhyme := range rhymes {
			sb.WriteString(fmt.Sprintf("  - %s\n", rhyme))
		}
	} else {
		// Show first 20 with indication there are more
		for i, rhyme := range rhymes[:20] {
			if i > 0 && i%5 == 0 {
				sb.WriteString("\n")
			}
			sb.WriteString(fmt.Sprintf("  %s", rhyme))
			if i < 19 {
				sb.WriteString(",")
			}
		}
		sb.WriteString(fmt.Sprintf("\n  ... and %d more\n", len(rhymes)-20))
	}
	
	return sb.String()
}

func analyzeContent(content string) string {
	lines := strings.Split(content, "\n")
	nonEmptyLines := 0
	wordCount := 0
	sections := 0
	
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		nonEmptyLines++
		wordCount += len(strings.Fields(line))
		
		// Check for section markers
		if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
			sections++
		}
	}
	
	var sb strings.Builder
	sb.WriteString("Lyrics Analysis:\n\n")
	sb.WriteString(fmt.Sprintf("Total lines: %d\n", nonEmptyLines))
	sb.WriteString(fmt.Sprintf("Word count: %d\n", wordCount))
	sb.WriteString(fmt.Sprintf("Sections found: %d\n", sections))
	
	if sections == 0 {
		sb.WriteString("\nTip: Consider adding section markers like [Verse 1], [Chorus], [Bridge]")
	}
	
	if wordCount < 50 {
		sb.WriteString("\nTip: Your lyrics are quite short. Consider expanding with more verses or a bridge.")
	} else if wordCount > 400 {
		sb.WriteString("\nTip: Your lyrics are quite long. Make sure each section serves the song's message.")
	}
	
	return sb.String()
}

func formatSessionStats(stats map[string]interface{}) string {
	var sb strings.Builder
	sb.WriteString("Session Statistics:\n\n")
	
	if current, ok := stats["current_session"].(map[string]interface{}); ok {
		if ww, ok := current["words_written"].(int); ok {
			sb.WriteString(fmt.Sprintf("Words written: %d\n", ww))
		}
		if ps, ok := current["progress_state"].(string); ok {
			sb.WriteString(fmt.Sprintf("Progress state: %s\n", ps))
		}
	}
	
	if ec, ok := stats["episode_count"].(int); ok {
		sb.WriteString(fmt.Sprintf("Episodes recorded: %d\n", ec))
	}
	if cc, ok := stats["conversation_count"].(int); ok {
		sb.WriteString(fmt.Sprintf("Conversations: %d\n", cc))
	}
	
	return sb.String()
}

func suggestChords(mood, key string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Chord suggestions in %s (%s mood):\n\n", key, mood))
	
	// Basic chord progressions based on mood
	progressions := map[string][]string{
		"happy":    {"I - IV - V - I", "I - V - vi - IV", "I - IV - I - V"},
		"sad":      {"vi - IV - I - V", "i - VI - III - VII", "i - iv - v - i"},
		"energetic": {"I - IV - V - IV", "I - III - IV - IV", "I - bVII - IV - I"},
		"calm":     {"I - vi - IV - V", "I - IV - vi - V", "Imaj7 - IVmaj7 - V7"},
		"neutral":  {"I - IV - V - I", "I - V - vi - IV", "I - vi - IV - V"},
	}
	
	progs, ok := progressions[strings.ToLower(mood)]
	if !ok {
		progs = progressions["neutral"]
	}
	
	for i, prog := range progs {
		sb.WriteString(fmt.Sprintf("%d. %s\n", i+1, prog))
	}
	
	sb.WriteString(fmt.Sprintf("\nIn the key of %s:\n", key))
	// This would translate to actual chords based on key
	sb.WriteString("(Chord translation coming soon)")
	
	return sb.String()
}

func getStructureSuggestion(style string) string {
	structures := map[string]string{
		"pop": `Pop Song Structure:

[Intro] - 4-8 bars
[Verse 1] - Tell the story, set the scene
[Pre-Chorus] - Build tension (optional)
[Chorus] - The hook, most memorable part
[Verse 2] - Develop the story
[Chorus]
[Bridge] - New perspective or key change
[Chorus] - Final, often with variation
[Outro] - 4-8 bars

Tip: Keep the chorus simple and singable!`,

		"rock": `Rock Song Structure:

[Intro] - Guitar riff or drums
[Verse 1] - Set up the theme
[Chorus] - High energy payoff
[Verse 2] - Continue the narrative
[Chorus]
[Solo] - Instrumental break
[Bridge] - Breakdown or key change
[Chorus] - Big finish
[Outro] - Fade or hard stop

Tip: Build energy toward the chorus!`,

		"ballad": `Ballad Structure:

[Intro] - Soft, sets mood
[Verse 1] - Intimate, detailed
[Verse 2] - Develop emotion
[Chorus] - Emotional peak
[Verse 3] - Deeper revelation
[Chorus] - With more intensity
[Outro] - Quiet resolution

Tip: Focus on emotional authenticity.`,

		"folk": `Folk Song Structure:

[Verse 1] - Introduce the story
[Verse 2] - Develop characters
[Chorus] - The moral or refrain
[Verse 3] - Conflict or change
[Chorus]
[Verse 4] - Resolution
[Chorus] - Final reflection

Tip: Tell a story with a clear beginning, middle, and end.`,
	}
	
	structure, ok := structures[strings.ToLower(style)]
	if !ok {
		structure = structures["pop"]
	}
	
	return structure
}
