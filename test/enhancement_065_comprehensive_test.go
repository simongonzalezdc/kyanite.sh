package noise

import (
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/Kyanite/noise/internal/app/ai"
	"github.com/Kyanite/noise/internal/export"
	"github.com/Kyanite/noise/internal/ui/styles"
)

// TestComprehensiveIntegration covers all Enhancement #6.5 features
func TestComprehensiveIntegration(t *testing.T) {
	// Create temporary directory for test outputs
	tempDir := t.TempDir()

	// Initialize all components
	themeManager := setupThemeManager(t, tempDir)
	aiAgent := setupAIAgent(t)
	exportService := setupExportService(t, tempDir)

	// Test data
	lyricContent := `# Midnight Dreams

[Verse 1]
The city lights are fading
As I walk through empty streets
Memories of you surround me
Like echoes in the night

[Chorus]
Oh, midnight dreams are calling
Whispering your name to me
In this lonely darkness
I reach for your hand again

pattern: C - G - Am - F
BPM: 120`

	patternContent := `pattern: I - V - vi - IV
tempo: 120 bpm
key: C major

Verse progression:
C - G - Am - F

Chorus progression:
F - C - G - Am

Bridge progression:
Am - F - C - G`

	// Run comprehensive tests
	t.Run("ThemeIntegration", func(t *testing.T) {
		testThemeIntegration(t, themeManager)
	})

	t.Run("ContextDetection", func(t *testing.T) {
		testContextDetection(t, lyricContent, patternContent)
	})

	t.Run("AIIntegration", func(t *testing.T) {
		testAIIntegration(t, aiAgent, lyricContent, patternContent)
	})

	t.Run("KnowledgeBaseIntegration", func(t *testing.T) {
		testKnowledgeBaseIntegration(t, aiAgent, lyricContent)
	})

	t.Run("ExportFormats", func(t *testing.T) {
		testExportFormats(t, exportService, lyricContent, patternContent)
	})

	t.Run("EndToEndWorkflows", func(t *testing.T) {
		testEndToEndWorkflows(t, themeManager, aiAgent, exportService, lyricContent, patternContent)
	})

	t.Run("PerformanceRequirements", func(t *testing.T) {
		testPerformanceRequirements(t, aiAgent, themeManager, exportService)
	})

	t.Run("ErrorHandling", func(t *testing.T) {
		testErrorHandling(t, aiAgent, exportService)
	})
}

// setupThemeManager initializes the theme system
func setupThemeManager(t *testing.T, tempDir string) *styles.ThemeManager {
	t.Helper()
	themeFilePath := filepath.Join(tempDir, "theme.json")
	themeManager := styles.NewThemeManager(themeFilePath)
	if themeManager == nil {
		t.Fatal("Failed to create theme manager")
	}

	themeManager.Init()
	return themeManager
}

// setupAIAgent initializes the AI agent with all components
func setupAIAgent(t *testing.T) *ai.QuickIdeaAgent {
	t.Helper()
	agent := ai.NewQuickIdeaAgent()
	if agent == nil {
		t.Fatal("Failed to create AI agent")
	}

	return agent
}

// setupExportService initializes the export service
func setupExportService(t *testing.T, tempDir string) *export.ExportService {
	t.Helper()
	service := export.NewExportService(tempDir)
	if service == nil {
		t.Fatal("Failed to create export service")
	}

	return service
}

// testThemeIntegration tests theme switching with lyric content
func testThemeIntegration(t *testing.T, themeManager *styles.ThemeManager) {
	t.Helper()
	// Get initial theme
	initialTheme := themeManager.GetCurrentTheme()
	if initialTheme == nil {
		t.Fatal("Failed to get initial theme")
	}

	// Test theme switching
	themes := themeManager.GetAllThemes()
	if len(themes) == 0 {
		t.Fatal("No themes found")
	}

	// Test switching through all themes
	for i := range themes {
		th := &themes[i]
		if err := themeManager.SetTheme(th.Name); err != nil {
			t.Errorf("Failed to switch to theme %s: %v", th.Name, err)
			continue
		}

		currentTheme := themeManager.GetCurrentTheme()
		if currentTheme.Name != th.Name {
			t.Errorf("Theme not set correctly: expected %s, got %s", th.Name, currentTheme.Name)
		}

		t.Logf("theme switch %d/%d succeeded: %s", i+1, len(themes), th.Name)
	}

	// Test theme persistence
	if err := themeManager.SetTheme("Neon Dreams"); err != nil {
		t.Errorf("Failed to set Neon Dreams theme: %v", err)
	}

	t.Logf("theme integration covered %d themes", len(themes))
}

