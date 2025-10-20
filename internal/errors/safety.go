package errors

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/Kyanite/noise/internal/domain"
	"github.com/Kyanite/noise/internal/logging"
)

// BackupManager handles content backup and recovery
type BackupManager struct {
	backupDir  string
	maxBackups int
	logger     *logging.Logger
	mu         sync.RWMutex
}

// BackupInfo contains information about a backup
type BackupInfo struct {
	ID        string    `json:"id"`
	SongID    int       `json:"song_id"`
	Timestamp time.Time `json:"timestamp"`
	Type      string    `json:"type"` // "auto", "manual", "pre_operation"
	Size      int64     `json:"size"`
	Path      string    `json:"path"`
}

// NewBackupManager creates a new backup manager
func NewBackupManager(backupDir string, maxBackups int, logger *logging.Logger) (*BackupManager, error) {
	// Create backup directory if it doesn't exist
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create backup directory: %w", err)
	}

	return &BackupManager{
		backupDir:  backupDir,
		maxBackups: maxBackups,
		logger:     logger,
	}, nil
}

// CreateBackup creates a backup of the song content
func (bm *BackupManager) CreateBackup(song *domain.Song, backupType string) (*BackupInfo, error) {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	if song == nil {
		return nil, NewValidationError("song cannot be nil for backup", nil)
	}

	// Create backup info
	backupInfo := &BackupInfo{
		ID:        generateBackupID(),
		SongID:    song.ID,
		Timestamp: time.Now(),
		Type:      backupType,
	}

	// Serialize song content
	content, err := bm.serializeSong(song)
	if err != nil {
		return nil, NewParsingError("serialize_song", err).WithOperation("CreateBackup").WithComponent("backup_manager")
	}

	// Create backup file
	backupPath := bm.getBackupPath(backupInfo.ID)
	if err := os.WriteFile(backupPath, content, 0644); err != nil {
		return nil, NewFileError("write_backup", backupPath, err).WithOperation("CreateBackup").WithComponent("backup_manager")
	}

	// Get file size
	fileInfo, err := os.Stat(backupPath)
	if err != nil {
		bm.logger.Warn("Failed to get backup file info", "path", backupPath, "error", err)
	} else {
		backupInfo.Size = fileInfo.Size()
	}

	backupInfo.Path = backupPath

	// Clean up old backups if needed
	if err := bm.cleanupOldBackups(); err != nil {
		bm.logger.Warn("Failed to cleanup old backups", "error", err)
	}

	bm.logger.Debug("Backup created successfully", "id", backupInfo.ID, "song_id", song.ID, "type", backupType)
	return backupInfo, nil
}

// RestoreBackup restores a song from backup
func (bm *BackupManager) RestoreBackup(backupID string) (*domain.Song, error) {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	backupPath := bm.getBackupPath(backupID)

	// Check if backup exists
	if _, err := os.Stat(backupPath); os.IsNotExist(err) {
		return nil, NewFileError("backup_not_found", backupPath, err).WithOperation("RestoreBackup").WithComponent("backup_manager")
	}

	// Read backup content
	content, err := os.ReadFile(backupPath)
	if err != nil {
		return nil, NewFileError("read_backup", backupPath, err).WithOperation("RestoreBackup").WithComponent("backup_manager")
	}

	// Deserialize song
	song, err := bm.deserializeSong(content)
	if err != nil {
		return nil, NewParsingError("deserialize_song", err).WithOperation("RestoreBackup").WithComponent("backup_manager")
	}

	bm.logger.Info("Backup restored successfully", "backup_id", backupID, "song_id", song.ID, "title", song.Metadata.Title)
	return song, nil
}

// ListBackups returns a list of available backups for a song
func (bm *BackupManager) ListBackups(songID int) ([]*BackupInfo, error) {
	bm.mu.RLock()
	defer bm.mu.RUnlock()

	entries, err := os.ReadDir(bm.backupDir)
	if err != nil {
		return nil, NewFileError("read_backup_dir", bm.backupDir, err).WithOperation("ListBackups").WithComponent("backup_manager")
	}

	var backups []*BackupInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		// Parse backup filename to extract info
		backupInfo, err := bm.parseBackupFilename(entry.Name())
		if err != nil {
			continue // Skip invalid backup files
		}

		// Filter by song ID if specified
		if songID > 0 && backupInfo.SongID != songID {
			continue
		}

		// Get file info
		fileInfo, err := entry.Info()
		if err != nil {
			bm.logger.Warn("Failed to get backup file info", "file", entry.Name(), "error", err)
			continue
		}

		backupInfo.Size = fileInfo.Size()
		backupInfo.Path = filepath.Join(bm.backupDir, entry.Name())

		backups = append(backups, backupInfo)
	}

	return backups, nil
}

