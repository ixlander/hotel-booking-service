package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"

	"github.com/ixlander/hotel-booking-service/internal/pkg/httputil"
	"github.com/ixlander/hotel-booking-service/internal/usecases"
)

type HotelController struct {
	hotelUsecase *usecases.HotelUsecase
}

func NewHotelController(hotelUsecase *usecases.HotelUsecase) *HotelController {
	return &HotelController{
		hotelUsecase: hotelUsecase,
	}
}

func (c *HotelController) RegisterRoutes(r chi.Router) {
	r.Get("/hotels", c.GetAllHotels)
	r.Get("/hotels/{id}/available", c.GetHotelWithAvailableRooms)
}

func (c *HotelController) GetAllHotels(w http.ResponseWriter, r *http.Request) {
	hotels, err := c.hotelUsecase.GetAllHotels(r.Context())
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputil.WriteJSON(w, http.StatusOK, hotels)
}

func (c *HotelController) GetHotelWithAvailableRooms(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid hotel ID")
		return
	}

	fromDateStr := r.URL.Query().Get("from_date")
	toDateStr := r.URL.Query().Get("to_date")

	if fromDateStr == "" || toDateStr == "" {
		httputil.WriteError(w, http.StatusBadRequest, "from_date and to_date parameters are required")
		return
	}

	fromDate, err := time.Parse("2006-01-02", fromDateStr)
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid from_date format")
		return
	}

	toDate, err := time.Parse("2006-01-02", toDateStr)
	if err != nil {
		httputil.WriteError(w, http.StatusBadRequest, "invalid to_date format")
		return
	}

	hotel, err := c.hotelUsecase.GetHotelWithAvailableRooms(r.Context(), id, fromDate, toDate)
	if err != nil {
		httputil.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if hotel == nil {
		httputil.WriteError(w, http.StatusNotFound, "hotel not found")
		return
	}

	httputil.WriteJSON(w, http.StatusOK, hotel)
}