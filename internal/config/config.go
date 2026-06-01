package config

import (
	"fmt"
	"os"
	"strings"
)

type Config struct {
	Port         string
	DatabaseURL  string
	BaseURL      string
}

func Load() (Config, error) {
	port := envOrDefault("PORT", "8080")
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return Config{}, fmt.Errorf("DATABASE_URL is required")
	}

	baseURL := os.Getenv("BASE_URL")
	if baseURL == "" {
		return Config{}, fmt.Errorf("BASE_URL is required")
	}
	baseURL = strings.TrimRight(baseURL, "/")

	return Config{
		Port:        port,
		DatabaseURL: databaseURL,
		BaseURL:     baseURL,
	}, nil
}

func envOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
