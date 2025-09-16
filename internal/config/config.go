package config

import (
	"log"
	"os"
	"strconv"
)

// Config holds the application configuration.
type Config struct {
	ServerPort  int
	DatabaseURL string
}

// Load loads configuration from environment variables.
func Load() *Config {
	portStr := os.Getenv("SERVER_PORT")
	if portStr == "" {
		portStr = "8080" // Default port
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Fatalf("Invalid server port: %v", err)
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Println("DATABASE_URL not set, using default value for development")
		dbURL = "postgres://user:password@localhost:5432/urlshortener?sslmode=disable"
	}

	return &Config{
		ServerPort:  port,
		DatabaseURL: dbURL,
	}
}
