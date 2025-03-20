package db

import (
	"database/sql"
	"fmt"
)

func removeOneRow(db *sql.DB, table string, requestID string) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE request_id = ?", table)
	_, err := db.Exec(query, requestID)
	return err
}