// DeleteBackup deletes a specific backup
func (bm *BackupManager) DeleteBackup(backupID string) error {
	bm.mu.Lock()
	defer bm.mu.Unlock()

	backupPath := bm.getBackupPath(backupID)

	if err := os.Remove(backupPath); err != nil {
		return NewFileError("delete_backup", backupPath, err).WithOperation("DeleteBackup").WithComponent("backup_manager")
	}

	bm.logger.Debug("Backup deleted successfully", "backup_id", backupID)
	return nil
}

// Helper methods

func (bm *BackupManager) serializeSong(song *domain.Song) ([]byte, error) {
	return json.MarshalIndent(song, "", "  ")
}

func (bm *BackupManager) deserializeSong(content []byte) (*domain.Song, error) {
	var song domain.Song
	if err := json.Unmarshal(content, &song); err != nil {
		return nil, err
	}
	return &song, nil
}

func (bm *BackupManager) getBackupPath(backupID string) string {
	return filepath.Join(bm.backupDir, fmt.Sprintf("backup_%s.json", backupID))
}

func (bm *BackupManager) parseBackupFilename(filename string) (*BackupInfo, error) {
	// Expected format: backup_{backupID}.json
	if len(filename) < 20 || filename[:7] != "backup_" || filename[len(filename)-5:] != ".json" {
		return nil, fmt.Errorf("invalid backup filename format")
	}

	backupID := filename[7 : len(filename)-5]

	// For now, return basic info - in a full implementation, this would parse more metadata
	backupInfo := &BackupInfo{
		ID: backupID,
		// SongID and other fields would need to be stored in the backup file or filename
	}

	// Log backup metadata for integrity tracking
	bm.logger.Debug("Parsed backup filename", "filename", filename, "backup_id", backupID)

	return backupInfo, nil
}

func (bm *BackupManager) cleanupOldBackups() error {
	entries, err := os.ReadDir(bm.backupDir)
	if err != nil {
		return err
	}

	if len(entries) <= bm.maxBackups {
		return nil
	}

	// Sort by modification time and remove oldest
	type backupEntry struct {
		name    string
		modTime time.Time
	}

	var backupEntries []backupEntry
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		backupEntries = append(backupEntries, backupEntry{
			name:    entry.Name(),
			modTime: info.ModTime(),
		})
	}

	// Sort by modification time (oldest first)
	for i := 0; i < len(backupEntries)-1; i++ {
		for j := i + 1; j < len(backupEntries); j++ {
			if backupEntries[i].modTime.After(backupEntries[j].modTime) {
				backupEntries[i], backupEntries[j] = backupEntries[j], backupEntries[i]
			}
		}
	}

	// Remove oldest backups
	toRemove := len(backupEntries) - bm.maxBackups
	for i := 0; i < toRemove; i++ {
		backupPath := filepath.Join(bm.backupDir, backupEntries[i].name)
		if err := os.Remove(backupPath); err != nil {
			bm.logger.Warn("Failed to remove old backup", "file", backupEntries[i].name, "error", err)
		} else {
			bm.logger.Debug("Removed old backup", "file", backupEntries[i].name)
		}
	}

	return nil
}

func generateBackupID() string {
	return fmt.Sprintf("%d_%d", time.Now().Unix(), time.Now().UnixNano()%1000000)
}

// RecoveryManager handles recovery from various failure scenarios
type RecoveryManager struct {
	backupManager      *BackupManager
	logger             *logging.Logger
	recoveryStrategies map[string]RecoveryStrategy
}

// RecoveryStrategy defines a strategy for recovering from errors
type RecoveryStrategy func(ctx context.Context, err error) error

// NewRecoveryManager creates a new recovery manager
func NewRecoveryManager(backupManager *BackupManager, logger *logging.Logger) *RecoveryManager {
	rm := &RecoveryManager{
		backupManager:      backupManager,
		logger:             logger,
		recoveryStrategies: make(map[string]RecoveryStrategy),
	}

	// Register default recovery strategies
	rm.registerDefaultStrategies()

	return rm
}

// RegisterStrategy registers a recovery strategy for a specific error type
func (rm *RecoveryManager) RegisterStrategy(errorType string, strategy RecoveryStrategy) {
	rm.recoveryStrategies[errorType] = strategy
}

