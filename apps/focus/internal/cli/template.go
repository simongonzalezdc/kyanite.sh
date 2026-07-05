package cli

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/kyanite/focus/internal/engine"
	"github.com/kyanite/focus/internal/repository"
	"github.com/kyanite/focus/internal/templates"
	"github.com/kyanite/focus/pkg/models"
	"github.com/kyanite/focus/pkg/styles"
	"github.com/kyanite/focus/pkg/utils"
	"github.com/spf13/cobra"
)

var (
	templateName       string
	templatePriority   string
	templateCategories string
)

var templateCmd = &cobra.Command{
	Use:   "template",
	Short: "Manage task templates",
	Long:  `Create, list, and use task templates for quick task creation.`,
}

var templateListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all task templates",
	Run: func(cmd *cobra.Command, args []string) {
		templateStore := templates.NewStore(getTemplatePath())
		templateList, err := templateStore.Load()
		if err != nil {
			showTemplateError(fmt.Sprintf("Failed to load templates: %s", err.Error()))
			return
		}

		if len(templateList) == 0 {
			fmt.Println(lipgloss.Style{}.
				Foreground(styles.SynthwaveYellow).
				Render("No templates found. Create one with: focus template create"))
			return
		}

		title := styles.SynthwaveTitle("📋 TASK TEMPLATES")
		fmt.Println(title)
		fmt.Println()

		for i, tmpl := range templateList {
			fmt.Printf("%d. %s\n", i+1, lipgloss.Style{}.
				Foreground(styles.SynthwaveCyan).
				Bold(true).
				Render(tmpl.Name))

			fmt.Printf("   %s\n", lipgloss.Style{}.
				Foreground(styles.SynthwavePink).
				Render(tmpl.Description))

			if tmpl.Priority != "" {
				fmt.Printf("   Priority: %s\n", lipgloss.Style{}.
					Foreground(getPriorityColor(tmpl.Priority)).
					Render(strings.ToUpper(tmpl.Priority)))
			}

			if len(tmpl.Categories) > 0 {
				fmt.Printf("   Categories: %s\n", lipgloss.Style{}.
					Foreground(styles.SynthwaveYellow).
					Render(strings.Join(tmpl.Categories, ", ")))
			}

			fmt.Printf("   ID: %s\n", lipgloss.Style{}.
				Foreground(styles.DarkVoid).
				Render(tmpl.ID[:8]+"..."))
			fmt.Println()
		}
	},
}

var templateCreateCmd = &cobra.Command{
	Use:   "create [description]",
	Short: "Create a new task template",
	Long: `Create a reusable task template.
Examples:
  focus template create "Daily standup" --name="Standup" --priority=high
  focus template create "Weekly review" --name="Review" --categories=work,planning`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			showTemplateError("Please provide a template description")
			return
		}

		description := strings.Join(args, " ")

		if templateName == "" {
			showTemplateError("Please provide a template name with --name")
			return
		}

		// Parse categories
		var categories []string
		if templateCategories != "" {
			categories = strings.Split(templateCategories, ",")
			for i := range categories {
				categories[i] = strings.TrimSpace(categories[i])
			}
		}

		// Create template
		template := models.TaskTemplate{
			Name:        templateName,
			Description: description,
			Priority:    templatePriority,
			Categories:  categories,
		}

		templateStore := templates.NewStore(getTemplatePath())
		created, err := templateStore.Add(template)
		if err != nil {
			showTemplateError(fmt.Sprintf("Failed to create template: %s", err.Error()))
			return
		}

		showTemplateSuccess(created)
	},
}

