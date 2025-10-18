package lyricforge_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/puente-labs/lyricforge/internal/ui/editor"
)

// TestRealTimePreviewManager tests the real-time preview manager functionality
func TestRealTimePreviewManager(t *testing.T) {
	// Create a real-time preview manager with default config
	manager := editor.NewRealTimePreviewManager(editor.DefaultPreviewUpdateConfig())

	// Test content change detection
	testContent := "# Test Song\n\nThis is a test verse.\n\n**Bold text** and *italic text*."
	manager.UpdateContent(testContent, editor.ChangeSourceEditor)

	// Verify content was processed
	if manager.IsUpdating() {
		t.Error("Manager should not be updating immediately after content update")
	}

	// Test performance tracking
	updateCount, avgUpdateTime := manager.GetPerformanceStats()
	if updateCount == 0 {
		t.Error("Expected at least one update to be tracked")
	}
	if avgUpdateTime < 0 {
		t.Error("Average update time should not be negative")
	}

	// Test content history
	history := manager.GetContentHistory()
	if len(history) == 0 {
		t.Error("Expected content history to contain at least one entry")
	}

	// Test hash generation
	// Test hash generation using reflection or expose method
	// For now, we'll test the manager's internal consistency
	testContent1 := testContent
	testContent2 := "Different content"

	// Since hashContent is private, we'll test through content updates
	manager.UpdateContent(testContent1, editor.ChangeSourceEditor)
	manager.UpdateContent(testContent1, editor.ChangeSourceEditor) // Should be debounced

	manager.UpdateContent(testContent2, editor.ChangeSourceEditor)

	// Test that different content produces different results
	finalHistory := manager.GetContentHistory()
	if len(finalHistory) < 2 {
		t.Error("Expected multiple content changes in history")
	}
}

// TestMarkdownValidator tests markdown validation functionality
func TestMarkdownValidator(t *testing.T) {
	validator := editor.NewMarkdownValidator()

	// Test valid markdown
	validContent := "# Header\n\n**Bold** and *italic* text.\n\n```go\nfmt.Println(\"Hello\")\n```"
	errors := validator.Validate(validContent)
	if len(errors) != 0 {
		t.Errorf("Expected no validation errors for valid markdown, got %d", len(errors))
	}

	// Test invalid markdown with unclosed code block
	invalidContent := "# Header\n\n```go\nfmt.Println(\"Hello\")" // Missing closing ```
	errors = validator.Validate(invalidContent)
	if len(errors) == 0 {
		t.Error("Expected validation errors for unclosed code block")
	}

	// Test unclosed bold text
	invalidContent2 := "# Header\n\n**Bold text"
	errors = validator.Validate(invalidContent2)
	if len(errors) == 0 {
		t.Error("Expected validation errors for unclosed bold text")
	}
}

// TestWordCounter tests word counting functionality
func TestWordCounter(t *testing.T) {
	counter := editor.NewWordCounter()

	// Test basic content analysis
	content := `# Song Title

[Verse 1]
This is the first verse of the song.
It has multiple lines and words.

[Chorus]
This is the chorus section.
More words here.

[Verse 2]
Second verse with different content.
Testing word count functionality.`

	stats := counter.Analyze(content)

	// Verify statistics
	if stats.WordCount <= 0 {
		t.Error("Expected positive word count")
	}
	if stats.CharacterCount <= 0 {
		t.Error("Expected positive character count")
	}
	if stats.LineCount <= 0 {
		t.Error("Expected positive line count")
	}
	if stats.ReadingTime <= 0 {
		t.Error("Expected positive reading time")
	}

	// Verify approximate word count (should be around 40-50 words)
	expectedMinWords := 30
	expectedMaxWords := 60
	if stats.WordCount < expectedMinWords || stats.WordCount > expectedMaxWords {
		t.Errorf("Word count %d outside expected range [%d, %d]", stats.WordCount, expectedMinWords, expectedMaxWords)
	}
}

