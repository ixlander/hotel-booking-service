package connections

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/ixlander/hotel-booking-service/internal/app/config"
)

func NewPostgresDB(config config.DatabaseConfig) (*sql.DB, error) {
	connectionString := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
		config.Host, config.Port, config.User, config.Password, config.DBName, config.SSLMode,
	)

	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}