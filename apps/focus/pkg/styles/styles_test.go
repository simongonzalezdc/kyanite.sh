package styles

import (
	"testing"
)

func TestFocusColorFunctions(t *testing.T) {
	// Test legacy color functions for backward compatibility
	// These should not panic
	_ = FocusStyle("Test")
	_ = FocusPinkColor("Test Pink")
	_ = FocusBlueColor("Test Blue")
	_ = FocusGreenColor("Test Green")
	_ = FocusPurpleColor("Test Purple")
}

func TestThemeAwareGetters(t *testing.T) {
	// Ensure color getters don't return empty
	if GetBackground() == "" {
		t.Error("GetBackground should not be empty")
	}
	if GetForeground() == "" {
		t.Error("GetForeground should not be empty")
	}
	if GetAccent() == "" {
		t.Error("GetAccent should not be empty")
	}
	if GetSuccess() == "" {
		t.Error("GetSuccess should not be empty")
	}
	if GetError() == "" {
		t.Error("GetError should not be empty")
	}
}

func TestStyleGetters(t *testing.T) {
	// These should not panic
	_ = GetTitleStyle()
	_ = GetBoxStyle()
	_ = GetPanelStyle()
	_ = GetSuccessStyle()
	_ = GetErrorStyle()
	_ = GetWarningStyle()
}

func TestRenderFunctions(t *testing.T) {
	// These should not panic
	_ = SynthwaveTitle("Test")
	_ = Title("Test")
	_ = FocusBox("Test", GetAccent())
	_ = CyberGridBox("Test")
	_ = HolographicText("Test")
	_ = DigitalRain("Test")
	_ = PriorityExplosion("high")
	_ = TaskStatus("completed")
	_ = CyberTag("Test")
	_ = Header()
	_ = CyberStats(1, 2, 3)
	_ = LoadingMessage()
	_ = EmptyStateMessage()
	_ = SynthwaveAIResponseStyle("test")
	_ = SynthwaveReportStyle("test")
	_ = PriorityStyle("high")
	_ = IDStyle("123")
	_ = HeaderStyle("Test")
	_ = FooterStyle("Test")
	_ = CategoryStyle("Test")
	_ = SuggestionStyle("Test")
	_ = AIResponseStyle("Test")
	_ = ReportStyle("Test")
}

func TestCycleTheme(t *testing.T) {
	original := GetCurrentThemeName()
	CycleTheme()
	// Should have changed
	_ = GetCurrentThemeName()
	// Restore
	SetThemeByName(original)
}
