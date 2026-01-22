package editor

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestFileDialogCreation tests file dialog creation and initialization
func TestFileDialogCreation(t *testing.T) {
	t.Run("OpenFileDialog", func(t *testing.T) {
		dialog := NewFileDialogModel(DialogOpen, "Open File", "./test", []string{".txt", ".md"})

		assert.Equal(t, DialogOpen, dialog.dialogType)
		assert.Equal(t, "Open File", dialog.title)
		assert.Equal(t, "./test", dialog.defaultPath)
		assert.Equal(t, []string{".txt", ".md"}, dialog.allowedExts)
		assert.False(t, dialog.visible)
		assert.False(t, dialog.showHidden)
	})

	t.Run("SaveFileDialog", func(t *testing.T) {
		dialog := NewFileDialogModel(DialogSave, "Save File", "./test", []string{".json"})

		assert.Equal(t, DialogSave, dialog.dialogType)
		assert.Equal(t, "Save File", dialog.title)
		assert.Equal(t, "./test", dialog.defaultPath)
		assert.Equal(t, []string{".json"}, dialog.allowedExts)
	})

	t.Run("SaveAsDialog", func(t *testing.T) {
		dialog := NewFileDialogModel(DialogSaveAs, "Save As", "./test", []string{".txt"})

		assert.Equal(t, DialogSaveAs, dialog.dialogType)
		assert.Equal(t, "Save As", dialog.title)
	})
}

// TestFileDialogVisibility tests file dialog visibility controls
func TestFileDialogVisibility(t *testing.T) {
	dialog := NewFileDialogModel(DialogOpen, "Test", "./test", []string{".txt"})

	assert.False(t, dialog.IsVisible(), "Dialog should not be visible initially")

	dialog.Show()
	assert.True(t, dialog.IsVisible(), "Dialog should be visible after Show")

	dialog.Hide()
	assert.False(t, dialog.IsVisible(), "Dialog should not be visible after Hide")
}

// TestFileDialogDimensions tests file dialog dimension handling
func TestFileDialogDimensions(t *testing.T) {
	dialog := NewFileDialogModel(DialogOpen, "Test", "./test", []string{".txt"})

	dialog.SetDimensions(80, 24)
	assert.Equal(t, 80, dialog.width)
	assert.Equal(t, 24, dialog.height)

	// Check that list dimensions are set correctly
	assert.Greater(t, dialog.list.Width(), 0, "List width should be set")
	assert.Greater(t, dialog.list.Height(), 0, "List height should be set")
}

// TestFileDialogCallbacks tests file dialog callback functionality
func TestFileDialogCallbacks(t *testing.T) {
	dialog := NewFileDialogModel(DialogOpen, "Test", "./test", []string{".txt"})

	cancelCalled := false

	dialog.SetCancelCallback(func() {
		cancelCalled = true
	})

	// Test cancel callback
	dialog.Show()
	dialog.Update(tea.KeyMsg{Type: tea.KeyEsc})
	assert.True(t, cancelCalled, "Cancel callback should be called after Esc key")
	assert.False(t, dialog.IsVisible(), "Dialog should be hidden after Esc key")
}

// TestFileDialogValidation tests file dialog validation
func TestFileDialogValidation(t *testing.T) {
	dialog := NewFileDialogModel(DialogOpen, "Test", "./test", []string{".txt"})

	t.Run("ValidFilename", func(t *testing.T) {
		err := dialog.validateFilename("valid_filename.txt")
		assert.NoError(t, err)
	})

	t.Run("EmptyFilename", func(t *testing.T) {
		err := dialog.validateFilename("")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "filename cannot be empty")
	})

	t.Run("InvalidCharacters", func(t *testing.T) {
		err := dialog.validateFilename("invalid<name>.txt")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid character")
	})

	t.Run("PathTraversal", func(t *testing.T) {
		err := dialog.validateFilename("../../../etc/passwd")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "parent directory references")
	})

	t.Run("AbsolutePath", func(t *testing.T) {
		err := dialog.validateFilename("/etc/passwd")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "must be relative")
	})

	t.Run("InvalidExtension", func(t *testing.T) {
		dialog.allowedExts = []string{".txt"}
		err := dialog.validateFilename("test.exe")
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "file extension .exe is not allowed")
	})
}

// TestFileDialogDirectoryLoading tests file dialog directory loading
func TestFileDialogDirectoryLoading(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "file_dialog_test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create some test files
	testFiles := []string{"test1.txt", "test2.md", "test3.txt"}
	for _, filename := range testFiles {
		filePath := filepath.Join(tempDir, filename)
		err := os.WriteFile(filePath, []byte("test content"), 0644)
		require.NoError(t, err)
	}

	// Create a subdirectory
	subDir := filepath.Join(tempDir, "subdir")
	err = os.MkdirAll(subDir, 0755)
	require.NoError(t, err)

	dialog := NewFileDialogModel(DialogOpen, "Test", tempDir, []string{".txt"})
	dialog.currentDir = tempDir

	// Load directory
	dialog.loadDirectory()

	// Check that files were loaded
	assert.Greater(t, len(dialog.items), 0, "Items should be loaded")

	// Check that only .txt files are loaded (based on extension filter)
	txtCount := 0
	for _, item := range dialog.items {
		if fileItem, ok := item.(FileItem); ok && !fileItem.isDir {
			ext := filepath.Ext(fileItem.name)
			if ext == ".txt" {
				txtCount++
			}
		}
	}
	assert.Equal(t, 2, txtCount, "Should load only .txt files")
}