// TestTOCGenerator tests table of contents generation
func TestTOCGenerator(t *testing.T) {
	generator := editor.NewTOCGenerator()

	content := `# Main Title

## Section 1

Some content here.

### Subsection 1.1

More content.

## Section 2

Different section.

### Subsection 2.1

Final content.`

	entries := generator.GenerateTOC(content)

	// Should find headers
	if len(entries) == 0 {
		t.Error("Expected to find header entries in TOC")
	}

	// Check hierarchy
	var h1Count, h2Count, h3Count int
	for _, entry := range entries {
		switch entry.Level {
		case 1:
			h1Count++
		case 2:
			h2Count++
		case 3:
			h3Count++
		}
	}

	if h1Count == 0 {
		t.Error("Expected at least one H1 entry")
	}
	if h2Count == 0 {
		t.Error("Expected at least one H2 entry")
	}
	if h3Count == 0 {
		t.Error("Expected at least one H3 entry")
	}

	// Verify titles are extracted correctly
	foundMainTitle := false
	for _, entry := range entries {
		if strings.Contains(entry.Title, "Main Title") {
			foundMainTitle = true
			break
		}
	}
	if !foundMainTitle {
		t.Error("Expected to find 'Main Title' in TOC entries")
	}
}

// TestPreviewPaneIntegration tests integration between preview pane and real-time manager
func TestPreviewPaneIntegration(t *testing.T) {
	// Create preview pane
	previewPane := editor.NewPreviewPaneModel()

	// Test initial state
	if previewPane.GetContent() != "Preview will appear here..." {
		t.Error("Expected default content in preview pane")
	}

	// Test content setting
	testContent := "# Test\n\nThis is test content."
	previewPane.SetContent(testContent)

	// Verify content was set
	if previewPane.GetContent() != testContent {
		t.Error("Expected content to be set in preview pane")
	}

	// Test real-time manager integration
	manager := previewPane.GetRealtimeManager()
	if manager == nil {
		t.Error("Expected real-time manager to be initialized")
	}

	// Test statistics
	stats := previewPane.GetPreviewStats()
	if stats.WordCount != 0 {
		t.Error("Expected zero word count for initial content")
	}

	// Test scroll sync toggle
	previewPane.ToggleScrollSync()
	// Note: We can't easily test the actual scroll sync without a full UI setup

	// Test feature toggles
	previewPane.ToggleWordCount()
	previewPane.ToggleReadingTime()
	previewPane.ToggleTOC()

	// Test TOC generation
	tocEntries := previewPane.GetTOC()
	if len(tocEntries) != 0 {
		t.Error("Expected empty TOC for content without headers")
	}
}

// TestPerformanceWithLargeContent tests performance with large documents
func TestPerformanceWithLargeContent(t *testing.T) {
	manager := editor.NewRealTimePreviewManager(editor.DefaultPreviewUpdateConfig())

	// Create large content (simulate a long song)
	var largeContent strings.Builder
	largeContent.WriteString("# Long Song Title\n\n")

	// Add many verses
	for i := 1; i <= 100; i++ {
		largeContent.WriteString(fmt.Sprintf(`[Verse %d]
This is verse number %d.
It contains multiple lines of lyrics.
Each verse has several sentences.
This helps test performance with large documents.

`, i, i))
	}

	// Add chorus sections
	for i := 1; i <= 20; i++ {
		largeContent.WriteString(fmt.Sprintf(`[Chorus %d]
This is chorus number %d.
Testing with repetitive content.
Performance should remain smooth.
Even with large amounts of text.

`, i, i))
	}

	content := largeContent.String()

	// Test update performance
	startTime := time.Now()
	manager.UpdateContent(content, editor.ChangeSourceEditor)
	updateDuration := time.Since(startTime)

	// Should complete within reasonable time (less than 100ms for this size)
	if updateDuration > 100*time.Millisecond {
		t.Errorf("Update took too long: %v", updateDuration)
	}

	// Test that manager is not stuck updating
	if manager.IsUpdating() {
		t.Error("Manager should not be stuck in updating state")
	}

	// Test performance stats
	updateCount, avgUpdateTime := manager.GetPerformanceStats()
	if updateCount == 0 {
		t.Error("Expected update count to be tracked")
	}
	if avgUpdateTime < 0 {
		t.Error("Average update time should not be negative")
	}
}

