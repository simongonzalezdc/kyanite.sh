package noise

import (
	"context"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/Kyanite/noise/internal/app"
	"github.com/Kyanite/noise/internal/infra/db"
	"github.com/Kyanite/noise/internal/ui/editor"
)

// TestEditorWorkflow verifies basic editor lifecycle and view synchronization.
func TestEditorWorkflow(t *testing.T) {
	tempDir := t.TempDir()
	database, err := db.New(db.Config{DataDir: tempDir})
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	defer database.Close()

	autoCfg := app.DefaultAutoSaveConfig()
	autoCfg.IntervalSeconds = 1
	autoSvc := app.NewAutoSaveService(database, autoCfg)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = autoSvc.Start(ctx) // ignore error in test environment
	if err := autoSvc.Stop(); err != nil {
		t.Logf("Warning: Failed to stop auto-save service: %v", err)
	}

	sp := editor.NewSplitPaneModel(database)
	// send a window size to initialize dimensions
	sp.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	content := "# Test Song\n\nVerse 1\nLine A\nLine B\n\n[Chorus]\nSing along"
	sp.SetEditorText(content)
	time.Sleep(100 * time.Millisecond)

	view := sp.View()
	if view == "" {
		t.Fatal("expected non-empty view")
	}

	if !strings.Contains(view, "Test Song") {
		t.Fatalf("expected view to contain title; got: %q", view)
	}

	// simulate tab to switch focus
	sp.Update(tea.KeyMsg{Type: tea.KeyTab, Runes: []rune{'\t'}})
	// simulate esc key
	sp.Update(tea.KeyMsg{Type: tea.KeyEsc, Runes: []rune{'\x1b'}})

	sp.Cleanup()
}

// TestPreviewSync ensures editor updates appear in preview portion of the view.
func TestPreviewSync(t *testing.T) {
	tempDir := t.TempDir()
	database, err := db.New(db.Config{DataDir: tempDir})
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	defer database.Close()

	sp := editor.NewSplitPaneModel(database)
	sp.Update(tea.WindowSizeMsg{Width: 100, Height: 30})

	md := "# Sync Test\n\n[Verse]\nHello\nWorld"
	sp.SetEditorText(md)
	time.Sleep(100 * time.Millisecond)

	view := sp.View()
	if !strings.Contains(view, "Sync Test") || !strings.Contains(view, "Hello") {
		t.Fatalf("expected view to contain updated markdown; got: %q", view)
	}

	updated := md + "\n\n[Bridge]\nBridge text"
	sp.SetEditorText(updated)
	time.Sleep(100 * time.Millisecond)
	view2 := sp.View()
	if !strings.Contains(view2, "Bridge text") {
		t.Fatalf("expected view to reflect updated content; got: %q", view2)
	}

	sp.Cleanup()
}

// TestDatabaseAndAutoSave validates basic DB versioning and autosave recovery.
func TestDatabaseAndAutoSave(t *testing.T) {
	tempDir := t.TempDir()
	database, err := db.New(db.Config{DataDir: tempDir})
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	defer database.Close()

	editorSvc := app.NewEditorService(database, database)
	song, err := editorSvc.CreateSong("DB Song", "Tester")
	if err != nil {
		t.Fatalf("failed to create song: %v", err)
	}

	_, err = database.SaveVersion(song.ID, "v1", false, "v1")
	if err != nil {
		t.Fatalf("failed to save version v1: %v", err)
	}
	_, err = database.SaveVersion(song.ID, "v2", true, "milestone")
	if err != nil {
		t.Fatalf("failed to save version v2: %v", err)
	}

	versions, err := editorSvc.GetVersions(song.ID, 10)
	if err != nil {
		t.Fatalf("failed to list versions: %v", err)
	}
	if len(versions) < 2 {
		t.Fatalf("expected >=2 versions, got %d", len(versions))
	}

	autoCfg := app.DefaultAutoSaveConfig()
	autoCfg.IntervalSeconds = 1
	autoSvc := app.NewAutoSaveService(database, autoCfg)
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = autoSvc.Start(ctx)
	if err := autoSvc.Stop(); err != nil {
		t.Logf("Warning: Failed to stop auto-save service: %v", err)
	}

	_, err = database.SaveVersion(song.ID, "autosave content", false, "auto-save")
	if err != nil {
		t.Fatalf("failed to save autosave version: %v", err)
	}

	rec, err := autoSvc.RecoverFromLastSave(song.ID)
	if err != nil {
		t.Fatalf("recover failed: %v", err)
	}
	if !strings.Contains(rec, "autosave content") {
		t.Fatalf("unexpected recovered content: %q", rec)
	}
}
