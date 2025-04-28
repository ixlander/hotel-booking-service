package deliveries

import (
	"encoding/json"
	"net/http"
	"strconv"
	
	"github.com/gorilla/mux"
	
	"hotel-booking-service/internal/data"
	"hotel-booking-service/internal/usecases"
)

type BookingController struct {
	bookingUsecase *usecases.BookingUsecase
}

func NewBookingController(bookingUsecase *usecases.BookingUsecase) *BookingController {
	return &BookingController{
		bookingUsecase: bookingUsecase,
	}
}

func (c *BookingController) CreateBooking(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	
	var req data.CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	booking, err := c.bookingUsecase.CreateBooking(userID, req.RoomID, req.FromDate, req.ToDate)
	if err != nil {
		if err.Error() == "room not available for the selected dates" {
			http.Error(w, err.Error(), http.StatusConflict)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(booking)
}

func (c *BookingController) CancelBooking(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	
	vars := mux.Vars(r)
	bookingID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}
	
	err = c.bookingUsecase.CancelBooking(userID, bookingID)
	if err != nil {
		if err.Error() == "booking not found" {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}
		if err.Error() == "booking does not belong to this user" {
			http.Error(w, err.Error(), http.StatusForbidden)
			return
		}
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Booking cancelled successfully"})
}


func (c *BookingController) GetUserBookings(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	
	bookings, err := c.bookingUsecase.GetUserBookings(userID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}