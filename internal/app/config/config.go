package config

import (
	"os"
	"strconv"
	"time"
	
	"github.com/joho/godotenv"
	
	"hotel-booking-service/internal/app/connections"
)

type Config struct {
	Server   ServerConfig
	Database connections.PostgresConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string
}

type JWTConfig struct {
	Secret       string
	TokenExpiry  time.Duration
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		if _, ok := os.LookupEnv("APP_ENV"); !ok {
			return nil, err
		}
	}
	
	dbPort, _ := strconv.Atoi(getEnv("DB_PORT", "5432"))
	
	tokenExpiryHours, _ := strconv.Atoi(getEnv("JWT_TOKEN_EXPIRY_HOURS", "24"))
	
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: connections.PostgresConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     dbPort,
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "hotel_booking"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:      getEnv("JWT_SECRET", "your_secret_key"),
			TokenExpiry: time.Duration(tokenExpiryHours) * time.Hour,
		},
	}, nil
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}