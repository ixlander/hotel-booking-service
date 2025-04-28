package services

import (
	"hotel-booking-service/internal/data"
	"hotel-booking-service/internal/repositories"
)

type UserService struct {
	userRepo *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{userRepo: userRepo}
}

func (s *UserService) GetAllUsers() ([]*data.User, error) {
	users, err := s.userRepo.GetAllUsers()
	if err != nil {
		return nil, err
	}

	var userPointers []*data.User
	for i := range users {
		userPointers = append(userPointers, &users[i])
	}

	return userPointers, nil
}

func (s *UserService) GetUserByID(id int) (*data.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) CreateUser(user *data.User) (*data.User, error) {
	createdUser, err := s.userRepo.Create(user.Email, user.Password)
	if err != nil {
		return nil, err
	}
	return createdUser, nil
}
