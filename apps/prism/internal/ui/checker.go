package ui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/design"
	"github.com/kyanite/prism/internal/clipboard"
	"github.com/kyanite/prism/internal/color"
	"github.com/kyanite/prism/internal/theme"
)

// ContrastResult holds the result of a WCAG contrast check.
type ContrastResult struct {
	Ratio     float64
	Level     string
	PassedAA  bool
	PassedAAA bool
}

// ContrastResultMsg carries the result of a contrast check back to Update.
type ContrastResultMsg struct {
	Result *ContrastResult
	Err    string
}

// CheckerModel represents the WCAG accessibility checker screen
type CheckerModel struct {
	themeManager *theme.Manager
	styles       Styles
	foreground   string
	background   string
	result       *ContrastResult
	err          string
	status       string
}

// NewCheckerModel creates a new checker model
func NewCheckerModel(tm *theme.Manager) CheckerModel {
	return CheckerModel{
		themeManager: tm,
		styles:       NewStyles(tm.CurrentTheme()),
		foreground:   "#FFFFFF",
		background:   "#000000",
	}
}

// Init initializes the checker
func (m CheckerModel) Init() tea.Cmd {
	return nil
}

// Update handles checker messages
func (m CheckerModel) Update(msg tea.Msg) (CheckerModel, tea.Cmd) {
	switch msg := msg.(type) {
	case ContrastResultMsg:
		if msg.Err != "" {
			m.err = msg.Err
			m.result = nil
		} else {
			m.result = msg.Result
			m.err = ""
		}
		return m, nil
	case tea.KeyMsg:
		switch msg.String() {
		case "enter", " ":
			return m, m.check()
		case "c":
			if m.result != nil {
				m.copyResult()
			}
		}
	}

	return m, nil
}

func (m *CheckerModel) copyResult() {
	text := fmt.Sprintf("Contrast: %.2f:1 - WCAG %s (FG: %s, BG: %s)",
		m.result.Ratio, m.result.Level, m.foreground, m.background)

	if err := clipboard.Copy(text); err == nil {
		m.status = "✓ Copied contrast result to clipboard"
	} else {
		m.status = "✗ Clipboard unavailable"
	}
	m.err = ""
}

// View renders the checker
func (m CheckerModel) View() string {
	styles := NewStyles(m.themeManager.CurrentTheme())

	var b strings.Builder

	// Title
	title := styles.Title.Render("WCAG Accessibility Checker")
	b.WriteString(title)
	b.WriteString("\n\n")

	// Input colors
	b.WriteString(styles.Primary.Render("Colors to Check:"))
	b.WriteString("\n")

	fgSwatch := lipgloss.Style{}.
		Background(lipgloss.Color(m.foreground)).
		Padding(design.SpacingNone, design.SpacingS).
		Render("██")
	fmt.Fprintf(&b, "Foreground: %s %s\n", fgSwatch, m.foreground)

	bgSwatch := lipgloss.Style{}.
		Background(lipgloss.Color(m.background)).
		Padding(design.SpacingNone, design.SpacingS).
		Render("██")
	fmt.Fprintf(&b, "Background: %s %s\n", bgSwatch, m.background)

	b.WriteString("\n")

	// Result
	if m.result != nil {
		b.WriteString(styles.Secondary.Render("Contrast Ratio:"))
		b.WriteString("\n")

		// Contrast ratio
		ratioText := fmt.Sprintf("%.2f:1", m.result.Ratio)
		var ratioStyle lipgloss.Style
		if m.result.PassedAAA {
			ratioStyle = styles.Success
		} else if m.result.PassedAA {
			ratioStyle = styles.Accent
		} else {
			ratioStyle = styles.Error
		}
		b.WriteString(ratioStyle.Bold(true).Render(ratioText))
		b.WriteString("\n\n")

		// WCAG levels
		b.WriteString("WCAG Compliance:\n")
		if m.result.PassedAAA {
			b.WriteString(styles.Success.Render("✓ AAA (7:1) - Excellent"))
		} else if m.result.PassedAA {
			b.WriteString(styles.Accent.Render("✓ AA (4.5:1) - Good"))
		} else {
			b.WriteString(styles.Error.Render("✗ FAIL - Insufficient contrast"))
		}
		b.WriteString("\n\n")

		// Sample text
		b.WriteString(styles.Secondary.Render("Sample:"))
		b.WriteString("\n")

		sample := lipgloss.Style{}.
			Foreground(lipgloss.Color(m.foreground)).
			Background(lipgloss.Color(m.background)).
			Padding(design.SpacingXS, design.SpacingS).
			Render("The quick brown fox jumps over the lazy dog")
		b.WriteString(sample)
		b.WriteString("\n\n")
	}

	// Error
	if m.err != "" {
		b.WriteString(styles.Error.Render("Error: " + m.err))
		b.WriteString("\n\n")
	}

	// Status
	if m.status != "" {
		b.WriteString(styles.Success.Render(m.status))
		b.WriteString("\n\n")
	}

	// Help
	var help string
	if m.result != nil {
		help = styles.Muted.Render("Enter: Check Contrast • C: Copy • Esc: Menu")
	} else {
		help = styles.Muted.Render("Enter: Check Contrast • Esc: Menu")
	}
	b.WriteString(help)

	// Wrap in border
	content := styles.Border.Width(ContentWidth).Render(b.String())

	return lipgloss.Place(ScreenWidth, ScreenHeight, lipgloss.Center, lipgloss.Center, content)
}

// check calculates contrast using the design module
func (m CheckerModel) check() tea.Cmd {
	return func() tea.Msg {
		fg, err := color.ParseHex(m.foreground)
		if err != nil {
			return ContrastResultMsg{Err: "Invalid foreground color"}
		}

		bg, err := color.ParseHex(m.background)
		if err != nil {
			return ContrastResultMsg{Err: "Invalid background color"}
		}

		ratio, err := design.ContrastRatio(fg.Hex, bg.Hex)
		if err != nil {
			return ContrastResultMsg{Err: "Contrast calculation failed"}
		}

		passedAA := ratio >= design.ContrastRatioAA
		passedAAA := ratio >= design.ContrastRatioAAA

		var level string
		if passedAAA {
			level = "AAA"
		} else if passedAA {
			level = "AA"
		} else {
			level = "FAIL"
		}

		result := ContrastResult{Ratio: ratio, Level: level, PassedAA: passedAA, PassedAAA: passedAAA}
		return ContrastResultMsg{Result: &result}
	}
}
