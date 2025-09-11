package config

import (
	"os"

	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Browsers []BrowserConfig `yaml:"browsers"`
}

type BrowserConfig struct {
	Type  string `yaml:"type"`
	Name  string `yaml:"name"`
	Proxy bool   `yaml:"proxy,omitempty"`
	Host  string `yaml:"host,omitempty"`
}

func LoadKumoConfig() *Config {
	path := "kumo-config.yaml"
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read config file %s: %v", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		log.Fatalf("Failed to parse YAML: %v", err)
	}

	return &cfg
}
