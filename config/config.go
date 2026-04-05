package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL string
	ServerPort  string
	PasetoKey   string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found. Using system variables.")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		log.Fatal("DATABASE_URL environment variable is missing.")
	}

	pasetoKey := os.Getenv("PASETO_KEY")
	if len(pasetoKey) != 32 {
		log.Fatal("FATAL: PASETO_KEY must be exactly 32 characters long.")
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	return &Config{
		DatabaseURL: dbURL,
		ServerPort:  port,
		PasetoKey:   pasetoKey,
	}
}