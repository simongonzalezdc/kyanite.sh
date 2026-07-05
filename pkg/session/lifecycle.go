// Package session provides shared session lifecycle management for kyanite.sh apps.
// It handles session ID generation, welcome messages, session persistence, and
// cross-app context saving — eliminating duplicate lifecycle code in each app.
package session

import (
	"context"
	"fmt"
	"time"

	"github.com/kyanite/ai"
	"github.com/kyanite/appnames"
)

// Session holds the current app session state.
type Session struct {
	App       string
	SessionID string
	brain     *ai.Brain
}

// OnStartup initializes a new session for the given app.
// It generates a session ID and prints a welcome message if a recent session exists.
// Returns nil session if brain is nil (offline mode — no session tracking).
func OnStartup(brain *ai.Brain, app string) *Session {
	if brain == nil {
		return &Session{App: app, SessionID: fmt.Sprintf("%s-%d", app, time.Now().Unix())}
	}

	s := &Session{
		App:       app,
		SessionID: fmt.Sprintf("%s-%d", app, time.Now().Unix()),
		brain:     brain,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	sessions, err := brain.GetRecentSessions(ctx, 1)
	if err == nil && len(sessions) > 0 {
		fmt.Printf("Welcome back — last session: %s\n", sessions[0].Title)
	}

	return s
}

// OnShutdown persists the session and cross-app context.
// Safe to call even if brain is nil (no-op).
func OnShutdown(s *Session, title string, state any) {
	if s == nil || s.brain == nil {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_ = s.brain.SaveSession(ctx, s.SessionID, title, state)

	// Save cross-app context for other apps
	otherApps := appnames.Others(s.App)
	for _, target := range otherApps {
		_ = s.brain.SaveCrossAppContext(ctx, target, "activity", title, 0.5)
	}
}

// ID returns the session ID.
func (s *Session) ID() string {
	if s == nil {
		return ""
	}
	return s.SessionID
}

// Brain returns the underlying Brain, or nil.
func (s *Session) Brain() *ai.Brain {
	if s == nil {
		return nil
	}
	return s.brain
}

