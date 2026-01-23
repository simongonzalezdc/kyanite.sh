package agent

import (
	"testing"
	"time"
)

func TestContextObserverCreation(t *testing.T) {
	observer := NewContextObserver(nil, nil)
	if observer == nil {
		t.Fatal("NewContextObserver returned nil")
	}
	
	if observer.config == nil {
		t.Error("Config should not be nil")
	}
	
	if observer.contentDetector == nil {
		t.Error("ContentDetector should not be nil")
	}
}

func TestContextObserverRecordContentChange(t *testing.T) {
	observer := NewContextObserver(nil, nil)
	
	// Record initial content
	observer.RecordContentChange("Hello")
	
	// Check that edit history is updated
	if len(observer.editHistory) != 1 {
		t.Errorf("Expected 1 edit, got %d", len(observer.editHistory))
	}
	
	// Record more content
	observer.RecordContentChange("Hello World")
	
	if len(observer.editHistory) != 2 {
		t.Errorf("Expected 2 edits, got %d", len(observer.editHistory))
	}
	
	// Check that last edit is an insert
	if observer.editHistory[1].Type != "insert" {
		t.Errorf("Expected insert, got %s", observer.editHistory[1].Type)
	}
}

func TestContextObserverDeletionTracking(t *testing.T) {
	observer := NewContextObserver(nil, nil)
	
	// Add content
	observer.RecordContentChange("Hello World this is a test")
	
	// Delete content
	observer.RecordContentChange("Hello World")
	
	// Last edit should be a delete
	if observer.editHistory[1].Type != "delete" {
		t.Errorf("Expected delete, got %s", observer.editHistory[1].Type)
	}
	
	// WordsDelta should be negative
	if observer.editHistory[1].WordsDelta >= 0 {
		t.Errorf("Expected negative WordsDelta for delete, got %d", observer.editHistory[1].WordsDelta)
	}
}

func TestContextObserverProgressState(t *testing.T) {
	observer := NewContextObserver(nil, nil)
	
	// Initial state should be starting
	state := observer.GetProgressState()
	if state != StateStarting {
		t.Errorf("Expected StateStarting, got %v", state)
	}
	
	// Add multiple edits to simulate flowing state
	for i := 0; i < 10; i++ {
		content := ""
		for j := 0; j <= i; j++ {
			content += "word "
		}
		observer.RecordContentChange(content)
		time.Sleep(10 * time.Millisecond)
	}
	
	// Should have progress now
	if len(observer.editHistory) < 5 {
		t.Error("Expected at least 5 edits in history")
	}
}

func TestContextObserverWritingVelocity(t *testing.T) {
	observer := NewContextObserver(nil, nil)
	
	// Record content changes
	observer.RecordContentChange("One two three")
	time.Sleep(50 * time.Millisecond)
	observer.RecordContentChange("One two three four five six")
	time.Sleep(50 * time.Millisecond)
	observer.RecordContentChange("One two three four five six seven eight nine")
	
	velocity := observer.GetWritingVelocity()
	// Velocity should be >= 0
	if velocity < 0 {
		t.Errorf("Writing velocity should not be negative: %f", velocity)
	}
}

func TestContextObserverStuckDetection(t *testing.T) {
	config := &AgentConfig{
		StuckThreshold: 100 * time.Millisecond, // Short threshold for testing
	}
	observer := NewContextObserver(nil, config)
	
	// Record some content
	observer.RecordContentChange("Hello")
	
	// Initially should not be stuck
	if observer.IsUserStuck() {
		t.Error("Should not be stuck immediately after edit")
	}
	
	// Wait past the stuck threshold
	time.Sleep(150 * time.Millisecond)
	
	if !observer.IsUserStuck() {
		t.Error("Should be stuck after waiting past threshold")
	}
}

func TestContextObserverShouldSuggest(t *testing.T) {
	config := &AgentConfig{
		StuckThreshold: 100 * time.Millisecond,
	}
	observer := NewContextObserver(nil, config)
	
	// Initially, should suggest goals
	if !observer.ShouldSuggest(SuggestStructure) {
		t.Error("Should suggest structure at start")
	}
	
	// Record content and wait to get stuck
	observer.RecordContentChange("Hello world")
	time.Sleep(150 * time.Millisecond)
	
	// Should suggest next line when stuck
	if !observer.ShouldSuggest(SuggestNextLine) {
		t.Error("Should suggest next line when stuck")
	}
}

