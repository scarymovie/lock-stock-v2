package postgres

import (
	"os"
	"strconv"
)

type Config struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
	SSLMode  string
}

func GetPostgresConfig() Config {
	return Config{
		Host:     getEnv("POSTGRES_HOST", "lock-stock-v2-postgres-development"),
		Port:     getEnvAsInt("POSTGRES_PORT", 5432),
		User:     getEnv("POSTGRES_USER", "db_user"),
		Password: getEnv("POSTGRES_PASSWORD", "db_password"),
		DBName:   getEnv("POSTGRES_DB", "db_database"),
		SSLMode:  getEnv("POSTGRES_SSLMODE", "disable"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return fallback
}
