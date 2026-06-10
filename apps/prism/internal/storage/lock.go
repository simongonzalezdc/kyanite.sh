package storage

import (
	"fmt"
	"os"
	"path/filepath"
)

// FileLock represents a file lock
type FileLock struct {
	path string
	file *os.File
}

// LockFile acquires an advisory lock on a file
func LockFile(path string) (*FileLock, error) {
	// For now, simple implementation using lock files

	lockPath := path + ".lock"
	lockFile, err := os.OpenFile(lockPath, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to acquire lock: %w", err)
	}

	return &FileLock{
		path: lockPath,
		file: lockFile,
	}, nil
}

// UnlockFile releases a file lock
func UnlockFile(lock *FileLock) error {
	if lock == nil {
		return nil
	}

	if err := lock.file.Close(); err != nil {
		return fmt.Errorf("failed to close lock file: %w", err)
	}
	return os.Remove(lock.path)
}

// AtomicWrite writes data atomically by writing to a temp file then renaming
func AtomicWrite(path string, data []byte) error {

	dir := filepath.Dir(path)
	tmpFile, err := os.CreateTemp(dir, ".tmp-*")
	if err != nil {
		return fmt.Errorf("failed to create temp file: %w", err)
	}
	tmpPath := tmpFile.Name()

	// Ensure cleanup on error
	var success bool
	defer func() {
		if !success {
			_ = tmpFile.Close()
			_ = os.Remove(tmpPath)
		}
	}()

	// Write data
	if _, err = tmpFile.Write(data); err != nil {
		return fmt.Errorf("failed to write temp file: %w", err)
	}

	// Close temp file
	if err = tmpFile.Close(); err != nil {
		return fmt.Errorf("failed to close temp file: %w", err)
	}

	// Atomic rename
	if err = os.Rename(tmpPath, path); err != nil {
		return fmt.Errorf("failed to rename temp file: %w", err)
	}

	success = true
	return nil
}
