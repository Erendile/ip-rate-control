package middleware

import (
	"database/sql"
	"log"
	"net/http"
	"time"
)

const maxRequestsPerHour = 10

func RateLimitMiddleware(db *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ipAddress := r.RemoteAddr
		now := time.Now()

		var requestCount int
		var lastRequest time.Time

		log.Printf("Received request from IP: %s", ipAddress)

		err := db.QueryRow("SELECT request_count, last_request FROM ip_requests WHERE ip_address = $1", ipAddress).Scan(&requestCount, &lastRequest)
		if err != nil {
			if err == sql.ErrNoRows {
				log.Printf("No record found for IP: %s, creating new entry.", ipAddress)
				_, err = db.Exec("INSERT INTO ip_requests (ip_address, request_count, last_request) VALUES ($1, $2, $3)",
					ipAddress, 1, now)
				if err != nil {
					log.Printf("Error inserting new request: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
			} else {
				log.Printf("Error querying request count: %v", err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		} else {
			log.Printf("Record found for IP: %s, request_count: %d, last_request: %s", ipAddress, requestCount, lastRequest)
			if now.Sub(lastRequest).Hours() < 1 {
				if requestCount >= maxRequestsPerHour {
					log.Printf("Rate limit exceeded for IP: %s", ipAddress)
					http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
					return
				}

				_, err = db.Exec("UPDATE ip_requests SET request_count = request_count + 1, last_request = $1 WHERE ip_address = $2", now, ipAddress)
				if err != nil {
					log.Printf("Error updating request count: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				log.Printf("Updated request count for IP: %s", ipAddress)
			} else {
				_, err = db.Exec("UPDATE ip_requests SET request_count = 1, last_request = $1 WHERE ip_address = $2", now, ipAddress)
				if err != nil {
					log.Printf("Error resetting request count: %v", err)
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				log.Printf("Reset request count for IP: %s", ipAddress)
			}
		}

		next.ServeHTTP(w, r)
	})
}
