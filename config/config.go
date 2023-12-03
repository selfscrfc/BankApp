package config

import (
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"os"
	"time"
)

// App config struct
type Config struct {
	GRPC     GRPCConfig
	Server   ServerConfig
	Postgres PostgresConfig

	Cookie  Cookie
	Store   Store
	Session Session
	Logger  Logger
}

type ServerConfig struct {
	AppVersion        string
	Port              string
	PprofPort         string
	Mode              string
	JwtSecretKey      string
	CookieName        string
	ReadTimeout       time.Duration
	WriteTimeout      time.Duration
	SSL               bool
	CtxDefaultTimeout time.Duration
	CSRF              bool
	Debug             bool
}

type GRPCConfig struct {
	Port    int
	Timeout time.Duration
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

// Postgresql config
type PostgresConfig struct {
	PostgresqlHost     string
	PostgresqlPort     string
	PostgresqlUser     string
	PostgresqlPassword string
	PostgresqlDbname   string
	PostgresqlSSLMode  bool
	PgDriver           string
}

// Cookie config
type Cookie struct {
	Name     string
	MaxAge   int
	Secure   bool
	HTTPOnly bool
}

// Session config
type Session struct {
	Prefix string
	Name   string
	Expire int
}

// Metrics config
type Metrics struct {
	URL         string
	ServiceName string
}

// Store config
type Store struct {
	ImagesFolder string
}

// AWS S3
type Jaeger struct {
	Host        string
	ServiceName string
	LogSpans    bool
}
