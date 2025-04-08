package main

import (
	"log"

	"github.com/ixlander/hotel-booking-service/internal/app/config"
	"github.com/ixlander/hotel-booking-service/internal/app/connections"
	"github.com/ixlander/hotel-booking-service/internal/app/start"
	"github.com/ixlander/hotel-booking-service/internal/app/store"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	
	db, err := connections.NewPostgresConnection(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	
	appStore := store.NewStore(db)
	
	router := start.SetupRouter(appStore, cfg.JWT.Secret, int(cfg.JWT.TTL.Hours()))
	
	log.Printf("Server is running on port %s", cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}