package agent

import (
	"strings"
	"sync"
	"time"

	"github.com/kyanite/noise/internal/app/ai"
)

// ContextObserver monitors user activity and detects writing context
type ContextObserver struct {
	memory          *MemoryManager
	contentDetector *ai.ContextDetector
	config          *AgentConfig

	// Tracking state
	lastContent     string
	lastContentTime time.Time
	editHistory     []EditEvent
	pauseStart      time.Time
	isPaused        bool

	// Computed metrics
	writingVelocity float64 // words per minute
	deletionRatio   float64 // deletions / total edits

	mutex sync.RWMutex
}

// NewContextObserver creates a new context observer
func NewContextObserver(memory *MemoryManager, config *AgentConfig) *ContextObserver {
	if config == nil {
		config = DefaultAgentConfig()
	}

	return &ContextObserver{
		memory:          memory,
		contentDetector: ai.NewContextDetector(),
		config:          config,
		editHistory:     make([]EditEvent, 0, 100),
		lastContentTime: time.Now(),
	}
}

// AnalyzeContent returns the detected content type (lyrics, patterns, mixed)
func (o *ContextObserver) AnalyzeContent(content string) ai.ContentType {
	if o.contentDetector == nil {
		return ai.ContentTypeUnknown
	}
	return o.contentDetector.AnalyzeContent(content)
}

// RecordContentChange records a content change and updates metrics
func (o *ContextObserver) RecordContentChange(newContent string) {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	now := time.Now()

	// Calculate edit delta
	oldWords := countWords(o.lastContent)
	newWords := countWords(newContent)
	wordsDelta := newWords - oldWords

	// Determine edit type
	editType := "replace"
	if wordsDelta > 0 {
		editType = "insert"
	} else if wordsDelta < 0 {
		editType = "delete"
	}

	// Create edit event
	edit := EditEvent{
		Timestamp:  now,
		Type:       editType,
		WordsDelta: wordsDelta,
		Content:    extractDiff(o.lastContent, newContent),
	}

	// Add to history
	o.editHistory = append(o.editHistory, edit)
	if len(o.editHistory) > 100 {
		o.editHistory = o.editHistory[len(o.editHistory)-100:]
	}

	// Update metrics
	o.updateMetrics()

	// Record to memory manager
	if o.memory != nil {
		o.memory.RecordEdit(edit)
	}

	// Update last content
	o.lastContent = newContent
	o.lastContentTime = now

	// Reset pause state since there was activity
	o.isPaused = false
}

// updateMetrics recalculates writing velocity and deletion ratio
func (o *ContextObserver) updateMetrics() {
	// Must be called with lock held

	if len(o.editHistory) < 2 {
		o.writingVelocity = 0
		o.deletionRatio = 0
		return
	}

	// Calculate metrics over last 5 minutes of edits
	cutoff := time.Now().Add(-5 * time.Minute)
	var recentEdits []EditEvent
	for _, edit := range o.editHistory {
		if edit.Timestamp.After(cutoff) {
			recentEdits = append(recentEdits, edit)
		}
	}

	if len(recentEdits) < 2 {
		return
	}

	// Calculate writing velocity (words per minute)
	firstEdit := recentEdits[0]
	lastEdit := recentEdits[len(recentEdits)-1]
	duration := lastEdit.Timestamp.Sub(firstEdit.Timestamp).Minutes()

	if duration > 0 {
		var totalWordsAdded int
		for _, edit := range recentEdits {
			if edit.WordsDelta > 0 {
				totalWordsAdded += edit.WordsDelta
			}
		}
		o.writingVelocity = float64(totalWordsAdded) / duration
	} else {
		o.writingVelocity = 0
	}

	// Calculate deletion ratio
	var insertCount, deleteCount int
	for _, edit := range recentEdits {
		if edit.Type == "insert" {
			insertCount++
		} else if edit.Type == "delete" {
			deleteCount++
		}
	}

	totalEdits := insertCount + deleteCount
	if totalEdits > 0 {
		o.deletionRatio = float64(deleteCount) / float64(totalEdits)
	}
}

// GetWritingVelocity returns the current writing velocity (words per minute)
func (o *ContextObserver) GetWritingVelocity() float64 {
	o.mutex.RLock()
	defer o.mutex.RUnlock()
	return o.writingVelocity
}

// GetDeletionRatio returns the ratio of deletions to total edits
func (o *ContextObserver) GetDeletionRatio() float64 {
	o.mutex.RLock()
	defer o.mutex.RUnlock()
	return o.deletionRatio
}

// GetPauseDuration returns how long the user has been paused
func (o *ContextObserver) GetPauseDuration() time.Duration {
	o.mutex.RLock()
	defer o.mutex.RUnlock()

	if o.isPaused {
		return time.Since(o.pauseStart)
	}
	return time.Since(o.lastContentTime)
}

// CheckForPause checks if the user has paused and updates state
func (o *ContextObserver) CheckForPause() {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	pauseDuration := time.Since(o.lastContentTime)

	if !o.isPaused && pauseDuration > 30*time.Second {
		o.isPaused = true
		o.pauseStart = o.lastContentTime
	}
}

// IsUserStuck returns true if the user appears to be stuck
func (o *ContextObserver) IsUserStuck() bool {
	o.mutex.RLock()
	defer o.mutex.RUnlock()

	// Check if paused too long
	pauseDuration := time.Since(o.lastContentTime)
	if pauseDuration > o.config.StuckThreshold {
		return true
	}

	// Check if deletion ratio is high (lots of deletions = struggling)
	if o.deletionRatio > 0.7 && len(o.editHistory) > 10 {
		return true
	}

	return false
}

