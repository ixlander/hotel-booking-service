// internal/data/models.go
package data

import (
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // Don't expose password in JSON responses
	CreatedAt time.Time `json:"created_at"`
}

type Hotel struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	City  string `json:"city"`
	Rooms []Room `json:"rooms,omitempty"`
}

type Room struct {
	ID       int64   `json:"id"`
	HotelID  int64   `json:"hotel_id"`
	Number   string  `json:"number"`
	Capacity int     `json:"capacity"`
	Price    float64 `json:"price"`
}

type Booking struct {
	ID        int64     `json:"id"`
	UserID    int64     `json:"user_id"`
	RoomID    int64     `json:"room_id"`
	FromDate  time.Time `json:"from_date"`
	ToDate    time.Time `json:"to_date"`
	CreatedAt time.Time `json:"created_at"`
	Status    string    `json:"status"` // active, cancelled
}

// Request/Response structures

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
	User  User   `json:"user"`
}

type BookingRequest struct {
	RoomID   int64     `json:"room_id" binding:"required"`
	FromDate time.Time `json:"from_date" binding:"required"`
	ToDate   time.Time `json:"to_date" binding:"required"`
}

type HotelResponse struct {
	Hotels []Hotel `json:"hotels"`
}

type ApiError struct {
	Error string `json:"error"`
}