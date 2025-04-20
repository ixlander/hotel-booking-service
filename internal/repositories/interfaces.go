package repositories

import (
	"context"
	"time"

	"github.com/ixlander/hotel-booking-service/internal/data"
)

type UserRepository interface {
	Create(ctx context.Context, user *data.User) error
	FindByID(ctx context.Context, id int64) (*data.User, error)
	FindByEmail(ctx context.Context, email string) (*data.User, error)
}

type HotelRepository interface {
	GetAll(ctx context.Context) ([]*data.Hotel, error)
	FindByID(ctx context.Context, id int64) (*data.Hotel, error)
}

type RoomRepository interface {
	GetByHotelID(ctx context.Context, hotelID int64) ([]*data.Room, error)
	FindByID(ctx context.Context, id int64) (*data.Room, error)
	CheckAvailability(ctx context.Context, roomID int64, fromDate, toDate time.Time) (bool, error)
	GetAvailableRooms(ctx context.Context, hotelID int64, fromDate, toDate time.Time) ([]*data.Room, error)
}

type BookingRepository interface {
	Create(ctx context.Context, booking *data.Booking) error
	FindByID(ctx context.Context, id int64) (*data.Booking, error)
	FindByUserID(ctx context.Context, userID int64) ([]*data.Booking, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
}