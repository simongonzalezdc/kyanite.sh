package export

import (
	"strings"
	"testing"

	"github.com/kyanite/syntax/internal/scene"
	"github.com/kyanite/syntax/internal/story"
)

func TestExportMarkdown(t *testing.T) {
	project := &story.Project{
		Title:  "Test Novel",
		Author: "Test Author",
		Genre:  "Fiction",
	}

	scenes := map[string]*scene.Scene{
		"scene1": {
			ID:          "scene1",
			Name:        "Opening Scene",
			Chapter:     1,
			SceneNumber: 1,
			Content:     "It was a dark and stormy night.",
			WordCount:   7,
		},
		"scene2": {
			ID:          "scene2",
			Name:        "The Plot Thickens",
			Chapter:     1,
			SceneNumber: 2,
			Content:     "The hero discovered a clue.",
			WordCount:   5,
		},
		"scene3": {
			ID:          "scene3",
			Name:        "Chapter Two Begins",
			Chapter:     2,
			SceneNumber: 1,
			Content:     "A new day dawned.",
			WordCount:   4,
		},
	}

	data, err := ExportMarkdown(project, scenes)
	if err != nil {
		t.Fatalf("ExportMarkdown failed: %v", err)
	}

	output := string(data)

	// Check title
	if !strings.Contains(output, "# Test Novel") {
		t.Error("Expected output to contain title")
	}

	// Check author
	if !strings.Contains(output, "by Test Author") {
		t.Error("Expected output to contain author")
	}

	// Check genre
	if !strings.Contains(output, "Fiction") {
		t.Error("Expected output to contain genre")
	}

	// Check chapters
	if !strings.Contains(output, "# Chapter 1") {
		t.Error("Expected output to contain Chapter 1")
	}

	if !strings.Contains(output, "# Chapter 2") {
		t.Error("Expected output to contain Chapter 2")
	}

	// Check scene names
	if !strings.Contains(output, "## Opening Scene") {
		t.Error("Expected output to contain scene name 'Opening Scene'")
	}

	if !strings.Contains(output, "## The Plot Thickens") {
		t.Error("Expected output to contain scene name 'The Plot Thickens'")
	}

	// Check scene content
	if !strings.Contains(output, "It was a dark and stormy night.") {
		t.Error("Expected output to contain scene content")
	}

	// Check statistics
	if !strings.Contains(output, "## Story Statistics") {
		t.Error("Expected output to contain statistics section")
	}

	if !strings.Contains(output, "Total Scenes: 3") {
		t.Error("Expected output to contain correct scene count")
	}

	if !strings.Contains(output, "Total Words: 16") {
		t.Error("Expected output to contain correct word count")
	}
}

func TestExportHTML(t *testing.T) {
	project := &story.Project{
		Title:  "Test Novel",
		Author: "Test Author",
		Genre:  "Sci-Fi",
	}

	scenes := map[string]*scene.Scene{
		"scene1": {
			ID:          "scene1",
			Name:        "First Scene",
			Chapter:     1,
			SceneNumber: 1,
			Content:     "In the beginning...",
			WordCount:   3,
		},
	}

	data, err := ExportHTML(project, scenes)
	if err != nil {
		t.Fatalf("ExportHTML failed: %v", err)
	}

	output := string(data)

	// Check HTML structure
	if !strings.Contains(output, "<!DOCTYPE html>") {
		t.Error("Expected output to contain DOCTYPE")
	}

	if !strings.Contains(output, "<html") {
		t.Error("Expected output to contain html tag")
	}

	if !strings.Contains(output, "</html>") {
		t.Error("Expected output to contain closing html tag")
	}

	// Check title in head
	if !strings.Contains(output, "<title>Test Novel</title>") {
		t.Error("Expected output to contain title in head")
	}

	// Check title in body
	if !strings.Contains(output, "<h1>Test Novel</h1>") {
		t.Error("Expected output to contain h1 title")
	}

	// Check author
	if !strings.Contains(output, "by Test Author") {
		t.Error("Expected output to contain author")
	}

	// Check genre
	if !strings.Contains(output, "Sci-Fi") {
		t.Error("Expected output to contain genre")
	}

	// Check chapter
	if !strings.Contains(output, "<h2>Chapter 1</h2>") {
		t.Error("Expected output to contain chapter heading")
	}

	// Check scene name
	if !strings.Contains(output, "<h3>First Scene</h3>") {
		t.Error("Expected output to contain scene name")
	}

	// Check CSS
	if !strings.Contains(output, "<style>") {
		t.Error("Expected output to contain CSS")
	}
}

func TestExportMarkdownEmptyScenes(t *testing.T) {
	project := &story.Project{
		Title:  "Empty Novel",
		Author: "Author",
		Genre:  "Mystery",
	}

	scenes := map[string]*scene.Scene{}

	data, err := ExportMarkdown(project, scenes)
	if err != nil {
		t.Fatalf("ExportMarkdown failed with empty scenes: %v", err)
	}

	output := string(data)

	// Should still have title and stats
	if !strings.Contains(output, "# Empty Novel") {
		t.Error("Expected output to contain title even with no scenes")
	}

	if !strings.Contains(output, "Total Scenes: 0") {
		t.Error("Expected output to show 0 scenes")
	}

	if !strings.Contains(output, "Total Words: 0") {
		t.Error("Expected output to show 0 words")
	}
}

func TestSceneSorting(t *testing.T) {
	project := &story.Project{
		Title: "Sorted Novel",
	}

	// Create scenes out of order
	scenes := map[string]*scene.Scene{
		"scene3": {
			ID:          "scene3",
			Chapter:     2,
			SceneNumber: 1,
			Content:     "Third in order",
		},
		"scene1": {
			ID:          "scene1",
			Chapter:     1,
			SceneNumber: 1,
			Content:     "First in order",
		},
		"scene2": {
			ID:          "scene2",
			Chapter:     1,
			SceneNumber: 2,
			Content:     "Second in order",
		},
	}

	data, err := ExportMarkdown(project, scenes)
	if err != nil {
		t.Fatalf("ExportMarkdown failed: %v", err)
	}

	output := string(data)

	// Find positions of content in output
	firstPos := strings.Index(output, "First in order")
	secondPos := strings.Index(output, "Second in order")
	thirdPos := strings.Index(output, "Third in order")

	if firstPos == -1 || secondPos == -1 || thirdPos == -1 {
		t.Fatal("Not all scene content found in output")
	}

	// Verify correct order
	if !(firstPos < secondPos && secondPos < thirdPos) {
		t.Error("Scenes are not in correct order")
	}
}
