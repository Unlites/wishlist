package config

import (
	"fmt"
	"time"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Postgres struct {
		Host     string `env:"POSTGRES_HOST"`
		Port     string `env:"POSTGRES_PORT"`
		Database string `env:"POSTGRES_DB"`
		Username string `env:"POSTGRES_USER"`
		Password string `env:"POSTGRES_PASSWORD"`
		DSN      string
	}
	JWT struct {
		SecretKey      string        `env:"JWT_SECRET_KEY"`
		ExpirationTime time.Duration `env:"JWT_EXPIRATION_TIME"`
	}
	HttpServer struct {
		Address      string        `env:"HTTP_SERVER_ADDRESS"`
		ReadTimeout  time.Duration `env:"HTTP_READ_TIMEOUT"`
		WriteTimeout time.Duration `env:"HTTP_WRITE_TIMEOUT"`
		IdleTimeout  time.Duration `env:"HTTP_IDLE_TIMEOUT"`
	}
	HashCost        int           `env:"HASH_COST"`
	ShutdownTimeout time.Duration `env:"SHUTDOWN_TIMEOUT"`
	LogLevel        string        `env:"LOG_LEVEL"`
}

func New() *Config {
	var cfg Config

	if err := cleanenv.ReadEnv(&cfg); err != nil {
		panic("failed to read config: " + err.Error())
	}

	cfg.Postgres.DSN = fmt.Sprintf("postgres://%s:%s@%s:%s/%s",
		cfg.Postgres.Username,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Database,
	)

	return &cfg
}
