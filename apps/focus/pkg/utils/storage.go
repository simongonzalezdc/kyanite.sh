package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	// Storage path cache
	cachedStoragePath string
	// Ensure migration runs only once
	migrationOnce sync.Once
)

// GetStoragePath returns the path to tasks.json (migrated to ~/.focus)
// Migration is performed only once, on first call.
func GetStoragePath() string {
	// Run migration exactly once
	migrationOnce.Do(func() {
		home, err := os.UserHomeDir()
		if err != nil {
			cachedStoragePath = "./tasks.json"
			return
		}

		// Attempt migration (best-effort, errors ignored)
		_ = migrateStorage()

		cachedStoragePath = filepath.Join(home, ".focus", "tasks.json")
	})

	return cachedStoragePath
}

// migrateStorage performs a best-effort migration:
// - creates ~/.focus/
// - backups existing ~/.focus/ into ~/.focus/backup-pre-migration-<ts> if not empty
// - copies contents from ~/.neon/ into ~/.focus/
// - copies important files from ~/.todo/ (e.g., ai_cache.json) into ~/.focus/
// Errors are returned but migration attempts to continue where reasonable.
func migrateStorage() error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	neonPath := filepath.Join(home, ".neon")
	todoPath := filepath.Join(home, ".todo")
	focusPath := filepath.Join(home, ".focus")

	// Ensure focus dir exists
	if err := os.MkdirAll(focusPath, 0o755); err != nil {
		return fmt.Errorf("create focus dir: %w", err)
	}

	var firstErr error

	// If focus already has content, create a timestamped backup copy
	hasFiles := false
	ents, err := os.ReadDir(focusPath)
	if err == nil && len(ents) > 0 {
		hasFiles = true
	}
	if hasFiles {
		backupPath := filepath.Join(focusPath, fmt.Sprintf("backup-pre-migration-%d", time.Now().Unix()))
		if err := copyDir(focusPath, backupPath); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("backup existing focus dir: %w", err)
		}
	}

	// Migrate from .neon if present
	if info, err := os.Stat(neonPath); err == nil && info.IsDir() {
		if err := copyDir(neonPath, focusPath); err != nil && firstErr == nil {
			firstErr = fmt.Errorf("migrate .neon -> .focus: %w", err)
		}
	}

	// Merge ai cache or specific files from .todo
	if info, err := os.Stat(todoPath); err == nil && info.IsDir() {
		// copy ai_cache.json if exists
		src := filepath.Join(todoPath, "ai_cache.json")
		if _, err := os.Stat(src); err == nil {
			dst := filepath.Join(focusPath, "ai_cache.json")
			if err := copyFile(src, dst); err != nil && firstErr == nil {
				firstErr = fmt.Errorf("migrate ai_cache.json: %w", err)
			}
		}
		// Optionally copy other useful files (safe, non-recursive)
	}

	return firstErr
}

// copyDir recursively copies src directory into dst (creates dst).
func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			// continue walking where possible
			return err
		}
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}
		dstPath := filepath.Join(dst, rel)

		if info.IsDir() {
			return os.MkdirAll(dstPath, info.Mode().Perm())
		}
		// Copy file
		return copyFile(path, dstPath)
	})
}

// copyFile copies a single file from src to dst, preserving mode.
func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	// Ensure destination directory exists
	if err := os.MkdirAll(filepath.Dir(dst), 0o755); err != nil {
		return err
	}

	out, err := os.OpenFile(dst, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	defer func() {
		_ = out.Close()
	}()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}
	// Try to copy mode from source
	if fi, err := os.Stat(src); err == nil {
		_ = os.Chmod(dst, fi.Mode())
	}
	return nil
}
