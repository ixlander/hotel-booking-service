package store

import (
	"database/sql"
	
	"hotel-booking-service/internal/repositories"
)

type Store struct {
	UserRepo    *repositories.UserRepository
	HotelRepo   *repositories.HotelRepository
	RoomRepo    *repositories.RoomRepository
	BookingRepo *repositories.BookingRepository
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		UserRepo:    repositories.NewUserRepository(db),
		HotelRepo:   repositories.NewHotelRepository(db),
		RoomRepo:    repositories.NewRoomRepository(db),
		BookingRepo: repositories.NewBookingRepository(db),
	}
}