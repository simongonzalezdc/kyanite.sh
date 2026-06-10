package editor

import (
	"strings"
	"time"
)

// Buffer represents an in-memory text buffer
type Buffer struct {
	lines      []string
	cursorLine int
	cursorCol  int
	modified   bool
	history    *EditHistory

	// Search state
	searchTerm    string
	searchResults []SearchResult
	searchIndex   int
}

// SearchResult represents a search match
type SearchResult struct {
	Line int
	Col  int
	Text string
}

// EditCommand represents a single edit operation
type EditCommand interface {
	Undo(b *Buffer)
	Redo(b *Buffer)
	GetCursor() (line, col int)
}

// InsertRuneCommand represents inserting a single character
type InsertRuneCommand struct {
	Line      int
	Col       int
	Rune      rune
	PrevLine  int
	PrevCol   int
	Timestamp time.Time
}

func (c *InsertRuneCommand) Undo(b *Buffer) {
	if c.Line < 0 || c.Line >= len(b.lines) {
		return
	}
	line := b.lines[c.Line]
	runeWidth := len(string(c.Rune))
	if c.Col >= 0 && c.Col+runeWidth <= len(line) {
		b.lines[c.Line] = line[:c.Col] + line[c.Col+runeWidth:]
	}
	b.cursorLine = c.PrevLine
	b.cursorCol = c.PrevCol
}

func (c *InsertRuneCommand) Redo(b *Buffer) {
	if c.Line < 0 || c.Line >= len(b.lines) {
		return
	}
	line := b.lines[c.Line]
	if c.Col >= 0 && c.Col <= len(line) {
		b.lines[c.Line] = line[:c.Col] + string(c.Rune) + line[c.Col:]
	}
	b.cursorLine = c.Line
	b.cursorCol = c.Col + len(string(c.Rune))
}

func (c *InsertRuneCommand) GetCursor() (line, col int) {
	return c.Line, c.Col + len(string(c.Rune))
}

// InsertNewlineCommand represents inserting a newline
type InsertNewlineCommand struct {
	Line      int
	Col       int
	After     string // Text after cursor that moved to new line
	PrevLine  int
	PrevCol   int
	Timestamp time.Time
}

func (c *InsertNewlineCommand) Undo(b *Buffer) {
	if c.Line < 0 || c.Line >= len(b.lines) || c.Line+1 >= len(b.lines) {
		return
	}
	// Merge the split lines back together
	b.lines[c.Line] = b.lines[c.Line] + c.After
	b.lines = append(b.lines[:c.Line+1], b.lines[c.Line+2:]...)
	b.cursorLine = c.PrevLine
	b.cursorCol = c.PrevCol
}

func (c *InsertNewlineCommand) Redo(b *Buffer) {
	if c.Line < 0 || c.Line >= len(b.lines) {
		return
	}
	line := b.lines[c.Line]
	before := line[:c.Col]
	b.lines[c.Line] = before
	b.lines = append(b.lines[:c.Line+1], append([]string{c.After}, b.lines[c.Line+1:]...)...)
	b.cursorLine = c.Line + 1
	b.cursorCol = 0
}

func (c *InsertNewlineCommand) GetCursor() (line, col int) {
	return c.Line + 1, 0
}

// DeleteCharCommand represents deleting a character
type DeleteCharCommand struct {
	Line      int
	Col       int
	Deleted   rune   // What was deleted
	WasMerge  bool   // Was this a line merge?
	PrevLine  string // Previous line content if merge
	PrevLine2 int    // Previous cursor line
	PrevCol2  int    // Previous cursor col
	Timestamp time.Time
}

// ReplaceAllCommand represents a bulk find/replace operation
type ReplaceAllCommand struct {
	Replacements []struct {
		Line        int
		Col         int
		OldText     string
		NewText     string
		OrigContent string // Original line content for undo
	}
	PrevLine  int
	PrevCol   int
	Timestamp time.Time
}

