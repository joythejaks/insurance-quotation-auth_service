package config

import "os"

type Config struct {
	AppPort    string
	DBHost     string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	DBSSLMode  string
	JWTSecret  string
}

func LoadConfig() *Config {
	return &Config{
		AppPort:    getEnv("APP_PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBName:     getEnv("DB_NAME", "auth_db"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),
		JWTSecret:  getEnv("JWT_SECRET", "super-secret-key"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
