package config

import "encoding/json"

type InstallationConfig struct {
	DBPath     string `yaml:"db_path"`
	DBPassword string `yaml:"db_password"`
	AESKey     string `yaml:"aes_key"`
}

type ServerConfig struct {
	Address string `yaml:"address"`
	Port    int    `yaml:"port"`
	Key     string `yaml:"key"`
}

type CommunicationConfig struct {
	Interval    string         `yaml:"interval"`
	Servers     []ServerConfig `yaml:"servers"`
	Certificate string         `yaml:"certificate"`
}

type HidingConfig struct {
	ProcessName string `yaml:"process_name"`
}

// AppDataConfig is the same as AppData
type AppDataConfig struct {
	Id     string          `json:"id"`
	Type   string          `json:"type"`
	Name   string          `json:"name"`
	Params json.RawMessage `json:"params"`
}
