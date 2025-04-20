package data

import (
	"time"
)

type User struct {
	ID        int64     `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Hotel struct {
	ID    int64  `json:"id" db:"id"`
	Name  string `json:"name" db:"name"`
	City  string `json:"city" db:"city"`
	Rooms []Room `json:"rooms,omitempty"`
}

type Room struct {
	ID       int64   `json:"id" db:"id"`
	HotelID  int64   `json:"hotel_id" db:"hotel_id"`
	Number   string  `json:"number" db:"number"`
	Capacity int     `json:"capacity" db:"capacity"`
	Price    float64 `json:"price" db:"price"`
}

type Booking struct {
	ID        int64     `json:"id" db:"id"`
	UserID    int64     `json:"user_id" db:"user_id"`
	RoomID    int64     `json:"room_id" db:"room_id"`
	FromDate  time.Time `json:"from_date" db:"from_date"`
	ToDate    time.Time `json:"to_date" db:"to_date"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	Status    string    `json:"status" db:"status"`
}

type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type BookingRequest struct {
	RoomID   int64     `json:"room_id" validate:"required"`
	FromDate time.Time `json:"from_date" validate:"required"`
	ToDate   time.Time `json:"to_date" validate:"required"`
}

type BookingResponse struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	RoomID    int64     `json:"room_id"`
	FromDate  time.Time `json:"from_date"`
	ToDate    time.Time `json:"to_date"`
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"`
}