package config

import (
	"errors"
	"fmt"
	"ghostminion/logger"
	"gopkg.in/yaml.v3"
	"os"
	"sync"
)

type InstallationConfig struct {
	DBPath      string `yaml:"DBPath"`
	LogFilePath string `yaml:"LogFilePath"`
	DBPassword  string `yaml:"DBPassword"`
	AESKey      string `yaml:"AESKey"`
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
	lgr      = logger.GetLogger()
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

	if loadError != nil {
		lgr.Error("Failed to load config", loadError)
	}

	return instance, loadError
}

func GetConfig() (*Config, error) {
	if instance == nil {
		return nil, errors.New("config not initialized. Call LoadConfig first")
	}
	return instance, nil
}

func SaveConfig(path string) error {
	data, err := yaml.Marshal(instance)
	if err != nil {
		lgr.Error("Got error while marshalling config", err)
		return err
	}
	return os.WriteFile(path, data, 0600)
}

func UpdateConfig(savePath string, updateFn func(c *Config)) error {
	mu.Lock()
	defer mu.Unlock()
	updateFn(instance)
	return SaveConfig(savePath)
}
