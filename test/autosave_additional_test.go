package noise

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"github.com/Kyanite/noise/internal/app"
	"github.com/Kyanite/noise/internal/infra/db"
)

// helper to create test DB in a temp dir
func newTestDB(t *testing.T) *db.DB {
	t.Helper()
	tmp := t.TempDir()
	dir := filepath.Join(tmp, "db")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("failed to create db dir: %v", err)
	}
	database, err := db.New(db.Config{DataDir: dir})
	if err != nil {
		t.Fatalf("failed to create db: %v", err)
	}
	return database
}

// TestLoadAndSaveConfig ensures the loaders/save helpers are callable and sane.
func TestLoadAndSaveConfig(t *testing.T) {
	cfg, err := app.LoadAutoSaveConfigFromFile("nonexistent.json")
	if err != nil {
		t.Fatalf("LoadAutoSaveConfigFromFile returned error: %v", err)
	}
	if cfg == nil {
		t.Fatalf("expected non-nil config from loader")
	}

	if err := app.SaveAutoSaveConfigToFile(cfg, "noop.json"); err != nil {
		t.Fatalf("SaveAutoSaveConfigToFile returned error: %v", err)
	}
}

// TestSaveContentTriggersCallbacksWhenNotStarted verifies that SaveContent
// triggers asynchronous saving and that status callbacks are invoked even when
// the service hasn't been started.
func TestSaveContentTriggersCallbacksWhenNotStarted(t *testing.T) {
	database := newTestDB(t)
	defer database.Close()

	// Use short retry/delay to keep test fast
	cfg := app.DefaultAutoSaveConfig()
	cfg.DebounceMs = 50
	cfg.RetryDelayMs = 10
	cfg.MaxRetries = 1

	service := app.NewAutoSaveService(database, cfg)

	statusCh := make(chan app.AutoSaveStatus, 10)
	errCh := make(chan error, 10)

	service.SetStatusChangeCallback(func(s app.AutoSaveStatus) {
		statusCh <- s
	})
	service.SetErrorCallback(func(err error) {
		errCh <- err
	})

	service.SaveContent("callback-test-content")

	// Wait for up to 2s for success status (allowing for internal small sleeps)
	timeout := time.After(2 * time.Second)
	gotSuccess := false
loop:
	for {
		select {
		case st := <-statusCh:
			if st == app.AutoSaveSuccess {
				gotSuccess = true
				break loop
			}
		case e := <-errCh:
			t.Fatalf("unexpected error callback: %v", e)
		case <-timeout:
			break loop
		}
	}

	if !gotSuccess {
		t.Fatalf("did not observe AutoSaveSuccess status; statuses may be: %v", drainStatusChannel(statusCh))
	}
}

// drainStatusChannel drains remaining statuses for debugging.
func drainStatusChannel(ch chan app.AutoSaveStatus) []app.AutoSaveStatus {
	var out []app.AutoSaveStatus
	for {
		select {
		case s := <-ch:
			out = append(out, s)
		default:
			return out
		}
	}
}

// TestRecoverFromLastSaveNoVersions asserts that RecoverFromLastSave returns an error
// when there are no versions available for the given song ID.
func TestRecoverFromLastSaveNoVersions(t *testing.T) {
	database := newTestDB(t)
	defer database.Close()

	service := app.NewAutoSaveService(database, nil)

	_, err := service.RecoverFromLastSave(9999) // unlikely to exist
	if err == nil {
		t.Fatalf("expected error recovering when no versions exist")
	}
}

