package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost             string
	DBName             string
	DBUsername         string
	DBPassword         string
	DBPort             string
	JWTSecret          string
	GoogleClientID     string
	GoogleClientSecret string
}

func NewConfig() (config Config) {
	if err := godotenv.Load(); err != nil {
		log.Println("ℹ️ No .env file found, using system environment variables.")
	}

	config = Config{
		DBHost:             getEnv("DB_HOST"),
		DBName:             getEnv("DB_DATABASE"),
		DBUsername:         getEnv("DB_USERNAME"),
		DBPassword:         getEnv("DB_PASSWORD"),
		DBPort:             getEnv("DB_PORT"),
		JWTSecret:          getEnv("JWT_SECRET"),
		GoogleClientID:     getEnv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret: getEnv("GOOGLE_CLIENT_SECRET"),
	}
	return
}

func getEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		log.Fatalf("❌ Missing required environment variable: %s", key)
	}
	return value
}
