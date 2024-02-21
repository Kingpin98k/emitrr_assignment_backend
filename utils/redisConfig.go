package utils

import (
	"os"

	"github.com/go-redis/redis/v8"
)

// NewClient creates a new redis client
func Client() *redis.Client {
	//Connect to Redis
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDRESS"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	// ping, err := client.Ping().Result()
	// fmt.Println(ping, err)

	return client
}