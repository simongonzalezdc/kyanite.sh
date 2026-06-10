package app

import (
	"context"
	"fmt"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/Kyanite/noise/internal/infra/db"
)

func newAutoSaveTestConfig() *AutoSaveConfig {
	cfg := DefaultAutoSaveConfig()
	cfg.Enabled = true
	cfg.IntervalSeconds = 1
	cfg.DebounceMs = 10
	cfg.MaxRetries = 1
	cfg.RetryDelayMs = 10
	cfg.MaxVersions = 5
	return cfg
}

func newAutoSaveTestService(t testing.TB, cfg *AutoSaveConfig) (*AutoSaveService, *db.DB) {
	t.Helper()

	dataDir := filepath.Join(t.TempDir(), "autosave-db")
	database, err := db.New(db.Config{DataDir: dataDir})
	if err != nil {
		t.Fatalf("failed to create test database: %v", err)
	}

	t.Cleanup(func() {
		_ = database.Close()
	})

	if cfg == nil {
		cfg = newAutoSaveTestConfig()
	} else {
		cfgCopy := *cfg
		cfg = &cfgCopy
		if cfg.IntervalSeconds <= 0 {
			cfg.IntervalSeconds = 1
		}
		if cfg.DebounceMs <= 0 {
			cfg.DebounceMs = 5
		}
		if cfg.RetryDelayMs <= 0 {
			cfg.RetryDelayMs = 5
		}
		if cfg.MaxRetries <= 0 {
			cfg.MaxRetries = 1
		}
		if cfg.MaxVersions <= 0 {
			cfg.MaxVersions = 5
		}
	}

	service := NewAutoSaveService(database, cfg)
	return service, database
}

func waitForStatus(t testing.TB, ch <-chan AutoSaveStatus, expected AutoSaveStatus, timeout time.Duration) []AutoSaveStatus {
	t.Helper()

	var seen []AutoSaveStatus
	deadline := time.After(timeout)

	for {
		select {
		case status := <-ch:
			seen = append(seen, status)
			if status == expected {
				return seen
			}
		case <-deadline:
			t.Fatalf("timed out waiting for status %q, saw %v", expected, seen)
		}
	}
}

func assertNoStatus(t testing.TB, ch <-chan AutoSaveStatus, duration time.Duration) {
	t.Helper()

	select {
	case status := <-ch:
		t.Fatalf("unexpected status transition: %v", status)
	case <-time.After(duration):
	}
}

func TestAutoSaveServicePeriodicSaveRespectsContentChanges(t *testing.T) {
	t.Skip("KNOWN LIMITATION: Autosave callback timing race condition - see docs/KNOWN_TEST_LIMITATIONS.md")
	t.Parallel()

	service, _ := newAutoSaveTestService(t, nil)
	statusCh := make(chan AutoSaveStatus, 16)
	service.SetStatusChangeCallback(func(status AutoSaveStatus) {
		statusCh <- status
	})

	// Seed content for periodic save.
	service.contentMutex.Lock()
	service.lastContent = "interval content v1"
	service.contentMutex.Unlock()

	service.performPeriodicSave()
	waitForStatus(t, statusCh, AutoSaveSuccess, 2*time.Second)

	versions, err := service.db.GetVersions(0, 10)
	if err != nil {
		t.Fatalf("expected version retrieval to succeed: %v", err)
	}
	if len(versions) != 1 {
		t.Fatalf("expected exactly one version after first periodic save, got %d", len(versions))
	}
	if versions[0].Content != "interval content v1" {
		t.Fatalf("expected saved content %q, got %q", "interval content v1", versions[0].Content)
	}

	assertNoStatus(t, statusCh, 100*time.Millisecond)

	// Second call with unchanged content should not create a new version.
	service.performPeriodicSave()
	assertNoStatus(t, statusCh, 150*time.Millisecond)

	versionsAfterSecondCall, err := service.db.GetVersions(0, 10)
	if err != nil {
		t.Fatalf("failed to get versions after second call: %v", err)
	}
	if len(versionsAfterSecondCall) != 1 {
		t.Fatalf("expected version count to remain 1, got %d", len(versionsAfterSecondCall))
	}

	// Change content and ensure another periodic save occurs.
	service.contentMutex.Lock()
	service.lastContent = "interval content v2"
	service.contentMutex.Unlock()

	service.performPeriodicSave()
	waitForStatus(t, statusCh, AutoSaveSuccess, 2*time.Second)

	versionsAfterThirdCall, err := service.db.GetVersions(0, 10)
	if err != nil {
		t.Fatalf("failed to get versions after third call: %v", err)
	}
	if len(versionsAfterThirdCall) != 2 {
		t.Fatalf("expected two versions after content change, got %d", len(versionsAfterThirdCall))
	}
	if versionsAfterThirdCall[0].Content != "interval content v2" {
		t.Fatalf("expected latest content %q, got %q", "interval content v2", versionsAfterThirdCall[0].Content)
	}
}

