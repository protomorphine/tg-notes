// Package config provides config    types and methods to configure application
package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Environment   string           `yaml:"environment" env-default:"local"`
	Logger        LoggerConfig     `yaml:"logger"`
	Bot           BotConfig        `yaml:"bot"`
	HTTPServer    HTTPServerConfig `yaml:"httpServer"`
	GitRepository GitRepository    `yaml:"gitRepository"`
}

type BotConfig struct {
	Key           string        `env:"TG_API_KEY" env-required:"true"`
	InitTimeout   time.Duration `yaml:"initTimeout" env-default:"1m"`
	WebHookURL    string        `yaml:"webHookURL" env-required:"true"`
	AllowedUserID int64         `yaml:"allowedUserID"`
}

type LoggerConfig struct {
	MinLevel string `yaml:"minLevel" env-default:"INFO"`
}

type HTTPServerConfig struct {
	Addr string `yaml:"addr" env-default:":80"`
}

type GitRepository struct {
	URL         string    `yaml:"url" env-required:"true"`
	Path        string    `yaml:"path"`
	KeyPath     string    `yaml:"keyPath"`
	KeyPassword string    `env:"KEY_PASSWD"`
	PathToSave  string    `yaml:"saveTo" env-required:"true"`
	Branch      string    `yaml:"branch"`
	Committer   Committer `yaml:"committer"`
}

type Committer struct {
	Name string `yaml:"name"`
}

func Load(configPath string) (*Config, error) {
	var config Config

	err := cleanenv.ReadConfig(configPath, &config)
	return &config, err
}
