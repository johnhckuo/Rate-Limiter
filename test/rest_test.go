package test

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/johnhckuo/Rate-Limiter/internal/environment"

	"github.com/go-resty/resty/v2"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	"github.com/stretchr/testify/assert"
)

func Test_Reset(t *testing.T) {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	client := resty.New()
	resp, _ := client.R().Get("http://localhost:4000/reset")
	assert.Equal(t, 200, resp.StatusCode())
}
func Test_Basic(t *testing.T) {

	client := resty.New()

	resp, _ := client.R().Get("http://localhost:4000")

	assert.Equal(t, 200, resp.StatusCode())
	resp, _ = client.R().Get("http://localhost:4000/reset")
	assert.Equal(t, 200, resp.StatusCode())

}

func Test_Multiple(t *testing.T) {

	client := resty.New()

	burstLimit, _ := strconv.Atoi(os.Getenv(environment.BurstLimit))
	for i := 0; i < burstLimit; i++ {
		resp, _ := client.R().Get("http://localhost:4000")

		assert.Equal(t, 200, resp.StatusCode())
	}
	resp, _ := client.R().Get("http://localhost:4000/reset")
	assert.Equal(t, 200, resp.StatusCode())

}

func Test_Exceed_Burst(t *testing.T) {
	client := resty.New()
	var wg sync.WaitGroup
	successCounter := int64(0)
	failCounter := int64(0)

	burstLimit, _ := strconv.Atoi(os.Getenv(environment.BurstLimit))
	wg.Add(burstLimit + 1)
	for i := 0; i <= burstLimit; i++ {
		go func(wg *sync.WaitGroup, i int) {
			resp, _ := client.R().Get("http://localhost:4000")
			log.Printf("Sending Request #%v \n", i)
			if resp.StatusCode() == 200 {
				atomic.AddInt64(&successCounter, 1)
			} else {
				atomic.AddInt64(&failCounter, 1)
			}
			wg.Done()
		}(&wg, i)
	}
	wg.Wait()
	assert.Equal(t, int64(burstLimit), successCounter)
	assert.Equal(t, int64(1), failCounter)

	fmt.Printf("Done\n")
	resp, _ := client.R().Get("http://localhost:4000/reset")
	assert.Equal(t, 200, resp.StatusCode())
}

func Test_One_Minute_Exceed(t *testing.T) {
	client := resty.New()

	rateLimit, _ := strconv.Atoi(os.Getenv(environment.RateLimit))

	for i := 0; i <= rateLimit; i++ {
		if i == rateLimit/2 {
			log.Println("30 sec has passed")
		}
		time.Sleep(time.Millisecond * 500)
		resp, _ := client.R().Get("http://localhost:4000")
		if i == rateLimit {
			assert.Equal(t, 429, resp.StatusCode())
		} else {
			assert.Equal(t, 200, resp.StatusCode())
		}
	}
	resp, _ := client.R().Get("http://localhost:4000/reset")
	assert.Equal(t, 200, resp.StatusCode())
}

func Test_Two_Minute(t *testing.T) {
	client := resty.New()

	rateLimit, _ := strconv.Atoi(os.Getenv(environment.RateLimit))

	for i := 0; i < rateLimit*2; i++ {
		if i == rateLimit {
			log.Println("1 minute has passed")
		}
		time.Sleep(time.Second * 1)
		resp, _ := client.R().Get("http://localhost:4000")
		assert.Equal(t, 200, resp.StatusCode())
	}
	resp, _ := client.R().Get("http://localhost:4000/reset")
	assert.Equal(t, 200, resp.StatusCode())
}
