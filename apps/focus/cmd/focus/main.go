package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/kyanite/focus/internal/cli"
	"github.com/kyanite/focus/internal/tui"
)

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