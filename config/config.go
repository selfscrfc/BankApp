package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"time"
)

// App config struct
type Config struct {
	GRPC   GRPCConfig
	Server ServerConfig
	Logger Logger
}

type ServerConfig struct {
	AppVersion string
	Port       string
	Mode       string
}

type GRPCConfig struct {
	CustomerServicePort int
	AccountsServicePort int
	Timeout             time.Duration
}

func LoadConfig() (*Config, error) {
	path := os.Getenv("CONFIG_PATH")
	var cfgName string
	if path == "" {
		path = "../"
		cfgName = "./config/config-local"
	} else {
		cfgName = "./config/config-docker"
	}

	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, err
	}

	v := viper.New()

	v.SetConfigName(cfgName)
	v.AddConfigPath(path)
	v.AutomaticEnv()
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return nil, errors.New("file not found")
		}
		return nil, err
	}

	var cfg *Config

	err := v.Unmarshal(&cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// Server config struct

// Logger config
type Logger struct {
	Development       bool
	DisableCaller     bool
	DisableStacktrace bool
	Encoding          string
	Level             string
}
