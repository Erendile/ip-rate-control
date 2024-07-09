package middleware

import (
	"database/sql"
	"ip-rate-control/pkg/ip"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

const maxRequestsPerHour = 10

func RateLimitMiddleware(db *sql.DB) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ipAddress := ip.GetClientIP(r)
			log.Println(ipAddress)
			now := time.Now().In(time.Local)
			log.Println("Current Time (Local):", now)

			var requestCount int
			var lastRequest time.Time

			err := db.QueryRow("SELECT request_count, last_request FROM ip_requests WHERE ip_address = $1", ipAddress).Scan(&requestCount, &lastRequest)
			if err != nil {
				if err == sql.ErrNoRows {
					_, err = db.Exec("INSERT INTO ip_requests (ip_address, request_count, last_request) VALUES ($1, $2, $3)",
						ipAddress, 1, now)
					if err != nil {
						http.Error(w, "Internal Server Error", http.StatusInternalServerError)
						return
					}
				} else {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			} else {
				log.Println("Last Request (Local):", lastRequest)
				duration := now.Sub(lastRequest)
				log.Println("Duration since last request (hours):", duration.Hours())
				if duration.Hours() < 1 {
					if requestCount >= maxRequestsPerHour {
						http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
						return
					}

					_, err = db.Exec("UPDATE ip_requests SET request_count = request_count + 1, last_request = $1 WHERE ip_address = $2", now, ipAddress)
					if err != nil {
						http.Error(w, "Internal Server Error", http.StatusInternalServerError)
						return
					}
				} else {
					_, err = db.Exec("UPDATE ip_requests SET request_count = 1, last_request = $1 WHERE ip_address = $2", now, ipAddress)
					if err != nil {
						http.Error(w, "Internal Server Error", http.StatusInternalServerError)
						return
					}
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
