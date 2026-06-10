package app

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/kyanite/syntax/internal/editor"
	"github.com/kyanite/syntax/internal/scene"
	"github.com/kyanite/syntax/internal/storage"
	"github.com/kyanite/syntax/internal/utils"
)

func (m Model) viewScenes() string {
	if m.CurrentProject == nil {
		return "No project loaded"
	}

	var b strings.Builder

	// Title bar
	titleBar := m.Styles.StatusBar.Render(fmt.Sprintf(" %s - Scenes ", m.CurrentProject.Title))
	b.WriteString(titleBar)
	b.WriteString("\n\n")

	// Data is loaded in Update via ensureDataLoaded()

	b.WriteString(m.Styles.Heading.Render("📝 Scenes"))
	b.WriteString("\n\n")

	if len(m.CurrentProject.Scenes) == 0 {
		b.WriteString(m.Styles.Text.Render("No scenes yet. Press 'n' to create one."))
	} else {
		for _, sc := range m.CurrentProject.Scenes {
			b.WriteString(m.Styles.Accent.Render(fmt.Sprintf("• Ch%d Sc%d: %s", sc.Chapter, sc.SceneNumber, sc.Name)))
			b.WriteString(m.Styles.Text.Render(fmt.Sprintf(" [%s]", sc.Status)))
			b.WriteString("\n")
			b.WriteString(m.Styles.Text.Render(fmt.Sprintf("  %d words", sc.WordCount)))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")

	// Show input prompt if in input mode
	if m.InputMode {
		b.WriteString(m.Styles.Accent.Render(fmt.Sprintf("Enter scene name: %s█", m.InputValue)))
		b.WriteString("\n")
	} else {
		b.WriteString(m.Styles.Text.Render("n - New Scene | v - Validate | d - Delete | Enter - Edit | ↑/↓ - Navigate | Esc - Back"))
	}

	if m.Message != "" {
		b.WriteString("\n\n")
		b.WriteString(m.Styles.Success.Render(m.Message))
	}

	return b.String()
}

func (m Model) handleScenesKeys(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Handle input mode for scene creation
	if m.InputMode {
		switch msg.String() {
		case "enter":
			// Create scene with user-provided name
			if m.InputValue != "" {
				sc := &scene.Scene{
					Chapter:     1, // Default chapter
					SceneNumber: len(m.CurrentProject.Scenes) + 1,
					Name:        m.InputValue,
					Status:      "draft",
				}

				err := storage.SaveScene(m.CurrentProject.Directory, sc)
				if err != nil {
					m.Error = err
				} else {
					if m.CurrentProject.Scenes == nil {
						m.CurrentProject.Scenes = make(map[string]*scene.Scene)
					}
					m.CurrentProject.Scenes[sc.ID] = sc
					m.CurrentProject.TotalScenes++
					m.Message = fmt.Sprintf("Created scene: %s", sc.Name)

					// Save updated project metadata
					if err := storage.SaveProjectMetadata(m.CurrentProject); err != nil {
						m.Error = fmt.Errorf("failed to save project: %w", err)
					}
				}
			}
			m.InputMode = false
			m.InputValue = ""
			return m, nil

		case "esc":
			m.InputMode = false
			m.InputValue = ""
			return m, nil

		case "backspace":
			if len(m.InputValue) > 0 {
				m.InputValue = m.InputValue[:len(m.InputValue)-1]
			}
			return m, nil

		default:
			// Add character to input
			if len(msg.String()) == 1 {
				m.InputValue += msg.String()
			}
			return m, nil
		}
	}

	// Normal mode
	switch msg.String() {
	case "up", "k":
		if m.SelectedIndex > 0 {
			m.SelectedIndex--
		}
		return m, nil

	case "down", "j":
		sceneCount := len(m.CurrentProject.Scenes)
		if m.SelectedIndex < sceneCount-1 {
			m.SelectedIndex++
		}
		return m, nil

	case "enter":
		// Open scene in editor
		sc := utils.GetSceneAtIndex(m.CurrentProject.Scenes, m.SelectedIndex)
		if sc != nil {
			m.CurrentScene = sc
			m.Buffer = editor.NewBuffer(sc.Content)
			m.EditorMode = EditorModeNormal
			m.CurrentScreen = ScreenTextEditor
		}
		return m, nil

	case "n":
		// Enter input mode for scene creation
		m.InputMode = true
		m.InputValue = ""
		return m, nil

	case "d":
		// Delete selected scene
		sceneID := utils.FindSceneIDAtIndex(m.CurrentProject.Scenes, m.SelectedIndex)
		if sceneID != "" {
			sc := m.CurrentProject.Scenes[sceneID]
			// Delete the scene
			err := storage.DeleteScene(m.CurrentProject.Directory, sceneID)
			if err != nil {
				m.Error = err
				return m, nil
			}

			// Remove from project's scene map
			delete(m.CurrentProject.Scenes, sceneID)
			m.CurrentProject.TotalScenes--

			// Update project metadata
			if err := storage.SaveProjectMetadata(m.CurrentProject); err != nil {
				m.Error = fmt.Errorf("failed to save project: %w", err)
			} else {
				m.Message = fmt.Sprintf("Deleted scene: %s", sc.Name)
			}

			// Adjust selected index
			if m.SelectedIndex >= len(m.CurrentProject.Scenes) {
				m.SelectedIndex = len(m.CurrentProject.Scenes) - 1
			}
			if m.SelectedIndex < 0 {
				m.SelectedIndex = 0
			}
		}
		return m, nil

	case "v":
		// Show scene validation
		m.CurrentScreen = ScreenSceneValidation
		return m, nil

	case "esc":
		m.CurrentScreen = ScreenEditor
		m.Message = ""
		return m, nil
	}

	return m, nil
}
