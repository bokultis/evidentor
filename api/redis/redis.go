package redis

import (
	"context"
	"os"

	"github.com/go-redis/redis/v8"
)

//Ctx global context
var Ctx = context.Background()

//RedisClient global
var RedisClient *redis.Client

// SetupRedis return redis cleint
func SetupRedis() *redis.Client {
	//Initializing redis
	dsn := os.Getenv("REDIS_DSN")
	if len(dsn) == 0 {
		dsn = "localhost:6379"
	}
	redisClient := redis.NewClient(&redis.Options{
		Addr:     dsn, //redis port
		Password: "",  // no password set
		DB:       0,   // use default DB
	})
	_, err := redisClient.Ping(Ctx).Result()
	if err != nil {
		panic(err)
	}

	return redisClient
}
