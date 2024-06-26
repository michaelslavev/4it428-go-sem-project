package utils

import (
	"github.com/joho/godotenv"
	"log"
	"os"
)

type ServerConfig struct {
	IP          string
	Port        string
	SupabaseURL string
	SupabaseKEY string
	DatabaseURL string
}

func LoadConfig(envPath string) ServerConfig {
	if err := godotenv.Load(envPath); err != nil {
		log.Printf("Error loading .env file from %s: %v", envPath, err)
	}

	return ServerConfig{
		IP:          getEnv("SERVER_IP", "0.0.0.0"),
		Port:        getEnv("SERVER_PORT", "8080"),
		SupabaseURL: getEnv("SUPABASE_URL", "supabase.com"),
		SupabaseKEY: getEnv("SUPABASE_KEY", "skrrr"),
		DatabaseURL: getEnv("DATABASE_URL", "skrrr"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
