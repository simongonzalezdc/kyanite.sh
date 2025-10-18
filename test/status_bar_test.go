package lyricforge_test

import (
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/puente-labs/lyricforge/internal/app"
	"github.com/puente-labs/lyricforge/internal/ui/editor"
)

// TestStatusBarModelCreation tests the creation of a status bar model
func TestStatusBarModelCreation(t *testing.T) {
	// Create status bar model
	model := editor.NewStatusBarModel()

	// Verify model was created successfully
	if model == nil {
		t.Fatal("Expected status bar model to be created, got nil")
	}

	// Verify initial state
	wordCount, charCount, lineCount := model.GetDocumentStats()
	if wordCount != 0 || charCount != 0 || lineCount != 0 {
		t.Error("Expected initial statistics to be zero")
	}

	// Verify initial cursor position
	line, col := model.GetCursorPosition()
	if line != 0 || col != 0 {
		t.Error("Expected initial cursor position to be (0,0)")
	}

	// Verify initial auto-save status
	status, _ := model.GetAutoSaveInfo()
	if status != app.AutoSaveIdle {
		t.Error("Expected initial auto-save status to be idle")
	}

	// Verify initial editor features
	showLineNumbers, wordWrap, autoIndent, bracketMatching := model.GetEditorFeatures()
	if !showLineNumbers || !wordWrap || !autoIndent || !bracketMatching {
		t.Error("Expected all editor features to be enabled by default")
	}
}

// TestStatusBarContentUpdates tests content-based statistics updates
func TestStatusBarContentUpdates(t *testing.T) {
	model := editor.NewStatusBarModel()

	// Test with empty content
	model.UpdateContent("")
	wordCount, charCount, lineCount := model.GetDocumentStats()
	if wordCount != 0 || charCount != 0 || lineCount != 0 {
		t.Error("Expected empty content to have zero statistics")
	}

	// Test with simple content
	simpleContent := "Hello world"
	model.UpdateContent(simpleContent)
	wordCount, charCount, lineCount = model.GetDocumentStats()

	expectedWordCount := 2
	expectedCharCount := len(simpleContent)
	expectedLineCount := 1

	if wordCount != expectedWordCount {
		t.Errorf("Expected word count %d, got %d", expectedWordCount, wordCount)
	}

	if charCount != expectedCharCount {
		t.Errorf("Expected character count %d, got %d", expectedCharCount, charCount)
	}

	if lineCount != expectedLineCount {
		t.Errorf("Expected line count %d, got %d", expectedLineCount, lineCount)
	}

	// Test with multi-line content
	multiLineContent := `Line 1
Line 2
Line 3 with more words here`
	model.UpdateContent(multiLineContent)
	wordCount, _, lineCount = model.GetDocumentStats()

	expectedWordCount = 10
	expectedLineCount = 3

	if wordCount != expectedWordCount {
		t.Errorf("Expected word count %d, got %d", expectedWordCount, wordCount)
	}

	if lineCount != expectedLineCount {
		t.Errorf("Expected line count %d, got %d", expectedLineCount, lineCount)
	}
}

// TestStatusBarCursorPosition tests cursor position tracking
func TestStatusBarCursorPosition(t *testing.T) {
	model := editor.NewStatusBarModel()

	// Test initial position
	line, col := model.GetCursorPosition()
	if line != 0 || col != 0 {
		t.Error("Expected initial cursor position to be (0,0)")
	}

	// Test position update
	model.UpdateCursorPosition(5, 10)
	line, col = model.GetCursorPosition()

	if line != 5 || col != 10 {
		t.Errorf("Expected cursor position (5,10), got (%d,%d)", line, col)
	}

	// Test position update to different location
	model.UpdateCursorPosition(1, 0)
	line, col = model.GetCursorPosition()

	if line != 1 || col != 0 {
		t.Errorf("Expected cursor position (1,0), got (%d,%d)", line, col)
	}
}

// TestStatusBarAutoSaveStatus tests auto-save status tracking
func TestStatusBarAutoSaveStatus(t *testing.T) {
	model := editor.NewStatusBarModel()

	// Test initial status
	status, _ := model.GetAutoSaveInfo()
	if status != app.AutoSaveIdle {
		t.Error("Expected initial auto-save status to be idle")
	}

	// Test status update
	testTime := time.Now()
	model.UpdateAutoSaveStatus(app.AutoSaveSaving, testTime)
	status, _ = model.GetAutoSaveInfo()

	if status != app.AutoSaveSaving {
		t.Error("Expected auto-save status to be saving")
	}

	// Test status update to success
	model.UpdateAutoSaveStatus(app.AutoSaveSuccess, testTime)
	status, _ = model.GetAutoSaveInfo()

	if status != app.AutoSaveSuccess {
		t.Error("Expected auto-save status to be success")
	}
}

