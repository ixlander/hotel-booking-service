package services

import (
	"hotel-booking-service/internal/data"
	"hotel-booking-service/internal/repositories"
)

type HotelService struct {
	hotelRepo *repositories.HotelRepository
}

func NewHotelService(hotelRepo *repositories.HotelRepository) *HotelService {
	return &HotelService{hotelRepo: hotelRepo}
}

func (s *HotelService) GetAllHotels() ([]*data.Hotel, error) {
	hotels, err := s.hotelRepo.GetAllHotels()
	if err != nil {
		return nil, err
	}

	var hotelPointers []*data.Hotel
	for i := range hotels {
		hotelPointers = append(hotelPointers, &hotels[i])
	}

	return hotelPointers, nil
}

func (s *HotelService) GetHotelByID(id int) (*data.Hotel, error) {
	hotel, err := s.hotelRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return hotel, nil
}

func (s *HotelService) CreateHotel(hotel *data.Hotel) error {
	return s.hotelRepo.CreateHotel(hotel)
}