func (c *DeleteCharCommand) Undo(b *Buffer) {
	if c.WasMerge {
		// Un-merge lines
		if c.Line < 0 || c.Line >= len(b.lines) {
			return
		}
		currLine := b.lines[c.Line]
		splitPos := len(c.PrevLine)
		before := currLine[:splitPos]
		after := currLine[splitPos:]
		b.lines[c.Line] = before
		b.lines = append(b.lines[:c.Line+1], append([]string{after}, b.lines[c.Line+1:]...)...)
		b.cursorLine = c.Line + 1
		b.cursorCol = 0
	} else {
		// Re-insert deleted character
		if c.Line < 0 || c.Line >= len(b.lines) {
			return
		}
		line := b.lines[c.Line]
		if c.Col >= 0 && c.Col <= len(line) {
			b.lines[c.Line] = line[:c.Col] + string(c.Deleted) + line[c.Col:]
		}
		b.cursorLine = c.PrevLine2
		b.cursorCol = c.PrevCol2
	}
}

func (c *DeleteCharCommand) Redo(b *Buffer) {
	if c.WasMerge {
		// Re-merge lines
		if c.Line < 0 || c.Line >= len(b.lines) || c.Line+1 >= len(b.lines) {
			return
		}
		prevLine := b.lines[c.Line]
		currLine := b.lines[c.Line+1]
		b.lines[c.Line] = prevLine + currLine
		b.lines = append(b.lines[:c.Line+1], b.lines[c.Line+2:]...)
		b.cursorLine = c.Line
		b.cursorCol = len(prevLine)
	} else {
		// Re-delete character
		if c.Line < 0 || c.Line >= len(b.lines) {
			return
		}
		line := b.lines[c.Line]
		if c.Col > 0 && c.Col <= len(line) {
			b.lines[c.Line] = line[:c.Col-1] + line[c.Col:]
			b.cursorLine = c.Line
			b.cursorCol = c.Col - 1
		}
	}
}

func (c *DeleteCharCommand) GetCursor() (line, col int) {
	if c.WasMerge {
		return c.Line, len(c.PrevLine)
	}
	return c.Line, c.Col - 1
}

func (c *ReplaceAllCommand) Undo(b *Buffer) {
	// Restore original line contents in reverse order
	for i := len(c.Replacements) - 1; i >= 0; i-- {
		repl := c.Replacements[i]
		if repl.Line >= 0 && repl.Line < len(b.lines) {
			b.lines[repl.Line] = repl.OrigContent
		}
	}
	b.cursorLine = c.PrevLine
	b.cursorCol = c.PrevCol
}

func (c *ReplaceAllCommand) Redo(b *Buffer) {
	// Re-apply all replacements
	for _, repl := range c.Replacements {
		if repl.Line >= 0 && repl.Line < len(b.lines) {
			line := b.lines[repl.Line]
			before := line[:repl.Col]
			after := line[repl.Col+len(repl.OldText):]
			b.lines[repl.Line] = before + repl.NewText + after
		}
	}
}

func (c *ReplaceAllCommand) GetCursor() (line, col int) {
	return c.PrevLine, c.PrevCol
}

// EditHistory manages undo/redo with delta-based commands
type EditHistory struct {
	past    []EditCommand
	future  []EditCommand
	maxSize int
}

// NewBuffer creates a new text buffer
func NewBuffer(content string) *Buffer {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 {
		lines = []string{""}
	}

	return &Buffer{
		lines:      lines,
		cursorLine: 0,
		cursorCol:  0,
		modified:   false,
		history: &EditHistory{
			past:    make([]EditCommand, 0, 100),
			future:  make([]EditCommand, 0, 100),
			maxSize: 100,
		},
	}
}

// GetContent returns the full buffer content
func (b *Buffer) GetContent() string {
	return strings.Join(b.lines, "\n")
}

// GetLines returns all lines
func (b *Buffer) GetLines() []string {
	return b.lines
}

// GetLine returns a specific line
func (b *Buffer) GetLine(lineNum int) string {
	if lineNum < 0 || lineNum >= len(b.lines) {
		return ""
	}
	return b.lines[lineNum]
}

// LineCount returns the number of lines
func (b *Buffer) LineCount() int {
	return len(b.lines)
}

// CursorPosition returns current cursor position
func (b *Buffer) CursorPosition() (line, col int) {
	return b.cursorLine, b.cursorCol
}

