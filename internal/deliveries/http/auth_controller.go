package http

import (
	"errors"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"github.com/ixlander/hotel-booking-service/internal/data"
	"github.com/ixlander/hotel-booking-service/internal/pkg/httputil"
	"github.com/ixlander/hotel-booking-service/internal/usecases"
)

type AuthController struct {
	authUsecase *usecases.AuthUsecase
	validate    *validator.Validate
}

func NewAuthController(authUsecase *usecases.AuthUsecase) *AuthController {
	return &AuthController{
		authUsecase: authUsecase,
		validate:    validator.New(),
	}
}

func (c *AuthController) RegisterRoutes(r chi.Router) {
	r.Post("/register", c.Register)
	r.Post("/login", c.Login)
}

func (c *AuthController) Register(w http.ResponseWriter, r *http.Request) {
	var req data.LoginRequest

	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := c.validate.Struct(req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	user, err := c.authUsecase.Register(r.Context(), req.Email, req.Password)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, user)
}

func (c *AuthController) Login(w http.ResponseWriter, r *http.Request) {
	var req data.LoginRequest

	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := c.validate.Struct(req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response, err := c.authUsecase.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if errors.Is(err, usecases.ErrUserNotFound) {
			httputil.WriteError(w, http.StatusNotFound, "user not found")
			return
		}
		if errors.Is(err, usecases.ErrInvalidCredentials) {
			httputil.WriteError(w, http.StatusUnauthorized, "invalid credentials")
			return
		}
		httputil.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusOK, response)
}