// TestGetMilestonesAndVersioning ensures that milestone creation is persisted and
// that GetMilestones filters milestone versions correctly.
func TestGetMilestonesAndVersioning(t *testing.T) {
	database := newTestDB(t)
	defer database.Close()

	service := app.NewAutoSaveService(database, nil)

	// Create several versions, some milestones and some not
	if err := service.SaveWithVersioning(42, "auto-1", false, ""); err != nil {
		t.Fatalf("SaveWithVersioning failed: %v", err)
	}
	if err := service.SaveWithVersioning(42, "mile-1", true, "My Milestone"); err != nil {
		t.Fatalf("SaveWithVersioning (milestone) failed: %v", err)
	}
	if err := service.SaveWithVersioning(42, "auto-2", false, ""); err != nil {
		t.Fatalf("SaveWithVersioning failed: %v", err)
	}
	if err := service.SaveWithVersioning(42, "mile-2", true, "Another Milestone"); err != nil {
		t.Fatalf("SaveWithVersioning (milestone) failed: %v", err)
	}

	milestones, err := service.GetMilestones(42)
	if err != nil {
		t.Fatalf("GetMilestones returned error: %v", err)
	}
	if len(milestones) < 2 {
		t.Fatalf("expected at least 2 milestones, got %d", len(milestones))
	}

	// Verify milestone names and IsMilestone flag
	found := map[string]bool{}
	for _, m := range milestones {
		if !m.IsMilestone {
			t.Fatalf("expected version %d to be a milestone", m.ID)
		}
		found[m.MilestoneName] = true
	}
	if !found["My Milestone"] || !found["Another Milestone"] {
		t.Fatalf("expected created milestone names to be present; found %v", found)
	}
}

// TestCleanupOldVersions ensures that CleanupOldVersions deletes older versions
// when the total exceeds the configured MaxVersions.
func TestCleanupOldVersions(t *testing.T) {
	database := newTestDB(t)
	defer database.Close()

	cfg := app.DefaultAutoSaveConfig()
	cfg.MaxVersions = 3
	service := app.NewAutoSaveService(database, cfg)

	// Create 7 versions for songID 0 (auto-save general)
	for i := 0; i < 7; i++ {
		if err := service.SaveWithVersioning(0, "v"+time.Now().Format("150405.000")+strconv.Itoa(i), false, ""); err != nil {
			t.Fatalf("failed to create version %d: %v", i, err)
		}
		// tiny sleep to ensure ordering timestamps differ
		time.Sleep(5 * time.Millisecond)
	}

	// Ensure we have more than MaxVersions prior to cleanup
	before, err := service.GetVersionHistory(0, 100)
	if err != nil {
		t.Fatalf("GetVersionHistory failed: %v", err)
	}
	if len(before) <= cfg.MaxVersions {
		t.Fatalf("setup failed; expected more than %d versions, got %d", cfg.MaxVersions, len(before))
	}

	if err := service.CleanupOldVersions(0); err != nil {
		t.Fatalf("CleanupOldVersions failed: %v", err)
	}

	after, err := service.GetVersionHistory(0, 100)
	if err != nil {
		t.Fatalf("GetVersionHistory failed: %v", err)
	}
	if len(after) > cfg.MaxVersions {
		t.Fatalf("expected at most %d versions after cleanup, got %d", cfg.MaxVersions, len(after))
	}
}

// TestPeriodicStartStop exercises Start/Stop with a short interval to ensure the
// periodic timer goroutine runs and does not leak when stopped.
func TestPeriodicStartStop(t *testing.T) {
	database := newTestDB(t)
	defer database.Close()

	cfg := app.DefaultAutoSaveConfig()
	cfg.IntervalSeconds = 1
	cfg.DebounceMs = 10

	service := app.NewAutoSaveService(database, cfg)

	ctx, cancel := context.WithCancel(context.Background())
	if err := service.Start(ctx); err != nil {
		t.Fatalf("failed to start service: %v", err)
	}

	// set content so periodic save has something to act on
	service.SaveContent("periodic-test-content")

	// let it run a short while
	time.Sleep(1200 * time.Millisecond)

	if err := service.Stop(); err != nil {
		t.Fatalf("failed to stop service: %v", err)
	}
	cancel()

	// Allow small time for shutdown
	time.Sleep(50 * time.Millisecond)
}

// TestGetSaveStatisticsEmpty ensures statistics handle empty state gracefully.
func TestGetSaveStatisticsEmpty(t *testing.T) {
	database := newTestDB(t)
	defer database.Close()

	service := app.NewAutoSaveService(database, nil)
	stats, err := service.GetSaveStatistics(9999)
	if err != nil {
		t.Fatalf("GetSaveStatistics returned error: %v", err)
	}
	if stats == nil {
		t.Fatalf("expected stats object even when there are no versions")
	}
}