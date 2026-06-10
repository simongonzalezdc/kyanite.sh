package examples

import (
	"strings"
	"testing"
	"time"

	"github.com/kyanite/noise/internal/plugins"
)

func TestChordAnalyzerPlugin_NewChordAnalyzerPlugin(t *testing.T) {
	plugin := NewChordAnalyzerPlugin()

	if plugin == nil {
		t.Fatal("Failed to create chord analyzer plugin")
	}

	metadata := plugin.Metadata()
	if metadata.ID != "chord_analyzer" {
		t.Errorf("Expected plugin ID 'chord_analyzer', got '%s'", metadata.ID)
	}

	if metadata.Name != "Chord Analyzer" {
		t.Errorf("Expected plugin name 'Chord Analyzer', got '%s'", metadata.Name)
	}

	if metadata.Version != "1.0.0" {
		t.Errorf("Expected plugin version '1.0.0', got '%s'", metadata.Version)
	}

	if len(metadata.Capabilities) != 2 {
		t.Errorf("Expected 2 capabilities, got %d", len(metadata.Capabilities))
	}

	// Check capabilities
	hasTheoryTool := false
	hasEditorTool := false
	for _, cap := range metadata.Capabilities {
		if cap == plugins.CapabilityTheoryTool {
			hasTheoryTool = true
		}
		if cap == plugins.CapabilityEditorTool {
			hasEditorTool = true
		}
	}

	if !hasTheoryTool {
		t.Error("Plugin should have TheoryTool capability")
	}

	if !hasEditorTool {
		t.Error("Plugin should have EditorTool capability")
	}

	if plugin.IsEnabled() {
		t.Error("Plugin should be disabled by default")
	}
}

func TestChordAnalyzerPlugin_Lifecycle(t *testing.T) {
	plugin := NewChordAnalyzerPlugin()
	ctx := plugins.TestPluginContext(t)

	// Test Initialize
	err := plugin.Initialize(ctx)
	if err != nil {
		t.Errorf("Initialize should not error: %v", err)
	}

	if plugin.Metadata().LoadTime.IsZero() {
		t.Error("Load time should be set after initialization")
	}

	// Test Enable
	err = plugin.Enable()
	if err != nil {
		t.Errorf("Enable should not error: %v", err)
	}

	if !plugin.IsEnabled() {
		t.Error("Plugin should be enabled after Enable()")
	}

	// Test Disable
	err = plugin.Disable()
	if err != nil {
		t.Errorf("Disable should not error: %v", err)
	}

	if plugin.IsEnabled() {
		t.Error("Plugin should be disabled after Disable()")
	}

	// Test Cleanup
	err = plugin.Cleanup()
	if err != nil {
		t.Errorf("Cleanup should not error: %v", err)
	}
}

func TestChordAnalyzerPlugin_AnalyzeChords(t *testing.T) {
	plugin := NewChordAnalyzerPlugin()

	tests := []struct {
		name           string
		content        string
		expectedTotal  int
		expectedChords []string
	}{
		{
			name:           "Simple chords",
			content:        "C G Am F",
			expectedTotal:  4,
			expectedChords: []string{"C", "G", "Am", "F"},
		},
		{
			name:           "Chords with text",
			content:        "Verse: C G Am F\nChorus: C G Am F",
			expectedTotal:  8,
			expectedChords: []string{"C", "G", "Am", "F"},
		},
		{
			name:           "Complex chords",
			content:        "Cmaj7 Gsus4 Am7 F#dim",
			expectedTotal:  4,
			expectedChords: []string{"Cmaj7", "Gsus4", "Am7", "F#dim"},
		},
		{
			name:           "No chords",
			content:        "This is just lyrics without chords",
			expectedTotal:  0,
			expectedChords: []string{},
		},
		{
			name:           "Mixed case and sharps/flats",
			content:        "C G Bb F",
			expectedTotal:  4,
			expectedChords: []string{"C", "G", "Bb", "F"},
		},
		{
			name:           "Chords with numbers",
			content:        "C2 Gsus4 Am7 Fmaj7",
			expectedTotal:  4,
			expectedChords: []string{"C2", "Gsus4", "Am7", "Fmaj7"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			analysis, err := plugin.AnalyzeChords(tt.content)
			if err != nil {
				t.Errorf("AnalyzeChords should not error: %v", err)
			}

			if analysis.TotalChords != tt.expectedTotal {
				t.Errorf("Expected %d total chords, got %d", tt.expectedTotal, analysis.TotalChords)
			}

			if len(analysis.UniqueChords) != len(tt.expectedChords) {
				t.Errorf("Expected %d unique chords, got %d", len(tt.expectedChords), len(analysis.UniqueChords))
			}

			// Check that all expected chords are present
			for _, expectedChord := range tt.expectedChords {
				found := false
				for _, actualChord := range analysis.UniqueChords {
					if actualChord == expectedChord {
						found = true
						break
					}
				}
				if !found {
					t.Errorf("Expected chord '%s' not found in analysis", expectedChord)
				}
			}
		})
	}
}