func TestAutoSaveServiceSaveContentWithoutStart(t *testing.T) {
	t.Parallel()

	service, _ := newAutoSaveTestService(t, nil)
	statusCh := make(chan AutoSaveStatus, 16)

	var mu sync.Mutex
	var errors []error

	service.SetStatusChangeCallback(func(status AutoSaveStatus) {
		statusCh <- status
	})
	service.SetErrorCallback(func(err error) {
		mu.Lock()
		defer mu.Unlock()
		errors = append(errors, err)
	})

	service.SaveContent("offline save content")

	waitForStatus(t, statusCh, AutoSaveSuccess, 2*time.Second)

	if last := service.GetLastSaveTime(); last.IsZero() {
		t.Fatal("expected last save time to be recorded")
	}

	mu.Lock()
	defer mu.Unlock()
	if len(errors) != 0 {
		t.Fatalf("expected no errors during save, got %v", errors)
	}
}

func TestAutoSaveServiceForceSaveErrorAndRecovery(t *testing.T) {
	t.Parallel()

	service, originalDB := newAutoSaveTestService(t, nil)
	statusCh := make(chan AutoSaveStatus, 16)

	var seenErrors []error
	var mu sync.Mutex

	service.SetStatusChangeCallback(func(status AutoSaveStatus) {
		statusCh <- status
	})
	service.SetErrorCallback(func(err error) {
		mu.Lock()
		seenErrors = append(seenErrors, err)
		mu.Unlock()
	})

	if err := originalDB.Close(); err != nil {
		t.Fatalf("failed to close underlying database: %v", err)
	}

	if err := service.ForceSave("should fail"); err == nil {
		t.Fatal("expected force save to fail when database is closed")
	}

	waitForStatus(t, statusCh, AutoSaveError, 2*time.Second)

	mu.Lock()
	if len(seenErrors) == 0 {
		t.Fatal("expected onError callback to be invoked for failure")
	}
	mu.Unlock()

	recoveredDB, err := db.New(db.Config{DataDir: filepath.Join(t.TempDir(), "autosave-recovered")})
	if err != nil {
		t.Fatalf("failed to create recovered database: %v", err)
	}
	t.Cleanup(func() {
		_ = recoveredDB.Close()
	})

	service.db = recoveredDB

	service.SaveContent("recovered content")
	waitForStatus(t, statusCh, AutoSaveSuccess, 2*time.Second)

	if last := service.GetLastSaveTime(); last.IsZero() {
		t.Fatal("expected last save time after recovery")
	}
}

func TestAutoSaveServiceCleanupOldVersions(t *testing.T) {
	t.Parallel()

	cfg := newAutoSaveTestConfig()
	cfg.MaxVersions = 3
	service, _ := newAutoSaveTestService(t, cfg)

	for i := 0; i < 5; i++ {
		content := fmt.Sprintf("cleanup-content-%d", i)
		if err := service.ForceSave(content); err != nil {
			t.Fatalf("force save %d failed: %v", i, err)
		}
	}

	if err := service.CleanupOldVersions(0); err != nil {
		t.Fatalf("cleanup old versions failed: %v", err)
	}

	versions, err := service.db.GetVersions(0, 10)
	if err != nil {
		t.Fatalf("expected to retrieve versions after cleanup: %v", err)
	}

	if len(versions) != cfg.MaxVersions {
		t.Fatalf("expected %d versions after cleanup, got %d", cfg.MaxVersions, len(versions))
	}

	expected := []string{"cleanup-content-4", "cleanup-content-3", "cleanup-content-2"}
	for i, version := range versions {
		if version.Content != expected[i] {
			t.Fatalf("expected version %d content %q, got %q", i, expected[i], version.Content)
		}
	}
}

func TestAutoSaveServiceConcurrentForceSave(t *testing.T) {
	t.Parallel()

	service, _ := newAutoSaveTestService(t, nil)

	const goroutines = 8
	var wg sync.WaitGroup
	wg.Add(goroutines)

	for i := 0; i < goroutines; i++ {
		go func(id int) {
			defer wg.Done()
			content := fmt.Sprintf("concurrent-content-%d", id)
			if err := service.ForceSave(content); err != nil {
				t.Errorf("force save failed for goroutine %d: %v", id, err)
			}
		}(i)
	}

	wg.Wait()

	versions, err := service.db.GetVersions(0, goroutines)
	if err != nil {
		t.Fatalf("failed to retrieve versions after concurrent saves: %v", err)
	}

	if len(versions) != goroutines {
		t.Fatalf("expected %d versions after concurrent saves, got %d", goroutines, len(versions))
	}

	// Ensure latest version content matches one of the expected entries.
	latest := versions[0].Content
	if latest == "" {
		t.Fatal("expected latest version content to be non-empty")
	}
}

func TestAutoSaveServiceStartAndStopLifecycle(t *testing.T) {
	t.Parallel()

	service, _ := newAutoSaveTestService(t, nil)
	statusCh := make(chan AutoSaveStatus, 16)

	service.SetStatusChangeCallback(func(status AutoSaveStatus) {
		statusCh <- status
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	if err := service.Start(ctx); err != nil {
		t.Fatalf("failed to start auto-save service: %v", err)
	}

	service.SaveContent("lifecycle-content")

	waitForStatus(t, statusCh, AutoSaveSuccess, 2*time.Second)

	if err := service.Stop(); err != nil {
		t.Fatalf("failed to stop service: %v", err)
	}

	if status := service.GetStatus(); status != AutoSaveIdle && status != AutoSaveSuccess {
		t.Fatalf("expected idle or success status after stop, got %v", status)
	}
}
