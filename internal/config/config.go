// Package config provides config    types and methods to configure application
package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Environment string       `yaml:"environment" env-default:"local"`
	Logger      LoggerConfig `yaml:"logger"`
	Bot         BotConfig    `yaml:"bot"`
}

type BotConfig struct {
	Key     string        `yaml:"key" env-required:"true"`
	Timeout time.Duration `yaml:"timeout" env-default:"1m"`
}

type LoggerConfig struct {
	MinLevel string `yaml:"minLevel" env-default:"INFO"`
}

func Load(configPath string) (*Config, error) {
	var config Config
	if err := cleanenv.ReadConfig(configPath, &config); err != nil {
		return nil, err
	}

	return &config, nil
}
