package agent

import (
	"os"
	"testing"
	"time"

	"github.com/Kyanite/noise/internal/infra/db"
)

func setupTestDB(t *testing.T) *db.DB {
	// Create a temporary directory for the test database
	tempDir, err := os.MkdirTemp("", "muse-test-*")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	t.Cleanup(func() {
		os.RemoveAll(tempDir)
	})

	database, err := db.New(db.Config{DataDir: tempDir})
	if err != nil {
		t.Fatalf("Failed to create test database: %v", err)
	}
	t.Cleanup(func() {
		database.Close()
	})

	return database
}

func TestMemoryManagerCreation(t *testing.T) {
	database := setupTestDB(t)
	
	mm := NewMemoryManager(database, nil)
	if mm == nil {
		t.Fatal("NewMemoryManager returned nil")
	}
	
	if mm.GetSessionID() == "" {
		t.Error("Session ID should not be empty")
	}
	
	wm := mm.GetWorkingMemory()
	if wm == nil {
		t.Fatal("Working memory should not be nil")
	}
	
	if wm.ProgressState != StateStarting {
		t.Errorf("Initial progress state should be StateStarting, got %v", wm.ProgressState)
	}
}

func TestMemoryManagerEditTracking(t *testing.T) {
	database := setupTestDB(t)
	mm := NewMemoryManager(database, nil)
	
	// Record some edits
	mm.RecordEdit(EditEvent{
		Type:       "insert",
		Position:   0,
		Length:     10,
		Content:    "Hello World",
		WordsDelta: 2,
	})
	
	wm := mm.GetWorkingMemory()
	if len(wm.RecentEdits) != 1 {
		t.Errorf("Expected 1 edit, got %d", len(wm.RecentEdits))
	}
	
	if wm.WordsWritten != 2 {
		t.Errorf("Expected 2 words written, got %d", wm.WordsWritten)
	}
	
	// Record more edits to test flowing state
	for i := 0; i < 10; i++ {
		mm.RecordEdit(EditEvent{
			Type:       "insert",
			Position:   i * 10,
			Length:     5,
			Content:    "test ",
			WordsDelta: 5,
		})
	}
	
	wm = mm.GetWorkingMemory()
	if wm.ProgressState != StateFlowing {
		t.Errorf("Expected StateFlowing after many inserts, got %v", wm.ProgressState)
	}
}

func TestMemoryManagerEpisodicMemory(t *testing.T) {
	database := setupTestDB(t)
	mm := NewMemoryManager(database, nil)
	
	// Record an episode
	err := mm.RecordEpisode(EpisodicEvent{
		EventType:      EventTypeBrainstorm,
		Section:        "verse1",
		ContentSnippet: "Testing the brainstorm feature",
		Outcome:        "accepted",
		Metadata: map[string]string{
			"theme": "love",
		},
	})
	if err != nil {
		t.Fatalf("Failed to record episode: %v", err)
	}
	
	// Retrieve episodes
	events, err := mm.GetRecentEpisodes(10)
	if err != nil {
		t.Fatalf("Failed to get recent episodes: %v", err)
	}
	
	if len(events) != 1 {
		t.Fatalf("Expected 1 episode, got %d", len(events))
	}
	
	if events[0].EventType != EventTypeBrainstorm {
		t.Errorf("Expected EventTypeBrainstorm, got %v", events[0].EventType)
	}
	
	if events[0].Metadata["theme"] != "love" {
		t.Errorf("Expected metadata theme 'love', got '%v'", events[0].Metadata["theme"])
	}
}

func TestMemoryManagerPreferences(t *testing.T) {
	database := setupTestDB(t)
	mm := NewMemoryManager(database, nil)
	
	// Set a preference
	err := mm.SetPreference("preferred_genre", "rock", 0.8, "explicit")
	if err != nil {
		t.Fatalf("Failed to set preference: %v", err)
	}
	
	// Get the preference
	pref, err := mm.GetPreference("preferred_genre")
	if err != nil {
		t.Fatalf("Failed to get preference: %v", err)
	}
	
	if pref == nil {
		t.Fatal("Preference should not be nil")
	}
	
	if pref.Key != "preferred_genre" {
		t.Errorf("Expected key 'preferred_genre', got '%s'", pref.Key)
	}
	
	if pref.Confidence != 0.8 {
		t.Errorf("Expected confidence 0.8, got %f", pref.Confidence)
	}
	
	if pref.Source != "explicit" {
		t.Errorf("Expected source 'explicit', got '%s'", pref.Source)
	}
	
	// Get non-existent preference
	pref2, err := mm.GetPreference("nonexistent")
	if err != nil {
		t.Fatalf("Error getting nonexistent preference: %v", err)
	}
	if pref2 != nil {
		t.Error("Expected nil for nonexistent preference")
	}
}