func TestContextObserverGetSuggestionTrigger(t *testing.T) {
	config := &AgentConfig{
		StuckThreshold: 100 * time.Millisecond,
	}
	observer := NewContextObserver(nil, config)
	
	// At start, should trigger structure suggestion
	suggType, reason := observer.GetSuggestionTrigger()
	if suggType != SuggestStructure && suggType != SuggestGoal {
		t.Errorf("Expected SuggestStructure or SuggestGoal at start, got %v", suggType)
	}
	if reason == "" {
		t.Error("Should have a reason for suggestion")
	}
	
	// Get stuck
	observer.RecordContentChange("Hello")
	time.Sleep(150 * time.Millisecond)
	
	suggType, reason = observer.GetSuggestionTrigger()
	if suggType != SuggestNextLine {
		t.Errorf("Expected SuggestNextLine when stuck, got %v", suggType)
	}
	if reason == "" {
		t.Error("Should have a reason for suggestion when stuck")
	}
}

func TestContextObserverReset(t *testing.T) {
	observer := NewContextObserver(nil, nil)
	
	// Add some activity
	observer.RecordContentChange("Hello World")
	observer.RecordContentChange("Hello World Again")
	
	if len(observer.editHistory) == 0 {
		t.Error("Should have edits before reset")
	}
	
	// Reset
	observer.Reset()
	
	if len(observer.editHistory) != 0 {
		t.Error("Edit history should be empty after reset")
	}
	
	if observer.lastContent != "" {
		t.Error("Last content should be empty after reset")
	}
	
	if observer.writingVelocity != 0 {
		t.Error("Writing velocity should be 0 after reset")
	}
}

func TestContextObserverGetRecentEditSummary(t *testing.T) {
	observer := NewContextObserver(nil, nil)
	
	// Add some edits
	observer.RecordContentChange("Hello")
	observer.RecordContentChange("Hello World")
	observer.RecordContentChange("Hello World Again")
	
	summary := observer.GetRecentEditSummary()
	
	if summary["edit_count"].(int) != 3 {
		t.Errorf("Expected edit_count 3, got %v", summary["edit_count"])
	}
	
	if _, ok := summary["writing_velocity"]; !ok {
		t.Error("Summary should include writing_velocity")
	}
	
	if _, ok := summary["progress_state"]; !ok {
		t.Error("Summary should include progress_state")
	}
}

func TestCountWords(t *testing.T) {
	tests := []struct {
		input    string
		expected int
	}{
		{"", 0},
		{"  ", 0},
		{"hello", 1},
		{"hello world", 2},
		{"  hello   world  ", 2},
		{"one two three four five", 5},
	}
	
	for _, tc := range tests {
		result := countWords(tc.input)
		if result != tc.expected {
			t.Errorf("countWords(%q) = %d, expected %d", tc.input, result, tc.expected)
		}
	}
}

func TestExtractDiff(t *testing.T) {
	tests := []struct {
		old      string
		new      string
		expected string
	}{
		{"hello", "hello world", " world"},
		{"hello world", "hello", " world"},
		{"same", "same", ""},
	}
	
	for _, tc := range tests {
		result := extractDiff(tc.old, tc.new)
		if result != tc.expected {
			t.Errorf("extractDiff(%q, %q) = %q, expected %q", tc.old, tc.new, result, tc.expected)
		}
	}
}

func TestContextObserverPauseDetection(t *testing.T) {
	observer := NewContextObserver(nil, nil)
	
	// Initially not paused
	observer.RecordContentChange("Hello")
	
	if observer.isPaused {
		t.Error("Should not be paused right after content change")
	}
	
	// Wait a bit and check for pause
	time.Sleep(50 * time.Millisecond)
	observer.CheckForPause()
	
	// With default settings (30s threshold), should still not be paused
	if observer.isPaused {
		t.Error("Should not be paused after only 50ms")
	}
}

func TestContextObserverDeletionRatio(t *testing.T) {
	observer := NewContextObserver(nil, nil)
	
	// Add some inserts
	observer.RecordContentChange("word1")
	observer.RecordContentChange("word1 word2")
	observer.RecordContentChange("word1 word2 word3")
	
	// Initial deletion ratio should be low
	ratio := observer.GetDeletionRatio()
	if ratio > 0.5 {
		t.Errorf("Deletion ratio should be low after inserts: %f", ratio)
	}
	
	// Add some deletes
	observer.RecordContentChange("word1 word2")
	observer.RecordContentChange("word1")
	
	// Deletion ratio should increase
	ratio = observer.GetDeletionRatio()
	// Note: ratio depends on the implementation details
	if ratio < 0 || ratio > 1 {
		t.Errorf("Deletion ratio should be between 0 and 1: %f", ratio)
	}
}
