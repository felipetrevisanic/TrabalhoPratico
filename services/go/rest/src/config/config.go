package config

import (
	"fmt"
	"os"
)

type Config struct {
	Host        string
	Port        string
	Database    string
	User        string
	Password    string
	SSLMode     string
	HTTPAddress string
}

func Load() Config {
	return Config{
		Host:        getEnv("DB_HOST", "localhost"),
		Port:        getEnv("DB_PORT", "5432"),
		Database:    getEnv("DB_NAME", "tcc_banco"),
		User:        getEnv("DB_USER", "postgres"),
		Password:    getEnv("DB_PASSWORD", "postgres"),
		SSLMode:     getEnv("DB_SSLMODE", "disable"),
		HTTPAddress: getEnv("HTTP_ADDRESS", ":8080"),
	}
}

func (c Config) DatabaseURL() string {
	return fmt.Sprintf(
		"host=%s port=%s dbname=%s user=%s password=%s sslmode=%s",
		c.Host,
		c.Port,
		c.Database,
		c.User,
		c.Password,
		c.SSLMode,
	)
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}

	return value
}
