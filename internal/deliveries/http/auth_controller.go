package deliveries

import (
	"encoding/json"
	"net/http"
	"time"
	
	"hotel-booking-service/internal/data"
	"hotel-booking-service/internal/usecases"
)

type AuthController struct {
	authUsecase *usecases.AuthUsecase
}

func NewAuthController(authUsecase *usecases.AuthUsecase) *AuthController {
	return &AuthController{
		authUsecase: authUsecase,
	}
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var req data.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}
	
	user, err := c.authUsecase.Register(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(user)
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req data.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	if req.Email == "" || req.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}
	
	user, token, err := c.authUsecase.Login(req.Email, req.Password)
	if err != nil {
		if err.Error() == "user not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data.LoginResponse{
		Token: token,
		User:  *user,
	})
}