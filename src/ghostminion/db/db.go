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

	db, err := GetDB(dbPath)
	if err != nil {
		return nil, err
	}

	if err = loadSchema(db, schemaFilePath); err != nil {
		return nil, err
	}

	if firstInstall {
		_, err := db.Exec("INSERT INTO metadata (install_time) VALUES (?)", time.Now())
		if err != nil {
			return nil, err
		}
	}

	return db, nil
}

func loadSchema(db *sql.DB, schemaFilePath string) error {
	schema, err := os.ReadFile(schemaFilePath)
	if err != nil {
		return fmt.Errorf("failed to read schema.sql: %v", err)
	}
	_, err = db.Exec(string(schema))
	if err != nil {
		return fmt.Errorf("failed to execute schema: %v", err)
	}
	return nil
}

func GetDB(dbPath string) (*sql.DB, error) {
	return sql.Open("sqlite3", dbPath)
}
