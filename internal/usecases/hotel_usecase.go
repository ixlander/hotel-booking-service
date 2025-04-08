package usecases

import (
	"context"
	"time"

	"github.com/ixlander/hotel-booking-service/internal/data"
	"github.com/ixlander/hotel-booking-service/internal/repositories"
)

type HotelUsecase struct {
	hotelRepo repositories.HotelRepository
	roomRepo  repositories.RoomRepository
}

func NewHotelUsecase(
	hotelRepo repositories.HotelRepository,
	roomRepo repositories.RoomRepository,
) *HotelUsecase {
	return &HotelUsecase{
		hotelRepo: hotelRepo,
		roomRepo:  roomRepo,
	}
}

func (u *HotelUsecase) GetAllHotels(ctx context.Context) ([]data.Hotel, error) {
	hotels, err := u.hotelRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}
	
	return hotels, nil
}

func (u *HotelUsecase) GetHotelWithRooms(ctx context.Context, hotelID int64) (*data.Hotel, error) {
	hotel, err := u.hotelRepo.FindByID(ctx, hotelID)
	if err != nil {
		return nil, err
	}
	
	if hotel == nil {
		return nil, nil
	}
	
	rooms, err := u.roomRepo.GetByHotelID(ctx, hotelID)
	if err != nil {
		return nil, err
	}
	
	hotel.Rooms = rooms
	
	return hotel, nil
}

func (u *HotelUsecase) GetAvailableRooms(ctx context.Context, hotelID int64, fromDate, toDate time.Time) ([]data.Room, error) {
	return u.roomRepo.GetAvailableRooms(ctx, hotelID, fromDate, toDate)
}