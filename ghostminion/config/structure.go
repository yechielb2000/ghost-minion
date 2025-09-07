package config

import "encoding/json"

type InstallationConfig struct {
	DBPath     string `yaml:"DBPath"`
	DBPassword string `yaml:"DBPassword"`
	AESKey     string `yaml:"AESKey"`
}

type ServerConfig struct {
	Address string `yaml:"Address"`
	Port    int    `yaml:"Port"`
	Key     string `yaml:"Key"`
}

type CommunicationConfig struct {
	Interval    string         `yaml:"Interval"`
	Servers     []ServerConfig `yaml:"Servers"`
	Certificate string         `yaml:"Certificate"`
}

type HiderConfig struct {
	NewProcessName string `yaml:"NewProcessName"`
}

// AppDataConfig is the same as AppData
type AppDataConfig struct {
	Id     string          `json:"id"`
	Type   string          `json:"type"`
	Name   string          `json:"name"`
	Params json.RawMessage `json:"params"`
}
