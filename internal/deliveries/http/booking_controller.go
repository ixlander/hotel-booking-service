package deliveries

import (
	"encoding/json"
	"fmt"
	"hotel-booking-service/internal/data"
	"hotel-booking-service/internal/usecases"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

const userIDContextKey = "userID"

type BookingController struct {
	bookingUsecase *usecases.BookingUsecase
}

func NewBookingController(bookingUsecase *usecases.BookingUsecase) *BookingController {
	return &BookingController{
		bookingUsecase: bookingUsecase,
	}
}

func (c *BookingController) CreateBooking(w http.ResponseWriter, r *http.Request) {

	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in CreateBooking: %v", r)
			sendErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
	}()

	userIDValue := r.Context().Value(userIDContextKey)

	if userIDValue == nil {
		log.Printf("userID not found in context")
		sendErrorResponse(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}

	userID, ok := userIDValue.(int)
	if !ok {
		log.Printf("userID is not an int: %T", userIDValue)
		sendErrorResponse(w, "Invalid user ID format", http.StatusInternalServerError)
		return
	}

	log.Printf("Processing booking creation for user ID: %d", userID)

	var req data.CreateBookingRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Error decoding request body: %v", err)
		sendErrorResponse(w, "Invalid request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	log.Printf("Booking request: Room ID: %d, From: %s, To: %s", 
		req.RoomID, req.FromDate.Format(time.RFC3339), req.ToDate.Format(time.RFC3339))

	booking, err := c.bookingUsecase.CreateBooking(userID, req.RoomID, req.FromDate, req.ToDate)
	if err != nil {
		log.Printf("Error creating booking: %v", err)
		if err.Error() == "room not available for the selected dates" {
			sendErrorResponse(w, err.Error(), http.StatusConflict)
			return
		}
		sendErrorResponse(w, fmt.Sprintf("Booking creation failed: %v", err), http.StatusBadRequest)
		return
	}

	log.Printf("Successfully created booking with ID: %d", booking.ID)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"message": "Booking created successfully",
		"booking_id": booking.ID,
	})
}

func (c *BookingController) CancelBooking(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in CancelBooking: %v", r)
			sendErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
	}()

	userIDValue := r.Context().Value(userIDContextKey)
	if userIDValue == nil {
		sendErrorResponse(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	
	userID, ok := userIDValue.(int)
	if !ok {
		sendErrorResponse(w, "Invalid user ID format", http.StatusInternalServerError)
		return
	}
	
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
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in GetBookingByID: %v", r)
			sendErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
	}()

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
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in UpdateBooking: %v", r)
			sendErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
	}()

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
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Recovered from panic in GetUserBookings: %v", r)
			sendErrorResponse(w, "Internal server error", http.StatusInternalServerError)
		}
	}()

	userIDValue := r.Context().Value(userIDContextKey)
	if userIDValue == nil {
		sendErrorResponse(w, "User ID not found in context", http.StatusUnauthorized)
		return
	}
	
	userID, ok := userIDValue.(int)
	if !ok {
		sendErrorResponse(w, "Invalid user ID format", http.StatusInternalServerError)
		return
	}
	
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