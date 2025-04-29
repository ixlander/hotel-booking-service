package repositories

import (
	"database/sql"
	"fmt"
	"hotel-booking-service/internal/data"
)

type HotelRepository struct {
	db *sql.DB
}

func NewHotelRepository(db *sql.DB) *HotelRepository {
	return &HotelRepository{db: db}
}

func (r *HotelRepository) GetAllHotels() ([]data.Hotel, error) {
	query := `SELECT id, name, city FROM hotels`
	
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	
	var hotels []data.Hotel
	for rows.Next() {
		var hotel data.Hotel
		if err := rows.Scan(
			&hotel.ID,
			&hotel.Name,
			&hotel.City,
		); err != nil {
			return nil, err
		}
		hotels = append(hotels, hotel)
	}
	
	if err = rows.Err(); err != nil {
		return nil, err
	}
	
	return hotels, nil
}

func (r *HotelRepository) GetByID(id int) (*data.Hotel, error) {
	query := `SELECT id, name, city FROM hotels WHERE id = $1`
	
	var hotel data.Hotel
	err := r.db.QueryRow(query, id).Scan(
		&hotel.ID,
		&hotel.Name,
		&hotel.City,
	)
	
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	
	return &hotel, nil
}

func (r *HotelRepository) GetRoomsByHotelID(hotelID int) ([]data.Room, error) {
	query := `
		SELECT id, hotel_id, number, capacity, price
		FROM rooms
		WHERE hotel_id = $1
	`
	
	rows, err := r.db.Query(query, hotelID)
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

func (r *HotelRepository) CreateHotel(hotel data.Hotel) (*data.Hotel, error) {
	query := `INSERT INTO hotels (name, city) VALUES ($1, $2) RETURNING id`
	err := r.db.QueryRow(query, hotel.Name, hotel.City).Scan(&hotel.ID)
	if err != nil {
		return nil, fmt.Errorf("could not insert hotel: %v", err)
	}
	return &hotel, nil
}

func (r *HotelRepository) UpdateHotel(hotel data.Hotel) (*data.Hotel, error) {
	query := `UPDATE hotels SET name=$1, city=$2 WHERE id=$3`
	_, err := r.db.Exec(query, hotel.Name, hotel.City, hotel.ID)
	if err != nil {
		return nil, err
	}
	return &hotel, nil
}

func (r *HotelRepository) DeleteHotel(id int) error {
	_, err := r.db.Exec(`DELETE FROM hotels WHERE id = $1`, id)
	return err
}