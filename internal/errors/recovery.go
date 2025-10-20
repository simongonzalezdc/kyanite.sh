package errors

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/Kyanite/noise/internal/domain"
	"github.com/Kyanite/noise/internal/logging"
)

// FileCorruptionDetector detects corrupted files
type FileCorruptionDetector struct {
	logger           *logging.Logger
	corruptionRules  map[string]CorruptionRule
	recoveryHandlers map[string]RecoveryHandler
	mu               sync.RWMutex
}

// CorruptionRule defines how to detect corruption in a file type
type CorruptionRule struct {
	FileExtension string
	Description   string
	Detector      func(filePath string, content []byte) error
}

// RecoveryHandler defines how to recover from corruption for a file type
type RecoveryHandler struct {
	Description string
	Handler     func(filePath string, backupPath string) error
}

// NewFileCorruptionDetector creates a new file corruption detector
func NewFileCorruptionDetector(logger *logging.Logger) *FileCorruptionDetector {
	detector := &FileCorruptionDetector{
		logger:           logger,
		corruptionRules:  make(map[string]CorruptionRule),
		recoveryHandlers: make(map[string]RecoveryHandler),
	}

	// Register default corruption rules and handlers
	detector.registerDefaultRules()

	return detector
}

// registerDefaultRules registers default corruption detection rules
func (fcd *FileCorruptionDetector) registerDefaultRules() {
	// JSON file corruption detection
	fcd.corruptionRules[".json"] = CorruptionRule{
		FileExtension: ".json",
		Description:   "JSON syntax validation",
		Detector: func(filePath string, content []byte) error {
			if !json.Valid(content) {
				return NewParsingError("json_syntax", fmt.Errorf("invalid JSON syntax in file: %s", filePath))
			}

			// Try to unmarshal to check structure
			var temp interface{}
			if err := json.Unmarshal(content, &temp); err != nil {
				return NewParsingError("json_structure", fmt.Errorf("invalid JSON structure in file: %s, error: %w", filePath, err))
			}

			return nil
		},
	}

	// Text file corruption detection (basic)
	fcd.corruptionRules[".txt"] = CorruptionRule{
		FileExtension: ".txt",
		Description:   "Text file encoding validation",
		Detector: func(filePath string, content []byte) error {
			// Check for null bytes which indicate binary corruption
			if len(content) > 0 {
				for i, b := range content {
					if b == 0 && i < len(content)-1 { // Allow trailing null
						return NewParsingError("binary_corruption", fmt.Errorf("binary corruption detected in text file: %s at position %d", filePath, i))
					}
				}
			}
			return nil
		},
	}

	// Markdown file corruption detection
	fcd.corruptionRules[".md"] = CorruptionRule{
		FileExtension: ".md",
		Description:   "Markdown structure validation",
		Detector: func(filePath string, content []byte) error {
			text := string(content)

			// Check for unmatched markdown syntax that might indicate corruption
			lines := strings.Split(text, "\n")
			openCodeBlocks := 0

			for i, line := range lines {
				// Count code blocks
				codeBlockCount := strings.Count(line, "```")
				if codeBlockCount%2 == 1 {
					openCodeBlocks++
				}

				// Check for extremely long lines that might indicate corruption
				if len(line) > 10000 {
					fcd.logger.Warn("Extremely long line detected", "file", filePath, "line", i+1, "length", len(line))
				}
			}

			// Check for unmatched code blocks
			if openCodeBlocks%2 == 1 {
				return NewParsingError("markdown_syntax", fmt.Errorf("unmatched code blocks in markdown file: %s", filePath))
			}

			return nil
		},
	}

	// Register recovery handlers
	fcd.recoveryHandlers[".json"] = RecoveryHandler{
		Description: "JSON file recovery using backup",
		Handler: func(filePath string, backupPath string) error {
			return fcd.recoverJSONFile(filePath, backupPath)
		},
	}

	fcd.recoveryHandlers[".txt"] = RecoveryHandler{
		Description: "Text file recovery using backup",
		Handler: func(filePath string, backupPath string) error {
			return fcd.recoverTextFile(filePath, backupPath)
		},
	}

	fcd.recoveryHandlers[".md"] = RecoveryHandler{
		Description: "Markdown file recovery using backup",
		Handler: func(filePath string, backupPath string) error {
			return fcd.recoverTextFile(filePath, backupPath)
		},
	}
}

// DetectCorruption checks if a file is corrupted
func (fcd *FileCorruptionDetector) DetectCorruption(filePath string) error {
	fcd.mu.RLock()
	defer fcd.mu.RUnlock()

	// Get file extension
	ext := strings.ToLower(filepath.Ext(filePath))

	// Find applicable rule
	rule, exists := fcd.corruptionRules[ext]
	if !exists {
		// No specific rule for this file type, perform basic checks
		return fcd.basicCorruptionCheck(filePath)
	}

	// Read file content
	content, err := os.ReadFile(filePath)
	if err != nil {
		return NewFileError("read_file", filePath, err).WithOperation("DetectCorruption").WithComponent("corruption_detector")
	}

	// Apply corruption detection rule
	return rule.Detector(filePath, content)
}