// Recover attempts to recover from an error using registered strategies
func (rm *RecoveryManager) Recover(ctx context.Context, err error, context map[string]interface{}) error {
	if err == nil {
		return nil
	}

	// Find appropriate recovery strategy
	var strategy RecoveryStrategy

	if appErr, ok := err.(*AppError); ok {
		strategy = rm.recoveryStrategies[string(appErr.Category)]
	}

	// If no specific strategy found, try generic strategies
	if strategy == nil {
		strategy = rm.recoveryStrategies["generic"]
	}

	// If still no strategy, return error
	if strategy == nil {
		return err
	}

	// Attempt recovery
	rm.logger.Info("Attempting error recovery", "error", err.Error())

	if err := strategy(ctx, err); err != nil {
		rm.logger.Error("Recovery strategy failed", "error", err)
		return err
	}

	rm.logger.Info("Recovery successful")
	return nil
}

// registerDefaultStrategies registers default recovery strategies
func (rm *RecoveryManager) registerDefaultStrategies() {
	// Generic recovery strategy
	rm.RegisterStrategy("generic", func(ctx context.Context, err error) error {
		// For generic errors, just log and continue
		rm.logger.Warn("Generic recovery: logging error and continuing", "error", err)
		return nil
	})

	// File system recovery strategy
	rm.RegisterStrategy("file", func(ctx context.Context, err error) error {
		// For file errors, attempt to restore from backup if available
		rm.logger.Info("Attempting file system recovery")

		// This would be implemented with specific file recovery logic
		return nil
	})

	// Database recovery strategy
	rm.RegisterStrategy("database", func(ctx context.Context, err error) error {
		// For database errors, attempt to reconnect or restore from backup
		rm.logger.Info("Attempting database recovery")

		// This would be implemented with specific database recovery logic
		return nil
	})

	// Memory recovery strategy
	rm.RegisterStrategy("resource", func(ctx context.Context, err error) error {
		// For resource errors, attempt garbage collection or graceful degradation
		rm.logger.Info("Attempting resource recovery")

		// Force garbage collection if possible
		// This is Go-specific and would need runtime.GC() call

		return nil
	})
}

// GracefulDegradation handles graceful degradation when errors occur
type GracefulDegradation struct {
	features map[string]bool
	logger   *logging.Logger
	mu       sync.RWMutex
}

// NewGracefulDegradation creates a new graceful degradation handler
func NewGracefulDegradation(logger *logging.Logger) *GracefulDegradation {
	return &GracefulDegradation{
		features: make(map[string]bool),
		logger:   logger,
	}
}

// DisableFeature disables a feature due to errors
func (gd *GracefulDegradation) DisableFeature(feature string, reason error) {
	gd.mu.Lock()
	defer gd.mu.Unlock()

	gd.features[feature] = false
	gd.logger.Warn("Feature disabled due to error", "feature", feature, "reason", reason)

	// Show user notification
	ShowGlobalWarning(
		"Feature Temporarily Disabled",
		fmt.Sprintf("The %s feature has been temporarily disabled due to an error. It will be re-enabled automatically when the issue is resolved.", feature),
	)
}

// EnableFeature enables a previously disabled feature
func (gd *GracefulDegradation) EnableFeature(feature string) {
	gd.mu.Lock()
	defer gd.mu.Unlock()

	gd.features[feature] = true
	gd.logger.Info("Feature re-enabled", "feature", feature)

	// Show user notification
	ShowGlobalInfo(
		"Feature Re-enabled",
		fmt.Sprintf("The %s feature has been re-enabled.", feature),
	)
}

// IsFeatureEnabled checks if a feature is enabled
func (gd *GracefulDegradation) IsFeatureEnabled(feature string) bool {
	gd.mu.RLock()
	defer gd.mu.RUnlock()

	enabled, exists := gd.features[feature]
	return !exists || enabled // Default to enabled if not in map
}

// GetDisabledFeatures returns a list of currently disabled features
func (gd *GracefulDegradation) GetDisabledFeatures() []string {
	gd.mu.RLock()
	defer gd.mu.RUnlock()

	var disabled []string
	for feature, enabled := range gd.features {
		if !enabled {
			disabled = append(disabled, feature)
		}
	}

	return disabled
}

// AutoRecovery handles automatic recovery on application startup
type AutoRecovery struct {
	backupManager *BackupManager
	logger        *logging.Logger
}

// NewAutoRecovery creates a new auto-recovery handler
func NewAutoRecovery(backupManager *BackupManager, logger *logging.Logger) *AutoRecovery {
	return &AutoRecovery{
		backupManager: backupManager,
		logger:        logger,
	}
}

