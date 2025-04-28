package repositories

import (
	"database/sql"
	"time"
	
	"hotel-booking-service/internal/data"
)

type BookingRepository struct {
	db *sql.DB
}

func NewBookingRepository(db *sql.DB) *BookingRepository {
	return &BookingRepository{db: db}
}

func (r *BookingRepository) CreateBooking(userID, roomID int, fromDate, toDate time.Time) (*data.Booking, error) {
	query := `
		INSERT INTO bookings (user_id, room_id, from_date, to_date, status)
		VALUES ($1, $2, $3, $4, 'confirmed')
		RETURNING id, user_id, room_id, from_date, to_date, status, created_at
	`
	
	var booking data.Booking
	err := r.db.QueryRow(
		query,
		userID,
		roomID,
		fromDate,
		toDate,
	).Scan(
		&booking.ID,
		&booking.UserID,
		&booking.RoomID,
		&booking.FromDate,
		&booking.ToDate,
		&booking.Status,
		&booking.CreatedAt,
	)
	
	if err != nil {
		return nil, err
	}
	
	return &booking, nil
}

func (r *BookingRepository) GetBooking(id int) (*data.Booking, error) {
	query := `
		SELECT id, user_id, room_id, from_date, to_date, status, created_at
		FROM bookings
		WHERE id = $1
	`
	
	var booking data.Booking
	err := r.db.QueryRow(query, id).Scan(
		&booking.ID,
		&booking.UserID,
		&booking.RoomID,
		&booking.FromDate,
		&booking.ToDate,
		&booking.Status,
		&booking.CreatedAt,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return &booking, nil
}

func (r *BookingRepository) UpdateBookingStatus(id int, status string) error {
	query := `UPDATE bookings SET status = $1 WHERE id = $2`
	
	_, err := r.db.Exec(query, status, id)
	return err
}

func (r *BookingRepository) GetUserBookings(userID int) ([]data.Booking, error) {
	query := `
		SELECT id, user_id, room_id, from_date, to_date, status, created_at
		FROM bookings
		WHERE user_id = $1
		ORDER BY created_at DESC
	`
	
	rows, err := r.db.Query(query, userID)
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
			&booking.Status,
			&booking.CreatedAt,
		); err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return bookings, nil
}

func (r *BookingRepository) GetAllBookings() ([]data.Booking, error) {
	query := `SELECT id, user_id, room_id, from_date, to_date, status, created_at FROM bookings`

	rows, err := r.db.Query(query)
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
			&booking.Status,
			&booking.CreatedAt,
		); err != nil {
			return nil, err
		}
		bookings = append(bookings, booking)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return bookings, nil
}