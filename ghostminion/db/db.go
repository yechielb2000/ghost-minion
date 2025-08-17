package db

import (
	"database/sql"
	"fmt"
	"ghostminion/cryptography"
	"ghostminion/db/dbDataTypes"
	_ "modernc.org/sqlite"
	"os"
	"time"
)

type TableConfig struct {
	Name      string
	BatchSize int
}

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

func FetchRows(cfg TableConfig) ([]map[string]interface{}, error) {
	query := fmt.Sprintf(
		"SELECT * FROM %s ORDER BY save_time ASC LIMIT ?",
		cfg.Name,
	)

	rows, err := dbInstance.Query(query, cfg.BatchSize)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var results []map[string]interface{}
	for rows.Next() {
		vals := make([]interface{}, len(cols))
		ptrs := make([]interface{}, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}

		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}

		rowMap := make(map[string]interface{})
		var requestID string
		for i, col := range cols {
			rowMap[col] = vals[i]
			if col == "request_id" {
				requestID, _ = vals[i].(string)
			}
		}

		if requestID != "" {
			if err := RemoveRow(cfg.Name, requestID); err != nil {
				return nil, err
			}
		}

		results = append(results, rowMap)
	}

	return results, nil
}

func RemoveRow(table string, requestID string) error {
	_, err := dbInstance.Exec(fmt.Sprintf("DELETE FROM %s WHERE request_id = ?", table), requestID)
	return err
}

func WriteData(requestID string, dataType dbDataTypes.DataType, data []byte) error {
	data, err := cryptography.EncryptData(data)
	if err != nil {
		return err
	}
	query := "INSERT INTO data (request_id, data, data_type) VALUES (?, ?, ?)"
	_, err = dbInstance.Exec(query, requestID, data, dataType)
	return err
}

func WriteLog(level string, message []byte) error {
	message, err := cryptography.EncryptData(message)
	if err != nil {
		return err
	}
	_, err = dbInstance.Exec("INSERT INTO logs (message, level) VALUES (?, ?)", message, level)
	return err
}
