package db

import (
	"database/sql"
	"fmt"
	"os"
	"time"
)

const (
	TableImages   = "images"
	TableFiles    = "files"
	TableCommands = "commands"
	TableKeylogs  = "keylogs"
)

func Init(schemaFilePath string, dbPath string) (*sql.DB, error) {
	firstInstall := false

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		firstInstall = true
	}

	// Open or create DB
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	// Read schema from file
	schema, err := os.ReadFile(schemaFilePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read schema.sql: %v", err)
	}

	// Execute schema SQL
	_, err = db.Exec(string(schema))
	if err != nil {
		return nil, fmt.Errorf("failed to execute schema: %v", err)
	}

	if firstInstall {
		_, err := db.Exec("INSERT INTO metadata (install_time) VALUES (?)", time.Now())
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}