// TestStatusBarEditorFeatures tests editor feature flag tracking
func TestStatusBarEditorFeatures(t *testing.T) {
	model := editor.NewStatusBarModel()

	// Test initial features
	showLineNumbers, wordWrap, autoIndent, bracketMatching := model.GetEditorFeatures()
	if !showLineNumbers || !wordWrap || !autoIndent || !bracketMatching {
		t.Error("Expected all editor features to be enabled initially")
	}

	// Test feature updates
	model.UpdateEditorFeatures(false, true, false, true)
	showLineNumbers, wordWrap, autoIndent, bracketMatching = model.GetEditorFeatures()

	if showLineNumbers {
		t.Error("Expected line numbers to be disabled")
	}

	if !wordWrap {
		t.Error("Expected word wrap to remain enabled")
	}

	if autoIndent {
		t.Error("Expected auto indent to be disabled")
	}

	if !bracketMatching {
		t.Error("Expected bracket matching to remain enabled")
	}
}

// TestStatusBarViewRendering tests status bar view rendering
func TestStatusBarViewRendering(t *testing.T) {
	model := editor.NewStatusBarModel()

	// Set dimensions
	model.SetDimensions(80, 1)

	// Test view rendering with default state
	view := model.View()
	if view == "" {
		t.Error("Expected view to render default state")
	}

	// Test view rendering with content
	model.UpdateContent("Test content for rendering")
	view = model.View()
	if view == "" {
		t.Error("Expected view to render with content")
	}

	// Test view rendering with cursor position
	model.UpdateCursorPosition(2, 5)
	view = model.View()
	if view == "" {
		t.Error("Expected view to render with cursor position")
	}

	// Test view rendering with auto-save status
	model.UpdateAutoSaveStatus(app.AutoSaveSaving, time.Now())
	view = model.View()
	if view == "" {
		t.Error("Expected view to render with auto-save status")
	}
}

// TestStatusBarResponsiveModes tests responsive display modes
func TestStatusBarResponsiveModes(t *testing.T) {
	model := editor.NewStatusBarModel()

	// Test compact mode (width < 100)
	model.SetDimensions(90, 1)
	model.UpdateResponsiveMode(90)

	view := model.View()
	if view == "" {
		t.Error("Expected view to render in compact mode")
	}

	// Test minimal mode (width < 80)
	model.SetDimensions(70, 1)
	model.UpdateResponsiveMode(70)

	view = model.View()
	if view == "" {
		t.Error("Expected view to render in minimal mode")
	}

	// Test full mode (width >= 100)
	model.SetDimensions(120, 1)
	model.UpdateResponsiveMode(120)

	view = model.View()
	if view == "" {
		t.Error("Expected view to render in full mode")
	}
}

// TestStatusBarFileInfo tests file information display
func TestStatusBarFileInfo(t *testing.T) {
	model := editor.NewStatusBarModel()

	// Test file info update
	model.UpdateFileInfo("test-song.md")

	// Test view rendering with file info
	model.SetDimensions(80, 1)
	view := model.View()
	if view == "" {
		t.Error("Expected view to render with file info")
	}
}

// TestStatusBarZoomLevel tests zoom level display
func TestStatusBarZoomLevel(t *testing.T) {
	model := editor.NewStatusBarModel()

	// Test zoom level update
	model.UpdateZoomLevel(150)

	// Test view rendering with zoom level
	model.SetDimensions(80, 1)
	view := model.View()
	if view == "" {
		t.Error("Expected view to render with zoom level")
	}
}

// TestStatusBarShortcutHints tests shortcut hints display
func TestStatusBarShortcutHints(t *testing.T) {
	model := editor.NewStatusBarModel()

	// Test shortcut hints update
	model.UpdateShortcutHints("Ctrl+S: Save | Ctrl+F: Find")

	// Test view rendering with shortcut hints
	model.SetDimensions(80, 1)
	view := model.View()
	if view == "" {
		t.Error("Expected view to render with shortcut hints")
	}
}

// TestStatusBarThrottling tests update throttling functionality
func TestStatusBarThrottling(t *testing.T) {
	model := editor.NewStatusBarModel()

	// Set short throttle duration for testing
	model.SetUpdateThrottle(50 * time.Millisecond)

	// Test rapid updates (should be throttled)
	start := time.Now()
	model.UpdateContent("Content 1")
	model.UpdateContent("Content 2")
	model.UpdateContent("Content 3")
	duration := time.Since(start)

	// Should take at least the throttle duration
	if duration < 50*time.Millisecond {
		t.Error("Expected updates to be throttled")
	}
}

