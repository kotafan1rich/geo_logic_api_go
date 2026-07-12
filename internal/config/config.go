package config

import (
	"log/slog"
	"time"

	"github.com/caarlos0/env/v9"
	"github.com/joho/godotenv"
	"github.com/kotafan1rich/geo_logic_api_go/internal/logger"
)

type Database struct {
	Host         string `env:"DB_HOST" envDefault:"localhost"`
	Port         string `env:"DB_PORT" envDefault:"5432"`
	User         string `env:"DB_USER" envDefault:"postgres"`
	Password     string `env:"DB_PASSWORD" envDefault:"postgres"`
	Name         string `env:"DB_NAME" envDefault:"postgres"`
	SSLMode      string `env:"DB_SSL_MODE" envDefault:"disable"`
	MaxIdleConns int    `env:"MAX_IDLE_CONNS" envDefault:"10"`
	MaxOpenConns int    `env:"MAX_OPEN_CONNS" envDefault:"100"`
}

type HttpServer struct {
	ServerPort string `env:"SERVER_PORT" envDefault:"8080"`
	GinMode    string `env:"GIN_MODE" envDefault:"debug"`

	ReadTimeout  time.Duration `env:"READ_TIMEOUT" envDefault:"5s"`
	WriteTimeout time.Duration `env:"WRITE_TIMEOUT" envDefault:"10s"`
	IdleTimeout  time.Duration `env:"IDLE_TIMEOUT" envDefault:"60s"`
}

type Logging struct {
	LogLevel     logger.LogLevel  `env:"LOG_LEVEL" envDefault:"info"`
	LogFormat    logger.LogFormat `env:"LOG_FORMAT" envDefault:"text"` // json or text
	LogAddSource bool             `env:"LOG_ADD_SOURCE" envDefault:"false"`
}

type Config struct {
	Database   Database
	HttpServer HttpServer
	Logging    Logging
}

var config *Config

func MustLoad() {
	err := godotenv.Load()
	if err != nil {
		slog.Warn("No .env file found, using environment variables")
	}

	cfg := &Config{}

	err = env.Parse(cfg)
	if err != nil {
		slog.Error("error load config", "err", err)
		panic(err)
	}

	config = cfg
}

func Get() *Config {
	if config == nil {
		panic("config is not initialized! call config.MustLoad() first")
	}
	return config
}

func (d *Database) DSN() string {
	return "host=" + d.Host +
		" port=" + d.Port +
		" user=" + d.User +
		" password=" + d.Password +
		" dbname=" + d.Name +
		" sslmode=" + d.SSLMode
}