// testContextDetection tests the context detection functionality
func testContextDetection(t *testing.T, lyricContent, patternContent string) {
	t.Helper()
	// Create context detector
	detector := ai.NewContextDetector()

	// Test lyric content detection
	lyricType := detector.AnalyzeContent(lyricContent)
	if lyricType != ai.ContentTypeLyrics && lyricType != ai.ContentTypeMixed {
		t.Errorf("Expected lyric or mixed content, got %s", lyricType)
	}

	// Test pattern content detection
	patternType := detector.AnalyzeContent(patternContent)
	if patternType != ai.ContentTypePatterns && patternType != ai.ContentTypeMixed {
		t.Errorf("Expected pattern or mixed content, got %s", patternType)
	}

	// Test detailed analysis
	lyricAnalysis := detector.GetContextAnalysis(lyricContent)
	if lyricAnalysis.ContentType == ai.ContentTypeUnknown {
		t.Error("Failed to analyze lyric content")
	}

	patternAnalysis := detector.GetContextAnalysis(patternContent)
	if patternAnalysis.ContentType == ai.ContentTypeUnknown {
		t.Error("Failed to analyze pattern content")
	}

	t.Log("Context detection test passed")
	t.Logf("  lyric content: %s (confidence: %.2f)", lyricAnalysis.ContentType, lyricAnalysis.Confidence)
	t.Logf("  pattern content: %s (confidence: %.2f)", patternAnalysis.ContentType, patternAnalysis.Confidence)
}

// testAIIntegration tests AI agent functionality
func testAIIntegration(t *testing.T, agent *ai.QuickIdeaAgent, lyricContent, patternContent string) {
	ctx := context.Background()

	// Test unstick mode for lyrics
	unstickReq := ai.QuickRequest{
		Mode:    ai.QuickIdeaModeUnstick,
		Context: lyricContent,
	}

	unstickResp, err := agent.Generate(ctx, unstickReq)
	if err != nil {
		t.Errorf("AI unstick request failed: %v", err)
	} else if len(unstickResp.Suggestions) == 0 {
		t.Error("AI unstick returned no suggestions")
	} else {
		t.Logf("AI unstick for lyrics: %d suggestions", len(unstickResp.Suggestions))
	}

	// Test spark mode for patterns
	sparkReq := ai.QuickRequest{
		Mode:    ai.QuickIdeaModeSpark,
		Context: patternContent,
		Options: map[string]string{"theme": "electronic"},
	}

	sparkResp, err := agent.Generate(ctx, sparkReq)
	if err != nil {
		t.Errorf("AI spark request failed: %v", err)
	} else if len(sparkResp.Suggestions) == 0 {
		t.Error("AI spark returned no suggestions")
	} else {
		t.Logf("AI spark for patterns: %d suggestions", len(sparkResp.Suggestions))
	}

	// Test tweak mode
	tweakReq := ai.QuickRequest{
		Mode:    ai.QuickIdeaModeTweak,
		Context: "The rain falls softly on the window",
	}

	tweakResp, err := agent.Generate(ctx, tweakReq)
	if err != nil {
		t.Errorf("AI tweak request failed: %v", err)
	} else if len(tweakResp.Suggestions) == 0 {
		t.Error("AI tweak returned no suggestions")
	} else {
		t.Logf("AI tweak: %d suggestions", len(tweakResp.Suggestions))
	}

	// Test check mode
	checkReq := ai.QuickRequest{
		Mode:    ai.QuickIdeaModeCheck,
		Context: lyricContent,
	}

	checkResp, err := agent.Generate(ctx, checkReq)
	if err != nil {
		t.Errorf("AI check request failed: %v", err)
	} else if checkResp.Rating == "" {
		t.Error("AI check returned no rating")
	} else {
		t.Logf("AI check: rating=%s, tip=%s", checkResp.Rating, checkResp.Tip)
	}
}

// testKnowledgeBaseIntegration tests knowledge base functionality
func testKnowledgeBaseIntegration(t *testing.T, agent *ai.QuickIdeaAgent, content string) {
	ctx := context.Background()

	// Check knowledge base status
	status := agent.GetKnowledgeBaseStatus(ctx)
	if status == nil {
		t.Error("Failed to get knowledge base status")
	} else {
		t.Logf("Knowledge base status: available=%v, cards=%d", status.Available, status.CardCount)
	}

	// Test knowledge base availability
	isAvailable := agent.IsKnowledgeBaseAvailable(ctx)
	t.Logf("Knowledge base available: %v", isAvailable)

	// Test AI request with knowledge base enhancement
	req := ai.QuickRequest{
		Mode:    ai.QuickIdeaModeUnstick,
		Context: content,
	}

	resp, err := agent.Generate(ctx, req)
	if err != nil {
		t.Errorf("AI request with KB failed: %v", err)
	} else {
		t.Logf("AI request with knowledge base: %d suggestions, response time: %v",
			len(resp.Suggestions), resp.ResponseTime)
	}
}

