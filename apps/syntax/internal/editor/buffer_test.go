package editor

import (
	"strings"
	"testing"
)

func TestNewBuffer(t *testing.T) {
	content := "Hello\nWorld"
	buf := NewBuffer(content)

	if buf.GetContent() != content {
		t.Errorf("Expected content '%s', got '%s'", content, buf.GetContent())
	}

	lines := buf.GetLines()
	if len(lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(lines))
	}

	if lines[0] != "Hello" {
		t.Errorf("Expected first line 'Hello', got '%s'", lines[0])
	}

	if lines[1] != "World" {
		t.Errorf("Expected second line 'World', got '%s'", lines[1])
	}
}

func TestInsertRune(t *testing.T) {
	buf := NewBuffer("")

	buf.InsertRune('H')
	buf.InsertRune('i')

	content := buf.GetContent()
	if content != "Hi" {
		t.Errorf("Expected 'Hi', got '%s'", content)
	}
}

func TestInsertNewline(t *testing.T) {
	buf := NewBuffer("Hello")

	// Move cursor to end
	buf.cursorCol = 5
	buf.InsertNewline()
	buf.InsertRune('W')
	buf.InsertRune('o')
	buf.InsertRune('r')
	buf.InsertRune('l')
	buf.InsertRune('d')

	content := buf.GetContent()
	if content != "Hello\nWorld" {
		t.Errorf("Expected 'Hello\\nWorld', got '%s'", content)
	}

	lines := buf.GetLines()
	if len(lines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(lines))
	}
}

func TestDeleteChar(t *testing.T) {
	buf := NewBuffer("Hello")

	// Move cursor to end
	buf.cursorCol = 5

	// Delete 'o'
	buf.DeleteChar()

	content := buf.GetContent()
	if content != "Hell" {
		t.Errorf("Expected 'Hell', got '%s'", content)
	}
}

func TestCursorMovement(t *testing.T) {
	buf := NewBuffer("Hello\nWorld")

	// Start at 0,0
	line, col := buf.CursorPosition()
	if line != 0 || col != 0 {
		t.Errorf("Expected cursor at (0,0), got (%d,%d)", line, col)
	}

	// Move right
	buf.MoveCursorRight()
	line, col = buf.CursorPosition()
	if line != 0 || col != 1 {
		t.Errorf("Expected cursor at (0,1), got (%d,%d)", line, col)
	}

	// Move down
	buf.MoveCursorDown()
	line, col = buf.CursorPosition()
	if line != 1 {
		t.Errorf("Expected cursor at line 1, got %d", line)
	}

	// Move left
	buf.MoveCursorLeft()
	line, col = buf.CursorPosition()
	if col != 0 {
		t.Errorf("Expected cursor at col 0, got %d", col)
	}

	// Move up
	buf.MoveCursorUp()
	line, col = buf.CursorPosition()
	if line != 0 {
		t.Errorf("Expected cursor at line 0, got %d", line)
	}
}

func TestUndo(t *testing.T) {
	buf := NewBuffer("")

	// Type "Hello"
	buf.InsertRune('H')
	buf.InsertRune('e')
	buf.InsertRune('l')
	buf.InsertRune('l')
	buf.InsertRune('o')

	content := buf.GetContent()
	if content != "Hello" {
		t.Errorf("Expected 'Hello', got '%s'", content)
	}

	// Undo once
	if !buf.Undo() {
		t.Error("Expected undo to succeed")
	}

	content = buf.GetContent()
	if content != "Hell" {
		t.Errorf("After undo, expected 'Hell', got '%s'", content)
	}
}

func TestUndoMultiByteRune(t *testing.T) {
	buf := NewBuffer("")

	buf.InsertRune('e')
	buf.InsertRune('é')

	if content := buf.GetContent(); content != "eé" {
		t.Fatalf("Expected 'eé', got '%s'", content)
	}

	if !buf.Undo() {
		t.Fatal("Expected undo to succeed")
	}

	if content := buf.GetContent(); content != "e" {
		t.Fatalf("After undo, expected 'e', got '%s'", content)
	}
}

func TestRedo(t *testing.T) {
	buf := NewBuffer("")

	// Type "Hi"
	buf.InsertRune('H')
	buf.InsertRune('i')

	// Undo
	buf.Undo()

	content := buf.GetContent()
	if content != "H" {
		t.Errorf("After undo, expected 'H', got '%s'", content)
	}

	// Redo
	if !buf.Redo() {
		t.Error("Expected redo to succeed")
	}

	content = buf.GetContent()
	if content != "Hi" {
		t.Errorf("After redo, expected 'Hi', got '%s'", content)
	}
}

func TestUndoRedoMultiple(t *testing.T) {
	buf := NewBuffer("")

	// Type "ABC"
	buf.InsertRune('A')
	buf.InsertRune('B')
	buf.InsertRune('C')

	// Undo all
	buf.Undo() // -> "AB"
	buf.Undo() // -> "A"
	buf.Undo() // -> ""

	content := buf.GetContent()
	if content != "" {
		t.Errorf("After undoing all, expected empty string, got '%s'", content)
	}

	// Redo all
	buf.Redo() // -> "A"
	buf.Redo() // -> "AB"
	buf.Redo() // -> "ABC"

	content = buf.GetContent()
	if content != "ABC" {
		t.Errorf("After redoing all, expected 'ABC', got '%s'", content)
	}
}

func TestModifiedFlag(t *testing.T) {
	buf := NewBuffer("Hello")

	if buf.IsModified() {
		t.Error("New buffer should not be modified")
	}

	buf.InsertRune('!')

	if !buf.IsModified() {
		t.Error("Buffer should be modified after insert")
	}

	buf.SetModified(false)

	if buf.IsModified() {
		t.Error("Buffer should not be modified after SetModified(false)")
	}
}

func TestLargeContent(t *testing.T) {
	// Test with large content
	lines := make([]string, 1000)
	for i := range lines {
		lines[i] = "This is line number"
	}
	content := strings.Join(lines, "\n")

	buf := NewBuffer(content)

	bufLines := buf.GetLines()
	if len(bufLines) != 1000 {
		t.Errorf("Expected 1000 lines, got %d", len(bufLines))
	}

	// Test cursor movement in large buffer
	for i := 0; i < 500; i++ {
		buf.MoveCursorDown()
	}

	line, _ := buf.CursorPosition()
	if line != 500 {
		t.Errorf("Expected cursor at line 500, got %d", line)
	}
}
