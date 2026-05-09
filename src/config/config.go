package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DBURL           string
	DBHost          string
	DBPort          string
	DBUser          string
	DBPassword      string
	DBName          string
	DBSSLMode       string
	JWTSecret       string
	XenditKey       string
	XenditWebToken  string
	Port            string
}

func Load() *Config {
	godotenv.Load()

	return &Config{
		DBURL:           os.Getenv("DATABASE_URL"),
		DBHost:          getEnv("DB_HOST", "localhost"),
		DBPort:          getEnv("DB_PORT", "5432"),
		DBUser:          getEnv("DB_USER", "postgres"),
		DBPassword:      getEnv("DB_PASSWORD", "postgres"),
		DBName:          getEnv("DB_NAME", "go_commerce"),
		DBSSLMode:       getEnv("DB_SSLMODE", "disable"),
		JWTSecret:       getEnv("JWT_SECRET", "secret"),
		XenditKey:       os.Getenv("XENDIT_SECRET_KEY"),
		XenditWebToken:  os.Getenv("XENDIT_WEBHOOK_TOKEN"),
		Port:            getEnv("PORT", "8080"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
