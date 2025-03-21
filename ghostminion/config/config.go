package config

import (
	"fmt"
	"ghostminion/apps"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
)

type InstallationConfig struct {
	DBPath     string `yaml:"DBPath"`
	ConfigFile string `yaml:"ConfigFile"`
	DBPassword string `yaml:"DBPassword"`
}

type ServerConfig struct {
	Address string `yaml:"Address"`
	Key     string `yaml:"Key"`
}

type CommunicationConfig struct {
	Interval string         `yaml:"Interval"`
	Servers  []ServerConfig `yaml:"Servers"`
}

type AppsConfig struct {
	Keylogger     *apps.KeyLoggerApp     `yaml:"Keylogger,omitempty"`
	Screenshot    *apps.ScreenshotApp    `yaml:"Screenshot,omitempty"`
	SecurityGuard *apps.SecurityGuardApp `yaml:"SecurityGuard,omitempty"`
}

type Config struct {
	Installation  InstallationConfig  `yaml:"Installation"`
	Communication CommunicationConfig `yaml:"Communication"`
	Apps          AppsConfig          `yaml:"Apps"`
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
