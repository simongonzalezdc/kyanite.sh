package noise

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	appErrors "github.com/Kyanite/noise/internal/errors"
	"github.com/Kyanite/noise/internal/logging"
)

func newTestLogger(t *testing.T) *logging.Logger {
	t.Helper()
	cfg := &logging.Config{
		Level:      logging.DEBUG,
		Output:     io.Discard,
		ShowCaller: false,
		LogFile:    "",
	}
	l, err := logging.New(cfg)
	if err != nil {
		t.Fatalf("failed to create test logger: %v", err)
	}
	return l
}

type mockReporter struct {
	last      *appErrors.ErrorReport
	reportErr error
	name      string
	called    bool
}

func (m *mockReporter) Report(ctx context.Context, report *appErrors.ErrorReport) error {
	m.called = true
	m.last = report
	if m.reportErr != nil {
		return m.reportErr
	}
	return nil
}
func (m *mockReporter) Name() string { return m.name }

func TestErrorManagerInitializationAndReporting(t *testing.T) {
	logger := newTestLogger(t)

	config := &appErrors.ErrorConfig{
		EnableReporting: true,
		LogToFile:       false,
		ReportFilters:   []string{},
		DefaultRetries:  2,
		RetryDelay:      10 * time.Millisecond,
		MaxRecoveryTime: 200 * time.Millisecond,
		NotifyUser:      false,
		ShowStackTrace:  false,
	}

	em := appErrors.NewErrorManager(logger, config)
	if em == nil {
		t.Fatal("expected ErrorManager instance")
	}

	mock := &mockReporter{name: "mock"}
	em.AddReporter(mock)

	report := em.HandleError(context.Background(), errors.New("database connection failed"))
	if report == nil {
		t.Fatal("expected non-nil report")
	}

	// Because EnableReporting == true and message contains "database" the reporter should be called
	if !mock.called {
		t.Fatal("expected reporter to be called")
	}

	// Validate report contents
	if report.Error == nil {
		t.Fatal("expected AppError in report")
	}
	if !report.Error.HasCategory(appErrors.CategoryDatabase) && report.Error.Category != appErrors.CategoryDatabase {
		t.Fatalf("expected report category to be database, got %s", report.Error.Category)
	}
}

func TestExternalErrorReporterHTTPIntegration(t *testing.T) {
	logger := newTestLogger(t)

	var receivedBody []byte
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var err error
		receivedBody, err = io.ReadAll(r.Body)
		if err != nil {
			t.Fatalf("failed read body: %v", err)
		}
		w.WriteHeader(200)
	}))
	defer ts.Close()

	reporter := appErrors.NewExternalErrorReporter(ts.URL, "dummy-key", logger)

	// Only high/critical severities are sent. Use High.
	appErr := appErrors.NewAppError("E1", "external test", nil, appErrors.CategoryUnknown, appErrors.SeverityHigh, appErrors.RecoveryNone)
	report := &appErrors.ErrorReport{Error: appErr, Handled: true, Timestamp: time.Now()}

	if err := reporter.Report(context.Background(), report); err != nil {
		t.Fatalf("expected successful report, got: %v", err)
	}

	// Validate payload contains the error code
	if !bytes.Contains(receivedBody, []byte("E1")) {
		t.Fatalf("expected payload to contain error code E1, got: %s", string(receivedBody))
	}

	// Test failing response
	ts2 := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		if _, err := w.Write([]byte("server error")); err != nil {
			panic(err)
		}
	}))
	defer ts2.Close()

	reporterBad := appErrors.NewExternalErrorReporter(ts2.URL, "k", logger)
	err := reporterBad.Report(context.Background(), report)
	if err == nil {
		t.Fatal("expected error when external service returns non-2xx")
	}
}

