package config

import (
	"os"

	"github.com/Khaym03/kumo/internal/pkg/utils"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Browsers             []BrowserConfig `yaml:"browsers"`
	NumOfPagesPerBrowser int             `yaml:"numOfPagesPerBrowser"`
	RetryCount           int             `yaml:"retryCount"`
	StorageDir           string          `yaml:"storageDir"`
	FingerprintsPath string `yaml:"fingerprintsPath"`
	AllowBadgerLogger    bool            `yaml:"allowBadgerLogger"`
	RateLimit RateLimitConfig `yaml:"rate_limit"`
}

type BrowserConfig struct {
	Type    string       `yaml:"type" json:"type"`
	Name    string       `yaml:"name" json:"name"`
	Proxy   bool         `yaml:"proxy,omitempty" json:"proxy,omitempty"`
	Address string       `yaml:"address,omitempty" json:"address,omitempty"`
	Pages   []PageConfig `yaml:"pages,omitempty" json:"pages"`
}

type PageConfig struct {
	Id string `yaml:"id" json:"id"`
}

type RateLimitConfig struct {
	Delay int  `yaml:"delay" json:"delay"`
	MaxBurst int  `yaml:"max_burst" json:"max_burst"`
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

func (c *Config) String() string {
	return utils.ToJSONString(c)
}