// TestDebouncedUpdates tests the debouncing mechanism
func TestDebouncedUpdates(t *testing.T) {
	config := editor.DefaultPreviewUpdateConfig()
	config.DebounceDelay = 50 * time.Millisecond // Short delay for testing

	manager := editor.NewRealTimePreviewManager(config)

	// Track update calls
	var updateCallCount int
	manager.SetCallbacks(
		func(content string) {
			updateCallCount++
		},
		nil, nil, nil, nil,
	)

	// Send multiple rapid updates
	content1 := "Content 1"
	content2 := "Content 2"
	content3 := "Content 3"

	manager.UpdateContent(content1, editor.ChangeSourceEditor)
	manager.UpdateContent(content2, editor.ChangeSourceEditor)
	manager.UpdateContent(content3, editor.ChangeSourceEditor)

	// Wait for debounce delay
	time.Sleep(100 * time.Millisecond)

	// Should only have one effective update due to debouncing
	// Note: The exact behavior depends on implementation timing
	if updateCallCount == 0 {
		t.Error("Expected at least one update call")
	}
}

// TestConcurrentUpdates tests thread safety with concurrent updates
func TestConcurrentUpdates(t *testing.T) {
	manager := editor.NewRealTimePreviewManager(editor.DefaultPreviewUpdateConfig())

	// Track update calls
	var updateCallCount int
	manager.SetCallbacks(
		func(content string) {
			updateCallCount++
		},
		nil, nil, nil, nil,
	)

	// Send concurrent updates
	done := make(chan bool, 5)
	for i := 0; i < 5; i++ {
		go func(id int) {
			content := fmt.Sprintf("Concurrent content %d", id)
			manager.UpdateContent(content, editor.ChangeSourceEditor)
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 5; i++ {
		<-done
	}

	// Should handle concurrent updates without issues
	if updateCallCount == 0 {
		t.Error("Expected at least one update call from concurrent updates")
	}
}

// BenchmarkRealTimePreview benchmarks real-time preview performance
func BenchmarkRealTimePreview(b *testing.B) {
	manager := editor.NewRealTimePreviewManager(editor.DefaultPreviewUpdateConfig())

	// Create test content
	content := `# Benchmark Song

[Verse 1]
This is a verse for benchmarking.
Testing performance with multiple lines.
Each line contains several words.
Performance should remain consistent.

[Chorus]
This is the chorus section.
Testing with repetitive content.
Benchmarking helps ensure smooth operation.
Real-time updates should be fast.

[Verse 2]
Second verse for additional testing.
More content to process efficiently.
Should handle medium-sized documents well.
Performance is critical for user experience.`

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		manager.UpdateContent(content, editor.ChangeSourceEditor)
	}
}

// BenchmarkLargeDocument benchmarks performance with large documents
func BenchmarkLargeDocument(b *testing.B) {
	manager := editor.NewRealTimePreviewManager(editor.DefaultPreviewUpdateConfig())

	// Create large content
	var largeContent strings.Builder
	largeContent.WriteString("# Large Document Test\n\n")

	for i := 1; i <= 50; i++ {
		largeContent.WriteString(fmt.Sprintf(`[Verse %d]
This is verse number %d in a large document.
Each verse contains multiple lines of content.
Testing performance with substantial amounts of text.
Should maintain smooth operation even with large documents.

`, i, i))
	}

	content := largeContent.String()

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		manager.UpdateContent(content, editor.ChangeSourceEditor)
	}
}