func TestMemoryManagerChatHistory(t *testing.T) {
	database := setupTestDB(t)
	mm := NewMemoryManager(database, nil)
	
	// Record chat messages
	err := mm.RecordChatMessage(ChatMessage{
		Role:    "user",
		Content: "Help me with this verse",
		Context: map[string]string{
			"section": "verse1",
		},
	})
	if err != nil {
		t.Fatalf("Failed to record user message: %v", err)
	}
	
	err = mm.RecordChatMessage(ChatMessage{
		Role:    "assistant",
		Content: "Here are some suggestions...",
	})
	if err != nil {
		t.Fatalf("Failed to record assistant message: %v", err)
	}
	
	// Get chat history
	history, err := mm.GetChatHistory(10)
	if err != nil {
		t.Fatalf("Failed to get chat history: %v", err)
	}
	
	if len(history) != 2 {
		t.Fatalf("Expected 2 messages, got %d", len(history))
	}
	
	// Should be in chronological order
	if history[0].Role != "user" {
		t.Errorf("First message should be from user, got %s", history[0].Role)
	}
	if history[1].Role != "assistant" {
		t.Errorf("Second message should be from assistant, got %s", history[1].Role)
	}
}

func TestMemoryManagerSession(t *testing.T) {
	database := setupTestDB(t)
	mm := NewMemoryManager(database, nil)
	
	originalSession := mm.GetSessionID()
	
	// Start a new session
	err := mm.StartSession()
	if err != nil {
		t.Fatalf("Failed to start session: %v", err)
	}
	
	newSession := mm.GetSessionID()
	if newSession == originalSession {
		t.Error("New session should have different ID")
	}
	
	// Record some activity
	mm.RecordEdit(EditEvent{
		Type:       "insert",
		WordsDelta: 10,
	})
	
	// End session
	err = mm.EndSession()
	if err != nil {
		t.Fatalf("Failed to end session: %v", err)
	}
}

func TestMemoryManagerStats(t *testing.T) {
	database := setupTestDB(t)
	mm := NewMemoryManager(database, nil)
	
	// Record some data
	mm.RecordEpisode(EpisodicEvent{EventType: EventTypeEdit})
	mm.SetPreference("test_key", "test_value", 0.5, "test")
	mm.RecordChatMessage(ChatMessage{Role: "user", Content: "test"})
	
	stats, err := mm.GetMemoryStats()
	if err != nil {
		t.Fatalf("Failed to get memory stats: %v", err)
	}
	
	if stats["episode_count"].(int) < 1 {
		t.Error("Expected at least 1 episode")
	}
	
	if stats["preference_count"].(int) < 1 {
		t.Error("Expected at least 1 preference")
	}
	
	if stats["conversation_count"].(int) < 1 {
		t.Error("Expected at least 1 conversation")
	}
}

func TestMemoryManagerClearMemory(t *testing.T) {
	database := setupTestDB(t)
	mm := NewMemoryManager(database, nil)
	
	// Record some data
	mm.RecordEpisode(EpisodicEvent{EventType: EventTypeEdit})
	mm.SetPreference("test_key", "test_value", 0.5, "test")
	mm.RecordChatMessage(ChatMessage{Role: "user", Content: "test"})
	
	// Clear all memory
	err := mm.ClearMemory(true, true, true)
	if err != nil {
		t.Fatalf("Failed to clear memory: %v", err)
	}
	
	// Verify cleared
	stats, err := mm.GetMemoryStats()
	if err != nil {
		t.Fatalf("Failed to get memory stats: %v", err)
	}
	
	if stats["episode_count"].(int) != 0 {
		t.Error("Episodes should be cleared")
	}
	
	if stats["preference_count"].(int) != 0 {
		t.Error("Preferences should be cleared")
	}
	
	if stats["conversation_count"].(int) != 0 {
		t.Error("Conversations should be cleared")
	}
}

func TestProgressStateString(t *testing.T) {
	tests := []struct {
		state    ProgressState
		expected string
	}{
		{StateStarting, "starting"},
		{StateFlowing, "flowing"},
		{StateStuck, "stuck"},
		{StateRefining, "refining"},
		{StateReviewing, "reviewing"},
		{ProgressState(99), "unknown"},
	}
	
	for _, tc := range tests {
		if tc.state.String() != tc.expected {
			t.Errorf("Expected %s for %d, got %s", tc.expected, tc.state, tc.state.String())
		}
	}
}

func TestSuggestionTypeString(t *testing.T) {
	tests := []struct {
		st       SuggestionType
		expected string
	}{
		{SuggestNextLine, "next_line"},
		{SuggestChord, "chord"},
		{SuggestBreak, "break"},
		{SuggestReference, "reference"},
		{SuggestStructure, "structure"},
		{SuggestRhyme, "rhyme"},
		{SuggestGoal, "goal"},
		{SuggestionType(99), "unknown"},
	}
	
	for _, tc := range tests {
		if tc.st.String() != tc.expected {
			t.Errorf("Expected %s for %d, got %s", tc.expected, tc.st, tc.st.String())
		}
	}
}

func TestDefaultAgentConfig(t *testing.T) {
	config := DefaultAgentConfig()
	
	if config == nil {
		t.Fatal("DefaultAgentConfig returned nil")
	}
	
	if config.MaxEpisodicEvents != 1000 {
		t.Errorf("Expected MaxEpisodicEvents 1000, got %d", config.MaxEpisodicEvents)
	}
	
	if config.SessionTimeout != 2*time.Hour {
		t.Errorf("Expected SessionTimeout 2h, got %v", config.SessionTimeout)
	}
	
	if config.MinConfidence != 0.6 {
		t.Errorf("Expected MinConfidence 0.6, got %f", config.MinConfidence)
	}
	
	if !config.LearningEnabled {
		t.Error("Expected LearningEnabled to be true")
	}
}
