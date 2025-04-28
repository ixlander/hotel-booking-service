package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"hotel-booking-service/internal/app/config"
	"hotel-booking-service/internal/app/connections"
)

func main() {
	upFlag := flag.Bool("up", false, "Run migrations up")
	downFlag := flag.Bool("down", false, "Roll back migrations")
	versionFlag := flag.Bool("version", false, "Show current migration version")
	forceFlag := flag.Int("force", -1, "Force migration to specific version")
	
	flag.Parse()
	
	cfg := config.Load()
	
	db, err := connections.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		log.Fatalf("Failed to create migration driver: %v", err)
	}
	
	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations",
		"postgres", driver)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}
	
	if *upFlag {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to run migrations: %v", err)
		}
		log.Println("Migrations applied successfully")
	} else if *downFlag {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to roll back migrations: %v", err)
		}
		log.Println("Migrations rolled back successfully")
	} else if *versionFlag {
		version, dirty, err := m.Version()
		if err != nil {
			log.Fatalf("Failed to get migration version: %v", err)
		}
		fmt.Printf("Current migration version: %d (dirty: %v)\n", version, dirty)
	} else if *forceFlag >= 0 {
		if err := m.Force(*forceFlag); err != nil {
			log.Fatalf("Failed to force migration version: %v", err)
		}
		log.Printf("Migration version forced to %d\n", *forceFlag)
	} else {
		fmt.Println("No action specified. Use -up, -down, -version, or -force.")
		os.Exit(1)
	}
}