package services

import (
	"time"

	"hotel-booking-service/internal/data"
	"hotel-booking-service/internal/repositories"
)

type BookingService struct {
	bookingRepo *repositories.BookingRepository
}

func NewBookingService(bookingRepo *repositories.BookingRepository) *BookingService {
	return &BookingService{bookingRepo: bookingRepo}
}

func (s *BookingService) GetAllBookings() ([]*data.Booking, error) {
	bookings, err := s.bookingRepo.GetAllBookings()
	if err != nil {
		return nil, err
	}

	var bookingPointers []*data.Booking
	for i := range bookings {
		bookingPointers = append(bookingPointers, &bookings[i])
	}

	return bookingPointers, nil
}

func (s *BookingService) GetBookingByID(id int) (*data.Booking, error) {
	booking, err := s.bookingRepo.GetBooking(id)
	if err != nil {
		return nil, err
	}
	return booking, nil
}

func (s *BookingService) CreateBooking(userID, roomID int, fromDate, toDate time.Time) (*data.Booking, error) {
	return s.bookingRepo.CreateBooking(userID, roomID, fromDate, toDate)
}