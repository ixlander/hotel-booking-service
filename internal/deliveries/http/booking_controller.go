package http

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"

	"github.com/ixlander/hotel-booking-service/internal/data"
	"github.com/ixlander/hotel-booking-service/internal/deliveries/http/middleware"
	"github.com/ixlander/hotel-booking-service/internal/pkg/httputil"
	"github.com/ixlander/hotel-booking-service/internal/usecases"
)

type BookingController struct {
	bookingUsecase *usecases.BookingUsecase
	validate       *validator.Validate
}

func NewBookingController(bookingUsecase *usecases.BookingUsecase) *BookingController {
	return &BookingController{
		bookingUsecase: bookingUsecase,
		validate:       validator.New(),
	}
}

func (c *BookingController) RegisterRoutes(r chi.Router) {
	r.Post("/bookings", c.CreateBooking)
	r.Get("/bookings", c.GetUserBookings)
	r.Delete("/bookings/{id}", c.CancelBooking)
}

func (c *BookingController) CreateBooking(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		httputil.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	var req data.BookingRequest

	if err := httputil.DecodeJSON(r, &req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid request body")
		return
	}

	if err := c.validate.Struct(req); err != nil {
		httputil.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	booking, err := c.bookingUsecase.CreateBooking(r.Context(), userID, req.RoomID, req.FromDate, req.ToDate)
	if err != nil {
		switch {
		case errors.Is(err, usecases.ErrRoomNotAvailable):
			httputil.WriteError(w, http.StatusConflict, err.Error())
		case errors.Is(err, usecases.ErrInvalidDateRange):
			httputil.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, usecases.ErrPastBooking):
			httputil.WriteError(w, http.StatusBadRequest, err.Error())
		default:
			httputil.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	httputil.WriteJSON(w, http.StatusCreated, booking)
}

func (c *BookingController) GetUserBookings(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		httputil.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	bookings, err := c.bookingUsecase.GetUserBookings(r.Context(), userID)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusOK, bookings)
}

func (c *BookingController) CancelBooking(w http.ResponseWriter, r *http.Request) {
	userID, ok := middleware.GetUserID(r.Context())
	if !ok {
		httputil.WriteError(w, http.StatusUnauthorized, "unauthorized")
		return
	}

	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid booking ID")
		return
	}

	err = c.bookingUsecase.CancelBooking(r.Context(), userID, id)
	if err != nil {
		switch {
		case errors.Is(err, usecases.ErrBookingNotFound):
			httputil.WriteError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, usecases.ErrBookingNotBelongsUser):
			httputil.WriteError(w, http.StatusForbidden, err.Error())
		default:
			httputil.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	httputil.WriteJSON(w, http.StatusOK, map[string]interface{}{
		"message": "Booking cancelled successfully",
	})
}