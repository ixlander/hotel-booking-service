package start

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	httpDelivery "github.com/ixlander/hotel-booking-service/internal/deliveries/http"
	httpMiddleware "github.com/ixlander/hotel-booking-service/internal/deliveries/http/middleware"
	"github.com/ixlander/hotel-booking-service/internal/usecases"
)

func NewRouter(
	authUsecase *usecases.AuthUsecase,
	hotelUsecase *usecases.HotelUsecase,
	bookingUsecase *usecases.BookingUsecase,
	jwtSecret string,
) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	authMiddleware := httpMiddleware.NewAuthMiddleware(jwtSecret)

	authController := httpDelivery.NewAuthController(authUsecase)
	hotelController := httpDelivery.NewHotelController(hotelUsecase)
	bookingController := httpDelivery.NewBookingController(bookingUsecase)

	r.Group(func(r chi.Router) {
		authController.RegisterRoutes(r)
		r.Get("/hotels", hotelController.GetAllHotels)
		r.Get("/hotels/{id}/available", hotelController.GetHotelWithAvailableRooms)
	})

	r.Group(func(r chi.Router) {
		r.Use(authMiddleware.Middleware)
		bookingController.RegisterRoutes(r)
	})

	return r
}