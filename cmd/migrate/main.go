package main

import (
	"flag"
	"log"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/ixlander/hotel-booking-service/internal/app/config"
)

func main() {
	migrationsPath := flag.String("path", "migrations", "Path to migration files")
	up := flag.Bool("up", false, "Run migrations up")
	down := flag.Bool("down", false, "Run migrations down")
	flag.Parse()
	
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}
	
	dsn := "postgres://" + cfg.Database.User + ":" + cfg.Database.Password + "@" + 
		cfg.Database.Host + ":" + cfg.Database.Port + "/" + cfg.Database.DBName + 
		"?sslmode=" + cfg.Database.SSLMode
	
	
	m, err := migrate.New("file://"+*migrationsPath, dsn)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}
	
	if *up {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to apply migrations: %v", err)
		}
		log.Println("Migrations applied successfully")
	} else if *down {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatalf("Failed to rollback migrations: %v", err)
		}
		log.Println("Migrations rolled back successfully")
	} else {
		log.Println("No action specified. Use -up or -down")
	}
}