// GetProgressState returns the current progress state
func (o *ContextObserver) GetProgressState() ProgressState {
	o.mutex.RLock()
	defer o.mutex.RUnlock()

	pauseDuration := time.Since(o.lastContentTime)

	// Check for stuck state
	if pauseDuration > o.config.StuckThreshold {
		return StateStuck
	}

	// Check for reviewing state (reading, no edits)
	if pauseDuration > 30*time.Second && pauseDuration < o.config.StuckThreshold {
		return StateReviewing
	}

	// Check edit patterns
	if len(o.editHistory) < 5 {
		return StateStarting
	}

	// High deletion ratio = refining
	if o.deletionRatio > 0.5 {
		return StateRefining
	}

	// Good velocity = flowing
	if o.writingVelocity > 5 { // 5+ words per minute
		return StateFlowing
	}

	return StateStarting
}

// GetStuckReason returns a reason why the user might be stuck
func (o *ContextObserver) GetStuckReason() string {
	o.mutex.RLock()
	defer o.mutex.RUnlock()

	pauseDuration := time.Since(o.lastContentTime)

	if pauseDuration > o.config.StuckThreshold {
		return "long_pause"
	}

	if o.deletionRatio > 0.7 {
		return "many_deletions"
	}

	if o.writingVelocity < 1 && len(o.editHistory) > 10 {
		return "slow_progress"
	}

	return ""
}

// GetRecentEditSummary returns a summary of recent editing activity
func (o *ContextObserver) GetRecentEditSummary() map[string]interface{} {
	o.mutex.RLock()
	defer o.mutex.RUnlock()

	summary := map[string]interface{}{
		"edit_count":       len(o.editHistory),
		"writing_velocity": o.writingVelocity,
		"deletion_ratio":   o.deletionRatio,
		"is_paused":        o.isPaused,
		"pause_duration":   time.Since(o.lastContentTime).String(),
		"progress_state":   o.GetProgressState().String(),
	}

	// Calculate words written in recent edits
	var wordsAdded, wordsDeleted int
	for _, edit := range o.editHistory {
		if edit.WordsDelta > 0 {
			wordsAdded += edit.WordsDelta
		} else {
			wordsDeleted += -edit.WordsDelta
		}
	}
	summary["words_added"] = wordsAdded
	summary["words_deleted"] = wordsDeleted
	summary["net_words"] = wordsAdded - wordsDeleted

	return summary
}

// ShouldSuggest returns true if it's appropriate to show a suggestion
func (o *ContextObserver) ShouldSuggest(suggestionType SuggestionType) bool {
	o.mutex.RLock()
	defer o.mutex.RUnlock()

	state := o.GetProgressState()
	pauseDuration := time.Since(o.lastContentTime)

	switch suggestionType {
	case SuggestNextLine:
		// Suggest next line when stuck or paused
		return state == StateStuck || (state == StateReviewing && pauseDuration > time.Minute)

	case SuggestBreak:
		// Suggest break after long session
		if o.memory != nil {
			wm := o.memory.GetWorkingMemory()
			sessionDuration := time.Since(wm.SessionStart)
			return sessionDuration > 2*time.Hour
		}
		return false

	case SuggestRhyme:
		// Suggest rhymes when writing lyrics
		return true // Always available for lyrics

	case SuggestChord:
		// Suggest chords when velocity is good
		return state == StateFlowing || state == StateRefining

	case SuggestStructure:
		// Suggest structure when starting or stuck
		return state == StateStarting || state == StateStuck

	case SuggestGoal:
		// Suggest goals at session start
		return state == StateStarting && len(o.editHistory) < 5

	default:
		return false
	}
}

// GetSuggestionTrigger returns the reason a suggestion should be triggered
func (o *ContextObserver) GetSuggestionTrigger() (SuggestionType, string) {
	state := o.GetProgressState()

	switch state {
	case StateStuck:
		reason := o.GetStuckReason()
		switch reason {
		case "long_pause":
			return SuggestNextLine, "You've been paused for a while"
		case "many_deletions":
			return SuggestNextLine, "Struggling with this section?"
		default:
			return SuggestNextLine, "Need some inspiration?"
		}

	case StateStarting:
		if o.memory != nil {
			wm := o.memory.GetWorkingMemory()
			if wm.CurrentSong == nil || wm.WordsWritten == 0 {
				return SuggestGoal, "What do you want to write about today?"
			}
		}
		return SuggestStructure, "How about starting with a verse?"

	case StateReviewing:
		return SuggestNextLine, "Ready to continue?"

	default:
		return SuggestionType(-1), ""
	}
}

// Reset clears the observer state
func (o *ContextObserver) Reset() {
	o.mutex.Lock()
	defer o.mutex.Unlock()

	o.lastContent = ""
	o.lastContentTime = time.Now()
	o.editHistory = make([]EditEvent, 0, 100)
	o.isPaused = false
	o.writingVelocity = 0
	o.deletionRatio = 0
}

// countWords counts the number of words in a string
func countWords(s string) int {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0
	}
	return len(strings.Fields(s))
}

// extractDiff extracts a simple diff between old and new content
func extractDiff(old, new string) string {
	// Simple implementation: return the part that changed
	if len(new) > len(old) {
		// Addition
		return new[len(old):]
	} else if len(new) < len(old) {
		// Deletion
		return old[len(new):]
	}
	return ""
}
