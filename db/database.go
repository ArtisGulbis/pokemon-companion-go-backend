package db

import (
	"database/sql"
	_ "embed"
	"fmt"
	"os"

	_ "github.com/glebarez/go-sqlite"
)

//go:embed schema2.sql
var schemaSQL string

type Database struct {
	*sql.DB
}

func New(dbPath string) (*Database, error) {
	dbExists := fileExists(dbPath)

	sqlDB, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	// CRITICAL: For in-memory databases, all operations must use the same connection
	if dbPath == ":memory:" {
		sqlDB.SetMaxOpenConns(1)
	}

	// Enable foreign keys (important for SQLite!)
	if _, err := sqlDB.Exec("PRAGMA foreign_keys = ON"); err != nil {
		sqlDB.Close()
		return nil, fmt.Errorf("failed to enable foreign keys: %w", err)
	}

	db := &Database{sqlDB}

	if !dbExists {
		if err := db.initSchema(); err != nil {
			db.Close()
			return nil, err
		}
	}

	return db, nil
}

func (db *Database) initSchema() error {
	_, err := db.Exec(schemaSQL)
	if err != nil {
		return fmt.Errorf("schema initialization failed: %w", err)
	}
	fmt.Println("Schema initialized successfully")
	return nil
}

func (db *Database) Reset() error {
	fmt.Println("🔄 Resetting database...")

	_, err := db.Exec(schemaSQL)
	if err != nil {
		return fmt.Errorf("schema reset failed: %w", err)
	}

	fmt.Println("✅ Database reset complete")
	return nil
}

func fileExists(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}
