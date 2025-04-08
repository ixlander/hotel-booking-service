package store

import (
	"database/sql"

	"github.com/ixlander/hotel-booking-service/internal/repositories"
)

type Store struct {
	UserRepo    repositories.UserRepository
	HotelRepo   repositories.HotelRepository
	RoomRepo    repositories.RoomRepository
	BookingRepo repositories.BookingRepository
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		UserRepo:    repositories.NewPostgresUserRepository(db),
		HotelRepo:   repositories.NewPostgresHotelRepository(db),
		RoomRepo:    repositories.NewPostgresRoomRepository(db),
		BookingRepo: repositories.NewPostgresBookingRepository(db),
	}
}