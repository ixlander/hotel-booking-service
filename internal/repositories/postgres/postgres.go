package postgres

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"

	"github.com/ixlander/hotel-booking-service/internal/repositories"
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

func (r *UserRepo) GetAllUsers() ([]*repositories.User, error) {
	rows, err := r.db.Query("SELECT id, name, email FROM users")
	if err != nil {
		log.Printf("Error fetching users: %v", err)
		return nil, err
	}
	defer rows.Close()

	var users []*repositories.User
	for rows.Next() {
		var user repositories.User
		if err := rows.Scan(&user.ID, &user.Name, &user.Email); err != nil {
			log.Printf("Error scanning user: %v", err)
			return nil, err
		}
		users = append(users, &user)
	}
	return users, nil
}

func (r *UserRepo) GetUserByID(id int) (*repositories.User, error) {
	row := r.db.QueryRow("SELECT id, name, email FROM users WHERE id = $1", id)
	var user repositories.User
	if err := row.Scan(&user.ID, &user.Name, &user.Email); err != nil {
		log.Printf("Error fetching user: %v", err)
		return nil, err
	}
	return &user, nil
}

func (r *UserRepo) CreateUser(user *repositories.User) error {
	_, err := r.db.Exec("INSERT INTO users (name, email) VALUES ($1, $2)", user.Name, user.Email)
	if err != nil {
		log.Printf("Error creating user: %v", err)
		return err
	}
	return nil
}

func (r *HotelRepo) GetAllHotels() ([]*repositories.Hotel, error) {
	rows, err := r.db.Query("SELECT id, name, location FROM hotels")
	if err != nil {
		log.Printf("Error fetching hotels: %v", err)
		return nil, err
	}
	defer rows.Close()

	var hotels []*repositories.Hotel
	for rows.Next() {
		var hotel repositories.Hotel
		if err := rows.Scan(&hotel.ID, &hotel.Name, &hotel.Location); err != nil {
			log.Printf("Error scanning hotel: %v", err)
			return nil, err
		}
		hotels = append(hotels, &hotel)
	}
	return hotels, nil
}

func (r *HotelRepo) GetHotelByID(id int) (*repositories.Hotel, error) {
	row := r.db.QueryRow("SELECT id, name, location FROM hotels WHERE id = $1", id)
	var hotel repositories.Hotel
	if err := row.Scan(&hotel.ID, &hotel.Name, &hotel.Location); err != nil {
		log.Printf("Error fetching hotel: %v", err)
		return nil, err
	}
	return &hotel, nil
}

func (r *HotelRepo) CreateHotel(hotel *repositories.Hotel) error {
	_, err := r.db.Exec("INSERT INTO hotels (name, location) VALUES ($1, $2)", hotel.Name, hotel.Location)
	if err != nil {
		log.Printf("Error creating hotel: %v", err)
		return err
	}
	return nil
}

func (r *RoomRepo) GetAllRooms() ([]*repositories.Room, error) {
	rows, err := r.db.Query("SELECT id, hotel_id, room_type, price FROM rooms")
	if err != nil {
		log.Printf("Error fetching rooms: %v", err)
		return nil, err
	}
	defer rows.Close()

	var rooms []*repositories.Room
	for rows.Next() {
		var room repositories.Room
		if err := rows.Scan(&room.ID, &room.HotelID, &room.RoomType, &room.Price); err != nil {
			log.Printf("Error scanning room: %v", err)
			return nil, err
		}
		rooms = append(rooms, &room)
	}
	return rooms, nil
}

func (r *RoomRepo) GetRoomByID(id int) (*repositories.Room, error) {
	row := r.db.QueryRow("SELECT id, hotel_id, room_type, price FROM rooms WHERE id = $1", id)
	var room repositories.Room
	if err := row.Scan(&room.ID, &room.HotelID, &room.RoomType, &room.Price); err != nil {
		log.Printf("Error fetching room: %v", err)
		return nil, err
	}
	return &room, nil
}

func (r *RoomRepo) CreateRoom(room *repositories.Room) error {
	_, err := r.db.Exec("INSERT INTO rooms (hotel_id, room_type, price) VALUES ($1, $2, $3)", room.HotelID, room.RoomType, room.Price)
	if err != nil {
		log.Printf("Error creating room: %v", err)
		return err
	}
	return nil
}

func (r *BookingRepo) GetAllBookings() ([]*repositories.Booking, error) {
	rows, err := r.db.Query("SELECT id, user_id, room_id, check_in, check_out FROM bookings")
	if err != nil {
		log.Printf("Error fetching bookings: %v", err)
		return nil, err
	}
	defer rows.Close()

	var bookings []*repositories.Booking
	for rows.Next() {
		var booking repositories.Booking
		if err := rows.Scan(&booking.ID, &booking.UserID, &booking.RoomID, &booking.CheckIn, &booking.CheckOut); err != nil {
			log.Printf("Error scanning booking: %v", err)
			return nil, err
		}
		bookings = append(bookings, &booking)
	}
	return bookings, nil
}

func (r *BookingRepo) GetBookingByID(id int) (*repositories.Booking, error) {
	row := r.db.QueryRow("SELECT id, user_id, room_id, check_in, check_out FROM bookings WHERE id = $1", id)
	var booking repositories.Booking
	if err := row.Scan(&booking.ID, &booking.UserID, &booking.RoomID, &booking.CheckIn, &booking.CheckOut); err != nil {
		log.Printf("Error fetching booking: %v", err)
		return nil, err
	}
	return &booking, nil
}

func (r *BookingRepo) CreateBooking(booking *repositories.Booking) error {
	_, err := r.db.Exec("INSERT INTO bookings (user_id, room_id, check_in, check_out) VALUES ($1, $2, $3, $4)", booking.UserID, booking.RoomID, booking.CheckIn, booking.CheckOut)
	if err != nil {
		log.Printf("Error creating booking: %v", err)
		return err
	}
	return nil
}