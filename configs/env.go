package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func GetEnvironment() string {
	environment := os.Getenv("ENV")
	if environment == "" {
		environment = "prod"
	}
	return environment
}

func LoadEnv() {
	environment := GetEnvironment()

	envFile := ".env." + environment

	log.Printf("Loading %s environment variables", envFile)
	err := godotenv.Load(envFile)
	if err != nil {
		log.Fatalf("Error loading %s file %v", envFile, err)
	}
	log.Printf("%s environment variables loaded", envFile)
}

func GetEnv(key string) string {
	return os.Getenv(key)
}
