package errors

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/Kyanite/noise/internal/domain"
	logging "github.com/Kyanite/noise/internal/logging"
)

// TestLogger provides a logger for testing that captures log messages
type TestLogger struct {
	writer *testLogWriter
	mu     sync.RWMutex
	t      *testing.T
	*logging.Logger
}

// NewTestLogger creates a new test logger
func NewTestLogger(t *testing.T) *TestLogger {
	// Create a writer that will capture all log messages
	writer := &testLogWriter{messages: make([]string, 0)}

	// Create a real logger that writes to the buffer
	config := &logging.Config{
		Level:      logging.DEBUG,
		Output:     writer,
		ShowCaller: false,
	}

	logger, err := logging.New(config)
	if err != nil {
		t.Fatalf("Failed to create test logger: %v", err)
	}

	return &TestLogger{
		writer: writer,
		t:      t,
		Logger: logger,
	}
}

// messages provides access to the captured log messages
func (tl *TestLogger) messages() []string {
	tl.writer.mu.RLock()
	defer tl.writer.mu.RUnlock()
	return tl.writer.messages
}

// testLogWriter is a writer that captures log messages
type testLogWriter struct {
	messages []string
	mu       sync.RWMutex
}

// Write implements io.Writer
func (tlw *testLogWriter) Write(p []byte) (n int, err error) {
	tlw.mu.Lock()
	defer tlw.mu.Unlock()

	message := string(p)
	tlw.messages = append(tlw.messages, message)
	return len(p), nil
}

// GetMessages returns all logged messages
func (tl *TestLogger) GetMessages() []string {
	tl.writer.mu.RLock()
	defer tl.writer.mu.RUnlock()

	messages := make([]string, len(tl.writer.messages))
	copy(messages, tl.writer.messages)
	return messages
}

// Clear clears all logged messages
func (tl *TestLogger) Clear() {
	tl.writer.mu.Lock()
	defer tl.writer.mu.Unlock()
	tl.writer.messages = make([]string, 0)
}

// ContainsMessage checks if a message containing the given text was logged
func (tl *TestLogger) ContainsMessage(text string) bool {
	for _, msg := range tl.GetMessages() {
		if contains(msg, text) {
			return true
		}
	}
	return false
}

// TestErrorReporter is a mock error reporter for testing
type TestErrorReporter struct {
	reports []*ErrorReport
	mu      sync.RWMutex
}

// NewTestErrorReporter creates a new test error reporter
func NewTestErrorReporter() *TestErrorReporter {
	return &TestErrorReporter{
		reports: make([]*ErrorReport, 0),
	}
}

// Report handles an error report
func (ter *TestErrorReporter) Report(ctx context.Context, report *ErrorReport) error {
	ter.mu.Lock()
	defer ter.mu.Unlock()
	ter.reports = append(ter.reports, report)
	return nil
}

// Name returns the reporter name
func (ter *TestErrorReporter) Name() string {
	return "test_reporter"
}

// GetReports returns all received reports
func (ter *TestErrorReporter) GetReports() []*ErrorReport {
	ter.mu.RLock()
	defer ter.mu.RUnlock()

	reports := make([]*ErrorReport, len(ter.reports))
	copy(reports, ter.reports)
	return reports
}

// Clear clears all reports
func (ter *TestErrorReporter) Clear() {
	ter.mu.Lock()
	defer ter.mu.Unlock()
	ter.reports = make([]*ErrorReport, 0)
}

// TestRecoveryFunc is a mock recovery function for testing
type TestRecoveryFunc struct {
	callCount  int
	lastError  error
	shouldFail bool
	failAfter  int
	mu         sync.RWMutex
}

// NewTestRecoveryFunc creates a new test recovery function
func NewTestRecoveryFunc() *TestRecoveryFunc {
	return &TestRecoveryFunc{
		callCount:  0,
		shouldFail: false,
		failAfter:  -1,
	}
}

// Execute executes the recovery function
func (trf *TestRecoveryFunc) Execute(ctx context.Context, err error) error {
	trf.mu.Lock()
	defer trf.mu.Unlock()

	trf.callCount++
	trf.lastError = err

	if trf.shouldFail && (trf.failAfter == -1 || trf.callCount > trf.failAfter) {
		return fmt.Errorf("recovery failed on attempt %d", trf.callCount)
	}

	return nil
}

// GetCallCount returns the number of times the function was called
func (trf *TestRecoveryFunc) GetCallCount() int {
	trf.mu.RLock()
	defer trf.mu.RUnlock()
	return trf.callCount
}

