package start

import (
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
	userUsecase := usecases.NewUserUsecase(store.UserRepo) 

	authController := deliveries.NewAuthController(authUsecase)
	hotelController := deliveries.NewHotelController(hotelUsecase)
	bookingController := deliveries.NewBookingController(bookingUsecase)
	userController := deliveries.NewUserController(userUsecase, cfg.JWT.Secret)

	auth := middleware.AuthMiddleware(cfg.JWT.Secret)

	router.HandleFunc("/register", authController.Register).Methods("POST")
	router.HandleFunc("/login", authController.Login).Methods("POST")
	router.HandleFunc("/hotels", hotelController.GetAllHotels).Methods("GET")
	router.HandleFunc("/hotels/{id:[0-9]+}", hotelController.GetHotelByID).Methods("GET")
	router.HandleFunc("/hotels/{id:[0-9]+}/rooms", hotelController.GetHotelRooms).Methods("GET")

	api := router.PathPrefix("/api").Subrouter()
	api.Use(auth)

	api.HandleFunc("/users/me", userController.GetCurrentUser).Methods("GET")
	api.HandleFunc("/users/{id:[0-9]+}", userController.UpdateUser).Methods("PUT")
	api.HandleFunc("/users/{id:[0-9]+}", userController.DeleteUser).Methods("DELETE")

	api.HandleFunc("/bookings", bookingController.CreateBooking).Methods("POST")
	api.HandleFunc("/bookings", bookingController.GetUserBookings).Methods("GET")
	api.HandleFunc("/bookings/{id:[0-9]+}", bookingController.GetBookingByID).Methods("GET")
	api.HandleFunc("/bookings/{id:[0-9]+}", bookingController.CancelBooking).Methods("DELETE")
	api.HandleFunc("/bookings/{id:[0-9]+}", bookingController.UpdateBooking).Methods("PUT")

	api.HandleFunc("/hotels", hotelController.CreateHotel).Methods("POST")
	api.HandleFunc("/hotels/{id:[0-9]+}", hotelController.UpdateHotel).Methods("PUT")
	api.HandleFunc("/hotels/{id:[0-9]+}", hotelController.DeleteHotel).Methods("DELETE")

	api.HandleFunc("/hotels/{hotelID:[0-9]+}/rooms", hotelController.CreateRoom).Methods("POST")
	api.HandleFunc("/rooms/{id:[0-9]+}", hotelController.UpdateRoom).Methods("PUT")
	api.HandleFunc("/rooms/{id:[0-9]+}", hotelController.DeleteRoom).Methods("DELETE")

	return router
}
