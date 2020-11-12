package limiter

import (
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/johnhckuo/Rate-Limiter/internal/environment"
	"github.com/johnhckuo/Rate-Limiter/internal/utils"
	"golang.org/x/time/rate"
)

//GoLimiter Creates a map to hold the rate limiters for each visitor and a mutex.
type GoLimiter struct {
	visitors map[string]*visitor
	mu       sync.Mutex
}

// Create a custom visitor struct which holds the rate limiter for each
// visitor and the last time that the visitor was seen.
type visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

//NewGoLimiter creates GoLimiter and return
func NewGoLimiter() *GoLimiter {
	l := &GoLimiter{visitors: make(map[string]*visitor)}
	// Run a background goroutine to remove old entries from the visitors map.
	go l.cleanupVisitors()
	return l
}

//Limit will check if current request is good to go or not
func (l *GoLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Get the IP address for the current user.
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		keyName := utils.NewKeyName(r.Method, r.RequestURI, ip)

		// Call the getRateLimiter function to retreive the rate limiter for the current user.
		rateLimiter := l.getRateLimiter(keyName)

		if rateLimiter.Allow() == false {
			http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// Retrieve and return the rate limiter for the current visitor if it
// already exists. Otherwise create a new rate limiter and add it to
// the visitors map, using the method + url + IP address as the key.
func (l *GoLimiter) getRateLimiter(keyName string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	_, exists := l.visitors[keyName]
	if !exists {
		burst, _ := strconv.Atoi(os.Getenv(environment.BurstLimit))
		limit, _ := strconv.Atoi(os.Getenv(environment.RateLimit))
		limiter := rate.NewLimiter(rate.Limit(limit), burst)
		l.visitors[keyName] = &visitor{limiter, time.Now()}
	} else {
		l.visitors[keyName].lastSeen = time.Now()
	}

	return l.visitors[keyName].limiter
}

// Every minute check the map for visitors that haven't been seen for
// more than 1 minutes and delete the entries.
func (l *GoLimiter) cleanupVisitors() {
	expiration, _ := strconv.ParseInt(os.Getenv(environment.RateLimitExpirationSecond), 10, 64)
	for {
		time.Sleep(time.Minute)

		l.mu.Lock()
		for ip, v := range l.visitors {
			if time.Since(v.lastSeen) > time.Duration(expiration)*time.Second {
				delete(l.visitors, ip)
			}
		}
		l.mu.Unlock()
	}
}
