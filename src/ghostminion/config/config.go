package config

import (
	"fmt"
	"ghostminion/apps"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
)

type InstallationConfig struct {
	RunPassword string
	DBPath      string
	ConfigFile  string
}

type ServerConfig struct {
	Address string
	Key     int
}

type CommunicationConfig struct {
	Interval string
	Servers  []ServerConfig
}

type AppsConfig struct {
	Keylogger     apps.KeyLoggerApp
	Screenshot    apps.ScreenshotApp
	SecurityGuard apps.SecurityGuardApp
}

type Config struct {
	Installation  InstallationConfig
	Communication CommunicationConfig
	Apps          AppsConfig
}

var (
	configInstance *Config
	once           sync.Once
)

func LoadConfig(path string) (*Config, error) {
	var loadError error

	once.Do(func() {
		data, readError := os.ReadFile(path)
		if readError != nil {
			loadError = fmt.Errorf("failed to read config file: %w", readError)
			return
		}
		configInstance = &Config{}
		if yamlError := yaml.Unmarshal(data, configInstance); yamlError != nil {
			loadError = fmt.Errorf("failed to parse YAML: %w", yamlError)
			return
		}
	})

	return configInstance, loadError
}

func GetConfig() *Config {
	if configInstance == nil {
		fmt.Println("Config not initialized. Call LoadConfig first.")
	}
	return configInstance
}
