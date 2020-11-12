package persist

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis"
)

//Redis struct holds connected client for furture operation
type Redis struct {
	client *redis.Client
	ctx    context.Context
}

//NewRedis create a new Redis instance and return
func NewRedis(connectionString string) *Redis {
	log.Printf("Connecting to Redis: %s", connectionString)
	opt, err := redis.ParseURL(connectionString)
	if err != nil {
		log.Fatalf("%v", err.Error())
	}

	ctx := context.Background()

	client := redis.NewClient(opt)

	_, err = client.Ping().Result()
	if err != nil {
		log.Fatalf("%v", err.Error())
	}

	return &Redis{client: client, ctx: ctx}
}

//SetNX will set a value to the given key if the key is not existed, or else do nothing
func (r *Redis) SetNX(key string, val interface{}, duration int64) error {
	if ok, err := r.client.SetNX(key, val, time.Duration(duration)*time.Second).Result(); !ok {
		return err
	}
	return nil
}

/*
func (r *Redis) HSetNX(key string, duration int64) error {

		pip := r.client.Pipeline()

		if _, err := pip.HGet(key, "counter").Result(); err == redis.Nil {
			if _, err := pip.HMSet(key, ).Result(); err != nil {
				return err
			}

			if _, err := pip.Expire(key, time.Duration(duration)*time.Second).Result(); err != nil {
				return err
			}
		}

		_, err := pip.Exec()
		if err != nil {

			return err
		}


	return nil
}
*/

//Incr increase the give key by one
func (r *Redis) Incr(key string) (int64, error) {

	incr, err := r.client.Incr(key).Result()
	if err != nil {
		return -1, err
	}
	return incr, nil
}

//Reset will reset the value of the given key to 0
func (r *Redis) Reset(key string, expiration int) error {

	_, err := r.client.Set(key, 0, time.Duration(expiration)*time.Second).Result()
	if err != nil {
		return err
	}
	return nil
}
