package communication

import "ghostminion/config"

type Telemetry struct {
	AgentID   string `yaml:"AgentID"`
	FirstTime bool   `yaml:"FirstTime"`
	AgentType string `yaml:"AgentName"`
	IsAlive   bool   `yaml:"isAlive"`
}

func NewTelemetry(firstTime, isAlive bool) (Telemetry, error) {
	c, err := config.GetConfig()
	if err != nil {
		return Telemetry{}, err
	}

	return Telemetry{
		AgentID:   c.AgentID,
		FirstTime: firstTime,
		AgentType: "GhostMinion",
		IsAlive:   isAlive,
	}, nil
}
