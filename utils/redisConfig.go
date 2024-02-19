package utils

import (
	"github.com/go-redis/redis/v8"
)

// NewClient creates a new redis client
func Client() *redis.Client {
	//Connect to Redis
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	// ping, err := client.Ping().Result()
	// fmt.Println(ping, err)

	return client
}