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
	ErrUserNotFound      = errors.New("user not found")
	ErrInvalidCredentials = errors.New("invalid credentials")
)

type AuthUsecase struct {
	userRepo repositories.UserRepository
	jwtSecret string
	jwtTTL    time.Duration
}

func NewAuthUsecase(
	userRepo repositories.UserRepository,
	jwtSecret string,
	jwtTTL time.Duration,
) *AuthUsecase {
	return &AuthUsecase{
		userRepo: userRepo,
		jwtSecret: jwtSecret,
		jwtTTL: jwtTTL,
	}
}

type JWTClaims struct {
	UserID int64 `json:"user_id"`
	jwt.RegisteredClaims
}

func (u *AuthUsecase) Login(ctx context.Context, email, password string) (string, *data.User, error) {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		return "", nil, err
	}
	
	if user == nil {
		return "", nil, ErrUserNotFound
	}
	
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", nil, ErrInvalidCredentials
	}
	
	token, err := u.generateToken(user.ID)
	if err != nil {
		return "", nil, err
	}
	
	user.Password = ""
	
	return token, user, nil
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
	
	createdUser, err := u.userRepo.Create(ctx, user)
	if err != nil {
		return nil, err
	}
	
	createdUser.Password = ""
	
	return createdUser, nil
}

func (u *AuthUsecase) VerifyToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(u.jwtSecret), nil
	})
	
	if err != nil {
		return nil, err
	}
	
	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}
	
	return nil, errors.New("invalid token")
}

func (u *AuthUsecase) generateToken(userID int64) (string, error) {
	claims := &JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(u.jwtTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	
	return token.SignedString([]byte(u.jwtSecret))
}