package middleware

import (
	"ip-rate-control/pkg/ip"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/mux"
)

const (
	maxRequestsPerHour = 10
	ttlDuration        = time.Hour
	//ttlDuration = 2 * time.Minute
)

var (
	ipRequests      = make(map[string]int)
	queue           []queueItem
	mutex           sync.Mutex
	cleanupInterval time.Duration
)

type queueItem struct {
	ip           string
	timeToRemove time.Time
}

func cleanupQueue() {
	for {
		mutex.Lock()
		now := time.Now()

		for len(queue) > 0 && queue[0].timeToRemove.Before(now) {
			ip := queue[0].ip
			delete(ipRequests, ip)
			log.Printf("CleanupQueue: Removed IP: %s\n", ip)
			queue = queue[1:]
		}

		if len(queue) > 0 {
			cleanupInterval = queue[0].timeToRemove.Sub(now)
			if cleanupInterval < 0 {
				cleanupInterval = 0
			}
			log.Printf("CleanupQueue: Waiting for %v\n", cleanupInterval)
		} else {
			cleanupInterval = time.Minute
			log.Println("CleanupQueue: Queue is empty. Waiting for 1 minute.")
		}
		mutex.Unlock()

		time.Sleep(cleanupInterval)
	}
}

func init() {
	go cleanupQueue()
}

func RateLimitMiddleware() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ipAddress := ip.GetClientIP(r)
			//ipAddress := getLocalClientIP(r)

			mutex.Lock()
			defer mutex.Unlock()
			count, exists := ipRequests[ipAddress]
			if !exists {
				count = 0
				ipRequests[ipAddress] = count

				timeToRemove := time.Now().Add(ttlDuration)
				queue = append(queue, queueItem{
					ip:           ipAddress,
					timeToRemove: timeToRemove,
				})
				log.Printf("RateLimitMiddleware: IP added to queue for rate limiting: %s\n", ipAddress)
			}

			if count >= maxRequestsPerHour {
				log.Printf("RateLimitMiddleware: Rate limit exceeded for IP: %s\n", ipAddress)
				http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
				return
			}

			ipRequests[ipAddress]++
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