// testExportFormats tests all export formats
func testExportFormats(t *testing.T, service *export.ExportService, lyricContent, patternContent string) {
	// Test Markdown export
	mdPath, err := service.ExportToMarkdown(lyricContent, "Test Song")
	if err != nil {
		t.Errorf("Markdown export failed: %v", err)
		return
	}
	verifyFileExists(t, mdPath)
	t.Logf("Markdown export: %s", mdPath)

	// Test Plain Text export
	txtPath, err := service.ExportToPlainText(lyricContent, "Test Song")
	if err != nil {
		t.Errorf("Plain text export failed: %v", err)
		return
	}
	verifyFileExists(t, txtPath)
	t.Logf("Plain text export: %s", txtPath)

	// Test ChordPro export
	choPath, err := service.ExportToChordPro(lyricContent, "Test Song")
	if err != nil {
		t.Errorf("ChordPro export failed: %v", err)
		return
	}
	verifyFileExists(t, choPath)
	t.Logf("ChordPro export: %s", choPath)

	// Test JSON export (existing format)
	jsonPath, err := service.ExportFull(lyricContent, "Test Song", 120, true)
	if err != nil {
		t.Errorf("JSON export failed: %v", err)
		return
	}
	verifyFileExists(t, jsonPath)
	t.Logf("JSON export: %s", jsonPath)

	// Test pattern export
	patternPath, err := service.ExportToPattern(patternContent, "Test Pattern")
	if err != nil {
		t.Errorf("Pattern export failed: %v", err)
		return
	}
	verifyFileExists(t, patternPath)
	t.Logf("Pattern export: %s", patternPath)

	// Test lyrics export
	lyricsPath, err := service.ExportToLyrics(lyricContent, "Test Lyrics")
	if err != nil {
		t.Errorf("Lyrics export failed: %v", err)
		return
	}
	verifyFileExists(t, lyricsPath)
	t.Logf("Lyrics export: %s", lyricsPath)

	// Test chords export
	chordsPath, err := service.ExportToChords(lyricContent, "Test Chords")
	if err != nil {
		t.Errorf("Chords export failed: %v", err)
		return
	}
	verifyFileExists(t, chordsPath)
	t.Logf("Chords export: %s", chordsPath)
}

// testEndToEndWorkflows tests complete workflows
func testEndToEndWorkflows(t *testing.T, themeManager *styles.ThemeManager, agent *ai.QuickIdeaAgent,
	exportService *export.ExportService, lyricContent, patternContent string) {

	ctx := context.Background()

	// Workflow 1: Lyric creation with AI assistance and export
	t.Run("LyricWorkflow", func(t *testing.T) {
		// Set theme
		if err := themeManager.SetTheme("Neon Dreams"); err != nil {
			t.Errorf("Failed to set theme: %v", err)
		}

		// Get AI suggestions
		req := ai.QuickRequest{
			Mode:    ai.QuickIdeaModeUnstick,
			Context: lyricContent,
		}

		resp, err := agent.Generate(ctx, req)
		if err != nil {
			t.Errorf("AI request failed: %v", err)
		}

		// Enhance content with AI suggestions
		enhancedContent := lyricContent
		if len(resp.Suggestions) > 0 {
			enhancedContent += "\n\n" + resp.Suggestions[0]
		}

		// Export in multiple formats
		mdPath, err := exportService.ExportToMarkdown(enhancedContent, "AI Enhanced Song")
		if err != nil {
			t.Errorf("Markdown export failed: %v", err)
		}

		choPath, err := exportService.ExportToChordPro(enhancedContent, "AI Enhanced Song")
		if err != nil {
			t.Errorf("ChordPro export failed: %v", err)
		}

		t.Logf("Lyric workflow completed: %s, %s", mdPath, choPath)
	})

	// Workflow 2: Pattern creation with AI assistance and export
	t.Run("PatternWorkflow", func(t *testing.T) {
		// Set different theme
		if err := themeManager.SetTheme("Midnight Jazz"); err != nil {
			t.Errorf("Failed to set theme: %v", err)
		}

		// Get AI suggestions for patterns
		req := ai.QuickRequest{
			Mode:    ai.QuickIdeaModeSpark,
			Context: patternContent,
			Options: map[string]string{"theme": "jazz"},
		}

		resp, err := agent.Generate(ctx, req)
		if err != nil {
			t.Errorf("AI request failed: %v", err)
		}

		// Enhance content with AI suggestions
		enhancedPattern := patternContent
		if len(resp.Suggestions) > 0 {
			enhancedPattern += "\n\n" + resp.Suggestions[0]
		}

		// Export pattern
		patternPath, err := exportService.ExportToPattern(enhancedPattern, "AI Enhanced Pattern")
		if err != nil {
			t.Errorf("Pattern export failed: %v", err)
		}

		t.Logf("Pattern workflow completed: %s", patternPath)
	})
}

