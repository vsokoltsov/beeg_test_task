package utils

import (
	"fmt"

	"github.com/go-redis/redis"
)

// RedisClient saves connection to redis in instance
var RedisClient *redis.Client

// InitRedis initializes redis instance
func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, redisErr := RedisClient.Ping().Result()
	if redisErr != nil {
		fmt.Println(redisErr)
	}
}
