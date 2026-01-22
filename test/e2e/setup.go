package e2e

import (
	"testing"

	"github.com/Kyanite/noise/internal/app"
	"github.com/Kyanite/noise/internal/app/ai"
	"github.com/Kyanite/noise/internal/collaboration"
	"github.com/Kyanite/noise/internal/config"
	"github.com/Kyanite/noise/internal/errors"
	"github.com/Kyanite/noise/internal/infra/db"
	"github.com/Kyanite/noise/internal/logging"
	"github.com/Kyanite/noise/internal/theme"
	"github.com/Kyanite/noise/internal/ui/dashboard"
	"github.com/Kyanite/noise/internal/ui/editor"
	tea "github.com/charmbracelet/bubbletea"
)

// TestLogger wraps a testing.T to implement the logging.Logger interface
type TestLogger struct {
	T      *testing.T
	Logger *logging.Logger
}

// NewTestLogger creates a test logger that writes to testing.T
func NewTestLogger(t *testing.T) *TestLogger {
	logger, _ := logging.New(logging.DefaultConfig())
	return &TestLogger{
		T:      t,
		Logger: logger,
	}
}

// E2ETestSetup provides a comprehensive test environment for end-to-end testing
type E2ETestSetup struct {
	T              *testing.T
	TempDir        string
	Database       *db.DB
	ErrorManager   *errors.ErrorManager
	ThemeManager   *theme.Manager
	AIAgent        *ai.QuickIdeaAgent
	SessionManager *collaboration.SessionManager
	EditorModel    *editor.SplitPaneModel
	DashboardModel *dashboard.DashboardModel
	Cleanup        func()
}

// NewE2ETestSetup creates a comprehensive test environment
func NewE2ETestSetup(t *testing.T) *E2ETestSetup {
	tempDir := t.TempDir()

	// Initialize database
	database, err := db.New(db.Config{DataDir: tempDir})
	if err != nil {
		t.Fatalf("Failed to create database: %v", err)
	}

	// Initialize error manager
	logger := NewTestLogger(t)
	errorConfig := errors.DefaultErrorConfig()
	errorManager := errors.NewErrorManager(logger.Logger, errorConfig)

	// Initialize theme manager
	themeManager := theme.GetManager()

	// Initialize AI agent with mock provider
	mockProvider := ai.NewMockEnhancementProvider()
	aiAgent := ai.NewQuickIdeaAgent()
	aiAgent = aiAgent.WithKnowledgeBase(mockProvider)

	// Initialize session manager
	mockSessionManager := collaboration.NewMockSessionManager()

	// Initialize AI service
	aiService := app.NewAIService(config.DefaultConfig())

	// Initialize editor model
	editorModel := editor.NewSplitPaneModel(database, aiService)
	editorModel.Update(tea.WindowSizeMsg{Width: 120, Height: 40})

	// Initialize dashboard model
	dashboardModel := dashboard.NewDashboardModel()

	cleanup := func() {
		if err := database.Close(); err != nil {
			t.Logf("Warning: Failed to close database: %v", err)
		}
		if err := errorManager.Close(); err != nil {
			t.Logf("Warning: Failed to close error manager: %v", err)
		}
		editorModel.Cleanup()
	}

	return &E2ETestSetup{
		T:              t,
		TempDir:        tempDir,
		Database:       database,
		ErrorManager:   errorManager,
		ThemeManager:   themeManager,
		AIAgent:        aiAgent,
		SessionManager: mockSessionManager.SessionManager,
		EditorModel:    editorModel,
		DashboardModel: dashboardModel,
		Cleanup:        cleanup,
	}
}
