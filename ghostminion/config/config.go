package config

import (
	"fmt"
	"gopkg.in/yaml.v3"
	"os"
	"path/filepath"
	"sync"
)

type Config struct {
	AgentID        string              `yaml:"AgentID"`
	Installation   InstallationConfig  `yaml:"Installation"`
	Communication  CommunicationConfig `yaml:"Communication"`
	Apps           AppsConfig          `yaml:"Apps"`
	mu             sync.Mutex
	configFilePath string
}

var (
	instance *Config
	once     sync.Once
)

func NewConfig(configFilePath string) *Config {
	config := &Config{
		mu:             sync.Mutex{},
		configFilePath: configFilePath,
	}
	once.Do(func() {
		data, readError := os.ReadFile(configFilePath)
		if readError != nil {
			panic(fmt.Errorf("failed to read config file: %w", readError))
		}
		if yamlError := yaml.Unmarshal(data, config); yamlError != nil {
			panic(fmt.Errorf("failed to parse YAML: %w", yamlError))
		}
	})
	return config
}

func GetInstance() *Config {
	once.Do(func() {
		configFilePath := GetConfigFilePath()
		instance = NewConfig(configFilePath)
	})
	return instance
}

func (cfg *Config) Save() error {
	data, err := yaml.Marshal(instance)
	if err != nil {
		return err
	}
	return os.WriteFile(cfg.configFilePath, data, 0600)
}

func (cfg *Config) Update(updateFn func(c *Config)) error {
	cfg.mu.Lock()
	defer cfg.mu.Unlock()
	updateFn(instance)
	return cfg.Save()
}

func GetConfigFilePath() string {
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	return filepath.Join(currentDir, "config.yaml")
}

func DeleteConfig() error {
	return os.Remove(GetConfigFilePath())
}
