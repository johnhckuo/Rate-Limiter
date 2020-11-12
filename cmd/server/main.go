package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/johnhckuo/Rate-Limiter/internal/environment"
	_ "github.com/joho/godotenv/autoload"

	"github.com/johnhckuo/Rate-Limiter/pkg/limiter"
)

func okHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("OK"))
}

func resetHandler(w http.ResponseWriter, r *http.Request) {

	w.Write([]byte("RESET"))
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc("/", okHandler)

	mux.HandleFunc("/reset", resetHandler)

	// Wrap the servemux with the limit middleware.
	log.Println("Listening on :4000...")
	var limiterClient limiter.Limiter

	val := strings.ToUpper(os.Getenv(environment.RateLimiter))
	if val == environment.GoLimiter {
		limiterClient = limiter.NewGoLimiter()
	} else if val == environment.RedisLimiter {
		limiterClient = limiter.NewRedisLimiter()
	} else {
		log.Fatalf("Unrecognized rate limiter option")
	}

	srv := &http.Server{
		Addr:    ":4000",
		Handler: limiterClient.Limit(mux),
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("%v", err)
		}
	}()
	c := make(chan os.Signal, 1)

	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

	// Block until we receive our signal.
	<-c
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}

	log.Println("Server exiting")

}

func init() {
	if _, exist := os.LookupEnv(environment.PersistStorage); !exist {
		log.Fatalf("env variables: %s is missing", environment.PersistStorage)
	}

	if _, exist := os.LookupEnv(environment.RateLimiter); !exist {
		log.Fatalf("env variables: %s is missing", environment.RateLimiter)
	}

	if _, exist := os.LookupEnv(environment.DbConnectionString); !exist {
		log.Fatalf("env variables: %s is missing", environment.DbConnectionString)
	}

	if _, exist := os.LookupEnv(environment.RateLimit); !exist {
		//setting default value as test requirement
		os.Setenv(environment.RateLimit, "60")
	}

	if _, exist := os.LookupEnv(environment.RateLimitExpirationSecond); !exist {
		//setting default value as test requirement
		os.Setenv(environment.RateLimitExpirationSecond, "60")
	}

	if _, exist := os.LookupEnv(environment.BurstLimit); !exist {
		os.Setenv(environment.BurstLimit, "10")
	}

}