func TestHandleErrorWithRecovery_Succeeds(t *testing.T) {
	logger := newTestLogger(t)

	config := &appErrors.ErrorConfig{
		EnableReporting: false,
		LogToFile:       false,
		DefaultRetries:  3,
		RetryDelay:      10 * time.Millisecond,
		MaxRecoveryTime: 500 * time.Millisecond,
	}
	em := appErrors.NewErrorManager(logger, config)

	// Create an underlying error that categorizes to network (has "network")
	err := errors.New("network timeout")

	// Recovery function will succeed on second attempt
	attempts := 0
	rec := func(ctx context.Context, e error) error {
		attempts++
		if attempts < 2 {
			return errors.New("transient")
		}
		return nil
	}

	report := em.HandleErrorWithRecovery(context.Background(), err, rec)
	if report == nil {
		t.Fatal("expected report")
	}
	// Should have attempted recovery and marked recovered
	if !report.Recovered {
		t.Fatalf("expected recovered true, got false; recovery: %+v", report.Recovery)
	}
	if report.Recovery == nil || report.Recovery.Attempts < 1 {
		t.Fatalf("expected recovery attempts recorded, got %+v", report.Recovery)
	}
}

func TestAttemptRecoveryTimeout(t *testing.T) {
	logger := newTestLogger(t)

	config := &appErrors.ErrorConfig{
		EnableReporting: false,
		LogToFile:       false,
		DefaultRetries:  5,
		RetryDelay:      10 * time.Millisecond,
		MaxRecoveryTime: 20 * time.Millisecond, // force timeout
	}
	em := appErrors.NewErrorManager(logger, config)

	err := errors.New("network issue")
	// Never succeeds
	rec := func(ctx context.Context, e error) error {
		time.Sleep(50 * time.Millisecond)
		return errors.New("still failing")
	}

	report := em.HandleErrorWithRecovery(context.Background(), err, rec)
	if report == nil {
		t.Fatal("expected report")
	}
	// Recovery should not be successful due to timeout
	if report.Recovered {
		t.Fatal("expected recovery to fail due to timeout")
	}
	if report.Recovery == nil || report.Recovery.Result == "" {
		t.Fatalf("expected recovery result to be set, got %+v", report.Recovery)
	}
}

func TestFileCorruptionDetectorDetectAndRecoverJSON(t *testing.T) {
	logger := newTestLogger(t)
	dir := t.TempDir()
	orig := filepath.Join(dir, "test.json")
	backup := orig + ".backup"

	// Write corrupted JSON
	if err := os.WriteFile(orig, []byte("{invalid json:"), 0o644); err != nil {
		t.Fatalf("write corrupted: %v", err)
	}
	// Write valid backup
	valid := map[string]interface{}{"hello": "world"}
	bs, _ := json.Marshal(valid)
	if err := os.WriteFile(backup, bs, 0o644); err != nil {
		t.Fatalf("write backup: %v", err)
	}

	fcd := appErrors.NewFileCorruptionDetector(logger)
	if err := fcd.DetectCorruption(orig); err == nil {
		t.Fatal("expected DetectCorruption to fail for corrupted json")
	}

	// RecoverFile should find .backup and restore
	if err := fcd.RecoverFile(orig); err != nil {
		t.Fatalf("expected RecoverFile to succeed, got %v", err)
	}

	// After recovery file should be valid JSON
	content, err := os.ReadFile(orig)
	if err != nil {
		t.Fatalf("read recovered: %v", err)
	}
	var out map[string]interface{}
	if err := json.Unmarshal(content, &out); err != nil {
		t.Fatalf("expected valid json after recovery, got %v", err)
	}
	if out["hello"] != "world" {
		t.Fatalf("unexpected recovered content: %v", out)
	}
}

func TestFileCorruptionDetectorNoBackup(t *testing.T) {
	logger := newTestLogger(t)
	dir := t.TempDir()
	orig := filepath.Join(dir, "no_backup.json")

	// Write corrupted JSON
	if err := os.WriteFile(orig, []byte("{invalid json:"), 0o644); err != nil {
		t.Fatalf("write corrupted: %v", err)
	}

	fcd := appErrors.NewFileCorruptionDetector(logger)
	err := fcd.RecoverFile(orig)
	if err == nil {
		t.Fatal("expected RecoverFile to fail when no backup exists")
	}
}

