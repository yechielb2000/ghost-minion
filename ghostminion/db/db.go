package db

import (
	"database/sql"
	"fmt"
	"ghostminion/config"
	"ghostminion/cryptography"
	"ghostminion/db/dbDataTypes"
	_ "modernc.org/sqlite"
	"os"
	"sync"
)

type FetchTableConfig struct {
	Name      string
	BatchSize int
}

type AgentDB struct {
	session    *sql.DB
	dbPath     string
	dbPassword string
}

var (
	dbInstance *AgentDB
	once       sync.Once
)

func GetInstance() *AgentDB {
	once.Do(func() {
		cfg := config.GetInstance()
		dbInstance = NewAgentDB(cfg.Installation.DBPath, cfg.Installation.DBPassword)
	})
	return dbInstance
}

func NewAgentDB(dbPath string, dbPassword string) *AgentDB {
	db := &AgentDB{
		session:    nil,
		dbPath:     dbPath,
		dbPassword: dbPassword,
	}
	once.Do(func() {
		if _, err := os.Stat(dbPath); err != nil {
			if err = db.loadSchema(); err != nil {
				return
			}
		}
	})

	if err := db.startSession(); err != nil {
		return nil
	}
	return db
}

func (db *AgentDB) startSession() error {
	connStr := fmt.Sprintf("%s?_pragma_key=%s", db.dbPath, db.dbPassword)
	if session, err := sql.Open("sqlite3", connStr); err != nil {
		return err
	} else {
		db.session = session
	}
	return nil
}

func (db *AgentDB) loadSchema() error {
	tx, err := db.session.Begin()
	if err != nil {
		return err
	}
	_, err = tx.Exec(schemas)
	if err != nil {
		return err
	}
	return tx.Commit()
}

func (db *AgentDB) FetchRows(cfg FetchTableConfig) ([]map[string]interface{}, error) {
	query := fmt.Sprintf(
		"SELECT * FROM %s ORDER BY save_time ASC LIMIT ?",
		cfg.Name,
	)

	rows, err := db.session.Query(query, cfg.BatchSize)
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
			val := vals[i]

			if (col == "data" || col == "message") && val != nil {
				if cipherBytes, ok := val.([]byte); ok {
					if plain, err := cryptography.DecryptData(cipherBytes); err == nil {
						val = plain
					} else {
						return nil, err
					}
				}
			}

			rowMap[col] = val

			if col == "request_id" {
				requestID, _ = val.(string)
			}
		}

		if requestID != "" {
			if err := db.RemoveRow(cfg.Name, requestID); err != nil {
				return nil, err
			}
		}

		results = append(results, rowMap)
	}

	return results, nil
}

func (db *AgentDB) RemoveRow(table string, requestID string) error {
	_, err := db.session.Exec(fmt.Sprintf("DELETE FROM %s WHERE request_id = ?", table), requestID)
	return err
}

func (db *AgentDB) WriteData(requestID string, dataType dbDataTypes.DataType, data []byte) error {
	data, err := cryptography.EncryptData(data)
	if err != nil {
		return err
	}
	query := "INSERT INTO data (request_id, data, data_type) VALUES (?, ?, ?)"
	_, err = db.session.Exec(query, requestID, data, string(dataType))
	return err
}

func (db *AgentDB) WriteLog(level string, message []byte) error {
	message, err := cryptography.EncryptData(message)
	if err != nil {
		return err
	}
	_, err = db.session.Exec("INSERT INTO logs (message, level) VALUES (?, ?)", message, level)
	return err
}
