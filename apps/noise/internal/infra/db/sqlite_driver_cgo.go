//go:build cgo

package db

import (
	// Import modernc.org/sqlite for CGO builds (no CGO required, provides better security and cross-platform support)
	_ "modernc.org/sqlite"
)

const sqliteDriverName = "sqlite"
