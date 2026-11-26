package config

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Postgres struct {
		Host     string
		Port     string
		User     string
		Password string
		DBName   string
	}
	JWT struct {
		Secret          string
		AccessDuration  time.Duration
		RefreshDuration time.Duration
	}
}

func NewConfig() *Config {
	envFiles := []string{
		".env",
		"../../.env",
	}

	for _, file := range envFiles {
		if err := godotenv.Load(file); err != nil {
			if errors.Is(err, os.ErrNotExist) {
				continue
			}

			log.Fatalf("Ошибка загрузки %s: %v", file, err)
		}
	}

	cfg := &Config{}
	cfg.Postgres.Host = os.Getenv("POSTGRES_HOST")
	cfg.Postgres.Port = os.Getenv("POSTGRES_PORT")
	cfg.Postgres.User = os.Getenv("POSTGRES_USER")
	cfg.Postgres.Password = os.Getenv("POSTGRES_PASSWORD")
	cfg.Postgres.DBName = os.Getenv("POSTGRES_DB")

	cfg.JWT.Secret = os.Getenv("JWT_SECRET")
	accessDuration, err := time.ParseDuration(os.Getenv("JWT_ACCESS_TOKEN_DURATION"))
	if err != nil {
		log.Fatal("invalid JWT_ACCESS_TOKEN_DURATION")
	}
	cfg.JWT.AccessDuration = accessDuration

	refreshDuration, err := time.ParseDuration(os.Getenv("JWT_REFRESH_TOKEN_DURATION"))
	if err != nil {
		log.Fatal("invalid JWT_REFRESH_TOKEN_DURATION")
	}
	cfg.JWT.RefreshDuration = refreshDuration
	return cfg
}
