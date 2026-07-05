package cli

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/focus/internal/journal"
	"github.com/kyanite/focus/pkg/models"
	"github.com/kyanite/focus/pkg/styles"
	"github.com/kyanite/focus/pkg/utils"
	"github.com/spf13/cobra"
)

var journalCmd = &cobra.Command{
	Use:   "journal",
	Short: "Manage journal entries",
	Long: styles.GetBoxStyle().Render("📝 Journal - Record thoughts, reflections, and insights\n\n" +
		"Commands:\n" +
		"  new     Create a new journal entry\n" +
		"  list    List all journal entries\n" +
		"  view    View a specific journal entry\n" +
		"  search  Search journal entries\n" +
		"  export  Export entry to syntax.sh"),
}

var (
	journalDate     string
	journalTitle    string
	journalMood     string
	journalTemplate string
	journalTags     []string
	journalID       string
	journalQuery    string
	exportTypeStr   string
	toSyntax        bool
)

var journalNewCmd = &cobra.Command{
	Use:   "new",
	Short: "Create a new journal entry",
	Long:  "Create a new journal entry with optional template, tags, and metadata.",
	Run: func(cmd *cobra.Command, args []string) {
		entry := createJournalEntry()
		if entry == nil {
			return
		}

		storage := utils.NewJournalStorage()
		if err := storage.AddEntry(entry); err != nil {
			fmt.Printf("❌ Error saving journal entry: %v\n", err)
			return
		}

		// Success message
		successStyle := lipgloss.Style{}.
			Foreground(styles.GetSuccess()).
			Bold(true).
			Render("✅ Journal entry created successfully")
		fmt.Println(successStyle)

		// Entry details
		detailStyle := lipgloss.Style{}.
			Foreground(styles.GetForeground()).
			Background(styles.GetPanel()).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.GetBorder())

		details := fmt.Sprintf(
			"📅 Date: %s\n📝 Title: %s\n😊 Mood: %s\n🏷️  Tags: %s\n📊 Words: %d",
			entry.Date,
			entry.Title,
			entry.Mood,
			strings.Join(entry.Tags, ", "),
			entry.WordCount,
		)

		fmt.Println(detailStyle.Render(details))
	},
}

var journalListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all journal entries",
	Long:  "List all journal entries with their dates, titles, and word counts.",
	Run: func(cmd *cobra.Command, args []string) {
		storage := utils.NewJournalStorage()
		entries, err := storage.LoadEntries()
		if err != nil {
			fmt.Printf("❌ Error loading journal entries: %v\n", err)
			return
		}

		if len(entries) == 0 {
			emptyStyle := lipgloss.Style{}.
				Foreground(styles.GetWarning()).
				Bold(true).
				Render("📭 No journal entries found")
			fmt.Println(emptyStyle)
			fmt.Println("\n💡 Create your first entry with: focus journal new")
			return
		}

		// Header
		headerStyle := lipgloss.Style{}.
			Foreground(styles.GetAccent()).
			Bold(true).
			Render("📚 Journal Entries")
		fmt.Println(headerStyle)
		fmt.Println()

		// List entries (show most recent 10)
		start := 0
		if len(entries) > 10 {
			start = len(entries) - 10
		}
		for i := start; i < len(entries); i++ {
			entry := entries[i]
			entryStyle := lipgloss.Style{}.
				Foreground(styles.GetForeground()).
				Background(styles.GetPanel()).
				Padding(0, 1).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(styles.GetBorder())

			title := entry.Title
			if title == "" {
				title = entry.Content[:min(50, len(entry.Content))]
				if len(entry.Content) > 50 {
					title += "..."
				}
			}

			line := fmt.Sprintf("📅 %s | 📝 %s | 📊 %d words | 😊 %s",
				entry.Date, title, entry.WordCount, entry.Mood)

			fmt.Println(entryStyle.Render(line))
		}

		if len(entries) > 10 {
			moreStyle := lipgloss.Style{}.
				Foreground(styles.GetWarning()).
				Italic(true).
				Render(fmt.Sprintf("... and %d more entries", len(entries)-10))
			fmt.Println(moreStyle)
		}
	},
}

var journalViewCmd = &cobra.Command{
	Use:   "view [date]",
	Short: "View a specific journal entry",
	Long:  "View a journal entry by date. If no date is provided, shows today's entry.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		date := journalDate
		if date == "" {
			date = time.Now().Format("2006-01-02")
		}

		storage := utils.NewJournalStorage()
		entry, err := storage.GetEntryByDate(date)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			fmt.Printf("💡 Available dates: focus journal list\n")
			return
		}

		// Display entry
		displayJournalEntry(entry)
	},
}

