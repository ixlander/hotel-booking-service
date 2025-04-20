package store

import (
	"database/sql"

	"github.com/ixlander/hotel-booking-service/internal/repositories"
	"github.com/ixlander/hotel-booking-service/internal/repositories/postgres"
)

type Store struct {
	UserRepo    repositories.UserRepository
	HotelRepo   repositories.HotelRepository
	RoomRepo    repositories.RoomRepository
	BookingRepo repositories.BookingRepository
}

func NewStore(db *sql.DB) *Store {
	return &Store{
		UserRepo:    postgres.NewUserRepo(db),
		HotelRepo:   postgres.NewHotelRepo(db),
		RoomRepo:    postgres.NewRoomRepo(db),
		BookingRepo: postgres.NewBookingRepo(db),
	}
}