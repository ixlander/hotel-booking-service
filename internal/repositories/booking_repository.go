package repositories

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ixlander/hotel-booking-service/internal/data"
)

type PostgresBookingRepository struct {
	db *sql.DB
}

func NewPostgresBookingRepository(db *sql.DB) *PostgresBookingRepository {
	return &PostgresBookingRepository{db: db}
}

func (r *PostgresBookingRepository) Create(ctx context.Context, booking *data.Booking) (*data.Booking, error) {
	query := `
		INSERT INTO bookings (user_id, room_id, from_date, to_date, status)
		VALUES ($1, $2, $3, $4, 'active')
		RETURNING id, created_at
	`
	
	err := r.db.QueryRowContext(
		ctx, 
		query, 
		booking.UserID, 
		booking.RoomID, 
		booking.FromDate, 
		booking.ToDate,
	).Scan(
		&booking.ID,
		&booking.CreatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	booking.Status = "active"
	return booking, nil
}

func (r *PostgresBookingRepository) FindByID(ctx context.Context, id int64) (*data.Booking, error) {
	query := `
		SELECT id, user_id, room_id, from_date, to_date, created_at, status 
		FROM bookings 
		WHERE id = $1
	`
	
	var booking data.Booking
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&booking.ID,
		&booking.UserID,
		&booking.RoomID,
		&booking.FromDate,
		&booking.ToDate,
		&booking.CreatedAt,
		&booking.Status,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Booking not found
		}
		return nil, err
	}
	
	return &booking, nil
}

func (r *PostgresBookingRepository) FindByUserID(ctx context.Context, userID int64) ([]data.Booking, error) {
	query := `
		SELECT id, user_id, room_id, from_date, to_date, created_at, status 
		FROM bookings 
		WHERE user_id = $1
	`
	
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var bookings []data.Booking
	for rows.Next() {
		var booking data.Booking
		if err := rows.Scan(
			&booking.ID,
			&booking.UserID,
			&booking.RoomID,
			&booking.FromDate,
			&booking.ToDate,
			&booking.CreatedAt,
			&booking.Status,
		); err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return bookings, nil
}

func (r *PostgresBookingRepository) UpdateStatus(ctx context.Context, id int64, status string) error {
	query := `UPDATE bookings SET status = $1 WHERE id = $2`
	
	_, err := r.db.ExecContext(ctx, query, status, id)
	return err
}