// GetLastError returns the last error passed to the function
func (trf *TestRecoveryFunc) GetLastError() error {
	trf.mu.RLock()
	defer trf.mu.RUnlock()
	return trf.lastError
}

// SetShouldFail sets whether the function should fail
func (trf *TestRecoveryFunc) SetShouldFail(shouldFail bool, failAfter int) {
	trf.mu.Lock()
	defer trf.mu.Unlock()
	trf.shouldFail = shouldFail
	trf.failAfter = failAfter
}

// Reset resets the function state
func (trf *TestRecoveryFunc) Reset() {
	trf.mu.Lock()
	defer trf.mu.Unlock()
	trf.callCount = 0
	trf.lastError = nil
	trf.shouldFail = false
	trf.failAfter = -1
}

// TestSetup provides common setup for error handling tests
type TestSetup struct {
	TempDir  string
	Logger   *TestLogger
	Manager  *ErrorManager
	Reporter *TestErrorReporter
	Recovery *TestRecoveryFunc
	Cleanup  func()
}

// NewTestSetup creates a new test setup
func NewTestSetup(t *testing.T) *TestSetup {
	// Create temporary directory
	tempDir, err := os.MkdirTemp("", "error_test_*")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	// Create test logger
	logger := NewTestLogger(t)

	// Create error manager
	config := DefaultErrorConfig()
	config.LogDirectory = tempDir
	config.EnableReporting = true // Enable reporting for tests
	manager := NewErrorManager(logger.Logger, config)

	// Create test reporter
	reporter := NewTestErrorReporter()
	manager.AddReporter(reporter)

	// Create test recovery function
	recovery := NewTestRecoveryFunc()

	// Setup cleanup function
	cleanup := func() {
		os.RemoveAll(tempDir)
		manager.Close()
	}

	return &TestSetup{
		TempDir:  tempDir,
		Logger:   logger,
		Manager:  manager,
		Reporter: reporter,
		Recovery: recovery,
		Cleanup:  cleanup,
	}
}

// CreateTestSong creates a test song for testing
func CreateTestSong(id int, title string) *domain.Song {
	return &domain.Song{
		ID:       id,
		Filepath: fmt.Sprintf("/test/song_%d.md", id),
		Metadata: domain.SongMetadata{
			Title: title,
		},
		Sections: []domain.Section{
			{
				Type:   domain.SectionVerse,
				Number: 1,
				Lines: []domain.Line{
					{Text: "Test line 1", Syllables: 3},
					{Text: "Test line 2", Syllables: 3},
				},
			},
		},
	}
}

// CreateTestFile creates a test file with the given content
func CreateTestFile(t *testing.T, dir, filename, content string) string {
	filePath := filepath.Join(dir, filename)
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file %s: %v", filePath, err)
	}
	return filePath
}

// CreateCorruptedTestFile creates a corrupted test file
func CreateCorruptedTestFile(t *testing.T, dir, filename string) string {
	filePath := filepath.Join(dir, filename)
	ext := filepath.Ext(filename)

	var corruptedContent []byte

	switch ext {
	case ".json":
		// Create invalid JSON for JSON files
		corruptedContent = []byte(`{"invalid": json content}`)
	case ".txt", ".md":
		// Create text with null bytes to trigger binary corruption detection
		corruptedContent = []byte("text with \x00 null byte corruption")
	default:
		// Default to invalid JSON
		corruptedContent = []byte(`{"invalid": json content}`)
	}

	err := os.WriteFile(filePath, corruptedContent, 0644)
	if err != nil {
		t.Fatalf("Failed to create corrupted test file %s: %v", filePath, err)
	}
	return filePath
}

// CreateBackupFile creates a backup file for testing
func CreateBackupFile(t *testing.T, backupDir, backupID string, song *domain.Song) string {
	backupPath := filepath.Join(backupDir, fmt.Sprintf("backup_%s.json", backupID))

	// Simple JSON serialization for testing
	content := fmt.Sprintf(`{
		"id": %d,
		"filepath": "%s",
		"metadata": {
			"title": "%s"
		},
		"sections": [
			{
				"type": "verse",
				"number": 1,
				"lines": [
					{"text": "Test line 1", "syllables": 3},
					{"text": "Test line 2", "syllables": 3}
				]
			}
		]
	}`, song.ID, song.Filepath, song.Metadata.Title)

	err := os.WriteFile(backupPath, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create backup file %s: %v", backupPath, err)
	}
	return backupPath
}

