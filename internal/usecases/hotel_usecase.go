package usecases

import (
	"time"
	
	"hotel-booking-service/internal/data"
	"hotel-booking-service/internal/repositories"
)

type HotelUsecase struct {
	hotelRepo *repositories.HotelRepository
	roomRepo  *repositories.RoomRepository
}

func NewHotelUsecase(
	hotelRepo *repositories.HotelRepository,
	roomRepo *repositories.RoomRepository,
) *HotelUsecase {
	return &HotelUsecase{
		hotelRepo: hotelRepo,
		roomRepo:  roomRepo,
	}
}

func (uc *HotelUsecase) GetAllHotels(fromDate, toDate time.Time) ([]data.Hotel, error) {
	hotels, err := uc.hotelRepo.GetAllHotels()
	if err != nil {
		return nil, err
	}
	
	for i := range hotels {
		availableRooms, err := uc.roomRepo.GetAvailableRoomsByHotelID(hotels[i].ID, fromDate, toDate)
		if err != nil {
			return nil, err
		}
		
		hotels[i].Rooms = availableRooms
	}
	
	return hotels, nil
}

func (uc *HotelUsecase) GetHotelByID(id int, fromDate, toDate time.Time) (*data.Hotel, error) {
	hotel, err := uc.hotelRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	
	if hotel == nil {
		return nil, nil
	}
	
	availableRooms, err := uc.roomRepo.GetAvailableRoomsByHotelID(id, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	
	hotel.Rooms = availableRooms
	
	return hotel, nil
}

func (uc *HotelUsecase) GetRoomsByHotelID(hotelID int, fromDate, toDate time.Time) ([]data.Room, error) {
	rooms, err := uc.roomRepo.GetAvailableRoomsByHotelID(hotelID, fromDate, toDate)
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

func (uc *HotelUsecase) CreateHotel(hotel data.Hotel) (*data.Hotel, error) {
	return uc.hotelRepo.CreateHotel(hotel)
}

func (uc *HotelUsecase) UpdateHotel(hotel data.Hotel) (*data.Hotel, error) {
	return uc.hotelRepo.UpdateHotel(hotel)
}

func (uc *HotelUsecase) DeleteHotel(hotelID int) error {
	return uc.hotelRepo.DeleteHotel(hotelID)
}

func (uc *HotelUsecase) CreateRoom(room data.Room) (*data.Room, error) {
	return uc.roomRepo.CreateRoom(&room)
}

func (uc *HotelUsecase) UpdateRoom(room data.Room) (*data.Room, error) {
	return uc.roomRepo.UpdateRoom(&room)
}

func (uc *HotelUsecase) DeleteRoom(roomID int) error {
	return uc.roomRepo.DeleteRoom(roomID)
}