func TestChordAnalyzerPlugin_ExtractChords(t *testing.T) {
	plugin := NewChordAnalyzerPlugin()

	content := "C G Am F C G Am F"
	chords := plugin.extractChords(content)

	if len(chords) != 8 {
		t.Errorf("Expected 8 chords, got %d", len(chords))
	}

	expectedChords := []string{"C", "G", "Am", "F", "C", "G", "Am", "F"}
	for i, expectedChord := range expectedChords {
		if chords[i] != expectedChord {
			t.Errorf("Expected chord '%s' at position %d, got '%s'", expectedChord, i, chords[i])
		}
	}
}

func TestChordAnalyzerPlugin_GetUniqueChords(t *testing.T) {
	plugin := NewChordAnalyzerPlugin()

	chords := []string{"C", "G", "Am", "F", "C", "G", "Am", "F"}
	uniqueChords := plugin.getUniqueChords(chords)

	if len(uniqueChords) != 4 {
		t.Errorf("Expected 4 unique chords, got %d", len(uniqueChords))
	}

	expectedUnique := []string{"C", "G", "Am", "F"}
	for _, expectedChord := range expectedUnique {
		found := false
		for _, actualChord := range uniqueChords {
			if actualChord == expectedChord {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected unique chord '%s' not found", expectedChord)
		}
	}
}

func TestChordAnalyzerPlugin_IdentifyProgressions(t *testing.T) {
	plugin := NewChordAnalyzerPlugin()

	// Test with enough chords to identify progression
	chords := []string{"C", "G", "Am", "F", "C", "G", "Am", "F"}
	progressions := plugin.identifyProgressions(chords)

	if len(progressions) == 0 {
		t.Error("Expected to identify at least one progression")
	}

	// Test with too few chords
	chords = []string{"C", "G"}
	progressions = plugin.identifyProgressions(chords)

	// Should not identify progression with too few chords
	if len(progressions) > 0 {
		t.Error("Should not identify progression with too few chords")
	}
}

func TestChordAnalyzerPlugin_DetectKeySignature(t *testing.T) {
	plugin := NewChordAnalyzerPlugin()

	// Test with C major chords
	chords := []string{"C", "G", "Am", "F", "C"}
	key := plugin.detectKeySignature(chords)

	if !strings.Contains(key, "C") {
		t.Errorf("Expected key to contain 'C', got '%s'", key)
	}

	// Test with G major chords
	chords = []string{"G", "D", "Em", "C", "G"}
	key = plugin.detectKeySignature(chords)

	if !strings.Contains(key, "G") {
		t.Errorf("Expected key to contain 'G', got '%s'", key)
	}

	// Test with empty chords
	chords = []string{}
	key = plugin.detectKeySignature(chords)

	if key != "Unknown" {
		t.Errorf("Expected 'Unknown' key for empty chords, got '%s'", key)
	}
}

func TestChordAnalyzerPlugin_CalculateComplexity(t *testing.T) {
	plugin := NewChordAnalyzerPlugin()

	// Test with simple chords
	chords := []string{"C", "G", "C", "G"}
	complexity := plugin.calculateComplexity(chords)

	expectedComplexity := len(plugin.getUniqueChords(chords)) * 10
	if complexity != expectedComplexity {
		t.Errorf("Expected complexity %d, got %d", expectedComplexity, complexity)
	}

	// Test with complex chords
	chords = []string{"Cmaj7", "Gsus4", "Am7", "F#dim", "Bm7b5", "E7#9"}
	complexity = plugin.calculateComplexity(chords)

	expectedComplexity = len(plugin.getUniqueChords(chords)) * 10
	if complexity != expectedComplexity {
		t.Errorf("Expected complexity %d, got %d", expectedComplexity, complexity)
	}

	// Test with empty chords
	chords = []string{}
	complexity = plugin.calculateComplexity(chords)

	if complexity != 0 {
		t.Errorf("Expected complexity 0 for empty chords, got %d", complexity)
	}
}

func TestChordAnalyzerPlugin_GetEditorTool(t *testing.T) {
	plugin := NewChordAnalyzerPlugin()
	tool := plugin.GetEditorTool()

	if tool == nil {
		t.Fatal("GetEditorTool should not return nil")
	}

	if tool.ID != "chord_analyzer" {
		t.Errorf("Expected tool ID 'chord_analyzer', got '%s'", tool.ID)
	}

	if tool.Name != "Analyze Chords" {
		t.Errorf("Expected tool name 'Analyze Chords', got '%s'", tool.Name)
	}

	if tool.Icon != "♪" {
		t.Errorf("Expected tool icon '♪', got '%s'", tool.Icon)
	}

	if tool.Handler == nil {
		t.Error("Tool handler should not be nil")
	}

	// Test the handler
	result, err := tool.Handler("C G Am F")
	if err != nil {
		t.Errorf("Tool handler should not error: %v", err)
	}

	if !strings.Contains(result, "Chord Analysis Results") {
		t.Error("Tool handler result should contain analysis header")
	}

	if !strings.Contains(result, "Total Chords: 4") {
		t.Error("Tool handler result should contain total chord count")
	}

	if !strings.Contains(result, "C, G, Am, F") {
		t.Error("Tool handler result should contain unique chords")
	}
}

func TestChordAnalyzerPlugin_GetManifest(t *testing.T) {
	plugin := NewChordAnalyzerPlugin()
	manifest, err := plugin.GetManifest()
	if err != nil {
		t.Errorf("GetManifest should not error: %v", err)
	}

	if len(manifest) == 0 {
		t.Error("Manifest should not be empty")
	}

	// Check that manifest contains expected fields
	manifestStr := string(manifest)
	if !strings.Contains(manifestStr, "chord_analyzer") {
		t.Error("Manifest should contain plugin ID")
	}

	if !strings.Contains(manifestStr, "Chord Analyzer") {
		t.Error("Manifest should contain plugin name")
	}

	if !strings.Contains(manifestStr, "theory_tool") {
		t.Error("Manifest should contain theory_tool capability")
	}

	if !strings.Contains(manifestStr, "editor_tool") {
		t.Error("Manifest should contain editor_tool capability")
	}
}

func TestMarkdownExportPlugin_NewMarkdownExportPlugin(t *testing.T) {
	plugin := NewMarkdownExportPlugin()

	if plugin == nil {
		t.Fatal("Failed to create markdown export plugin")
	}

	metadata := plugin.Metadata()
	if metadata.ID != "markdown_export" {
		t.Errorf("Expected plugin ID 'markdown_export', got '%s'", metadata.ID)
	}

	if metadata.Name != "Markdown Export" {
		t.Errorf("Expected plugin name 'Markdown Export', got '%s'", metadata.Name)
	}

	if metadata.Version != "1.0.0" {
		t.Errorf("Expected plugin version '1.0.0', got '%s'", metadata.Version)
	}

	if len(metadata.Capabilities) != 1 {
		t.Errorf("Expected 1 capability, got %d", len(metadata.Capabilities))
	}

	// Check capability
	if metadata.Capabilities[0] != plugins.CapabilityExportFormat {
		t.Errorf("Expected capability '%s', got '%s'", plugins.CapabilityExportFormat, metadata.Capabilities[0])
	}

	if plugin.IsEnabled() {
		t.Error("Plugin should be disabled by default")
	}
}

func TestMarkdownExportPlugin_Lifecycle(t *testing.T) {
	plugin := NewMarkdownExportPlugin()
	ctx := plugins.TestPluginContext(t)

	// Test Initialize
	err := plugin.Initialize(ctx)
	if err != nil {
		t.Errorf("Initialize should not error: %v", err)
	}

	if plugin.Metadata().LoadTime.IsZero() {
		t.Error("Load time should be set after initialization")
	}

	// Test Enable
	err = plugin.Enable()
	if err != nil {
		t.Errorf("Enable should not error: %v", err)
	}

	if !plugin.IsEnabled() {
		t.Error("Plugin should be enabled after Enable()")
	}

	// Test Disable
	err = plugin.Disable()
	if err != nil {
		t.Errorf("Disable should not error: %v", err)
	}

	if plugin.IsEnabled() {
		t.Error("Plugin should be disabled after Disable()")
	}

	// Test Cleanup
	err = plugin.Cleanup()
	if err != nil {
		t.Errorf("Cleanup should not error: %v", err)
	}
}

func TestMarkdownExportPlugin_ExportToMarkdown(t *testing.T) {
	plugin := NewMarkdownExportPlugin()

	tests := []struct {
		name     string
		content  string
		contains []string
	}{
		{
			name:    "Simple lyrics",
			content: "Verse 1\nLine 1\nLine 2\n\nChorus\nLine 3\nLine 4",
			contains: []string{
				"# Exported Lyrics",
				"Exported on:",
				"Verse 1",
				"Line 1",
				"Line 2",
				"Chorus",
				"Line 3",
				"Line 4",
			},
		},
		{
			name:    "Lyrics with chorus indicators",
			content: "Verse\nLine 1\n\n(Chorus)\nChorus line 1\nChorus line 2\n\nVerse\nLine 2",
			contains: []string{
				"# Exported Lyrics",
				"Verse",
				"Line 1",
				"**(Chorus)**",
				"**Chorus line 1**",
				"**Chorus line 2**",
			},
		},
		{
			name:    "Empty content",
			content: "",
			contains: []string{
				"# Exported Lyrics",
				"Exported on:",
			},
		},
		{
			name:    "Content with multiple empty lines",
			content: "Line 1\n\n\n\nLine 2",
			contains: []string{
				"# Exported Lyrics",
				"Line 1",
				"Line 2",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := plugin.exportToMarkdown(tt.content)
			if err != nil {
				t.Errorf("exportToMarkdown should not error: %v", err)
			}

			resultStr := string(result)
			for _, expectedContent := range tt.contains {
				if !strings.Contains(resultStr, expectedContent) {
					t.Errorf("Expected result to contain '%s', but it didn't. Result: %s", expectedContent, resultStr)
				}
			}
		})
	}
}

func TestMarkdownExportPlugin_IsChorusLine(t *testing.T) {
	plugin := NewMarkdownExportPlugin()

	tests := []struct {
		line     string
		expected bool
	}{
		{"chorus line", true},
		{"Chorus line", true},
		{"(chorus)", true},
		{"[chorus]", true},
		{"refrain line", true},
		{"(refrain)", true},
		{"[refrain]", true},
		{"hook line", true},
		{"(hook)", true},
		{"[hook]", true},
		{"verse line", false},
		{"bridge line", false},
		{"intro line", false},
		{"outro line", false},
		{"", false},
		{"some random text", false},
	}

	for _, tt := range tests {
		t.Run(tt.line, func(t *testing.T) {
			result := plugin.isChorusLine(tt.line)
			if result != tt.expected {
				t.Errorf("Expected isChorusLine('%s') to be %v, got %v", tt.line, tt.expected, result)
			}
		})
	}
}

func TestMarkdownExportPlugin_GetExportFormat(t *testing.T) {
	plugin := NewMarkdownExportPlugin()
	format := plugin.GetExportFormat()

	if format == nil {
		t.Fatal("GetExportFormat should not return nil")
	}

	if format.ID != "markdown" {
		t.Errorf("Expected format ID 'markdown', got '%s'", format.ID)
	}

	if format.Name != "Markdown" {
		t.Errorf("Expected format name 'Markdown', got '%s'", format.Name)
	}

	if format.Extension != ".md" {
		t.Errorf("Expected format extension '.md', got '%s'", format.Extension)
	}

	if format.MimeType != "text/markdown" {
		t.Errorf("Expected format MIME type 'text/markdown', got '%s'", format.MimeType)
	}

	if format.Handler == nil {
		t.Error("Format handler should not be nil")
	}

	// Test the handler
	result, err := format.Handler("Test content\nLine 2")
	if err != nil {
		t.Errorf("Format handler should not error: %v", err)
	}

	resultStr := string(result)
	if !strings.Contains(resultStr, "# Exported Lyrics") {
		t.Error("Format handler result should contain header")
	}

	if !strings.Contains(resultStr, "Test content") {
		t.Error("Format handler result should contain original content")
	}
}

func TestMarkdownExportPlugin_GetManifest(t *testing.T) {
	plugin := NewMarkdownExportPlugin()
	manifest, err := plugin.GetManifest()
	if err != nil {
		t.Errorf("GetManifest should not error: %v", err)
	}

	if len(manifest) == 0 {
		t.Error("Manifest should not be empty")
	}

	// Check that manifest contains expected fields
	manifestStr := string(manifest)
	if !strings.Contains(manifestStr, "markdown_export") {
		t.Error("Manifest should contain plugin ID")
	}

	if !strings.Contains(manifestStr, "Markdown Export") {
		t.Error("Manifest should contain plugin name")
	}

	if !strings.Contains(manifestStr, "export_format") {
		t.Error("Manifest should contain export_format capability")
	}
}

func TestExamplePlugins_SecurityValidation(t *testing.T) {
	// Test that example plugins pass security validation
	sm := plugins.TestSecurityManager(t)

	// Test Chord Analyzer Plugin
	chordPlugin := NewChordAnalyzerPlugin()
	err := sm.ValidatePluginManifest(chordPlugin.Metadata())
	if err != nil {
		t.Errorf("Chord analyzer plugin should pass security validation: %v", err)
	}

	// Test Markdown Export Plugin
	markdownPlugin := NewMarkdownExportPlugin()
	err = sm.ValidatePluginManifest(markdownPlugin.Metadata())
	if err != nil {
		t.Errorf("Markdown export plugin should pass security validation: %v", err)
	}
}

func TestExamplePlugins_Performance(t *testing.T) {
	// Test performance of chord analysis with large content
	plugin := NewChordAnalyzerPlugin()

	// Generate large content with many chords
	var contentBuilder strings.Builder
	chords := []string{"C", "G", "Am", "F", "D", "A", "Bm", "G"}
	for i := 0; i < 1000; i++ {
		contentBuilder.WriteString(chords[i%len(chords)])
		contentBuilder.WriteString(" ")
		if i%10 == 0 {
			contentBuilder.WriteString("\n")
		}
	}
	content := contentBuilder.String()

	// Measure time taken
	start := time.Now()
	_, err := plugin.AnalyzeChords(content)
	duration := time.Since(start)

	if err != nil {
		t.Errorf("AnalyzeChords should not error: %v", err)
	}

	// Should complete within reasonable time (adjust threshold as needed)
	if !relaxPerfBudgets() && duration > 100*time.Millisecond {
		t.Errorf("Analysis took too long: %v", duration)
	}
}

func TestExamplePlugins_ConcurrentAccess(t *testing.T) {
	// Test that plugins can be used concurrently
	plugin := NewChordAnalyzerPlugin()

	// Enable the plugin first
	if err := plugin.Enable(); err != nil {
		t.Fatalf("Failed to enable plugin: %v", err)
	}

	// Test concurrent chord analysis
	const numGoroutines = 10
	const numAnalyses = 100

	done := make(chan bool, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func(id int) {
			for j := 0; j < numAnalyses; j++ {
				content := "C G Am F"
				_, err := plugin.AnalyzeChords(content)
				if err != nil {
					t.Errorf("Goroutine %d: AnalyzeChords should not error: %v", id, err)
				}
			}
			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < numGoroutines; i++ {
		<-done
	}
}

// Table-driven tests for edge cases
func TestExamplePlugins_EdgeCases(t *testing.T) {
	chordPlugin := NewChordAnalyzerPlugin()
	markdownPlugin := NewMarkdownExportPlugin()

	tests := []struct {
		name     string
		testFunc func(t *testing.T)
	}{
		{
			name: "Chord analysis with special characters",
			testFunc: func(t *testing.T) {
				content := "C/G G/B Am/F# D/A"
				analysis, err := chordPlugin.AnalyzeChords(content)
				if err != nil {
					t.Errorf("AnalyzeChords should not error: %v", err)
				}
				if analysis.TotalChords == 0 {
					t.Error("Should find chords with special characters")
				}
			},
		},
		{
			name: "Chord analysis with unicode",
			testFunc: func(t *testing.T) {
				content := "C G Am F ♪ ♫ ♬"
				analysis, err := chordPlugin.AnalyzeChords(content)
				if err != nil {
					t.Errorf("AnalyzeChords should not error: %v", err)
				}
				if analysis.TotalChords != 4 {
					t.Errorf("Expected 4 chords, got %d", analysis.TotalChords)
				}
			},
		},
		{
			name: "Markdown export with very long lines",
			testFunc: func(t *testing.T) {
				longLine := strings.Repeat("This is a very long line ", 10) // Reduced length
				content := "Short line\n" + longLine + "\nAnother short line"

				result, err := markdownPlugin.exportToMarkdown(content)
				if err != nil {
					t.Errorf("exportToMarkdown should not error: %v", err)
				}

				resultStr := string(result)
				// Check if at least part of the long line is present
				if !strings.Contains(resultStr, "This is a very long line") {
					t.Errorf("Result should contain long line. Result: %s", resultStr)
				}
			},
		},
		{
			name: "Markdown export with special markdown characters",
			testFunc: func(t *testing.T) {
				content := "Line with *bold* and _italic_ text\nLine with [link](url) and `code`"

				result, err := markdownPlugin.exportToMarkdown(content)
				if err != nil {
					t.Errorf("exportToMarkdown should not error: %v", err)
				}

				resultStr := string(result)
				if !strings.Contains(resultStr, "*bold*") {
					t.Error("Result should preserve markdown formatting")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, tt.testFunc)
	}
}
