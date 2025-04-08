package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"

	"github.com/ixlander/hotel-booking-service/internal/app/connections"
)

type Config struct {
	Server struct {
		Port string
	}
	
	Database connections.PostgresConfig
	
	JWT struct {
		Secret string
		TTL    time.Duration
	}
}

func LoadConfig() (*Config, error) {
	_ = godotenv.Load()
	
	config := &Config{}
	
	config.Server.Port = getEnv("SERVER_PORT", "8080")
	
	config.Database = connections.PostgresConfig{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "postgres"),
		DBName:   getEnv("DB_NAME", "hotel_booking"),
		SSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
	
	// JWT config
	config.JWT.Secret = getEnv("JWT_SECRET", "your-secret-key")
	jwtTTLHours, _ := strconv.Atoi(getEnv("JWT_TTL_HOURS", "24"))
	config.JWT.TTL = time.Duration(jwtTTLHours) * time.Hour
	
	return config, nil
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}