// TestFileDialogSelection tests file dialog selection functionality
func TestFileDialogSelection(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "file_dialog_test_selection")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a test file
	testFile := filepath.Join(tempDir, "test.txt")
	err = os.WriteFile(testFile, []byte("test content"), 0644)
	require.NoError(t, err)

	dialog := NewFileDialogModel(DialogOpen, "Test", tempDir, []string{".txt"})
	dialog.currentDir = tempDir
	dialog.loadDirectory()

	// Find the test file in the items
	var testItem FileItem
	found := false
	for _, item := range dialog.items {
		if fileItem, ok := item.(FileItem); ok && fileItem.name == "test.txt" {
			testItem = fileItem
			found = true
			break
		}
	}
	require.True(t, found, "Test file should be found in items")

	// Test selection
	dialog.selectItem(testItem)
	assert.Equal(t, testFile, dialog.selectedFile, "Selected file should be set")
}

// TestFileItem tests the FileItem implementation
func TestFileItem(t *testing.T) {
	t.Run("FileItem", func(t *testing.T) {
		item := FileItem{
			name:    "test.txt",
			path:    "/path/to/test.txt",
			isDir:   false,
			size:    1024,
			modTime: time.Now(),
		}

		assert.Equal(t, "test.txt (1.0KB)", item.Title())
		assert.NotEmpty(t, item.Description())
		assert.Equal(t, "test.txt", item.FilterValue())
	})

	t.Run("DirectoryItem", func(t *testing.T) {
		item := FileItem{
			name:    "testdir",
			path:    "/path/to/testdir",
			isDir:   true,
			size:    0,
			modTime: time.Now(),
		}

		assert.Equal(t, "testdir/", item.Title())
		assert.Equal(t, "Directory", item.Description())
		assert.Equal(t, "testdir", item.FilterValue())
	})
}

// TestFileDialogNavigation tests file dialog navigation
func TestFileDialogNavigation(t *testing.T) {
	// Create a temporary directory structure for testing
	tempDir, err := os.MkdirTemp("", "file_dialog_test_nav")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a subdirectory
	subDir := filepath.Join(tempDir, "subdir")
	err = os.MkdirAll(subDir, 0755)
	require.NoError(t, err)

	dialog := NewFileDialogModel(DialogOpen, "Test", tempDir, []string{".txt"})
	dialog.currentDir = subDir

	// Test going up directory
	dialog.goUpDirectory()
	assert.Equal(t, tempDir, dialog.currentDir, "Should navigate to parent directory")
}

// TestFileDialogHiddenFiles tests file dialog hidden file handling
func TestFileDialogHiddenFiles(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "file_dialog_test_hidden")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create a regular file and a hidden file
	regularFile := filepath.Join(tempDir, "regular.txt")
	hiddenFile := filepath.Join(tempDir, ".hidden.txt")

	err = os.WriteFile(regularFile, []byte("test content"), 0644)
	require.NoError(t, err)

	err = os.WriteFile(hiddenFile, []byte("hidden content"), 0644)
	require.NoError(t, err)

	dialog := NewFileDialogModel(DialogOpen, "Test", tempDir, []string{".txt"})
	dialog.currentDir = tempDir

	// Load directory with showHidden = false
	dialog.showHidden = false
	dialog.loadDirectory()

	// Count visible items
	visibleCount := len(dialog.items)

	// Load directory with showHidden = true
	dialog.showHidden = true
	dialog.loadDirectory()

	// Should have more items when hidden files are shown
	assert.Greater(t, len(dialog.items), visibleCount, "Should show more files when hidden files are enabled")
}

// TestFileDialogRendering tests file dialog rendering
func TestFileDialogRendering(t *testing.T) {
	dialog := NewFileDialogModel(DialogOpen, "Test Dialog", "./test", []string{".txt"})
	dialog.SetDimensions(80, 24)
	dialog.visible = true

	// Test rendering when visible
	view := dialog.View()
	assert.NotEmpty(t, view, "View should not be empty when visible")
	assert.Contains(t, view, "Test Dialog", "View should contain dialog title")

	// Test rendering when hidden
	dialog.visible = false
	view = dialog.View()
	assert.Empty(t, view, "View should be empty when hidden")
}

// TestFileDialogErrorHandling tests file dialog error handling
func TestFileDialogErrorHandling(t *testing.T) {
	dialog := NewFileDialogModel(DialogOpen, "Test", "./nonexistent", []string{".txt"})
	dialog.currentDir = "./nonexistent"

	// Load directory that doesn't exist
	dialog.loadDirectory()

	// Should have error set
	assert.Error(t, dialog.err, "Should have error when loading nonexistent directory")

	// Error should be displayed in view
	dialog.visible = true
	view := dialog.View()
	assert.Contains(t, view, "Error:", "View should contain error message")
}

// BenchmarkFileDialogLoading benchmarks file dialog directory loading
func BenchmarkFileDialogLoading(b *testing.B) {
	// Create a temporary directory with many files for testing
	tempDir, err := os.MkdirTemp("", "file_dialog_bench")
	require.NoError(b, err)
	defer os.RemoveAll(tempDir)

	// Create many test files
	for i := 0; i < 100; i++ {
		filename := filepath.Join(tempDir, fmt.Sprintf("test%d.txt", i))
		err := os.WriteFile(filename, []byte("test content"), 0644)
		require.NoError(b, err)
	}

	dialog := NewFileDialogModel(DialogOpen, "Test", tempDir, []string{".txt"})
	dialog.currentDir = tempDir

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		dialog.loadDirectory()
	}
}
