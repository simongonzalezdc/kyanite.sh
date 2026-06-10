//go:build cgo

package db

import (
	// Import SQLite driver for CGO builds
	_ "github.com/mattn/go-sqlite3"
)

const sqliteDriverName = "sqlite3"
