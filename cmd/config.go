package cmd

import (
	"fmt"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

// Global koanf instance. Use . as the key path delimiter. This can be / or anything.
var (
	koanfConfig = koanf.New(".")
	config      *Config
)

type Config struct {
	Logger struct {
		Level    string `yaml:"port"`
		Encoding string `yaml:"encoding"`
		Color    bool   `yaml:"color"`
		Output   string `yaml:"output"`
	} `yaml:"logger"`

	Metrics struct {
		Enabled bool   `yaml:"enabled"`
		Host    string `yaml:"host"`
		Port    int    `yaml:"port"`
	} `yaml:"metrics"`

	Profiler struct {
		Enabled bool   `yaml:"enabled"`
		Pidfile string `yaml:"pidfile"`
	} `yaml:"profiler"`

	Server struct {
		Host     string `yaml:"host"`
		Port     string `yaml:"port"`
		Tls      bool   `yaml:"tls"`
		Devcert  bool   `yaml:"devcert"`
		Certfile string `yaml:"certfile"`
		Keyfile  string `yaml:"keyfile"`

		Log struct {
			Enabled      bool     `yaml:"enabled"`
			Level        string   `yaml:"level"`
			RequestBody  bool     `yaml:"request_body"`
			ResponseBody bool     `yaml:"response_body"`
			IgnorePaths  []string `yaml:"ignore_paths"`
		} `yaml:"log"`

		CORS struct {
			Enabled          bool     `yaml:"enabled"`
			AllowedOrigins   []string `yaml:"allowed_origins"`
			AllowedMethods   []string `yaml:"allowed_methods"`
			AllowedHeaders   []string `yaml:"allowed_headers"`
			AllowCredentials bool     `yaml:"allow_credentials"`
			MaxAge           int      `yaml:"max_age"`
		} `yaml:"cors"`

		Metrics struct {
			Enabled     bool     `yaml:"enabled"`
			IgnorePaths []string `yaml:"ignore_paths"`
		} `yaml:"metrics"`
	} `yaml:"server"`
}

func Load() (*koanf.Koanf, *Config, error) {
	if config != nil {
		return koanfConfig, config, nil
	}

	koanfConfig.Load(file.Provider("configs/defaults.yaml"), yaml.Parser())
	if koanfConfig.Raw() == nil || len(koanfConfig.Raw()) == 0 {
		return nil, nil, fmt.Errorf("could not load config: %s", "defaults.yaml")
	}

	config = &Config{}
	if err := koanfConfig.Unmarshal("", config); err != nil {
		return nil, nil, err
	}

	return koanfConfig, config, nil
}
