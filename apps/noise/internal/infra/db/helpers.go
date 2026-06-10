package db

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/Kyanite/noise/internal/constants"
	"github.com/Kyanite/noise/internal/domain"
	appErrors "github.com/Kyanite/noise/internal/errors"
	errutil "github.com/Kyanite/noise/internal/errutil"
	"github.com/Kyanite/noise/internal/logging"
)

// isRetryableError determines if a database error is retryable
func (db *DB) isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	errStr := strings.ToLower(err.Error())

	// Retryable errors
	retryablePatterns := []string{
		"timeout",
		"connection reset",
		"connection refused",
		"temporary failure",
		"server closed",
		"database is locked",
		"database locked",
		"deadlock",
		"lock wait timeout",
	}

	for _, pattern := range retryablePatterns {
		if strings.Contains(errStr, pattern) {
			return true
		}
	}

	return false
}

// validateConnection checks if the database connection is healthy
func (db *DB) validateConnection() error {
	if db.conn == nil {
		return appErrors.NewDatabaseError("connection_nil", fmt.Errorf("database connection is nil")).WithComponent("repository")
	}

	ctx, cancel := context.WithTimeout(context.Background(), constants.DBValidationTimeout)
	defer cancel()

	if err := db.conn.PingContext(ctx); err != nil {
		return appErrors.NewDatabaseError("ping_failed", err).WithComponent("repository")
	}

	return nil
}

// beginTransaction starts a database transaction with error handling
func (db *DB) beginTransaction(ctx context.Context) (*sql.Tx, error) {
	if err := db.validateConnection(); err != nil {
		return nil, err
	}

	tx, err := db.conn.BeginTx(ctx, nil)
	if err != nil {
		return nil, appErrors.NewDatabaseError("begin_transaction", err).WithOperation("BeginTransaction").WithComponent("repository")
	}

	return tx, nil
}

// rollbackTransaction rolls back a transaction with error handling
func (db *DB) rollbackTransaction(tx *sql.Tx, operation string) {
	if tx == nil {
		return
	}

	if err := tx.Rollback(); err != nil {
		logging.GetDefaultLogger().Error("Failed to rollback transaction", "operation", operation, "error", err)
	}
}

// commitTransaction commits a transaction with error handling
func (db *DB) commitTransaction(tx *sql.Tx, operation string) error {
	if tx == nil {
		return appErrors.NewDatabaseError("commit_nil_tx", fmt.Errorf("transaction is nil")).WithOperation(operation).WithComponent("repository")
	}

	if err := tx.Commit(); err != nil {
		return appErrors.NewDatabaseError("commit_transaction", err).WithOperation(operation).WithComponent("repository")
	}

	return nil
}

// handleDatabaseError creates appropriate error types based on database error conditions
func (db *DB) handleDatabaseError(operation string, err error) error {
	if err == nil {
		return nil
	}

	dbErr := appErrors.NewDatabaseError(operation, err).WithComponent("repository")

	// Categorize specific database errors
	errStr := strings.ToLower(err.Error())

	switch {
	case strings.Contains(errStr, "no such table") || strings.Contains(errStr, "doesn't exist"):
		return appErrors.NewDatabaseError("table_not_found", err).WithOperation(operation).WithComponent("repository")
	case strings.Contains(errStr, "unique constraint") || strings.Contains(errStr, "duplicate entry"):
		return appErrors.NewDatabaseError("duplicate_entry", err).WithOperation(operation).WithComponent("repository")
	case strings.Contains(errStr, "foreign key constraint"):
		return appErrors.NewDatabaseError("foreign_key_violation", err).WithOperation(operation).WithComponent("repository")
	case strings.Contains(errStr, "check constraint"):
		return appErrors.NewDatabaseError("check_constraint", err).WithOperation(operation).WithComponent("repository")
	case strings.Contains(errStr, "syntax error") || strings.Contains(errStr, "near"):
		return appErrors.NewDatabaseError("sql_syntax_error", err).WithOperation(operation).WithComponent("repository")
	case err == sql.ErrNoRows:
		return appErrors.NewDatabaseError("no_rows", err).WithOperation(operation).WithComponent("repository")
	case strings.Contains(errStr, "locked"):
		return appErrors.NewDatabaseError("database_locked", err).WithOperation(operation).WithComponent("repository")
	default:
		return dbErr
	}
}

// getProjectInTx retrieves a project within a transaction
func (db *DB) getProjectInTx(tx *sql.Tx, id int) (*domain.Project, error) {
	query := `
		SELECT id, name, description, song_ids, created_at, updated_at
		FROM projects WHERE id = ?`

	row := tx.QueryRow(query, id)

	var project domain.Project
	var songIDsJSON string

	err := row.Scan(
		&project.ID,
		&project.Name,
		&project.Description,
		&songIDsJSON,
		&project.CreatedAt,
		&project.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("project with ID %d not found", id)
		}
		return nil, errutil.Wrap(err, "get project")
	}

	// Unmarshal song IDs
	project.SongIDs, err = unmarshalIntArray(songIDsJSON)
	if err != nil {
		return nil, errutil.Wrap(err, "unmarshal song IDs")
	}

	return &project, nil
}

// updateProjectInTx updates a project within a transaction
func (db *DB) updateProjectInTx(tx *sql.Tx, project *domain.Project) error {
	songIDsJSON, err := marshalIntArray(project.SongIDs)
	if err != nil {
		return errutil.Wrap(err, "marshal song IDs")
	}

	query := `
		UPDATE projects
		SET name = ?, description = ?, song_ids = ?, updated_at = ?
		WHERE id = ?`

	_, err = tx.Exec(query,
		project.Name,
		project.Description,
		songIDsJSON,
		time.Now(),
		project.ID,
	)
	if err != nil {
		return errutil.Wrap(err, "update project")
	}

	return nil
}
