package start

import (
	"github.com/gin-gonic/gin"

	"github.com/ixlander/hotel-booking-service/internal/app/store"
	"github.com/ixlander/hotel-booking-service/internal/deliveries/http"
	"github.com/ixlander/hotel-booking-service/internal/usecases"
)

func SetupRouter(store *store.Store, jwtSecret string, jwtTTL int) *gin.Engine {
	router := gin.Default()
	
	authUsecase := usecases.NewAuthUsecase(store.UserRepo, jwtSecret, jwtTTL)
	hotelUsecase := usecases.NewHotelUsecase(store.HotelRepo, store.RoomRepo)
	bookingUsecase := usecases.NewBookingUsecase(store.BookingRepo, store.RoomRepo)
	
	authController := http.NewAuthController(authUsecase)
	hotelController := http.NewHotelController(hotelUsecase)
	bookingController := http.NewBookingController(bookingUsecase)
	
	router.POST("/register", authController.Register())
	router.POST("/login", authController.Login())     
	
	router.GET("/hotels", hotelController.GetAllHotels())           
	router.GET("/hotels/:id", hotelController.GetHotelWithRooms())   
	router.GET("/hotels/:id/available-rooms", hotelController.GetAvailableRooms()) 
	
	authorized := router.Group("/")
	authorized.Use(authController.AuthMiddleware())
	{
		authorized.POST("/bookings", bookingController.CreateBooking())   
		authorized.GET("/bookings", bookingController.GetUserBookings()) 
		authorized.DELETE("/bookings/:id", bookingController.CancelBooking())
	}
	
	return router
}
