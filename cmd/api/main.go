package main

import (
	"github.com/gorilla/mux"
	"ip-rate-control/internal/handler"
	"ip-rate-control/internal/middleware"
	"ip-rate-control/pkg/config"
	"log"
	"net/http"
)

func main() {
	cfg := config.NewConfiguration()

	router := mux.NewRouter()
	router.Use(middleware.RateLimitMiddleware())
	router.HandleFunc("/", handler.RootHandler).Methods("GET")

	log.Printf("Server is running on port %s", cfg.Server.Port)
	if err := http.ListenAndServe(":"+cfg.Server.Port, router); err != nil {
		log.Fatalf("Could not start server: %v\n", err)
	}
}
