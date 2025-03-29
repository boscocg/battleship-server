package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnvironment() string {
	environment := os.Getenv("ENV")
	if environment == "" {
		environment = "local"
	}
	return environment
}

func LoadEnv() {
	envFile := ".env"
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Error loading %s file %v", envFile, err)
	}
	log.Printf("%s environment variables loaded", envFile)
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
