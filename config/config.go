package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBHost                string
	DBName                string
	DBUsername            string
	DBPassword            string
	DBPort                string
	JWTSecret             string
	GoogleClientID        string
	GoogleClientSecret    string
	FileUploadSizeLimitMB string
	S3BucketName          string
	S3Endpoint            string
	S3AccessKeyId         string
	S3SecretAccessKey     string
	S3UseSSL              bool
	StripeSecretKey       string
	StripeWebhookSecret   string
}

func NewConfig() (config Config) {
	if err := godotenv.Load(); err != nil {
		log.Println("ℹ️ No .env file found, using system environment variables.")
	}

	config = Config{
		DBHost:                getEnv("DB_HOST"),
		DBName:                getEnv("DB_DATABASE"),
		DBUsername:            getEnv("DB_USERNAME"),
		DBPassword:            getEnv("DB_PASSWORD"),
		DBPort:                getEnv("DB_PORT"),
		JWTSecret:             getEnv("JWT_SECRET"),
		GoogleClientID:        getEnv("GOOGLE_CLIENT_ID"),
		GoogleClientSecret:    getEnv("GOOGLE_CLIENT_SECRET"),
		FileUploadSizeLimitMB: getEnv("FILE_UPLOAD_SIZE_LIMIT_MB"),
		S3BucketName:          getEnv("S3_BUCKET_NAME"),
		S3Endpoint:            getEnv("S3_ENDPOINT"),
		S3AccessKeyId:         getEnv("S3_ACCESS_KEY_ID"),
		S3SecretAccessKey:     getEnv("S3_SECRET_ACCESS_KEY"),
		S3UseSSL:              true,
		StripeSecretKey:       getEnv("STRIPE_SECRET_KEY"),
		StripeWebhookSecret:   getEnv("STRIPE_WEBHOOK_SECRET"),
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