func TestNotificationManagerAndUIIntegration(t *testing.T) {
	logger := newTestLogger(t)
	nm := appErrors.NewNotificationManager(*logger)

	// Show error notification
	appErr := appErrors.NewAppError("UIERR", "UI fail", nil, appErrors.CategoryUI, appErrors.SeverityMedium, appErrors.RecoveryGraceful)
	id := nm.ShowError("Test Error", "something bad", appErr)
	if id == "" {
		t.Fatal("expected non-empty notification id")
	}

	active := nm.GetActiveNotifications()
	if len(active) == 0 {
		t.Fatal("expected at least one active notification")
	}

	// Dismiss it and ensure it's no longer active
	nm.DismissNotification(id)
	active = nm.GetActiveNotifications()
	if len(active) != 0 {
		t.Fatalf("expected zero active notifications after dismiss, got %d", len(active))
	}

	// Test UI model rendering shows details when toggled
	nm2 := appErrors.NewNotificationManager(*logger)
	_ = nm2.ShowError("Detailed", "with details", appErr)
	ui := appErrors.NewNotificationUIModel(nm2)

	view := ui.View()
	if view == "" {
		t.Fatal("expected UI view to be non-empty")
	}
}

func TestFileAndConsoleReportersAndComposite(t *testing.T) {
	logger := newTestLogger(t)
	dir := t.TempDir()
	filePath := filepath.Join(dir, "errors.log")

	// File reporter
	frep, err := appErrors.NewFileErrorReporter(filePath, logger)
	if err != nil {
		t.Fatalf("failed create file reporter: %v", err)
	}

	appErr := appErrors.NewAppError("F1", "file report", nil, appErrors.CategoryUnknown, appErrors.SeverityMedium, appErrors.RecoveryNone)
	report := &appErrors.ErrorReport{Error: appErr, Handled: true, Timestamp: time.Now()}

	if err := frep.Report(context.Background(), report); err != nil {
		t.Fatalf("expected file report to succeed, got %v", err)
	}
	// Check contents in file
	if err := frep.Close(); err != nil {
		t.Fatalf("failed to close file reporter: %v", err)
	}
	data, err := os.ReadFile(filePath)
	if err != nil {
		t.Fatalf("failed read log file: %v", err)
	}
	if !bytes.Contains(data, []byte("Error Report")) {
		t.Fatalf("expected log file to contain Error Report, got %s", string(data))
	}

	// Console reporter
	var buf bytes.Buffer
	crep := appErrors.NewConsoleErrorReporter(&buf, logger)
	highErr := appErrors.NewAppError("C1", "critical", nil, appErrors.CategoryUnknown, appErrors.SeverityHigh, appErrors.RecoveryNone)
	rep2 := &appErrors.ErrorReport{Error: highErr, Handled: true, Timestamp: time.Now()}
	if err := crep.Report(context.Background(), rep2); err != nil {
		t.Fatalf("console report failed: %v", err)
	}
	if !bytes.Contains(buf.Bytes(), []byte("ERROR:")) {
		t.Fatalf("expected console output to include ERROR:, got %s", buf.String())
	}

	// Composite reporter with one failing reporter
	bad := &mockReporter{name: "bad", reportErr: errors.New("fail")}
	ok := &mockReporter{name: "ok"}
	composite := appErrors.NewCompositeErrorReporter([]appErrors.ErrorReporter{bad, ok}, logger)
	err = composite.Report(context.Background(), report)
	if err == nil {
		t.Fatal("expected composite report to return an error when one reporter fails")
	}
	// Ensure ok reporter still was invoked
	if !ok.called {
		t.Fatal("expected ok reporter to be called even if one fails")
	}
}

func TestErrorMetricsRecording(t *testing.T) {
	m := appErrors.NewErrorMetrics()
	appErr := appErrors.NewAppError("M1", "metrics test", nil, appErrors.CategoryDatabase, appErrors.SeverityHigh, appErrors.RecoveryRetry)
	m.RecordError(appErr)
	m.RecordRecovery(true)

	got := m.GetMetrics()
	if got.TotalErrors != 1 {
		t.Fatalf("expected total errors 1, got %d", got.TotalErrors)
	}
	if got.ErrorsByCategory[string(appErrors.CategoryDatabase)] != 1 {
		t.Fatalf("expected category count 1, got %d", got.ErrorsByCategory[string(appErrors.CategoryDatabase)])
	}
	if got.RecoverySuccess != 1 {
		t.Fatalf("expected recovery success 1, got %d", got.RecoverySuccess)
	}
}
