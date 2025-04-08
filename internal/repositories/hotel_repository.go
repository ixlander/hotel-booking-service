package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ixlander/hotel-booking-service/internal/data"
)

type PostgresHotelRepository struct {
	db *sql.DB
}

func NewPostgresHotelRepository(db *sql.DB) *PostgresHotelRepository {
	return &PostgresHotelRepository{db: db}
}

func (r *PostgresHotelRepository) GetAll(ctx context.Context) ([]data.Hotel, error) {
	query := `SELECT id, name, city FROM hotels`
	
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var hotels []data.Hotel
	for rows.Next() {
		var hotel data.Hotel
		if err := rows.Scan(&hotel.ID, &hotel.Name, &hotel.City); err != nil {
			return nil, err
		}
		hotels = append(hotels, hotel)
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return hotels, nil
}

func (r *PostgresHotelRepository) FindByID(ctx context.Context, id int64) (*data.Hotel, error) {
	query := `SELECT id, name, city FROM hotels WHERE id = $1`
	
	var hotel data.Hotel
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&hotel.ID,
		&hotel.Name,
		&hotel.City,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Hotel not found
		}
		return nil, err
	}
	
	return &hotel, nil
}