package communication

import (
	"encoding/json"
	"ghostminion/config"
	"time"
)

type Telemetry struct {
	AgentID          string `yaml:"AgentID"`
	AgentType        string `yaml:"AgentName"`
	IsAlive          bool   `yaml:"isAlive"`
	CurrentTimestamp int64  `yaml:"CurrentTimestamp"`
}

func NewTelemetry(isAlive bool) (Telemetry, error) {
	return Telemetry{
		AgentID:          config.GetInstance().AgentID,
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
	return SendRequest(
		POST,
		CreateRoute(serverConfig, "receive"),
		map[string]string{
			"Content-Type": "application/json",
		},
		jsonTelemetry,
	)
}
