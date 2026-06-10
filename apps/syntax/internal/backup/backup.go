package backup

import (
	"archive/zip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Backup represents a project backup
type Backup struct {
	Timestamp time.Time
	Path      string
	Size      int64
}

// GetBackupDir returns the backup directory for a project
func GetBackupDir(projectDir string) string {
	return filepath.Join(projectDir, ".syntax", "backups")
}

// CreateBackup creates a backup of the entire project
func CreateBackup(projectDir string) error {
	backupDir := GetBackupDir(projectDir)

	// Ensure backup directory exists
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return fmt.Errorf("failed to create backup directory: %w", err)
	}

	// Create backup filename with timestamp
	timestamp := time.Now().Format("2006-01-02_15-04-05")
	backupPath := filepath.Join(backupDir, fmt.Sprintf("backup_%s.zip", timestamp))

	// Create zip file
	zipFile, err := os.Create(backupPath)
	if err != nil {
		return fmt.Errorf("failed to create backup file: %w", err)
	}
	defer zipFile.Close()

	zipWriter := zip.NewWriter(zipFile)
	defer zipWriter.Close()

	// Walk through project directory and add files to zip
	err = filepath.Walk(projectDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip backup directory itself
		if strings.Contains(path, filepath.Join(".syntax", "backups")) {
			if info.IsDir() {
				return filepath.SkipDir
			}
			return nil
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(projectDir, path)
		if err != nil {
			return err
		}

		// Create zip entry
		zipEntry, err := zipWriter.Create(relPath)
		if err != nil {
			return err
		}

		// Copy file contents
		file, err := os.Open(path)
		if err != nil {
			return err
		}
		defer file.Close()

		_, err = io.Copy(zipEntry, file)
		return err
	})

	if err != nil {
		os.Remove(backupPath)
		return fmt.Errorf("failed to create backup: %w", err)
	}

	return nil
}

// ListBackups returns a list of all backups for a project
func ListBackups(projectDir string) ([]Backup, error) {
	backupDir := GetBackupDir(projectDir)

	// Check if backup directory exists
	if _, err := os.Stat(backupDir); os.IsNotExist(err) {
		return []Backup{}, nil
	}

	// Read backup directory
	entries, err := os.ReadDir(backupDir)
	if err != nil {
		return nil, fmt.Errorf("failed to read backup directory: %w", err)
	}

	backups := []Backup{}
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".zip") {
			continue
		}

		// Parse timestamp from filename
		name := strings.TrimSuffix(entry.Name(), ".zip")
		name = strings.TrimPrefix(name, "backup_")
		timestamp, err := time.Parse("2006-01-02_15-04-05", name)
		if err != nil {
			// Skip files with invalid format
			continue
		}

		// Get file info
		info, err := entry.Info()
		if err != nil {
			continue
		}

		backups = append(backups, Backup{
			Timestamp: timestamp,
			Path:      filepath.Join(backupDir, entry.Name()),
			Size:      info.Size(),
		})
	}

	// Sort by timestamp (newest first)
	sort.Slice(backups, func(i, j int) bool {
		return backups[i].Timestamp.After(backups[j].Timestamp)
	})

	return backups, nil
}

// RestoreBackup restores a project from a backup
func RestoreBackup(backupPath, projectDir string) error {
	// Open zip file
	zipReader, err := zip.OpenReader(backupPath)
	if err != nil {
		return fmt.Errorf("failed to open backup: %w", err)
	}
	defer zipReader.Close()

	// Create temporary directory for extraction
	tempDir := projectDir + "_restore_temp"
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}

	// Extract all files
	for _, file := range zipReader.File {
		// Create file path
		filePath := filepath.Join(tempDir, file.Name)

		// Create parent directories
		if err := os.MkdirAll(filepath.Dir(filePath), 0755); err != nil {
			os.RemoveAll(tempDir)
			return fmt.Errorf("failed to create directory: %w", err)
		}

		// Skip if it's a directory
		if file.FileInfo().IsDir() {
			continue
		}

		// Extract file
		srcFile, err := file.Open()
		if err != nil {
			os.RemoveAll(tempDir)
			return fmt.Errorf("failed to open file in backup: %w", err)
		}

		dstFile, err := os.Create(filePath)
		if err != nil {
			srcFile.Close()
			os.RemoveAll(tempDir)
			return fmt.Errorf("failed to create file: %w", err)
		}

		_, err = io.Copy(dstFile, srcFile)
		srcFile.Close()
		dstFile.Close()

		if err != nil {
			os.RemoveAll(tempDir)
			return fmt.Errorf("failed to extract file: %w", err)
		}
	}

	// Backup current project (in case restore fails)
	currentBackup := projectDir + "_pre_restore_backup"
	if err := os.Rename(projectDir, currentBackup); err != nil {
		os.RemoveAll(tempDir)
		return fmt.Errorf("failed to backup current project: %w", err)
	}

	// Move restored project to original location
	if err := os.Rename(tempDir, projectDir); err != nil {
		// Restore from backup on error
		os.Rename(currentBackup, projectDir)
		os.RemoveAll(tempDir)
		return fmt.Errorf("failed to restore project: %w", err)
	}

	// Remove pre-restore backup
	os.RemoveAll(currentBackup)

	return nil
}

// DeleteBackup deletes a specific backup
func DeleteBackup(backupPath string) error {
	return os.Remove(backupPath)
}

// CleanOldBackups removes backups older than a specified duration
func CleanOldBackups(projectDir string, maxAge time.Duration) error {
	backups, err := ListBackups(projectDir)
	if err != nil {
		return err
	}

	now := time.Now()
	for _, backup := range backups {
		if now.Sub(backup.Timestamp) > maxAge {
			if err := DeleteBackup(backup.Path); err != nil {
				return err
			}
		}
	}

	return nil
}