// InsertRune inserts a character at cursor position
func (b *Buffer) InsertRune(r rune) {
	// Validate cursor position
	if b.cursorLine < 0 || b.cursorLine >= len(b.lines) {
		return
	}

	line := b.lines[b.cursorLine]

	// Clamp cursor column to valid range
	if b.cursorCol < 0 {
		b.cursorCol = 0
	}
	if b.cursorCol > len(line) {
		b.cursorCol = len(line)
	}

	// Record command for undo
	cmd := &InsertRuneCommand{
		Line:      b.cursorLine,
		Col:       b.cursorCol,
		Rune:      r,
		PrevLine:  b.cursorLine,
		PrevCol:   b.cursorCol,
		Timestamp: time.Now(),
	}
	b.recordCommand(cmd)

	before := line[:b.cursorCol]
	after := line[b.cursorCol:]
	b.lines[b.cursorLine] = before + string(r) + after
	b.cursorCol += len(string(r))
	b.modified = true
}

// InsertNewline inserts a newline at cursor position
func (b *Buffer) InsertNewline() {
	// Validate cursor position
	if b.cursorLine < 0 || b.cursorLine >= len(b.lines) {
		return
	}

	line := b.lines[b.cursorLine]

	// Clamp cursor column to valid range
	if b.cursorCol < 0 {
		b.cursorCol = 0
	}
	if b.cursorCol > len(line) {
		b.cursorCol = len(line)
	}

	before := line[:b.cursorCol]
	after := line[b.cursorCol:]

	// Record command for undo
	cmd := &InsertNewlineCommand{
		Line:      b.cursorLine,
		Col:       b.cursorCol,
		After:     after,
		PrevLine:  b.cursorLine,
		PrevCol:   b.cursorCol,
		Timestamp: time.Now(),
	}
	b.recordCommand(cmd)

	b.lines[b.cursorLine] = before
	// Insert new line after current
	b.lines = append(b.lines[:b.cursorLine+1], append([]string{after}, b.lines[b.cursorLine+1:]...)...)

	b.cursorLine++
	b.cursorCol = 0
	b.modified = true
}

// DeleteChar deletes character before cursor (backspace)
func (b *Buffer) DeleteChar() {
	// Validate cursor position
	if b.cursorLine < 0 || b.cursorLine >= len(b.lines) {
		return
	}

	if b.cursorCol == 0 {
		// At start of line - merge with previous
		if b.cursorLine > 0 {
			prevLine := b.lines[b.cursorLine-1]
			currLine := b.lines[b.cursorLine]

			// Record command for undo
			cmd := &DeleteCharCommand{
				Line:      b.cursorLine - 1,
				Col:       b.cursorCol,
				Deleted:   '\n', // Newline deletion
				WasMerge:  true,
				PrevLine:  prevLine,
				PrevLine2: b.cursorLine,
				PrevCol2:  b.cursorCol,
				Timestamp: time.Now(),
			}
			b.recordCommand(cmd)

			b.lines[b.cursorLine-1] = prevLine + currLine
			b.lines = append(b.lines[:b.cursorLine], b.lines[b.cursorLine+1:]...)
			b.cursorLine--
			b.cursorCol = len(prevLine)
			b.modified = true
		}
	} else {
		// Ensure cursorCol is within bounds
		line := b.lines[b.cursorLine]
		if b.cursorCol > len(line) {
			b.cursorCol = len(line)
		}
		if b.cursorCol > 0 {
			// Get the character being deleted
			deleted := rune(line[b.cursorCol-1])

			// Record command for undo
			cmd := &DeleteCharCommand{
				Line:      b.cursorLine,
				Col:       b.cursorCol,
				Deleted:   deleted,
				WasMerge:  false,
				PrevLine2: b.cursorLine,
				PrevCol2:  b.cursorCol,
				Timestamp: time.Now(),
			}
			b.recordCommand(cmd)

			before := line[:b.cursorCol-1]
			after := line[b.cursorCol:]
			b.lines[b.cursorLine] = before + after
			b.cursorCol--
			b.modified = true
		}
	}
}

// MoveCursor moves the cursor
func (b *Buffer) MoveCursorUp() {
	if b.cursorLine > 0 {
		b.cursorLine--
		if b.cursorCol > len(b.lines[b.cursorLine]) {
			b.cursorCol = len(b.lines[b.cursorLine])
		}
	}
}