var templateUseCmd = &cobra.Command{
	Use:   "use [template-name]",
	Short: "Create a task from a template",
	Long: `Use a template to quickly create a new task.
Examples:
  focus template use "Standup"
  focus template use "Review"`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			showTemplateError("Please provide a template name")
			return
		}

		name := strings.Join(args, " ")

		templateStore := templates.NewStore(getTemplatePath())
		template, err := templateStore.GetByName(name)
		if err != nil {
			showTemplateError(fmt.Sprintf("Template '%s' not found", name))
			return
		}

		// Create task from template
		repo := repository.NewStoreRepository(utils.GetStoragePath())
		taskEngine := engine.New(repo)

		task, err := taskEngine.AddTask(template.ToTask())
		if err != nil {
			showTemplateError(fmt.Sprintf("Failed to create task: %s", err.Error()))
			return
		}

		fmt.Println(styles.SynthwaveTitle("✨ TASK CREATED FROM TEMPLATE"))
		fmt.Println()
		fmt.Printf("%s %s\n", lipgloss.Style{}.Foreground(styles.SynthwaveCyan).Render("▸ ID:"),
			task.ID[:8]+"...")
		fmt.Printf("%s %s\n", lipgloss.Style{}.Foreground(styles.SynthwavePink).Render("▸ Description:"),
			task.Description)
		fmt.Printf("%s %s\n", lipgloss.Style{}.Foreground(styles.SynthwaveYellow).Render("▸ Template:"),
			template.Name)
	},
}

var templateDeleteCmd = &cobra.Command{
	Use:   "delete [template-name]",
	Short: "Delete a task template",
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			showTemplateError("Please provide a template name")
			return
		}

		name := strings.Join(args, " ")

		templateStore := templates.NewStore(getTemplatePath())
		template, err := templateStore.GetByName(name)
		if err != nil {
			showTemplateError(fmt.Sprintf("Template '%s' not found", name))
			return
		}

		if err := templateStore.Delete(template.ID); err != nil {
			showTemplateError(fmt.Sprintf("Failed to delete template: %s", err.Error()))
			return
		}

		fmt.Println(lipgloss.Style{}.
			Foreground(styles.SynthwaveGreen).
			Render(fmt.Sprintf("✅ Template '%s' deleted successfully", name)))
	},
}

func getTemplatePath() string {
	focusDir := filepath.Join(utils.GetStoragePath(), "..")
	return filepath.Join(focusDir, "templates.json")
}

func showTemplateSuccess(template models.TaskTemplate) {
	title := styles.SynthwaveTitle("📋 TEMPLATE CREATED")
	fmt.Println(title)
	fmt.Println()

	details := []struct {
		label string
		value string
		color lipgloss.Color
	}{
		{"NAME", template.Name, styles.SynthwaveCyan},
		{"DESCRIPTION", template.Description, styles.SynthwavePink},
	}

	if template.Priority != "" {
		details = append(details, struct {
			label string
			value string
			color lipgloss.Color
		}{"PRIORITY", strings.ToUpper(template.Priority), getPriorityColor(template.Priority)})
	}

	if len(template.Categories) > 0 {
		details = append(details, struct {
			label string
			value string
			color lipgloss.Color
		}{"CATEGORIES", strings.Join(template.Categories, ", "), styles.SynthwaveYellow})
	}

	for _, detail := range details {
		detailLine := lipgloss.Style{}.
			Foreground(detail.color).
			Background(styles.DarkVoid).
			Bold(true).
			Padding(0, 1).
			Render(fmt.Sprintf("▸ %s: %s", detail.label, detail.value))
		fmt.Println(detailLine)
	}
	fmt.Println()

	info := lipgloss.Style{}.
		Foreground(styles.SynthwaveCyan).
		Background(styles.DarkVoid).
		Padding(1, 2).
		Render(fmt.Sprintf("💡 Use this template with: focus template use \"%s\"", template.Name))
	fmt.Println(info)
}

func showTemplateError(message string) {
	errorBox := lipgloss.Style{}.
		Foreground(styles.SynthwaveRed).
		Background(styles.DarkVoid).
		Bold(true).
		Padding(1, 2).
		Border(lipgloss.RoundedBorder()).
		BorderForeground(styles.SynthwaveRed).
		Render(fmt.Sprintf("❌ TEMPLATE ERROR\n\n%s", message))
	fmt.Println(errorBox)
}

func init() {
	rootCmd.AddCommand(templateCmd)
	templateCmd.AddCommand(templateListCmd)
	templateCmd.AddCommand(templateCreateCmd)
	templateCmd.AddCommand(templateUseCmd)
	templateCmd.AddCommand(templateDeleteCmd)

	// Flags for create command
	templateCreateCmd.Flags().StringVar(&templateName, "name", "", "Template name (required)")
	templateCreateCmd.Flags().StringVar(&templatePriority, "priority", "medium", "Default priority")
	templateCreateCmd.Flags().StringVar(&templateCategories, "categories", "", "Default categories (comma-separated)")
}
