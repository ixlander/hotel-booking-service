package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/yourusername/hotel-booking-service/internal/data"
	"github.com/yourusername/hotel-booking-service/internal/usecases"
)

type HotelController struct {
	hotelUsecase *usecases.HotelUsecase
}

func NewHotelController(hotelUsecase *usecases.HotelUsecase) *HotelController {
	return &HotelController{
		hotelUsecase: hotelUsecase,
	}
}

func (c *HotelController) GetAllHotels() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		hotels, err := c.hotelUsecase.GetAllHotels(ctx)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, data.ApiError{Error: err.Error()})
			return
		}
		
		ctx.JSON(http.StatusOK, data.HotelResponse{Hotels: hotels})
	}
}

func (c *HotelController) GetHotelWithRooms() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		id, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, data.ApiError{Error: "invalid hotel id"})
			return
		}
		
		hotel, err := c.hotelUsecase.GetHotelWithRooms(ctx, id)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, data.ApiError{Error: err.Error()})
			return
		}
		
		if hotel == nil {
			ctx.JSON(http.StatusNotFound, data.ApiError{Error: "hotel not found"})
			return
		}
		
		ctx.JSON(http.StatusOK, hotel)
	}
}

func (c *HotelController) GetAvailableRooms() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		idParam := ctx.Param("id")
		hotelID, err := strconv.ParseInt(idParam, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, data.ApiError{Error: "invalid hotel id"})
			return
		}
		
		fromDateStr := ctx.Query("from_date")
		toDateStr := ctx.Query("to_date")
		
		if fromDateStr == "" || toDateStr == "" {
			ctx.JSON(http.StatusBadRequest, data.ApiError{Error: "from_date and to_date query parameters are required"})
			return
		}
		
		fromDate, err := time.Parse("2006-01-02", fromDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, data.ApiError{Error: "invalid from_date format, use YYYY-MM-DD"})
			return
		}
		
		toDate, err := time.Parse("2006-01-02", toDateStr)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, data.ApiError{Error: "invalid to_date format, use YYYY-MM-DD"})
			return
		}
		
		if fromDate.After(toDate) || fromDate.Equal(toDate) {
			ctx.JSON(http.StatusBadRequest, data.ApiError{Error: "from_date must be before to_date"})
			return
		}
		
		rooms, err := c.hotelUsecase.GetAvailableRooms(ctx, hotelID, fromDate, toDate)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, data.ApiError{Error: err.Error()})
			return
		}
		
		ctx.JSON(http.StatusOK, gin.H{"rooms": rooms})
	}
}