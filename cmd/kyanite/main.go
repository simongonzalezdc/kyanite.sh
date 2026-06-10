package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/design"
	"github.com/kyanite/design/icons"
	"github.com/kyanite/config"
)

var version = "dev"

type appEntry struct {
	name       string
	desc       string
	icon       string
	binaryName string
}

var apps = []appEntry{
	{name: "Focus", desc: "Task management & productivity — AI-powered mission control", icon: "focus", binaryName: "focus"},
	{name: "Noise", desc: "Music production notebook — lyrics, ideas & creative workflow", icon: "music", binaryName: "noise"},
	{name: "Syntax", desc: "Distraction-free Markdown editor with live preview", icon: "file", binaryName: "syntax"},
	{name: "Prism", desc: "Color palette designer — harmony generation & WCAG checker", icon: "photo", binaryName: "prism"},
}

// ── Model ───────────────────────────────────────────────────────────

type model struct {
	cursor  int
	width   int
	height  int
	quitting bool
}

func newModel() model {
	return model{width: 70, height: 22}
}

func (m model) Init() tea.Cmd { return nil }

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down", "j":
			if m.cursor < len(apps)-1 {
				m.cursor++
			}
		case "enter":
			return m, launchApp(apps[m.cursor])
		case "q", "esc", "ctrl+c":
			m.quitting = true
			return m, tea.Quit
		}
	}
	return m, nil
}

func (m model) View() string {
	if m.quitting {
		return ""
	}

	t := design.DefaultTheme()

	// ── Frame style ──
	frame := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(t.Border).
		Foreground(t.Text).
		Background(t.Background).
		Padding(1, 3)

	// ── Header ──
	icon := lipgloss.NewStyle().Foreground(t.Primary).Render(icons.GetIcon("tools"))
	bullet := lipgloss.NewStyle().Foreground(t.Muted).Render(icons.GetIcon("bullet"))
	title := lipgloss.NewStyle().Foreground(t.Accent).Bold(true).Render("kyanite.sh")
	ver := lipgloss.NewStyle().Foreground(t.Muted).Render("v" + version)
	header := lipgloss.JoinHorizontal(lipgloss.Center, icon, " ", title, " ", bullet, " ", ver)

	// ── Separator ──
	innerW := m.width - frame.GetHorizontalFrameSize()
	sep := lipgloss.NewStyle().Foreground(t.Border).Render(strings.Repeat("─", max(1, innerW)))

	// ── App items ──
	var items []string
	for i, a := range apps {
		var numStyle, nameStyle, descStyle lipgloss.Style
		if i == m.cursor {
			numStyle = lipgloss.NewStyle().Foreground(t.Accent).Bold(true).Background(t.Panel).Width(4)
			nameStyle = lipgloss.NewStyle().Foreground(t.Accent).Bold(true).Background(t.Panel)
			descStyle = lipgloss.NewStyle().Foreground(t.Text).Background(t.Panel)
		} else {
			numStyle = lipgloss.NewStyle().Foreground(t.Muted).Width(4)
			nameStyle = lipgloss.NewStyle().Foreground(t.Text).Bold(true)
			descStyle = lipgloss.NewStyle().Foreground(t.Muted)
		}

		iconStr := icons.GetIcon(a.icon)
		num := numStyle.Render(fmt.Sprintf(" %d.", i+1))
		name := nameStyle.Render(fmt.Sprintf(" %s  %s", iconStr, a.name))
		desc := descStyle.Render("     " + a.desc)

		items = append(items, num+name, desc)
		if i < len(apps)-1 {
			items = append(items, "")
		}
	}
	itemBlock := strings.Join(items, "\n")

	// ── Footer ──
	help := "↑/k ↓/j  navigate  •  enter  launch  •  q  quit"
	helpRendered := lipgloss.NewStyle().Foreground(t.Muted).Render(help)
	themeRendered := lipgloss.NewStyle().Foreground(t.Border).Render(t.Name)
	padW := max(1, innerW-lipgloss.Width(helpRendered)-lipgloss.Width(themeRendered))
	footer := helpRendered + strings.Repeat(" ", padW) + themeRendered

	// ── Assemble ──
	content := header + "\n" + sep + "\n\n" + itemBlock + "\n" + sep + "\n" + footer
	return frame.Render(content)
}

// ── App launching ───────────────────────────────────────────────────

func launchApp(entry appEntry) tea.Cmd {
	binPath, err := findBinary(entry.binaryName)
	if err != nil {
		return tea.Quit
	}
	return func() tea.Msg {
		cmd := exec.Command(binPath)
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		_ = cmd.Run()
		return nil
	}
}

func findBinary(name string) (string, error) {
	self, _ := os.Executable()
	selfDir := filepath.Dir(self)

	// Same directory
	if isExec(filepath.Join(selfDir, name)) {
		return filepath.Join(selfDir, name), nil
	}

	// dist/ relative to workspace
	for _, rel := range []string{"../../dist", "../dist", "dist"} {
		candidate := filepath.Join(selfDir, rel, name)
		if abs, err := filepath.Abs(candidate); err == nil && isExec(abs) {
			return abs, nil
		}
	}

	// PATH
	if p, err := exec.LookPath(name); err == nil && isExec(p) {
		return p, nil
	}

	return "", fmt.Errorf("binary %q not found — run: make build", name)
}

func isExec(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}
	if runtime.GOOS == "windows" {
		return strings.ToLower(filepath.Ext(path)) == ".exe" && info.Mode().IsRegular()
	}
	return info.Mode().IsRegular() && info.Mode()&0111 != 0
}

// ── Main ────────────────────────────────────────────────────────────

func main() {
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "config":
			handleConfigCommand()
			return
		case "version", "-v", "--version":
			fmt.Printf("kyanite %s\n", version)
			return
		case "help", "-h", "--help":
			fmt.Println("kyanite — unified TUI launcher for the kyanite.sh suite")
			fmt.Println()
			fmt.Println("Usage:")
			fmt.Println("  kyanite              Launch app selector")
			fmt.Println("  kyanite config init  Create default config file")
			fmt.Println("  kyanite config show  Show resolved configuration")
			fmt.Println("  kyanite version      Print version")
			return
		default:
			fmt.Fprintf(os.Stderr, "unknown command: %s\nRun 'kyanite help' for usage.\n", os.Args[1])
			os.Exit(1)
		}
	}

	p := tea.NewProgram(newModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "kyanite: %v\n", err)
		os.Exit(1)
	}
}

func handleConfigCommand() {
	if len(os.Args) < 3 {
		fmt.Fprintln(os.Stderr, "Usage: kyanite config <init|show>")
		os.Exit(1)
	}
	switch os.Args[2] {
	case "init":
		if err := config.Init(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	case "show":
		if err := config.Show(); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "unknown config command: %s\n", os.Args[2])
		os.Exit(1)
	}
}