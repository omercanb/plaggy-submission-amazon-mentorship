package db

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
)

func InitDB(dbPath string) (*sql.DB, error) {
	dir := filepath.Dir(dbPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open DB: %w", err)
	}

	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to communicate with DB: %w", err)
	}
	return db, nil
}
