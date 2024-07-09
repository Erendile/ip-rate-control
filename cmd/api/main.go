package main

import (
	"ip-rate-control/internal/database"
	"ip-rate-control/internal/handler"
	"ip-rate-control/internal/middleware"
	"ip-rate-control/pkg/config"
	"log"
	"net/http"
)

func main() {
	cfg := config.NewConfiguration()
	db, err := database.NewPostgresDB(cfg.Database)
	if err != nil {
		log.Fatalf("Error initializing database: %v", err)
	}

	mux := http.NewServeMux()
	fs := http.FileServer(http.Dir("static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	mux.Handle("/", middleware.RateLimitMiddleware(db, http.HandlerFunc(handler.RootHandler)))

	log.Printf("Server is running on port %s", cfg.Server.Port)
	if err := http.ListenAndServe(":"+cfg.Server.Port, mux); err != nil {
		log.Fatalf("Could not start server: %v\n", err)
	}
}
