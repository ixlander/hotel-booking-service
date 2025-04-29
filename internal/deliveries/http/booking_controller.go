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
		sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	
	booking, err := c.bookingUsecase.CreateBooking(userID, req.RoomID, req.FromDate, req.ToDate)
	if err != nil {
		if err.Error() == "room not available for the selected dates" {
			sendErrorResponse(w, err.Error(), http.StatusConflict)
			return
		}
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Booking created successfully",
		"booking_id": booking.ID,
	})
}

func (c *BookingController) CancelBooking(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	
	vars := mux.Vars(r)
	bookingID, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendErrorResponse(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}
	
	err = c.bookingUsecase.CancelBooking(userID, bookingID)
	if err != nil {
		switch err.Error() {
		case "booking not found":
			sendErrorResponse(w, err.Error(), http.StatusNotFound)
		case "booking does not belong to this user":
			sendErrorResponse(w, err.Error(), http.StatusForbidden)
		default:
			sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		}
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Booking cancelled successfully"})
}

func (c *BookingController) GetBookingByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookingID, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendErrorResponse(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	booking, err := c.bookingUsecase.GetBookingByID(bookingID)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusNotFound)
		return
	}

	if booking == nil {
		sendErrorResponse(w, "Booking not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(booking)
}

func (c *BookingController) UpdateBooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookingID, err := strconv.Atoi(vars["id"])
	if err != nil {
		sendErrorResponse(w, "Invalid booking ID", http.StatusBadRequest)
		return
	}

	var req data.UpdateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	err = c.bookingUsecase.UpdateBooking(bookingID, req.Status)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Booking updated successfully"})
}

func (c *BookingController) GetUserBookings(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(int)
	
	bookings, err := c.bookingUsecase.GetUserBookings(userID)
	if err != nil {
		sendErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookings)
}

func sendErrorResponse(w http.ResponseWriter, message string, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	response := map[string]string{"message": message}
	json.NewEncoder(w).Encode(response)
}