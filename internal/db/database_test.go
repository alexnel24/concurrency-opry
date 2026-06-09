package db

import (
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	_ "modernc.org/sqlite"
)

func TestOpenDB_ReturnsConnectedDB(t *testing.T) {
	// :memory: required for SQLite; hardcoded in openDB() -> sql.Open("sqlite", dbPath)
	t.Setenv("DATABASE_PATH", ":memory:")

	db, err := openDB()
	require.NoError(t, err)
	require.NotNil(t, db)
	t.Cleanup(func() { db.Close() })

	assert.Equal(t, 1, db.Stats().MaxOpenConnections)
	assert.NoError(t, db.Ping())
}

func TestOpenDB_RespectsEnvVar(t *testing.T) {
	path := filepath.Join(t.TempDir(), "test.db")
	t.Setenv("DATABASE_PATH", path)

	db, err := openDB()
	require.NoError(t, err)
	require.NotNil(t, db)
	t.Cleanup(func() { db.Close() })

	assert.NoError(t, db.Ping())
	_, err = os.Stat(path)
	assert.NoError(t, err, "expected database file to exist at env var path")
}

func TestApplySchemas_CreatesExpectedTables(t *testing.T) {
	// :memory: is SQLite-specific; used here to avoid writing a file to disk during tests
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	applySchemas(db)

	for _, table := range []string{"events", "artists", "performances"} {
		var count int
		err := db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?`, table).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count, "expected table %q to exist", table)
	}
}

func TestApplySchemas_IsIdempotent(t *testing.T) {
	// :memory: is SQLite-specific; used here to avoid writing a file to disk during tests
	db, err := sql.Open("sqlite", ":memory:")
	require.NoError(t, err)
	t.Cleanup(func() { db.Close() })

	applySchemas(db)
	// second call verifies IF NOT EXISTS makes schema application a safe no-op on repeat runs
	applySchemas(db)

	assert.NoError(t, db.Ping())
}

func TestInitDB_ReturnsReadyDB(t *testing.T) {
	t.Setenv("DATABASE_PATH", ":memory:")

	db, err := InitDB()
	require.NoError(t, err)
	require.NotNil(t, db)
	t.Cleanup(func() { db.Close() })

	for _, table := range []string{"events", "artists", "performances"} {
		var count int
		err := db.QueryRow(`SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name=?`, table).Scan(&count)
		require.NoError(t, err)
		assert.Equal(t, 1, count, "expected table %q to exist", table)
	}
}
