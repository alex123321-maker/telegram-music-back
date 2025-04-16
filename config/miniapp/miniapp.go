package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        int
	DatabaseURL string
}

func Load() (*Config, error) {
	_ = godotenv.Load()

	portStr := os.Getenv("APP_PORT")
	if portStr == "" {
		portStr = "3000"
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("некорректный порт: %v", err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		dbURL = "postgres://postgres:postgres@localhost:5430/telegram_music?sslmode=disable"
	}

	return &Config{
		Port:        port,
		DatabaseURL: dbURL,
	}, nil
}
