package db

import (
	"database/sql"
	"fmt"
	_ "modernc.org/sqlite"
	"os"
	"time"
)

const (
	TableImages   = "images"
	TableFiles    = "files"
	TableCommands = "commands"
	TableKeylogs  = "keylogs"
)

const dbSchemaFilePath = "./db/schema.sql"

func Init(dbPath string) error {
	firstInstall := false

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		firstInstall = true
	}

	db, err := GetDB(dbPath)
	if err != nil {
		return err
	}

	if err = loadSchema(db); err != nil {
		return err
	}

	if firstInstall {
		_, err = db.Exec("INSERT INTO metadata (install_time) VALUES (?)", time.Now())
		if err != nil {
			return err
		}
	}

	return nil
}

func GetDB(dbPath string) (*sql.DB, error) {
	return sql.Open("sqlite", dbPath)
}

func loadSchema(db *sql.DB) error {
	schema, err := os.ReadFile(dbSchemaFilePath)
	if err != nil {
		return fmt.Errorf("failed to read schema.sql: %v", err)
	}
	_, err = db.Exec(string(schema))
	if err != nil {
		return fmt.Errorf("failed to execute schema: %v", err)
	}
	return nil
}
