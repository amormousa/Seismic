package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Config holds all the settings that seismic needs, loaded at startup from the env vars
type Config struct {
	DatabaseURL string
	JWTSecret   string
	Port        string
}

// Load reads the .env file (if exist) and returns a Config struct
// with all values filled in. It exit immediately if a required value is missing
func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, reading from environment directly")
	}

	cfg := &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		Port:        os.Getenv("PORT"),
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("DATABASE_URL is not set in .env")
	}

	if cfg.DatabaseURL == "" {
		log.Fatal("JWT_SECRET is not set in .env")
	}

	if cfg.DatabaseURL == "" {
		cfg.Port = "5024"
	}

	return cfg
}
