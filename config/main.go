package config

import (
	"os"

	"gopkg.in/yaml.v2"
)

// Config is a parsed version of the configuration file.
type Config struct {
	BaseURL       string `yaml:"base_url"`
	Port          int    `yaml:"port"`
	DatabaseURL   string `yaml:"database_url"`
	SessionSecret string `yaml:"session_secret"`
	DataPath      string `yaml:"data_path"`
}

// LoadConfig reads the file at the given path and parses it as a YAML
// configuration file.
func LoadConfig(configPath string) (*Config, error) {
	config := &Config{}

	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	d := yaml.NewDecoder(file)
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}
