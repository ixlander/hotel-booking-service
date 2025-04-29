package deliveries

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	
	"github.com/gorilla/mux"
	
	"hotel-booking-service/internal/usecases"
	"hotel-booking-service/internal/data" 
)

type HotelController struct {
	hotelUsecase *usecases.HotelUsecase
}

func NewHotelController(hotelUsecase *usecases.HotelUsecase) *HotelController {
	return &HotelController{
		hotelUsecase: hotelUsecase,
	}
}

func (c *HotelController) GetAllHotels(w http.ResponseWriter, r *http.Request) {
	fromDateStr := r.URL.Query().Get("from_date")
	toDateStr := r.URL.Query().Get("to_date")
	
	fromDate := time.Now()
	toDate := fromDate.AddDate(0, 0, 1)
	
	if fromDateStr != "" {
		parsedFromDate, err := time.Parse("2006-01-02", fromDateStr)
		if err == nil {
			fromDate = parsedFromDate
		}
	}
	
	if toDateStr != "" {
		parsedToDate, err := time.Parse("2006-01-02", toDateStr)
		if err == nil {
			toDate = parsedToDate
		}
	}
	
	hotels, err := c.hotelUsecase.GetAllHotels(fromDate, toDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hotels)
}

func (c *HotelController) GetHotelByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hotelID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid hotel ID", http.StatusBadRequest)
		return
	}
	
	fromDateStr := r.URL.Query().Get("from_date")
	toDateStr := r.URL.Query().Get("to_date")
	
	fromDate := time.Now()
	toDate := fromDate.AddDate(0, 0, 1)
	
	if fromDateStr != "" {
		parsedFromDate, err := time.Parse("2006-01-02", fromDateStr)
		if err == nil {
			fromDate = parsedFromDate
		}
	}
	
	if toDateStr != "" {
		parsedToDate, err := time.Parse("2006-01-02", toDateStr)
		if err == nil {
			toDate = parsedToDate
		}
	}
	
	hotel, err := c.hotelUsecase.GetHotelByID(hotelID, fromDate, toDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	if hotel == nil {
		http.Error(w, "Hotel not found", http.StatusNotFound)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(hotel)
}

func (c *HotelController) GetHotelRooms(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hotelID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid hotel ID", http.StatusBadRequest)
		return
	}

	fromDateStr := r.URL.Query().Get("from_date")
	toDateStr := r.URL.Query().Get("to_date")

	fromDate := time.Now()
	toDate := fromDate.AddDate(0, 0, 1)

	if fromDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", fromDateStr); err == nil {
			fromDate = parsed
		}
	}

	if toDateStr != "" {
		if parsed, err := time.Parse("2006-01-02", toDateStr); err == nil {
			toDate = parsed
		}
	}

	rooms, err := c.hotelUsecase.GetRoomsByHotelID(hotelID, fromDate, toDate)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rooms)
}

func (c *HotelController) CreateHotel(w http.ResponseWriter, r *http.Request) {
	var hotel data.Hotel
	if err := json.NewDecoder(r.Body).Decode(&hotel); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdHotel, err := c.hotelUsecase.CreateHotel(hotel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdHotel)
}

func (c *HotelController) UpdateHotel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hotelID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid hotel ID", http.StatusBadRequest)
		return
	}

	var hotel data.Hotel
	if err := json.NewDecoder(r.Body).Decode(&hotel); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	hotel.ID = hotelID

	updatedHotel, err := c.hotelUsecase.UpdateHotel(hotel)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedHotel)
}

func (c *HotelController) DeleteHotel(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	hotelID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid hotel ID", http.StatusBadRequest)
		return
	}

	err = c.hotelUsecase.DeleteHotel(hotelID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (c *HotelController) CreateRoom(w http.ResponseWriter, r *http.Request) {
	var room data.Room
	if err := json.NewDecoder(r.Body).Decode(&room); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdRoom, err := c.hotelUsecase.CreateRoom(room)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(createdRoom)
}

func (c *HotelController) UpdateRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}

	var room data.Room
	if err := json.NewDecoder(r.Body).Decode(&room); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	room.ID = roomID

	updatedRoom, err := c.hotelUsecase.UpdateRoom(room)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(updatedRoom)
}

func (c *HotelController) DeleteRoom(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	roomID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid room ID", http.StatusBadRequest)
		return
	}

	err = c.hotelUsecase.DeleteRoom(roomID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}