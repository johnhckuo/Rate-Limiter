package limiter

import (
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/johnhckuo/Rate-Limiter/internal/environment"
	"github.com/johnhckuo/Rate-Limiter/internal/utils"
	"github.com/johnhckuo/Rate-Limiter/pkg/persist"
)

//RedisLimiter Creates a map to hold the rate limiters for each visitor and a mutex.
type RedisLimiter struct {
	mu       sync.Mutex
	visitors map[string]*redisVisitor
	client   persist.Db
}

// Create a custom redisVisitor struct which holds the rate limiter for each redisVisitor
type redisVisitor struct {
	limit    int64
	burst    int
	duration time.Duration
}

//NewRedisLimiter creates a new Redis Limiter and return
func NewRedisLimiter() *RedisLimiter {
	visitors := make(map[string]*redisVisitor)
	r := &RedisLimiter{visitors: visitors}

	val := strings.ToUpper(os.Getenv(environment.PersistStorage))
	if val == environment.RedisStorage {
		r.client = persist.NewRedis(os.Getenv(environment.DbConnectionString))
	}

	return r
}

//Limit will check if current request is good to go or not
func (r *RedisLimiter) Limit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {

		url := req.RequestURI
		// Get the IP address for the current user.
		ip, _, err := net.SplitHostPort(req.RemoteAddr)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// only path "/" will be rate limited
		if url == "/" {

			// Call the getRateLimiter function to retreive the rate limiter for the current user.
			keyName := utils.NewKeyName(req.Method, req.RequestURI, ip)
			r.getRateLimiter(keyName)

			if r.allow(keyName) == false {
				http.Error(w, http.StatusText(429), http.StatusTooManyRequests)
				return
			}
		} else if url == "/reset" {

			keyName := utils.NewKeyName("get", "/", ip)
			err := r.reset(keyName)
			if err != nil {
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
				return
			}
		}

		next.ServeHTTP(w, req)
	})
}

// Retrieve and return the rate limiter for the current redisVisitor if it
// already exists. Otherwise create a new rate limit counter in redis and using the method + url + IP address as the key.
func (r *RedisLimiter) getRateLimiter(hashKey string) {
	r.mu.Lock()
	//defer req.mu.Unlock()
	expiration, _ := strconv.ParseInt(os.Getenv(environment.RateLimitExpirationSecond), 10, 64)

	_, exists := r.visitors[hashKey]
	if !exists {
		burst, _ := strconv.Atoi(os.Getenv(environment.BurstLimit))
		limit, _ := strconv.ParseInt(os.Getenv(environment.RateLimit), 10, 64)
		r.visitors[hashKey] = &redisVisitor{limit, burst, time.Duration(expiration) * time.Second}

	}
	r.mu.Unlock()

	//set burst rate limit counter
	err := r.client.SetNX(hashKey+"_burst", 0, 1)
	if err != nil {
		log.Printf("Error %v", err)
	}

	// set per minute rate limit counter
	err = r.client.SetNX(hashKey, 0, expiration)
	if err != nil {
		log.Printf("Error %v", err)
	}

	return
}

func (r *RedisLimiter) allow(hashKey string) bool {
	log.Printf("Checking API usage of key: %v", hashKey)
	counter, err := r.client.Incr(hashKey)
	log.Printf("per_minute_counter: %v \n", counter)
	if counter > r.visitors[hashKey].limit || err != nil {
		log.Printf("key %v per minute limit exceeded", hashKey)
		return false
	}

	burstCounter, err := r.client.Incr(hashKey + "_burst")
	log.Printf("burst_counter: %v \n", burstCounter)

	if burstCounter > int64(r.visitors[hashKey].burst) || err != nil {
		log.Printf("key %v burst limit exceeded", hashKey)
		return false
	}

	return true
}

func (r *RedisLimiter) reset(hashKey string) error {
	log.Printf("Resetting key %v", hashKey)
	expiration, _ := strconv.Atoi(os.Getenv(environment.RateLimitExpirationSecond))

	err := r.client.Reset(hashKey, expiration)
	if err != nil {
		return err
	}

	err = r.client.Reset(hashKey+"_burst", 1)
	if err != nil {
		return err
	}
	return nil
}
