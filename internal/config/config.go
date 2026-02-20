// Package config provides config types and methods to configure application
package config

import (
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

// Config represents the application's configuration.
type Config struct {
	Environment   string           `yaml:"environment" env-default:"prod"` // current environment
	Logger        LoggerConfig     `yaml:"logger"`                         // logger configuration
	Bot           BotConfig        `yaml:"bot"`                            // Telegram bot configuration
	HTTPServer    HTTPServerConfig `yaml:"httpServer"`                     // HTTP server configuration
	GitRepository GitRepository    `yaml:"gitRepository"`                  // git repository configuration
}

// BotConfig represents the Telegram bot's configuration.
type BotConfig struct {
	Key           string        `env:"TG_API_KEY" env-required:"true"`                    // bot API key
	InitTimeout   time.Duration `yaml:"initTimeout" env-default:"1m"`                     // bot init timeout
	WebHookURL    string        `yaml:"webHookURL" env:"WEBHOOK_URL" env-required:"true"` // URL where Telegram will send updates
	AllowedUserID int64         `yaml:"allowedUserID"`                                    // user ID, which allowed to perform actions
}

// LoggerConfig represents the logger's configuration.
type LoggerConfig struct {
	MinLevel string `yaml:"minLevel" env-default:"INFO"` // minimal log level
}

// HTTPServerConfig represents the HTTP server's configuration.
type HTTPServerConfig struct {
	Addr string `yaml:"addr" env-default:":80"` // address to bind
}

// GitRepository represents the Git repository's configuration.
type GitRepository struct {
	URL             string        `yaml:"url" env-required:"true"`            // remote repo URL
	Path            string        `yaml:"path"`                               // local path to clone repo
	Auth            GitAuth       `yaml:"auth"`                               // git authentication config
	PathToSave      string        `yaml:"saveTo" env-required:"true"`         // path to save notes inside repo
	Branch          string        `yaml:"branch"`                             // repo working branch
	RemoteName      string        `yaml:"remoteName" env-default:"origin"`    // git remote name
	Committer       Committer     `yaml:"committer"`                          // committer info
	BufSize         int           `yaml:"bufSize" env-required:"true"`        // notes buffer size
	UpdateDuratiion time.Duration `yaml:"updateDuration" env-required:"true"` // duration to fill buffer; save occurs when buffer is full or last save was specified time ago
}

// GitAuth represents the Git authentication configuration.
type GitAuth struct {
	Key         string `env:"KEY"`        // ssh key to access repo
	KeyPassword string `env:"KEY_PASSWD"` // password to ssh key
}

// Committer represents the committer's information.
type Committer struct {
	Name string `yaml:"name"` // commiter name
}

// Load loads the configuration from the given path.
func Load(configPath string) (*Config, error) {
	var config Config

	err := cleanenv.ReadConfig(configPath, &config)
	return &config, err
}