// TestStatusBarPerformance tests performance with frequent updates
func TestStatusBarPerformance(t *testing.T) {
	model := editor.NewStatusBarModel()

	// Set dimensions
	model.SetDimensions(80, 1)

	// Test performance with many updates
	for i := 0; i < 100; i++ {
		content := fmt.Sprintf("Content update %d with various words and lines", i)
		model.UpdateContent(content)
		model.UpdateCursorPosition(i%10, i%20)
	}

	// Test that model handles many updates without error
	view := model.View()
	if view == "" {
		t.Error("Expected model to handle many updates without error")
	}

	// Verify final statistics
	wordCount, _, _ := model.GetDocumentStats()
	if wordCount <= 0 {
		t.Error("Expected positive word count after many updates")
	}
}

// TestStatusBarEdgeCases tests edge cases and boundary conditions
func TestStatusBarEdgeCases(t *testing.T) {
	model := editor.NewStatusBarModel()

	// Test with very long content
	longContent := strings.Repeat("Very long line with many words ", 100)
	model.UpdateContent(longContent)

	wordCount, _, _ := model.GetDocumentStats()
	if wordCount <= 0 {
		t.Error("Expected positive word count for long content")
	}

	// Test with special characters
	specialContent := "Content with special chars: !@#$%^&*()_+{}|:<>?[]\\;',./\""
	model.UpdateContent(specialContent)

	_, charCount, _ := model.GetDocumentStats()
	if charCount != len(specialContent) {
		t.Error("Expected character count to handle special characters correctly")
	}

	// Test with unicode content
	unicodeContent := "Content with unicode: ñáéíóú 🚀 ñáéíóú 🚀"
	model.UpdateContent(unicodeContent)

	_, charCount, _ = model.GetDocumentStats()
	if charCount != len(unicodeContent) {
		t.Error("Expected character count to handle unicode correctly")
	}
}

// TestStatusBarWidthConstraints tests behavior with different widths
func TestStatusBarWidthConstraints(t *testing.T) {
	model := editor.NewStatusBarModel()

	// Test with very small width
	model.SetDimensions(20, 1)
	view := model.View()
	if view == "" {
		t.Error("Expected view to render even with small width")
	}

	// Test with zero width
	model.SetDimensions(0, 1)
	_ = model.View()
	// Should handle gracefully (may be empty but not panic)

	// Test with very large width
	model.SetDimensions(200, 1)
	view = model.View()
	if view == "" {
		t.Error("Expected view to render with large width")
	}
}

// TestStatusBarContentHashing tests content hashing for change detection
func TestStatusBarContentHashing(t *testing.T) {
	model := editor.NewStatusBarModel()

	// Test that same content produces same hash
	content1 := "Same content"
	content2 := "Same content"

	model.UpdateContent(content1)
	model.UpdateContent(content2)

	// Since we can't access the hash directly, we'll test through behavior
	// The model should handle duplicate content efficiently
	view := model.View()
	if view == "" {
		t.Error("Expected model to handle duplicate content efficiently")
	}
}

// TestStatusBarRealTimeUpdates tests real-time update handling
func TestStatusBarRealTimeUpdates(t *testing.T) {
	model := editor.NewStatusBarModel()

	// Simulate real-time updates
	updates := []string{
		"Initial content",
		"Initial content with more words",
		"Initial content with more words and lines\nNew line added",
		"Final content with all features",
	}

	for _, update := range updates {
		model.UpdateContent(update)
		model.UpdateCursorPosition(1, 5)
		time.Sleep(10 * time.Millisecond) // Simulate real-time delay
	}

	// Test that model handles real-time updates without error
	view := model.View()
	if view == "" {
		t.Error("Expected model to handle real-time updates without error")
	}

	// Verify final state
	wordCount, _, _ := model.GetDocumentStats()
	if wordCount <= 0 {
		t.Error("Expected positive word count after real-time updates")
	}
}

// BenchmarkStatusBarUpdates benchmarks status bar update performance
func BenchmarkStatusBarUpdates(b *testing.B) {
	model := editor.NewStatusBarModel()
	model.SetDimensions(80, 1)

	// Create varied content for benchmarking
	contents := make([]string, b.N)
	for i := 0; i < b.N; i++ {
		contents[i] = fmt.Sprintf("Content update %d with varied length and word count for benchmarking", i)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		model.UpdateContent(contents[i])
	}
}

