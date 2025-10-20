package editor

import (
	"strings"
	"testing"

	"github.com/Kyanite/noise/internal/app"
	"github.com/Kyanite/noise/internal/theme"
	"github.com/Kyanite/noise/internal/ui/styles"
	"github.com/charmbracelet/bubbles/textarea"
)

func TestThemeIntegrationWithLyricFeatures(t *testing.T) {
	originalID := theme.GetManager().CurrentID()
	defer theme.ApplyThemeByID(originalID)

	themeIDs := getThemeCycleOrder()
	if len(themeIDs) == 0 {
		t.Fatal("expected at least one theme to be registered")
	}

	theoryService := app.NewTheoryService()

	for _, themeID := range themeIDs {
		theme.ApplyThemeByID(themeID)
		expectedTheme := theme.GetTheme(themeID)

		if styles.Primary != expectedTheme.Primary {
			t.Fatalf("expected primary color %s for theme %s, got %s", expectedTheme.Primary, themeID, styles.Primary)
		}

		rhymes, err := theoryService.FindRhymes("love")
		if err != nil {
			t.Fatalf("find rhymes failed for theme %s: %v", themeID, err)
		}
		if len(rhymes) == 0 {
			t.Fatalf("expected rhymes for 'love' under theme %s", themeID)
		}

		syllables, err := theoryService.CountSyllables("melody")
		if err != nil {
			t.Fatalf("count syllables failed for theme %s: %v", themeID, err)
		}
		if syllables == 0 {
			t.Fatalf("expected non-zero syllable count under theme %s", themeID)
		}

		textArea := textarea.New()
		editorModel := NewEditorPaneModel(textArea)
		chords := []string{"Cmaj7", "Am7", "Fmaj7", "G7"}

		editorModel.insertChords(chords)

		text := editorModel.GetText()
		for _, chord := range chords {
			if !strings.Contains(text, chord) {
				t.Fatalf("expected chord %s to be inserted under theme %s", chord, themeID)
			}
		}
	}
}

func TestThemeShortcutCycling(t *testing.T) {
	originalID := theme.GetManager().CurrentID()
	defer theme.ApplyThemeByID(originalID)

	themeIDs := getThemeCycleOrder()
	if len(themeIDs) < 2 {
		t.Skip("not enough themes registered to test cycling")
	}

	theme.ApplyThemeByID(themeIDs[0])

	model := NewSplitPaneModel(nil)
	defer model.Cleanup()

	if _, cmd := model.handleShortcutAction(ShortcutAction{Type: ActionNextTheme}); cmd != nil {
		t.Fatalf("unexpected command from ActionNextTheme: %#v", cmd)
	}

	expectedNext := themeIDs[1]
	if current := theme.GetManager().CurrentID(); current != expectedNext {
		t.Fatalf("expected theme %s after ActionNextTheme, got %s", expectedNext, current)
	}

	if _, cmd := model.handleShortcutAction(ShortcutAction{Type: ActionPreviousTheme}); cmd != nil {
		t.Fatalf("unexpected command from ActionPreviousTheme: %#v", cmd)
	}
	expectedPrev := themeIDs[0]
	if current := theme.GetManager().CurrentID(); current != expectedPrev {
		t.Fatalf("expected theme %s after ActionPreviousTheme, got %s", expectedPrev, current)
	}

	theme.ApplyThemeByID(themeIDs[len(themeIDs)-1])

	if _, cmd := model.handleShortcutAction(ShortcutAction{Type: ActionThemeCycle}); cmd != nil {
		t.Fatalf("unexpected command from ActionThemeCycle: %#v", cmd)
	}
	expectedCycle := themeIDs[0]
	if current := theme.GetManager().CurrentID(); current != expectedCycle {
		t.Fatalf("expected theme %s after ActionThemeCycle, got %s", expectedCycle, current)
	}
}
