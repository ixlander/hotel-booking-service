package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/yourusername/hotel-booking-service/internal/data"
	"github.com/yourusername/hotel-booking-service/internal/repositories"
)

var (
	ErrRoomNotFound      = errors.New("room not found")
	ErrRoomNotAvailable  = errors.New("room not available for selected dates")
	ErrBookingNotFound   = errors.New("booking not found")
	ErrUnauthorized      = errors.New("user is not authorized to perform this action")
	ErrInvalidDateRange  = errors.New("invalid date range")
)

type BookingUsecase struct {
	bookingRepo repositories.BookingRepository
	roomRepo    repositories.RoomRepository
}

func NewBookingUsecase(
	bookingRepo repositories.BookingRepository,
	roomRepo repositories.RoomRepository,
) *BookingUsecase {
	return &BookingUsecase{
		bookingRepo: bookingRepo,
		roomRepo:    roomRepo,
	}
}

func (u *BookingUsecase) CreateBooking(ctx context.Context, userID, roomID int64, fromDate, toDate time.Time) (*data.Booking, error) {
	if fromDate.After(toDate) || fromDate.Equal(toDate) {
		return nil, ErrInvalidDateRange
	}
	
	room, err := u.roomRepo.FindByID(ctx, roomID)
	if err != nil {
		return nil, err
	}
	
	if room == nil {
		return nil, ErrRoomNotFound
	}
	
	available, err := u.roomRepo.CheckAvailability(ctx, roomID, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	
	if !available {
		return nil, ErrRoomNotAvailable
	}
	
	booking := &data.Booking{
		UserID:   userID,
		RoomID:   roomID,
		FromDate: fromDate,
		ToDate:   toDate,
		Status:   "active",
	}
	
	return u.bookingRepo.Create(ctx, booking)
}

func (u *BookingUsecase) GetUserBookings(ctx context.Context, userID int64) ([]data.Booking, error) {
	return u.bookingRepo.FindByUserID(ctx, userID)
}

func (u *BookingUsecase) CancelBooking(ctx context.Context, userID, bookingID int64) error {
	booking, err := u.bookingRepo.FindByID(ctx, bookingID)
	if err != nil {
		return err
	}
	
	if booking == nil {
		return ErrBookingNotFound
	}
	
	if booking.UserID != userID {
		return ErrUnauthorized
	}
	
	return u.bookingRepo.UpdateStatus(ctx, bookingID, "cancelled")
}