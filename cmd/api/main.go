package main

import (
	"github.com/gorilla/mux"
	"ip-rate-control/internal/database"
	"ip-rate-control/internal/handler"
	"ip-rate-control/internal/middleware"
	"ip-rate-control/internal/redis"
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

	redisClient := redis.InitializeRedis(cfg.Redis)

	router := mux.NewRouter()
	router.Use(middleware.RateLimitMiddleware(redisClient, db))
	router.HandleFunc("/", handler.RootHandler).Methods("GET")

	log.Printf("Server is running on port %s", cfg.Server.Port)
	if err := http.ListenAndServe(":"+cfg.Server.Port, router); err != nil {
		log.Fatalf("Could not start server: %v\n", err)
	}
}
