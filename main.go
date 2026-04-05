package main

import (
	"context"
	"log"
	"net/http"

	"github.com/ravdreamin/core-ledger-api/config"
	"github.com/ravdreamin/core-ledger-api/handlers"
	"github.com/ravdreamin/core-ledger-api/middleware"
	
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	// 1. Load the environment variables
	cfg := config.LoadConfig()

	// 2. Connect to the database pool
	pool, err := pgxpool.New(context.Background(), cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	defer pool.Close()

	log.Println("Successfully connected to the database!")
	
	// 3. Register Routes
	http.HandleFunc("POST /login", handlers.LoginHandler(pool, cfg.PasetoKey))
	
	http.HandleFunc("GET /dashboard/summary", middleware.RequireAuth(cfg.PasetoKey, 
		middleware.RequireRole(handlers.DashboardSummaryHandler(pool), "admin", "analyst", "viewer"),
	))

	http.HandleFunc("GET /dashboard/categories", middleware.RequireAuth(cfg.PasetoKey, 
		middleware.RequireRole(handlers.CategoryReportHandler(pool), "admin", "analyst"),
	))

	http.HandleFunc("POST /records", middleware.RequireAuth(cfg.PasetoKey, 
		middleware.RequireRole(handlers.CreateRecordHandler(pool), "admin"),
	))

	// 4. Start the HTTP Server
	log.Printf("Starting server on port %s...\n", cfg.ServerPort)
	if err := http.ListenAndServe(":"+cfg.ServerPort, nil); err != nil {
		log.Fatalf("Server crashed: %v\n", err)
	}
}