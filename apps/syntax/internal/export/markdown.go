package export

import (
	"fmt"
	"strings"
	"time"

	"github.com/kyanite/syntax/internal/scene"
	"github.com/kyanite/syntax/internal/story"
)

// ExportMarkdown exports the project to a single Markdown file
func ExportMarkdown(project *story.Project, scenes map[string]*scene.Scene) ([]byte, error) {
	var b strings.Builder

	// Title page
	b.WriteString(fmt.Sprintf("# %s\n\n", project.Title))
	if project.Author != "" {
		b.WriteString(fmt.Sprintf("**by %s**\n\n", project.Author))
	}
	if project.Genre != "" {
		b.WriteString(fmt.Sprintf("*%s*\n\n", project.Genre))
	}
	b.WriteString("---\n\n")

	// Sort scenes by chapter and scene number
	sceneList := scene.SortScenes(scenes)

	// Group by chapter
	currentChapter := 0
	for _, sc := range sceneList {
		if sc.Chapter != currentChapter {
			currentChapter = sc.Chapter
			b.WriteString(fmt.Sprintf("# Chapter %d\n\n", currentChapter))
		}

		// Scene heading (optional)
		if sc.Name != "" {
			b.WriteString(fmt.Sprintf("## %s\n\n", sc.Name))
		}

		// Scene content
		b.WriteString(sc.Content)
		b.WriteString("\n\n")
		b.WriteString("---\n\n")
	}

	// Stats footer
	b.WriteString("\n\n")
	b.WriteString("## Story Statistics\n\n")
	b.WriteString(fmt.Sprintf("- Total Scenes: %d\n", len(scenes)))
	totalWords := 0
	for _, sc := range scenes {
		totalWords += sc.WordCount
	}
	b.WriteString(fmt.Sprintf("- Total Words: %d\n", totalWords))
	b.WriteString(fmt.Sprintf("- Exported: %s\n", time.Now().Format("2006-01-02 15:04:05")))

	return []byte(b.String()), nil
}

// ExportHTML exports the project to HTML
func ExportHTML(project *story.Project, scenes map[string]*scene.Scene) ([]byte, error) {
	var b strings.Builder

	// HTML header
	b.WriteString("<!DOCTYPE html>\n")
	b.WriteString("<html lang=\"en\">\n")
	b.WriteString("<head>\n")
	b.WriteString("  <meta charset=\"UTF-8\">\n")
	b.WriteString("  <meta name=\"viewport\" content=\"width=device-width, initial-scale=1.0\">\n")
	b.WriteString(fmt.Sprintf("  <title>%s</title>\n", project.Title))
	b.WriteString("  <style>\n")
	b.WriteString(`
    body {
      font-family: Georgia, serif;
      line-height: 1.6;
      max-width: 800px;
      margin: 0 auto;
      padding: 20px;
      background: #fefefe;
      color: #333;
    }
    h1 { font-size: 2.5em; margin-bottom: 0.5em; }
    h2 { font-size: 2em; margin-top: 2em; }
    h3 { font-size: 1.5em; margin-top: 1.5em; }
    .author { font-style: italic; color: #666; }
    .genre { color: #888; }
    hr { border: none; border-top: 1px solid #ddd; margin: 2em 0; }
    .chapter-break { page-break-before: always; }
  `)
	b.WriteString("  </style>\n")
	b.WriteString("</head>\n")
	b.WriteString("<body>\n")

	// Title page
	b.WriteString(fmt.Sprintf("  <h1>%s</h1>\n", project.Title))
	if project.Author != "" {
		b.WriteString(fmt.Sprintf("  <p class=\"author\">by %s</p>\n", project.Author))
	}
	if project.Genre != "" {
		b.WriteString(fmt.Sprintf("  <p class=\"genre\">%s</p>\n", project.Genre))
	}
	b.WriteString("  <hr>\n\n")

	// Sort and output scenes
	sceneList := scene.SortScenes(scenes)

	currentChapter := 0
	for _, sc := range sceneList {
		if sc.Chapter != currentChapter {
			currentChapter = sc.Chapter
			if currentChapter > 1 {
				b.WriteString("  <div class=\"chapter-break\"></div>\n")
			}
			b.WriteString(fmt.Sprintf("  <h2>Chapter %d</h2>\n", currentChapter))
		}

		if sc.Name != "" {
			b.WriteString(fmt.Sprintf("  <h3>%s</h3>\n", sc.Name))
		}

		// Convert markdown to HTML (simple version)
		content := simpleMarkdownToHTML(sc.Content)
		b.WriteString(content)
		b.WriteString("  <hr>\n\n")
	}

	b.WriteString("</body>\n")
	b.WriteString("</html>\n")

	return []byte(b.String()), nil
}

// simpleMarkdownToHTML converts basic markdown to HTML
func simpleMarkdownToHTML(md string) string {
	lines := strings.Split(md, "\n")
	var result strings.Builder

	for _, line := range lines {
		if strings.HasPrefix(line, "# ") {
			result.WriteString(fmt.Sprintf("  <h1>%s</h1>\n", strings.TrimPrefix(line, "# ")))
		} else if strings.HasPrefix(line, "## ") {
			result.WriteString(fmt.Sprintf("  <h2>%s</h2>\n", strings.TrimPrefix(line, "## ")))
		} else if strings.HasPrefix(line, "### ") {
			result.WriteString(fmt.Sprintf("  <h3>%s</h3>\n", strings.TrimPrefix(line, "### ")))
		} else if strings.TrimSpace(line) == "" {
			result.WriteString("  <br>\n")
		} else {
			result.WriteString(fmt.Sprintf("  <p>%s</p>\n", line))
		}
	}

	return result.String()
}
