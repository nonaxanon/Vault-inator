package main

import (
	"log"
	"net/http"
	"os"

	"github.com/nonaxanon/vault-inator/internal/api"
	"github.com/nonaxanon/vault-inator/internal/config"
	"github.com/nonaxanon/vault-inator/internal/services"
	"github.com/nonaxanon/vault-inator/internal/storage"
)

func main() {
	// Load configuration
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// Initialize services
	passwordService := services.NewPasswordService()
	authService := services.NewAuthService(passwordService)

	// Get database connection string from environment variable
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://admin:admin@192.168.1.8:5432/dev?sslmode=disable"
	}

	// Create database connection
	db, err := storage.NewDB(connStr, cfg.MasterPassword)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize database
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}

	// Create API server with services
	server := api.NewServer(db, authService, passwordService)

	// Start server
	log.Printf("Starting server on :8080")
	if err := http.ListenAndServe(":8080", server); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
