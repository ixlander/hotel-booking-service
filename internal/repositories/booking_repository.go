package postgres

import (
	"context"
	"database/sql"
	"errors"

	"github.com/ixlander/hotel-booking-service/internal/data"
)

type BookingRepo struct {
	db *sql.DB
}

func NewBookingRepo(db *sql.DB) *BookingRepo {
	return &BookingRepo{db: db}
}

func (r *BookingRepo) Create(ctx context.Context, booking *data.Booking) error {
	query := `
		INSERT INTO bookings (user_id, room_id, from_date, to_date, status)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at
	`
	
	if booking.Status == "" {
		booking.Status = "confirmed"
	}
	
	err := r.db.QueryRowContext(
		ctx, query, 
		booking.UserID, booking.RoomID, booking.FromDate, booking.ToDate, booking.Status,
	).Scan(&booking.ID, &booking.CreatedAt)
	
	if err != nil {
		return err
	}
	
	return nil
}

func (r *BookingRepo) FindByID(ctx context.Context, id int64) (*data.Booking, error) {
	query := `
		SELECT id, user_id, room_id, from_date, to_date, created_at, status
		FROM bookings WHERE id = $1
	`
	
	var booking data.Booking
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&booking.ID, &booking.UserID, &booking.RoomID,
		&booking.FromDate, &booking.ToDate, &booking.CreatedAt, &booking.Status,
	)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Booking not found
		}
		return nil, err
	}
	
	return &booking, nil
}

func (r *BookingRepo) FindByUserID(ctx context.Context, userID int64) ([]*data.Booking, error) {
	query := `
		SELECT id, user_id, room_id, from_date, to_date, created_at, status
		FROM bookings WHERE user_id = $1
	`
	
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var bookings []*data.Booking
	for rows.Next() {
		var booking data.Booking
		if err := rows.Scan(
			&booking.ID, &booking.UserID, &booking.RoomID,
			&booking.FromDate, &booking.ToDate, &booking.CreatedAt, &booking.Status,
		); err != nil {
			return nil, err
		}
		bookings = append(bookings, &booking)
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return bookings, nil
}

func (r *BookingRepo) UpdateStatus(ctx context.Context, id int64, status string) error {
	query := `UPDATE bookings SET status = $1 WHERE id = $2`
	
	_, err := r.db.ExecContext(ctx, query, status, id)
	if err != nil {
		return err
	}
	
	return nil
}