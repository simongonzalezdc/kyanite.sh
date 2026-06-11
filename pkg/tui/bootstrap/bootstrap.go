// Package bootstrap provides shared initialization for kyanite.sh TUI apps.
//
// Every app follows the same pattern:
//
//	config.Load() → ai.New() → session.OnStartup() → run TUI → session.OnShutdown()
//
// This package extracts the first three steps into a single Init call,
// standardizing brain creation, warning output, and session lifecycle setup.
package bootstrap

import (
	"context"
	"fmt"
	"time"

	"github.com/kyanite/ai"
	"github.com/kyanite/appnames"
	"github.com/kyanite/config"
	"github.com/kyanite/session"
)

// App holds the initialized state returned by Init.
type App struct {
	Brain     *ai.Brain
	Session   *session.Session
	AppName   string
	SessionID string
}

// Init loads config, creates brain, and starts a session.
//
// If the config file is missing or brain creation fails, it prints a warning
// and returns with Brain == nil (offline mode). The caller should proceed —
// all kyanite apps are designed to work without AI.
//
// appName must be one of appnames.Focus, appnames.Noise, etc.
func Init(appName string) *App {
	root, _ := config.Load()
	return initWithConfig(root, appName)
}

// InitWithConfig is like Init but accepts a pre-loaded config.
func InitWithConfig(root *config.Root, appName string) *App {
	return initWithConfig(root, appName)
}

func initWithConfig(root *config.Root, appName string) *App {
	cfg := ai.ConfigFromRoot(root, appName)
	brain, err := ai.New(cfg)
	if err != nil {
		fmt.Fprintf(fmtWrapper{}, "warning: brain init failed (AI features offline): %v\n", err)
	}

	sess := session.OnStartup(brain, appName)

	return &App{
		Brain:     brain,
		Session:   sess,
		AppName:   appName,
		SessionID: sess.ID(),
	}
}

// Shutdown saves session state and cross-app context, then closes the brain.
// Safe to call even if app or brain is nil.
func (a *App) Shutdown(title string, state any) {
	if a == nil {
		return
	}
	session.OnShutdown(a.Session, title, state)
}

// Others returns all app names except the current one.
// Convenience wrapper around appnames.Others.
func (a *App) Others() []string {
	if a == nil {
		return appnames.All
	}
	return appnames.Others(a.AppName)
}

// SaveCrossAppContext saves context for each other app.
// Convenience wrapper that iterates over app.Others().
// Errors are logged but not propagated — cross-app context is best-effort.
func (a *App) SaveCrossAppContext(ctx context.Context, contextType, summary string, score float32) {
	if a == nil || a.Brain == nil || a.Session == nil {
		return
	}
	for _, target := range a.Others() {
		if err := a.Brain.SaveCrossAppContext(ctx, target, contextType, summary, score); err != nil {
			fmt.Fprintf(fmtWrapper{}, "warning: cross-app context save failed for %s: %v\n", target, err)
		}
	}
}

// fmtWrapper wraps fmt.Fprintf to satisfy the io.Writer interface via a String() method.
type fmtWrapper struct{}

func (fmtWrapper) Write(p []byte) (int, error) {
	fmt.Print(string(p))
	return len(p), nil
}

// WelcomeBack prints a welcome message if a recent session exists.
// Returns the session title if found, empty string otherwise.
func (a *App) WelcomeBack() string {
	if a == nil || a.Brain == nil {
		return ""
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if !a.Brain.IsMemoryAvailable(ctx) {
		return ""
	}

	sessions, err := a.Brain.GetRecentSessions(ctx, 1)
	if err != nil || len(sessions) == 0 {
		return ""
	}

	s := sessions[0]
	title := s.Title
	if title == "" {
		title = s.SessionID
	}
	ago := formatTimeAgo(time.Since(s.UpdatedAt))
	fmt.Printf("Welcome back — last session: %s (%s)\n", title, ago)
	return title
}

// formatTimeAgo returns a human-readable duration string.
func formatTimeAgo(d time.Duration) string {
	switch {
	case d < time.Minute:
		return "just now"
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd ago", int(d.Hours()/24))
	}
}
