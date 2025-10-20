package styles

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/pterm/pterm"
	"math/rand"
	"time"
)

// Synthwave Color Palette - Maximum Visual Impact
var (
	// Core Synthwave Colors
	SynthwavePink    = lipgloss.Color("#FF10F0")
	SynthwaveCyan    = lipgloss.Color("#00FFF0")
	SynthwavePurple  = lipgloss.Color("#BD10E0")
	SynthwaveYellow  = lipgloss.Color("#FFF01F")
	SynthwaveGreen   = lipgloss.Color("#39FF14")
	SynthwaveOrange  = lipgloss.Color("#FF6B35")
	SynthwaveRed     = lipgloss.Color("#FF0040")
	
	// Dark Gradient Colors
	DeepSpace    = lipgloss.Color("#0A0014")
	DarkVoid     = lipgloss.Color("#1A0033")
	CyberGrid    = lipgloss.Color("#2D1B69")
	FocusGrid    = lipgloss.Color("#3D2B8F")
	
	// Glitch Colors
	GlitchRed    = lipgloss.Color("#FF0000")
	GlitchGreen  = lipgloss.Color("#00FF00")
	GlitchBlue   = lipgloss.Color("#0000FF")
	
	// Metallic Finishes
	ChromeSilver = lipgloss.Color("#C0C0C0")
	Platinum     = lipgloss.Color("#E5E4E2")
	GoldAccent   = lipgloss.Color("#FFD700")
)

// Glitch Effect Generator
var glitchRand = rand.New(rand.NewSource(time.Now().UnixNano()))

func GlitchText(text string) string {
	glitchChars := []string{"@", "#", "$", "%", "&", "*", "■", "▪", "▫", "◆", "◇", "○", "●", "□", "■", "△", "▽"}
	
	result := ""
	for _, char := range text {
		if glitchRand.Float32() < 0.1 { // 10% chance of glitch
			result += string(glitchChars[glitchRand.Intn(len(glitchChars))])
		} else {
			result += string(char)
		}
	}
	return result
}

func RGBSplitText(text string) string {
	colors := []lipgloss.Color{GlitchRed, GlitchGreen, GlitchBlue}
	result := ""
	
	for i, char := range text {
		color := colors[i%len(colors)]
		result += lipgloss.NewStyle().Foreground(color).Bold(true).Render(string(char))
	}
	return result
}

// Advanced Styles with Maximum Impact
func SynthwaveTitle(text string) string {
	return lipgloss.NewStyle().
		Foreground(SynthwavePink).
		Background(DeepSpace).
		Bold(true).
		Italic(true).
		Underline(true).
		Padding(1, 3).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(SynthwaveCyan).
		AlignHorizontal(lipgloss.Center).
		Render(text)
}

func GlitchTitle(text string) string {
	normal := lipgloss.NewStyle().
		Foreground(SynthwavePink).
		Bold(true).
		Render(text)
	
	glitch1 := lipgloss.NewStyle().
		Foreground(SynthwaveCyan).
		Bold(true).
		Background(lipgloss.Color("#FF00FF")).
		Render(GlitchText(text))
	
	glitch2 := lipgloss.NewStyle().
		Foreground(SynthwaveYellow).
		Bold(true).
		Render(RGBSplitText(text))
	
	// Randomly return glitched version
	switch glitchRand.Intn(4) {
	case 0:
		return glitch1
	case 1:
		return glitch2
	default:
		return normal
	}
}

func FocusBox(text string, borderColor lipgloss.Color) string {
	return lipgloss.NewStyle().
		Foreground(SynthwaveCyan).
		Background(DarkVoid).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(borderColor).
		Bold(true).
		Render(text)
}

func CyberGridBox(text string) string {
	return lipgloss.NewStyle().
		Foreground(SynthwaveGreen).
		Background(CyberGrid).
		Padding(0, 1).
		Border(lipgloss.NormalBorder()).
		BorderForeground(SynthwaveCyan).
		Render("▶ " + text)
}

func HolographicText(text string) string {
	hologramColors := []lipgloss.Color{
		SynthwavePink, SynthwaveCyan, SynthwavePurple, 
		SynthwaveYellow, SynthwaveGreen, SynthwaveOrange,
	}
	
	result := ""
	for i, char := range text {
		color := hologramColors[i%len(hologramColors)]
		result += lipgloss.NewStyle().
			Foreground(color).
			Bold(true).
			Background(DarkVoid).
			Render(string(char))
	}
	return result
}

func DigitalRain(text string) string {
	return lipgloss.NewStyle().
		Foreground(SynthwaveGreen).
		Background(DeepSpace).
		Bold(true).
		Italic(true).
		Render("▓ " + text + " ▓")
}

func PriorityExplosion(priority string) string {
	switch priority {
	case "high":
		return lipgloss.NewStyle().
			Foreground(SynthwaveRed).
			Background(DarkVoid).
			Bold(true).
			Blink(true).
			Render("🔥 HIGH PRIORITY 🔥")
	case "medium":
		return lipgloss.NewStyle().
			Foreground(SynthwaveYellow).
			Background(CyberGrid).
			Bold(true).
			Render("⚡ MEDIUM ⚡")
	case "low":
		return lipgloss.NewStyle().
			Foreground(SynthwaveGreen).
			Background(DeepSpace).
			Bold(true).
			Render("💤 LOW PRIORITY 💤")
	default:
		return lipgloss.NewStyle().
			Foreground(SynthwavePurple).
			Render("◉ NORMAL")
	}
}

