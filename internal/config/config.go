package config

import "os"

type Config struct {
	HostAddr    string
	LogLevel    string
	DatabaseURL string
}

func New() *Config {
	return &Config{
		HostAddr:    getEnv("HOST_ADDR", ":8080"),
		LogLevel:    getEnv("LOG_LEVEL", "debug"),
		DatabaseURL: getEnv("DATABASE_URL", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
