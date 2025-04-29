package usecases

import (
	"errors"
	"time"
	
	"hotel-booking-service/internal/data"
	"hotel-booking-service/internal/repositories"
)

type BookingUsecase struct {
	bookingRepo *repositories.BookingRepository
	roomRepo    *repositories.RoomRepository
}

func NewBookingUsecase(
	bookingRepo *repositories.BookingRepository,
	roomRepo *repositories.RoomRepository,
) *BookingUsecase {
	return &BookingUsecase{
		bookingRepo: bookingRepo,
		roomRepo:    roomRepo,
	}
}

func (uc *BookingUsecase) CreateBooking(userID, roomID int, fromDate, toDate time.Time) (*data.Booking, error) {
	if fromDate.After(toDate) {
		return nil, errors.New("from date must be before to date")
	}
	
	if fromDate.Before(time.Now()) {
		return nil, errors.New("from date must be in the future")
	}
	
	room, err := uc.roomRepo.GetByID(roomID)
	if err != nil {
		return nil, err
	}
	
	if room == nil {
		return nil, errors.New("room not found")
	}
	
	available, err := uc.roomRepo.CheckRoomAvailability(roomID, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	
	if !available {
		return nil, errors.New("room not available for the selected dates")
	}
	
	booking, err := uc.bookingRepo.CreateBooking(userID, roomID, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	
	return booking, nil
}

func (uc *BookingUsecase) CancelBooking(userID, bookingID int) error {
	booking, err := uc.bookingRepo.GetBooking(bookingID)
	if err != nil {
		return err
	}
	
	if booking == nil {
		return errors.New("booking not found")
	}
	
	if booking.UserID != userID {
		return errors.New("booking does not belong to this user")
	}
	
	if booking.Status == "cancelled" {
		return errors.New("booking is already cancelled")
	}
	
	return uc.bookingRepo.UpdateBookingStatus(bookingID, "cancelled")
}

func (uc *BookingUsecase) GetUserBookings(userID int) ([]data.Booking, error) {
	return uc.bookingRepo.GetUserBookings(userID)
}

func (uc *BookingUsecase) GetBookingByID(bookingID int) (*data.Booking, error) {
	return uc.bookingRepo.GetBooking(bookingID)
}

func (uc *BookingUsecase) UpdateBooking(bookingID int, status string) error {
	return uc.bookingRepo.UpdateBookingStatus(bookingID, status)
}