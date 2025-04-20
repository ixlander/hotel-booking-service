package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/ixlander/hotel-booking-service/internal/data"
)

type RoomRepo struct {
	db *sql.DB
}

func NewRoomRepo(db *sql.DB) *RoomRepo {
	return &RoomRepo{db: db}
}

func (r *RoomRepo) GetByHotelID(ctx context.Context, hotelID int64) ([]*data.Room, error) {
	query := `SELECT id, hotel_id, number, capacity, price FROM rooms WHERE hotel_id = $1`
	
	rows, err := r.db.QueryContext(ctx, query, hotelID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var rooms []*data.Room
	for rows.Next() {
		var room data.Room
		if err := rows.Scan(&room.ID, &room.HotelID, &room.Number, &room.Capacity, &room.Price); err != nil {
			return nil, err
		}
		rooms = append(rooms, &room)
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return rooms, nil
}

func (r *RoomRepo) FindByID(ctx context.Context, id int64) (*data.Room, error) {
	query := `SELECT id, hotel_id, number, capacity, price FROM rooms WHERE id = $1`
	
	var room data.Room
	err := r.db.QueryRowContext(ctx, query, id).
		Scan(&room.ID, &room.HotelID, &room.Number, &room.Capacity, &room.Price)
	
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil 
		}
		return nil, err
	}
	
	return &room, nil
}

func (r *RoomRepo) CheckAvailability(ctx context.Context, roomID int64, fromDate, toDate time.Time) (bool, error) {
	query := `
		SELECT COUNT(*) FROM bookings 
		WHERE room_id = $1 AND status != 'cancelled'
		AND (
			(from_date <= $2 AND to_date >= $2) OR
			(from_date <= $3 AND to_date >= $3) OR
			(from_date >= $2 AND to_date <= $3)
		)
	`
	
	var count int
	err := r.db.QueryRowContext(ctx, query, roomID, fromDate, toDate).Scan(&count)
	if err != nil {
		return false, err
	}
	
	return count == 0, nil
}

func (r *RoomRepo) GetAvailableRooms(ctx context.Context, hotelID int64, fromDate, toDate time.Time) ([]*data.Room, error) {
	query := `
		SELECT r.id, r.hotel_id, r.number, r.capacity, r.price
		FROM rooms r
		WHERE r.hotel_id = $1
		AND NOT EXISTS (
			SELECT 1 FROM bookings b
			WHERE b.room_id = r.id AND b.status != 'cancelled'
			AND (
				(b.from_date <= $2 AND b.to_date >= $2) OR
				(b.from_date <= $3 AND b.to_date >= $3) OR
				(b.from_date >= $2 AND b.to_date <= $3)
			)
		)
	`
	
	rows, err := r.db.QueryContext(ctx, query, hotelID, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var rooms []*data.Room
	for rows.Next() {
		var room data.Room
		if err := rows.Scan(&room.ID, &room.HotelID, &room.Number, &room.Capacity, &room.Price); err != nil {
			return nil, err
		}
		rooms = append(rooms, &room)
	}
	
	if err := rows.Err(); err != nil {
		return nil, err
	}
	
	return rooms, nil
}