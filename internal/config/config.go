package config

import (
	"os"
	"time"

	"github.com/goccy/go-yaml"
	"gopkg.in/yaml.v3"
)

// Config — корневая структура конфигурации.
type Config struct {
	IntervalMinutes int `yaml:"interval_minutes"`
	IncludeKeywords []string `yaml:"include_keywords"`
	ExcludeKeywords []string `yaml:"exclude_keywords"`
	Kwork KworkConfig `yaml:"kwork"`
	Server ServerConfig `yaml:"server"`
}

type KworkConfig struct {
	Login string `yaml:"login"`
	Password string `yaml:"password"`
	PhoneLast string `yaml:"phone_last"`
	ProxyURL string `yaml:"proxy_url"`
}

type ServerConfig struct {
	Addr string `yaml:"addr"`
}

// Load читает YAML-файл и возвращает Config.
// Если interval_minutes не задан, по умолчанию используется 180 минут (3 часа).
func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config 
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err 
	}
	if cfg.IntervalMinutes <= 0 {
		cfg.IntervalMinutes = 100 
	}

	return &cfg, nil
}

// Interval возвращает длительность интервала в виде time.Duration.
func (c *Config) Interval() time.Duration {
	return time.Duration(c.IntervalMinutes) * time.Minute 
}