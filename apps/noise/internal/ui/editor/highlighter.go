package editor

import (
	"regexp"
	"sort"

	"github.com/kyanite/noise/internal/theme"
	"github.com/charmbracelet/lipgloss"
)

// MarkdownElement represents a highlighted markdown element.
type MarkdownElement struct {
	Type    ElementType
	Content string
	Start   int
	End     int
	Style   lipgloss.Style
}

// ElementType represents different types of markdown elements.
type ElementType int

const (
	ElementText ElementType = iota
	ElementHeader
	ElementBold
	ElementItalic
	ElementCode
	ElementCodeBlock
	ElementLink
	ElementList
	ElementBlockquote
)

// SyntaxHighlighter handles markdown parsing and styling.
type SyntaxHighlighter struct {
	patterns map[ElementType]*regexp.Regexp
	styles   map[ElementType]lipgloss.Style
}

// NewSyntaxHighlighter creates a new syntax highlighter with predefined styles.
func NewSyntaxHighlighter() *SyntaxHighlighter {
	t := theme.GetManager().Current()
	sh := &SyntaxHighlighter{
		patterns: make(map[ElementType]*regexp.Regexp),
		styles:   make(map[ElementType]lipgloss.Style),
	}

	// Define regex patterns for markdown elements.
	sh.patterns[ElementHeader] = regexp.MustCompile(`^(#{1,6})\s+(.+)$`)
	sh.patterns[ElementBold] = regexp.MustCompile(`\*\*(.+?)\*\*`)
	sh.patterns[ElementItalic] = regexp.MustCompile(`\*([^*]+?)\*`)
	sh.patterns[ElementCode] = regexp.MustCompile("`([^`]+?)`")
	sh.patterns[ElementCodeBlock] = regexp.MustCompile("(?s)```[\\s\\S]*?```")
	sh.patterns[ElementLink] = regexp.MustCompile(`\[([^\]]+)\]\(([^)]+)\)`)
	sh.patterns[ElementList] = regexp.MustCompile(`(?m)^(\s*)([-*+]|\d+\.)\s+(.+)$`)
	sh.patterns[ElementBlockquote] = regexp.MustCompile(`(?m)^>\s+(.+)$`)

	// Define styles for each element type using current theme.
	sh.styles[ElementHeader] = lipgloss.Style{}.Bold(true).Foreground(t.Accent)
	sh.styles[ElementBold] = lipgloss.Style{}.Bold(true).Foreground(t.Primary)
	sh.styles[ElementItalic] = lipgloss.Style{}.Italic(true).Foreground(t.Secondary)
	sh.styles[ElementCode] = lipgloss.Style{}.Background(t.Background).Foreground(t.Success).Padding(0, 1)
	sh.styles[ElementCodeBlock] = lipgloss.Style{}.Background(t.Background).Foreground(t.Text)
	sh.styles[ElementLink] = lipgloss.Style{}.Foreground(t.Accent).Underline(true)
	sh.styles[ElementList] = lipgloss.Style{}.Foreground(t.Secondary)
	sh.styles[ElementBlockquote] = lipgloss.Style{}.Foreground(t.Secondary).Italic(true)

	return sh
}

// ParseMarkdown parses markdown content and returns highlighted elements.
func (sh *SyntaxHighlighter) ParseMarkdown(content string) []MarkdownElement {
	if content == "" {
		return nil
	}

	var elements []MarkdownElement
	processed := make([]bool, len(content))

	// Process code blocks first (they take precedence).
	if pattern := sh.patterns[ElementCodeBlock]; pattern != nil {
		codeBlockMatches := pattern.FindAllStringIndex(content, -1)
		for _, match := range codeBlockMatches {
			start, end := match[0], match[1]
			if start >= 0 && end <= len(content) && !processed[start] && !processed[end-1] {
				elements = append(elements, MarkdownElement{
					Type:    ElementCodeBlock,
					Content: content[start:end],
					Start:   start,
					End:     end,
					Style:   sh.styles[ElementCodeBlock],
				})
				for i := start; i < end; i++ {
					processed[i] = true
				}
			}
		}
	}

	// Process other elements.
	for elementType, pattern := range sh.patterns {
		if elementType == ElementCodeBlock || pattern == nil {
			continue
		}

		matches := pattern.FindAllStringSubmatchIndex(content, -1)
		for _, match := range matches {
			if len(match) < 4 {
				continue
			}

			start, end := match[0], match[1]
			if start < 0 || end > len(content) || processed[start] || processed[end-1] {
				continue
			}

			subStart, subEnd := match[2], match[3]
			if subStart < 0 || subEnd > len(content) {
				subStart, subEnd = start, end
			}

			element := MarkdownElement{
				Type:    elementType,
				Content: content[subStart:subEnd],
				Start:   start,
				End:     end,
				Style:   sh.styles[elementType],
			}

			elements = append(elements, element)
			for i := start; i < end; i++ {
				processed[i] = true
			}
		}
	}

	// Sort elements by start position to ensure deterministic output.
	sort.Slice(elements, func(i, j int) bool {
		return elements[i].Start < elements[j].Start
	})

	return elements
}
