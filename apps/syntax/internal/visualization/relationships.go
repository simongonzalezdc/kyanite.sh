package visualization

import (
	"fmt"
	"strings"

	"github.com/kyanite/syntax/internal/character"
)

// GenerateRelationshipMap creates an ASCII art relationship map
func GenerateRelationshipMap(characters map[string]*character.Character) string {
	if len(characters) == 0 {
		return "No characters to visualize"
	}

	var b strings.Builder

	b.WriteString("┌────────────────────────────────────────────────────────────┐\n")
	b.WriteString("│             CHARACTER RELATIONSHIP MAP                     │\n")
	b.WriteString("└────────────────────────────────────────────────────────────┘\n\n")

	// Create list of characters
	charList := make([]*character.Character, 0, len(characters))
	for _, char := range characters {
		charList = append(charList, char)
	}

	// Show each character and their relationships
	for i, char := range charList {
		// Character box
		b.WriteString("┌")
		b.WriteString(strings.Repeat("─", 30))
		b.WriteString("┐\n")

		// Name
		name := truncate(char.Name, 28)
		padding := (28 - len(name)) / 2
		b.WriteString("│")
		b.WriteString(strings.Repeat(" ", padding))
		b.WriteString(name)
		b.WriteString(strings.Repeat(" ", 28-padding-len(name)))
		b.WriteString("│\n")

		// Role if exists
		if char.Role != "" {
			role := truncate(char.Role, 28)
			b.WriteString("│ ")
			b.WriteString(role)
			b.WriteString(strings.Repeat(" ", 28-len(role)))
			b.WriteString("│\n")
		}

		b.WriteString("└")
		b.WriteString(strings.Repeat("─", 30))
		b.WriteString("┘\n")

		// Relationships
		if len(char.Relationships) > 0 {
			for _, rel := range char.Relationships {
				// Find the related character
				relatedChar, exists := characters[rel.CharacterID]
				if !exists {
					continue
				}

				// Draw arrow based on tension
				var arrow, tension string
				switch rel.Tension {
				case "low":
					arrow = "  ───▶"
					tension = "✓ Low"
				case "medium":
					arrow = "  ═══▶"
					tension = "! Med"
				case "high":
					arrow = "  ━━━▶"
					tension = "✗ High"
				default:
					arrow = "  ───▶"
					tension = "? Unknown"
				}

				b.WriteString(arrow)
				b.WriteString(" ")

				// Relationship type and target
				relType := truncate(rel.Type, 15)
				targetName := truncate(relatedChar.Name, 15)

				b.WriteString(fmt.Sprintf("%s: %s [%s]\n", relType, targetName, tension))
			}
			b.WriteString("\n")
		} else {
			b.WriteString("  (No relationships)\n\n")
		}

		// Add spacing between characters
		if i < len(charList)-1 {
			b.WriteString("\n")
		}
	}

	// Add legend
	b.WriteString("\n")
	b.WriteString("┌─ LEGEND ───────────────────────────────┐\n")
	b.WriteString("│  ───▶  Low Tension    (Friendly)       │\n")
	b.WriteString("│  ═══▶  Medium Tension (Complicated)    │\n")
	b.WriteString("│  ━━━▶  High Tension   (Hostile)        │\n")
	b.WriteString("└────────────────────────────────────────┘\n")

	return b.String()
}

// GenerateRelationshipMatrix creates a matrix view of all relationships
func GenerateRelationshipMatrix(characters map[string]*character.Character) string {
	if len(characters) == 0 {
		return "No characters to visualize"
	}

	charList := make([]*character.Character, 0, len(characters))
	for _, char := range characters {
		charList = append(charList, char)
	}

	var b strings.Builder

	b.WriteString("CHARACTER RELATIONSHIP MATRIX\n\n")

	// Header
	b.WriteString("         ")
	for _, char := range charList {
		abbrev := getAbbreviation(char.Name)
		b.WriteString(fmt.Sprintf("│ %s ", abbrev))
	}
	b.WriteString("│\n")

	// Separator
	b.WriteString("─────────")
	for range charList {
		b.WriteString("┼────")
	}
	b.WriteString("┤\n")

	// Rows
	for _, fromChar := range charList {
		abbrev := getAbbreviation(fromChar.Name)
		b.WriteString(fmt.Sprintf("%-8s ", abbrev))

		for _, toChar := range charList {
			b.WriteString("│ ")

			if fromChar.ID == toChar.ID {
				b.WriteString("── ") // Self
			} else {
				// Find relationship
				var tension string
				found := false
				for _, rel := range fromChar.Relationships {
					if rel.CharacterID == toChar.ID {
						switch rel.Tension {
						case "low":
							tension = "✓ "
						case "medium":
							tension = "! "
						case "high":
							tension = "✗ "
						default:
							tension = "? "
						}
						found = true
						break
					}
				}
				if !found {
					tension = "·  "
				}
				b.WriteString(tension)
			}
		}
		b.WriteString("│ ")
		b.WriteString(fromChar.Name)
		b.WriteString("\n")
	}

	// Legend
	b.WriteString("\n")
	b.WriteString("Legend: ✓ Low | ! Medium | ✗ High | · None | ── Self\n")

	return b.String()
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func getAbbreviation(name string) string {
	parts := strings.Fields(name)
	if len(parts) == 0 {
		return "??"
	}
	if len(parts) == 1 {
		if len(name) >= 2 {
			return name[:2]
		}
		return name
	}
	// Use initials for multiple names
	abbrev := ""
	for _, part := range parts {
		if len(part) > 0 {
			abbrev += string(part[0])
		}
	}
	if len(abbrev) > 2 {
		abbrev = abbrev[:2]
	}
	return abbrev
}
