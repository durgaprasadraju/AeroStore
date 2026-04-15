package config

import (
	"os"
)

type Config struct {
	GRPCPort    string
	PostgresURL string
	RedisURL    string
}

func LoadConfig() *Config {
	return &Config{
		GRPCPort:    getEnv("GRPC_PORT", "50051"),
		PostgresURL: getEnv("POSTGRES_URL", "postgres://admin:password123@localhost:5432/aerostore_db?sslmode=disable"),
		RedisURL:    getEnv("REDIS_URL", "localhost:6379"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
