package ui

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/prism/internal/theme"
)

// TheoryModel represents the color theory education screen
type TheoryModel struct {
	themeManager *theme.Manager
	styles       Styles
	lessons      []theoryLesson
	selected     int
	viewing      bool
}

type theoryLesson struct {
	title   string
	content string
}

// NewTheoryModel creates a new theory model
func NewTheoryModel(tm *theme.Manager) TheoryModel {
	lessons := []theoryLesson{
		{
			title: "Complementary Colors",
			content: "Complementary colors are opposite each other on the color wheel (180° apart).\n\n" +
				"Examples: Red & Green, Blue & Orange, Yellow & Purple\n\n" +
				"They create high contrast and vibrant looks when used together.",
		},
		{
			title: "Analogous Colors",
			content: "Analogous colors are adjacent to each other on the color wheel (±30°).\n\n" +
				"They create harmonious, serene, and comfortable designs.\n\n" +
				"Example: Blue, Blue-Green, Green",
		},
		{
			title: "Triadic Colors",
			content: "Triadic colors are evenly spaced around the color wheel (120° apart).\n\n" +
				"They create vibrant color schemes even with muted tones.\n\n" +
				"Example: Red, Yellow, Blue (primary colors)",
		},
		{
			title: "Warm vs Cool Colors",
			content: "Warm colors (red, orange, yellow) evoke warmth and energy.\n\n" +
				"Cool colors (blue, green, purple) evoke calmness and serenity.\n\n" +
				"Understanding temperature helps create mood in designs.",
		},
		{
			title: "Tints, Shades, and Tones",
			content: "Tint: Color + White (lighter)\n" +
				"Shade: Color + Black (darker)\n" +
				"Tone: Color + Gray (desaturated)\n\n" +
				"These create monochromatic palettes with depth.",
		},
	}

	return TheoryModel{
		themeManager: tm,
		styles:       NewStyles(tm.CurrentTheme()),
		lessons:      lessons,
		selected:     0,
		viewing:      false,
	}
}

// Init initializes the theory screen
func (m TheoryModel) Init() tea.Cmd {
	return nil
}

// Update handles theory messages
func (m TheoryModel) Update(msg tea.Msg) (TheoryModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if !m.viewing && m.selected > 0 {
				m.selected--
			}
		case "down", "j":
			if !m.viewing && m.selected < len(m.lessons)-1 {
				m.selected++
			}
		case "enter", " ":
			m.viewing = !m.viewing
		case "esc":
			if m.viewing {
				m.viewing = false
				return m, nil
			}
		}
	}

	return m, nil
}

// View renders the theory screen
func (m TheoryModel) View() string {
	styles := NewStyles(m.themeManager.CurrentTheme())

	var b strings.Builder

	// Title
	title := styles.Title.Render("Color Theory Lessons")
	b.WriteString(title)
	b.WriteString("\n\n")

	if m.viewing {
		// Show lesson content
		lesson := m.lessons[m.selected]

		lessonTitle := styles.Primary.Bold(true).Render(lesson.title)
		b.WriteString(lessonTitle)
		b.WriteString("\n\n")

		b.WriteString(lesson.content)
		b.WriteString("\n\n")

		help := styles.Muted.Render("Enter: Back to list • Esc: Menu")
		b.WriteString(help)
	} else {
		// Show lesson list
		b.WriteString(styles.Secondary.Render("Select a lesson:"))
		b.WriteString("\n")

		for i, lesson := range m.lessons {
			style := styles.Unselected
			cursor := "  "
			if i == m.selected {
				style = styles.Selected
				cursor = "▸ "
			}

			line := cursor + lesson.title
			b.WriteString(style.Render(line))
			b.WriteString("\n")
		}

		b.WriteString("\n")
		help := styles.Muted.Render("↑/↓: Navigate • Enter: View • Esc: Menu")
		b.WriteString(help)
	}

	// Wrap in border
	content := styles.Border.Width(ContentWidth).Render(b.String())

	return lipgloss.Place(ScreenWidth, ScreenHeight, lipgloss.Center, lipgloss.Center, content)
}
