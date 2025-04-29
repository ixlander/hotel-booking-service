package services

import (
	"hotel-booking-service/internal/data"
	"hotel-booking-service/internal/repositories"
	"time"
)

type RoomService struct {
	roomRepo *repositories.RoomRepository
}

func NewRoomService(roomRepo *repositories.RoomRepository) *RoomService {
	return &RoomService{roomRepo: roomRepo}
}

func (s *RoomService) GetAllRooms() ([]*data.Room, error) {
	rooms, err := s.roomRepo.GetAllRooms()
	if err != nil {
		return nil, err
	}

	var roomPointers []*data.Room
	for i := range rooms {
		roomPointers = append(roomPointers, &rooms[i])
	}

	return roomPointers, nil
}

func (s *RoomService) GetRoomByID(id int) (*data.Room, error) {
	room, err := s.roomRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return room, nil
}

func (s *RoomService) CreateRoom(room *data.Room) (*data.Room, error) {
	createdRoom, err := s.roomRepo.CreateRoom(room)
	if err != nil {
		return nil, err
	}
	return createdRoom, nil
}

func (s *RoomService) GetAvailableRoomsByHotelID(hotelID int, fromDate, toDate time.Time) ([]*data.Room, error) {
	rooms, err := s.roomRepo.GetAvailableRoomsByHotelID(hotelID, fromDate, toDate)
	if err != nil {
		return nil, err
	}

	var roomPointers []*data.Room
	for i := range rooms {
		roomPointers = append(roomPointers, &rooms[i])
	}

	return roomPointers, nil
}