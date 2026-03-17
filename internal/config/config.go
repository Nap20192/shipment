package config

import (
	"fmt"
	"os"
	"strconv"
)

type Config struct {
	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	LogLevel   string
	LogDir     string
	GRPCPort   int
}

func LoadConfig() (*Config, error) {
	return &Config{
		DBHost:     getEnv("POSTGRES_HOST", "localhost"),
		DBPort:     getEnvInt("POSTGRES_PORT", 5432),
		DBUser:     getEnv("POSTGRES_USER", "postgres"),
		DBPassword: getEnv("POSTGRES_PASSWORD", "postgres"),
		DBName:     getEnv("POSTGRES_DB", "shipment_db"),
		LogLevel:   getEnv("LOG_LEVEL", "INFO"),
		LogDir:     getEnv("LOG_DIR", "./logs"),
		GRPCPort:   getEnvInt("GRPC_PORT", 50051),
	}, nil
}

func (c *Config) DBConnString() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		c.DBUser, c.DBPassword, c.DBHost, c.DBPort, c.DBName)
}

func (c *Config) GrpcPortString() string {
	return fmt.Sprintf(":%d", c.GRPCPort)
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if i, err := strconv.Atoi(value); err == nil {
			return i
		}
	}
	return fallback
}
