// Package config provides config    types and methods to configure application
package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Environment string           `yaml:"environment" env-default:"local"`
	Logger      LoggerConfig     `yaml:"logger"`
	Bot         BotConfig        `yaml:"bot"`
	HTTPServer  HTTPServerConfig `yaml:"httpServer"`
}

type BotConfig struct {
	Key         string        `env:"TG_API_KEY" env-required:"true"`
	InitTimeout time.Duration `yaml:"initTimeout" env-default:"1m"`
	WebHookURL  string        `yaml:"webHookURL" env-required:"true"`
}

type LoggerConfig struct {
	MinLevel string `yaml:"minLevel" env-default:"INFO"`
}

type HTTPServerConfig struct {
	Addr string `yaml:"addr" env-default:":80"`
}

func Load(configPath string) (*Config, error) {
	var config Config

	err := cleanenv.ReadConfig(configPath, &config)
	return &config, err
}
