package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"github.com/kyanite/focus/internal/cli"
	"github.com/kyanite/focus/internal/tui"
)

// ollamaCmd holds a reference to the ollama serve process for cleanup.
var ollamaCmd *exec.Cmd

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

	if err := setupOllama(); err != nil {
		fmt.Printf("⚠️  Ollama setup issue: %v\n", err)
	}
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

func setupOllama() error {
	// Check if ollama is available
	if _, err := exec.LookPath("ollama"); err != nil {
		fmt.Println("🤖 Ollama not found. Installing...")
		if err := installOllama(); err != nil {
			return fmt.Errorf("could not install ollama: %w", err)
		}
	}

	// Check if ollama is running
	if !isOllamaRunning() {
		fmt.Println("🚀 Starting Ollama service...")
		if err := startOllama(); err != nil {
			fmt.Printf("⚠️  Could not start Ollama: %v\n", err)
			fmt.Println("💡 You can start it manually with: ollama serve")
		} else {
			fmt.Println("✅ Ollama started successfully!")
		}
	}

	// Ensure required model is available
	if !isModelAvailable("qwen2.5:1.5b") {
		fmt.Println("📥 Downloading qwen2.5:1.5b model (this may take a moment)...")
		if err := pullModel("qwen2.5:1.5b"); err != nil {
			fmt.Printf("⚠️  Model download issue: %v\n", err)
		} else {
			fmt.Println("✅ AI model ready!")
		}
	} else {
		fmt.Println("✅ AI model already available")
	}

	return nil
}

func installOllama() error {
	fmt.Println("💡 Installing Ollama automatically...")

	// Detect OS and install accordingly
	switch runtime.GOOS {
	case "windows":
		// Download Windows installer
		return downloadFile("https://ollama.com/download/OllamaSetup.exe", "OllamaSetup.exe", func() error {
			cmd := exec.Command("OllamaSetup.exe", "/S")
			return cmd.Run()
		})
	case "darwin":
		// Use Homebrew or download macOS installer
		return exec.Command("brew", "install", "ollama").Run()
	default:
		// Linux download
		return downloadFile("https://ollama.com/download/ollama-linux-amd64.tgz", "ollama.tgz", func() error {
			return exec.Command("tar", "xzf", "ollama.tgz").Run()
		})
	}
}

func isOllamaRunning() bool {
	resp, err := http.Get("http://localhost:11434/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()
	return resp.StatusCode == 200
}

func startOllama() error {
	cmd := exec.Command("ollama", "serve")
	if err := cmd.Start(); err != nil {
		return err
	}
	// ollama serve is designed to run as a persistent daemon.
	// Track the process so we can clean up on exit.
	ollamaCmd = cmd
	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		if ollamaCmd != nil && ollamaCmd.Process != nil {
			_ = ollamaCmd.Process.Kill()
		}
		os.Exit(0)
	}()
	return nil
}

func isModelAvailable(model string) bool {
	resp, err := http.Get("http://localhost:11434/api/tags")
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	var tags struct {
		Models []struct {
			Name string `json:"name"`
		} `json:"models"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&tags); err != nil {
		return false
	}

	for _, m := range tags.Models {
		if strings.Contains(m.Name, model) {
			return true
		}
	}
	return false
}

func pullModel(model string) error {
	cmd := exec.Command("ollama", "pull", model)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func downloadFile(url, filename string, postInstall func() error) error {
	fmt.Printf("📥 Downloading from %s...\n", url)

	// Use curl or http client for download
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	out, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}

	// Run post-install step
	if postInstall != nil {
		return postInstall()
	}

	return nil
}

