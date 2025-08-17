package communication

import (
	"encoding/json"
	"ghostminion/config"
	"time"
)

type Telemetry struct {
	AgentID          string `yaml:"AgentID"`
	FirstTime        bool   `yaml:"FirstTime"`
	AgentType        string `yaml:"AgentName"`
	IsAlive          bool   `yaml:"isAlive"`
	CurrentTimestamp int64  `yaml:"CurrentTimestamp"`
}

func NewTelemetry(firstTime bool, isAlive bool) (Telemetry, error) {
	c, err := config.GetConfig()
	if err != nil {
		return Telemetry{}, err
	}

	return Telemetry{
		AgentID:          c.AgentID,
		FirstTime:        firstTime,
		AgentType:        "GhostMinion",
		IsAlive:          isAlive,
		CurrentTimestamp: time.Now().Unix(),
	}, nil
}

func SendTelemetry(serverConfig config.ServerConfig, telemetry Telemetry) ([]byte, int, error) {
	jsonTelemetry, err := json.Marshal(telemetry)
	if err != nil {
		return []byte{}, 0, err
	}
	route := CreateRoute(serverConfig, "receive")
	return SendRequest(POST,
		route,
		map[string]string{
			"Content-Type": "application/json",
		}, jsonTelemetry)
}