var journalSearchCmd = &cobra.Command{
	Use:   "search [keyword]",
	Short: "Search journal entries",
	Long:  "Search journal entries by keyword in content, title, or tags.",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		keyword := strings.Join(args, " ")
		storage := utils.NewJournalStorage()

		entries, err := storage.SearchEntries(keyword)
		if err != nil {
			fmt.Printf("❌ Error searching entries: %v\n", err)
			return
		}

		if len(entries) == 0 {
			noResultsStyle := lipgloss.Style{}.
				Foreground(styles.GetWarning()).
				Bold(true).
				Render(fmt.Sprintf("🔍 No results found for: %s", keyword))
			fmt.Println(noResultsStyle)
			return
		}

		// Results header
		resultsStyle := lipgloss.Style{}.
			Foreground(styles.GetSuccess()).
			Bold(true).
			Render(fmt.Sprintf("🔍 Found %d entries for: %s", len(entries), keyword))
		fmt.Println(resultsStyle)
		fmt.Println()

		// Display results
		for _, entry := range entries {
			resultStyle := lipgloss.Style{}.
				Foreground(styles.GetForeground()).
				Background(styles.GetPanel()).
				Padding(0, 1).
				Border(lipgloss.RoundedBorder()).
				BorderForeground(styles.GetBorder())

			line := fmt.Sprintf("📅 %s | 📝 %s | 📊 %d words",
				entry.Date, entry.Title, entry.WordCount)
			fmt.Println(resultStyle.Render(line))
		}
	},
}

var journalExportCmd = &cobra.Command{
	Use:   "export [date] --to-syntax --type [type]",
	Short: "Export journal entry to syntax.sh",
	Long:  "Export a journal entry to syntax.sh format for story development.",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		date := journalDate
		if date == "" {
			date = time.Now().Format("2006-01-02")
		}

		exportTypeStr, err := cmd.Flags().GetString("type")
		if err != nil {
			exportTypeStr = "dialogue" // default
		}
		toSyntax, err := cmd.Flags().GetBool("to-syntax")
		if err != nil {
			toSyntax = false
		}

		if !toSyntax {
			fmt.Println("❌ Error: --to-syntax flag is required")
			fmt.Println("💡 Usage: focus journal export <date> --to-syntax --type <type>")
			return
		}

		// Parse export type
		var exportType models.ExportType
		switch exportTypeStr {
		case "character":
			exportType = models.ExportCharacter
		case "dialogue":
			exportType = models.ExportDialogue
		case "scene":
			exportType = models.ExportScene
		case "research":
			exportType = models.ExportResearch
		default:
			fmt.Printf("❌ Error: Invalid export type '%s'\n", exportTypeStr)
			fmt.Println("💡 Valid types: character, dialogue, scene, research")
			return
		}

		// Get entry
		storage := utils.NewJournalStorage()
		entry, err := storage.GetEntryByDate(date)
		if err != nil {
			fmt.Printf("❌ Error: %v\n", err)
			return
		}

		// Export to syntax.sh
		exporter := journal.NewExporter()
		if err := exporter.ExportToSyntax(entry, exportType); err != nil {
			fmt.Printf("❌ Error exporting to syntax.sh: %v\n", err)
			return
		}

		// Success message
		successStyle := lipgloss.Style{}.
			Foreground(styles.GetSuccess()).
			Bold(true).
			Render("✅ Exported to syntax.sh")

		filename := entry.GetExportFilename(exportType)
		detailsStyle := lipgloss.Style{}.
			Foreground(styles.GetForeground()).
			Background(styles.GetPanel()).
			Padding(1, 2).
			Border(lipgloss.RoundedBorder()).
			BorderForeground(styles.GetBorder())

		details := fmt.Sprintf(
			"📁 File: %s\n📂 Path: %s\n📝 Type: %s",
			filename, exporter.GetExportPath(), exportTypeStr,
		)

		fmt.Println(successStyle)
		fmt.Println(detailsStyle.Render(details))
	},
}