// basicCorruptionCheck performs basic corruption checks for unknown file types
func (fcd *FileCorruptionDetector) basicCorruptionCheck(filePath string) error {
	// Check if file exists and is readable
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return NewFileError("stat_file", filePath, err).WithOperation("basicCorruptionCheck").WithComponent("corruption_detector")
	}

	// Check file size (0 bytes might indicate corruption)
	if fileInfo.Size() == 0 {
		return NewFileError("empty_file", filePath, fmt.Errorf("file is empty")).WithOperation("basicCorruptionCheck").WithComponent("corruption_detector")
	}

	// Check if file size is reasonable (not extremely large)
	maxSize := int64(100 * 1024 * 1024) // 100MB
	if fileInfo.Size() > maxSize {
		fcd.logger.Warn("File size exceeds reasonable limits", "file", filePath, "size", fileInfo.Size())
	}

	return nil
}

// RecoverFile attempts to recover a corrupted file using available backups
func (fcd *FileCorruptionDetector) RecoverFile(filePath string) error {
	fcd.mu.RLock()
	defer fcd.mu.RUnlock()

	ext := strings.ToLower(filepath.Ext(filePath))

	// Find applicable recovery handler
	handler, exists := fcd.recoveryHandlers[ext]
	if !exists {
		return NewAppError("NO_RECOVERY_HANDLER", fmt.Sprintf("No recovery handler available for file type: %s", ext), nil, CategoryFile, SeverityMedium, RecoveryManual)
	}

	// Look for backup files
	backupPath := filePath + ".backup"
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		// Try alternative backup locations
		backupPath = filepath.Join(filepath.Dir(filePath), "backups", filepath.Base(filePath))
		if _, err := os.Stat(backupPath); os.IsNotExist(err) {
			return NewAppError("NO_BACKUP", "No backup file found for recovery", nil, CategoryFile, SeverityHigh, RecoveryManual)
		}
	}

	// Attempt recovery
	fcd.logger.Info("Attempting file recovery", "file", filePath, "backup", backupPath)

	if err := handler.Handler(filePath, backupPath); err != nil {
		return NewFileError("recovery_failed", filePath, err).WithOperation("RecoverFile").WithComponent("corruption_detector")
	}

	fcd.logger.Info("File recovery completed successfully", "file", filePath)

	// Show success notification
	ShowGlobalSuccess("File Recovered", fmt.Sprintf("Successfully recovered corrupted file: %s", filepath.Base(filePath)))

	return nil
}

// recoverJSONFile recovers a JSON file from backup
func (fcd *FileCorruptionDetector) recoverJSONFile(filePath string, backupPath string) error {
	// Read backup content
	backupContent, err := os.ReadFile(backupPath)
	if err != nil {
		return NewFileError("read_backup", backupPath, err)
	}

	// Validate backup content
	if !json.Valid(backupContent) {
		return NewParsingError("invalid_backup", fmt.Errorf("backup file is also corrupted: %s", backupPath))
	}

	// Create recovery copy of corrupted file
	corruptedCopy := filePath + ".corrupted"
	if err := os.Rename(filePath, corruptedCopy); err != nil {
		fcd.logger.Warn("Failed to create corrupted file copy", "error", err)
	}

	// Restore from backup
	if err := os.WriteFile(filePath, backupContent, 0644); err != nil {
		return NewFileError("write_recovered", filePath, err)
	}

	fcd.logger.Info("JSON file recovered successfully", "original", filePath, "backup", backupPath)
	return nil
}

// recoverTextFile recovers a text file from backup
func (fcd *FileCorruptionDetector) recoverTextFile(filePath string, backupPath string) error {
	// Read backup content
	backupContent, err := os.ReadFile(backupPath)
	if err != nil {
		return NewFileError("read_backup", backupPath, err)
	}

	// Create recovery copy of corrupted file
	corruptedCopy := filePath + ".corrupted"
	if err := os.Rename(filePath, corruptedCopy); err != nil {
		fcd.logger.Warn("Failed to create corrupted file copy", "error", err)
	}

	// Restore from backup
	if err := os.WriteFile(filePath, backupContent, 0644); err != nil {
		return NewFileError("write_recovered", filePath, err)
	}

	fcd.logger.Info("Text file recovered successfully", "original", filePath, "backup", backupPath)
	return nil
}

