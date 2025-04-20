package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ixlander/hotel-booking-service/internal/app/config"
	"github.com/ixlander/hotel-booking-service/internal/app/connections"
	"github.com/ixlander/hotel-booking-service/internal/app/start"
	"github.com/ixlander/hotel-booking-service/internal/app/store"
	"github.com/ixlander/hotel-booking-service/internal/usecases"
)

func main() {
	cfg := config.Load()

	db, err := connections.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	s := store.NewStore(db)

	authUsecase := usecases.NewAuthUsecase(s.UserRepo, cfg.JWT.Secret, cfg.JWT.Duration)
	hotelUsecase := usecases.NewHotelUsecase(s.HotelRepo, s.RoomRepo)
	bookingUsecase := usecases.NewBookingUsecase(s.RoomRepo, s.BookingRepo)

	router := start.NewRouter(authUsecase, hotelUsecase, bookingUsecase, cfg.JWT.Secret)

	server := &http.Server{
		Addr:         ":" + cfg.Server.Port,
		Handler:      router,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited properly")
}