// AssertError asserts that an error has the expected properties
func AssertError(t *testing.T, err error, expectedCode string, expectedCategory ErrorCategory) {
	t.Helper()

	if err == nil {
		t.Fatalf("Expected error but got nil")
	}

	appErr, ok := err.(*AppError)
	if !ok {
		t.Fatalf("Expected *AppError but got %T", err)
	}

	if appErr.Code != expectedCode {
		t.Errorf("Expected error code %s but got %s", expectedCode, appErr.Code)
	}

	if appErr.Category != expectedCategory {
		t.Errorf("Expected error category %s but got %s", expectedCategory, appErr.Category)
	}
}

// AssertNoError asserts that there is no error
func AssertNoError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

// AssertContains asserts that a slice contains an element
func AssertContains[T comparable](t *testing.T, slice []T, element T) {
	t.Helper()

	for _, item := range slice {
		if item == element {
			return
		}
	}

	t.Errorf("Expected slice to contain %v but it didn't", element)
}

// AssertNotContains asserts that a slice does not contain an element
func AssertNotContains[T comparable](t *testing.T, slice []T, element T) {
	t.Helper()

	for _, item := range slice {
		if item == element {
			t.Errorf("Expected slice to not contain %v but it did", element)
			return
		}
	}
}

// WaitForCondition waits for a condition to be true or times out
func WaitForCondition(t *testing.T, condition func() bool, timeout time.Duration, message string) {
	t.Helper()

	start := time.Now()
	for time.Since(start) < timeout {
		if condition() {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}

	t.Fatalf("Condition not met within timeout: %s", message)
}

// MockNotificationHandler is a mock notification handler for testing
type MockNotificationHandler struct {
	notifications []*Notification
	mu            sync.RWMutex
}

// NewMockNotificationHandler creates a new mock notification handler
func NewMockNotificationHandler() *MockNotificationHandler {
	return &MockNotificationHandler{
		notifications: make([]*Notification, 0),
	}
}

// HandleNotification handles a notification
func (mnh *MockNotificationHandler) HandleNotification(notification *Notification) {
	mnh.mu.Lock()
	defer mnh.mu.Unlock()
	mnh.notifications = append(mnh.notifications, notification)
}

// GetNotifications returns all received notifications
func (mnh *MockNotificationHandler) GetNotifications() []*Notification {
	mnh.mu.RLock()
	defer mnh.mu.RUnlock()

	notifications := make([]*Notification, len(mnh.notifications))
	copy(notifications, mnh.notifications)
	return notifications
}

// Clear clears all notifications
func (mnh *MockNotificationHandler) Clear() {
	mnh.mu.Lock()
	defer mnh.mu.Unlock()
	mnh.notifications = make([]*Notification, 0)
}

// CountNotificationsByType returns the count of notifications by type
func (mnh *MockNotificationHandler) CountNotificationsByType(notificationType NotificationType) int {
	mnh.mu.RLock()
	defer mnh.mu.RUnlock()

	count := 0
	for _, notification := range mnh.notifications {
		if notification.Type == notificationType {
			count++
		}
	}
	return count
}

// MockUIComponent is a mock UI component for testing
type MockUIComponent struct {
	visible bool
	updates []interface{}
	mu      sync.RWMutex
}

// NewMockUIComponent creates a new mock UI component
func NewMockUIComponent() *MockUIComponent {
	return &MockUIComponent{
		visible: false,
		updates: make([]interface{}, 0),
	}
}

// Show shows the UI component
func (muc *MockUIComponent) Show() {
	muc.mu.Lock()
	defer muc.mu.Unlock()
	muc.visible = true
}

// Hide hides the UI component
func (muc *MockUIComponent) Hide() {
	muc.mu.Lock()
	defer muc.mu.Unlock()
	muc.visible = false
}

// IsVisible returns whether the component is visible
func (muc *MockUIComponent) IsVisible() bool {
	muc.mu.RLock()
	defer muc.mu.RUnlock()
	return muc.visible
}

// Update updates the UI component
func (muc *MockUIComponent) Update(data interface{}) {
	muc.mu.Lock()
	defer muc.mu.Unlock()
	muc.updates = append(muc.updates, data)
}

// GetUpdates returns all updates
func (muc *MockUIComponent) GetUpdates() []interface{} {
	muc.mu.RLock()
	defer muc.mu.RUnlock()

	updates := make([]interface{}, len(muc.updates))
	copy(updates, muc.updates)
	return updates
}

// Clear clears all updates
func (muc *MockUIComponent) Clear() {
	muc.mu.Lock()
	defer muc.mu.Unlock()
	muc.updates = make([]interface{}, 0)
}
