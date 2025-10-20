package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
	"math/rand"

	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/focus/pkg/styles"
	"github.com/kyanite/focus/internal/tui"
	"github.com/kyanite/focus/internal/cli"
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
		cli.Execute()
		return
	}

	if err := setupOllama(); err != nil {
		fmt.Printf("⚠️  Ollama setup issue: %v\n", err)
	}
	fmt.Println("🌌 Loading focus.sh TUI System...")
	fmt.Println("   ✨ AI-powered task management with professional interface")
	fmt.Println()
	fmt.Println("🚀 Launching TUI dashboard...")
	runTUIDirectly()
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
	if runtime.GOOS == "windows" {
		// Download Windows installer
		return downloadFile("https://ollama.com/download/OllamaSetup.exe", "OllamaSetup.exe", func() error {
			cmd := exec.Command("OllamaSetup.exe", "/S")
			return cmd.Run()
		})
	} else if runtime.GOOS == "darwin" {
		// Use Homebrew or download macOS installer
		return exec.Command("brew", "install", "ollama").Run()
	} else {
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
	resp.Body.Close()
	return resp.StatusCode == 200
}

func startOllama() error {
	cmd := exec.Command("ollama", "serve")
	return cmd.Start()
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

func showEpicIntro() {
	// Clear screen and start the show
	fmt.Print("\033[2J\033[H")
	
	introFrames := []struct {
		art     string
		color   lipgloss.Color
		delay   time.Duration
	}{
		{
			`
╔══════════════════════════════════════════════════════════════╗
║                                                              ║
║    ◢◤◢◤◢◤◢◤◢◤◢◤◢◤◢◤◢◤◢◤◢◤◢◤◢◤◢◤◢◤◢◤◢◤◢◤◢◤◢◤◢◤    ║
║    ◥◣◥◣◥◣◥◣◥◣◥◣◥◣◥◣◥◣◥◣◥◣◥◣◥◣◥◣◥◣◥◣◥◣◥◣◥◣    ║
║       ◢◤◢◤◢◤◢◤◢◤◢◤◢◤◢◤◢◤◢◤◢◤◢◤◢◤◢◤◢◤◢◤       ║
║       ◥◣◥◣◥◣◥◣◥◣◥◣◥◣◥◣◥◣◥◣◥◣◥◣◥◣◥◣◥◣◥◣◥◣◥◣       ║
║       ◢◤◢◤◢◤◢◤◢◤◢◤◢◤◢◤◢◤◢◤◢◤◢◤◢◤◢◤◢◤◢◤       ║
║       ◥◣◥◣◥◣◥◣◥◣◥◣◥◣◥◣◥◣◥◣◥◣◥◣◥◣◥◣◥◣◥◣◥◣◥◣       ║
║                                                              ║
║            ⚡⚡⚡⚡⚡⚡⚡⚡⚡⚡⚡⚡⚡⚡⚡⚡⚡⚡⚡⚡⚡            ║
║                 ◈◈◈◈◈◈◈◈◈◈◈◈◈◈◈◈◈◈◈◈                 ║
╚══════════════════════════════════════════════════════════════╝`,
			styles.SynthwavePink,
			time.Millisecond * 800,
		},
		{
			`
▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰
▰                                                              ▰
▰    ⚡ INITIALIZING NEURAL INTERFACE...                        ▰
▰    🤖 LOADING AI ASSISTANT...                               ▰
▰    💫 ACTIVATING VISUAL ENHANCEMENTS...                       ▰
▰    🌐 CONNECTING TO CHARM LIBRARIES...                       ▰
▰                                                              ▰
▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰`,
			styles.SynthwaveCyan,
			time.Millisecond * 600,
		},
		{
			`
██████████████████████████████████████████████████████████████
██████████████████████████████████████████████████████████████
██████████████████████████████████████████████████████████████
████ ▀▄▀▄▀▄ AI ENHANCED PRODUCTIVITY ▄▀▄▀▄▀ ███████████████████
██████████████████████████████████████████████████████████████
██████████████████████████████████████████████████████████████`,
			styles.SynthwaveGreen,
			time.Millisecond * 400,
		},
	}

	for _, frame := range introFrames {
		// Clear screen
		fmt.Print("\033[2J\033[H")
		
		// Render frame with epic styling
		styled := lipgloss.NewStyle().
			Foreground(frame.color).
			Background(styles.DeepSpace).
			Bold(true).
			AlignHorizontal(lipgloss.Center).
			Render(frame.art)
		
		fmt.Println(styled)
		time.Sleep(frame.delay)
	}

	// Final glitch effect
	glitchSymbols := []string{
		"⚡⚡⚡ ◈◈◈ ◆◆◆ ◊◊◊ ⚡⚡⚡",
		"◆◆◆ ◊◊◊ ⚡⚡⚡ ◈◈◈ ◆◆◆",
		"◈◈◈ ◆◆◆ ◊◊◊ ⚡⚡⚡ ◈◈◈",
		"◊◊◊ ⚡⚡⚡ ◈◈◈ ◆◆◆ ◊◊◊",
		"⚡◆◈◊ ◈⚡◆◈ ◊◈◆⚡ ◊⚡◆◈",
		"◆◈⚡◊ ◆◈◊⚡ ◈⚡◆◊ ◈◆⚡◊",
		"◈◊⚡◆ ◊◈⚡◆ ◊⚡◆◈ ⚡◈◊◆",
	}
	
	for range 10 {
		fmt.Print("\033[2J\033[H")
		
		glitchText := glitchSymbols[rand.Intn(len(glitchSymbols))]
		glitchStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(fmt.Sprintf("#%02X%02X%02X", 
				rand.Intn(255), rand.Intn(255), rand.Intn(255)))).
			Background(styles.DeepSpace).
			Bold(true).
			AlignHorizontal(lipgloss.Center).
			Render(glitchText)
		
		fmt.Println(glitchStyle)
		time.Sleep(time.Millisecond * 100)
	}

	// Clear and show ready state
	fmt.Print("\033[2J\033[H")
	
	readyText := styles.SynthwaveTitle("🚀 FOCUS.SH SYSTEMS READY")
	fmt.Println(readyText)
	fmt.Println()
	
	readyMsg := styles.HolographicText("AI-powered task management with maximum visual impact achieved.")
	fmt.Println(readyMsg)
	fmt.Println()
	
	controlHint := lipgloss.NewStyle().
		Foreground(styles.SynthwaveCyan).
		Background(styles.DarkVoid).
		Italic(true).
		Render("💫 Type 'focus --help' to begin your productivity mission")
	fmt.Println(controlHint)
	
	time.Sleep(time.Second * 2)
	
	// Clear for main interface
	fmt.Print("\033[2J\033[H")
}