// ScanDirectory scans a directory for corrupted files
func (fcd *FileCorruptionDetector) ScanDirectory(dirPath string) ([]string, error) {
	var corruptedFiles []string

	err := filepath.Walk(dirPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fcd.logger.Warn("Error accessing file during scan", "file", path, "error", err)
			return nil // Continue scanning other files
		}

		// Skip directories
		if info.IsDir() {
			return nil
		}

		// Check for corruption
		if err := fcd.DetectCorruption(path); err != nil {
			fcd.logger.Warn("Corrupted file detected", "file", path, "error", err)
			corruptedFiles = append(corruptedFiles, path)
		}

		return nil
	})

	if err != nil {
		return nil, NewFileError("scan_directory", dirPath, err).WithOperation("ScanDirectory").WithComponent("corruption_detector")
	}

	return corruptedFiles, nil
}

// AutoRecoverAll attempts to auto-recover all corrupted files in a directory
func (fcd *FileCorruptionDetector) AutoRecoverAll(dirPath string) (int, error) {
	corruptedFiles, err := fcd.ScanDirectory(dirPath)
	if err != nil {
		return 0, err
	}

	recovered := 0
	for _, filePath := range corruptedFiles {
		if err := fcd.RecoverFile(filePath); err != nil {
			fcd.logger.Error("Failed to auto-recover file", "file", filePath, "error", err)
		} else {
			recovered++
		}
	}

	if recovered > 0 {
		ShowGlobalSuccess("Auto-Recovery Completed", fmt.Sprintf("Successfully recovered %d corrupted files.", recovered))
	}

	return recovered, nil
}

// EnhancedBackupManager extends the existing backup manager with corruption detection
type EnhancedBackupManager struct {
	*BackupManager
	corruptionDetector *FileCorruptionDetector
	logger             *logging.Logger
}

// NewEnhancedBackupManager creates a new enhanced backup manager
func NewEnhancedBackupManager(backupDir string, maxBackups int, logger *logging.Logger) (*EnhancedBackupManager, error) {
	baseManager, err := NewBackupManager(backupDir, maxBackups, logger)
	if err != nil {
		return nil, err
	}

	return &EnhancedBackupManager{
		BackupManager:      baseManager,
		corruptionDetector: NewFileCorruptionDetector(logger),
		logger:             logger,
	}, nil
}

// CreateVerifiedBackup creates a backup and verifies its integrity
func (ebm *EnhancedBackupManager) CreateVerifiedBackup(song *domain.Song, backupType string) (*BackupInfo, error) {
	// Create regular backup
	backupInfo, err := ebm.BackupManager.CreateBackup(song, backupType)
	if err != nil {
		return nil, err
	}

	// Verify backup integrity
	backupPath := ebm.getBackupPath(backupInfo.ID)
	if err := ebm.corruptionDetector.DetectCorruption(backupPath); err != nil {
		ebm.logger.Error("Backup verification failed", "backup_id", backupInfo.ID, "error", err)

		// Remove corrupted backup
		os.Remove(backupPath)

		return nil, NewFileError("backup_corrupted", backupPath, err).WithOperation("CreateVerifiedBackup").WithComponent("enhanced_backup_manager")
	}

	ebm.logger.Debug("Backup verified successfully", "backup_id", backupInfo.ID)
	return backupInfo, nil
}

// RecoverCorruptedSong attempts to recover a corrupted song file
func (ebm *EnhancedBackupManager) RecoverCorruptedSong(songFilePath string) (*domain.Song, error) {
	ebm.logger.Info("Attempting to recover corrupted song file", "file", songFilePath)

	// Check if file is corrupted
	if err := ebm.corruptionDetector.DetectCorruption(songFilePath); err == nil {
		// File is not corrupted
		return nil, NewAppError("FILE_NOT_CORRUPTED", "File is not corrupted", nil, CategoryFile, SeverityLow, RecoveryNone)
	}

	// Attempt recovery
	if err := ebm.corruptionDetector.RecoverFile(songFilePath); err != nil {
		return nil, err
	}

	// Try to load the recovered song
	content, err := os.ReadFile(songFilePath)
	if err != nil {
		return nil, NewFileError("read_recovered", songFilePath, err)
	}

	song, err := ebm.deserializeSong(content)
	if err != nil {
		return nil, NewParsingError("deserialize_recovered", err)
	}

	ebm.logger.Info("Song file recovered successfully", "file", songFilePath, "song_id", song.ID)
	return song, nil
}

// Global corruption detector instance
var globalCorruptionDetector *FileCorruptionDetector

// InitializeCorruptionDetector initializes the global corruption detector
func InitializeCorruptionDetector(logger *logging.Logger) {
	globalCorruptionDetector = NewFileCorruptionDetector(logger)
	logger.Info("Corruption detector initialized")
}

// GetGlobalCorruptionDetector returns the global corruption detector
func GetGlobalCorruptionDetector() *FileCorruptionDetector {
	return globalCorruptionDetector
}
