package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port         string
	DatastoreURL string
}

func init() {
	godotenv.Load()
}

func LoadEnv() *Config {
	return &Config{
		Port:         os.Getenv("PORT"),
		DatastoreURL: os.Getenv("DATASTORE_URL"),
	}
}
