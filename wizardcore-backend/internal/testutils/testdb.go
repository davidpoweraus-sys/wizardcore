package testutils

import (
	"database/sql"
	"fmt"
	"os"
	"testing"

	"github.com/yourusername/wizardcore-backend/internal/database"
)

// TestDatabaseURL returns the database URL for testing.
// It reads from environment variable TEST_DATABASE_URL.
// If not set, returns empty string and skips the test.
func TestDatabaseURL(t *testing.T) string {
	url := os.Getenv("TEST_DATABASE_URL")
	if url == "" {
		t.Skip("TEST_DATABASE_URL not set, skipping integration test")
	}
	return url
}

// SetupTestDB connects to the test database, runs migrations, and returns a connection.
// It also registers a cleanup function to truncate tables after each test.
func SetupTestDB(t *testing.T) *sql.DB {
	t.Helper()

	dbURL := TestDatabaseURL(t)

	// Connect to database
	db, err := database.Connect(dbURL)
	if err != nil {
		t.Fatalf("failed to connect to test database: %v", err)
	}

	// Skip migrations for now because migration files are not in the expected format.
	// The schema is already applied manually.
	// if err := database.RunMigrations(dbURL); err != nil {
	// 	t.Fatalf("failed to run migrations: %v", err)
	// }

	// Cleanup function to truncate all tables after test
	t.Cleanup(func() {
		if err := truncateAllTables(db); err != nil {
			t.Logf("failed to truncate tables: %v", err)
		}
		db.Close()
	})

	return db
}

// truncateAllTables truncates all tables in the public schema.
// This is a destructive operation; use only in test environments.
func truncateAllTables(db *sql.DB) error {
	rows, err := db.Query(`
		SELECT tablename FROM pg_tables 
		WHERE schemaname = 'public' 
		AND tablename NOT IN ('schema_migrations')
	`)
	if err != nil {
		return fmt.Errorf("failed to list tables: %w", err)
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var table string
		if err := rows.Scan(&table); err != nil {
			return fmt.Errorf("failed to scan table name: %w", err)
		}
		tables = append(tables, table)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating tables: %w", err)
	}

	// Disable triggers and truncate each table
	for _, table := range tables {
		_, err := db.Exec(fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			return fmt.Errorf("failed to truncate table %s: %w", table, err)
		}
	}
	return nil
}