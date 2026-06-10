package export

import (
	"fmt"
	"strings"

	"github.com/jung-kurt/gofpdf"
	"github.com/kyanite/syntax/internal/scene"
	"github.com/kyanite/syntax/internal/story"
)

// ExportPDF exports the project to a PDF file
func ExportPDF(project *story.Project, scenes map[string]*scene.Scene) ([]byte, error) {
	// Create new PDF with letter size
	pdf := gofpdf.New("P", "mm", "Letter", "")
	pdf.SetAutoPageBreak(true, 15)

	// Add a page
	pdf.AddPage()

	// Title page
	pdf.SetFont("Arial", "B", 24)
	pdf.CellFormat(0, 20, "", "", 0, "C", false, 0, "")
	pdf.Ln(40)
	pdf.CellFormat(0, 15, project.Title, "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "I", 16)
	pdf.CellFormat(0, 10, "by "+project.Author, "", 1, "C", false, 0, "")

	pdf.SetFont("Arial", "", 12)
	pdf.Ln(10)
	if project.Genre != "" {
		pdf.CellFormat(0, 8, "Genre: "+project.Genre, "", 1, "C", false, 0, "")
	}

	// Sort scenes by chapter and scene number
	sortedScenes := scene.SortScenes(scenes)

	// Add scenes
	currentChapter := 0
	for _, sc := range sortedScenes {
		// Start new chapter
		if sc.Chapter != currentChapter {
			pdf.AddPage()
			currentChapter = sc.Chapter

			pdf.SetFont("Arial", "B", 18)
			pdf.CellFormat(0, 15, fmt.Sprintf("Chapter %d", currentChapter), "", 1, "L", false, 0, "")
			pdf.Ln(5)
		}

		// Scene heading
		if sc.Name != "" {
			pdf.SetFont("Arial", "B", 14)
			pdf.CellFormat(0, 10, sc.Name, "", 1, "L", false, 0, "")
			pdf.Ln(3)
		}

		// Scene content
		pdf.SetFont("Arial", "", 12)

		// Split content into paragraphs
		paragraphs := strings.Split(sc.Content, "\n\n")
		for _, para := range paragraphs {
			if strings.TrimSpace(para) == "" {
				continue
			}

			// Write paragraph with proper line wrapping
			pdf.MultiCell(0, 6, strings.TrimSpace(para), "", "L", false)
			pdf.Ln(4)
		}

		pdf.Ln(6)
	}

	// Add final page with statistics
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.CellFormat(0, 12, "Story Statistics", "", 1, "L", false, 0, "")
	pdf.Ln(5)

	pdf.SetFont("Arial", "", 12)

	totalWords := 0
	for _, sc := range sortedScenes {
		totalWords += sc.WordCount
	}

	stats := []string{
		fmt.Sprintf("Total Scenes: %d", len(sortedScenes)),
		fmt.Sprintf("Total Words: %d", totalWords),
		fmt.Sprintf("Total Characters: %d", project.TotalCharacters),
		fmt.Sprintf("Total Locations: %d", project.TotalLocations),
	}

	for _, stat := range stats {
		pdf.CellFormat(0, 8, stat, "", 1, "L", false, 0, "")
	}

	// Output to bytes
	var buf strings.Builder
	err := pdf.Output(&stringWriterWrapper{&buf})
	if err != nil {
		return nil, fmt.Errorf("failed to generate PDF: %w", err)
	}

	return []byte(buf.String()), nil
}

// stringWriterWrapper wraps strings.Builder to implement io.Writer
type stringWriterWrapper struct {
	*strings.Builder
}

func (w *stringWriterWrapper) Write(p []byte) (n int, err error) {
	return w.Builder.Write(p)
}