// AttemptRecovery attempts to recover from a previous crash or error
func (ar *AutoRecovery) AttemptRecovery() ([]*domain.Song, error) {
	ar.logger.Info("Starting auto-recovery check")

	var recoveredSongs []*domain.Song

	// Look for emergency backup files
	emergencyDir := filepath.Join(ar.backupManager.backupDir, "emergency")
	if _, err := os.Stat(emergencyDir); !os.IsNotExist(err) {
		songs, err := ar.recoverFromEmergencyBackups(emergencyDir)
		if err != nil {
			ar.logger.Error("Failed to recover from emergency backups", "error", err)
		} else {
			recoveredSongs = append(recoveredSongs, songs...)
		}
	}

	// Look for auto-save files that might not have been committed
	autoSaveDir := filepath.Join(ar.backupManager.backupDir, "autosave")
	if _, err := os.Stat(autoSaveDir); !os.IsNotExist(err) {
		songs, err := ar.recoverFromAutoSave(autoSaveDir)
		if err != nil {
			ar.logger.Error("Failed to recover from auto-save", "error", err)
		} else {
			recoveredSongs = append(recoveredSongs, songs...)
		}
	}

	if len(recoveredSongs) > 0 {
		ar.logger.Info("Auto-recovery completed", "songs_recovered", len(recoveredSongs))

		// Show user notification
		ShowGlobalSuccess(
			"Recovery Completed",
			fmt.Sprintf("Successfully recovered %d song(s) from previous session.", len(recoveredSongs)),
		)
	} else {
		ar.logger.Debug("No songs to recover")
	}

	return recoveredSongs, nil
}

// recoverFromEmergencyBackups recovers songs from emergency backup files
func (ar *AutoRecovery) recoverFromEmergencyBackups(emergencyDir string) ([]*domain.Song, error) {
	entries, err := os.ReadDir(emergencyDir)
	if err != nil {
		return nil, err
	}

	var songs []*domain.Song
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		backupPath := filepath.Join(emergencyDir, entry.Name())
		content, err := os.ReadFile(backupPath)
		if err != nil {
			ar.logger.Warn("Failed to read emergency backup", "file", entry.Name(), "error", err)
			continue
		}

		song, err := ar.backupManager.deserializeSong(content)
		if err != nil {
			ar.logger.Warn("Failed to deserialize emergency backup", "file", entry.Name(), "error", err)
			continue
		}

		songs = append(songs, song)
		ar.logger.Info("Recovered song from emergency backup", "file", entry.Name(), "song_id", song.ID)
	}

	return songs, nil
}

// recoverFromAutoSave recovers songs from auto-save files
func (ar *AutoRecovery) recoverFromAutoSave(autoSaveDir string) ([]*domain.Song, error) {
	// This would implement auto-save recovery logic
	// For now, return empty slice
	return []*domain.Song{}, nil
}

// Global instances
var (
	globalBackupManager       *BackupManager
	globalRecoveryManager     *RecoveryManager
	globalGracefulDegradation *GracefulDegradation
	globalAutoRecovery        *AutoRecovery
)

// InitializeSafetySystems initializes all safety systems
func InitializeSafetySystems(dataDir string, logger *logging.Logger) error {
	backupDir := filepath.Join(dataDir, "backups")

	// Initialize backup manager
	backupManager, err := NewBackupManager(backupDir, 50, logger)
	if err != nil {
		return fmt.Errorf("failed to initialize backup manager: %w", err)
	}

	// Initialize recovery manager
	recoveryManager := NewRecoveryManager(backupManager, logger)

	// Initialize graceful degradation
	gracefulDegradation := NewGracefulDegradation(logger)

	// Initialize auto-recovery
	autoRecovery := NewAutoRecovery(backupManager, logger)

	// Set global instances
	globalBackupManager = backupManager
	globalRecoveryManager = recoveryManager
	globalGracefulDegradation = gracefulDegradation
	globalAutoRecovery = autoRecovery

	logger.Info("Safety systems initialized successfully")
	return nil
}

// GetGlobalBackupManager returns the global backup manager
func GetGlobalBackupManager() *BackupManager {
	return globalBackupManager
}

// GetGlobalRecoveryManager returns the global recovery manager
func GetGlobalRecoveryManager() *RecoveryManager {
	return globalRecoveryManager
}

// GetGlobalGracefulDegradation returns the global graceful degradation handler
func GetGlobalGracefulDegradation() *GracefulDegradation {
	return globalGracefulDegradation
}

// GetGlobalAutoRecovery returns the global auto-recovery handler
func GetGlobalAutoRecovery() *AutoRecovery {
	return globalAutoRecovery
}
