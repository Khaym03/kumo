package config

import (
	"database/sql"
	"log"
	"os"

	db "github.com/Khaym03/kumo/db/sqlite/gen"
	"github.com/Khaym03/kumo/ports"
	"github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	Queries    *db.Queries
	Logger     *logrus.Logger
	TaskStatus ports.TaskStatus
	DB         *sql.DB
}

func NewAppConfig(
	DB *sql.DB,
	queries *db.Queries,
	ts ports.TaskStatus,
) AppConfig {
	return AppConfig{
		Queries:    queries,
		TaskStatus: ts,
		DB:         DB,
	}
}

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
