package repositories

import (
	"time"

	"database/sql"
	"hotel-booking-service/internal/data"
)

type RoomRepository struct {
	db *sql.DB
}

func NewRoomRepository(db *sql.DB) *RoomRepository {
	return &RoomRepository{db: db}
}

func (r *RoomRepository) GetByID(id int) (*data.Room, error) {
	query := `SELECT id, hotel_id, number, capacity, price FROM rooms WHERE id = $1`
	
	var room data.Room
	err := r.db.QueryRow(query, id).Scan(
		&room.ID,
		&room.HotelID,
		&room.Number,
		&room.Capacity,
		&room.Price,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return &room, nil
}

func (r *RoomRepository) CheckRoomAvailability(roomID int, fromDate, toDate time.Time) (bool, error) {
	query := `
		SELECT COUNT(*) FROM bookings 
		WHERE room_id = $1 
		AND status != 'cancelled'
		AND (
			(from_date <= $2 AND to_date >= $2) OR
			(from_date <= $3 AND to_date >= $3) OR
			(from_date >= $2 AND to_date <= $3)
		)
	`
	
	var count int
	err := r.db.QueryRow(query, roomID, fromDate, toDate).Scan(&count)
	if err != nil {
		return false, err
	}
	
	return count == 0, nil
}

func (r *RoomRepository) GetAvailableRoomsByHotelID(hotelID int, fromDate, toDate time.Time) ([]data.Room, error) {
	query := `
		SELECT r.id, r.hotel_id, r.number, r.capacity, r.price
		FROM rooms r
		WHERE r.hotel_id = $1
		AND NOT EXISTS (
			SELECT 1 FROM bookings b
			WHERE b.room_id = r.id
			AND b.status != 'cancelled'
			AND (
				(b.from_date <= $2 AND b.to_date >= $2) OR
				(b.from_date <= $3 AND b.to_date >= $3) OR
				(b.from_date >= $2 AND b.to_date <= $3)
			)
		)
	`
	
	rows, err := r.db.Query(query, hotelID, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var rooms []data.Room
	for rows.Next() {
		var room data.Room
		if err := rows.Scan(
			&room.ID,
			&room.HotelID,
			&room.Number,
			&room.Capacity,
			&room.Price,
		); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return rooms, nil
}

func (r *RoomRepository) GetAllRooms() ([]data.Room, error) {
	query := `SELECT id, hotel_id, number, capacity, price FROM rooms`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rooms []data.Room
	for rows.Next() {
		var room data.Room
		if err := rows.Scan(&room.ID, &room.HotelID, &room.Number, &room.Capacity, &room.Price); err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}
	return rooms, nil
}

func (r *RoomRepository) CreateRoom(room *data.Room) (*data.Room, error) {
	query := `INSERT INTO rooms (hotel_id, number, capacity, price) 
	          VALUES ($1, $2, $3, $4) 
	          RETURNING id`
	
	err := r.db.QueryRow(query, room.HotelID, room.Number, room.Capacity, room.Price).Scan(&room.ID)
	if err != nil {
		return nil, err
	}

	return room, nil
}

func (r *RoomRepository) UpdateRoom(room *data.Room) (*data.Room, error) {
	query := `UPDATE rooms SET hotel_id = $1, number = $2, capacity = $3, price = $4 WHERE id = $5`
	_, err := r.db.Exec(query, room.HotelID, room.Number, room.Capacity, room.Price, room.ID)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (r *RoomRepository) DeleteRoom(roomID int) error {
	query := `DELETE FROM rooms WHERE id = $1`
	_, err := r.db.Exec(query, roomID)
	if err != nil {
		return err
	}
	return nil
}