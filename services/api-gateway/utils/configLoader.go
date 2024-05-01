package utils

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type ServerConfig struct {
	IP                     string
	Port                   string
	JWTSecret              string
	AuthServiceURL         string
	NewsletterServiceURL   string
	SubscriptionServiceURL string
	PublishingServiceURL   string
}

func LoadConfig(envPath string) ServerConfig {
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Error loading .env file from %s: %v", envPath, err)
	}

	return ServerConfig{
		IP:                     getEnv("SERVER_IP", "0.0.0.0"),
		Port:                   getEnv("SERVER_PORT", "8080"),
		JWTSecret:              getEnv("JWT_SECRET", "abc"),
		AuthServiceURL:         getEnv("AUTH_SERVICE_URL", "localhost"),
		NewsletterServiceURL:   getEnv("NEWSLETTER_SERVICE_URL", "localhost"),
		SubscriptionServiceURL: getEnv("SUBSCRIPTION_SERVICE_URL", "localhost"),
		PublishingServiceURL:   getEnv("PUBLISHING_SERVICE_URL", "localhost"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
