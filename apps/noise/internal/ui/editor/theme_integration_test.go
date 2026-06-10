package editor

import (
	"sort"
	"strings"
	"testing"

	"github.com/Kyanite/noise/internal/app"
	"github.com/Kyanite/noise/internal/config"
	"github.com/Kyanite/noise/internal/theme"
	"github.com/charmbracelet/bubbles/textarea"
)

func TestThemeIntegrationWithLyricFeatures(t *testing.T) {
	themeManager := theme.GetManager()
	originalTheme := themeManager.Current()

	// Restore original theme after test
	defer func() {
		// Find the ID of the original theme
		for id, th := range theme.Registry {
			if th.Name == originalTheme.Name {
				themeManager.SetTheme(id)
				break
			}
		}
	}()

	themeIDs := theme.ListThemes()
	if len(themeIDs) == 0 {
		t.Fatal("expected at least one theme to be registered")
	}

	theoryService := app.NewTheoryService()

	for _, themeID := range themeIDs {
		themeManager.SetTheme(themeID)

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
		editorModel := NewEditorPaneModel(textArea, nil)
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
	themeManager := theme.GetManager()
	originalTheme := themeManager.Current()

	// Restore original theme after test
	defer func() {
		// Find the ID of the original theme
		for id, th := range theme.Registry {
			if th.Name == originalTheme.Name {
				themeManager.SetTheme(id)
				break
			}
		}
	}()

	themeIDs := theme.ListThemes()
	sort.Strings(themeIDs)
	if len(themeIDs) < 2 {
		t.Skip("not enough themes registered to test cycling")
	}

	themeManager.SetTheme(themeIDs[0])

	aiService := app.NewAIService(config.DefaultConfig())
	model := NewSplitPaneModel(nil, aiService)
	defer model.Cleanup()

	if _, cmd := model.handleShortcutAction(ShortcutAction{Type: ActionNextTheme}); cmd != nil {
		t.Fatalf("unexpected command from ActionNextTheme: %#v", cmd)
	}

	expectedNext := themeIDs[1]
	currentTheme := themeManager.Current()
	// Find the ID of the current theme
	var currentID string
	for id, th := range theme.Registry {
		if th.Name == currentTheme.Name {
			currentID = id
			break
		}
	}
	if currentID != expectedNext {
		t.Fatalf("expected theme %s after ActionNextTheme, got %s", expectedNext, currentID)
	}

	if _, cmd := model.handleShortcutAction(ShortcutAction{Type: ActionPreviousTheme}); cmd != nil {
		t.Fatalf("unexpected command from ActionPreviousTheme: %#v", cmd)
	}
	expectedPrev := themeIDs[0]
	currentTheme = themeManager.Current()
	// Find the ID of the current theme
	currentID = ""
	for id, th := range theme.Registry {
		if th.Name == currentTheme.Name {
			currentID = id
			break
		}
	}
	if currentID != expectedPrev {
		t.Fatalf("expected theme %s after ActionPreviousTheme, got %s", expectedPrev, currentID)
	}

	themeManager.SetTheme(themeIDs[len(themeIDs)-1])

	if _, cmd := model.handleShortcutAction(ShortcutAction{Type: ActionThemeCycle}); cmd != nil {
		t.Fatalf("unexpected command from ActionThemeCycle: %#v", cmd)
	}
	expectedCycle := themeIDs[0]
	currentTheme = themeManager.Current()
	// Find the ID of the current theme
	currentID = ""
	for id, th := range theme.Registry {
		if th.Name == currentTheme.Name {
			currentID = id
			break
		}
	}
	if currentID != expectedCycle {
		t.Fatalf("expected theme %s after ActionThemeCycle, got %s", expectedCycle, currentID)
	}
}
