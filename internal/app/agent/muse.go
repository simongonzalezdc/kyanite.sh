package agent

import (
	"fmt"
	"sync"
	"time"

	"github.com/Kyanite/noise/internal/app"
	"github.com/Kyanite/noise/internal/infra/db"
	tea "github.com/charmbracelet/bubbletea"
)

// Muse is the main AI companion agent
type Muse struct {
	memory   *MemoryManager
	observer *ContextObserver
	config   *AgentConfig

	// AI service for generating suggestions
	aiService *app.AIService

	// Current suggestion
	currentSuggestion  *Suggestion
	lastSuggestionTime time.Time

	// State
	isActive bool

	mutex sync.RWMutex
}

// NewMuse creates a new Muse agent
func NewMuse(database *db.DB, aiService *app.AIService, config *AgentConfig) *Muse {
	if config == nil {
		config = DefaultAgentConfig()
	}

	memory := NewMemoryManager(database, config)
	observer := NewContextObserver(memory, config)

	return &Muse{
		memory:    memory,
		observer:  observer,
		config:    config,
		aiService: aiService,
		isActive:  true,
	}
}

// GetMemory returns the memory manager
func (m *Muse) GetMemory() *MemoryManager {
	return m.memory
}

// GetObserver returns the context observer
func (m *Muse) GetObserver() *ContextObserver {
	return m.observer
}

// OnContentChange is called when the editor content changes
func (m *Muse) OnContentChange(content string) {
	if !m.isActive {
		return
	}

	m.observer.RecordContentChange(content)
}

// Tick is called periodically to check for suggestions
func (m *Muse) Tick() tea.Cmd {
	return tea.Tick(5*time.Second, func(t time.Time) tea.Msg {
		return MuseTickMsg{Time: t}
	})
}

// MuseTickMsg is sent periodically to trigger suggestion checks
type MuseTickMsg struct {
	Time time.Time
}

// MuseSuggestionMsg is sent when a suggestion is available
type MuseSuggestionMsg struct {
	Suggestion *Suggestion
}

// MuseDismissSuggestionMsg is sent to dismiss the current suggestion
type MuseDismissSuggestionMsg struct{}

// Update handles Muse-related messages
func (m *Muse) Update(msg tea.Msg) tea.Cmd {
	switch msg.(type) {
	case MuseTickMsg:
		// Check if we should show a suggestion
		m.observer.CheckForPause()

		if m.shouldShowSuggestion() {
			sugg := m.generateSuggestion()
			if sugg != nil {
				return func() tea.Msg {
					return MuseSuggestionMsg{Suggestion: sugg}
				}
			}
		}
		return m.Tick()

	case MuseDismissSuggestionMsg:
		m.dismissSuggestion()
		return nil
	}

	return nil
}

// shouldShowSuggestion checks if it's appropriate to show a suggestion
func (m *Muse) shouldShowSuggestion() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	// Check cooldown
	if time.Since(m.lastSuggestionTime) < m.config.SuggestionCooldown {
		return false
	}

	// Check if there's already an active suggestion
	if m.currentSuggestion != nil && time.Now().Before(m.currentSuggestion.ExpiresAt) {
		return false
	}

	// Check if user is stuck or needs help
	return m.observer.IsUserStuck()
}

// generateSuggestion creates a new suggestion based on current context
func (m *Muse) generateSuggestion() *Suggestion {
	suggType, reason := m.observer.GetSuggestionTrigger()
	if suggType == SuggestionType(-1) {
		return nil
	}

	m.mutex.Lock()
	defer m.mutex.Unlock()

	suggestion := &Suggestion{
		ID:         fmt.Sprintf("sugg-%d", time.Now().UnixNano()),
		Type:       suggType,
		Content:    m.getSuggestionContent(suggType),
		Confidence: 0.8,
		Reason:     reason,
		ExpiresAt:  time.Now().Add(2 * time.Minute),
		Actions: []SuggestionAction{
			{Label: "Accept", Key: "tab", ActionID: "accept"},
			{Label: "Dismiss", Key: "esc", ActionID: "dismiss"},
		},
	}

	m.currentSuggestion = suggestion
	m.lastSuggestionTime = time.Now()

	return suggestion
}

// getSuggestionContent generates content for a suggestion type
func (m *Muse) getSuggestionContent(suggType SuggestionType) string {
	state := m.observer.GetProgressState()

	switch suggType {
	case SuggestNextLine:
		if state == StateStuck {
			return "Try continuing with a new thought or exploring a different direction."
		}
		return "Ready to continue? Let me help with the next line."

	case SuggestBreak:
		return "You've been writing for a while. Consider taking a short break."

	case SuggestStructure:
		return "Consider starting with: [Verse 1] to establish your song structure."

	case SuggestGoal:
		return "What theme or emotion do you want to explore in this song?"

	case SuggestRhyme:
		return "Need a rhyme? Press Ctrl+R to open the rhyme finder."

	case SuggestChord:
		return "Looking for chord ideas? Press Ctrl+K to explore progressions."

	default:
		return "Need help? Press Ctrl+Space to chat with me."
	}
}

// AcceptSuggestion accepts the current suggestion
func (m *Muse) AcceptSuggestion() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.currentSuggestion != nil {
		// Record acceptance
		if m.memory != nil {
			_ = m.memory.RecordEpisode(EpisodicEvent{
				EventType: EventTypeSuggestionAccepted,
				Metadata: map[string]string{
					"suggestion_type": m.currentSuggestion.Type.String(),
					"suggestion_id":   m.currentSuggestion.ID,
				},
			})
		}
		m.currentSuggestion = nil
	}
}

// dismissSuggestion dismisses the current suggestion
func (m *Muse) dismissSuggestion() {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	if m.currentSuggestion != nil {
		// Record dismissal
		if m.memory != nil {
			_ = m.memory.RecordEpisode(EpisodicEvent{
				EventType: EventTypeSuggestionDismissed,
				Metadata: map[string]string{
					"suggestion_type": m.currentSuggestion.Type.String(),
					"suggestion_id":   m.currentSuggestion.ID,
				},
			})
		}
		m.currentSuggestion = nil
	}
}

// GetCurrentSuggestion returns the current active suggestion
func (m *Muse) GetCurrentSuggestion() *Suggestion {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.currentSuggestion
}

// GetProgressState returns the current progress state
func (m *Muse) GetProgressState() ProgressState {
	return m.observer.GetProgressState()
}

// GetSessionStats returns current session statistics
func (m *Muse) GetSessionStats() map[string]interface{} {
	if m.memory == nil {
		return nil
	}

	wm := m.memory.GetWorkingMemory()

	return map[string]interface{}{
		"session_id":       m.memory.GetSessionID(),
		"words_written":    wm.WordsWritten,
		"progress_state":   wm.ProgressState.String(),
		"session_start":    wm.SessionStart,
		"session_duration": time.Since(wm.SessionStart).String(),
		"writing_velocity": m.observer.GetWritingVelocity(),
	}
}

// SetActive enables or disables the Muse agent
func (m *Muse) SetActive(active bool) {
	m.mutex.Lock()
	defer m.mutex.Unlock()
	m.isActive = active
}

// IsActive returns whether the Muse agent is active
func (m *Muse) IsActive() bool {
	m.mutex.RLock()
	defer m.mutex.RUnlock()
	return m.isActive
}

// StartSession starts a new session
func (m *Muse) StartSession() error {
	if m.memory != nil {
		return m.memory.StartSession()
	}
	return nil
}

// EndSession ends the current session
func (m *Muse) EndSession() error {
	if m.memory != nil {
		return m.memory.EndSession()
	}
	return nil
}