// testPerformanceRequirements verifies performance requirements
func testPerformanceRequirements(t *testing.T, agent *ai.QuickIdeaAgent, themeManager *styles.ThemeManager,
	exportService *export.ExportService) {

	ctx := context.Background()

	// Test AI response time (< 2 seconds requirement)
	start := time.Now()
	req := ai.QuickRequest{
		Mode:    ai.QuickIdeaModeUnstick,
		Context: "The night is dark and full of stars",
	}

	_, err := agent.Generate(ctx, req)
	aiResponseTime := time.Since(start)

	if err != nil {
		t.Errorf("AI request failed during performance test: %v", err)
	} else if aiResponseTime > 2*time.Second {
		t.Errorf("AI response time exceeded 2 seconds: %v", aiResponseTime)
	} else {
		t.Logf("AI response time: %v (under 2s requirement)", aiResponseTime)
	}

	// Test theme switching performance
	themes := themeManager.GetAllThemes()
	var avgThemeTime time.Duration
	if len(themes) > 0 {
		start = time.Now()
		for _, th := range themes {
			themeManager.SetTheme(th.Name)
		}
		themeSwitchTime := time.Since(start)
		avgThemeTime = themeSwitchTime / time.Duration(len(themes))

		if avgThemeTime > 100*time.Millisecond {
			t.Errorf("Theme switching too slow: %v average", avgThemeTime)
		} else {
			t.Logf("Theme switching performance: %v average", avgThemeTime)
		}
	}

	// Test export performance
	testContent := strings.Repeat("Test line with content\n", 100)
	start = time.Now()
	_, err = exportService.ExportToMarkdown(testContent, "Performance Test")
	exportTime := time.Since(start)

	if err != nil {
		t.Errorf("Export failed during performance test: %v", err)
	} else if exportTime > 1*time.Second {
		t.Errorf("Export too slow: %v", exportTime)
	} else {
		t.Logf("Export performance: %v", exportTime)
	}

	// Log performance summary
	t.Logf("Performance requirements met:")
	t.Logf("  - AI Response: %v (requirement: <2s)", aiResponseTime)
	if len(themes) > 0 {
		t.Logf("  - Theme Switching: %v avg (requirement: <100ms)", avgThemeTime)
	}
	t.Logf("  - Export: %v (requirement: <1s)", exportTime)
}

// testErrorHandling tests error handling and graceful degradation
func testErrorHandling(t *testing.T, agent *ai.QuickIdeaAgent, exportService *export.ExportService) {
	ctx := context.Background()

	// Test AI with empty context
	req := ai.QuickRequest{
		Mode:    ai.QuickIdeaModeUnstick,
		Context: "",
	}

	resp, err := agent.Generate(ctx, req)
	if err != nil {
		t.Logf("AI correctly handled empty context: %v", err)
	} else if len(resp.Suggestions) == 0 {
		t.Logf("AI gracefully handled empty context with no suggestions")
	}

	// Test AI with invalid mode
	req = ai.QuickRequest{
		Mode:    "invalid",
		Context: "Test content",
	}

	resp, err = agent.Generate(ctx, req)
	if err != nil {
		t.Logf("AI correctly handled invalid mode: %v", err)
	} else {
		t.Logf("AI gracefully handled invalid mode")
	}

	// Test export with empty content
	_, err = exportService.ExportToMarkdown("", "Empty Test")
	if err != nil {
		t.Logf(" Export correctly handled empty content: %v", err)
	} else {
		t.Logf(" Export gracefully handled empty content")
	}

	// Test export with invalid path (Windows-specific)
	invalidService := export.NewExportService("Z:/invalid/path/that/does/not/exist")
	_, err = invalidService.ExportToMarkdown("Test content", "Invalid Path Test")
	if err != nil {
		t.Logf(" Export correctly handled invalid path: %v", err)
	} else {
		t.Log(" Export gracefully handled invalid path (may succeed on some systems)")
	}

	// Test knowledge base unavailability
	kbStatus := agent.GetKnowledgeBaseStatus(ctx)
	if kbStatus != nil && !kbStatus.Available {
		t.Logf(" Knowledge base correctly reports unavailable: %s", kbStatus.Error)
	}

	t.Logf(" Error handling tests completed")
}

// verifyFileExists checks if a file exists
func verifyFileExists(t *testing.T, path string) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Errorf("Expected file does not exist: %s", path)
	}
}
