package postgres

import (
	"database/sql"
	"log"
	"hotel-booking-service/internal/data"
	_ "github.com/lib/pq"
)

type UserRepo struct {
	db *sql.DB
}

type HotelRepo struct {
	db *sql.DB
}

type RoomRepo struct {
	db *sql.DB
}

type BookingRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *UserRepo {
	return &UserRepo{db: db}
}

func NewHotelRepo(db *sql.DB) *HotelRepo {
	return &HotelRepo{db: db}
}

func NewRoomRepo(db *sql.DB) *RoomRepo {
	return &RoomRepo{db: db}
}

func NewBookingRepo(db *sql.DB) *BookingRepo {
	return &BookingRepo{db: db}
}

// UserRepo Methods
func (r *UserRepo) GetAllUsers() ([]*data.User, error) {
	rows, err := r.db.Query("SELECT id, email, created_at FROM users")
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		return nil, err
	}
	defer rows.Close()

	var users []*data.User
	for rows.Next() {
		var user data.User
		if err := rows.Scan(&user.ID, &user.Email, &user.CreatedAt); err != nil {
			log.Printf("Error scanning user: %v", err)
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *UserRepo) GetUserByID(id int) (*data.User, error) {
	row := r.db.QueryRow("SELECT id, email, created_at FROM users WHERE id = $1", id)
	var user data.User
	if err := row.Scan(&user.ID, &user.Email, &user.CreatedAt); err != nil {
		log.Printf("Error fetching user: %v", err)
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) CreateUser(user *data.User) error {
	_, err := r.db.Exec("INSERT INTO users (email, created_at) VALUES ($1, $2)", user.Email, user.CreatedAt)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return err
	}
	return nil
}

func (r *HotelRepo) GetAllHotels() ([]*data.Hotel, error) {
	rows, err := r.db.Query("SELECT id, name, city FROM hotels")
	if err != nil {
		log.Printf("Error fetching hotels: %v", err)
		return nil, err
	}
	defer rows.Close()

	var hotels []*data.Hotel
	for rows.Next() {
		var hotel data.Hotel
		if err := rows.Scan(&hotel.ID, &hotel.Name, &hotel.City); err != nil {
			log.Printf("Error scanning hotel: %v", err)
			return nil, err
		}
		hotels = append(hotels, &hotel)
	}
	return hotels, nil
}

func (r *HotelRepo) GetHotelByID(id int) (*data.Hotel, error) {
	row := r.db.QueryRow("SELECT id, name, city FROM hotels WHERE id = $1", id)
	var hotel data.Hotel
	if err := row.Scan(&hotel.ID, &hotel.Name, &hotel.City); err != nil {
		log.Printf("Error fetching hotel: %v", err)
		return nil, err
	}
	return &hotel, nil
}

func (r *HotelRepo) CreateHotel(hotel *data.Hotel) error {
	_, err := r.db.Exec("INSERT INTO hotels (name, city) VALUES ($1, $2)", hotel.Name, hotel.City)
	if err != nil {
		log.Printf("Error creating hotel: %v", err)
		return err
	}
	return nil
}

func (r *RoomRepo) GetAllRooms() ([]*data.Room, error) {
	rows, err := r.db.Query("SELECT id, hotel_id, number, capacity, price FROM rooms")
	if err != nil {
		log.Printf("Error fetching rooms: %v", err)
		return nil, err
	}
	defer rows.Close()

	var rooms []*data.Room
	for rows.Next() {
		var room data.Room
		if err := rows.Scan(&room.ID, &room.HotelID, &room.Number, &room.Capacity, &room.Price); err != nil {
			log.Printf("Error scanning room: %v", err)
			return nil, err
		}
		rooms = append(rooms, &room)
	}
	return rooms, nil
}

func (r *RoomRepo) GetRoomByID(id int) (*data.Room, error) {
	row := r.db.QueryRow("SELECT id, hotel_id, number, capacity, price FROM rooms WHERE id = $1", id)
	var room data.Room
	if err := row.Scan(&room.ID, &room.HotelID, &room.Number, &room.Capacity, &room.Price); err != nil {
		log.Printf("Error fetching room: %v", err)
		return nil, err
	}
	return &room, nil
}

func (r *RoomRepo) CreateRoom(room *data.Room) error {
	_, err := r.db.Exec("INSERT INTO rooms (hotel_id, number, capacity, price) VALUES ($1, $2, $3, $4)", room.HotelID, room.Number, room.Capacity, room.Price)
	if err != nil {
		log.Printf("Error creating room: %v", err)
		return err
	}
	return nil
}

func (r *BookingRepo) GetAllBookings() ([]*data.Booking, error) {
	rows, err := r.db.Query("SELECT id, user_id, room_id, from_date, to_date, status, created_at FROM bookings")
	if err != nil {
		log.Printf("Error fetching bookings: %v", err)
		return nil, err
	}
	defer rows.Close()

	var bookings []*data.Booking
	for rows.Next() {
		var booking data.Booking
		if err := rows.Scan(&booking.ID, &booking.UserID, &booking.RoomID, &booking.FromDate, &booking.ToDate, &booking.Status, &booking.CreatedAt); err != nil {
			log.Printf("Error scanning booking: %v", err)
			return nil, err
		}
		bookings = append(bookings, &booking)
	}
	return bookings, nil
}

func (r *BookingRepo) GetBookingByID(id int) (*data.Booking, error) {
	row := r.db.QueryRow("SELECT id, user_id, room_id, from_date, to_date, status, created_at FROM bookings WHERE id = $1", id)
	var booking data.Booking
	if err := row.Scan(&booking.ID, &booking.UserID, &booking.RoomID, &booking.FromDate, &booking.ToDate, &booking.Status, &booking.CreatedAt); err != nil {
		log.Printf("Error fetching booking: %v", err)
		return nil, err
	}
	return &booking, nil
}

func (r *BookingRepo) CreateBooking(booking *data.Booking) error {
	_, err := r.db.Exec("INSERT INTO bookings (user_id, room_id, from_date, to_date, status, created_at) VALUES ($1, $2, $3, $4, $5, $6)", booking.UserID, booking.RoomID, booking.FromDate, booking.ToDate, booking.Status, booking.CreatedAt)
	if err != nil {
		log.Printf("Error creating booking: %v", err)
		return err
	}
	return nil
}