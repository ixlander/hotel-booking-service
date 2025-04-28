package usecases

import (
	"errors"
	"time"
	
	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v4"
	
	"hotel-booking-service/internal/data"
	"hotel-booking-service/internal/repositories"
)

type AuthUsecase struct {
	userRepo *repositories.UserRepository
	jwtSecret string
	tokenExpiry time.Duration
}

func NewAuthUsecase(userRepo *repositories.UserRepository, jwtSecret string, tokenExpiry time.Duration) *AuthUsecase {
	return &AuthUsecase{
		userRepo: userRepo,
		jwtSecret: jwtSecret,
		tokenExpiry: tokenExpiry,
	}
}

func (uc *AuthUsecase) Register(email, password string) (*data.User, error) {
	existingUser, err := uc.userRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	
	if existingUser != nil {
		return nil, errors.New("user with this email already exists")
	}
	
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}
	
	user, err := uc.userRepo.Create(email, string(hashedPassword))
	if err != nil {
		return nil, err
	}
	
	return user, nil
}

func (uc *AuthUsecase) Login(email, password string) (*data.User, string, error) {
	user, err := uc.userRepo.FindByEmail(email)
	if err != nil {
		return nil, "", err
	}
	
	if user == nil {
		return nil, "", errors.New("user not found")
	}
	
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, "", errors.New("invalid credentials")
	}
	
	token, err := uc.generateJWT(user.ID)
	if err != nil {
		return nil, "", err
	}
	
	user.Password = ""
	
	return user, token, nil
}

func (uc *AuthUsecase) generateJWT(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(uc.tokenExpiry).Unix(),
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	return token.SignedString([]byte(uc.jwtSecret))
}