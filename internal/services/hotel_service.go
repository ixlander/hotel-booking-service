package services

import (
	"hotel-booking-service/internal/data"
	"hotel-booking-service/internal/repositories"
	"time"
)

type HotelService struct {
	hotelRepo *repositories.HotelRepository
	roomRepo  *repositories.RoomRepository 
}

func NewHotelService(hotelRepo *repositories.HotelRepository, roomRepo *repositories.RoomRepository) *HotelService {
	return &HotelService{
		hotelRepo: hotelRepo,
		roomRepo:  roomRepo,
	}
}

func (s *HotelService) GetAllHotels(fromDate, toDate time.Time) ([]*data.Hotel, error) {
	hotels, err := s.hotelRepo.GetAllHotels()
	if err != nil {
		return nil, err
	}

	var hotelPointers []*data.Hotel
	for i := range hotels {
		availableRooms, err := s.roomRepo.GetAvailableRoomsByHotelID(hotels[i].ID, fromDate, toDate)
		if err != nil {
			return nil, err
		}
		hotels[i].Rooms = availableRooms
		hotelPointers = append(hotelPointers, &hotels[i])
	}

	return hotelPointers, nil
}

func (s *HotelService) GetHotelByID(id int, fromDate, toDate time.Time) (*data.Hotel, error) {
	hotel, err := s.hotelRepo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if hotel == nil {
		return nil, nil
	}

	availableRooms, err := s.roomRepo.GetAvailableRoomsByHotelID(id, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	hotel.Rooms = availableRooms

	return hotel, nil
}