// createJournalEntry creates a new journal entry interactively
func createJournalEntry() *models.JournalEntry {
	// Get date
	date := journalDate
	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	// Get template
	template := journalTemplate
	if template == "" {
		template = "daily_log"
	}

	// Find template
	var selectedTemplate models.JournalTemplate
	templateFound := false
	for _, t := range models.JournalTemplates {
		if t.Name == template {
			selectedTemplate = t
			templateFound = true
			break
		}
	}

	if !templateFound {
		fmt.Printf("⚠️ Template '%s' not found, using daily_log\n", template)
		for _, t := range models.JournalTemplates {
			if t.Name == "daily_log" {
				selectedTemplate = t
				break
			}
		}
	}

	// Get content
	fmt.Printf("📝 %s - %s\n", strings.ToUpper(template), selectedTemplate.Description)
	fmt.Println("Enter your journal entry (press Enter on an empty line to finish):")

	var content strings.Builder
	for {
		var line string
		fmt.Print("> ")
		_, _ = fmt.Scanln(&line) // Ignore error for user input
		if line == "" {
			break
		}
		content.WriteString(line + "\n")
	}

	contentStr := strings.TrimSpace(content.String())
	if contentStr == "" {
		fmt.Println("❌ Empty entry cancelled")
		return nil
	}

	// Create entry
	entry := models.NewJournalEntry(
		date,
		journalTitle,
		contentStr,
		journalMood,
		template,
		journalTags,
	)

	return entry
}

// displayJournalEntry displays a journal entry in a formatted way
func displayJournalEntry(entry *models.JournalEntry) {
	// Header
	headerStyle := lipgloss.Style{}.
		Foreground(styles.GetAccent()).
		Bold(true).
		Render(fmt.Sprintf("📅 Journal Entry - %s", entry.Date))
	fmt.Println(headerStyle)

	// Metadata
	metadataStyle := lipgloss.Style{}.
		Foreground(styles.GetForeground()).
		Background(styles.GetPanel()).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.GetBorder())

	metadata := []string{}
	if entry.Title != "" {
		metadata = append(metadata, fmt.Sprintf("📝 Title: %s", entry.Title))
	}
	if entry.Mood != "" {
		metadata = append(metadata, fmt.Sprintf("😊 Mood: %s", entry.Mood))
	}
	if len(entry.Tags) > 0 {
		metadata = append(metadata, fmt.Sprintf("🏷️  Tags: %s", strings.Join(entry.Tags, ", ")))
	}
	metadata = append(metadata, fmt.Sprintf("📊 Words: %d", entry.WordCount))
	metadata = append(metadata, fmt.Sprintf("🕐 Updated: %s", entry.UpdatedAt.Format("2006-01-02 15:04")))

	fmt.Println(metadataStyle.Render(strings.Join(metadata, "\n")))

	// Content
	contentStyle := lipgloss.Style{}.
		Foreground(styles.GetForeground()).
		Background(styles.GetPanel()).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.GetBorder()).
		Width(80)

	fmt.Println(contentStyle.Render(entry.Content))
}

// min returns the minimum of two integers
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func init() {
	journalCmd.AddCommand(journalNewCmd)
	journalCmd.AddCommand(journalListCmd)
	journalCmd.AddCommand(journalViewCmd)
	journalCmd.AddCommand(journalSearchCmd)
	journalCmd.AddCommand(journalExportCmd)

	// Flags for journal new
	journalNewCmd.Flags().StringVarP(&journalDate, "date", "d", "", "Date of the entry (YYYY-MM-DD, default: today)")
	journalNewCmd.Flags().StringVarP(&journalTitle, "title", "t", "", "Title of the entry")
	journalNewCmd.Flags().StringVarP(&journalMood, "mood", "m", "", "Mood of the entry")
	journalNewCmd.Flags().StringVarP(&journalTemplate, "template", "", "", "Template to use (morning_pages, evening_reflection, daily_log)")
	journalNewCmd.Flags().StringSliceVarP(&journalTags, "tags", "g", []string{}, "Tags for the entry")

	// Flags for journal view
	journalViewCmd.Flags().StringVarP(&journalDate, "date", "d", "", "Date of the entry (YYYY-MM-DD, default: today)")
	journalViewCmd.Flags().StringVarP(&journalID, "id", "i", "", "ID of the entry")

	// Flags for journal search
	journalSearchCmd.Flags().StringVarP(&journalQuery, "query", "q", "", "Search query")
	journalSearchCmd.Flags().StringSliceVarP(&journalTags, "tags", "g", []string{}, "Filter by tags")
	journalSearchCmd.Flags().StringVarP(&journalMood, "mood", "m", "", "Filter by mood")

	// Flags for journal export
	journalExportCmd.Flags().StringVarP(&journalID, "id", "i", "", "ID of the entry to export")
	journalExportCmd.Flags().StringVarP(&exportTypeStr, "type", "t", "markdown", "Export type (markdown, syntax)")
	journalExportCmd.Flags().BoolVarP(&toSyntax, "to-syntax", "s", false, "Export to syntax.sh format")

	// journalCmd is now registered in root.go to avoid duplication
}
