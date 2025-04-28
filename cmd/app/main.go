package main

import (
	"fmt"
	"log"
	"net/http"
	
	"hotel-booking-service/internal/app/config"
	"hotel-booking-service/internal/app/connections"
	"hotel-booking-service/internal/app/start"
	"hotel-booking-service/internal/app/store"
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
	log.Println("Connected to database")
	
	appStore := store.NewStore(db)
	
	router := start.SetupRoutes(cfg, appStore)
	
	addr := fmt.Sprintf(":%s", cfg.Server.Port)
	log.Printf("Starting server on %s", addr)
	if err := http.ListenAndServe(addr, router); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}