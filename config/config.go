package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DSN  string
	Host string
	Level        string
}

func New() (*Config, error) {
	err := godotenv.Load()
	if err != nil {
		log.Print("No .env file found")
		log.Print(err)
		return nil, err
	}

	return &Config{
		DSN:          getEnv("DSN"),
		Host:         getEnv("HOST"),
		Level:        getEnv("LEVEL"),
	}, nil
}

func getEnv(key string) string {
	val, exists := os.LookupEnv(key)
	if exists {
		return val
	}
	return ""
}
