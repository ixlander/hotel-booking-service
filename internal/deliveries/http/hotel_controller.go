package deliveries

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"
	
	"github.com/gorilla/mux"
	
	"hotel-booking-service/internal/usecases"
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