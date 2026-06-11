package errors

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/kyanite/noise/internal/domain"
	logging "github.com/kyanite/noise/internal/logging"
)

// testLogger wraps a real Logger that writes to a buffer for assertions.
type testLogger struct {
	*logging.Logger
	buf *bytes.Buffer
	mu  sync.Mutex
}

// NewTestLogger creates a test logger that captures log output for assertions.
func NewTestLogger(t *testing.T) *testLogger {
	t.Helper()
	var buf bytes.Buffer
	cfg := logging.DefaultConfig()
	cfg.Level = logging.DEBUG
	cfg.Output = &buf
	logger, err := logging.New(cfg)
	if err != nil {
		t.Fatalf("failed to create test logger: %v", err)
	}
	return &testLogger{Logger: logger, buf: &buf}
}

// ContainsMessage reports whether any captured log output contains substr.
func (tl *testLogger) ContainsMessage(substr string) bool {
	tl.mu.Lock()
	defer tl.mu.Unlock()
	return strings.Contains(tl.buf.String(), substr)
}

// GetMessages returns all captured log output split into lines.
func (tl *testLogger) GetMessages() []string {
	tl.mu.Lock()
	defer tl.mu.Unlock()
	lines := strings.Split(tl.buf.String(), "\n")
	// Filter empty trailing line
	out := make([]string, 0, len(lines))
	for _, l := range lines {
		if l != "" {
			out = append(out, l)
		}
	}
	return out
}

// Clear resets the captured output.
func (tl *testLogger) Clear() {
	tl.mu.Lock()
	defer tl.mu.Unlock()
	tl.buf.Reset()
}

// TestErrorReporter captures error reports for test assertions.
type TestErrorReporter struct {
	mu      sync.Mutex
	reports []ErrorReport
}

// NewTestErrorReporter creates a reporter that captures reports in memory.
func NewTestErrorReporter() *TestErrorReporter {
	return &TestErrorReporter{}
}

// Name implements the ErrorReporter interface.
func (ter *TestErrorReporter) Name() string { return "test" }

// Report implements the ErrorReporter interface.
func (ter *TestErrorReporter) Report(ctx context.Context, report *ErrorReport) error {
	ter.mu.Lock()
	defer ter.mu.Unlock()
	ter.reports = append(ter.reports, *report)
	return nil
}

// GetReports returns all captured reports.
func (ter *TestErrorReporter) GetReports() []ErrorReport {
	ter.mu.Lock()
	defer ter.mu.Unlock()
	out := make([]ErrorReport, len(ter.reports))
	copy(out, ter.reports)
	return out
}

// Clear removes all captured reports.
func (ter *TestErrorReporter) Clear() {
	ter.mu.Lock()
	defer ter.mu.Unlock()
	ter.reports = nil
}

// testSetup holds common test dependencies.
type testSetup struct {
	Logger   *testLogger
	Manager  *ErrorManager
	Reporter *TestErrorReporter
}

// NewTestSetup creates a standard test setup with manager and reporter.
func NewTestSetup(t *testing.T) *testSetup {
	t.Helper()
	logger := NewTestLogger(t)
	reporter := NewTestErrorReporter()
	config := DefaultErrorConfig()
	config.EnableReporting = true
	manager := NewErrorManager(logger.Logger, config)
	manager.AddReporter(reporter)
	return &testSetup{
		Logger:   logger,
		Manager:  manager,
		Reporter: reporter,
	}
}

// Cleanup is a no-op kept for backward compatibility with existing tests.
func (ts *testSetup) Cleanup() {}

// CreateTestFile creates a file with the given content in dir and returns its path.
func CreateTestFile(t *testing.T, dir, name, content string) string {
	t.Helper()
	path := filepath.Join(dir, name)
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatalf("failed to create test file %s: %v", name, err)
	}
	return path
}

// CreateCorruptedTestFile creates a file with intentionally corrupted JSON-like content.
func CreateCorruptedTestFile(t *testing.T, dir, name string) string {
	t.Helper()
	return CreateTestFile(t, dir, name, `{"invalid": json content`)
}

// CreateTestSong creates a minimal test song with the given ID and title.
func CreateTestSong(id int, title string) *domain.Song {
	return &domain.Song{
		ID:       id,
		Filepath: filepath.Join(os.TempDir(), title+".song"),
		Metadata: domain.SongMetadata{
			Title: title,
		},
	}
}

// relaxPerfBudgets returns true in CI or slow environments to skip strict timing checks.
func relaxPerfBudgets() bool {
	return os.Getenv("CI") != "" || os.Getenv("NOISE_RELAX_PERF") != ""
}
