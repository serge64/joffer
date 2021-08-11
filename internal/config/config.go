package config

import (
	"os"
	"strconv"
)

type Config struct {
	HostAddr      string
	LogLevel      string
	DatabaseURL   string
	SessionAddr   string
	SessionKey    string
	SessionMaxAge int
	ClientID      string
	ClientSecret  string
	Token         string
}

func New() *Config {
	maxAge, err := strconv.Atoi(getEnv("SESSION_MAX_AGE", "604800"))
	if err != nil {
		maxAge = 604800
	}

	return &Config{
		HostAddr:      getEnv("HOST_ADDR", ":8080"),
		LogLevel:      getEnv("LOG_LEVEL", "debug"),
		DatabaseURL:   getEnv("DATABASE_URL", ""),
		SessionAddr:   getEnv("SESSION_ADDR", ":6379"),
		SessionKey:    getEnv("SESSION_KEY", "secret"),
		SessionMaxAge: maxAge,
		ClientID:      getEnv("CLIENT_ID", ""),
		ClientSecret:  getEnv("CLIENT_SECRET", ""),
		Token:         getEnv("TOKEN", ""),
	}
}

func getEnv(key, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return defaultValue
}
