package middleware

import (
	"context"
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/redis/go-redis/v9"
	"ip-rate-control/pkg/ip"
	"log"
	"net/http"
	"time"
)

const (
	maxRequestsPerHour = 10
	expirationTime     = time.Hour
	//expirationTime = 2 * time.Minute
)

func RateLimitMiddleware(redisClient *redis.Client, db *sql.DB) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.Background()
			ipAddress := ip.GetClientIP(r)
			//ipAddress := getLocalClientIP(r)
			log.Println(ipAddress)

			now := time.Now().In(time.Local)
			log.Println("Current Time (Local):", now)

			requestKey := "request_count:" + ipAddress

			// Redis üzerinde gelen ip kontrolü
			requestCount, err := redisClient.Get(ctx, requestKey).Int64()
			if err == redis.Nil {
				// Redise kayıt
				err = redisClient.Set(ctx, requestKey, 1, expirationTime).Err()
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}
				requestCount = 1

				// ip postgre kontrolü
				var exists bool
				err = db.QueryRow("SELECT EXISTS(SELECT 1 FROM ip_requests WHERE ip_address = $1)", ipAddress).Scan(&exists)
				if err != nil {
					http.Error(w, "Internal Server Error", http.StatusInternalServerError)
					return
				}

				if !exists {
					_, err = db.Exec("INSERT INTO ip_requests (ip_address, first_request) VALUES ($1, $2)", ipAddress, now)
					if err != nil {
						http.Error(w, "Internal Server Error", http.StatusInternalServerError)
						return
					}
				}
			} else if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			// requestCount güncellemesi
			requestCount, err = redisClient.Incr(ctx, requestKey).Result()
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}

			if requestCount > maxRequestsPerHour {
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

/*
func getLocalClientIP(r *http.Request) string {
	queryIP := r.URL.Query().Get("ip")
	return queryIP
}
*/
