package usecases

import (
	"errors"
	"hotel-booking-service/internal/data"
	"hotel-booking-service/internal/repositories"
)

type UserUsecase struct {
	userRepo *repositories.UserRepository
}

func NewUserUsecase(userRepo *repositories.UserRepository) *UserUsecase {
	return &UserUsecase{userRepo: userRepo}
}

func (uc *UserUsecase) GetUserByID(id int) (*data.User, error) {
	user, err := uc.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (uc *UserUsecase) GetAllUsers() ([]data.User, error) {
	return uc.userRepo.GetAllUsers()
}

func (uc *UserUsecase) UpdateUser(user *data.User) error {
	return uc.userRepo.Update(user)
}

// Удаление пользователя
func (uc *UserUsecase) DeleteUser(id int) error {
	return uc.userRepo.Delete(id)
}