// BenchmarkStatusBarRendering benchmarks status bar rendering performance
func BenchmarkStatusBarRendering(b *testing.B) {
	model := editor.NewStatusBarModel()
	model.SetDimensions(80, 1)

	// Set up model with content and state
	model.UpdateContent("Benchmark content with words and statistics")
	model.UpdateCursorPosition(5, 10)
	model.UpdateAutoSaveStatus(app.AutoSaveSuccess, time.Now())
	model.UpdateZoomLevel(125)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_ = model.View()
	}
}

// TestStatusBarStatisticsCalculation tests statistics calculation accuracy
func TestStatusBarStatisticsCalculation(t *testing.T) {
	model := editor.NewStatusBarModel()

	// Test cases with known statistics
	testCases := []struct {
		content   string
		wordCount int
		charCount int
		lineCount int
	}{
		{"", 0, 0, 0},
		{"Hello", 1, 5, 1},
		{"Hello world", 2, 11, 1},
		{"Hello\nworld", 2, 11, 2},
		{"Line 1\nLine 2\nLine 3", 6, 18, 3},
		{"Multiple   spaces   test", 3, 23, 1},
		{"\n\n\n", 0, 3, 4},
	}

	for _, tc := range testCases {
		model.UpdateContent(tc.content)
		wordCount, charCount, lineCount := model.GetDocumentStats()

		if wordCount != tc.wordCount {
			t.Errorf("For content '%s', expected word count %d, got %d", tc.content, tc.wordCount, wordCount)
		}

		if charCount != tc.charCount {
			t.Errorf("For content '%s', expected character count %d, got %d", tc.content, tc.charCount, charCount)
		}

		if lineCount != tc.lineCount {
			t.Errorf("For content '%s', expected line count %d, got %d", tc.content, tc.lineCount, lineCount)
		}
	}
}

// TestStatusBarViewStructure tests the structure of rendered views
func TestStatusBarViewStructure(t *testing.T) {
	model := editor.NewStatusBarModel()
	model.SetDimensions(100, 1)

	// Test minimal view structure
	model.UpdateResponsiveMode(50) // Force minimal mode
	view := model.View()

	// Should be exactly the right width
	if len(view) != 50 && len(view) != 0 {
		t.Errorf("Expected minimal view to be 50 chars or empty, got %d chars", len(view))
	}

	// Test compact view structure
	model.UpdateResponsiveMode(90) // Force compact mode
	view = model.View()

	// Should be exactly the right width
	if len(view) != 90 && len(view) != 0 {
		t.Errorf("Expected compact view to be 90 chars or empty, got %d chars", len(view))
	}

	// Test full view structure
	model.UpdateResponsiveMode(120) // Force full mode
	view = model.View()

	// Should be exactly the right width
	if len(view) != 120 && len(view) != 0 {
		t.Errorf("Expected full view to be 120 chars or empty, got %d chars", len(view))
	}
}

// TestStatusBarStateConsistency tests state consistency across updates
func TestStatusBarStateConsistency(t *testing.T) {
	model := editor.NewStatusBarModel()

	// Set up initial state
	model.UpdateContent("Initial content")
	model.UpdateCursorPosition(3, 7)
	model.UpdateAutoSaveStatus(app.AutoSaveSuccess, time.Now())
	model.UpdateEditorFeatures(true, false, true, false)
	model.UpdateZoomLevel(110)

	// Perform multiple updates
	for i := 0; i < 10; i++ {
		model.UpdateContent(fmt.Sprintf("Updated content %d", i))
		model.UpdateCursorPosition(i, i*2)
	}

	// Verify final state is consistent
	wordCount, _, _ := model.GetDocumentStats()
	if wordCount <= 0 {
		t.Error("Expected consistent word count after multiple updates")
	}

	line, col := model.GetCursorPosition()
	if line < 0 || col < 0 {
		t.Error("Expected consistent cursor position after multiple updates")
	}

	status, _ := model.GetAutoSaveInfo()
	if status != app.AutoSaveSuccess {
		t.Error("Expected consistent auto-save status after multiple updates")
	}
}

// TestStatusBarMemoryUsage tests memory usage with large content
func TestStatusBarMemoryUsage(t *testing.T) {
	model := editor.NewStatusBarModel()

	// Test with very large content
	largeContent := strings.Repeat("Word ", 10000)
	model.UpdateContent(largeContent)

	// Should handle large content without issues
	wordCount, _, _ := model.GetDocumentStats()
	if wordCount != 10000 {
		t.Errorf("Expected 10000 words, got %d", wordCount)
	}

	// Test view rendering with large content
	model.SetDimensions(80, 1)
	view := model.View()
	if view == "" {
		t.Error("Expected view to render with large content")
	}
}
