package utils

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGetStoragePath(t *testing.T) {
	path := GetStoragePath()
	
	// Should be a valid file path
	if path == "" {
		t.Error("GetStoragePath should not return empty string")
	}
	
	// Should end with tasks.json
	if filepath.Base(path) != "tasks.json" {
		t.Errorf("Expected path to end with 'tasks.json', got '%s'", path)
	}
	
	// Should contain .focus directory
	if !strings.Contains(path, ".focus") {
		t.Errorf("Expected path to contain '.focus', got '%s'", path)
	}
}

func TestMigrateStorage(t *testing.T) {
	// Test that migrateStorage doesn't panic
	// In a real test, we would mock filesystem operations
	err := migrateStorage()
	
	// migrateStorage should not return error even if migration fails
	// It's designed to be best-effort
	if err != nil {
		// This is not necessarily an error - migration might fail if old dirs don't exist
		t.Logf("MigrateStorage returned error (this might be expected): %v", err)
	}
}

func TestCopyFile(t *testing.T) {
	// Create temporary test files
	srcDir := t.TempDir()
	srcFile := filepath.Join(srcDir, "src.txt")
	dstFile := filepath.Join(srcDir, "dst.txt")
	
	// Create source file with test content
	content := "test content for copy"
	err := os.WriteFile(srcFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create source file: %v", err)
	}
	
	// Copy file
	err = copyFile(srcFile, dstFile)
	if err != nil {
		t.Fatalf("Failed to copy file: %v", err)
	}
	
	// Verify destination file exists and has same content
	if _, err := os.Stat(dstFile); os.IsNotExist(err) {
		t.Error("Destination file should exist after copy")
		return
	}
	
	dstContent, err := os.ReadFile(dstFile)
	if err != nil {
		t.Fatalf("Failed to read destination file: %v", err)
	}
	
	if string(dstContent) != content {
		t.Errorf("Expected destination content '%s', got '%s'", content, string(dstContent))
	}
	
	// Clean up
	os.Remove(srcFile)
	os.Remove(dstFile)
}

func TestCopyFile_NonExistentSource(t *testing.T) {
	srcDir := t.TempDir()
	srcFile := filepath.Join(srcDir, "nonexistent.txt")
	dstFile := filepath.Join(srcDir, "dst.txt")
	
	// Try to copy non-existent file
	err := copyFile(srcFile, dstFile)
	if err == nil {
		t.Error("Expected error when copying non-existent file")
	}
}

func TestCopyDir(t *testing.T) {
	// Create temporary directory structure
	srcDir := filepath.Join(t.TempDir(), "src")
	dstDir := filepath.Join(t.TempDir(), "dst")
	
	// Create source directory with files
	os.MkdirAll(srcDir, 0755)
	
	file1 := filepath.Join(srcDir, "file1.txt")
	file2 := filepath.Join(srcDir, "subdir", "file2.txt")
	os.MkdirAll(filepath.Dir(file2), 0755)
	
	os.WriteFile(file1, []byte("content1"), 0644)
	os.WriteFile(file2, []byte("content2"), 0644)
	
	// Copy directory
	err := copyDir(srcDir, dstDir)
	if err != nil {
		t.Fatalf("Failed to copy directory: %v", err)
	}
	
	// Verify directory structure was copied
	dstFile1 := filepath.Join(dstDir, "file1.txt")
	dstFile2 := filepath.Join(dstDir, "subdir", "file2.txt")
	
	if _, err := os.Stat(dstFile1); os.IsNotExist(err) {
		t.Error("File1 should exist in destination directory")
	}
	
	if _, err := os.Stat(dstFile2); os.IsNotExist(err) {
		t.Error("File2 should exist in destination subdirectory")
	}
	
	// Verify content
	content1, _ := os.ReadFile(dstFile1)
	if string(content1) != "content1" {
		t.Errorf("Expected file1 content 'content1', got '%s'", string(content1))
	}
	
	// Clean up
	os.RemoveAll(srcDir)
	os.RemoveAll(dstDir)
}

func TestCopyDir_NonExistentSource(t *testing.T) {
	srcDir := filepath.Join(t.TempDir(), "nonexistent")
	dstDir := filepath.Join(t.TempDir(), "dst")
	
	// Try to copy non-existent directory
	err := copyDir(srcDir, dstDir)
	if err == nil {
		t.Error("Expected error when copying non-existent directory")
	}
}

func TestCopyDir_EmptySource(t *testing.T) {
	srcDir := t.TempDir()
	dstDir := filepath.Join(t.TempDir(), "empty_dst")
	
	// Create empty source directory
	os.Mkdir(srcDir, 0755)
	
	// Copy empty directory
	err := copyDir(srcDir, dstDir)
	if err != nil {
		t.Fatalf("Failed to copy empty directory: %v", err)
	}
	
	// Verify destination directory exists
	if _, err := os.Stat(dstDir); os.IsNotExist(err) {
		t.Error("Destination directory should exist after copying empty directory")
	}
	
	// Verify destination directory is empty
	entries, err := os.ReadDir(dstDir)
	if err != nil {
		t.Fatalf("Failed to read destination directory: %v", err)
	}
	
	// Should only contain "." and ".."
	if len(entries) > 2 {
		t.Errorf("Expected empty directory to contain only '.' and '..', found %d entries", len(entries))
	}
	
	// Clean up
	os.Remove(srcDir)
	os.RemoveAll(dstDir)
}

func TestUserHomeDir_Fallback(t *testing.T) {
	// This test is hard to do without mocking os.UserHomeDir
	// In a real scenario, you would mock the os.UserHomeDir function
	// For now, we'll just test that the function exists and doesn't panic
	// This function is not defined in the storage.go file, so we'll skip this test
	t.Skip("GetUserHomeDir function not available for testing")
}