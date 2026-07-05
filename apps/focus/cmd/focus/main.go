package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/kyanite/ai"
	"github.com/kyanite/appnames"
	"github.com/kyanite/config"
	"github.com/kyanite/focus/internal/cli"
	"github.com/kyanite/focus/internal/tui"
)

// sessionState holds serializable focus session state for Brain persistence.
type sessionState struct {
	TaskCount    int        `json:"task_count"`
	CurrentView  string     `json:"current_view"`
	StartedAt    time.Time  `json:"started_at"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
}

func main() {
	if len(os.Args) > 1 {
		if os.Args[1] == "mcp-server" {
			if err := runMCPServer(); err != nil {
				fmt.Fprintf(os.Stderr, "mcp-server error: %v\n", err)
			}
			return
		}
		// Any other subcommand should be handled by the CLI layer
		if err := cli.Execute(); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		return
	}

	// AI model is now hosted on NUCBox via pkg/ai Brain — no local Ollama setup needed.
	fmt.Println("🌌 Loading focus.sh TUI System...")
	fmt.Println("   ✨ AI-powered task management with professional interface")
	fmt.Println()
	fmt.Println("🚀 Launching TUI dashboard...")
	if err := runTUIDirectly(); err != nil {
		fmt.Fprintf(os.Stderr, "TUI error: %v\n", err)
		os.Exit(1)
	}
}

func runMCPServer() error {
	root := findRepoRoot()
	if root == "" {
		root, _ = os.Getwd()
	}
	exePath := filepath.Join(root, "mcp-server.exe")
	if _, err := os.Stat(exePath); err == nil {
		cmd := exec.Command(exePath)
		cmd.Dir = root
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		return cmd.Run()
	}
	cmd := exec.Command("go", "run", "./cmd/focus-mcp")
	cmd.Dir = root
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func findRepoRoot() string {
	d, _ := os.Getwd()
	for range 10 {
		gm := filepath.Join(d, "go.mod")
		mcpmain := filepath.Join(d, "cmd", "focus-mcp", "main.go")
		if _, err := os.Stat(gm); err == nil {
			if _, err2 := os.Stat(mcpmain); err2 == nil {
				return d
			}
		}
		n := filepath.Dir(d)
		if n == d {
			break
		}
		d = n
	}
	return ""
}

// Run TUI directly without CLI interference
func runTUIDirectly() error {
	// Create Brain directly for session lifecycle (pkg/ai).
	root, _ := config.Load()
	cfg := ai.ConfigFromRoot(root, "focus")
	brain, _ := ai.New(cfg)

	ctx := context.Background()
	sessionID := fmt.Sprintf("focus-%d", time.Now().Unix())
	startedAt := time.Now()

	// On startup: attempt to resume the most recent session
	if brain != nil && brain.IsMemoryAvailable(ctx) {
		sessions, err := brain.GetRecentSessions(ctx, 1)
		if err == nil && len(sessions) > 0 {
			s := sessions[0]
			ago := formatTimeAgo(time.Since(s.UpdatedAt))
			title := s.Title
			if title == "" {
				title = s.SessionID
			}
			fmt.Printf("Welcome back — last session: %s (%s)\n", title, ago)
		}
	}

	// Create sample TUI tasks in the correct DashboardTask format
	tasks := []tui.DashboardTask{
		{
			ID:          "1",
			Description: "Complete synthwave dashboard project",
			Priority:    "high",
			Status:      "pending",
			CreatedAt:   time.Now(),
			Deadline:    nil,
			Categories:  []string{"coding", "focus"},
			Notes:       "Focus on focus.sh UI components and retro-futuristic design",
		},
		{
			ID:          "2",
			Description: "Design retro-futuristic UI components",
			Priority:    "medium",
			Status:      "pending",
			CreatedAt:   time.Now(),
			Deadline:    nil,
			Categories:  []string{"design", "synthwave"},
			Notes:       "Professional AI assistant with helpful personality",
		},
		{
			ID:          "3",
			Description: "Integrate AI chat with synthwave personality",
			Priority:    "medium",
			Status:      "pending",
			CreatedAt:   time.Now(),
			Deadline:    nil,
			Categories:  []string{"ai", "productivity"},
			Notes:       "Professional and helpful AI assistant",
		},
	}

	// Deferred shutdown: save session and cross-app context
	defer func() {
		if brain == nil {
			return
		}

		completed := 0
		for _, t := range tasks {
			if t.Status == "completed" {
				completed++
			}
		}

		// Save session
		title := fmt.Sprintf("%d tasks, %d completed", len(tasks), completed)
		state := sessionState{
			TaskCount:   len(tasks),
			CurrentView: "dashboard",
			StartedAt:   startedAt,
		}
		if completed > 0 {
			now := time.Now()
			state.CompletedAt = &now
		}
		if err := brain.SaveSession(ctx, sessionID, title, state); err != nil {
			log.Printf("focus: save session: %v", err)
		}

		// Save cross-app context for other kyanite apps
		if completed > 0 {
			summary := fmt.Sprintf("Completed %d tasks in focus", completed)
			for _, targetApp := range appnames.Others(appnames.Focus) {
				if err := brain.SaveCrossAppContext(ctx, targetApp, "task_completion", summary, 0.6); err != nil {
					log.Printf("focus: save cross-app context for %s: %v", targetApp, err)
				}
			}
		}

		brain.Close()
	}()

	fmt.Printf("📋 Loaded %d tasks into focus.sh system...\n", len(tasks))

	// Launch the TUI
	fmt.Println("🎮 focus.sh Task Management System Starting...")
	fmt.Println("   🌨 Professional interface loading...")
	fmt.Println("   ✨ Kyanite theme system activating...")
	fmt.Println("   ⚡ System ready...")
	fmt.Println()

	// Launch actual TUI dashboard
	return tui.StartMainDashboard(tasks)
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