package config

import (
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	IP   string `yaml:"ip"`
	Port int    `yaml:"port"`
}

func ReadCfg() (Config, error) {
	file, err := os.ReadFile("config.yml")
	var cfg Config
	if err != nil {
		return cfg, err
	}
	err = yaml.Unmarshal(file, &cfg)
	if err != nil {
		return cfg, err
	}
	return cfg, nil
}