func (b *Buffer) MoveCursorDown() {
	if b.cursorLine < len(b.lines)-1 {
		b.cursorLine++
		if b.cursorCol > len(b.lines[b.cursorLine]) {
			b.cursorCol = len(b.lines[b.cursorLine])
		}
	}
}

func (b *Buffer) MoveCursorLeft() {
	if b.cursorCol > 0 {
		b.cursorCol--
	} else if b.cursorLine > 0 {
		// Move to end of previous line
		b.cursorLine--
		b.cursorCol = len(b.lines[b.cursorLine])
	}
}

func (b *Buffer) MoveCursorRight() {
	if b.cursorCol < len(b.lines[b.cursorLine]) {
		b.cursorCol++
	} else if b.cursorLine < len(b.lines)-1 {
		// Move to start of next line
		b.cursorLine++
		b.cursorCol = 0
	}
}

// MoveCursorTo moves cursor to specific position
func (b *Buffer) MoveCursorTo(line, col int) {
	if line >= 0 && line < len(b.lines) {
		b.cursorLine = line
		if col >= 0 && col <= len(b.lines[line]) {
			b.cursorCol = col
		}
	}
}

// IsModified returns whether buffer has been modified
func (b *Buffer) IsModified() bool {
	return b.modified
}

// SetModified sets the modified flag
func (b *Buffer) SetModified(modified bool) {
	b.modified = modified
}

// recordCommand records an edit command for undo/redo
func (b *Buffer) recordCommand(cmd EditCommand) {
	// Only keep last maxSize commands
	if len(b.history.past) >= b.history.maxSize {
		b.history.past = b.history.past[1:]
	}

	b.history.past = append(b.history.past, cmd)
	b.history.future = nil // Clear redo stack
}

// Undo reverts the last edit operation
func (b *Buffer) Undo() bool {
	if len(b.history.past) == 0 {
		return false
	}

	// Get the last command
	cmd := b.history.past[len(b.history.past)-1]
	b.history.past = b.history.past[:len(b.history.past)-1]

	// Undo the command
	cmd.Undo(b)

	// Move command to future for redo
	b.history.future = append(b.history.future, cmd)
	b.modified = true

	return true
}

// Redo re-applies a previously undone edit operation
func (b *Buffer) Redo() bool {
	if len(b.history.future) == 0 {
		return false
	}

	// Get the last undone command
	cmd := b.history.future[len(b.history.future)-1]
	b.history.future = b.history.future[:len(b.history.future)-1]

	// Redo the command
	cmd.Redo(b)

	// Move command back to past
	b.history.past = append(b.history.past, cmd)
	b.modified = true

	return true
}

// Find searches for the given term and stores results
func (b *Buffer) Find(term string, caseSensitive bool) int {
	b.searchTerm = term
	b.searchResults = nil
	b.searchIndex = -1

	if term == "" {
		return 0
	}

	searchTerm := term
	if !caseSensitive {
		searchTerm = strings.ToLower(term)
	}

	for lineNum, line := range b.lines {
		searchLine := line
		if !caseSensitive {
			searchLine = strings.ToLower(line)
		}

		col := 0
		for {
			idx := strings.Index(searchLine[col:], searchTerm)
			if idx == -1 {
				break
			}

			actualCol := col + idx
			b.searchResults = append(b.searchResults, SearchResult{
				Line: lineNum,
				Col:  actualCol,
				Text: line[actualCol : actualCol+len(term)],
			})
			col = actualCol + 1
		}
	}

	if len(b.searchResults) > 0 {
		b.searchIndex = 0
		// Jump to first result
		b.cursorLine = b.searchResults[0].Line
		b.cursorCol = b.searchResults[0].Col
	}

	return len(b.searchResults)
}

// FindNext jumps to the next search result
func (b *Buffer) FindNext() bool {
	if len(b.searchResults) == 0 {
		return false
	}

	b.searchIndex = (b.searchIndex + 1) % len(b.searchResults)
	result := b.searchResults[b.searchIndex]
	b.cursorLine = result.Line
	b.cursorCol = result.Col
	return true
}

