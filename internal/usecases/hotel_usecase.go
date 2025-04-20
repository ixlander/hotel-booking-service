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

func NewHotelUsecase(hotelRepo repositories.HotelRepository, roomRepo repositories.RoomRepository) *HotelUsecase {
	return &HotelUsecase{
		hotelRepo: hotelRepo,
		roomRepo:  roomRepo,
	}
}

func (u *HotelUsecase) GetAllHotels(ctx context.Context) ([]*data.Hotel, error) {
	hotels, err := u.hotelRepo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	for _, hotel := range hotels {
		rooms, err := u.roomRepo.GetByHotelID(ctx, hotel.ID)
		if err != nil {
			return nil, err
		}
		hotel.Rooms = make([]data.Room, len(rooms))
		for i, room := range rooms {
			hotel.Rooms[i] = *room
		}
	}

	return hotels, nil
}

func (u *HotelUsecase) GetHotelWithAvailableRooms(ctx context.Context, hotelID int64, fromDate, toDate time.Time) (*data.Hotel, error) {
	hotel, err := u.hotelRepo.FindByID(ctx, hotelID)
	if err != nil {
		return nil, err
	}
	if hotel == nil {
		return nil, nil 
	}

	rooms, err := u.roomRepo.GetAvailableRooms(ctx, hotelID, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	hotel.Rooms = make([]data.Room, len(rooms))
	for i, room := range rooms {
		hotel.Rooms[i] = *room
	}

	return hotel, nil
}