func TaskStatus(status string) string {
	if status == "completed" {
		return lipgloss.NewStyle().
			Foreground(SynthwaveGreen).
			Background(DarkVoid).
			Bold(true).
			Render("✅ MISSION COMPLETE ✅")
	} else {
		return lipgloss.NewStyle().
			Foreground(SynthwaveCyan).
			Background(CyberGrid).
			Bold(true).
			Render("◯ ACTIVE MISSION ◯")
	}
}

func CyberTag(tag string) string {
	return lipgloss.NewStyle().
		Foreground(SynthwavePink).
		Background(DarkVoid).
		Padding(0, 1).
		Border(lipgloss.NormalBorder()).
		BorderForeground(SynthwavePurple).
		Render("🏷️ " + tag)
}

func GlitchFooter(text string) string {
	artifacts := []string{"▓", "░", "▒", "█", "▄", "▀", "■", "□"}
	artifact := artifacts[glitchRand.Intn(len(artifacts))]
	
	return lipgloss.NewStyle().
		Foreground(SynthwavePurple).
		Italic(true).
		Render(artifact + " " + text + " " + artifact)
}

func MatrixHeader() string {
	return lipgloss.NewStyle().
		Foreground(SynthwaveGreen).
		Background(DeepSpace).
		Bold(true).
		AlignHorizontal(lipgloss.Center).
		Padding(1).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(SynthwaveCyan).
		Render(`
╔══════════════════════════════════════╗
║     SYNTHWAVE MISSION MATRIX v2.0     ║
║  ▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰  ║
╚══════════════════════════════════════╝`)
}

func CyberStats(active, completed, total int) string {
	return lipgloss.NewStyle().
		Foreground(SynthwaveCyan).
		Background(DarkVoid).
		Bold(true).
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(SynthwavePink).
		Render(
			"📊 GRID STATUS: " +
			lipgloss.NewStyle().Foreground(SynthwaveYellow).Render("⚡ "+string(rune(active))) +
			" ACTIVE | " +
			lipgloss.NewStyle().Foreground(SynthwaveGreen).Render("✅ "+string(rune(completed))) +
			" COMPLETED | " +
			lipgloss.NewStyle().Foreground(SynthwavePurple).Render("🌟 "+string(rune(total))) +
			" TOTAL",
		)
}

func LoadingAnimation() string {
	frames := []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}
	frame := frames[glitchRand.Intn(len(frames))]
	
	return lipgloss.NewStyle().
		Foreground(SynthwaveCyan).
		Bold(true).
		Render(frame + " LOADING SYNTHWAVE INTERFACE " + frame)
}

func EmptyStateMessage() string {
	return lipgloss.NewStyle().
		Foreground(SynthwaveCyan).
		Background(DeepSpace).
		Bold(true).
		AlignHorizontal(lipgloss.Center).
		Padding(2).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(SynthwavePink).
		Render(`
🔮 THE MATRIX IS EMPTY 🔮
🌌 Digital void awaits your missions
✨ Initialize with 'focus add "mission description"'
▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰▰`)
}

// Pterm Integration for Maximum Charm Library Showcase
func SynthwaveGetPriorityColor(priority string) *pterm.Style {
	switch priority {
	case "high":
		return &pterm.Style{
			pterm.FgRed,
			pterm.Bold,
			pterm.BgLightBlue,
		}
	case "medium":
		return &pterm.Style{
			pterm.FgYellow,
			pterm.Bold,
			pterm.BgMagenta,
		}
	case "low":
		return &pterm.Style{
			pterm.FgGreen,
			pterm.Bold,
			pterm.BgCyan,
		}
	default:
		return &pterm.Style{
			pterm.FgCyan,
			pterm.Bold,
			pterm.BgBlue,
		}
	}
}

func SynthwaveAIResponseStyle(response string) string {
	return lipgloss.NewStyle().
		Foreground(SynthwaveCyan).
		Background(DarkVoid).
		Bold(true).
		Italic(true).
		Padding(1).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(SynthwavePurple).
		Render("🤖 AI SYNTHESIS: " + response)
}

func SynthwaveReportStyle(text string) string {
	return lipgloss.NewStyle().
		Foreground(SynthwavePink).
		Background(CyberGrid).
		Padding(2).
		Border(lipgloss.DoubleBorder()).
		BorderForeground(SynthwaveCyan).
		Bold(true).
		Render("📋 ANALYSIS REPORT:\n\n" + text)
}

// Utility Functions for Visual Effects
func RandomGlitchChar() string {
	chars := []string{"@", "#", "$", "%", "&", "*", "■", "▪", "▫", "◆", "◇"}
	return chars[glitchRand.Intn(len(chars))]
}

func CreateDigitalArtifact() string {
	artifacts := []string{
		"▓▒░", "░▒▓", "█▀▄", "▄▀█", "◆◇◆", "◇◆◇", 
		"▰▱▰", "▱▰▱", "⚡⚡⚡", "✨✨✨", "🔥🔥🔥", "💫💫💫",
	}
	return artifacts[glitchRand.Intn(len(artifacts))]
}
