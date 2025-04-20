package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"

	"github.com/ixlander/hotel-booking-service/internal/data"
	"github.com/ixlander/hotel-booking-service/internal/repositories"
)

var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserNotFound       = errors.New("user not found")
)

type AuthUsecase struct {
	userRepo    repositories.UserRepository
	jwtSecret   []byte
	jwtDuration time.Duration
}

func NewAuthUsecase(userRepo repositories.UserRepository, jwtSecret string, jwtDuration time.Duration) *AuthUsecase {
	return &AuthUsecase{
		userRepo:    userRepo,
		jwtSecret:   []byte(jwtSecret),
		jwtDuration: jwtDuration,
	}
}

func (u *AuthUsecase) Register(ctx context.Context, email, password string) (*data.User, error) {
	existingUser, err := u.userRepo.FindByEmail(ctx, email)
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

	user := &data.User{
		Email:    email,
		Password: string(hashedPassword),
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return nil, err
	}

	user.Password = ""
	return user, nil
}

func (u *AuthUsecase) Login(ctx context.Context, email, password string) (*data.LoginResponse, error) {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, ErrUserNotFound
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	token, err := u.generateJWT(user.ID)
	if err != nil {
		return nil, err
	}

	user.Password = ""

	return &data.LoginResponse{
		Token: token,
		User:  *user,
	}, nil
}

func (u *AuthUsecase) generateJWT(userID int64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(u.jwtDuration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(u.jwtSecret)
}