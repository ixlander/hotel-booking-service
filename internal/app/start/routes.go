package start

import (
	"net/http"
	
	"github.com/gorilla/mux"
	
	"hotel-booking-service/internal/app/config"
	"hotel-booking-service/internal/app/store"
	"hotel-booking-service/internal/deliveries/http"
	"hotel-booking-service/internal/deliveries/http/middleware"
	"hotel-booking-service/internal/usecases"
)

func SetupRoutes(cfg *config.Config, store *store.Store) *mux.Router {
	router := mux.NewRouter()
	
	authUsecase := usecases.NewAuthUsecase(store.UserRepo, cfg.JWT.Secret, cfg.JWT.TokenExpiry)
	hotelUsecase := usecases.NewHotelUsecase(store.HotelRepo, store.RoomRepo)
	bookingUsecase := usecases.NewBookingUsecase(store.BookingRepo, store.RoomRepo)

	authController := deliveries.NewAuthController(authUsecase)
	hotelController := deliveries.NewHotelController(hotelUsecase)
	bookingController := deliveries.NewBookingController(bookingUsecase)
	
	auth := middleware.AuthMiddleware(cfg.JWT.Secret)
	
	router.HandleFunc("/register", authController.Register).Methods("POST")
	router.HandleFunc("/login", authController.Login).Methods("POST")
	
	router.HandleFunc("/hotels", hotelController.GetAllHotels).Methods("GET")
	router.HandleFunc("/hotels/{id}", hotelController.GetHotelByID).Methods("GET")
	
	bookingRouter := router.PathPrefix("/bookings").Subrouter()
	bookingRouter.Use(auth)
	bookingRouter.HandleFunc("", bookingController.CreateBooking).Methods("POST")
	bookingRouter.HandleFunc("", bookingController.GetUserBookings).Methods("GET")
	bookingRouter.HandleFunc("/{id}", bookingController.CancelBooking).Methods("DELETE")
	
	return router
}