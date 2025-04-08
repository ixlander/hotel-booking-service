package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/ixlander/hotel-booking-service/internal/data"
	"github.com/ixlander/hotel-booking-service/internal/usecases"
)

type BookingController struct {
	bookingUsecase *usecases.BookingUsecase
}

func NewBookingController(bookingUsecase *usecases.BookingUsecase) *BookingController {
	return &BookingController{
		bookingUsecase: bookingUsecase,
	}
}

func (c *BookingController) CreateBooking() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := GetUserID(ctx)
		if !exists {
			ctx.JSON(http.StatusUnauthorized, data.ApiError{Error: "unauthorized"})
			return
		}
		
		var req data.BookingRequest
		if err := ctx.ShouldBindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, data.ApiError{Error: err.Error()})
			return
		}

		if req.FromDate.After(req.ToDate) || req.FromDate.Equal(req.ToDate) {
			ctx.JSON(http.StatusBadRequest, data.ApiError{Error: "from_date must be before to_date"})
			return
		}
		
		booking, err := c.bookingUsecase.CreateBooking(ctx, userID, req.RoomID, req.FromDate, req.ToDate)
		if err != nil {
			status := http.StatusInternalServerError
			
			switch err {
			case usecases.ErrRoomNotFound:
				status = http.StatusNotFound
			case usecases.ErrRoomNotAvailable:
				status = http.StatusConflict
			case usecases.ErrInvalidDateRange:
				status = http.StatusBadRequest
			}
			
			ctx.JSON(status, data.ApiError{Error: err.Error()})
			return
		}
		
		ctx.JSON(http.StatusCreated, booking)
	}
}

func (c *BookingController) GetUserBookings() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := GetUserID(ctx)
		if !exists {
			ctx.JSON(http.StatusUnauthorized, data.ApiError{Error: "unauthorized"})
			return
		}
		
		bookings, err := c.bookingUsecase.GetUserBookings(ctx, userID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, data.ApiError{Error: err.Error()})
			return
		}
		
		ctx.JSON(http.StatusOK, gin.H{"bookings": bookings})
	}
}

func (c *BookingController) CancelBooking() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		userID, exists := GetUserID(ctx)
		if !exists {
			ctx.JSON(http.StatusUnauthorized, data.ApiError{Error: "unauthorized"})
			return
		}
		
		bookingIDParam := ctx.Param("id")
		bookingID, err := strconv.ParseInt(bookingIDParam, 10, 64)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, data.ApiError{Error: "invalid booking id"})
			return
		}
		
		err = c.bookingUsecase.CancelBooking(ctx, userID, bookingID)
		if err != nil {
			status := http.StatusInternalServerError
			
			switch err {
			case usecases.ErrBookingNotFound:
				status = http.StatusNotFound
			case usecases.ErrUnauthorized:
				status = http.StatusForbidden
			}
			
			ctx.JSON(status, data.ApiError{Error: err.Error()})
			return
		}
		
		ctx.JSON(http.StatusOK, gin.H{"message": "booking cancelled successfully"})
	}
}