// FindPrevious jumps to the previous search result
func (b *Buffer) FindPrevious() bool {
	if len(b.searchResults) == 0 {
		return false
	}

	b.searchIndex--
	if b.searchIndex < 0 {
		b.searchIndex = len(b.searchResults) - 1
	}

	result := b.searchResults[b.searchIndex]
	b.cursorLine = result.Line
	b.cursorCol = result.Col
	return true
}

// ClearSearch clears the search results
func (b *Buffer) ClearSearch() {
	b.searchTerm = ""
	b.searchResults = nil
	b.searchIndex = -1
}

// GetSearchInfo returns current search information
func (b *Buffer) GetSearchInfo() (term string, current int, total int) {
	if len(b.searchResults) == 0 {
		return b.searchTerm, 0, 0
	}
	return b.searchTerm, b.searchIndex + 1, len(b.searchResults)
}

// ReplaceCurrent replaces the current search match and moves to the next one
func (b *Buffer) ReplaceCurrent(replacement string) bool {
	if len(b.searchResults) == 0 || b.searchIndex < 0 || b.searchIndex >= len(b.searchResults) {
		return false
	}

	result := b.searchResults[b.searchIndex]
	line := b.lines[result.Line]

	// Create replace command for undo
	cmd := &ReplaceAllCommand{
		Replacements: make([]struct {
			Line        int
			Col         int
			OldText     string
			NewText     string
			OrigContent string
		}, 1),
		PrevLine:  b.cursorLine,
		PrevCol:   b.cursorCol,
		Timestamp: time.Now(),
	}

	cmd.Replacements[0].Line = result.Line
	cmd.Replacements[0].Col = result.Col
	cmd.Replacements[0].OldText = b.searchTerm
	cmd.Replacements[0].NewText = replacement
	cmd.Replacements[0].OrigContent = line

	// Record command
	b.recordCommand(cmd)

	// Perform replacement
	before := line[:result.Col]
	after := line[result.Col+len(b.searchTerm):]
	b.lines[result.Line] = before + replacement + after

	b.modified = true

	// Update cursor position to end of replacement
	b.cursorLine = result.Line
	b.cursorCol = result.Col + len(replacement)

	// Remove this result from search results
	b.searchResults = append(b.searchResults[:b.searchIndex], b.searchResults[b.searchIndex+1:]...)

	// Adjust search index
	if b.searchIndex >= len(b.searchResults) {
		b.searchIndex = 0
	}

	// If no more results, clear search
	if len(b.searchResults) == 0 {
		b.ClearSearch()
		return false
	}

	// Move to next occurrence
	if len(b.searchResults) > 0 && b.searchIndex < len(b.searchResults) {
		next := b.searchResults[b.searchIndex]
		b.cursorLine = next.Line
		b.cursorCol = next.Col
	}

	return true
}

// ReplaceAll replaces all occurrences of the search term
func (b *Buffer) ReplaceAll(searchTerm, replacement string, caseSensitive bool) int {
	count := b.Find(searchTerm, caseSensitive)
	if count == 0 {
		return 0
	}

	// Build the replace command
	cmd := &ReplaceAllCommand{
		Replacements: make([]struct {
			Line        int
			Col         int
			OldText     string
			NewText     string
			OrigContent string
		}, len(b.searchResults)),
		PrevLine:  b.cursorLine,
		PrevCol:   b.cursorCol,
		Timestamp: time.Now(),
	}

	// Store replacement info
	for i, result := range b.searchResults {
		cmd.Replacements[i].Line = result.Line
		cmd.Replacements[i].Col = result.Col
		cmd.Replacements[i].OldText = searchTerm
		cmd.Replacements[i].NewText = replacement
		cmd.Replacements[i].OrigContent = b.lines[result.Line]
	}

	// Record the command
	b.recordCommand(cmd)

	// Replace in reverse order to avoid index shifting
	for i := len(b.searchResults) - 1; i >= 0; i-- {
		result := b.searchResults[i]
		line := b.lines[result.Line]
		before := line[:result.Col]
		after := line[result.Col+len(searchTerm):]
		b.lines[result.Line] = before + replacement + after
	}

	b.modified = true
	b.ClearSearch()
	return count
}
